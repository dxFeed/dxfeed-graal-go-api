package candle

import "github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"

// CandleExchange is the exchange attribute of a candle symbol: the exchange code after "&" in the symbol, or the
// composite exchange. The zero value is the composite exchange, the default.
type CandleExchange struct {
	code byte
}

var (
	// CandleExchangeComposite is the composite exchange: the data of all exchanges. It is the default.
	CandleExchangeComposite = CandleExchange{}
	// CandleExchangeDefault is the default exchange: the composite one.
	CandleExchangeDefault = CandleExchangeComposite
)

// NewCandleExchange returns the exchange with the given code; the code 0 is the composite exchange.
func NewCandleExchange(exchangeCode byte) CandleExchange {
	return CandleExchange{code: exchangeCode}
}

// ExchangeCode returns the exchange code: a character, or 0 for the composite exchange.
func (e CandleExchange) ExchangeCode() byte {
	return e.code
}

// String returns the exchange code as a string, or "COMPOSITE".
func (e CandleExchange) String() string {
	if e.code == 0 {
		return "COMPOSITE"
	}
	return string([]byte{e.code})
}

// ChangeAttributeForSymbol returns the candle symbol string with this exchange set; the composite one removes it.
func (e CandleExchange) ChangeAttributeForSymbol(symbol string) string {
	return eventsymbols.ChangeExchangeCode(symbol, e.code)
}

// CandleExchangeFromSymbol returns the exchange of the candle symbol string; the composite one if there is none.
func CandleExchangeFromSymbol(symbol string) CandleExchange {
	return NewCandleExchange(eventsymbols.GetExchangeCode(symbol))
}
