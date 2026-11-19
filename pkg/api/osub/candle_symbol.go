// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package osub

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"

// CandleSymbol is a candle symbol as a subscription symbol: it subscribes to the candles of the period, type and other
// attributes of the candle symbol, such as "AAPL{=5m}".
type CandleSymbol struct {
	symbol *candle.CandleSymbol
}

// Candle returns the candle symbol as a subscription symbol:
//
//	err := subscription.AddSymbol(osub.Candle(candle.NewCandleSymbol("AAPL{=5m}")))
//
// It can also be the symbol of a TimeSeriesSubscriptionSymbol. A nil candle symbol is rejected when it is added.
func Candle(symbol *candle.CandleSymbol) *CandleSymbol {
	return &CandleSymbol{symbol}
}

// CandleSymbol returns the wrapped candle symbol.
func (symbol *CandleSymbol) CandleSymbol() *candle.CandleSymbol {
	return symbol.symbol
}

// String returns the string form of the candle symbol, such as "AAPL{=5m}".
func (symbol *CandleSymbol) String() string {
	if symbol == nil || symbol.symbol == nil {
		return "<nil>"
	}
	return symbol.symbol.String()
}

func (*CandleSymbol) isSymbol() {}

var _ Symbol = (*CandleSymbol)(nil)
