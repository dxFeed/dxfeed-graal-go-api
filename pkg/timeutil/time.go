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

func GetSecondsFromTime(timeMillis int64) int64 {
	if timeMillis >= 0 {
		return mathutil.MinInt64(timeMillis/SECOND, math.MaxInt32)
	}
	return mathutil.MinInt64((timeMillis+1)/SECOND-1, math.MinInt32)
}

func GetMillisFromTime(timeMillis int64) int32 {
	return int32(mathutil.FloorModInt64(timeMillis, SECOND))
}

func GetYearMonthDayByDayId(dayId int32) int32 {
	j := dayId + 2472632 // this shifts the epoch back htto astronomical year -4800
	g := mathutil.Div(j, 146097)
	dg := j - g*146097
	c := (dg/36524 + 1) * 3 / 4
	dc := dg - c*36524
	b := dc / 1461
	db := dc - b*1461
	a := (db/365 + 1) * 3 / 4
	da := db - a*365
	y := g*400 + c*100 + b*4 + a // this is the integer number of full years elapsed since March 1, 4801 BC at 00:00 UTC
	m := (da*5+308)/153 - 2      // this is the integer number of full months elapsed since the last March 1 at 00:00 UTC
	d := da - (m+4)*153/5 + 122  // this is the number of days elapsed since day 1 of the month at 00:00 UTC
	yyyy := y - 4800 + (m+2)/12
	mm := (m+2)%12 + 1
	dd := d + 1
	yyyymmdd := mathutil.Abs(yyyy)*10000 + mm*100 + dd
	if yyyy >= 0 {
		return yyyymmdd
	} else {
		return -yyyymmdd
	}
}
