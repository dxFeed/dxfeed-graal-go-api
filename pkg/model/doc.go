// Package model holds the models of the events of a feed, as the C++ and Java APIs have them.
//
// IndexedTxModel passes the complete transactions of the indexed events of a symbol, such as the orders of an order
// book, to its listener:
//
//	m, err := model.NewIndexedTxModelBuilder[*order.Order]().
//		WithFeed(feed).
//		WithSymbol("AAPL").
//		WithSources(ntv).
//		WithListener(func(source events.IndexedEventSourceInterface, orders []*order.Order, isSnapshot bool) {
//			...
//		}).
//		Build()
//	...
//	defer m.Close()
//
// TimeSeriesTxModel does the same for the time series events of a symbol from a time, such as candles:
//
//	m, err := model.NewTimeSeriesTxModelBuilder[*candle.Candle]().
//		WithFeed(feed).
//		WithSymbol("AAPL{=1d}").
//		WithFromTime(from).
//		WithListener(func(candles []*candle.Candle, isSnapshot bool) {
//			...
//		}).
//		Build()
package model
