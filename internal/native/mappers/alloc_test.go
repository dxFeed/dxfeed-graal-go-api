package mappers_test

import (
	"fmt"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/misc"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/series"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/summary"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/theoprice"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
)

// eventsWithAllStrings returns one event of every type that has a mapper, with every string field set.
func eventsWithAllStrings() []events.EventType {
	str := func(s string) *string { return &s }
	o := order.NewOrder("ORDER")
	o.SetMarketMaker(str("MM"))
	ao := order.NewAnalyticOrder("ANALYTIC")
	ao.SetMarketMaker(str("MM"))
	so := order.NewSpreadOrder("SPREAD")
	so.SetSpreadSymbol(str("A-B"))
	p := profile.NewProfile("PROFILE")
	p.SetDescription(str("description"))
	p.SetStatusReason(str("reason"))
	ts := timeandsale.NewTimeAndSale("TNS")
	ts.SetExchangeSaleConditions(str("conditions"))
	ts.SetBuyer(str("buyer"))
	ts.SetSeller(str("seller"))
	msg := misc.NewMessage("MESSAGE")
	msg.SetAttachment(str("attachment"))
	tm := misc.NewTextMessage("TEXT_MESSAGE")
	tm.SetText(str("text"))
	cfg := misc.NewConfiguration("CONFIGURATION")
	cfg.SetAttachment(str("attachment"))
	tc := misc.NewTextConfiguration("TEXT_CONFIGURATION")
	tc.SetText(str("text"))
	return []events.EventType{
		candle.NewCandle("CANDLE{=d}"), greeks.NewGreeks("GREEKS"), ao, o, so, p, quote.NewQuote("QUOTE"),
		summary.NewSummary("SUMMARY"), ts, trade.NewTrade("TRADE"), trade.NewTradeETH("TRADE_ETH"),
		underlying.NewUnderlying("UNDERLYING"), candle.NewDailyCandle("DAILY{=d}"), series.NewSeries("SERIES"),
		theoprice.NewTheoPrice("THEO"), msg, tm, cfg, tc,
	}
}

// TestFreeCEventFreesAllAllocations checks that FreeCEvent frees everything CEvent allocated, for every mapper.
func TestFreeCEventFreesAllAllocations(t *testing.T) {
	available := mappers.AvailableMappers()
	covered := map[int32]bool{}
	for _, event := range eventsWithAllStrings() {
		code := event.Type().NativeCode()
		mapper := available[code]
		if mapper == nil {
			t.Fatalf("no mapper for %T", event)
		}
		covered[code] = true
		before := mappers.LiveAllocations()
		native := mapper.CEvent(event)
		if mappers.LiveAllocations() == before {
			t.Errorf("%T: CEvent allocated nothing", event)
		}
		mapper.FreeCEvent(native)
		if after := mappers.LiveAllocations(); after != before {
			t.Errorf("%T: %d allocations left after FreeCEvent", event, after-before)
		}
	}
	for code := range available {
		if !covered[code] {
			t.Errorf("no event for the mapper of native code %d", code)
		}
	}
}

// TestMapperRoundTrip maps every event to C memory and back and checks that no field is lost; this covers the GoEvent
// side of the mappers whose events a local hub delivers as another type, such as DailyCandle.
func TestMapperRoundTrip(t *testing.T) {
	available := mappers.AvailableMappers()
	for _, event := range eventsWithAllStrings() {
		mapper := available[event.Type().NativeCode()]
		native := mapper.CEvent(event)
		got := mapper.GoEvent(native)
		mapper.FreeCEvent(native)
		if fmt.Sprint(got) != fmt.Sprint(event) {
			t.Errorf("%T changed in the round trip:\n got %v\nwant %v", event, got, event)
		}
	}
}
