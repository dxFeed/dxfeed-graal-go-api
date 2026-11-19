// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_catch_exception.h and the v4 names of its types (see the group @ref V4_Exceptions).
 */

#ifndef DXFG_V4_CATCH_EXCEPTION_H
#define DXFG_V4_CATCH_EXCEPTION_H

#include "dxfg_catch_exception.h"
#include "dxfg_v4_common.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Exceptions Exceptions (v4)
 *  @brief The v4 functions and type names of the group @ref Exceptions. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_stack_trace_element_t */
typedef dxfg_stack_trace_element_t dxfg_v4_stack_trace_element_t;

/** @copydoc dxfg_stack_trace_element_list */
typedef dxfg_stack_trace_element_list dxfg_v4_stack_trace_element_list;

/** @copydoc dxfg_exception_t */
typedef dxfg_exception_t dxfg_v4_exception_t;

/**
 * @brief Writes the exception of the last failed function call of the calling thread and clears it.
 *
 * Call it when a function returned #DXFG_EXECUTE_FAIL to learn what went wrong. The exception is stored per thread:
 * this function takes the exception of the thread that calls it (the `thread` parameter is not used to choose it), so
 * call it on the thread on which the function failed. The stored exception is cleared by this call: a second call
 * writes NULL until another function fails.
 *
 * A function that completes successfully does not clear the stored exception, so call this function right after the
 * error: otherwise it can return the exception of an earlier failure that was not taken. When several functions fail
 * on a thread before the exception is taken, only the exception of the last failure is stored.
 *
 * The function reads the stored exception, converts it and then clears it. If the conversion fails (e.g. there is no
 * memory left), the function returns #DXFG_EXECUTE_FAIL and the stored exception is replaced by the exception of the
 * conversion. If `exception` is NULL, the function returns #DXFG_EXECUTE_FAIL and the stored exception is replaced by
 * an `IllegalArgumentException`.
 *
 * Java: `JNIThreadLocalPendingException.get()` and `JNIThreadLocalPendingException.clear()` (GraalVM native image),
 * the result is converted from `Throwable`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] exception The pointer to the variable that receives the exception (release it with
 * dxfg_v4_Exception_release()), or NULL if no exception is stored for the calling thread; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) when the exception or NULL was written, #DXFG_EXECUTE_FAIL (-1) when the
 * function failed (call it again to get the reason).
 */
int32_t dxfg_v4_get_and_clear_thread_exception(graal_isolatethread_t *thread,
                                               DXFG_OUT dxfg_v4_exception_t **exception);

/**
 * @brief Releases an exception that was written by dxfg_v4_get_and_clear_thread_exception().
 *
 * Frees the structure together with its strings, its stack trace and its whole chain of causes. A NULL pointer is
 * ignored. Do not use the structure or the pointers inside it afterwards. Releasing the same exception twice is
 * undefined behaviour, as for every release function.
 *
 * Java: not applicable (the structure is allocated by the SDK outside the Java heap).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] exception The exception to release, or NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0): the function does not fail.
 */
int32_t dxfg_v4_Exception_release(graal_isolatethread_t *thread, dxfg_v4_exception_t *exception);

/** @} */ // end of V4_Exceptions

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_CATCH_EXCEPTION_H
