package dxtime_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/dxtime"
)

var moment = time.Date(2024, 5, 1, 10, 30, 0, 123e6, time.UTC)

func newFormat(t *testing.T, zone string) *dxtime.Format {
	t.Helper()
	loc, err := time.LoadLocation(zone)
	if err != nil {
		t.Skip(err)
	}
	f, err := dxtime.NewFormat(loc)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestFormatAndParse(t *testing.T) {
	for _, tc := range []struct {
		zone                         string
		plain, millis, withZone, iso string
	}{
		{"UTC", "20240501-103000", "20240501-103000.123", "20240501-103000+0000", "2024-05-01T10:30:00.123Z"},
		{"America/New_York", "20240501-063000", "20240501-063000.123", "20240501-063000-0400",
			"2024-05-01T06:30:00.123-04:00"},
	} {
		t.Run(tc.zone, func(t *testing.T) {
			f := newFormat(t, tc.zone)
			for _, c := range []struct {
				f    *dxtime.Format
				want string
			}{{f, tc.plain}, {f.WithMillis(), tc.millis}, {f.WithTimeZone(), tc.withZone}, {f.AsFullISO(), tc.iso}} {
				got, err := c.f.Format(moment)
				if err != nil || got != c.want {
					t.Errorf("Format = %q, %v; want %q", got, err, c.want)
				}
				// Each string is parsed back in the zone of the format.
				parsed, err := f.Parse(got)
				want := moment
				if c.want == tc.plain || c.want == tc.withZone {
					want = moment.Truncate(time.Second)
				}
				if err != nil || !parsed.Equal(want) {
					t.Errorf("Parse(%q) = %v, %v; want %v", got, parsed, err, want)
				}
			}
			// A string with a zone does not depend on the zone of the format.
			if parsed, err := f.Parse("2024-05-01T10:30:00Z"); err != nil || !parsed.Equal(moment.Truncate(time.Second)) {
				t.Errorf("Parse of UTC = %v, %v", parsed, err)
			}
		})
	}
}

func TestParseSpecial(t *testing.T) {
	f := dxtime.GMTFormat()
	if parsed, err := f.Parse("0"); err != nil || parsed.UnixMilli() != 0 {
		t.Errorf("Parse(0) = %v, %v", parsed, err)
	}
	if parsed, err := f.Parse("1714559400123"); err != nil || !parsed.Equal(moment) {
		t.Errorf("Parse of millis = %v, %v", parsed, err)
	}
	before := time.Now().Add(-time.Hour)
	if parsed, err := f.Parse("-1h"); err != nil || parsed.Sub(before).Abs() > time.Minute {
		t.Errorf("Parse(-1h) = %v, %v; want about %v", parsed, err, before)
	}
	if _, err := f.Parse("bad"); err == nil {
		t.Error("Parse(bad) gives no error")
	}
	if s, err := f.Format(time.Time{}); err != nil || s != "0" {
		t.Errorf("Format of the zero time = %q, %v", s, err)
	}
}

func TestDefaultFormatIsLocal(t *testing.T) {
	local, err := dxtime.NewFormat(time.Local)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []*dxtime.Format{dxtime.DefaultFormat(), local} {
		got, err := f.WithTimeZone().Format(moment)
		if want := moment.Local().Format("20060102-150405-0700"); err != nil || got != want {
			t.Errorf("Format = %q, %v; want %q", got, err, want)
		}
	}
}

func TestUnknownZone(t *testing.T) {
	if _, err := dxtime.NewFormat(time.FixedZone("Nowhere", 3600)); !errors.Is(err, dxtime.ErrUnknownTimeZone) {
		t.Errorf("error = %v, want ErrUnknownTimeZone", err)
	}
	if _, err := dxtime.NewFormat(nil); err == nil {
		t.Error("NewFormat(nil) gives no error")
	}
}

func TestPeriods(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want time.Duration
	}{
		{"PT1H30M", 90 * time.Minute}, {"1d", 24 * time.Hour}, {"0.5", 500 * time.Millisecond}, {"0", 0},
		{"inf", dxtime.Unlimited},
	} {
		if got, err := dxtime.ParsePeriod(tc.in); err != nil || got != tc.want {
			t.Errorf("ParsePeriod(%q) = %v, %v; want %v", tc.in, got, err, tc.want)
		}
	}
	for _, in := range []string{"-1h", "x"} {
		if _, err := dxtime.ParsePeriod(in); err == nil {
			t.Errorf("ParsePeriod(%q) gives no error", in)
		}
	}
	for _, tc := range []struct {
		in   time.Duration
		want string
	}{
		{0, "PT0S"}, {90 * time.Minute, "PT1H30M0S"}, {25 * time.Hour, "P1DT1H0S"}, {1500 * time.Millisecond, "PT1.5S"},
		{dxtime.Unlimited, "inf"},
	} {
		got, err := dxtime.FormatPeriod(tc.in)
		if err != nil || got != tc.want {
			t.Errorf("FormatPeriod(%v) = %q, %v; want %q", tc.in, got, err, tc.want)
		}
		if back, err := dxtime.ParsePeriod(got); err != nil || back != tc.in {
			t.Errorf("ParsePeriod(%q) = %v, %v; want %v", got, back, err, tc.in)
		}
	}
	if _, err := dxtime.FormatPeriod(-time.Second); err == nil {
		t.Error("FormatPeriod of a negative period gives no error")
	}
}
