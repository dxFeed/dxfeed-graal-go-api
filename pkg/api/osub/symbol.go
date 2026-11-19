package osub

// Symbol is a subscription symbol. It is implemented only by the types of this package: String, *CandleSymbol,
// *WildcardSymbol, *TimeSeriesSubscriptionSymbol and *IndexedEventSubscriptionSymbol.
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
