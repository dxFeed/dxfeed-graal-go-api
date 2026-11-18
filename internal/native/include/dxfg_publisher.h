// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Publisher functions and types declarations
 */

#ifndef DXFG_PUBLISHER_H
#define DXFG_PUBLISHER_H

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

/** @defgroup Publisher Publisher
 *  @brief Publishing of events to the local and remote feeds, and the observable subscriptions that tell which
 *  symbols the feeds are subscribed to.
 *
 *  The publisher (dxfg_publisher_t) publishes events to the feeds that are connected to it locally or over
 *  the network. The observable subscription of the publisher (dxfg_observable_subscription_t) is the set of symbols
 *  that the feeds are subscribed to for one event type: add a change listener
 *  (dxfg_ObservableSubscriptionChangeListener_new()) to know when to publish the events.
 *
 *  The ownership rules (handles, lists) are the ones of the group @ref Javac. A list that is passed to a function is
 *  only read by the SDK: the caller keeps its ownership. A list that a function returns is released by the caller
 *  with the function given in the description of the function. The value that a function returns on error is given
 *  in its `@return`. dxfg_get_and_clear_thread_exception_t() gives the exception.
 *  @{
 */

/// @cond
// Forward declaration: the type is described in dxfg_javac.h.
typedef struct dxfg_java_object_handler dxfg_java_object_handler;
/// @endcond

/**
 * @brief A handle of a Java DXPublisher object; release it with dxfg_JavaObjectHandler_release().
 *
 * Get it with dxfg_DXPublisher_getInstance() or dxfg_DXEndpoint_getPublisher(). Releasing the handle does not close
 * the endpoint of the publisher.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXPublisher.html">Javadoc</a>
 */
typedef struct dxfg_publisher_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_publisher_t;

/**
 * @brief A handle of a Java ObservableSubscription object; release it with dxfg_JavaObjectHandler_release().
 *
 * It is the observable set of the symbols that the feeds are subscribed to, for one event type. Get it with
 * dxfg_DXPublisher_getSubscription(). The subscription is closed when the endpoint of the publisher is closed.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/ObservableSubscription.html">Javadoc</a>
 */
typedef struct dxfg_observable_subscription_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_observable_subscription_t;

/**
 * Returns the default application-wide singleton publisher
 * (<a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXPublisher.html#getInstance--">
 * DXPublisher.getInstance()</a>): the publisher of the default endpoint with the publisher role. It is configured by
 * the system properties (see dxfg_system_set_property()) when it is first used.
 *
 * All callers share this publisher and its endpoint: closing the endpoint closes it for everyone in the process.
 *
 * @deprecated Use dxfg_v4_DXPublisher_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The publisher (release the handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXPublisher.getInstance()`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXPublisher_getInstance()")
dxfg_publisher_t*                   dxfg_DXPublisher_getInstance(graal_isolatethread_t *thread);

/**
 * Publishes the events to the corresponding feeds.
 *
 * If the endpoint of the publisher has the publisher role and is connected, the events are delivered to the remote
 * endpoints. The local feed of the endpoint always receives them. A feed receives an event only if it is subscribed
 * to its symbol (or to all symbols with a wildcard, or, for a time series event, with a time series subscription
 * for the time of the event). The published events are not stored and are lost if there is no subscription; the last
 * events of the lasting event types are kept while the subscription to their symbol exists.
 *
 * The events are serialized by the function: the list and the events are only read, the caller keeps its ownership
 * and can reuse or release them after the function returns.
 *
 * @deprecated Use dxfg_v4_DXPublisher_publishEvents().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] publisher The publisher handle.
 * @param[in] events The events to publish.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `DXPublisher.publishEvents(Collection)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXPublisher_publishEvents()")
int32_t                             dxfg_DXPublisher_publishEvents(graal_isolatethread_t *thread, dxfg_publisher_t *publisher, dxfg_event_type_list *events);

/**
 * Returns the observable set of the subscribed symbols for the event type.
 *
 * The symbols of the set are the symbols that the feeds are subscribed to; the set contains the wildcard symbol if
 * there is a subscription to all symbols, and, for a time series subscription, the symbols wrapped with the time of
 * the subscription (dxfg_time_series_subscription_symbol_t). The set may notify its change listeners of the same
 * symbols several times without a removal between the notifications: it happens when the subscription disappears,
 * the cached data is lost and the subscription appears again. On each such notification the data provider is to
 * publish the most recent events of the symbols (dxfg_DXPublisher_publishEvents()).
 *
 * @deprecated Use dxfg_v4_DXPublisher_getSubscription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] publisher The publisher handle.
 * @param[in] eventClazz The event type.
 * @return The observable subscription (release the handle with dxfg_JavaObjectHandler_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `DXPublisher.getSubscription(Class)`.
 */
DXFG_DEPRECATED("use dxfg_v4_DXPublisher_getSubscription()")
dxfg_observable_subscription_t*     dxfg_DXPublisher_getSubscription(graal_isolatethread_t *thread, dxfg_publisher_t *publisher, dxfg_event_clazz_t eventClazz);


/**
 * Checks whether the observable subscription is closed.
 *
 * The subscription of a publisher is closed when its endpoint is closed.
 *
 * @deprecated Use dxfg_v4_ObservableSubscription_isClosed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @return 1 if the subscription is closed, 0 if it is not, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `ObservableSubscription.isClosed()`.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableSubscription_isClosed()")
int32_t                    dxfg_ObservableSubscription_isClosed(graal_isolatethread_t *thread, dxfg_observable_subscription_t *sub);

/**
 * Returns the event types of the observable subscription.
 *
 * The subscription of a publisher (dxfg_DXPublisher_getSubscription()) has the single event type that was requested.
 *
 * @deprecated Use dxfg_v4_ObservableSubscription_getEventTypes().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @return The list of the event types (release it with dxfg_CList_EventClazz_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `ObservableSubscription.getEventTypes()`.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableSubscription_getEventTypes()")
dxfg_event_clazz_list_t*   dxfg_ObservableSubscription_getEventTypes(graal_isolatethread_t *thread, dxfg_observable_subscription_t *sub);

/**
 * Checks whether the observable subscription contains the event type.
 *
 * @deprecated Use dxfg_v4_ObservableSubscription_containsEventType().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[in] eventClazz The event type.
 * @return 1 if the subscription contains the event type, 0 if it does not, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `ObservableSubscription.containsEventType(Class)`.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableSubscription_containsEventType()")
int32_t                    dxfg_ObservableSubscription_containsEventType(graal_isolatethread_t *thread, dxfg_observable_subscription_t *sub, dxfg_event_clazz_t eventClazz);

/**
 * Adds the change listener to the observable subscription.
 *
 * Create the listener with dxfg_ObservableSubscriptionChangeListener_new(); see its description for the callbacks of
 * the listener: they are called on a Java thread, not necessarily the one that added the listener, and the symbols
 * belong to the SDK. If the listener is added to an open subscription, it is notified immediately of the existing
 * subscription state (its `symbolsAdded` function) before the function returns, and after that of every change. The
 * calls of the listener for one subscription do not overlap. Adding a listener that is already added, or to a closed
 * subscription, does nothing.
 *
 * @deprecated Use dxfg_v4_ObservableSubscription_addChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `ObservableSubscription.addChangeListener(ObservableSubscriptionChangeListener)`.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableSubscription_addChangeListener()")
int32_t                    dxfg_ObservableSubscription_addChangeListener(graal_isolatethread_t *thread, dxfg_observable_subscription_t *sub, dxfg_observable_subscription_change_listener_t *listener);

/**
 * Removes the change listener from the observable subscription.
 *
 * If the listener was added, it is notified of the removal (its `subscriptionClosed` function) before the function
 * returns. Removing a listener that is not added does nothing. To be sure that the functions of the listener are not
 * called anymore, release it (dxfg_JavaObjectHandler_release()) or deactivate it (dxfg_NativeListener_deactivate()).
 *
 * @deprecated Use dxfg_v4_ObservableSubscription_removeChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `ObservableSubscription.removeChangeListener(ObservableSubscriptionChangeListener)`.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableSubscription_removeChangeListener()")
int32_t                    dxfg_ObservableSubscription_removeChangeListener(graal_isolatethread_t *thread, dxfg_observable_subscription_t *sub, dxfg_observable_subscription_change_listener_t *listener);

/** @} */ // end of Publisher

#ifdef __cplusplus
}
#endif

#endif // DXFG_PUBLISHER_H