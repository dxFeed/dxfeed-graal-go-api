// Package schedule tells the trading days and sessions of an instrument: a [Schedule] comes from the trading hours of
// an instrument profile or from a schedule definition, and splits time into [Day]s, each made of [Session]s of a
// [SessionType]. [DayFilter] and [SessionFilter] choose the days and sessions to find, such as the next trading day.
//
//	s, _ := schedule.FromProfile(profile)
//	day, _ := s.DayByTime(time.Now())
//	next, ok, _ := day.NextDay(schedule.TradingDays)
//
// The schedules, days and sessions are immutable values; they can be compared with == and kept in maps. The default
// schedule data that the definitions such as NewYorkETH() refer to are built in; DownloadDefaults and SetDefaults change
// them.
package schedule
