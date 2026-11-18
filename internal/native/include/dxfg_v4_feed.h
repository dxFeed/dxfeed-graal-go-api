// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_feed.h and the v4 names of its types (see the group @ref V4_Feed).
 */

#ifndef DXFG_V4_FEED_H
#define DXFG_V4_FEED_H

#include "dxfg_feed.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_events.h"
#include "dxfg_v4_catch_exception.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Feed Feed (v4)
 *  @brief The v4 functions and type names of the group @ref Feed. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_subscription_t */
typedef dxfg_subscription_t dxfg_v4_subscription_t;

/** @copydoc dxfg_time_series_subscription_t */
typedef dxfg_time_series_subscription_t dxfg_v4_time_series_subscription_t;

/** @copydoc dxfg_executor_t */
typedef dxfg_executor_t dxfg_v4_executor_t;

/** @copydoc dxfg_feed_t */
typedef dxfg_feed_t dxfg_v4_feed_t;

/** @copydoc dxfg_promise_t */
typedef dxfg_promise_t dxfg_v4_promise_t;

/** @copydoc dxfg_promise_list */
typedef dxfg_promise_list dxfg_v4_promise_list;

/** @copydoc dxfg_promise_events_t */
typedef dxfg_promise_events_t dxfg_v4_promise_events_t;

/** @copydoc dxfg_promise_event_t */
typedef dxfg_promise_event_t dxfg_v4_promise_event_t;

/** @copydoc dxfg_promise_handler_function */
typedef dxfg_promise_handler_function dxfg_v4_promise_handler_function;

/**
 * @brief Returns the feed of the default endpoint of the process.
 *
 * Applies the SDK defaults of the QD system properties (a property that the user has set keeps its value) and registers
 * the default endpoint, so that dxfg_system_close_all_and_await_termination() closes it. Every call returns a new
 * handle of the same feed.
 *
 * Java: `DXFeed.getInstance()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] feed The pointer to the variable that receives the feed (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getInstance(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_feed_t **feed);

/**
 * @brief Creates a subscription for the events of one class and attaches it to the feed.
 *
 * The subscription has no symbols and no listeners yet: set them with the dxfg_DXFeedSubscription_* functions. The call
 * is dxfg_v4_DXFeedSubscription_new() followed by dxfg_v4_DXFeed_attachSubscription().
 *
 * Java: `DXFeed.createSubscription(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events to subscribe to.
 * @param[out] subscription The pointer to the variable that receives the subscription (release its handle with
 * dxfg_v4_JavaObjectHandler_release(); that neither closes nor detaches the subscription); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_createSubscription(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                          dxfg_v4_event_clazz_t eventClazz,
                                          DXFG_OUT dxfg_v4_subscription_t **subscription);

/**
 * @brief Creates a subscription for the events of several classes and attaches it to the feed.
 *
 * See dxfg_v4_DXFeed_createSubscription(). The list is only read: the caller keeps its ownership. An empty list is an
 * error.
 *
 * Java: `DXFeed.createSubscription(Class...)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazzes The classes of the events to subscribe to.
 * @param[out] subscription The pointer to the variable that receives the subscription (release its handle with
 * dxfg_v4_JavaObjectHandler_release(); that neither closes nor detaches the subscription); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_createSubscriptionWithEventTypes(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                                        dxfg_v4_event_clazz_list *eventClazzes,
                                                        DXFG_OUT dxfg_v4_subscription_t **subscription);

/**
 * @brief Creates a time series subscription for the events of one class and attaches it to the feed.
 *
 * Use the classes of the time series events (see dxfg_v4_event_clazz_t). The subscription has no symbols and no
 * listeners yet. Its `fromTime` is initially INT64_MAX, which means that no events are received until
 * dxfg_v4_DXFeedTimeSeriesSubscription_setFromTime() is called.
 *
 * The result contains dxfg_v4_subscription_t as its first field `sub`: pass `&result->sub` (or the pointer cast) to the
 * functions that take a dxfg_v4_subscription_t.
 *
 * Java: `DXFeed.createTimeSeriesSubscription(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the time series events to subscribe to.
 * @param[out] timeSeriesSubscription The pointer to the variable that receives the subscription (release its handle
 * with dxfg_v4_JavaObjectHandler_release(); that neither closes nor detaches the subscription); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_createTimeSeriesSubscription(
    graal_isolatethread_t *thread, dxfg_v4_feed_t *feed, dxfg_v4_event_clazz_t eventClazz,
    DXFG_OUT dxfg_v4_time_series_subscription_t **timeSeriesSubscription);

/**
 * @brief Creates a time series subscription for the events of several classes and attaches it to the feed.
 *
 * See dxfg_v4_DXFeed_createTimeSeriesSubscription(). The list is only read: the caller keeps its ownership. An empty
 * list is an error.
 *
 * Java: `DXFeed.createTimeSeriesSubscription(Class...)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazzes The classes of the time series events to subscribe to.
 * @param[out] timeSeriesSubscription The pointer to the variable that receives the subscription (release its handle
 * with dxfg_v4_JavaObjectHandler_release(); that neither closes nor detaches the subscription); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_createTimeSeriesSubscriptionWithEventTypes(
    graal_isolatethread_t *thread, dxfg_v4_feed_t *feed, dxfg_v4_event_clazz_list *eventClazzes,
    DXFG_OUT dxfg_v4_time_series_subscription_t **timeSeriesSubscription);

/**
 * @brief Attaches the subscription to the feed.
 *
 * From now on the feed publishes the data to the subscription. The call does nothing if the subscription is already
 * attached to this feed. A NULL feed or subscription is an error.
 *
 * Java: `DXFeed.attachSubscription(DXFeedSubscription)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] sub The subscription, created by dxfg_v4_DXFeedSubscription_new() or
 * dxfg_v4_DXFeedSubscription_newWithEventTypes() (or any other subscription handle).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_attachSubscription(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                          dxfg_v4_subscription_t *sub);

/**
 * @brief Detaches the subscription from the feed.
 *
 * The call does nothing if the subscription is not attached to this feed. The data that the subscription has already
 * received stays in it. A NULL feed or subscription is an error.
 *
 * Java: `DXFeed.detachSubscription(DXFeedSubscription)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] sub The subscription.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_detachSubscription(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                          dxfg_v4_subscription_t *sub);

/**
 * @brief Detaches the subscription from the feed and clears the data that was delivered to the subscription.
 *
 * The data is cleared by publishing empty events to the subscription. The call does nothing if the subscription is not
 * attached to this feed. A NULL feed or subscription is an error. For a subscription attached to an event queue
 * (dxfg_v4_EventQueue_new(), dxfg_v4_EventQueue_attach()) the call does not wait for a poll: the events delivered while
 * it runs, the empty ones among them, arrive in the poll in progress or the next one. For a subscription whose
 * executor is the executor of a shared queue (dxfg_v4_EventQueue_getExecutor()) the call waits until another thread
 * runs its pending work in a poll or the queue is closed: do not call it on the polling thread.
 *
 * Java: `DXFeed.detachSubscriptionAndClear(DXFeedSubscription)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] sub The subscription.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_detachSubscriptionAndClear(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                                  dxfg_v4_subscription_t *sub);

/**
 * @brief Returns the last event for the event class and symbol if there is a subscription for them.
 *
 * Works only for the lasting events (Quote, Profile, etc.). The function does not make remote calls, it reads the local
 * cache of the feed. The events are in the cache only if an attached subscription is subscribed to the event class and
 * symbol (a wildcard subscription does not count; a subscription can also be defined by the endpoint properties). Use
 * dxfg_v4_DXFeed_getLastEventPromise() to request an event without a subscription.
 *
 * If there is a subscription, but the event has not arrived yet, the event has only its symbol set and the default
 * values of the other fields. If there is no subscription, the function returns NULL: it is not an error, so
 * dxfg_v4_get_and_clear_thread_exception() returns NULL then. With the #DXFG_ENDPOINT_ROLE_STREAM_FEED role the
 * function always returns NULL.
 *
 * Java: `DXFeed.getLastEventIfSubscribed(Class, Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the event, a lasting one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership.
 * @param[out] event The pointer to the variable that receives the new event (release it with
 * dxfg_v4_EventType_release()), or NULL if there is no subscription for the event class and symbol; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getLastEventIfSubscribed(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                                dxfg_v4_event_clazz_t eventClazz, dxfg_v4_symbol_t *symbol,
                                                DXFG_OUT dxfg_v4_event_type_t **event);

/**
 * @brief Returns the indexed events for the event class, symbol and source if there is a subscription for them.
 *
 * Works only for the indexed events. The function does not make remote calls, it reads the local cache of the feed. The
 * events are in the cache only if an attached subscription is subscribed to the event class, symbol and source (a
 * subscription can also be defined by the endpoint properties). The events are ordered by index, their event flags are
 * always 0. The list is not a consistent snapshot: use dxfg_v4_DXFeed_getIndexedEventsPromise() for it. With the
 * #DXFG_ENDPOINT_ROLE_STREAM_FEED role there are never events.
 *
 * If there is no subscription (Java returns null) or the events have not arrived yet, the result is an empty list
 * (`size` is 0, `elements` is NULL); the two cases cannot be told apart by the result. The empty list must be released
 * as well.
 *
 * Java: `DXFeed.getIndexedEventsIfSubscribed(Class, Object, IndexedEventSource)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, an indexed one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership.
 * @param[in] source The name of the order source (OrderSource.valueOf(String)), a NUL-terminated UTF-8 string: a
 * predefined name such as `DEFAULT`, or at most 4 alphanumeric characters. A malformed or NULL name is an error. The
 * string is only read.
 * @param[out] events The pointer to the variable that receives the new list of new events (release it with
 * dxfg_v4_event_type_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getIndexedEventsIfSubscribed(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                                    dxfg_v4_event_clazz_t eventClazz, dxfg_v4_symbol_t *symbol,
                                                    const char *source, DXFG_OUT dxfg_v4_event_type_list **events);

/**
 * @brief Returns the time series events for the event class, symbol and time range if there is a subscription for them.
 *
 * Works only for the time series events. The function does not make remote calls, it reads the local cache of the feed.
 * The events are in the cache only if an attached subscription is subscribed to the event class, symbol and time (a
 * subscription can also be defined by the endpoint properties). The events are ordered by time, their event flags are
 * always 0. The list is not a consistent snapshot: use dxfg_v4_DXFeed_getTimeSeriesPromise() for it. With the
 * #DXFG_ENDPOINT_ROLE_STREAM_FEED role there are never events.
 *
 * If there is no subscription (Java returns null) or the events have not arrived yet, the result is an empty list
 * (`size` is 0, `elements` is NULL); the two cases cannot be told apart by the result. The empty list must be released
 * as well.
 *
 * Java: `DXFeed.getTimeSeriesIfSubscribed(Class, Object, long, long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, a time series one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership. It must not be a
 * dxfg_v4_time_series_subscription_symbol_t.
 * @param[in] from_time The time (inclusive) to return the events from, in milliseconds since 1970-01-01 00:00:00 UTC
 * (see the time of a time series event).
 * @param[in] to_time The time (inclusive) to return the events to, in the same unit; INT64_MAX means no upper limit.
 * @param[out] events The pointer to the variable that receives the new list of new events (release it with
 * dxfg_v4_event_type_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getTimeSeriesIfSubscribed(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                                 dxfg_v4_event_clazz_t eventClazz, dxfg_v4_symbol_t *symbol,
                                                 int64_t from_time, int64_t to_time,
                                                 DXFG_OUT dxfg_v4_event_type_list **events);

/**
 * @brief Gets the last event for the type and symbol of the given event.
 *
 * The result is a new event with the values of the last event, or a copy of the given event when the last event is not
 * available (no subscription, the data have not arrived yet, etc.), since Java leaves the event unchanged then.
 *
 * The given event is only read, as by dxfg_v4_DXPublisher_publishEvents(): the caller allocates it (with any
 * allocator), fills it and frees it; the SDK copies it during the call. The event must be of a lasting type (Quote,
 * Profile, etc.).
 *
 * The function does not make remote calls, it reads the local cache of the feed. The events are in the cache only if an
 * attached subscription is subscribed to the event type and symbol (a wildcard subscription does not count). With the
 * #DXFG_ENDPOINT_ROLE_STREAM_FEED role the last event is never available.
 *
 * Java: `DXFeed.getLastEvent(LastingEvent)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] event The event: its type and symbol, and the values to return when the last event is not available.
 * @param[out] lastEvent The pointer to the new event (NULL on error). Free the event with dxfg_v4_EventType_release().
 * The pointer itself must not be NULL (an error then).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the event is NULL or of a type that is not lasting). Use dxfg_v4_get_and_clear_thread_exception() to determine if an
 * exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getLastEvent(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                    const dxfg_v4_event_type_t *event, DXFG_OUT dxfg_v4_event_type_t **lastEvent);

/**
 * @brief The bulk version of dxfg_v4_DXFeed_getLastEvent(): gets a new list with the last event for every event of the
 * given list, in the same order.
 *
 * The given list and its events are only read (see dxfg_v4_DXFeed_getLastEvent()).
 *
 * Java: `DXFeed.getLastEvents(Collection)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] events The events.
 * @param[out] lastEvents The pointer to the new list of new events (NULL on error). Free the list with
 * dxfg_v4_event_type_list_release(). The pointer itself must not be NULL (an error then).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the list or one of its events is NULL, or an event is of a type that is not lasting). Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getLastEvents(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                     const dxfg_v4_event_type_list *events,
                                     DXFG_OUT dxfg_v4_event_type_list **lastEvents);

/**
 * @brief Requests the last event for the event class and symbol and returns the promise of the result.
 *
 * Works only for the lasting events (Quote, Profile, etc.). The promise completes with a new event. The request is sent
 * to the uplink data provider; it works also without a subscription. The result is a snapshot: use a subscription for
 * the updates in real time. The promise is cancelled when the endpoint of the feed is closed. If the data is not
 * available for a transient reason (no subscription, no connection to the uplink, etc.), the promise completes when the
 * issue is resolved, which may take arbitrarily long: wait with a timeout (dxfg_v4_Promise_awaitWithTimeout()). If the
 * data is permanently not available (not supported), the promise completes exceptionally with IllegalArgumentException.
 * With the #DXFG_ENDPOINT_ROLE_STREAM_FEED role the promise completes exceptionally.
 *
 * A handler added by dxfg_v4_Promise_whenDone() is called by a thread of the executor of the endpoint of the feed (see
 * dxfg_v4_DXEndpoint_Builder_withExecutor()).
 *
 * Java: `DXFeed.getLastEventPromise(Class, Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the event, a lasting one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership.
 * @param[out] promise The pointer to the variable that receives the promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()). Pass it to the dxfg_Promise_* functions as `&result->handler` (a
 * dxfg_v4_promise_t) and read the event with dxfg_v4_Promise_EventType_getResult(); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getLastEventPromise(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                           dxfg_v4_event_clazz_t eventClazz, dxfg_v4_symbol_t *symbol,
                                           DXFG_OUT dxfg_v4_promise_event_t **promise);

/**
 * @brief The bulk version of dxfg_v4_DXFeed_getLastEventPromise(): requests the last events for one event class and
 * several symbols.
 *
 * It is much faster than the requests of the single events. Every promise of the result is a promise of an event (see
 * dxfg_v4_promise_event_t) and is used as in dxfg_v4_DXFeed_getLastEventPromise(). Use dxfg_v4_Promises_allOf() to wait
 * for all of them with one timeout.
 *
 * Java: `DXFeed.getLastEventsPromises(Class, Collection)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, a lasting one.
 * @param[in] symbols The symbols; the list is only read, the caller keeps its ownership. A NULL or empty list gives an
 * empty result.
 * @param[out] promises The pointer to the variable that receives the list of the promises, one per symbol, in the order
 * of the symbols (release it with dxfg_v4_java_object_handler_list_release(), passing its `list` field; that releases
 * every promise handle); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getLastEventsPromises(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                             dxfg_v4_event_clazz_t eventClazz, dxfg_v4_symbol_list *symbols,
                                             DXFG_OUT dxfg_v4_promise_list **promises);

/**
 * @brief Requests the indexed events for the event class, symbol and source and returns the promise of the result.
 *
 * Works only for the indexed events. The promise completes with a new list of new events ordered by index; read it with
 * dxfg_v4_Promise_List_EventType_getResult(). The promise completes only when a consistent snapshot has been received;
 * the event flags of the events are always 0, so the list must not be passed to dxfg_v4_DXPublisher_publishEvents().
 * The request is sent to the uplink data provider; it works also without a subscription. The result is a snapshot: use
 * a subscription for the updates in real time. The promise is cancelled when the endpoint of the feed is closed. If the
 * data is not available for a transient reason (no subscription, no connection to the uplink, etc.), the promise
 * completes when the issue is resolved, which may take arbitrarily long: wait with a timeout
 * (dxfg_v4_Promise_awaitWithTimeout()). If the data is permanently not available (not supported), the promise completes
 * exceptionally with IllegalArgumentException. With the #DXFG_ENDPOINT_ROLE_STREAM_FEED role the promise completes
 * exceptionally.
 *
 * Use the source with the identifier 0 (the default source) for the events that have no multiple sources (e.g., Series)
 * and an order source for the events with multiple sources (Order, etc.). The synthetic order sources are not
 * supported.
 *
 * A handler added by dxfg_v4_Promise_whenDone() is called by a thread of the executor of the endpoint of the feed (see
 * dxfg_v4_DXEndpoint_Builder_withExecutor()).
 *
 * Java: `DXFeed.getIndexedEventsPromise(Class, Object, IndexedEventSource)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] eventClazz The class of the events, an indexed one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership. It must not be a
 * dxfg_v4_indexed_event_subscription_symbol_t.
 * @param[in] source The source of the events; it is only read, the caller keeps its ownership.
 * @param[out] promise The pointer to the variable that receives the promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()). Pass it to the dxfg_Promise_* functions as `&result->base`; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getIndexedEventsPromise(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                               dxfg_v4_event_clazz_t eventClazz, dxfg_v4_symbol_t *symbol,
                                               dxfg_v4_indexed_event_source_t* source,
                                               DXFG_OUT dxfg_v4_promise_events_t **promise);

/**
 * @brief Requests the time series events for the event class, symbol and time range and returns the promise of the
 * result.
 *
 * Works only for the time series events. The promise completes with a new list of new events ordered by time; read it
 * with dxfg_v4_Promise_List_EventType_getResult(). The range and the depth of the events are limited by the uplink data
 * provider. The promise completes only when a consistent snapshot has been received; the event flags of the events are
 * always 0, so the list must not be passed to dxfg_v4_DXPublisher_publishEvents(). The request is sent to the uplink
 * data provider; it works also without a subscription. The result is a snapshot: use a subscription for the updates in
 * real time. The promise is cancelled when the endpoint of the feed is closed. If the data is not available for a
 * transient reason (no subscription, no connection to the uplink, etc.), the promise completes when the issue is
 * resolved, which may take arbitrarily long: wait with a timeout (dxfg_v4_Promise_awaitWithTimeout()). If the data is
 * permanently not available (not supported), the promise completes exceptionally with IllegalArgumentException. With
 * the #DXFG_ENDPOINT_ROLE_STREAM_FEED role the promise completes exceptionally.
 *
 * A handler added by dxfg_v4_Promise_whenDone() is called by a thread of the executor of the endpoint of the feed (see
 * dxfg_v4_DXEndpoint_Builder_withExecutor()).
 *
 * Java: `DXFeed.getTimeSeriesPromise(Class, Object, long, long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] feed The feed.
 * @param[in] clazz The class of the events, a time series one.
 * @param[in] symbol The symbol; it is only read, the caller keeps its ownership. It must not be a
 * dxfg_v4_time_series_subscription_symbol_t.
 * @param[in] fromTime The time (inclusive) to request the events from, in milliseconds since 1970-01-01 00:00:00 UTC
 * (see the time of a time series event).
 * @param[in] toTime The time (inclusive) to request the events to, in the same unit; INT64_MAX means no upper limit.
 * @param[out] promise The pointer to the variable that receives the promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()). Pass it to the dxfg_Promise_* functions as `&result->base`; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeed_getTimeSeriesPromise(graal_isolatethread_t *thread, dxfg_v4_feed_t *feed,
                                            dxfg_v4_event_clazz_t clazz, dxfg_v4_symbol_t *symbol, int64_t fromTime,
                                            int64_t toTime, DXFG_OUT dxfg_v4_promise_events_t **promise);

/**
 * @brief Sets the earliest time from which the time series subscription receives the events.
 *
 * The initial value is INT64_MAX, which means that no events are received until the time is set.
 *
 * Java: `DXFeedTimeSeriesSubscription.setFromTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The time series subscription.
 * @param[in] fromTime The time in milliseconds since 1970-01-01 00:00:00 UTC.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedTimeSeriesSubscription_setFromTime(graal_isolatethread_t *thread,
                                                         dxfg_v4_time_series_subscription_t *sub, int64_t fromTime);

/**
 * @brief Returns the earliest time from which the time series subscription receives the events.
 *
 * It is the value set by dxfg_v4_DXFeedTimeSeriesSubscription_setFromTime(), INT64_MAX for a new subscription. The
 * handle must be one of a time series subscription (dxfg_v4_DXFeed_createTimeSeriesSubscription()): the handle of
 * another subscription, cast to the type, fails the call with `IllegalArgumentException`.
 *
 * Java: `DXFeedTimeSeriesSubscription.getFromTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The time series subscription.
 * @param[out] fromTime The pointer to the variable that receives the time in milliseconds since 1970-01-01 00:00:00
 * UTC; must not be NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_DXFeedTimeSeriesSubscription_getFromTime(graal_isolatethread_t *thread,
                                                         dxfg_v4_time_series_subscription_t *sub,
                                                         DXFG_OUT int64_t *fromTime);

/**
 * @brief Tells whether the promise is done: completed with a result, completed exceptionally or cancelled.
 *
 * Java: `Promise.isDone()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[out] done The pointer to the variable that receives 1 if the promise is done, 0 if not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the promise is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_isDone(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise, DXFG_OUT int32_t *done);

/**
 * @brief Tells whether the promise has completed normally, with a result.
 *
 * Java: `Promise.hasResult()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[out] hasResult The pointer to the variable that receives 1 if the promise has completed normally, 0 if not;
 * must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the promise is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_hasResult(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                  DXFG_OUT int32_t *hasResult);

/**
 * @brief Tells whether the promise has completed exceptionally or was cancelled.
 *
 * Use dxfg_v4_Promise_getException() to get the exception.
 *
 * Java: `Promise.hasException()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[out] hasException The pointer to the variable that receives 1 if the promise has completed exceptionally or
 * was cancelled, 0 if not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the promise is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_hasException(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                     DXFG_OUT int32_t *hasException);

/**
 * @brief Tells whether the promise was cancelled.
 *
 * Java: `Promise.isCancelled()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[out] cancelled The pointer to the variable that receives 1 if the promise was cancelled, 0 if not; must not be
 * NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the promise is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_isCancelled(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                    DXFG_OUT int32_t *cancelled);

/**
 * @brief Returns the result of a promise of an event.
 *
 * Every call returns a new event. The function is for the promises that complete with an event (see
 * dxfg_v4_promise_event_t), a promise of another kind makes it fail. Use dxfg_v4_Promise_hasResult() to tell whether
 * the promise has a result.
 *
 * Java: `Promise.getResult()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise of an event.
 * @param[out] event The pointer to the variable that receives the new event (release it with
 * dxfg_v4_EventType_release()), or NULL if the promise has no result (it is not done, completed exceptionally or was
 * cancelled); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_EventType_getResult(graal_isolatethread_t *thread, dxfg_v4_promise_event_t *promise,
                                            DXFG_OUT dxfg_v4_event_type_t **event);

/**
 * @brief Returns the result of a promise of a list of events.
 *
 * Every call returns a new list of new events. The function is for the promises that complete with a list of events
 * (see dxfg_v4_promise_events_t), a promise of another kind makes it fail. If the promise has no result (it is not
 * done, completed exceptionally or was cancelled), the result is an empty list (`size` is 0, `elements` is NULL), which
 * must be released as well: use dxfg_v4_Promise_hasResult() to tell the cases apart.
 *
 * Java: `Promise.getResult()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise of a list of events.
 * @param[out] events The pointer to the variable that receives the new list of new events (release it with
 * dxfg_v4_event_type_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_List_EventType_getResult(graal_isolatethread_t *thread, dxfg_v4_promise_events_t *promise,
                                                 DXFG_OUT dxfg_v4_event_type_list **events);

/**
 * @brief Gets the exception of the promise.
 *
 * The exception is a new structure; for a cancelled promise it describes CancellationException. If the promise has no
 * exception (it is not done or has completed normally), the function writes NULL.
 *
 * Java: `Promise.getException()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[out] exception The pointer to the variable that receives the exception (release it with
 * dxfg_v4_Exception_release()), or NULL if the promise has no exception; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_getException(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                     DXFG_OUT dxfg_v4_exception_t **exception);

/**
 * @brief Waits for the promise to complete.
 *
 * The calling thread is blocked, without a time limit. If the wait is interrupted, the promise is cancelled. The
 * exception is PromiseException (it wraps the exception of the promise) or CancellationException. Read the result with
 * dxfg_v4_Promise_EventType_getResult() or dxfg_v4_Promise_List_EventType_getResult().
 *
 * Java: `Promise.await()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if the promise has completed normally or #DXFG_EXECUTE_FAIL (-1) if it has
 * completed exceptionally, was cancelled or the wait was interrupted. Use dxfg_v4_get_and_clear_thread_exception() to
 * get the exception.
 */
int32_t dxfg_v4_Promise_await(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise);

/**
 * @brief Waits for the promise to complete or the timeout to elapse.
 *
 * The calling thread is blocked. If the wait times out, the promise is cancelled and the function fails with
 * `CancellationException`. A timeout of 0 or less does not wait. Use dxfg_v4_Promise_awaitWithoutException() to
 * distinguish the timeout without an exception.
 *
 * Java: `Promise.await(long, TimeUnit)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] timeoutInMilliseconds The timeout in milliseconds.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if the promise has completed normally or #DXFG_EXECUTE_FAIL (-1) if it has
 * completed exceptionally, was cancelled, the wait timed out or was interrupted. Use
 * dxfg_v4_get_and_clear_thread_exception() to get the exception.
 */
int32_t dxfg_v4_Promise_awaitWithTimeout(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                         int32_t timeoutInMilliseconds);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] timeoutInMilliseconds The timeout in milliseconds.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if the promise has completed normally, #DXFG_PROMISE_AWAIT_TIMED_OUT (1) if
 * the wait timed out (the promise is cancelled), or #DXFG_EXECUTE_FAIL (-1) if the promise has completed exceptionally,
 * was cancelled or the wait was interrupted. Use dxfg_v4_get_and_clear_thread_exception() to get the exception.
 */
int32_t dxfg_v4_Promise_awaitWithoutException(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                              int32_t timeoutInMilliseconds);

/**
 * @brief Cancels the promise.
 *
 * The call does nothing if the promise is already done. Otherwise the promise becomes cancelled: its exception is
 * CancellationException, the waiters fail, and the handlers added by dxfg_v4_Promise_whenDone() are called: for the
 * promises of the feed requests by a thread of the executor of the endpoint of the feed, for the other promises (such
 * as the one of dxfg_v4_Promises_allOf()) before the function returns, by the calling thread.
 *
 * Java: `Promise.cancel()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_cancel(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise);

/**
 * @brief Completes the promise of a list of events normally, with the given list.
 *
 * The call does nothing if the promise is already done (and returns success). Otherwise the handlers added by
 * dxfg_v4_Promise_whenDone() are called by the same threads as for dxfg_v4_Promise_cancel(). The promise must be a
 * promise of a list of events (see dxfg_v4_promise_events_t), the SDK does not check it.
 *
 * The list and its events are only read: the SDK converts them into Java events, the caller keeps the ownership.
 *
 * Java: `Promise.complete(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] events The list of events, the result of the promise. A NULL list is an empty list.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_List_EventType_complete(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                                dxfg_v4_event_type_list* events);

/**
 * @brief Completes the promise of an event normally, with the given event.
 *
 * The call does nothing if the promise is already done (and returns success). Otherwise the handlers added by
 * dxfg_v4_Promise_whenDone() are called by the same threads as for dxfg_v4_Promise_cancel(). The promise must be a
 * promise of an event (see dxfg_v4_promise_event_t), the SDK does not check it.
 *
 * The event is only read: the SDK converts it into a Java event, the caller keeps the ownership. A NULL event completes
 * the promise with a null result.
 *
 * Java: `Promise.complete(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] event The event, the result of the promise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_EventType_complete(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                           dxfg_v4_event_type_t* event);

/**
 * @brief Completes the promise exceptionally, with the given exception.
 *
 * The function does not work in this version: the SDK cannot convert a dxfg_v4_exception_t into a Java exception, so a
 * call with a non-NULL exception fails with IllegalStateException, and a call with NULL fails with
 * NullPointerException. The promise is not changed. Use
 * dxfg_v4_Promise_completeExceptionallyWithClassNameAndMessage() to fail a promise.
 *
 * Java: `Promise.completeExceptionally(Throwable)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] exception The exception.
 * @return #DXFG_EXECUTE_FAIL (-1) always in this version. Use dxfg_v4_get_and_clear_thread_exception() to get the
 * exception.
 */
int32_t dxfg_v4_Promise_completeExceptionally(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                              dxfg_v4_exception_t* exception);

/**
 * @brief Adds a handler that is called when the promise is done.
 *
 * The handler is called exactly once: when the promise completes normally, completes exceptionally or is cancelled.
 * There is no way to remove it. If the promise is already done, the handler is called before the function returns, by
 * the calling thread. Otherwise it is called by the thread that completes or cancels the promise. For the promises of
 * the feed requests it is a thread of the executor of the endpoint of the feed; for the others it is the thread that
 * called dxfg_v4_Promise_cancel(), dxfg_v4_Promise_EventType_complete(), etc., or the thread that waits with
 * dxfg_v4_Promise_awaitWithTimeout() when its timeout cancels the promise. Use dxfg_v4_Promise_whenDoneAsync() to run
 * the handler on a given executor.
 *
 * The handler receives the pointer `promise` that is passed here: the SDK does not copy or release the structure, so it
 * must stay valid until the handler has been called. `user_data` is passed to the handler as is, the SDK does not use
 * it.
 *
 * Java: `Promise.whenDone(PromiseHandler)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] promise_handler_function The function to call when the promise is done, see
 * dxfg_v4_promise_handler_function.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_whenDone(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                 dxfg_v4_promise_handler_function promise_handler_function, void *user_data);

/**
 * @brief Adds a handler that is called by the given executor when the promise is done.
 *
 * As dxfg_v4_Promise_whenDone(), but the handler is never called by the calling thread or the thread that completes the
 * promise: it is called by a thread of the executor, as a task that the SDK gives to it. The handler is called exactly
 * once, and there is no way to remove it. The pointer `promise` must stay valid until the handler has been called, see
 * dxfg_v4_Promise_whenDone().
 *
 * The executor handle is only read during the call: the caller keeps it and may release it after the function returns
 * (the promise keeps the executor object).
 *
 * Java: `Promise.whenDoneAsync(PromiseHandler, Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise.
 * @param[in] promise_handler_function The function to call when the promise is done, see
 * dxfg_v4_promise_handler_function.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it.
 * @param[in] executor The executor that calls the function, e.g. from dxfg_v4_Executors_newFixedThreadPool(). A NULL
 * executor is an error.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_whenDoneAsync(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                      dxfg_v4_promise_handler_function promise_handler_function, void *user_data,
                                      dxfg_v4_executor_t* executor);

/**
 * @brief Creates a promise that is already completed normally, with the given Java object as the result.
 *
 * The Java method is static: the SDK does not use the state of `promise`, the result does not depend on it. The result
 * of the new promise is the object of the handle `handler`.
 *
 * Java: `Promise.completed(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise A promise handle; it is not used.
 * @param[in] handler The handle of the Java object that is the result of the new promise; it is only read, the caller
 * keeps the ownership of the handle.
 * @param[out] newPromise The pointer to the variable that receives the new promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_completed(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                  dxfg_v4_java_object_handler *handler, DXFG_OUT dxfg_v4_promise_t **newPromise);

/**
 * @brief Creates a promise that is already failed, with the given exception.
 *
 * The function does not work in this version: the SDK cannot convert a dxfg_v4_exception_t into a Java exception, so a
 * call with a non-NULL exception fails with IllegalStateException and returns NULL. A call with a NULL exception does
 * not fail either, but it creates a promise that is completed normally with a null result, not a failed one (the
 * implementation calls `Promise.completed(Object)`). The Java method is static: the SDK does not use the state of
 * `promise`. Use dxfg_v4_Promise_failedWithClassNameAndMessage() to create a failed promise.
 *
 * Java: `Promise.failed(Throwable)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise A promise handle; it is not used.
 * @param[in] exception The exception.
 * @param[out] newPromise The pointer to the variable that receives the new promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promise_failed(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                               dxfg_v4_exception_t* exception, DXFG_OUT dxfg_v4_promise_t **newPromise);

/**
 * @brief Creates a new promise that is not done yet.
 *
 * The application completes the promise with dxfg_v4_Promise_complete(),
 * dxfg_v4_Promise_completeExceptionallyWithClassNameAndMessage() or dxfg_v4_Promise_cancel(), and may give it to the
 * functions that take a promise. Its result is a Java object of any class; read it with dxfg_v4_Promise_getResult().
 * The handlers added by dxfg_v4_Promise_whenDone() before the promise is done are called by the thread that completes
 * or cancels it, before that function returns.
 *
 * Java: `new Promise()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] promise The pointer to the variable that receives the new promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_Promise_new(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_promise_t **promise);

/**
 * @brief Completes the promise normally, with the given Java object as the result.
 *
 * The first result wins: if the promise is already done (completed, failed or cancelled), the call does nothing and
 * returns success. Otherwise the handlers added by dxfg_v4_Promise_whenDone() are called by the calling thread before
 * the function returns (the handlers of dxfg_v4_Promise_whenDoneAsync() by their executor), and the waiters wake up.
 *
 * The promise must be one whose result may be any object, such as one from dxfg_v4_Promise_new(): the SDK does not
 * check the class of the result, so completing a promise of an event with an object of another class makes
 * dxfg_v4_Promise_EventType_getResult() fail later. For the promises of events use dxfg_v4_Promise_EventType_complete()
 * and dxfg_v4_Promise_List_EventType_complete().
 * A promise of a feed request is completed by the SDK as well: completing it from C races with the SDK, and the first
 * result wins.
 *
 * Java: `Promise.complete(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise; must not be NULL.
 * @param[in] result The handle of the Java object that is the result, or NULL for a null result. It is only read: the
 * caller keeps the ownership of the handle and may release it after the function returns.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution, also when the promise was already done, or
 * #DXFG_EXECUTE_FAIL (-1) on error (e.g., the promise is NULL). Use dxfg_v4_get_and_clear_thread_exception() to
 * determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_Promise_complete(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                 dxfg_v4_java_object_handler *result);

/**
 * @brief Gets the result of the promise as a handle of a Java object.
 *
 * The result is NULL when the promise is not done, completed exceptionally, was cancelled, or completed with a null
 * result: use dxfg_v4_Promise_hasResult() to tell the cases apart. The function works for a promise of any kind; for
 * the promises of events, dxfg_v4_Promise_EventType_getResult() and dxfg_v4_Promise_List_EventType_getResult() give
 * the events themselves.
 *
 * Java: `Promise.getResult()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise; must not be NULL.
 * @param[out] result The pointer to the variable that receives a new handle of the result (release it with
 * dxfg_v4_JavaObjectHandler_release()), or NULL; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g.,
 * the promise is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_Promise_getResult(graal_isolatethread_t *thread, dxfg_v4_promise_t *promise,
                                  DXFG_OUT dxfg_v4_java_object_handler **result);

/**
 * @brief Completes the promise exceptionally, with a new Java exception of the named class and the given message.
 *
 * The SDK creates the exception from a fixed set of classes, without reflection (the native image cannot create a
 * class by its name):
 *  - `java.lang.RuntimeException`
 *  - `java.lang.IllegalArgumentException`
 *  - `java.lang.IllegalStateException`
 *  - `java.lang.UnsupportedOperationException`
 *  - `java.io.IOException`
 *  - `java.util.concurrent.CancellationException`
 *  - `java.util.concurrent.TimeoutException`
 *
 * The names are the fully qualified ones, as in the field `class_name` of dxfg_v4_exception_t. Any other name (another
 * class, a class that is not an exception, or no class at all) gives a `java.lang.RuntimeException` whose message is
 * `className + ": " + message`, or `className` when the message is NULL. A NULL or empty class name gives a
 * `java.lang.RuntimeException` with the message as is. The exception has no cause. A
 * `java.util.concurrent.CancellationException` fails the promise but does not cancel it: dxfg_v4_Promise_isCancelled()
 * gives 0; use dxfg_v4_Promise_cancel() to cancel a promise.
 *
 * The first result wins: if the promise is already done (completed, failed or cancelled), the call does nothing and
 * returns success. Otherwise the handlers are called as for dxfg_v4_Promise_complete(), and the waiters fail
 * (dxfg_v4_Promise_getException() gives the exception itself).
 *
 * Java: `Promise.completeExceptionally(Throwable)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promise The promise; must not be NULL.
 * @param[in] className The fully qualified name of the class of the exception (a UTF-8 null-terminated string), or
 * NULL.
 * @param[in] message The message of the exception (a UTF-8 null-terminated string), or NULL for no message.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution, also when the promise was already done, or
 * #DXFG_EXECUTE_FAIL (-1) on error (e.g., the promise is NULL). Use dxfg_v4_get_and_clear_thread_exception() to
 * determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_Promise_completeExceptionallyWithClassNameAndMessage(graal_isolatethread_t *thread,
                                                                     dxfg_v4_promise_t *promise,
                                                                     const char *className, const char *message);

/**
 * @brief Creates a promise that is already failed, with a new Java exception of the named class and the given
 * message.
 *
 * The exception is created as in dxfg_v4_Promise_completeExceptionallyWithClassNameAndMessage(), from the same set
 * of classes and with the same fallback to `java.lang.RuntimeException`. A handler added to the new promise by
 * dxfg_v4_Promise_whenDone() is called before that function returns, by the calling thread.
 *
 * Java: `Promise.failed(Throwable)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] className The fully qualified name of the class of the exception (a UTF-8 null-terminated string), or
 * NULL.
 * @param[in] message The message of the exception (a UTF-8 null-terminated string), or NULL for no message.
 * @param[out] newPromise The pointer to the variable that receives the new promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_Promise_failedWithClassNameAndMessage(graal_isolatethread_t *thread, const char *className,
                                                      const char *message, DXFG_OUT dxfg_v4_promise_t **newPromise);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promises The list of the promises; it is only read, the caller keeps its ownership.
 * @param[out] promise The pointer to the variable that receives the new promise (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Promises_allOf(graal_isolatethread_t *thread, dxfg_v4_promise_list *promises,
                               DXFG_OUT dxfg_v4_promise_t **promise);

/**
 * Releases a list of promises that a function returned (e.g. dxfg_v4_DXFeed_getLastEventsPromises()): releases every
 * promise handle, the array of the handles and the list itself. It does the same as
 * dxfg_v4_java_object_handler_list_release() with the field `list` of the list.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] promises The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_promise_list_release(graal_isolatethread_t *thread, dxfg_v4_promise_list *promises);

/** @} */ // end of V4_Feed

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_FEED_H
