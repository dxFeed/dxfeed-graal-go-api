// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_logging.h and the v4 names of its types (see the group @ref V4_Logging).
 */

#ifndef DXFG_V4_LOGGING_H
#define DXFG_V4_LOGGING_H

#include "dxfg_logging.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_catch_exception.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Logging Logging (v4)
 *  @brief The v4 functions and type names of the group @ref Logging. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_logging_level_t */
typedef dxfg_logging_level_t dxfg_v4_logging_level_t;

/** @copydoc dxfg_logging_listener_t */
typedef dxfg_logging_listener_t dxfg_v4_logging_listener_t;

/** @copydoc dxfg_logging_listener_function_t */
typedef dxfg_logging_listener_function_t dxfg_v4_logging_listener_function_t;

/**
 * @brief Creates a new logging listener.
 *
 * The listener does nothing until it is installed with dxfg_v4_logging_set_listener(). The function is called as
 * described at dxfg_v4_logging_listener_function_t.
 *
 * Releasing the last handle of the listener with dxfg_v4_JavaObjectHandler_release() (or
 * dxfg_v4_NativeListener_deactivate()) stops the calls of the function: after the release returns the function is not
 * running and is not called anymore, so the user data can be freed. An installed listener cannot be replaced or removed
 * (see dxfg_v4_logging_set_listener()), it can only be deactivated this way.
 *
 * Java: `InterceptableLoggingListener.onLog(Level, long, String, long, String, String, Throwable, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The listener callback.
 * @param[in] user_data The user data that will be passed to the callback; the SDK does not use it. It must stay valid
 * until the listener is deactivated.
 * @param[out] listener The pointer to the variable that receives the new listener (release its handle with
 * dxfg_v4_JavaObjectHandler_release()). It is written only on success. A NULL pointer is an error.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_logging_listener_new(graal_isolatethread_t *thread, dxfg_v4_logging_listener_function_t user_func,
                                     void *user_data, DXFG_OUT dxfg_v4_logging_listener_t **listener);

/**
 * @brief Installs the logging listener.
 *
 * It toggles the logging: the existing handlers of the logging system are removed, so the messages are not logged to
 * the console and to files anymore, and the messages are given to the listener. The level of the messages that reach
 * the listener is taken from the system property `log.level` (`INFO` by default) and can be changed with
 * dxfg_v4_logging_set_log_level().
 *
 * Only the first installed listener works: if a listener is already installed, the function does nothing. A NULL
 * listener does nothing as well. In both cases the function returns #DXFG_EXECUTE_SUCCESSFULLY. The installed
 * listener cannot be removed; release its handle (see dxfg_v4_logging_listener_new()) to stop its calls.
 *
 * Java: `InterceptableLogging.setListener(InterceptableLoggingListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] listener The logging listener. The caller keeps its handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_logging_set_listener(graal_isolatethread_t *thread, dxfg_v4_logging_listener_t *listener);

/**
 * @brief Sets the logging level for the listener or the console or the log file.
 *
 * The level is set for the first handler of the logging system: the handler of the listener if one is installed
 * (see dxfg_v4_logging_set_listener()), otherwise the console or the log file (`log.file`). If there is no such
 * handler, the function does nothing and returns #DXFG_EXECUTE_SUCCESSFULLY.
 *
 * With the handlers of `InterceptableLogging` it also sets the level of the root logger to the more verbose of the
 * two handler levels, so that the messages below it are not created at all.
 *
 * Java: `InterceptableLogging.setLogLevel(Level)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] level The new logging level.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_logging_set_log_level(graal_isolatethread_t *thread, dxfg_v4_logging_level_t level);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] level The new logging level.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_logging_set_err_level(graal_isolatethread_t *thread, dxfg_v4_logging_level_t level);

/** @} */ // end of V4_Logging

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_LOGGING_H
