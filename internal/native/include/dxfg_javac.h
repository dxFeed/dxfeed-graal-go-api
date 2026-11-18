// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Java functions and types declarations
 */

#ifndef DXFG_JAVAC_H
#define DXFG_JAVAC_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"

/** @defgroup Javac Javac
 *  @brief Java helper types and functions that every other header uses: object handles, lists, strings, executors,
 *  time zones, time periods and formats, authorization tokens.
 *
 *  This group is the reference for the ownership rules of the C API:
 *  - A Java object that a function returns is returned as a handle (dxfg_java_object_handler, or a typed handle
 *    structure that contains it, e.g. dxfg_executor_t). The handle is allocated by the SDK and owned by the caller:
 *    release it with dxfg_JavaObjectHandler_release(). A handle keeps its Java object alive.
 *  - A string (`const char *`) that a function returns is allocated by the SDK and owned by the caller: release it
 *    with dxfg_String_release(). A string that is a field of a structure belongs to that structure and is freed
 *    with it.
 *  - A list that a function returns is allocated by the SDK together with its array of elements and the elements:
 *    release it with dxfg_CList_String_release() (dxfg_string_list) or dxfg_CList_JavaObjectHandler_release()
 *    (dxfg_java_object_handler_list).
 *  - The value returned on error is given in the `@return` of each function and differs between functions. In this
 *    header a function that returns `int32_t` status returns #DXFG_EXECUTE_SUCCESSFULLY (0) or #DXFG_EXECUTE_FAIL
 *    (-1) on error, and a function that returns a pointer returns NULL on error. A function of this header that
 *    returns data as `int32_t`, like dxfg_Object_equals() or dxfg_TimeZone_getRawOffset(), returns INT32_MIN
 *    (-2147483648) on error. dxfg_get_and_clear_thread_exception_t() gives the exception.
 *  @{
 */

/**
 * @brief The handle of a Java object: the reference by which the C code refers to an object that lives in Java.
 *
 * The handle is allocated by the SDK and owned by the caller: release it with dxfg_JavaObjectHandler_release(),
 * which frees the structure and drops the reference to the Java object. As long as any handle of an object
 * exists, the object is not garbage collected. Use dxfg_JavaObjectHandler_clone() to get another handle of the same
 * object.
 *
 * The typed handle structures (e.g. dxfg_executor_t, dxfg_time_zone_t) consist of this structure only (the field
 * `handler`), so a pointer to a typed handle is passed to the functions that take a dxfg_java_object_handler * as
 * `&typed->handler` (or by a cast).
 *
 * Use dxfg_Object_toString() to get the string representation of the object, dxfg_Object_equals() to compare
 * objects and dxfg_Object_hashCode() to get the hash code.
 */
typedef struct dxfg_java_object_handler {
    /** The reference to the Java object (an opaque value of the SDK): do not interpret or change it. */
    void *java_object_handle;
} dxfg_java_object_handler;

/**
 * @brief A handle of a Java object with the same layout as dxfg_java_object_handler.
 *
 * No function of this API takes it: the functions use dxfg_java_object_handler.
 */
typedef struct dxfg_java_object_t {
    /** The reference to the Java object (an opaque value of the SDK): do not interpret or change it. */
    void *java_object_handle;
} dxfg_java_object_t;

/**
 * @brief The common layout of the lists of the SDK: the number of elements and the array of pointers to them.
 *
 * The concrete list types (dxfg_string_list, dxfg_java_object_handler_list, ...) have their own element types and
 * their own release functions; there is no function that frees a dxfg_list as such.
 */
typedef struct dxfg_list {
    /** The number of elements; 0 for an empty list. */
    int32_t size;
    /** The array of `size` pointers to the elements; NULL for an empty list that the SDK created. */
    void **elements;
} dxfg_list;

/**
 * @brief A list of handles of Java objects.
 *
 * A list that a function returns is allocated by the SDK: release it with dxfg_CList_JavaObjectHandler_release(),
 * which releases every handle of the list, frees the array of the elements and the list itself. A list that is
 * passed to a function is only read by the SDK: the caller keeps its ownership.
 */
typedef struct dxfg_java_object_handler_list {
    /** The number of handles; 0 for an empty list. */
    int32_t size;
    /** The array of `size` pointers to the handles; NULL for an empty list that the SDK created. */
    dxfg_java_object_handler **elements;
} dxfg_java_object_handler_list;

/**
 * @brief A list of strings.
 *
 * A list that a function returns is allocated by the SDK: release it with dxfg_CList_String_release(), which frees
 * every string, the array of the strings and the list itself. A list that is passed to a function is only read by
 * the SDK: the caller keeps its ownership.
 */
typedef struct dxfg_string_list {
    /** The number of strings; 0 for an empty list. */
    int32_t size;
    /** The array of `size` pointers to the NUL-terminated UTF-8 strings; NULL for an empty list that the SDK
     * created. */
    const char **elements;
} dxfg_string_list;

/**
 * @brief A handle of a Java Executor object; release it with dxfg_JavaObjectHandler_release().
 *
 * An executor runs the tasks that the SDK gives it, e.g. the calls of the listeners of a subscription or a model.
 * Create it with dxfg_Executors_newFixedThreadPool(), dxfg_Executors_newScheduledThreadPool() or
 * dxfg_ExecutorBaseOnConcurrentLinkedQueue_new().
 *
 * https://docs.oracle.com/en/java/javase/22/docs/api/java.base/java/util/concurrent/Executor.html
 */
typedef struct dxfg_executor_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_executor_t;

/**
 * @brief A handle of a Java InputStream object; release it with dxfg_JavaObjectHandler_release().
 *
 * A stream of bytes that the SDK reads, e.g. the instrument profile readers. Create it with
 * dxfg_ByteArrayInputStream_new().
 *
 * https://docs.oracle.com/en/java/javase/22/docs/api/java.base/java/io/InputStream.html
 */
typedef struct dxfg_input_stream_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_input_stream_t;

/**
 * @brief A handle of a Java TimePeriod object; release it with dxfg_JavaObjectHandler_release().
 *
 * A period of time with a precision of milliseconds that is never negative. Create it with
 * dxfg_TimePeriod_ZERO(), dxfg_TimePeriod_UNLIMITED() or dxfg_TimePeriod_valueOf().
 */
typedef struct dxfg_time_period_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_time_period_t;

/**
 * @brief A handle of a Java TimeFormat object; release it with dxfg_JavaObjectHandler_release().
 *
 * A parser and a formatter of dates and times in an ISO-compatible format for a time zone. Create it with
 * dxfg_TimeFormat_DEFAULT(), dxfg_TimeFormat_GMT() or dxfg_TimeFormat_getInstance().
 */
typedef struct dxfg_time_format_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_time_format_t;

/**
 * @brief A handle of a Java TimeZone object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_TimeZone_getTimeZone(), dxfg_TimeZone_getDefault() or dxfg_TimeFormat_getTimeZone().
 *
 * https://docs.oracle.com/en/java/javase/22/docs/api/java.base/java/util/TimeZone.html
 */
typedef struct dxfg_time_zone_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_time_zone_t;

/**
 * @brief A handle of a Java AuthToken object; release it with dxfg_JavaObjectHandler_release().
 *
 * An authorization token: an authorization scheme (e.g. "Basic" or "Bearer") and its value. Create it with
 * dxfg_AuthToken_valueOf() or one of the dxfg_AuthToken_create*() functions.
 */
typedef struct dxfg_auth_token_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_auth_token_t;

/**
 * @brief The rounding behavior of numerical operations that can discard precision (java.math.RoundingMode), e.g.
 * of dxfg_PriceIncrements_roundPrice3().
 *
 * Specifies a <i>rounding behavior</i> for numerical operations capable of discarding precision. Each rounding mode
 * indicates how the least significant returned digit of a rounded result is to be calculated. If fewer digits are
 * returned than the digits needed to represent the exact numerical result, the discarded digits will be referred to
 * as the <i>discarded fraction</i> regardless the digits' contribution to the value of the number. In other words,
 * considered as a numerical value, the discarded fraction could have an absolute value greater than one.
 *
 * https://docs.oracle.com/en/java/javase/22/docs/api/java.base/java/math/RoundingMode.html
 */
typedef enum dxfg_rounding_mode_t {
    /**
     * Rounding mode to round away from zero. Always increments the digit prior to a non-zero discarded fraction.
     * Note that this rounding mode never decreases the magnitude of the calculated value.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode UP Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_UP` rounding</th>
     *<tr align=right><td>5.5</td>  <td>6</td>
     *<tr align=right><td>2.5</td>  <td>3</td>
     *<tr align=right><td>1.6</td>  <td>2</td>
     *<tr align=right><td>1.1</td>  <td>2</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>-2</td>
     *<tr align=right><td>-1.6</td> <td>-2</td>
     *<tr align=right><td>-2.5</td> <td>-3</td>
     *<tr align=right><td>-5.5</td> <td>-6</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_UP = 0,

    /**
     * Rounding mode to round towards zero. Never increments the digit prior to a discarded fraction (i.e., truncates).
     * Note that this rounding mode never increases the magnitude of the calculated value.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode DOWN Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_DOWN` rounding</th>
     *<tr align=right><td>5.5</td>  <td>5</td>
     *<tr align=right><td>2.5</td>  <td>2</td>
     *<tr align=right><td>1.6</td>  <td>1</td>
     *<tr align=right><td>1.1</td>  <td>1</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>-1</td>
     *<tr align=right><td>-1.6</td> <td>-1</td>
     *<tr align=right><td>-2.5</td> <td>-2</td>
     *<tr align=right><td>-5.5</td> <td>-5</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_DOWN,

    /**
     * Rounding mode to round towards positive infinity. If the result is positive, behaves as for
     * #DXFG_ROUNDING_MODE_UP; if negative, behaves as for #DXFG_ROUNDING_MODE_DOWN. Note that this rounding mode never
     * decreases the calculated value.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode CEILING Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_CEILING` rounding</th>
     *<tr align=right><td>5.5</td>  <td>6</td>
     *<tr align=right><td>2.5</td>  <td>3</td>
     *<tr align=right><td>1.6</td>  <td>2</td>
     *<tr align=right><td>1.1</td>  <td>2</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>-1</td>
     *<tr align=right><td>-1.6</td> <td>-1</td>
     *<tr align=right><td>-2.5</td> <td>-2</td>
     *<tr align=right><td>-5.5</td> <td>-5</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_CEILING,

    /**
     * Rounding mode to round towards negative infinity. If the result is positive, behave as for
     * #DXFG_ROUNDING_MODE_DOWN; if negative, behave as for #DXFG_ROUNDING_MODE_UP. Note that this rounding mode never
     * increases the calculated value.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode FLOOR Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_FLOOR` rounding</th>
     *<tr align=right><td>5.5</td>  <td>5</td>
     *<tr align=right><td>2.5</td>  <td>2</td>
     *<tr align=right><td>1.6</td>  <td>1</td>
     *<tr align=right><td>1.1</td>  <td>1</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>-2</td>
     *<tr align=right><td>-1.6</td> <td>-2</td>
     *<tr align=right><td>-2.5</td> <td>-3</td>
     *<tr align=right><td>-5.5</td> <td>-6</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_FLOOR,

    /**
     * Rounding mode to round towards "nearest neighbor" unless both neighbors are equidistant, in which case round up.
     * Behaves as for #DXFG_ROUNDING_MODE_UP if the discarded fraction is <= 0.5; otherwise, behaves as for
     * #DXFG_ROUNDING_MODE_DOWN. Note that this is the rounding mode commonly taught at school.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode HALF_UP Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_HALF_UP` rounding</th>
     *<tr align=right><td>5.5</td>  <td>6</td>
     *<tr align=right><td>2.5</td>  <td>3</td>
     *<tr align=right><td>1.6</td>  <td>2</td>
     *<tr align=right><td>1.1</td>  <td>1</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>-1</td>
     *<tr align=right><td>-1.6</td> <td>-2</td>
     *<tr align=right><td>-2.5</td> <td>-3</td>
     *<tr align=right><td>-5.5</td> <td>-6</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_HALF_UP,

    /**
     * Rounding mode to round towards "nearest neighbor" unless both neighbors are equidistant, in which case round
     * down. Behaves as for #DXFG_ROUNDING_MODE_UP if the discarded fraction is > 0.5; otherwise, behaves as for
     * #DXFG_ROUNDING_MODE_DOWN.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode HALF_DOWN Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_HALF_DOWN` rounding</th>
     *<tr align=right><td>5.5</td>  <td>5</td>
     *<tr align=right><td>2.5</td>  <td>2</td>
     *<tr align=right><td>1.6</td>  <td>2</td>
     *<tr align=right><td>1.1</td>  <td>1</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>-1</td>
     *<tr align=right><td>-1.6</td> <td>-2</td>
     *<tr align=right><td>-2.5</td> <td>-2</td>
     *<tr align=right><td>-5.5</td> <td>-5</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_HALF_DOWN,

    /**
     * Rounding mode to round towards the "nearest neighbor" unless both neighbors are equidistant, in which case, round
     * towards the even neighbor. Behaves as for #DXFG_ROUNDING_MODE_HALF_UP if the digit to the left of the
     * discarded fraction is odd; behaves as for #DXFG_ROUNDING_MODE_HALF_DOWN if it's even. Note that this is the
     * rounding mode that statistically minimizes cumulative error when applied repeatedly over a sequence of
     * calculations. It is sometimes known as "Banker's rounding" and is chiefly used in the USA. This rounding mode is
     * analogous to the rounding policy used for `float` and `double` arithmetic in Java.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode HALF_EVEN Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_HALF_EVEN` rounding</th>
     *<tr align=right><td>5.5</td>  <td>6</td>
     *<tr align=right><td>2.5</td>  <td>2</td>
     *<tr align=right><td>1.6</td>  <td>2</td>
     *<tr align=right><td>1.1</td>  <td>1</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>-1</td>
     *<tr align=right><td>-1.6</td> <td>-2</td>
     *<tr align=right><td>-2.5</td> <td>-2</td>
     *<tr align=right><td>-5.5</td> <td>-6</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_HALF_EVEN,

    /**
     * Rounding mode to assert that the requested operation has an exact result, hence no rounding is necessary. If
     * this rounding mode is specified on an operation that yields an inexact result, a Java `ArithmeticException` is
     * thrown.
     *
     *<p>Example:
     *<table border>
     * <caption><b>Rounding mode UNNECESSARY Examples</b></caption>
     *<tr valign=top><th>Input Number</th>
     *    <th>Input rounded to one digit<br> with `DXFG_ROUNDING_MODE_UNNECESSARY` rounding</th>
     *<tr align=right><td>5.5</td>  <td>throw `ArithmeticException`</td>
     *<tr align=right><td>2.5</td>  <td>throw `ArithmeticException`</td>
     *<tr align=right><td>1.6</td>  <td>throw `ArithmeticException`</td>
     *<tr align=right><td>1.1</td>  <td>throw `ArithmeticException`</td>
     *<tr align=right><td>1.0</td>  <td>1</td>
     *<tr align=right><td>-1.0</td> <td>-1</td>
     *<tr align=right><td>-1.1</td> <td>throw `ArithmeticException`</td>
     *<tr align=right><td>-1.6</td> <td>throw `ArithmeticException`</td>
     *<tr align=right><td>-2.5</td> <td>throw `ArithmeticException`</td>
     *<tr align=right><td>-5.5</td> <td>throw `ArithmeticException`</td>
     *</table>
     */
    DXFG_ROUNDING_MODE_UNNECESSARY,
} dxfg_rounding_mode_t;

/**
 * @brief The function registered by dxfg_Object_finalize(): called after the Java object has been garbage collected.
 *
 * It is called by a daemon Java thread of the SDK, not by the thread that registered it. The calls are made one after
 * another on that thread, so a function that blocks delays the other ones. A registration is called once. The SDK
 * does not use the user data and does not free it.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_data The pointer that was passed to dxfg_Object_finalize().
 */
typedef void (*dxfg_finalize_function)(graal_isolatethread_t *thread, void *user_data);

/**
 * Registers a function that is called with the user data after the Java object is garbage collected.
 *
 * The Java object is not collected while a handle of it exists, so the function is called after all the handles of
 * the object have been released and nothing else in Java refers to the object. The collected objects are checked
 * about once a second. The function is called by a daemon Java thread of the SDK (see dxfg_finalize_function), once
 * for each registration.
 *
 * It does not stop the callbacks of a listener: see dxfg_NativeListener_deactivate() for that.
 *
 * @deprecated Use dxfg_v4_Object_finalize().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] handler The handle of the object to watch.
 * @param[in] finalize The function to call.
 * @param[in] user_data The user data to pass to the function; the SDK does not free it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Object_finalize()")
int32_t dxfg_Object_finalize(graal_isolatethread_t *thread, dxfg_java_object_handler *handler,
                             dxfg_finalize_function finalize, void *user_data);

/**
 * Frees the memory of the handle structure and releases the reference to the Java object. Do not use the handle after
 * that. Passing NULL does nothing.
 *
 * If the handle is the last handle (see dxfg_JavaObjectHandler_clone()) of a listener created by a
 * `dxfg_*Listener_new` function, the listener is deactivated (see dxfg_NativeListener_deactivate()): the function
 * waits for the listener's callbacks that are in progress on other threads and after it returns the listener's
 * function is not called anymore, so its user data can be freed. Remove the listener (e.g.
 * dxfg_DXFeedSubscription_removeEventListener()) before releasing it, otherwise it stays added and does nothing.
 *
 * Do not hold the locks that the listener's function may take while releasing the listener: it may deadlock.
 * Releasing a listener from its own callback is allowed (the function does not wait for that call).
 *
 * @deprecated Use dxfg_v4_JavaObjectHandler_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] handler The handle to release (e.g. `&executor->handler` for a dxfg_executor_t).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_JavaObjectHandler_release()")
int32_t dxfg_JavaObjectHandler_release(graal_isolatethread_t *thread, dxfg_java_object_handler *handler);

/**
 * Deactivates a listener created by a `dxfg_*Listener_new` function (e.g. dxfg_DXFeedEventListener_new(),
 * dxfg_PropertyChangeListener_new()) regardless of the number of its handles: waits for the listener's callbacks
 * that are in progress on other threads, after that the listener's function is not called anymore, so its user data
 * can be freed. The listener stays added where it was added, but does nothing, and cannot be activated again.
 *
 * Releasing the last handle of the listener with dxfg_JavaObjectHandler_release() deactivates it as well.
 * Do not hold the locks that the listener's function may take while deactivating the listener: it may deadlock.
 * Deactivating a listener from its own callback is allowed (the function does not wait for that call).
 *
 * @deprecated Use dxfg_v4_NativeListener_deactivate().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error
 * (e.g. the object is not a listener created by a `dxfg_*Listener_new` function).
 * Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_NativeListener_deactivate()")
int32_t dxfg_NativeListener_deactivate(graal_isolatethread_t *thread, dxfg_java_object_handler *listener);

/**
 * Creates another handle of the Java object: allocates the memory of the new handle structure. Both handles refer to
 * the same object and are released independently: the object stays alive until all of its handles are released.
 *
 * Sample:
 * ```c
 * dxfg_cfi_value_t *values = {};
 * int32_t size = 0;
 * dxfg_cfi_value_t *clone_of_options_cfi_value = {};
 *
 * // ...
 *
 * dxfg_CFI_decipher(isolate_thread, cfi, &values, &size);
 *
 * for (int32_t i = 0; i < size; i++) {
 * // ...
 *
 *     dxfg_cfi_value_t *cfi_value = (dxfg_cfi_value_t *)(&values[i].handler);
 * // ...
 *     dxfg_JavaObjectHandler_clone(
 *         isolate_thread,
 *         (dxfg_java_object_handler *)cfi_value,
 *         (dxfg_java_object_handler **)(&clone_of_options_cfi_value));
 * }
 * ```
 *
 * @deprecated Use dxfg_v4_JavaObjectHandler_clone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] javaObjectHandler The handle to clone; it stays valid.
 * @param[out] javaObjectHandlerClone The pointer to which the pointer to the new handle is written; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the object clone handle.
 */
DXFG_DEPRECATED("use dxfg_v4_JavaObjectHandler_clone()")
int32_t dxfg_JavaObjectHandler_clone(graal_isolatethread_t *thread, dxfg_java_object_handler *javaObjectHandler,
                                     dxfg_java_object_handler **javaObjectHandlerClone);

/**
 * Releases an array of handles (or simple structures containing only a handle) that a function returned as one block
 * of consecutive structures: releases every handle of the array and frees the memory of the array.
 *
 * Sample:
 * ```c
 * // dxfg_cfi_value_t *values = {};
 * // int32_t size = {};
 * dxfg_JavaObjectHandler_array_release(isolate_thread, (const dxfg_java_object_handler **)values, size);
 * ```
 *
 * @deprecated Use dxfg_v4_JavaObjectHandler_array_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] handlersArray The array of handles. If it is NULL or `size` is 0, nothing is done.
 * @param[in] size The number of handles in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_JavaObjectHandler_array_release()")
int32_t dxfg_JavaObjectHandler_array_release(graal_isolatethread_t *thread,
                                             const dxfg_java_object_handler **handlersArray, int32_t size);

/**
 * Creates a new Java Object.
 *
 * @deprecated Use dxfg_v4_Object_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] object The pointer to which the pointer to the handle of the new object is written; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 * Use dxfg_JavaObjectHandler_release() to free the object handle.
 *
 * Java: `Object()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Object_new()")
int32_t dxfg_Object_new(graal_isolatethread_t *thread, DXFG_OUT dxfg_java_object_handler **object);

/**
 * Returns the string representation of the Java object. It works for the handle of any Java object, including the
 * typed handles (e.g. dxfg_time_zone_t, dxfg_auth_token_t).
 *
 * @deprecated Use dxfg_v4_Object_toString().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the object. A NULL handle is an error.
 * @return The UTF-8 string, NUL-terminated (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Object.toString()`.
 */
DXFG_DEPRECATED("use dxfg_v4_Object_toString()")
const char *dxfg_Object_toString(graal_isolatethread_t *thread, dxfg_java_object_handler *object);

/**
 * Checks whether the Java objects are equal: the handles of the same object, and the handles of the objects that are
 * equal in the sense of the Java `equals` method of the first object. Two NULL handles are equal, a NULL handle and
 * a non-NULL handle are not.
 *
 * @deprecated Use dxfg_v4_Object_equals().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the first object.
 * @param[in] other The handle of the second object.
 * @return 1 if the objects are equal, 0 if they are not, or INT32_MIN (-2147483648) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `Objects.equals(Object, Object)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Object_equals()")
int32_t dxfg_Object_equals(graal_isolatethread_t *thread, dxfg_java_object_handler *object,
                           dxfg_java_object_handler *other);

/**
 * Returns the hash code of the Java object: the result of its Java `hashCode` method, or 0 for a NULL handle.
 *
 * @deprecated Use dxfg_v4_Object_hashCode().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the object.
 * @return The hash code, or INT32_MIN (-2147483648) on error. INT32_MIN can also be a hash code: use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `Objects.hashCode(Object)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Object_hashCode()")
int32_t dxfg_Object_hashCode(graal_isolatethread_t *thread, dxfg_java_object_handler *object);

/**
 * Compares the Java objects by their natural order. The objects must implement the Java `Comparable` interface and
 * be comparable to each other, otherwise the function fails. The handles of the same object are equal; a NULL handle
 * is less than a non-NULL handle.
 *
 * @deprecated Use dxfg_v4_Comparable_compareTo().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the first object.
 * @param[in] other The handle of the second object.
 * @return A negative value if the first object is less than the second one, 0 if they are equal, a positive value if
 * it is greater, or INT32_MIN (-2147483648) on error. Use dxfg_get_and_clear_thread_exception_t() to determine if
 * an exception was thrown.
 *
 * Java: `Comparable.compareTo(Object)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Comparable_compareTo()")
int32_t dxfg_Comparable_compareTo(graal_isolatethread_t *thread, dxfg_java_object_handler *object,
                                  dxfg_java_object_handler *other);

/**
 * Releases a list of handles that a function returned: releases every handle of the list, frees the array of the
 * elements and the list itself.
 *
 * @deprecated Use dxfg_v4_java_object_handler_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_java_object_handler_list_release()")
int32_t dxfg_CList_JavaObjectHandler_release(graal_isolatethread_t *thread, dxfg_java_object_handler_list *list);

/**
 * Frees a string that a function returned (e.g. dxfg_Object_toString()). Do not pass a string that is a field of
 * a structure, such as a string of an event: it is freed with the structure.
 *
 * @deprecated Use dxfg_v4_String_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The string to free. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_String_release()")
int32_t dxfg_String_release(graal_isolatethread_t *thread, const char *string);

/**
 * Releases a list of strings that a function returned: frees every string of the list, the array of the strings and
 * the list itself.
 *
 * @deprecated Use dxfg_v4_string_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] stringList The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_string_list_release()")
int32_t dxfg_CList_String_release(graal_isolatethread_t *thread, dxfg_string_list *stringList);

/**
 * Creates an executor with a fixed number of threads that runs the submitted tasks. The threads are daemon threads.
 * The SDK has no function that shuts the executor down.
 *
 * Read the "Threads and locks" section of
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedSubscription.html for how the executor is used by the
 * listeners.
 *
 * @deprecated Use dxfg_v4_Executors_newFixedThreadPool().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] nThreads The number of threads of the pool; must be positive.
 * @param[in] nameThreads The prefix of the names of the threads: they are named `<nameThreads>-1`,
 * `<nameThreads>-2`, and so on.
 * @return The executor (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Executors.newFixedThreadPool(int, ThreadFactory)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Executors_newFixedThreadPool()")
dxfg_executor_t *dxfg_Executors_newFixedThreadPool(graal_isolatethread_t *thread, int nThreads,
                                                   const char *nameThreads);

/**
 * Creates an executor that runs the submitted tasks, also with a delay or periodically. The threads are daemon
 * threads. The SDK has no function that shuts the executor down.
 *
 * See dxfg_Executors_newFixedThreadPool() for how the listeners use the executor.
 *
 * @deprecated Use dxfg_v4_Executors_newScheduledThreadPool().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] nThreads The number of threads to keep in the pool, even if they are idle; must not be negative.
 * @param[in] nameThreads The prefix of the names of the threads: they are named `<nameThreads>-1`,
 * `<nameThreads>-2`, and so on.
 * @return The executor (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `Executors.newScheduledThreadPool(int, ThreadFactory)`.
 */
DXFG_DEPRECATED("use dxfg_v4_Executors_newScheduledThreadPool()")
dxfg_executor_t *dxfg_Executors_newScheduledThreadPool(graal_isolatethread_t *thread, int nThreads,
                                                       const char *nameThreads);

/**
 * Creates an executor that does not run the tasks by itself: a submitted task is added to a queue and runs when
 * dxfg_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks() is called, on the thread that calls it. Use it
 * to have the calls of the listeners made on a thread of your own.
 *
 * @deprecated Use dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The executor (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_new()")
dxfg_executor_t *dxfg_ExecutorBaseOnConcurrentLinkedQueue_new(graal_isolatethread_t *thread);

/**
 * Runs the tasks that are in the queue of the executor, on the calling thread, until the queue is empty (the tasks
 * added by the running tasks are run as well). If a task throws an exception, the function fails and the rest of
 * the tasks stay in the queue.
 *
 * @deprecated Use dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] executor The executor created by dxfg_ExecutorBaseOnConcurrentLinkedQueue_new(); any other executor
 * makes the function fail.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks()")
int32_t dxfg_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks(graal_isolatethread_t *thread,
                                                                        dxfg_executor_t *executor);

/**
 * Creates an input stream that reads the bytes of the memory block. The bytes are copied: the caller keeps the memory
 * block and may free it after the function returns.
 *
 * @deprecated Use dxfg_v4_ByteArrayInputStream_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] bytes The bytes of the stream; must point to at least `size` bytes.
 * @param[in] size The number of bytes; must not be negative.
 * @return The input stream (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `ByteArrayInputStream(byte[])`.
 */
DXFG_DEPRECATED("use dxfg_v4_ByteArrayInputStream_new()")
dxfg_input_stream_t *dxfg_ByteArrayInputStream_new(graal_isolatethread_t *thread, const char *bytes, int32_t size);

/**
 * Requests a garbage collection. It is only for testing, e.g. of the functions registered by dxfg_Object_finalize().
 *
 * @deprecated Use dxfg_v4_gc().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `System.gc()`.
 */
DXFG_DEPRECATED("use dxfg_v4_gc()")
int32_t dxfg_gc(graal_isolatethread_t *thread);

/**
 * Always fails: throws a Java exception (a RuntimeException "recursion" with a cause "cause", thrown several calls
 * deep, so that its stack trace has several frames). It is only for testing of dxfg_get_and_clear_thread_exception_t().
 *
 * @deprecated Use dxfg_v4_throw_exception().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return Always NULL. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_throw_exception()")
dxfg_java_object_handler *dxfg_throw_exception(graal_isolatethread_t *thread);

/**
 * Returns the time zone for the ID. If the ID is not known, the result is the GMT zone.
 *
 * @deprecated Use dxfg_v4_TimeZone_getTimeZone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ID The ID of the time zone, e.g. "America/New_York" or "GMT+03:00".
 * @return The time zone (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeZone.getTimeZone(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getTimeZone()")
dxfg_time_zone_t *dxfg_TimeZone_getTimeZone(graal_isolatethread_t *thread, const char *ID);

/**
 * Returns the default time zone of the Java virtual machine.
 *
 * @deprecated Use dxfg_v4_TimeZone_getDefault().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The time zone (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeZone.getDefault()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getDefault()")
dxfg_time_zone_t *dxfg_TimeZone_getDefault(graal_isolatethread_t *thread);

/**
 * Returns the ID of the time zone.
 *
 * @deprecated Use dxfg_v4_TimeZone_getID().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @return The ID (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeZone.getID()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getID()")
const char *dxfg_TimeZone_getID(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone);

/**
 * Returns the long name of the time zone (for the standard time) in the default locale of the Java virtual machine.
 *
 * @deprecated Use dxfg_v4_TimeZone_getDisplayName().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @return The name (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeZone.getDisplayName()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getDisplayName()")
const char *dxfg_TimeZone_getDisplayName(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone);

/**
 * Returns the name of the time zone in the default locale of the Java virtual machine, for the standard or the
 * daylight saving time, in the long or the short style.
 *
 * @deprecated Use dxfg_v4_TimeZone_getDisplayNameWithStyle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] daylight 1 for the name of the daylight saving time, any other value for the name of the standard time.
 * @param[in] style 0 for the short name (`TimeZone.SHORT`), 1 for the long name (`TimeZone.LONG`); any other value is
 * an error.
 * @return The name (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeZone.getDisplayName(boolean, int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getDisplayNameWithStyle()")
const char *dxfg_TimeZone_getDisplayName2(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone, int32_t daylight,
                                          int32_t style);

/**
 * Returns the amount of time to be added to the local standard time to get the local wall clock time, in
 * milliseconds. It is 0 if the time zone does not observe the daylight saving time.
 *
 * @deprecated Use dxfg_v4_TimeZone_getDSTSavings().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @return The savings in milliseconds, or INT32_MIN (-2147483648) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.getDSTSavings()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getDSTSavings()")
int32_t dxfg_TimeZone_getDSTSavings(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone);

/**
 * Checks whether the time zone uses the daylight saving time.
 *
 * @deprecated Use dxfg_v4_TimeZone_useDaylightTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @return 1 if the time zone uses the daylight saving time, 0 if it does not, or INT32_MIN (-2147483648) on error.
 * Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.useDaylightTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_useDaylightTime()")
int32_t dxfg_TimeZone_useDaylightTime(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone);

/**
 * Checks whether the time zone is currently in the daylight saving time, or a transition from the standard time to the
 * daylight saving time occurs at any future time.
 *
 * @deprecated Use dxfg_v4_TimeZone_observesDaylightTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @return 1 if the time zone observes the daylight saving time, 0 if it does not, or INT32_MIN (-2147483648) on
 * error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.observesDaylightTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_observesDaylightTime()")
int32_t dxfg_TimeZone_observesDaylightTime(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone);

/**
 * Returns the offset of the time zone from UTC at the date, in milliseconds. If the daylight saving time is in effect
 * at the date, the offset includes the saving.
 *
 * @deprecated Use dxfg_v4_TimeZone_getOffset().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] date The date: the number of milliseconds since 1970-01-01T00:00:00Z.
 * @return The offset in milliseconds, or INT32_MIN (-2147483648) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.getOffset(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getOffset()")
int32_t dxfg_TimeZone_getOffset(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone, int64_t date);

/**
 * Returns the offset of the time zone from UTC at the local date, in milliseconds, adjusted for the daylight saving
 * time.
 *
 * @deprecated Use dxfg_v4_TimeZone_getOffsetByCalendarFields().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] era The era of the date: 1 for AD (`GregorianCalendar.AD`) or 0 for BC (`GregorianCalendar.BC`).
 * @param[in] year The year in the era.
 * @param[in] month The month of the year, 0-based: 0 for January.
 * @param[in] day The day of the month.
 * @param[in] dayOfWeek The day of the week: 1 for Sunday (`Calendar.SUNDAY`) to 7 for Saturday.
 * @param[in] milliseconds The milliseconds in the day, in the standard local time.
 * @return The offset in milliseconds, or INT32_MIN (-2147483648) on error (e.g. a parameter is out of range). Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.getOffset(int, int, int, int, int, int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getOffsetByCalendarFields()")
int32_t dxfg_TimeZone_getOffset2(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone, int32_t era, int32_t year,
                                 int32_t month, int32_t day, int32_t dayOfWeek, int32_t milliseconds);

/**
 * Returns the offset of the time zone from UTC, in milliseconds, that does not include the daylight saving time.
 *
 * @deprecated Use dxfg_v4_TimeZone_getRawOffset().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @return The raw offset in milliseconds, or INT32_MIN (-2147483648) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.getRawOffset()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_getRawOffset()")
int32_t dxfg_TimeZone_getRawOffset(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone);

/**
 * Checks whether the time zone has the same rules and offset as the other time zone (they may differ only in the ID).
 *
 * @deprecated Use dxfg_v4_TimeZone_hasSameRules().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] other The other time zone.
 * @return 1 if the rules are the same, 0 if they are not (also if `other` is NULL), or INT32_MIN (-2147483648) on
 * error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.hasSameRules(TimeZone)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_hasSameRules()")
int32_t dxfg_TimeZone_hasSameRules(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone, dxfg_time_zone_t *other);

/**
 * Checks whether the date is in the daylight saving time of the time zone.
 *
 * @deprecated Use dxfg_v4_TimeZone_inDaylightTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] date The date: the number of milliseconds since 1970-01-01T00:00:00Z.
 * @return 1 if the date is in the daylight saving time, 0 if it is not, or INT32_MIN (-2147483648) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.inDaylightTime(Date)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_inDaylightTime()")
int32_t dxfg_TimeZone_inDaylightTime(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone, int64_t date);

/**
 * Sets the ID of the time zone. It changes the Java object, so all the handles of the time zone see the new ID.
 *
 * @deprecated Use dxfg_v4_TimeZone_setID().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] ID The new ID.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or INT32_MIN (-2147483648), not
 * #DXFG_EXECUTE_FAIL, on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.setID(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_setID()")
int32_t dxfg_TimeZone_setID(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone, const char *ID);

/**
 * Sets the offset of the time zone from UTC that does not include the daylight saving time. It changes the Java
 * object, so all the handles of the time zone see the new offset.
 *
 * @deprecated Use dxfg_v4_TimeZone_setRawOffset().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] offsetMillis The raw offset in milliseconds.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or INT32_MIN (-2147483648), not
 * #DXFG_EXECUTE_FAIL, on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeZone.setRawOffset(int)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeZone_setRawOffset()")
int32_t dxfg_TimeZone_setRawOffset(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone, int32_t offsetMillis);

/**
 * Returns the time period of zero.
 *
 * @deprecated Use dxfg_v4_TimePeriod_ZERO().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The time period (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimePeriod.ZERO`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimePeriod_ZERO()")
dxfg_time_period_t *dxfg_TimePeriod_ZERO(graal_isolatethread_t *thread);

/**
 * Returns the time period of "infinity": its time is the maximum value of `int64_t` milliseconds.
 *
 * @deprecated Use dxfg_v4_TimePeriod_UNLIMITED().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The time period (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimePeriod.UNLIMITED`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimePeriod_UNLIMITED()")
dxfg_time_period_t *dxfg_TimePeriod_UNLIMITED(graal_isolatethread_t *thread);

/**
 * Returns the time period of the given number of milliseconds.
 *
 * @deprecated Use dxfg_v4_TimePeriod_valueOf().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] value The value in milliseconds; a negative value is an error.
 * @return The time period (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimePeriod.valueOf(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimePeriod_valueOf()")
dxfg_time_period_t *dxfg_TimePeriod_valueOf(graal_isolatethread_t *thread, int64_t value);

/**
 * Returns the time period that the string represents. The string is an ISO8601 duration with these simplifications:
 * the letters are case insensitive; the letters "P" and "T" can be omitted; the letter "S" can be omitted, then the
 * last number is the seconds; the seconds can be fractional (milliseconds precision); an omitted part is zero; the
 * string "inf" is the unlimited period.
 *
 * @deprecated Use dxfg_v4_TimePeriod_valueOfString().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] value The string, e.g. "PT1H30M", "1d" or "0.5".
 * @return The time period (release it with dxfg_JavaObjectHandler_release()), or NULL on error (e.g. the string
 * cannot be parsed). Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimePeriod.valueOf(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimePeriod_valueOfString()")
dxfg_time_period_t *dxfg_TimePeriod_valueOf2(graal_isolatethread_t *thread, const char *value);

/**
 * Returns the time period in milliseconds.
 *
 * @deprecated Use dxfg_v4_TimePeriod_getTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timePeriod The time period.
 * @return The time period in milliseconds, or -1 on error (a time period is never negative). Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimePeriod.getTime()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimePeriod_getTime()")
int64_t dxfg_TimePeriod_getTime(graal_isolatethread_t *thread, dxfg_time_period_t *timePeriod);

/**
 * Returns the time period in whole seconds (rounded down). A period longer than about 68 years, e.g.
 * dxfg_TimePeriod_UNLIMITED(), does not fit into `int32_t`: the result is truncated to its lower 32 bits.
 *
 * @deprecated Use dxfg_v4_TimePeriod_getSeconds().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timePeriod The time period.
 * @return The time period in seconds, or INT32_MIN (-2147483648) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimePeriod.getSeconds()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimePeriod_getSeconds()")
int32_t dxfg_TimePeriod_getSeconds(graal_isolatethread_t *thread, dxfg_time_period_t *timePeriod);

/**
 * Returns the time period in nanoseconds. A period longer than about 292 years, e.g. dxfg_TimePeriod_UNLIMITED(),
 * overflows `int64_t` and gives a wrong value, which may be negative (-1000000 for the unlimited period); use
 * dxfg_TimePeriod_getTime() for such periods.
 *
 * @deprecated Use dxfg_v4_TimePeriod_getNanos().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timePeriod The time period.
 * @return The time period in nanoseconds, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimePeriod.getNanos()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimePeriod_getNanos()")
int64_t dxfg_TimePeriod_getNanos(graal_isolatethread_t *thread, dxfg_time_period_t *timePeriod);

/**
 * Returns the time format of the default time zone of the Java virtual machine.
 *
 * @deprecated Use dxfg_v4_TimeFormat_DEFAULT().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The time format (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.DEFAULT`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_DEFAULT()")
dxfg_time_format_t *dxfg_TimeFormat_DEFAULT(graal_isolatethread_t *thread);

/**
 * Returns the time format of the GMT time zone.
 *
 * @deprecated Use dxfg_v4_TimeFormat_GMT().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The time format (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.GMT`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_GMT()")
dxfg_time_format_t *dxfg_TimeFormat_GMT(graal_isolatethread_t *thread);

/**
 * Returns the time format of the time zone. The format keeps a copy of the time zone: later changes of the time zone
 * do not affect it.
 *
 * @deprecated Use dxfg_v4_TimeFormat_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone; NULL is an error.
 * @return The time format (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.getInstance(TimeZone)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_getInstance()")
dxfg_time_format_t *dxfg_TimeFormat_getInstance(graal_isolatethread_t *thread, dxfg_time_zone_t *timeZone);

/**
 * Returns the time format that also includes the time zone into the string of dxfg_TimeFormat_format().
 *
 * @deprecated Use dxfg_v4_TimeFormat_withTimeZone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @return The time format (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.withTimeZone()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_withTimeZone()")
dxfg_time_format_t *dxfg_TimeFormat_withTimeZone(graal_isolatethread_t *thread, dxfg_time_format_t *timeFormat);

/**
 * Returns the time format that also includes the milliseconds into the string of dxfg_TimeFormat_format().
 *
 * @deprecated Use dxfg_v4_TimeFormat_withMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @return The time format (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.withMillis()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_withMillis()")
dxfg_time_format_t *dxfg_TimeFormat_withMillis(graal_isolatethread_t *thread, dxfg_time_format_t *timeFormat);

/**
 * Returns the time format that produces the full ISO8601 string of `yyyy-MM-dd'T'HH:mm:ss.SSSX`.
 *
 * @deprecated Use dxfg_v4_TimeFormat_asFullIso().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @return The time format (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.asFullIso()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_asFullIso()")
dxfg_time_format_t *dxfg_TimeFormat_asFullIso(graal_isolatethread_t *thread, dxfg_time_format_t *timeFormat);

/**
 * Parses the date and the time from the string in an ISO 8601 format. The string can be:
 * - `0`: the zero time;
 * - `-<time period>`: the current time minus the time period (see dxfg_TimePeriod_valueOf2());
 * - a number of milliseconds since 1970-01-01T00:00:00Z: at least 9 digits;
 * - `<date>[(T|t|-| )<time>][<time zone>]` or `[T|t]<time>[<time zone>]`: the date is `yyyy-MM-dd` or `yyyyMMdd`,
 *   the time is `HH:mm:ss[.sss]`, `HHmmss[.sss]`, `HH:mm`, `HHmm` or `HH` (a missing time is 00:00:00, a missing date
 *   is the current date), the time zone is `[+-]HH:mm`, `[+-]HHmm`, `[+-]HH`, `Z` or a name of a time zone. A string
 *   without a time zone is in the time zone of the time format.
 *
 * @deprecated Use dxfg_v4_TimeFormat_parse().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[in] value The string to parse; NULL is an error.
 * @return The number of milliseconds since 1970-01-01T00:00:00Z, or -1 on error. -1 is also the time of
 * 1969-12-31T23:59:59.999Z: use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 *
 * Java: `TimeFormat.parse(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_parse()")
int64_t dxfg_TimeFormat_parse(graal_isolatethread_t *thread, dxfg_time_format_t *timeFormat, const char *value);

/**
 * Converts the time to a string like `yyyyMMdd-HHmmss` in the time zone of the time format. The time zone
 * (RFC 822, e.g. `+0300`) and the milliseconds (`.sss`) are included if the format was obtained with
 * dxfg_TimeFormat_withTimeZone() and dxfg_TimeFormat_withMillis(). The time 0 is formatted as "0". If the year of the
 * time is outside of 0001 to 9999, the result is the decimal number of the milliseconds.
 *
 * @deprecated Use dxfg_v4_TimeFormat_format().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[in] value The time: the number of milliseconds since 1970-01-01T00:00:00Z.
 * @return The string (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.format(long)`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_format()")
const char *dxfg_TimeFormat_format(graal_isolatethread_t *thread, dxfg_time_format_t *timeFormat, int64_t value);

/**
 * Returns the time zone of the time format.
 *
 * @deprecated Use dxfg_v4_TimeFormat_getTimeZone().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @return A copy of the time zone (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `TimeFormat.getTimeZone()`.
 */
DXFG_DEPRECATED("use dxfg_v4_TimeFormat_getTimeZone()")
dxfg_time_zone_t *dxfg_TimeFormat_getTimeZone(graal_isolatethread_t *thread, dxfg_time_format_t *timeFormat);

/**
 * Creates the authorization token from the string of the scheme and the value separated by a space, e.g.
 * `Bearer <token>` or `Basic <Base64 of user:password>`. The value of the "Basic" scheme must be the Base64-encoded
 * `user:password`.
 *
 * @deprecated Use dxfg_v4_AuthToken_valueOf().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The string; NULL is an error.
 * @return The token (release it with dxfg_JavaObjectHandler_release()), or NULL on error (e.g. the string is
 * malformed). Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.valueOf(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_valueOf()")
dxfg_auth_token_t *dxfg_AuthToken_valueOf(graal_isolatethread_t *thread, const char *string);

/**
 * Creates the "Basic" authorization token (RFC 2617) from the string `user:password`. The user name and the password
 * can be empty.
 *
 * @deprecated Use dxfg_v4_AuthToken_createBasicToken().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] userPassword The string `user:password`; a string without a colon or NULL is an error.
 * @return The token (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.createBasicToken(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_createBasicToken()")
dxfg_auth_token_t *dxfg_AuthToken_createBasicToken(graal_isolatethread_t *thread, const char *userPassword);

/**
 * Creates the "Basic" authorization token (RFC 2617) from the user name and the password. They can be empty.
 *
 * @deprecated Use dxfg_v4_AuthToken_createBasicTokenWithUserAndPassword().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user The user name; NULL is an error.
 * @param[in] password The password; NULL is an error.
 * @return The token (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.createBasicToken(String, String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_createBasicTokenWithUserAndPassword()")
dxfg_auth_token_t *dxfg_AuthToken_createBasicToken2(graal_isolatethread_t *thread, const char *user,
                                                    const char *password);

/**
 * Creates the "Basic" authorization token (RFC 2617) from the user name and the password, or returns NULL if both of
 * them are NULL or empty (it is not an error). A NULL user name or password is taken as empty.
 *
 * @deprecated Use dxfg_v4_AuthToken_createBasicTokenOrNull().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user The user name.
 * @param[in] password The password.
 * @return The token (release it with dxfg_JavaObjectHandler_release()), or NULL if the user name and the password are
 * both NULL or empty, or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.createBasicTokenOrNull(String, String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_createBasicTokenOrNull()")
dxfg_auth_token_t *dxfg_AuthToken_createBasicTokenOrNull(graal_isolatethread_t *thread, const char *user,
                                                         const char *password);

/**
 * Creates the "Bearer" authorization token (RFC 6750) from the access token.
 *
 * @deprecated Use dxfg_v4_AuthToken_createBearerToken().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] token The access token; NULL or an empty string is an error.
 * @return The token (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.createBearerToken(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_createBearerToken()")
dxfg_auth_token_t *dxfg_AuthToken_createBearerToken(graal_isolatethread_t *thread, const char *token);

/**
 * Creates the "Bearer" authorization token (RFC 6750) from the access token, or returns NULL if the access token is
 * NULL or empty (it is not an error).
 *
 * @deprecated Use dxfg_v4_AuthToken_createBearerTokenOrNull().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] token The access token.
 * @return The token (release it with dxfg_JavaObjectHandler_release()), or NULL if the access token is NULL or empty,
 * or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.createBearerTokenOrNull(String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_createBearerTokenOrNull()")
dxfg_auth_token_t *dxfg_AuthToken_createBearerTokenOrNull(graal_isolatethread_t *thread, const char *token);

/**
 * Creates the authorization token with a custom scheme and value.
 *
 * @deprecated Use dxfg_v4_AuthToken_createCustomToken().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] scheme The scheme; NULL or an empty string is an error.
 * @param[in] value The value; NULL or an empty string is an error.
 * @return The token (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.createCustomToken(String, String)`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_createCustomToken()")
dxfg_auth_token_t *dxfg_AuthToken_createCustomToken(graal_isolatethread_t *thread, const char *scheme,
                                                    const char *value);

/**
 * Returns the value of the HTTP authorization header: the scheme and the value separated by a space.
 *
 * @deprecated Use dxfg_v4_AuthToken_getHttpAuthorization().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @return The string (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.getHttpAuthorization()`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_getHttpAuthorization()")
const char *dxfg_AuthToken_getHttpAuthorization(graal_isolatethread_t *thread, dxfg_auth_token_t *authToken);

/**
 * Returns the user name of the token.
 *
 * @deprecated Use dxfg_v4_AuthToken_getUser().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @return The user name (release it with dxfg_String_release()), or NULL if the token has no user name (it is not a
 * "Basic" token), or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.getUser()`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_getUser()")
const char *dxfg_AuthToken_getUser(graal_isolatethread_t *thread, dxfg_auth_token_t *authToken);

/**
 * Returns the password of the token.
 *
 * @deprecated Use dxfg_v4_AuthToken_getPassword().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @return The password (release it with dxfg_String_release()), or NULL if the token has no password (it is not a
 * "Basic" token), or on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.getPassword()`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_getPassword()")
const char *dxfg_AuthToken_getPassword(graal_isolatethread_t *thread, dxfg_auth_token_t *authToken);

/**
 * Returns the authorization scheme of the token, e.g. "Basic" or "Bearer".
 *
 * @deprecated Use dxfg_v4_AuthToken_getScheme().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @return The scheme (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.getScheme()`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_getScheme()")
const char *dxfg_AuthToken_getScheme(graal_isolatethread_t *thread, dxfg_auth_token_t *authToken);

/**
 * Returns the value of the token: the access token of the "Bearer" scheme (RFC 6750) or the Base64-encoded
 * `user:password` of the "Basic" scheme (RFC 2617).
 *
 * @deprecated Use dxfg_v4_AuthToken_getValue().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @return The value (release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 *
 * Java: `AuthToken.getValue()`.
 */
DXFG_DEPRECATED("use dxfg_v4_AuthToken_getValue()")
const char *dxfg_AuthToken_getValue(graal_isolatethread_t *thread, dxfg_auth_token_t *authToken);

/** @} */ // end of Javac

#ifdef __cplusplus
}
#endif

#endif // DXFG_JAVAC_H
