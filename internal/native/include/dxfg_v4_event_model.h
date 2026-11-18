// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_event_model.h and the v4 names of its types (see the group @ref V4_EventModel).
 */

#ifndef DXFG_V4_EVENT_MODEL_H
#define DXFG_V4_EVENT_MODEL_H

#include "dxfg_event_model.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"
#include "dxfg_v4_feed.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_EventModel EventModel (v4)
 *  @brief The v4 functions and type names of the group @ref EventModel. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_order_book_model_t */
typedef dxfg_order_book_model_t dxfg_v4_order_book_model_t;

/** @copydoc dxfg_indexed_event_model_t */
typedef dxfg_indexed_event_model_t dxfg_v4_indexed_event_model_t;

/** @copydoc dxfg_time_series_event_model_t */
typedef dxfg_time_series_event_model_t dxfg_v4_time_series_event_model_t;

/** @copydoc dxfg_observable_list_model_t */
typedef dxfg_observable_list_model_t dxfg_v4_observable_list_model_t;

/** @copydoc dxfg_observable_list_model_listener_t */
typedef dxfg_observable_list_model_listener_t dxfg_v4_observable_list_model_listener_t;

/** @copydoc dxfg_observable_list_model_listener_function */
typedef dxfg_observable_list_model_listener_function dxfg_v4_observable_list_model_listener_function;

/** @copydoc dxfg_order_book_model_listener_t */
typedef dxfg_order_book_model_listener_t dxfg_v4_order_book_model_listener_t;

/** @copydoc dxfg_order_book_model_listener_function */
typedef dxfg_order_book_model_listener_function dxfg_v4_order_book_model_listener_function;

/** @copydoc dxfg_order_book_model_filter_t */
typedef dxfg_order_book_model_filter_t dxfg_v4_order_book_model_filter_t;

/** @copydoc dxfg_TimeSeriesTxModelSortOrder_t */
typedef dxfg_TimeSeriesTxModelSortOrder_t dxfg_v4_TimeSeriesTxModelSortOrder_t;

/** @copydoc dxfg_indexed_tx_model_t */
typedef dxfg_indexed_tx_model_t dxfg_v4_indexed_tx_model_t;

/** @copydoc dxfg_indexed_tx_model_builder_t */
typedef dxfg_indexed_tx_model_builder_t dxfg_v4_indexed_tx_model_builder_t;

/** @copydoc dxfg_time_series_tx_model_t */
typedef dxfg_time_series_tx_model_t dxfg_v4_time_series_tx_model_t;

/** @copydoc dxfg_time_series_tx_model_builder_t */
typedef dxfg_time_series_tx_model_builder_t dxfg_v4_time_series_tx_model_builder_t;

/** @copydoc dxfg_tx_model_listener_t */
typedef dxfg_tx_model_listener_t dxfg_v4_tx_model_listener_t;

/** @copydoc dxfg_TxModelListener_function_eventsReceived */
typedef dxfg_TxModelListener_function_eventsReceived dxfg_v4_TxModelListener_function_eventsReceived;

/**
 * @brief Creates a new order book model.
 *
 * The model is not attached to any feed and not subscribed to any symbol. Its filter is
 * ALL and its lot size is 1. Set the symbol with dxfg_v4_OrderBookModel_setSymbol() and attach the model
 * with dxfg_v4_OrderBookModel_attach() to start receiving orders.
 *
 * Java: `new OrderBookModel()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] orderBookModel The pointer to the variable that receives the model (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_new(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_order_book_model_t **orderBookModel);

/**
 * @brief Creates a new model of indexed events of the event type.
 *
 * The model is not attached to any feed and not subscribed to any symbol. Set the symbol with
 * dxfg_v4_IndexedEventModel_setSymbol() and attach the model with dxfg_v4_IndexedEventModel_attach() to start
 * receiving events.
 *
 * Java: `new IndexedEventModel(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] clazz The event type of the model; it has to be an indexed event type.
 * @param[out] indexedEventModel The pointer to the variable that receives the model (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_new(graal_isolatethread_t *thread, dxfg_v4_event_clazz_t clazz,
                                      DXFG_OUT dxfg_v4_indexed_event_model_t **indexedEventModel);

/**
 * @brief Creates a new model of time-series events of the event type.
 *
 * The model is not attached to any feed and not subscribed to any symbol (its from-time is INT64_MAX). Set
 * the symbol and the from-time with dxfg_v4_TimeSeriesEventModel_setSymbol() and
 * dxfg_v4_TimeSeriesEventModel_setFromTime() and attach the model with dxfg_v4_TimeSeriesEventModel_attach() to start
 * receiving events.
 *
 * Java: `new TimeSeriesEventModel(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] clazz The event type of the model; it has to be a time-series event type.
 * @param[out] timeSeriesEventModel The pointer to the variable that receives the model (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_new(graal_isolatethread_t *thread, dxfg_v4_event_clazz_t clazz,
                                         DXFG_OUT dxfg_v4_time_series_event_model_t **timeSeriesEventModel);

/**
 * @brief Attaches the order book model to the feed: the model starts receiving the events of its symbol.
 *
 * The feed keeps an attached model: releasing the handle of the model does not detach it. Detach the
 * model with dxfg_v4_OrderBookModel_detach() or close it with dxfg_v4_OrderBookModel_close().
 *
 * Java: `OrderBookModel.attach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_attach(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                      dxfg_v4_feed_t *feed);

/**
 * @brief Detaches the order book model from the feed.
 *
 * Java: `OrderBookModel.detach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_detach(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                      dxfg_v4_feed_t *feed);

/**
 * @brief Closes the order book model: it becomes permanently detached from all feeds and removes all its listeners.
 *
 * Java: `OrderBookModel.close()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_close(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model);

/**
 * @brief Returns the executor that calls the listeners of the order book model.
 *
 * Java: `OrderBookModel.getExecutor()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[out] executor The pointer to the variable that receives the executor (release it with
 * dxfg_v4_JavaObjectHandler_release()), or NULL if the default executor of the feed is used; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_getExecutor(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                           DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * @brief Sets the executor that calls the listeners of the order book model.
 *
 * The executor of a shared event queue (dxfg_v4_EventQueue_getExecutor()) makes the polls of the queue call the
 * listeners; the executor of a subscription attached to a queue is refused with `IllegalStateException` (see
 * @ref V4_EventQueue).
 *
 * Java: `OrderBookModel.setExecutor(Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] executor The executor, or NULL to use the default executor of the feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_setExecutor(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                           dxfg_v4_executor_t *executor);

/**
 * @brief Clears the subscription symbol and, with it, all the orders of the model.
 *
 * This is a shortcut for dxfg_v4_OrderBookModel_setSymbol() with NULL.
 *
 * Java: `OrderBookModel.clear()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_clear(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model);

/**
 * @brief Attaches the indexed event model to the feed: the model starts receiving the events of its symbol.
 *
 * The feed keeps an attached model: releasing the handle of the model does not detach it. Detach the
 * model with dxfg_v4_IndexedEventModel_detach() or close it with dxfg_v4_IndexedEventModel_close().
 *
 * Java: `IndexedEventModel.attach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_attach(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                         dxfg_v4_feed_t *feed);

/**
 * @brief Detaches the indexed event model from the feed.
 *
 * Java: `IndexedEventModel.detach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_detach(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                         dxfg_v4_feed_t *feed);

/**
 * @brief Closes the indexed event model: it becomes permanently detached from all feeds and removes all its listeners.
 *
 * Java: `IndexedEventModel.close()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_close(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model);

/**
 * @brief Returns the executor that calls the listeners of the indexed event model.
 *
 * Java: `IndexedEventModel.getExecutor()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[out] executor The pointer to the variable that receives the executor (release it with
 * dxfg_v4_JavaObjectHandler_release()), or NULL if the default executor of the feed is used; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_getExecutor(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                              DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * @brief Sets the executor that calls the listeners of the indexed event model.
 *
 * The executor of a shared event queue (dxfg_v4_EventQueue_getExecutor()) makes the polls of the queue call the
 * listeners; the executor of a subscription attached to a queue is refused with `IllegalStateException` (see
 * @ref V4_EventQueue).
 *
 * Java: `IndexedEventModel.setExecutor(Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] executor The executor, or NULL to use the default executor of the feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_setExecutor(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                              dxfg_v4_executor_t *executor);

/**
 * @brief Clears the subscription symbol and, with it, all the events of the model.
 *
 * This is a shortcut for dxfg_v4_IndexedEventModel_setSymbol() with NULL.
 *
 * Java: `IndexedEventModel.clear()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_clear(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model);

/**
 * @brief Attaches the time-series event model to the feed: the model starts receiving the events of its symbol.
 *
 * The feed keeps an attached model: releasing the handle of the model does not detach it. Detach the
 * model with dxfg_v4_TimeSeriesEventModel_detach() or close it with dxfg_v4_TimeSeriesEventModel_close().
 *
 * Java: `TimeSeriesEventModel.attach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_attach(graal_isolatethread_t *thread, dxfg_v4_time_series_event_model_t *model,
                                            dxfg_v4_feed_t *feed);

/**
 * @brief Detaches the time-series event model from the feed.
 *
 * Java: `TimeSeriesEventModel.detach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_detach(graal_isolatethread_t *thread, dxfg_v4_time_series_event_model_t *model,
                                            dxfg_v4_feed_t *feed);

/**
 * @brief Closes the time-series event model: it becomes permanently detached from all feeds and removes all its
 * listeners.
 *
 * Java: `TimeSeriesEventModel.close()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_close(graal_isolatethread_t *thread, dxfg_v4_time_series_event_model_t *model);

/**
 * @brief Returns the executor that calls the listeners of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.getExecutor()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[out] executor The pointer to the variable that receives the executor (release it with
 * dxfg_v4_JavaObjectHandler_release()), or NULL if the default executor of the feed is used; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_getExecutor(graal_isolatethread_t *thread,
                                                 dxfg_v4_time_series_event_model_t *model,
                                                 DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * @brief Sets the executor that calls the listeners of the time-series event model.
 *
 * The executor of a shared event queue (dxfg_v4_EventQueue_getExecutor()) makes the polls of the queue call the
 * listeners; the executor of a subscription attached to a queue is refused with `IllegalStateException` (see
 * @ref V4_EventQueue).
 *
 * Java: `TimeSeriesEventModel.setExecutor(Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] executor The executor, or NULL to use the default executor of the feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_setExecutor(graal_isolatethread_t *thread,
                                                 dxfg_v4_time_series_event_model_t *model,
                                                 dxfg_v4_executor_t *executor);

/**
 * @brief Clears the subscription symbol, sets the from-time to INT64_MAX and, with it, clears all the events
 * of the model.
 *
 * Java: `TimeSeriesEventModel.clear()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_clear(graal_isolatethread_t *thread, dxfg_v4_time_series_event_model_t *model);

/**
 * Returns the filter of an order book model.
 *
 * Java: `OrderBookModel.getFilter()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model.
 * @param[out] filter The pointer to the variable that receives the filter of the model, a
 * #dxfg_v4_order_book_model_filter_t value; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_getFilter(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                         DXFG_OUT dxfg_v4_order_book_model_filter_t *filter);

/**
 * @brief Sets the filter of the order book model.
 *
 * Java: `OrderBookModel.setFilter(OrderBookModelFilter)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] filter The filter.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_setFilter(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                         dxfg_v4_order_book_model_filter_t filter);

/**
 * @brief Returns the subscription symbol of the order book model.
 *
 * Java: `OrderBookModel.getSymbol()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[out] symbol The pointer to the variable that receives the symbol (a UTF-8 string; release it with
 * dxfg_v4_String_release()), or NULL if the model is not subscribed (the default); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_getSymbol(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                         DXFG_OUT const char **symbol);

/**
 * @brief Sets the subscription symbol of the order book model.
 *
 * Java: `OrderBookModel.setSymbol(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] symbol The symbol, a UTF-8 string, or NULL to unsubscribe (that removes the orders of the model).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_setSymbol(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                         const char* symbol);

/**
 * @brief Returns the subscription symbol of the indexed event model.
 *
 * Java: `IndexedEventModel.getSymbol()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[out] symbol The pointer to the variable that receives the symbol (release it with dxfg_v4_Symbol_release()),
 * or NULL if the model is not subscribed (the default); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_getSymbol(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                            DXFG_OUT dxfg_v4_symbol_t **symbol);

/**
 * @brief Sets the subscription symbol of the indexed event model.
 *
 * Java: `IndexedEventModel.setSymbol(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] symbol The symbol, or NULL to unsubscribe (that removes the events of the model).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_setSymbol(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                            dxfg_v4_symbol_t* symbol);

/**
 * @brief Returns the subscription symbol of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.getSymbol()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[out] symbol The pointer to the variable that receives the symbol (release it with dxfg_v4_Symbol_release()),
 * or NULL if the model is not subscribed (the default); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_getSymbol(graal_isolatethread_t *thread, dxfg_v4_time_series_event_model_t *model,
                                               DXFG_OUT dxfg_v4_symbol_t **symbol);

/**
 * @brief Sets the subscription symbol of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.setSymbol(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] symbol The symbol, or NULL to unsubscribe (that removes the events of the model).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_setSymbol(graal_isolatethread_t *thread, dxfg_v4_time_series_event_model_t *model,
                                               dxfg_v4_symbol_t* symbol);

/**
 * @brief Returns the lot size of the order book model.
 *
 * The lot size is a multiplier applied to the sizes of the composite, regional and aggregate orders (1 by
 * default).
 *
 * Java: `OrderBookModel.getLotSize()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[out] lotSize The pointer to the variable that receives the lot size; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_getLotSize(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                          DXFG_OUT int32_t *lotSize);

/**
 * @brief Sets the lot size of the order book model.
 *
 * A lot size less than 1 is an error.
 *
 * Java: `OrderBookModel.setLotSize(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] lotSize The lot size multiplier, 1 or greater.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_setLotSize(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                          int32_t lotSize);

/**
 * @brief Returns the size limit of the indexed event model.
 *
 * Java: `IndexedEventModel.getSizeLimit()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[out] sizeLimit The pointer to the variable that receives the size limit (INT32_MAX means no limit, it is the
 * default); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_getSizeLimit(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                               DXFG_OUT int32_t *sizeLimit);

/**
 * @brief Sets the size limit of the indexed event model.
 *
 * A negative size limit is an error.
 *
 * Java: `IndexedEventModel.setSizeLimit(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] sizeLimit The size limit, not negative. When the limit is exceeded, the first events of the list
 * are dropped.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_setSizeLimit(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                               int32_t sizeLimit);

/**
 * @brief Returns the size limit of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.getSizeLimit()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[out] sizeLimit The pointer to the variable that receives the size limit (INT32_MAX means no limit, it is the
 * default); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_getSizeLimit(graal_isolatethread_t *thread,
                                                  dxfg_v4_time_series_event_model_t *model,
                                                  DXFG_OUT int32_t *sizeLimit);

/**
 * @brief Sets the size limit of the time-series event model.
 *
 * A negative size limit is an error.
 *
 * Java: `TimeSeriesEventModel.setSizeLimit(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] sizeLimit The size limit, not negative. When the limit is exceeded, the first events of the list
 * are dropped.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_setSizeLimit(graal_isolatethread_t *thread,
                                                  dxfg_v4_time_series_event_model_t *model, int32_t sizeLimit);

/**
 * @brief Returns the time from which the time-series event model subscribes for events.
 *
 * Java: `TimeSeriesEventModel.getFromTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[out] fromTime The pointer to the variable that receives the time in milliseconds since the epoch (INT64_MAX
 * means not subscribed, it is the default); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_getFromTime(graal_isolatethread_t *thread,
                                                 dxfg_v4_time_series_event_model_t *model, DXFG_OUT int64_t *fromTime);

/**
 * @brief Sets the time from which the time-series event model subscribes for events.
 *
 * Java: `TimeSeriesEventModel.setFromTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] fromTime The time in milliseconds since the epoch; INT64_MAX unsubscribes.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_setFromTime(graal_isolatethread_t *thread,
                                                 dxfg_v4_time_series_event_model_t *model, int64_t fromTime);

/**
 * @brief Returns the buy orders (the bid side) of the order book model.
 *
 * The orders are arranged by price: the first order is the top of the bid side (the highest price).
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_v4_ObservableListModel_toArray().
 *
 * Java: `OrderBookModel.getBuyOrders()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[out] buyOrders The pointer to the variable that receives the list (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_getBuyOrders(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                            DXFG_OUT dxfg_v4_observable_list_model_t **buyOrders);

/**
 * @brief Returns the sell orders (the offer side) of the order book model.
 *
 * The orders are arranged by price: the first order is the top of the offer side (the lowest price).
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_v4_ObservableListModel_toArray().
 *
 * Java: `OrderBookModel.getSellOrders()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[out] sellOrders The pointer to the variable that receives the list (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_getSellOrders(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                             DXFG_OUT dxfg_v4_observable_list_model_t **sellOrders);

/**
 * @brief Returns the list of events of the indexed event model.
 *
 * The events are ordered by the id of their source and by their index.
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_v4_ObservableListModel_toArray().
 *
 * Java: `IndexedEventModel.getEventsList()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[out] eventsList The pointer to the variable that receives the list (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventModel_getEventsList(graal_isolatethread_t *thread, dxfg_v4_indexed_event_model_t *model,
                                                DXFG_OUT dxfg_v4_observable_list_model_t **eventsList);

/**
 * @brief Returns the list of events of the time-series event model.
 *
 * The events are ordered by their time.
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_v4_ObservableListModel_toArray().
 *
 * Java: `TimeSeriesEventModel.getEventsList()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[out] eventsList The pointer to the variable that receives the list (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeSeriesEventModel_getEventsList(graal_isolatethread_t *thread,
                                                   dxfg_v4_time_series_event_model_t *model,
                                                   DXFG_OUT dxfg_v4_observable_list_model_t **eventsList);

/**
 * @brief Creates a listener of the changes of a list that calls the function with the user data.
 *
 * Add the listener to a list with dxfg_v4_ObservableListModel_addListener(). After the last handle of the listener
 * is released the function is not called anymore and none of its calls is in progress, so the user data can be
 * freed. Remove the listener from the list before releasing it, otherwise it stays added and does nothing, see
 * dxfg_v4_JavaObjectHandler_release().
 *
 * Java: `ObservableListModelListener.modelChanged(ObservableListModelListener.Change)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call after the list has changed; see
 * #dxfg_v4_observable_list_model_listener_function for the thread and the ownership of its data.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it. It has to stay valid
 * until the listener is released.
 * @param[out] observableListModelListener The pointer to the variable that receives the listener (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableListModelListener_new(
    graal_isolatethread_t *thread, dxfg_v4_observable_list_model_listener_function user_func, void *user_data,
    DXFG_OUT dxfg_v4_observable_list_model_listener_t **observableListModelListener);

/**
 * @brief Adds the listener to the list.
 *
 * Java: `ObservableListModel.addListener(ObservableListModelListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableListModel_addListener(graal_isolatethread_t *thread, dxfg_v4_observable_list_model_t *list,
                                                dxfg_v4_observable_list_model_listener_t *listener);

/**
 * @brief Removes the listener from the list. If the listener is not added to the list, nothing happens.
 *
 * Java: `ObservableListModel.removeListener(ObservableListModelListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableListModel_removeListener(graal_isolatethread_t *thread, dxfg_v4_observable_list_model_t *list,
                                                   dxfg_v4_observable_list_model_listener_t *listener);

/**
 * @brief Creates a listener of the changes of an order book model that calls the function with the user data.
 *
 * Add the listener to a model with dxfg_v4_OrderBookModel_addListener(). After the last handle of the listener
 * is released the function is not called anymore and none of its calls is in progress, so the user data can be
 * freed. Remove the listener from the model before releasing it, otherwise it stays added and does nothing, see
 * dxfg_v4_JavaObjectHandler_release().
 *
 * Java: `OrderBookModelListener.modelChanged(OrderBookModelListener.Change)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call after the model has changed; see
 * #dxfg_v4_order_book_model_listener_function for the thread and the ownership of its data.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it. It has to stay valid
 * until the listener is released.
 * @param[out] orderBookModelListener The pointer to the variable that receives the listener (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModelListener_new(graal_isolatethread_t *thread,
                                           dxfg_v4_order_book_model_listener_function user_func, void *user_data,
                                           DXFG_OUT dxfg_v4_order_book_model_listener_t **orderBookModelListener);

/**
 * @brief Adds the listener to the order book model. A closed model does not add listeners.
 *
 * Java: `OrderBookModel.addListener(OrderBookModelListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_addListener(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                           dxfg_v4_order_book_model_listener_t *listener);

/**
 * @brief Removes the listener from the order book model. If the listener is not added to the model, nothing happens.
 *
 * Java: `OrderBookModel.removeListener(OrderBookModelListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_OrderBookModel_removeListener(graal_isolatethread_t *thread, dxfg_v4_order_book_model_t *model,
                                              dxfg_v4_order_book_model_listener_t *listener);

/**
 * @brief Returns the events of the list as a new list of events.
 *
 * The result is a snapshot: it does not change when the list changes. An empty list gives a list of size 0.
 *
 * Java: the elements of `ObservableListModel` (a `java.util.List`).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list handle.
 * @param[out] events The pointer to the variable that receives the events (release them with
 * dxfg_v4_event_type_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableListModel_toArray(graal_isolatethread_t *thread, dxfg_v4_observable_list_model_t *list,
                                            DXFG_OUT dxfg_v4_event_type_list **events);

/**
 * @brief Creates a new builder of a model of the indexed events of the type.
 *
 * Java: `IndexedTxModel.newBuilder(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of the indexed events of the model.
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_newBuilder(graal_isolatethread_t* thread, dxfg_v4_event_clazz_t eventType,
                                          DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Returns the sources of the model.
 *
 * Java: `IndexedTxModel.getSources()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[out] sources The pointer to the variable that receives the sources (release them with
 * dxfg_v4_indexed_event_source_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_getSources(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source,
                                          DXFG_OUT dxfg_v4_indexed_event_source_list **sources);

/**
 * @brief Sets the sources of the model: it subscribes to the new ones and unsubscribes from the old ones.
 *
 * An empty list means all the sources.
 *
 * Java: `IndexedTxModel.setSources(Set)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] sources The new sources of the model; an empty list means all the sources.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_setSources(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source,
                                          dxfg_v4_indexed_event_source_list* sources);

/**
 * @brief Returns whether the listener gets all the pending transactions in one call.
 *
 * Java: `AbstractTxModel.isBatchProcessing()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[out] batchProcessing The pointer to the variable that receives 1 if enabled, 0 otherwise; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_isBatchProcessing(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source,
                                                 DXFG_OUT int32_t *batchProcessing);

/**
 * @brief Returns whether a snapshot is delivered as its resulting state (without the removed events).
 *
 * Java: `AbstractTxModel.isSnapshotProcessing()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[out] snapshotProcessing The pointer to the variable that receives 1 if enabled, 0 otherwise; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_isSnapshotProcessing(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source,
                                                    DXFG_OUT int32_t *snapshotProcessing);

/**
 * @brief Attaches the model to the feed.
 *
 * Java: `AbstractTxModel.attach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_attach(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source,
                                      dxfg_v4_feed_t* feed);

/**
 * @brief Detaches the model from the feed.
 *
 * Java: `AbstractTxModel.detach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_detach(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source,
                                      dxfg_v4_feed_t* feed);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.setAggregationPeriod(TimePeriod)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] aggregationPeriod The new aggregation period of the events.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_setAggregationPeriod(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source,
                                                    dxfg_v4_time_period_t* aggregationPeriod);

/**
 * @brief Closes the model: it is detached from all the feeds for good. Closing it again does nothing.
 *
 * Java: `AbstractTxModel.close()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_close(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_t* source);

/**
 * @brief Sets the sources (an empty list means all the sources).
 *
 * Java: `IndexedTxModel.Builder.withSources(Collection)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] sources The sources of the events; an empty list means all the sources.
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withSources(graal_isolatethread_t* thread,
                                                   dxfg_v4_indexed_tx_model_builder_t* source,
                                                   dxfg_v4_indexed_event_source_list* sources,
                                                   DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Builds the model (the symbol must be set).
 *
 * Java: `IndexedTxModel.Builder.build()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[out] indexedTxModel The pointer to the variable that receives the model (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_build(graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_builder_t* source,
                                             DXFG_OUT dxfg_v4_indexed_tx_model_t **indexedTxModel);

/**
 * @brief Sets whether the listener gets all the pending transactions in one call (1, default) or one by one (0).
 *
 * Java: `AbstractTxModel.Builder.withBatchProcessing(boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isBatchProcessing 1 to deliver all the pending transactions in one call, 0 to deliver them one by one.
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withBatchProcessing(
    graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_builder_t* source, int32_t isBatchProcessing,
    DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Sets whether a snapshot is delivered as its resulting state (1) or as it comes (0, default).
 *
 * Java: `AbstractTxModel.Builder.withSnapshotProcessing(boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isSnapshotProcessing 1 to deliver a snapshot as its resulting state (without the removed events), 0 to
 * deliver
 * it as it comes.
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withSnapshotProcessing(
    graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_builder_t* source, int32_t isSnapshotProcessing,
    DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Sets the feed the model is attached to.
 *
 * Java: `AbstractTxModel.Builder.withFeed(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] feed The feed the model is attached to.
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withFeed(graal_isolatethread_t* thread,
                                                dxfg_v4_indexed_tx_model_builder_t* source, dxfg_v4_feed_t* feed,
                                                DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.Builder.withAggregationPeriod(TimePeriod)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] aggregationPeriod The aggregation period of the events.
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withAggregationPeriod(
    graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_builder_t* source, dxfg_v4_time_period_t* aggregationPeriod,
    DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Sets the symbol of the subscription.
 *
 * Java: `AbstractTxModel.Builder.withSymbol(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] symbol The symbol to subscribe for (a dxfg_v4_symbol_t).
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withSymbol(graal_isolatethread_t* thread,
                                                  dxfg_v4_indexed_tx_model_builder_t* source, dxfg_v4_symbol_t* symbol,
                                                  DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Sets the listener of the transactions.
 *
 * Java: `AbstractTxModel.Builder.withListener(TxModelListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] listener The listener of the transactions (see dxfg_v4_TxModelListener_new()).
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withListener(
    graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_builder_t* source, dxfg_v4_tx_model_listener_t* listener,
    DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Sets the executor that calls the listener (by default, the executor of the feed).
 *
 * The executor of a shared event queue (dxfg_v4_EventQueue_getExecutor()) makes the polls of the queue call the
 * listeners; the executor of a subscription attached to a queue is refused with `IllegalStateException` (see
 * @ref V4_EventQueue).
 *
 * Java: `AbstractTxModel.Builder.withExecutor(Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] executor The executor that calls the listener.
 * @param[out] indexedTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_IndexedTxModel_Builder_withExecutor(
    graal_isolatethread_t* thread, dxfg_v4_indexed_tx_model_builder_t* source, dxfg_v4_executor_t* executor,
    DXFG_OUT dxfg_v4_indexed_tx_model_builder_t **indexedTxModelBuilder);

/**
 * @brief Creates a new builder of a model of the time-series events of the type.
 *
 * Java: `TimeSeriesTxModel.newBuilder(Class)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of the time-series events of the model.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_newBuilder(
    graal_isolatethread_t* thread, dxfg_v4_event_clazz_t eventType,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Returns the time the model is subscribed from.
 *
 * Java: `TimeSeriesTxModel.getFromTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[out] fromTime The pointer to the variable that receives the time in milliseconds since the epoch (INT64_MAX if
 * the model is not subscribed); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_getFromTime(graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_t* source,
                                              DXFG_OUT int64_t *fromTime);

/**
 * @brief Sets the time the model is subscribed from (INT64_MAX unsubscribes); the model gets a new snapshot.
 *
 * Java: `TimeSeriesTxModel.setFromTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] fromTime The time to subscribe from, in milliseconds since the epoch; INT64_MAX unsubscribes.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_setFromTime(graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_t* source,
                                              int64_t fromTime);

/**
 * @brief Returns whether the listener gets all the pending transactions in one call.
 *
 * Java: `AbstractTxModel.isBatchProcessing()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[out] batchProcessing The pointer to the variable that receives 1 if enabled, 0 otherwise; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_isBatchProcessing(graal_isolatethread_t* thread,
                                                    dxfg_v4_time_series_tx_model_t* source,
                                                    DXFG_OUT int32_t *batchProcessing);

/**
 * @brief Returns whether a snapshot is delivered as its resulting state (without the removed events).
 *
 * Java: `AbstractTxModel.isSnapshotProcessing()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[out] snapshotProcessing The pointer to the variable that receives 1 if enabled, 0 otherwise; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_isSnapshotProcessing(graal_isolatethread_t* thread,
                                                       dxfg_v4_time_series_tx_model_t* source,
                                                       DXFG_OUT int32_t *snapshotProcessing);

/**
 * @brief Attaches the model to the feed.
 *
 * Java: `AbstractTxModel.attach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_attach(graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_t* source,
                                         dxfg_v4_feed_t* feed);

/**
 * @brief Detaches the model from the feed.
 *
 * Java: `AbstractTxModel.detach(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_detach(graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_t* source,
                                         dxfg_v4_feed_t* feed);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.setAggregationPeriod(TimePeriod)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] aggregationPeriod The new aggregation period of the events.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_setAggregationPeriod(graal_isolatethread_t* thread,
                                                       dxfg_v4_time_series_tx_model_t* source,
                                                       dxfg_v4_time_period_t* aggregationPeriod);

/**
 * @brief Closes the model: it is detached from all the feeds for good. Closing it again does nothing.
 *
 * Java: `AbstractTxModel.close()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_close(graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_t* source);

/**
 * @brief Sets the time to subscribe from, in milliseconds since the epoch (by default INT64_MAX: not subscribed).
 *
 * Java: `TimeSeriesTxModel.Builder.withFromTime(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] fromTime The time to subscribe from, in milliseconds since the epoch (INT64_MAX: not subscribed).
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withFromTime(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, int64_t fromTime,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets the order of the events of a snapshot.
 *
 * Java: `TimeSeriesTxModel.Builder.withSorting(TimeSeriesTxModel.SortOrder)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] sortOrder The order of the events of a snapshot.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withSorting(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source,
    dxfg_v4_TimeSeriesTxModelSortOrder_t sortOrder,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets whether the removed events are dropped (1) or delivered (0, default).
 *
 * Java: `TimeSeriesTxModel.Builder.ignoreRemoveEvents(boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] ignoreRemoveEvents 1 to drop the removed events, 0 to deliver them.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_ignoreRemoveEvents(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, int32_t ignoreRemoveEvents,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets whether the events with an index not greater than the last processed one are dropped (1)
 * or delivered (0, default).
 *
 * Java: `TimeSeriesTxModel.Builder.ignoreEventsFromPast(boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] ignoreEventFromPast 1 to drop the events with an index not greater than the last processed one, 0 to
 * deliver
 * them.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_ignoreEventsFromPast(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, int32_t ignoreEventFromPast,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets whether a new bar comes after the previous (closing) one, to confirm it (1), or alone (0, default).
 *
 * With 1, the snapshot processing and the ascending sorting are on.
 *
 * Java: `TimeSeriesTxModel.Builder.withConfirmationTick(boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] confirmationTick 1 to deliver a new bar after the previous (closing) one, to confirm it, 0 to deliver it
 * alone.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withConfirmationTick(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, int32_t confirmationTick,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Accepted for compatibility: the pointer is ignored, the model gets no event peek (a Java `Consumer`).
 *
 * For an event peek that works, use dxfg_v4_TimeSeriesTxModel_Builder_withEventPeekFunction().
 *
 * Java: `TimeSeriesTxModel.Builder.withEventPeek(Consumer)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] eventPeek Ignored: pass NULL.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withEventPeek(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, void* eventPeek,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief The event peek function of a time series model, see
 * dxfg_v4_TimeSeriesTxModel_Builder_withEventPeekFunction().
 *
 * The model calls it for each event of a transaction or a snapshot, in the order of the list, before its listener gets
 * the list; it runs on the thread of the model's executor (the thread that calls the listener), not concurrently with
 * itself for one model. The event is valid during the call only: the SDK frees it after the function returns.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] event The event; do not keep or release it.
 * @param[in] userData The user data given to dxfg_v4_TimeSeriesTxModel_Builder_withEventPeekFunction().
 * @experimental
 * @since 3.13.0
 */
typedef void (*dxfg_v4_time_series_tx_model_event_peek_function)(graal_isolatethread_t *thread,
                                                                  dxfg_v4_event_type_t *event, void *userData);

/**
 * @brief Sets the event peek of the model: a function that sees each event before the listener gets the list.
 *
 * The function is called as #dxfg_v4_time_series_tx_model_event_peek_function says. NULL removes the event peek set
 * before, and so does dxfg_v4_TimeSeriesTxModel_Builder_withEventPeek(); the last call before
 * dxfg_v4_TimeSeriesTxModel_Builder_build() wins. The library reports the function with
 * #DXFG_V4_FEATURE_TIME_SERIES_TX_MODEL_EVENT_PEEK.
 *
 * The event peek is a native listener without a handle of its own: it is not called inside the tasks of a closed event
 * queue (as the listeners are not) and after dxfg_v4_NativeListener_deactivateAll(); the release of the builder or
 * model handles does not stop it. The function pointer and the user data must stay valid while the model can still
 * deliver events: until the model is closed and its executor has no pending notifications, or until
 * dxfg_v4_NativeListener_deactivateAll() returns.
 *
 * Java: `TimeSeriesTxModel.Builder.withEventPeek(Consumer)` with a `Consumer` that calls `eventPeek`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle; must not be NULL.
 * @param[in] eventPeek The event peek function, or NULL for none.
 * @param[in] userData The user data passed to `eventPeek`, may be NULL.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withEventPeekFunction(
    graal_isolatethread_t *thread, dxfg_v4_time_series_tx_model_builder_t *source,
    dxfg_v4_time_series_tx_model_event_peek_function eventPeek, void *userData,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Builds the model (the symbol must be set).
 *
 * Java: `TimeSeriesTxModel.Builder.build()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[out] timeSeriesTxModel The pointer to the variable that receives the model (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_build(graal_isolatethread_t* thread,
                                                dxfg_v4_time_series_tx_model_builder_t* source,
                                                DXFG_OUT dxfg_v4_time_series_tx_model_t **timeSeriesTxModel);

/**
 * @brief Sets whether the listener gets all the pending transactions in one call (1, default) or one by one (0).
 *
 * Java: `AbstractTxModel.Builder.withBatchProcessing(boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isBatchProcessing 1 to deliver all the pending transactions in one call, 0 to deliver them one by one.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withBatchProcessing(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, int32_t isBatchProcessing,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets whether a snapshot is delivered as its resulting state (1) or as it comes (0, default).
 *
 * Java: `AbstractTxModel.Builder.withSnapshotProcessing(boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isSnapshotProcessing 1 to deliver a snapshot as its resulting state (without the removed events), 0 to
 * deliver
 * it as it comes.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withSnapshotProcessing(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, int32_t isSnapshotProcessing,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets the feed the model is attached to.
 *
 * Java: `AbstractTxModel.Builder.withFeed(DXFeed)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] feed The feed the model is attached to.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withFeed(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, dxfg_v4_feed_t* feed,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.Builder.withAggregationPeriod(TimePeriod)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] aggregationPeriod The aggregation period of the events.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withAggregationPeriod(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source,
    dxfg_v4_time_period_t* aggregationPeriod,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets the symbol of the subscription.
 *
 * Java: `AbstractTxModel.Builder.withSymbol(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] symbol The symbol to subscribe for (a dxfg_v4_symbol_t).
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withSymbol(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, dxfg_v4_symbol_t* symbol,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets the listener of the transactions.
 *
 * Java: `AbstractTxModel.Builder.withListener(TxModelListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] listener The listener of the transactions (see dxfg_v4_TxModelListener_new()).
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withListener(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source,
    dxfg_v4_tx_model_listener_t* listener, DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Sets the executor that calls the listener (by default, the executor of the feed).
 *
 * The executor of a shared event queue (dxfg_v4_EventQueue_getExecutor()) makes the polls of the queue call the
 * listeners; the executor of a subscription attached to a queue is refused with `IllegalStateException` (see
 * @ref V4_EventQueue).
 *
 * Java: `AbstractTxModel.Builder.withExecutor(Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] executor The executor that calls the listener.
 * @param[out] timeSeriesTxModelBuilder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TimeSeriesTxModel_Builder_withExecutor(
    graal_isolatethread_t* thread, dxfg_v4_time_series_tx_model_builder_t* source, dxfg_v4_executor_t* executor,
    DXFG_OUT dxfg_v4_time_series_tx_model_builder_t **timeSeriesTxModelBuilder);

/**
 * @brief Creates a listener that calls the function with the user data.
 *
 * After the last handle of the listener is released (dxfg_v4_JavaObjectHandler_release()), the function is not called
 * anymore and none of its calls is in progress, so the user data can be freed.
 *
 * Java: `TxModelListener.eventsReceived(IndexedEventSource, List, boolean)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] function_eventsReceived The function to call when transactions are received; see
 * #dxfg_v4_TxModelListener_function_eventsReceived for the thread and the ownership of its data.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it. It has to stay valid
 * until the listener is released.
 * @param[out] txModelListener The pointer to the variable that receives the listener (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_TxModelListener_new(graal_isolatethread_t* thread,
                                    dxfg_v4_TxModelListener_function_eventsReceived function_eventsReceived,
                                    void* user_data, DXFG_OUT dxfg_v4_tx_model_listener_t **txModelListener);

/** @} */ // end of V4_EventModel

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_EVENT_MODEL_H
