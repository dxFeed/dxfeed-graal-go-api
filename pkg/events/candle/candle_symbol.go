package candle

import "github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"

// CandleSymbol is the symbol of a candle event: a base symbol with optional attributes such as period, price type and
// session, for example "AAPL&Q{=5m,price=bid}". The string form is normalized: the attributes are sorted by key, have
// their canonical form, and the default ones are removed.
type CandleSymbol struct {
	symbol     string
	baseSymbol string
	exchange   CandleExchange
	price      CandlePrice
	session    CandleSession
	period     CandlePeriod
	alignment  CandleAlignment
	priceLevel CandlePriceLevel
	dataType   CandleDataType
}

// NormalizeCandleSymbol returns the candle symbol string with all attributes in their canonical form. The attributes
// that cannot be parsed stay as they are.
func NormalizeCandleSymbol(symbol string) string {
	symbol = NormalizeCandlePriceAttribute(symbol)
	symbol = NormalizeCandleSessionAttribute(symbol)
	symbol = NormalizeCandlePeriodAttribute(symbol)
	symbol = NormalizeCandleAlignmentAttribute(symbol)
	symbol = NormalizeCandlePriceLevelAttribute(symbol)
	return NormalizeCandleDataTypeAttribute(symbol)
}

// initFields reads the attributes of the already normalized symbol. With lenient set, an attribute that cannot be
// parsed gets its default value; otherwise it is an error.
func (c *CandleSymbol) initFields(lenient bool) error {
	var err error
	// check keeps the first error, in the order of the C++ API, and gives the default value to a bad attribute.
	check := func(e error, setDefault func()) {
		if e == nil {
			return
		}
		setDefault()
		if !lenient && err == nil {
			err = e
		}
	}
	var e error
	c.baseSymbol = eventsymbols.GetBaseSymbol(c.symbol)
	c.exchange = CandleExchangeFromSymbol(c.symbol)
	c.price, e = CandlePriceFromSymbol(c.symbol)
	check(e, func() { c.price = CandlePriceDefault })
	c.session = CandleSessionFromSymbol(c.symbol)
	c.period, e = CandlePeriodFromSymbol(c.symbol)
	check(e, func() { c.period = CandlePeriodDefault })
	c.alignment, e = CandleAlignmentFromSymbol(c.symbol)
	check(e, func() { c.alignment = CandleAlignmentDefault })
	c.priceLevel, e = CandlePriceLevelFromSymbol(c.symbol)
	check(e, func() { c.priceLevel = CandlePriceLevelDefault })
	c.dataType = CandleDataTypeFromSymbol(c.symbol)
	return err
}

// ParseCandleSymbol converts the string into a candle symbol and normalizes it. It is an error if an attribute of
// the string is not valid, for example "AAPL{=garbage}".
func ParseCandleSymbol(symbol string) (*CandleSymbol, error) {
	c := &CandleSymbol{symbol: NormalizeCandleSymbol(symbol)}
	if err := c.initFields(false); err != nil {
		return nil, err
	}
	return c, nil
}

// NewCandleSymbol converts the string into a candle symbol and normalizes it. It never fails: an attribute that
// cannot be parsed stays in the string as it is and has its default value in the getters. Use ParseCandleSymbol to
// get an error instead.
func NewCandleSymbol(symbol string) *CandleSymbol {
	c := &CandleSymbol{symbol: NormalizeCandleSymbol(symbol)}
	_ = c.initFields(true)
	return c
}

// NewCandleSymbolWithAttributes converts the string into a candle symbol with the attributes set, in the given order,
// and normalizes it. It is an error if an attribute of the resulting string is not valid.
func NewCandleSymbolWithAttributes(symbol string, attributes ...CandleSymbolAttribute) (*CandleSymbol, error) {
	for _, attribute := range attributes {
		symbol = attribute.ChangeAttributeForSymbol(symbol)
	}
	return ParseCandleSymbol(symbol)
}

// BaseSymbol returns the base market symbol without the exchange code and the attributes.
func (c *CandleSymbol) BaseSymbol() string {
	return c.baseSymbol
}

// Exchange returns the exchange attribute of the symbol.
func (c *CandleSymbol) Exchange() CandleExchange {
	return c.exchange
}

// Price returns the price type attribute of the symbol.
func (c *CandleSymbol) Price() CandlePrice {
	return c.price
}

// Session returns the session attribute of the symbol.
func (c *CandleSymbol) Session() CandleSession {
	return c.session
}

// Period returns the aggregation period attribute of the symbol.
func (c *CandleSymbol) Period() CandlePeriod {
	return c.period
}

// Alignment returns the alignment attribute of the symbol.
func (c *CandleSymbol) Alignment() CandleAlignment {
	return c.alignment
}

// PriceLevel returns the price level attribute of the symbol.
func (c *CandleSymbol) PriceLevel() CandlePriceLevel {
	return c.priceLevel
}

// DataType returns the data type attribute of the symbol.
func (c *CandleSymbol) DataType() CandleDataType {
	return c.dataType
}

// Equal reports whether the symbols have the same normalized string.
func (c *CandleSymbol) Equal(other *CandleSymbol) bool {
	if c == nil || other == nil {
		return c == other
	}
	return c.symbol == other.symbol
}

// Symbol returns the normalized string form of the candle symbol.
//
// Deprecated: use String.
func (c *CandleSymbol) Symbol() *string {
	symbol := c.symbol
	return &symbol
}

// SetSymbol sets the candle symbol from its string form, as NewCandleSymbol does; nil is an empty string.
//
// Deprecated: a candle symbol is a value, create a new one with NewCandleSymbol or ParseCandleSymbol.
func (c *CandleSymbol) SetSymbol(symbol *string) {
	value := ""
	if symbol != nil {
		value = *symbol
	}
	*c = *NewCandleSymbol(value)
}

// String returns the normalized string form of the candle symbol, which ParseCandleSymbol converts back.
func (c *CandleSymbol) String() string {
	return c.symbol
}
