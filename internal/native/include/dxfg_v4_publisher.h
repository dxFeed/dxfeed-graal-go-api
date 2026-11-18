// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_publisher.h and the v4 names of its types (see the group @ref V4_Publisher).
 */

#ifndef DXFG_V4_PUBLISHER_H
#define DXFG_V4_PUBLISHER_H

#include "dxfg_publisher.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_events.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Publisher Publisher (v4)
 *  @brief The v4 functions and type names of the group @ref Publisher. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_java_object_handler */
typedef dxfg_java_object_handler dxfg_v4_java_object_handler;

/** @copydoc dxfg_publisher_t */
typedef dxfg_publisher_t dxfg_v4_publisher_t;

/** @copydoc dxfg_observable_subscription_t */
typedef dxfg_observable_subscription_t dxfg_v4_observable_subscription_t;

/**
 * Returns the default application-wide singleton publisher
 * (<a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXPublisher.html#getInstance--">
 * DXPublisher.getInstance()</a>): the publisher of the default endpoint with the publisher role. It is configured by
 * the system properties (see dxfg_system_set_property()) when it is first used.
 *
 * All callers share this publisher and its endpoint: closing the endpoint closes it for everyone in the process.
 *
 * Java: `DXPublisher.getInstance()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] publisher The pointer to the variable that receives the publisher (release the handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXPublisher_getInstance(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_publisher_t **publisher);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] publisher The publisher handle.
 * @param[in] events The events to publish.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `DXPublisher.publishEvents(Collection)`.
 */
int32_t dxfg_v4_DXPublisher_publishEvents(graal_isolatethread_t *thread, dxfg_v4_publisher_t *publisher,
                                          dxfg_v4_event_type_list *events);

/**
 * Returns the observable set of the subscribed symbols for the event type.
 *
 * The symbols of the set are the symbols that the feeds are subscribed to; the set contains the wildcard symbol if
 * there is a subscription to all symbols, and, for a time series subscription, the symbols wrapped with the time of
 * the subscription (dxfg_v4_time_series_subscription_symbol_t). The set may notify its change listeners of the same
 * symbols several times without a removal between the notifications: it happens when the subscription disappears,
 * the cached data is lost and the subscription appears again. On each such notification the data provider is to
 * publish the most recent events of the symbols (dxfg_v4_DXPublisher_publishEvents()).
 *
 * Java: `DXPublisher.getSubscription(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] publisher The publisher handle.
 * @param[in] eventClazz The event type.
 * @param[out] subscription The pointer to the variable that receives the observable subscription (release the handle
 * with dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXPublisher_getSubscription(graal_isolatethread_t *thread, dxfg_v4_publisher_t *publisher,
                                            dxfg_v4_event_clazz_t eventClazz,
                                            DXFG_OUT dxfg_v4_observable_subscription_t **subscription);

/**
 * Checks whether the observable subscription is closed.
 *
 * The subscription of a publisher is closed when its endpoint is closed.
 *
 * Java: `ObservableSubscription.isClosed()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[out] closed The pointer to the variable that receives 1 if the subscription is closed, 0 if it is not; must
 * not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableSubscription_isClosed(graal_isolatethread_t *thread, dxfg_v4_observable_subscription_t *sub,
                                                DXFG_OUT int32_t *closed);

/**
 * Returns the event types of the observable subscription.
 *
 * The subscription of a publisher (dxfg_v4_DXPublisher_getSubscription()) has the single event type that was requested.
 *
 * Java: `ObservableSubscription.getEventTypes()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[out] eventTypes The pointer to the variable that receives the list of the event types (release it with
 * dxfg_v4_event_clazz_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableSubscription_getEventTypes(graal_isolatethread_t *thread,
                                                     dxfg_v4_observable_subscription_t *sub,
                                                     DXFG_OUT dxfg_v4_event_clazz_list **eventTypes);

/**
 * Checks whether the observable subscription contains the event type.
 *
 * Java: `ObservableSubscription.containsEventType(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[in] eventClazz The event type.
 * @param[out] contains The pointer to the variable that receives 1 if the subscription contains the event type, 0 if it
 * does not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableSubscription_containsEventType(graal_isolatethread_t *thread,
                                                         dxfg_v4_observable_subscription_t *sub,
                                                         dxfg_v4_event_clazz_t eventClazz, DXFG_OUT int32_t *contains);

/**
 * Adds the change listener to the observable subscription.
 *
 * Create the listener with dxfg_v4_ObservableSubscriptionChangeListener_new(); see its description for the callbacks of
 * the listener: they are called on a Java thread, not necessarily the one that added the listener, and the symbols
 * belong to the SDK. If the listener is added to an open subscription, it is notified immediately of the existing
 * subscription state (its `symbolsAdded` function) before the function returns, and after that of every change. The
 * calls of the listener for one subscription do not overlap. Adding a listener that is already added, or to a closed
 * subscription, does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `ObservableSubscription.addChangeListener(ObservableSubscriptionChangeListener)`.
 */
int32_t dxfg_v4_ObservableSubscription_addChangeListener(graal_isolatethread_t *thread,
                                                         dxfg_v4_observable_subscription_t *sub,
                                                         dxfg_v4_observable_subscription_change_listener_t *listener);

/**
 * Removes the change listener from the observable subscription.
 *
 * If the listener was added, it is notified of the removal (its `subscriptionClosed` function) before the function
 * returns. Removing a listener that is not added does nothing. To be sure that the functions of the listener are not
 * called anymore, release it (dxfg_v4_JavaObjectHandler_release()) or deactivate it
 * (dxfg_v4_NativeListener_deactivate()).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sub The observable subscription handle.
 * @param[in] listener The change listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `ObservableSubscription.removeChangeListener(ObservableSubscriptionChangeListener)`.
 */
int32_t dxfg_v4_ObservableSubscription_removeChangeListener(
    graal_isolatethread_t *thread, dxfg_v4_observable_subscription_t *sub,
    dxfg_v4_observable_subscription_change_listener_t *listener);

/** @} */ // end of V4_Publisher

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_PUBLISHER_H
