// SPDX-License-Identifier: MPL-2.0

#ifndef DXFEED_GRAAL_NATIVE_SDK_FEED_H_
#define DXFEED_GRAAL_NATIVE_SDK_FEED_H_

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_events.h"
#include "dxfg_catch_exception.h"
#include "dxfg_javac.h"
#include "graal_isolate.h"

/** @defgroup Feed
 *  @{
 */

/**
 * @brief Forward declarations.
 */
typedef struct dxfg_subscription_t dxfg_subscription_t;
typedef struct dxfg_time_series_subscription_t dxfg_time_series_subscription_t;
typedef struct dxfg_executor_t dxfg_executor_t;

/**
 * @brief The DXFeed.
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html">Javadoc</a>
 */
typedef struct dxfg_feed_t {
    dxfg_java_object_handler handler;
} dxfg_feed_t;

/**
 * @brief The DXFeed.
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/promise/Promise.html">Javadoc</a>
 */
typedef struct dxfg_promise_t {
    dxfg_java_object_handler handler;
} dxfg_promise_t;


typedef struct dxfg_promise_list {
    dxfg_java_object_handler_list list;
} dxfg_promise_list;

typedef struct dxfg_promise_events_t {
    dxfg_promise_t base;
} dxfg_promise_events_t;


typedef struct dxfg_promise_event_t {
    dxfg_promise_t handler;
} dxfg_promise_event_t;

dxfg_feed_t*                      dxfg_DXFeed_getInstance(graal_isolatethread_t *thread);
dxfg_subscription_t*              dxfg_DXFeed_createSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz);
dxfg_subscription_t*              dxfg_DXFeed_createSubscription2(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_list_t *eventClazzes);
dxfg_time_series_subscription_t*  dxfg_DXFeed_createTimeSeriesSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz);
dxfg_time_series_subscription_t*  dxfg_DXFeed_createTimeSeriesSubscription2(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_list_t *eventClazzes);
int32_t                           dxfg_DXFeed_attachSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_subscription_t *sub);
int32_t                           dxfg_DXFeed_detachSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_subscription_t *sub);
int32_t                           dxfg_DXFeed_detachSubscriptionAndClear(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_subscription_t *sub);
dxfg_event_type_t*                dxfg_DXFeed_getLastEventIfSubscribed(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol);
dxfg_event_type_list*             dxfg_DXFeed_getIndexedEventsIfSubscribed(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol, const char *source);
dxfg_event_type_list*             dxfg_DXFeed_getTimeSeriesIfSubscribed(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol, int64_t from_time, int64_t to_time);
// use dxfg_EventType_new to create an empty structure so that java tries to free up memory when replacing subjects
int32_t                           dxfg_DXFeed_getLastEvent(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_type_t *event);
int32_t                           dxfg_DXFeed_getLastEvents(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_type_list *events);
dxfg_promise_event_t*             dxfg_DXFeed_getLastEventPromise(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol);
dxfg_promise_list*                dxfg_DXFeed_getLastEventsPromises(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_list *symbols);
dxfg_promise_events_t*            dxfg_DXFeed_getIndexedEventsPromise(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol, dxfg_indexed_event_source_t* source);
dxfg_promise_events_t*            dxfg_DXFeed_getTimeSeriesPromise(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t clazz, dxfg_symbol_t *symbol, int64_t fromTime, int64_t toTime);

int32_t                           dxfg_DXFeedTimeSeriesSubscription_setFromTime(graal_isolatethread_t *thread, dxfg_time_series_subscription_t *sub, int64_t fromTime);

typedef void (*dxfg_promise_handler_function)(graal_isolatethread_t *thread, dxfg_promise_t *promise, void *user_data);

int32_t               dxfg_Promise_isDone(graal_isolatethread_t *thread, dxfg_promise_t *promise);
int32_t               dxfg_Promise_hasResult(graal_isolatethread_t *thread, dxfg_promise_t *promise);
int32_t               dxfg_Promise_hasException(graal_isolatethread_t *thread, dxfg_promise_t *promise);
int32_t               dxfg_Promise_isCancelled(graal_isolatethread_t *thread, dxfg_promise_t *promise);
dxfg_event_type_t*    dxfg_Promise_EventType_getResult(graal_isolatethread_t *thread, dxfg_promise_event_t *promise);
dxfg_event_type_list* dxfg_Promise_List_EventType_getResult(graal_isolatethread_t *thread, dxfg_promise_events_t *promise);
dxfg_exception_t*     dxfg_Promise_getException(graal_isolatethread_t *thread, dxfg_promise_t *promise);
int32_t               dxfg_Promise_await(graal_isolatethread_t *thread, dxfg_promise_t *promise);
int32_t               dxfg_Promise_await2(graal_isolatethread_t *thread, dxfg_promise_t *promise, int32_t timeoutInMilliseconds);
int32_t               dxfg_Promise_awaitWithoutException(graal_isolatethread_t *thread, dxfg_promise_t *promise, int32_t timeoutInMilliseconds);
int32_t               dxfg_Promise_cancel(graal_isolatethread_t *thread, dxfg_promise_t *promise);
int32_t               dxfg_Promise_List_EventType_complete(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_event_type_list* events);
int32_t               dxfg_Promise_EventType_complete(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_event_type_t* event);
int32_t               dxfg_Promise_completeExceptionally(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_exception_t* exception);
int32_t               dxfg_Promise_whenDone(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_promise_handler_function promise_handler_function, void *user_data);
int32_t               dxfg_Promise_whenDoneAsync(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_promise_handler_function promise_handler_function, void *user_data, dxfg_executor_t* executor);
dxfg_promise_t*       dxfg_Promise_completed(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_java_object_handler *handler);
dxfg_promise_t*       dxfg_Promise_failed(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_exception_t* exception);

dxfg_promise_t*       dxfg_Promises_allOf(graal_isolatethread_t *thread, dxfg_promise_list *promises);


/** @} */ // end of Feed

#ifdef __cplusplus
}
#endif

#endif // DXFEED_GRAAL_NATIVE_SDK_FEED_H_
