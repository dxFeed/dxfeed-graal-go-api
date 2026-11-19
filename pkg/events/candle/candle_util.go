package candle

import (
	"math"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"
)

// CandleSymbolAttribute is an attribute of a candle symbol: it knows how to set itself in a symbol string.
type CandleSymbolAttribute interface {
	// ChangeAttributeForSymbol returns the candle symbol string with this attribute set. The default value of an
	// attribute removes it from the symbol.
	ChangeAttributeForSymbol(symbol string) string
}

func lowerASCII(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

// equalFoldASCII compares the strings ignoring the case of the Latin letters.
func equalFoldASCII(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if lowerASCII(a[i]) != lowerASCII(b[i]) {
			return false
		}
	}
	return true
}

// hasPrefixFoldASCII reports whether s starts with prefix, ignoring the case of the Latin letters.
func hasPrefixFoldASCII(s, prefix string) bool {
	return len(s) >= len(prefix) && equalFoldASCII(s[:len(prefix)], prefix)
}

// doubleEquals compares the numbers as Java Double.compare does: NaN equals NaN, 0.0 and -0.0 differ.
func doubleEquals(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.IsNaN(a) && math.IsNaN(b)
	}
	if a != b {
		return false
	}
	return a != 0 || math.Signbit(a) == math.Signbit(b)
}

// removeDefaultAttribute removes the attribute from the symbol, unless that changes the base symbol or the exchange
// code (the symbol "A{=t}{=t}" has the base symbol "A{=t}"): then the attribute is set to its default value instead.
func removeDefaultAttribute(symbol, key, defaultValue string) string {
	removed := eventsymbols.RemoveAttributeStringByKey(symbol, key)
	if eventsymbols.GetBaseSymbol(removed) == eventsymbols.GetBaseSymbol(symbol) &&
		eventsymbols.GetExchangeCode(removed) == eventsymbols.GetExchangeCode(symbol) {
		return removed
	}
	return eventsymbols.ChangeAttributeStringByKey(symbol, key, defaultValue)
}

// normalizeAttribute brings the attribute of the symbol to its canonical form. The canonical function returns the
// canonical string of an attribute value and whether the value is the default one; an error leaves the symbol as is.
func normalizeAttribute(symbol, key, defaultValue string, canonical func(value string) (string, bool, error)) string {
	value, ok := eventsymbols.GetAttributeStringByKey(symbol, key)
	if !ok {
		return symbol
	}
	normal, isDefault, err := canonical(value)
	if err != nil {
		return symbol
	}
	if isDefault {
		return removeDefaultAttribute(symbol, key, defaultValue)
	}
	if value != normal {
		return eventsymbols.ChangeAttributeStringByKey(symbol, key, normal)
	}
	return symbol
}
