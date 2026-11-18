// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_common.h and the v4 names of its types (see the group @ref V4_Common).
 */

#ifndef DXFG_V4_COMMON_H
#define DXFG_V4_COMMON_H

#include "dxfg_common.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Common Common (v4)
 *  @brief The v4 functions and type names of the group @ref Common. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_string_to_double_map_entry_t */
typedef dxfg_string_to_double_map_entry_t dxfg_v4_string_to_double_map_entry_t;

/**
 * @brief Frees a block of the unmanaged (GraalVM) memory that a function of the SDK allocated, e.g. the array of
 * doubles that dxfg_PriceIncrements_getPriceIncrements() returns.
 *
 * Only a block of memory that a function of the SDK allocated and documents to be released with this function can be
 * passed. Do not use the memory after that.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] pointer The pointer to the unmanaged memory that was returned by the SDK.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_free(graal_isolatethread_t *thread, void *pointer);

/** @} */ // end of V4_Common

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_COMMON_H
