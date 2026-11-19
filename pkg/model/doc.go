// Package model holds the models of the events of a feed, as the C++ and Java APIs have them: IndexedTxModel passes
// the complete transactions of the indexed events of a symbol, such as the orders of an order book, to its listener:
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
package model
