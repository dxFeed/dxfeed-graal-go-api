package api_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
)

// eventChannel is an event listener that passes the events to a channel.
type eventChannel chan interface{}

func (c eventChannel) Update(events []interface{}) {
	for _, event := range events {
		c <- event
	}
}

// localHub creates a local hub endpoint, its feed and its publisher.
func localHub() (*api.DXEndpoint, *api.DXFeed, *api.DXPublisher) {
	endpoint, err := api.NewEndpoint(api.LocalHub)
	if err != nil {
		log.Fatal(err)
	}
	feed, err := endpoint.GetFeed()
	if err != nil {
		log.Fatal(err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		log.Fatal(err)
	}
	return endpoint, feed, publisher
}

// receive returns the next event of the channel, or ends the program after 10 s.
func receive(events eventChannel) interface{} {
	select {
	case event := <-events:
		return event
	case <-time.After(10 * time.Second):
		log.Fatal("no event received")
		return nil
	}
}

// A time series subscription gets the events of a symbol from a given time on: here the daily candles of a local
// hub. The history comes as a snapshot, newest first, marked with the snapshot flags of package events.
func ExampleDXFeed_CreateSubscription_timeSeries() {
	endpoint, feed, publisher := localHub()
	defer func() { _ = endpoint.CloseAndAwaitTermination(context.Background()) }()

	subscription, err := feed.CreateSubscription(eventcodes.Candle)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = subscription.Close() }()
	received := make(eventChannel, 16)
	if err := subscription.AddListener(received); err != nil {
		log.Fatal(err)
	}
	// The snapshot of a local hub reaches the subscription when it covers the whole subscribed time range.
	day := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	from := day.AddDate(0, 0, -2).UnixMilli()
	if err := subscription.AddSymbol(osub.NewTimeSeriesSubscriptionSymbol("EXAMPLE{=1d}", from)); err != nil {
		log.Fatal(err)
	}

	// A data source publishes the history of the symbol as a snapshot.
	closes := []float64{102, 101, 100}
	snapshot := make([]interface{}, len(closes))
	for i, price := range closes {
		c := candle.NewCandle("EXAMPLE{=1d}")
		c.SetTime(day.AddDate(0, 0, -i).UnixMilli())
		c.SetClose(price)
		snapshot[i] = c
	}
	snapshot[0].(*candle.Candle).SetEventFlags(events.SnapshotBegin.Flag())
	snapshot[len(snapshot)-1].(*candle.Candle).SetEventFlags(events.SnapshotEnd.Flag())
	if err := publisher.Publish(snapshot); err != nil {
		log.Fatal(err)
	}

	for {
		c := receive(received).(*candle.Candle)
		fmt.Println(time.UnixMilli(c.Time()).UTC().Format(time.DateOnly), c.Close())
		if events.IsSnapshotEndOrSnip(c) {
			break
		}
	}
	// Output:
	// 2026-01-05 102
	// 2026-01-04 101
	// 2026-01-03 100
}

// Indexed events have several values per symbol, told apart by their index: here two orders of a book. A
// subscription with an osub.IndexedEventSubscriptionSymbol gets the events of one source only. The events of
// different indices may come in any order.
func ExampleDXFeed_CreateSubscription_indexed() {
	endpoint, feed, publisher := localHub()
	defer func() { _ = endpoint.CloseAndAwaitTermination(context.Background()) }()

	subscription, err := feed.CreateSubscription(eventcodes.Order)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = subscription.Close() }()
	received := make(eventChannel, 16)
	if err := subscription.AddListener(received); err != nil {
		log.Fatal(err)
	}
	if err := subscription.AddSymbol(osub.NewIndexedEventSubscriptionSymbol("EXAMPLE", order.Default())); err != nil {
		log.Fatal(err)
	}

	bid := order.NewOrder("EXAMPLE")
	if err := bid.SetIndex(1); err != nil {
		log.Fatal(err)
	}
	bid.SetSide(side.Buy)
	bid.SetPrice(99.5)
	bid.SetSize(300)
	ask := order.NewOrder("EXAMPLE")
	if err := ask.SetIndex(2); err != nil {
		log.Fatal(err)
	}
	ask.SetSide(side.Sell)
	ask.SetPrice(100.5)
	ask.SetSize(200)
	if err := publisher.Publish([]interface{}{bid, ask}); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		o := receive(received).(*order.Order)
		fmt.Println(o.Index(), o.Side(), o.Price(), o.Size())
	}
	// Unordered output:
	// 1 Buy 99.5 300
	// 2 Sell 100.5 200
}

// stateLogger is a connection state listener that prints the changes.
type stateLogger struct{}

func (stateLogger) UpdateState(old, new common.ConnectionState) {
	fmt.Println("connection:", old, "->", new)
}

// A connection state listener reports the state of the endpoint as it connects, loses and restores the connection.
// This example needs the network, so it is not run by the tests.
func ExampleDXEndpoint_AddListener() {
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()

	endpoint.AddListener(stateLogger{})
	if err := endpoint.Connect("demo.dxfeed.com:7300"); err != nil {
		log.Fatal(err)
	}
	// The listener prints "Not connected -> Connecting", then "Connecting -> Connected".
	time.Sleep(5 * time.Second)
}
