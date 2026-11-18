// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_glossary.h and the v4 names of its types (see the group @ref V4_Glossary).
 */

#ifndef DXFG_V4_GLOSSARY_H
#define DXFG_V4_GLOSSARY_H

#include "dxfg_glossary.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Glossary Glossary (v4)
 *  @brief The v4 functions and type names of the group @ref Glossary. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_additional_underlyings_t */
typedef dxfg_additional_underlyings_t dxfg_v4_additional_underlyings_t;

/** @copydoc dxfg_cfi_t */
typedef dxfg_cfi_t dxfg_v4_cfi_t;

/** @copydoc dxfg_cfi_attribute_t */
typedef dxfg_cfi_attribute_t dxfg_v4_cfi_attribute_t;

/** @copydoc dxfg_cfi_value_t */
typedef dxfg_cfi_value_t dxfg_v4_cfi_value_t;

/** @copydoc dxfg_price_increments_t */
typedef dxfg_price_increments_t dxfg_v4_price_increments_t;

/**
 * @brief A list of doubles: the array holds the values themselves, not pointers to them.
 *
 * A list that a function returns (dxfg_v4_PriceIncrements_getPriceIncrements()) is released with
 * dxfg_v4_double_list_release(). A list that the caller passes to a function
 * (dxfg_v4_PriceIncrements_valueOfIncrements()) is only read by it and stays with the caller.
 */
typedef struct dxfg_v4_double_list {
    /** The number of elements; 0 for an empty list. */
    int32_t size;
    /** The array of `size` values; may be NULL for an empty list. */
    double *elements;
} dxfg_v4_double_list;

/**
 * @brief A list of string to double map entries: the array holds the entries themselves, not pointers to them.
 *
 * A list that a function returns (dxfg_v4_AdditionalUnderlyings_getMap()) is released with
 * dxfg_v4_string_to_double_map_entry_list_release(), which frees the keys too. A list that the caller passes to a
 * function (dxfg_v4_AdditionalUnderlyings_valueOfMap()) is only read by it and stays with the caller.
 */
typedef struct dxfg_v4_string_to_double_map_entry_list {
    /** The number of elements; 0 for an empty list. */
    int32_t size;
    /** The array of `size` entries; may be NULL for an empty list. */
    dxfg_v4_string_to_double_map_entry_t *elements;
} dxfg_v4_string_to_double_map_entry_list;

/**
 * @brief A list of CFI value handles.
 *
 * A list that a function returns (dxfg_v4_CFI_decipher(), dxfg_v4_CFI_Attribute_getValues()) is released with
 * dxfg_v4_cfi_value_list_release(), which releases every handle of the list too.
 */
typedef struct dxfg_v4_cfi_value_list {
    /** The number of elements; 0 for an empty list. */
    int32_t size;
    /** The array of `size` CFI value handles; may be NULL for an empty list. */
    dxfg_v4_cfi_value_t **elements;
} dxfg_v4_cfi_value_list;

/**
 * Returns an empty additional underlyings
 * (https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#EMPTY it has empty text and empty
 * map). The handle of this instance will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] emptyAdditionalUnderlyings The pointer to result's handle (an empty additional underlyings
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#EMPTY)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the additional underlyings handle.
 */
int32_t dxfg_v4_AdditionalUnderlyings_EMPTY(graal_isolatethread_t *thread,
                                            DXFG_OUT dxfg_v4_additional_underlyings_t **emptyAdditionalUnderlyings);

/**
 * Returns an instance of additional underlyings for specified textual representation.
 * The handle of this instance will be written to the passed pointer.
 *
 * See dxfg_v4_AdditionalUnderlyings_getText() for format specification.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] text The textual representation.
 * @param[out] additionalUnderlyings The pointer to result's handle (additional underlyings for specified
 * textual representation)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the additional underlyings handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#valueOf-java.lang.String-
 */
int32_t dxfg_v4_AdditionalUnderlyings_valueOf(graal_isolatethread_t *thread, const char *text,
                                              DXFG_OUT dxfg_v4_additional_underlyings_t **additionalUnderlyings);

/**
 * Returns an instance of additional underlyings for a map from underlying symbols to their SPC (the internal
 * representation).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] map The map entries ([string, double] pairs); the list stays with the caller. A NULL list or a list of
 * size 0 is an empty map. A list with a positive size and NULL `elements` is an error.
 * @param[out] additionalUnderlyings The pointer to the variable that receives the additional underlyings (release it
 * with dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#valueOf-java.util.Map-
 */
int32_t dxfg_v4_AdditionalUnderlyings_valueOfMap(graal_isolatethread_t *thread,
                                                 const dxfg_v4_string_to_double_map_entry_list *map,
                                                 DXFG_OUT dxfg_v4_additional_underlyings_t **additionalUnderlyings);

/**
 * Returns SPC for specified underlying symbol or 0 is specified symbol is not found.
 * The spc will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] text The textual representation of the additional underlyings.
 * @param[in] symbol The symbol by which spc will be obtained.
 * @param[out] spc The pointer to spc.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getSPC-java.lang.String-java.lang.String-
 */
int32_t dxfg_v4_AdditionalUnderlyings_getSPCFromText(graal_isolatethread_t *thread, const char *text,
                                                     const char *symbol, DXFG_OUT double *spc);

/**
 * Returns textual representation of additional underlyings in the format:
 * Where UNDERLYING is a symbol of underlying instrument and SPC is a number of shares per contract of that underlying.
 * All additional underlyings are listed in the alphabetical order of underlying symbol. In cases when option settles
 * with additional cash the underlying symbol will specify cash symbol and SPC will specify amount of cash.
 * The textual info will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] additionalUnderlyings The additional underlyings handle.
 * @param[out] text A pointer to the textual representation of the additional underlyings.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the textual representation.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getText--
 */
int32_t dxfg_v4_AdditionalUnderlyings_getText(graal_isolatethread_t *thread,
                                              dxfg_v4_additional_underlyings_t *additionalUnderlyings,
                                              DXFG_OUT char **text);

/**
 * Returns internal representation of additional underlyings as a map from underlying symbol to its SPC.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] additionalUnderlyings The additional underlyings handle.
 * @param[out] map The pointer to the variable that receives the list of the map entries ([string, double] pairs;
 * release it with dxfg_v4_string_to_double_map_entry_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getMap--
 */
int32_t dxfg_v4_AdditionalUnderlyings_getMap(graal_isolatethread_t *thread,
                                             dxfg_v4_additional_underlyings_t *additionalUnderlyings,
                                             DXFG_OUT dxfg_v4_string_to_double_map_entry_list **map);

/**
 * Returns SPC for specified underlying symbol or 0 is specified symbol is not found.
 * The spc will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] additionalUnderlyings The additional underlyings handle.
 * @param[in] symbol The symbol by which spc will be obtained.
 * @param[out] spc The pointer to spc.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/AdditionalUnderlyings.html#getSPC-java.lang.String-
 */
int32_t dxfg_v4_AdditionalUnderlyings_getSPC(graal_isolatethread_t *thread,
                                             dxfg_v4_additional_underlyings_t *additionalUnderlyings,
                                             const char *symbol, DXFG_OUT double *spc);

/**
 * Returns an empty CFI (it has code "XXXXXX").
 * (https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#EMPTY).
 * The handle of this instance will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] emptyCfi The pointer to result's handle (an empty CFI
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#EMPTY)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the CFI handle.
 */
int32_t dxfg_v4_CFI_EMPTY(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_cfi_t **emptyCfi);

/**
 * Returns an instance of CFI for specified CFI code. (https://en.wikipedia.org/wiki/ISO_10962)
 * Accepts short code and expands it to 6 letters by appending "X" at the end.
 * The handle of this instance will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] code The CFI code.
 * @param[out] cfi The pointer to result's handle (for specified CFI code)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the CFI handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#valueOf-java.lang.String-
 */
int32_t dxfg_v4_CFI_valueOf(graal_isolatethread_t *thread, const char *code, DXFG_OUT dxfg_v4_cfi_t **cfi);

/**
 * Returns an instance of CFI for specified integer representation of CFI code.
 * The handle of this instance will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] intCode The integer representation of CFI code.
 * @param[out] cfi The pointer to result's handle (for specified CFI code)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the CFI handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#valueOf-int-
 */
int32_t dxfg_v4_CFI_valueOfIntCode(graal_isolatethread_t *thread, int32_t intCode, DXFG_OUT dxfg_v4_cfi_t **cfi);

/**
 * Returns CFI code. The code always has length of 6 characters.
 * The code will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] code A pointer to the CFI code.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the CFI code.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getCode--
 */
int32_t dxfg_v4_CFI_getCode(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT char **code);

/**
 * Returns integer representation of CFI code.
 * The code will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] intCode A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getIntCode--
 */
int32_t dxfg_v4_CFI_getIntCode(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int32_t *intCode);

/**
 * Returns single UTF16 character for instrument category - the first character of the CFI code.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] category A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getCategory--
 */
int32_t dxfg_v4_CFI_getCategory(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int16_t *category);

/**
 * Returns single UTF16 character for instrument group - the second character of the CFI code.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] group A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#getGroup--
 */
int32_t dxfg_v4_CFI_getGroup(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int16_t *group);

/**
 * Returns 1 if corresponding instrument is an equity, 0 otherwise.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isEquity A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isEquity--
 */
int32_t dxfg_v4_CFI_isEquity(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int32_t *isEquity);

/**
 * Returns 1 if corresponding instrument is a debt instrument, 0 otherwise.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isDebtInstrument A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isDebtInstrument--
 */
int32_t dxfg_v4_CFI_isDebtInstrument(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi,
                                     DXFG_OUT int32_t *isDebtInstrument);

/**
 * Returns 1 if corresponding instrument is an entitlement (right), 0 otherwise.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isEntitlement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isEntitlement--
 */
int32_t dxfg_v4_CFI_isEntitlement(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int32_t *isEntitlement);

/**
 * Returns 1 if corresponding instrument is an option, 0 otherwise.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isOption A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isOption--
 */
int32_t dxfg_v4_CFI_isOption(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int32_t *isOption);

/**
 * Returns 1 if corresponding instrument is a future, 0 otherwise.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isFuture A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isFuture--
 */
int32_t dxfg_v4_CFI_isFuture(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int32_t *isFuture);

/**
 * Returns 1 if corresponding instrument is an other (miscellaneous) instrument, 0 otherwise.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] isOther A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#isOther--
 */
int32_t dxfg_v4_CFI_isOther(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT int32_t *isOther);

/**
 * Returns the values that explain the meaning of each character of the CFI code.
 * The list always has 6 elements, and each value explains the corresponding character.
 * If certain character is inapplicable, unknown or unrecognized - corresponding value will contain reference to this
 * fact.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] values The pointer to the variable that receives the list of the values (release it with
 * dxfg_v4_cfi_value_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#decipher--
 */
int32_t dxfg_v4_CFI_decipher(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi,
                             DXFG_OUT dxfg_v4_cfi_value_list **values);

/**
 * Returns short textual description of this CFI code by listing names of all values for the characters in this CFI
 * code.
 * @see dxfg_CFI_decipher method.
 *
 * The code will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfi The CFI handle.
 * @param[out] description A pointer to the description.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the description.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.html#describe--
 */
int32_t dxfg_v4_CFI_describe(graal_isolatethread_t *thread, dxfg_v4_cfi_t *cfi, DXFG_OUT char **description);

/**
 * Returns short name of this attribute.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiAttribute The CFI attribute handle.
 * @param[out] name A pointer to the attribute name.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the attribute name.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Attribute.html#getName--
 */
int32_t dxfg_v4_CFI_Attribute_getName(graal_isolatethread_t *thread, dxfg_v4_cfi_attribute_t *cfiAttribute,
                                      DXFG_OUT char **name);

/**
 * Returns description of this attribute.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiAttribute The CFI attribute handle.
 * @param[out] description A pointer to the attribute description.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the attribute description.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Attribute.html#getDescription--
 */
int32_t dxfg_v4_CFI_Attribute_getDescription(graal_isolatethread_t *thread, dxfg_v4_cfi_attribute_t *cfiAttribute,
                                             DXFG_OUT char **description);

/**
 * Returns values of this attribute.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiAttribute The CFI attribute handle.
 * @param[out] values The pointer to the variable that receives the list of the values (release it with
 * dxfg_v4_cfi_value_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Attribute.html#getValues--
 */
int32_t dxfg_v4_CFI_Attribute_getValues(graal_isolatethread_t *thread, dxfg_v4_cfi_attribute_t *cfiAttribute,
                                        DXFG_OUT dxfg_v4_cfi_value_list **values);

/**
 * Returns attribute that contains this value.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] attribute A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the attribute handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getAttribute--
 */
int32_t dxfg_v4_CFI_Value_getAttribute(graal_isolatethread_t *thread, dxfg_v4_cfi_value_t *cfiValue,
                                       DXFG_OUT dxfg_v4_cfi_attribute_t **attribute);

/**
 * Returns single UTF16 character code of this value.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] code A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getCode--
 */
int32_t dxfg_v4_CFI_Value_getCode(graal_isolatethread_t *thread, dxfg_v4_cfi_value_t *cfiValue, DXFG_OUT int16_t *code);

/**
 * Returns short name of this value.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] name A pointer to the value name.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the value name.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getName--
 */
int32_t dxfg_v4_CFI_Value_getName(graal_isolatethread_t *thread, dxfg_v4_cfi_value_t *cfiValue, DXFG_OUT char **name);

/**
 * Returns description of this value.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] cfiValue The CFI value handle.
 * @param[out] description A pointer to the value description.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the value description.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/CFI.Value.html#getDescription--
 */
int32_t dxfg_v4_CFI_Value_getDescription(graal_isolatethread_t *thread, dxfg_v4_cfi_value_t *cfiValue,
                                         DXFG_OUT char **description);

/**
 * Returns an empty price increments - it has empty text and sole increment with value 0.
 * The handle of this instance will be written to the passed pointer.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] emptyPriceIncrements The pointer to result's handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the price increments handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#EMPTY
 */
int32_t dxfg_v4_PriceIncrements_EMPTY(graal_isolatethread_t *thread,
                                      DXFG_OUT dxfg_v4_price_increments_t **emptyPriceIncrements);

/**
 * Returns an instance of price increments for specified textual representation.
 * The handle of this instance will be written to the passed pointer.
 *
 * See dxfg_v4_PriceIncrements_getText() for format specification.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] text The textual representation.
 * @param[out] priceIncrements The pointer to result's handle (price increments for specified
 * textual representation)
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the price increments handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#valueOf-java.lang.String-
 */
int32_t dxfg_v4_PriceIncrements_valueOf(graal_isolatethread_t *thread, const char *text,
                                        DXFG_OUT dxfg_v4_price_increments_t **priceIncrements);

/**
 * Returns an instance of price increments for specified single increment.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] increment The single increment.
 * @param[out] priceIncrements The pointer to result's handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the price increments handle.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#valueOf-double-
 */
int32_t dxfg_v4_PriceIncrements_valueOfIncrement(graal_isolatethread_t *thread, double increment,
                                                 DXFG_OUT dxfg_v4_price_increments_t **priceIncrements);

/**
 * Returns an instance of price increments for specified internal representation.
 * See dxfg_v4_PriceIncrements_getPriceIncrements() for details about internal representation.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] increments The increments in internal representation; the list stays with the caller. A NULL list or a
 * list of size 0 gives @ref dxfg_v4_PriceIncrements_EMPTY() "EMPTY". A list with a positive size and NULL `elements`
 * is an error.
 * @param[out] priceIncrements The pointer to the variable that receives the price increments (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#valueOf-double:A-
 */
int32_t dxfg_v4_PriceIncrements_valueOfIncrements(graal_isolatethread_t *thread, const dxfg_v4_double_list *increments,
                                                  DXFG_OUT dxfg_v4_price_increments_t **priceIncrements);

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
 * The case with empty text is a special stub used for @ref dxfg_v4_PriceIncrements_EMPTY() "EMPTY" value, it uses sole
 * increment with value 0.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] text A pointer to the textual representation of the price increments.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_String_release() to free the textual representation.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getText--
 */
int32_t dxfg_v4_PriceIncrements_getText(graal_isolatethread_t *thread, dxfg_v4_price_increments_t *priceIncrements,
                                        DXFG_OUT char **text);

/**
 * Returns internal representation of price increments as a single list of double values.
 * This list specifies all numbers from textual representation (see dxfg_v4_PriceIncrements_getText()) in the same
 * order. Therefore numbers at even positions are increments and numbers at odd positions are upper limits.
 * The list always has odd length - the infinite upper limit of last range is always omitted and
 * the first increment (for price range adjacent to 0) is always included even for @ref dxfg_v4_PriceIncrements_EMPTY()
 * "EMPTY" value.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] increments The pointer to the variable that receives the list of the values (release it with
 * dxfg_v4_double_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrements--
 */
int32_t dxfg_v4_PriceIncrements_getPriceIncrements(graal_isolatethread_t *thread,
                                                   dxfg_v4_price_increments_t *priceIncrements,
                                                   DXFG_OUT dxfg_v4_double_list **increments);

/**
 * Returns first price increment (for price range adjacent to 0), usually the smallest one.
 * Returns 0 for @ref dxfg_v4_PriceIncrements_EMPTY() "EMPTY" price increments.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] priceIncrement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrement--
 */
int32_t dxfg_v4_PriceIncrements_getPriceIncrement(graal_isolatethread_t *thread,
                                                  dxfg_v4_price_increments_t *priceIncrements,
                                                  DXFG_OUT double *priceIncrement);

/**
 * Returns price increment which shall be applied to the specified price.
 * If price is Not-a-Number (NaN) then first price increment is returned.
 * If price is a breakpoint between two ranges then minimum of upward and downward increments is returned.
 * This method is equivalent to calling @ref dxfg_v4_PriceIncrements_getPriceIncrementForPriceAndDirection()
 * "dxfg_PriceIncrements_getPriceIncrement3(thread, price_increments, price, 0, &price_increment)".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[out] priceIncrement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrement-double-
 */
int32_t dxfg_v4_PriceIncrements_getPriceIncrementForPrice(graal_isolatethread_t *thread,
                                                          dxfg_v4_price_increments_t *priceIncrements, double price,
                                                          DXFG_OUT double *priceIncrement);

/**
 * Returns price increment which shall be applied to the specified price in the specified direction.
 * If price is Not-a-Number (NaN) then first price increment is returned.
 * If price is a breakpoint between two ranges and direction is 0 then minimum of upward and downward increments is
 * returned.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[out] priceIncrement A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPriceIncrement-double-int-
 */
int32_t dxfg_v4_PriceIncrements_getPriceIncrementForPriceAndDirection(graal_isolatethread_t *thread,
                                                                      dxfg_v4_price_increments_t *priceIncrements,
                                                                      double price, int32_t direction,
                                                                      DXFG_OUT double *priceIncrement);

/**
 * Returns first price precision (for price range adjacent to 0), usually the largest one.
 * Returns 0 for @ref dxfg_v4_PriceIncrements_EMPTY() "EMPTY" price increments.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[out] pricePrecision A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPricePrecision--
 */
int32_t dxfg_v4_PriceIncrements_getPricePrecision(graal_isolatethread_t *thread,
                                                  dxfg_v4_price_increments_t *priceIncrements,
                                                  DXFG_OUT int32_t *pricePrecision);

/**
 * Returns price precision for the price range which contains specified price.
 * Price precision is a number of decimal digits after decimal point that are needed
 * to represent all valid prices in the given price range.
 * This method returns price precision in the interval [0, 18] inclusive.
 * If price is Not-a-Number (NaN) then first price precision is returned.
 * If price is a breakpoint between two ranges then precision of lower range is returned.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[out] pricePrecision A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#getPricePrecision-double-
 */
int32_t dxfg_v4_PriceIncrements_getPricePrecisionForPrice(graal_isolatethread_t *thread,
                                                          dxfg_v4_price_increments_t *priceIncrements, double price,
                                                          DXFG_OUT int32_t *pricePrecision);

/**
 * Returns specified price rounded to nearest valid value.
 * If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 * This method is equivalent to calling @ref #dxfg_v4_PriceIncrements_roundPriceWithDirection()
 * "dxfg_PriceIncrements_roundPrice2(isolate_thread, price_increments, price, 0, &rounded_price)".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[out] roundedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#roundPrice-double-
 */
int32_t dxfg_v4_PriceIncrements_roundPrice(graal_isolatethread_t *thread, dxfg_v4_price_increments_t *priceIncrements,
                                           double price, DXFG_OUT double *roundedPrice);

/**
 * Returns specified price rounded in the specified direction to nearest value
 * that is valid according to price increment rules.
 * If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 * If direction is 0 then price is rounded to nearest valid value.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[out] roundedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#roundPrice-double-int-
 */
int32_t dxfg_v4_PriceIncrements_roundPriceWithDirection(graal_isolatethread_t *thread,
                                                        dxfg_v4_price_increments_t *priceIncrements, double price,
                                                        int32_t direction, DXFG_OUT double *roundedPrice);

/**
 * Returns specified price rounded according to specified rounding mode to nearest value
 * that is valid according to price increment rules.
 * If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] roundingMode The rounding mode.
 * @param[out] roundedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#roundPrice-double-java.math.RoundingMode-
 */
int32_t dxfg_v4_PriceIncrements_roundPriceWithRoundingMode(graal_isolatethread_t *thread,
                                                           dxfg_v4_price_increments_t *priceIncrements, double price,
                                                           dxfg_v4_rounding_mode_t roundingMode,
                                                           DXFG_OUT double *roundedPrice);

/**
 * Returns specified price incremented in the specified direction by appropriate increment and then rounded to nearest
 * valid value. If price is Not-a-Number (NaN) then NaN is returned.
 * If appropriate price increment is 0 then specified price is returned as is.
 * This method is equivalent to calling @ref dxfg_v4_PriceIncrements_incrementPriceWithStep()
 * "dxfg_PriceIncrements_incrementPrice2(thread, priceIncrements, price, direction, 0, &incrementedPrice)".
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[out] incrementedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error.
 * Throws IllegalArgumentException if direction is 0.
 * Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#incrementPrice-double-int-
 */
int32_t dxfg_v4_PriceIncrements_incrementPrice(graal_isolatethread_t *thread,
                                               dxfg_v4_price_increments_t *priceIncrements, double price,
                                               int32_t direction, DXFG_OUT double *incrementedPrice);

/**
 * Returns specified price incremented in the specified direction by the maximum of specified step and appropriate
 * increment, and then rounded to nearest valid value.
 * If price is Not-a-Number (NaN) then NaN is returned. If both step and appropriate price increment are 0 then
 * specified price is returned as is.
 * Note that step must be positive even for negative directions.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] priceIncrements The price increments handle.
 * @param[in] price The price.
 * @param[in] direction The direction.
 * @param[in] step The step.
 * @param[out] incrementedPrice A pointer to the result.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error.
 * Throws IllegalArgumentException if direction is 0 or step uses invalid value.
 * Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @see https://docs.dxfeed.com/dxfeed/api/com/dxfeed/glossary/PriceIncrements.html#incrementPrice-double-int-double-
 */
int32_t dxfg_v4_PriceIncrements_incrementPriceWithStep(graal_isolatethread_t *thread,
                                                       dxfg_v4_price_increments_t *priceIncrements, double price,
                                                       int32_t direction, double step,
                                                       DXFG_OUT double *incrementedPrice);

/**
 * Releases a list of doubles that a function returned (dxfg_v4_PriceIncrements_getPriceIncrements()): frees the
 * array of the values and the list itself.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_double_list_release(graal_isolatethread_t *thread, dxfg_v4_double_list *list);

/**
 * Releases a list of string to double map entries that a function returned (dxfg_v4_AdditionalUnderlyings_getMap()):
 * frees the key of every entry, the array of the entries and the list itself.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_string_to_double_map_entry_list_release(graal_isolatethread_t *thread,
                                                         dxfg_v4_string_to_double_map_entry_list *list);

/**
 * Releases a list of CFI values that a function returned (dxfg_v4_CFI_decipher(), dxfg_v4_CFI_Attribute_getValues()):
 * releases every CFI value handle, frees the array of the handles and the list itself. To keep a value after that,
 * clone its handle first with dxfg_v4_JavaObjectHandler_clone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] values The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_cfi_value_list_release(graal_isolatethread_t *thread, dxfg_v4_cfi_value_list *values);

/** @} */ // end of V4_Glossary

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_GLOSSARY_H
