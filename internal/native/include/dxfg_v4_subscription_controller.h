// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_subscription_controller.h and the v4 names of its types (see the group @ref V4_SubscriptionController).
 */

#ifndef DXFG_V4_SUBSCRIPTION_CONTROLLER_H
#define DXFG_V4_SUBSCRIPTION_CONTROLLER_H

#include "dxfg_subscription_controller.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_SubscriptionController SubscriptionController (v4)
 *  @brief The v4 functions and type names of the group @ref SubscriptionController. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_feed_t */
typedef dxfg_feed_t dxfg_v4_feed_t;

/** @copydoc dxfg_subscription_controller_t */
typedef dxfg_subscription_controller_t dxfg_v4_subscription_controller_t;

/**
 * Attaches the underlying subscription to the specified feed.
 *
 * The attached subscription receives the events of the feed. It can be attached to several feeds at the same time,
 * but this is rarely necessary and is not recommended.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.attach(DXFeed)`.
 *
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/SubscriptionController.html#attach-com.dxfeed.api.DXFeed-
 */
int32_t dxfg_v4_SubscriptionController_attach(graal_isolatethread_t *thread,
                                              dxfg_v4_subscription_controller_t *subscriptionController,
                                              dxfg_v4_feed_t *feed);

/**
 * Detaches the underlying subscription from the specified feed.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.detach(DXFeed)`.
 *
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/SubscriptionController.html#detach-com.dxfeed.api.DXFeed-
 */
int32_t dxfg_v4_SubscriptionController_detach(graal_isolatethread_t *thread,
                                              dxfg_v4_subscription_controller_t *subscriptionController,
                                              dxfg_v4_feed_t *feed);

/**
 * Returns maximum number of events in the single notification of the listener.
 * This is not a strict guarantee, some implementations may exceed the limit.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] eventsBatchLimit The maximum number of events in the single notification of the listener: 0 for the
 * optimal limit (the default), INT32_MAX for the maximum limit. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `eventsBatchLimit` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.getEventsBatchLimit()`.
 */
int32_t dxfg_v4_SubscriptionController_getEventsBatchLimit(graal_isolatethread_t *thread,
                                                           dxfg_v4_subscription_controller_t *subscriptionController,
                                                           DXFG_OUT int32_t *eventsBatchLimit);

/**
 * Sets maximum number of events in the single notification of the listener.
 * This is not a strict guarantee, some implementations may exceed the limit.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] eventsBatchLimit The notification events limit: 0 for the optimal limit (the default), INT32_MAX for the
 * maximum limit. It must not be negative.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, for a negative limit). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.setEventsBatchLimit(int)`.
 */
int32_t dxfg_v4_SubscriptionController_setEventsBatchLimit(graal_isolatethread_t *thread,
                                                           dxfg_v4_subscription_controller_t *subscriptionController,
                                                           int32_t eventsBatchLimit);

/**
 * Returns the aggregation period for data for this underlying subscription.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] aggregationPeriod The aggregation period for data, represented as a TimePeriod object. It is NULL if no
 * aggregation period was set. Release the handle with dxfg_v4_JavaObjectHandler_release(). It is written on success
 * only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `aggregationPeriod` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.getAggregationPeriod()`.
 */
int32_t dxfg_v4_SubscriptionController_getAggregationPeriod(graal_isolatethread_t *thread,
                                                            dxfg_v4_subscription_controller_t *subscriptionController,
                                                            DXFG_OUT dxfg_v4_time_period_t **aggregationPeriod);

/**
 * Returns the aggregation period for data for this underlying subscription in millis.
 *
 * The function fails if no aggregation period was set (the Java side throws a `NullPointerException`): use
 * dxfg_v4_SubscriptionController_getAggregationPeriod(), which gives NULL in this case.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] aggregationPeriod The aggregation period for data, represented as a timestamp in millis. It is written
 * on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `aggregationPeriod` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.getAggregationPeriod()`.
 */
int32_t dxfg_v4_SubscriptionController_getAggregationPeriodMillis(
    graal_isolatethread_t *thread, dxfg_v4_subscription_controller_t *subscriptionController,
    DXFG_OUT int64_t *aggregationPeriod);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `0.1s` limits notification to once every 100ms (at most 10 per second).
 * The new aggregation period will take effect during the next iteration of data notification.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] aggregationPeriod The aggregation period for data. The controller does not take the handle: release it
 * with dxfg_v4_JavaObjectHandler_release() when it is no longer needed. For a subscription the period must not exceed
 * 24 hours.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.setAggregationPeriod(TimePeriod)`.
 */
int32_t dxfg_v4_SubscriptionController_setAggregationPeriod(graal_isolatethread_t *thread,
                                                            dxfg_v4_subscription_controller_t *subscriptionController,
                                                            dxfg_v4_time_period_t *aggregationPeriod);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `100` limits notification to once every 100ms (at most 10 per second).
 * The new aggregation period will take effect during the next iteration of data notification.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] aggregationPeriod The aggregation period for data in millis. For a subscription it must not exceed
 * 24 hours (86400000).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.setAggregationPeriod(TimePeriod)`, called with `TimePeriod.valueOf(long)`.
 */
int32_t dxfg_v4_SubscriptionController_setAggregationPeriodMillis(
    graal_isolatethread_t *thread, dxfg_v4_subscription_controller_t *subscriptionController,
    int64_t aggregationPeriod);

/**
 * Returns executor for processing event notifications on this underlying subscription.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] executor The executor for processing event notifications on this underlying subscription.
 * Or `NULL` if the default executor is used. Release the handle with dxfg_v4_JavaObjectHandler_release(). It is written
 * on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `executor` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.getExecutor()`.
 */
int32_t dxfg_v4_SubscriptionController_getExecutor(graal_isolatethread_t *thread,
                                                   dxfg_v4_subscription_controller_t *subscriptionController,
                                                   DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * Changes executor for processing event notifications on this underlying subscription.
 *
 * The controller does not take the handle of the executor: release it with dxfg_v4_JavaObjectHandler_release() when it
 * is no longer needed. As dxfg_v4_DXFeedSubscription_setExecutor(), the function fails with `IllegalStateException`
 * for a subscription attached to an event queue and for the executor of an attached subscription, and accepts the
 * executor of a shared queue (dxfg_v4_EventQueue_getExecutor()).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] executor The executor for processing event notifications on this underlying subscription.
 * Or `NULL` if the default executor is used.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.setExecutor(Executor)`.
 */
int32_t dxfg_v4_SubscriptionController_setExecutor(graal_isolatethread_t *thread,
                                                   dxfg_v4_subscription_controller_t *subscriptionController,
                                                   dxfg_v4_executor_t *executor);

/** @} */ // end of V4_SubscriptionController

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_SUBSCRIPTION_CONTROLLER_H
