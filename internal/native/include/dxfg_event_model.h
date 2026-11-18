// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK EventModel functions and types declarations
 */

#ifndef DXFG_EVENT_MODEL_H
#define DXFG_EVENT_MODEL_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"
#include "dxfg_javac.h"
#include "dxfg_feed.h"

/** @defgroup EventModel EventModel
 *  @brief The models of the events of one or more symbols: the order book, indexed, time series and transaction
 *  models.
 *
 *  Threads: dxfg_order_book_model_t, dxfg_indexed_event_model_t and dxfg_time_series_event_model_t are not
 *  thread-safe. Only their attach, detach and close functions may be called from any thread. Call any other function
 *  of an attached model, and of its lists, only from inside a listener call or from the thread of the executor that
 *  calls the listeners, or synchronize the calls. The transaction models (dxfg_indexed_tx_model_t and
 *  dxfg_time_series_tx_model_t) are thread-safe.
 *  @{
 */

/**
 * @brief A handle of a Java OrderBookModel object; release it with dxfg_JavaObjectHandler_release().
 *
 * The model of an order book: it keeps the buy and the sell orders of its symbol in two lists, see
 * dxfg_OrderBookModel_getBuyOrders() and dxfg_OrderBookModel_getSellOrders().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/market/OrderBookModel.html">Javadoc</a>
 */
typedef struct dxfg_order_book_model_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_order_book_model_t;

/**
 * @brief A handle of a Java IndexedEventModel object; release it with dxfg_JavaObjectHandler_release().
 *
 * The model of a list of indexed events of its symbol, see dxfg_IndexedEventModel_getEventsList().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/IndexedEventModel.html">Javadoc</a>
 */
typedef struct dxfg_indexed_event_model_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_indexed_event_model_t;

/**
 * @brief A handle of a Java TimeSeriesEventModel object; release it with dxfg_JavaObjectHandler_release().
 *
 * The model of a list of time-series events of its symbol, see dxfg_TimeSeriesEventModel_getEventsList().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/TimeSeriesEventModel.html">Javadoc</a>
 */
typedef struct dxfg_time_series_event_model_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_time_series_event_model_t;

/**
 * @brief A handle of a Java ObservableListModel object; release it with dxfg_JavaObjectHandler_release().
 *
 * A live list of events that notifies its listeners about its changes. The handle refers to the same list that
 * the model keeps: the list changes as the model receives events.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/ObservableListModel.html">Javadoc</a>
 */
typedef struct dxfg_observable_list_model_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_observable_list_model_t;

/**
 * @brief A handle of a Java ObservableListModelListener object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_ObservableListModelListener_new(). The last handle of the listener that is released
 * deactivates the listener: its function is not called anymore.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/ObservableListModelListener.html">Javadoc</a>
 */
typedef struct dxfg_observable_list_model_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_observable_list_model_listener_t;

/**
 * @brief The function of a dxfg_observable_list_model_listener_t: called after a list of events has changed.
 *
 * The function is called by the executor of the model that owns the list (by default the executor of the feed the
 * model is attached to, see dxfg_IndexedEventModel_setExecutor()), not by the thread that added the listener. Calls
 * of one model are never concurrent. It is safe to query the model and its lists from inside the call.
 *
 * The SDK frees the passed list right after the function returns: the callee must not release it and must copy
 * whatever it needs to keep.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] orders All the events of the changed list, not only the changed ones. The list belongs to the SDK and
 * is valid during the call only.
 * @param[in] user_data The pointer that was passed to dxfg_ObservableListModelListener_new(); the SDK does not use it.
 */
typedef void (*dxfg_observable_list_model_listener_function)(graal_isolatethread_t *thread, dxfg_event_type_list* orders, void *user_data);

/**
 * @brief A handle of a Java OrderBookModelListener object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_OrderBookModelListener_new(). The last handle of the listener that is released
 * deactivates the listener: its function is not called anymore.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/market/OrderBookModelListener.html">Javadoc</a>
 */
typedef struct dxfg_order_book_model_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_order_book_model_listener_t;

/**
 * @brief The function of a dxfg_order_book_model_listener_t: called after an order book model has changed.
 *
 * The function is called by the executor of the model (by default the executor of the feed the model is attached to,
 * see dxfg_OrderBookModel_setExecutor()), not by the thread that added the listener. The calls of one model are never
 * concurrent. The notification says which model changed, not what changed: query the lists
 * dxfg_OrderBookModel_getBuyOrders() and dxfg_OrderBookModel_getSellOrders() to get the new orders.
 *
 * The handle of the model is created for the call and released by the SDK right after the function returns: the callee
 * must not release it. To keep the model, clone the handle with dxfg_JavaObjectHandler_clone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The model that has changed. The handle is valid during the call only.
 * @param[in] user_data The pointer that was passed to dxfg_OrderBookModelListener_new(); the SDK does not use it.
 */
typedef void (*dxfg_order_book_model_listener_function)(graal_isolatethread_t *thread, dxfg_order_book_model_t* model, void *user_data);

/**
 * @brief Creates a new order book model.
 *
 * The model is not attached to any feed and not subscribed to any symbol. Its filter is
 * ALL and its lot size is 1. Set the symbol with dxfg_OrderBookModel_setSymbol() and attach the model
 * with dxfg_OrderBookModel_attach() to start receiving orders.
 *
 * Java: `new OrderBookModel()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The model (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_new()")
dxfg_order_book_model_t*          dxfg_OrderBookModel_new(graal_isolatethread_t *thread);

/**
 * @brief Creates a new model of indexed events of the event type.
 *
 * The model is not attached to any feed and not subscribed to any symbol. Set the symbol with
 * dxfg_IndexedEventModel_setSymbol() and attach the model with dxfg_IndexedEventModel_attach() to start
 * receiving events.
 *
 * Java: `new IndexedEventModel(Class)`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] clazz The event type of the model; it has to be an indexed event type.
 * @return The model (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_new()")
dxfg_indexed_event_model_t*       dxfg_IndexedEventModel_new(graal_isolatethread_t *thread, dxfg_event_clazz_t clazz);

/**
 * @brief Creates a new model of time-series events of the event type.
 *
 * The model is not attached to any feed and not subscribed to any symbol (its from-time is INT64_MAX). Set
 * the symbol and the from-time with dxfg_TimeSeriesEventModel_setSymbol() and
 * dxfg_TimeSeriesEventModel_setFromTime() and attach the model with dxfg_TimeSeriesEventModel_attach() to start
 * receiving events.
 *
 * Java: `new TimeSeriesEventModel(Class)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] clazz The event type of the model; it has to be a time-series event type.
 * @return The model (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_new()")
dxfg_time_series_event_model_t*   dxfg_TimeSeriesEventModel_new(graal_isolatethread_t *thread, dxfg_event_clazz_t clazz);

/**
 * @brief The filter of an order book model: which scopes of the order events are shown in the book.
 *
 * If more than one scope is allowed, the events of a more detailed scope hide the events of a less detailed one: an
 * aggregate order of an exchange and a market maker hides the regional order of that exchange, which in turn hides the
 * composite order.
 *
 * Java: `OrderBookModelFilter`.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/market/OrderBookModelFilter.html">Javadoc</a>
 */
typedef enum dxfg_order_book_model_filter_t {
    /** Only the composite orders (the best bid or offer of the market) are shown. */
    COMPOSITE,
    /** Only the regional orders (the best bid or offer of an exchange) are shown. */
    REGIONAL,
    /** Only the aggregate orders (a price level or a market maker's best) are shown. */
    AGGREGATE,
    /** Only the individual orders are shown. */
    ORDER,
    /** The composite and the regional orders are shown. */
    COMPOSITE_REGIONAL,
    /** All the orders except the individual ones are shown. */
    COMPOSITE_REGIONAL_AGGREGATE,
    /** The orders of all scopes are shown. */
    ALL,
} dxfg_order_book_model_filter_t;

/**
 * @brief Attaches the order book model to the feed: the model starts receiving the events of its symbol.
 *
 * The feed keeps an attached model: releasing the handle of the model does not detach it. Detach the
 * model with dxfg_OrderBookModel_detach() or close it with dxfg_OrderBookModel_close().
 *
 * Java: `OrderBookModel.attach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_attach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_attach()")
int32_t                                       dxfg_OrderBookModel_attach(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_feed_t *feed);

/**
 * @brief Detaches the order book model from the feed.
 *
 * Java: `OrderBookModel.detach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_detach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_detach()")
int32_t                                       dxfg_OrderBookModel_detach(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_feed_t *feed);

/**
 * @brief Closes the order book model: it becomes permanently detached from all feeds and removes all its listeners.
 *
 * Java: `OrderBookModel.close()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_close()")
int32_t                                       dxfg_OrderBookModel_close(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Returns the executor that calls the listeners of the order book model.
 *
 * Java: `OrderBookModel.getExecutor()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_getExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return The executor (release it with dxfg_JavaObjectHandler_release()), or NULL if the default executor
 * of the feed is used or on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception
 * was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_getExecutor()")
dxfg_executor_t*                              dxfg_OrderBookModel_getExecutor(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Sets the executor that calls the listeners of the order book model.
 *
 * Java: `OrderBookModel.setExecutor(Executor)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_setExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] executor The executor, or NULL to use the default executor of the feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_setExecutor()")
int32_t                                       dxfg_OrderBookModel_setExecutor(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_executor_t *executor);

/**
 * @brief Clears the subscription symbol and, with it, all the orders of the model.
 *
 * This is a shortcut for dxfg_OrderBookModel_setSymbol() with NULL.
 *
 * Java: `OrderBookModel.clear()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_clear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_clear()")
int32_t                                       dxfg_OrderBookModel_clear(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Attaches the indexed event model to the feed: the model starts receiving the events of its symbol.
 *
 * The feed keeps an attached model: releasing the handle of the model does not detach it. Detach the
 * model with dxfg_IndexedEventModel_detach() or close it with dxfg_IndexedEventModel_close().
 *
 * Java: `IndexedEventModel.attach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_attach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_attach()")
int32_t                                       dxfg_IndexedEventModel_attach(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_feed_t *feed);

/**
 * @brief Detaches the indexed event model from the feed.
 *
 * Java: `IndexedEventModel.detach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_detach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_detach()")
int32_t                                       dxfg_IndexedEventModel_detach(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_feed_t *feed);

/**
 * @brief Closes the indexed event model: it becomes permanently detached from all feeds and removes all its listeners.
 *
 * Java: `IndexedEventModel.close()`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_close()")
int32_t                                       dxfg_IndexedEventModel_close(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);

/**
 * @brief Returns the executor that calls the listeners of the indexed event model.
 *
 * Java: `IndexedEventModel.getExecutor()`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_getExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return The executor (release it with dxfg_JavaObjectHandler_release()), or NULL if the default executor
 * of the feed is used or on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception
 * was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_getExecutor()")
dxfg_executor_t*                              dxfg_IndexedEventModel_getExecutor(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);

/**
 * @brief Sets the executor that calls the listeners of the indexed event model.
 *
 * Java: `IndexedEventModel.setExecutor(Executor)`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_setExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] executor The executor, or NULL to use the default executor of the feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_setExecutor()")
int32_t                                       dxfg_IndexedEventModel_setExecutor(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_executor_t *executor);

/**
 * @brief Clears the subscription symbol and, with it, all the events of the model.
 *
 * This is a shortcut for dxfg_IndexedEventModel_setSymbol() with NULL.
 *
 * Java: `IndexedEventModel.clear()`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_clear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_clear()")
int32_t                                       dxfg_IndexedEventModel_clear(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);

/**
 * @brief Attaches the time-series event model to the feed: the model starts receiving the events of its symbol.
 *
 * The feed keeps an attached model: releasing the handle of the model does not detach it. Detach the
 * model with dxfg_TimeSeriesEventModel_detach() or close it with dxfg_TimeSeriesEventModel_close().
 *
 * Java: `TimeSeriesEventModel.attach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_attach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_attach()")
int32_t                                       dxfg_TimeSeriesEventModel_attach(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_feed_t *feed);

/**
 * @brief Detaches the time-series event model from the feed.
 *
 * Java: `TimeSeriesEventModel.detach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_detach().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_detach()")
int32_t                                       dxfg_TimeSeriesEventModel_detach(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_feed_t *feed);

/**
 * @brief Closes the time-series event model: it becomes permanently detached from all feeds and removes all its
 * listeners.
 *
 * Java: `TimeSeriesEventModel.close()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_close()")
int32_t                                       dxfg_TimeSeriesEventModel_close(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);

/**
 * @brief Returns the executor that calls the listeners of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.getExecutor()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_getExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return The executor (release it with dxfg_JavaObjectHandler_release()), or NULL if the default executor
 * of the feed is used or on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception
 * was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_getExecutor()")
dxfg_executor_t*                              dxfg_TimeSeriesEventModel_getExecutor(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);

/**
 * @brief Sets the executor that calls the listeners of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.setExecutor(Executor)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_setExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] executor The executor, or NULL to use the default executor of the feed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_setExecutor()")
int32_t                                       dxfg_TimeSeriesEventModel_setExecutor(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_executor_t *executor);

/**
 * @brief Clears the subscription symbol, sets the from-time to INT64_MAX and, with it, clears all the events
 * of the model.
 *
 * Java: `TimeSeriesEventModel.clear()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_clear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_clear()")
int32_t                                       dxfg_TimeSeriesEventModel_clear(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);

/**
 * @brief Returns the filter of the order book model.
 *
 * Java: `OrderBookModel.getFilter()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_getFilter().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return The filter, a value of #dxfg_order_book_model_filter_t, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_getFilter()")
int32_t                                       dxfg_OrderBookModel_getFilter(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Sets the filter of the order book model.
 *
 * Java: `OrderBookModel.setFilter(OrderBookModelFilter)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_setFilter().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] filter The filter.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_setFilter()")
int32_t                                       dxfg_OrderBookModel_setFilter(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_order_book_model_filter_t filter);

/**
 * @brief Returns the subscription symbol of the order book model.
 *
 * Java: `OrderBookModel.getSymbol()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_getSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return The symbol (a UTF-8 string; release it with dxfg_String_release()), or NULL if the model is not
 * subscribed (the default) or on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception
 * was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_getSymbol()")
const char*                                   dxfg_OrderBookModel_getSymbol(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Sets the subscription symbol of the order book model.
 *
 * Java: `OrderBookModel.setSymbol(String)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_setSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] symbol The symbol, a UTF-8 string, or NULL to unsubscribe (that removes the orders of the model).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_setSymbol()")
int32_t                                       dxfg_OrderBookModel_setSymbol(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, const char* symbol);

/**
 * @brief Returns the subscription symbol of the indexed event model.
 *
 * Java: `IndexedEventModel.getSymbol()`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_getSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return The symbol (release it with dxfg_Symbol_release()), or NULL if the model is not subscribed (the
 * default). On error the result is `(dxfg_symbol_t *)-1`, not NULL: use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_getSymbol()")
dxfg_symbol_t*                                dxfg_IndexedEventModel_getSymbol(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);

/**
 * @brief Sets the subscription symbol of the indexed event model.
 *
 * Java: `IndexedEventModel.setSymbol(Object)`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_setSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] symbol The symbol, or NULL to unsubscribe (that removes the events of the model).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_setSymbol()")
int32_t                                       dxfg_IndexedEventModel_setSymbol(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_symbol_t* symbol);

/**
 * @brief Returns the subscription symbol of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.getSymbol()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_getSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return The symbol (release it with dxfg_Symbol_release()), or NULL if the model is not subscribed (the
 * default). On error the result is `(dxfg_symbol_t *)-1`, not NULL: use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_getSymbol()")
dxfg_symbol_t*                                dxfg_TimeSeriesEventModel_getSymbol(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);

/**
 * @brief Sets the subscription symbol of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.setSymbol(Object)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_setSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] symbol The symbol, or NULL to unsubscribe (that removes the events of the model).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_setSymbol()")
int32_t                                       dxfg_TimeSeriesEventModel_setSymbol(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_symbol_t* symbol);

/**
 * @brief Returns the lot size of the order book model.
 *
 * The lot size is a multiplier applied to the sizes of the composite, regional and aggregate orders (1 by
 * default).
 *
 * Java: `OrderBookModel.getLotSize()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_getLotSize().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return The lot size, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception
 * was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_getLotSize()")
int32_t                                       dxfg_OrderBookModel_getLotSize(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Sets the lot size of the order book model.
 *
 * A lot size less than 1 is an error.
 *
 * Java: `OrderBookModel.setLotSize(int)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_setLotSize().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] lotSize The lot size multiplier, 1 or greater.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_setLotSize()")
int32_t                                       dxfg_OrderBookModel_setLotSize(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, int32_t lotSize);

/**
 * @brief Returns the size limit of the indexed event model.
 *
 * Java: `IndexedEventModel.getSizeLimit()`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_getSizeLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return The size limit (INT32_MAX means no limit, it is the default), or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_getSizeLimit()")
int32_t                                       dxfg_IndexedEventModel_getSizeLimit(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);

/**
 * @brief Sets the size limit of the indexed event model.
 *
 * A negative size limit is an error.
 *
 * Java: `IndexedEventModel.setSizeLimit(int)`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_setSizeLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @param[in] sizeLimit The size limit, not negative. When the limit is exceeded, the first events of the list
 * are dropped.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_setSizeLimit()")
int32_t                                       dxfg_IndexedEventModel_setSizeLimit(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, int32_t sizeLimit);

/**
 * @brief Returns the size limit of the time-series event model.
 *
 * Java: `TimeSeriesEventModel.getSizeLimit()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_getSizeLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return The size limit (INT32_MAX means no limit, it is the default), or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_getSizeLimit()")
int32_t                                       dxfg_TimeSeriesEventModel_getSizeLimit(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);

/**
 * @brief Sets the size limit of the time-series event model.
 *
 * A negative size limit is an error.
 *
 * Java: `TimeSeriesEventModel.setSizeLimit(int)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_setSizeLimit().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] sizeLimit The size limit, not negative. When the limit is exceeded, the first events of the list
 * are dropped.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_setSizeLimit()")
int32_t                                       dxfg_TimeSeriesEventModel_setSizeLimit(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, int32_t sizeLimit);

/**
 * @brief Returns the time from which the time-series event model subscribes for events.
 *
 * Java: `TimeSeriesEventModel.getFromTime()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_getFromTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return The time in milliseconds since the epoch (INT64_MAX means not subscribed, it is the default), or -1 on
 * error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_getFromTime()")
int64_t                                       dxfg_TimeSeriesEventModel_getFromTime(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);

/**
 * @brief Sets the time from which the time-series event model subscribes for events.
 *
 * Java: `TimeSeriesEventModel.setFromTime(long)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_setFromTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @param[in] fromTime The time in milliseconds since the epoch; INT64_MAX unsubscribes.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_setFromTime()")
int32_t                                       dxfg_TimeSeriesEventModel_setFromTime(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, int64_t fromTime);

/**
 * @brief Returns the buy orders (the bid side) of the order book model.
 *
 * The orders are arranged by price: the first order is the top of the bid side (the highest price).
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_ObservableListModel_toArray().
 *
 * Java: `OrderBookModel.getBuyOrders()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_getBuyOrders().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return The list (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_getBuyOrders()")
dxfg_observable_list_model_t*                 dxfg_OrderBookModel_getBuyOrders(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Returns the sell orders (the offer side) of the order book model.
 *
 * The orders are arranged by price: the first order is the top of the offer side (the lowest price).
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_ObservableListModel_toArray().
 *
 * Java: `OrderBookModel.getSellOrders()`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_getSellOrders().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @return The list (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_getSellOrders()")
dxfg_observable_list_model_t*                 dxfg_OrderBookModel_getSellOrders(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);

/**
 * @brief Returns the list of events of the indexed event model.
 *
 * The events are ordered by the id of their source and by their index.
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_ObservableListModel_toArray().
 *
 * Java: `IndexedEventModel.getEventsList()`.
 *
 * @deprecated Use dxfg_v4_IndexedEventModel_getEventsList().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The indexed event model handle.
 * @return The list (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventModel_getEventsList()")
dxfg_observable_list_model_t*                 dxfg_IndexedEventModel_getEventsList(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);

/**
 * @brief Returns the list of events of the time-series event model.
 *
 * The events are ordered by their time.
 * The handle refers to the live list of the model: its content changes as the model receives events. The function
 * creates a new handle on every call. Read the list with dxfg_ObservableListModel_toArray().
 *
 * Java: `TimeSeriesEventModel.getEventsList()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesEventModel_getEventsList().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The time-series event model handle.
 * @return The list (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesEventModel_getEventsList()")
dxfg_observable_list_model_t*                 dxfg_TimeSeriesEventModel_getEventsList(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);

/**
 * @brief Creates a listener of the changes of a list that calls the function with the user data.
 *
 * Add the listener to a list with dxfg_ObservableListModel_addListener(). After the last handle of the listener
 * is released the function is not called anymore and none of its calls is in progress, so the user data can be
 * freed. Remove the listener from the list before releasing it, otherwise it stays added and does nothing, see
 * dxfg_JavaObjectHandler_release().
 *
 * Java: `ObservableListModelListener.modelChanged(ObservableListModelListener.Change)`.
 *
 * @deprecated Use dxfg_v4_ObservableListModelListener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call after the list has changed; see
 * #dxfg_observable_list_model_listener_function for the thread and the ownership of its data.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it. It has to stay valid
 * until the listener is released.
 * @return The listener (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableListModelListener_new()")
dxfg_observable_list_model_listener_t*        dxfg_ObservableListModelListener_new(graal_isolatethread_t *thread, dxfg_observable_list_model_listener_function user_func, void *user_data);

/**
 * @brief Adds the listener to the list.
 *
 * Java: `ObservableListModel.addListener(ObservableListModelListener)`.
 *
 * @deprecated Use dxfg_v4_ObservableListModel_addListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableListModel_addListener()")
int32_t                                       dxfg_ObservableListModel_addListener(graal_isolatethread_t *thread, dxfg_observable_list_model_t *list, dxfg_observable_list_model_listener_t *listener);

/**
 * @brief Removes the listener from the list. If the listener is not added to the list, nothing happens.
 *
 * Java: `ObservableListModel.removeListener(ObservableListModelListener)`.
 *
 * @deprecated Use dxfg_v4_ObservableListModel_removeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableListModel_removeListener()")
int32_t                                       dxfg_ObservableListModel_removeListener(graal_isolatethread_t *thread, dxfg_observable_list_model_t *list, dxfg_observable_list_model_listener_t *listener);

/**
 * @brief Creates a listener of the changes of an order book model that calls the function with the user data.
 *
 * Add the listener to a model with dxfg_OrderBookModel_addListener(). After the last handle of the listener
 * is released the function is not called anymore and none of its calls is in progress, so the user data can be
 * freed. Remove the listener from the model before releasing it, otherwise it stays added and does nothing, see
 * dxfg_JavaObjectHandler_release().
 *
 * Java: `OrderBookModelListener.modelChanged(OrderBookModelListener.Change)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModelListener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call after the model has changed; see
 * #dxfg_order_book_model_listener_function for the thread and the ownership of its data.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it. It has to stay valid
 * until the listener is released.
 * @return The listener (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModelListener_new()")
dxfg_order_book_model_listener_t*             dxfg_OrderBookModelListener_new(graal_isolatethread_t *thread, dxfg_order_book_model_listener_function user_func, void *user_data);

/**
 * @brief Adds the listener to the order book model. A closed model does not add listeners.
 *
 * Java: `OrderBookModel.addListener(OrderBookModelListener)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_addListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_addListener()")
int32_t                                       dxfg_OrderBookModel_addListener(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_order_book_model_listener_t *listener);

/**
 * @brief Removes the listener from the order book model. If the listener is not added to the model, nothing happens.
 *
 * Java: `OrderBookModel.removeListener(OrderBookModelListener)`.
 *
 * @deprecated Use dxfg_v4_OrderBookModel_removeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] model The order book model handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OrderBookModel_removeListener()")
int32_t                                       dxfg_OrderBookModel_removeListener(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_order_book_model_listener_t *listener);

/**
 * @brief Returns the events of the list as a new list of events.
 *
 * The result is a snapshot: it does not change when the list changes. An empty list gives a list of size 0.
 *
 * Java: the elements of `ObservableListModel` (a `java.util.List`).
 *
 * @deprecated Use dxfg_v4_ObservableListModel_toArray().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list handle.
 * @return The events (release them with dxfg_CList_EventType_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableListModel_toArray()")
dxfg_event_type_list*                         dxfg_ObservableListModel_toArray(graal_isolatethread_t *thread, dxfg_observable_list_model_t *list);

/**
 * @brief The order in which the events of a snapshot are delivered by a time-series transactional model.
 *
 * Java: `TimeSeriesTxModel.SortOrder`.
 */
typedef enum dxfg_TimeSeriesTxModelSortOrder_t {
    /// Indicates that no sorting should be applied. Events are passed in the order they were received.
    NONE,
    /// Indicates that sorting should be in ascending order by index.
    ASCENDING,
    /// Indicates that sorting should be in descending order by index.
    DESCENDING,
} dxfg_TimeSeriesTxModelSortOrder_t;

/**
 * @name Transactional models
 * IndexedTxModel and TimeSeriesTxModel deliver the complete transactions of indexed events to a listener.
 * Every builder function returns a new handle of the same builder: release every handle the functions return.
 * @{
 */

/**
 * @brief A handle of a Java IndexedTxModel object; release it with dxfg_JavaObjectHandler_release().
 *
 * The model is an incremental model of indexed events: it delivers the complete transactions of the events of its
 * symbol and sources to its listener. Close the model with dxfg_IndexedTxModel_close() before releasing the handle:
 * an attached model stays attached to the feed.
 */
typedef struct dxfg_indexed_tx_model_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_indexed_tx_model_t;

/**
 * @brief A handle of a Java IndexedTxModel.Builder object; release it with dxfg_JavaObjectHandler_release().
 *
 * The builder of dxfg_indexed_tx_model_t, see dxfg_IndexedTxModel_newBuilder().
 */
typedef struct dxfg_indexed_tx_model_builder_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_indexed_tx_model_builder_t;

/**
 * @brief A handle of a Java TimeSeriesTxModel object; release it with dxfg_JavaObjectHandler_release().
 *
 * The model is an incremental model of time-series events: it delivers the complete transactions of the events of
 * its symbol from its time to its listener. Close the model with dxfg_TimeSeriesTxModel_close() before releasing the
 * handle: an attached model stays attached to the feed.
 */
typedef struct dxfg_time_series_tx_model_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_time_series_tx_model_t;

/**
 * @brief A handle of a Java TimeSeriesTxModel.Builder object; release it with dxfg_JavaObjectHandler_release().
 *
 * The builder of dxfg_time_series_tx_model_t, see dxfg_TimeSeriesTxModel_newBuilder().
 */
typedef struct dxfg_time_series_tx_model_builder_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_time_series_tx_model_builder_t;

/**
 * @brief A handle of a Java TxModelListener object; release it with dxfg_JavaObjectHandler_release().
 *
 * The listener of the transactional models, create it with dxfg_TxModelListener_new().
 */
typedef struct dxfg_tx_model_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_tx_model_listener_t;

/**
 * @brief The function of a dxfg_tx_model_listener_t: called when one or more complete transactions are received.
 *
 * The function is called by the executor of the model (set with dxfg_IndexedTxModel_Builder_withExecutor() or
 * dxfg_TimeSeriesTxModel_Builder_withExecutor(); by default the executor of the feed), not by the thread that built
 * the model. The calls of one model are never concurrent. The events of one call have the same source and symbol.
 *
 * The SDK creates the source and the list of events for the call and frees them right after the function returns:
 * the callee must not release them and must copy whatever it needs to keep.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The source of the events. It is valid during the call only.
 * @param[in] events The events of one or more transactions. They belong to the SDK and are valid during the call
 * only. The list can be empty if an empty snapshot was received and the snapshot processing is on.
 * @param[in] isSnapshot 1 if the events are a snapshot (the state kept for the source has to be cleared), 0 otherwise.
 * @param[in] user_data The pointer that was passed to dxfg_TxModelListener_new(); the SDK does not use it.
 */
typedef void (*dxfg_TxModelListener_function_eventsReceived)(graal_isolatethread_t* thread, dxfg_indexed_event_source_t* source, dxfg_event_type_list* events, int32_t isSnapshot, void* user_data);

/**
 * @brief Creates a new builder of a model of the indexed events of the type.
 *
 * Java: `IndexedTxModel.newBuilder(Class)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_newBuilder() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of the indexed events of the model.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_newBuilder() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_newBuilder(graal_isolatethread_t* thread, dxfg_event_clazz_t eventType);

/**
 * @brief Returns the sources of the model.
 *
 * Java: `IndexedTxModel.getSources()`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_getSources() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return The sources (release them with dxfg_CList_IndexedEventSource_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_getSources() (experimental)")
dxfg_indexed_event_source_list* dxfg_IndexedTxModel_getSources(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source);

/**
 * @brief Sets the sources of the model: it subscribes to the new ones and unsubscribes from the old ones.
 *
 * An empty list means all the sources.
 *
 * Java: `IndexedTxModel.setSources(Set)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_setSources() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] sources The new sources of the model; an empty list means all the sources.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_setSources() (experimental)")
int32_t dxfg_IndexedTxModel_setSources(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source, dxfg_indexed_event_source_list* sources);

/**
 * @brief Returns whether the listener gets all the pending transactions in one call.
 *
 * Java: `AbstractTxModel.isBatchProcessing()`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_isBatchProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return 1 if enabled, 0 otherwise, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to
 * determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_isBatchProcessing() (experimental)")
int32_t dxfg_IndexedTxModel_isBatchProcessing(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source);

/**
 * @brief Returns whether a snapshot is delivered as its resulting state (without the removed events).
 *
 * Java: `AbstractTxModel.isSnapshotProcessing()`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_isSnapshotProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return 1 if enabled, 0 otherwise, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to
 * determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_isSnapshotProcessing() (experimental)")
int32_t dxfg_IndexedTxModel_isSnapshotProcessing(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source);

/**
 * @brief Attaches the model to the feed.
 *
 * Java: `AbstractTxModel.attach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_attach() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_attach() (experimental)")
int32_t dxfg_IndexedTxModel_attach(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source, dxfg_feed_t* feed);

/**
 * @brief Detaches the model from the feed.
 *
 * Java: `AbstractTxModel.detach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_detach() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_detach() (experimental)")
int32_t dxfg_IndexedTxModel_detach(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source, dxfg_feed_t* feed);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.setAggregationPeriod(TimePeriod)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_setAggregationPeriod() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] aggregationPeriod The new aggregation period of the events.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_setAggregationPeriod() (experimental)")
int32_t dxfg_IndexedTxModel_setAggregationPeriod(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source, dxfg_time_period_t* aggregationPeriod);

/**
 * @brief Closes the model: it is detached from all the feeds for good. Closing it again does nothing.
 *
 * Java: `AbstractTxModel.close()`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_close() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_close() (experimental)")
int32_t dxfg_IndexedTxModel_close(graal_isolatethread_t* thread, dxfg_indexed_tx_model_t* source);

/**
 * @brief Sets the sources (an empty list means all the sources).
 *
 * Java: `IndexedTxModel.Builder.withSources(Collection)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withSources() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] sources The sources of the events; an empty list means all the sources.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withSources() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withSources(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, dxfg_indexed_event_source_list* sources);

/**
 * @brief Builds the model (the symbol must be set).
 *
 * Java: `IndexedTxModel.Builder.build()`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_build() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @return The model (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_build() (experimental)")
dxfg_indexed_tx_model_t* dxfg_IndexedTxModel_Builder_build(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source);

/**
 * @brief Sets whether the listener gets all the pending transactions in one call (1, default) or one by one (0).
 *
 * Java: `AbstractTxModel.Builder.withBatchProcessing(boolean)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withBatchProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isBatchProcessing 1 to deliver all the pending transactions in one call, 0 to deliver them one by one.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withBatchProcessing() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withBatchProcessing(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, int32_t isBatchProcessing);

/**
 * @brief Sets whether a snapshot is delivered as its resulting state (1) or as it comes (0, default).
 *
 * Java: `AbstractTxModel.Builder.withSnapshotProcessing(boolean)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withSnapshotProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isSnapshotProcessing 1 to deliver a snapshot as its resulting state (without the removed events), 0 to
 * deliver
 * it as it comes.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withSnapshotProcessing() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withSnapshotProcessing(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, int32_t isSnapshotProcessing);

/**
 * @brief Sets the feed the model is attached to.
 *
 * Java: `AbstractTxModel.Builder.withFeed(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withFeed() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] feed The feed the model is attached to.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withFeed() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withFeed(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, dxfg_feed_t* feed);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.Builder.withAggregationPeriod(TimePeriod)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withAggregationPeriod() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] aggregationPeriod The aggregation period of the events.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withAggregationPeriod() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withAggregationPeriod(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, dxfg_time_period_t* aggregationPeriod);

/**
 * @brief Sets the symbol of the subscription.
 *
 * Java: `AbstractTxModel.Builder.withSymbol(Object)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withSymbol() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] symbol The symbol to subscribe for (a dxfg_symbol_t).
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withSymbol() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withSymbol(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, dxfg_symbol_t* symbol);

/**
 * @brief Sets the listener of the transactions.
 *
 * Java: `AbstractTxModel.Builder.withListener(TxModelListener)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withListener() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] listener The listener of the transactions (see dxfg_TxModelListener_new()).
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withListener() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withListener(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, dxfg_tx_model_listener_t* listener);

/**
 * @brief Sets the executor that calls the listener (by default, the executor of the feed).
 *
 * Java: `AbstractTxModel.Builder.withExecutor(Executor)`.
 *
 * @deprecated Use dxfg_v4_IndexedTxModel_Builder_withExecutor() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] executor The executor that calls the listener.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedTxModel_Builder_withExecutor() (experimental)")
dxfg_indexed_tx_model_builder_t* dxfg_IndexedTxModel_Builder_withExecutor(graal_isolatethread_t* thread, dxfg_indexed_tx_model_builder_t* source, dxfg_executor_t* executor);

/**
 * @brief Creates a new builder of a model of the time-series events of the type.
 *
 * Java: `TimeSeriesTxModel.newBuilder(Class)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_newBuilder() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The type of the time-series events of the model.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_newBuilder() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_newBuilder(graal_isolatethread_t* thread, dxfg_event_clazz_t eventType);

/**
 * @brief Returns the time the model is subscribed from.
 *
 * Java: `TimeSeriesTxModel.getFromTime()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_getFromTime() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return The time in milliseconds since the epoch (INT64_MAX if the model is not subscribed), or -1 on
 * error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_getFromTime() (experimental)")
int64_t dxfg_TimeSeriesTxModel_getFromTime(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source);

/**
 * @brief Sets the time the model is subscribed from (INT64_MAX unsubscribes); the model gets a new snapshot.
 *
 * Java: `TimeSeriesTxModel.setFromTime(long)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_setFromTime() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] fromTime The time to subscribe from, in milliseconds since the epoch; INT64_MAX unsubscribes.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_setFromTime() (experimental)")
int32_t dxfg_TimeSeriesTxModel_setFromTime(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source, int64_t fromTime);

/**
 * @brief Returns whether the listener gets all the pending transactions in one call.
 *
 * Java: `AbstractTxModel.isBatchProcessing()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_isBatchProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return 1 if enabled, 0 otherwise, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to
 * determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_isBatchProcessing() (experimental)")
int32_t dxfg_TimeSeriesTxModel_isBatchProcessing(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source);

/**
 * @brief Returns whether a snapshot is delivered as its resulting state (without the removed events).
 *
 * Java: `AbstractTxModel.isSnapshotProcessing()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_isSnapshotProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return 1 if enabled, 0 otherwise, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to
 * determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_isSnapshotProcessing() (experimental)")
int32_t dxfg_TimeSeriesTxModel_isSnapshotProcessing(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source);

/**
 * @brief Attaches the model to the feed.
 *
 * Java: `AbstractTxModel.attach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_attach() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to attach to.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_attach() (experimental)")
int32_t dxfg_TimeSeriesTxModel_attach(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source, dxfg_feed_t* feed);

/**
 * @brief Detaches the model from the feed.
 *
 * Java: `AbstractTxModel.detach(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_detach() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] feed The feed to detach from.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_detach() (experimental)")
int32_t dxfg_TimeSeriesTxModel_detach(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source, dxfg_feed_t* feed);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.setAggregationPeriod(TimePeriod)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_setAggregationPeriod() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @param[in] aggregationPeriod The new aggregation period of the events.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_setAggregationPeriod() (experimental)")
int32_t dxfg_TimeSeriesTxModel_setAggregationPeriod(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source, dxfg_time_period_t* aggregationPeriod);

/**
 * @brief Closes the model: it is detached from all the feeds for good. Closing it again does nothing.
 *
 * Java: `AbstractTxModel.close()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_close() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The model handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_close() (experimental)")
int32_t dxfg_TimeSeriesTxModel_close(graal_isolatethread_t* thread, dxfg_time_series_tx_model_t* source);

/**
 * @brief Sets the time to subscribe from, in milliseconds since the epoch (by default INT64_MAX: not subscribed).
 *
 * Java: `TimeSeriesTxModel.Builder.withFromTime(long)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withFromTime() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] fromTime The time to subscribe from, in milliseconds since the epoch (INT64_MAX: not subscribed).
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withFromTime() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withFromTime(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, int64_t fromTime);

/**
 * @brief Sets the order of the events of a snapshot.
 *
 * Java: `TimeSeriesTxModel.Builder.withSorting(TimeSeriesTxModel.SortOrder)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withSorting() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] sortOrder The order of the events of a snapshot.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withSorting() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withSorting(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, dxfg_TimeSeriesTxModelSortOrder_t sortOrder);

/**
 * @brief Sets whether the removed events are dropped (1) or delivered (0, default).
 *
 * Java: `TimeSeriesTxModel.Builder.ignoreRemoveEvents(boolean)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_ignoreRemoveEvents() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] ignoreRemoveEvents 1 to drop the removed events, 0 to deliver them.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_ignoreRemoveEvents() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_ignoreRemoveEvents(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, int32_t ignoreRemoveEvents);

/**
 * @brief Sets whether the events with an index not greater than the last processed one are dropped (1)
 * or delivered (0, default).
 *
 * Java: `TimeSeriesTxModel.Builder.ignoreEventsFromPast(boolean)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_ignoreEventsFromPast() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] ignoreEventFromPast 1 to drop the events with an index not greater than the last processed one, 0 to
 * deliver
 * them.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_ignoreEventsFromPast() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_ignoreEventsFromPast(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, int32_t ignoreEventFromPast);

/**
 * @brief Sets whether a new bar comes after the previous (closing) one, to confirm it (1), or alone (0, default).
 *
 * With 1, the snapshot processing and the ascending sorting are on.
 *
 * Java: `TimeSeriesTxModel.Builder.withConfirmationTick(boolean)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withConfirmationTick() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] confirmationTick 1 to deliver a new bar after the previous (closing) one, to confirm it, 0 to deliver it
 * alone.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withConfirmationTick() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withConfirmationTick(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, int32_t confirmationTick);

/**
 * @brief Accepted for compatibility: the pointer is ignored, the model gets no event peek (a Java `Consumer`).
 *
 * Java: `TimeSeriesTxModel.Builder.withEventPeek(Consumer)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withEventPeek() (experimental), or
 * dxfg_v4_TimeSeriesTxModel_Builder_withEventPeekFunction() (experimental) for an event peek that works.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] eventPeek Ignored: pass NULL.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withEventPeek() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withEventPeek(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, void* eventPeek);

/**
 * @brief Builds the model (the symbol must be set).
 *
 * Java: `TimeSeriesTxModel.Builder.build()`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_build() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @return The model (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_build() (experimental)")
dxfg_time_series_tx_model_t* dxfg_TimeSeriesTxModel_Builder_build(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source);

/**
 * @brief Sets whether the listener gets all the pending transactions in one call (1, default) or one by one (0).
 *
 * Java: `AbstractTxModel.Builder.withBatchProcessing(boolean)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withBatchProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isBatchProcessing 1 to deliver all the pending transactions in one call, 0 to deliver them one by one.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withBatchProcessing() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withBatchProcessing(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, int32_t isBatchProcessing);

/**
 * @brief Sets whether a snapshot is delivered as its resulting state (1) or as it comes (0, default).
 *
 * Java: `AbstractTxModel.Builder.withSnapshotProcessing(boolean)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withSnapshotProcessing() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] isSnapshotProcessing 1 to deliver a snapshot as its resulting state (without the removed events), 0 to
 * deliver
 * it as it comes.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withSnapshotProcessing() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withSnapshotProcessing(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, int32_t isSnapshotProcessing);

/**
 * @brief Sets the feed the model is attached to.
 *
 * Java: `AbstractTxModel.Builder.withFeed(DXFeed)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withFeed() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] feed The feed the model is attached to.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withFeed() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withFeed(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, dxfg_feed_t* feed);

/**
 * @brief Sets the aggregation period of the events.
 *
 * Java: `AbstractTxModel.Builder.withAggregationPeriod(TimePeriod)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withAggregationPeriod() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] aggregationPeriod The aggregation period of the events.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withAggregationPeriod() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withAggregationPeriod(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, dxfg_time_period_t* aggregationPeriod);

/**
 * @brief Sets the symbol of the subscription.
 *
 * Java: `AbstractTxModel.Builder.withSymbol(Object)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withSymbol() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] symbol The symbol to subscribe for (a dxfg_symbol_t).
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withSymbol() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withSymbol(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, dxfg_symbol_t* symbol);

/**
 * @brief Sets the listener of the transactions.
 *
 * Java: `AbstractTxModel.Builder.withListener(TxModelListener)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withListener() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] listener The listener of the transactions (see dxfg_TxModelListener_new()).
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withListener() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withListener(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, dxfg_tx_model_listener_t* listener);

/**
 * @brief Sets the executor that calls the listener (by default, the executor of the feed).
 *
 * Java: `AbstractTxModel.Builder.withExecutor(Executor)`.
 *
 * @deprecated Use dxfg_v4_TimeSeriesTxModel_Builder_withExecutor() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The builder handle.
 * @param[in] executor The executor that calls the listener.
 * @return The builder (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeSeriesTxModel_Builder_withExecutor() (experimental)")
dxfg_time_series_tx_model_builder_t* dxfg_TimeSeriesTxModel_Builder_withExecutor(graal_isolatethread_t* thread, dxfg_time_series_tx_model_builder_t* source, dxfg_executor_t* executor);

/**
 * @brief Creates a listener that calls the function with the user data.
 *
 * After the last handle of the listener is released (dxfg_JavaObjectHandler_release()), the function is not called
 * anymore and none of its calls is in progress, so the user data can be freed.
 *
 * Java: `TxModelListener.eventsReceived(IndexedEventSource, List, boolean)`.
 *
 * @deprecated Use dxfg_v4_TxModelListener_new() (experimental).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] function_eventsReceived The function to call when transactions are received; see
 * #dxfg_TxModelListener_function_eventsReceived for the thread and the ownership of its data.
 * @param[in] user_data The pointer that is passed to the function as is; the SDK does not use it. It has to stay valid
 * until the listener is released.
 * @return The listener (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_TxModelListener_new() (experimental)")
dxfg_tx_model_listener_t* dxfg_TxModelListener_new(graal_isolatethread_t* thread, dxfg_TxModelListener_function_eventsReceived  function_eventsReceived, void* user_data);

/** @} */

/** @} */ // end of EventModel

#ifdef __cplusplus
}
#endif

#endif // DXFG_EVENT_MODEL_H
