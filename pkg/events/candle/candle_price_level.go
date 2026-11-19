// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package candle

import (
	"fmt"
	"math"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/javautil"
)

// CandlePriceLevelAttributeKey is the key of the price level in a candle symbol string: "pl".
const CandlePriceLevelAttributeKey = "pl"

// CandlePriceLevel is the price level attribute of a candle symbol: the price increment of the price types of the
// candle period. The default level is NaN, which is the zero value of the type; use the constructors for any other.
type CandlePriceLevel struct {
	value   float64
	defined bool
}

// CandlePriceLevelDefault is the default price level, NaN.
var CandlePriceLevelDefault = CandlePriceLevel{}

// NewCandlePriceLevel returns the price level with the given value. NaN is the default level. It is an error for a
// negative value (including -0.0) and for an infinity.
func NewCandlePriceLevel(value float64) (CandlePriceLevel, error) {
	if math.IsInf(value, 0) || (!math.IsNaN(value) && math.Signbit(value)) {
		return CandlePriceLevel{}, fmt.Errorf("incorrect candle price level: %s", javautil.DoubleToString(value))
	}
	if math.IsNaN(value) {
		return CandlePriceLevelDefault, nil
	}
	return CandlePriceLevel{value: value, defined: true}, nil
}

// Value returns the value of the level: NaN for the default one.
func (l CandlePriceLevel) Value() float64 {
	if !l.defined {
		return math.NaN()
	}
	return l.value
}

// String returns the level as a long if its value is whole, otherwise as Java Double.toString does; "NaN" for the
// default level. ParseCandlePriceLevel converts it back.
func (l CandlePriceLevel) String() string {
	return javautil.LongOrDoubleToString(l.Value())
}

// Equal compares the levels as Java Double.compare does: NaN equals NaN, 0.0 and -0.0 differ.
func (l CandlePriceLevel) Equal(other CandlePriceLevel) bool {
	return doubleEquals(l.Value(), other.Value())
}

// IsDefault reports whether the level is the default one, NaN.
func (l CandlePriceLevel) IsDefault() bool {
	return !l.defined
}

// ChangeAttributeForSymbol returns the candle symbol string with this level set; the default one removes it.
func (l CandlePriceLevel) ChangeAttributeForSymbol(symbol string) string {
	if l.IsDefault() {
		return eventsymbols.RemoveAttributeStringByKey(symbol, CandlePriceLevelAttributeKey)
	}
	return eventsymbols.ChangeAttributeStringByKey(symbol, CandlePriceLevelAttributeKey, l.String())
}

// ParseCandlePriceLevel parses the string form of a price level as Java Double.parseDouble does. It is an error if
// the string is not a number or the level is not valid, see NewCandlePriceLevel.
func ParseCandlePriceLevel(s string) (CandlePriceLevel, error) {
	value, ok := javautil.ParseDouble(s)
	if !ok {
		return CandlePriceLevel{}, fmt.Errorf("incorrect candle price level: %s", s)
	}
	return NewCandlePriceLevel(value)
}

// CandlePriceLevelFromSymbol returns the price level of the candle symbol string; the default one if there is none.
// It is an error if the attribute is not valid.
func CandlePriceLevelFromSymbol(symbol string) (CandlePriceLevel, error) {
	value, ok := eventsymbols.GetAttributeStringByKey(symbol, CandlePriceLevelAttributeKey)
	if !ok {
		return CandlePriceLevelDefault, nil
	}
	return ParseCandlePriceLevel(value)
}

// NormalizeCandlePriceLevelAttribute returns the candle symbol string with the price level in its canonical form; the
// default one is removed. A symbol with an invalid attribute is returned as is.
func NormalizeCandlePriceLevelAttribute(symbol string) string {
	return normalizeAttribute(symbol, CandlePriceLevelAttributeKey, CandlePriceLevelDefault.String(),
		func(value string) (string, bool, error) {
			level, err := ParseCandlePriceLevel(value)
			return level.String(), level.IsDefault(), err
		})
}
