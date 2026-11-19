// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_orcs.h and the v4 names of its types (see the group @ref V4_OrcsApi).
 */

#ifndef DXFG_V4_ORCS_H
#define DXFG_V4_ORCS_H

#include "dxfg_orcs.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_events.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_OrcsApi ORCS API (v4)
 *  @brief The v4 functions and type names of the group @ref OrcsApi. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_price_level_service_t */
typedef dxfg_price_level_service_t dxfg_v4_price_level_service_t;

/** @copydoc dxfg_auth_order_source_t */
typedef dxfg_auth_order_source_t dxfg_v4_auth_order_source_t;

/** @copydoc dxfg_symbols_by_order_source_id_map_entry_t */
typedef dxfg_symbols_by_order_source_id_map_entry_t dxfg_v4_symbols_by_order_source_id_map_entry_t;

/** @copydoc dxfg_symbols_by_order_source_id_map_entry_list_t */
typedef dxfg_symbols_by_order_source_id_map_entry_list_t dxfg_v4_symbols_by_order_source_id_map_entry_list;

/** @copydoc dxfg_symbols_by_order_source_map_entry_t */
typedef dxfg_symbols_by_order_source_map_entry_t dxfg_v4_symbols_by_order_source_map_entry_t;

/** @copydoc dxfg_symbols_by_order_source_map_entry_list_t */
typedef dxfg_symbols_by_order_source_map_entry_list_t dxfg_v4_symbols_by_order_source_map_entry_list;

/**
 * @brief Creates a client of the price level service and connects it to an RMI endpoint.
 *
 * The client is an RMI client endpoint that is connected to the address; the service is its proxy. The handle of the
 * client is written to `service`. Close the connection with dxfg_v4_PriceLevelService_close().
 *
 * Java: `new PriceLevelServiceHolder(String)` (the SDK class that wraps `RMIEndpoint` and `PriceLevelService`).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] address The RMI endpoint's address: a UTF-8 null-terminated string; the caller keeps its ownership.
 * @param[out] service The new client (release the handle with dxfg_v4_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelService_new(graal_isolatethread_t *thread, const char *address,
                                      DXFG_OUT dxfg_v4_price_level_service_t **service);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_v4_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] orderSource The order source to request: a dxfg_v4_indexed_event_source_t of the type ORDER_SOURCE; the
 * call fails for another kind of source. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] caller The caller identifier: a UTF-8 null-terminated string that is passed to the service as is.
 * @param[out] orders The list of the order events (dxfg_v4_event_type_list), possibly empty. Release it with
 * dxfg_v4_event_type_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelService_getOrdersWithCaller(graal_isolatethread_t *thread,
                                                      dxfg_v4_price_level_service_t *service,
                                                      dxfg_v4_symbol_t *candleSymbol,
                                                      dxfg_v4_indexed_event_source_t *orderSource, int64_t from,
                                                      int64_t to, const char *caller,
                                                      DXFG_OUT dxfg_v4_event_type_list **orders);

/**
 * @brief Returns the price levels (orders) of a candle symbol within a time range for the default caller.
 *
 * It is the same as dxfg_v4_PriceLevelService_getOrdersWithCaller() with the caller identifier "dxfg": see that
 * function for the requirements on the candle symbol and the granularity, and for the empty list.
 *
 * Java: `PriceLevelService.getOrders(CandleSymbol, int, long, long, String)` with the caller "dxfg" (the SDK method
 * `PriceLevelServiceHolder.getOrders(CandleSymbol, OrderSource, long, long)`).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_v4_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] orderSource The order source to request: a dxfg_v4_indexed_event_source_t of the type ORDER_SOURCE; the
 * call fails for another kind of source. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[out] orders The list of the order events (dxfg_v4_event_type_list), possibly empty. Release it with
 * dxfg_v4_event_type_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelService_getOrders(graal_isolatethread_t *thread, dxfg_v4_price_level_service_t *service,
                                            dxfg_v4_symbol_t *candleSymbol, dxfg_v4_indexed_event_source_t *orderSource,
                                            int64_t from, int64_t to, DXFG_OUT dxfg_v4_event_type_list **orders);

/**
 * @brief Returns the order sources and the symbols of each order source that are available to the client.
 *
 * The order sources and the symbols are filtered according to the permissions of the client. The view is built as of
 * now: the result contains only the existing data (for example, no symbols that were delisted). Read the result with
 * dxfg_v4_AuthOrderSource_getByIds() or dxfg_v4_AuthOrderSource_getByOrderSources().
 *
 * Java: `PriceLevelService.getAuthOrderSource()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[out] authOrderSource The AuthOrderSource (release the handle with dxfg_v4_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelService_getAuthOrderSource(graal_isolatethread_t *thread,
                                                     dxfg_v4_price_level_service_t *service,
                                                     DXFG_OUT dxfg_v4_auth_order_source_t **authOrderSource);

/**
 * @brief Returns the quotes of a candle symbol within a time range.
 *
 * The quotes are sorted in ascending order by time (see the Java documentation).
 *
 * Java: `PriceLevelService.getQuotes(CandleSymbol, long, long, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_v4_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] caller The caller identifier: a UTF-8 null-terminated string that is passed to the service as is.
 * @param[out] quotes The list of the quote events (dxfg_v4_event_type_list), possibly empty. Release it with
 * dxfg_v4_event_type_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelService_getQuotesWithCaller(graal_isolatethread_t *thread,
                                                      dxfg_v4_price_level_service_t *service,
                                                      dxfg_v4_symbol_t *candleSymbol, int64_t from, int64_t to,
                                                      const char *caller, DXFG_OUT dxfg_v4_event_type_list **quotes);

/**
 * @brief Returns the quotes of a candle symbol within a time range for the default caller.
 *
 * It is the same as dxfg_v4_PriceLevelService_getQuotesWithCaller() with the caller identifier "dxfg".
 *
 * Java: `PriceLevelService.getQuotes(CandleSymbol, long, long, String)` with the caller "dxfg" (the SDK method
 * `PriceLevelServiceHolder.getQuotes(CandleSymbol, long, long)`).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @param[in] candleSymbol The candle symbol to request: a dxfg_v4_candle_symbol_t (type #CANDLE); the call fails for
 * another kind of symbol. The caller keeps its ownership.
 * @param[in] from The "from" time in UTC (the unit is not stated by the Java documentation).
 * @param[in] to The "to" time in UTC (the unit is not stated by the Java documentation).
 * @param[out] quotes The list of the quote events (dxfg_v4_event_type_list), possibly empty. Release it with
 * dxfg_v4_event_type_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelService_getQuotes(graal_isolatethread_t *thread, dxfg_v4_price_level_service_t *service,
                                            dxfg_v4_symbol_t *candleSymbol, int64_t from, int64_t to,
                                            DXFG_OUT dxfg_v4_event_type_list **quotes);

/**
 * @brief Closes the client of the price level service and its RMI endpoint.
 *
 * Do not use the client after the call; release its handle with dxfg_v4_JavaObjectHandler_release().
 *
 * Java: `PriceLevelServiceHolder.close()` (the SDK class that closes its `RMIEndpoint`).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The PriceLevelService.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelService_close(graal_isolatethread_t *thread, dxfg_v4_price_level_service_t *service);

/**
 * @brief Returns the entries (the identifier of an order source, the symbols) of an AuthOrderSource.
 *
 * The list is a copy: it does not depend on the handle of the AuthOrderSource.
 *
 * Java: `AuthOrderSource.getByIds()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authOrderSource The AuthOrderSource.
 * @param[out] symbolsByOrderSourceIdMapEntryList The list of the entries. Release it with
 * dxfg_v4_symbols_by_order_source_id_map_entry_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthOrderSource_getByIds(
    graal_isolatethread_t *thread, dxfg_v4_auth_order_source_t *authOrderSource,
    DXFG_OUT dxfg_v4_symbols_by_order_source_id_map_entry_list **symbolsByOrderSourceIdMapEntryList);

/**
 * @brief Frees a list of the entries (the identifier of an order source, the symbols).
 *
 * Frees the list, its entries and their symbol lists. A NULL list is ignored.
 *
 * Java: not applicable (the list is allocated by the SDK outside the Java heap).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list of the entries that dxfg_v4_AuthOrderSource_getByIds() returned.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_symbols_by_order_source_id_map_entry_list_release(
    graal_isolatethread_t *thread, dxfg_v4_symbols_by_order_source_id_map_entry_list *list);

/**
 * @brief Creates an AuthOrderSource from the entries (the identifier of an order source, the symbols).
 *
 * The list is only read and copied: the caller keeps its ownership and may free it after the call (a list that
 * dxfg_v4_AuthOrderSource_getByIds() returned is released with
 * dxfg_v4_symbols_by_order_source_id_map_entry_list_release()). The whole list is checked before the AuthOrderSource
 * is created: a NULL list, a NULL entry, a NULL symbol list, a NULL symbol or two entries with the same identifier
 * fail the call with `IllegalArgumentException`. An empty list gives an AuthOrderSource without order sources.
 * dxfg_v4_AuthOrderSource_getByIds() of the result returns the same entries.
 *
 * Java: `AuthOrderSource.create(Map)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbolsByOrderSourceIdMapEntryList The entries.
 * @param[out] authOrderSource The pointer to the variable that receives the AuthOrderSource (release the handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_AuthOrderSource_create(
    graal_isolatethread_t *thread,
    const dxfg_v4_symbols_by_order_source_id_map_entry_list *symbolsByOrderSourceIdMapEntryList,
    DXFG_OUT dxfg_v4_auth_order_source_t **authOrderSource);

/**
 * @brief Returns the entries (an order source, the symbols) of an AuthOrderSource.
 *
 * The list is a copy: it does not depend on the handle of the AuthOrderSource.
 *
 * Java: `AuthOrderSource.getByOrderSources()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authOrderSource The AuthOrderSource.
 * @param[out] symbolsByOrderSourceMapEntryList The list of the entries. Release it with
 * dxfg_v4_symbols_by_order_source_map_entry_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthOrderSource_getByOrderSources(
    graal_isolatethread_t *thread, dxfg_v4_auth_order_source_t *authOrderSource,
    DXFG_OUT dxfg_v4_symbols_by_order_source_map_entry_list **symbolsByOrderSourceMapEntryList);

/**
 * @brief Frees a list of the entries (an order source, the symbols).
 *
 * Frees the list, its entries, their order sources and their symbol lists. A NULL list is ignored.
 *
 * Java: not applicable (the list is allocated by the SDK outside the Java heap).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list of the entries that dxfg_v4_AuthOrderSource_getByOrderSources() returned.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_symbols_by_order_source_map_entry_list_release(graal_isolatethread_t *thread,
                                                               dxfg_v4_symbols_by_order_source_map_entry_list *list);

/**
 * @brief Checks the consistency of a list of orders.
 *
 * PriceLevelChecker is a utility to check that in each moment of time the bid price of the orders is lower than the
 * ask price. The check goes through the orders in the order of the list, so pass a list that is sorted by time, as
 * dxfg_v4_PriceLevelService_getOrdersWithCaller() returns it. Besides the result, the check reports to the Java log
 * (not to the caller) the additional information:
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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] orders The list of the orders (events that are orders, for example the list that
 * dxfg_v4_PriceLevelService_getOrdersWithCaller() returned); the call fails for another kind of event. The caller keeps
 * its ownership.
 * @param[in] timeGapBound The greatest allowed gap between the times of two subsequent orders, in milliseconds (the
 * unit of the time of an order).
 * @param[in] printQuotes Non-zero to log the quotes during the check, 0 not to.
 * @param[out] isValid Set to 1 if the orders are valid, 0 otherwise.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PriceLevelChecker_validate(graal_isolatethread_t *thread, dxfg_v4_event_type_list *orders,
                                           int64_t timeGapBound, int32_t /*bool*/ printQuotes,
                                           DXFG_OUT int32_t *isValid);

/** @} */ // end of V4_OrcsApi

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_ORCS_H
