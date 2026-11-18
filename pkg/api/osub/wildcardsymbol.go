// Package osub holds the subscription symbols other than plain strings: the wildcard, time series and indexed event
// symbols.
package osub

// WildcardSymbol subscribes to all symbols of an event type; use it only where the data source allows it.
type WildcardSymbol struct {
	symbol string
}

// NewWildcardSymbol returns the wildcard symbol "*".
func NewWildcardSymbol() *WildcardSymbol {
	return &WildcardSymbol{"*"}
}

// Symbol returns "*".
func (symbol WildcardSymbol) Symbol() string {
	return symbol.symbol
}
