// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func TestEventCodesMatchSdk(t *testing.T) {
	codes := map[string]eventcodes.EventCode{
		"QUOTE":              eventcodes.Quote,
		"PROFILE":            eventcodes.Profile,
		"SUMMARY":            eventcodes.Summary,
		"GREEKS":             eventcodes.Greeks,
		"CANDLE":             eventcodes.Candle,
		"DAILY_CANDLE":       eventcodes.DailyCandle,
		"UNDERLYING":         eventcodes.Underlying,
		"THEO_PRICE":         eventcodes.TheoPrice,
		"TRADE":              eventcodes.Trade,
		"TRADE_ETH":          eventcodes.TradeETH,
		"CONFIGURATION":      eventcodes.Configuration,
		"MESSAGE":            eventcodes.Message,
		"TIME_AND_SALE":      eventcodes.TimeAndSale,
		"ORDER_BASE":         eventcodes.OrderBase,
		"ORDER":              eventcodes.Order,
		"ANALYTIC_ORDER":     eventcodes.AnalyticOrder,
		"OTC_MARKETS_ORDER":  eventcodes.OtcMarketsOrder,
		"SPREAD_ORDER":       eventcodes.SpreadOrder,
		"SERIES":             eventcodes.Series,
		"OPTION_SALE":        eventcodes.OptionSale,
		"TEXT_MESSAGE":       eventcodes.TextMessage,
		"MARKET_MAKER":       eventcodes.MarketMaker,
		"TEXT_CONFIGURATION": eventcodes.TextConfiguration,
		"NUAM_ORDER":         eventcodes.NuamOrder,
		"NUAM_TIME_AND_SALE": eventcodes.NuamTimeAndSale,
		"NUAM_TRADE":         eventcodes.NuamTrade,
		"ORDER_IMBALANCE":    eventcodes.OrderImbalance,
	}
	sdk := sdkEventClazzes()
	for name, value := range sdk {
		code, ok := codes[name]
		if !ok {
			t.Errorf("DXFG_EVENT_%s = %d has no Go event code", name, value)
			continue
		}
		if code.NativeCode() != value {
			t.Errorf("Go code for DXFG_EVENT_%s is %d, SDK value is %d", name, code.NativeCode(), value)
		}
	}
	for name := range codes {
		if _, ok := sdk[name]; !ok {
			t.Errorf("Go event code %s has no SDK constant", name)
		}
	}
}

// TestSdkEventClazzesCoverHeader fails when a new SDK adds an event class that sdkEventClazzes does not list yet.
func TestSdkEventClazzesCoverHeader(t *testing.T) {
	header, err := os.ReadFile(filepath.Join("include", "dxfg_events.h"))
	if err != nil {
		t.Fatal(err)
	}
	enum := regexp.MustCompile(`(?s)typedef enum dxfg_event_clazz_t \{(.*?)\} dxfg_event_clazz_t;`).FindSubmatch(header)
	if enum == nil {
		t.Fatal("dxfg_event_clazz_t enum not found in dxfg_events.h")
	}
	names := regexp.MustCompile(`(?m)^\s*DXFG_EVENT_(\w+)`).FindAllSubmatch(enum[1], -1)
	sdk := sdkEventClazzes()
	if len(names) != len(sdk) {
		t.Errorf("dxfg_event_clazz_t has %d constants, sdkEventClazzes lists %d", len(names), len(sdk))
	}
	for _, name := range names {
		if _, ok := sdk[string(name[1])]; !ok {
			t.Errorf("DXFG_EVENT_%s is missing from sdkEventClazzes", name[1])
		}
	}
}
