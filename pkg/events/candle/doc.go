// Package candle provides the [Candle] and [DailyCandle] events, the open, high, low and close prices and the volume of
// a symbol over a period.
//
// A candle symbol is a string with the attributes of the candle in braces, for example "AAPL{=1d}" for daily candles
// or "AAPL{=5m,tho=true}" for 5-minute candles of the regular trading hours. Candles are time series events: subscribe
// with the string inside an osub.TimeSeriesSubscriptionSymbol, such as
// osub.NewTimeSeriesSubscriptionSymbol("AAPL{=1d}", fromTime), to get their history. [Candle.EventSymbol] of a
// received candle is the parsed [CandleSymbol].
package candle
