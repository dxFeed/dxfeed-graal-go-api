package api

import (
	"sort"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

func newSubscription(t *testing.T, eventTypes ...eventcodes.EventCode) *DXFeedSubscription {
	t.Helper()
	feed, _ := newLocalHub(t)
	sub, err := feed.CreateSubscription(eventTypes...)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	t.Cleanup(func() { _ = sub.Close() })
	return sub
}

// stringSymbols returns the osub.String symbols as sorted strings and fails on a symbol of another kind.
func stringSymbols(t *testing.T, symbols []osub.Symbol) []string {
	t.Helper()
	result := make([]string, 0, len(symbols))
	for _, s := range symbols {
		text, ok := s.(osub.String)
		if !ok {
			t.Fatalf("symbol %#v is not an osub.String", s)
		}
		result = append(result, string(text))
	}
	sort.Strings(result)
	return result
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSubscriptionSetAndGetSymbols(t *testing.T) {
	sub := newSubscription(t, eventcodes.Quote)
	symbols, err := sub.Symbols()
	if err != nil || len(symbols) != 0 {
		t.Fatalf("Symbols() of a new subscription = %v, %v; want none", symbols, err)
	}
	if err := sub.AddSymbols(osub.Strings("AAPL", "IBM", "AAPL{=d}")...); err != nil {
		t.Fatalf("AddSymbols: %v", err)
	}
	symbols, err = sub.Symbols()
	if err != nil {
		t.Fatalf("Symbols: %v", err)
	}
	if got, want := stringSymbols(t, symbols), []string{"AAPL", "AAPL{=d}", "IBM"}; !equalStrings(got, want) {
		t.Errorf("Symbols() = %v, want %v", got, want)
	}
	decorated, err := sub.DecoratedSymbols()
	if err != nil {
		t.Fatalf("DecoratedSymbols: %v", err)
	}
	if got, want := stringSymbols(t, decorated), []string{"AAPL", "AAPL{=d}", "IBM"}; !equalStrings(got, want) {
		t.Errorf("DecoratedSymbols() = %v, want %v", got, want)
	}

	if err := sub.SetSymbols(osub.Strings("MSFT")...); err != nil {
		t.Fatalf("SetSymbols: %v", err)
	}
	symbols, err = sub.Symbols()
	if err != nil {
		t.Fatalf("Symbols: %v", err)
	}
	if got := stringSymbols(t, symbols); !equalStrings(got, []string{"MSFT"}) {
		t.Errorf("Symbols() after SetSymbols = %v, want [MSFT]", got)
	}
	if err := sub.SetSymbols(); err != nil {
		t.Fatalf("SetSymbols(): %v", err)
	}
	if symbols, err = sub.Symbols(); err != nil || len(symbols) != 0 {
		t.Errorf("Symbols() after SetSymbols() = %v, %v; want none", symbols, err)
	}
}

// TestSubscriptionSymbolKindsRoundTrip checks that every symbol kind comes back from the SDK as it was added.
func TestSubscriptionSymbolKindsRoundTrip(t *testing.T) {
	sub := newSubscription(t, eventcodes.Order, eventcodes.TimeAndSale)
	err := sub.SetSymbols(
		osub.NewWildcardSymbol(),
		osub.NewTimeSeriesSubscriptionSymbol("IBM", 1000),
		osub.NewIndexedEventSubscriptionSymbol(osub.String("MSFT"), order.NtvL2()),
		osub.NewIndexedEventSubscriptionSymbol("AAPL", events.NewIndexedEventSource(77, "SRC77")),
	)
	if err != nil {
		t.Fatalf("SetSymbols: %v", err)
	}
	symbols, err := sub.Symbols()
	if err != nil {
		t.Fatalf("Symbols: %v", err)
	}
	if len(symbols) != 4 {
		t.Fatalf("Symbols() = %#v; want 4 symbols", symbols)
	}
	var wildcard, timeSeries, orderSource, indexedSource bool
	for _, s := range symbols {
		switch v := s.(type) {
		case *osub.WildcardSymbol:
			wildcard = true
		case *osub.TimeSeriesSubscriptionSymbol:
			timeSeries = v.Symbol() == "IBM" && v.FromTime() == 1000
		case *osub.IndexedEventSubscriptionSymbol:
			source := v.Source()
			name := source.Name()
			if name == nil {
				t.Fatalf("source %#v without a name", source)
			}
			switch {
			case v.Symbol() == "MSFT" && source.Type() == events.OrderSourceType:
				orderSource = source.Id() == order.NtvL2().Id() && *name == *order.NtvL2().Name()
			case v.Symbol() == "AAPL" && source.Type() == events.IndexedEventSourceType:
				indexedSource = source.Id() == 77 && *name == "SRC77"
			}
		default:
			t.Errorf("unexpected symbol %#v", s)
		}
	}
	if !wildcard || !timeSeries || !orderSource || !indexedSource {
		t.Errorf("Symbols() = %#v; wildcard %v, time series %v, order source %v, indexed source %v", symbols,
			wildcard, timeSeries, orderSource, indexedSource)
	}
}

func TestSubscriptionEventTypes(t *testing.T) {
	sub := newSubscription(t, eventcodes.Quote, eventcodes.Trade)
	types, err := sub.EventTypes()
	if err != nil {
		t.Fatalf("EventTypes: %v", err)
	}
	if len(types) != 2 || !(types[0] == eventcodes.Quote && types[1] == eventcodes.Trade ||
		types[0] == eventcodes.Trade && types[1] == eventcodes.Quote) {
		t.Errorf("EventTypes() = %v, want Quote and Trade", types)
	}
	for code, want := range map[eventcodes.EventCode]bool{eventcodes.Quote: true, eventcodes.Trade: true,
		eventcodes.Candle: false} {
		if got, err := sub.ContainsEventType(code); err != nil || got != want {
			t.Errorf("ContainsEventType(%v) = %v, %v; want %v", code, got, err, want)
		}
	}
}

func TestSubscriptionAggregationPeriod(t *testing.T) {
	sub := newSubscription(t, eventcodes.Quote)
	if period, ok, err := sub.AggregationPeriod(); err != nil || ok || period != 0 {
		t.Fatalf("AggregationPeriod() of a new subscription = %v, %v, %v; want none", period, ok, err)
	}
	if err := sub.SetAggregationPeriod(1500 * time.Millisecond); err != nil {
		t.Fatalf("SetAggregationPeriod: %v", err)
	}
	if period, ok, err := sub.AggregationPeriod(); err != nil || !ok || period != 1500*time.Millisecond {
		t.Errorf("AggregationPeriod() = %v, %v, %v; want 1.5s", period, ok, err)
	}
	if err := sub.SetAggregationPeriod(time.Microsecond); err != nil {
		t.Fatalf("SetAggregationPeriod(1µs): %v", err)
	}
	if period, _, err := sub.AggregationPeriod(); err != nil || period != time.Millisecond {
		t.Errorf("AggregationPeriod() after 1µs = %v, %v; want 1ms", period, err)
	}
	for _, bad := range []time.Duration{-time.Second, 25 * time.Hour} {
		if err := sub.SetAggregationPeriod(bad); err == nil {
			t.Errorf("SetAggregationPeriod(%v): no error", bad)
		}
	}
	if err := sub.ClearAggregationPeriod(); err != nil {
		t.Fatalf("ClearAggregationPeriod: %v", err)
	}
	if _, ok, err := sub.AggregationPeriod(); err != nil || ok {
		t.Errorf("AggregationPeriod() after Clear: ok %v, %v; want none", ok, err)
	}
	if info, err := sub.AggregationPeriodInfo(); err != nil || info.Known {
		t.Errorf("AggregationPeriodInfo() before any events = %+v, %v; want unknown", info, err)
	}
}

func TestSubscriptionEventsBatchLimit(t *testing.T) {
	sub := newSubscription(t, eventcodes.Quote)
	if limit, err := sub.EventsBatchLimit(); err != nil || limit != EventsBatchLimitOptimal {
		t.Fatalf("EventsBatchLimit() of a new subscription = %d, %v; want optimal", limit, err)
	}
	for _, limit := range []int32{10, EventsBatchLimitMax, EventsBatchLimitOptimal} {
		if err := sub.SetEventsBatchLimit(limit); err != nil {
			t.Fatalf("SetEventsBatchLimit(%d): %v", limit, err)
		}
		if got, err := sub.EventsBatchLimit(); err != nil || got != limit {
			t.Errorf("EventsBatchLimit() = %d, %v; want %d", got, err, limit)
		}
	}
	if err := sub.SetEventsBatchLimit(-1); err == nil {
		t.Error("SetEventsBatchLimit(-1): no error")
	}
	if dropped, err := sub.TotalDroppedEvents(); err != nil || dropped != 0 {
		t.Errorf("TotalDroppedEvents() = %d, %v; want 0", dropped, err)
	}
}

// TestSubscriptionExecutor checks that the listeners of a subscription run in the executor set on it.
func TestSubscriptionExecutor(t *testing.T) {
	feed, publisher := newLocalHub(t)
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer sub.Close()
	executor, err := NewInPlaceExecutor()
	if err != nil {
		t.Fatalf("NewInPlaceExecutor: %v", err)
	}
	defer executor.Close()
	if err := sub.SetExecutor(executor); err != nil {
		t.Fatalf("SetExecutor: %v", err)
	}
	collector := &eventCollector{events: make(chan interface{}, 16)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatalf("AddListener: %v", err)
	}
	if err := sub.AddSymbol(osub.String("EXECUTOR")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	if err := publisher.Publish([]interface{}{quote.NewQuote("EXECUTOR")}); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	// The events wait in the executor's queue until ProcessAllPendingTasks runs them.
	select {
	case got := <-collector.events:
		t.Fatalf("event delivered before ProcessAllPendingTasks: %v", got)
	case <-time.After(500 * time.Millisecond):
	}
	for deadline := time.Now().Add(10 * time.Second); len(collector.events) == 0; {
		if err := executor.ProcessAllPendingTasks(); err != nil {
			t.Fatalf("ProcessAllPendingTasks: %v", err)
		}
		if time.Now().After(deadline) {
			t.Fatal("no event within 10 s of ProcessAllPendingTasks")
		}
		time.Sleep(10 * time.Millisecond)
	}

	closed, err := NewInPlaceExecutor()
	if err != nil {
		t.Fatalf("NewInPlaceExecutor: %v", err)
	}
	_ = closed.Close()
	if err := sub.SetExecutor(closed); err == nil {
		t.Error("SetExecutor of a closed executor: no error")
	}
	if err := sub.SetExecutor(nil); err != nil {
		t.Errorf("SetExecutor(nil): %v", err)
	}
}
