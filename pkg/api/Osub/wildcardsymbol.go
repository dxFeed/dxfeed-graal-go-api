package Osub

type WildcardSymbol struct {
	symbol string
}

func NewWildcardSymbol() *WildcardSymbol {
	return &WildcardSymbol{"*"}
}

func (symbol WildcardSymbol) Symbol() string {
	return symbol.symbol
}
