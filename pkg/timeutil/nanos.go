package timeutil

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

const (
	NanosInMillis = 1_000_000
)

func GetNanosFromMillisAndNanoPart(timeMillis int64, timeNanoPart int32) int64 {
	return timeMillis*NanosInMillis + int64(timeNanoPart)
}

func GetMillisFromNanos(nanos int64) int64 {
	return mathutil.FloorDivInt64(nanos, NanosInMillis)
}

func GetNanoPartFromNanos(nanos int64) int64 {
	return mathutil.FloorModInt64(nanos, NanosInMillis)
}
