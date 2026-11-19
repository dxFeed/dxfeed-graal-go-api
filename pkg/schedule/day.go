package schedule

import (
	"fmt"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// Day is a day of a schedule: a continuous period of time, usually a calendar day, that may start and end at other
// times than midnight, made of sessions that cover it without gaps.
type Day struct {
	schedule Schedule
	f        native.DayFields
}

// Schedule returns the schedule of the day.
func (d Day) Schedule() Schedule { return d.schedule }

// DayID returns the number of days since January 1, 1970.
func (d Day) DayID() int { return int(d.f.DayID) }

// YearMonthDay returns the date of the day as a yyyymmdd number, such as 20241231.
func (d Day) YearMonthDay() int { return int(d.f.YearMonthDay) }

// Year returns the year of the day.
func (d Day) Year() int { return int(d.f.YearMonthDay) / 10000 }

// MonthOfYear returns the month of the day, 1 for January.
func (d Day) MonthOfYear() time.Month { return time.Month(d.f.YearMonthDay / 100 % 100) }

// DayOfMonth returns the day of the month, from 1.
func (d Day) DayOfMonth() int { return int(d.f.YearMonthDay % 100) }

// DayOfWeek returns the day of the week.
func (d Day) DayOfWeek() time.Weekday {
	// Day 0, January 1, 1970, is a Thursday.
	return time.Weekday(((int(d.f.DayID)+int(time.Thursday))%7 + 7) % 7)
}

// IsHoliday reports whether the day is a holiday.
func (d Day) IsHoliday() bool { return d.f.Holiday }

// IsShortDay reports whether the day is a short day.
func (d Day) IsShortDay() bool { return d.f.ShortDay }

// IsTrading reports whether the day has a trading session.
func (d Day) IsTrading() bool { return d.f.Trading }

// StartTime returns the start of the day, inclusive.
func (d Day) StartTime() time.Time { return timeOf(d.f.StartTime) }

// EndTime returns the end of the day, exclusive.
func (d Day) EndTime() time.Time { return timeOf(d.f.EndTime) }

// ResetTime returns the time of the reset of the trading data for the day, such as the daily statistics.
func (d Day) ResetTime() time.Time { return timeOf(d.f.ResetTime) }

// ContainsTime reports whether the day contains a time.
func (d Day) ContainsTime(t time.Time) bool {
	millis := t.UnixMilli()
	return d.f.StartTime <= millis && millis < d.f.EndTime
}

// String returns the date of the day and its kind, such as "Day(2024-12-31, trading)".
func (d Day) String() string {
	kind := "non-trading"
	switch {
	case d.f.Holiday:
		kind = "holiday"
	case d.f.ShortDay:
		kind = "short day"
	case d.f.Trading:
		kind = "trading"
	}
	return fmt.Sprintf("Day(%04d-%02d-%02d, %s)", d.Year(), d.MonthOfYear(), d.DayOfMonth(), kind)
}

// Sessions returns the sessions of the day in their order.
func (d Day) Sessions() ([]Session, error) {
	fields, err := native.DaySessions(d.schedule.f.Definition, d.f.DayID)
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, 0, len(fields))
	for _, f := range fields {
		sessions = append(sessions, d.schedule.session(f))
	}
	return sessions, nil
}

// SessionByTime returns the session of the day that contains a time; it fails for a time outside the day.
func (d Day) SessionByTime(t time.Time) (Session, error) {
	f, err := native.DaySessionByTime(d.schedule.f.Definition, d.f.DayID, t.UnixMilli())
	return d.schedule.session(f), err
}

// FirstSession returns the first session of the day that the filter accepts; ok is false if there is none.
func (d Day) FirstSession(filter SessionFilter) (session Session, ok bool, err error) {
	f, ok, err := native.DayEdgeSession(d.schedule.f.Definition, d.f.DayID, false, filter.spec(d.schedule))
	return d.schedule.session(f), ok, err
}

// LastSession returns the last session of the day that the filter accepts; ok is false if there is none.
func (d Day) LastSession(filter SessionFilter) (session Session, ok bool, err error) {
	f, ok, err := native.DayEdgeSession(d.schedule.f.Definition, d.f.DayID, true, filter.spec(d.schedule))
	return d.schedule.session(f), ok, err
}

// PrevDay returns the nearest day before this one that the filter accepts; ok is false if there is none within
// about a year.
func (d Day) PrevDay(filter DayFilter) (day Day, ok bool, err error) {
	f, ok, err := native.DayNeighbor(d.schedule.f.Definition, d.f.DayID, false, filter.spec(d.schedule))
	return d.schedule.day(f), ok, err
}

// NextDay returns the nearest day after this one that the filter accepts; ok is false if there is none within about
// a year.
func (d Day) NextDay(filter DayFilter) (day Day, ok bool, err error) {
	f, ok, err := native.DayNeighbor(d.schedule.f.Definition, d.f.DayID, true, filter.spec(d.schedule))
	return d.schedule.day(f), ok, err
}
