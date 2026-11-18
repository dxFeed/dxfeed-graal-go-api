// Package candle provides the [Candle] and [DailyCandle] events, the open, high, low and close prices and the volume of
// a symbol over a period, and the [CandleSymbol] that names the aggregation.
//
// A candle symbol is a symbol with the attributes of the candle in braces, for example "AAPL{=1d}" for daily candles or
// "AAPL{=5m,tho=true}" for 5-minute candles of the regular trading hours. Candles are time series events: subscribe to
// them with an osub.TimeSeriesSubscriptionSymbol to get their history.
package candle
