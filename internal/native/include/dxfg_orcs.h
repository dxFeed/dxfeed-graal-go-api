// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK ORCS API functions and types declarations
 */

#ifndef DXFG_ORCS_H
#define DXFG_ORCS_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_events.h"
#include "dxfg_javac.h"
#include "graal_isolate.h"

/**
 * @defgroup OrcsApi ORCS API
 * @brief The client of the ORCS (order and quote history) service: price levels and quotes of a candle symbol within a
 * time range, the order sources and symbols that are available to the client, and the checker of an order list.
 *
 * All the functions of this API return #DXFG_EXECUTE_SUCCESSFULLY (0) or #DXFG_EXECUTE_FAIL (-1) on error, and the
 * result is passed through an `[out]` parameter; it is not changed on error. A NULL `[out]` pointer is an error too.
 * dxfg_get_and_clear_thread_exception_t() gives the exception. A handle that a function returns is owned by the
 * caller: release it with dxfg_JavaObjectHandler_release() (see the group Javac for the ownership rules). A list of
 * events is released with dxfg_CList_EventType_release(), the other lists have their own `_free` functions.
 */

/**
 * @defgroup PriceLevelService PriceLevelService
 * @brief The handle of the client of the price level service.
 * @ingroup OrcsApi
 *
 * @see PriceLevelServiceFunctions
 * @{
 */

/**
 * @brief A handle of a client of the price level service; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_PriceLevelService_new(), which connects it to an RMI endpoint. Releasing the handle does not
 * close the connection: call dxfg_PriceLevelService_close().
 *
 * https://docs.dxfeed.com/orcs/com/dxfeed/orcs/api/PriceLevelService.html
 */
typedef struct dxfg_price_level_service_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_price_level_service_t;

/// @} end of PriceLevelService

/**
 * @defgroup AuthOrderSource AuthOrderSource
 * @brief The order sources and symbols that are available to the client, and the structures that describe them.
 * @ingroup OrcsApi
 *
 * @see AuthOrderSourceFunctions
 * @{
 */

/**
 * @brief A handle of Java AuthOrderSource: the information about the symbols that are available to the client for the
 * entitled collection of order sources; release it with dxfg_JavaObjectHandler_release().
 *
 * Get it with dxfg_PriceLevelService_getAuthOrderSource(); read it with dxfg_AuthOrderSource_getByIds() or
 * dxfg_AuthOrderSource_getByOrderSources().
 *
 * https://docs.dxfeed.com/orcs/com/dxfeed/orcs/api/AuthOrderSource.html
 */
typedef struct dxfg_auth_order_source_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_auth_order_source_t;

/**
 * @brief An entry of the symbols by the identifier of an order source (Java `Map.Entry<Integer, Set<String>>`).
 *
 * The structure belongs to the list dxfg_symbols_by_order_source_id_map_entry_list_t that contains it.
 */
typedef struct dxfg_symbols_by_order_source_id_map_entry_t {
    /** The identifier of the order source (the `id` of a dxfg_indexed_event_source_t of the type ORDER_SOURCE). */
    int32_t order_source_id;
    /** The symbols that are available to the client for the order source (UTF-8 null-terminated strings, in an
     * unspecified order); the list belongs to the entry and is freed with it. */
    dxfg_string_list *symbols;
} dxfg_symbols_by_order_source_id_map_entry_t;

/**
 * @brief A list of the entries of the symbols by the identifier of an order source (Java
 * `Collection<? extends Map.Entry<Integer, Set<String>>>`).
 *
 * Returned by dxfg_AuthOrderSource_getByIds(); release it with dxfg_symbols_by_order_source_id_map_entry_list_free().
 */
typedef struct dxfg_symbols_by_order_source_id_map_entry_list_t {
    /** The number of the entries; 0 for an empty list. */
    int32_t size;
    /** The array of `size` pointers to the entries; NULL for an empty list. */
    dxfg_symbols_by_order_source_id_map_entry_t **elements;
} dxfg_symbols_by_order_source_id_map_entry_list_t;

/**
 * @brief An entry of the symbols by an order source (Java `Map.Entry<? extends IndexedEventSource, Set<String>>`).
 *
 * The structure belongs to the list dxfg_symbols_by_order_source_map_entry_list_t that contains it.
 */
typedef struct dxfg_symbols_by_order_source_map_entry_t {
    /** The order source, a dxfg_indexed_event_source_t of the type ORDER_SOURCE; it belongs to the entry and is freed
     * with it. */
    dxfg_indexed_event_source_t *order_source;
    /** The symbols that are available to the client for the order source (UTF-8 null-terminated strings, in an
     * unspecified order); the list belongs to the entry and is freed with it. */
    dxfg_string_list *symbols;
} dxfg_symbols_by_order_source_map_entry_t;

/**
 * @brief A list of the entries of the symbols by an order source (Java
 * `Collection<? extends Map.Entry<? extends IndexedEventSource, Set<String>>>`).
 *
 * Returned by dxfg_AuthOrderSource_getByOrderSources(); release it with
 * dxfg_symbols_by_order_source_map_entry_list_free().
 */
typedef struct dxfg_symbols_by_order_source_map_entry_list_t {
    /** The number of the entries; 0 for an empty list. */
    int32_t size;
    /** The array of `size` pointers to the entries; NULL for an empty list. */
    dxfg_symbols_by_order_source_map_entry_t **elements;
} dxfg_symbols_by_order_source_map_entry_list_t;

/// @} end of AuthOrderSource

/**
 * @defgroup PriceLevelServiceFunctions PriceLevelService' functions
 * @brief The functions that create, query and close the client of the price level service.
 * @ingroup PriceLevelService
 * @{
 */

/**
 * @brief Creates a client of the price level service and connects it to an RMI endpoint.
 *
 * The client is an RMI client endpoint that is connected to the address; the service is its proxy. The handle of the
 * client is written to `service`. Close the connection with dxfg_PriceLevelService_close().
 *
 * Java: `new PriceLevelServiceHolder(String)` (the SDK class that wraps `RMIEndpoint` and `PriceLevelService`).
 *
 * @deprecated Use dxfg_v4_PriceLevelService_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] address The RMI endpoint's address: a UTF-8 null-terminated string; the caller keeps its ownership.
 * @param[out] service The new client (release the handle with dxfg_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelService_new()")
int32_t dxfg_PriceLevelService_new(graal_isolatethread_t *thread, const char *address,
                                   DXFG_OUT dxfg_price_level_service_t **service);

/**
 * @brief Returns the price levels (orders) of a candle symbol within a time range.
 *
 * The orders are sorted in ascending order by time.
 *
 * The candle symbol shall contain a supported CandlePeriod and a custom attribute called granularity with the key
 * `GRANULARITY_ATTRIBUTE_KEY` ("gr"). The granularity value shall be an integer in seconds; the minimal value is 1
 * second. If the CandlePeriod or the granularity is not supported by the service, the result is an empty list.
 *
 * Java: `PriceLevelService.getOrders(CandleSymbol, int, long, long, String)`, the order source is passed by its id.
 *
 * @deprecated Use dxfg_v4_PriceLevelService_getOrdersWithCaller().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] orderSource The order source to request: a dxfg_indexed_event_source_t of the type ORDER_SOURCE; the call
 * fails for another kind of source. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] caller The caller identifier: a UTF-8 null-terminated string that is passed to the service as is.
 * @param[out] orders The list of the order events (dxfg_event_type_list), possibly empty. Release it with
 * dxfg_CList_EventType_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelService_getOrdersWithCaller()")
int32_t dxfg_PriceLevelService_getOrders(graal_isolatethread_t *thread, dxfg_price_level_service_t *service,
                                         dxfg_symbol_t *candleSymbol, dxfg_indexed_event_source_t *orderSource,
                                         int64_t from, int64_t to, const char *caller,
                                         DXFG_OUT dxfg_event_type_list **orders);

/**
 * @brief Returns the price levels (orders) of a candle symbol within a time range for the default caller.
 *
 * It is the same as dxfg_PriceLevelService_getOrders() with the caller identifier "dxfg": see that function for the
 * requirements on the candle symbol and the granularity, and for the empty list.
 *
 * Java: `PriceLevelService.getOrders(CandleSymbol, int, long, long, String)` with the caller "dxfg" (the SDK method
 * `PriceLevelServiceHolder.getOrders(CandleSymbol, OrderSource, long, long)`).
 *
 * @deprecated Use dxfg_v4_PriceLevelService_getOrders().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] orderSource The order source to request: a dxfg_indexed_event_source_t of the type ORDER_SOURCE; the call
 * fails for another kind of source. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[out] orders The list of the order events (dxfg_event_type_list), possibly empty. Release it with
 * dxfg_CList_EventType_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelService_getOrders()")
int32_t dxfg_PriceLevelService_getOrders2(graal_isolatethread_t *thread, dxfg_price_level_service_t *service,
                                          dxfg_symbol_t *candleSymbol, dxfg_indexed_event_source_t *orderSource,
                                          int64_t from, int64_t to, DXFG_OUT dxfg_event_type_list **orders);

/**
 * @brief Returns the order sources and the symbols of each order source that are available to the client.
 *
 * The order sources and the symbols are filtered according to the permissions of the client. The view is built as of
 * now: the result contains only the existing data (for example, no symbols that were delisted). Read the result with
 * dxfg_AuthOrderSource_getByIds() or dxfg_AuthOrderSource_getByOrderSources().
 *
 * Java: `PriceLevelService.getAuthOrderSource()`.
 *
 * @deprecated Use dxfg_v4_PriceLevelService_getAuthOrderSource().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[out] authOrderSource The AuthOrderSource (release the handle with dxfg_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelService_getAuthOrderSource()")
int32_t dxfg_PriceLevelService_getAuthOrderSource(graal_isolatethread_t *thread, dxfg_price_level_service_t *service,
                                                  DXFG_OUT dxfg_auth_order_source_t **authOrderSource);

/**
 * @brief Returns the quotes of a candle symbol within a time range.
 *
 * The quotes are sorted in ascending order by time (see the Java documentation).
 *
 * Java: `PriceLevelService.getQuotes(CandleSymbol, long, long, String)`.
 *
 * @deprecated Use dxfg_v4_PriceLevelService_getQuotesWithCaller().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] caller The caller identifier: a UTF-8 null-terminated string that is passed to the service as is.
 * @param[out] quotes The list of the quote events (dxfg_event_type_list), possibly empty. Release it with
 * dxfg_CList_EventType_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelService_getQuotesWithCaller()")
int32_t dxfg_PriceLevelService_getQuotes(graal_isolatethread_t *thread, dxfg_price_level_service_t *service,
                                         dxfg_symbol_t *candleSymbol, int64_t from, int64_t to, const char *caller,
                                         DXFG_OUT dxfg_event_type_list **quotes);

/**
 * @brief Returns the quotes of a candle symbol within a time range for the default caller.
 *
 * It is the same as dxfg_PriceLevelService_getQuotes() with the caller identifier "dxfg".
 *
 * Java: `PriceLevelService.getQuotes(CandleSymbol, long, long, String)` with the caller "dxfg" (the SDK method
 * `PriceLevelServiceHolder.getQuotes(CandleSymbol, long, long)`).
 *
 * @deprecated Use dxfg_v4_PriceLevelService_getQuotes().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[out] quotes The list of the quote events (dxfg_event_type_list), possibly empty. Release it with
 * dxfg_CList_EventType_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelService_getQuotes()")
int32_t dxfg_PriceLevelService_getQuotes2(graal_isolatethread_t *thread, dxfg_price_level_service_t *service,
                                          dxfg_symbol_t *candleSymbol, int64_t from, int64_t to,
                                          DXFG_OUT dxfg_event_type_list **quotes);

/**
 * @brief Closes the client of the price level service and its RMI endpoint.
 *
 * Do not use the client after the call; release its handle with dxfg_JavaObjectHandler_release().
 *
 * Java: `PriceLevelServiceHolder.close()` (the SDK class that closes its `RMIEndpoint`).
 *
 * @deprecated Use dxfg_v4_PriceLevelService_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelService_close()")
int32_t dxfg_PriceLevelService_close(graal_isolatethread_t *thread, dxfg_price_level_service_t *service);

/// @} end of PriceLevelServiceFunctions

/**
 * @defgroup AuthOrderSourceFunctions AuthOrderSource' functions
 * @brief The functions that read an AuthOrderSource and free the lists that they return.
 * @ingroup AuthOrderSource
 * @{
 */

/**
 * @brief Returns the entries (the identifier of an order source, the symbols) of an AuthOrderSource.
 *
 * The list is a copy: it does not depend on the handle of the AuthOrderSource.
 *
 * Java: `AuthOrderSource.getByIds()`.
 *
 * @deprecated Use dxfg_v4_AuthOrderSource_getByIds().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authOrderSource The AuthOrderSource.
 * @param[out] symbolsByOrderSourceIdMapEntryList The list of the entries. Release it with
 * dxfg_symbols_by_order_source_id_map_entry_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthOrderSource_getByIds()")
int32_t dxfg_AuthOrderSource_getByIds(
    graal_isolatethread_t *thread, dxfg_auth_order_source_t *authOrderSource,
    DXFG_OUT dxfg_symbols_by_order_source_id_map_entry_list_t **symbolsByOrderSourceIdMapEntryList);

/**
 * @brief Frees a list of the entries (the identifier of an order source, the symbols).
 *
 * Frees the list, its entries and their symbol lists. A NULL list is ignored.
 *
 * Java: not applicable (the list is allocated by the SDK outside the Java heap).
 *
 * @deprecated Use dxfg_v4_symbols_by_order_source_id_map_entry_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list of the entries that dxfg_AuthOrderSource_getByIds() returned.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_symbols_by_order_source_id_map_entry_list_release()")
int32_t dxfg_symbols_by_order_source_id_map_entry_list_free(graal_isolatethread_t *thread,
                                                            dxfg_symbols_by_order_source_id_map_entry_list_t *list);

/**
 * @brief Returns the entries (an order source, the symbols) of an AuthOrderSource.
 *
 * The list is a copy: it does not depend on the handle of the AuthOrderSource.
 *
 * Java: `AuthOrderSource.getByOrderSources()`.
 *
 * @deprecated Use dxfg_v4_AuthOrderSource_getByOrderSources().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authOrderSource The AuthOrderSource.
 * @param[out] symbolsByOrderSourceMapEntryList The list of the entries. Release it with
 * dxfg_symbols_by_order_source_map_entry_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthOrderSource_getByOrderSources()")
int32_t dxfg_AuthOrderSource_getByOrderSources(
    graal_isolatethread_t *thread, dxfg_auth_order_source_t *authOrderSource,
    DXFG_OUT dxfg_symbols_by_order_source_map_entry_list_t **symbolsByOrderSourceMapEntryList);

/**
 * @brief Frees a list of the entries (an order source, the symbols).
 *
 * Frees the list, its entries, their order sources and their symbol lists. A NULL list is ignored.
 *
 * Java: not applicable (the list is allocated by the SDK outside the Java heap).
 *
 * @deprecated Use dxfg_v4_symbols_by_order_source_map_entry_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list of the entries that dxfg_AuthOrderSource_getByOrderSources() returned.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_symbols_by_order_source_map_entry_list_release()")
int32_t dxfg_symbols_by_order_source_map_entry_list_free(graal_isolatethread_t *thread,
                                                         dxfg_symbols_by_order_source_map_entry_list_t *list);

/// @} end of AuthOrderSourceFunctions

/**
 * @defgroup PriceLevelCheckerFunctions PriceLevelChecker' functions
 * @brief The checker of the consistency of a list of orders.
 * @ingroup OrcsApi
 * @{
 */

/**
 * @brief Checks the consistency of a list of orders.
 *
 * PriceLevelChecker is a utility to check that in each moment of time the bid price of the orders is lower than the
 * ask price. The check goes through the orders in the order of the list, so pass a list that is sorted by time, as
 * dxfg_PriceLevelService_getOrders() returns it. Besides the result, the check reports to the Java log (not to the
 * caller) the additional information:
 * - a gap between two subsequent orders that is greater than `timeGapBound`;
 * - the last bid and ask, changed with the minimal period of 1 second, when `printQuotes` is enabled;
 * - a spike in the quotes, when `printQuotes` is enabled.
 *
 * The orders are not valid if the check finds a crossing of the bid and ask prices, or if the greatest gap between two
 * subsequent orders is greater than `timeGapBound`. A list without orders is valid. Note that in some cases it is
 * normal when the bid is greater than or equal to the ask.
 *
 * Java: `PriceLevelChecker.validate(List, long, boolean)`.
 *
 * @deprecated Use dxfg_v4_PriceLevelChecker_validate().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] orders The list of the orders (events that are orders, for example the list that
 * dxfg_PriceLevelService_getOrders() returned); the call fails for another kind of event. The caller keeps its
 * ownership.
 * @param[in] timeGapBound The greatest allowed gap between the times of two subsequent orders, in milliseconds (the
 * unit of the time of an order).
 * @param[in] printQuotes Non-zero to log the quotes during the check, 0 not to.
 * @param[out] isValid Set to 1 if the orders are valid, 0 otherwise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_PriceLevelChecker_validate()")
int32_t dxfg_PriceLevelChecker_validate(graal_isolatethread_t *thread, dxfg_event_type_list *orders,
                                        int64_t timeGapBound, int32_t /*bool*/ printQuotes, DXFG_OUT int32_t *isValid);

/// @} end of PriceLevelCheckerFunctions

#ifdef __cplusplus
}
#endif

#endif // DXFG_ORCS_H
