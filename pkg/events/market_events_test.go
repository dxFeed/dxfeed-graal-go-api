// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package events_test

import (
	"strings"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/marketmaker"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/optionsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/orderimbalance"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
)

func TestMarketMakerIndex(t *testing.T) {
	m := marketmaker.NewMarketMaker("AAPL")
	if m.MarketMaker() != nil {
		t.Fatalf("new market maker ID = %q, want nil", *m.MarketMaker())
	}
	m.SetExchangeCode('Q')
	mmid := "NSDQ"
	if err := m.SetMarketMaker(&mmid); err != nil {
		t.Fatal(err)
	}
	// 'Q' in bits 32 to 47, "NSDQ" one byte per character in bits 0 to 31.
	if want := int64(0x51_4E534451); m.Index() != want {
		t.Fatalf("index = %#x, want %#x", m.Index(), want)
	}
	if m.ExchangeCode() != 'Q' || *m.MarketMaker() != "NSDQ" {
		t.Fatalf("exchange = %c, market maker = %q", m.ExchangeCode(), *m.MarketMaker())
	}

	// Only the last 4 characters are stored, as in the Java API; the exchange code is kept.
	long := "ABCDEF"
	if err := m.SetMarketMaker(&long); err != nil {
		t.Fatal(err)
	}
	if *m.MarketMaker() != "CDEF" || m.ExchangeCode() != 'Q' {
		t.Fatalf("market maker = %q, exchange = %c", *m.MarketMaker(), m.ExchangeCode())
	}

	for _, bad := range []string{"ABCDEFGHI", "Ж"} {
		if err := m.SetMarketMaker(&bad); err == nil {
			t.Errorf("SetMarketMaker(%q) returned no error", bad)
		}
	}
	if err := m.SetMarketMaker(nil); err != nil || m.MarketMaker() != nil {
		t.Fatalf("SetMarketMaker(nil): err %v, market maker %v", err, m.MarketMaker())
	}
	if err := m.SetIndex(-1); err == nil {
		t.Error("SetIndex(-1) returned no error")
	}

	m.SetBidTime(1000)
	m.SetAskTime(2000)
	if m.Time() != 2000 {
		t.Errorf("time = %d, want the later of bid and ask time", m.Time())
	}
}

func TestOtcMarketsOrderFlags(t *testing.T) {
	o := order.NewOtcMarketsOrder("AAPL")
	o.SetOpen(true)
	o.SetUnsolicited(true)
	o.SetOtcMarketsPriceType(order.OtcMarketsPriceTypeWanted)
	o.SetSaturated(true)
	o.SetAutoExecution(true)
	o.SetNmsConditional(true)
	if want := int32(1 | 2 | 2<<2 | 1<<4 | 1<<5 | 1<<6); o.OtcMarketsFlags() != want {
		t.Fatalf("flags = %#x, want %#x", o.OtcMarketsFlags(), want)
	}
	o.SetUnsolicited(false)
	o.SetOtcMarketsPriceType(order.OtcMarketsPriceTypeActual)
	if o.IsUnsolicited() || !o.IsOpen() || o.OtcMarketsPriceType() != order.OtcMarketsPriceTypeActual ||
		!o.IsSaturated() || !o.IsAutoExecution() || !o.IsNmsConditional() {
		t.Fatalf("flags = %#x after the changes", o.OtcMarketsFlags())
	}
	if s := o.String(); !strings.HasPrefix(s, "OtcMarketsOrder{AAPL") || !strings.Contains(s, "priceType=ACTUAL") {
		t.Errorf("String() = %s", s)
	}
}

func TestOrderImbalanceFlags(t *testing.T) {
	o := orderimbalance.NewOrderImbalance("AAPL")
	o.SetImbalanceSide(orderimbalance.ImbalanceSidePaused)
	o.SetAuctionType(orderimbalance.AuctionTypeReopening)
	if want := int32(5 | 3<<3); o.Flags() != want {
		t.Fatalf("flags = %#x, want %#x", o.Flags(), want)
	}
	ntv, err := order.ValueOfName("NTV")
	if err != nil {
		t.Fatal(err)
	}
	o.SetSource(ntv)
	if s := o.String(); !strings.Contains(s, "source=NTV") || !strings.Contains(s, "imbalanceSide=PAUSED") ||
		!strings.Contains(s, "auctionType=REOPENING") {
		t.Errorf("String() = %s", s)
	}
	if err := o.SetSequence(-1); err == nil {
		t.Error("SetSequence(-1) returned no error")
	}
}

func TestOptionSaleFlags(t *testing.T) {
	o := optionsale.NewOptionSale("AAPL")
	o.SetTradeThroughExempt('X')
	o.SetAggressorSide(side.Sell)
	o.SetIsSpreadLeg(true)
	o.SetIsExtendedTradingHours(true)
	o.SetIsValidTick(true)
	o.SetSaleType(timeandsale.TypeCancel)
	// The same layout as the flags of TimeAndSale.
	ts := timeandsale.NewTimeAndSale("AAPL")
	ts.SetTradeThroughExempt('X')
	ts.SetAggressorSide(side.Sell)
	ts.SetIsSpreadLeg(true)
	ts.SetIsExtendedTradingHours(true)
	ts.SetIsValidTick(true)
	ts.SetTimeAndSaleType(timeandsale.TypeCancel)
	if o.Flags() != ts.Flags() {
		t.Fatalf("flags = %#x, TimeAndSale flags = %#x", o.Flags(), ts.Flags())
	}
	o.SetIsSpreadLeg(false)
	if o.IsSpreadLeg() || !o.IsCancel() || o.AggressorSide() != side.Sell || o.TradeThroughExempt() != 'X' {
		t.Fatalf("flags = %#x after the changes", o.Flags())
	}
	o.SetTimeNanos(1_700_000_000_123_456_789)
	if o.TimeNanos() != 1_700_000_000_123_456_789 || o.TimeNanoPart() != 456_789 {
		t.Errorf("time nanos = %d, nano part = %d", o.TimeNanos(), o.TimeNanoPart())
	}
}

func TestNuamOrderFlags(t *testing.T) {
	o := order.NewNuamOrder("AAPL")
	o.SetNuamOrderType(order.NuamOrderTypeImbalance)
	o.SetTimeInForce(order.NuamTimeInForceTypeNumberOfDays)
	if want := int32(5 | 6<<4); o.NuamFlags() != want {
		t.Fatalf("flags = %#x, want %#x", o.NuamFlags(), want)
	}
	o.SetNuamOrderType(order.NuamOrderTypeLimit)
	if o.NuamOrderType() != order.NuamOrderTypeLimit || o.TimeInForce() != order.NuamTimeInForceTypeNumberOfDays {
		t.Fatalf("flags = %#x after the change", o.NuamFlags())
	}
	if s := o.String(); !strings.HasPrefix(s, "NuamOrder{AAPL") || !strings.Contains(s, "timeInForce=NUMBER_OF_DAYS") ||
		!strings.Contains(s, "orderType=LIMIT") {
		t.Errorf("String() = %s", s)
	}
}

func TestNuamTimeAndSaleAndTradeStrings(t *testing.T) {
	ts := timeandsale.NewNuamTimeAndSale("AAPL")
	ts.SetMatchId(42)
	if s := ts.String(); !strings.HasPrefix(s, "NuamTimeAndSale{AAPL") || !strings.HasSuffix(s, ", matchId=42}") {
		t.Errorf("String() = %s", s)
	}
	tr := trade.NewNuamTrade("AAPL")
	tr.SetNumberOfTrades(3)
	if s := tr.String(); !strings.HasPrefix(s, "NuamTrade{AAPL") || !strings.Contains(s, "numberOfTrades=3") {
		t.Errorf("String() = %s", s)
	}
}
