package formatutil

import (
	"fmt"
	"strconv"
	"time"
)

func FormatChar(c rune) string {
	if c >= 32 && c <= 126 {
		return string(c)
	}
	if c == 0 {
		return "\\0"
	}
	return fmt.Sprintf("\\u%04x", c)
}

func FormatString(s *string) string {
	if s == nil {
		return "nil"
	}
	return *s
}

func FormatNullableString(format string, s *string, defaultValue string) string {
	if s == nil {
		return defaultValue
	}
	return fmt.Sprintf(format, s)
}

func FormatTime(timeMillis int64) string {
	const defaultTimeFormat = "20060102-150405.000-07:00"
	if timeMillis == 0 {
		return "0"
	}
	return time.UnixMilli(timeMillis).Format(defaultTimeFormat)
}

func FormatFloat64(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func FormatInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}

func FormatBool(v bool) string {
	return strconv.FormatBool(v)
}

func HexFormat(value int64) string {
	return `0x` + strconv.FormatInt(value, 16)
}
