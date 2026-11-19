// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package trade

import (
	"math"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// NuamTrade is a Trade of the Nuam Exchange with its trade statistics.
type NuamTrade struct {
	Trade
	tradeStatTime        int64
	lastSignificantPrice float64
	lastPriceForAll      float64
	numberOfTrades       int32
	vwap                 float64
}

// NewNuamTrade creates a NuamTrade for the given event symbol.
func NewNuamTrade(eventSymbol string) *NuamTrade {
	return &NuamTrade{
		Trade:                *NewTrade(eventSymbol),
		lastSignificantPrice: math.NaN(),
		lastPriceForAll:      math.NaN(),
		vwap:                 math.NaN(),
	}
}

// Type returns the event code of this event.
func (t *NuamTrade) Type() eventcodes.EventCode {
	return eventcodes.NuamTrade
}

// TradeStatTime returns the time of the trade statistics in milliseconds since the Unix epoch.
func (t *NuamTrade) TradeStatTime() int64 {
	return t.tradeStatTime
}

// SetTradeStatTime sets the time of the trade statistics in milliseconds since the Unix epoch.
func (t *NuamTrade) SetTradeStatTime(tradeStatTime int64) {
	t.tradeStatTime = tradeStatTime
}

// LastSignificantPrice returns the last significant price.
func (t *NuamTrade) LastSignificantPrice() float64 {
	return t.lastSignificantPrice
}

// SetLastSignificantPrice sets the last significant price.
func (t *NuamTrade) SetLastSignificantPrice(lastSignificantPrice float64) {
	t.lastSignificantPrice = lastSignificantPrice
}

// LastPriceForAll returns the last price for all.
func (t *NuamTrade) LastPriceForAll() float64 {
	return t.lastPriceForAll
}

// SetLastPriceForAll sets the last price for all.
func (t *NuamTrade) SetLastPriceForAll(lastPriceForAll float64) {
	t.lastPriceForAll = lastPriceForAll
}

// NumberOfTrades returns the number of trades.
func (t *NuamTrade) NumberOfTrades() int32 {
	return t.numberOfTrades
}

// SetNumberOfTrades sets the number of trades.
func (t *NuamTrade) SetNumberOfTrades(numberOfTrades int32) {
	t.numberOfTrades = numberOfTrades
}

// Vwap returns the volume-weighted average price.
func (t *NuamTrade) Vwap() float64 {
	return t.vwap
}

// SetVwap sets the volume-weighted average price.
func (t *NuamTrade) SetVwap(vwap float64) {
	t.vwap = vwap
}

// String returns a string representation of this event.
func (t *NuamTrade) String() string {
	return "NuamTrade{" + t.TradeBase.String() +
		", tradeStatTime=" + formatutil.FormatTime(t.tradeStatTime) +
		", lastSignificantPrice=" + formatutil.FormatFloat64(t.lastSignificantPrice) +
		", lastPriceForAll=" + formatutil.FormatFloat64(t.lastPriceForAll) +
		", numberOfTrades=" + formatutil.FormatInt64(int64(t.numberOfTrades)) +
		", vwap=" + formatutil.FormatFloat64(t.vwap) +
		"}"
}
