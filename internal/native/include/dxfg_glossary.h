// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Glossary functions and types declarations
 */

#ifndef DXFG_GLOSSARY_H
#define DXFG_GLOSSARY_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_javac.h"
#include "graal_isolate.h"

/**
 * @defgroup Glossary Glossary
 */

/**
 * @defgroup AdditionalUnderlyings AdditionalUnderlyings
 * @ingroup Glossary
 * @{
 */

/**
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html
 *
 * @see AdditionalUnderlyingsFunctions
 *
 * Use dxfg_Object_toString() to retrieve a string representation.
 *
 * Use dxfg_Object_equals() to compare for equality.
 *
 * Use dxfg_Object_hashCode() to get the hash code.
 *
 * @see Javac
 */
typedef struct dxfg_additional_underlyings_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_additional_underlyings_t;

/// @} end of AdditionalUnderlyings

/**
 * @defgroup CFI CFI
 * @ingroup Glossary
 * @{
 */

/**
 * A handle of a wrapper class for Classification of Financial Instruments code as defined in ISO 10962 standard.
 * Main purpose is to provide code validity checks and to construct textual explanation of CFI code as defined in the
 * standard via API. This handle does not provide API-accessible constants for specific instrument attributes and
 * values.
 *
 * https://en.wikipedia.org/wiki/ISO_10962
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html
 *
 * @see CFIFunctions
 */
typedef struct dxfg_cfi_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_cfi_t;

/**
 * @defgroup CFIAttribute CFI.Attribute
 * @{
 */

/**
 * Describes single attribute with all values as defined in the ISO 10962 standard.
 *
 * https://en.wikipedia.org/wiki/ISO_10962
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Attribute.html
 *
 * @see CFIAttributeFunctions
 */
typedef struct dxfg_cfi_attribute_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_cfi_attribute_t;

/// @} end of CFIAttribute

/**
 * @defgroup CFIValue CFI.Value
 * @{
 */

/**
 * Describes single value of single character of CFI code as defined in the ISO 10962 standard.
 *
 * https://en.wikipedia.org/wiki/ISO_10962
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html
 *
 * @see CFIValueFunctions
 */
typedef struct dxfg_cfi_value_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_cfi_value_t;

/// @} end of CFIValue

/// @} end of CFI

/**
 * @defgroup PriceIncrements PriceIncrements
 * @ingroup Glossary
 * @{
 */

/**
 * Represents rules for valid price quantization for a given instrument on a certain exchange.
 * These rules are defined as a set of price ranges with associated price increments.
 * Each price increment defines what price values are valid for corresponding price range -
 * valid prices shall be divisible by corresponding price increment.
 * <p>
 * All price ranges shall be mutually exclusive and they shall cover entire space from 0 to infinity.
 * Therefore all ranges can be represented as a sequence of numbers where increments are interleaved
 * with range limits, with extreme limits (0 and infinity) omitted for short.
 * Negative space (from negative infinity to 0) uses symmetrical price ranges.
 * <p>
 * There is a special value @ref dxfg_PriceIncrements_EMPTY() "EMPTY" which is used to represent unknown or undefined
 * rules. This value has empty textual representation and is considered to have sole increment with value 0. <p> See
 * dxfg_PriceIncrements_getText() and dxfg_PriceIncrements_getPriceIncrements() for details about used formats and
 * representations.
 *
 * <p><b>NOTE:</b>This class produces precise results for decimal numbers with at most 14
 * significant digits and at most 14 digits after decimal point.
 *
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html
 *
 * @see PriceIncrementsFunctions
 */
typedef struct dxfg_price_increments_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_price_increments_t;

/// @} end of PriceIncrements

/**
 * @defgroup AdditionalUnderlyingsFunctions AdditionalUnderlyings' functions
 * @ingroup AdditionalUnderlyings
 * @{
 */

/**
 * Returns an empty additional underlyings
 * (https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#EMPTY it has empty text and empty
 * map). The handle of this instance will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * dxfg_additional_underlyings_t* additional_underlyings = {};
 * int32_t result = dxfg_AdditionalUnderlyings_EMPTY(isolate_thread, &additional_underlyings);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with additionalUnderlyings
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &additional_underlyings->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_AdditionalUnderlyings_EMPTY().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] emptyAdditionalUnderlyings The pointer to result's handle (an empty additional underlyings
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#EMPTY)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the additional underlyings handle.
 */
DXFG_DEPRECATED("use dxfg_v4_AdditionalUnderlyings_EMPTY()")
int32_t dxfg_AdditionalUnderlyings_EMPTY(graal_isolatethread_t *thread,
                                         DXFG_OUT dxfg_additional_underlyings_t **emptyAdditionalUnderlyings);

/**
 * Returns an instance of additional underlyings for specified textual representation.
 * The handle of this instance will be written to the passed pointer.
 *
 * See dxfg_AdditionalUnderlyings_getText() for format specification.
 *
 * Sample:
 * ```c
 * dxfg_additional_underlyings_t* additional_underlyings = {};
 * int32_t result = dxfg_AdditionalUnderlyings_valueOf(isolate_thread, "FIS 53; US$ 45.46", &additional_underlyings);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with additionalUnderlyings
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &additional_underlyings->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_AdditionalUnderlyings_valueOf().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] text The textual representation.
 * @param[out] additionalUnderlyings The pointer to result's handle (additional underlyings for specified
 * textual representation)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the additional underlyings handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#valueOf-java.lang.String-
 */
DXFG_DEPRECATED("use dxfg_v4_AdditionalUnderlyings_valueOf()")
int32_t dxfg_AdditionalUnderlyings_valueOf(graal_isolatethread_t *thread, const char *text,
                                           DXFG_OUT dxfg_additional_underlyings_t **additionalUnderlyings);

/**
 * Returns an instance of additional underlyings for specified internal representation.
 * The handle of this instance will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * dxfg_string_to_double_map_entry_t map_entries[] = {{"FIS", 53}, {"US$", 45.46}, {"ABC", 123}};
 * dxfg_additional_underlyings_t* additional_underlyings = {};
 * int32_t result = dxfg_AdditionalUnderlyings_valueOf2(isolate_thread, map_entries,
 *                                                      sizeof(map_entries) / sizeof(map_entries[0]),
 *                                                      &additional_underlyings);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with additionalUnderlyings
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &additional_underlyings->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_AdditionalUnderlyings_valueOfMap().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] mapEntries The array of map entries ([string, double] pair).
 * @param[in] size A size of the array.
 * @param[out] additionalUnderlyings The pointer to result's handle (additional underlyings for specified
 * internal representation)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the additional underlyings handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#valueOf-java.util.Map-
 */
DXFG_DEPRECATED("use dxfg_v4_AdditionalUnderlyings_valueOfMap()")
int32_t dxfg_AdditionalUnderlyings_valueOf2(graal_isolatethread_t *thread,
                                            const dxfg_string_to_double_map_entry_t *mapEntries, int32_t size,
                                            DXFG_OUT dxfg_additional_underlyings_t **additionalUnderlyings);

/**
 * Returns SPC for specified underlying symbol or 0 is specified symbol is not found.
 * The spc will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * double spc = 0.0;
 * int32_t result = dxfg_AdditionalUnderlyings_getSPC(isolate_thread, "FIS 53; US$ 45.46", "US$", &spc);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("spc = %g\n", spc); // spc = 45.46
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_AdditionalUnderlyings_getSPCFromText().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] text The textual representation of the additional underlyings.
 * @param[in] symbol The symbol by which spc will be obtained.
 * @param[out] spc The pointer to spc.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getSPC-java.lang.String-java.lang.String-
 */
DXFG_DEPRECATED("use dxfg_v4_AdditionalUnderlyings_getSPCFromText()")
int32_t dxfg_AdditionalUnderlyings_getSPC(graal_isolatethread_t *thread, const char *text, const char *symbol,
                                          DXFG_OUT double *spc);

/**
 * Returns textual representation of additional underlyings in the format:
 * ```
 * TEXT ::= "" | LIST
 * LIST ::= AU | AU "; " LIST
 * AU ::= UNDERLYING " " SPC
 * ```
 * Where UNDERLYING is a symbol of underlying instrument and SPC is a number of shares per contract of that underlying.
 * All additional underlyings are listed in the alphabetical order of underlying symbol. In cases when option settles
 * with additional cash the underlying symbol will specify cash symbol and SPC will specify amount of cash.
 * The textual info will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * char *text = {};
 * int32_t result = dxfg_AdditionalUnderlyings_getText(isolate_thread, additional_underlyings, &text);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  text = %s\n", text);
 *     dxfg_String_release(isolate_thread, text);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_AdditionalUnderlyings_getText().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] additionalUnderlyings The additional underlyings handle.
 * @param[out] text A pointer to the textual representation of the additional underlyings.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the textual representation.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getText--
 */
DXFG_DEPRECATED("use dxfg_v4_AdditionalUnderlyings_getText()")
int32_t dxfg_AdditionalUnderlyings_getText(graal_isolatethread_t *thread,
                                           dxfg_additional_underlyings_t *additionalUnderlyings, DXFG_OUT char **text);

/**
 * Returns internal representation of additional underlyings as a map from underlying symbol to its SPC.
 *
 * Sample:
 * ```c
 * dxfg_string_to_double_map_entry_t *map_entries = {};
 * int32_t size = {};
 * int32_t result = dxfg_AdditionalUnderlyings_getMap(isolate_thread, additional_underlyings, &map_entries, &size);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     puts("  map entries:");
 *
 *     for (int i = 0; i < size; i++) {
 *         printf("  [%d]: key = '%s', value = %g\n", i, map_entries[i].key, map_entries[i].value);
 *     }
 *
 *     dxfg_free_string_to_double_map_entries(isolate_thread, map_entries, size);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_AdditionalUnderlyings_getMap().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] additionalUnderlyings The additional underlyings handle.
 * @param[out] mapEntries A pointer to the array of map entries ([string, double] pair).
 * @param[out] size A pointer to the array size.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_free_string_to_double_map_entries() to free the array of map entries.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getMap--
 */
DXFG_DEPRECATED("use dxfg_v4_AdditionalUnderlyings_getMap()")
int32_t dxfg_AdditionalUnderlyings_getMap(graal_isolatethread_t *thread,
                                          dxfg_additional_underlyings_t *additionalUnderlyings,
                                          DXFG_OUT dxfg_string_to_double_map_entry_t **mapEntries,
                                          DXFG_OUT int32_t *size);

/**
 * Returns SPC for specified underlying symbol or 0 is specified symbol is not found.
 * The spc will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * double spc = 0.0;
 * int32_t result = dxfg_AdditionalUnderlyings_getSPC2(isolate_thread, additional_underlyings, "US$", &spc);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("spc = %g\n", spc);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_AdditionalUnderlyings_getSPC().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] additionalUnderlyings The additional underlyings handle.
 * @param[in] symbol The symbol by which spc will be obtained.
 * @param[out] spc The pointer to spc.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getSPC-java.lang.String-
 */
DXFG_DEPRECATED("use dxfg_v4_AdditionalUnderlyings_getSPC()")
int32_t dxfg_AdditionalUnderlyings_getSPC2(graal_isolatethread_t *thread,
                                           dxfg_additional_underlyings_t *additionalUnderlyings, const char *symbol,
                                           DXFG_OUT double *spc);

/** @} */ // end of AdditionalUnderlyingsFunctions

/**
 * @defgroup CFIFunctions CFI' functions
 * @ingroup CFI
 * @{
 */

/**
 * Returns an empty CFI (it has code "XXXXXX").
 * (https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#EMPTY).
 * The handle of this instance will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * dxfg_cfi_t *cfi = {};
 * int32_t result = dxfg_CFI_EMPTY(isolate_thread, &cfi);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with CFI
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &cfi->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_EMPTY().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] emptyCfi The pointer to result's handle (an empty CFI
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#EMPTY)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the CFI handle.
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_EMPTY()")
int32_t dxfg_CFI_EMPTY(graal_isolatethread_t *thread, DXFG_OUT dxfg_cfi_t **emptyCfi);

/**
 * Returns an instance of CFI for specified CFI code. (https://en.wikipedia.org/wiki/ISO_10962)
 * Accepts short code and expands it to 6 letters by appending "X" at the end.
 * The handle of this instance will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * dxfg_cfi_t *cfi = {};
 * int32_t result = dxfg_CFI_valueOf(isolate_thread, "OPEFPS", &cfi);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with CFI
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &cfi->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_valueOf().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] code The CFI code.
 * @param[out] cfi The pointer to result's handle (for specified CFI code)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the CFI handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#valueOf-java.lang.String-
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_valueOf()")
int32_t dxfg_CFI_valueOf(graal_isolatethread_t *thread, const char *code, DXFG_OUT dxfg_cfi_t **cfi);

/**
 * Returns an instance of CFI for specified integer representation of CFI code.
 * The handle of this instance will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * dxfg_cfi_t *cfi = {};
 * int32_t result = dxfg_CFI_valueOf2(isolate_thread, 520264211, &cfi); // OPEFPS
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with CFI
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &cfi->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_valueOfIntCode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] intCode The integer representation of CFI code.
 * @param[out] cfi The pointer to result's handle (for specified CFI code)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the CFI handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#valueOf-int-
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_valueOfIntCode()")
int32_t dxfg_CFI_valueOf2(graal_isolatethread_t *thread, int32_t intCode, DXFG_OUT dxfg_cfi_t **cfi);

/**
 * Returns CFI code. The code always has length of 6 characters.
 * The code will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * char *code = {};
 * int32_t result = dxfg_CFI_getCode(isolate_thread, cfi, &code); // cfi(520264211)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  code = %s\n", code); // code = OPEFPS
 *     dxfg_String_release(isolate_thread, code);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_getCode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] code A pointer to the CFI code.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the CFI code.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getCode--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_getCode()")
int32_t dxfg_CFI_getCode(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT char **code);

/**
 * Returns integer representation of CFI code.
 * The code will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * int32_t intCode = {};
 * int32_t result = dxfg_CFI_getIntCode(isolate_thread, cfi, &intCode); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  intCode = %d\n", intCode); // intCode = 520264211
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_getIntCode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] intCode A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getIntCode--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_getIntCode()")
int32_t dxfg_CFI_getIntCode(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int32_t *intCode);

/**
 * Returns single UTF16 character for instrument category - the first character of the CFI code.
 *
 * Sample:
 * ```c
 * int16_t category = 0;
 * int32_t result = dxfg_CFI_getCategory(isolate_thread, cfi, &category); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  category = 0x%x\n", category); // category = 0x4f
 *     printf("  category = '%c'\n", (char)((uint16_t)category & 0xFF)); // category = 'O'
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_getCategory().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] category A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getCategory--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_getCategory()")
int32_t dxfg_CFI_getCategory(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int16_t *category);

/**
 * Returns single UTF16 character for instrument group - the second character of the CFI code.
 *
 * Sample:
 * ```c
 * int16_t group = 0;
 * int32_t result = dxfg_CFI_getGroup(isolate_thread, cfi, &group); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  group = 0x%x\n", group); // group = 0x50
 *     printf("  group = '%c'\n", (char)((uint16_t)group & 0xFF)); // group = 'P'
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_getGroup().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] group A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getGroup--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_getGroup()")
int32_t dxfg_CFI_getGroup(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int16_t *group);

/**
 * Returns 1 if corresponding instrument is an equity, 0 otherwise.
 *
 * Sample:
 * ```c
 * int32_t is_equity = 0;
 * int32_t result = dxfg_CFI_isEquity(isolate_thread, cfi, &is_equity); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  is_equity = %d\n", is_equity); // is_equity = 0
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_isEquity().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isEquity A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isEquity--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_isEquity()")
int32_t dxfg_CFI_isEquity(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int32_t *isEquity);

/**
 * Returns 1 if corresponding instrument is a debt instrument, 0 otherwise.
 *
 * Sample:
 * ```c
 * int32_t is_debt_instrument = 0;
 * int32_t result = dxfg_CFI_isDebtInstrument(isolate_thread, cfi, &is_debt_instrument); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  is_debt_instrument = %d\n", is_debt_instrument); // is_debt_instrument = 0
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_isDebtInstrument().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isDebtInstrument A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isDebtInstrument--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_isDebtInstrument()")
int32_t dxfg_CFI_isDebtInstrument(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int32_t *isDebtInstrument);

/**
 * Returns 1 if corresponding instrument is an entitlement (right), 0 otherwise.
 *
 * Sample:
 * ```c
 * int32_t is_entitlement = 0;
 * int32_t result = dxfg_CFI_isEntitlement(isolate_thread, cfi, &is_entitlement); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  is_entitlement = %d\n", is_entitlement); // is_entitlement = 0
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_isEntitlement().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isEntitlement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isEntitlement--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_isEntitlement()")
int32_t dxfg_CFI_isEntitlement(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int32_t *isEntitlement);

/**
 * Returns 1 if corresponding instrument is an option, 0 otherwise.
 *
 * Sample:
 * ```c
 * int32_t is_option = 0;
 * int32_t result = dxfg_CFI_isOption(isolate_thread, cfi, &is_option); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  is_option = %d\n", is_option); // is_option = 1
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_isOption().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isOption A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isOption--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_isOption()")
int32_t dxfg_CFI_isOption(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int32_t *isOption);

/**
 * Returns 1 if corresponding instrument is a future, 0 otherwise.
 *
 * Sample:
 * ```c
 * int32_t is_future = 0;
 * int32_t result = dxfg_CFI_isFuture(isolate_thread, cfi, &is_future); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  is_future = %d\n", is_future); // is_future = 0
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_isFuture().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isFuture A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isFuture--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_isFuture()")
int32_t dxfg_CFI_isFuture(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int32_t *isFuture);

/**
 * Returns 1 if corresponding instrument is an other (miscellaneous) instrument, 0 otherwise.
 *
 * Sample:
 * ```c
 * int32_t is_other = 0;
 * int32_t result = dxfg_CFI_isOther(isolate_thread, cfi, &is_other); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  is_other = %d\n", is_other); // is_other = 0
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_isOther().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isOther A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isOther--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_isOther()")
int32_t dxfg_CFI_isOther(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT int32_t *isOther);

/**
 * Returns array of values (dxfg_cfi_value_t) that explain meaning of each character in the CFI code.
 * Array always has length of 6 and each value explains corresponding character.
 * If certain character is inapplicable, unknown or unrecognized - corresponding value will contain reference to this
 * fact.
 *
 * Sample:
 * ```c
 * dxfg_cfi_value_t *values = {};
 * int32_t size = 0;
 * int32_t result = dxfg_CFI_decipher(isolate_thread, cfi, &values, &size); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with values.
 *     // Release.
 *     dxfg_JavaObjectHandler_array_release(isolate_thread, (const dxfg_java_object_handler **)values, size);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_decipher().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] values A pointer to the result.
 * @param[out] size A pointer to the result size.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_array_release() to free the array of values.
 * If there is a need to extend the life of some Value, then you can clone the handle using the
 * dxfg_JavaObjectHandler_clone() function.
 *
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#decipher--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_decipher()")
int32_t dxfg_CFI_decipher(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT dxfg_cfi_value_t **values,
                          DXFG_OUT int32_t *size);

/**
 * Returns short textual description of this CFI code by listing names of all values for the characters in this CFI
 * code.
 * @see dxfg_CFI_decipher method.
 *
 * The code will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * char *description = {};
 * int32_t result = dxfg_CFI_describe(isolate_thread, cfi, &description); // cfi(OPEFPS)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     // description = 'Options; Put options; European; Futures; Physical; Standardized'
 *     printf("  description = '%s'\n", description);
 *     dxfg_String_release(isolate_thread, description);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_describe().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] description A pointer to the description.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the description.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#describe--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_describe()")
int32_t dxfg_CFI_describe(graal_isolatethread_t *thread, dxfg_cfi_t *cfi, DXFG_OUT char **description);

/** @} */ // end of CFIFunctions

/**
 * @defgroup CFIAttributeFunctions CFI.Attribute' functions
 * @ingroup CFIAttribute
 * @{
 */

/**
 * Returns short name of this attribute.
 *
 * Sample:
 * ```c
 * char *name = {};
 * int32_t result = dxfg_CFI_Attribute_getName(isolate_thread, attribute, &name);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  name = %s\n", name);
 *     dxfg_String_release(isolate_thread, name);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_Attribute_getName().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiAttribute The CFI attribute handle.
 * @param[out] name A pointer to the attribute name.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the attribute name.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Attribute.html#getName--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_Attribute_getName()")
int32_t dxfg_CFI_Attribute_getName(graal_isolatethread_t *thread, dxfg_cfi_attribute_t *cfiAttribute,
                                   DXFG_OUT char **name);

/**
 * Returns description of this attribute.
 *
 * Sample:
 * ```c
 * char *description = {};
 * int32_t result = dxfg_CFI_Attribute_getDescription(isolate_thread, attribute, &description);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  description = %s\n", description);
 *     dxfg_String_release(isolate_thread, description);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_Attribute_getDescription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiAttribute The CFI attribute handle.
 * @param[out] description A pointer to the attribute description.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the attribute description.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Attribute.html#getDescription--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_Attribute_getDescription()")
int32_t dxfg_CFI_Attribute_getDescription(graal_isolatethread_t *thread, dxfg_cfi_attribute_t *cfiAttribute,
                                          DXFG_OUT char **description);

/**
 * Returns values of this attribute.
 *
 * Sample:
 * ```c
 * dxfg_cfi_value_t *values = {};
 * int32_t size = 0;
 * int32_t result = dxfg_CFI_Attribute_getValues(isolate_thread, attributes, &values, &size);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with values.
 *     // Release.
 *     dxfg_JavaObjectHandler_array_release(isolate_thread, (const dxfg_java_object_handler **)values, size);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_Attribute_getValues().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiAttribute The CFI attribute handle.
 * @param[out] values A pointer to the result.
 * @param[out] size A pointer to the result size.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_array_release() to free the array of values.
 * If there is a need to extend the life of some Value, then you can clone the handle using the
 * dxfg_JavaObjectHandler_clone() function.
 *
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Attribute.html#getValues--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_Attribute_getValues()")
int32_t dxfg_CFI_Attribute_getValues(graal_isolatethread_t *thread, dxfg_cfi_attribute_t *cfiAttribute,
                                     DXFG_OUT dxfg_cfi_value_t **values, DXFG_OUT int32_t *size);

/** @} */ // end of CFIAttributeFunctions

/**
 * @defgroup CFIValueFunctions CFI.Value' functions
 * @ingroup CFIValue
 * @{
 */

/**
 * Returns attribute that contains this value.
 *
 * Sample:
 * ```c
 * dxfg_cfi_attribute_t *attribute = {};
 * int32_t result = dxfg_CFI_Value_getAttribute(isolate_thread, value, &attribute);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with CFI attribute.
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &attribute->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_Value_getAttribute().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] attribute A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the attribute handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getAttribute--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_Value_getAttribute()")
int32_t dxfg_CFI_Value_getAttribute(graal_isolatethread_t *thread, dxfg_cfi_value_t *cfiValue,
                                    DXFG_OUT dxfg_cfi_attribute_t **attribute);

/**
 * Returns single UTF16 character code of this value.
 *
 * Sample:
 * ```c
 * int16_t code = 0;
 * int32_t result = dxfg_CFI_Value_getCode(isolate_thread, value, &code); // cfi(OPEFPS).values[0]
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  code = 0x%x\n", code); // code = 0x4f
 *     printf("  code = '%c'\n", (char)((uint16_t)code & 0xFF)); // code = 'O'
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_Value_getCode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] code A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getCode--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_Value_getCode()")
int32_t dxfg_CFI_Value_getCode(graal_isolatethread_t *thread, dxfg_cfi_value_t *cfiValue, DXFG_OUT int16_t *code);

/**
 * Returns short name of this value.
 *
 * Sample:
 * ```c
 * char *name = {};
 * int32_t result = dxfg_CFI_Value_getName(isolate_thread, value, &name);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  name = %s\n", name);
 *     dxfg_String_release(isolate_thread, name);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_Value_getName().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] name A pointer to the value name.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the value name.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getName--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_Value_getName()")
int32_t dxfg_CFI_Value_getName(graal_isolatethread_t *thread, dxfg_cfi_value_t *cfiValue, DXFG_OUT char **name);

/**
 * Returns description of this value.
 *
 * Sample:
 * ```c
 * char *description = {};
 * int32_t result = dxfg_CFI_Value_getDescription(isolate_thread, value, &description);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  description = %s\n", description);
 *     dxfg_String_release(isolate_thread, description);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_CFI_Value_getDescription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] description A pointer to the value description.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the value description.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getDescription--
 */
DXFG_DEPRECATED("use dxfg_v4_CFI_Value_getDescription()")
int32_t dxfg_CFI_Value_getDescription(graal_isolatethread_t *thread, dxfg_cfi_value_t *cfiValue,
                                      DXFG_OUT char **description);

/** @} */ // end of CFIValueFunctions

/**
 * @defgroup PriceIncrementsFunctions PriceIncrements' functions
 * @ingroup PriceIncrements
 * @{
 */

/**
 * Returns an empty price increments - it has empty text and sole increment with value 0.
 * The handle of this instance will be written to the passed pointer.
 *
 * Sample:
 * ```c
 * dxfg_price_increments_t* empty_price_increments = {};
 * int32_t result = dxfg_PriceIncrements_EMPTY(isolate_thread, &empty_price_increments);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with empty_price_increments
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &empty_price_increments->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_EMPTY().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] emptyPriceIncrements The pointer to result's handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the price increments handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#EMPTY
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_EMPTY()")
int32_t dxfg_PriceIncrements_EMPTY(graal_isolatethread_t *thread,
                                   DXFG_OUT dxfg_price_increments_t **emptyPriceIncrements);

/**
 * Returns an instance of price increments for specified textual representation.
 * The handle of this instance will be written to the passed pointer.
 *
 * See dxfg_PriceIncrements_getText() for format specification.
 *
 * Sample:
 * ```c
 * dxfg_price_increments_t* price_increments = {};
 * int32_t result = dxfg_PriceIncrements_valueOf(isolate_thread, "0.0001 1; 0.01", &price_increments);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with price_increments
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &price_increments->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_valueOf().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] text The textual representation.
 * @param[out] priceIncrements The pointer to result's handle (price increments for specified
 * textual representation)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the price increments handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#valueOf-java.lang.String-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_valueOf()")
int32_t dxfg_PriceIncrements_valueOf(graal_isolatethread_t *thread, const char *text,
                                     DXFG_OUT dxfg_price_increments_t **priceIncrements);

/**
 * Returns an instance of price increments for specified single increment.
 *
 * Sample:
 * ```c
 * dxfg_price_increments_t* price_increments = {};
 * int32_t result = dxfg_PriceIncrements_valueOf2(isolate_thread, 0.01, &price_increments);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with price_increments
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &price_increments->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_valueOfIncrement().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] increment The single increment.
 * @param[out] priceIncrements The pointer to result's handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the price increments handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#valueOf-double-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_valueOfIncrement()")
int32_t dxfg_PriceIncrements_valueOf2(graal_isolatethread_t *thread, double increment,
                                      DXFG_OUT dxfg_price_increments_t **priceIncrements);

/**
 * Returns an instance of price increments for specified internal representation.
 * See dxfg_PriceIncrements_getPriceIncrements() for details about internal representation.
 *
 * Sample:
 * ```c
 * dxfg_price_increments_t* price_increments = {};
 * double increments[] = {0.01, 20, 0.02, 50, 0.05, 100, 0.1, 250, 0.25, 500, 0.5, 1000, 1, 2500, 2.5};
 * int32_t result = dxfg_PriceIncrements_valueOf3(isolate_thread, increments, sizeof(increments) /
 * sizeof(increments[0]), &price_increments);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do some work with price_increments
 *     // Release
 *     dxfg_JavaObjectHandler_release(isolate_thread, &price_increments->handler);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_valueOfIncrements().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] increments The array of increments in internal representation.
 * @param[in] size A size of the array.
 * @param[out] priceIncrements The pointer to result's handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the price increments handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#valueOf-double:A-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_valueOfIncrements()")
int32_t dxfg_PriceIncrements_valueOf3(graal_isolatethread_t *thread, const double *increments, int32_t size,
                                      DXFG_OUT dxfg_price_increments_t **priceIncrements);

/**
 * Returns textual representation of price increments in the format:
 * <pre>
 * TEXT ::= "" | LIST
 * LIST ::= INCREMENT | RANGE "; " LIST
 * RANGE ::= INCREMENT " " UPPER_LIMIT
 * </pre>
 * Where INCREMENT is a price increment in the given price range and UPPER_LIMIT is the upper bound of that range.
 * All ranges are listed in the ascending order of upper limits and the last range is considered to extend toward
 * infinity and is therefore specified without upper limit. All increments and limits are finite positive numbers.
 * The case with empty text is a special stub used for @ref dxfg_PriceIncrements_EMPTY() "EMPTY" value, it uses sole
 * increment with value 0.
 *
 * Sample:
 * ```c
 * char *text = {};
 * int32_t result = dxfg_PriceIncrements_getText(isolate_thread, price_increments, &text); // valueOf3({0.01, 20, 0.02},
 * 3)
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 *     printf("  text = '%s'\n", text); // text = '0.01 20; 0.02'
 *     dxfg_String_release(isolate_thread, text);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_getText().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] text A pointer to the textual representation of the price increments.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_String_release() to free the textual representation.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getText--
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_getText()")
int32_t dxfg_PriceIncrements_getText(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                     DXFG_OUT char **text);

/**
 * Returns internal representation of price increments as a single array of double values.
 * This array specifies all numbers from textual representation (see dxfg_PriceIncrements_getText()) in the same order.
 * Therefore numbers at even positions are increments and numbers at odd positions are upper limits.
 * The array always has odd length - the infinite upper limit of last range is always omitted and
 * the first increment (for price range adjacent to 0) is always included even for @ref dxfg_PriceIncrements_EMPTY()
 * "EMPTY" value.
 *
 * Sample:
 * ```c
 * double* increments = {};
 * int32_t size = 0;
 * int32_t result = dxfg_PriceIncrements_getPriceIncrements(isolate_thread, price_increments, &increments, &size);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with increments.
 *     // Free the array.
 *     dxfg_free(isolate_thread, increments);
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_getPriceIncrements().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] increments A pointer to the result.
 * @param[out] size A pointer to the result size.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_free() to free the array of increments.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrements--
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_getPriceIncrements()")
int32_t dxfg_PriceIncrements_getPriceIncrements(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                                DXFG_OUT double **increments, DXFG_OUT int32_t *size);

/**
 * Returns first price increment (for price range adjacent to 0), usually the smallest one.
 * Returns 0 for @ref dxfg_PriceIncrements_EMPTY() "EMPTY" price increments.
 *
 * Sample:
 * ```c
 * double price_increment = 0.0;
 * int32_t result = dxfg_PriceIncrements_getPriceIncrement(isolate_thread, price_increments, &price_increment);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with price_increment.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_getPriceIncrement().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] priceIncrement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrement--
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_getPriceIncrement()")
int32_t dxfg_PriceIncrements_getPriceIncrement(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                               DXFG_OUT double *priceIncrement);

/**
 * Returns price increment which shall be applied to the specified price.
 * If price is Not-a-Number (NaN) then first price increment is returned.
 * If price is a breakpoint between two ranges then minimum of upward and downward increments is returned.
 * This method is equivalent to calling @ref dxfg_PriceIncrements_getPriceIncrement3()
 * "dxfg_PriceIncrements_getPriceIncrement3(thread, price_increments, price, 0, &price_increment)".
 *
 * Sample:
 * ```c
 * double price_increment = 0.0;
 * int32_t result = dxfg_PriceIncrements_getPriceIncrement2(isolate_thread, price_increments, 100.34, &price_increment);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with price_increment.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_getPriceIncrementForPrice().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[out] priceIncrement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrement-double-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_getPriceIncrementForPrice()")
int32_t dxfg_PriceIncrements_getPriceIncrement2(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                                double price, DXFG_OUT double *priceIncrement);

/**
 * Returns price increment which shall be applied to the specified price in the specified direction.
 * If price is Not-a-Number (NaN) then first price increment is returned.
 * If price is a breakpoint between two ranges and direction is 0 then minimum of upward and downward increments is
 * returned.
 *
 * Sample:
 * ```c
 * double price_increment = 0.0;
 * int32_t result = dxfg_PriceIncrements_getPriceIncrement3(isolate_thread, price_increments, 100.34, 1,
 *     &price_increment);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with price_increment.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_getPriceIncrementForPriceAndDirection().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[out] priceIncrement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrement-double-int-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_getPriceIncrementForPriceAndDirection()")
int32_t dxfg_PriceIncrements_getPriceIncrement3(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                                double price, int32_t direction, DXFG_OUT double *priceIncrement);

/**
 * Returns first price precision (for price range adjacent to 0), usually the largest one.
 * Returns 0 for @ref dxfg_PriceIncrements_EMPTY() "EMPTY" price increments.
 *
 * Sample:
 * ```c
 * int32_t price_precision = 0;
 * int32_t result = dxfg_PriceIncrements_getPricePrecision(isolate_thread, price_increments, &price_precision);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with price_precision.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_getPricePrecision().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] pricePrecision A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPricePrecision--
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_getPricePrecision()")
int32_t dxfg_PriceIncrements_getPricePrecision(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                               DXFG_OUT int32_t *pricePrecision);

/**
 * Returns price precision for the price range which contains specified price.
 * Price precision is a number of decimal digits after decimal point that are needed
 * to represent all valid prices in the given price range.
 * This method returns price precision in the interval [0, 18] inclusive.
 * If price is Not-a-Number (NaN) then first price precision is returned.
 * If price is a breakpoint between two ranges then precision of lower range is returned.
 *
 * Sample:
 * ```c
 * int32_t price_precision = 0;
 * int32_t result = dxfg_PriceIncrements_getPricePrecision2(isolate_thread, price_increments, 23, &price_precision);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with price_precision.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_getPricePrecisionForPrice().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[out] pricePrecision A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPricePrecision-double-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_getPricePrecisionForPrice()")
int32_t dxfg_PriceIncrements_getPricePrecision2(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                                double price, DXFG_OUT int32_t *pricePrecision);

/**
 * Returns specified price rounded to nearest valid value.
 * If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 * This method is equivalent to calling @ref #dxfg_PriceIncrements_roundPrice2()
 * "dxfg_PriceIncrements_roundPrice2(isolate_thread, price_increments, price, 0, &rounded_price)".
 *
 * Sample:
 * ```c
 * double rounded_price = 0.0;
 * int32_t result = dxfg_PriceIncrements_roundPrice(isolate_thread, price_increments, 100.34, &rounded_price);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with rounded_price.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_roundPrice().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[out] roundedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#roundPrice-double-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_roundPrice()")
int32_t dxfg_PriceIncrements_roundPrice(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                        double price, DXFG_OUT double *roundedPrice);

/**
 * Returns specified price rounded in the specified direction to nearest value
 * that is valid according to price increment rules.
 * If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 * If direction is 0 then price is rounded to nearest valid value.
 *
 * Sample:
 * ```c
 * double rounded_price = 0.0;
 * int32_t result = dxfg_PriceIncrements_roundPrice(isolate_thread, price_increments, 100.34, 1, &rounded_price);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with rounded_price.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_roundPriceWithDirection().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[out] roundedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#roundPrice-double-int-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_roundPriceWithDirection()")
int32_t dxfg_PriceIncrements_roundPrice2(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                         double price, int32_t direction, DXFG_OUT double *roundedPrice);

/**
 * Returns specified price rounded according to specified rounding mode to nearest value
 * that is valid according to price increment rules.
 * If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 *
 * Sample:
 * ```c
 * double rounded_price = 0.0;
 * int32_t result = dxfg_PriceIncrements_roundPrice(isolate_thread, price_increments, 100.34,
 *     DXFG_ROUNDING_MODE_CEILING, &rounded_price);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with rounded_price.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_roundPriceWithRoundingMode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] roundingMode The rounding mode.
 * @param[out] roundedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#roundPrice-double-java.math.RoundingMode-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_roundPriceWithRoundingMode()")
int32_t dxfg_PriceIncrements_roundPrice3(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                         double price, dxfg_rounding_mode_t roundingMode,
                                         DXFG_OUT double *roundedPrice);

/**
 * Returns specified price incremented in the specified direction by appropriate increment and then rounded to nearest
 * valid value. If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 * This method is equivalent to calling @ref dxfg_PriceIncrements_incrementPrice2()
 * "dxfg_PriceIncrements_incrementPrice2(thread, priceIncrements, price, direction, 0, &incrementedPrice)".
 *
 * Sample:
 * ```c
 * double incremented_price = 0.0;
 * int32_t result = dxfg_PriceIncrements_incrementPrice(isolate_thread, price_increments, 100.34, 1,
 *     &incremented_price);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with incremented_price.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_incrementPrice().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[out] incrementedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error.
 * Throws IllegalArgumentException if direction is 0.
 * Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#incrementPrice-double-int-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_incrementPrice()")
int32_t dxfg_PriceIncrements_incrementPrice(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                            double price, int32_t direction, DXFG_OUT double *incrementedPrice);

/**
 * Returns specified price incremented in the specified direction by the maximum of specified step and appropriate
 * increment, and then rounded to nearest valid value.
 * If price is Not-a-Number (NaN) then NaN is returned. If both step and appropriate price increment are 0 then
 * specified price is returned as is.
 * Note that step must be positive even for negative directions.
 *
 * Sample:
 * ```c
 * double incremented_price = 0.0;
 * int32_t result = dxfg_PriceIncrements_incrementPrice2(isolate_thread, price_increments, 100.34, 1, 0.5,
 *     &incremented_price);
 *
 * if (result == DXFG_EXECUTE_SUCCESSFULLY) {
 * // Do something with incremented_price.
 * } else {
 * // Use the dxfg_get_and_clear_thread_exception_t to print the stack trace.
 * }
 * ```
 * @deprecated Use dxfg_v4_PriceIncrements_incrementPriceWithStep().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[in] step The step.
 * @param[out] incrementedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error.
 * Throws IllegalArgumentException if direction is 0 or step uses invalid value.
 * Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#incrementPrice-double-int-double-
 */
DXFG_DEPRECATED("use dxfg_v4_PriceIncrements_incrementPriceWithStep()")
int32_t dxfg_PriceIncrements_incrementPrice2(graal_isolatethread_t *thread, dxfg_price_increments_t *priceIncrements,
                                             double price, int32_t direction, double step,
                                             DXFG_OUT double *incrementedPrice);

/** @} */ // end of PriceIncrementsFunctions

#ifdef __cplusplus
}
#endif

#endif // DXFG_GLOSSARY_H
