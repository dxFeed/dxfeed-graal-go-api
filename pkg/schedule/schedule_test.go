package schedule_test

import (
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/schedule"
)

const hour = time.Hour

// dayStart returns the start of a day id in UTC.
func dayStart(dayID int) time.Time {
	return time.Unix(0, 0).UTC().Add(time.Duration(dayID) * 24 * hour)
}

func ok[T any](t *testing.T) func(v T, err error) T {
	return func(v T, err error) T {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}

func found[T any](t *testing.T, want bool) func(v T, ok bool, err error) T {
	return func(v T, ok bool, err error) T {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		if ok != want {
			t.Fatalf("found = %v, want %v", ok, want)
		}
		return v
	}
}

func TestSessions(t *testing.T) {
	s := ok[schedule.Schedule](t)(schedule.FromDefinition("(tz=GMT;0=01000200)"))
	if s.TimeZoneID() != "GMT" {
		t.Errorf("TimeZoneID = %q", s.TimeZoneID())
	}
	day := ok[schedule.Day](t)(s.DayByID(42))
	if day.DayID() != 42 || day.YearMonthDay() != 19700212 || day.DayOfWeek() != time.Thursday {
		t.Fatalf("day = %d %d %v", day.DayID(), day.YearMonthDay(), day.DayOfWeek())
	}
	if !day.StartTime().Equal(dayStart(42)) || !day.EndTime().Equal(dayStart(43)) || !day.IsTrading() {
		t.Errorf("day = %v %v-%v", day, day.StartTime(), day.EndTime())
	}
	sessions := ok[[]schedule.Session](t)(day.Sessions())
	want := []struct {
		typ   schedule.SessionType
		start time.Duration
	}{{schedule.NoTrading, 0}, {schedule.Regular, hour}, {schedule.NoTrading, 2 * hour}}
	if len(sessions) != len(want) {
		t.Fatalf("sessions = %v", sessions)
	}
	for i, w := range want {
		got := sessions[i]
		if got.Type() != w.typ || !got.StartTime().Equal(dayStart(42).Add(w.start)) || got.Day() != day {
			t.Errorf("session %d = %v", i, got)
		}
		if got.IsTrading() != w.typ.IsTrading() {
			t.Errorf("session %d trading = %v", i, got.IsTrading())
		}
	}

	at := dayStart(42).Add(90 * time.Minute)
	regular := ok[schedule.Session](t)(s.SessionByTime(at))
	if regular != sessions[1] || !regular.ContainsTime(at) || regular.ContainsTime(regular.EndTime()) {
		t.Errorf("SessionByTime = %v", regular)
	}
	if got := ok[schedule.Session](t)(day.SessionByTime(at)); got != regular {
		t.Errorf("Day.SessionByTime = %v", got)
	}
	if got := ok[schedule.Day](t)(s.DayByTime(at)); got != day {
		t.Errorf("DayByTime = %v", got)
	}
	if got := ok[schedule.Day](t)(s.DayByYearMonthDay(19700212)); got != day {
		t.Errorf("DayByYearMonthDay = %v", got)
	}
	if _, err := day.SessionByTime(dayStart(50)); err == nil {
		t.Error("Day.SessionByTime outside the day did not fail")
	}

	next := found[schedule.Session](t, true)(regular.NextSession(schedule.TradingSessions))
	if next.Day().DayID() != 43 || next.Type() != schedule.Regular {
		t.Errorf("NextSession = %v on %v", next, next.Day())
	}
	prev := found[schedule.Session](t, true)(regular.PrevSession(schedule.AnySession))
	if prev != sessions[0] {
		t.Errorf("PrevSession = %v", prev)
	}
	first := found[schedule.Session](t, true)(day.FirstSession(schedule.RegularSessions))
	last := found[schedule.Session](t, true)(day.LastSession(schedule.NoTradingSessions))
	if first != regular || last != sessions[2] {
		t.Errorf("FirstSession = %v, LastSession = %v", first, last)
	}
	found[schedule.Session](t, false)(day.FirstSession(schedule.PreMarketSessions))
	found[schedule.Session](t, false)(day.FirstSession(schedule.SessionFilter{Type: 9}))

	afterTwo := schedule.SessionFilter{Accept: func(x schedule.Session) bool {
		return x.Day() == day && !x.StartTime().Before(dayStart(42).Add(2*hour))
	}}
	if got := found[schedule.Session](t, true)(regular.NextSession(afterTwo)); got != sessions[2] {
		t.Errorf("NextSession(custom) = %v", got)
	}
	found[schedule.Session](t, false)(day.FirstSession(schedule.SessionFilter{Trading: schedule.Yes,
		Accept: func(schedule.Session) bool { return false }}))

	nearest := found[schedule.Session](t, true)(s.NearestSessionByTime(dayStart(42).Add(5*hour),
		schedule.TradingSessions))
	if nearest != regular {
		t.Errorf("NearestSessionByTime = %v", nearest)
	}
}

func TestDays(t *testing.T) {
	s := ok[schedule.Schedule](t)(schedule.FromDefinition("(tz=GMT;td=12345;0=09001700)"))
	friday := ok[schedule.Day](t)(s.DayByID(43))
	if friday.DayOfWeek() != time.Friday || !friday.IsTrading() {
		t.Fatalf("day 43 = %v %v", friday, friday.DayOfWeek())
	}
	monday := found[schedule.Day](t, true)(friday.NextDay(schedule.TradingDays))
	if monday.DayID() != 46 || monday.DayOfWeek() != time.Monday {
		t.Errorf("NextDay(TradingDays) = %v", monday)
	}
	saturday := found[schedule.Day](t, true)(friday.NextDay(schedule.NonTradingDays))
	if saturday.DayID() != 44 || saturday.IsTrading() {
		t.Errorf("NextDay(NonTradingDays) = %v", saturday)
	}
	sunday := found[schedule.Day](t, true)(monday.PrevDay(schedule.Sundays))
	if sunday.DayID() != 45 || sunday.DayOfWeek() != time.Sunday {
		t.Errorf("PrevDay(Sundays) = %v", sunday)
	}
	weekEnd := found[schedule.Day](t, true)(monday.NextDay(schedule.WeekEnds))
	if weekEnd.DayID() != 51 {
		t.Errorf("NextDay(WeekEnds) = %v", weekEnd)
	}
	custom := schedule.DayFilter{Trading: schedule.Yes, Accept: func(d schedule.Day) bool { return d.DayOfMonth()%2 == 0 }}
	even := found[schedule.Day](t, true)(friday.NextDay(custom))
	if even.YearMonthDay() != 19700216 {
		t.Errorf("NextDay(custom) = %v", even)
	}
	found[schedule.Day](t, false)(friday.NextDay(schedule.DayFilter{Weekdays: []time.Weekday{9}}))
	if got := found[schedule.Session](t, true)(saturday.FirstSession(schedule.AnySession)); got.Type() !=
		schedule.NoTrading {
		t.Errorf("weekend session = %v", got)
	}

	holiday := ok[schedule.Schedule](t)(schedule.FromDefinition("(tz=GMT;td=12345;hd=19700216;0=09001700)"))
	friday = ok[schedule.Day](t)(holiday.DayByID(43))
	tuesday := found[schedule.Day](t, true)(friday.NextDay(schedule.TradingDays))
	if tuesday.DayID() != 47 {
		t.Errorf("NextDay over a holiday = %v", tuesday)
	}
	h := found[schedule.Day](t, true)(friday.NextDay(schedule.Holidays))
	if h.DayID() != 46 || !h.IsHoliday() || h.IsTrading() {
		t.Errorf("NextDay(Holidays) = %v", h)
	}
}

func TestDefaults(t *testing.T) {
	s := ok[schedule.Schedule](t)(schedule.FromDefinition("NewYorkETH()"))
	if s.Name() != "US ETH" || s.TimeZoneID() != "America/New_York" {
		t.Fatalf("schedule = %q %q", s.Name(), s.TimeZoneID())
	}
	eve := ok[schedule.Day](t)(s.DayByYearMonthDay(20241224))
	if !eve.IsShortDay() || !eve.IsTrading() {
		t.Errorf("2024-12-24 = %v", eve)
	}
	christmas := ok[schedule.Day](t)(s.DayByYearMonthDay(20241225))
	if !christmas.IsHoliday() || christmas.IsTrading() {
		t.Errorf("2024-12-25 = %v", christmas)
	}
	regular := found[schedule.Session](t, true)(eve.FirstSession(schedule.RegularSessions))
	ny := ok[*time.Location](t)(time.LoadLocation("America/New_York"))
	if got := regular.EndTime().In(ny).Format("15:04"); got != "13:00" {
		t.Errorf("short day regular session ends at %s", got)
	}
}

func TestNoTrading(t *testing.T) {
	s := ok[schedule.Schedule](t)(schedule.FromDefinition("(tz=GMT;0=)"))
	day := ok[schedule.Day](t)(s.DayByID(42))
	if day.IsTrading() {
		t.Errorf("day = %v", day)
	}
	found[schedule.Day](t, false)(day.NextDay(schedule.TradingDays))
	found[schedule.Session](t, false)(s.NearestSessionByTime(dayStart(42), schedule.TradingSessions))
}

func TestFromProfile(t *testing.T) {
	profile := events.NewInstrumentProfile()
	hours := "UTC(tz=UTC;td=12345;de=+0000;0=0000+0000)"
	profile.SetTradingHours(&hours)
	s := ok[schedule.Schedule](t)(schedule.FromProfile(profile))
	day := ok[schedule.Day](t)(s.DayByID(0))
	if day.DayOfWeek() != time.Thursday || day.YearMonthDay() != 19700101 || !day.IsTrading() {
		t.Errorf("day 0 = %v %v", day, day.DayOfWeek())
	}
	if _, err := schedule.FromProfile(nil); err == nil {
		t.Error("FromProfile(nil) did not fail")
	}
	if _, err := schedule.FromProfileAndVenue(profile, ""); err == nil {
		t.Error("FromProfileAndVenue with an empty venue did not fail")
	}
	venues := ok[[]string](t)(schedule.TradingVenues(profile))
	t.Logf("venues = %v", venues)
}

func TestInvalid(t *testing.T) {
	for _, def := range []string{"", "a\x00b", "(tz=GMT;0=badvalue)"} {
		if _, err := schedule.FromDefinition(def); err == nil {
			t.Errorf("FromDefinition(%q) did not fail", def)
		}
	}
	if s, err := schedule.FromDefinition("(tz=GMT;0=)"); err == nil {
		if _, err := s.DayByYearMonthDay(1 << 40); err == nil {
			t.Error("DayByYearMonthDay out of range did not fail")
		}
	}
}

func TestSessionTypeString(t *testing.T) {
	if schedule.Regular.String() != "REGULAR" || schedule.SessionType(0).String() != "SessionType(0)" {
		t.Error(schedule.Regular.String(), schedule.SessionType(0).String())
	}
	if !schedule.PreMarket.IsTrading() || schedule.NoTrading.IsTrading() {
		t.Error("IsTrading")
	}
}
