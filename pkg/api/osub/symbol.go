// Package osub holds the subscription symbols: String for a plain symbol such as "AAPL", the wildcard, and the time
// series and indexed event symbols. They all implement Symbol, the type that DXFeedSubscription.AddSymbols and the
// other subscription methods of package api take.
//
// # Why a sealed interface
//
// Symbol is a sealed interface: only the types of this package implement it, so passing a value of another type, such
// as a plain string or an int, is a compile-time error rather than an error at run time. The other dxFeed APIs check
// symbols at compile time too: the C++ API through SymbolWrapper, which a string converts to implicitly, and the Swift
// API through a protocol that String conforms to. Go has no implicit conversions and no way to add a method to string, so a plain "AAPL" cannot
// satisfy an interface. Hence String, a string type that implements Symbol, and Strings, which turns several strings
// into symbols in one call:
//
//	err := subscription.AddSymbols(osub.Strings("AAPL", "IBM")...)
//
// A generic method such as AddSymbols[S ~string] is impossible in Go, since methods have no type parameters, and a
// union type constraint cannot be the element type of a slice, while tools need a mixed []osub.Symbol. A sealed
// interface keeps the compile-time check and one method set for all symbol kinds, close to the Java and C++ APIs.
package osub

// Symbol is a subscription symbol. It is implemented only by the types of this package: String, *WildcardSymbol,
// *TimeSeriesSubscriptionSymbol and *IndexedEventSubscriptionSymbol.
type Symbol interface {
	isSymbol()
}

// String is a plain symbol, such as "AAPL", or a candle symbol, such as "AAPL{=d}".
type String string

func (String) isSymbol() {}

// Strings returns the given strings as symbols, so that a subscription takes them in one call:
// subscription.AddSymbols(osub.Strings("AAPL", "IBM")...). It accepts any string type.
func Strings[S ~string](symbols ...S) []Symbol {
	result := make([]Symbol, len(symbols))
	for i, symbol := range symbols {
		result[i] = String(symbol)
	}
	return result
}

func (*WildcardSymbol) isSymbol() {}

func (*TimeSeriesSubscriptionSymbol) isSymbol() {}

func (*IndexedEventSubscriptionSymbol) isSymbol() {}

// The symbol kinds that the subscription accepts.
var (
	_ Symbol = String("")
	_ Symbol = (*WildcardSymbol)(nil)
	_ Symbol = (*TimeSeriesSubscriptionSymbol)(nil)
	_ Symbol = (*IndexedEventSubscriptionSymbol)(nil)
)
