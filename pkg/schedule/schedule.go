// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package schedule

import (
	"errors"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// Schedule is the trading schedule of an instrument: its days and sessions in its time zone.
type Schedule struct {
	f native.ScheduleFields
}

// FromDefinition returns the schedule of a definition: a name and properties in parentheses, such as
// "NewYorkETH()" or "(tz=GMT;td=12345;0=09001700)"; the properties of a known name, such as NewYorkETH, come from the
// defaults.
func FromDefinition(definition string) (Schedule, error) {
	f, err := native.ScheduleOfDefinition(definition)
	return Schedule{f}, err
}

// FromProfile returns the schedule of the trading hours of an instrument profile.
func FromProfile(profile *events.InstrumentProfile) (Schedule, error) {
	if profile == nil {
		return Schedule{}, errors.New("schedule: nil profile")
	}
	f, err := native.ScheduleOfProfile(profile, "")
	return Schedule{f}, err
}

// FromProfileAndVenue returns the schedule of a trading venue in the trading hours of an instrument profile.
func FromProfileAndVenue(profile *events.InstrumentProfile, venue string) (Schedule, error) {
	if profile == nil {
		return Schedule{}, errors.New("schedule: nil profile")
	}
	if venue == "" {
		return Schedule{}, errors.New("schedule: empty venue")
	}
	f, err := native.ScheduleOfProfile(profile, venue)
	return Schedule{f}, err
}

// TradingVenues returns the trading venues in the trading hours of an instrument profile.
func TradingVenues(profile *events.InstrumentProfile) ([]string, error) {
	if profile == nil {
		return nil, errors.New("schedule: nil profile")
	}
	return native.ScheduleTradingVenues(profile)
}

// DownloadDefaults downloads the default schedule data. The config is a URL to download once, "<period>@<URL>" (such
// as "1d@https://example.com/schedule.zip") to download now and then every period, or "" to stop the periodic
// download.
func DownloadDefaults(config string) error {
	return native.ScheduleDownloadDefaults(config)
}

// SetDefaults sets the default schedule data, plain or compressed.
func SetDefaults(data []byte) error {
	return native.ScheduleSetDefaults(data)
}

// Definition returns the definition of the schedule, such as "NewYorkETH()".
func (s Schedule) Definition() string { return s.f.Definition }

// Name returns the display name of the schedule, such as "US ETH" for NewYorkETH().
func (s Schedule) Name() string { return s.f.Name }

// TimeZone returns the display name of the time zone of the schedule.
func (s Schedule) TimeZone() string { return s.f.TimeZone }

// TimeZoneID returns the identifier of the time zone of the schedule, such as "America/New_York".
func (s Schedule) TimeZoneID() string { return s.f.TimeZoneID }

// String returns the definition of the schedule.
func (s Schedule) String() string { return s.f.Definition }

func (s Schedule) day(f native.DayFields) Day {
	return Day{schedule: s, f: f}
}

func (s Schedule) session(f native.SessionFields) Session {
	return Session{schedule: s, f: f}
}

// DayByTime returns the day that contains a time.
func (s Schedule) DayByTime(t time.Time) (Day, error) {
	f, err := native.ScheduleDayByTime(s.f.Definition, t.UnixMilli())
	return s.day(f), err
}

// DayByID returns the day of an id: the number of days since January 1, 1970.
func (s Schedule) DayByID(dayID int) (Day, error) {
	id, err := int32Of(dayID)
	if err != nil {
		return Day{}, err
	}
	f, err := native.ScheduleDayByID(s.f.Definition, id)
	return s.day(f), err
}

// DayByYearMonthDay returns the day of a yyyymmdd number, such as 20241231.
func (s Schedule) DayByYearMonthDay(yearMonthDay int) (Day, error) {
	ymd, err := int32Of(yearMonthDay)
	if err != nil {
		return Day{}, err
	}
	f, err := native.ScheduleDayByYearMonthDay(s.f.Definition, ymd)
	return s.day(f), err
}

// SessionByTime returns the session that contains a time.
func (s Schedule) SessionByTime(t time.Time) (Session, error) {
	f, err := native.ScheduleSessionByTime(s.f.Definition, t.UnixMilli())
	return s.session(f), err
}

// NearestSessionByTime returns the session nearest to a time that the filter accepts: the session that contains the
// time, or else the nearest one before or after it. It searches about a year each way; ok is false if it finds none.
func (s Schedule) NearestSessionByTime(t time.Time, filter SessionFilter) (session Session, ok bool, err error) {
	f, ok, err := native.ScheduleNearestSession(s.f.Definition, t.UnixMilli(), filter.spec(s))
	return s.session(f), ok, err
}

func int32Of(v int) (int32, error) {
	if int(int32(v)) != v {
		return 0, errors.New("schedule: the value is out of range")
	}
	return int32(v), nil
}

// timeOf converts milliseconds since the epoch to a time.
func timeOf(millis int64) time.Time {
	return time.UnixMilli(millis)
}
