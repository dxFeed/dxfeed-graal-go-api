package native

import (
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

// TestGoSymbolRoundTrip converts symbols to C and back without the SDK.
func TestGoSymbolRoundTrip(t *testing.T) {
	cases := []osub.Symbol{
		osub.String("AAPL"),
		osub.NewWildcardSymbol(),
		osub.NewTimeSeriesSubscriptionSymbol(osub.String("IBM"), 5),
		osub.NewIndexedEventSubscriptionSymbol("MSFT", order.NtvL2()),
		osub.NewIndexedEventSubscriptionSymbol("AAPL", events.NewIndexedEventSource(77, "SRC77")),
	}
	for _, symbol := range cases {
		c, err := newCSymbol(symbol)
		if err != nil {
			t.Fatalf("newCSymbol(%#v): %v", symbol, err)
		}
		got, err := goSymbol(c)
		freeCSymbol(c)
		if err != nil {
			t.Fatalf("goSymbol(%#v): %v", symbol, err)
		}
		switch want := symbol.(type) {
		case osub.String:
			if got != want {
				t.Errorf("goSymbol = %#v, want %#v", got, want)
			}
		case *osub.WildcardSymbol:
			if _, ok := got.(*osub.WildcardSymbol); !ok {
				t.Errorf("goSymbol = %#v, want the wildcard", got)
			}
		case *osub.TimeSeriesSubscriptionSymbol:
			ts, ok := got.(*osub.TimeSeriesSubscriptionSymbol)
			if !ok || ts.Symbol() != "IBM" || ts.FromTime() != 5 {
				t.Errorf("goSymbol = %#v, want IBM from 5 with a plain string inside", got)
			}
		case *osub.IndexedEventSubscriptionSymbol:
			indexed, ok := got.(*osub.IndexedEventSubscriptionSymbol)
			if !ok || indexed.Symbol() != want.Symbol() || indexed.Source().Type() != want.Source().Type() ||
				indexed.Source().Id() != want.Source().Id() || *indexed.Source().Name() != *want.Source().Name() {
				t.Errorf("goSymbol = %#v, want %#v", got, want)
			}
		}
	}
}

func TestGoSymbolRejectsNil(t *testing.T) {
	if _, err := goSymbol(nil); err == nil {
		t.Error("goSymbol(nil): no error")
	}
}
