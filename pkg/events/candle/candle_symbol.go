package candle

type CandleSymbol struct {
	symbol *string
}

func NewCandleSymbol(symbol string) *CandleSymbol {
	return &CandleSymbol{symbol: &symbol}
}

func (c *CandleSymbol) Symbol() *string {
	return c.symbol
}

func (c *CandleSymbol) SetSymbol(symbol *string) {
	c.symbol = symbol
}
