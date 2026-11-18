package candle

import (
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"
)

// CandleAlignmentAttributeKey is the key of the alignment in a candle symbol string: "a".
const CandleAlignmentAttributeKey = "a"

// CandleAlignment is the alignment attribute of a candle symbol: how the candles are aligned in time.
// The zero value is CandleAlignmentMidnight, the default.
type CandleAlignment uint8

const (
	// CandleAlignmentMidnight aligns the candles on midnight.
	CandleAlignmentMidnight CandleAlignment = iota
	// CandleAlignmentSession aligns the candles on the trading sessions.
	CandleAlignmentSession

	// CandleAlignmentDefault is the default alignment: CandleAlignmentMidnight.
	CandleAlignmentDefault = CandleAlignmentMidnight
)

var candleAlignments = [...]string{
	CandleAlignmentMidnight: "m",
	CandleAlignmentSession:  "s",
}

// String returns the string form of the alignment: "m" or "s".
func (a CandleAlignment) String() string {
	return candleAlignments[a]
}

// ChangeAttributeForSymbol returns the candle symbol string with this alignment set; the default one removes it.
func (a CandleAlignment) ChangeAttributeForSymbol(symbol string) string {
	if a == CandleAlignmentDefault {
		return eventsymbols.RemoveAttributeStringByKey(symbol, CandleAlignmentAttributeKey)
	}
	return eventsymbols.ChangeAttributeStringByKey(symbol, CandleAlignmentAttributeKey, a.String())
}

// ParseCandleAlignment parses the string form of an alignment ("m" or "s", in any case). It is an error for any other
// string.
func ParseCandleAlignment(s string) (CandleAlignment, error) {
	for i, str := range candleAlignments {
		if str == s {
			return CandleAlignment(i), nil
		}
	}
	for i, str := range candleAlignments {
		if equalFoldASCII(str, s) {
			return CandleAlignment(i), nil
		}
	}
	return 0, fmt.Errorf("unknown candle alignment: %s", s)
}

// CandleAlignmentFromSymbol returns the alignment of the candle symbol string; the default one if there is none.
// It is an error if the attribute is not valid.
func CandleAlignmentFromSymbol(symbol string) (CandleAlignment, error) {
	value, ok := eventsymbols.GetAttributeStringByKey(symbol, CandleAlignmentAttributeKey)
	if !ok {
		return CandleAlignmentDefault, nil
	}
	return ParseCandleAlignment(value)
}

// NormalizeCandleAlignmentAttribute returns the candle symbol string with the alignment in its canonical form; the
// default one is removed. A symbol with an invalid attribute is returned as is.
func NormalizeCandleAlignmentAttribute(symbol string) string {
	return normalizeAttribute(symbol, CandleAlignmentAttributeKey, CandleAlignmentDefault.String(),
		func(value string) (string, bool, error) {
			alignment, err := ParseCandleAlignment(value)
			return alignment.String(), alignment == CandleAlignmentDefault, err
		})
}
