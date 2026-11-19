// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package candle provides the [Candle] and [DailyCandle] events, the open, high, low and close prices and the volume of
// a symbol over a period.
//
// A candle symbol is a string with the attributes of the candle in braces, for example "AAPL{=1d}" for daily candles
// or "AAPL{=5m,tho=true}" for 5-minute candles of the regular trading hours. Candles are time series events: subscribe
// with the string inside an osub.TimeSeriesSubscriptionSymbol, such as
// osub.NewTimeSeriesSubscriptionSymbol("AAPL{=1d}", fromTime), to get their history. [CandleSymbol] parses such a
// string and builds one from attributes ([NewCandleSymbolWithAttributes]); its String form is the symbol to subscribe
// to. [Candle.EventSymbol] of a received candle is its CandleSymbol.
package candle
