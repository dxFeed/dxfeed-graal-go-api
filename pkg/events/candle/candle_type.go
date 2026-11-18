package candle

import "fmt"

// CandleType is the type of the candle aggregation period. Together with a value it makes a CandlePeriod.
// The zero value is CandleTypeTick.
type CandleType uint8

const (
	// CandleTypeTick is a certain number of ticks.
	CandleTypeTick CandleType = iota
	// CandleTypeSecond is a certain number of seconds.
	CandleTypeSecond
	// CandleTypeMinute is a certain number of minutes.
	CandleTypeMinute
	// CandleTypeHour is a certain number of hours.
	CandleTypeHour
	// CandleTypeDay is a certain number of days.
	CandleTypeDay
	// CandleTypeWeek is a certain number of weeks.
	CandleTypeWeek
	// CandleTypeMonth is a certain number of months (30 days).
	CandleTypeMonth
	// CandleTypeOptexp is a certain number of option expirations.
	CandleTypeOptexp
	// CandleTypeYear is a certain number of years (365 days).
	CandleTypeYear
	// CandleTypeVolume is a certain trade volume.
	CandleTypeVolume
	// CandleTypePrice is a certain price change, in points.
	CandleTypePrice
	// CandleTypePriceMomentum is a certain price change, with the momentum rules of the Java API.
	CandleTypePriceMomentum
	// CandleTypePriceRenko is a certain price change, with the Renko rules of the Java API.
	CandleTypePriceRenko
)

type candleTypeInfo struct {
	name           string
	str            string
	intervalMillis int64
}

const (
	millisPerSecond = 1000
	millisPerMinute = 60 * millisPerSecond
	millisPerHour   = 60 * millisPerMinute
	millisPerDay    = 24 * millisPerHour
)

var candleTypes = [...]candleTypeInfo{
	CandleTypeTick:          {"TICK", "t", 0},
	CandleTypeSecond:        {"SECOND", "s", millisPerSecond},
	CandleTypeMinute:        {"MINUTE", "m", millisPerMinute},
	CandleTypeHour:          {"HOUR", "h", millisPerHour},
	CandleTypeDay:           {"DAY", "d", millisPerDay},
	CandleTypeWeek:          {"WEEK", "w", 7 * millisPerDay},
	CandleTypeMonth:         {"MONTH", "mo", 30 * millisPerDay},
	CandleTypeOptexp:        {"OPTEXP", "o", 30 * millisPerDay},
	CandleTypeYear:          {"YEAR", "y", 365 * millisPerDay},
	CandleTypeVolume:        {"VOLUME", "v", 0},
	CandleTypePrice:         {"PRICE", "p", 0},
	CandleTypePriceMomentum: {"PRICE_MOMENTUM", "pm", 0},
	CandleTypePriceRenko:    {"PRICE_RENKO", "pr", 0},
}

// CandleTypes returns all candle types.
func CandleTypes() []CandleType {
	result := make([]CandleType, len(candleTypes))
	for i := range candleTypes {
		result[i] = CandleType(i)
	}
	return result
}

// PeriodIntervalMillis returns the length of one unit of the type in milliseconds as closely as possible; it is 0
// for the types that are not time intervals (ticks, volume and the price types).
func (t CandleType) PeriodIntervalMillis() int64 {
	return candleTypes[t].intervalMillis
}

// Name returns the name of the type: "TICK", "MINUTE", "PRICE_RENKO".
func (t CandleType) Name() string {
	return candleTypes[t].name
}

// String returns the short form of the type: "t", "m", "pr". ParseCandleType converts it back.
func (t CandleType) String() string {
	return candleTypes[t].str
}

// ParseCandleType parses the string form of a candle type: the short form ("m"), a name in any case or its prefix
// ("Minute", "min"), or a plural of a name ("Minutes"). It is an error for an empty or an unknown string.
func ParseCandleType(s string) (CandleType, error) {
	if s == "" {
		return 0, fmt.Errorf("missing candle type")
	}
	for i, info := range candleTypes {
		if info.str == s {
			return CandleType(i), nil
		}
	}
	for i, info := range candleTypes {
		// Tick|TICK|tick, Minute|MINUTE|minute, Second|SECOND|second, etc.
		if hasPrefixFoldASCII(info.name, s) {
			return CandleType(i), nil
		}
		// Ticks, Minutes, Seconds, etc.
		if s[len(s)-1] == 's' && equalFoldASCII(info.name, s[:len(s)-1]) {
			return CandleType(i), nil
		}
	}
	return 0, fmt.Errorf("unknown candle type: %s", s)
}
