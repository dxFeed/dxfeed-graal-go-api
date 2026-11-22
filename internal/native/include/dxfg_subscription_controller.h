// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK SubscriptionController functions and types declarations
 */

#ifndef DXFG_SUBSCRIPTION_CONTROLLER_H
#define DXFG_SUBSCRIPTION_CONTROLLER_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_javac.h"
#include "graal_isolate.h"

/** @defgroup SubscriptionController SubscriptionController
 *  @brief The controller of the underlying subscription of a subscription or of a model: attaching it to feeds,
 *  the events batch limit, the aggregation period and the executor.
 *
 *  A subscription controller is returned by dxfg_DXFeedSubscription_getSubscriptionController() and by the models
 *  (e.g. dxfg_IndexedEventTxModel_getSubscriptionController()). It manages the underlying subscription of that
 *  subscription or model. It is marked as experimental in the Java API.
 *
 *  The ownership rules (handles) are the ones of the group @ref Javac. A handle that a function returns (in an
 *  `[out]` parameter) is owned by the caller and is released with dxfg_JavaObjectHandler_release(). A handle that is
 *  passed to a function is only read by the SDK: the caller keeps its ownership. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY (0), or #DXFG_EXECUTE_FAIL (-1) on error (an `[out]` parameter that is NULL is an
 *  error as well, and an `[out]` parameter is written on success only). dxfg_get_and_clear_thread_exception_t()
 *  gives the exception.
 *  @{
 */

/// @cond
// Forward declaration: the type is described in dxfg_feed.h.
typedef struct dxfg_feed_t dxfg_feed_t;
/// @endcond

/**
 * @brief A handle of a Java SubscriptionController object; release it with dxfg_JavaObjectHandler_release().
 *
 * Get it with dxfg_DXFeedSubscription_getSubscriptionController() or from a model. Releasing the handle does not
 * affect the subscription.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/SubscriptionController.html">Javadoc</a>
 */
typedef struct dxfg_subscription_controller_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_subscription_controller_t;

/**
 * Attaches the underlying subscription to the specified feed.
 *
 * The attached subscription receives the events of the feed. It can be attached to several feeds at the same time,
 * but this is rarely necessary and is not recommended.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_attach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.attach(DXFeed)`.
 *
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/SubscriptionController.html#attach-com.dxfeed.api.DXFeed-
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_attach()")
int32_t dxfg_SubscriptionController_attach(graal_isolatethread_t *thread,
                                           dxfg_subscription_controller_t *subscriptionController, dxfg_feed_t *feed);

/**
 * Detaches the underlying subscription from the specified feed.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_detach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.detach(DXFeed)`.
 *
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/SubscriptionController.html#detach-com.dxfeed.api.DXFeed-
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_detach()")
int32_t dxfg_SubscriptionController_detach(graal_isolatethread_t *thread,
                                           dxfg_subscription_controller_t *subscriptionController, dxfg_feed_t *feed);

/**
 * Returns maximum number of events in the single notification of the listener.
 * This is not a strict guarantee, some implementations may exceed the limit.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_getEventsBatchLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] eventsBatchLimit The maximum number of events in the single notification of the listener: 0 for the
 * optimal limit (the default), INT32_MAX for the maximum limit. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `eventsBatchLimit` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.getEventsBatchLimit()`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_getEventsBatchLimit()")
int32_t dxfg_SubscriptionController_getEventsBatchLimit(graal_isolatethread_t *thread,
                                                        dxfg_subscription_controller_t *subscriptionController,
                                                        DXFG_OUT int32_t *eventsBatchLimit);

/**
 * Sets maximum number of events in the single notification of the listener.
 * This is not a strict guarantee, some implementations may exceed the limit.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_setEventsBatchLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] eventsBatchLimit The notification events limit: 0 for the optimal limit (the default), INT32_MAX for the
 * maximum limit. It must not be negative.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, for a negative limit). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.setEventsBatchLimit(int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_setEventsBatchLimit()")
int32_t dxfg_SubscriptionController_setEventsBatchLimit(graal_isolatethread_t *thread,
                                                        dxfg_subscription_controller_t *subscriptionController,
                                                        int32_t eventsBatchLimit);

/**
 * Returns the aggregation period for data for this underlying subscription.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_getAggregationPeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] aggregationPeriod The aggregation period for data, represented as a TimePeriod object. It is NULL if no
 * aggregation period was set. Release the handle with dxfg_JavaObjectHandler_release(). It is written on success
 * only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `aggregationPeriod` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.getAggregationPeriod()`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_getAggregationPeriod()")
int32_t dxfg_SubscriptionController_getAggregationPeriod(graal_isolatethread_t *thread,
                                                         dxfg_subscription_controller_t *subscriptionController,
                                                         DXFG_OUT dxfg_time_period_t **aggregationPeriod);

/**
 * Returns the aggregation period for data for this underlying subscription in millis.
 *
 * The function fails if no aggregation period was set (the Java side throws a `NullPointerException`): use
 * dxfg_SubscriptionController_getAggregationPeriod(), which gives NULL in this case.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_getAggregationPeriodMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] aggregationPeriod The aggregation period for data, represented as a timestamp in millis. It is written
 * on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `aggregationPeriod` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 *
 * Java: `SubscriptionController.getAggregationPeriod()`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_getAggregationPeriodMillis()")
int32_t dxfg_SubscriptionController_getAggregationPeriodMillis(graal_isolatethread_t *thread,
                                                               dxfg_subscription_controller_t *subscriptionController,
                                                               DXFG_OUT int64_t *aggregationPeriod);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `0.1s` limits notification to once every 100ms (at most 10 per second).
 * The new aggregation period will take effect during the next iteration of data notification.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_setAggregationPeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] aggregationPeriod The aggregation period for data. The controller does not take the handle: release it
 * with dxfg_JavaObjectHandler_release() when it is no longer needed. For a subscription the period must not exceed
 * 24 hours.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.setAggregationPeriod(TimePeriod)`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_setAggregationPeriod()")
int32_t dxfg_SubscriptionController_setAggregationPeriod(graal_isolatethread_t *thread,
                                                         dxfg_subscription_controller_t *subscriptionController,
                                                         dxfg_time_period_t *aggregationPeriod);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `100` limits notification to once every 100ms (at most 10 per second).
 * The new aggregation period will take effect during the next iteration of data notification.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_setAggregationPeriodMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] aggregationPeriod The aggregation period for data in millis. For a subscription it must not exceed
 * 24 hours (86400000).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.setAggregationPeriod(TimePeriod)`, called with `TimePeriod.valueOf(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_setAggregationPeriodMillis()")
int32_t dxfg_SubscriptionController_setAggregationPeriodMillis(graal_isolatethread_t *thread,
                                                               dxfg_subscription_controller_t *subscriptionController,
                                                               int64_t aggregationPeriod);

/**
 * Returns executor for processing event notifications on this underlying subscription.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_getExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[out] executor The executor for processing event notifications on this underlying subscription.
 * Or `NULL` if the default executor is used. Release the handle with dxfg_JavaObjectHandler_release(). It is written
 * on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `executor` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.getExecutor()`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_getExecutor()")
int32_t dxfg_SubscriptionController_getExecutor(graal_isolatethread_t *thread,
                                                dxfg_subscription_controller_t *subscriptionController,
                                                DXFG_OUT dxfg_executor_t **executor);

/**
 * Changes executor for processing event notifications on this underlying subscription.
 *
 * The controller does not take the handle of the executor: release it with dxfg_JavaObjectHandler_release() when it
 * is no longer needed.
 *
 * @deprecated Use dxfg_v4_SubscriptionController_setExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] subscriptionController The subscription controller handle.
 * @param[in] executor The executor for processing event notifications on this underlying subscription.
 * Or `NULL` if the default executor is used.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `SubscriptionController.setExecutor(Executor)`.
 */
DXFG_DEPRECATED("use dxfg_v4_SubscriptionController_setExecutor()")
int32_t dxfg_SubscriptionController_setExecutor(graal_isolatethread_t *thread,
                                                dxfg_subscription_controller_t *subscriptionController,
                                                dxfg_executor_t *executor);

/** @} */ // end of SubscriptionController

#ifdef __cplusplus
}
#endif

#endif // DXFG_SUBSCRIPTION_CONTROLLER_H