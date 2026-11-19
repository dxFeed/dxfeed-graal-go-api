// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Exceptions functions and types declarations
 */

#ifndef DXFG_CATCH_EXCEPTION_H
#define DXFG_CATCH_EXCEPTION_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"

/**
 * @defgroup Exceptions Exceptions
 * @brief The exception that a failed function of the C API leaves for the calling thread.
 *
 * A function of the C API that fails returns an error value (it is given in the `@return` of each function, see
 * dxfg_common.h for #DXFG_EXECUTE_FAIL) and stores the Java exception that caused the failure for the thread that
 * called it. dxfg_get_and_clear_thread_exception_t() returns the stored exception as a dxfg_exception_t, and
 * dxfg_Exception_release() releases it.
 * @{
 */

/**
 * @brief An element of the stack trace of an exception: one frame, the same data as a Java StackTraceElement.
 *
 * The structure and its strings belong to the dxfg_exception_t that contains it and are freed with it by
 * dxfg_Exception_release(). A string that the Java side has no value for is NULL.
 *
 * <a href="https://docs.oracle.com/javase/8/docs/api/java/lang/StackTraceElement.html">Javadoc</a>
 */
typedef struct dxfg_stack_trace_element_t {
    /** The fully qualified name of the class that contains the execution point (a UTF-8 null-terminated string). */
    const char *class_name;
    /** The name of the class loader of that class, NULL if the class loader has no name. */
    const char *class_loader_name;
    /** The name of the source file of that class, NULL if the information is not available. */
    const char *file_name;
    /** The name of the method that contains the execution point (a UTF-8 null-terminated string). */
    const char *method_name;
    /** The line number in the source file, a negative number if the information is not available (-2 for a native
     * method). */
    int32_t line_number;
    /** 1 if the method is a native method, 0 otherwise. */
    int32_t is_native_method;
    /** The name of the module that contains the class, NULL if the class is not in a named module. */
    const char *module_name;
    /** The version of that module, NULL if the module has no version. */
    const char *module_version;
} dxfg_stack_trace_element_t;

/**
 * @brief The stack trace of an exception: a list of its frames, the innermost call first.
 *
 * The structure belongs to the dxfg_exception_t that contains it and is freed with it by dxfg_Exception_release().
 */
typedef struct dxfg_stack_trace_element_list {
    /** The number of the frames; 0 for an empty stack trace. */
    int32_t size;
    /** The array of `size` pointers to the frames; NULL for an empty stack trace. */
    dxfg_stack_trace_element_t **elements;
} dxfg_stack_trace_element_list;

/**
 * @brief An exception that was thrown by the Java side: the class, the message, the stack trace and the cause.
 *
 * It is returned by dxfg_get_and_clear_thread_exception_t() and owned by the caller: release it with
 * dxfg_Exception_release(), which frees the structure, its strings, its stack trace and the whole chain of causes.
 * Every pointer inside the structure belongs to it; copy what has to outlive the structure.
 *
 * <a href="https://docs.oracle.com/javase/8/docs/api/java/lang/Exception.html">Javadoc</a>
 */
typedef struct dxfg_exception_t {
    /** The canonical name of the exception class (for example `java.lang.IllegalArgumentException`), a UTF-8
     * null-terminated string; NULL if the class has no canonical name (for example, an anonymous class). */
    const char *class_name;
    /** The detail message of the exception (a UTF-8 null-terminated string), NULL if the exception has no message. */
    const char *message;
    /**
     * The printed stack trace of the exception (a UTF-8 null-terminated string): the output of Java's
     * `printStackTrace()` without its first line (the class name and the message, see `class_name` and `message`).
     * The lines of the frames and the "Caused by" sections of the causes follow.
     */
    const char *print_stack_trace;
    /** The stack trace of this exception (not of its causes); never NULL, an empty list if the exception has no
     * stack trace. */
    dxfg_stack_trace_element_list *stack_trace;
    /** The exception that caused this one, NULL if there is no cause. The cause has its own cause, and so on. */
    struct dxfg_exception_t *cause;
} dxfg_exception_t;

/**
 * @brief Returns the exception of the last failed function call of the calling thread and clears it.
 *
 * Call it when a function reported an error (see the `@return` of that function) to learn what went wrong. The
 * exception is stored per thread: this function returns the exception of the thread that calls it (the `thread`
 * parameter is not used to choose it), so call it on the thread on which the function failed. The stored exception is
 * cleared by this call: a second call returns NULL until another function fails.
 *
 * A function that completes successfully does not clear the stored exception, so call this function right after the
 * error: otherwise it can return the exception of an earlier failure that was not taken. When several functions fail
 * on a thread before the exception is taken, only the exception of the last failure is returned.
 *
 * Java: `JNIThreadLocalPendingException.get()` and `JNIThreadLocalPendingException.clear()` (GraalVM native image),
 * the result is converted from `Throwable`.
 *
 * @deprecated Use dxfg_v4_get_and_clear_thread_exception().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The exception (release it with dxfg_Exception_release()), or NULL if no exception is stored for the
 * calling thread.
 */
DXFG_DEPRECATED("use dxfg_v4_get_and_clear_thread_exception()")
dxfg_exception_t* dxfg_get_and_clear_thread_exception_t(graal_isolatethread_t *thread);

/**
 * @brief Releases an exception that was returned by dxfg_get_and_clear_thread_exception_t().
 *
 * Frees the structure together with its strings, its stack trace and its whole chain of causes. A NULL pointer is
 * ignored. Do not use the structure or the pointers inside it afterwards. The function does not report errors.
 *
 * Java: not applicable (the structure is allocated by the SDK outside the Java heap).
 *
 * @deprecated Use dxfg_v4_Exception_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] exception The exception to release.
 */
DXFG_DEPRECATED("use dxfg_v4_Exception_release()")
void dxfg_Exception_release(graal_isolatethread_t *thread, dxfg_exception_t *exception);

/** @} */ // end of Exceptions

#ifdef __cplusplus
}
#endif

#endif // DXFG_CATCH_EXCEPTION_H
