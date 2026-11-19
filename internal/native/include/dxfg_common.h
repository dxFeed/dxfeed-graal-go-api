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
 *  The experimental functions are the twins of the Java API that is experimental (`@Experimental` in QD, and the
 *  transaction models taken from the experimental Java API), and the functions of the advanced headers: the
 *  `dxfg_v4_*.h` headers that dxfg_v4_api.h does not include, each one included on its own. The twins came out in
 *  3.12.0, before the tiers, and are experimental since 3.13.0; their older functions stay until 4.0.0, as the other
 *  deprecated functions do, and their deprecation messages say that the replacement is experimental. An advanced header
 *  may include the basic headers, never the other way round. The experimental event types say so in their descriptions;
 *  their structures are delivered as the other events are. Every declaration added after 3.12.0 names the release that
 *  added it (`@since`).
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
