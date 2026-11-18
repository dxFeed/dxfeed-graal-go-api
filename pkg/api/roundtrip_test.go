package api

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/summary"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
)

// eventCollector passes the received events to a channel.
type eventCollector struct {
	events chan interface{}
}

func (c *eventCollector) Update(events []interface{}) {
	for _, event := range events {
		c.events <- event
	}
}

// TestPublishSubscribeRoundTrip publishes one event of every type that has a native mapper through a local hub and
// checks that the subscription receives it unchanged: the event goes through both mappers and the native layer.
func TestPublishSubscribeRoundTrip(t *testing.T) {
	description := "round trip"
	cases := []struct {
		code  eventcodes.EventCode
		event func(symbol string) events.EventType
	}{
		{eventcodes.Quote, func(symbol string) events.EventType {
			e := quote.NewQuote(symbol)
			e.SetBidPrice(10.5)
			e.SetAskPrice(11.25)
			e.SetBidSize(100)
			return e
		}},
		{eventcodes.Profile, func(symbol string) events.EventType {
			e := profile.NewProfile(symbol)
			e.SetDescription(&description)
			e.SetHighLimitPrice(99.5)
			return e
		}},
		{eventcodes.Summary, func(symbol string) events.EventType {
			e := summary.NewSummary(symbol)
			e.SetDayOpenPrice(12.5)
			e.SetOpenInterest(42)
			return e
		}},
		{eventcodes.Greeks, func(symbol string) events.EventType {
			e := greeks.NewGreeks(symbol)
			e.SetPrice(1.5)
			e.SetDelta(0.25)
			return e
		}},
		{eventcodes.Candle, func(symbol string) events.EventType {
			e := candle.NewCandle(symbol)
			e.SetCount(1)
			e.SetOpen(13.5)
			e.SetHigh(13.5)
			e.SetLow(13.5)
			e.SetClose(13.5)
			e.SetVolume(7)
			return e
		}},
		{eventcodes.Underlying, func(symbol string) events.EventType {
			e := underlying.NewUnderlying(symbol)
			e.SetVolatility(0.5)
			return e
		}},
		{eventcodes.Trade, func(symbol string) events.EventType {
			e := trade.NewTrade(symbol)
			e.SetPrice(14.5)
			e.SetSize(3)
			return e
		}},
		{eventcodes.TradeETH, func(symbol string) events.EventType {
			e := trade.NewTradeETH(symbol)
			e.SetPrice(15.5)
			e.SetSize(4)
			return e
		}},
		{eventcodes.TimeAndSale, func(symbol string) events.EventType {
			e := timeandsale.NewTimeAndSale(symbol)
			e.SetPrice(16.5)
			e.SetSize(5)
			return e
		}},
		{eventcodes.Order, func(symbol string) events.EventType {
			e := order.NewOrder(symbol)
			e.SetSide(side.Buy)
			e.SetPrice(17.5)
			e.SetSize(6)
			return e
		}},
		{eventcodes.SpreadOrder, func(symbol string) events.EventType {
			e := order.NewSpreadOrder(symbol)
			e.SetSide(side.Buy)
			e.SetPrice(18.5)
			e.SetSize(7)
			return e
		}},
		{eventcodes.AnalyticOrder, func(symbol string) events.EventType {
			e := order.NewAnalyticOrder(symbol)
			e.SetSide(side.Buy)
			e.SetPrice(19.5)
			e.SetSize(8)
			// The default data scheme does not carry the iceberg fields; they come back as NaN.
			e.SetIcebergPeakSize(math.NaN())
			e.SetIcebergHiddenSize(math.NaN())
			e.SetIcebergExecutedSize(math.NaN())
			return e
		}},
	}

	feed, publisher := newLocalHub(t)
	for _, c := range cases {
		// Each type gets its own symbol, so no event of an earlier case can reach a later one.
		symbol := fmt.Sprintf("ROUND/TRIP/%d", c.code)
		event := c.event(symbol)
		t.Run(fmt.Sprintf("%T", event), func(t *testing.T) {
			sub, err := feed.CreateSubscription(c.code)
			if err != nil {
				t.Fatalf("create subscription: %v", err)
			}
			defer sub.Close()
			collector := &eventCollector{events: make(chan interface{}, 16)}
			if err := sub.AddListener(collector); err != nil {
				t.Fatalf("add listener: %v", err)
			}
			if err := sub.AddSymbol(osub.String(symbol)); err != nil {
				t.Fatalf("add symbol: %v", err)
			}
			if err := publisher.Publish([]interface{}{event}); err != nil {
				t.Fatalf("publish: %v", err)
			}
			want := fmt.Sprint(event)
			select {
			case got := <-collector.events:
				if fmt.Sprint(got) != want {
					t.Fatalf("received event differs:\n got %v\nwant %v", got, want)
				}
			case <-time.After(10 * time.Second):
				t.Fatalf("no event received within 10 s; published %v", want)
			}
		})
	}
}

// TestRemoveSymbolsStopsDelivery adds symbols as a list, receives their events, removes them one by one and as a list,
// and checks that later events of the removed symbols are not delivered.
func TestRemoveSymbolsStopsDelivery(t *testing.T) {
	feed, publisher := newLocalHub(t)
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer sub.Close()
	collector := &eventCollector{events: make(chan interface{}, 16)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	symbols := []string{"REMOVE/A", "REMOVE/B", "REMOVE/C"}
	if err := sub.AddSymbols(osub.Strings(symbols...)...); err != nil {
		t.Fatalf("add symbols: %v", err)
	}
	publish := func() {
		t.Helper()
		quotes := make([]interface{}, 0, len(symbols))
		for _, s := range symbols {
			quotes = append(quotes, quote.NewQuote(s))
		}
		if err := publisher.Publish(quotes); err != nil {
			t.Fatalf("publish: %v", err)
		}
	}
	publish()
	for range symbols {
		select {
		case <-collector.events:
		case <-time.After(10 * time.Second):
			t.Fatal("not all subscribed symbols delivered an event within 10 s")
		}
	}
	if err := sub.RemoveSymbol(osub.String(symbols[0])); err != nil {
		t.Fatalf("remove symbol: %v", err)
	}
	if err := sub.RemoveSymbols(osub.Strings(symbols[1:]...)...); err != nil {
		t.Fatalf("remove symbols: %v", err)
	}
	publish()
	select {
	case got := <-collector.events:
		t.Fatalf("event of a removed symbol delivered: %v", got)
	case <-time.After(time.Second):
	}
}

// newLocalHub returns the feed and the publisher of a local hub endpoint that is closed when the test ends.
func newLocalHub(t *testing.T) (*DXFeed, *DXPublisher) {
	t.Helper()
	endpoint, err := NewEndpoint(LocalHub)
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	t.Cleanup(func() { _ = endpoint.Close() })
	feed, err := endpoint.GetFeed()
	if err != nil {
		t.Fatalf("get feed: %v", err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		t.Fatalf("get publisher: %v", err)
	}
	return feed, publisher
}
