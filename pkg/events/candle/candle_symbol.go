package candle

// CandleSymbol is the symbol of a candle event: a base symbol with optional attributes such as period, price type and session.
type CandleSymbol struct {
	symbol *string
}

// NewCandleSymbol creates a CandleSymbol from its string representation.
func NewCandleSymbol(symbol string) *CandleSymbol {
	return &CandleSymbol{symbol: &symbol}
}

// Symbol returns the string form of the candle symbol.
func (c *CandleSymbol) Symbol() *string {
	return c.symbol
}

// SetSymbol sets the string form of the candle symbol.
func (c *CandleSymbol) SetSymbol(symbol *string) {
	c.symbol = symbol
}

// String returns the candle symbol as a string.
func (c *CandleSymbol) String() string {
	return *c.symbol
}
