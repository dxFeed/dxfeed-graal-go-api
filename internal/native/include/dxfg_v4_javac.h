// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_javac.h and the v4 names of its types (see the group @ref V4_Javac).
 */

#ifndef DXFG_V4_JAVAC_H
#define DXFG_V4_JAVAC_H

#include "dxfg_javac.h"
#include "dxfg_v4_common.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Javac Javac (v4)
 *  @brief The v4 functions and type names of the group @ref Javac. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_java_object_handler */
typedef dxfg_java_object_handler dxfg_v4_java_object_handler;

/** @copydoc dxfg_java_object_t */
typedef dxfg_java_object_t dxfg_v4_java_object_t;

/** @copydoc dxfg_list */
typedef dxfg_list dxfg_v4_list;

/** @copydoc dxfg_java_object_handler_list */
typedef dxfg_java_object_handler_list dxfg_v4_java_object_handler_list;

/** @copydoc dxfg_string_list */
typedef dxfg_string_list dxfg_v4_string_list;

/** @copydoc dxfg_executor_t */
typedef dxfg_executor_t dxfg_v4_executor_t;

/** @copydoc dxfg_input_stream_t */
typedef dxfg_input_stream_t dxfg_v4_input_stream_t;

/** @copydoc dxfg_time_period_t */
typedef dxfg_time_period_t dxfg_v4_time_period_t;

/** @copydoc dxfg_time_format_t */
typedef dxfg_time_format_t dxfg_v4_time_format_t;

/** @copydoc dxfg_time_zone_t */
typedef dxfg_time_zone_t dxfg_v4_time_zone_t;

/** @copydoc dxfg_auth_token_t */
typedef dxfg_auth_token_t dxfg_v4_auth_token_t;

/** @copydoc dxfg_rounding_mode_t */
typedef dxfg_rounding_mode_t dxfg_v4_rounding_mode_t;

/** @copydoc dxfg_finalize_function */
typedef dxfg_finalize_function dxfg_v4_finalize_function;

/**
 * @brief A list of bytes that a function returned: the array holds the bytes themselves.
 *
 * The functions that write data to memory return it, e.g. dxfg_v4_InstrumentProfileWriter_writeToBytes(). Release it
 * with dxfg_v4_byte_list_release(). The bytes can be read back with dxfg_v4_ByteArrayInputStream_new().
 * @since 3.13.0
 */
typedef struct dxfg_v4_byte_list {
    /** The number of bytes; 0 for an empty list. */
    int32_t size;
    /** The array of `size` bytes; NULL for an empty list. */
    char *elements;
} dxfg_v4_byte_list;

/**
 * Registers a function that is called with the user data after the Java object is garbage collected.
 *
 * The Java object is not collected while a handle of it exists, so the function is called after all the handles of
 * the object have been released and nothing else in Java refers to the object. The collected objects are checked
 * about once a second. The function is called by a daemon Java thread of the SDK (see dxfg_v4_finalize_function), once
 * for each registration.
 *
 * It does not stop the callbacks of a listener: see dxfg_v4_NativeListener_deactivate() for that.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] handler The handle of the object to watch.
 * @param[in] finalize The function to call.
 * @param[in] user_data The user data to pass to the function; the SDK does not free it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Object_finalize(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *handler,
                                dxfg_v4_finalize_function finalize, void *user_data);

/**
 * Frees the memory of the handle structure and releases the reference to the Java object. Do not use the handle after
 * that. Passing NULL does nothing.
 *
 * If the handle is the last handle (see dxfg_v4_JavaObjectHandler_clone()) of a listener created by a
 * `dxfg_*Listener_new` function, the listener is deactivated (see dxfg_v4_NativeListener_deactivate()): the function
 * waits for the listener's callbacks that are in progress on other threads and after it returns the listener's
 * function is not called anymore, so its user data can be freed. Remove the listener (e.g.
 * dxfg_v4_DXFeedSubscription_removeEventListener()) before releasing it, otherwise it stays added and does nothing.
 *
 * Do not hold the locks that the listener's function may take while releasing the listener: it may deadlock.
 * Releasing a listener from its own callback is allowed (the function does not wait for that call).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] handler The handle to release (e.g. `&executor->handler` for a dxfg_v4_executor_t).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_JavaObjectHandler_release(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *handler);

/**
 * Deactivates a listener created by a `dxfg_*Listener_new` function (e.g. dxfg_v4_DXFeedEventListener_new(),
 * dxfg_v4_PropertyChangeListener_new()) regardless of the number of its handles: waits for the listener's callbacks
 * that are in progress on other threads, after that the listener's function is not called anymore, so its user data
 * can be freed. The listener stays added where it was added, but does nothing, and cannot be activated again.
 *
 * Releasing the last handle of the listener with dxfg_v4_JavaObjectHandler_release() deactivates it as well.
 * Do not hold the locks that the listener's function may take while deactivating the listener: it may deadlock.
 * Deactivating a listener from its own callback is allowed (the function does not wait for that call).
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error
 * (e.g. the object is not a listener created by a `dxfg_*Listener_new` function).
 * Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_NativeListener_deactivate(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *listener);

/**
 * Deactivates all the listeners of this isolate created by the `dxfg_*Listener_new` functions, as
 * dxfg_v4_NativeListener_deactivate() does for one: waits for the listeners' callbacks that are in progress on other
 * threads, after that no listener's function is called anymore, so the user data of the listeners can be freed. The
 * listeners created afterwards are inactive from the start. The deactivation cannot be undone for the isolate.
 *
 * Then the function closes all the event queues of this isolate (dxfg_v4_EventQueue_close()) and waits for the polls
 * in progress on other threads: they return #DXFG_EVENT_QUEUE_CLOSED or the events they have already retrieved. The
 * queues created afterwards are closed from the start.
 *
 * The handles of the listeners stay valid and must be released as usual. Use dxfg_v4_system_shutdown() to deactivate
 * the listeners and then close the endpoints. Only the listeners are deactivated: the handlers of
 * dxfg_v4_Promise_whenDone(), dxfg_v4_Promise_whenDoneAsync(), the finalizers of dxfg_v4_Object_finalize() and the
 * filters of dxfg_v4_SessionFilter_new() and dxfg_v4_DayFilter_new() are still called.
 *
 * Calling it from a listener's callback is allowed (the function does not wait for that call), but not from the
 * callbacks of two threads at once: each would wait for the other. Do not hold the locks that the listeners' functions
 * may take while calling it, nor call it from an exit handler of a DLL (`DLL_PROCESS_DETACH`, a destructor of a
 * static object of a DLL): a callback that waits for the loader lock never finishes.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_NativeListener_deactivateAll(graal_isolatethread_t *thread);

/**
 * Creates another handle of the Java object: allocates the memory of the new handle structure. Both handles refer to
 * the same object and are released independently: the object stays alive until all of its handles are released.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] javaObjectHandler The handle to clone; it stays valid.
 * @param[out] javaObjectHandlerClone The pointer to which the pointer to the new handle is written; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the object clone handle.
 */
int32_t dxfg_v4_JavaObjectHandler_clone(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *javaObjectHandler,
                                        DXFG_OUT dxfg_v4_java_object_handler **javaObjectHandlerClone);

/**
 * Releases an array of handles (or simple structures containing only a handle) that a function returned as one block
 * of consecutive structures: releases every handle of the array and frees the memory of the array.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] handlersArray The array of handles. If it is NULL or `size` is 0, nothing is done.
 * @param[in] size The number of handles in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_JavaObjectHandler_array_release(graal_isolatethread_t *thread,
                                                const dxfg_v4_java_object_handler *handlersArray, int32_t size);

/**
 * Creates a new Java Object.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] object The pointer to which the pointer to the handle of the new object is written; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * Use dxfg_v4_JavaObjectHandler_release() to free the object handle.
 *
 * Java: `Object()`.
 */
int32_t dxfg_v4_Object_new(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_java_object_handler **object);

/**
 * Returns the string representation of the Java object. It works for the handle of any Java object, including the
 * typed handles (e.g. dxfg_v4_time_zone_t, dxfg_v4_auth_token_t).
 *
 * Java: `Object.toString()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the object. A NULL handle is an error.
 * @param[out] string The pointer to the variable that receives the UTF-8 string, NUL-terminated (release it with
 * dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Object_toString(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *object,
                                DXFG_OUT const char **string);

/**
 * Checks whether the Java objects are equal: the handles of the same object, and the handles of the objects that are
 * equal in the sense of the Java `equals` method of the first object. Two NULL handles are equal, a NULL handle and
 * a non-NULL handle are not.
 *
 * Java: `Objects.equals(Object, Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the first object.
 * @param[in] other The handle of the second object.
 * @param[out] equal The pointer to the variable that receives 1 if the objects are equal, 0 if they are not; must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Object_equals(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *object,
                              dxfg_v4_java_object_handler *other, DXFG_OUT int32_t *equal);

/**
 * Returns the hash code of the Java object: the result of its Java `hashCode` method, or 0 for a NULL handle.
 *
 * Java: `Objects.hashCode(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the object.
 * @param[out] hashCode The pointer to the variable that receives the hash code; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Object_hashCode(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *object,
                                DXFG_OUT int32_t *hashCode);

/**
 * Compares the Java objects by their natural order. The objects must implement the Java `Comparable` interface and
 * be comparable to each other, otherwise the function fails. The handles of the same object are equal; a NULL handle
 * is less than a non-NULL handle.
 *
 * Java: `Comparable.compareTo(Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] object The handle of the first object.
 * @param[in] other The handle of the second object.
 * @param[out] comparison The pointer to the variable that receives a negative value if the first object is less than
 * the second one, 0 if they are equal, a positive value if it is greater; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Comparable_compareTo(graal_isolatethread_t *thread, dxfg_v4_java_object_handler *object,
                                     dxfg_v4_java_object_handler *other, DXFG_OUT int32_t *comparison);

/**
 * Releases a list of handles that a function returned: releases every handle of the list, frees the array of the
 * elements and the list itself.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_java_object_handler_list_release(graal_isolatethread_t *thread, dxfg_v4_java_object_handler_list *list);

/**
 * Frees a string that a function returned (e.g. dxfg_v4_Object_toString()). Do not pass a string that is a field of
 * a structure, such as a string of an event: it is freed with the structure.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The string to free. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_String_release(graal_isolatethread_t *thread, const char *string);

/**
 * Releases a list of strings that a function returned: frees every string of the list, the array of the strings and
 * the list itself.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] stringList The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_string_list_release(graal_isolatethread_t *thread, dxfg_v4_string_list *stringList);

/**
 * Creates an executor with a fixed number of threads that runs the submitted tasks. The threads are daemon threads.
 * The SDK has no function that shuts the executor down.
 *
 * Read the "Threads and locks" section of
 * https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedSubscription.html for how the executor is used by the
 * listeners.
 *
 * Java: `Executors.newFixedThreadPool(int, ThreadFactory)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] nThreads The number of threads of the pool; must be positive.
 * @param[in] nameThreads The prefix of the names of the threads: they are named `<nameThreads>-1`,
 * `<nameThreads>-2`, and so on.
 * @param[out] executor The pointer to the variable that receives the executor (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Executors_newFixedThreadPool(graal_isolatethread_t *thread, int nThreads, const char *nameThreads,
                                             DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * Creates an executor that runs the submitted tasks, also with a delay or periodically. The threads are daemon
 * threads. The SDK has no function that shuts the executor down.
 *
 * See dxfg_v4_Executors_newFixedThreadPool() for how the listeners use the executor.
 *
 * Java: `Executors.newScheduledThreadPool(int, ThreadFactory)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] nThreads The number of threads to keep in the pool, even if they are idle; must not be negative.
 * @param[in] nameThreads The prefix of the names of the threads: they are named `<nameThreads>-1`,
 * `<nameThreads>-2`, and so on.
 * @param[out] executor The pointer to the variable that receives the executor (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Executors_newScheduledThreadPool(graal_isolatethread_t *thread, int nThreads, const char *nameThreads,
                                                 DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * Creates an executor that does not run the tasks by itself: a submitted task is added to a queue and runs when
 * dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks() is called, on the thread that calls it. Use it
 * to have the calls of the listeners made on a thread of your own.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] executor The pointer to the variable that receives the executor (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_new(graal_isolatethread_t *thread,
                                                        DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * Runs the tasks that are in the queue of the executor, on the calling thread, until the queue is empty (the tasks
 * added by the running tasks are run as well). If a task throws an exception, the function fails and the rest of
 * the tasks stay in the queue.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] executor The executor created by dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_new(); any other executor
 * (e.g. the executor of an event queue, dxfg_v4_EventQueue_getExecutor()) makes the function fail with
 * `IllegalArgumentException`.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks(graal_isolatethread_t *thread,
                                                                           dxfg_v4_executor_t *executor);

/**
 * Creates an input stream that reads the bytes of the memory block. The bytes are copied: the caller keeps the memory
 * block and may free it after the function returns.
 *
 * Java: `ByteArrayInputStream(byte[])`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] bytes The bytes of the stream; must point to at least `size` bytes.
 * @param[in] size The number of bytes; must not be negative.
 * @param[out] inputStream The pointer to the variable that receives the input stream (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_ByteArrayInputStream_new(graal_isolatethread_t *thread, const char *bytes, int32_t size,
                                         DXFG_OUT dxfg_v4_input_stream_t **inputStream);

/**
 * Releases a list of bytes that a function returned (e.g. dxfg_v4_InstrumentProfileWriter_writeToBytes()): frees the
 * array of the bytes and the list itself.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] list The list to release. If it is NULL, nothing is done.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_byte_list_release(graal_isolatethread_t *thread, dxfg_v4_byte_list *list);

/**
 * Requests a garbage collection. It is only for testing, e.g. of the functions registered by dxfg_v4_Object_finalize().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `System.gc()`.
 */
int32_t dxfg_v4_gc(graal_isolatethread_t *thread);

/**
 * Always fails: throws a Java exception (a RuntimeException "recursion" with a cause "cause", thrown several calls
 * deep, so that its stack trace has several frames). It is only for testing of
 * dxfg_v4_get_and_clear_thread_exception().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] object Not written: the function always fails; must not be NULL.
 * @return #DXFG_EXECUTE_FAIL (-1). Use dxfg_v4_get_and_clear_thread_exception() to get the exception.
 */
int32_t dxfg_v4_throw_exception(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_java_object_handler **object);

/**
 * Returns the time zone for the ID. If the ID is not known, the result is the GMT zone.
 *
 * Java: `TimeZone.getTimeZone(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ID The ID of the time zone, e.g. "America/New_York" or "GMT+03:00".
 * @param[out] timeZone The pointer to the variable that receives the time zone (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getTimeZone(graal_isolatethread_t *thread, const char *ID,
                                     DXFG_OUT dxfg_v4_time_zone_t **timeZone);

/**
 * Returns the default time zone of the Java virtual machine.
 *
 * Java: `TimeZone.getDefault()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] timeZone The pointer to the variable that receives the time zone (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getDefault(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_time_zone_t **timeZone);

/**
 * Returns the ID of the time zone.
 *
 * Java: `TimeZone.getID()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[out] id The pointer to the variable that receives the ID (release it with dxfg_v4_String_release()); must not
 * be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getID(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone, DXFG_OUT const char **id);

/**
 * Returns the long name of the time zone (for the standard time) in the default locale of the Java virtual machine.
 *
 * Java: `TimeZone.getDisplayName()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[out] displayName The pointer to the variable that receives the name (release it with
 * dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getDisplayName(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                        DXFG_OUT const char **displayName);

/**
 * Returns the name of the time zone in the default locale of the Java virtual machine, for the standard or the
 * daylight saving time, in the long or the short style.
 *
 * Java: `TimeZone.getDisplayName(boolean, int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] daylight 1 for the name of the daylight saving time, any other value for the name of the standard time.
 * @param[in] style 0 for the short name (`TimeZone.SHORT`), 1 for the long name (`TimeZone.LONG`); any other value is
 * an error.
 * @param[out] displayName The pointer to the variable that receives the name (release it with
 * dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getDisplayNameWithStyle(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                                 int32_t daylight, int32_t style, DXFG_OUT const char **displayName);

/**
 * Returns the amount of time to be added to the local standard time to get the local wall clock time, in
 * milliseconds. It is 0 if the time zone does not observe the daylight saving time.
 *
 * Java: `TimeZone.getDSTSavings()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[out] dstSavings The pointer to the variable that receives the savings in milliseconds; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getDSTSavings(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                       DXFG_OUT int32_t *dstSavings);

/**
 * Checks whether the time zone uses the daylight saving time.
 *
 * Java: `TimeZone.useDaylightTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[out] useDaylightTime The pointer to the variable that receives 1 if the time zone uses the daylight saving
 * time, 0 if it does not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_useDaylightTime(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                         DXFG_OUT int32_t *useDaylightTime);

/**
 * Checks whether the time zone is currently in the daylight saving time, or a transition from the standard time to the
 * daylight saving time occurs at any future time.
 *
 * Java: `TimeZone.observesDaylightTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[out] observesDaylightTime The pointer to the variable that receives 1 if the time zone observes the daylight
 * saving time, 0 if it does not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_observesDaylightTime(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                              DXFG_OUT int32_t *observesDaylightTime);

/**
 * Returns the offset of the time zone from UTC at the date, in milliseconds. If the daylight saving time is in effect
 * at the date, the offset includes the saving.
 *
 * Java: `TimeZone.getOffset(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] date The date: the number of milliseconds since 1970-01-01T00:00:00Z.
 * @param[out] offset The pointer to the variable that receives the offset in milliseconds; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getOffset(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone, int64_t date,
                                   DXFG_OUT int32_t *offset);

/**
 * Returns the offset of the time zone from UTC at the local date, in milliseconds, adjusted for the daylight saving
 * time.
 *
 * Java: `TimeZone.getOffset(int, int, int, int, int, int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] era The era of the date: 1 for AD (`GregorianCalendar.AD`) or 0 for BC (`GregorianCalendar.BC`).
 * @param[in] year The year in the era.
 * @param[in] month The month of the year, 0-based: 0 for January.
 * @param[in] day The day of the month.
 * @param[in] dayOfWeek The day of the week: 1 for Sunday (`Calendar.SUNDAY`) to 7 for Saturday.
 * @param[in] milliseconds The milliseconds in the day, in the standard local time.
 * @param[out] offset The pointer to the variable that receives the offset in milliseconds; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g. a
 * parameter is out of range). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getOffsetByCalendarFields(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                                   int32_t era, int32_t year, int32_t month, int32_t day,
                                                   int32_t dayOfWeek, int32_t milliseconds, DXFG_OUT int32_t *offset);

/**
 * Returns the offset of the time zone from UTC, in milliseconds, that does not include the daylight saving time.
 *
 * Java: `TimeZone.getRawOffset()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[out] rawOffset The pointer to the variable that receives the raw offset in milliseconds; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_getRawOffset(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                      DXFG_OUT int32_t *rawOffset);

/**
 * Checks whether the time zone has the same rules and offset as the other time zone (they may differ only in the ID).
 *
 * Java: `TimeZone.hasSameRules(TimeZone)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] other The other time zone.
 * @param[out] hasSameRules The pointer to the variable that receives 1 if the rules are the same, 0 if they are not
 * (also if `other` is NULL); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_hasSameRules(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                      dxfg_v4_time_zone_t *other, DXFG_OUT int32_t *hasSameRules);

/**
 * Checks whether the date is in the daylight saving time of the time zone.
 *
 * Java: `TimeZone.inDaylightTime(Date)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] date The date: the number of milliseconds since 1970-01-01T00:00:00Z.
 * @param[out] inDaylightTime The pointer to the variable that receives 1 if the date is in the daylight saving time, 0
 * if it is not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeZone_inDaylightTime(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone, int64_t date,
                                        DXFG_OUT int32_t *inDaylightTime);

/**
 * Sets the ID of the time zone. It changes the Java object, so all the handles of the time zone see the new ID.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] ID The new ID.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `TimeZone.setID(String)`.
 */
int32_t dxfg_v4_TimeZone_setID(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone, const char *ID);

/**
 * Sets the offset of the time zone from UTC that does not include the daylight saving time. It changes the Java
 * object, so all the handles of the time zone see the new offset.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone.
 * @param[in] offsetMillis The raw offset in milliseconds.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 *
 * Java: `TimeZone.setRawOffset(int)`.
 */
int32_t dxfg_v4_TimeZone_setRawOffset(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                      int32_t offsetMillis);

/**
 * Returns the time period of zero.
 *
 * Java: `TimePeriod.ZERO`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] timePeriod The pointer to the variable that receives the time period (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimePeriod_ZERO(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_time_period_t **timePeriod);

/**
 * Returns the time period of "infinity": its time is the maximum value of `int64_t` milliseconds.
 *
 * Java: `TimePeriod.UNLIMITED`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] timePeriod The pointer to the variable that receives the time period (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimePeriod_UNLIMITED(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_time_period_t **timePeriod);

/**
 * Returns the time period of the given number of milliseconds.
 *
 * Java: `TimePeriod.valueOf(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] value The value in milliseconds; a negative value is an error.
 * @param[out] timePeriod The pointer to the variable that receives the time period (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimePeriod_valueOf(graal_isolatethread_t *thread, int64_t value,
                                   DXFG_OUT dxfg_v4_time_period_t **timePeriod);

/**
 * Returns the time period that the string represents. The string is an ISO8601 duration with these simplifications:
 * the letters are case insensitive; the letters "P" and "T" can be omitted; the letter "S" can be omitted, then the
 * last number is the seconds; the seconds can be fractional (milliseconds precision); an omitted part is zero; the
 * string "inf" is the unlimited period.
 *
 * Java: `TimePeriod.valueOf(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] value The string, e.g. "PT1H30M", "1d" or "0.5".
 * @param[out] timePeriod The pointer to the variable that receives the time period (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g. the
 * string cannot be parsed). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimePeriod_valueOfString(graal_isolatethread_t *thread, const char *value,
                                         DXFG_OUT dxfg_v4_time_period_t **timePeriod);

/**
 * Returns the time period in milliseconds.
 *
 * Java: `TimePeriod.getTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timePeriod The time period.
 * @param[out] time The pointer to the variable that receives the time period in milliseconds; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimePeriod_getTime(graal_isolatethread_t *thread, dxfg_v4_time_period_t *timePeriod,
                                   DXFG_OUT int64_t *time);

/**
 * Returns the time period in whole seconds (rounded down). A period longer than about 68 years, e.g.
 * dxfg_v4_TimePeriod_UNLIMITED(), does not fit into `int32_t`: the result is truncated to its lower 32 bits.
 *
 * Java: `TimePeriod.getSeconds()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timePeriod The time period.
 * @param[out] seconds The pointer to the variable that receives the time period in seconds; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimePeriod_getSeconds(graal_isolatethread_t *thread, dxfg_v4_time_period_t *timePeriod,
                                      DXFG_OUT int32_t *seconds);

/**
 * Returns the time period in nanoseconds. A period longer than about 292 years, e.g. dxfg_v4_TimePeriod_UNLIMITED(),
 * overflows `int64_t` and gives a wrong value, which may be negative (-1000000 for the unlimited period); use
 * dxfg_v4_TimePeriod_getTime() for such periods.
 *
 * Java: `TimePeriod.getNanos()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timePeriod The time period.
 * @param[out] nanos The pointer to the variable that receives the time period in nanoseconds; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimePeriod_getNanos(graal_isolatethread_t *thread, dxfg_v4_time_period_t *timePeriod,
                                    DXFG_OUT int64_t *nanos);

/**
 * Returns the time format of the default time zone of the Java virtual machine.
 *
 * Java: `TimeFormat.DEFAULT`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] timeFormat The pointer to the variable that receives the time format (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_DEFAULT(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_time_format_t **timeFormat);

/**
 * Returns the time format of the GMT time zone.
 *
 * Java: `TimeFormat.GMT`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] timeFormat The pointer to the variable that receives the time format (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_GMT(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_time_format_t **timeFormat);

/**
 * Returns the time format of the time zone. The format keeps a copy of the time zone: later changes of the time zone
 * do not affect it.
 *
 * Java: `TimeFormat.getInstance(TimeZone)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeZone The time zone; NULL is an error.
 * @param[out] timeFormat The pointer to the variable that receives the time format (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_getInstance(graal_isolatethread_t *thread, dxfg_v4_time_zone_t *timeZone,
                                       DXFG_OUT dxfg_v4_time_format_t **timeFormat);

/**
 * Returns the time format that also includes the time zone into the string of dxfg_v4_TimeFormat_format().
 *
 * Java: `TimeFormat.withTimeZone()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[out] formatWithTimeZone The pointer to the variable that receives the time format (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_withTimeZone(graal_isolatethread_t *thread, dxfg_v4_time_format_t *timeFormat,
                                        DXFG_OUT dxfg_v4_time_format_t **formatWithTimeZone);

/**
 * Returns the time format that also includes the milliseconds into the string of dxfg_v4_TimeFormat_format().
 *
 * Java: `TimeFormat.withMillis()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[out] formatWithMillis The pointer to the variable that receives the time format (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_withMillis(graal_isolatethread_t *thread, dxfg_v4_time_format_t *timeFormat,
                                      DXFG_OUT dxfg_v4_time_format_t **formatWithMillis);

/**
 * Returns the time format that produces the full ISO8601 string of `yyyy-MM-dd'T'HH:mm:ss.SSSX`.
 *
 * Java: `TimeFormat.asFullIso()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[out] fullIsoFormat The pointer to the variable that receives the time format (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_asFullIso(graal_isolatethread_t *thread, dxfg_v4_time_format_t *timeFormat,
                                     DXFG_OUT dxfg_v4_time_format_t **fullIsoFormat);

/**
 * Parses the date and the time from the string in an ISO 8601 format. The string can be:
 * - `0`: the zero time;
 * - `-<time period>`: the current time minus the time period (see dxfg_v4_TimePeriod_valueOfString());
 * - a number of milliseconds since 1970-01-01T00:00:00Z: at least 9 digits;
 * - `<date>[(T|t|-| )<time>][<time zone>]` or `[T|t]<time>[<time zone>]`: the date is `yyyy-MM-dd` or `yyyyMMdd`,
 *   the time is `HH:mm:ss[.sss]`, `HHmmss[.sss]`, `HH:mm`, `HHmm` or `HH` (a missing time is 00:00:00, a missing date
 *   is the current date), the time zone is `[+-]HH:mm`, `[+-]HHmm`, `[+-]HH`, `Z` or a name of a time zone. A string
 *   without a time zone is in the time zone of the time format.
 *
 * Java: `TimeFormat.parse(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[in] value The string to parse; NULL is an error.
 * @param[out] time The pointer to the variable that receives the number of milliseconds since 1970-01-01T00:00:00Z;
 * must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_parse(graal_isolatethread_t *thread, dxfg_v4_time_format_t *timeFormat, const char *value,
                                 DXFG_OUT int64_t *time);

/**
 * Converts the time to a string like `yyyyMMdd-HHmmss` in the time zone of the time format. The time zone
 * (RFC 822, e.g. `+0300`) and the milliseconds (`.sss`) are included if the format was obtained with
 * dxfg_v4_TimeFormat_withTimeZone() and dxfg_v4_TimeFormat_withMillis(). The time 0 is formatted as "0". If the year of
 * the time is outside of 0001 to 9999, the result is the decimal number of the milliseconds.
 *
 * Java: `TimeFormat.format(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[in] value The time: the number of milliseconds since 1970-01-01T00:00:00Z.
 * @param[out] string The pointer to the variable that receives the string (release it with dxfg_v4_String_release());
 * must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_format(graal_isolatethread_t *thread, dxfg_v4_time_format_t *timeFormat, int64_t value,
                                  DXFG_OUT const char **string);

/**
 * Returns the time zone of the time format.
 *
 * Java: `TimeFormat.getTimeZone()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] timeFormat The time format.
 * @param[out] timeZone The pointer to the variable that receives a copy of the time zone (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_TimeFormat_getTimeZone(graal_isolatethread_t *thread, dxfg_v4_time_format_t *timeFormat,
                                       DXFG_OUT dxfg_v4_time_zone_t **timeZone);

/**
 * Creates the authorization token from the string of the scheme and the value separated by a space, e.g.
 * `Bearer <token>` or `Basic <Base64 of user:password>`. The value of the "Basic" scheme must be the Base64-encoded
 * `user:password`.
 *
 * Java: `AuthToken.valueOf(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The string; NULL is an error.
 * @param[out] authToken The pointer to the variable that receives the token (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (e.g. the
 * string is malformed). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_valueOf(graal_isolatethread_t *thread, const char *string,
                                  DXFG_OUT dxfg_v4_auth_token_t **authToken);

/**
 * Creates the "Basic" authorization token (RFC 2617) from the string `user:password`. The user name and the password
 * can be empty.
 *
 * Java: `AuthToken.createBasicToken(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] userPassword The string `user:password`; a string without a colon or NULL is an error.
 * @param[out] authToken The pointer to the variable that receives the token (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_createBasicToken(graal_isolatethread_t *thread, const char *userPassword,
                                           DXFG_OUT dxfg_v4_auth_token_t **authToken);

/**
 * Creates the "Basic" authorization token (RFC 2617) from the user name and the password. They can be empty.
 *
 * Java: `AuthToken.createBasicToken(String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user The user name; NULL is an error.
 * @param[in] password The password; NULL is an error.
 * @param[out] authToken The pointer to the variable that receives the token (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_createBasicTokenWithUserAndPassword(graal_isolatethread_t *thread, const char *user,
                                                              const char *password,
                                                              DXFG_OUT dxfg_v4_auth_token_t **authToken);

/**
 * Creates the "Basic" authorization token (RFC 2617) from the user name and the password, or returns NULL if both of
 * them are NULL or empty (it is not an error). A NULL user name or password is taken as empty.
 *
 * Java: `AuthToken.createBasicTokenOrNull(String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user The user name.
 * @param[in] password The password.
 * @param[out] authToken The pointer to the variable that receives the token (release it with
 * dxfg_v4_JavaObjectHandler_release()), or NULL if the user name and the password are both NULL or empty; must not be
 * NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_createBasicTokenOrNull(graal_isolatethread_t *thread, const char *user, const char *password,
                                                 DXFG_OUT dxfg_v4_auth_token_t **authToken);

/**
 * Creates the "Bearer" authorization token (RFC 6750) from the access token.
 *
 * Java: `AuthToken.createBearerToken(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] token The access token; NULL or an empty string is an error.
 * @param[out] authToken The pointer to the variable that receives the token (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_createBearerToken(graal_isolatethread_t *thread, const char *token,
                                            DXFG_OUT dxfg_v4_auth_token_t **authToken);

/**
 * Creates the "Bearer" authorization token (RFC 6750) from the access token, or returns NULL if the access token is
 * NULL or empty (it is not an error).
 *
 * Java: `AuthToken.createBearerTokenOrNull(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] token The access token.
 * @param[out] authToken The pointer to the variable that receives the token (release it with
 * dxfg_v4_JavaObjectHandler_release()), or NULL if the access token is NULL or empty; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_createBearerTokenOrNull(graal_isolatethread_t *thread, const char *token,
                                                  DXFG_OUT dxfg_v4_auth_token_t **authToken);

/**
 * Creates the authorization token with a custom scheme and value.
 *
 * Java: `AuthToken.createCustomToken(String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] scheme The scheme; NULL or an empty string is an error.
 * @param[in] value The value; NULL or an empty string is an error.
 * @param[out] authToken The pointer to the variable that receives the token (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_createCustomToken(graal_isolatethread_t *thread, const char *scheme, const char *value,
                                            DXFG_OUT dxfg_v4_auth_token_t **authToken);

/**
 * Returns the value of the HTTP authorization header: the scheme and the value separated by a space.
 *
 * Java: `AuthToken.getHttpAuthorization()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @param[out] httpAuthorization The pointer to the variable that receives the string (release it with
 * dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_getHttpAuthorization(graal_isolatethread_t *thread, dxfg_v4_auth_token_t *authToken,
                                               DXFG_OUT const char **httpAuthorization);

/**
 * Returns the user name of the token.
 *
 * Java: `AuthToken.getUser()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @param[out] user The pointer to the variable that receives the user name (release it with dxfg_v4_String_release()),
 * or NULL if the token has no user name (it is not a "Basic" token); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_getUser(graal_isolatethread_t *thread, dxfg_v4_auth_token_t *authToken,
                                  DXFG_OUT const char **user);

/**
 * Returns the password of the token.
 *
 * Java: `AuthToken.getPassword()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @param[out] password The pointer to the variable that receives the password (release it with
 * dxfg_v4_String_release()), or NULL if the token has no password (it is not a "Basic" token); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_getPassword(graal_isolatethread_t *thread, dxfg_v4_auth_token_t *authToken,
                                      DXFG_OUT const char **password);

/**
 * Returns the authorization scheme of the token, e.g. "Basic" or "Bearer".
 *
 * Java: `AuthToken.getScheme()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @param[out] scheme The pointer to the variable that receives the scheme (release it with dxfg_v4_String_release());
 * must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_getScheme(graal_isolatethread_t *thread, dxfg_v4_auth_token_t *authToken,
                                    DXFG_OUT const char **scheme);

/**
 * Returns the value of the token: the access token of the "Bearer" scheme (RFC 6750) or the Base64-encoded
 * `user:password` of the "Basic" scheme (RFC 2617).
 *
 * Java: `AuthToken.getValue()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] authToken The token.
 * @param[out] value The pointer to the variable that receives the value (release it with dxfg_v4_String_release());
 * must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_AuthToken_getValue(graal_isolatethread_t *thread, dxfg_v4_auth_token_t *authToken,
                                   DXFG_OUT const char **value);

/** @} */ // end of V4_Javac

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_JAVAC_H
