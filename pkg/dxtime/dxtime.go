// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package dxtime formats and parses the times and the time periods the way the dxFeed API does.
//
// The SDK takes times in the strings of its addresses and properties, for example the "fromTime" of a
// subscription or "tape[storagetime=1d]", and writes them in its log. A [Format] gives the same strings as the Java
// TimeFormat:
//
//	f := dxtime.GMTFormat()
//	t, err := f.Parse("2024-05-01T10:30:00Z") // also "20240501-103000", "-1h", millis since the epoch
//	s, err := f.WithMillis().Format(t)         // "20240501-103000.000"
//
// A time period is a [time.Duration]: [ParsePeriod] reads the strings of the Java TimePeriod ("PT1H30M", "1d", "0.5",
// "inf"), [FormatPeriod] writes them.
package dxtime

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// Format is a format of the times. Its methods are safe for concurrent use. The zero value is not a format; use the
// functions of the package to get one.
//
// The functions that return a format do not return errors: an error of the SDK (it cannot start, for example) is kept
// in the format and returned by [Format.Format] and [Format.Parse].
type Format struct {
	f   *native.TimeFormat
	err error

	// The derived formats, created once.
	derivedMu sync.Mutex
	derived   [3]*Format
}

// ErrUnknownTimeZone is the error of [NewFormat] for a time zone that the SDK does not know.
var ErrUnknownTimeZone = native.ErrUnknownTimeZone

var errNoFormat = errors.New("dxtime: the zero Format is not a format")

func newFormat(f *native.TimeFormat, err error) *Format {
	return &Format{f: f, err: err}
}

func (f *Format) check() error {
	switch {
	case f.err != nil:
		return f.err
	case f.f == nil:
		return errNoFormat
	}
	return nil
}

// DefaultFormat returns the format in the default time zone of the SDK: the local time zone or the system property
// user.timezone.
func DefaultFormat() *Format {
	return newFormat(native.NewTimeFormat(native.TimeFormatDefault))
}

// GMTFormat returns the format in GMT.
func GMTFormat() *Format {
	return newFormat(native.NewTimeFormat(native.TimeFormatGMT))
}

// NewFormat returns the format in the time zone of loc: [time.UTC], [time.Local] (the default time zone of the SDK, as
// in [DefaultFormat]) or a location of [time.LoadLocation] such as "America/New_York". A zone that the SDK does not
// know, and a zone of [time.FixedZone] whose name does not give its offset, is [ErrUnknownTimeZone].
func NewFormat(loc *time.Location) (*Format, error) {
	if loc == nil {
		return nil, errors.New("dxtime: nil location")
	}
	id := loc.String()
	switch {
	case loc == time.Local:
		id = ""
	case loc == time.UTC:
	default:
		// The SDK gets the name only: it must be a known zone with the same rules, not a fixed zone named "EST".
		loaded, err := time.LoadLocation(id)
		if err != nil || id == "" || id == "Local" || !sameOffsets(loc, loaded) {
			return nil, fmt.Errorf("%w %q", ErrUnknownTimeZone, id)
		}
	}
	f, err := native.NewTimeFormatInZone(id)
	if err != nil {
		return nil, err
	}
	return newFormat(f, nil), nil
}

// sameOffsets reports whether the zones have the same offsets in winter and in summer of this year.
func sameOffsets(a, b *time.Location) bool {
	year := time.Now().Year()
	for _, month := range []time.Month{time.January, time.July} {
		t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		_, offsetA := t.In(a).Zone()
		_, offsetB := t.In(b).Zone()
		if offsetA != offsetB {
			return false
		}
	}
	return true
}

func (f *Format) derive(kind int) *Format {
	if err := f.check(); err != nil {
		return newFormat(nil, err)
	}
	f.derivedMu.Lock()
	defer f.derivedMu.Unlock()
	if f.derived[kind] == nil {
		f.derived[kind] = newFormat(f.f.Derive(kind))
	}
	return f.derived[kind]
}

// WithTimeZone returns the format that also writes the time zone, for example "20240501-103000+0000".
func (f *Format) WithTimeZone() *Format {
	return f.derive(native.TimeFormatWithTimeZone)
}

// WithMillis returns the format that also writes the milliseconds, for example "20240501-103000.123".
func (f *Format) WithMillis() *Format {
	return f.derive(native.TimeFormatWithMillis)
}

// AsFullISO returns the format of the full ISO 8601 strings, "2024-05-01T10:30:00.123Z".
func (f *Format) AsFullISO() *Format {
	return f.derive(native.TimeFormatFullISO)
}

// Format writes the time with the precision of milliseconds. The Unix epoch is "0"; so is the zero [time.Time],
// which [Format.Parse] reads back as the epoch.
func (f *Format) Format(t time.Time) (string, error) {
	if err := f.check(); err != nil {
		return "", err
	}
	if t.IsZero() {
		t = time.UnixMilli(0)
	}
	if sec := t.Unix(); sec > math.MaxInt64/1000 || sec < math.MinInt64/1000 {
		return "", fmt.Errorf("dxtime: the time %v does not fit into the milliseconds since the epoch", t)
	}
	return f.f.Format(t.UnixMilli())
}

// Parse reads a time: "0" (the Unix epoch), a time before now ("-1h", "-PT5M"), milliseconds since the epoch (9 and
// more digits), or a date and a time such as "20240501-103000", "20240501-103000-0400" or "2024-05-01T10:30:00.123Z".
// A time without a zone offset is in the zone of the format. The result is in the local time zone.
func (f *Format) Parse(s string) (time.Time, error) {
	if err := f.check(); err != nil {
		return time.Time{}, err
	}
	millis, err := f.f.Parse(s)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(millis), nil
}

// Unlimited is the "inf" time period of the SDK, the maximum [time.Duration].
const Unlimited time.Duration = math.MaxInt64

// ParsePeriod reads a time period: an ISO 8601 duration ("PT1H30M", "P1D") where the letters are case insensitive and
// "P", "T" and the last "S" can be left out ("1h30m", "1d", "0.5"), or "inf" for [Unlimited]. The precision is
// milliseconds; a period of more than about 292 years is [Unlimited].
func ParsePeriod(s string) (time.Duration, error) {
	millis, err := native.ParseTimePeriod(s)
	if err != nil {
		return 0, err
	}
	if millis < 0 {
		return 0, fmt.Errorf("dxtime: the time period %q is negative", s)
	}
	if millis > math.MaxInt64/int64(time.Millisecond) {
		return Unlimited, nil
	}
	return time.Duration(millis) * time.Millisecond, nil
}

// FormatPeriod writes the time period as the SDK does, for example "PT1H30M0S", "P1DT1H0S", "PT1.5S" or "inf" for
// [Unlimited]. The period is truncated to milliseconds; a negative one is an error.
func FormatPeriod(d time.Duration) (string, error) {
	if d < 0 {
		return "", errors.New("dxtime: negative time period")
	}
	millis := d.Milliseconds()
	if d == Unlimited {
		millis = math.MaxInt64
	}
	return native.FormatTimePeriod(millis)
}
