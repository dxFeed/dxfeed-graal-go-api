// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package candle

import (
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/javautil"
)

// CandlePeriodAttributeKey is the key of the period in a candle symbol string: an empty string, which is a valid key.
const CandlePeriodAttributeKey = ""

// CandlePeriod is the aggregation period of a candle: a value and a CandleType, such as 5 minutes ("5m"). The value
// can be fractional ("0.5m"). Compare periods with Equal, not with ==, because the value is a double.
type CandlePeriod struct {
	value float64
	typ   CandleType
}

var (
	// CandlePeriodTick is one tick.
	CandlePeriodTick = CandlePeriod{1, CandleTypeTick}
	// CandlePeriodDay is one day.
	CandlePeriodDay = CandlePeriod{1, CandleTypeDay}
	// CandlePeriodDefault is the default period: one tick.
	CandlePeriodDefault = CandlePeriodTick
)

// NewCandlePeriod returns the candle period with the given value and type.
func NewCandlePeriod(value float64, typ CandleType) CandlePeriod {
	return CandlePeriod{value: value, typ: typ}
}

// Value returns the value of the period: the number of the units of its type.
func (p CandlePeriod) Value() float64 {
	return p.value
}

// Type returns the type of the period.
func (p CandlePeriod) Type() CandleType {
	return p.typ
}

// PeriodIntervalMillis returns the period in milliseconds as closely as possible (the Java (long) cast of the
// product); it is 0 for the types that are not time intervals.
func (p CandlePeriod) PeriodIntervalMillis() int64 {
	return javautil.LongCast(float64(p.typ.PeriodIntervalMillis()) * p.value)
}

// String returns the string form of the period as Java does: the type alone for the value 1 ("d"), otherwise the
// value (a long if it is whole, otherwise as Double.toString does) followed by the type ("5m", "0.5m").
// ParseCandlePeriod converts it back.
func (p CandlePeriod) String() string {
	if p.value == 1 {
		return p.typ.String()
	}
	return javautil.LongOrDoubleToString(p.value) + p.typ.String()
}

// Equal compares the periods as Java does: the values as Double.compare does (NaN equals NaN, 0.0 and -0.0 differ)
// and the types.
func (p CandlePeriod) Equal(other CandlePeriod) bool {
	return doubleEquals(p.value, other.value) && p.typ == other.typ
}

// ChangeAttributeForSymbol returns the candle symbol string with this period set; the default one removes it.
func (p CandlePeriod) ChangeAttributeForSymbol(symbol string) string {
	if p.Equal(CandlePeriodDefault) {
		return eventsymbols.RemoveAttributeStringByKey(symbol, CandlePeriodAttributeKey)
	}
	return eventsymbols.ChangeAttributeStringByKey(symbol, CandlePeriodAttributeKey, p.String())
}

// ParseCandlePeriod parses the string form of a period: an optional value followed by a type, such as "5m", "m",
// "0.5Day", "1e3s" or "Ticks"; the value is parsed as Java Double.parseDouble does and is 1 if it is omitted. It is
// an error for a bad value or type, see ParseCandleType.
func ParseCandlePeriod(s string) (CandlePeriod, error) {
	if s == CandleTypeDay.String() {
		return CandlePeriodDay, nil
	}
	if s == CandleTypeTick.String() {
		return CandlePeriodTick, nil
	}
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && c != '.' && c != '-' && c != '+' && c != 'e' && c != 'E' {
			break
		}
	}
	valueStr, typeStr := s[:i], s[i:]
	value := 1.0
	if valueStr != "" {
		var ok bool
		if value, ok = javautil.ParseDouble(valueStr); !ok {
			return CandlePeriod{}, fmt.Errorf("invalid candle period value: %s", valueStr)
		}
	}
	typ, err := ParseCandleType(typeStr)
	if err != nil {
		return CandlePeriod{}, err
	}
	return NewCandlePeriod(value, typ), nil
}

// CandlePeriodFromSymbol returns the period of the candle symbol string; the default one if there is none.
// It is an error if the attribute is not valid.
func CandlePeriodFromSymbol(symbol string) (CandlePeriod, error) {
	value, ok := eventsymbols.GetAttributeStringByKey(symbol, CandlePeriodAttributeKey)
	if !ok {
		return CandlePeriodDefault, nil
	}
	return ParseCandlePeriod(value)
}

// NormalizeCandlePeriodAttribute returns the candle symbol string with the period in its canonical form ("1Day" is
// "d"); the default one is removed. A symbol with an invalid attribute is returned as is.
func NormalizeCandlePeriodAttribute(symbol string) string {
	return normalizeAttribute(symbol, CandlePeriodAttributeKey, CandlePeriodDefault.String(),
		func(value string) (string, bool, error) {
			period, err := ParseCandlePeriod(value)
			return period.String(), period.Equal(CandlePeriodDefault), err
		})
}
