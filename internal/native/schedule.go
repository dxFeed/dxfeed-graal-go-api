// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>

extern int32_t OnScheduleDayFilter(graal_isolatethread_t *thread, dxfg_v4_day_t *day, void *user_data);
extern int32_t OnScheduleSessionFilter(graal_isolatethread_t *thread, dxfg_v4_session_t *session, void *user_data);
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// A schedule is kept in Go as its definition; every call gets the native schedule of the definition (the SDK caches
// the schedules by definition), finds the day by its id and the session by its index in the day, and releases them
// afterwards.

// FilterAny is the value of a filter condition that is not checked.
const FilterAny = int32(C.DXFG_V4_FILTER_ANY)

// ScheduleFields are the values of a schedule.
type ScheduleFields struct {
	Definition, Name, TimeZone, TimeZoneID string
}

// DayFields are the values of a day of a schedule.
type DayFields struct {
	DayID, YearMonthDay           int32
	Holiday, ShortDay, Trading    bool
	StartTime, EndTime, ResetTime int64
}

// SessionFields are the values of a session of a schedule: its day and its index in the day.
type SessionFields struct {
	Day                DayFields
	Index, Type        int32
	Trading, Empty     bool
	StartTime, EndTime int64
}

// DayFilterSpec is a day filter: the conditions (FilterAny or 0 for the mask are not checked) and a function that
// tests the days that pass them; a nil function accepts them.
type DayFilterSpec struct {
	DayOfWeekMask, Holiday, ShortDay, Trading int32
	Accept                                    func(DayFields) bool
}

// SessionFilterSpec is a session filter: the conditions and a function that tests the sessions that pass them; a
// nil function accepts them.
type SessionFilterSpec struct {
	Type, Trading int32
	Accept        func(SessionFields) bool
}

func scheduleString(thread *isolateThread, get func(**C.char) C.int32_t) (string, error) {
	var s *C.char
	if err := checkStatus(thread, get(&s)); err != nil {
		return "", err
	}
	if s == nil {
		return "", nil
	}
	defer C.dxfg_v4_String_release(thread.ptr, s)
	return C.GoString(s), nil
}

func scheduleInt32(thread *isolateThread, get func(*C.int32_t) C.int32_t) (int32, error) {
	var v C.int32_t
	err := checkStatus(thread, get(&v))
	return int32(v), err
}

func scheduleInt64(thread *isolateThread, get func(*C.int64_t) C.int32_t) (int64, error) {
	var v C.int64_t
	err := checkStatus(thread, get(&v))
	return int64(v), err
}

func scheduleFields(thread *isolateThread, s *C.dxfg_v4_schedule_t) (ScheduleFields, error) {
	var f ScheduleFields
	var err error
	if f.Definition, err = scheduleString(thread, func(p **C.char) C.int32_t {
		return C.dxfg_v4_Schedule_getDefinition(thread.ptr, s, p)
	}); err != nil {
		return f, err
	}
	if f.Name, err = scheduleString(thread, func(p **C.char) C.int32_t {
		return C.dxfg_v4_Schedule_getName(thread.ptr, s, p)
	}); err != nil {
		return f, err
	}
	if f.TimeZone, err = scheduleString(thread, func(p **C.char) C.int32_t {
		return C.dxfg_v4_Schedule_getTimeZone(thread.ptr, s, p)
	}); err != nil {
		return f, err
	}
	f.TimeZoneID, err = scheduleString(thread, func(p **C.char) C.int32_t {
		return C.dxfg_v4_Schedule_getTimeZone_getID(thread.ptr, s, p)
	})
	return f, err
}

func dayFields(thread *isolateThread, d *C.dxfg_v4_day_t) (DayFields, error) {
	var f DayFields
	var holiday, shortDay, trading int32
	var err error
	for _, get := range []struct {
		dst *int32
		f   func(*C.int32_t) C.int32_t
	}{
		{&f.DayID, func(v *C.int32_t) C.int32_t { return C.dxfg_v4_Day_getDayId(thread.ptr, d, v) }},
		{&f.YearMonthDay, func(v *C.int32_t) C.int32_t { return C.dxfg_v4_Day_getYearMonthDay(thread.ptr, d, v) }},
		{&holiday, func(v *C.int32_t) C.int32_t { return C.dxfg_v4_Day_isHoliday(thread.ptr, d, v) }},
		{&shortDay, func(v *C.int32_t) C.int32_t { return C.dxfg_v4_Day_isShortDay(thread.ptr, d, v) }},
		{&trading, func(v *C.int32_t) C.int32_t { return C.dxfg_v4_Day_isTrading(thread.ptr, d, v) }},
	} {
		if *get.dst, err = scheduleInt32(thread, get.f); err != nil {
			return f, err
		}
	}
	f.Holiday, f.ShortDay, f.Trading = holiday != 0, shortDay != 0, trading != 0
	for _, get := range []struct {
		dst *int64
		f   func(*C.int64_t) C.int32_t
	}{
		{&f.StartTime, func(v *C.int64_t) C.int32_t { return C.dxfg_v4_Day_getStartTime(thread.ptr, d, v) }},
		{&f.EndTime, func(v *C.int64_t) C.int32_t { return C.dxfg_v4_Day_getEndTime(thread.ptr, d, v) }},
		{&f.ResetTime, func(v *C.int64_t) C.int32_t { return C.dxfg_v4_Day_getResetTime(thread.ptr, d, v) }},
	} {
		if *get.dst, err = scheduleInt64(thread, get.f); err != nil {
			return f, err
		}
	}
	return f, nil
}

// sessionFieldsAt reads a session whose day and index are known.
func sessionFieldsAt(thread *isolateThread, s *C.dxfg_v4_session_t, day DayFields, index int) (SessionFields, error) {
	f := SessionFields{Day: day, Index: int32(index)}
	var sessionType C.dxfg_v4_session_type_t
	if err := checkStatus(thread, C.dxfg_v4_Session_getType(thread.ptr, s, &sessionType)); err != nil {
		return f, err
	}
	f.Type = int32(sessionType)
	trading, err := scheduleInt32(thread, func(v *C.int32_t) C.int32_t {
		return C.dxfg_v4_Session_isTrading(thread.ptr, s, v)
	})
	if err != nil {
		return f, err
	}
	empty, err := scheduleInt32(thread, func(v *C.int32_t) C.int32_t {
		return C.dxfg_v4_Session_isEmpty(thread.ptr, s, v)
	})
	if err != nil {
		return f, err
	}
	f.Trading, f.Empty = trading != 0, empty != 0
	if f.StartTime, err = scheduleInt64(thread, func(v *C.int64_t) C.int32_t {
		return C.dxfg_v4_Session_getStartTime(thread.ptr, s, v)
	}); err != nil {
		return f, err
	}
	f.EndTime, err = scheduleInt64(thread, func(v *C.int64_t) C.int32_t {
		return C.dxfg_v4_Session_getEndTime(thread.ptr, s, v)
	})
	return f, err
}

func daySessions(thread *isolateThread, d *C.dxfg_v4_day_t, f func([]*C.dxfg_v4_session_t) error) error {
	var list *C.dxfg_v4_session_list
	if err := checkStatus(thread, C.dxfg_v4_Day_getSessions(thread.ptr, d, &list)); err != nil {
		return err
	}
	defer C.dxfg_v4_session_list_release(thread.ptr, list)
	if list == nil || list.size == 0 {
		return f(nil)
	}
	return f(unsafe.Slice(list.elements, int(list.size)))
}

// sessionFields reads a session, finding its day and its index in the day.
func sessionFields(thread *isolateThread, s *C.dxfg_v4_session_t) (SessionFields, error) {
	var d *C.dxfg_v4_day_t
	if err := checkStatus(thread, C.dxfg_v4_Session_getDay(thread.ptr, s, &d)); err != nil {
		return SessionFields{}, err
	}
	defer C.dxfg_v4_Day_release(thread.ptr, d)
	day, err := dayFields(thread, d)
	if err != nil {
		return SessionFields{}, err
	}
	var result SessionFields
	err = daySessions(thread, d, func(sessions []*C.dxfg_v4_session_t) error {
		for i, other := range sessions {
			equal, err := scheduleInt32(thread, func(v *C.int32_t) C.int32_t {
				return C.dxfg_v4_Session_equals(thread.ptr, s, other, v)
			})
			if err != nil {
				return err
			}
			if equal != 0 {
				result, err = sessionFieldsAt(thread, s, day, i)
				return err
			}
		}
		return errors.New("native: the session is not among the sessions of its day")
	})
	return result, err
}

func withSchedule(thread *isolateThread, definition string, f func(s *C.dxfg_v4_schedule_t) error) error {
	ptr, free := cString(definition)
	defer free()
	var s *C.dxfg_v4_schedule_t
	if err := checkStatus(thread, C.dxfg_v4_Schedule_getInstanceByDefinition(thread.ptr, ptr, &s)); err != nil {
		return err
	}
	defer releaseJavaObject(thread, unsafe.Pointer(s))
	return f(s)
}

func withDay(thread *isolateThread, definition string, dayID int32, f func(d *C.dxfg_v4_day_t) error) error {
	return withSchedule(thread, definition, func(s *C.dxfg_v4_schedule_t) error {
		var d *C.dxfg_v4_day_t
		if err := checkStatus(thread, C.dxfg_v4_Schedule_getDayById(thread.ptr, s, C.int32_t(dayID), &d)); err != nil {
			return err
		}
		defer C.dxfg_v4_Day_release(thread.ptr, d)
		return f(d)
	})
}

func withSession(thread *isolateThread, definition string, dayID, index int32,
	f func(s *C.dxfg_v4_session_t) error) error {
	return withDay(thread, definition, dayID, func(d *C.dxfg_v4_day_t) error {
		return daySessions(thread, d, func(sessions []*C.dxfg_v4_session_t) error {
			if index < 0 || int(index) >= len(sessions) {
				return fmt.Errorf("native: no session %d in the day %d", index, dayID)
			}
			return f(sessions[index])
		})
	})
}

func checkDefinition(definition string) error {
	if definition == "" {
		return errors.New("native: empty schedule definition")
	}
	return checkNoNul(definition)
}

// ScheduleOfDefinition returns the schedule of a definition.
func ScheduleOfDefinition(definition string) (ScheduleFields, error) {
	if err := checkDefinition(definition); err != nil {
		return ScheduleFields{}, err
	}
	var result ScheduleFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withSchedule(thread, definition, func(s *C.dxfg_v4_schedule_t) (err error) {
			result, err = scheduleFields(thread, s)
			return err
		})
	})
	return result, err
}

// ScheduleOfProfile returns the schedule of the trading hours of a profile, of the venue if it is not empty.
func ScheduleOfProfile(profile *events.InstrumentProfile, venue string) (ScheduleFields, error) {
	profiles := []*events.InstrumentProfile{profile}
	if err := checkProfiles(profiles); err != nil {
		return ScheduleFields{}, err
	}
	if err := checkNoNul(venue); err != nil {
		return ScheduleFields{}, err
	}
	var result ScheduleFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		p, err := newCProfile(thread, profile)
		if err != nil {
			return err
		}
		defer freeCProfile(thread, p)
		var s *C.dxfg_v4_schedule_t
		if venue == "" {
			err = checkStatus(thread, C.dxfg_v4_Schedule_getInstance(thread.ptr, p, &s))
		} else {
			venuePtr, free := cString(venue)
			defer free()
			err = checkStatus(thread, C.dxfg_v4_Schedule_getInstanceWithVenue(thread.ptr, p, venuePtr, &s))
		}
		if err != nil {
			return err
		}
		defer releaseJavaObject(thread, unsafe.Pointer(s))
		result, err = scheduleFields(thread, s)
		return err
	})
	return result, err
}

// ScheduleTradingVenues returns the venues of the trading hours of a profile.
func ScheduleTradingVenues(profile *events.InstrumentProfile) ([]string, error) {
	if err := checkProfiles([]*events.InstrumentProfile{profile}); err != nil {
		return nil, err
	}
	var result []string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		p, err := newCProfile(thread, profile)
		if err != nil {
			return err
		}
		defer freeCProfile(thread, p)
		var list *C.dxfg_v4_string_list
		if err := checkStatus(thread, C.dxfg_v4_Schedule_getTradingVenues(thread.ptr, p, &list)); err != nil {
			return err
		}
		defer releaseStringList(thread, list)
		result = goStrings(list)
		return nil
	})
	return result, err
}

// ScheduleDownloadDefaults downloads the default schedule data: an empty config stops the periodic download, a URL
// downloads once, and "<period>@<URL>" downloads periodically.
func ScheduleDownloadDefaults(config string) error {
	if err := checkNoNul(config); err != nil {
		return err
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		ptr, free := cString(config)
		defer free()
		return checkStatus(thread, C.dxfg_v4_Schedule_downloadDefaults(thread.ptr, ptr))
	})
}

// ScheduleSetDefaults sets the default schedule data, plain or compressed.
func ScheduleSetDefaults(data []byte) error {
	if len(data) > 1<<31-1 {
		return errors.New("native: the schedule defaults are too large")
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		ptr := C.CBytes(data)
		defer C.free(ptr)
		return checkStatus(thread, C.dxfg_v4_Schedule_setDefaults(thread.ptr, (*C.char)(ptr), C.int32_t(len(data))))
	})
}

// ScheduleDayByID returns the day of a schedule by its id.
func ScheduleDayByID(definition string, dayID int32) (DayFields, error) {
	if err := checkDefinition(definition); err != nil {
		return DayFields{}, err
	}
	var result DayFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withDay(thread, definition, dayID, func(d *C.dxfg_v4_day_t) (err error) {
			result, err = dayFields(thread, d)
			return err
		})
	})
	return result, err
}

// ScheduleDayByTime returns the day of a schedule that contains a time (milliseconds since the epoch).
func ScheduleDayByTime(definition string, time int64) (DayFields, error) {
	return scheduleDay(definition, func(thread *isolateThread, s *C.dxfg_v4_schedule_t, d **C.dxfg_v4_day_t) C.int32_t {
		return C.dxfg_v4_Schedule_getDayByTime(thread.ptr, s, C.int64_t(time), d)
	})
}

// ScheduleDayByYearMonthDay returns the day of a schedule by its yyyymmdd number.
func ScheduleDayByYearMonthDay(definition string, yearMonthDay int32) (DayFields, error) {
	return scheduleDay(definition, func(thread *isolateThread, s *C.dxfg_v4_schedule_t, d **C.dxfg_v4_day_t) C.int32_t {
		return C.dxfg_v4_Schedule_getDayByYearMonthDay(thread.ptr, s, C.int32_t(yearMonthDay), d)
	})
}

func scheduleDay(definition string,
	get func(thread *isolateThread, s *C.dxfg_v4_schedule_t, d **C.dxfg_v4_day_t) C.int32_t) (DayFields, error) {
	if err := checkDefinition(definition); err != nil {
		return DayFields{}, err
	}
	var result DayFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withSchedule(thread, definition, func(s *C.dxfg_v4_schedule_t) error {
			var d *C.dxfg_v4_day_t
			if err := checkStatus(thread, get(thread, s, &d)); err != nil {
				return err
			}
			defer C.dxfg_v4_Day_release(thread.ptr, d)
			var err error
			result, err = dayFields(thread, d)
			return err
		})
	})
	return result, err
}

// ScheduleSessionByTime returns the session of a schedule that contains a time.
func ScheduleSessionByTime(definition string, time int64) (SessionFields, error) {
	if err := checkDefinition(definition); err != nil {
		return SessionFields{}, err
	}
	var result SessionFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withSchedule(thread, definition, func(s *C.dxfg_v4_schedule_t) error {
			var session *C.dxfg_v4_session_t
			if err := checkStatus(thread, C.dxfg_v4_Schedule_getSessionByTime(thread.ptr, s, C.int64_t(time),
				&session)); err != nil {
				return err
			}
			defer C.dxfg_v4_Session_release(thread.ptr, session)
			var err error
			result, err = sessionFields(thread, session)
			return err
		})
	})
	return result, err
}

// ScheduleNearestSession returns the session of a schedule nearest to a time that the filter accepts; ok is false
// if there is none.
func ScheduleNearestSession(definition string, time int64, filter SessionFilterSpec) (SessionFields, bool, error) {
	if err := checkDefinition(definition); err != nil {
		return SessionFields{}, false, err
	}
	return findSession(func(thread *isolateThread, f *C.dxfg_v4_session_filter_t,
		get func(func(**C.dxfg_v4_session_t) C.int32_t) error) error {
		return withSchedule(thread, definition, func(s *C.dxfg_v4_schedule_t) error {
			return get(func(session **C.dxfg_v4_session_t) C.int32_t {
				return C.dxfg_v4_Schedule_findNearestSessionByTime(thread.ptr, s, C.int64_t(time), f, session)
			})
		})
	}, filter)
}

// DaySessions returns the sessions of a day in their order.
func DaySessions(definition string, dayID int32) ([]SessionFields, error) {
	if err := checkDefinition(definition); err != nil {
		return nil, err
	}
	var result []SessionFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withDay(thread, definition, dayID, func(d *C.dxfg_v4_day_t) error {
			day, err := dayFields(thread, d)
			if err != nil {
				return err
			}
			return daySessions(thread, d, func(sessions []*C.dxfg_v4_session_t) error {
				result = make([]SessionFields, 0, len(sessions))
				for i, s := range sessions {
					f, err := sessionFieldsAt(thread, s, day, i)
					if err != nil {
						return err
					}
					result = append(result, f)
				}
				return nil
			})
		})
	})
	return result, err
}

// DaySessionByTime returns the session of a day that contains a time.
func DaySessionByTime(definition string, dayID int32, time int64) (SessionFields, error) {
	if err := checkDefinition(definition); err != nil {
		return SessionFields{}, err
	}
	var result SessionFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withDay(thread, definition, dayID, func(d *C.dxfg_v4_day_t) error {
			var session *C.dxfg_v4_session_t
			if err := checkStatus(thread, C.dxfg_v4_Day_getSessionByTime(thread.ptr, d, C.int64_t(time),
				&session)); err != nil {
				return err
			}
			defer C.dxfg_v4_Session_release(thread.ptr, session)
			var err error
			result, err = sessionFields(thread, session)
			return err
		})
	})
	return result, err
}

// DayEdgeSession returns the first (or the last) session of a day that the filter accepts; ok is false if there is
// none.
func DayEdgeSession(definition string, dayID int32, last bool, filter SessionFilterSpec) (SessionFields, bool, error) {
	if err := checkDefinition(definition); err != nil {
		return SessionFields{}, false, err
	}
	return findSession(func(thread *isolateThread, f *C.dxfg_v4_session_filter_t,
		get func(func(**C.dxfg_v4_session_t) C.int32_t) error) error {
		return withDay(thread, definition, dayID, func(d *C.dxfg_v4_day_t) error {
			return get(func(session **C.dxfg_v4_session_t) C.int32_t {
				if last {
					return C.dxfg_v4_Day_findLastSession(thread.ptr, d, f, session)
				}
				return C.dxfg_v4_Day_findFirstSession(thread.ptr, d, f, session)
			})
		})
	}, filter)
}

// SessionNeighbor returns the next (or the previous) session that the filter accepts; ok is false if there is none.
func SessionNeighbor(definition string, dayID, index int32, next bool, filter SessionFilterSpec) (SessionFields, bool,
	error) {
	if err := checkDefinition(definition); err != nil {
		return SessionFields{}, false, err
	}
	return findSession(func(thread *isolateThread, f *C.dxfg_v4_session_filter_t,
		get func(func(**C.dxfg_v4_session_t) C.int32_t) error) error {
		return withSession(thread, definition, dayID, index, func(s *C.dxfg_v4_session_t) error {
			return get(func(session **C.dxfg_v4_session_t) C.int32_t {
				if next {
					return C.dxfg_v4_Session_findNextSession(thread.ptr, s, f, session)
				}
				return C.dxfg_v4_Session_findPrevSession(thread.ptr, s, f, session)
			})
		})
	}, filter)
}

// findSession runs a search with a native session filter; the search calls get with the native find function.
func findSession(search func(thread *isolateThread, f *C.dxfg_v4_session_filter_t,
	get func(func(**C.dxfg_v4_session_t) C.int32_t) error) error, filter SessionFilterSpec) (SessionFields, bool,
	error) {
	var result SessionFields
	found := false
	err := executeInIsolateThread(func(thread *isolateThread) error {
		f, release, err := newSessionFilter(thread, filter)
		if err != nil {
			return err
		}
		defer release()
		return search(thread, f, func(find func(**C.dxfg_v4_session_t) C.int32_t) error {
			var session *C.dxfg_v4_session_t
			if err := checkStatus(thread, find(&session)); err != nil {
				return err
			}
			if session == nil {
				return nil
			}
			defer C.dxfg_v4_Session_release(thread.ptr, session)
			found = true
			var err error
			result, err = sessionFields(thread, session)
			return err
		})
	})
	return result, found && err == nil, err
}

// DayNeighbor returns the next (or the previous) day that the filter accepts; ok is false if there is none.
func DayNeighbor(definition string, dayID int32, next bool, filter DayFilterSpec) (DayFields, bool, error) {
	if err := checkDefinition(definition); err != nil {
		return DayFields{}, false, err
	}
	var result DayFields
	found := false
	err := executeInIsolateThread(func(thread *isolateThread) error {
		f, release, err := newDayFilter(thread, filter)
		if err != nil {
			return err
		}
		defer release()
		return withDay(thread, definition, dayID, func(d *C.dxfg_v4_day_t) error {
			var other *C.dxfg_v4_day_t
			var status C.int32_t
			if next {
				status = C.dxfg_v4_Day_findNextDay(thread.ptr, d, f, &other)
			} else {
				status = C.dxfg_v4_Day_findPrevDay(thread.ptr, d, f, &other)
			}
			if err := checkStatus(thread, status); err != nil {
				return err
			}
			if other == nil {
				return nil
			}
			defer C.dxfg_v4_Day_release(thread.ptr, other)
			found = true
			var err error
			result, err = dayFields(thread, other)
			return err
		})
	})
	return result, found && err == nil, err
}

// dayAcceptor and sessionAcceptor are the user data of the native filters.
type dayAcceptor func(thread *isolateThread, d *C.dxfg_v4_day_t) bool
type sessionAcceptor func(thread *isolateThread, s *C.dxfg_v4_session_t) bool

//export OnScheduleDayFilter
func OnScheduleDayFilter(thread *C.graal_isolatethread_t, day *C.dxfg_v4_day_t, userData unsafe.Pointer) (
	accepted C.int32_t) {
	defer recoverCallback("schedule day filter")
	if accept, ok := Restore(userData).(dayAcceptor); ok && accept(&isolateThread{ptr: thread}, day) {
		return 1
	}
	return 0
}

//export OnScheduleSessionFilter
func OnScheduleSessionFilter(thread *C.graal_isolatethread_t, session *C.dxfg_v4_session_t,
	userData unsafe.Pointer) (accepted C.int32_t) {
	defer recoverCallback("schedule session filter")
	if accept, ok := Restore(userData).(sessionAcceptor); ok && accept(&isolateThread{ptr: thread}, session) {
		return 1
	}
	return 0
}

// newDayFilter creates a native day filter that is used during this call only; release it with the returned
// function after the search.
func newDayFilter(thread *isolateThread, spec DayFilterSpec) (*C.dxfg_v4_day_filter_t, func(), error) {
	userData := Save(dayAcceptor(func(thread *isolateThread, d *C.dxfg_v4_day_t) bool {
		if spec.Accept == nil {
			return true
		}
		fields, err := dayFields(thread, d)
		if err != nil {
			logGo(LogLevelError, "schedule day filter: %v", err)
			return false
		}
		return spec.Accept(fields)
	}))
	var f *C.dxfg_v4_day_filter_t
	if err := checkStatus(thread, C.dxfg_v4_DayFilter_newWithUserData(thread.ptr, C.int32_t(spec.DayOfWeekMask),
		C.int32_t(spec.Holiday), C.int32_t(spec.ShortDay), C.int32_t(spec.Trading),
		(*[0]byte)(C.OnScheduleDayFilter), userData, &f)); err != nil {
		Unref(userData)
		return nil, nil, err
	}
	return f, func() {
		releaseJavaObject(thread, unsafe.Pointer(f))
		Unref(userData)
	}, nil
}

// newSessionFilter creates a native session filter like newDayFilter.
func newSessionFilter(thread *isolateThread, spec SessionFilterSpec) (*C.dxfg_v4_session_filter_t, func(), error) {
	userData := Save(sessionAcceptor(func(thread *isolateThread, s *C.dxfg_v4_session_t) bool {
		if spec.Accept == nil {
			return true
		}
		fields, err := sessionFields(thread, s)
		if err != nil {
			logGo(LogLevelError, "schedule session filter: %v", err)
			return false
		}
		return spec.Accept(fields)
	}))
	var f *C.dxfg_v4_session_filter_t
	if err := checkStatus(thread, C.dxfg_v4_SessionFilter_newWithUserData(thread.ptr, C.int32_t(spec.Type),
		C.int32_t(spec.Trading), (*[0]byte)(C.OnScheduleSessionFilter), userData, &f)); err != nil {
		Unref(userData)
		return nil, nil, err
	}
	return f, func() {
		releaseJavaObject(thread, unsafe.Pointer(f))
		Unref(userData)
	}, nil
}
