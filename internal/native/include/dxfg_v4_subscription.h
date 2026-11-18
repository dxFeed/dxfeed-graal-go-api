// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_subscription.h and the v4 names of its types (see the group @ref V4_Subscription).
 */

#ifndef DXFG_V4_SUBSCRIPTION_H
#define DXFG_V4_SUBSCRIPTION_H

#include "dxfg_subscription.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_events.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Subscription Subscription (v4)
 *  @brief The v4 functions and type names of the group @ref Subscription. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_feed_t */
typedef dxfg_feed_t dxfg_v4_feed_t;

/** @copydoc dxfg_subscription_controller_t */
typedef dxfg_subscription_controller_t dxfg_v4_subscription_controller_t;

/** @copydoc dxfg_subscription_t */
typedef dxfg_subscription_t dxfg_v4_subscription_t;

/** @copydoc dxfg_time_series_subscription_t */
typedef dxfg_time_series_subscription_t dxfg_v4_time_series_subscription_t;

/** @copydoc dxfg_feed_event_listener_t */
typedef dxfg_feed_event_listener_t dxfg_v4_feed_event_listener_t;

/** @copydoc dxfg_feed_event_listener_function */
typedef dxfg_feed_event_listener_function dxfg_v4_feed_event_listener_function;

/**
 * Creates a detached subscription for one event type.
 *
 * The subscription is not attached to any feed: attach it with dxfg_v4_DXFeedSubscription_attach(). The event type
 * cannot be changed later.
 *
 * Java: `DXFeedSubscription(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventClazz The event type.
 * @param[out] subscription The pointer to the variable that receives the subscription (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_new(graal_isolatethread_t *thread, dxfg_v4_event_clazz_t eventClazz,
                                       DXFG_OUT dxfg_v4_subscription_t **subscription);

/**
 * Creates a detached subscription for the given event types.
 *
 * The subscription is not attached to any feed: attach it with dxfg_v4_DXFeedSubscription_attach(). The event types
 * cannot be changed later. The list must not be empty.
 *
 * Java: `DXFeedSubscription(Class...)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventClazzes The event types. The list is only read: the caller keeps its ownership.
 * @param[out] subscription The pointer to the variable that receives the subscription (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, for an empty list). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_newWithEventTypes(graal_isolatethread_t *thread,
                                                     dxfg_v4_event_clazz_list *eventClazzes,
                                                     DXFG_OUT dxfg_v4_subscription_t **subscription);

/**
 * Returns the subscription controller for this subscription.
 *
 * The controller manages the same subscription: attaching to feeds, the batch limit, the aggregation period and the
 * executor (see the group @ref SubscriptionController).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] subscriptionController The pointer to the handle of the subscription controller of this subscription.
 * It is written on success only; release the handle with dxfg_v4_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `subscriptionController` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception
 * was thrown.
 *
 * Java: `DXFeedSubscription.getSubscriptionController()`.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedSubscription_getSubscriptionController(
    graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
    DXFG_OUT dxfg_v4_subscription_controller_t **subscriptionController);

/**
 * Attaches the subscription to the specified feed.
 *
 * The attached subscription receives the events of the feed for its symbols. It can be attached to several feeds at
 * the same time, but it cannot be known which feed an event came from.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.attach(DXFeed)`.
 */
int32_t dxfg_v4_DXFeedSubscription_attach(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                          dxfg_v4_feed_t *feed);

/**
 * Detaches the subscription from the specified feed.
 *
 * The subscription keeps its symbols and listeners, but does not receive the events of this feed any more.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.detach(DXFeed)`.
 */
int32_t dxfg_v4_DXFeedSubscription_detach(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                          dxfg_v4_feed_t *feed);

/**
 * Checks whether the subscription is closed (see dxfg_v4_DXFeedSubscription_close()).
 *
 * Java: `DXFeedSubscription.isClosed()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] closed The pointer to the variable that receives 1 if the subscription is closed, 0 if it is not; must
 * not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_isClosed(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                            DXFG_OUT int32_t *closed);

/**
 * Closes the subscription and makes it permanently detached from all feeds.
 *
 * The subscription removes all its event listeners and change listeners and does not accept new ones. Every change
 * listener is notified of the close (its `subscriptionClosed` function, see
 * dxfg_v4_ObservableSubscriptionChangeListener_new()) on the calling thread before the function returns. The handle is
 * not released: release it with dxfg_v4_JavaObjectHandler_release(). Closing a closed subscription does nothing.
 *
 * An event listener may still receive one more batch of events after the function returns. To be sure that its
 * function is not called anymore, release the listener (dxfg_v4_JavaObjectHandler_release()) or deactivate it
 * (dxfg_v4_NativeListener_deactivate()).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.close()`.
 */
int32_t dxfg_v4_DXFeedSubscription_close(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub);

/**
 * Returns the event types of the subscription.
 *
 * Java: `DXFeedSubscription.getEventTypes()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] eventTypes The pointer to the variable that receives the list of the event types (release it with
 * dxfg_v4_event_clazz_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_getEventTypes(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                 DXFG_OUT dxfg_v4_event_clazz_list **eventTypes);

/**
 * Checks whether the subscription contains the event type.
 *
 * Java: `DXFeedSubscription.containsEventType(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] eventClazz The event type.
 * @param[out] contains The pointer to the variable that receives 1 if the subscription contains the event type, 0 if it
 * does not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_containsEventType(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                     dxfg_v4_event_clazz_t eventClazz, DXFG_OUT int32_t *contains);

/**
 * Clears the set of the subscribed symbols.
 *
 * It is the same as setting an empty list of symbols with dxfg_v4_DXFeedSubscription_setSymbols().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.clear()`.
 */
int32_t dxfg_v4_DXFeedSubscription_clear(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub);

/**
 * Returns the subscribed symbols as they were set.
 *
 * For a time series subscription these are the original symbols, without the time of the subscription (see
 * dxfg_v4_DXFeedSubscription_getDecoratedSymbols()). The list is a copy of the symbols at the moment of the call.
 *
 * Java: `DXFeedSubscription.getSymbols()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] symbols The pointer to the variable that receives the list of the symbols (release it with
 * dxfg_v4_symbol_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_getSymbols(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                              DXFG_OUT dxfg_v4_symbol_list **symbols);

/**
 * Returns the subscribed symbols as the change listeners see them.
 *
 * For a time series subscription every symbol is wrapped with the time of the subscription
 * (dxfg_v4_time_series_subscription_symbol_t); for any other subscription the result is the same as of
 * dxfg_v4_DXFeedSubscription_getSymbols(). The list is a copy of the symbols at the moment of the call.
 *
 * Java: `DXFeedSubscription.getDecoratedSymbols()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] symbols The pointer to the variable that receives the list of the symbols (release it with
 * dxfg_v4_symbol_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_getDecoratedSymbols(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                       DXFG_OUT dxfg_v4_symbol_list **symbols);

/**
 * Changes the set of the subscribed symbols so that it contains only the given symbol.
 *
 * The event listeners receive the last events of the symbols that were added. The change listeners are notified of
 * the added and removed symbols on the calling thread before the function returns.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbol The symbol. It is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.setSymbols(Object...)`.
 */
int32_t dxfg_v4_DXFeedSubscription_setSymbol(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                             dxfg_v4_symbol_t *symbol);

/**
 * Changes the set of the subscribed symbols so that it contains only the symbols of the given list.
 *
 * The event listeners receive the last events of the symbols that were added. The change listeners are notified of
 * the added and removed symbols on the calling thread before the function returns. Prefer this function to
 * dxfg_v4_DXFeedSubscription_setSymbol() for several symbols.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbols The symbols. The list is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.setSymbols(Collection)`.
 */
int32_t dxfg_v4_DXFeedSubscription_setSymbols(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                              dxfg_v4_symbol_list *symbols);

/**
 * Adds the symbol to the set of the subscribed symbols.
 *
 * The event listeners receive the last event of the symbol. The change listeners are notified of the added symbol
 * on the calling thread before the function returns.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbol The symbol. It is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.addSymbols(Object)`.
 */
int32_t dxfg_v4_DXFeedSubscription_addSymbol(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                             dxfg_v4_symbol_t *symbol);

/**
 * Adds the symbols of the list to the set of the subscribed symbols.
 *
 * The event listeners receive the last events of the symbols that were added. The change listeners are notified of
 * the added symbols on the calling thread before the function returns. Prefer this function to
 * dxfg_v4_DXFeedSubscription_addSymbol() for several symbols.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbols The symbols. The list is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.addSymbols(Collection)`.
 */
int32_t dxfg_v4_DXFeedSubscription_addSymbols(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                              dxfg_v4_symbol_list *symbols);

/**
 * Removes the symbol from the set of the subscribed symbols.
 *
 * The change listeners are notified of the removed symbol on the calling thread before the function returns.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbol The symbol. It is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.removeSymbols(Object...)`.
 */
int32_t dxfg_v4_DXFeedSubscription_removeSymbol(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                dxfg_v4_symbol_t *symbol);

/**
 * Removes the symbols of the list from the set of the subscribed symbols.
 *
 * The change listeners are notified of the removed symbols on the calling thread before the function returns. Prefer
 * this function to dxfg_v4_DXFeedSubscription_removeSymbol() for several symbols.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbols The symbols. The list is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.removeSymbols(Collection)`.
 */
int32_t dxfg_v4_DXFeedSubscription_removeSymbols(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                 dxfg_v4_symbol_list *symbols);

/**
 * Returns the executor that processes the event notifications of the subscription.
 *
 * Java: `DXFeedSubscription.getExecutor()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] executor The pointer to the variable that receives the executor (release the handle with
 * dxfg_v4_JavaObjectHandler_release()), or NULL if the default executor of the attached feed is used; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_getExecutor(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                               DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * Changes the executor that processes the event notifications of the subscription.
 *
 * The event listeners of the subscription are called by the threads of this executor. The subscription does not
 * take the handle of the executor: release it with dxfg_v4_JavaObjectHandler_release() when it is no longer needed.
 * The executor of a subscription attached to an event queue (dxfg_v4_EventQueue_new(), dxfg_v4_EventQueue_attach())
 * cannot be changed, and it cannot be set on another subscription: the function fails with `IllegalStateException`.
 * The executor of a shared queue (dxfg_v4_EventQueue_getExecutor()) is accepted: the polls of the queue call the
 * listeners, but the subscription is not attached and its events do not reach the lists of the polls (see
 * @ref V4_EventQueue for the limits).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] executor The executor, or NULL to use the default executor of the attached feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.setExecutor(Executor)`.
 */
int32_t dxfg_v4_DXFeedSubscription_setExecutor(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                               dxfg_v4_executor_t *executor);

/**
 * Returns the aggregation period for data of this subscription.
 *
 * Java: `DXFeedSubscription.getAggregationPeriod()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] period The pointer to the variable that receives the aggregation period (release the handle with
 * dxfg_v4_JavaObjectHandler_release()), or NULL if no period was set; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_getAggregationPeriod(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                        DXFG_OUT dxfg_v4_time_period_t **period);

/**
 * Returns the aggregation period (in millis) for data for this subscription instance.
 *
 * The function fails if no aggregation period was set (the Java side throws a `NullPointerException`): use
 * dxfg_v4_DXFeedSubscription_getAggregationPeriod(), which returns NULL in this case.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] aggregationPeriod A pointer to the aggregation period in millis. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `aggregationPeriod` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 *
 * Java: `DXFeedSubscription.getAggregationPeriod()`.
 */
int32_t dxfg_v4_DXFeedSubscription_getAggregationPeriodMillis(graal_isolatethread_t *thread,
                                                              dxfg_v4_subscription_t *sub,
                                                              DXFG_OUT int64_t *aggregationPeriod);

/**
 * Sets the aggregation period for data.
 *
 * The period limits the rate of the notifications of the event listeners. The new period takes effect only after the
 * current one elapses. The effective period is the maximum of this period and the global aggregation period of the
 * endpoint.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] period The aggregation period, or NULL to clear the previous setting. It must not exceed 24 hours.
 * The subscription does not take the handle: release it with dxfg_v4_JavaObjectHandler_release() when it is no longer
 * needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_v4_get_and_clear_thread_exception() to determine if an
 * exception was thrown.
 *
 * Java: `DXFeedSubscription.setAggregationPeriod(TimePeriod)`.
 */
int32_t dxfg_v4_DXFeedSubscription_setAggregationPeriod(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                        dxfg_v4_time_period_t *period);

/**
 * Sets the aggregation period for data.
 * This method sets a new aggregation period for data, which will only take effect on the next iteration of data
 * notification. For example, if the current aggregation period is 5 seconds, and it is changed to 1 second, the next
 * notification may come after up to 5 seconds, after which the new aggregation period will take effect. The
 * effective period is the maximum of this period and the global aggregation period of the endpoint.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] aggregationPeriod The aggregation period in millis. It must not exceed 24 hours (86400000).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_v4_get_and_clear_thread_exception() to determine if an
 * exception was thrown.
 *
 * Java: `DXFeedSubscription.setAggregationPeriod(TimePeriod)`, called with `TimePeriod.valueOf(long)`.
 */
int32_t dxfg_v4_DXFeedSubscription_setAggregationPeriodMillis(graal_isolatethread_t *thread,
                                                              dxfg_v4_subscription_t *sub, int64_t aggregationPeriod);

/**
 * Returns the aggregation period that was applied to the events delivered to this subscription.
 *
 * The value is updated when the events are delivered to the listeners: both bounds are negative (unknown) before the
 * first events arrive.
 *
 * Java: `DXFeedSubscription.getAggregationPeriodInfo()`, `TimePeriodInfo.getMin()`, `TimePeriodInfo.getMax()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] min The minimum aggregation period in millis, negative if unknown. It is written on success only.
 * @param[out] max The maximum aggregation period in millis, negative if unknown. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `min` or `max` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_getAggregationPeriodInfo(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                            DXFG_OUT int64_t *min, DXFG_OUT int64_t *max);

/**
 * Returns the number of the events that were dropped because of a buffer overflow since the subscription was created.
 *
 * The function is experimental: the Java method is marked `@Experimental` and may change.
 *
 * Java: `DXFeedSubscription.getTotalDroppedEvents()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] droppedEvents The number of the dropped events. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `droppedEvents` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedSubscription_getTotalDroppedEvents(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                         DXFG_OUT int64_t *droppedEvents);

/**
 * Creates a new event listener.
 *
 * The listener calls the function with the events that the subscription receives, see
 * dxfg_v4_feed_event_listener_function for the thread of the call and the ownership of the events. Add the listener to
 * a subscription with dxfg_v4_DXFeedSubscription_addEventListener().
 *
 * Java: `DXFeedEventListener` (eventsReceived(List)).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call with the received events.
 * @param[in] user_data The user data to pass to the function; the SDK does not use or free it.
 * @param[out] listener The pointer to the variable that receives the listener (release it with
 * dxfg_v4_JavaObjectHandler_release()). Releasing the last handle of the listener with
 * dxfg_v4_JavaObjectHandler_release() (or dxfg_v4_NativeListener_deactivate()) stops the calls of the function: the
 * user data can be freed after that; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedEventListener_new(graal_isolatethread_t *thread, dxfg_v4_feed_event_listener_function user_func,
                                        void *user_data, DXFG_OUT dxfg_v4_feed_event_listener_t **listener);

/**
 * Adds the event listener to the subscription. Add the listeners before the symbols: QD fails with
 * `IllegalStateException` when a listener is added to an attached subscription that already has symbols.
 *
 * The listener is added to the subscription, not owned by it: release its handle with
 * dxfg_v4_JavaObjectHandler_release() when it is no longer needed. Adding a listener to a closed subscription does
 * nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription.
 * @param[in] listener The listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on success or #DXFG_EXECUTE_FAIL (-1) on error.
 * Use dxfg_v4_get_and_clear_thread_exception() to get the exception.
 *
 * Java: `DXFeedSubscription.addEventListener(DXFeedEventListener)`.
 */
int32_t dxfg_v4_DXFeedSubscription_addEventListener(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                    dxfg_v4_feed_event_listener_t *listener);

/**
 * Removes the event listener from the subscription.
 *
 * The listener may still receive one more batch of events after this function (or dxfg_v4_DXFeedSubscription_close())
 * returns. Release the listener (dxfg_v4_JavaObjectHandler_release()) or deactivate it
 * (dxfg_v4_NativeListener_deactivate()) to be sure that its function is not called anymore.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription.
 * @param[in] listener The listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on success or #DXFG_EXECUTE_FAIL (-1) on error.
 * Use dxfg_v4_get_and_clear_thread_exception() to get the exception.
 *
 * Java: `DXFeedSubscription.removeEventListener(DXFeedEventListener)`.
 */
int32_t dxfg_v4_DXFeedSubscription_removeEventListener(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                       dxfg_v4_feed_event_listener_t *listener);

/**
 * Adds the change listener to the subscription.
 *
 * Create the listener with dxfg_v4_ObservableSubscriptionChangeListener_new(). If the listener was not added yet, is
 * added to an open subscription and the subscription has symbols, the listener is notified immediately of the existing
 * symbols (its `symbolsAdded` function) on the calling thread before the function returns. After that it is notified
 * of every change on the thread that changes the subscription. The calls of the listener for one subscription do not
 * overlap. Adding a listener that is already added, or to a closed subscription, does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.addChangeListener(ObservableSubscriptionChangeListener)`.
 */
int32_t dxfg_v4_DXFeedSubscription_addChangeListener(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                     dxfg_v4_observable_subscription_change_listener_t *listener);

/**
 * Removes the change listener from the subscription.
 *
 * If the listener was added, it is notified of the removal (its `subscriptionClosed` function) on the calling thread
 * before the function returns. Removing a listener that is not added does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.removeChangeListener(ObservableSubscriptionChangeListener)`.
 */
int32_t dxfg_v4_DXFeedSubscription_removeChangeListener(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                        dxfg_v4_observable_subscription_change_listener_t *listener);

/**
 * Returns the maximum number of events in a single notification of the event listeners.
 *
 * It is not a strict guarantee: the notification may exceed the limit. 0 means the optimal limit (the default),
 * INT32_MAX means the maximum limit.
 *
 * Java: `DXFeedSubscription.getEventsBatchLimit()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] eventsBatchLimit The pointer to the variable that receives the events batch limit; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXFeedSubscription_getEventsBatchLimit(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                       DXFG_OUT int32_t *eventsBatchLimit);

/**
 * Sets the maximum number of events in a single notification of the event listeners.
 *
 * It is not a strict guarantee: the notification may exceed the limit. 0 means the optimal limit (the default),
 * INT32_MAX means the maximum limit.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] eventsBatchLimit The notification events limit; it must not be negative.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, for a negative limit). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 *
 * Java: `DXFeedSubscription.setEventsBatchLimit(int)`.
 */
int32_t dxfg_v4_DXFeedSubscription_setEventsBatchLimit(graal_isolatethread_t *thread, dxfg_v4_subscription_t *sub,
                                                       int32_t eventsBatchLimit);

/** @} */ // end of V4_Subscription

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_SUBSCRIPTION_H
