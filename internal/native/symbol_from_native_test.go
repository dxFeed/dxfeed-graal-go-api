package native

import (
	"testing"
	"unsafe"

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

func TestGoSymbolRoundTripTimeSeriesOfWildcard(t *testing.T) {
	c, err := newCSymbol(osub.NewTimeSeriesSubscriptionSymbol(osub.NewWildcardSymbol(), 9))
	if err != nil {
		t.Fatalf("newCSymbol: %v", err)
	}
	defer freeCSymbol(c)
	got, err := goSymbol(c)
	if err != nil {
		t.Fatalf("goSymbol: %v", err)
	}
	ts, ok := got.(*osub.TimeSeriesSubscriptionSymbol)
	if !ok || ts.FromTime() != 9 {
		t.Fatalf("goSymbol = %#v, want a time series symbol from 9", got)
	}
	if _, ok := ts.Symbol().(*osub.WildcardSymbol); !ok {
		t.Errorf("the nested symbol = %#v, want the wildcard", ts.Symbol())
	}
}

func TestGoSymbolRejectsUnknownType(t *testing.T) {
	c, err := newCSymbol(osub.String("AAPL"))
	if err != nil {
		t.Fatalf("newCSymbol: %v", err)
	}
	known := c._type
	defer func() {
		c._type = known
		freeCSymbol(c)
	}()
	c._type = 1000
	if got, err := goSymbol(c); err == nil {
		t.Errorf("goSymbol(unknown type) = %#v, no error", got)
	}
}

func TestGoIndexedEventSourceErrors(t *testing.T) {
	if _, err := goIndexedEventSource(nil); err == nil {
		t.Error("goIndexedEventSource(nil): no error")
	}

	unknown, err := newCIndexedEventSource(events.NewIndexedEventSource(1, "SRC1"))
	if err != nil {
		t.Fatalf("newCIndexedEventSource: %v", err)
	}
	known := unknown._type
	unknown._type = 1000
	_, err = goIndexedEventSource(unknown)
	unknown._type = known
	freeCIndexedEventSource(unknown)
	if err == nil {
		t.Error("goIndexedEventSource(unknown source type): no error")
	}

	// An order source whose identifier does not decode to a name: 0x21 is "!", which is not alphanumeric, and it is
	// not a registered source.
	bad, err := newCIndexedEventSource(order.NtvL2())
	if err != nil {
		t.Fatalf("newCIndexedEventSource: %v", err)
	}
	bad.id = 0x21
	_, err = goIndexedEventSource(bad)
	freeCIndexedEventSource(bad)
	if err == nil {
		t.Error("goIndexedEventSource(unknown order source id): no error")
	}
}

func TestGoSymbolsList(t *testing.T) {
	list, err := newCSymbolList([]osub.Symbol{osub.String("AAPL"), osub.String("IBM")})
	if err != nil {
		t.Fatalf("newCSymbolList: %v", err)
	}
	defer freeCSymbolList(list)

	got, err := goSymbols(list)
	if err != nil || len(got) != 2 || got[0] != osub.String("AAPL") || got[1] != osub.String("IBM") {
		t.Fatalf("goSymbols = %#v, %v, want AAPL and IBM", got, err)
	}

	// The conversion stops at the first bad element and returns no partial result.
	elements := unsafe.Slice(list.elements, int(list.size))
	second := elements[1]
	known := second._type
	second._type = 1000
	got, err = goSymbols(list)
	second._type = known
	if err == nil || got != nil {
		t.Errorf("goSymbols(unknown type) = %#v, %v, want an error and no symbols", got, err)
	}

	elements[1] = nil
	got, err = goSymbols(list)
	if err == nil || got != nil {
		t.Errorf("goSymbols(nil element) = %#v, %v, want an error and no symbols", got, err)
	}
	elements[1] = second // freeCSymbolList frees it
}

func TestGoSymbolsEmpty(t *testing.T) {
	if got, err := goSymbols(nil); got != nil || err != nil {
		t.Errorf("goSymbols(nil) = %#v, %v, want nil, nil", got, err)
	}
	list, err := newCSymbolList(nil)
	if err != nil {
		t.Fatalf("newCSymbolList: %v", err)
	}
	defer freeCSymbolList(list)
	if got, err := goSymbols(list); got != nil || err != nil {
		t.Errorf("goSymbols(empty) = %#v, %v, want nil, nil", got, err)
	}
}
