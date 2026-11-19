// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package eventcodes defines the codes that identify dxFeed event types.
package eventcodes

import "strconv"

// EventCode identifies an event type; each value equals the matching dxfg_event_clazz_t constant of the native SDK.
type EventCode int32

const (
	// Quote is the event code of the Quote event.
	Quote EventCode = 0
	// Profile is the event code of the Profile event.
	Profile EventCode = 1
	// Summary is the event code of the Summary event.
	Summary EventCode = 2
	// Greeks is the event code of the Greeks event.
	Greeks EventCode = 3
	// Candle is the event code of the Candle event.
	Candle EventCode = 4
	// DailyCandle is the event code of the DailyCandle event.
	DailyCandle EventCode = 5
	// Underlying is the event code of the Underlying event.
	Underlying EventCode = 6
	// TheoPrice is the event code of the TheoPrice event.
	TheoPrice EventCode = 7
	// Trade is the event code of the Trade event.
	Trade EventCode = 8
	// TradeETH is the event code of the TradeETH event.
	TradeETH EventCode = 9
	// Configuration is the event code of the Configuration event.
	Configuration EventCode = 10
	// Message is the event code of the Message event.
	Message EventCode = 11
	// TimeAndSale is the event code of the TimeAndSale event.
	TimeAndSale EventCode = 12
	// OrderBase is the event code of the OrderBase event.
	OrderBase EventCode = 13
	// Order is the event code of the Order event.
	Order EventCode = 14
	// AnalyticOrder is the event code of the AnalyticOrder event.
	AnalyticOrder EventCode = 15
	// OtcMarketsOrder is the event code of the OtcMarketsOrder event.
	OtcMarketsOrder EventCode = 16
	// SpreadOrder is the event code of the SpreadOrder event.
	SpreadOrder EventCode = 17
	// Series is the event code of the Series event.
	Series EventCode = 18
	// OptionSale is the event code of the OptionSale event.
	OptionSale EventCode = 19
	// TextMessage is the event code of the TextMessage event.
	TextMessage EventCode = 20
	// MarketMaker is the event code of the MarketMaker event.
	MarketMaker EventCode = 21
	// TextConfiguration is the event code of the TextConfiguration event.
	TextConfiguration EventCode = 22
	// NuamOrder is the event code of the NuamOrder event.
	NuamOrder EventCode = 23
	// NuamTimeAndSale is the event code of the NuamTimeAndSale event.
	NuamTimeAndSale EventCode = 24
	// NuamTrade is the event code of the NuamTrade event.
	NuamTrade EventCode = 25
	// OrderImbalance is the event code of the OrderImbalance event.
	OrderImbalance EventCode = 26
)

// names holds the name of each event type, which is the name of its constant, by the code.
var names = [...]string{
	Quote:             "Quote",
	Profile:           "Profile",
	Summary:           "Summary",
	Greeks:            "Greeks",
	Candle:            "Candle",
	DailyCandle:       "DailyCandle",
	Underlying:        "Underlying",
	TheoPrice:         "TheoPrice",
	Trade:             "Trade",
	TradeETH:          "TradeETH",
	Configuration:     "Configuration",
	Message:           "Message",
	TimeAndSale:       "TimeAndSale",
	OrderBase:         "OrderBase",
	Order:             "Order",
	AnalyticOrder:     "AnalyticOrder",
	OtcMarketsOrder:   "OtcMarketsOrder",
	SpreadOrder:       "SpreadOrder",
	Series:            "Series",
	OptionSale:        "OptionSale",
	TextMessage:       "TextMessage",
	MarketMaker:       "MarketMaker",
	TextConfiguration: "TextConfiguration",
	NuamOrder:         "NuamOrder",
	NuamTimeAndSale:   "NuamTimeAndSale",
	NuamTrade:         "NuamTrade",
	OrderImbalance:    "OrderImbalance",
}

// String returns the name of the event type, such as "Quote" or "TimeAndSale", or "EventCode(N)" for a value that is
// not an event code.
func (e EventCode) String() string {
	if e >= 0 && int(e) < len(names) {
		return names[e]
	}
	return "EventCode(" + strconv.Itoa(int(e)) + ")"
}

// All returns every event code in the order of the values. OrderBase is the common base of the order events and has
// no events of its own, so a subscription to it fails.
func All() []EventCode {
	codes := make([]EventCode, len(names))
	for i := range names {
		codes[i] = EventCode(i)
	}
	return codes
}

// ByName returns the event code of the event type with the name, such as "Quote" or "TimeAndSale"; the name is
// case-sensitive. It reports false for an unknown name.
func ByName(name string) (EventCode, bool) {
	for i, n := range names {
		if n == name {
			return EventCode(i), true
		}
	}
	return 0, false
}

// NativeCode returns the numeric event class constant of the native SDK.
func (e EventCode) NativeCode() int32 {
	return int32(e)
}
