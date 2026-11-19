package orcs_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/orcs"
)

// unreachable is an address without a server; the client connects in the background, so creating it succeeds.
const unreachable = "127.0.0.1:1"

func ntv(t *testing.T) *order.Source {
	t.Helper()
	source, err := order.ValueOfName("NTV")
	if err != nil {
		t.Fatal(err)
	}
	return source
}

func priceLevel(t *testing.T, time int64, s side.Side, price, size float64) *order.Order {
	t.Helper()
	o := order.NewOrder("AAPL")
	o.SetOrderSource(ntv(t))
	o.SetTime(time)
	o.SetSide(s)
	o.SetPrice(price)
	o.SetSize(size)
	return o
}

func TestGranularityAttributeKey(t *testing.T) {
	if orcs.GranularityAttributeKey != "gr" {
		t.Fatalf("GranularityAttributeKey = %q, want \"gr\"", orcs.GranularityAttributeKey)
	}
	symbol := candle.NewCandleSymbol("AAPL{=d,gr=s}")
	if symbol.BaseSymbol() != "AAPL" {
		t.Errorf("BaseSymbol = %q, want AAPL", symbol.BaseSymbol())
	}
}

func TestValidatePriceLevels(t *testing.T) {
	tests := []struct {
		name   string
		orders []*order.Order
		bound  time.Duration
		want   bool
	}{
		{"the bid below the ask", []*order.Order{
			priceLevel(t, 1_000, side.Buy, 100, 10),
			priceLevel(t, 1_000, side.Sell, 101, 10),
			priceLevel(t, 2_000, side.Buy, 100.5, 5),
		}, time.Minute, true},
		// The orders of one time are checked when the orders of the next time begin.
		{"the bid above the ask", []*order.Order{
			priceLevel(t, 1_000, side.Buy, 102, 10),
			priceLevel(t, 1_000, side.Sell, 101, 10),
			priceLevel(t, 2_000, side.Buy, 100, 5),
		}, time.Minute, false},
		{"a gap above the bound", []*order.Order{
			priceLevel(t, 1_000, side.Buy, 100, 10),
			priceLevel(t, 1_000+2*60_000, side.Sell, 101, 10),
		}, time.Minute, false},
		{"a gap within the bound", []*order.Order{
			priceLevel(t, 1_000, side.Buy, 100, 10),
			priceLevel(t, 1_000+2*60_000, side.Sell, 101, 10),
		}, 5 * time.Minute, true},
		{"no orders", nil, time.Minute, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid, err := orcs.ValidatePriceLevels(test.orders, test.bound, false)
			if err != nil || valid != test.want {
				t.Fatalf("ValidatePriceLevels = %v, %v; want %v", valid, err, test.want)
			}
		})
	}
	// printQuotes only adds the quotes to the log of the SDK.
	if valid, err := orcs.ValidatePriceLevels(tests[0].orders, time.Minute, true); err != nil || !valid {
		t.Errorf("ValidatePriceLevels with printQuotes = %v, %v; want true", valid, err)
	}
	if _, err := orcs.ValidatePriceLevels([]*order.Order{priceLevel(t, 1_000, side.Buy, 100, 10), nil}, time.Minute,
		false); err == nil {
		t.Error("ValidatePriceLevels with a nil order: no error")
	}
}

func TestAuthOrderSource(t *testing.T) {
	nyse, err := order.ValueOfName("NYSE")
	if err != nil {
		t.Fatal(err)
	}
	aos, err := orcs.NewAuthOrderSource(map[int32][]string{
		int32(ntv(t).Id()): {"MSFT", "AAPL"},
		int32(nyse.Id()):   {"IBM"},
	})
	if err != nil {
		t.Fatalf("NewAuthOrderSource: %v", err)
	}
	wantByIDs := map[int32][]string{int32(ntv(t).Id()): {"AAPL", "MSFT"}, int32(nyse.Id()): {"IBM"}}
	if got := aos.ByIDs(); !reflect.DeepEqual(got, wantByIDs) {
		t.Errorf("ByIDs = %v, want %v", got, wantByIDs)
	}
	bySources := aos.ByOrderSources()
	if len(bySources) != 2 {
		t.Fatalf("ByOrderSources = %v, want 2 entries", bySources)
	}
	for _, e := range bySources {
		want, ok := wantByIDs[int32(e.Source.Id())]
		if !ok || !reflect.DeepEqual(e.Symbols, want) {
			t.Errorf("ByOrderSources: %v has %v, want %v", e.Source, e.Symbols, want)
		}
	}
	if bySources[0].Source.Id() > bySources[1].Source.Id() {
		t.Errorf("ByOrderSources is not ordered by the identifier: %v", bySources)
	}
	// The results are copies.
	aos.ByIDs()[int32(nyse.Id())][0] = "changed"
	if got := aos.ByIDs()[int32(nyse.Id())][0]; got != "IBM" {
		t.Errorf("ByIDs changed through a copy: %q", got)
	}

	empty, err := orcs.NewAuthOrderSource(nil)
	if err != nil || len(empty.ByIDs()) != 0 || len(empty.ByOrderSources()) != 0 {
		t.Errorf("NewAuthOrderSource(nil) = %v, %v; want an empty AuthOrderSource", empty, err)
	}
}

func TestPriceLevelServiceCreateAndClose(t *testing.T) {
	service, err := orcs.NewPriceLevelService(unreachable)
	if err != nil {
		t.Fatalf("NewPriceLevelService: %v", err)
	}
	if err := service.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := service.Close(); err != nil {
		t.Fatalf("Close again: %v", err)
	}
	symbol := candle.NewCandleSymbol("AAPL{=d,gr=s}")
	now := time.Now()
	if _, err := service.GetOrders(symbol, ntv(t), now.Add(-time.Hour), now); err == nil {
		t.Error("GetOrders after Close: no error")
	}
	if _, err := service.GetQuotes(symbol, now.Add(-time.Hour), now); err == nil {
		t.Error("GetQuotes after Close: no error")
	}
	if _, err := service.GetAuthOrderSource(); err == nil {
		t.Error("GetAuthOrderSource after Close: no error")
	}
}

func TestPriceLevelServiceNilArguments(t *testing.T) {
	service, err := orcs.NewPriceLevelService(unreachable)
	if err != nil {
		t.Fatalf("NewPriceLevelService: %v", err)
	}
	defer service.Close()
	symbol := candle.NewCandleSymbol("AAPL{=d,gr=s}")
	now := time.Now()
	if _, err := service.GetOrders(nil, ntv(t), now, now); err == nil {
		t.Error("GetOrders with a nil symbol: no error")
	}
	if _, err := service.GetOrdersWithCaller(symbol, nil, now, now, "test"); err == nil {
		t.Error("GetOrdersWithCaller with a nil source: no error")
	}
	if _, err := service.GetQuotesWithCaller(nil, now, now, "test"); err == nil {
		t.Error("GetQuotesWithCaller with a nil symbol: no error")
	}
	if _, err := service.GetOrders(symbol, ntv(t), now, now.Add(-time.Hour)); err == nil {
		t.Error("GetOrders with the end before the start: no error")
	}
	if _, err := service.GetQuotes(symbol, now, now.Add(-time.Hour)); err == nil {
		t.Error("GetQuotes with the end before the start: no error")
	}
}
