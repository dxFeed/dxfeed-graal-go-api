// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Feed functions and types declarations
 */

#ifndef DXFG_FEED_H
#define DXFG_FEED_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_events.h"
#include "dxfg_catch_exception.h"
#include "dxfg_javac.h"
#include "graal_isolate.h"

/**
 *  @defgroup Feed Feed
 *  @brief The feed (Java `DXFeed`), its subscriptions, the requests of events, and the promises that the requests
 *  return.
 *
 *  The ownership rules of the handles, the strings and the lists are described in the group @ref Javac. The events that
 *  a function returns are released with dxfg_EventType_release(), the lists of events with
 *  dxfg_CList_EventType_release().
 *
 *  In this header a function that returns an `int32_t` status returns #DXFG_EXECUTE_SUCCESSFULLY (0) or
 *  #DXFG_EXECUTE_FAIL (-1) on error, and a function that returns a pointer returns NULL on error, except
 *  dxfg_Promise_getException() that returns `(dxfg_exception_t *)-1`. The functions that tell the state of a promise
 *  return 1, 0 or -1 on error.
 *
 *  A promise is the result of a request that completes later. The functions that start a request return a handle of the
 *  promise: wait for it (dxfg_Promise_await2()) or add a handler (dxfg_Promise_whenDone()), then read the result
 *  (dxfg_Promise_EventType_getResult(), dxfg_Promise_List_EventType_getResult()) and release the handle.
 *  @{
 */

/// @cond
// Forward declarations: the types are described in dxfg_subscription.h and dxfg_javac.h.
typedef struct dxfg_subscription_t dxfg_subscription_t;
typedef struct dxfg_time_series_subscription_t dxfg_time_series_subscription_t;
typedef struct dxfg_executor_t dxfg_executor_t;
/// @endcond

/**
 * @brief A handle of a Java DXFeed object; release it with dxfg_JavaObjectHandler_release().
 *
 * Get the feed of the default endpoint with dxfg_DXFeed_getInstance() or the feed of an endpoint with
 * dxfg_DXEndpoint_getFeed(). Releasing the handle does not close the endpoint of the feed.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html">Javadoc</a>
 */
typedef struct dxfg_feed_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_feed_t;

/**
 * @brief A handle of a Java Promise object; release it with dxfg_JavaObjectHandler_release().
 *
 * A promise is the result of a computation that completes later, normally or exceptionally. The functions of the feed
 * return the more specific dxfg_promise_event_t and dxfg_promise_events_t, which contain this structure.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/promise/Promise.html">Javadoc</a>
 */
typedef struct dxfg_promise_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_promise_t;


/**
 * @brief A list of promise handles.
 *
 * A list that a function returns is allocated by the SDK: release it with dxfg_v4_promise_list_release(), or with
 * dxfg_CList_JavaObjectHandler_release() passing the field `list`. That releases every promise handle, the array of the
 * handles and the list itself.
 */
typedef struct dxfg_promise_list {
    /** The list of the promise handles: the number of the handles and the array of the pointers to them. */
    dxfg_java_object_handler_list list;
} dxfg_promise_list;

/**
 * @brief A handle of a Java Promise object whose result is a list of events; release it with
 * dxfg_JavaObjectHandler_release().
 *
 * It is returned by dxfg_DXFeed_getIndexedEventsPromise() and dxfg_DXFeed_getTimeSeriesPromise(). Read the result with
 * dxfg_Promise_List_EventType_getResult().
 */
typedef struct dxfg_promise_events_t {
    /** The base structure: the promise. Pass `&promise->base` to the functions that take a dxfg_promise_t. */
    dxfg_promise_t base;
} dxfg_promise_events_t;


/**
 * @brief A handle of a Java Promise object whose result is an event; release it with dxfg_JavaObjectHandler_release().
 *
 * It is returned by dxfg_DXFeed_getLastEventPromise(). Read the result with dxfg_Promise_EventType_getResult().
 */
typedef struct dxfg_promise_event_t {
    /** The promise (the field is named `handler`, sic): pass `&promise->handler` to the functions that take
     * a dxfg_promise_t. */
    dxfg_promise_t handler;
} dxfg_promise_event_t;

/**
 * @brief Returns the feed of the default endpoint of the process.
 *
 * Applies the SDK defaults of the QD system properties (a property that the user has set keeps its value) and registers
 * the default endpoint, so that dxfg_system_close_all_and_await_termination() closes it. Every call returns a new
 * handle of the same feed.
 *
 * Java: `DXFeed.getInstance()`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The feed (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getInstance()")
dxfg_feed_t*                      dxfg_DXFeed_getInstance(graal_isolatethread_t *thread);

/**
 * @brief Creates a subscription for the events of one class and attaches it to the feed.
 *
 * The subscription has no symbols and no listeners yet: set them with the dxfg_DXFeedSubscription_* functions. The call
 * is dxfg_DXFeedSubscription_new() followed by dxfg_DXFeed_attachSubscription().
 *
 * Java: `DXFeed.createSubscription(Class)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_createSubscription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events to subscribe to.
 * @return The subscription (release its handle with dxfg_JavaObjectHandler_release(); that neither closes nor detaches
 * the subscription), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_createSubscription()")
dxfg_subscription_t*              dxfg_DXFeed_createSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz);

/**
 * @brief Creates a subscription for the events of several classes and attaches it to the feed.
 *
 * See dxfg_DXFeed_createSubscription(). The list is only read: the caller keeps its ownership. An empty list is an
 * error.
 *
 * Java: `DXFeed.createSubscription(Class...)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_createSubscriptionWithEventTypes().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazzes The classes of the events to subscribe to.
 * @return The subscription (release its handle with dxfg_JavaObjectHandler_release(); that neither closes nor detaches
 * the subscription), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_createSubscriptionWithEventTypes()")
dxfg_subscription_t*              dxfg_DXFeed_createSubscription2(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_list_t *eventClazzes);

/**
 * @brief Creates a time series subscription for the events of one class and attaches it to the feed.
 *
 * Use the classes of the time series events (see dxfg_event_clazz_t). The subscription has no symbols and no listeners
 * yet. Its `fromTime` is initially INT64_MAX, which means that no events are received until
 * dxfg_DXFeedTimeSeriesSubscription_setFromTime() is called.
 *
 * The result contains dxfg_subscription_t as its first field `sub`: pass `&result->sub` (or the pointer cast) to the
 * functions that take a dxfg_subscription_t.
 *
 * Java: `DXFeed.createTimeSeriesSubscription(Class)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_createTimeSeriesSubscription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the time series events to subscribe to.
 * @return The subscription (release its handle with dxfg_JavaObjectHandler_release(); that neither closes nor detaches
 * the subscription), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_createTimeSeriesSubscription()")
dxfg_time_series_subscription_t*  dxfg_DXFeed_createTimeSeriesSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz);

/**
 * @brief Creates a time series subscription for the events of several classes and attaches it to the feed.
 *
 * See dxfg_DXFeed_createTimeSeriesSubscription(). The list is only read: the caller keeps its ownership. An empty list
 * is an error.
 *
 * Java: `DXFeed.createTimeSeriesSubscription(Class...)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_createTimeSeriesSubscriptionWithEventTypes().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazzes The classes of the time series events to subscribe to.
 * @return The subscription (release its handle with dxfg_JavaObjectHandler_release(); that neither closes nor detaches
 * the subscription), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_createTimeSeriesSubscriptionWithEventTypes()")
dxfg_time_series_subscription_t*  dxfg_DXFeed_createTimeSeriesSubscription2(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_list_t *eventClazzes);

/**
 * @brief Attaches the subscription to the feed.
 *
 * From now on the feed publishes the data to the subscription. The call does nothing if the subscription is already
 * attached to this feed. A NULL feed or subscription is an error.
 *
 * Java: `DXFeed.attachSubscription(DXFeedSubscription)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_attachSubscription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] sub The subscription, created by dxfg_DXFeedSubscription_new() or dxfg_DXFeedSubscription_new2() (or any
 * other subscription handle).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_attachSubscription()")
int32_t                           dxfg_DXFeed_attachSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_subscription_t *sub);

/**
 * @brief Detaches the subscription from the feed.
 *
 * The call does nothing if the subscription is not attached to this feed. The data that the subscription has already
 * received stays in it. A NULL feed or subscription is an error.
 *
 * Java: `DXFeed.detachSubscription(DXFeedSubscription)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_detachSubscription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] sub The subscription.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_detachSubscription()")
int32_t                           dxfg_DXFeed_detachSubscription(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_subscription_t *sub);

/**
 * @brief Detaches the subscription from the feed and clears the data that was delivered to the subscription.
 *
 * The data is cleared by publishing empty events to the subscription. The call does nothing if the subscription is not
 * attached to this feed. A NULL feed or subscription is an error.
 *
 * Java: `DXFeed.detachSubscriptionAndClear(DXFeedSubscription)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_detachSubscriptionAndClear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] sub The subscription.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_detachSubscriptionAndClear()")
int32_t                           dxfg_DXFeed_detachSubscriptionAndClear(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_subscription_t *sub);

/**
 * @brief Returns the last event for the event class and symbol if there is a subscription for them.
 *
 * Works only for the lasting events (Quote, Profile, etc.). The function does not make remote calls, it reads the local
 * cache of the feed. The events are in the cache only if an attached subscription is subscribed to the event class and
 * symbol (a wildcard subscription does not count; a subscription can also be defined by the endpoint properties). Use
 * dxfg_DXFeed_getLastEventPromise() to request an event without a subscription.
 *
 * If there is a subscription, but the event has not arrived yet, the event has only its symbol set and the default
 * values of the other fields. If there is no subscription, the function returns NULL: it is not an error, so
 * dxfg_get_and_clear_thread_exception_t() returns NULL then. With the #DXFG_ENDPOINT_ROLE_STREAM_FEED role the function
 * always returns NULL.
 *
 * Java: `DXFeed.getLastEventIfSubscribed(Class, Object)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getLastEventIfSubscribed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the event, a lasting one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership.
 * @return The new event (release it with dxfg_EventType_release()), or NULL if there is no subscription for the event
 * class and symbol, or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getLastEventIfSubscribed()")
dxfg_event_type_t*                dxfg_DXFeed_getLastEventIfSubscribed(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol);

/**
 * @brief Returns the indexed events for the event class, symbol and source if there is a subscription for them.
 *
 * Works only for the indexed events. The function does not make remote calls, it reads the local cache of the feed. The
 * events are in the cache only if an attached subscription is subscribed to the event class, symbol and source (a
 * subscription can also be defined by the endpoint properties). The events are ordered by index, their event flags are
 * always 0. The list is not a consistent snapshot: use dxfg_DXFeed_getIndexedEventsPromise() for it. With the
 * #DXFG_ENDPOINT_ROLE_STREAM_FEED role there are never events.
 *
 * If there is no subscription (Java returns null) or the events have not arrived yet, the result is an empty list
 * (`size` is 0, `elements` is NULL); the two cases cannot be told apart by the result. The empty list must be released
 * as well.
 *
 * Java: `DXFeed.getIndexedEventsIfSubscribed(Class, Object, IndexedEventSource)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getIndexedEventsIfSubscribed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, an indexed one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership.
 * @param[in] source The name of the order source (OrderSource.valueOf(String)), a NUL-terminated UTF-8 string: a
 * predefined name such as `DEFAULT`, or at most 4 alphanumeric characters. A malformed or NULL name is an error. The
 * string is only read.
 * @return The new list of new events (release it with dxfg_CList_EventType_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getIndexedEventsIfSubscribed()")
dxfg_event_type_list*             dxfg_DXFeed_getIndexedEventsIfSubscribed(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol, const char *source);

/**
 * @brief Returns the time series events for the event class, symbol and time range if there is a subscription for them.
 *
 * Works only for the time series events. The function does not make remote calls, it reads the local cache of the feed.
 * The events are in the cache only if an attached subscription is subscribed to the event class, symbol and time (a
 * subscription can also be defined by the endpoint properties). The events are ordered by time, their event flags are
 * always 0. The list is not a consistent snapshot: use dxfg_DXFeed_getTimeSeriesPromise() for it. With the
 * #DXFG_ENDPOINT_ROLE_STREAM_FEED role there are never events.
 *
 * If there is no subscription (Java returns null) or the events have not arrived yet, the result is an empty list
 * (`size` is 0, `elements` is NULL); the two cases cannot be told apart by the result. The empty list must be released
 * as well.
 *
 * Java: `DXFeed.getTimeSeriesIfSubscribed(Class, Object, long, long)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getTimeSeriesIfSubscribed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, a time series one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership. It must not be a
 * dxfg_time_series_subscription_symbol_t.
 * @param[in] from_time The time (inclusive) to return the events from, in milliseconds since 1970-01-01 00:00:00 UTC
 * (see the time of a time series event).
 * @param[in] to_time The time (inclusive) to return the events to, in the same unit; INT64_MAX means no upper limit.
 * @return The new list of new events (release it with dxfg_CList_EventType_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getTimeSeriesIfSubscribed()")
dxfg_event_type_list*             dxfg_DXFeed_getTimeSeriesIfSubscribed(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol, int64_t from_time, int64_t to_time);

/**
 * @brief Fills the given event with the last event for its type and symbol.
 *
 * See dxfg_DXFeed_getLastEvent2() for what the last event is and when it is available.
 *
 * @deprecated Use dxfg_v4_DXFeed_getLastEvent(). This function writes the result into the given event: it frees the
 * strings of the event with the SDK's allocator and writes the strings allocated by the SDK there. So the event must be
 * created by dxfg_EventType_new() (or come from the SDK) and be freed by dxfg_EventType_release(); an event allocated
 * by the caller leaks the SDK's strings or corrupts memory. When the last event is not available, the event keeps the
 * values it had (as in Java), but the event created by dxfg_EventType_new() has only its symbol, so the values of the
 * caller's event are lost.
 *
 * Java: `DXFeed.getLastEvent(LastingEvent)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in,out] event The event created by dxfg_EventType_new(): its type and symbol, and the result. It must not be
 * NULL and must be of a lasting type, otherwise the function fails.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getLastEvent(), it does not write into the given event")
int32_t                           dxfg_DXFeed_getLastEvent(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_type_t *event);

/**
 * @brief The bulk version of dxfg_DXFeed_getLastEvent(): fills every event of the list.
 *
 * @deprecated Use dxfg_v4_DXFeed_getLastEvents(), see dxfg_DXFeed_getLastEvent() for the reasons. Every event of the
 * list must be created by dxfg_EventType_new(). The list must not be NULL.
 *
 * Java: `DXFeed.getLastEvent(LastingEvent)` for every event of the list.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in,out] events The events created by dxfg_EventType_new(), the results are written into them. An event of a
 * type that is not lasting makes the function fail; the events before it are already filled.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getLastEvents(), it does not write into the given events")
int32_t                           dxfg_DXFeed_getLastEvents(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_type_list *events);

/**
 * @brief Gets the last event for the type and symbol of the given event.
 *
 * The result is a new event with the values of the last event, or a copy of the given event when the last event is not
 * available (no subscription, the data have not arrived yet, etc.), since Java leaves the event unchanged then.
 *
 * The given event is only read, as by dxfg_DXPublisher_publishEvents(): the caller allocates it (with any allocator),
 * fills it and frees it; the SDK copies it during the call. The event must be of a lasting type (Quote, Profile, etc.).
 *
 * The function does not make remote calls, it reads the local cache of the feed. The events are in the cache only if an
 * attached subscription is subscribed to the event type and symbol (a wildcard subscription does not count). With the
 * #DXFG_ENDPOINT_ROLE_STREAM_FEED role the last event is never available.
 *
 * Java: `DXFeed.getLastEvent(LastingEvent)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getLastEvent().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] event The event: its type and symbol, and the values to return when the last event is not available.
 * @param[out] lastEvent The pointer to the new event (NULL on error). Free the event with dxfg_EventType_release(). The
 * pointer itself must not be NULL (an error then).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the event is NULL or of a type that is not lasting). Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getLastEvent()")
int32_t                           dxfg_DXFeed_getLastEvent2(graal_isolatethread_t *thread, dxfg_feed_t *feed, const dxfg_event_type_t *event, DXFG_OUT dxfg_event_type_t **lastEvent);

/**
 * @brief The bulk version of dxfg_DXFeed_getLastEvent2(): gets a new list with the last event for every event of the
 * given list, in the same order.
 *
 * The given list and its events are only read (see dxfg_DXFeed_getLastEvent2()).
 *
 * Java: `DXFeed.getLastEvents(Collection)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getLastEvents().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] events The events.
 * @param[out] lastEvents The pointer to the new list of new events (NULL on error). Free the list with
 * dxfg_CList_EventType_release(). The pointer itself must not be NULL (an error then).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the list or one of its events is NULL, or an event is of a type that is not lasting). Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getLastEvents()")
int32_t                           dxfg_DXFeed_getLastEvents2(graal_isolatethread_t *thread, dxfg_feed_t *feed, const dxfg_event_type_list *events, DXFG_OUT dxfg_event_type_list **lastEvents);

/**
 * @brief Requests the last event for the event class and symbol and returns the promise of the result.
 *
 * Works only for the lasting events (Quote, Profile, etc.). The promise completes with a new event. The request is sent
 * to the uplink data provider; it works also without a subscription. The result is a snapshot: use a subscription for
 * the updates in real time. The promise is cancelled when the endpoint of the feed is closed. If the data is not
 * available for a transient reason (no subscription, no connection to the uplink, etc.), the promise completes when the
 * issue is resolved, which may take arbitrarily long: wait with a timeout (dxfg_Promise_await2()). If the data is
 * permanently not available (not supported), the promise completes exceptionally with IllegalArgumentException. With
 * the #DXFG_ENDPOINT_ROLE_STREAM_FEED role the promise completes exceptionally.
 *
 * A handler added by dxfg_Promise_whenDone() is called by a thread of the executor of the endpoint of the feed (see
 * dxfg_DXEndpoint_Builder_withExecutor()).
 *
 * Java: `DXFeed.getLastEventPromise(Class, Object)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getLastEventPromise().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the event, a lasting one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership.
 * @return The promise (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Pass it to the
 * dxfg_Promise_* functions as `&result->handler` (a dxfg_promise_t) and read the event with
 * dxfg_Promise_EventType_getResult(). Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getLastEventPromise()")
dxfg_promise_event_t*             dxfg_DXFeed_getLastEventPromise(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol);

/**
 * @brief The bulk version of dxfg_DXFeed_getLastEventPromise(): requests the last events for one event class and
 * several symbols.
 *
 * It is much faster than the requests of the single events. Every promise of the result is a promise of an event (see
 * dxfg_promise_event_t) and is used as in dxfg_DXFeed_getLastEventPromise(). Use dxfg_Promises_allOf() to wait for all
 * of them with one timeout.
 *
 * Java: `DXFeed.getLastEventsPromises(Class, Collection)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getLastEventsPromises().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, a lasting one.
 * @param[in] symbols The symbols; the list is only read, the caller keeps its ownership. A NULL or empty list gives an
 * empty result.
 * @return The list of the promises, one per symbol, in the order of the symbols (release it with
 * dxfg_CList_JavaObjectHandler_release(), passing its `list` field; that releases every promise handle), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getLastEventsPromises()")
dxfg_promise_list*                dxfg_DXFeed_getLastEventsPromises(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_list *symbols);

/**
 * @brief Requests the indexed events for the event class, symbol and source and returns the promise of the result.
 *
 * Works only for the indexed events. The promise completes with a new list of new events ordered by index; read it with
 * dxfg_Promise_List_EventType_getResult(). The promise completes only when a consistent snapshot has been received; the
 * event flags of the events are always 0, so the list must not be passed to dxfg_DXPublisher_publishEvents(). The
 * request is sent to the uplink data provider; it works also without a subscription. The result is a snapshot: use a
 * subscription for the updates in real time. The promise is cancelled when the endpoint of the feed is closed. If the
 * data is not available for a transient reason (no subscription, no connection to the uplink, etc.), the promise
 * completes when the issue is resolved, which may take arbitrarily long: wait with a timeout (dxfg_Promise_await2()).
 * If the data is permanently not available (not supported), the promise completes exceptionally with
 * IllegalArgumentException. With the #DXFG_ENDPOINT_ROLE_STREAM_FEED role the promise completes exceptionally.
 *
 * Use the source with the identifier 0 (the default source) for the events that have no multiple sources (e.g., Series)
 * and an order source for the events with multiple sources (Order, etc.). The synthetic order sources are not
 * supported.
 *
 * A handler added by dxfg_Promise_whenDone() is called by a thread of the executor of the endpoint of the feed (see
 * dxfg_DXEndpoint_Builder_withExecutor()).
 *
 * Java: `DXFeed.getIndexedEventsPromise(Class, Object, IndexedEventSource)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getIndexedEventsPromise().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, an indexed one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership. It must not be a
 * dxfg_indexed_event_subscription_symbol_t.
 * @param[in] source The source of the events; it is only read, the caller keeps its ownership.
 * @return The promise (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Pass it to the
 * dxfg_Promise_* functions as `&result->base`. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getIndexedEventsPromise()")
dxfg_promise_events_t*            dxfg_DXFeed_getIndexedEventsPromise(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol, dxfg_indexed_event_source_t* source);

/**
 * @brief Requests the time series events for the event class, symbol and time range and returns the promise of the
 * result.
 *
 * Works only for the time series events. The promise completes with a new list of new events ordered by time; read it
 * with dxfg_Promise_List_EventType_getResult(). The range and the depth of the events are limited by the uplink data
 * provider. The promise completes only when a consistent snapshot has been received; the event flags of the events are
 * always 0, so the list must not be passed to dxfg_DXPublisher_publishEvents(). The request is sent to the uplink data
 * provider; it works also without a subscription. The result is a snapshot: use a subscription for the updates in real
 * time. The promise is cancelled when the endpoint of the feed is closed. If the data is not available for a transient
 * reason (no subscription, no connection to the uplink, etc.), the promise completes when the issue is resolved, which
 * may take arbitrarily long: wait with a timeout (dxfg_Promise_await2()). If the data is permanently not available (not
 * supported), the promise completes exceptionally with IllegalArgumentException. With the
 * #DXFG_ENDPOINT_ROLE_STREAM_FEED role the promise completes exceptionally.
 *
 * A handler added by dxfg_Promise_whenDone() is called by a thread of the executor of the endpoint of the feed (see
 * dxfg_DXEndpoint_Builder_withExecutor()).
 *
 * Java: `DXFeed.getTimeSeriesPromise(Class, Object, long, long)`.
 *
 * @deprecated Use dxfg_v4_DXFeed_getTimeSeriesPromise().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] clazz The class of the events, a time series one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership. It must not be a
 * dxfg_time_series_subscription_symbol_t.
 * @param[in] fromTime The time (inclusive) to request the events from, in milliseconds since 1970-01-01 00:00:00 UTC
 * (see the time of a time series event).
 * @param[in] toTime The time (inclusive) to request the events to, in the same unit; INT64_MAX means no upper limit.
 * @return The promise (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Pass it to the
 * dxfg_Promise_* functions as `&result->base`. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeed_getTimeSeriesPromise()")
dxfg_promise_events_t*            dxfg_DXFeed_getTimeSeriesPromise(graal_isolatethread_t *thread, dxfg_feed_t *feed, dxfg_event_clazz_t clazz, dxfg_symbol_t *symbol, int64_t fromTime, int64_t toTime);

/**
 * @brief Sets the earliest time from which the time series subscription receives the events.
 *
 * The initial value is INT64_MAX, which means that no events are received until the time is set.
 *
 * Java: `DXFeedTimeSeriesSubscription.setFromTime(long)`.
 *
 * @deprecated Use dxfg_v4_DXFeedTimeSeriesSubscription_setFromTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The time series subscription.
 * @param[in] fromTime The time in milliseconds since 1970-01-01 00:00:00 UTC.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedTimeSeriesSubscription_setFromTime()")
int32_t                           dxfg_DXFeedTimeSeriesSubscription_setFromTime(graal_isolatethread_t *thread, dxfg_time_series_subscription_t *sub, int64_t fromTime);

/**
 * @brief The function that is called when a promise is done.
 *
 * It is added by dxfg_Promise_whenDone() or dxfg_Promise_whenDoneAsync() and called exactly once, when the promise
 * completes normally, completes exceptionally or is cancelled. The thread that calls it is described in
 * dxfg_Promise_whenDone() and dxfg_Promise_whenDoneAsync(). The function can call the SDK functions with the given
 * `thread`, e.g. read the result with dxfg_Promise_EventType_getResult().
 *
 * `promise` is the pointer that was passed to dxfg_Promise_whenDone() or dxfg_Promise_whenDoneAsync(): the caller owns
 * it, the SDK does not release it and does not copy it, so the caller must keep it valid until the function has been
 * called.
 *
 * Java: `PromiseHandler.promiseDone(Promise)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise that is done: the pointer given when the function was added.
 * @param[in] user_data The pointer that was passed with the function as is; the SDK does not use it.
 */
typedef void (*dxfg_promise_handler_function)(graal_isolatethread_t *thread, dxfg_promise_t *promise, void *user_data);

/**
 * @brief Tells whether the promise is done: completed with a result, completed exceptionally or cancelled.
 *
 * Java: `Promise.isDone()`.
 *
 * @deprecated Use dxfg_v4_Promise_isDone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return 1 if the promise is done, 0 if not, or -1 on error (e.g., the promise is NULL). Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_isDone()")
int32_t               dxfg_Promise_isDone(graal_isolatethread_t *thread, dxfg_promise_t *promise);

/**
 * @brief Tells whether the promise has completed normally, with a result.
 *
 * Java: `Promise.hasResult()`.
 *
 * @deprecated Use dxfg_v4_Promise_hasResult().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return 1 if the promise has completed normally, 0 if not, or -1 on error (e.g., the promise is NULL). Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_hasResult()")
int32_t               dxfg_Promise_hasResult(graal_isolatethread_t *thread, dxfg_promise_t *promise);

/**
 * @brief Tells whether the promise has completed exceptionally or was cancelled.
 *
 * Use dxfg_Promise_getException() to get the exception.
 *
 * Java: `Promise.hasException()`.
 *
 * @deprecated Use dxfg_v4_Promise_hasException().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return 1 if the promise has completed exceptionally or was cancelled, 0 if not, or -1 on error (e.g., the promise is
 * NULL). Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_hasException()")
int32_t               dxfg_Promise_hasException(graal_isolatethread_t *thread, dxfg_promise_t *promise);

/**
 * @brief Tells whether the promise was cancelled.
 *
 * Java: `Promise.isCancelled()`.
 *
 * @deprecated Use dxfg_v4_Promise_isCancelled().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return 1 if the promise was cancelled, 0 if not, or -1 on error (e.g., the promise is NULL). Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_isCancelled()")
int32_t               dxfg_Promise_isCancelled(graal_isolatethread_t *thread, dxfg_promise_t *promise);

/**
 * @brief Returns the result of a promise of an event.
 *
 * Every call returns a new event. The function is for the promises that complete with an event (see
 * dxfg_promise_event_t), a promise of another kind makes it fail. Use dxfg_Promise_hasResult() to tell whether the
 * promise has a result.
 *
 * Java: `Promise.getResult()`.
 *
 * @deprecated Use dxfg_v4_Promise_EventType_getResult().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise of an event.
 * @return The new event (release it with dxfg_EventType_release()), or NULL if the promise has no result (it is not
 * done, completed exceptionally or was cancelled), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get
 * the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_EventType_getResult()")
dxfg_event_type_t*    dxfg_Promise_EventType_getResult(graal_isolatethread_t *thread, dxfg_promise_event_t *promise);

/**
 * @brief Returns the result of a promise of a list of events.
 *
 * Every call returns a new list of new events. The function is for the promises that complete with a list of events
 * (see dxfg_promise_events_t), a promise of another kind makes it fail. If the promise has no result (it is not done,
 * completed exceptionally or was cancelled), the result is an empty list (`size` is 0, `elements` is NULL), which must
 * be released as well: use dxfg_Promise_hasResult() to tell the cases apart.
 *
 * Java: `Promise.getResult()`.
 *
 * @deprecated Use dxfg_v4_Promise_List_EventType_getResult().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise of a list of events.
 * @return The new list of new events (release it with dxfg_CList_EventType_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_List_EventType_getResult()")
dxfg_event_type_list* dxfg_Promise_List_EventType_getResult(graal_isolatethread_t *thread, dxfg_promise_events_t *promise);

/**
 * @brief Returns the exception of the promise.
 *
 * The exception is a new structure; for a cancelled promise it describes CancellationException. If the promise has no
 * exception (it is not done or has completed normally), the function returns NULL.
 *
 * On error the function returns `(dxfg_exception_t *)-1`, not NULL (the SDK uses the same error value, -1, for all the
 * functions that return a status, and it becomes this pointer here): compare the result with it before using the
 * result. The exception is available by dxfg_get_and_clear_thread_exception_t() then.
 *
 * Java: `Promise.getException()`.
 *
 * @deprecated Use dxfg_v4_Promise_getException().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return The exception (release it with dxfg_Exception_release()), NULL if the promise has no exception, or
 * `(dxfg_exception_t *)-1` on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_getException()")
dxfg_exception_t*     dxfg_Promise_getException(graal_isolatethread_t *thread, dxfg_promise_t *promise);

/**
 * @brief The result of dxfg_Promise_awaitWithoutException() when the wait timed out (the promise is cancelled then).
 *
 * The value is 1. It differs from #DXFG_EXECUTE_SUCCESSFULLY (0), the result for a promise that has completed normally,
 * and from #DXFG_EXECUTE_FAIL (-1).
 */
#define DXFG_PROMISE_AWAIT_TIMED_OUT ((int32_t)1)

/**
 * @brief Waits for the promise to complete.
 *
 * The calling thread is blocked, without a time limit. If the wait is interrupted, the promise is cancelled. The
 * exception is PromiseException (it wraps the exception of the promise) or CancellationException. Read the result with
 * dxfg_Promise_EventType_getResult() or dxfg_Promise_List_EventType_getResult().
 *
 * Java: `Promise.await()`.
 *
 * @deprecated Use dxfg_v4_Promise_await().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if the promise has completed normally or #DXFG_EXECUTE_FAIL (-1) if it has
 * completed exceptionally, was cancelled or the wait was interrupted. Use dxfg_get_and_clear_thread_exception_t() to
 * get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_await()")
int32_t               dxfg_Promise_await(graal_isolatethread_t *thread, dxfg_promise_t *promise);

/**
 * @brief Waits for the promise to complete or the timeout to elapse.
 *
 * The calling thread is blocked. If the wait times out, the promise is cancelled and the function fails with
 * `CancellationException`. A timeout of 0 or less does not wait. Use dxfg_Promise_awaitWithoutException() to
 * distinguish the timeout without an exception.
 *
 * Java: `Promise.await(long, TimeUnit)`.
 *
 * @deprecated Use dxfg_v4_Promise_awaitWithTimeout().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] timeoutInMilliseconds The timeout in milliseconds.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if the promise has completed normally or #DXFG_EXECUTE_FAIL (-1) if it has
 * completed exceptionally, was cancelled, the wait timed out or was interrupted. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_awaitWithTimeout()")
int32_t               dxfg_Promise_await2(graal_isolatethread_t *thread, dxfg_promise_t *promise, int32_t timeoutInMilliseconds);

/**
 * @brief Waits for the promise to complete or the timeout to elapse, and reports the timeout without an exception.
 *
 * The calling thread is blocked. If the wait times out, the promise is cancelled and the function returns
 * #DXFG_PROMISE_AWAIT_TIMED_OUT.
 *
 * Before v3.5.0 the function returned #DXFG_EXECUTE_SUCCESSFULLY (0) on timeout as well.
 *
 * Java: `Promise.awaitWithoutException(long, TimeUnit)`.
 *
 * @deprecated Use dxfg_v4_Promise_awaitWithoutException().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] timeoutInMilliseconds The timeout in milliseconds.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if the promise has completed normally, #DXFG_PROMISE_AWAIT_TIMED_OUT (1) if
 * the wait timed out (the promise is cancelled), or #DXFG_EXECUTE_FAIL (-1) if the promise has completed exceptionally,
 * was cancelled or the wait was interrupted. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_awaitWithoutException()")
int32_t               dxfg_Promise_awaitWithoutException(graal_isolatethread_t *thread, dxfg_promise_t *promise, int32_t timeoutInMilliseconds);

/**
 * @brief Cancels the promise.
 *
 * The call does nothing if the promise is already done. Otherwise the promise becomes cancelled: its exception is
 * CancellationException, the waiters fail, and the handlers added by dxfg_Promise_whenDone() are called: for the
 * promises of the feed requests by a thread of the executor of the endpoint of the feed, for the other promises (such
 * as the one of dxfg_Promises_allOf()) before the function returns, by the calling thread.
 *
 * Java: `Promise.cancel()`.
 *
 * @deprecated Use dxfg_v4_Promise_cancel().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_cancel()")
int32_t               dxfg_Promise_cancel(graal_isolatethread_t *thread, dxfg_promise_t *promise);

/**
 * @brief Completes the promise of a list of events normally, with the given list.
 *
 * The call does nothing if the promise is already done (and returns success). Otherwise the handlers added by
 * dxfg_Promise_whenDone() are called by the same threads as for dxfg_Promise_cancel(). The promise must be a promise of
 * a list of events (see dxfg_promise_events_t), the SDK does not check it.
 *
 * The list and its events are only read: the SDK converts them into Java events, the caller keeps the ownership.
 *
 * Java: `Promise.complete(Object)`.
 *
 * @deprecated Use dxfg_v4_Promise_List_EventType_complete().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] events The list of events, the result of the promise. A NULL list is an empty list.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_List_EventType_complete()")
int32_t               dxfg_Promise_List_EventType_complete(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_event_type_list* events);

/**
 * @brief Completes the promise of an event normally, with the given event.
 *
 * The call does nothing if the promise is already done (and returns success). Otherwise the handlers added by
 * dxfg_Promise_whenDone() are called by the same threads as for dxfg_Promise_cancel(). The promise must be a promise of
 * an event (see dxfg_promise_event_t), the SDK does not check it.
 *
 * The event is only read: the SDK converts it into a Java event, the caller keeps the ownership. A NULL event completes
 * the promise with a null result.
 *
 * Java: `Promise.complete(Object)`.
 *
 * @deprecated Use dxfg_v4_Promise_EventType_complete().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] event The event, the result of the promise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_EventType_complete()")
int32_t               dxfg_Promise_EventType_complete(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_event_type_t* event);

/**
 * @brief Completes the promise exceptionally, with the given exception.
 *
 * The function does not work in this version: the SDK cannot convert a dxfg_exception_t into a Java exception, so a
 * call with a non-NULL exception fails with IllegalStateException, and a call with NULL fails with
 * NullPointerException. The promise is not changed. Use
 * dxfg_v4_Promise_completeExceptionallyWithClassNameAndMessage() to fail a promise.
 *
 * Java: `Promise.completeExceptionally(Throwable)`.
 *
 * @deprecated Use dxfg_v4_Promise_completeExceptionally().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] exception The exception.
 * @return #DXFG_EXECUTE_FAIL (-1) always in this version. Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_completeExceptionally()")
int32_t               dxfg_Promise_completeExceptionally(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_exception_t* exception);

/**
 * @brief Adds a handler that is called when the promise is done.
 *
 * The handler is called exactly once: when the promise completes normally, completes exceptionally or is cancelled.
 * There is no way to remove it. If the promise is already done, the handler is called before the function returns, by
 * the calling thread. Otherwise it is called by the thread that completes or cancels the promise. For the promises of
 * the feed requests it is a thread of the executor of the endpoint of the feed; for the others it is the thread that
 * called dxfg_Promise_cancel(), dxfg_Promise_EventType_complete(), etc., or the thread that waits with
 * dxfg_Promise_await2() when its timeout cancels the promise. Use dxfg_Promise_whenDoneAsync() to run the handler on a
 * given executor.
 *
 * The handler receives the pointer `promise` that is passed here: the SDK does not copy or release the structure, so it
 * must stay valid until the handler has been called. `user_data` is passed to the handler as is, the SDK does not use
 * it.
 *
 * Java: `Promise.whenDone(PromiseHandler)`.
 *
 * @deprecated Use dxfg_v4_Promise_whenDone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] promise_handler_function The function to call when the promise is done, see dxfg_promise_handler_function.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_whenDone()")
int32_t               dxfg_Promise_whenDone(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_promise_handler_function promise_handler_function, void *user_data);

/**
 * @brief Adds a handler that is called by the given executor when the promise is done.
 *
 * As dxfg_Promise_whenDone(), but the handler is never called by the calling thread or the thread that completes the
 * promise: it is called by a thread of the executor, as a task that the SDK gives to it. The handler is called exactly
 * once, and there is no way to remove it. The pointer `promise` must stay valid until the handler has been called, see
 * dxfg_Promise_whenDone().
 *
 * The executor handle is only read during the call: the caller keeps it and may release it after the function returns
 * (the promise keeps the executor object).
 *
 * Java: `Promise.whenDoneAsync(PromiseHandler, Executor)`.
 *
 * @deprecated Use dxfg_v4_Promise_whenDoneAsync().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] promise_handler_function The function to call when the promise is done, see dxfg_promise_handler_function.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it.
 * @param[in] executor The executor that calls the function, e.g. from dxfg_Executors_newFixedThreadPool(). A NULL
 * executor is an error.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_whenDoneAsync()")
int32_t               dxfg_Promise_whenDoneAsync(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_promise_handler_function promise_handler_function, void *user_data, dxfg_executor_t* executor);

/**
 * @brief Creates a promise that is already completed normally, with the given Java object as the result.
 *
 * The Java method is static: the SDK does not use the state of `promise`, the result does not depend on it. The result
 * of the new promise is the object of the handle `handler`.
 *
 * Java: `Promise.completed(Object)`.
 *
 * @deprecated Use dxfg_v4_Promise_completed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise A promise handle; it is not used.
 * @param[in] handler The handle of the Java object that is the result of the new promise; it is only read, the caller
 * keeps the ownership of the handle.
 * @return The new promise (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_completed()")
dxfg_promise_t*       dxfg_Promise_completed(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_java_object_handler *handler);

/**
 * @brief Creates a promise that is already failed, with the given exception.
 *
 * The function does not work in this version: the SDK cannot convert a dxfg_exception_t into a Java exception, so a
 * call with a non-NULL exception fails with IllegalStateException and returns NULL. A call with a NULL exception does
 * not fail either, but it creates a promise that is completed normally with a null result, not a failed one (the
 * implementation calls `Promise.completed(Object)`). The Java method is static: the SDK does not use the state of
 * `promise`. Use dxfg_v4_Promise_failedWithClassNameAndMessage() to create a failed promise.
 *
 * Java: `Promise.failed(Throwable)`.
 *
 * @deprecated Use dxfg_v4_Promise_failed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise A promise handle; it is not used.
 * @param[in] exception The exception.
 * @return The new promise (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promise_failed()")
dxfg_promise_t*       dxfg_Promise_failed(graal_isolatethread_t *thread, dxfg_promise_t *promise, dxfg_exception_t* exception);

/**
 * @brief Creates a promise that is done when all the given promises are done.
 *
 * The given promises may complete normally or exceptionally: the result of the new promise is null, their results are
 * not reflected in it, so inspect them one by one. If the list is empty, the new promise is already completed. When the
 * new promise completes for any reason (it is cancelled, for example), all the given promises are cancelled. A NULL
 * list is an empty list; a NULL promise in the list is an error.
 *
 * Java: `Promises.allOf(Collection)`.
 *
 * @deprecated Use dxfg_v4_Promises_allOf().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promises The list of the promises; it is only read, the caller keeps its ownership.
 * @return The new promise (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Promises_allOf()")
dxfg_promise_t*       dxfg_Promises_allOf(graal_isolatethread_t *thread, dxfg_promise_list *promises);


/** @} */ // end of Feed

#ifdef __cplusplus
}
#endif

#endif // DXFG_FEED_H
