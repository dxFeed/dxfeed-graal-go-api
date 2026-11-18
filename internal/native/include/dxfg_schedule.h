// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Schedule functions and types declarations
 */

#ifndef DXFG_SCHEDULE_H
#define DXFG_SCHEDULE_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"
#include "dxfg_ipf.h"

/** @defgroup Schedule Schedule
 *  Trading schedules of exchanges and instruments: schedules, days, sessions and the filters to search them.
 *
 *  All times are UTC milliseconds since the Unix epoch (1970-01-01T00:00:00Z). The supported date range is
 *  from 0001-01-02 to 9999-12-30.
 *
 *  Every function that returns a handle (a schedule, a day, a session or a filter) creates a new handle that the caller
 *  owns, even if the same Java object has been returned before: release it with the function named in its description.
 *  @{
 */

/**
 * @brief A handle of a Java Schedule object; release it with dxfg_JavaObjectHandler_release().
 *
 * A schedule describes the trading schedule of some class of instruments, for example NYSE stocks or CME corn
 * futures. It splits the whole time scale into consecutive days.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/Schedule.html">Javadoc</a>
 */
typedef struct dxfg_schedule_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_schedule_t;

/**
 * The type of a trading session: what kind of trading activity is allowed in it, if any.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/SessionType.html">Javadoc</a>
 */
typedef enum dxfg_session_type_t {
  /** A non-trading session: a period of time during which trading is not allowed. */
  DXFG_SESSION_TYPE_NO_TRADING = 0,
  /** A pre-market session: extended trading before the regular trading hours. */
  DXFG_SESSION_TYPE_PRE_MARKET,
  /** A regular session: the regular trading hours. */
  DXFG_SESSION_TYPE_REGULAR,
  /** An after-market session: extended trading after the regular trading hours. */
  DXFG_SESSION_TYPE_AFTER_MARKET
} dxfg_session_type_t;

/**
 * @brief A handle of a Java Session object; release it with dxfg_Session_release().
 *
 * A session is a continuous period of time with the same rules for trading activity. It is either trading or
 * non-trading and lies inside one day. Sessions do not overlap: they form a chain of adjacent periods, and the point
 * on the border belongs to the session that starts there. A session may have zero duration (an empty session).
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/Session.html">Javadoc</a>
 */
typedef struct dxfg_session_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_session_t;

/**
 * A list of session handles, returned by dxfg_Day_getSessions().
 *
 * Release the sessions with dxfg_Session_release() one by one and then the list itself with
 * dxfg_SessionList_wrapper_release(), or release the list together with its sessions with
 * dxfg_v4_session_list_release().
 */
typedef struct dxfg_session_list {
    /** The number of elements; 0 for an empty list. */
    int32_t size;
    /** The array of `size` session handles; NULL for an empty list. Each element is owned by the caller. */
    dxfg_session_t **elements;
} dxfg_session_list;

/**
 * @brief A handle of a Java Day object; release it with dxfg_Day_release().
 *
 * A day is a continuous period of time about 24 hours long, aligned to the business hours of the schedule, so days
 * start and end at various local times. Days do not overlap: they form a chain of adjacent periods, and the point on
 * the border belongs to the day that starts there. A day consists of sessions and is a trading day if at least one
 * of them is a trading session.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/Day.html">Javadoc</a>
 */
typedef struct dxfg_day_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_day_t;

/**
 * A custom session filter function: tests whether a session is an acceptable result of a search.
 *
 * The function is called by the search function that uses the filter (for example dxfg_Day_findFirstSession()),
 * on the thread that called that search function, possibly several times during the call.
 *
 * The session handle is created for this call and the SDK does not release it afterwards: release it with
 * dxfg_Session_release() unless it is needed later.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to test.
 * @return 1 if the session is accepted; any other value rejects it.
 */
typedef int32_t (*dxfg_session_filter_function)(graal_isolatethread_t *thread, dxfg_session_t *session);
/**
 * @brief A handle of a Java SessionFilter object; release it with dxfg_JavaObjectHandler_release().
 *
 * A filter for sessions used by the search functions.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/SessionFilter.html">Javadoc</a>
 */
typedef struct dxfg_session_filter_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_session_filter_t;
/**
 * The predefined session filters, see dxfg_SessionFilter_getInstance().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/SessionFilter.html">Javadoc</a>
 */
typedef enum dxfg_session_filter_prepare_t {
    /** Accepts any session. */
    DXFG_SESSION_FILTER_ANY = 0,
    /** Accepts trading sessions only. */
    DXFG_SESSION_FILTER_TRADING,
    /** Accepts non-trading sessions only. */
    DXFG_SESSION_FILTER_NON_TRADING,
    /** Accepts the sessions of type #DXFG_SESSION_TYPE_NO_TRADING. */
    DXFG_SESSION_FILTER_NO_TRADING,
    /** Accepts the sessions of type #DXFG_SESSION_TYPE_PRE_MARKET. */
    DXFG_SESSION_FILTER_PRE_MARKET,
    /** Accepts the sessions of type #DXFG_SESSION_TYPE_REGULAR. */
    DXFG_SESSION_FILTER_REGULAR,
    /** Accepts the sessions of type #DXFG_SESSION_TYPE_AFTER_MARKET. */
    DXFG_SESSION_FILTER_AFTER_MARKET
} dxfg_session_filter_prepare_t;

/**
 * Creates a session filter that calls a custom function to test sessions.
 *
 * The function pointer must stay valid as long as the filter may be used by search functions.
 *
 * @deprecated Use dxfg_v4_SessionFilter_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The function that tests a session: it returns 1 to accept the session.
 * @return The new session filter (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `new SessionFilter(SessionType, Boolean)` with `SessionFilter.accept(Session)` overridden to call `filter`.
 */
DXFG_DEPRECATED("use dxfg_v4_SessionFilter_new()")
dxfg_session_filter_t*    dxfg_SessionFilter_new(graal_isolatethread_t *thread, dxfg_session_filter_function filter);
/**
 * Returns a predefined session filter.
 *
 * @deprecated Use dxfg_v4_SessionFilter_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The predefined filter to return.
 * @return The session filter (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `SessionFilter.ANY`, `SessionFilter.TRADING`, `SessionFilter.NON_TRADING`, `SessionFilter.NO_TRADING`,
 * `SessionFilter.PRE_MARKET`, `SessionFilter.REGULAR`, `SessionFilter.AFTER_MARKET`.
 */
DXFG_DEPRECATED("use dxfg_v4_SessionFilter_getInstance()")
dxfg_session_filter_t*    dxfg_SessionFilter_getInstance(graal_isolatethread_t *thread, dxfg_session_filter_prepare_t filter);

/**
 * A custom day filter function: tests whether a day is an acceptable result of a search.
 *
 * The function is called by the search function that uses the filter (for example dxfg_Day_findNextDay()),
 * on the thread that called that search function, possibly several times during the call.
 *
 * The day handle is created for this call and the SDK does not release it afterwards: release it with
 * dxfg_Day_release() unless it is needed later.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to test.
 * @return 1 if the day is accepted; any other value rejects it.
 */
typedef int32_t (*dxfg_day_filter_function)(graal_isolatethread_t *thread, dxfg_day_t *day);
/**
 * @brief A handle of a Java DayFilter object; release it with dxfg_JavaObjectHandler_release().
 *
 * A filter for days used by the search functions.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/DayFilter.html">Javadoc</a>
 */
typedef struct dxfg_day_filter_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_day_filter_t;
/**
 * The predefined day filters, see dxfg_DayFilter_getInstance().
 *
 * The days of the week are numbered from 1 (Monday) to 7 (Sunday).
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/DayFilter.html">Javadoc</a>
 */
typedef enum dxfg_day_filter_prepare_t {
    /** Accepts any day. */
    DXFG_DAY_FILTER_ANY = 0,
    /** Accepts trading days only. */
    DXFG_DAY_FILTER_TRADING,
    /** Accepts non-trading days only. */
    DXFG_DAY_FILTER_NON_TRADING,
    /** Accepts holidays only. */
    DXFG_DAY_FILTER_HOLIDAY,
    /** Accepts short days only. */
    DXFG_DAY_FILTER_SHORT_DAY,
    /** Accepts Mondays only. */
    DXFG_DAY_FILTER_MONDAY,
    /** Accepts Tuesdays only. */
    DXFG_DAY_FILTER_TUESDAY,
    /** Accepts Wednesdays only. */
    DXFG_DAY_FILTER_WEDNESDAY,
    /** Accepts Thursdays only. */
    DXFG_DAY_FILTER_THURSDAY,
    /** Accepts Fridays only. */
    DXFG_DAY_FILTER_FRIDAY,
    /** Accepts Saturdays only. */
    DXFG_DAY_FILTER_SATURDAY,
    /** Accepts Sundays only. */
    DXFG_DAY_FILTER_SUNDAY,
    /** Accepts week days only: Monday to Friday. */
    DXFG_DAY_FILTER_WEEK_DAY,
    /** Accepts weekends only: Saturday and Sunday. */
    DXFG_DAY_FILTER_WEEK_END
} dxfg_day_filter_prepare_t;

/**
 * Creates a day filter that calls a custom function to test days.
 *
 * The function pointer must stay valid as long as the filter may be used by search functions.
 *
 * @deprecated Use dxfg_v4_DayFilter_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The function that tests a day: it returns 1 to accept the day.
 * @return The new day filter (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `new DayFilter(int, Boolean, Boolean, Boolean)` with `DayFilter.accept(Day)` overridden to call `filter`.
 */
DXFG_DEPRECATED("use dxfg_v4_DayFilter_new()")
dxfg_day_filter_t*    dxfg_DayFilter_new(graal_isolatethread_t *thread, dxfg_day_filter_function filter);
/**
 * Returns a predefined day filter.
 *
 * @deprecated Use dxfg_v4_DayFilter_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The predefined filter to return.
 * @return The day filter (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DayFilter.ANY`, `DayFilter.TRADING`, `DayFilter.NON_TRADING`, `DayFilter.HOLIDAY`, `DayFilter.SHORT_DAY`,
 * `DayFilter.MONDAY` ... `DayFilter.SUNDAY`, `DayFilter.WEEK_DAY`, `DayFilter.WEEK_END`.
 */
DXFG_DEPRECATED("use dxfg_v4_DayFilter_getInstance()")
dxfg_day_filter_t*    dxfg_DayFilter_getInstance(graal_isolatethread_t *thread, dxfg_day_filter_prepare_t filter);


/**
 * Returns the default schedule for an instrument profile.
 *
 * The schedule is defined by the trading hours of the profile. The function fails if the trading hours contain no
 * parenthesized schedule definition.
 *
 * @deprecated Use dxfg_v4_Schedule_getInstanceByProfileHandle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile whose schedule is requested.
 * @return The schedule (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getInstance(InstrumentProfile)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getInstanceByProfileHandle()")
dxfg_schedule_t*      dxfg_Schedule_getInstance(graal_isolatethread_t *thread, dxfg_instrument_profile_t *instrumentProfile);
/**
 * Returns the default schedule for a schedule definition.
 *
 * The definition has the form of the trading hours of an instrument profile. The first parenthesized part of it is
 * used. The function fails if the definition contains no parenthesized part.
 *
 * @deprecated Use dxfg_v4_Schedule_getInstanceByDefinition().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] scheduleDefinition The schedule definition, a null-terminated UTF-8 string.
 * @return The schedule (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getInstance(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getInstanceByDefinition()")
dxfg_schedule_t*      dxfg_Schedule_getInstance2(graal_isolatethread_t *thread, const char *scheduleDefinition);
/**
 * Returns the schedule for an instrument profile and a trading venue.
 *
 * The venue is the name that precedes the schedule definition in parentheses in the trading hours of the profile.
 * The function fails if the trading hours contain no schedule for the venue.
 *
 * @deprecated Use dxfg_v4_Schedule_getInstanceByProfileHandleWithVenue().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile whose schedule is requested.
 * @param[in] venue The trading venue, a null-terminated UTF-8 string.
 * @return The schedule (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getInstance(InstrumentProfile, String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getInstanceByProfileHandleWithVenue()")
dxfg_schedule_t*      dxfg_Schedule_getInstance3(graal_isolatethread_t *thread, dxfg_instrument_profile_t *instrumentProfile, const char *venue);
/**
 * Returns the trading venues of an instrument profile.
 *
 * The list is empty (size 0, elements NULL) if the trading hours of the profile contain no venues.
 *
 * @deprecated Use dxfg_v4_Schedule_getTradingVenuesByProfileHandle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile whose venues are requested.
 * @return The list of the venue names (release it with dxfg_CList_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getTradingVenues(InstrumentProfile)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getTradingVenuesByProfileHandle()")
dxfg_string_list*     dxfg_Schedule_getTradingVenues(graal_isolatethread_t *thread, dxfg_instrument_profile_t *instrumentProfile);
/**
 * Downloads the default schedule data and optionally starts the periodic download.
 *
 * The first download is done before the function returns. A failure of the download itself is logged and does not
 * make the function fail.
 *
 * @deprecated Use dxfg_v4_Schedule_downloadDefaults().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] downloadConfig The download config, a null-terminated UTF-8 string, one of:
 * an empty string or NULL (stop the periodic download); a URL (download once from the URL and stop the periodic
 * download); `URL,period` (start the periodic download from the URL); `auto` (start the periodic download from the
 * default location).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `Schedule.downloadDefaults(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_downloadDefaults()")
int32_t               dxfg_Schedule_downloadDefaults(graal_isolatethread_t *thread, const char *downloadConfig);
/**
 * Sets the shared defaults that are used by the schedule instances.
 *
 * @deprecated Use dxfg_v4_Schedule_setDefaults().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] data The content of the default data. Compressed data is detected by its header.
 * @param[in] size The size of `data` in bytes.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `Schedule.setDefaults(byte[])`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_setDefaults()")
int32_t               dxfg_Schedule_setDefaults(graal_isolatethread_t *thread, const char *data, int32_t size);
/**
 * Returns the session of a schedule that contains a time.
 *
 * The function fails if the time is outside the supported date range.
 *
 * @deprecated Use dxfg_v4_Schedule_getSessionByTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @return The session (release it with dxfg_Session_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getSessionByTime(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getSessionByTime()")
dxfg_session_t*       dxfg_Schedule_getSessionByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time);
/**
 * Returns the day of a schedule that contains a time.
 *
 * The function fails if the time is outside the supported date range.
 *
 * @deprecated Use dxfg_v4_Schedule_getDayByTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @return The day (release it with dxfg_Day_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getDayByTime(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getDayByTime()")
dxfg_day_t*           dxfg_Schedule_getDayByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time);
/**
 * Returns the day of a schedule with a day identifier.
 *
 * The function fails if the day identifier is outside the supported date range.
 *
 * @deprecated Use dxfg_v4_Schedule_getDayById().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] dayId The day identifier: the number of the day since 1970-01-01, which has the identifier 0; the
 * previous days have negative identifiers. See dxfg_Day_getDayId().
 * @return The day (release it with dxfg_Day_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getDayById(int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getDayById()")
dxfg_day_t*           dxfg_Schedule_getDayById(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int32_t dayId);
/**
 * Returns the day of a schedule with year, month and day numbers.
 *
 * If the day does not exist, the day with the lowest valid yearMonthDay that is greater than the specified one is
 * returned. The function fails if the value is outside the supported date range.
 *
 * @deprecated Use dxfg_v4_Schedule_getDayByYearMonthDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] yearMonthDay The year, month and day numbers packed as `year * 10000 + month * 100 + day`; for
 * example, 19770928 is September 28, 1977. See dxfg_Day_getYearMonthDay().
 * @return The day (release it with dxfg_Day_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getDayByYearMonthDay(int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getDayByYearMonthDay()")
dxfg_day_t*           dxfg_Schedule_getDayByYearMonthDay(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int32_t yearMonthDay);
/**
 * Returns the session of a schedule that is nearest to a time and is accepted by a filter.
 *
 * The function fails if the time is outside the supported date range or if no session accepted by the filter is found
 * within one year; in both cases it returns NULL and an exception is set. To get NULL without an exception when
 * no session is found, use dxfg_Schedule_findNearestSessionByTime().
 *
 * @deprecated Use dxfg_v4_Schedule_getNearestSessionByTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @param[in] filter The filter that tests the sessions.
 * @return The session (release it with dxfg_Session_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getNearestSessionByTime(long, SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getNearestSessionByTime()")
dxfg_session_t*       dxfg_Schedule_getNearestSessionByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time, dxfg_session_filter_t* filter);
/**
 * Returns the session of a schedule that is nearest to a time and is accepted by a filter, or NULL if there is none.
 *
 * NULL is returned without an exception if no session accepted by the filter is found within one year. The function
 * fails (NULL and an exception) if the time is outside the supported date range.
 *
 * @deprecated Use dxfg_v4_Schedule_findNearestSessionByTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @param[in] filter The filter that tests the sessions.
 * @return The session (release it with dxfg_Session_release()), or NULL if no session is found or on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception: it is NULL if no session is found.
 *
 * Java: `Schedule.findNearestSessionByTime(long, SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_findNearestSessionByTime()")
dxfg_session_t*       dxfg_Schedule_findNearestSessionByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time, dxfg_session_filter_t* filter);
/**
 * Returns the name of a schedule.
 *
 * @deprecated Use dxfg_v4_Schedule_getName().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @return The name, a null-terminated UTF-8 string (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getName()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getName()")
const char*           dxfg_Schedule_getName(graal_isolatethread_t *thread, dxfg_schedule_t* schedule);
/**
 * Returns the display name of the time zone in which a schedule is defined.
 *
 * To get the identifier of the time zone, use dxfg_Schedule_getTimeZone_getID().
 *
 * @deprecated Use dxfg_v4_Schedule_getTimeZone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @return The display name, a null-terminated UTF-8 string (release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getTimeZone()`, `TimeZone.getDisplayName()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getTimeZone()")
const char*           dxfg_Schedule_getTimeZone(graal_isolatethread_t *thread, dxfg_schedule_t* schedule);
/**
 * Returns the identifier of the time zone in which a schedule is defined.
 *
 * @deprecated Use dxfg_v4_Schedule_getTimeZone_getID().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @return The time zone identifier, a null-terminated UTF-8 string (release it with dxfg_String_release()), or NULL
 * on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Schedule.getTimeZone()`, `TimeZone.getID()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getTimeZone_getID()")
const char*           dxfg_Schedule_getTimeZone_getID(graal_isolatethread_t *thread, dxfg_schedule_t* schedule);
/**
 * Returns the definition of a schedule: its name and its properties in parentheses, e.g. `NYSE(...)`. It is the part
 * of the string given to dxfg_Schedule_getInstance2(), or of the trading hours of the instrument profile, that defines
 * this schedule.
 *
 * @deprecated Use dxfg_v4_Schedule_getDefinition().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[out] definition The definition, a null-terminated UTF-8 string (release it with dxfg_String_release()). It is
 * written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `definition` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `Schedule.getDefinition()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getDefinition()")
int32_t               dxfg_Schedule_getDefinition(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, DXFG_OUT char **definition);


/**
 * Returns the day to which a session belongs.
 *
 * @deprecated Use dxfg_v4_Session_getDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return The day (release it with dxfg_Day_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.getDay()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_getDay()")
dxfg_day_t*           dxfg_Session_getDay(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Returns the type of a session.
 *
 * @deprecated Use dxfg_v4_Session_getType().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return The type of the session, a #dxfg_session_type_t value, or -1 on error (it is not a valid session type). Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.getType()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_getType()")
int32_t               dxfg_Session_getType(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Checks whether trading activity is allowed within a session.
 *
 * @deprecated Use dxfg_v4_Session_isTrading().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return 1 if trading activity is allowed within the session, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.isTrading()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_isTrading()")
int32_t               dxfg_Session_isTrading(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Checks whether a session has zero duration.
 *
 * Empty sessions are used for example for indices that post a value once a day. They can be of any type.
 *
 * @deprecated Use dxfg_v4_Session_isEmpty().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return 1 if the session has zero duration, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.isEmpty()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_isEmpty()")
int32_t               dxfg_Session_isEmpty(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Returns the start time of a session (inclusive).
 *
 * For an empty session the start time equals the end time.
 *
 * @deprecated Use dxfg_v4_Session_getStartTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return The start time in UTC milliseconds, or -1 on error. -1 is also a valid time (1969-12-31T23:59:59.999Z), so
 * use dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Session.getStartTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_getStartTime()")
int64_t               dxfg_Session_getStartTime(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Returns the end time of a session (exclusive).
 *
 * For an empty session the end time equals the start time.
 *
 * @deprecated Use dxfg_v4_Session_getEndTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return The end time in UTC milliseconds, or -1 on error. -1 is also a valid time (1969-12-31T23:59:59.999Z), so
 * use dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Session.getEndTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_getEndTime()")
int64_t               dxfg_Session_getEndTime(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Checks whether a time belongs to a session.
 *
 * @deprecated Use dxfg_v4_Session_containsTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[in] time The time to check, in UTC milliseconds.
 * @return 1 if the time belongs to the session, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.containsTime(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_containsTime()")
int32_t               dxfg_Session_containsTime(graal_isolatethread_t *thread, dxfg_session_t* session, int64_t time);
/**
 * Returns the previous session accepted by a filter.
 *
 * The search may cross the day boundary and look up to a year back in time. The function fails if no such session
 * is found; it returns NULL and an exception is set. To get NULL without an exception when no session is found, use
 * dxfg_Session_findPrevSession().
 *
 * @deprecated Use dxfg_v4_Session_getPrevSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @return The nearest previous session accepted by the filter (release it with dxfg_Session_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.getPrevSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_getPrevSession()")
dxfg_session_t*       dxfg_Session_getPrevSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
/**
 * Returns the following session accepted by a filter.
 *
 * The search may cross the day boundary and look up to a year ahead. The function fails if no such session is
 * found; it returns NULL and an exception is set. To get NULL without an exception when no session is found, use
 * dxfg_Session_findNextSession().
 *
 * @deprecated Use dxfg_v4_Session_getNextSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @return The nearest following session accepted by the filter (release it with dxfg_Session_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.getNextSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_getNextSession()")
dxfg_session_t*       dxfg_Session_getNextSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
/**
 * Returns the previous session accepted by a filter, or NULL if there is none.
 *
 * The search may cross the day boundary and look up to a year back in time. NULL is returned without an exception
 * if no such session is found.
 *
 * @deprecated Use dxfg_v4_Session_findPrevSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @return The nearest previous session accepted by the filter (release it with dxfg_Session_release()), or NULL if no
 * session is found or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception: it is NULL if no
 * session is found.
 *
 * Java: `Session.findPrevSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_findPrevSession()")
dxfg_session_t*       dxfg_Session_findPrevSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
/**
 * Returns the following session accepted by a filter, or NULL if there is none.
 *
 * The search may cross the day boundary and look up to a year ahead. NULL is returned without an exception
 * if no such session is found.
 *
 * @deprecated Use dxfg_v4_Session_findNextSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @return The nearest following session accepted by the filter (release it with dxfg_Session_release()), or NULL if
 * no session is found or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception: it is NULL if no
 * session is found.
 *
 * Java: `Session.findNextSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_findNextSession()")
dxfg_session_t*       dxfg_Session_findNextSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
/**
 * Returns the hash code of a session.
 *
 * @deprecated Use dxfg_v4_Session_hashCode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return The hash code, or -1 on error. -1 can also be a valid hash code, so use
 * dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Session.hashCode()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_hashCode()")
int32_t               dxfg_Session_hashCode(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Checks whether two sessions are equal.
 *
 * Sessions are equal if they belong to equal days and have the same type, start time and end time.
 *
 * @deprecated Use dxfg_v4_Session_equals().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[in] otherSession The session to compare with.
 * @return 1 if the sessions are equal, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.equals(Object)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_equals()")
int32_t               dxfg_Session_equals(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_t* otherSession);
/**
 * Returns the string representation of a session.
 *
 * @deprecated Use dxfg_v4_Session_toString().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @return The string representation, a null-terminated UTF-8 string (release it with dxfg_String_release()), or NULL
 * on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Session.toString()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_toString()")
const char*           dxfg_Session_toString(graal_isolatethread_t *thread, dxfg_session_t* session);


/**
 * Returns the schedule to which a day belongs.
 *
 * @deprecated Use dxfg_v4_Day_getSchedule().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The schedule (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getSchedule()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getSchedule()")
dxfg_schedule_t*      dxfg_Day_getSchedule(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the identifier of a day: the number of the day since 1970-01-01, which has the identifier 0. The previous
 * days have negative identifiers.
 *
 * @deprecated Use dxfg_v4_Day_getDayId().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The day identifier, or -1 on error. -1 is also a valid identifier (1969-12-31), so use
 * dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Day.getDayId()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getDayId()")
int32_t               dxfg_Day_getDayId(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the year, month and day numbers of a day packed as `year * 10000 + month * 100 + day`.
 *
 * For example, September 28, 1977 has the value 19770928.
 *
 * @deprecated Use dxfg_v4_Day_getYearMonthDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The packed year, month and day numbers, or -1 on error (it is not a valid value). Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getYearMonthDay()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getYearMonthDay()")
int32_t               dxfg_Day_getYearMonthDay(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the calendar year of a day, for example 1977.
 *
 * @deprecated Use dxfg_v4_Day_getYear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The year, or -1 on error (it is not a valid year). Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 *
 * Java: `Day.getYear()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getYear()")
int32_t               dxfg_Day_getYear(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the calendar month number of a day, from 1 (January) to 12 (December).
 *
 * @deprecated Use dxfg_v4_Day_getMonthOfYear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The month number, or -1 on error (it is not a valid month). Use dxfg_get_and_clear_thread_exception_t() to
 * get the exception.
 *
 * Java: `Day.getMonthOfYear()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getMonthOfYear()")
int32_t               dxfg_Day_getMonthOfYear(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the ordinal number of a day in its month, starting with 1 for the first day of the month.
 *
 * @deprecated Use dxfg_v4_Day_getDayOfMonth().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The day of the month, or -1 on error (it is not a valid value). Use dxfg_get_and_clear_thread_exception_t()
 * to get the exception.
 *
 * Java: `Day.getDayOfMonth()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getDayOfMonth()")
int32_t               dxfg_Day_getDayOfMonth(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the ordinal number of a day in its week, from 1 (Monday) to 7 (Sunday).
 *
 * @deprecated Use dxfg_v4_Day_getDayOfWeek().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The day of the week, or -1 on error (it is not a valid value). Use dxfg_get_and_clear_thread_exception_t()
 * to get the exception.
 *
 * Java: `Day.getDayOfWeek()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getDayOfWeek()")
int32_t               dxfg_Day_getDayOfWeek(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Checks whether a day is an exchange holiday.
 *
 * Usually no trading takes place on a holiday.
 *
 * @deprecated Use dxfg_v4_Day_isHoliday().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return 1 if the day is a holiday, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.isHoliday()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_isHoliday()")
int32_t               dxfg_Day_isHoliday(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Checks whether a day is a short day.
 *
 * Usually trading stops earlier on a short day.
 *
 * @deprecated Use dxfg_v4_Day_isShortDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return 1 if the day is a short day, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.isShortDay()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_isShortDay()")
int32_t               dxfg_Day_isShortDay(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Checks whether trading activity is allowed within a day.
 *
 * A trading day has at least one trading session.
 *
 * @deprecated Use dxfg_v4_Day_isTrading().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return 1 if trading activity is allowed within the day, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.isTrading()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_isTrading()")
int32_t               dxfg_Day_isTrading(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the start time of a day (inclusive).
 *
 * @deprecated Use dxfg_v4_Day_getStartTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The start time in UTC milliseconds, or -1 on error. -1 is also a valid time (1969-12-31T23:59:59.999Z), so
 * use dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Day.getStartTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getStartTime()")
int64_t               dxfg_Day_getStartTime(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the end time of a day (exclusive).
 *
 * @deprecated Use dxfg_v4_Day_getEndTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The end time in UTC milliseconds, or -1 on error. -1 is also a valid time (1969-12-31T23:59:59.999Z), so
 * use dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Day.getEndTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getEndTime()")
int64_t               dxfg_Day_getEndTime(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Checks whether a time belongs to a day.
 *
 * @deprecated Use dxfg_v4_Day_containsTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] time The time to check, in UTC milliseconds.
 * @return 1 if the time belongs to the day, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.containsTime(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_containsTime()")
int32_t               dxfg_Day_containsTime(graal_isolatethread_t *thread, dxfg_day_t* day, int64_t time);
/**
 * Returns the reset time of a day.
 *
 * The daily data is reset on trading days only; the result has no meaning for non-trading days.
 *
 * @deprecated Use dxfg_v4_Day_getResetTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The reset time in UTC milliseconds, or -1 on error. -1 is also a valid time (1969-12-31T23:59:59.999Z), so
 * use dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Day.getResetTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getResetTime()")
int64_t               dxfg_Day_getResetTime(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the sessions that constitute a day, in the order in which they follow one another.
 *
 * @deprecated Use dxfg_v4_Day_getSessions().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The list of the sessions, or NULL on error. Release each session with dxfg_Session_release() and then the
 * list with dxfg_SessionList_wrapper_release(). Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getSessions()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getSessions()")
dxfg_session_list*    dxfg_Day_getSessions(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Returns the session of a day that contains a time.
 *
 * The function fails if no such session is found within the day; it returns NULL and an exception is set.
 *
 * @deprecated Use dxfg_v4_Day_getSessionByTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @return The session (release it with dxfg_Session_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getSessionByTime(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getSessionByTime()")
dxfg_session_t*       dxfg_Day_getSessionByTime(graal_isolatethread_t *thread, dxfg_day_t* day, int64_t time);
/**
 * Returns the first session of a day that is accepted by a filter.
 *
 * The search does not cross the day boundary. The function fails if no such session is found within the day; it
 * returns NULL and an exception is set. To get NULL without an exception when no session is found, use
 * dxfg_Day_findFirstSession().
 *
 * @deprecated Use dxfg_v4_Day_getFirstSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @return The session (release it with dxfg_Session_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getFirstSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getFirstSession()")
dxfg_session_t*       dxfg_Day_getFirstSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
/**
 * Returns the last session of a day that is accepted by a filter.
 *
 * The search does not cross the day boundary. The function fails if no such session is found within the day; it
 * returns NULL and an exception is set. To get NULL without an exception when no session is found, use
 * dxfg_Day_findLastSession().
 *
 * @deprecated Use dxfg_v4_Day_getLastSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @return The session (release it with dxfg_Session_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getLastSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getLastSession()")
dxfg_session_t*       dxfg_Day_getLastSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
/**
 * Returns the first session of a day that is accepted by a filter, or NULL if there is none.
 *
 * The search does not cross the day boundary. NULL is returned without an exception if no such session is found
 * within the day.
 *
 * @deprecated Use dxfg_v4_Day_findFirstSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @return The session (release it with dxfg_Session_release()), or NULL if no session is found or on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception: it is NULL if no session is found.
 *
 * Java: `Day.findFirstSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_findFirstSession()")
dxfg_session_t*       dxfg_Day_findFirstSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
/**
 * Returns the last session of a day that is accepted by a filter, or NULL if there is none.
 *
 * The search does not cross the day boundary. NULL is returned without an exception if no such session is found
 * within the day.
 *
 * @deprecated Use dxfg_v4_Day_findLastSession().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @return The session (release it with dxfg_Session_release()), or NULL if no session is found or on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception: it is NULL if no session is found.
 *
 * Java: `Day.findLastSession(SessionFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_findLastSession()")
dxfg_session_t*       dxfg_Day_findLastSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
/**
 * Returns the previous day accepted by a filter.
 *
 * The search looks up to a year back in time. The function fails if no such day is found; it returns NULL and an
 * exception is set. To get NULL without an exception when no day is found, use dxfg_Day_findPrevDay().
 *
 * @deprecated Use dxfg_v4_Day_getPrevDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @return The nearest previous day accepted by the filter (release it with dxfg_Day_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getPrevDay(DayFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getPrevDay()")
dxfg_day_t*           dxfg_Day_getPrevDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
/**
 * Returns the following day accepted by a filter.
 *
 * The search looks up to a year ahead. The function fails if no such day is found; it returns NULL and an exception
 * is set. To get NULL without an exception when no day is found, use dxfg_Day_findNextDay().
 *
 * @deprecated Use dxfg_v4_Day_getNextDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @return The nearest following day accepted by the filter (release it with dxfg_Day_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.getNextDay(DayFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_getNextDay()")
dxfg_day_t*           dxfg_Day_getNextDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
/**
 * Returns the previous day accepted by a filter, or NULL if there is none.
 *
 * The search looks up to a year back in time. NULL is returned without an exception if no such day is found.
 *
 * @deprecated Use dxfg_v4_Day_findPrevDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @return The nearest previous day accepted by the filter (release it with dxfg_Day_release()), or NULL if no day is
 * found or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception: it is NULL if no day is found.
 *
 * Java: `Day.findPrevDay(DayFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_findPrevDay()")
dxfg_day_t*           dxfg_Day_findPrevDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
/**
 * Returns the following day accepted by a filter, or NULL if there is none.
 *
 * The search looks up to a year ahead. NULL is returned without an exception if no such day is found.
 *
 * @deprecated Use dxfg_v4_Day_findNextDay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @return The nearest following day accepted by the filter (release it with dxfg_Day_release()), or NULL if no day is
 * found or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception: it is NULL if no day is found.
 *
 * Java: `Day.findNextDay(DayFilter)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_findNextDay()")
dxfg_day_t*           dxfg_Day_findNextDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
/**
 * Returns the hash code of a day.
 *
 * @deprecated Use dxfg_v4_Day_hashCode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The hash code, or -1 on error. -1 can also be a valid hash code, so use
 * dxfg_get_and_clear_thread_exception_t() to find out whether an exception was thrown.
 *
 * Java: `Day.hashCode()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_hashCode()")
int32_t               dxfg_Day_hashCode(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Checks whether two days are equal.
 *
 * Days are equal if they belong to equal schedules and have the same day identifier.
 *
 * @deprecated Use dxfg_v4_Day_equals().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] otherDay The day to compare with.
 * @return 1 if the days are equal, 0 if not, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.equals(Object)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_equals()")
int32_t               dxfg_Day_equals(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_t* otherDay);
/**
 * Returns the string representation of a day.
 *
 * @deprecated Use dxfg_v4_Day_toString().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @return The string representation, a null-terminated UTF-8 string (release it with dxfg_String_release()), or NULL
 * on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Day.toString()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_toString()")
const char*           dxfg_Day_toString(graal_isolatethread_t *thread, dxfg_day_t* day);


/**
 * Releases a day handle.
 *
 * @deprecated Use dxfg_v4_Day_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day handle to release. Do not use it afterwards.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Day_release()")
int32_t               dxfg_Day_release(graal_isolatethread_t *thread, dxfg_day_t* day);
/**
 * Releases a session handle.
 *
 * @deprecated Use dxfg_v4_Session_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session handle to release. Do not use it afterwards.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Session_release()")
int32_t               dxfg_Session_release(graal_isolatethread_t *thread, dxfg_session_t* session);
/**
 * Releases the dxfg_session_list structure and its array of elements and keeps the sessions themselves.
 *
 * Release the sessions with dxfg_Session_release() separately; the list must not be used afterwards.
 *
 * @deprecated Use dxfg_v4_session_list_wrapper_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sessions The list to release; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_session_list_wrapper_release()")
int32_t               dxfg_SessionList_wrapper_release(graal_isolatethread_t *thread, dxfg_session_list* sessions);

/**
 * Creates a schedule instance based on a instrument profile.
 *
 * The schedule is defined by the trading hours of the profile. The function fails if the trading hours contain no
 * parenthesized schedule definition.
 *
 * @deprecated Use dxfg_v4_Schedule_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile.
 * @param[out] schedule The new schedule instance.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the object handle.
 *
 * Java: `Schedule.getInstance(InstrumentProfile)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getInstance()")
int32_t dxfg_Schedule_getInstance4(graal_isolatethread_t *thread, dxfg_instrument_profile2_t *instrumentProfile, DXFG_OUT dxfg_schedule_t** schedule);

/**
 * Creates a schedule instance based on a instrument profile and a venue.
 *
 * The venue is the name that precedes the schedule definition in parentheses in the trading hours of the profile.
 * The function fails if the trading hours contain no schedule for the venue.
 *
 * @deprecated Use dxfg_v4_Schedule_getInstanceWithVenue().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile.
 * @param[in] venue The venue, a null-terminated UTF-8 string.
 * @param[out] schedule The new schedule instance.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the object handle.
 *
 * Java: `Schedule.getInstance(InstrumentProfile, String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getInstanceWithVenue()")
int32_t dxfg_Schedule_getInstance5(graal_isolatethread_t *thread, dxfg_instrument_profile2_t *instrumentProfile, const char *venue, DXFG_OUT dxfg_schedule_t** schedule);

/**
 * Returns trading venues by an instrument profile.
 *
 * The list is empty (size 0, elements NULL) if the trading hours of the profile contain no venues.
 *
 * @deprecated Use dxfg_v4_Schedule_getTradingVenues().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile.
 * @param[out] venues Returned list of venues.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_CList_String_release() to free the result.
 *
 * Java: `Schedule.getTradingVenues(InstrumentProfile)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Schedule_getTradingVenues()")
int32_t dxfg_Schedule_getTradingVenues2(graal_isolatethread_t *thread, dxfg_instrument_profile2_t *instrumentProfile, DXFG_OUT dxfg_string_list** venues);

/** @} */ // end of Schedule

#ifdef __cplusplus
}
#endif

#endif // DXFG_SCHEDULE_H




