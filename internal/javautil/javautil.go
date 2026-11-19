// Package javautil reproduces the text conversions of numbers that the Java API uses in candle symbols:
// Double.parseDouble, Double.toString and the (long) cast.
package javautil

import (
	"math"
	"strconv"
	"strings"
)

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isTypeSuffix(c byte) bool {
	return c == 'f' || c == 'F' || c == 'd' || c == 'D'
}

// skipMantissa skips `digits* ('.' digits*)?` from i and reports whether there was a digit.
func skipMantissa(s string, i *int, hex bool) bool {
	isDigitOf := isDigit
	if hex {
		isDigitOf = isHexDigit
	}
	digits := 0
	for ; *i < len(s) && isDigitOf(s[*i]); *i++ {
		digits++
	}
	if *i < len(s) && s[*i] == '.' {
		*i++
		for ; *i < len(s) && isDigitOf(s[*i]); *i++ {
			digits++
		}
	}
	return digits > 0
}

// skipExponentValue skips `[+-]? digits+` from i and reports whether there was a digit.
func skipExponentValue(s string, i *int) bool {
	if *i < len(s) && (s[*i] == '+' || s[*i] == '-') {
		*i++
	}
	start := *i
	for ; *i < len(s) && isDigit(s[*i]); *i++ {
	}
	return *i > start
}

// ParseDouble parses a string as Java Double.parseDouble does: the characters up to a space at both ends are
// trimmed, NaN and Infinity are accepted, so are hexadecimal literals and the type suffixes f, F, d and D. The second
// result is false if Java would throw NumberFormatException. A value out of the range gives an infinity or zero.
func ParseDouble(s string) (float64, bool) {
	for len(s) > 0 && s[0] <= ' ' {
		s = s[1:]
	}
	for len(s) > 0 && s[len(s)-1] <= ' ' {
		s = s[:len(s)-1]
	}
	i := 0
	negative := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		negative = s[i] == '-'
		i++
	}
	unsigned := s[i:]
	if unsigned == "NaN" {
		return math.NaN(), true
	}
	if unsigned == "Infinity" {
		if negative {
			return math.Inf(-1), true
		}
		return math.Inf(1), true
	}
	hex := len(unsigned) > 1 && unsigned[0] == '0' && (unsigned[1] == 'x' || unsigned[1] == 'X')
	if hex {
		i += 2
	}
	if !skipMantissa(s, &i, hex) {
		return 0, false
	}
	if hex {
		// The binary exponent is mandatory in a hexadecimal literal.
		if i >= len(s) || (s[i] != 'p' && s[i] != 'P') {
			return 0, false
		}
		i++
		if !skipExponentValue(s, &i) {
			return 0, false
		}
	} else if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if !skipExponentValue(s, &i) {
			return 0, false
		}
	}
	numberEnd := i
	if i < len(s) && isTypeSuffix(s[i]) {
		i++
	}
	if i != len(s) {
		return 0, false
	}
	// The text is checked against the Java grammar; a range error still gives the value Java gives (an infinity or 0).
	value, err := strconv.ParseFloat(s[:numberEnd], 64)
	if err != nil {
		if ne, ok := err.(*strconv.NumError); !ok || ne.Err != strconv.ErrRange {
			return 0, false
		}
	}
	return value, true
}

// LongCast converts a double to int64 as the Java (long) cast does: NaN is 0, out of range values saturate.
func LongCast(value float64) int64 {
	const twoPow63 = 9223372036854775808.0
	switch {
	case math.IsNaN(value):
		return 0
	case value >= twoPow63:
		return math.MaxInt64
	case value <= -twoPow63:
		return math.MinInt64
	}
	return int64(value)
}

// splitDigits returns the significant digits of a positive finite number as strconv writes it in the 'e' format and
// the decimal exponent of the first digit: "1.5e+07" is {"15", 7}.
func splitDigits(formatted string) (string, int) {
	ePos := strings.IndexAny(formatted, "eE")
	mantissa := formatted[:ePos]
	exponent, _ := strconv.Atoi(formatted[ePos+1:])
	digits := strings.Replace(mantissa, ".", "", 1)
	digits = strings.TrimRight(digits, "0")
	if digits == "" {
		digits = "0"
	}
	return digits, exponent
}

// DoubleToString formats a double as Java Double.toString does (the shortest digits that identify the value, at
// least two of them, with a decimal exponent outside of 1e-3 <= |x| < 1e7).
func DoubleToString(value float64) string {
	if math.IsNaN(value) {
		return "NaN"
	}
	if math.IsInf(value, 1) {
		return "Infinity"
	}
	if math.IsInf(value, -1) {
		return "-Infinity"
	}
	sign := ""
	if math.Signbit(value) {
		sign = "-"
	}
	magnitude := math.Abs(value)
	if magnitude == 0 {
		return sign + "0.0"
	}
	digits, exponent := splitDigits(strconv.FormatFloat(magnitude, 'e', -1, 64))
	// Java writes at least two significant digits; with only one shortest digit it takes the two-digit decimal
	// closest to the value if that still converts back to it (4.9E-324 rather than 5.0E-324).
	if len(digits) == 1 {
		twoDigits := strconv.FormatFloat(magnitude, 'e', 1, 64)
		if back, err := strconv.ParseFloat(twoDigits, 64); err == nil && back == magnitude {
			digits, exponent = splitDigits(twoDigits)
		}
	}
	var sb strings.Builder
	sb.WriteString(sign)
	if exponent >= -3 && exponent < 7 {
		if exponent >= 0 {
			intDigits := exponent + 1
			if len(digits) <= intDigits {
				sb.WriteString(digits + strings.Repeat("0", intDigits-len(digits)) + ".0")
			} else {
				sb.WriteString(digits[:intDigits] + "." + digits[intDigits:])
			}
		} else {
			sb.WriteString("0." + strings.Repeat("0", -exponent-1) + digits)
		}
	} else {
		fraction := "0"
		if len(digits) > 1 {
			fraction = digits[1:]
		}
		sb.WriteString(digits[:1] + "." + fraction + "E" + strconv.Itoa(exponent))
	}
	return sb.String()
}

// LongOrDoubleToString formats a value as a long if it equals its long cast, otherwise as Double.toString does.
func LongOrDoubleToString(value float64) string {
	longValue := LongCast(value)
	if value == float64(longValue) {
		return strconv.FormatInt(longValue, 10)
	}
	return DoubleToString(value)
}
