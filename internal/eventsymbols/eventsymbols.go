// Package eventsymbols works with the string form of market event symbols: a base symbol, an optional exchange code
// ("IBM&N") and optional attributes in curly braces ("IBM{price=bid,tho=true}"). It is the Go port of
// MarketEventSymbols of the Java and C++ APIs, so that every event with such symbols can reuse it.
package eventsymbols

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"regexp/syntax"
	"sort"
	"strconv"
	"strings"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/javautil"
)

const (
	// ExchangesPattern is the character class that all exchange codes of regional events match.
	ExchangesPattern = "[0-9A-Za-z]"
	// SupportedExchanges lists all exchange codes that match ExchangesPattern, in the order of their codes.
	SupportedExchanges = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	exchangeSeparator   = '&'
	attributesOpen      = '{'
	attributesClose     = '}'
	attributesSeparator = ','
	attributeValue      = '='
)

// HasExchangeCode reports whether the symbol has an exchange code ("GE&N").
func HasExchangeCode(symbol string) bool {
	return hasExchangeCode(symbol, lengthWithoutAttributes(symbol))
}

// GetExchangeCode returns the exchange code of the symbol, or 0 if there is none.
func GetExchangeCode(symbol string) byte {
	length := lengthWithoutAttributes(symbol)
	if hasExchangeCode(symbol, length) {
		return symbol[length-1]
	}
	return 0
}

// ChangeExchangeCode returns the symbol with the exchange code replaced; the code 0 removes it. The attributes stay.
func ChangeExchangeCode(symbol string, exchangeCode byte) string {
	length := lengthWithoutAttributes(symbol)
	result := baseSymbol(symbol, length)
	if exchangeCode != 0 {
		result += string(rune(exchangeSeparator)) + string([]byte{exchangeCode})
	}
	if length == len(symbol) {
		return result
	}
	return result + symbol[length:]
}

// GetBaseSymbol returns the symbol without the exchange code and the attributes.
func GetBaseSymbol(symbol string) string {
	return baseSymbol(symbol, lengthWithoutAttributes(symbol))
}

// ChangeBaseSymbol returns the symbol with the base symbol replaced; the exchange code and the attributes stay.
func ChangeBaseSymbol(symbol, base string) string {
	length := lengthWithoutAttributes(symbol)
	if hasExchangeCode(symbol, length) {
		return base + string(rune(exchangeSeparator)) + symbol[length-1:length] + symbol[length:]
	}
	if length == len(symbol) {
		return base
	}
	return base + symbol[length:]
}

// HasAttributes reports whether the symbol ends with a non-empty list of attributes in curly braces.
func HasAttributes(symbol string) bool {
	return lengthWithoutAttributes(symbol) < len(symbol)
}

// GetAttributeStringByKey returns the value of the attribute with the given key; the second result is false if the
// symbol has no such attribute. The key may be empty.
func GetAttributeStringByKey(symbol, key string) (string, bool) {
	return getAttribute(symbol, lengthWithoutAttributes(symbol), key)
}

// ChangeAttributeStringByKey returns the symbol with the attribute set; the attributes stay sorted by key.
func ChangeAttributeStringByKey(symbol, key, value string) string {
	length := lengthWithoutAttributes(symbol)
	if length == len(symbol) {
		return symbol + "{" + key + "=" + value + "}"
	}
	return addAttribute(symbol, length, key, value)
}

// RemoveAttributeStringByKey returns the symbol without the attribute with the given key.
func RemoveAttributeStringByKey(symbol, key string) string {
	return removeAttribute(symbol, lengthWithoutAttributes(symbol), key)
}

func hasExchangeCode(symbol string, length int) bool {
	return length >= 2 && symbol[length-2] == exchangeSeparator
}

func baseSymbol(symbol string, length int) string {
	if hasExchangeCode(symbol, length) {
		return symbol[:length-2]
	}
	return symbol[:length]
}

func hasAttributes(symbol string) bool {
	n := len(symbol)
	if n >= 3 && symbol[n-1] == attributesClose {
		open := strings.LastIndexByte(symbol[:n-1], attributesOpen)
		return open >= 0 && open < n-2
	}
	return false
}

func lengthWithoutAttributes(symbol string) int {
	if hasAttributes(symbol) {
		return strings.LastIndexByte(symbol, attributesOpen)
	}
	return len(symbol)
}

// indexFrom returns the index of c in symbol[from:] as an index of symbol, or -1.
func indexFrom(symbol string, c byte, from int) int {
	if from > len(symbol) {
		return -1
	}
	index := strings.IndexByte(symbol[from:], c)
	if index < 0 {
		return -1
	}
	return from + index
}

func getKey(symbol string, i int) (string, bool) {
	found := indexFrom(symbol, attributeValue, i)
	if found < 0 {
		return "", false
	}
	return symbol[i:found], true
}

func getNextKey(symbol string, i int) int {
	valuePos := indexFrom(symbol, attributeValue, i) + 1
	separatorPos := indexFrom(symbol, attributesSeparator, valuePos)
	if separatorPos < 0 {
		return len(symbol)
	}
	return separatorPos + 1
}

func getValue(symbol string, i, j int) string {
	valueOffset := indexFrom(symbol, attributeValue, i) + 1
	return symbol[valueOffset : j-1]
}

func dropKeyAndValue(symbol string, length, i, j int) string {
	if j == len(symbol) {
		if i == length+1 {
			return symbol[:length]
		}
		return symbol[:i-1] + symbol[j-1:]
	}
	return symbol[:i] + symbol[j:]
}

func getAttribute(symbol string, length int, key string) (string, bool) {
	if length == len(symbol) {
		return "", false
	}
	for i := length + 1; i < len(symbol); {
		currentKey, ok := getKey(symbol, i)
		if !ok {
			break
		}
		j := getNextKey(symbol, i)
		if key == currentKey {
			return getValue(symbol, i, j), true
		}
		i = j
	}
	return "", false
}

func removeAttribute(symbol string, length int, key string) string {
	if length == len(symbol) {
		return symbol
	}
	for i := length + 1; i < len(symbol); {
		currentKey, ok := getKey(symbol, i)
		if !ok {
			break
		}
		j := getNextKey(symbol, i)
		if key == currentKey {
			symbol = dropKeyAndValue(symbol, length, i, j)
		} else {
			i = j
		}
	}
	return symbol
}

func addAttribute(symbol string, length int, key, value string) string {
	i := length + 1
	added := false
	pair := key + "=" + value
	for i < len(symbol) {
		currentKey, ok := getKey(symbol, i)
		if !ok {
			break
		}
		j := getNextKey(symbol, i)
		cmp := strings.Compare(currentKey, key)
		switch {
		case cmp == 0:
			if added {
				// Drop it, since the key is already added.
				symbol = dropKeyAndValue(symbol, length, i, j)
			} else {
				// Replace the value.
				symbol = symbol[:i] + pair + symbol[j-1:]
				added = true
				i += len(key) + len(value) + 2
			}
		case cmp > 0 && !added:
			// Insert the pair here.
			symbol = symbol[:i] + pair + "," + symbol[i:]
			added = true
			i += len(key) + len(value) + 2
		default:
			i = j
		}
	}
	if added {
		return symbol
	}
	return symbol[:i-1] + "," + pair + symbol[i-1:]
}

// BuildSpreadSymbol builds the canonical spread symbol for the legs (the symbols with their ratios): "=AAPL-MSFT".
// The legs with a zero ratio are ignored; no legs give an empty string, one leg with the ratio 1 gives the leg symbol.
// The legs with a NaN ratio, which Java cannot order, go last, by symbol.
func BuildSpreadSymbol(spread map[string]float64) string {
	type leg struct {
		symbol string
		ratio  float64
	}
	legs := make([]leg, 0, len(spread))
	for symbol, ratio := range spread {
		if ratio != 0 {
			legs = append(legs, leg{symbol, ratio})
		}
	}
	if len(legs) == 0 {
		return ""
	}
	if len(legs) == 1 && legs[0].ratio == 1 {
		return legs[0].symbol
	}
	// Java SpreadLeg.compareTo: a larger ratio first, equal ratios by symbol (in reverse order for a negative ratio).
	sort.Slice(legs, func(x, y int) bool {
		a, b := legs[x], legs[y]
		aNaN, bNaN := math.IsNaN(a.ratio), math.IsNaN(b.ratio)
		if aNaN != bNaN {
			return !aNaN
		}
		if !aNaN && a.ratio != b.ratio {
			return a.ratio > b.ratio
		}
		if aNaN || a.ratio >= 0 {
			return a.symbol < b.symbol
		}
		return a.symbol > b.symbol
	})
	var sb strings.Builder
	sb.WriteByte('=')
	for _, l := range legs {
		if l.ratio < 0 {
			sb.WriteByte('-')
		} else if sb.Len() > 1 {
			sb.WriteByte('+') // Plus only when not first.
		}
		if ratio := math.Abs(l.ratio); ratio != 1 {
			// Java prints the ratio as an int when the (int) cast keeps its value, otherwise as Double.toString does.
			if ratio <= math.MaxInt32 && ratio == math.Trunc(ratio) {
				sb.WriteString(strconv.FormatInt(int64(ratio), 10))
			} else {
				sb.WriteString(javautil.DoubleToString(ratio))
			}
			sb.WriteByte('*')
		}
		sb.WriteString(l.symbol)
	}
	return sb.String()
}

// GetExchangesByPattern returns all exchange codes (the characters with the codes 0..127) that match the character
// class, an RE2 regular expression; its outer square brackets can be omitted. An empty class gives an empty string.
// It is an error if the class is not valid or a matching code is not in SupportedExchanges. The Java-specific syntax,
// such as \p{Upper}, is not supported.
func GetExchangesByPattern(characterClass string) (string, error) {
	if characterClass == "" {
		return "", nil
	}
	pattern := characterClass
	if pattern[0] != '[' && pattern[0] != '\\' {
		pattern = "[" + pattern + "]"
	}
	regex, err := regexp.Compile("^(?:" + pattern + ")$")
	if err != nil {
		var syntaxError *syntax.Error
		if errors.As(err, &syntaxError) && syntaxError.Code == syntax.ErrInvalidCharRange {
			return "", fmt.Errorf("invalid character class: %s: the end of a range is less than its start", characterClass)
		}
		return "", fmt.Errorf("invalid character class: %s: %w", characterClass, err)
	}
	var exchanges, unsupported strings.Builder
	for code := 0; code <= 127; code++ {
		exchange := string(rune(code))
		if !regex.MatchString(exchange) {
			continue
		}
		exchanges.WriteString(exchange)
		if !strings.Contains(SupportedExchanges, exchange) {
			if unsupported.Len() > 0 {
				unsupported.WriteString(", ")
			}
			unsupported.WriteString(exchange)
		}
	}
	if unsupported.Len() > 0 {
		return "", fmt.Errorf("unsupported exchange codes: %s", unsupported.String())
	}
	return exchanges.String(), nil
}
