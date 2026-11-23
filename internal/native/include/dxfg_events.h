// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Events functions and types declarations
 */

#ifndef DXFG_EVENTS_H
#define DXFG_EVENTS_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_javac.h"

/**
 * @brief The kind of an event source, see dxfg_indexed_event_source_t.
 */
typedef enum dxfg_indexed_event_source_type_t {
    /** A plain IndexedEventSource (for example, the default source with the identifier 0). */
    INDEXED_EVENT_SOURCE = 0,
    /** An OrderSource: the source of order events. */
    ORDER_SOURCE
} dxfg_indexed_event_source_type_t;

/**
 * @brief The source of indexed events.
 *
 * Create it with dxfg_IndexedEventSource_new() or dxfg_IndexedEventSource_new2(), release it with
 * dxfg_IndexedEventSource_release(). The sources that the SDK returns (for example, by
 * dxfg_IndexedEvent_getSource()) are released the same way.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/IndexedEventSource.html">Javadoc</a>
 */
typedef struct dxfg_indexed_event_source_t {
    /** The kind of the source. */
    dxfg_indexed_event_source_type_t type;
    /** The source identifier, a non-negative number; 0 is the default source. */
    int32_t id;
    /** The name of the source (a UTF-8 null-terminated string) that belongs to the structure and is freed with it. */
    const char *name;
} dxfg_indexed_event_source_t;

/**
 * @brief A list of event sources.
 *
 * The list returned by the SDK is released with dxfg_CList_IndexedEventSource_release().
 */
typedef struct dxfg_indexed_event_source_list {
    /** The number of the elements. */
    int32_t size;
    /** The array of `size` pointers to the sources; NULL when `size` is 0 in a list returned by the SDK. */
    dxfg_indexed_event_source_t **elements;
} dxfg_indexed_event_source_list;

/**
 * @brief The kind of a symbol, see dxfg_symbol_t.
 */
typedef enum dxfg_symbol_type_t {
    /** A string symbol, the structure is dxfg_string_symbol_t. */
    STRING = 0,
    /** A candle symbol, the structure is dxfg_candle_symbol_t. */
    CANDLE,
    /** The wildcard symbol (the subscription to all symbols), the structure is dxfg_wildcard_symbol_t. */
    WILDCARD,
    /** The subscription to the events of one source, the structure is dxfg_indexed_event_subscription_symbol_t. */
    INDEXED_EVENT_SUBSCRIPTION,
    /** The subscription to the time series from a time, the structure is dxfg_time_series_subscription_symbol_t. */
    TIME_SERIES_SUBSCRIPTION
} dxfg_symbol_type_t;

/**
 * @brief The base structure of every symbol: the first field of all the symbol structures.
 *
 * `type` tells which of the structures dxfg_string_symbol_t, dxfg_candle_symbol_t, dxfg_wildcard_symbol_t,
 * dxfg_indexed_event_subscription_symbol_t and dxfg_time_series_subscription_symbol_t the pointer really points
 * to; cast the pointer to it to read the other fields.
 *
 * A symbol created by dxfg_Symbol_new() or returned by the SDK is released with dxfg_Symbol_release() (a list of
 * symbols with dxfg_CList_symbol_release()): that frees the structure, its strings and the nested symbols and sources.
 */
typedef struct dxfg_symbol_t {
    /** The kind of the symbol. */
    dxfg_symbol_type_t type;
} dxfg_symbol_t;

/**
 * @brief A string symbol (for example, a ticker).
 */
typedef struct dxfg_string_symbol_t {
    /** The base structure: the fields of dxfg_symbol_t. */
    dxfg_symbol_t supper;
    /** The symbol (a UTF-8 null-terminated string) that belongs to the structure and is freed with it. */
    const char *symbol;
} dxfg_string_symbol_t;

/**
 * @brief The wildcard symbol: the subscription to all symbols of an event type.
 *
 * It has no fields besides the base structure.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/WildcardSymbol.html">Javadoc</a>
 */
typedef struct dxfg_wildcard_symbol_t {
    /** The base structure: the fields of dxfg_symbol_t. */
    dxfg_symbol_t supper;
} dxfg_wildcard_symbol_t;

/**
 * @brief A candle symbol.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/candle/CandleSymbol.html">Javadoc</a>
 */
typedef struct dxfg_candle_symbol_t {
    /** The base structure: the fields of dxfg_symbol_t. */
    dxfg_symbol_t supper;
    /**
     * The string form of the candle symbol, as CandleSymbol.toString() returns it (a UTF-8 null-terminated string);
     * it belongs to the structure and is freed with it.
     */
    const char *symbol;
} dxfg_candle_symbol_t;

/**
 * @brief The subscription to the indexed events of one source.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/IndexedEventSubscriptionSymbol.html">Javadoc</a>
 */
typedef struct dxfg_indexed_event_subscription_symbol_t {
    /** The base structure: the fields of dxfg_symbol_t. */
    dxfg_symbol_t supper;
    /** The event symbol; it belongs to this structure and is freed with it. */
    dxfg_symbol_t *symbol;
    /** The source of the events; it belongs to this structure and is freed with it. */
    dxfg_indexed_event_source_t *source;
} dxfg_indexed_event_subscription_symbol_t;

/**
 * @brief The subscription to the time series of events from a time.
 *
 * The source of such a subscription is always the default one.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/TimeSeriesSubscriptionSymbol.html">Javadoc</a>
 */
typedef struct dxfg_time_series_subscription_symbol_t {
    /** The base structure: the fields of dxfg_symbol_t. */
    dxfg_symbol_t supper;
    /** The event symbol; it belongs to this structure and is freed with it. */
    dxfg_symbol_t *symbol;
    /** The time of the subscription, in milliseconds since the Unix epoch (UTC). */
    int64_t from_time;
} dxfg_time_series_subscription_symbol_t;

/**
 * @brief A list of symbols.
 *
 * The list returned by the SDK (or passed by the SDK to a listener) is released with dxfg_CList_symbol_release().
 */
typedef struct dxfg_symbol_list {
    /** The number of the elements. */
    int32_t size;
    /** The array of `size` pointers to the symbols; NULL when `size` is 0 in a list created by the SDK. */
    dxfg_symbol_t **elements;
} dxfg_symbol_list;

/**
 * @brief The class of an event: tells which structure an event (dxfg_event_type_t) is.
 *
 * In the Java implementation, all events are lined up in a hierarchy from the underlying EventType interface. In the
 * C implementation, we specify the type in the event in a separate field because we cannot get it from the class
 * name.
 *
 * The comment of every value names the Java event class and its kind: a lasting event (LastingEvent: only the latest
 * event of a symbol is kept), an indexed event (IndexedEvent: events of a symbol are told apart by an index and a
 * source) or a time series event (TimeSeriesEvent: an indexed event with a time). The kind of the base class is
 * inherited by its subclasses. The structure of a value is named by the value (DXFG_EVENT_QUOTE: dxfg_quote_t).
 */
typedef enum dxfg_event_clazz_t {
    DXFG_EVENT_QUOTE = 0,    ///< Quote, a lasting event: the best bid and ask.
    DXFG_EVENT_PROFILE,      ///< Profile, a lasting event: the description of an instrument.
    DXFG_EVENT_SUMMARY,      ///< Summary, a lasting event: the trading session information.
    DXFG_EVENT_GREEKS,       ///< Greeks, a lasting event and a time series event: the option price and greeks.
    DXFG_EVENT_CANDLE,       ///< Candle, a lasting event and a time series event: OHLC and other values of a period.
    DXFG_EVENT_DAILY_CANDLE, ///< DailyCandle (deprecated in Java, use Candle): a Candle with the same fields.
    DXFG_EVENT_UNDERLYING,   ///< Underlying, a lasting event and a time series event: the values of an underlying.
    DXFG_EVENT_THEO_PRICE,   ///< TheoPrice, a lasting event and a time series event: the theoretical option price.
    // abstract DXFG_EVENT_TRADE_BASE,     // LASTING
    DXFG_EVENT_TRADE,              ///< Trade (extends TradeBase), a lasting event: the last trade in RTH.
    DXFG_EVENT_TRADE_ETH,          ///< TradeETH (extends TradeBase), a lasting event: the last trade in ETH.
    DXFG_EVENT_CONFIGURATION,      ///< Configuration, a lasting event: an attachment with a version.
    DXFG_EVENT_MESSAGE,            ///< Message: a notification with an attachment (neither lasting nor indexed).
    DXFG_EVENT_TIME_AND_SALE,      ///< TimeAndSale, a time series event: a trade or another market event.
    DXFG_EVENT_ORDER_BASE,         ///< OrderBase, an indexed event: the common fields of the order events.
    DXFG_EVENT_ORDER,              ///< Order (extends OrderBase), an indexed event: an order or a price level.
    DXFG_EVENT_ANALYTIC_ORDER,     ///< AnalyticOrder (extends Order), an indexed event: an Order with iceberg data.
    DXFG_EVENT_OTC_MARKETS_ORDER,  ///< OtcMarketsOrder (extends Order), an indexed event: an OTC Markets quote.
    DXFG_EVENT_SPREAD_ORDER,       ///< SpreadOrder (extends OrderBase), an indexed event: a spread order.
    DXFG_EVENT_SERIES,             ///< Series, an indexed event: the values of an option series.
    DXFG_EVENT_OPTION_SALE,        ///< OptionSale, an indexed event: a trade of an option.
    DXFG_EVENT_TEXT_MESSAGE,       ///< TextMessage: a notification with a text (neither lasting nor indexed).
    DXFG_EVENT_MARKET_MAKER,       ///< MarketMaker, an indexed event: the top quotes of a market participant.
    DXFG_EVENT_TEXT_CONFIGURATION, ///< TextConfiguration, a lasting event: a text with a version.
    DXFG_EVENT_NUAM_ORDER,         ///< NuamOrder (extends Order), an indexed event: an order of the Nuam Exchange.
    DXFG_EVENT_NUAM_TIME_AND_SALE, ///< NuamTimeAndSale (extends TimeAndSale), a time series event (Nuam Exchange).
    DXFG_EVENT_NUAM_TRADE,         ///< NuamTrade (extends Trade), a lasting event: Nuam Exchange trade statistics.
    DXFG_EVENT_ORDER_IMBALANCE,    ///< OrderImbalance, a market event (neither lasting nor indexed): auction imbalance.
} dxfg_event_clazz_t;

/**
 * @brief The base structure of every event: the first field of all the event structures.
 *
 * `clazz` tells which event structure the pointer really points to (for example, #DXFG_EVENT_QUOTE: dxfg_quote_t);
 * cast the pointer to that structure to read the other fields.
 *
 * Who creates and releases an event:
 * - dxfg_EventType_new() creates an event that the caller owns; all its fields besides `clazz` and the symbol are
 *   zero (not NaN).
 * - The SDK creates the events that it passes to a listener or returns as a result. An event passed to a listener
 *   belongs to the SDK: it is freed after the listener returns, copy what you need to keep. The events returned as
 *   a result belong to the caller.
 * - An event of either origin, and a list of events, is released with dxfg_EventType_release() and
 *   dxfg_CList_EventType_release(): that frees the structure together with all its strings (the `const char *`
 *   fields, such as the event symbol).
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/EventType.html">Javadoc</a>
 */
typedef struct dxfg_event_type_t {
    /** The class of the event: which event structure this is. */
    dxfg_event_clazz_t clazz;
} dxfg_event_type_t;

/**
 * @brief The base structure of the market events: the event symbol and the event time.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/MarketEvent.html">Javadoc</a>
 */
typedef struct dxfg_market_event_t {
    /** The base structure: the fields of dxfg_event_type_t. */
    dxfg_event_type_t event_type;
    /** The symbol of the event (a UTF-8 null-terminated string) that belongs to the event and is freed with it. */
    const char *event_symbol;
    /**
     * The time when the event was created, in milliseconds since the Unix epoch (UTC), or 0 if not available.
     * It is available only if the endpoint is created with the event time property and the data source has embedded
     * event times, which is typically the case for historical tape files and OnDemandService only.
     */
    int64_t event_time;
} dxfg_market_event_t;

///**
// * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/LastingEvent.html">Javadoc</a>
// */
// typedef struct dxfg_lasting_event_t {
//} dxfg_lasting_event_t;
//
///**
// * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/IndexedEvent.html">Javadoc</a>
// */
// typedef struct dxfg_indexed_event_t {
//    //    dxfg_indexed_event_source_t source;
//    //    int32_t event_flags;
//    //    int64_t index;
//} dxfg_indexed_event_t;
//
///**
// * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/TimeSeriesEvent.html">Javadoc</a>
// */
// typedef struct dxfg_time_series_event_t {
//    dxfg_indexed_event_t indexed_event;
//    //    int64_t time;
//} dxfg_time_series_event_t;

/**
 * @brief The quote event: a snapshot of the best bid and ask prices and sizes.
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_QUOTE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Quote.html">Javadoc</a>
 */
typedef struct dxfg_quote_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_lasting_event_t lasting_event;
    /**
     * The milliseconds part of the time of the last bid or ask change (bits 22 to 31) and the sequence (bits 0 to 21)
     * packed into one number; see Quote.getTime() and Quote.getSequence() in Java.
     */
    int32_t time_millis_sequence;
    /** The microseconds and nanoseconds part of the time of the last bid or ask change. */
    int32_t time_nano_part;
    /**
     * The time of the last bid change, in milliseconds since the Unix epoch (UTC). By default it has a precision of
     * seconds, so it is usually a multiple of 1000.
     */
    int64_t bid_time;
    /** The bid exchange code: the character code (Java `char`) of the exchange, 0 if not set. */
    int16_t bid_exchange_code;
    /** The bid price, NaN if unknown. */
    double bid_price;
    /** The bid size (with fractions), NaN if unknown. */
    double bid_size;
    /**
     * The time of the last ask change, in milliseconds since the Unix epoch (UTC). By default it has a precision of
     * seconds, so it is usually a multiple of 1000.
     */
    int64_t ask_time;
    /** The ask exchange code: the character code (Java `char`) of the exchange, 0 if not set. */
    int16_t ask_exchange_code;
    /** The ask price, NaN if unknown. */
    double ask_price;
    /** The ask size (with fractions), NaN if unknown. */
    double ask_size;
} dxfg_quote_t;

/**
 * @brief The profile event: a snapshot of the description of a security instrument.
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_PROFILE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Profile.html">Javadoc</a>
 */
typedef struct dxfg_profile_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_lasting_event_t lasting_event;
    /** The description of the instrument (a UTF-8 string, NULL if none) that belongs to the event. */
    const char *description;
    /** The reason that trading was halted (a UTF-8 string, NULL if none) that belongs to the event. */
    const char *status_reason;
    /** The start of the trading halt interval, in milliseconds since the Unix epoch (UTC), 0 if not set. */
    int64_t halt_start_time;
    /** The end of the trading halt interval, in milliseconds since the Unix epoch (UTC), 0 if not set. */
    int64_t halt_end_time;
    /** The maximal (high) allowed price, NaN if unknown. */
    double high_limit_price;
    /** The minimal (low) allowed price, NaN if unknown. */
    double low_limit_price;
    /** The maximal (high) price in the last 52 weeks, NaN if unknown. */
    double high_52_week_price;
    /** The minimal (low) price in the last 52 weeks, NaN if unknown. */
    double low_52_week_price;
    /** The correlation coefficient of the instrument to the S&P500 index, NaN if unknown. */
    double beta;
    /** The earnings per share (the profits of the company divided by the number of shares), NaN if unknown. */
    double earnings_per_share;
    /** The frequency of cash dividend payments per year (calculated), NaN if unknown. */
    double dividend_frequency;
    /** The amount of the last paid dividend, NaN if unknown. */
    double ex_dividend_amount;
    /** The identifier of the day of the last dividend payment (the ex-dividend date): the days since 1 Jan 1970. */
    int32_t ex_dividend_day_id;
    /** The number of shares outstanding, NaN if unknown. */
    double shares;
    /** The free float: the number of shares outstanding that are available to the public for trade, NaN if unknown. */
    double free_float;
    /**
     * The implementation-specific flags: the trading status in bits 0 and 1 and the short sale restriction in bits
     * 2 and 3; see Profile.getTradingStatus() and Profile.getShortSaleRestriction() in Java.
     */
    int32_t flags;
} dxfg_profile_t;

/**
 * @brief The summary event: a snapshot of the trading session information (highs, lows and so on).
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_SUMMARY.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Summary.html">Javadoc</a>
 */
typedef struct dxfg_summary_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_lasting_event_t lasting_event;
    /** The identifier of the day that the summary represents: the number of days since 1 Jan 1970. */
    int32_t day_id;
    /** The first (open) price of the day, NaN if unknown. */
    double day_open_price;
    /** The maximal (high) price of the day, NaN if unknown. */
    double day_high_price;
    /** The minimal (low) price of the day, NaN if unknown. */
    double day_low_price;
    /** The last (close) price of the day, NaN if unknown. */
    double day_close_price;
    /** The identifier of the previous day: the number of days since 1 Jan 1970. */
    int32_t prev_day_id;
    /** The last (close) price of the previous day, NaN if unknown. */
    double prev_day_close_price;
    /** The total volume traded in the previous day, NaN if unknown. */
    double prev_day_volume;
    /** The open interest of the symbol: the number of open contracts. */
    int64_t open_interest;
    /**
     * The implementation-specific flags: the price type of the previous day close price in bits 0 and 1 and the price
     * type of the day close price in bits 2 and 3; see Summary.getDayClosePriceType() in Java.
     */
    int32_t flags;
} dxfg_summary_t;

/**
 * @brief The Greeks event: a snapshot of the option price, the Black-Scholes volatility and the greeks.
 *
 * It is a lasting event and a time series event, its `clazz` is #DXFG_EVENT_GREEKS.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/Greeks.html">Javadoc</a>
 */
typedef struct dxfg_greeks_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_time_series_event_t time_series_event;
    //    dxfg_lasting_event_t lasting_event;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;

    /**
     * The unique per-symbol index of the event: the time in seconds since the Unix epoch in the higher 32 bits, then
     * the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see Greeks.getIndex() in Java.
     */
    int64_t index;
    /** The option market price, NaN if unknown. */
    double price;
    /** The Black-Scholes implied volatility of the option, NaN if unknown. */
    double volatility;
    /** The delta of the option: the first derivative of the option price by the underlying price, NaN if unknown. */
    double delta;
    /** The gamma of the option: the second derivative of the option price by the underlying price, NaN if unknown. */
    double gamma;
    /** The theta of the option: the first derivative of the option price by the days to expiration, NaN if unknown. */
    double theta;
    /**
     * The rho of the option: the first derivative of the option price by the percentage interest rate,
     * NaN if unknown.
     */
    double rho;
    /**
     * The vega of the option: the first derivative of the option price by the percentage volatility,
     * NaN if unknown.
     */
    double vega;
} dxfg_greeks_t;

/**
 * @brief The candle event: OHLC prices and other information for a specific period.
 *
 * It is a lasting event and a time series event, its `clazz` is #DXFG_EVENT_CANDLE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/candle/Candle.html">Javadoc</a>
 */
typedef struct dxfg_candle_t {
    /** The base structure: the fields of dxfg_event_type_t. */
    dxfg_event_type_t event_type;
    //    dxfg_time_series_event_t time_series_event;
    //    dxfg_lasting_event_t lasting_event;
    /**
     * The candle symbol in its string form, as CandleSymbol.toString() returns it (a UTF-8 null-terminated string);
     * it belongs to the event and is freed with it.
     */
    const char *event_symbol;
    /** The time when the event was created, in milliseconds since the Unix epoch (UTC), or 0 if not available. */
    int64_t event_time;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;
    /**
     * The unique per-symbol index of the candle: the time in seconds since the Unix epoch in the higher 32 bits, then
     * the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see Candle.getIndex() in Java.
     */
    int64_t index;
    /** The total number of the original trade (or quote) events in the candle. */
    int64_t count;
    /** The first (open) price of the candle, NaN if unknown. */
    double open;
    /** The maximal (high) price of the candle, NaN if unknown. */
    double high;
    /** The minimal (low) price of the candle, NaN if unknown. */
    double low;
    /** The last (close) price of the candle, NaN if unknown. */
    double close;
    /** The total volume of the candle (with fractions), NaN if unknown. */
    double volume;
    /** The volume-weighted average price of the candle, NaN if unknown. */
    double vwap;
    /** The bid volume of the candle (with fractions), NaN if unknown. */
    double bid_volume;
    /** The ask volume of the candle (with fractions), NaN if unknown. */
    double ask_volume;
    /** The implied volatility, NaN if unknown. */
    double imp_volatility;
    /** The open interest (with fractions), NaN if unknown. */
    double open_interest;
} dxfg_candle_t;

/**
 * @brief The daily candle event (deprecated in Java, use dxfg_candle_t): a candle with the same fields.
 *
 * Its `clazz` is #DXFG_EVENT_DAILY_CANDLE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/candle/DailyCandle.html">Javadoc</a>
 */
typedef struct dxfg_daily_candle_t {
    /** The base structure: the fields of dxfg_candle_t. */
    dxfg_candle_t candle;
} dxfg_daily_candle_t;

/**
 * @brief The underlying event: a snapshot of the values computed for an option underlying symbol.
 *
 * It is a lasting event and a time series event, its `clazz` is #DXFG_EVENT_UNDERLYING.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/Underlying.html">Javadoc</a>
 */
typedef struct dxfg_underlying_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_time_series_event_t time_series_event;
    //    dxfg_lasting_event_t lasting_event;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;

    /**
     * The unique per-symbol index of the event: the time in seconds since the Unix epoch in the higher 32 bits, then
     * the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see Underlying.getIndex()
     * in Java.
     */
    int64_t index;
    /** The 30-day implied volatility of the underlying based on the VIX methodology, NaN if unknown. */
    double volatility;
    /** The front month implied volatility of the underlying based on the VIX methodology, NaN if unknown. */
    double front_volatility;
    /** The back month implied volatility of the underlying based on the VIX methodology, NaN if unknown. */
    double back_volatility;
    /** The traded volume of the call options for a day, NaN if unknown. */
    double call_volume;
    /** The traded volume of the put options for a day, NaN if unknown. */
    double put_volume;
    /** The ratio of the traded volume of the put options to the traded volume of the call options, NaN if unknown. */
    double put_call_ratio;
} dxfg_underlying_t;

/**
 * @brief The theo price event: a snapshot of the theoretical option price computation.
 *
 * It is a lasting event and a time series event, its `clazz` is #DXFG_EVENT_THEO_PRICE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/TheoPrice.html">Javadoc</a>
 */
typedef struct dxfg_theo_price_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_time_series_event_t time_series_event;
    //    dxfg_lasting_event_t lasting_event;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;

    /**
     * The unique per-symbol index of the event: the time in seconds since the Unix epoch in the higher 32 bits, then
     * the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see TheoPrice.getIndex()
     * in Java.
     */
    int64_t index;
    /** The theoretical option price, NaN if unknown. */
    double price;
    /** The underlying price at the time of the computation, NaN if unknown. */
    double underlying_price;
    /** The delta of the theoretical price: its first derivative by the underlying price, NaN if unknown. */
    double delta;
    /** The gamma of the theoretical price: its second derivative by the underlying price, NaN if unknown. */
    double gamma;
    /** The implied simple dividend return of the corresponding option series, NaN if unknown. */
    double dividend;
    /** The implied simple interest return of the corresponding option series, NaN if unknown. */
    double interest;
} dxfg_theo_price_t;

/**
 * @brief The common fields of the trade events: dxfg_trade_t and dxfg_trade_eth_t.
 *
 * It is a lasting event, the information about the last trade. There is no event class of its own in the C API.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/TradeBase.html">Javadoc</a>
 */
typedef struct dxfg_trade_base_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_lasting_event_t lasting_event;
    /**
     * The time of the last trade and the sequence packed into one number: the time in seconds since the Unix epoch in
     * the higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see
     * TradeBase.getTimeSequence() in Java.
     */
    int64_t time_sequence;
    /** The microseconds and nanoseconds part of the time of the last trade. */
    int32_t time_nano_part;
    /** The exchange code of the last trade: the character code (Java `char`) of the exchange, 0 if not set. */
    int16_t exchange_code;
    /** The price of the last trade, NaN if unknown. */
    double price;
    /** The change of the last trade, NaN if unknown. */
    double change;
    /** The size of the last trade (with fractions), NaN if unknown. */
    double size;
    /** The identifier of the current trading day: the number of days since 1 Jan 1970. */
    int32_t day_id;
    /** The total volume traded for the day (with fractions), NaN if unknown. */
    double day_volume;
    /** The total turnover traded for the day, NaN if unknown. */
    double day_turnover;
    /**
     * The implementation-specific flags: the extended trading hours indicator in bit 0 and the tick direction in
     * bits 1 to 3; see TradeBase.getTickDirection() and TradeBase.isExtendedTradingHours() in Java.
     */
    int32_t flags;
    /** The trade ID for the events with a trade-related action (execution, correction, cancellation), 0 if none. */
    int64_t trade_id;
} dxfg_trade_base_t;

/**
 * @brief The trade event: the last trade during the regular trading hours (RTH) and the volume and turnover of the day.
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_TRADE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Trade.html">Javadoc</a>
 */
typedef struct dxfg_trade_t {
    /** The base structure: the fields of dxfg_trade_base_t. */
    dxfg_trade_base_t trade_base;
} dxfg_trade_t;

/**
 * @brief The TradeETH event: the last trade during the extended trading hours (ETH) and their volume and turnover.
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_TRADE_ETH.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/TradeETH.html">Javadoc</a>
 */
typedef struct dxfg_trade_eth_t {
    /** The base structure: the fields of dxfg_trade_base_t. */
    dxfg_trade_base_t trade_base;
} dxfg_trade_eth_t;

/**
 * @brief The configuration event: an application-specific attachment with a version.
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_CONFIGURATION.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/misc/Configuration.html">Javadoc</a>
 *
 * The `attachment` is handled as in dxfg_message_t. For a string payload prefer dxfg_text_configuration_t.
 */
typedef struct dxfg_configuration_t {
    /** The base structure: the fields of dxfg_event_type_t. */
    dxfg_event_type_t event_type;
    //    dxfg_lasting_event_t lasting_event;
    /** The symbol of the event (a UTF-8 null-terminated string) that belongs to the event and is freed with it. */
    const char *event_symbol;
    /** The time when the event was created, in milliseconds since the Unix epoch (UTC), or 0 if not available. */
    int64_t event_time;
    /** The version of the event: a higher value is a more recent update. */
    int32_t version;
    /** The attachment: a UTF-8 null-terminated string (`const char *`) or NULL, see dxfg_message_t. */
    void *attachment;
} dxfg_configuration_t;

/**
 * @brief The message event: a notification with an application-specific attachment.
 *
 * It is neither a lasting nor an indexed event: messages are never conflated, its `clazz` is #DXFG_EVENT_MESSAGE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/misc/Message.html">Javadoc</a>
 *
 * The `attachment` (the same for dxfg_configuration_t) is a UTF-8 null-terminated string (`const char *`, declared
 * as `void *`) or NULL when there is no attachment:
 * - Java -> C: a `String` attachment is passed as is, the other objects are passed as their string representation
 *   (`toString()`), no attachment is NULL. An attachment that Java cannot deserialize (e.g. its class is not
 *   registered for the native image) is NULL too, the SDK logs a warning then.
 * - C -> Java: the attachment becomes a `String`, NULL is no attachment.
 *
 * So a non-string attachment loses its type when it passes through C: if the event is published back
 * (or forwarded to another endpoint), the Java receivers get the string instead of the original object.
 *
 * The C API does not support arbitrary Java attachments: C can neither create nor read them. For a string payload
 * prefer dxfg_text_message_t and dxfg_text_configuration_t (TextMessage, TextConfiguration): their `text` is a plain
 * string that passes both ways without loss.
 *
 * Memory:
 * - The events that the SDK gives to the caller (to a listener, from dxfg_DXFeed_getLastEvent(), a promise, etc.)
 *   own their strings: the SDK frees them with the event (after the listener returns, or with dxfg_EventType_release()
 *   and dxfg_CList_EventType_release()). Copy the attachment to keep it longer.
 * - The events that the caller allocates itself and passes to the SDK (dxfg_DXPublisher_publishEvents(), etc.) are
 *   only read: the SDK copies the attachment during the call, the caller allocates and frees it.
 * - dxfg_EventType_release() of an event created by dxfg_EventType_new() frees its strings with the SDK's allocator.
 *   If the caller sets the attachment of such an event to its own string, it resets the attachment to NULL before
 *   the release and frees the string itself (the allocators of the caller and of the SDK may differ, e.g. on Windows).
 */
typedef struct dxfg_message_t {
    /** The base structure: the fields of dxfg_event_type_t. */
    dxfg_event_type_t event_type;
    /** The symbol of the event (a UTF-8 null-terminated string) that belongs to the event and is freed with it. */
    const char *event_symbol;
    /** The time when the event was created, in milliseconds since the Unix epoch (UTC), or 0 if not available. */
    int64_t event_time;
    /** The attachment: a UTF-8 null-terminated string (`const char *`) or NULL, see dxfg_message_t. */
    void *attachment;
} dxfg_message_t;

/**
 * @brief The time and sale event: a trade or another market event with a price, in a continuous time slice.
 *
 * It is a time series event, its `clazz` is #DXFG_EVENT_TIME_AND_SALE. Its time is the one in `index`, there is
 * no separate time field.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/TimeAndSale.html">Javadoc</a>
 */
typedef struct dxfg_time_and_sale_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_time_series_event_t time_series_event;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;
    /**
     * The unique per-symbol index of the event: the time of the original event in seconds since the Unix epoch in the
     * higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see
     * TimeAndSale.getIndex() in Java. It can be used to correct or cancel the event later.
     */
    int64_t index;
    /** The microseconds and nanoseconds part of the time of the event. */
    int32_t time_nano_part;
    /** The exchange code: the character code (Java `char`) of the exchange, 0 if not set. */
    int16_t exchange_code;
    /** The price, NaN if unknown. */
    double price;
    /** The size (with fractions), NaN if unknown. */
    double size;
    /** The bid price on the market at the time of the event, NaN if unknown. */
    double bid_price;
    /** The ask price on the market at the time of the event, NaN if unknown. */
    double ask_price;
    /**
     * The sale conditions provided by the data feed (a UTF-8 string, its format is specific for the feed, NULL if none)
     * that belong to the event and are freed with it.
     */
    const char *exchange_sale_conditions;
    /**
     * The implementation-specific flags: the type (new, correction, cancellation) in bits 0 and 1, the valid tick
     * indicator in bit 2, the extended trading hours indicator in bit 3, the spread leg indicator in bit 4,
     * the aggressor side in bits 5 and 6 and the TradeThroughExempt character in bits 8 to 15; see the getters of
     * TimeAndSale in Java (getType(), isValidTick(), getAggressorSide() and so on).
     */
    int32_t flags;
    /** The buyer (a UTF-8 string, NULL if none) that belongs to the event and is freed with it. */
    const char *buyer;
    /** The seller (a UTF-8 string, NULL if none) that belongs to the event and is freed with it. */
    const char *seller;
    /** The trade ID for the events with a trade-related action (execution, correction, cancellation), 0 if none. */
    int64_t trade_id;
} dxfg_time_and_sale_t;

/**
 * @brief The option sale event: a trade or another market event with a price for an option of an underlying.
 *
 * It is an indexed event with the default source, its `clazz` is #DXFG_EVENT_OPTION_SALE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/OptionSale.html">Javadoc</a>
 */
typedef struct dxfg_option_sale_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_indexed_event_t indexed_event;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;

    /**
     * The unique per-symbol index of the event, see OptionSale.getIndex() in Java. It identifies the event for a later
     * correction or cancellation.
     */
    int64_t index;
    /**
     * The time of the event and the sequence packed into one number: the time in seconds since the Unix epoch in the
     * higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see
     * OptionSale.getTimeSequence() in Java.
     */
    int64_t time_sequence;
    /** The microseconds and nanoseconds part of the time of the event. */
    int32_t time_nano_part;
    /** The exchange code: the character code (Java `char`) of the exchange, 0 if not set. */
    int16_t exchange_code;
    /** The price, NaN if unknown. */
    double price;
    /** The size, NaN if unknown. */
    double size;
    /** The bid price on the market at the time of the event, NaN if unknown. */
    double bid_price;
    /** The ask price on the market at the time of the event, NaN if unknown. */
    double ask_price;
    /**
     * The sale conditions provided by the data feed (a UTF-8 string, its format is specific for the feed, NULL if none)
     * that belong to the event and are freed with it.
     */
    const char *exchange_sale_conditions;
    /**
     * The implementation-specific flags, the same as in dxfg_time_and_sale_t: the type (new, correction,
     * cancellation), the valid tick, extended trading hours and spread leg indicators, the aggressor side and the
     * TradeThroughExempt character; see the getters of OptionSale in Java.
     */
    int32_t flags;
    /** The underlying price at the time of the event, NaN if unknown. */
    double underlying_price;
    /** The Black-Scholes implied volatility of the option at the time of the event, NaN if unknown. */
    double volatility;
    /** The delta of the option at the time of the event, NaN if unknown. */
    double delta;
    /** The option symbol of the event (a UTF-8 string, NULL if none) that belongs to the event and is freed with it. */
    const char *option_symbol;
} dxfg_option_sale_t;

/**
 * @brief The market maker event: a snapshot of the aggregated top quotes of a market participant.
 *
 * It is an indexed event, its `clazz` is #DXFG_EVENT_MARKET_MAKER.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/MarketMaker.html">Javadoc</a>
 */
typedef struct dxfg_market_maker_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_indexed_event_t indexed_event;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;

    /** The unique per-symbol index of the event, a non-negative number; see MarketMaker.getIndex() in Java. */
    int64_t index;
    /**
     * The time of the last bid change, in milliseconds since the Unix epoch (UTC). By default it has a precision of
     * seconds, so it is usually a multiple of 1000.
     */
    int64_t bid_time;
    /** The bid price, NaN if unknown. */
    double bid_price;
    /** The bid size, NaN if unknown. */
    double bid_size;
    /** The number of the individual bid orders in the event. */
    int64_t bid_count;
    /**
     * The time of the last ask change, in milliseconds since the Unix epoch (UTC). By default it has a precision of
     * seconds, so it is usually a multiple of 1000.
     */
    int64_t ask_time;
    /** The ask price, NaN if unknown. */
    double ask_price;
    /** The ask size, NaN if unknown. */
    double ask_size;
    /** The number of the individual ask orders in the event. */
    int64_t ask_count;
} dxfg_market_maker_t;

/**
 * @brief The common fields of the order events: dxfg_order_t, dxfg_spread_order_t and the structures built on them.
 *
 * It is an indexed event: a snapshot of the market depth of a symbol. Its `clazz` is #DXFG_EVENT_ORDER_BASE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/OrderBase.html">Javadoc</a>
 */
typedef struct dxfg_order_base_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_indexed_event_t indexed_event;
    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;

    /**
     * The unique per-symbol index of the order, a non-negative number. The source of the order is stored in its
     * highest bits, dxfg_IndexedEvent_getSource() returns it; see OrderBase.getIndex() in Java.
     */
    int64_t index;
    /**
     * The time of the order and the sequence packed into one number: the time in seconds since the Unix epoch in the
     * higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see
     * OrderBase.getTimeSequence() in Java.
     */
    int64_t time_sequence;
    /** The microseconds and nanoseconds part of the time of the order. */
    int32_t time_nano_part;

    /** The time of the last order action, in milliseconds since the Unix epoch (UTC), 0 if not available. */
    int64_t action_time;
    /** The order ID, 0 if not available (some actions, such as a trade or a bust, have no order ID). */
    int64_t order_id;
    /** The auxiliary order ID (for example, the ID of the order replaced by this new order), 0 if not available. */
    int64_t aux_order_id;

    /** The price of the order, NaN if unknown. */
    double price;
    /** The size of the order (with fractions), NaN if unknown. */
    double size;
    /** The executed size of the order, NaN if unknown. */
    double executed_size;
    /** The number of the individual orders in this aggregate order. */
    int64_t count;
    /**
     * The implementation-specific flags: the scope in bits 0 and 1, the side in bits 2 and 3, the exchange code in
     * bits 4 to 10 and the order action in bits 11 to 14; see the getters of OrderBase in Java (getScope(),
     * getOrderSide(), getExchangeCode(), getAction()).
     */
    int32_t flags;

    /** The trade (order execution) ID for the events that contain a trade-related action, 0 if not available. */
    int64_t trade_id;
    /** The trade price for the events that contain a trade-related action, NaN if unknown. */
    double trade_price;
    /** The trade size for the events that contain a trade-related action, NaN if unknown. */
    double trade_size;
} dxfg_order_base_t;

/**
 * @brief The order event: a snapshot of the market depth (an order or a price level) of a symbol.
 *
 * It is an indexed event, its `clazz` is #DXFG_EVENT_ORDER.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Order.html">Javadoc</a>
 */
typedef struct dxfg_order_t {
    /** The base structure: the fields of dxfg_order_base_t. */
    dxfg_order_base_t order_base;
    /**
     * The market maker or another aggregate identifier of the order (a UTF-8 string, NULL if none) that belongs to the
     * event and is freed with it. It is defined for the aggregate and individual orders.
     */
    const char *market_maker;
} dxfg_order_t;

/**
 * @brief The analytic order event: an order with the iceberg information.
 *
 * It is an indexed event, its `clazz` is #DXFG_EVENT_ANALYTIC_ORDER.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/AnalyticOrder.html">Javadoc</a>
 */
typedef struct dxfg_analytic_order_t {
    /** The base structure: the fields of dxfg_order_t. */
    dxfg_order_t order_base;
    /** The size of the peak: the visible part of the iceberg, NaN if unknown. */
    double iceberg_peak_size;
    /** The hidden size of the iceberg, NaN if unknown. */
    double iceberg_hidden_size;
    /** The executed size of the iceberg, NaN if unknown. */
    double iceberg_executed_size;
    /** The implementation-specific flags that hold the iceberg type; see AnalyticOrder.getIcebergType() in Java. */
    int32_t iceberg_flags;
} dxfg_analytic_order_t;

/**
 * @brief The OTC Markets order event: an order with the OTC Markets specific quote data.
 *
 * It is an indexed event, its `clazz` is #DXFG_EVENT_OTC_MARKETS_ORDER. The original OTC quote is presented as two
 * events: the bid and the ask side.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/OtcMarketsOrder.html">Javadoc</a>
 */
typedef struct dxfg_otc_markets_order_t {
    /** The base structure: the fields of dxfg_order_t. */
    dxfg_order_t order_base;
    /**
     * The Quote Access Payment (QAP): from 1 to 30 is a rebate, from -1 to -30 is an access fee, 0 is neither.
     */
    int32_t quote_access_payment;
    /**
     * The implementation-specific flags: the OTC Markets price type and the open, unsolicited, saturated, auto
     * execution and NMS conditional indicators; see the getters of OtcMarketsOrder in Java (getOtcMarketsPriceType(),
     * isOpen(), isUnsolicited() and so on).
     */
    int32_t otc_markets_flags;
} dxfg_otc_markets_order_t;

/**
 * @brief The spread order event: a snapshot of the market depth of the spreads on an underlying symbol.
 *
 * It is an indexed event, its `clazz` is #DXFG_EVENT_SPREAD_ORDER.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/SpreadOrder.html">Javadoc</a>
 */
typedef struct dxfg_spread_order_t {
    /** The base structure: the fields of dxfg_order_base_t. */
    dxfg_order_base_t order_base;
    /**
     * The symbol of the spread that the order represents (a UTF-8 string, NULL if none) that belongs to the event and
     * is freed with it. The event symbol is the underlying symbol of the subscription.
     */
    const char *spread_symbol;
} dxfg_spread_order_t;

/**
 * @brief The series event: a snapshot of the values computed for an option series of an underlying symbol.
 *
 * It is an indexed event, its `clazz` is #DXFG_EVENT_SERIES.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/Series.html">Javadoc</a>
 */
typedef struct dxfg_series_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    //    dxfg_indexed_event_t indexed_event;

    /** The transactional event flags (TX_PENDING, REMOVE_EVENT, SNAPSHOT_BEGIN and so on), see IndexedEvent in Java. */
    int32_t event_flags;

    /** The unique per-symbol index of the series, see Series.getIndex() in Java. */
    int64_t index;
    /**
     * The time of the series and the sequence packed into one number: the time in seconds since the Unix epoch in the
     * higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see
     * Series.getTimeSequence() in Java.
     */
    int64_t time_sequence;
    /** The day ID of the expiration: the number of days since 1 Jan 1970. */
    int32_t expiration;
    /** The implied volatility index of the series based on the VIX methodology, NaN if unknown. */
    double volatility;
    /** The traded volume of the call options for a day, NaN if unknown. */
    double call_volume;
    /** The traded volume of the put options for a day, NaN if unknown. */
    double put_volume;
    /** The ratio of the traded volume of the put options to the traded volume of the call options, NaN if unknown. */
    double put_call_ratio;
    /** The implied forward price of the series, NaN if unknown. */
    double forward_price;
    /** The implied simple dividend return of the series, NaN if unknown. */
    double dividend;
    /** The implied simple interest return of the series, NaN if unknown. */
    double interest;
} dxfg_series_t;

/**
 * @brief The text message event: a notification with a text.
 *
 * It is neither a lasting nor an indexed event: messages are never conflated, its `clazz` is
 * #DXFG_EVENT_TEXT_MESSAGE.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/misc/TextMessage.html">Javadoc</a>
 */
typedef struct dxfg_text_message_t {
    /** The base structure: the fields of dxfg_event_type_t. */
    dxfg_event_type_t event_type;
    /** The symbol of the event (a UTF-8 null-terminated string) that belongs to the event and is freed with it. */
    const char *event_symbol;
    /** The time when the event was created, in milliseconds since the Unix epoch (UTC), or 0 if not available. */
    int64_t event_time;
    /**
     * The time of the message and the sequence packed into one number: the time in seconds since the Unix epoch in the
     * higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21); see
     * TextMessage.getTimeSequence() in Java.
     */
    int64_t time_sequence;
    /** The text (a UTF-8 string, NULL if none) that belongs to the event and is freed with it. */
    const char *text;
} dxfg_text_message_t;

/**
 * @brief The text configuration event: a text with a version.
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_TEXT_CONFIGURATION.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/misc/TextConfiguration.html">Javadoc</a>
 */
typedef struct dxfg_text_configuration_t {
    /** The base structure: the fields of dxfg_event_type_t. */
    dxfg_event_type_t event_type;
    /** The symbol of the event (a UTF-8 null-terminated string) that belongs to the event and is freed with it. */
    const char *event_symbol;
    /** The time when the event was created, in milliseconds since the Unix epoch (UTC), or 0 if not available. */
    int64_t event_time;
    /**
     * The time of the configuration and the sequence packed into one number: the time in seconds since the Unix epoch
     * in the higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21).
     */
    int64_t time_sequence;
    /** The version of the event: a higher value is a more recent update. */
    int32_t version;
    /** The text (a UTF-8 string, NULL if none) that belongs to the event and is freed with it. */
    const char *text;
} dxfg_text_configuration_t;

/**
 * @brief The Nuam order event: an order with the data of the Nuam Exchange matching engine.
 *
 * It is an indexed event, its `clazz` is #DXFG_EVENT_NUAM_ORDER.
 *
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/custom/NuamOrder.html
 */
typedef struct dxfg_nuam_order_t {
    /** The base structure: the fields of dxfg_order_t. */
    dxfg_order_t order_base;
    /** The actor owning the order: the one that entered it or the owner appointed for an order entered on behalf. */
    int32_t actor_id;
    /** The participant owning the order. */
    int32_t participant_id;
    /** The actor submitting the message that changed the order. */
    int32_t submitter_id;
    /** The actor submitting an order message on behalf of someone else. */
    int32_t on_behalf_of_submitter_id;
    /** The client generated order identity (ClOrdID of the FIX protocol), a UTF-8 string that belongs to the event. */
    const char *client_order_id;
    /** The account to use for the order, a UTF-8 string that belongs to the event. */
    const char *customer_account;
    /** The free text filled in by the client entering the transaction, a UTF-8 string that belongs to the event. */
    const char *customer_info;
    /** The exchange-specific information (may contain unprintable characters), a UTF-8 string of the event. */
    const char *exchange_info;
    /** The additional information that defines the validity period of the order; see NuamOrder.getTimeInForceData(). */
    int32_t time_in_force_data;
    /** The order book ID the order triggers on (only for the trigger on price orders). */
    int32_t trigger_order_book_id;
    /** The trigger price of the order, NaN if unknown. */
    double trigger_price;
    /** The session type to trigger on. */
    int32_t trigger_session_type;
    /** The actual full order quantity (it does not change when the order is matched), NaN if unknown. */
    double order_quantity;
    /** The actual display quantity of a reserve order at the end of the transaction, NaN if not a reserve order. */
    double display_quantity;
    /** The initial display quantity of a reserve order, NaN if not a reserve order. */
    double refresh_quantity;
    /** The remaining (open) quantity of the order, NaN if unknown. */
    double leaves_quantity;
    /** The quantity matched within the transaction that updated the order, NaN if unknown. */
    double matched_quantity;
    /** The implementation-specific flags that hold the order type and the time in force type; see NuamOrder in Java. */
    int32_t nuam_flags;
} dxfg_nuam_order_t;

/**
 * @brief The Nuam time and sale event: a time and sale with the data of the Nuam Exchange.
 *
 * It is a time series event, its `clazz` is #DXFG_EVENT_NUAM_TIME_AND_SALE.
 *
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/custom/NuamTimeAndSale.html
 */
typedef struct dxfg_nuam_time_and_sale_t {
    /** The base structure: the fields of dxfg_time_and_sale_t. */
    dxfg_time_and_sale_t time_and_sale_base;
    /** The ID of the match between the buyer and the seller, available to the participants of the deal. */
    int64_t match_id;
} dxfg_nuam_time_and_sale_t;

/**
 * @brief The Nuam trade event: a trade with the trade statistics of the Nuam Exchange.
 *
 * It is a lasting event, its `clazz` is #DXFG_EVENT_NUAM_TRADE.
 *
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/custom/NuamTrade.html
 */
typedef struct dxfg_nuam_trade_t {
    /** The base structure: the fields of dxfg_trade_t. */
    dxfg_trade_t trade_base;
    /** The time when the trade statistics were updated, in milliseconds since the Unix epoch (UTC). */
    int64_t trade_stat_time;
    /** The price of the last trade reported by the exchange as significant, NaN if unknown. */
    double last_significant_price;
    /** The price of the last trade (for all deals) reported by the exchange, NaN if unknown. */
    double last_price_for_all;
    /** The number of the trades reported by the exchange. */
    int32_t number_of_trades;
    /** The volume-weighted average price (VWAP) reported by the exchange, NaN if unknown. */
    double vwap;
} dxfg_nuam_trade_t;

/**
 * @brief The order imbalance event: the order book statistics during an auction.
 *
 * It is a market event (neither lasting nor indexed), its `clazz` is #DXFG_EVENT_ORDER_IMBALANCE.
 *
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/OrderImbalance.html
 */
typedef struct dxfg_order_imbalance_t {
    /** The base structure: the fields of dxfg_market_event_t. */
    dxfg_market_event_t market_event;
    /** The identifier of the source of the event (OrderImbalance.getSource() in Java). */
    int32_t source_id;
    /**
     * The time of the event and the sequence packed into one number: the time in seconds since the Unix epoch in the
     * higher 32 bits, then the milliseconds part of the time (bits 22 to 31) and the sequence (bits 0 to 21).
     */
    int64_t time_sequence;
    /** The reference price: the price for which the paired and imbalance sizes are calculated, NaN if unknown. */
    double ref_price;
    /** The paired size: the number of shares that can be matched at the reference price, NaN if unknown. */
    double paired_size;
    /** The imbalance size: the number of shares that are not matched at the reference price, NaN if unknown. */
    double imbalance_size;
    /** The near price: the hypothetical price at which all the interest on the book could trade now, NaN if unknown. */
    double near_price;
    /** The far price: the hypothetical price at which only the auction interest could trade now, NaN if unknown. */
    double far_price;
    /**
     * The implementation-specific flags: the imbalance side in bits 0 to 2 and the auction type in bits 3 to 5; see
     * OrderImbalance.getImbalanceSide() and OrderImbalance.getAuctionType() in Java.
     */
    int32_t flags;
} dxfg_order_imbalance_t;

/**
 * @brief A list of events.
 *
 * The list returned by the SDK (or passed by the SDK to a listener) is released with
 * dxfg_CList_EventType_release().
 */
typedef struct dxfg_event_type_list {
    /** The number of the elements. */
    int32_t size;
    /** The array of `size` pointers to the events; NULL when `size` is 0 in a list created by the SDK. */
    dxfg_event_type_t **elements;
} dxfg_event_type_list;

/**
 * @brief A list of event classes.
 *
 * The list returned by the SDK is released with dxfg_CList_EventClazz_release().
 */
typedef struct dxfg_event_clazz_list_t {
    /** The number of the elements. */
    int32_t size;
    /** The array of `size` pointers to the event classes; NULL when `size` is 0 in a list created by the SDK. */
    dxfg_event_clazz_t **elements;
} dxfg_event_clazz_list_t;

/**
 * Creates a symbol of the given type.
 *
 * Only #STRING, #CANDLE and #WILDCARD can be created here: for any other type the function fails.
 * For #CANDLE the string is parsed as a candle symbol (an invalid string is an error) and the symbol of the created
 * structure is its normalized string form. For #WILDCARD the string is ignored.
 *
 * @deprecated Use dxfg_v4_Symbol_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbol The symbol: a UTF-8 null-terminated string.
 * @param[in] symbolType The type of the symbol to create.
 * @return The symbol (release it with dxfg_Symbol_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `CandleSymbol.valueOf(String)` for #CANDLE, `WildcardSymbol.ALL` for #WILDCARD.
 */
DXFG_DEPRECATED("use dxfg_v4_Symbol_new()")
dxfg_symbol_t *dxfg_Symbol_new(graal_isolatethread_t *thread, const char *symbol, dxfg_symbol_type_t symbolType);

/**
 * Releases a symbol: frees the structure, its strings and the nested symbols and sources. Passing NULL does nothing.
 *
 * @deprecated Use dxfg_v4_Symbol_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbol The symbol created by dxfg_Symbol_new() or returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Symbol_release()")
int32_t dxfg_Symbol_release(graal_isolatethread_t *thread, dxfg_symbol_t *symbol);

/**
 * Creates an event of the given class for the given symbol.
 *
 * The event has the given `clazz` and event symbol (a copy of the string), all the other fields are zero (not NaN
 * and not the defaults of the Java event). Fill the fields, then pass the event to a function that needs it.
 *
 * @deprecated Use dxfg_v4_EventType_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbolName The event symbol: a UTF-8 null-terminated string. For #DXFG_EVENT_CANDLE and
 * #DXFG_EVENT_DAILY_CANDLE it is the string form of a candle symbol.
 * @param[in] clazz The class of the event to create.
 * @return The event (release it with dxfg_EventType_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_EventType_new()")
dxfg_event_type_t *dxfg_EventType_new(graal_isolatethread_t *thread, const char *symbolName, dxfg_event_clazz_t clazz);

/**
 * Releases an event: frees the structure together with all its strings. Passing NULL does nothing.
 *
 * Use it for the events created by dxfg_EventType_new() and for the events that the SDK returns to the caller as a
 * result. Do not release an event that the SDK passed to a listener: the SDK frees it after the listener returns.
 *
 * @deprecated Use dxfg_v4_EventType_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The event to release.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_EventType_release()")
int32_t dxfg_EventType_release(graal_isolatethread_t *thread, dxfg_event_type_t *eventType);

/**
 * Releases a list of events: frees the list, its array and every event of the list together with its strings.
 * Passing NULL does nothing.
 *
 * @deprecated Use dxfg_v4_event_type_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventTypes The list of events returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_event_type_list_release()")
int32_t dxfg_CList_EventType_release(
    graal_isolatethread_t *thread,
    dxfg_event_type_list *eventTypes);

/**
 * Releases a list of event classes: frees the list, its array and every element. Passing NULL does nothing.
 *
 * @deprecated Use dxfg_v4_event_clazz_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventClazzes The list of event classes returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_event_clazz_list_release()")
int32_t dxfg_CList_EventClazz_release(graal_isolatethread_t *thread,
                                      dxfg_event_clazz_list_t *eventClazzes);

/**
 * Releases a list of symbols: frees the list, its array and every symbol of the list with its strings and the nested
 * symbols and sources.
 *
 * @deprecated Use dxfg_v4_symbol_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbols The list of symbols returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_symbol_list_release()")
int32_t dxfg_CList_symbol_release(
    graal_isolatethread_t *thread,
    dxfg_symbol_list *symbols);

/**
 * Returns the source of an indexed event.
 *
 * It works for the events of #DXFG_EVENT_GREEKS, #DXFG_EVENT_CANDLE, #DXFG_EVENT_DAILY_CANDLE,
 * #DXFG_EVENT_TIME_AND_SALE, #DXFG_EVENT_UNDERLYING, #DXFG_EVENT_THEO_PRICE and #DXFG_EVENT_SERIES (the default
 * source) and of #DXFG_EVENT_ORDER_BASE, #DXFG_EVENT_ORDER, #DXFG_EVENT_ANALYTIC_ORDER and #DXFG_EVENT_SPREAD_ORDER
 * (the source stored in the index of the order). For the events of any other class the function fails.
 *
 * @deprecated Use dxfg_v4_IndexedEvent_getSource().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] eventType The event.
 * @return The source (release it with dxfg_IndexedEventSource_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `IndexedEvent.getSource()`.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEvent_getSource()")
dxfg_indexed_event_source_t *dxfg_IndexedEvent_getSource(graal_isolatethread_t *thread, dxfg_event_type_t *eventType);

/**
 * Creates an event source by its name.
 *
 * @deprecated Use dxfg_v4_IndexedEventSource_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The name of the order source: a UTF-8 null-terminated string, either predefined or at most four
 * alphanumeric characters. If it is NULL, the default source (the identifier 0) is returned instead.
 * @return The source (release it with dxfg_IndexedEventSource_release()), or NULL on error (for example, for an invalid
 * name). Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `OrderSource.valueOf(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventSource_new()")
dxfg_indexed_event_source_t *dxfg_IndexedEventSource_new(
    graal_isolatethread_t *thread,
    const char *source);

/**
 * Creates an event source by its identifier.
 *
 * @deprecated Use dxfg_v4_IndexedEventSource_newWithId().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] sourceId The identifier of the source. If it is 0, the default source is returned, otherwise it is the
 * identifier of an order source (a negative identifier is an error).
 * @param[out] source The pointer to the source (release it with dxfg_IndexedEventSource_release()). It is not written
 * on error.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `OrderSource.valueOf(int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventSource_newWithId()")
int32_t dxfg_IndexedEventSource_new2(
    graal_isolatethread_t *thread,
    int32_t sourceId, DXFG_OUT dxfg_indexed_event_source_t** source);

/**
 * Releases an event source: frees the structure and its name. Passing NULL does nothing.
 *
 * @deprecated Use dxfg_v4_IndexedEventSource_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The source created by dxfg_IndexedEventSource_new(), dxfg_IndexedEventSource_new2() or returned by
 * the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_IndexedEventSource_release()")
int32_t dxfg_IndexedEventSource_release(graal_isolatethread_t *thread, dxfg_indexed_event_source_t *source);

/**
 * Releases a list of event sources: frees the list, its array and every source of the list with its name.
 * Passing NULL does nothing.
 *
 * @deprecated Use dxfg_v4_indexed_event_source_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list of sources returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_indexed_event_source_list_release()")
int32_t dxfg_CList_IndexedEventSource_release(graal_isolatethread_t *thread, dxfg_indexed_event_source_list* list);

/**
 * @brief A handle of a Java ObservableSubscriptionChangeListener object; release it with
 * dxfg_JavaObjectHandler_release().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/ObservableSubscriptionChangeListener.html">
 * Javadoc</a>
 */
typedef struct dxfg_observable_subscription_change_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_observable_subscription_change_listener_t;

/**
 * The callback of a subscription change listener: symbols were added to the subscription.
 *
 * It is called by the SDK after the symbols are added to the subscription of the Java side, on a Java thread (not
 * necessarily the thread that registered the listener), while the SDK holds the lock of the subscription: the
 * notifications of one subscription do not run at the same time. The symbols are decorated depending on the
 * subscription: for a time series subscription they are dxfg_time_series_subscription_symbol_t.
 *
 * The list belongs to the SDK and is freed after the function returns: do not release it, copy the symbols that you
 * need to keep.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbols The added symbols.
 * @param[in] user_data The user data passed to dxfg_ObservableSubscriptionChangeListener_new().
 */
typedef void (*dxfg_ObservableSubscriptionChangeListener_function_symbolsAdded)(graal_isolatethread_t *thread, dxfg_symbol_list *symbols, void *user_data);

/**
 * The callback of a subscription change listener: symbols were removed from the subscription.
 *
 * It is called in the same way as dxfg_ObservableSubscriptionChangeListener_function_symbolsAdded, after the symbols
 * are removed from the subscription.
 *
 * The list belongs to the SDK and is freed after the function returns: do not release it, copy the symbols that you
 * need to keep.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbols The removed symbols.
 * @param[in] user_data The user data passed to dxfg_ObservableSubscriptionChangeListener_new().
 */
typedef void (*dxfg_ObservableSubscriptionChangeListener_function_symbolsRemoved)(graal_isolatethread_t *thread, dxfg_symbol_list *symbols, void *user_data);

/**
 * The callback of a subscription change listener: the subscription was closed, or the listener was removed from it.
 *
 * It is called in the same way as dxfg_ObservableSubscriptionChangeListener_function_symbolsAdded.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_data The user data passed to dxfg_ObservableSubscriptionChangeListener_new().
 */
typedef void (*dxfg_ObservableSubscriptionChangeListener_function_subscriptionClosed)(graal_isolatethread_t *thread, void *user_data);

/**
 * Creates a listener of the changes of a subscription.
 *
 * Add the listener to a subscription (for example, with dxfg_DXFeedSubscription_addChangeListener()). The SDK calls
 * the functions without a check for NULL, so pass all three. Releasing the last handle of the listener with
 * dxfg_JavaObjectHandler_release() (or dxfg_NativeListener_deactivate()) stops the calls: after that none of the
 * functions is called or running, and the user data can be freed.
 *
 * @deprecated Use dxfg_v4_ObservableSubscriptionChangeListener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] function_symbolsAdded The function called when symbols are added to the subscription.
 * @param[in] function_symbolsRemoved The function called when symbols are removed from the subscription.
 * @param[in] function_subscriptionClosed The function called when the subscription is closed or the listener is
 * removed.
 * @param[in] user_data The user data passed to every function; the SDK does not use or free it.
 * @return The listener handle (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `ObservableSubscriptionChangeListener` (symbolsAdded(Set), symbolsRemoved(Set), subscriptionClosed()).
 */
DXFG_DEPRECATED("use dxfg_v4_ObservableSubscriptionChangeListener_new()")
dxfg_observable_subscription_change_listener_t* dxfg_ObservableSubscriptionChangeListener_new(
    graal_isolatethread_t *thread,
    dxfg_ObservableSubscriptionChangeListener_function_symbolsAdded function_symbolsAdded,
    dxfg_ObservableSubscriptionChangeListener_function_symbolsRemoved function_symbolsRemoved,
    dxfg_ObservableSubscriptionChangeListener_function_subscriptionClosed function_subscriptionClosed, void *user_data);

#ifdef __cplusplus
}
#endif

#endif // DXFG_EVENTS_H
