// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_indexed_event_tx_model.h and the v4 names of its types (see the group @ref V4_IndexedEventTxModel).
 */

#ifndef DXFG_V4_INDEXED_EVENT_TX_MODEL_H
#define DXFG_V4_INDEXED_EVENT_TX_MODEL_H

#include "dxfg_indexed_event_tx_model.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_catch_exception.h"
#include "dxfg_v4_events.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_IndexedEventTxModel IndexedEventTxModel (v4)
 *  @brief The v4 functions and type names of the group @ref IndexedEventTxModel. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_feed_t */
typedef dxfg_feed_t dxfg_v4_feed_t;

/** @copydoc dxfg_subscription_controller_t */
typedef dxfg_subscription_controller_t dxfg_v4_subscription_controller_t;

/** @copydoc dxfg_indexed_event_tx_model_t */
typedef dxfg_indexed_event_tx_model_t dxfg_v4_indexed_event_tx_model_t;

/** @copydoc dxfg_indexed_event_tx_model_builder_t */
typedef dxfg_indexed_event_tx_model_builder_t dxfg_v4_indexed_event_tx_model_builder_t;

/** @copydoc dxfg_indexed_event_tx_model_listener_t */
typedef dxfg_indexed_event_tx_model_listener_t dxfg_v4_indexed_event_tx_model_listener_t;

/** @copydoc dxfg_IndexedEventTxModel_Listener_eventsReceived_f */
typedef dxfg_IndexedEventTxModel_Listener_eventsReceived_f dxfg_v4_IndexedEventTxModel_Listener_eventsReceived_f;

/**
 * Factory method to create a new builder for this model.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of indexed events processed by the model.
 * @param[out] builder A new builder instance handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the builder handle.
 */
int32_t dxfg_v4_IndexedEventTxModel_newBuilder(graal_isolatethread_t *thread, dxfg_v4_event_clazz_t eventType,
                                               DXFG_OUT dxfg_v4_indexed_event_tx_model_builder_t **builder);

/**
 * Returns the indexed event subscription symbol associated with this model.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] symbol The indexed event subscription symbol associated with this model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_Symbol_release() to free the symbol.
 */
int32_t dxfg_v4_IndexedEventTxModel_getIndexedEventSubscriptionSymbol(graal_isolatethread_t *thread,
                                                                      dxfg_v4_indexed_event_tx_model_t *model,
                                                                      DXFG_OUT dxfg_v4_symbol_t **symbol);

/**
 * Returns the symbol associated with this model.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] symbol The symbol associated with this model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the symbol.
 */
int32_t dxfg_v4_IndexedEventTxModel_getSymbol(graal_isolatethread_t *thread, dxfg_v4_indexed_event_tx_model_t *model,
                                              DXFG_OUT const char **symbol);

/**
 * Returns the source associated with this model.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] source The source associated with this model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_IndexedEventSource_release() to free the source.
 */
int32_t dxfg_v4_IndexedEventTxModel_getSource(graal_isolatethread_t *thread, dxfg_v4_indexed_event_tx_model_t *model,
                                              DXFG_OUT dxfg_v4_indexed_event_source_t **source);

/**
 * Returns the subscription controller for the underlying subscription.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] subscriptionController The subscription controller for the underlying subscription.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the SubscriptionController handle.
 */
int32_t dxfg_v4_IndexedEventTxModel_getSubscriptionController(
    graal_isolatethread_t *thread, dxfg_v4_indexed_event_tx_model_t *model,
    DXFG_OUT dxfg_v4_subscription_controller_t **subscriptionController);

/**
 * Returns whether this model is closed.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] isClosed `true` (1) if this model is closed; `false` (0) otherwise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_isClosed(graal_isolatethread_t *thread, dxfg_v4_indexed_event_tx_model_t *model,
                                             DXFG_OUT int32_t *isClosed);

/**
 * Closes this model: it is detached from the feed for good and stops the notifications.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_close(graal_isolatethread_t *thread, dxfg_v4_indexed_event_tx_model_t *model);

/**
 * Constructs a new builder for the specified event type.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of indexed events processed by the model.
 * @param[out] builder A new builder instance handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the builder handle.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_new(graal_isolatethread_t *thread, dxfg_v4_event_clazz_t eventType,
                                                DXFG_OUT dxfg_v4_indexed_event_tx_model_builder_t **builder);

/**
 * Sets the feed for the model being created.
 *
 * The feed can also be attached later, after the model has been built using the
 * @ref dxfg_v4_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] feed The feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withFeed(graal_isolatethread_t *thread,
                                                     dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                     dxfg_v4_feed_t *feed);

/**
 * Sets the subscription symbol and its source (dxfg_v4_indexed_event_subscription_symbol_t) for the model being
 * created.
 *
 * This symbol and source cannot be added or changed after the model has been built.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] symbol The subscription symbol (dxfg_v4_indexed_event_subscription_symbol_t).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withSymbol(graal_isolatethread_t *thread,
                                                       dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                       dxfg_v4_symbol_t *symbol);

/**
 * Sets the subscription symbol for the model being created.
 *
 * The symbol cannot be added or changed after the model has been built.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] symbol The subscription symbol.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withStringSymbol(graal_isolatethread_t *thread,
                                                             dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                             const char *symbol);

/**
 * Sets the source to subscribe to.
 *
 * The source cannot be added or changed after the model has been built.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] source The specified source.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withSource(graal_isolatethread_t *thread,
                                                       dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                       dxfg_v4_indexed_event_source_t *source);

/**
 * Creates a new listener for the model.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventsReceived The callback function pointer.
 * @param[in] userData The user data which will be passed to callback.
 * @param[out] listener The new listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the listener handle.
 */
int32_t dxfg_v4_IndexedEventTxModel_Listener_new(graal_isolatethread_t *thread,
                                                 dxfg_v4_IndexedEventTxModel_Listener_eventsReceived_f eventsReceived,
                                                 void *userData,
                                                 DXFG_OUT dxfg_v4_indexed_event_tx_model_listener_t **listener);

/**
 * Sets the listener for transaction notifications.
 * The notification is invoked from a separate thread via the @ref dxfg_v4_IndexedEventTxModel_Builder_withExecutor()
 * "executor". The listener cannot be added or changed after the model has been built.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] listener The model listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withListener(graal_isolatethread_t *thread,
                                                         dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                         dxfg_v4_indexed_event_tx_model_listener_t *listener);

/**
 * Sets the maximum number of events in the single notification of
 * dxfg_v4_IndexedEventTxModel_Listener_eventsReceived_f. This value can be exceeded if the size of the received
 * transaction exceeds the set limit. This value can be changed later, after the model has been built using the
 * @ref dxfg_v4_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] eventsBatchLimit The notification events limit.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withEventsBatchLimit(graal_isolatethread_t *thread,
                                                                 dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                                 int32_t eventsBatchLimit);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `0.1s` limits notification to once every 100ms (at most 10 per second).
 * This value can be changed later, after the model has been built using the
 * @ref dxfg_v4_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] aggregationPeriod The aggregation period for data.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withAggregationPeriod(graal_isolatethread_t *thread,
                                                                  dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                                  dxfg_v4_time_period_t *aggregationPeriod);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `100` limits notification to once every 100ms (at most 10 per second).
 * This value can be changed later, after the model has been built using the
 * @ref dxfg_v4_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] aggregationPeriod The aggregation period for data.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withAggregationPeriodMillis(
    graal_isolatethread_t *thread, dxfg_v4_indexed_event_tx_model_builder_t *builder, int64_t aggregationPeriod);

/**
 * Sets the executor for processing listener notifications.
 * The default executor for all models is the endpoint's one (see dxfg_v4_DXEndpoint_Builder_withExecutor()).
 * This executor can be changed later, after the model has been built using the
 * @ref dxfg_v4_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] executor The executor instance.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_withExecutor(graal_isolatethread_t *thread,
                                                         dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                         dxfg_v4_executor_t *executor);

/**
 * Builds a new instance of dxfg_v4_indexed_event_tx_model_t with the provided configuration.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[out] model The created model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the model handle.
 */
int32_t dxfg_v4_IndexedEventTxModel_Builder_build(graal_isolatethread_t *thread,
                                                  dxfg_v4_indexed_event_tx_model_builder_t *builder,
                                                  DXFG_OUT dxfg_v4_indexed_event_tx_model_t **model);

/** @} */ // end of V4_IndexedEventTxModel

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_INDEXED_EVENT_TX_MODEL_H
