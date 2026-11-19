package schedule

import (
	"fmt"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// SessionType is the type of a session. Its zero value is not a type: a [SessionFilter] with it accepts any type.
type SessionType int

// The session types.
const (
	// NoTrading is a session without trading, such as a night or a weekend.
	NoTrading SessionType = iota + 1
	// PreMarket is a trading session before the regular one.
	PreMarket
	// Regular is the regular trading session.
	Regular
	// AfterMarket is a trading session after the regular one.
	AfterMarket
)

var sessionTypeNames = [...]string{"", "NO_TRADING", "PRE_MARKET", "REGULAR", "AFTER_MARKET"}

// String returns the name of the type, such as "REGULAR".
func (t SessionType) String() string {
	if t >= NoTrading && t <= AfterMarket {
		return sessionTypeNames[t]
	}
	return fmt.Sprintf("SessionType(%d)", int(t))
}

// IsTrading reports whether the sessions of the type allow trading.
func (t SessionType) IsTrading() bool {
	return t >= PreMarket && t <= AfterMarket
}

// Session is a continuous period of time of a day with the same trading rules, such as the regular trading hours.
type Session struct {
	schedule Schedule
	f        native.SessionFields
}

// Day returns the day of the session.
func (s Session) Day() Day { return s.schedule.day(s.f.Day) }

// Type returns the type of the session.
func (s Session) Type() SessionType { return SessionType(s.f.Type + 1) }

// IsTrading reports whether the session allows trading.
func (s Session) IsTrading() bool { return s.f.Trading }

// IsEmpty reports whether the session has no duration: its start time equals its end time.
func (s Session) IsEmpty() bool { return s.f.Empty }

// StartTime returns the start of the session, inclusive.
func (s Session) StartTime() time.Time { return timeOf(s.f.StartTime) }

// EndTime returns the end of the session, exclusive.
func (s Session) EndTime() time.Time { return timeOf(s.f.EndTime) }

// ContainsTime reports whether the session contains a time.
func (s Session) ContainsTime(t time.Time) bool {
	millis := t.UnixMilli()
	return s.f.StartTime <= millis && millis < s.f.EndTime
}

// String returns the type and the period of the session in UTC, such as
// "Session(REGULAR, 2024-12-31T14:30:00Z-2024-12-31T21:00:00Z)".
func (s Session) String() string {
	return fmt.Sprintf("Session(%v, %s-%s)", s.Type(), s.StartTime().UTC().Format(time.RFC3339),
		s.EndTime().UTC().Format(time.RFC3339))
}

// PrevSession returns the nearest session before this one that the filter accepts, possibly on an earlier day; ok is
// false if there is none within about a year.
func (s Session) PrevSession(filter SessionFilter) (session Session, ok bool, err error) {
	f, ok, err := native.SessionNeighbor(s.schedule.f.Definition, s.f.Day.DayID, s.f.Index, false,
		filter.spec(s.schedule))
	return s.schedule.session(f), ok, err
}

// NextSession returns the nearest session after this one that the filter accepts, possibly on a later day; ok is
// false if there is none within about a year.
func (s Session) NextSession(filter SessionFilter) (session Session, ok bool, err error) {
	f, ok, err := native.SessionNeighbor(s.schedule.f.Definition, s.f.Day.DayID, s.f.Index, true,
		filter.spec(s.schedule))
	return s.schedule.session(f), ok, err
}
