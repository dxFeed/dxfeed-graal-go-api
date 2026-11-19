package schedule

import (
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// Cond is a condition of a filter on a yes-or-no property. Its zero value, Any, does not check the property.
type Cond int

// The conditions.
const (
	// Any accepts any value of the property.
	Any Cond = iota
	// Yes accepts the days or sessions that have the property.
	Yes
	// No accepts the days or sessions that do not have the property.
	No
)

func (c Cond) native() int32 {
	switch c {
	case Yes:
		return 1
	case No:
		return 0
	default:
		return native.FilterAny
	}
}

// DayFilter chooses days. A day passes when it falls on one of the Weekdays (any day if empty), meets the
// conditions, and Accept, if set, returns true for it. The zero value accepts any day.
type DayFilter struct {
	Weekdays                   []time.Weekday
	Holiday, ShortDay, Trading Cond
	Accept                     func(Day) bool
}

// The predefined day filters.
var (
	AnyDay         = DayFilter{}
	TradingDays    = DayFilter{Trading: Yes}
	NonTradingDays = DayFilter{Trading: No}
	Holidays       = DayFilter{Holiday: Yes}
	ShortDays      = DayFilter{ShortDay: Yes}
	Mondays        = DayFilter{Weekdays: []time.Weekday{time.Monday}}
	Tuesdays       = DayFilter{Weekdays: []time.Weekday{time.Tuesday}}
	Wednesdays     = DayFilter{Weekdays: []time.Weekday{time.Wednesday}}
	Thursdays      = DayFilter{Weekdays: []time.Weekday{time.Thursday}}
	Fridays        = DayFilter{Weekdays: []time.Weekday{time.Friday}}
	Saturdays      = DayFilter{Weekdays: []time.Weekday{time.Saturday}}
	Sundays        = DayFilter{Weekdays: []time.Weekday{time.Sunday}}
	WeekDays       = DayFilter{Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday}}
	WeekEnds = DayFilter{Weekdays: []time.Weekday{time.Saturday, time.Sunday}}
)

func (f DayFilter) spec(s Schedule) native.DayFilterSpec {
	spec := native.DayFilterSpec{
		Holiday:  f.Holiday.native(),
		ShortDay: f.ShortDay.native(),
		Trading:  f.Trading.native(),
	}
	for _, w := range f.Weekdays {
		// The native mask has bit N for the ISO day of the week N: 1 for Monday, ..., 7 for Sunday.
		iso := int(w)
		if w == time.Sunday {
			iso = 7
		}
		if iso >= 1 && iso <= 7 {
			spec.DayOfWeekMask |= 1 << iso
		}
	}
	if accept := f.Accept; accept != nil {
		spec.Accept = func(d native.DayFields) bool { return accept(s.day(d)) }
	}
	return spec
}

// SessionFilter chooses sessions. A session passes when it has the Type (any type if zero), meets the condition,
// and Accept, if set, returns true for it. The zero value accepts any session.
type SessionFilter struct {
	Type    SessionType
	Trading Cond
	Accept  func(Session) bool
}

// The predefined session filters.
var (
	AnySession          = SessionFilter{}
	TradingSessions     = SessionFilter{Trading: Yes}
	NonTradingSessions  = SessionFilter{Trading: No}
	NoTradingSessions   = SessionFilter{Type: NoTrading}
	PreMarketSessions   = SessionFilter{Type: PreMarket}
	RegularSessions     = SessionFilter{Type: Regular}
	AfterMarketSessions = SessionFilter{Type: AfterMarket}
)

func (f SessionFilter) spec(s Schedule) native.SessionFilterSpec {
	spec := native.SessionFilterSpec{Type: native.FilterAny, Trading: f.Trading.native()}
	switch {
	case f.Type >= NoTrading && f.Type <= AfterMarket:
		spec.Type = int32(f.Type) - 1
	case f.Type != 0:
		// An unknown type matches no session.
		spec.Accept = func(native.SessionFields) bool { return false }
		return spec
	}
	if accept := f.Accept; accept != nil {
		spec.Accept = func(fields native.SessionFields) bool { return accept(s.session(fields)) }
	}
	return spec
}
