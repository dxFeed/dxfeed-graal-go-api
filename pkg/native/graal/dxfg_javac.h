// SPDX-License-Identifier: MPL-2.0

#ifndef DXFEED_GRAAL_NATIVE_SDK_JAVAC_H_
#define DXFEED_GRAAL_NATIVE_SDK_JAVAC_H_

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "graal_isolate.h"

/** @defgroup Javac
 *  @{
 */

typedef struct dxfg_java_object_handler {
    void *java_object_handle;
} dxfg_java_object_handler;

typedef struct dxfg_list {
    int32_t size;
    void **elements;
} dxfg_list;

typedef struct dxfg_java_object_handler_list {
    int32_t size;
    dxfg_java_object_handler **elements;
} dxfg_java_object_handler_list;

typedef struct dxfg_string_list {
    int32_t size;
    const char **elements;
} dxfg_string_list;

typedef struct dxfg_executor_t {
    dxfg_java_object_handler handler;
} dxfg_executor_t;

typedef struct dxfg_input_stream_t {
    dxfg_java_object_handler handler;
} dxfg_input_stream_t;

typedef struct dxfg_time_period_t {
    dxfg_java_object_handler handler;
} dxfg_time_period_t;

typedef struct dxfg_time_format_t {
    dxfg_java_object_handler handler;
} dxfg_time_format_t;

typedef struct dxfg_time_zone_t {
    dxfg_java_object_handler handler;
} dxfg_time_zone_t;

typedef void (*dxfg_finalize_function)(graal_isolatethread_t *thread, void *user_data);

int32_t                     dxfg_Object_finalize(graal_isolatethread_t *thread, dxfg_java_object_handler* handler, dxfg_finalize_function finalize, void *user_data);

// free the memory occupied by the с data structure and release the reference to the java object
int32_t                     dxfg_JavaObjectHandler_release(graal_isolatethread_t *thread, dxfg_java_object_handler* handler);
const char*                 dxfg_Object_toString(graal_isolatethread_t *thread, dxfg_java_object_handler* object);
int32_t                     dxfg_Object_equals(graal_isolatethread_t *thread, dxfg_java_object_handler* object, dxfg_java_object_handler* other);
int32_t                     dxfg_Object_hashCode(graal_isolatethread_t *thread, dxfg_java_object_handler* object);

// free the memory occupied by the с data structure (list and all elements) and release the reference to the java object for all elements
int32_t                     dxfg_CList_JavaObjectHandler_release(graal_isolatethread_t *thread, dxfg_java_object_handler_list* list);

int32_t                     dxfg_String_release(graal_isolatethread_t *thread, const char* string);
int32_t                     dxfg_CList_String_release(graal_isolatethread_t *thread, dxfg_string_list* string);

// read the "Threads and locks" sections at https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedSubscription.html
dxfg_executor_t*            dxfg_Executors_newFixedThreadPool(graal_isolatethread_t *thread, int nThreads, const char* nameThreads);
dxfg_executor_t*            dxfg_Executors_newScheduledThreadPool(graal_isolatethread_t *thread, int nThreads, const char* nameThreads);
dxfg_executor_t*            dxfg_ExecutorBaseOnConcurrentLinkedQueue_new(graal_isolatethread_t *thread);
int32_t                     dxfg_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks(graal_isolatethread_t *thread, dxfg_executor_t *executor);

dxfg_input_stream_t*        dxfg_ByteArrayInputStream_new(graal_isolatethread_t *thread, const char* bytes, int32_t size);

int32_t                     dxfg_gc(graal_isolatethread_t *thread); // only for testing
dxfg_java_object_handler*   dxfg_throw_exception(graal_isolatethread_t *thread); // only for testing

dxfg_time_zone_t*           dxfg_TimeZone_getTimeZone(graal_isolatethread_t *thread, const char* ID);
dxfg_time_zone_t*           dxfg_TimeZone_getDefault(graal_isolatethread_t *thread);
const char*                 dxfg_TimeZone_getID(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone);
const char*                 dxfg_TimeZone_getDisplayName(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone);
const char*                 dxfg_TimeZone_getDisplayName2(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone, int32_t daylight, int32_t style);
int32_t                     dxfg_TimeZone_getDSTSavings(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone);
int32_t                     dxfg_TimeZone_useDaylightTime(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone);
int32_t                     dxfg_TimeZone_observesDaylightTime(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone);
int32_t                     dxfg_TimeZone_getOffset(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone, int64_t date);
int32_t                     dxfg_TimeZone_getOffset2(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone, int32_t era, int32_t year, int32_t month, int32_t day, int32_t dayOfWeek, int32_t milliseconds);
int32_t                     dxfg_TimeZone_getRawOffset(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone);
int32_t                     dxfg_TimeZone_hasSameRules(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone, dxfg_time_zone_t* other);
int32_t                     dxfg_TimeZone_inDaylightTime(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone, int64_t date);
int32_t                     dxfg_TimeZone_setID(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone, const char* ID);
int32_t                     dxfg_TimeZone_setRawOffset(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone, int32_t offsetMillis);
dxfg_time_period_t*         dxfg_TimePeriod_ZERO(graal_isolatethread_t *thread);
dxfg_time_period_t*         dxfg_TimePeriod_UNLIMITED(graal_isolatethread_t *thread);
dxfg_time_period_t*         dxfg_TimePeriod_valueOf(graal_isolatethread_t *thread, int64_t value);
dxfg_time_period_t*         dxfg_TimePeriod_valueOf2(graal_isolatethread_t *thread, const char* value);
int64_t                     dxfg_TimePeriod_getTime(graal_isolatethread_t *thread, dxfg_time_period_t* timePeriod);
int32_t                     dxfg_TimePeriod_getSeconds(graal_isolatethread_t *thread, dxfg_time_period_t* timePeriod);
int64_t                     dxfg_TimePeriod_getNanos(graal_isolatethread_t *thread, dxfg_time_period_t* timePeriod);
dxfg_time_format_t*         dxfg_TimeFormat_DEFAULT(graal_isolatethread_t *thread);
dxfg_time_format_t*         dxfg_TimeFormat_GMT(graal_isolatethread_t *thread);
dxfg_time_format_t*         dxfg_TimeFormat_getInstance(graal_isolatethread_t *thread, dxfg_time_zone_t* timeZone);
dxfg_time_format_t*         dxfg_TimeFormat_withTimeZone(graal_isolatethread_t *thread, dxfg_time_format_t* timeFormat);
dxfg_time_format_t*         dxfg_TimeFormat_withMillis(graal_isolatethread_t *thread, dxfg_time_format_t* timeFormat);
dxfg_time_format_t*         dxfg_TimeFormat_asFullIso(graal_isolatethread_t *thread, dxfg_time_format_t* timeFormat);
int64_t                     dxfg_TimeFormat_parse(graal_isolatethread_t *thread, dxfg_time_format_t* timeFormat, const char* value);
const char*                 dxfg_TimeFormat_format(graal_isolatethread_t *thread, dxfg_time_format_t* timeFormat, int64_t value);
dxfg_time_zone_t*           dxfg_TimeFormat_getTimeZone(graal_isolatethread_t *thread, dxfg_time_format_t* timeFormat);

/** @} */ // end of Javac

#ifdef __cplusplus
}
#endif

#endif // DXFEED_GRAAL_NATIVE_SDK_JAVAC_H_




