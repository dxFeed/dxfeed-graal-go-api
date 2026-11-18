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
	"math"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// Format is a format of the times. Its methods are safe for concurrent use.
type Format struct {
	f *native.TimeFormat
}

// ErrUnknownTimeZone is the error of [NewFormat] for a time zone that the SDK does not know.
var ErrUnknownTimeZone = native.ErrUnknownTimeZone

func must(f *native.TimeFormat, err error) *Format {
	if err != nil {
		// The constants of the Java API fail only if the isolate cannot be created.
		panic(err)
	}
	return &Format{f}
}

// DefaultFormat returns the format in the default time zone of the SDK: the local time zone or the system property
// user.timezone. It panics if the SDK cannot start.
func DefaultFormat() *Format {
	return must(native.NewTimeFormat(native.TimeFormatDefault))
}

// GMTFormat returns the format in GMT. It panics if the SDK cannot start.
func GMTFormat() *Format {
	return must(native.NewTimeFormat(native.TimeFormatGMT))
}

// NewFormat returns the format in the time zone of loc: an IANA name such as "America/New_York", or [time.UTC].
// [time.Local] is the default time zone of the SDK, as in [DefaultFormat]. A zone that the SDK does not know is
// [ErrUnknownTimeZone].
func NewFormat(loc *time.Location) (*Format, error) {
	if loc == nil {
		return nil, errors.New("dxtime: nil location")
	}
	id := loc.String()
	if loc == time.Local {
		id = ""
	}
	f, err := native.NewTimeFormatInZone(id)
	if err != nil {
		return nil, err
	}
	return &Format{f}, nil
}

func (f *Format) derive(kind int) *Format {
	return must(f.f.Derive(kind))
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

// Format writes the time with the precision of milliseconds; the Unix epoch (and the zero [time.Time]) is "0".
func (f *Format) Format(t time.Time) (string, error) {
	if t.IsZero() {
		t = time.UnixMilli(0)
	}
	return f.f.Format(t.UnixMilli())
}

// Parse reads a time: "0" (the Unix epoch), a time before now ("-1h", "-PT5M"), milliseconds since the epoch (9 and
// more digits), or a date and a time such as "20240501-103000", "20240501-103000-0400" or "2024-05-01T10:30:00.123Z".
// A time without a zone offset is in the zone of the format. The result is in the local time zone.
func (f *Format) Parse(s string) (time.Time, error) {
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
	if millis > math.MaxInt64/int64(time.Millisecond) {
		return Unlimited, nil
	}
	return time.Duration(millis) * time.Millisecond, nil
}

// FormatPeriod writes the time period as the SDK does, for example "PT1H30M0S", "P1DT1H0S", "PT1.5S" or "inf" for
// [Unlimited]. The period is
// truncated to milliseconds; a negative one is an error.
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
