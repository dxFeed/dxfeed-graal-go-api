// SPDX-License-Identifier: MPL-2.0

#ifndef DXFEED_GRAAL_NATIVE_SDK_SCHEDULE_H_
#define DXFEED_GRAAL_NATIVE_SDK_SCHEDULE_H_

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "graal_isolate.h"
#include "dxfg_ipf.h"

/** @defgroup Schedule
 *  @{
 */

typedef struct dxfg_schedule_t {
    dxfg_java_object_handler handler;
} dxfg_schedule_t;

typedef enum dxfg_session_type_t {
  DXFG_SESSION_TYPE_NO_TRADING = 0,
  DXFG_SESSION_TYPE_PRE_MARKET,
  DXFG_SESSION_TYPE_REGULAR,
  DXFG_SESSION_TYPE_AFTER_MARKET
} dxfg_session_type_t;

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/Session.html">Javadoc</a>
 */
typedef struct dxfg_session_t {
    dxfg_java_object_handler handler;
} dxfg_session_t;

typedef struct dxfg_session_list {
    int32_t size;
    dxfg_session_t **elements;
} dxfg_session_list;

typedef struct dxfg_day_t {
    dxfg_java_object_handler handler;
} dxfg_day_t;

typedef int32_t (*dxfg_session_filter_function)(graal_isolatethread_t *thread, dxfg_session_t *session);
typedef struct dxfg_session_filter_t {
    dxfg_java_object_handler handler;
} dxfg_session_filter_t;
typedef enum dxfg_session_filter_prepare_t {
    DXFG_SESSION_FILTER_ANY = 0,
    DXFG_SESSION_FILTER_TRADING,
    DXFG_SESSION_FILTER_NON_TRADING,
    DXFG_SESSION_FILTER_NO_TRADING,
    DXFG_SESSION_FILTER_PRE_MARKET,
    DXFG_SESSION_FILTER_REGULAR,
    DXFG_SESSION_FILTER_AFTER_MARKET
} dxfg_session_filter_prepare_t;
dxfg_session_filter_t*    dxfg_SessionFilter_new(graal_isolatethread_t *thread, dxfg_session_filter_function filter);
dxfg_session_filter_t*    dxfg_SessionFilter_getInstance(graal_isolatethread_t *thread, dxfg_session_filter_prepare_t filter);

typedef int32_t (*dxfg_day_filter_function)(graal_isolatethread_t *thread, dxfg_day_t *day);
typedef struct dxfg_day_filter_t {
    dxfg_java_object_handler handler;
} dxfg_day_filter_t;
typedef enum dxfg_day_filter_prepare_t {
    DXFG_DAY_FILTER_ANY = 0,
    DXFG_DAY_FILTER_TRADING,
    DXFG_DAY_FILTER_NON_TRADING,
    DXFG_DAY_FILTER_HOLIDAY,
    DXFG_DAY_FILTER_SHORT_DAY,
    DXFG_DAY_FILTER_MONDAY,
    DXFG_DAY_FILTER_TUESDAY,
    DXFG_DAY_FILTER_WEDNESDAY,
    DXFG_DAY_FILTER_THURSDAY,
    DXFG_DAY_FILTER_FRIDAY,
    DXFG_DAY_FILTER_SATURDAY,
    DXFG_DAY_FILTER_SUNDAY,
    DXFG_DAY_FILTER_WEEK_DAY,
    DXFG_DAY_FILTER_WEEK_END
} dxfg_day_filter_prepare_t;
dxfg_day_filter_t*    dxfg_DayFilter_new(graal_isolatethread_t *thread, dxfg_day_filter_function filter);
dxfg_day_filter_t*    dxfg_DayFilter_getInstance(graal_isolatethread_t *thread, dxfg_day_filter_prepare_t filter);


dxfg_schedule_t*      dxfg_Schedule_getInstance(graal_isolatethread_t *thread, dxfg_instrument_profile_t *instrumentProfile);
dxfg_schedule_t*      dxfg_Schedule_getInstance2(graal_isolatethread_t *thread, const char *scheduleDefinition);
dxfg_schedule_t*      dxfg_Schedule_getInstance3(graal_isolatethread_t *thread, dxfg_instrument_profile_t *instrumentProfile, const char *venue);
dxfg_string_list*     dxfg_Schedule_getTradingVenues(graal_isolatethread_t *thread, dxfg_instrument_profile_t *instrumentProfile);
int32_t               dxfg_Schedule_downloadDefaults(graal_isolatethread_t *thread, const char *downloadConfig);
int32_t               dxfg_Schedule_setDefaults(graal_isolatethread_t *thread, const char *data, int32_t size);
dxfg_session_t*       dxfg_Schedule_getSessionByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time);
dxfg_day_t*           dxfg_Schedule_getDayByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time);
dxfg_day_t*           dxfg_Schedule_getDayById(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int32_t dayId);
dxfg_day_t*           dxfg_Schedule_getDayByYearMonthDay(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int32_t yearMonthDay);
dxfg_session_t*       dxfg_Schedule_getNearestSessionByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time, dxfg_session_filter_t* filter);
dxfg_session_t*       dxfg_Schedule_findNearestSessionByTime(graal_isolatethread_t *thread, dxfg_schedule_t* schedule, int64_t time, dxfg_session_filter_t* filter);
const char*           dxfg_Schedule_getName(graal_isolatethread_t *thread, dxfg_schedule_t* schedule);
const char*           dxfg_Schedule_getTimeZone(graal_isolatethread_t *thread, dxfg_schedule_t* schedule);
const char*           dxfg_Schedule_getTimeZone_getID(graal_isolatethread_t *thread, dxfg_schedule_t* schedule);


dxfg_day_t*           dxfg_Session_getDay(graal_isolatethread_t *thread, dxfg_session_t* session);
int32_t               dxfg_Session_getType(graal_isolatethread_t *thread, dxfg_session_t* session);
int32_t               dxfg_Session_isTrading(graal_isolatethread_t *thread, dxfg_session_t* session);
int32_t               dxfg_Session_isEmpty(graal_isolatethread_t *thread, dxfg_session_t* session);
int64_t               dxfg_Session_getStartTime(graal_isolatethread_t *thread, dxfg_session_t* session);
int64_t               dxfg_Session_getEndTime(graal_isolatethread_t *thread, dxfg_session_t* session);
int32_t               dxfg_Session_containsTime(graal_isolatethread_t *thread, dxfg_session_t* session, int64_t time);
dxfg_session_t*       dxfg_Session_getPrevSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
dxfg_session_t*       dxfg_Session_getNextSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
dxfg_session_t*       dxfg_Session_findPrevSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
dxfg_session_t*       dxfg_Session_findNextSession(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_filter_t* filter);
int32_t               dxfg_Session_hashCode(graal_isolatethread_t *thread, dxfg_session_t* session);
int32_t               dxfg_Session_equals(graal_isolatethread_t *thread, dxfg_session_t* session, dxfg_session_t* otherSession);
const char*           dxfg_Session_toString(graal_isolatethread_t *thread, dxfg_session_t* session);


dxfg_schedule_t*      dxfg_Day_getSchedule(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_getDayId(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_getYearMonthDay(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_getYear(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_getMonthOfYear(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_getDayOfMonth(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_getDayOfWeek(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_isHoliday(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_isShortDay(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_isTrading(graal_isolatethread_t *thread, dxfg_day_t* day);
int64_t               dxfg_Day_getStartTime(graal_isolatethread_t *thread, dxfg_day_t* day);
int64_t               dxfg_Day_getEndTime(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_containsTime(graal_isolatethread_t *thread, dxfg_day_t* day, int64_t time);
int64_t               dxfg_Day_getResetTime(graal_isolatethread_t *thread, dxfg_day_t* day);
dxfg_session_list*    dxfg_Day_getSessions(graal_isolatethread_t *thread, dxfg_day_t* day);
dxfg_session_t*       dxfg_Day_getSessionByTime(graal_isolatethread_t *thread, dxfg_day_t* day, int64_t time);
dxfg_session_t*       dxfg_Day_getFirstSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
dxfg_session_t*       dxfg_Day_getLastSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
dxfg_session_t*       dxfg_Day_findFirstSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
dxfg_session_t*       dxfg_Day_findLastSession(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_session_filter_t* filter);
dxfg_day_t*           dxfg_Day_getPrevDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
dxfg_day_t*           dxfg_Day_getNextDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
dxfg_day_t*           dxfg_Day_findPrevDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
dxfg_day_t*           dxfg_Day_findNextDay(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_filter_t* filter);
int32_t               dxfg_Day_hashCode(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Day_equals(graal_isolatethread_t *thread, dxfg_day_t* day, dxfg_day_t* otherDay);
const char*           dxfg_Day_toString(graal_isolatethread_t *thread, dxfg_day_t* day);


int32_t               dxfg_Day_release(graal_isolatethread_t *thread, dxfg_day_t* day);
int32_t               dxfg_Session_release(graal_isolatethread_t *thread, dxfg_session_t* session);

/** @} */ // end of Schedule

#ifdef __cplusplus
}
#endif

#endif // DXFEED_GRAAL_NATIVE_SDK_SCHEDULE_H_




