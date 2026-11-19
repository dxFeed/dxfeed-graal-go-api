package candle

import "github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"

// CandleDataTypeAttributeKey is the key of the data type in a candle symbol string: "dt".
const CandleDataTypeAttributeKey = "dt"

// CandleDataType is the data type attribute of a candle symbol: which data the candles are aggregated from. Besides
// the constants, it holds any other value. The zero value is CandleDataTypeUndefined, the default.
type CandleDataType struct {
	value string
}

var (
	// CandleDataTypeSource aggregates by the rules defined per data source (for example, an exchange).
	CandleDataTypeSource = CandleDataType{"source"}
	// CandleDataTypeTHO selects the same data as CandleSessionRegular.
	CandleDataTypeTHO = CandleDataType{"tho"}
	// CandleDataTypeAll uses all valid TimeAndSales or Quotes to build the candles.
	CandleDataTypeAll = CandleDataType{"all"}
	// CandleDataTypeRTH uses all valid non-ETH TimeAndSales or Quotes from a regular trading hours.
	CandleDataTypeRTH = CandleDataType{"rth"}
	// CandleDataTypeETH uses all valid ETH TimeAndSales or Quotes from an extended trading hours.
	CandleDataTypeETH = CandleDataType{"eth"}
	// CandleDataTypeUndefined leaves the data type to be selected by a chart server. It is the default.
	CandleDataTypeUndefined = CandleDataType{""}
	// CandleDataTypeDefault is the default data type: CandleDataTypeUndefined.
	CandleDataTypeDefault = CandleDataTypeUndefined
)

// NewCandleDataType returns the candle data type with the given value (case-sensitive). An empty value is
// CandleDataTypeDefault.
func NewCandleDataType(value string) CandleDataType {
	return CandleDataType{value: value}
}

// Value returns the string form of the data type, a lower case string such as "tho"; empty for the undefined one.
func (d CandleDataType) Value() string {
	return d.value
}

// String returns the string form of the data type, the same as Value.
func (d CandleDataType) String() string {
	return d.value
}

// ChangeAttributeForSymbol returns the candle symbol string with this data type set; the default one removes it.
func (d CandleDataType) ChangeAttributeForSymbol(symbol string) string {
	if d == CandleDataTypeDefault {
		return eventsymbols.RemoveAttributeStringByKey(symbol, CandleDataTypeAttributeKey)
	}
	return eventsymbols.ChangeAttributeStringByKey(symbol, CandleDataTypeAttributeKey, d.value)
}

// ParseCandleDataType returns the data type with the given string form, the same as NewCandleDataType: any string is
// a valid data type.
func ParseCandleDataType(s string) CandleDataType {
	return NewCandleDataType(s)
}

// CandleDataTypeFromSymbol returns the data type of the candle symbol string; the default one if there is none.
func CandleDataTypeFromSymbol(symbol string) CandleDataType {
	value, ok := eventsymbols.GetAttributeStringByKey(symbol, CandleDataTypeAttributeKey)
	if !ok {
		return CandleDataTypeDefault
	}
	return ParseCandleDataType(value)
}

// NormalizeCandleDataTypeAttribute returns the candle symbol string with an empty data type removed; any other value
// stays as is.
func NormalizeCandleDataTypeAttribute(symbol string) string {
	return normalizeAttribute(symbol, CandleDataTypeAttributeKey, CandleDataTypeDefault.String(),
		func(value string) (string, bool, error) {
			return value, value == "", nil
		})
}
