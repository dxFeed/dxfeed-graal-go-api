// Package osub holds the subscription symbols: [String] for a plain symbol such as "AAPL", [CandleSymbol] for a
// candle symbol, the wildcard ([WildcardSymbol]), and the time series ([TimeSeriesSubscriptionSymbol]) and indexed event
// ([IndexedEventSubscriptionSymbol]) symbols. They all implement [Symbol], the type that
// [github.com/dxfeed/dxfeed-graal-go-api/pkg/api.DXFeedSubscription.AddSymbols] and the other subscription methods take.
//
// # Symbol kinds
//
//   - [String]: a symbol of the dxFeed symbology, such as "AAPL", "ETH/USD:GDAX" or a candle symbol "AAPL{=1d}".
//   - [CandleSymbol]: a [github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle.CandleSymbol] made with [Candle],
//     such as osub.Candle(candle.NewCandleSymbol("AAPL{=5m}")).
//   - [WildcardSymbol]: all the symbols of the subscribed event types; it is accepted by some feeds only.
//   - [TimeSeriesSubscriptionSymbol]: a symbol and the time from which the history of a time series event is wanted.
//   - [IndexedEventSubscriptionSymbol]: a symbol and the source of an indexed event, such as the orders of one exchange.
//
// # Why a sealed interface
//
// Symbol is a sealed interface: only the types of this package implement it, so passing a value of another type, such
// as a plain string or an int, is a compile-time error rather than an error at run time. The other dxFeed APIs check
// symbols at compile time too: the C++ API through SymbolWrapper, which a string converts to implicitly, and the Swift
// API through a protocol that String conforms to. Go has no implicit conversions and no way to add a method to string,
// so a plain "AAPL" cannot satisfy an interface. Hence String, a string type that implements Symbol, and [Strings],
// which turns several strings into symbols in one call:
//
//	err := subscription.AddSymbols(osub.Strings("AAPL", "IBM")...)
//
// A generic method such as AddSymbols[S ~string] is impossible in Go, since methods have no type parameters, and a
// union type constraint cannot be the element type of a slice, while tools need a mixed []osub.Symbol. A sealed
// interface keeps the compile-time check and one method set for all symbol kinds, close to the Java and C++ APIs.
package osub
