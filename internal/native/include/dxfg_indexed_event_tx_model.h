// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK IndexedEventTxModel functions and type declarations
 */

#ifndef DXFG_INDEXED_EVENT_TX_MODEL_H
#define DXFG_INDEXED_EVENT_TX_MODEL_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_catch_exception.h"
#include "dxfg_events.h"
#include "dxfg_javac.h"
#include "graal_isolate.h"

/** @defgroup IndexedEventTxModel IndexedEventTxModel
 *  @{
 */

/// @cond
// Forward declarations: the types are described in dxfg_feed.h and dxfg_subscription_controller.h.
typedef struct dxfg_feed_t dxfg_feed_t;
typedef struct dxfg_subscription_controller_t dxfg_subscription_controller_t;
/// @endcond

/**
 * The IndexedEventTxModel instance's handle.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/IndexedEventTxModel.html">Javadoc</a>
 */
typedef struct dxfg_indexed_event_tx_model_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_indexed_event_tx_model_t;

/**
 * The IndexedEventTxModel.Builder instance's handle.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/IndexedEventTxModel.Builder.html">Javadoc</a>
 */
typedef struct dxfg_indexed_event_tx_model_builder_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_indexed_event_tx_model_builder_t;

/**
 * The IndexedEventTxModel.Listener instance's handle.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/IndexedEventTxModel.Listener.html">Javadoc</a>
 */
typedef struct dxfg_indexed_event_tx_model_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_indexed_event_tx_model_listener_t;

/**
 * Invoked when a complete transaction (one or more) is received.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] events The list of received events representing one or more completed transactions.
 * @param[in] isSnapshot `true` (1) if the events form a snapshot; `false` (0) otherwise.
 * @param[in] userData Passed to the created listener user data.
 */
typedef void (*dxfg_IndexedEventTxModel_Listener_eventsReceived_f)(graal_isolatethread_t *thread,
                                                                   dxfg_event_type_list *events, int32_t isSnapshot,
                                                                   void *userData);

/**
 * Factory method to create a new builder for this model.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_newBuilder().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of indexed events processed by the model.
 * @param[out] builder A new builder instance handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the builder handle.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_newBuilder()")
int32_t dxfg_IndexedEventTxModel_newBuilder(graal_isolatethread_t *thread, dxfg_event_clazz_t eventType,
                                            DXFG_OUT dxfg_indexed_event_tx_model_builder_t **builder);

/**
 * Returns the indexed event subscription symbol associated with this model.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_getIndexedEventSubscriptionSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] symbol The indexed event subscription symbol associated with this model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_Symbol_release() to free the symbol.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_getIndexedEventSubscriptionSymbol()")
int32_t dxfg_IndexedEventTxModel_getIndexedEventSubscriptionSymbol(graal_isolatethread_t *thread,
                                                                   dxfg_indexed_event_tx_model_t *model,
                                                                   DXFG_OUT dxfg_symbol_t **symbol);

/**
 * Returns the symbol associated with this model.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_getSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] symbol The symbol associated with this model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the symbol.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_getSymbol()")
int32_t dxfg_IndexedEventTxModel_getSymbol(graal_isolatethread_t *thread, dxfg_indexed_event_tx_model_t *model,
                                           DXFG_OUT const char **symbol);

/**
 * Returns the source associated with this model.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_getSource().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] source The source associated with this model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_IndexedEventSource_release() to free the source.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_getSource()")
int32_t dxfg_IndexedEventTxModel_getSource(graal_isolatethread_t *thread, dxfg_indexed_event_tx_model_t *model,
                                           DXFG_OUT dxfg_indexed_event_source_t **source);

/**
 * Returns the subscription controller for the underlying subscription.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_getSubscriptionController().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] subscriptionController The subscription controller for the underlying subscription.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the SubscriptionController handle.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_getSubscriptionController()")
int32_t
dxfg_IndexedEventTxModel_getSubscriptionController(graal_isolatethread_t *thread, dxfg_indexed_event_tx_model_t *model,
                                                   DXFG_OUT dxfg_subscription_controller_t **subscriptionController);

/**
 * Returns whether this model is closed.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_isClosed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @param[out] isClosed `true` (1) if this model is closed; `false` (0) otherwise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_isClosed()")
int32_t dxfg_IndexedEventTxModel_isClosed(graal_isolatethread_t *thread, dxfg_indexed_event_tx_model_t *model,
                                          DXFG_OUT int32_t *isClosed);

/**
 * Closes this model: it is detached from the feed for good and stops the notifications.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The current model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_close()")
int32_t dxfg_IndexedEventTxModel_close(graal_isolatethread_t *thread, dxfg_indexed_event_tx_model_t *model);

/**
 * Constructs a new builder for the specified event type.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of indexed events processed by the model.
 * @param[out] builder A new builder instance handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the builder handle.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_new()")
int32_t dxfg_IndexedEventTxModel_Builder_new(graal_isolatethread_t *thread, dxfg_event_clazz_t eventType,
                                             DXFG_OUT dxfg_indexed_event_tx_model_builder_t **builder);

/**
 * Sets the feed for the model being created.
 *
 * The feed can also be attached later, after the model has been built using the
 * @ref dxfg_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withFeed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] feed The feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withFeed()")
int32_t dxfg_IndexedEventTxModel_Builder_withFeed(graal_isolatethread_t *thread,
                                                  dxfg_indexed_event_tx_model_builder_t *builder, dxfg_feed_t *feed);

/**
 * Sets the subscription symbol and its source (dxfg_indexed_event_subscription_symbol_t) for the model being created.
 *
 * This symbol and source cannot be added or changed after the model has been built.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] symbol The subscription symbol (dxfg_indexed_event_subscription_symbol_t).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withSymbol()")
int32_t dxfg_IndexedEventTxModel_Builder_withSymbol(graal_isolatethread_t *thread,
                                                    dxfg_indexed_event_tx_model_builder_t *builder,
                                                    dxfg_symbol_t *symbol);

/**
 * Sets the subscription symbol for the model being created.
 *
 * The symbol cannot be added or changed after the model has been built.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withStringSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] symbol The subscription symbol.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withStringSymbol()")
int32_t dxfg_IndexedEventTxModel_Builder_withSymbol2(graal_isolatethread_t *thread,
                                                     dxfg_indexed_event_tx_model_builder_t *builder,
                                                     const char *symbol);

/**
 * Sets the source to subscribe to.
 *
 * The source cannot be added or changed after the model has been built.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withSource().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] source The specified source.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withSource()")
int32_t dxfg_IndexedEventTxModel_Builder_withSource(graal_isolatethread_t *thread,
                                                    dxfg_indexed_event_tx_model_builder_t *builder,
                                                    dxfg_indexed_event_source_t *source);

/**
 * Creates a new listener for the model.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Listener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventsReceived The callback function pointer.
 * @param[in] userData The user data which will be passed to callback.
 * @param[out] listener The new listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the listener handle.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Listener_new()")
int32_t dxfg_IndexedEventTxModel_Listener_new(graal_isolatethread_t *thread,
                                              dxfg_IndexedEventTxModel_Listener_eventsReceived_f eventsReceived,
                                              void *userData,
                                              DXFG_OUT dxfg_indexed_event_tx_model_listener_t **listener);

/**
 * Sets the listener for transaction notifications.
 * The notification is invoked from a separate thread via the @ref dxfg_IndexedEventTxModel_Builder_withExecutor()
 * "executor". The listener cannot be added or changed after the model has been built.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] listener The model listener.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withListener()")
int32_t dxfg_IndexedEventTxModel_Builder_withListener(graal_isolatethread_t *thread,
                                                      dxfg_indexed_event_tx_model_builder_t *builder,
                                                      dxfg_indexed_event_tx_model_listener_t *listener);

/**
 * Sets the maximum number of events in the single notification of dxfg_IndexedEventTxModel_Listener_eventsReceived_f.
 * This value can be exceeded if the size of the received transaction exceeds the set limit.
 * This value can be changed later, after the model has been built using the
 * @ref dxfg_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withEventsBatchLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] eventsBatchLimit The notification events limit.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withEventsBatchLimit()")
int32_t dxfg_IndexedEventTxModel_Builder_withEventsBatchLimit(graal_isolatethread_t *thread,
                                                              dxfg_indexed_event_tx_model_builder_t *builder,
                                                              int32_t eventsBatchLimit);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `0.1s` limits notification to once every 100ms (at most 10 per second).
 * This value can be changed later, after the model has been built using the
 * @ref dxfg_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withAggregationPeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] aggregationPeriod The aggregation period for data.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withAggregationPeriod()")
int32_t dxfg_IndexedEventTxModel_Builder_withAggregationPeriod(graal_isolatethread_t *thread,
                                                               dxfg_indexed_event_tx_model_builder_t *builder,
                                                               dxfg_time_period_t *aggregationPeriod);

/**
 * Sets the aggregation period for data that limits the rate of data notifications.
 * For example, setting the value to `100` limits notification to once every 100ms (at most 10 per second).
 * This value can be changed later, after the model has been built using the
 * @ref dxfg_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withAggregationPeriodMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] aggregationPeriod The aggregation period for data.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withAggregationPeriodMillis()")
int32_t dxfg_IndexedEventTxModel_Builder_withAggregationPeriodMillis(graal_isolatethread_t *thread,
                                                                     dxfg_indexed_event_tx_model_builder_t *builder,
                                                                     int64_t aggregationPeriod);

/**
 * Sets the executor for processing listener notifications.
 * The default executor for all models is the endpoint's one (see dxfg_DXEndpoint_Builder_withExecutor()).
 * This executor can be changed later, after the model has been built using the
 * @ref dxfg_IndexedEventTxModel_getSubscriptionController() "subscription controller".
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_withExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[in] executor The executor instance.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_withExecutor()")
int32_t dxfg_IndexedEventTxModel_Builder_withExecutor(graal_isolatethread_t *thread,
                                                      dxfg_indexed_event_tx_model_builder_t *builder,
                                                      dxfg_executor_t *executor);

/**
 * Builds a new instance of dxfg_indexed_event_tx_model_t with the provided configuration.
 *
 * @deprecated Use dxfg_v4_IndexedEventTxModel_Builder_build().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The current builder.
 * @param[out] model The created model.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the model handle.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventTxModel_Builder_build()")
int32_t dxfg_IndexedEventTxModel_Builder_build(graal_isolatethread_t *thread,
                                               dxfg_indexed_event_tx_model_builder_t *builder,
                                               DXFG_OUT dxfg_indexed_event_tx_model_t **model);

/** @} */ // end of Feed

#ifdef __cplusplus
}
#endif

#endif // DXFG_INDEXED_EVENT_TX_MODEL_H
