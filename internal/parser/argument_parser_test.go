// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package parser

import (
	"reflect"
	"slices"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func TestParseEventTypes(t *testing.T) {
	for value, want := range map[string][]eventcodes.EventCode{
		"Quote":                         {eventcodes.Quote},
		"quote, TIMEANDSALE":            {eventcodes.Quote, eventcodes.TimeAndSale},
		"DailyCandle,Candle":            {eventcodes.DailyCandle, eventcodes.Candle},
		"OrderImbalance,nuamtrade":      {eventcodes.OrderImbalance, eventcodes.NuamTrade},
		"TradeETH,tradeeth,TextMessage": {eventcodes.TradeETH, eventcodes.TradeETH, eventcodes.TextMessage},
	} {
		got, err := ParseEventTypes(value)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("ParseEventTypes(%q) = %v, %v; want %v", value, got, err, want)
		}
	}
}

func TestParseEventTypesAll(t *testing.T) {
	for _, value := range []string{"all", "ALL", " All "} {
		got, err := ParseEventTypes(value)
		if err != nil {
			t.Fatalf("ParseEventTypes(%q): %v", value, err)
		}
		if len(got) != len(eventcodes.All())-1 || slices.Contains(got, eventcodes.OrderBase) ||
			!slices.Contains(got, eventcodes.DailyCandle) {
			t.Errorf("ParseEventTypes(%q) = %v, want every type except OrderBase", value, got)
		}
	}
}

func TestParseEventTypesUnknown(t *testing.T) {
	for value, want := range map[string]string{
		"Qoute":             `unknown event type "Qoute"`,
		"Quote,Qoute,Trdae": `unknown event type "Qoute", "Trdae"`,
		"Quote,":            `unknown event type ""`,
		"":                  `unknown event type ""`,
	} {
		got, err := ParseEventTypes(value)
		if err == nil || err.Error() != want {
			t.Errorf("ParseEventTypes(%q) = %v, %v; want error %s", value, got, err, want)
		}
	}
}
