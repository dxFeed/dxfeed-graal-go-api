// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_schedule.h and the v4 names of its types (see the group @ref V4_Schedule).
 */

#ifndef DXFG_V4_SCHEDULE_H
#define DXFG_V4_SCHEDULE_H

#include "dxfg_schedule.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_ipf.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Schedule Schedule (v4)
 *  @brief The v4 functions and type names of the group @ref Schedule. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_schedule_t */
typedef dxfg_schedule_t dxfg_v4_schedule_t;

/** @copydoc dxfg_session_type_t */
typedef dxfg_session_type_t dxfg_v4_session_type_t;

/** @copydoc dxfg_session_t */
typedef dxfg_session_t dxfg_v4_session_t;

/**
 * A list of session handles, returned by dxfg_v4_Day_getSessions().
 *
 * Release the list together with its sessions with dxfg_v4_session_list_release(). To keep the sessions, release
 * the list itself with dxfg_v4_session_list_wrapper_release() and then each session with dxfg_v4_Session_release().
 */
typedef dxfg_session_list dxfg_v4_session_list;

/** @copydoc dxfg_day_t */
typedef dxfg_day_t dxfg_v4_day_t;

/** @copydoc dxfg_session_filter_function */
typedef dxfg_session_filter_function dxfg_v4_session_filter_function;

/** @copydoc dxfg_session_filter_t */
typedef dxfg_session_filter_t dxfg_v4_session_filter_t;

/** @copydoc dxfg_session_filter_prepare_t */
typedef dxfg_session_filter_prepare_t dxfg_v4_session_filter_prepare_t;

/** @copydoc dxfg_day_filter_function */
typedef dxfg_day_filter_function dxfg_v4_day_filter_function;

/** @copydoc dxfg_day_filter_t */
typedef dxfg_day_filter_t dxfg_v4_day_filter_t;

/** @copydoc dxfg_day_filter_prepare_t */
typedef dxfg_day_filter_prepare_t dxfg_v4_day_filter_prepare_t;

/**
 * Creates a session filter that calls a custom function to test sessions.
 *
 * The function pointer must stay valid as long as the filter may be used by search functions.
 *
 * Java: `new SessionFilter(SessionType, Boolean)` with `SessionFilter.accept(Session)` overridden to call `filter`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The function that tests a session: it returns 1 to accept the session.
 * @param[out] sessionFilter The pointer to the variable that receives the new session filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_SessionFilter_new(graal_isolatethread_t *thread, dxfg_v4_session_filter_function filter,
                                  DXFG_OUT dxfg_v4_session_filter_t **sessionFilter);

/**
 * Returns a predefined session filter.
 *
 * Java: `SessionFilter.ANY`, `SessionFilter.TRADING`, `SessionFilter.NON_TRADING`, `SessionFilter.NO_TRADING`,
 * `SessionFilter.PRE_MARKET`, `SessionFilter.REGULAR`, `SessionFilter.AFTER_MARKET`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The predefined filter to return.
 * @param[out] sessionFilter The pointer to the variable that receives the session filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_SessionFilter_getInstance(graal_isolatethread_t *thread, dxfg_v4_session_filter_prepare_t filter,
                                          DXFG_OUT dxfg_v4_session_filter_t **sessionFilter);

/**
 * Creates a day filter that calls a custom function to test days.
 *
 * The function pointer must stay valid as long as the filter may be used by search functions.
 *
 * Java: `new DayFilter(int, Boolean, Boolean, Boolean)` with `DayFilter.accept(Day)` overridden to call `filter`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The function that tests a day: it returns 1 to accept the day.
 * @param[out] dayFilter The pointer to the variable that receives the new day filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DayFilter_new(graal_isolatethread_t *thread, dxfg_v4_day_filter_function filter,
                              DXFG_OUT dxfg_v4_day_filter_t **dayFilter);

/**
 * Returns a predefined day filter.
 *
 * Java: `DayFilter.ANY`, `DayFilter.TRADING`, `DayFilter.NON_TRADING`, `DayFilter.HOLIDAY`, `DayFilter.SHORT_DAY`,
 * `DayFilter.MONDAY` ... `DayFilter.SUNDAY`, `DayFilter.WEEK_DAY`, `DayFilter.WEEK_END`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The predefined filter to return.
 * @param[out] dayFilter The pointer to the variable that receives the day filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DayFilter_getInstance(graal_isolatethread_t *thread, dxfg_v4_day_filter_prepare_t filter,
                                      DXFG_OUT dxfg_v4_day_filter_t **dayFilter);

/**
 * The value of a condition of dxfg_v4_DayFilter_newWithUserData() and dxfg_v4_SessionFilter_newWithUserData() that
 * is not checked (Java `null`): for a flag, -1 instead of 0 (false) or 1 (true); for the session type, -1 instead of
 * a value of #dxfg_session_type_t.
 * @since 3.13.0
 */
#define DXFG_V4_FILTER_ANY (-1)

/**
 * A custom day filter function with user data: tests whether a day is an acceptable result of a search, see
 * dxfg_v4_DayFilter_newWithUserData().
 *
 * The function is called by the search function that uses the filter (for example dxfg_v4_Day_findNextDay()), on
 * the thread that called that search function, possibly several times during the call, and only for the days that
 * pass the conditions given to dxfg_v4_DayFilter_newWithUserData().
 *
 * The day handle is valid during the call only: the SDK releases it after the function returns, so the function
 * must not release it.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to test.
 * @param[in] userData The user data given to dxfg_v4_DayFilter_newWithUserData().
 * @return 1 if the day is accepted; any other value rejects it.
 * @since 3.13.0
 */
typedef int32_t (*dxfg_v4_day_filter_user_data_function)(graal_isolatethread_t *thread, dxfg_v4_day_t *day,
                                                         void *userData);

/**
 * A custom session filter function with user data: tests whether a session is an acceptable result of a search, see
 * dxfg_v4_SessionFilter_newWithUserData().
 *
 * The function is called by the search function that uses the filter (for example dxfg_v4_Day_findFirstSession()),
 * on the thread that called that search function, possibly several times during the call, and only for the sessions
 * that pass the conditions given to dxfg_v4_SessionFilter_newWithUserData().
 *
 * The session handle is valid during the call only: the SDK releases it after the function returns, so the function
 * must not release it.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to test.
 * @param[in] userData The user data given to dxfg_v4_SessionFilter_newWithUserData().
 * @return 1 if the session is accepted; any other value rejects it.
 * @since 3.13.0
 */
typedef int32_t (*dxfg_v4_session_filter_user_data_function)(graal_isolatethread_t *thread,
                                                             dxfg_v4_session_t *session, void *userData);

/**
 * Creates a day filter with the conditions of the Java constructor and a custom function with user data.
 *
 * The filter accepts a day if the day passes the conditions and then the function returns 1 for it: the function is
 * not called for a day that the conditions reject. A condition of #DXFG_V4_FILTER_ANY (or a mask of 0) is not
 * checked.
 *
 * The filter is a plain Java object, not a listener: neither the release of its handle nor
 * dxfg_v4_NativeListener_deactivateAll() stops the calls of the function. The function pointer and the user data must
 * stay valid while Java can still call the filter, which may outlive the handle (for example while a search function
 * that got the filter runs on another thread).
 *
 * Java: `new DayFilter(int, Boolean, Boolean, Boolean)` with `DayFilter.accept(Day)` overridden to return
 * `super.accept(day)` and then the result of `filter`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] dayOfWeekMask The bit mask of the accepted days of the week: bit N for the day of the week N (2 for
 * Monday, 4 for Tuesday, ..., 128 for Sunday); 0 accepts any day of the week.
 * @param[in] holiday 1 to accept holidays only, 0 to accept the other days only, #DXFG_V4_FILTER_ANY for any day.
 * @param[in] shortDay 1 to accept short days only, 0 to accept the other days only, #DXFG_V4_FILTER_ANY for any day.
 * @param[in] trading 1 to accept trading days only, 0 to accept the non-trading days only, #DXFG_V4_FILTER_ANY for
 * any day.
 * @param[in] filter The function that tests a day; must not be NULL.
 * @param[in] userData The user data passed to `filter`, may be NULL.
 * @param[out] dayFilter The pointer to the variable that receives the new day filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `holiday`, `shortDay` or `trading` is not -1, 0 or 1, or `filter` is NULL: `IllegalArgumentException`). Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_DayFilter_newWithUserData(graal_isolatethread_t *thread, int32_t dayOfWeekMask, int32_t holiday,
                                          int32_t shortDay, int32_t trading,
                                          dxfg_v4_day_filter_user_data_function filter, void *userData,
                                          DXFG_OUT dxfg_v4_day_filter_t **dayFilter);

/**
 * Creates a session filter with the conditions of the Java constructor and a custom function with user data.
 *
 * The filter accepts a session if the session passes the conditions and then the function returns 1 for it: the
 * function is not called for a session that the conditions reject. A condition of #DXFG_V4_FILTER_ANY is not
 * checked.
 *
 * The filter is a plain Java object, not a listener: neither the release of its handle nor
 * dxfg_v4_NativeListener_deactivateAll() stops the calls of the function. The function pointer and the user data must
 * stay valid while Java can still call the filter, which may outlive the handle (for example while a search function
 * that got the filter runs on another thread).
 *
 * Java: `new SessionFilter(SessionType, Boolean)` with `SessionFilter.accept(Session)` overridden to return
 * `super.accept(session)` and then the result of `filter`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] type The accepted session type, a value of #dxfg_session_type_t, or #DXFG_V4_FILTER_ANY for any type.
 * @param[in] trading 1 to accept trading sessions only, 0 to accept the non-trading sessions only,
 * #DXFG_V4_FILTER_ANY for any session.
 * @param[in] filter The function that tests a session; must not be NULL.
 * @param[in] userData The user data passed to `filter`, may be NULL.
 * @param[out] sessionFilter The pointer to the variable that receives the new session filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `type` is neither -1 nor a value of #dxfg_session_type_t, `trading` is not -1, 0 or 1, or `filter` is NULL:
 * `IllegalArgumentException`). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_SessionFilter_newWithUserData(graal_isolatethread_t *thread, int32_t type, int32_t trading,
                                              dxfg_v4_session_filter_user_data_function filter, void *userData,
                                              DXFG_OUT dxfg_v4_session_filter_t **sessionFilter);

/**
 * Returns the default schedule for an instrument profile.
 *
 * The schedule is defined by the trading hours of the profile. The function fails if the trading hours contain no
 * parenthesized schedule definition.
 *
 * Java: `Schedule.getInstance(InstrumentProfile)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile whose schedule is requested.
 * @param[out] schedule The pointer to the variable that receives the schedule (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getInstanceByProfileHandle(graal_isolatethread_t *thread,
                                                    dxfg_v4_instrument_profile_t *instrumentProfile,
                                                    DXFG_OUT dxfg_v4_schedule_t **schedule);

/**
 * Returns the default schedule for a schedule definition.
 *
 * The definition has the form of the trading hours of an instrument profile. The first parenthesized part of it is
 * used. The function fails if the definition contains no parenthesized part.
 *
 * Java: `Schedule.getInstance(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] scheduleDefinition The schedule definition, a null-terminated UTF-8 string.
 * @param[out] schedule The pointer to the variable that receives the schedule (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getInstanceByDefinition(graal_isolatethread_t *thread, const char *scheduleDefinition,
                                                 DXFG_OUT dxfg_v4_schedule_t **schedule);

/**
 * Returns the schedule for an instrument profile and a trading venue.
 *
 * The venue is the name that precedes the schedule definition in parentheses in the trading hours of the profile.
 * The function fails if the trading hours contain no schedule for the venue.
 *
 * Java: `Schedule.getInstance(InstrumentProfile, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile whose schedule is requested.
 * @param[in] venue The trading venue, a null-terminated UTF-8 string.
 * @param[out] schedule The pointer to the variable that receives the schedule (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getInstanceByProfileHandleWithVenue(graal_isolatethread_t *thread,
                                                             dxfg_v4_instrument_profile_t *instrumentProfile,
                                                             const char *venue, DXFG_OUT dxfg_v4_schedule_t **schedule);

/**
 * Returns the trading venues of an instrument profile.
 *
 * The list is empty (size 0, elements NULL) if the trading hours of the profile contain no venues.
 *
 * Java: `Schedule.getTradingVenues(InstrumentProfile)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile whose venues are requested.
 * @param[out] venues The pointer to the variable that receives the list of the venue names (release it with
 * dxfg_v4_string_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getTradingVenuesByProfileHandle(graal_isolatethread_t *thread,
                                                         dxfg_v4_instrument_profile_t *instrumentProfile,
                                                         DXFG_OUT dxfg_v4_string_list **venues);

/**
 * Downloads the default schedule data and optionally starts the periodic download.
 *
 * The first download is done before the function returns. A failure of the download itself is logged and does not
 * make the function fail.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] downloadConfig The download config, a null-terminated UTF-8 string, one of:
 * an empty string or NULL (stop the periodic download); a URL (download once from the URL and stop the periodic
 * download); `URL,period` (start the periodic download from the URL); `auto` (start the periodic download from the
 * default location).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `Schedule.downloadDefaults(String)`.
 */
int32_t dxfg_v4_Schedule_downloadDefaults(graal_isolatethread_t *thread, const char *downloadConfig);

/**
 * Sets the shared defaults that are used by the schedule instances.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] data The content of the default data. Compressed data is detected by its header.
 * @param[in] size The size of `data` in bytes.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `Schedule.setDefaults(byte[])`.
 */
int32_t dxfg_v4_Schedule_setDefaults(graal_isolatethread_t *thread, const char *data, int32_t size);

/**
 * Returns the session of a schedule that contains a time.
 *
 * The function fails if the time is outside the supported date range.
 *
 * Java: `Schedule.getSessionByTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getSessionByTime(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule, int64_t time,
                                          DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the day of a schedule that contains a time.
 *
 * The function fails if the time is outside the supported date range.
 *
 * Java: `Schedule.getDayByTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @param[out] day The pointer to the variable that receives the day (release it with dxfg_v4_Day_release()); must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getDayByTime(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule, int64_t time,
                                      DXFG_OUT dxfg_v4_day_t **day);

/**
 * Returns the day of a schedule with a day identifier.
 *
 * The function fails if the day identifier is outside the supported date range.
 *
 * Java: `Schedule.getDayById(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] dayId The day identifier: the number of the day since 1970-01-01, which has the identifier 0; the
 * previous days have negative identifiers. See dxfg_v4_Day_getDayId().
 * @param[out] day The pointer to the variable that receives the day (release it with dxfg_v4_Day_release()); must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getDayById(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule, int32_t dayId,
                                    DXFG_OUT dxfg_v4_day_t **day);

/**
 * Returns the day of a schedule with year, month and day numbers.
 *
 * If the day does not exist, the day with the lowest valid yearMonthDay that is greater than the specified one is
 * returned. The function fails if the value is outside the supported date range.
 *
 * Java: `Schedule.getDayByYearMonthDay(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] yearMonthDay The year, month and day numbers packed as `year * 10000 + month * 100 + day`; for
 * example, 19770928 is September 28, 1977. See dxfg_v4_Day_getYearMonthDay().
 * @param[out] day The pointer to the variable that receives the day (release it with dxfg_v4_Day_release()); must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getDayByYearMonthDay(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule,
                                              int32_t yearMonthDay, DXFG_OUT dxfg_v4_day_t **day);

/**
 * Returns the session of a schedule that is nearest to a time and is accepted by a filter.
 *
 * The function fails if the time is outside the supported date range or if no session accepted by the filter is found
 * within one year; in both cases it returns #DXFG_EXECUTE_FAIL. To get NULL instead of a failure when
 * no session is found, use dxfg_v4_Schedule_findNearestSessionByTime().
 *
 * Java: `Schedule.getNearestSessionByTime(long, SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getNearestSessionByTime(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule,
                                                 int64_t time, dxfg_v4_session_filter_t* filter,
                                                 DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the session of a schedule that is nearest to a time and is accepted by a filter, or NULL if there is none.
 *
 * NULL is written and the function succeeds if no session accepted by the filter is found within one year. The function
 * fails (#DXFG_EXECUTE_FAIL) if the time is outside the supported date range.
 *
 * Java: `Schedule.findNearestSessionByTime(long, SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()), or NULL if no session is found; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_findNearestSessionByTime(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule,
                                                  int64_t time, dxfg_v4_session_filter_t* filter,
                                                  DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the name of a schedule.
 *
 * Java: `Schedule.getName()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[out] name The pointer to the variable that receives the name, a null-terminated UTF-8 string (release it with
 * dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getName(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule,
                                 DXFG_OUT const char **name);

/**
 * Returns the display name of the time zone in which a schedule is defined.
 *
 * To get the identifier of the time zone, use dxfg_v4_Schedule_getTimeZone_getID().
 *
 * Java: `Schedule.getTimeZone()`, `TimeZone.getDisplayName()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[out] timeZoneDisplayName The pointer to the variable that receives the display name, a null-terminated UTF-8
 * string (release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getTimeZone(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule,
                                     DXFG_OUT const char **timeZoneDisplayName);

/**
 * Returns the identifier of the time zone in which a schedule is defined.
 *
 * Java: `Schedule.getTimeZone()`, `TimeZone.getID()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[out] timeZoneId The pointer to the variable that receives the time zone identifier, a null-terminated UTF-8
 * string (release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Schedule_getTimeZone_getID(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule,
                                           DXFG_OUT const char **timeZoneId);

/**
 * Returns the definition of a schedule: its name and its properties in parentheses, e.g. `NYSE(...)`. It is the part
 * of the string given to dxfg_v4_Schedule_getInstanceByDefinition(), or of the trading hours of the instrument profile,
 * that defines this schedule.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] schedule The schedule.
 * @param[out] definition The definition, a null-terminated UTF-8 string (release it with dxfg_v4_String_release()). It
 * is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `definition` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `Schedule.getDefinition()`.
 */
int32_t dxfg_v4_Schedule_getDefinition(graal_isolatethread_t *thread, dxfg_v4_schedule_t* schedule,
                                       DXFG_OUT char **definition);

/**
 * Returns the day to which a session belongs.
 *
 * Java: `Session.getDay()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] day The pointer to the variable that receives the day (release it with dxfg_v4_Day_release()); must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_getDay(graal_isolatethread_t *thread, dxfg_v4_session_t* session, DXFG_OUT dxfg_v4_day_t **day);

/**
 * Returns the type of a session.
 *
 * Java: `Session.getType()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] type The pointer to the variable that receives the type of the session, a #dxfg_v4_session_type_t value;
 * must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_getType(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                DXFG_OUT dxfg_v4_session_type_t *type);

/**
 * Checks whether trading activity is allowed within a session.
 *
 * Java: `Session.isTrading()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] trading The pointer to the variable that receives 1 if trading activity is allowed within the session, 0
 * if not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_isTrading(graal_isolatethread_t *thread, dxfg_v4_session_t* session, DXFG_OUT int32_t *trading);

/**
 * Checks whether a session has zero duration.
 *
 * Empty sessions are used for example for indices that post a value once a day. They can be of any type.
 *
 * Java: `Session.isEmpty()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] empty The pointer to the variable that receives 1 if the session has zero duration, 0 if not; must not be
 * NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_isEmpty(graal_isolatethread_t *thread, dxfg_v4_session_t* session, DXFG_OUT int32_t *empty);

/**
 * Returns the start time of a session (inclusive).
 *
 * For an empty session the start time equals the end time.
 *
 * Java: `Session.getStartTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] startTime The pointer to the variable that receives the start time in UTC milliseconds. -1 is also a
 * valid time (1969-12-31T23:59:59.999Z); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_getStartTime(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                     DXFG_OUT int64_t *startTime);

/**
 * Returns the end time of a session (exclusive).
 *
 * For an empty session the end time equals the start time.
 *
 * Java: `Session.getEndTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] endTime The pointer to the variable that receives the end time in UTC milliseconds. -1 is also a valid
 * time (1969-12-31T23:59:59.999Z); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_getEndTime(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                   DXFG_OUT int64_t *endTime);

/**
 * Checks whether a time belongs to a session.
 *
 * Java: `Session.containsTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[in] time The time to check, in UTC milliseconds.
 * @param[out] contains The pointer to the variable that receives 1 if the time belongs to the session, 0 if not; must
 * not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_containsTime(graal_isolatethread_t *thread, dxfg_v4_session_t* session, int64_t time,
                                     DXFG_OUT int32_t *contains);

/**
 * Returns the previous session accepted by a filter.
 *
 * The search may cross the day boundary and look up to a year back in time. The function fails if no such session
 * is found; it returns #DXFG_EXECUTE_FAIL. To get NULL instead of a failure when no session is found, use
 * dxfg_v4_Session_findPrevSession().
 *
 * Java: `Session.getPrevSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] prevSession The pointer to the variable that receives the nearest previous session accepted by the filter
 * (release it with dxfg_v4_Session_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_getPrevSession(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                       dxfg_v4_session_filter_t* filter, DXFG_OUT dxfg_v4_session_t **prevSession);

/**
 * Returns the following session accepted by a filter.
 *
 * The search may cross the day boundary and look up to a year ahead. The function fails if no such session is
 * found; it returns #DXFG_EXECUTE_FAIL. To get NULL instead of a failure when no session is found, use
 * dxfg_v4_Session_findNextSession().
 *
 * Java: `Session.getNextSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] nextSession The pointer to the variable that receives the nearest following session accepted by the
 * filter (release it with dxfg_v4_Session_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_getNextSession(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                       dxfg_v4_session_filter_t* filter, DXFG_OUT dxfg_v4_session_t **nextSession);

/**
 * Returns the previous session accepted by a filter, or NULL if there is none.
 *
 * The search may cross the day boundary and look up to a year back in time. NULL is written and the function succeeds
 * if no such session is found.
 *
 * Java: `Session.findPrevSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] prevSession The pointer to the variable that receives the nearest previous session accepted by the filter
 * (release it with dxfg_v4_Session_release()), or NULL if no session is found; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_findPrevSession(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                        dxfg_v4_session_filter_t* filter, DXFG_OUT dxfg_v4_session_t **prevSession);

/**
 * Returns the following session accepted by a filter, or NULL if there is none.
 *
 * The search may cross the day boundary and look up to a year ahead. NULL is written and the function succeeds
 * if no such session is found.
 *
 * Java: `Session.findNextSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session to start from.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] nextSession The pointer to the variable that receives the nearest following session accepted by the
 * filter (release it with dxfg_v4_Session_release()), or NULL if no session is found; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_findNextSession(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                        dxfg_v4_session_filter_t* filter, DXFG_OUT dxfg_v4_session_t **nextSession);

/**
 * Returns the hash code of a session.
 *
 * Java: `Session.hashCode()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] hashCode The pointer to the variable that receives the hash code. -1 can also be a valid hash code; must
 * not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_hashCode(graal_isolatethread_t *thread, dxfg_v4_session_t* session, DXFG_OUT int32_t *hashCode);

/**
 * Checks whether two sessions are equal.
 *
 * Sessions are equal if they belong to equal days and have the same type, start time and end time.
 *
 * Java: `Session.equals(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[in] otherSession The session to compare with.
 * @param[out] equal The pointer to the variable that receives 1 if the sessions are equal, 0 if not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_equals(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                               dxfg_v4_session_t* otherSession, DXFG_OUT int32_t *equal);

/**
 * Returns the string representation of a session.
 *
 * Java: `Session.toString()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session.
 * @param[out] string The pointer to the variable that receives the string representation, a null-terminated UTF-8
 * string (release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_toString(graal_isolatethread_t *thread, dxfg_v4_session_t* session,
                                 DXFG_OUT const char **string);

/**
 * Returns the schedule to which a day belongs.
 *
 * Java: `Day.getSchedule()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] schedule The pointer to the variable that receives the schedule (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getSchedule(graal_isolatethread_t *thread, dxfg_v4_day_t* day,
                                DXFG_OUT dxfg_v4_schedule_t **schedule);

/**
 * Returns the identifier of a day: the number of the day since 1970-01-01, which has the identifier 0. The previous
 * days have negative identifiers.
 *
 * Java: `Day.getDayId()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] dayId The pointer to the variable that receives the day identifier. -1 is also a valid identifier
 * (1969-12-31); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getDayId(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *dayId);

/**
 * Returns the year, month and day numbers of a day packed as `year * 10000 + month * 100 + day`.
 *
 * For example, September 28, 1977 has the value 19770928.
 *
 * Java: `Day.getYearMonthDay()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] yearMonthDay The pointer to the variable that receives the packed year, month and day numbers; must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getYearMonthDay(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *yearMonthDay);

/**
 * Returns the calendar year of a day, for example 1977.
 *
 * Java: `Day.getYear()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] year The pointer to the variable that receives the year; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getYear(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *year);

/**
 * Returns the calendar month number of a day, from 1 (January) to 12 (December).
 *
 * Java: `Day.getMonthOfYear()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] monthOfYear The pointer to the variable that receives the month number; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getMonthOfYear(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *monthOfYear);

/**
 * Returns the ordinal number of a day in its month, starting with 1 for the first day of the month.
 *
 * Java: `Day.getDayOfMonth()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] dayOfMonth The pointer to the variable that receives the day of the month; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getDayOfMonth(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *dayOfMonth);

/**
 * Returns the ordinal number of a day in its week, from 1 (Monday) to 7 (Sunday).
 *
 * Java: `Day.getDayOfWeek()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] dayOfWeek The pointer to the variable that receives the day of the week; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getDayOfWeek(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *dayOfWeek);

/**
 * Checks whether a day is an exchange holiday.
 *
 * Usually no trading takes place on a holiday.
 *
 * Java: `Day.isHoliday()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] holiday The pointer to the variable that receives 1 if the day is a holiday, 0 if not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_isHoliday(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *holiday);

/**
 * Checks whether a day is a short day.
 *
 * Usually trading stops earlier on a short day.
 *
 * Java: `Day.isShortDay()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] shortDay The pointer to the variable that receives 1 if the day is a short day, 0 if not; must not be
 * NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_isShortDay(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *shortDay);

/**
 * Checks whether trading activity is allowed within a day.
 *
 * A trading day has at least one trading session.
 *
 * Java: `Day.isTrading()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] trading The pointer to the variable that receives 1 if trading activity is allowed within the day, 0 if
 * not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_isTrading(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *trading);

/**
 * Returns the start time of a day (inclusive).
 *
 * Java: `Day.getStartTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] startTime The pointer to the variable that receives the start time in UTC milliseconds. -1 is also a
 * valid time (1969-12-31T23:59:59.999Z); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getStartTime(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int64_t *startTime);

/**
 * Returns the end time of a day (exclusive).
 *
 * Java: `Day.getEndTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] endTime The pointer to the variable that receives the end time in UTC milliseconds. -1 is also a valid
 * time (1969-12-31T23:59:59.999Z); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getEndTime(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int64_t *endTime);

/**
 * Checks whether a time belongs to a day.
 *
 * Java: `Day.containsTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] time The time to check, in UTC milliseconds.
 * @param[out] contains The pointer to the variable that receives 1 if the time belongs to the day, 0 if not; must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_containsTime(graal_isolatethread_t *thread, dxfg_v4_day_t* day, int64_t time,
                                 DXFG_OUT int32_t *contains);

/**
 * Returns the reset time of a day.
 *
 * The daily data is reset on trading days only; the result has no meaning for non-trading days.
 *
 * Java: `Day.getResetTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] resetTime The pointer to the variable that receives the reset time in UTC milliseconds. -1 is also a
 * valid time (1969-12-31T23:59:59.999Z); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getResetTime(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int64_t *resetTime);

/**
 * Returns the sessions that constitute a day, in the order in which they follow one another.
 *
 * Java: `Day.getSessions()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] sessions The pointer to the variable that receives the list of the sessions (release it with
 * dxfg_v4_session_list_release(), or release each session with dxfg_v4_Session_release() and then the list with
 * dxfg_v4_session_list_wrapper_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getSessions(graal_isolatethread_t *thread, dxfg_v4_day_t* day,
                                DXFG_OUT dxfg_v4_session_list **sessions);

/**
 * Returns the session of a day that contains a time.
 *
 * The function fails if no such session is found within the day; it returns #DXFG_EXECUTE_FAIL.
 *
 * Java: `Day.getSessionByTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] time The time to search for, in UTC milliseconds.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getSessionByTime(graal_isolatethread_t *thread, dxfg_v4_day_t* day, int64_t time,
                                     DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the first session of a day that is accepted by a filter.
 *
 * The search does not cross the day boundary. The function fails if no such session is found within the day; it
 * returns #DXFG_EXECUTE_FAIL. To get NULL instead of a failure when no session is found, use
 * dxfg_v4_Day_findFirstSession().
 *
 * Java: `Day.getFirstSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getFirstSession(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_session_filter_t* filter,
                                    DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the last session of a day that is accepted by a filter.
 *
 * The search does not cross the day boundary. The function fails if no such session is found within the day; it
 * returns #DXFG_EXECUTE_FAIL. To get NULL instead of a failure when no session is found, use
 * dxfg_v4_Day_findLastSession().
 *
 * Java: `Day.getLastSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getLastSession(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_session_filter_t* filter,
                                   DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the first session of a day that is accepted by a filter, or NULL if there is none.
 *
 * The search does not cross the day boundary. NULL is written and the function succeeds if no such session is found
 * within the day.
 *
 * Java: `Day.findFirstSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()), or NULL if no session is found; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_findFirstSession(graal_isolatethread_t *thread, dxfg_v4_day_t* day,
                                     dxfg_v4_session_filter_t* filter, DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the last session of a day that is accepted by a filter, or NULL if there is none.
 *
 * The search does not cross the day boundary. NULL is written and the function succeeds if no such session is found
 * within the day.
 *
 * Java: `Day.findLastSession(SessionFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] filter The filter that tests the sessions.
 * @param[out] session The pointer to the variable that receives the session (release it with
 * dxfg_v4_Session_release()), or NULL if no session is found; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_findLastSession(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_session_filter_t* filter,
                                    DXFG_OUT dxfg_v4_session_t **session);

/**
 * Returns the previous day accepted by a filter.
 *
 * The search looks up to a year back in time. The function fails if no such day is found; it returns
 * #DXFG_EXECUTE_FAIL. To get NULL instead of a failure when no day is found, use dxfg_v4_Day_findPrevDay().
 *
 * Java: `Day.getPrevDay(DayFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @param[out] prevDay The pointer to the variable that receives the nearest previous day accepted by the filter
 * (release it with dxfg_v4_Day_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getPrevDay(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_day_filter_t* filter,
                               DXFG_OUT dxfg_v4_day_t **prevDay);

/**
 * Returns the following day accepted by a filter.
 *
 * The search looks up to a year ahead. The function fails if no such day is found; it returns #DXFG_EXECUTE_FAIL. To
 * get NULL instead of a failure when no day is found, use dxfg_v4_Day_findNextDay().
 *
 * Java: `Day.getNextDay(DayFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @param[out] nextDay The pointer to the variable that receives the nearest following day accepted by the filter
 * (release it with dxfg_v4_Day_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_getNextDay(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_day_filter_t* filter,
                               DXFG_OUT dxfg_v4_day_t **nextDay);

/**
 * Returns the previous day accepted by a filter, or NULL if there is none.
 *
 * The search looks up to a year back in time. NULL is written and the function succeeds if no such day is found.
 *
 * Java: `Day.findPrevDay(DayFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @param[out] prevDay The pointer to the variable that receives the nearest previous day accepted by the filter
 * (release it with dxfg_v4_Day_release()), or NULL if no day is found; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_findPrevDay(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_day_filter_t* filter,
                                DXFG_OUT dxfg_v4_day_t **prevDay);

/**
 * Returns the following day accepted by a filter, or NULL if there is none.
 *
 * The search looks up to a year ahead. NULL is written and the function succeeds if no such day is found.
 *
 * Java: `Day.findNextDay(DayFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day to start from.
 * @param[in] filter The filter that tests the days.
 * @param[out] nextDay The pointer to the variable that receives the nearest following day accepted by the filter
 * (release it with dxfg_v4_Day_release()), or NULL if no day is found; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_findNextDay(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_day_filter_t* filter,
                                DXFG_OUT dxfg_v4_day_t **nextDay);

/**
 * Returns the hash code of a day.
 *
 * Java: `Day.hashCode()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] hashCode The pointer to the variable that receives the hash code. -1 can also be a valid hash code; must
 * not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_hashCode(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT int32_t *hashCode);

/**
 * Checks whether two days are equal.
 *
 * Days are equal if they belong to equal schedules and have the same day identifier.
 *
 * Java: `Day.equals(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[in] otherDay The day to compare with.
 * @param[out] equal The pointer to the variable that receives 1 if the days are equal, 0 if not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_equals(graal_isolatethread_t *thread, dxfg_v4_day_t* day, dxfg_v4_day_t* otherDay,
                           DXFG_OUT int32_t *equal);

/**
 * Returns the string representation of a day.
 *
 * Java: `Day.toString()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day.
 * @param[out] string The pointer to the variable that receives the string representation, a null-terminated UTF-8
 * string (release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_toString(graal_isolatethread_t *thread, dxfg_v4_day_t* day, DXFG_OUT const char **string);

/**
 * Releases a day handle.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] day The day handle to release. Do not use it afterwards.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Day_release(graal_isolatethread_t *thread, dxfg_v4_day_t* day);

/**
 * Releases a session handle.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] session The session handle to release. Do not use it afterwards.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Session_release(graal_isolatethread_t *thread, dxfg_v4_session_t* session);

/**
 * Releases the dxfg_v4_session_list structure and its array of elements and keeps the sessions themselves.
 *
 * Release the sessions with dxfg_v4_Session_release() separately; the list must not be used afterwards.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sessions The list to release; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_session_list_wrapper_release(graal_isolatethread_t *thread, dxfg_v4_session_list* sessions);

/**
 * Creates a schedule instance based on a instrument profile.
 *
 * The schedule is defined by the trading hours of the profile. The function fails if the trading hours contain no
 * parenthesized schedule definition.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile.
 * @param[out] schedule The new schedule instance.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the object handle.
 *
 * Java: `Schedule.getInstance(InstrumentProfile)`.
 */
int32_t dxfg_v4_Schedule_getInstance(graal_isolatethread_t *thread,
                                     dxfg_v4_instrument_profile_data_t *instrumentProfile,
                                     DXFG_OUT dxfg_v4_schedule_t** schedule);

/**
 * Creates a schedule instance based on a instrument profile and a venue.
 *
 * The venue is the name that precedes the schedule definition in parentheses in the trading hours of the profile.
 * The function fails if the trading hours contain no schedule for the venue.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile.
 * @param[in] venue The venue, a null-terminated UTF-8 string.
 * @param[out] schedule The new schedule instance.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the object handle.
 *
 * Java: `Schedule.getInstance(InstrumentProfile, String)`.
 */
int32_t dxfg_v4_Schedule_getInstanceWithVenue(graal_isolatethread_t *thread,
                                              dxfg_v4_instrument_profile_data_t *instrumentProfile, const char *venue,
                                              DXFG_OUT dxfg_v4_schedule_t** schedule);

/**
 * Returns trading venues by an instrument profile.
 *
 * The list is empty (size 0, elements NULL) if the trading hours of the profile contain no venues.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrumentProfile The instrument profile.
 * @param[out] venues Returned list of venues.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_string_list_release() to free the result.
 *
 * Java: `Schedule.getTradingVenues(InstrumentProfile)`.
 */
int32_t dxfg_v4_Schedule_getTradingVenues(graal_isolatethread_t *thread,
                                          dxfg_v4_instrument_profile_data_t *instrumentProfile,
                                          DXFG_OUT dxfg_v4_string_list** venues);

/**
 * Releases a list of sessions that a function returned (dxfg_v4_Day_getSessions()): releases every session handle,
 * frees the array of the handles and the list itself. To keep the sessions, release the list with
 * dxfg_v4_session_list_wrapper_release() instead.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sessions The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_session_list_release(graal_isolatethread_t *thread, dxfg_v4_session_list *sessions);

/** @} */ // end of V4_Schedule

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_SCHEDULE_H
