package api

import (
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
)

// subscribeAndPublish subscribes the feed to the symbols of the events, publishes the events and waits until the
// subscription has received count events, so that the feed has them in its cache. The subscription is closed when the
// test ends.
func subscribeAndPublish(t *testing.T, feed *DXFeed, publisher *DXPublisher, eventType eventcodes.EventCode,
	symbols []osub.Symbol, events []interface{}, count int) *DXFeedSubscription {
	t.Helper()
	sub, err := feed.CreateSubscription(eventType)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	t.Cleanup(func() { _ = sub.Close() })
	collector := &eventCollector{events: make(chan interface{}, 64)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	if err := sub.AddSymbols(symbols...); err != nil {
		t.Fatalf("add symbols: %v", err)
	}
	if err := publisher.Publish(events); err != nil {
		t.Fatalf("publish: %v", err)
	}
	for i := 0; i < count; i++ {
		select {
		case <-collector.events:
		case <-time.After(10 * time.Second):
			t.Fatalf("%d of %d events delivered within 10 s", i, count)
		}
	}
	return sub
}

func TestGetLastEvent(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "LAST/A"

	got, err := feed.GetLastEventIfSubscribed(eventcodes.Quote, osub.String(symbol))
	if err != nil || got != nil {
		t.Fatalf("GetLastEventIfSubscribed before a subscription = %v, %v; want nil, nil", got, err)
	}
	// Without the event the result is a copy of the given event.
	probe := quote.NewQuote(symbol)
	probe.SetBidPrice(1.25)
	got, err = feed.GetLastEvent(probe)
	if err != nil {
		t.Fatalf("GetLastEvent without data: %v", err)
	}
	if q, ok := got.(*quote.Quote); !ok || q == probe || q.BidPrice() != 1.25 {
		t.Fatalf("GetLastEvent without data = %v, want a copy of the given quote", got)
	}

	published := quote.NewQuote(symbol)
	published.SetBidPrice(10.5)
	subscribeAndPublish(t, feed, publisher, eventcodes.Quote, osub.Strings(symbol), []interface{}{published}, 1)

	got, err = feed.GetLastEvent(quote.NewQuote(symbol))
	if err != nil {
		t.Fatalf("GetLastEvent: %v", err)
	}
	if q, ok := got.(*quote.Quote); !ok || q.BidPrice() != 10.5 || *q.EventSymbol() != symbol {
		t.Fatalf("GetLastEvent = %v, want the published quote", got)
	}
	got, err = feed.GetLastEventIfSubscribed(eventcodes.Quote, osub.String(symbol))
	if err != nil {
		t.Fatalf("GetLastEventIfSubscribed: %v", err)
	}
	if q, ok := got.(*quote.Quote); !ok || q.BidPrice() != 10.5 {
		t.Fatalf("GetLastEventIfSubscribed = %v, want the published quote", got)
	}
}

func TestGetLastEvents(t *testing.T) {
	feed, publisher := newLocalHub(t)
	symbols := []string{"LAST/B", "LAST/C"}
	published := make([]interface{}, len(symbols))
	for i, s := range symbols {
		q := quote.NewQuote(s)
		q.SetBidPrice(float64(i + 1))
		published[i] = q
	}
	subscribeAndPublish(t, feed, publisher, eventcodes.Quote, osub.Strings(symbols...), published, len(symbols))

	// The order of the result follows the order of the request.
	got, err := feed.GetLastEvents([]interface{}{quote.NewQuote(symbols[1]), quote.NewQuote(symbols[0])})
	if err != nil {
		t.Fatalf("GetLastEvents: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("GetLastEvents returned %d events, want 2", len(got))
	}
	for i, want := range []float64{2, 1} {
		if q, ok := got[i].(*quote.Quote); !ok || q.BidPrice() != want {
			t.Fatalf("GetLastEvents[%d] = %v, want bid price %v", i, got[i], want)
		}
	}
}

func TestGetLastEventRejectsNonEvents(t *testing.T) {
	feed, _ := newLocalHub(t)
	if _, err := feed.GetLastEvent("not an event"); err == nil {
		t.Fatal("GetLastEvent accepted a string")
	}
	if _, err := feed.GetLastEvents([]interface{}{quote.NewQuote("X"), 42}); err == nil {
		t.Fatal("GetLastEvents accepted an int")
	}
}

func TestGetIndexedEventsIfSubscribed(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "INDEXED/A"

	got, err := feed.GetIndexedEventsIfSubscribed(eventcodes.Order, osub.String(symbol), order.Default())
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("GetIndexedEventsIfSubscribed before a subscription = %v, %v; want an empty list", got, err)
	}

	published := make([]interface{}, 2)
	for i := range published {
		o := order.NewOrder(symbol)
		if err := o.SetIndex(int64(i + 1)); err != nil {
			t.Fatalf("set index: %v", err)
		}
		o.SetSide(side.Buy)
		o.SetPrice(float64(10 + i))
		o.SetSize(1)
		published[i] = o
	}
	subscribeAndPublish(t, feed, publisher, eventcodes.Order, osub.Strings(symbol), published, len(published))

	got, err = feed.GetIndexedEventsIfSubscribed(eventcodes.Order, osub.String(symbol), order.Default())
	if err != nil {
		t.Fatalf("GetIndexedEventsIfSubscribed: %v", err)
	}
	if len(got) != len(published) {
		t.Fatalf("GetIndexedEventsIfSubscribed returned %d events, want %d: %v", len(got), len(published), got)
	}
	for i, e := range got {
		if o, ok := e.(*order.Order); !ok || o.Index() != int64(i+1) {
			t.Fatalf("GetIndexedEventsIfSubscribed[%d] = %v, want the order with index %d", i, e, i+1)
		}
	}
	if _, err := feed.GetIndexedEventsIfSubscribed(eventcodes.Order, osub.String(symbol), nil); err == nil {
		t.Fatal("GetIndexedEventsIfSubscribed accepted a nil source")
	}
}

func TestGetTimeSeriesIfSubscribed(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "SERIES/A"
	base := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC).UnixMilli()

	published := make([]interface{}, 3)
	for i := range published {
		e := timeandsale.NewTimeAndSale(symbol)
		e.SetTime(base + int64(i)*1000)
		e.SetPrice(float64(100 + i))
		e.SetSize(1)
		published[i] = e
	}
	subscribeAndPublish(t, feed, publisher, eventcodes.TimeAndSale,
		[]osub.Symbol{osub.NewTimeSeriesSubscriptionSymbol(symbol, base)}, published, len(published))

	all, err := feed.GetTimeSeriesIfSubscribed(eventcodes.TimeAndSale, osub.String(symbol), base, math.MaxInt64)
	if err != nil {
		t.Fatalf("GetTimeSeriesIfSubscribed: %v", err)
	}
	if len(all) != len(published) {
		t.Fatalf("GetTimeSeriesIfSubscribed returned %d events, want %d: %v", len(all), len(published), all)
	}
	// The from-time filters the events out.
	got, err := feed.GetTimeSeriesIfSubscribed(eventcodes.TimeAndSale, osub.String(symbol), base+1000, math.MaxInt64)
	if err != nil {
		t.Fatalf("GetTimeSeriesIfSubscribed from a later time: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("GetTimeSeriesIfSubscribed from a later time returned %d events, want 2: %v", len(got), got)
	}
	for _, e := range got {
		if ts, ok := e.(*timeandsale.TimeAndSale); !ok || ts.Time() < base+1000 {
			t.Fatalf("GetTimeSeriesIfSubscribed returned %v, before the from-time", e)
		}
	}
	if _, err := feed.GetTimeSeriesIfSubscribed(eventcodes.TimeAndSale,
		osub.NewTimeSeriesSubscriptionSymbol(symbol, base), base, math.MaxInt64); err == nil {
		t.Fatal("GetTimeSeriesIfSubscribed accepted a time series subscription symbol")
	}
}
