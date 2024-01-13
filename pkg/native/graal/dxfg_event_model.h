// SPDX-License-Identifier: MPL-2.0

#ifndef DXFEED_GRAAL_NATIVE_SDK_EVENT_MODEL_H_
#define DXFEED_GRAAL_NATIVE_SDK_EVENT_MODEL_H_

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "graal_isolate.h"
#include "dxfg_javac.h"
#include "dxfg_feed.h"

/** @defgroup EventModel
 *  @{
 */

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/market/OrderBookModel.html">Javadoc</a>
 */
typedef struct dxfg_order_book_model_t {
    dxfg_java_object_handler handler;
} dxfg_order_book_model_t;

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/IndexedEventModel.html">Javadoc</a>
 */
typedef struct dxfg_indexed_event_model_t {
    dxfg_java_object_handler handler;
} dxfg_indexed_event_model_t;

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/TimeSeriesEventModel.html">Javadoc</a>
 */
typedef struct dxfg_time_series_event_model_t {
    dxfg_java_object_handler handler;
} dxfg_time_series_event_model_t;

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/ObservableListModel.html">Javadoc</a>
 */
typedef struct dxfg_observable_list_model_t {
    dxfg_java_object_handler handler;
} dxfg_observable_list_model_t;

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/ObservableListModelListener.html">Javadoc</a>
 */
typedef struct dxfg_observable_list_model_listener_t {
    dxfg_java_object_handler handler;
} dxfg_observable_list_model_listener_t;

typedef void (*dxfg_observable_list_model_listener_function)(graal_isolatethread_t *thread, dxfg_event_type_list* orders, void *user_data);

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/market/OrderBookModelListener.html">Javadoc</a>
 */
typedef struct dxfg_order_book_model_listener_t {
    dxfg_java_object_handler handler;
} dxfg_order_book_model_listener_t;

typedef void (*dxfg_order_book_model_listener_function)(graal_isolatethread_t *thread, dxfg_order_book_model_t* model, void *user_data);

dxfg_order_book_model_t*          dxfg_OrderBookModel_new(graal_isolatethread_t *thread);
dxfg_indexed_event_model_t*       dxfg_IndexedEventModel_new(graal_isolatethread_t *thread, dxfg_event_clazz_t clazz);
dxfg_time_series_event_model_t*   dxfg_TimeSeriesEventModel_new(graal_isolatethread_t *thread, dxfg_event_clazz_t clazz);

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/market/OrderBookModelFilter.html">Javadoc</a>
 */
typedef enum dxfg_order_book_model_filter_t {
    COMPOSITE,
    REGIONAL,
    AGGREGATE,
    ORDER,
    COMPOSITE_REGIONAL,
    COMPOSITE_REGIONAL_AGGREGATE,
    ALL,
} dxfg_order_book_model_filter_t;

int32_t                                       dxfg_OrderBookModel_attach(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_feed_t *feed);
int32_t                                       dxfg_OrderBookModel_detach(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_feed_t *feed);
int32_t                                       dxfg_OrderBookModel_close(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);
dxfg_executor_t*                              dxfg_OrderBookModel_getExecutor(graal_isolatethread_t *thread, dxfg_order_book_model_t *model); // use dxfg_JavaObjectHandler_release to release the allocated memory
int32_t                                       dxfg_OrderBookModel_setExecutor(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_executor_t *executor);
int32_t                                       dxfg_OrderBookModel_clear(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);
int32_t                                       dxfg_IndexedEventModel_attach(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_feed_t *feed);
int32_t                                       dxfg_IndexedEventModel_detach(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_feed_t *feed);
int32_t                                       dxfg_IndexedEventModel_close(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);
dxfg_executor_t*                              dxfg_IndexedEventModel_getExecutor(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model); // use dxfg_JavaObjectHandler_release to release the allocated memory
int32_t                                       dxfg_IndexedEventModel_setExecutor(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_executor_t *executor);
int32_t                                       dxfg_IndexedEventModel_clear(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);
int32_t                                       dxfg_TimeSeriesEventModel_attach(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_feed_t *feed);
int32_t                                       dxfg_TimeSeriesEventModel_detach(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_feed_t *feed);
int32_t                                       dxfg_TimeSeriesEventModel_close(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);
dxfg_executor_t*                              dxfg_TimeSeriesEventModel_getExecutor(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model); // use dxfg_JavaObjectHandler_release to release the allocated memory
int32_t                                       dxfg_TimeSeriesEventModel_setExecutor(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_executor_t *executor);
int32_t                                       dxfg_TimeSeriesEventModel_clear(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);
int32_t                                       dxfg_OrderBookModel_getFilter(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);
int32_t                                       dxfg_OrderBookModel_setFilter(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_order_book_model_filter_t filter);
const char*                                   dxfg_OrderBookModel_getSymbol(graal_isolatethread_t *thread, dxfg_order_book_model_t *model); // use dxfg_String_release to release the allocated memory
int32_t                                       dxfg_OrderBookModel_setSymbol(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, const char* symbol);
dxfg_symbol_t*                                dxfg_IndexedEventModel_getSymbol(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model); // use dxfg_Symbol_release to release the allocated memory
int32_t                                       dxfg_IndexedEventModel_setSymbol(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, dxfg_symbol_t* symbol);
dxfg_symbol_t*                                dxfg_TimeSeriesEventModel_getSymbol(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model); // use dxfg_Symbol_release to release the allocated memory
int32_t                                       dxfg_TimeSeriesEventModel_setSymbol(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, dxfg_symbol_t* symbol);
int32_t                                       dxfg_OrderBookModel_getLotSize(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);
int32_t                                       dxfg_OrderBookModel_setLotSize(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, int32_t lotSize);
int32_t                                       dxfg_IndexedEventModel_getSizeLimit(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);
int32_t                                       dxfg_IndexedEventModel_setSizeLimit(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model, int32_t sizeLimit);
int32_t                                       dxfg_TimeSeriesEventModel_getSizeLimit(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);
int32_t                                       dxfg_TimeSeriesEventModel_setSizeLimit(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, int32_t sizeLimit);
int64_t                                       dxfg_TimeSeriesEventModel_getFromTime(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);
int32_t                                       dxfg_TimeSeriesEventModel_setFromTime(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model, int64_t fromTime);
dxfg_observable_list_model_t*                 dxfg_OrderBookModel_getBuyOrders(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);
dxfg_observable_list_model_t*                 dxfg_OrderBookModel_getSellOrders(graal_isolatethread_t *thread, dxfg_order_book_model_t *model);
dxfg_observable_list_model_t*                 dxfg_IndexedEventModel_getEventsList(graal_isolatethread_t *thread, dxfg_indexed_event_model_t *model);
dxfg_observable_list_model_t*                 dxfg_TimeSeriesEventModel_getEventsList(graal_isolatethread_t *thread, dxfg_time_series_event_model_t *model);
dxfg_observable_list_model_listener_t*        dxfg_ObservableListModelListener_new(graal_isolatethread_t *thread, dxfg_observable_list_model_listener_function user_func, void *user_data);
int32_t                                       dxfg_ObservableListModel_addListener(graal_isolatethread_t *thread, dxfg_observable_list_model_t *list, dxfg_observable_list_model_listener_t *listener);
int32_t                                       dxfg_ObservableListModel_removeListener(graal_isolatethread_t *thread, dxfg_observable_list_model_t *list, dxfg_observable_list_model_listener_t *listener);
dxfg_order_book_model_listener_t*             dxfg_OrderBookModelListener_new(graal_isolatethread_t *thread, dxfg_order_book_model_listener_function user_func, void *user_data);
int32_t                                       dxfg_OrderBookModel_addListener(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_order_book_model_listener_t *listener);
int32_t                                       dxfg_OrderBookModel_removeListener(graal_isolatethread_t *thread, dxfg_order_book_model_t *model, dxfg_order_book_model_listener_t *listener);
dxfg_event_type_list*                         dxfg_ObservableListModel_toArray(graal_isolatethread_t *thread, dxfg_observable_list_model_t *list); // use dxfg_CList_EventType_release to release the allocated memory

/** @} */ // end of EventModel

#ifdef __cplusplus
}
#endif

#endif // DXFEED_GRAAL_NATIVE_SDK_EVENT_MODEL_H_