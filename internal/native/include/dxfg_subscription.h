// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Subscription functions and types declarations
 */

#ifndef DXFG_SUBSCRIPTION_H
#define DXFG_SUBSCRIPTION_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"
#include "dxfg_events.h"
#include "dxfg_javac.h"

/** @defgroup Subscription Subscription
 *  @brief Subscriptions to the events of a feed: the event types and the symbols, the listeners of events and of
 *  the changes of the subscription, the executor, the aggregation period and the batch limit.
 *
 *  A subscription (dxfg_subscription_t) is interested in the events whose type is in its set of event types and whose
 *  symbol is in its set of symbols. The event types are set when the subscription is created, the symbols can be
 *  changed at any time. A subscription created by dxfg_DXFeedSubscription_new() or
 *  dxfg_DXFeedSubscription_new2() is detached: it is not attached to any feed and receives no events. Attach it
 *  with dxfg_DXFeedSubscription_attach(). A subscription created by a feed (e.g. dxfg_DXFeed_createSubscription())
 *  is attached to that feed. A time series subscription (dxfg_time_series_subscription_t) is a subscription for
 *  the time series events that is created by a feed with dxfg_DXFeed_createTimeSeriesSubscription().
 *
 *  The events are passed to the event listeners (see dxfg_DXFeedEventListener_new()) by the executor of the
 *  subscription, one batch after another. Add the event listeners before the symbols: the listeners receive the
 *  last events of the symbols that were added. The changes of the subscription are passed to the change listeners
 *  (see dxfg_ObservableSubscriptionChangeListener_new()).
 *
 *  A subscription is closed only by dxfg_DXFeedSubscription_close(): releasing its handle does not close it. A
 *  subscription that stays attached to a feed keeps receiving events and is not freed, so close it when it is not
 *  needed any more.
 *
 *  The ownership rules (handles, strings, lists) are the ones of the group @ref Javac. A symbol or a list that is
 *  passed to a function is only read by the SDK: the caller keeps its ownership. A list that a function returns is
 *  released by the caller with the function given in the description of the function. The value that a function
 *  returns on error is given in its `@return`. dxfg_get_and_clear_thread_exception_t() gives the exception.
 *  @{
 */

/// @cond
// Forward declarations: the types are described in dxfg_feed.h and dxfg_subscription_controller.h.
typedef struct dxfg_feed_t dxfg_feed_t;
typedef struct dxfg_subscription_controller_t dxfg_subscription_controller_t;
/// @endcond

/**
 * @brief A handle of a Java DXFeedSubscription object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_DXFeedSubscription_new(), dxfg_DXFeedSubscription_new2() or with a feed (e.g.
 * dxfg_DXFeed_createSubscription()). Releasing the handle does not close the subscription: call
 * dxfg_DXFeedSubscription_close().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedSubscription.html">Javadoc</a>
 */
typedef struct dxfg_subscription_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_subscription_t;

/**
 * @brief A handle of a Java DXFeedTimeSeriesSubscription object; release it with dxfg_JavaObjectHandler_release().
 *
 * A time series subscription is a subscription (dxfg_subscription_t) for the events that are time series: pass
 * `&sub->sub` to the functions of this header. The symbols that are set are wrapped, for the listeners of the
 * changes, with the time from which the events are subscribed (see dxfg_DXFeedTimeSeriesSubscription_setFromTime()).
 * Create it with dxfg_DXFeed_createTimeSeriesSubscription() or dxfg_DXFeed_createTimeSeriesSubscription2().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedTimeSeriesSubscription.html">Javadoc</a>
 */
typedef struct dxfg_time_series_subscription_t {
    /** The base structure: the fields of dxfg_subscription_t. */
    dxfg_subscription_t sub;
} dxfg_time_series_subscription_t;

/**
 * @brief A handle of a Java DXFeedEventListener object; release it with dxfg_JavaObjectHandler_release().
 *
 * The listener calls a dxfg_feed_event_listener_function with the received events. Create it with
 * dxfg_DXFeedEventListener_new(). The last handle of the listener that is released deactivates the listener: its
 * function is not called anymore.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedEventListener.html">Javadoc</a>
 */
typedef struct dxfg_feed_event_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_feed_event_listener_t;

/**
 * @brief The function of a dxfg_feed_event_listener_t: called when events are received.
 *
 * The function is called by a thread of the executor of the subscription (the default executor of the feed's
 * endpoint, or the executor set by dxfg_DXFeedSubscription_setExecutor()), not by the thread that added the
 * listener. The calls for one subscription do not overlap: the next batch of events is passed after the function
 * returns for the previous one. The number of events in a batch is limited by the events batch limit (see
 * dxfg_DXFeedSubscription_setEventsBatchLimit()).
 *
 * The list and the events belong to the SDK: it frees them (with their strings) after the function returns. Do not
 * release them, copy what you need to keep.
 *
 * Java: `DXFeedEventListener.eventsReceived(List)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] events The received events.
 * @param[in] user_data The pointer that was passed to dxfg_DXFeedEventListener_new(); the SDK does not use it.
 */
typedef void (*dxfg_feed_event_listener_function)(graal_isolatethread_t *thread, dxfg_event_type_list *events, void *user_data);

/**
 * Creates a detached subscription for one event type.
 *
 * The subscription is not attached to any feed: attach it with dxfg_DXFeedSubscription_attach(). The event type
 * cannot be changed later.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventClazz The event type.
 * @return The subscription (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXFeedSubscription(Class)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_new()")
dxfg_subscription_t*      dxfg_DXFeedSubscription_new(graal_isolatethread_t *thread, dxfg_event_clazz_t eventClazz);

/**
 * Creates a detached subscription for the given event types.
 *
 * The subscription is not attached to any feed: attach it with dxfg_DXFeedSubscription_attach(). The event types
 * cannot be changed later. The list must not be empty.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_newWithEventTypes().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventClazzes The event types. The list is only read: the caller keeps its ownership.
 * @return The subscription (release it with dxfg_JavaObjectHandler_release()), or NULL on error (for example, for
 * an empty list). Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXFeedSubscription(Class...)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_newWithEventTypes()")
dxfg_subscription_t*      dxfg_DXFeedSubscription_new2(graal_isolatethread_t *thread, dxfg_event_clazz_list_t *eventClazzes);

/**
 * Returns the subscription controller for this subscription.
 *
 * The controller manages the same subscription: attaching to feeds, the batch limit, the aggregation period and the
 * executor (see the group @ref SubscriptionController).
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getSubscriptionController() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] subscriptionController The pointer to the handle of the subscription controller of this subscription.
 * It is written on success only; release the handle with dxfg_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `subscriptionController` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception
 * was thrown.
 *
 * Java: `DXFeedSubscription.getSubscriptionController()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getSubscriptionController() (experimental)")
int32_t                   dxfg_DXFeedSubscription_getSubscriptionController(graal_isolatethread_t *thread, dxfg_subscription_t *sub, DXFG_OUT dxfg_subscription_controller_t **subscriptionController);

/**
 * Attaches the subscription to the specified feed.
 *
 * The attached subscription receives the events of the feed for its symbols. It can be attached to several feeds at
 * the same time, but it cannot be known which feed an event came from.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_attach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.attach(DXFeed)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_attach()")
int32_t                   dxfg_DXFeedSubscription_attach(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_feed_t *feed);

/**
 * Detaches the subscription from the specified feed.
 *
 * The subscription keeps its symbols and listeners, but does not receive the events of this feed any more.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_detach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.detach(DXFeed)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_detach()")
int32_t                   dxfg_DXFeedSubscription_detach(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_feed_t *feed);

/**
 * Checks whether the subscription is closed (see dxfg_DXFeedSubscription_close()).
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_isClosed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return 1 if the subscription is closed, 0 if it is not, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.isClosed()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_isClosed()")
int32_t                   dxfg_DXFeedSubscription_isClosed(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Closes the subscription and makes it permanently detached from all feeds.
 *
 * The subscription removes all its event listeners and change listeners and does not accept new ones. Every change
 * listener is notified of the close (its `subscriptionClosed` function, see
 * dxfg_ObservableSubscriptionChangeListener_new()) on the calling thread before the function returns. The handle is
 * not released: release it with dxfg_JavaObjectHandler_release(). Closing a closed subscription does nothing.
 *
 * An event listener may still receive one more batch of events after the function returns. To be sure that its
 * function is not called anymore, release the listener (dxfg_JavaObjectHandler_release()) or deactivate it
 * (dxfg_NativeListener_deactivate()).
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.close()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_close()")
int32_t                   dxfg_DXFeedSubscription_close(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Returns the event types of the subscription.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getEventTypes().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return The list of the event types (release it with dxfg_CList_EventClazz_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXFeedSubscription.getEventTypes()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getEventTypes()")
dxfg_event_clazz_list_t*  dxfg_DXFeedSubscription_getEventTypes(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Checks whether the subscription contains the event type.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_containsEventType().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] eventClazz The event type.
 * @return 1 if the subscription contains the event type, 0 if it does not, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.containsEventType(Class)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_containsEventType()")
int32_t                   dxfg_DXFeedSubscription_containsEventType(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_event_clazz_t eventClazz);

/**
 * Clears the set of the subscribed symbols.
 *
 * It is the same as setting an empty list of symbols with dxfg_DXFeedSubscription_setSymbols().
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_clear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.clear()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_clear()")
int32_t                   dxfg_DXFeedSubscription_clear(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Returns the subscribed symbols as they were set.
 *
 * For a time series subscription these are the original symbols, without the time of the subscription (see
 * dxfg_DXFeedSubscription_getDecoratedSymbols()). The list is a copy of the symbols at the moment of the call.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getSymbols().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return The list of the symbols (release it with dxfg_CList_symbol_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXFeedSubscription.getSymbols()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getSymbols()")
dxfg_symbol_list*         dxfg_DXFeedSubscription_getSymbols(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Returns the subscribed symbols as the change listeners see them.
 *
 * For a time series subscription every symbol is wrapped with the time of the subscription
 * (dxfg_time_series_subscription_symbol_t); for any other subscription the result is the same as of
 * dxfg_DXFeedSubscription_getSymbols(). The list is a copy of the symbols at the moment of the call.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getDecoratedSymbols().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return The list of the symbols (release it with dxfg_CList_symbol_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXFeedSubscription.getDecoratedSymbols()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getDecoratedSymbols()")
dxfg_symbol_list*         dxfg_DXFeedSubscription_getDecoratedSymbols(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Changes the set of the subscribed symbols so that it contains only the given symbol.
 *
 * The event listeners receive the last events of the symbols that were added. The change listeners are notified of
 * the added and removed symbols on the calling thread before the function returns.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_setSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbol The symbol. It is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.setSymbols(Object...)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_setSymbol()")
int32_t                   dxfg_DXFeedSubscription_setSymbol(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_symbol_t *symbol);

/**
 * Changes the set of the subscribed symbols so that it contains only the symbols of the given list.
 *
 * The event listeners receive the last events of the symbols that were added. The change listeners are notified of
 * the added and removed symbols on the calling thread before the function returns. Prefer this function to
 * dxfg_DXFeedSubscription_setSymbol() for several symbols.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_setSymbols().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbols The symbols. The list is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.setSymbols(Collection)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_setSymbols()")
int32_t                   dxfg_DXFeedSubscription_setSymbols(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_symbol_list *symbols);

/**
 * Adds the symbol to the set of the subscribed symbols.
 *
 * The event listeners receive the last event of the symbol. The change listeners are notified of the added symbol
 * on the calling thread before the function returns.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_addSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbol The symbol. It is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.addSymbols(Object)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_addSymbol()")
int32_t                   dxfg_DXFeedSubscription_addSymbol(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_symbol_t *symbol);

/**
 * Adds the symbols of the list to the set of the subscribed symbols.
 *
 * The event listeners receive the last events of the symbols that were added. The change listeners are notified of
 * the added symbols on the calling thread before the function returns. Prefer this function to
 * dxfg_DXFeedSubscription_addSymbol() for several symbols.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_addSymbols().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbols The symbols. The list is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.addSymbols(Collection)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_addSymbols()")
int32_t                   dxfg_DXFeedSubscription_addSymbols(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_symbol_list *symbols);

/**
 * Removes the symbol from the set of the subscribed symbols.
 *
 * The change listeners are notified of the removed symbol on the calling thread before the function returns.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_removeSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbol The symbol. It is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.removeSymbols(Object...)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_removeSymbol()")
int32_t                   dxfg_DXFeedSubscription_removeSymbol(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_symbol_t *symbol);

/**
 * Removes the symbols of the list from the set of the subscribed symbols.
 *
 * The change listeners are notified of the removed symbols on the calling thread before the function returns. Prefer
 * this function to dxfg_DXFeedSubscription_removeSymbol() for several symbols.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_removeSymbols().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] symbols The symbols. The list is only read: the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.removeSymbols(Collection)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_removeSymbols()")
int32_t                   dxfg_DXFeedSubscription_removeSymbols(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_symbol_list *symbols);

/**
 * Returns the executor that processes the event notifications of the subscription.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return The executor (release the handle with dxfg_JavaObjectHandler_release()), or NULL if the default executor of
 * the attached feed is used. NULL is also returned on error: use dxfg_get_and_clear_thread_exception_t() to get the
 * exception and to tell the two cases apart.
 *
 * Java: `DXFeedSubscription.getExecutor()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getExecutor()")
dxfg_executor_t*          dxfg_DXFeedSubscription_getExecutor(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Changes the executor that processes the event notifications of the subscription.
 *
 * The event listeners of the subscription are called by the threads of this executor. The subscription does not
 * take the handle of the executor: release it with dxfg_JavaObjectHandler_release() when it is no longer needed.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_setExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] executor The executor, or NULL to use the default executor of the attached feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.setExecutor(Executor)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_setExecutor()")
int32_t                   dxfg_DXFeedSubscription_setExecutor(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_executor_t *executor);

/**
 * Returns the aggregation period for data of this subscription.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getAggregationPeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return The aggregation period (release the handle with dxfg_JavaObjectHandler_release()), or NULL if no period was
 * set. NULL is also returned on error: use dxfg_get_and_clear_thread_exception_t() to get the exception and to tell
 * the two cases apart.
 *
 * Java: `DXFeedSubscription.getAggregationPeriod()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getAggregationPeriod()")
dxfg_time_period_t*       dxfg_DXFeedSubscription_getAggregationPeriod(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Returns the aggregation period (in millis) for data for this subscription instance.
 *
 * The function fails if no aggregation period was set (the Java side throws a `NullPointerException`): use
 * dxfg_DXFeedSubscription_getAggregationPeriod(), which returns NULL in this case.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getAggregationPeriodMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] aggregationPeriod A pointer to the aggregation period in millis. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `aggregationPeriod` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 *
 * Java: `DXFeedSubscription.getAggregationPeriod()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getAggregationPeriodMillis()")
int32_t  dxfg_DXFeedSubscription_getAggregationPeriodMillis(graal_isolatethread_t *thread, dxfg_subscription_t *sub, DXFG_OUT int64_t *aggregationPeriod);

/**
 * Sets the aggregation period for data.
 *
 * The period limits the rate of the notifications of the event listeners. The new period takes effect only after the
 * current one elapses. The effective period is the maximum of this period and the global aggregation period of the
 * endpoint.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_setAggregationPeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] period The aggregation period, or NULL to clear the previous setting. It must not exceed 24 hours.
 * The subscription does not take the handle: release it with dxfg_JavaObjectHandler_release() when it is no longer
 * needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 *
 * Java: `DXFeedSubscription.setAggregationPeriod(TimePeriod)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_setAggregationPeriod()")
int32_t                   dxfg_DXFeedSubscription_setAggregationPeriod(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_time_period_t *period);

/**
 * Sets the aggregation period for data.
 * This method sets a new aggregation period for data, which will only take effect on the next iteration of data
 * notification. For example, if the current aggregation period is 5 seconds, and it is changed to 1 second, the next
 * notification may come after up to 5 seconds, after which the new aggregation period will take effect. The
 * effective period is the maximum of this period and the global aggregation period of the endpoint.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_setAggregationPeriodMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] aggregationPeriod The aggregation period in millis. It must not exceed 24 hours (86400000).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 *
 * Java: `DXFeedSubscription.setAggregationPeriod(TimePeriod)`, called with `TimePeriod.valueOf(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_setAggregationPeriodMillis()")
int32_t dxfg_DXFeedSubscription_setAggregationPeriodMillis(graal_isolatethread_t *thread, dxfg_subscription_t *sub, int64_t aggregationPeriod);

/**
 * Returns the aggregation period that was applied to the events delivered to this subscription.
 *
 * The value is updated when the events are delivered to the listeners: both bounds are negative (unknown) before the
 * first events arrive.
 *
 * Java: `DXFeedSubscription.getAggregationPeriodInfo()`, `TimePeriodInfo.getMin()`, `TimePeriodInfo.getMax()`.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getAggregationPeriodInfo().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] min The minimum aggregation period in millis, negative if unknown. It is written on success only.
 * @param[out] max The maximum aggregation period in millis, negative if unknown. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `min` or `max` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getAggregationPeriodInfo()")
int32_t dxfg_DXFeedSubscription_getAggregationPeriodInfo(graal_isolatethread_t *thread, dxfg_subscription_t *sub,
                                                         DXFG_OUT int64_t *min, DXFG_OUT int64_t *max);

/**
 * Returns the number of the events that were dropped because of a buffer overflow since the subscription was created.
 *
 * The function is experimental: the Java method is marked `@Experimental` and may change.
 *
 * Java: `DXFeedSubscription.getTotalDroppedEvents()`.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getTotalDroppedEvents() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[out] droppedEvents The number of the dropped events. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `droppedEvents` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getTotalDroppedEvents() (experimental)")
int32_t dxfg_DXFeedSubscription_getTotalDroppedEvents(graal_isolatethread_t *thread, dxfg_subscription_t *sub,
                                                      DXFG_OUT int64_t *droppedEvents);

/**
 * Creates a new event listener.
 *
 * The listener calls the function with the events that the subscription receives, see
 * dxfg_feed_event_listener_function for the thread of the call and the ownership of the events. Add the listener to
 * a subscription with dxfg_DXFeedSubscription_addEventListener().
 *
 * @deprecated Use dxfg_v4_DXFeedEventListener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call with the received events.
 * @param[in] user_data The user data to pass to the function; the SDK does not use or free it.
 * @return The listener (release it with dxfg_JavaObjectHandler_release()) or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 * Releasing the last handle of the listener with dxfg_JavaObjectHandler_release() (or dxfg_NativeListener_deactivate())
 * stops the calls of the function: the user data can be freed after that.
 *
 * Java: `DXFeedEventListener` (eventsReceived(List)).
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedEventListener_new()")
dxfg_feed_event_listener_t* dxfg_DXFeedEventListener_new(graal_isolatethread_t *thread, dxfg_feed_event_listener_function user_func, void *user_data);

/**
 * Adds the event listener to the subscription. Add the listeners before the symbols: QD fails with
 * `IllegalStateException` when a listener is added to an attached subscription that already has symbols.
 *
 * The listener is added to the subscription, not owned by it: release its handle with
 * dxfg_JavaObjectHandler_release() when it is no longer needed. Adding a listener to a closed subscription does
 * nothing.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_addEventListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription.
 * @param[in] listener The listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on success or #DXFG_EXECUTE_FAIL (-1) on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXFeedSubscription.addEventListener(DXFeedEventListener)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_addEventListener()")
int32_t                   dxfg_DXFeedSubscription_addEventListener(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_feed_event_listener_t *listener);

/**
 * Removes the event listener from the subscription.
 *
 * The listener may still receive one more batch of events after this function (or dxfg_DXFeedSubscription_close())
 * returns. Release the listener (dxfg_JavaObjectHandler_release()) or deactivate it (dxfg_NativeListener_deactivate())
 * to be sure that its function is not called anymore.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_removeEventListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription.
 * @param[in] listener The listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on success or #DXFG_EXECUTE_FAIL (-1) on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXFeedSubscription.removeEventListener(DXFeedEventListener)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_removeEventListener()")
int32_t                   dxfg_DXFeedSubscription_removeEventListener(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_feed_event_listener_t *listener);

/**
 * Adds the change listener to the subscription.
 *
 * Create the listener with dxfg_ObservableSubscriptionChangeListener_new(). If the listener was not added yet, is added
 * to an open subscription and the subscription has symbols, the listener is notified immediately of the existing
 * symbols (its `symbolsAdded` function) on the calling thread before the function returns. After that it is notified
 * of every change on the thread that changes the subscription. The calls of the listener for one subscription do not
 * overlap. Adding a listener that is already added, or to a closed subscription, does nothing.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_addChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.addChangeListener(ObservableSubscriptionChangeListener)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_addChangeListener()")
int32_t                   dxfg_DXFeedSubscription_addChangeListener(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_observable_subscription_change_listener_t *listener);

/**
 * Removes the change listener from the subscription.
 *
 * If the listener was added, it is notified of the removal (its `subscriptionClosed` function) on the calling thread
 * before the function returns. Removing a listener that is not added does nothing.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_removeChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.removeChangeListener(ObservableSubscriptionChangeListener)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_removeChangeListener()")
int32_t                   dxfg_DXFeedSubscription_removeChangeListener(graal_isolatethread_t *thread, dxfg_subscription_t *sub, dxfg_observable_subscription_change_listener_t *listener);

/**
 * Returns the maximum number of events in a single notification of the event listeners.
 *
 * It is not a strict guarantee: the notification may exceed the limit. 0 means the optimal limit (the default),
 * INT32_MAX means the maximum limit.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_getEventsBatchLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @return The events batch limit, or -1 on error (a valid limit is never negative). Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXFeedSubscription.getEventsBatchLimit()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_getEventsBatchLimit()")
int32_t                   dxfg_DXFeedSubscription_getEventsBatchLimit(graal_isolatethread_t *thread, dxfg_subscription_t *sub);

/**
 * Sets the maximum number of events in a single notification of the event listeners.
 *
 * It is not a strict guarantee: the notification may exceed the limit. 0 means the optimal limit (the default),
 * INT32_MAX means the maximum limit.
 *
 * @deprecated Use dxfg_v4_DXFeedSubscription_setEventsBatchLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The subscription handle.
 * @param[in] eventsBatchLimit The notification events limit; it must not be negative.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, for a negative limit). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 *
 * Java: `DXFeedSubscription.setEventsBatchLimit(int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXFeedSubscription_setEventsBatchLimit()")
int32_t                   dxfg_DXFeedSubscription_setEventsBatchLimit(graal_isolatethread_t *thread, dxfg_subscription_t *sub, int32_t eventsBatchLimit);

/** @} */ // end of Subscription

#ifdef __cplusplus
}
#endif

#endif // DXFG_SUBSCRIPTION_H