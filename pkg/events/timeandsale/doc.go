// Package timeandsale provides the [TimeAndSale] event: a trade or another price-bearing event of the market (a
// correction or a cancellation, see [Type]), with its exchange, buyer and seller. TimeAndSale is a time series event:
// subscribe with an osub.TimeSeriesSubscriptionSymbol to get the history of a day.
package timeandsale
