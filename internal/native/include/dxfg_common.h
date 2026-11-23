// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Common declarations
 */

#ifndef DXFG_COMMON_H
#define DXFG_COMMON_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "graal_isolate.h"

/** @defgroup Common Common
 *  @brief The macros and the functions that every header of the C API uses: the marker of an output parameter, the
 *  deprecation marker, the status results and the functions that free memory allocated by the SDK.
 *
 *  The ownership rules of the handles, the strings and the lists are described in the group @ref Javac.
 *  @{
 */

/**
 * Marks an output parameter of a function: a pointer to a variable that the function writes the result to, e.g.
 * `DXFG_OUT dxfg_endpoint_t **endpoint`. The caller passes the address of its own variable, and the variable is
 * written only if the function succeeds (returns #DXFG_EXECUTE_SUCCESSFULLY). The macro expands to nothing: it only
 * documents the direction of the parameter in the declaration (the `[out]` of the documentation).
 */
#ifndef DXFG_OUT
#    ifndef DOXYGEN_SHOULD_SKIP_THIS
#        define DXFG_OUT
#    endif // DOXYGEN_SHOULD_SKIP_THIS
#endif     /* DXFG_OUT */

/**
 * Marks a deprecated function: the compilers warn where it is used. The function still works, see its description for
 * the replacement.
 */
#ifndef DXFG_DEPRECATED
#    if defined(__cplusplus) && __cplusplus >= 201402L
#        define DXFG_DEPRECATED(message) [[deprecated(message)]]
#    elif defined(__GNUC__) || defined(__clang__)
#        define DXFG_DEPRECATED(message) __attribute__((deprecated(message)))
#    elif defined(_MSC_VER)
#        define DXFG_DEPRECATED(message) __declspec(deprecated(message))
#    else
#        define DXFG_DEPRECATED(message)
#    endif
#endif /* DXFG_DEPRECATED */

/**
 * Marks an experimental function: it may change or be removed in a minor release (see the tiers in
 * @ref V4_Conventions), and its description says so. The macro expands to nothing, unless the application defines it
 * before the headers are included, e.g. `#define DXFG_EXPERIMENTAL [[deprecated("experimental")]]` in C++ or
 * `-D'DXFG_EXPERIMENTAL=__attribute__((deprecated("experimental")))'` with GCC and Clang, so that the compilers warn
 * where the experimental functions are used.
 */
#ifndef DXFG_EXPERIMENTAL
#    define DXFG_EXPERIMENTAL
#endif /* DXFG_EXPERIMENTAL */

/**
 * A successful result that any dxFeed Graal Native SDK function can return: the status of the functions that return
 * `int32_t` status, #DXFG_EXECUTE_SUCCESSFULLY (0) means that the function was executed without an exception.
 */
#define DXFG_EXECUTE_SUCCESSFULLY ((int32_t)0)

/**
 * A result that signals an error that any dxFeed Graal Native SDK function can return: a function that returns
 * `int32_t` status returns #DXFG_EXECUTE_FAIL (-1) when it was terminated by an exception. For information about the
 * thrown exception, see the dxfg_get_and_clear_thread_exception_t() function. Other functions signal an error with
 * another value, which is given in the `@return` of the function.
 */
#define DXFG_EXECUTE_FAIL ((int32_t)-1)

/**
 * @brief A map entry that consists of a string and a double: an element of an array that represents a map from
 * strings to doubles.
 *
 * An array that a function returns (e.g. dxfg_AdditionalUnderlyings_getMap()) is allocated by the SDK together with
 * the strings of the keys: release it with dxfg_free_string_to_double_map_entries(). An array that is passed to a
 * function (e.g. dxfg_AdditionalUnderlyings_valueOf2()) is only read by the SDK: the caller keeps its ownership.
 */
typedef struct dxfg_string_to_double_map_entry_t {
    /** The key: a NUL-terminated UTF-8 string. It belongs to the entry and is freed with the array. */
    const char *key;
    /** The value that corresponds to the key. */
    double value;
} dxfg_string_to_double_map_entry_t;

/**
 * @brief Frees an array of string to double map entries that a function of the SDK returned.
 *
 * It frees the key of every entry and the array itself. Do not use the array and the keys after that. If the array
 * is NULL or `size` is not positive, the function does nothing (an empty array is not freed) and returns
 * #DXFG_EXECUTE_SUCCESSFULLY.
 *
 * @deprecated Use dxfg_v4_string_to_double_map_entry_list_release(): dxfg_v4_AdditionalUnderlyings_getMap() returns
 * a list.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] mapEntries The array of map entries ([string, double] pair) that was returned by the SDK.
 * @param[in] size The number of entries in the array, as returned together with the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_string_to_double_map_entry_list_release()")
int32_t dxfg_free_string_to_double_map_entries(graal_isolatethread_t *thread,
                                               dxfg_string_to_double_map_entry_t *mapEntries, int32_t size);

/**
 * @brief Frees a block of the unmanaged (GraalVM) memory that a function of the SDK allocated, e.g. the array of
 * doubles that dxfg_PriceIncrements_getPriceIncrements() returns.
 *
 * Only a block of memory that a function of the SDK allocated and documents to be released with this function can be
 * passed. Do not use the memory after that.
 *
 * @deprecated Use dxfg_v4_free().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] pointer The pointer to the unmanaged memory that was returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_free()")
int32_t dxfg_free(graal_isolatethread_t *thread, void *pointer);

/**
 * @brief Frees an array of C strings that a function of the SDK returned.
 *
 * It frees every string and the array itself. Do not use the array and the strings after that. If the array is NULL
 * or `size` is not positive, the function does nothing (an empty array is not freed) and returns
 * #DXFG_EXECUTE_SUCCESSFULLY. A list of strings (dxfg_string_list) is not an array of this kind: it is released with
 * dxfg_CList_String_release().
 *
 * @deprecated Use dxfg_v4_string_list_release(): the v4 functions return strings in lists.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] strings The array of the NUL-terminated strings that was returned by the SDK.
 * @param[in] size The number of strings in the array, as returned together with the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_string_list_release()")
int32_t dxfg_free_strings(graal_isolatethread_t *thread, char** strings, int32_t size);

/** @} */ // end of Common

/** @defgroup V4_Conventions The v4 functions
 *  @brief The `dxfg_v4_` functions: the operations of the older functions of the C API in one uniform form.
 *
 *  A function of the C API that has a `dxfg_v4_` twin is deprecated in favour of it. The twin is declared in a
 *  `dxfg_v4_<area>.h` header next to the header of the older function (dxfg_v4_api.h and dxfg_api.h include all of
 *  them). dxfg_free_strings() and dxfg_free_string_to_double_map_entries() have no twin and are deprecated: the v4
 *  functions return lists, released by their `dxfg_v4_<element>_list_release()` functions. The deprecated functions
 *  are removed in 4.0.0, the next major release; their deprecation warnings name the functions to move to. A twin
 *  does the same as the older function, only its form differs:
 *  - It returns `int32_t`: #DXFG_EXECUTE_SUCCESSFULLY (0) on success, #DXFG_EXECUTE_FAIL (-1) on error. The status is
 *    the only error signal: after -1, dxfg_v4_get_and_clear_thread_exception() gives the exception. A successful call
 *    does not clear an exception left by an earlier failed call that was not taken, so take it right after the error.
 *    A function that waits may also return documented positive codes for outcomes that are not errors (e.g.
 *    #DXFG_PROMISE_AWAIT_TIMED_OUT, #DXFG_EVENT_QUEUE_NO_EVENTS).
 *  - It writes its result to the last parameter, marked #DXFG_OUT: the caller passes the address of its own variable.
 *    The variable is written only if the function succeeds (dxfg_v4_DXFeed_getLastEvent() and
 *    dxfg_v4_DXFeed_getLastEvents() write NULL on error, as their older functions do); after an error its value must
 *    not be used. A NULL out pointer is an error (`IllegalArgumentException`), unless the parameter is documented as
 *    optional.
 *  - A NULL handle, string or list written to the out variable means "no value" (e.g. no session is found), it is
 *    never an error. A number is written as it is: -1 or 0 are ordinary values, not error markers.
 *  - It returns a collection as a list `dxfg_v4_<element>_list` (`size` and `elements`; `elements` of an empty list
 *    may be NULL), which is released with its `dxfg_v4_<element>_list_release()` function. A list that the caller
 *    builds and passes to a function stays the caller's: it is never passed to a release function.
 *  - Its name has no digit suffix: an overload is named by what it takes (e.g. dxfg_v4_TimePeriod_valueOfString()).
 *  - The types of its parameters have `dxfg_v4_` names: typedefs of the older types with the same layout. The
 *    constants of the enumerations keep their names.
 *
 *  Migrate a call of an older function to the function that its deprecation message names, which is also the
 *  `replacements` of the function in `dxfg-api.json`; never by the name. Since an overload is named by what it takes,
 *  the same name may stand for another overload in the two generations, and a migration by the name compiles, or binds,
 *  and calls another function. Without the digit suffix, the `dxfg_v4_` name of an older function may be its
 *  replacement (dxfg_v4_AdditionalUnderlyings_getSPC() replaces dxfg_AdditionalUnderlyings_getSPC2()) or another
 *  function: dxfg_DXEndpoint_getInstance2() is replaced by dxfg_v4_DXEndpoint_createWithRole(), while
 *  dxfg_v4_DXEndpoint_getInstance() replaces dxfg_DXEndpoint_getInstance(). These older functions without a digit
 *  suffix have a `dxfg_v4_` function of the same name that is another overload, each listed with its replacement and
 *  the older function that the `dxfg_v4_` function of the same name replaces (the build checks the list against
 *  `dxfg-api.json`):
 *  <!-- dxfg-renamed:begin -->
 *  - dxfg_AdditionalUnderlyings_getSPC() -> dxfg_v4_AdditionalUnderlyings_getSPCFromText();
 *    dxfg_v4_AdditionalUnderlyings_getSPC() replaces dxfg_AdditionalUnderlyings_getSPC2().
 *  - dxfg_InstrumentProfileCollector_updateInstrumentProfile() ->
 *    dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandle();
 *    dxfg_v4_InstrumentProfileCollector_updateInstrumentProfile() replaces
 *    dxfg_InstrumentProfileCollector_updateInstrumentProfile2().
 *  - dxfg_InstrumentProfileCollector_updateInstrumentProfiles() ->
 *    dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandles();
 *    dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles() replaces
 *    dxfg_InstrumentProfileCollector_updateInstrumentProfiles2() and
 *    dxfg_InstrumentProfileCollector_updateInstrumentProfiles3().
 *  - dxfg_InstrumentProfileReader_read() -> dxfg_v4_InstrumentProfileReader_readAsHandles();
 *    dxfg_v4_InstrumentProfileReader_read() replaces dxfg_InstrumentProfileReader_read3() and
 *    dxfg_InstrumentProfileReader_read5().
 *  - dxfg_InstrumentProfileReader_readCompressed() -> dxfg_v4_InstrumentProfileReader_readCompressedAsHandles();
 *    dxfg_v4_InstrumentProfileReader_readCompressed() replaces dxfg_InstrumentProfileReader_readCompressed2() and
 *    dxfg_InstrumentProfileReader_readCompressed3().
 *  - dxfg_InstrumentProfileReader_readFromFile() -> dxfg_v4_InstrumentProfileReader_readFromFileAsHandles();
 *    dxfg_v4_InstrumentProfileReader_readFromFile() replaces dxfg_InstrumentProfileReader_readFromFile4() and
 *    dxfg_InstrumentProfileReader_readFromFile7().
 *  - dxfg_Iterable_InstrumentProfile_next() -> dxfg_v4_Iterable_InstrumentProfile_nextHandle();
 *    dxfg_v4_Iterable_InstrumentProfile_next() replaces dxfg_Iterable_InstrumentProfile_next2().
 *  - dxfg_PriceLevelService_getOrders() -> dxfg_v4_PriceLevelService_getOrdersWithCaller();
 *    dxfg_v4_PriceLevelService_getOrders() replaces dxfg_PriceLevelService_getOrders2().
 *  - dxfg_PriceLevelService_getQuotes() -> dxfg_v4_PriceLevelService_getQuotesWithCaller();
 *    dxfg_v4_PriceLevelService_getQuotes() replaces dxfg_PriceLevelService_getQuotes2().
 *  - dxfg_Schedule_getInstance() -> dxfg_v4_Schedule_getInstanceByProfileHandle();
 *    dxfg_v4_Schedule_getInstance() replaces dxfg_Schedule_getInstance4().
 *  - dxfg_Schedule_getTradingVenues() -> dxfg_v4_Schedule_getTradingVenuesByProfileHandle();
 *    dxfg_v4_Schedule_getTradingVenues() replaces dxfg_Schedule_getTradingVenues2().
 *  <!-- dxfg-renamed:end -->
 *
 *  The replacements of a few older functions are alternatives of another shape rather than twins, so no `dxfg_v4_`
 *  function has their names:
 *  <!-- dxfg-alternatives:begin -->
 *  - dxfg_DXEndpoint_executor() -> dxfg_v4_DXEndpoint_Builder_withExecutor(): the executor is set on the builder,
 *    before the endpoint is built, not on an endpoint.
 *  - dxfg_InstrumentProfileCustomFields_addNonEmptyFieldNames() ->
 *    dxfg_v4_InstrumentProfileCustomFields_getNonEmptyFieldNames(): the same list, without the `updated` flag.
 *  - dxfg_free_strings() -> dxfg_v4_string_list_release(): the v4 functions return lists instead of arrays and sizes.
 *  - dxfg_free_string_to_double_map_entries() -> dxfg_v4_string_to_double_map_entry_list_release(): the same, for the
 *    entries of a map.
 *  <!-- dxfg-alternatives:end -->
 *
 *  The `dxfg_v4_` prefix names a generation of the C API, not a release of the SDK: the rules above hold for every
 *  function and type with this prefix, and later releases, major ones included, add functions with the same prefix.
 *  A new generation appears only when these rules change. It is named after the major release in which it becomes
 *  the only one: the v4 functions become the only ones in 4.0.0.
 *
 *  The v4 functions have two tiers:
 *  - A stable function keeps its name, parameters and behaviour. It changes only by deprecation (a new function
 *    replaces it, and the old one is marked `DXFG_DEPRECATED`), and it is removed only in a major release.
 *  - An experimental function is marked with `DXFG_EXPERIMENTAL` and listed on the page @ref experimental. It may
 *    change or be removed in any minor release, also between the release candidates of one release. An incompatible
 *    change gives it a new name (named by what it takes, with no digit suffix) and increments the advanced revision
 *    (dxfg_v4_system_get_advanced_revision(), #DXFG_V4_ADVANCED_REVISION), so a wrapper bound by name misses the old
 *    symbol instead of calling it with a changed signature. A compatible addition adds a feature
 *    (dxfg_v4_system_is_feature_supported()) and keeps the revision. A function promoted to stable keeps its name.
 *
 *  The experimental functions are the ones whose shape is not settled yet (the feed filters, the ready signal of an
 *  event queue, the event peek of the time series model, the listener of the message connectors and the last connection
 *  errors), and the functions of the advanced headers: the `dxfg_v4_*.h` headers that dxfg_v4_api.h does not include,
 *  each one included on its own. An advanced header may include the basic headers, never the other way round. Every
 *  declaration added after 3.12.0 names the release that added it (`@since`).
 *
 *  A stable function may call a Java member that is `@Experimental` in QD (`IndexedEventTxModel`,
 *  `SubscriptionController`, `DXFeedSubscription.getSubscriptionController()` and `getTotalDroppedEvents()`, the Nuam
 *  and `OrderImbalance` events). A change of the Java member does not reach the C function. If the member is removed
 *  without a replacement, the function stays: its body is replaced by one that fails with
 *  `UnsupportedOperationException`, and the function is deprecated and removed only in a major release. The layout of
 *  a stable event structure does not change. These twins and the ones of the transaction models (`IndexedTxModel`,
 *  `TimeSeriesTxModel`, part of the SDK) were experimental from 3.13.0 to 3.18.0 and are stable now.
 *
 *  The older and the v4 functions can be mixed on the same objects: a handle, a string or a list that one of them
 *  returned can be passed to the other one and released by the release function of either.
 *
 *  The older functions that have twins are marked with `DXFG_DEPRECATED`, so the compilers warn where they are used.
 *  To compile code that still uses them without the warnings, define the macro empty before the headers are included:
 *  `#define DXFG_DEPRECATED(message)` in the code, or `-DDXFG_DEPRECATED(message)=` on the command line of GCC and
 *  Clang. With MSVC, the warning is C4996 (`/wd4996` disables it).
 *
 *  A wrapper for another language builds on these rules: @ref wrapper_guide states what it may rely on and must do
 *  (threads, ownership, outcomes, versioning, delivery), and the conformance kit that checks it.
 */

#ifdef __cplusplus
}
#endif

#endif // DXFG_COMMON_H
