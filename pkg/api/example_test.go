package api_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// quotePrinter is an event listener that prints quotes and signals the first one.
type quotePrinter struct {
	received chan struct{}
}

func (p quotePrinter) Update(events []interface{}) {
	for _, event := range events {
		if q, ok := event.(*quote.Quote); ok {
			fmt.Printf("%s bid %.2f ask %.2f\n", *q.EventSymbol(), q.BidPrice(), q.AskPrice())
			select {
			case p.received <- struct{}{}:
			default:
			}
		}
	}
}

// A local hub connects a feed and a publisher in one process: the quote published to it reaches the subscription.
// A feed of a real data source is created with api.Feed and connected with Connect, for example
// endpoint.Connect("demo.dxfeed.com:7300").
func Example() {
	endpoint, err := api.NewEndpoint(api.LocalHub)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = endpoint.CloseAndAwaitTermination(context.Background()) }()

	feed, err := endpoint.GetFeed()
	if err != nil {
		log.Fatal(err)
	}
	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = subscription.Close() }()

	printer := quotePrinter{received: make(chan struct{}, 1)}
	if err := subscription.AddListener(printer); err != nil {
		log.Fatal(err)
	}
	if err := subscription.AddSymbol("EXAMPLE"); err != nil {
		log.Fatal(err)
	}

	publisher, err := endpoint.GetPublisher()
	if err != nil {
		log.Fatal(err)
	}
	q := quote.NewQuote("EXAMPLE")
	q.SetBidPrice(100.25)
	q.SetAskPrice(100.5)
	if err := publisher.Publish([]interface{}{q}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-printer.received:
	case <-time.After(10 * time.Second):
		log.Fatal("no quote received")
	}
	// Output: EXAMPLE bid 100.25 ask 100.50
}

// Time series and indexed events take symbols of package osub; a plain string or a wildcard works for the others.
func ExampleDXFeedSubscription_AddSymbols() {
	endpoint, err := api.NewEndpoint(api.LocalHub)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()
	feed, err := endpoint.GetFeed()
	if err != nil {
		log.Fatal(err)
	}
	subscription, err := feed.CreateSubscription(eventcodes.Quote, eventcodes.Candle)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = subscription.Close() }()

	weekAgo := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()
	err = subscription.AddSymbols(
		"AAPL",
		osub.NewWildcardSymbol(),
		osub.NewTimeSeriesSubscriptionSymbol("AAPL{=d}", weekAgo),
	)
	fmt.Println(err)
	// Output: <nil>
}

// A wait on the endpoint ends with the context; the SDK wait goes on in the background.
func ExampleDXEndpoint_AwaitProcessed() {
	endpoint, err := api.NewEndpoint(api.LocalHub)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println(endpoint.AwaitProcessed(ctx))
	// Output: <nil>
}
