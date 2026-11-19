// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_events.h and the v4 names of its types (see the group @ref V4_Events).
 */

#ifndef DXFG_V4_EVENTS_H
#define DXFG_V4_EVENTS_H

#include "dxfg_events.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Events Events (v4)
 *  @brief The v4 functions and type names of dxfg_events.h. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_indexed_event_source_type_t */
typedef dxfg_indexed_event_source_type_t dxfg_v4_indexed_event_source_type_t;

/** @copydoc dxfg_indexed_event_source_t */
typedef dxfg_indexed_event_source_t dxfg_v4_indexed_event_source_t;

/** @copydoc dxfg_indexed_event_source_list */
typedef dxfg_indexed_event_source_list dxfg_v4_indexed_event_source_list;

/** @copydoc dxfg_symbol_type_t */
typedef dxfg_symbol_type_t dxfg_v4_symbol_type_t;

/** @copydoc dxfg_symbol_t */
typedef dxfg_symbol_t dxfg_v4_symbol_t;

/** @copydoc dxfg_string_symbol_t */
typedef dxfg_string_symbol_t dxfg_v4_string_symbol_t;

/** @copydoc dxfg_wildcard_symbol_t */
typedef dxfg_wildcard_symbol_t dxfg_v4_wildcard_symbol_t;

/** @copydoc dxfg_candle_symbol_t */
typedef dxfg_candle_symbol_t dxfg_v4_candle_symbol_t;

/** @copydoc dxfg_indexed_event_subscription_symbol_t */
typedef dxfg_indexed_event_subscription_symbol_t dxfg_v4_indexed_event_subscription_symbol_t;

/** @copydoc dxfg_time_series_subscription_symbol_t */
typedef dxfg_time_series_subscription_symbol_t dxfg_v4_time_series_subscription_symbol_t;

/** @copydoc dxfg_symbol_list */
typedef dxfg_symbol_list dxfg_v4_symbol_list;

/** @copydoc dxfg_event_clazz_t */
typedef dxfg_event_clazz_t dxfg_v4_event_clazz_t;

/** @copydoc dxfg_event_type_t */
typedef dxfg_event_type_t dxfg_v4_event_type_t;

/** @copydoc dxfg_market_event_t */
typedef dxfg_market_event_t dxfg_v4_market_event_t;

/** @copydoc dxfg_quote_t */
typedef dxfg_quote_t dxfg_v4_quote_t;

/** @copydoc dxfg_profile_t */
typedef dxfg_profile_t dxfg_v4_profile_t;

/** @copydoc dxfg_summary_t */
typedef dxfg_summary_t dxfg_v4_summary_t;

/** @copydoc dxfg_greeks_t */
typedef dxfg_greeks_t dxfg_v4_greeks_t;

/** @copydoc dxfg_candle_t */
typedef dxfg_candle_t dxfg_v4_candle_t;

/** @copydoc dxfg_daily_candle_t */
typedef dxfg_daily_candle_t dxfg_v4_daily_candle_t;

/** @copydoc dxfg_underlying_t */
typedef dxfg_underlying_t dxfg_v4_underlying_t;

/** @copydoc dxfg_theo_price_t */
typedef dxfg_theo_price_t dxfg_v4_theo_price_t;

/** @copydoc dxfg_trade_base_t */
typedef dxfg_trade_base_t dxfg_v4_trade_base_t;

/** @copydoc dxfg_trade_t */
typedef dxfg_trade_t dxfg_v4_trade_t;

/** @copydoc dxfg_trade_eth_t */
typedef dxfg_trade_eth_t dxfg_v4_trade_eth_t;

/** @copydoc dxfg_configuration_t */
typedef dxfg_configuration_t dxfg_v4_configuration_t;

/** @copydoc dxfg_message_t */
typedef dxfg_message_t dxfg_v4_message_t;

/** @copydoc dxfg_time_and_sale_t */
typedef dxfg_time_and_sale_t dxfg_v4_time_and_sale_t;

/** @copydoc dxfg_option_sale_t */
typedef dxfg_option_sale_t dxfg_v4_option_sale_t;

/** @copydoc dxfg_market_maker_t */
typedef dxfg_market_maker_t dxfg_v4_market_maker_t;

/** @copydoc dxfg_order_base_t */
typedef dxfg_order_base_t dxfg_v4_order_base_t;

/** @copydoc dxfg_order_t */
typedef dxfg_order_t dxfg_v4_order_t;

/** @copydoc dxfg_analytic_order_t */
typedef dxfg_analytic_order_t dxfg_v4_analytic_order_t;

/** @copydoc dxfg_otc_markets_order_t */
typedef dxfg_otc_markets_order_t dxfg_v4_otc_markets_order_t;

/** @copydoc dxfg_spread_order_t */
typedef dxfg_spread_order_t dxfg_v4_spread_order_t;

/** @copydoc dxfg_series_t */
typedef dxfg_series_t dxfg_v4_series_t;

/** @copydoc dxfg_text_message_t */
typedef dxfg_text_message_t dxfg_v4_text_message_t;

/** @copydoc dxfg_text_configuration_t */
typedef dxfg_text_configuration_t dxfg_v4_text_configuration_t;

/** @copydoc dxfg_nuam_order_t */
typedef dxfg_nuam_order_t dxfg_v4_nuam_order_t;

/** @copydoc dxfg_nuam_time_and_sale_t */
typedef dxfg_nuam_time_and_sale_t dxfg_v4_nuam_time_and_sale_t;

/** @copydoc dxfg_nuam_trade_t */
typedef dxfg_nuam_trade_t dxfg_v4_nuam_trade_t;

/** @copydoc dxfg_order_imbalance_t */
typedef dxfg_order_imbalance_t dxfg_v4_order_imbalance_t;

/** @copydoc dxfg_event_type_list */
typedef dxfg_event_type_list dxfg_v4_event_type_list;

/** @copydoc dxfg_event_clazz_list_t */
typedef dxfg_event_clazz_list_t dxfg_v4_event_clazz_list;

/** @copydoc dxfg_observable_subscription_change_listener_t */
typedef dxfg_observable_subscription_change_listener_t dxfg_v4_observable_subscription_change_listener_t;

/** @copydoc dxfg_ObservableSubscriptionChangeListener_function_symbolsAdded */
typedef dxfg_ObservableSubscriptionChangeListener_function_symbolsAdded
    dxfg_v4_ObservableSubscriptionChangeListener_function_symbolsAdded;

/** @copydoc dxfg_ObservableSubscriptionChangeListener_function_symbolsRemoved */
typedef dxfg_ObservableSubscriptionChangeListener_function_symbolsRemoved
    dxfg_v4_ObservableSubscriptionChangeListener_function_symbolsRemoved;

/** @copydoc dxfg_ObservableSubscriptionChangeListener_function_subscriptionClosed */
typedef dxfg_ObservableSubscriptionChangeListener_function_subscriptionClosed
    dxfg_v4_ObservableSubscriptionChangeListener_function_subscriptionClosed;

/**
 * @brief The callback of a subscription change listener: the configuration of the subscription changed.
 *
 * The configuration is the non-functional settings of the subscription, such as its aggregation period: for example,
 * dxfg_v4_DXFeedSubscription_setAggregationPeriod() with a value other than the current one calls it, the same value
 * does not. It is called by the thread that changed the configuration, before that function returns, while the SDK
 * holds the lock of the subscription.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_data The user data passed to
 * dxfg_v4_ObservableSubscriptionChangeListener_newWithConfigurationChanged().
 * @since 3.13.0
 */
typedef void (*dxfg_v4_ObservableSubscriptionChangeListener_function_configurationChanged)(
    graal_isolatethread_t *thread, void *user_data);

/**
 * Creates a symbol of the given type.
 *
 * Only #STRING, #CANDLE and #WILDCARD can be created here: for any other type the function fails.
 * For #CANDLE the string is parsed as a candle symbol (an invalid string is an error) and the symbol of the created
 * structure is its normalized string form. For #WILDCARD the string is ignored.
 *
 * Java: `CandleSymbol.valueOf(String)` for #CANDLE, `WildcardSymbol.ALL` for #WILDCARD.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbol The symbol: a UTF-8 null-terminated string.
 * @param[in] symbolType The type of the symbol to create.
 * @param[out] newSymbol The pointer to the variable that receives the symbol (release it with
 * dxfg_v4_Symbol_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Symbol_new(graal_isolatethread_t *thread, const char *symbol, dxfg_v4_symbol_type_t symbolType,
                           DXFG_OUT dxfg_v4_symbol_t **newSymbol);

/**
 * Releases a symbol: frees the structure, its strings and the nested symbols and sources. Passing NULL does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbol The symbol created by dxfg_v4_Symbol_new() or returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Symbol_release(graal_isolatethread_t *thread, dxfg_v4_symbol_t *symbol);

/**
 * Creates an event of the given class for the given symbol.
 *
 * The event has the given `clazz` and event symbol (a copy of the string), all the other fields are zero (not NaN
 * and not the defaults of the Java event). Fill the fields, then pass the event to a function that needs it.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbolName The event symbol: a UTF-8 null-terminated string. For #DXFG_EVENT_CANDLE and
 * #DXFG_EVENT_DAILY_CANDLE it is the string form of a candle symbol.
 * @param[in] clazz The class of the event to create.
 * @param[out] event The pointer to the variable that receives the event (release it with dxfg_v4_EventType_release());
 * must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_EventType_new(graal_isolatethread_t *thread, const char *symbolName, dxfg_v4_event_clazz_t clazz,
                              DXFG_OUT dxfg_v4_event_type_t **event);

/**
 * Releases an event: frees the structure together with all its strings. Passing NULL does nothing.
 *
 * Use it for the events created by dxfg_v4_EventType_new() and for the events that the SDK returns to the caller as a
 * result. Do not release an event that the SDK passed to a listener: the SDK frees it after the listener returns.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The event to release.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_EventType_release(graal_isolatethread_t *thread, dxfg_v4_event_type_t *eventType);

/**
 * Releases a list of events: frees the list, its array and every event of the list together with its strings.
 * Passing NULL does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventTypes The list of events returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_event_type_list_release(graal_isolatethread_t *thread, dxfg_v4_event_type_list *eventTypes);

/**
 * Releases a list of event classes: frees the list, its array and every element. Passing NULL does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventClazzes The list of event classes returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_event_clazz_list_release(graal_isolatethread_t *thread, dxfg_v4_event_clazz_list *eventClazzes);

/**
 * Releases a list of symbols: frees the list, its array and every symbol of the list with its strings and the nested
 * symbols and sources.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbols The list of symbols returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_symbol_list_release(graal_isolatethread_t *thread, dxfg_v4_symbol_list *symbols);

/**
 * Returns the source of an indexed event.
 *
 * It works for the events of #DXFG_EVENT_GREEKS, #DXFG_EVENT_CANDLE, #DXFG_EVENT_DAILY_CANDLE,
 * #DXFG_EVENT_TIME_AND_SALE, #DXFG_EVENT_UNDERLYING, #DXFG_EVENT_THEO_PRICE and #DXFG_EVENT_SERIES (the default
 * source) and of #DXFG_EVENT_ORDER_BASE, #DXFG_EVENT_ORDER, #DXFG_EVENT_ANALYTIC_ORDER and #DXFG_EVENT_SPREAD_ORDER
 * (the source stored in the index of the order). For the events of any other class the function fails.
 *
 * Java: `IndexedEvent.getSource()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The event.
 * @param[out] source The pointer to the variable that receives the source (release it with
 * dxfg_v4_IndexedEventSource_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEvent_getSource(graal_isolatethread_t *thread, dxfg_v4_event_type_t *eventType,
                                       DXFG_OUT dxfg_v4_indexed_event_source_t **source);

/**
 * Creates an event source by its name.
 *
 * Java: `OrderSource.valueOf(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The name of the order source: a UTF-8 null-terminated string, either predefined or at most four
 * alphanumeric characters. If it is NULL, the default source (the identifier 0) is returned instead.
 * @param[out] newSource The pointer to the variable that receives the source (release it with
 * dxfg_v4_IndexedEventSource_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, for an invalid name). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventSource_new(graal_isolatethread_t *thread, const char *source,
                                       DXFG_OUT dxfg_v4_indexed_event_source_t **newSource);

/**
 * Creates an event source by its identifier.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sourceId The identifier of the source. If it is 0, the default source is returned, otherwise it is the
 * identifier of an order source (a negative identifier is an error).
 * @param[out] source The pointer to the source (release it with dxfg_v4_IndexedEventSource_release()). It is not
 * written on error.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `OrderSource.valueOf(int)`.
 */
int32_t dxfg_v4_IndexedEventSource_newWithId(graal_isolatethread_t *thread, int32_t sourceId,
                                             DXFG_OUT dxfg_v4_indexed_event_source_t** source);

/**
 * Releases an event source: frees the structure and its name. Passing NULL does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The source created by dxfg_v4_IndexedEventSource_new(), dxfg_v4_IndexedEventSource_newWithId() or
 * returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IndexedEventSource_release(graal_isolatethread_t *thread, dxfg_v4_indexed_event_source_t *source);

/**
 * Releases a list of event sources: frees the list, its array and every source of the list with its name.
 * Passing NULL does nothing.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list of sources returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_indexed_event_source_list_release(graal_isolatethread_t *thread,
                                                  dxfg_v4_indexed_event_source_list* list);

/**
 * Creates a listener of the changes of a subscription.
 *
 * Add the listener to a subscription (for example, with dxfg_v4_DXFeedSubscription_addChangeListener()). The SDK calls
 * the functions without a check for NULL, so pass all three. Releasing the last handle of the listener with
 * dxfg_v4_JavaObjectHandler_release() (or dxfg_v4_NativeListener_deactivate()) stops the calls: after that none of the
 * functions is called or running, and the user data can be freed.
 *
 * Java: `ObservableSubscriptionChangeListener` (symbolsAdded(Set), symbolsRemoved(Set), subscriptionClosed()).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] function_symbolsAdded The function called when symbols are added to the subscription.
 * @param[in] function_symbolsRemoved The function called when symbols are removed from the subscription.
 * @param[in] function_subscriptionClosed The function called when the subscription is closed or the listener is
 * removed.
 * @param[in] user_data The user data passed to every function; the SDK does not use or free it.
 * @param[out] listener The pointer to the variable that receives the listener handle (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ObservableSubscriptionChangeListener_new(
    graal_isolatethread_t *thread,
    dxfg_v4_ObservableSubscriptionChangeListener_function_symbolsAdded function_symbolsAdded,
    dxfg_v4_ObservableSubscriptionChangeListener_function_symbolsRemoved function_symbolsRemoved,
    dxfg_v4_ObservableSubscriptionChangeListener_function_subscriptionClosed function_subscriptionClosed,
    void *user_data, DXFG_OUT dxfg_v4_observable_subscription_change_listener_t **listener);

/**
 * @brief Creates a listener of the changes of a subscription that is also told about the changes of its configuration.
 *
 * It is dxfg_v4_ObservableSubscriptionChangeListener_new() with the fourth callback,
 * `function_configurationChanged`. Each of the four callbacks may be NULL: the SDK does not call a NULL one. Add the
 * listener to a subscription with dxfg_v4_DXFeedSubscription_addChangeListener(). Releasing the last handle of the
 * listener with dxfg_v4_JavaObjectHandler_release(), or dxfg_v4_NativeListener_deactivate() or
 * dxfg_v4_NativeListener_deactivateAll(), stops the calls: after that none of the functions is called or running
 * (except a call on the deactivating thread itself), and the user data can be freed.
 *
 * Java: `ObservableSubscriptionChangeListener` (symbolsAdded(Set), symbolsRemoved(Set), subscriptionClosed(),
 * configurationChanged()).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] function_symbolsAdded The function called when symbols are added to the subscription, or NULL.
 * @param[in] function_symbolsRemoved The function called when symbols are removed from the subscription, or NULL.
 * @param[in] function_subscriptionClosed The function called when the subscription is closed or the listener is
 * removed, or NULL.
 * @param[in] function_configurationChanged The function called when the configuration of the subscription changes, or
 * NULL.
 * @param[in] user_data The user data passed to every function; the SDK does not use or free it.
 * @param[out] listener The pointer to the variable that receives the listener handle (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_ObservableSubscriptionChangeListener_newWithConfigurationChanged(
    graal_isolatethread_t *thread,
    dxfg_v4_ObservableSubscriptionChangeListener_function_symbolsAdded function_symbolsAdded,
    dxfg_v4_ObservableSubscriptionChangeListener_function_symbolsRemoved function_symbolsRemoved,
    dxfg_v4_ObservableSubscriptionChangeListener_function_subscriptionClosed function_subscriptionClosed,
    dxfg_v4_ObservableSubscriptionChangeListener_function_configurationChanged function_configurationChanged,
    void *user_data, DXFG_OUT dxfg_v4_observable_subscription_change_listener_t **listener);

/** @} */ // end of V4_Events

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_EVENTS_H
