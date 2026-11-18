package api

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/marketmaker"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/misc"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/optionsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/orderimbalance"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/series"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/summary"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/theoprice"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
)

// roundTripTime is the time set on the events of the round trip test: a whole second, so that no part is lost.
var roundTripTime = time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC).UnixMilli()

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
		{eventcodes.DailyCandle, func(symbol string) events.EventType {
			e := candle.NewDailyCandle(symbol + "{=d}")
			e.SetTime(roundTripTime)
			e.SetCount(2)
			e.SetOpen(20.5)
			e.SetHigh(21.5)
			e.SetLow(19.5)
			e.SetClose(20.75)
			e.SetVolume(9)
			e.SetVwap(20.25)
			e.SetBidVolume(4)
			e.SetAskVolume(5)
			e.SetImpVolatility(0.3)
			e.SetOpenInterest(11)
			return e
		}},
		{eventcodes.Series, func(symbol string) events.EventType {
			e := series.NewSeries(symbol)
			e.SetEventFlags(events.TxPending.Flag())
			e.SetIndex(7)
			e.SetTime(roundTripTime)
			_ = e.SetSequence(3)
			e.SetExpiration(20600)
			e.SetVolatility(0.25)
			e.SetCallVolume(100)
			e.SetPutVolume(50)
			e.SetPutCallRatio(0.5)
			e.SetForwardPrice(101.5)
			e.SetDividend(0.01)
			e.SetInterest(0.02)
			return e
		}},
		{eventcodes.TheoPrice, func(symbol string) events.EventType {
			e := theoprice.NewTheoPrice(symbol)
			e.SetTime(roundTripTime)
			_ = e.SetSequence(4)
			e.SetPrice(2.5)
			e.SetUnderlyingPrice(100.5)
			e.SetDelta(0.4)
			e.SetGamma(0.05)
			e.SetDividend(0.01)
			e.SetInterest(0.02)
			return e
		}},
		{eventcodes.Message, func(symbol string) events.EventType {
			e := misc.NewMessage(symbol)
			e.SetAttachment(&description)
			return e
		}},
		{eventcodes.TextMessage, func(symbol string) events.EventType {
			e := misc.NewTextMessage(symbol)
			e.SetTime(roundTripTime)
			_ = e.SetSequence(5)
			e.SetText(&description)
			return e
		}},
		{eventcodes.Configuration, func(symbol string) events.EventType {
			e := misc.NewConfiguration(symbol)
			e.SetVersion(6)
			e.SetAttachment(&description)
			return e
		}},
		{eventcodes.TextConfiguration, func(symbol string) events.EventType {
			e := misc.NewTextConfiguration(symbol)
			e.SetTime(roundTripTime)
			_ = e.SetSequence(7)
			e.SetVersion(8)
			e.SetText(&description)
			return e
		}},
		{eventcodes.OptionSale, func(symbol string) events.EventType {
			e := optionsale.NewOptionSale(symbol)
			e.SetIndex(9)
			e.SetTime(roundTripTime)
			_ = e.SetSequence(6)
			e.SetExchangeCode('C')
			e.SetPrice(3.5)
			e.SetSize(10)
			e.SetBidPrice(3.25)
			e.SetAskPrice(3.75)
			// The data scheme keeps the sale conditions as a short string of up to 8 characters.
			conditions := "FTI"
			e.SetExchangeSaleConditions(&conditions)
			e.SetAggressorSide(side.Buy)
			e.SetIsValidTick(true)
			e.SetUnderlyingPrice(150.5)
			e.SetVolatility(0.35)
			e.SetDelta(0.45)
			e.SetOptionSymbol(&description)
			return e
		}},
		{eventcodes.MarketMaker, func(symbol string) events.EventType {
			e := marketmaker.NewMarketMaker(symbol)
			e.SetExchangeCode('Q')
			mmid := "NSDQ"
			_ = e.SetMarketMaker(&mmid)
			e.SetBidTime(roundTripTime)
			e.SetBidPrice(22.5)
			e.SetBidSize(100)
			e.SetBidCount(2)
			e.SetAskTime(roundTripTime + 1000)
			e.SetAskPrice(23.5)
			e.SetAskSize(200)
			e.SetAskCount(3)
			return e
		}},
		{eventcodes.OtcMarketsOrder, func(symbol string) events.EventType {
			e := order.NewOtcMarketsOrder(symbol)
			e.SetSide(side.Sell)
			e.SetPrice(24.5)
			e.SetSize(9)
			e.SetQuoteAccessPayment(-3)
			e.SetOtcMarketsPriceType(order.OtcMarketsPriceTypeActual)
			e.SetOpen(true)
			e.SetAutoExecution(true)
			return e
		}},
		{eventcodes.OrderImbalance, func(symbol string) events.EventType {
			e := orderimbalance.NewOrderImbalance(symbol)
			e.SetTime(roundTripTime)
			_ = e.SetSequence(8)
			e.SetRefPrice(25.5)
			e.SetPairedSize(1000)
			e.SetImbalanceSize(300)
			e.SetNearPrice(25.25)
			e.SetFarPrice(25.75)
			e.SetImbalanceSide(orderimbalance.ImbalanceSideBuy)
			e.SetAuctionType(orderimbalance.AuctionTypeClosing)
			return e
		}},
	}

	covered := map[int32]bool{}
	for _, c := range cases {
		covered[c.code.NativeCode()] = true
	}
	for code := range mappers.AvailableMappers() {
		if !covered[code] {
			t.Errorf("no round trip case for the mapper of native code %d", code)
		}
	}

	feed, publisher := newLocalHub(t)
	for _, c := range cases {
		// Each type gets its own symbol, so no event of an earlier case can reach a later one.
		symbol := fmt.Sprintf("ROUND/TRIP/%d", c.code)
		event := c.event(symbol)
		t.Run(fmt.Sprintf("%T", event), func(t *testing.T) {
			code, subSymbol, want := c.code, symbol, fmt.Sprint(event)
			if c.code == eventcodes.DailyCandle {
				// The hub delivers a published DailyCandle as a Candle with the same fields.
				code, subSymbol = eventcodes.Candle, symbol+"{=d}"
				want = strings.Replace(want, "DailyCandle{", "Candle{", 1)
			}
			sub, err := feed.CreateSubscription(code)
			if err != nil {
				t.Fatalf("create subscription: %v", err)
			}
			defer sub.Close()
			collector := &eventCollector{events: make(chan interface{}, 16)}
			if err := sub.AddListener(collector); err != nil {
				t.Fatalf("add listener: %v", err)
			}
			if err := sub.AddSymbol(osub.String(subSymbol)); err != nil {
				t.Fatalf("add symbol: %v", err)
			}
			if err := publisher.Publish([]interface{}{event}); err != nil {
				t.Fatalf("publish: %v", err)
			}
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
