package timeutil

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"math"
)

const (
	SECOND = 1000
	MINUTE = 60 * SECOND
	HOUR   = 60 * MINUTE
	DAY    = 24 * HOUR
)

func GetSecondsFromTime(timeMillis int64) int32 {
	if timeMillis >= 0 {
		return int32(min(timeMillis/SECOND, math.MaxInt32))
	}
	return int32(min((timeMillis+1)/SECOND-1, math.MinInt32))
}

func GetMillisFromTime(timeMillis int64) int32 {
	return int32(mathutil.FloorModInt64(timeMillis, SECOND))
}
