package candle

import (
	"fmt"
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"
)

// CandlePriceAttributeKey is the key of the price type in a candle symbol string: "price".
const CandlePriceAttributeKey = "price"

// CandlePrice is the price type attribute of a candle symbol: which price the candles are built from.
// The zero value is CandlePriceLast, the default.
type CandlePrice uint8

const (
	// CandlePriceLast is the last trading price.
	CandlePriceLast CandlePrice = iota
	// CandlePriceBid is the quote bid price.
	CandlePriceBid
	// CandlePriceAsk is the quote ask price.
	CandlePriceAsk
	// CandlePriceMark is the market price defined as an average between quote bid and ask prices.
	CandlePriceMark
	// CandlePriceSettlement is the official settlement price that is defined by an exchange, or the last trading
	// price otherwise.
	CandlePriceSettlement

	// CandlePriceDefault is the default price type: CandlePriceLast.
	CandlePriceDefault = CandlePriceLast
)

var candlePrices = [...]string{
	CandlePriceLast:       "last",
	CandlePriceBid:        "bid",
	CandlePriceAsk:        "ask",
	CandlePriceMark:       "mark",
	CandlePriceSettlement: "s",
}

// String returns the string form of the price type: "last", "bid", "ask", "mark" or "s".
func (p CandlePrice) String() string {
	if int(p) >= len(candlePrices) {
		return "CandlePrice(" + strconv.Itoa(int(p)) + ")"
	}
	return candlePrices[p]
}

// ChangeAttributeForSymbol returns the candle symbol string with this price type set; the default one removes it.
func (p CandlePrice) ChangeAttributeForSymbol(symbol string) string {
	if p == CandlePriceDefault {
		return eventsymbols.RemoveAttributeStringByKey(symbol, CandlePriceAttributeKey)
	}
	return eventsymbols.ChangeAttributeStringByKey(symbol, CandlePriceAttributeKey, p.String())
}

// ParseCandlePrice parses the string form of a price type: the exact string or its prefix in any case ("b", "ASK").
// It is an error for an empty or an unknown string.
func ParseCandlePrice(s string) (CandlePrice, error) {
	if s == "" {
		return 0, fmt.Errorf("missing candle price")
	}
	for i, str := range candlePrices {
		if str == s {
			return CandlePrice(i), nil
		}
	}
	for i, str := range candlePrices {
		if hasPrefixFoldASCII(str, s) {
			return CandlePrice(i), nil
		}
	}
	return 0, fmt.Errorf("unknown candle price: %s", s)
}

// CandlePriceFromSymbol returns the price type of the candle symbol string; the default one if there is none.
// It is an error if the attribute is not valid.
func CandlePriceFromSymbol(symbol string) (CandlePrice, error) {
	value, ok := eventsymbols.GetAttributeStringByKey(symbol, CandlePriceAttributeKey)
	if !ok {
		return CandlePriceDefault, nil
	}
	return ParseCandlePrice(value)
}

// NormalizeCandlePriceAttribute returns the candle symbol string with the price type in its canonical form; the
// default one is removed. A symbol with an invalid attribute is returned as is.
func NormalizeCandlePriceAttribute(symbol string) string {
	return normalizeAttribute(symbol, CandlePriceAttributeKey, CandlePriceDefault.String(),
		func(value string) (string, bool, error) {
			price, err := ParseCandlePrice(value)
			return price.String(), price == CandlePriceDefault, err
		})
}
