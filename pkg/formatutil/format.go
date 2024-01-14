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
