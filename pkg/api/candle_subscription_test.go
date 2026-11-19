package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// TestCandleSymbolSubscription subscribes to candles by a candle symbol, alone and inside a time series symbol, through
// a local hub and checks that a published candle of that symbol arrives.
func TestCandleSymbolSubscription(t *testing.T) {
	cases := []struct {
		name   string
		symbol func(s *candle.CandleSymbol) osub.Symbol
	}{
		{"candle symbol", func(s *candle.CandleSymbol) osub.Symbol { return osub.Candle(s) }},
		{"time series of osub.Candle", func(s *candle.CandleSymbol) osub.Symbol {
			return osub.NewTimeSeriesSubscriptionSymbol(osub.Candle(s), 0)
		}},
		{"time series of candle.CandleSymbol", func(s *candle.CandleSymbol) osub.Symbol {
			return osub.NewTimeSeriesSubscriptionSymbol(s, 0)
		}},
	}
	feed, publisher := newLocalHub(t)
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// The default price attribute is dropped: the subscription and the candle use the normalized symbol.
			symbol := candle.NewCandleSymbol(fmt.Sprintf("CANDLE/SUB/%d{=5m,price=last}", i))
			want := fmt.Sprintf("CANDLE/SUB/%d{=5m}", i)
			if symbol.String() != want {
				t.Fatalf("candle symbol %q, want %q", symbol, want)
			}
			sub, err := feed.CreateSubscription(eventcodes.Candle)
			if err != nil {
				t.Fatalf("create subscription: %v", err)
			}
			defer sub.Close()
			collector := &eventCollector{events: make(chan interface{}, 16)}
			if err := sub.AddListener(collector); err != nil {
				t.Fatalf("add listener: %v", err)
			}
			if err := sub.AddSymbol(c.symbol(symbol)); err != nil {
				t.Fatalf("add symbol: %v", err)
			}
			event := candle.NewCandle(want)
			event.SetTime(roundTripTime)
			event.SetCount(1)
			event.SetClose(42.5)
			if err := publisher.Publish([]interface{}{event}); err != nil {
				t.Fatalf("publish: %v", err)
			}
			select {
			case got := <-collector.events:
				received, ok := got.(*candle.Candle)
				if !ok {
					t.Fatalf("received %T, want *candle.Candle", got)
				}
				if received.EventSymbol().String() != want || received.Close() != 42.5 {
					t.Fatalf("received %v, want the candle of %s with close 42.5", received, want)
				}
			case <-time.After(10 * time.Second):
				t.Fatalf("no candle of %s received within 10 s", want)
			}
		})
	}
}
