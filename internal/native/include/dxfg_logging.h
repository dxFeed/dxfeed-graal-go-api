// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Logging functions and types declarations
 */

#ifndef DXFG_LOGGING_H
#define DXFG_LOGGING_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_catch_exception.h"
#include "dxfg_javac.h"
#include "graal_isolate.h"

/** @defgroup Logging Logging
 *  @brief Interception of the messages that the QD library and the SDK log, and the logging levels.
 *
 *  The messages are logged by the Java `com.devexperts.logging` system of the isolate. Instead of the console and the
 *  log files the messages can be given to a C function: create a listener with dxfg_logging_listener_new() and
 *  install it with dxfg_logging_set_listener(). The logging implementation is chosen when the first message is logged,
 *  so the property `log.className` (see dxfg_system_set_property()) must be set before that, as in the sample of
 *  dxfg_logging_listener_new().
 *
 *  The classes `InterceptableLogging` and `InterceptableLoggingListener` that the functions of this group call belong
 *  to the SDK (package `com.devexperts.logging`), not to the QD library.
 *
 *  With `InterceptableLogging` the level `log.level` (the listener or the console or the log file, `INFO` by default)
 *  and `err.level` (the err file) also set the level of the root logger of `java.util.logging` to the more verbose of
 *  the two, so that the messages below it are not created at all and the more verbose ones (`DEBUG`, `ALL`) are.
 *
 *  The last connection errors of the message connectors (dxfg_v4_MessageConnector_getLastError()) do not depend on the
 *  listener or on the levels: the SDK captures them with the default logging and with `InterceptableLogging`, also when
 *  the level is `OFF`. With another `log.className` nothing is captured, and a level set on QD's loggers of the
 *  connectors themselves wins (see @ref V4_MessageConnector).
 *
 *  The ownership rules of the handles are described in the group @ref Javac.
 *  @{
 */

/**
 * @brief The logging levels that can be used to control the logging output.
 *
 * The values correspond to the levels of `java.util.logging.Level` that the logging system uses.
 *
 * Java: `java.util.logging.Level`.
 */
typedef enum dxfg_logging_level_t {
    /** Indicates that all messages should be logged (`Level.ALL`). */
    DXFG_LOGGING_LEVEL_ALL = 0,

    /** Indicates a highly detailed tracing message (`Level.FINEST`). */
    DXFG_LOGGING_LEVEL_TRACE,

    /** A message level providing tracing debug information (`Level.FINE`). */
    DXFG_LOGGING_LEVEL_DEBUG,

    /** A message level for informational messages (`Level.INFO`). */
    DXFG_LOGGING_LEVEL_INFO,

    /** A message level indicating a potential problem (`Level.WARNING`). */
    DXFG_LOGGING_LEVEL_WARN,

    /** A message level indicating a serious failure (`Level.SEVERE`). */
    DXFG_LOGGING_LEVEL_ERROR,

    /** A special level that can be used to turn off logging (`Level.OFF`). */
    DXFG_LOGGING_LEVEL_OFF,
} dxfg_logging_level_t;

/**
 * @brief A handle of a Java InterceptableLoggingListener object, a wrapper of a listener that allows to intercept the
 * logged messages; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_logging_listener_new() and install it with dxfg_logging_set_listener(). The last handle of the
 * listener that is released deactivates the listener: its function is not called anymore.
 */
typedef struct dxfg_logging_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_logging_listener_t;

/**
 * @brief The function of a dxfg_logging_listener_t: called for every message that is logged while the listener is
 * installed and the level of the message is not below the logging level (see dxfg_logging_set_log_level()).
 *
 * The function is called synchronously by the Java thread that logs the message, which may be any thread of the
 * isolate (including the threads of the endpoints), not the thread that installed the listener. The calls are
 * serialized by a lock of the logging system: the function is never called concurrently, and while it runs the other
 * threads that log wait. So it should return quickly and should not block.
 *
 * The strings and the exception are allocated by the SDK for the call and released by the SDK after the function
 * returns: copy what is needed, do not keep the pointers and do not release them. A string or the exception is NULL
 * when the corresponding value is absent.
 *
 * Java: `InterceptableLoggingListener.onLog(Level, long, String, long, String, String, Throwable, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] level The level of the message. A Java level without a constant here is given as the next lower one:
 * `FINER` as #DXFG_LOGGING_LEVEL_TRACE, `CONFIG` as #DXFG_LOGGING_LEVEL_DEBUG (since 3.16.0; before, a call that
 * logged at such a level failed).
 * @param[in] timestamp The time of the message in milliseconds since the Unix epoch (UTC).
 * @param[in] thread_name The name of the Java thread that logged the message, a NUL-terminated UTF-8 string.
 * @param[in] thread_id The identifier of the Java thread that logged the message.
 * @param[in] logger_name The name of the logger, a NUL-terminated UTF-8 string (NULL if the logger has no name).
 * @param[in] message The text of the message without the time, the level and other parts of the layout, a
 * NUL-terminated UTF-8 string.
 * @param[in] exception The exception that is logged with the message, or NULL if there is none. See
 * dxfg_exception_t.
 * @param[in] formatted_message The whole line that the logging system would write: the message with the time, the
 * level, the thread and the exception formatted with the layout of the logging system, a NUL-terminated UTF-8
 * string.
 * @param[in] user_data The pointer that was passed to dxfg_logging_listener_new(); the SDK does not use it.
 */
typedef void (*dxfg_logging_listener_function_t)(graal_isolatethread_t *thread, dxfg_logging_level_t level,
                                                 int64_t timestamp, const char *thread_name, int64_t thread_id,
                                                 const char *logger_name, const char *message,
                                                 dxfg_exception_t *exception, const char *formatted_message,
                                                 void *user_data);

/**
 * @brief Creates a new logging listener.
 *
 * The listener does nothing until it is installed with dxfg_logging_set_listener(). The function is called as
 * described at dxfg_logging_listener_function_t.
 *
 * Releasing the last handle of the listener with dxfg_JavaObjectHandler_release() (or dxfg_NativeListener_deactivate())
 * stops the calls of the function: after the release returns the function is not running and is not called anymore,
 * so the user data can be freed. An installed listener cannot be replaced or removed (see
 * dxfg_logging_set_listener()), it can only be deactivated this way.
 *
 * Sample:
 * ```c
 * void callback(graal_isolatethread_t *thread, dxfg_logging_level_t level, int64_t timestamp, const char *thread_name,
 *     int64_t thread_id, const char *logger_name, const char *message, dxfg_exception_t *exception,
 *     const char *formatted_message, void *user_data) {
 *     printf("[%" PRId64 "]: %s", thread_id, formatted_message);
 * }
 *
 * // ...
 *
 * // Set the logging Java class
 * dxfg_system_set_property(isolate_thread, "log.className", "com.devexperts.logging.InterceptableLogging");
 * // Turn off logging.
 * dxfg_system_set_property(isolate_thread, "log.level", "OFF");
 *
 * dxfg_logging_listener_t *listener = {};
 *
 * dxfg_logging_listener_new(isolate_thread, &callback, NULL, &listener);
 * dxfg_logging_set_listener(isolate_thread, listener);
 *
 * dxfg_endpoint_t* ep = dxfg_DXEndpoint_getInstance(isolate_thread);
 *
 * // Turn on logging
 * dxfg_logging_set_log_level(isolate_thread, DXFG_LOGGING_LEVEL_INFO);
 * dxfg_DXEndpoint_close(isolate_thread, ep);
 * dxfg_JavaObjectHandler_release(isolate_thread, &listener->handler);
 * dxfg_JavaObjectHandler_release(isolate_thread, &ep->handler);
 * ```
 *
 * Java: `InterceptableLoggingListener.onLog(Level, long, String, long, String, String, Throwable, String)`.
 *
 * @deprecated Use dxfg_v4_logging_listener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The listener callback.
 * @param[in] user_data The user data that will be passed to the callback; the SDK does not use it. It must stay valid
 * until the listener is deactivated.
 * @param[out] listener The pointer to the variable that receives the new listener (release its handle with
 * dxfg_JavaObjectHandler_release()). It is written only on success. A NULL pointer is an error.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_logging_listener_new()")
int32_t dxfg_logging_listener_new(graal_isolatethread_t *thread, dxfg_logging_listener_function_t user_func,
                                  void *user_data, DXFG_OUT dxfg_logging_listener_t **listener);

/**
 * @brief Installs the logging listener.
 *
 * It toggles the logging: the existing handlers of the logging system are removed, so the messages are not logged to
 * the console and to files anymore, and the messages are given to the listener. The level of the messages that reach
 * the listener is taken from the system property `log.level` (`INFO` by default) and can be changed with
 * dxfg_logging_set_log_level().
 *
 * Only the first installed listener works: if a listener is already installed, the function does nothing. A NULL
 * listener does nothing as well. In both cases the function returns #DXFG_EXECUTE_SUCCESSFULLY. The installed
 * listener cannot be removed; release its handle (see dxfg_logging_listener_new()) to stop its calls.
 *
 * Java: `InterceptableLogging.setListener(InterceptableLoggingListener)`.
 *
 * @deprecated Use dxfg_v4_logging_set_listener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] listener The logging listener. The caller keeps its handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_logging_set_listener()")
int32_t dxfg_logging_set_listener(graal_isolatethread_t *thread, dxfg_logging_listener_t *listener);

/**
 * @brief Sets the logging level for the listener or the console or the log file.
 *
 * The level is set for the first handler of the logging system: the handler of the listener if one is installed
 * (see dxfg_logging_set_listener()), otherwise the console or the log file (`log.file`). If there is no such handler,
 * the function does nothing and returns #DXFG_EXECUTE_SUCCESSFULLY.
 *
 * With the handlers of `InterceptableLogging` it also sets the level of the root logger to the more verbose of the
 * two handler levels, so that the messages below it are not created at all.
 *
 * Java: `InterceptableLogging.setLogLevel(Level)`.
 *
 * @deprecated Use dxfg_v4_logging_set_log_level().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] level The new logging level.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_logging_set_log_level()")
int32_t dxfg_logging_set_log_level(graal_isolatethread_t *thread, dxfg_logging_level_t level);

/**
 * @brief Sets the logging level for the err file.
 *
 * The level is set for the second handler of the logging system: the file of the errors (`err.file`). If there is no
 * such handler (e.g. a listener is installed or `err.file` is not set), the function does nothing and returns
 * #DXFG_EXECUTE_SUCCESSFULLY. The default level of the err file is taken from the system property `err.level`
 * (warnings and errors by default). The level of the root logger is set as by the function of the log level.
 *
 * Java: `InterceptableLogging.setErrLevel(Level)`.
 *
 * @deprecated Use dxfg_v4_logging_set_err_level().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] level The new logging level.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_logging_set_err_level()")
int32_t dxfg_logging_set_err_level(graal_isolatethread_t *thread, dxfg_logging_level_t level);

/** @} */ // end of Logging

#ifdef __cplusplus
}
#endif

#endif // DXFG_LOGGING_H