// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK System functions and types declarations
 */

#ifndef DXFG_SYSTEM_H
#define DXFG_SYSTEM_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"

/** @defgroup System System
 *  @brief The Java system properties of the isolate and the shutdown of the SDK's resources.
 *
 *  The ownership rules of the strings are described in the group @ref Javac.
 *  @{
 */

/**
 * @brief Sets the Java system property indicated by the specified key.
 *
 * The property is a property of the Java system of the isolate (`System.setProperty`): the QD library and the SDK
 * read their settings from it, e.g. the properties of the endpoints (`dxfeed.address`, `dxfeed.threadPoolSize`, ...)
 * and of the logging (`log.className`, `log.level`, see the group @ref Logging). The properties that configure
 * endpoints must be set before the first endpoint is created. For example:
 * - `com.dxfeed.sdk.TimeSyncTracker.enable` = `true` enables the QD `TimeSyncTracker`, which sends UDP multicast
 *   packets (239.192.51.45:5145 by default, see `com.devexperts.qd.tools.TimeSyncTracker.addr` and `.port`)
 *   to compare the clocks of the hosts. The SDK disables it by default.
 *
 * Java: `System.setProperty(String, String)`.
 *
 * @deprecated Use dxfg_v4_system_set_property().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] key The name of the property, a NUL-terminated UTF-8 string; the SDK copies it. It must not be NULL or
 * empty.
 * @param[in] value The value of the property, a NUL-terminated UTF-8 string; the SDK copies it. It must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_system_set_property()")
int32_t dxfg_system_set_property(graal_isolatethread_t *thread, const char *key, const char *value);

/**
 * @brief Gets the Java system property indicated by the specified key.
 *
 * Java: `System.getProperty(String)`.
 *
 * @deprecated Use dxfg_v4_system_get_property().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] key The name of the property, a NUL-terminated UTF-8 string. It must not be NULL or empty.
 * @return The value of the property as a NUL-terminated UTF-8 string allocated by the SDK (release it with
 * dxfg_system_release_property()), or NULL if the property is not set. NULL is also returned on error: use
 * dxfg_get_and_clear_thread_exception_t() to get the exception, it returns NULL if there was none.
 */
DXFG_DEPRECATED("use dxfg_v4_system_get_property()")
const char *dxfg_system_get_property(graal_isolatethread_t *thread, const char *key);

/**
 * @brief Releases a property value that dxfg_system_get_property() returned.
 *
 * The string is freed: do not use it after that. Passing NULL does nothing.
 *
 * @deprecated Use dxfg_v4_system_release_property().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] value The value returned by dxfg_system_get_property().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_system_release_property()")
int32_t dxfg_system_release_property(graal_isolatethread_t *thread, const char *value);

/**
 * @brief Closes all the endpoints and instrument profile connections created through the SDK in this isolate
 * and waits for the termination of the endpoints.
 *
 * It includes the default endpoints of dxfg_DXEndpoint_getInstance(), dxfg_DXFeed_getInstance(),
 * dxfg_DXPublisher_getInstance() and dxfg_OnDemandService_getInstance(). The handles of the objects stay valid
 * and must be released as usual.
 *
 * Call it before graal_tear_down_isolate() (or before the process exits): the threads of an open endpoint
 * ("DXEndpoint-DXExecutorThread-N", etc.) do not stop on `Thread.interrupt()`, so graal_tear_down_isolate()
 * waits for them forever while any endpoint is open. To stop the calls of the native listeners, release them
 * (see dxfg_JavaObjectHandler_release()) before tearing down the isolate. To find the threads that block
 * the tear-down, pass the `-XX:TearDownWarningSeconds=<seconds>` option in `graal_create_isolate_params_t.argv`.
 *
 * Recommended shutdown sequence:
 * 1. Release the listeners (dxfg_JavaObjectHandler_release() or dxfg_NativeListener_deactivate()).
 * 2. dxfg_system_close_all_and_await_termination().
 * 3. Release the remaining handles.
 * 4. graal_tear_down_isolate().
 *
 * The function blocks the calling thread until the endpoints are terminated: do not call it from a listener of an
 * endpoint, it would wait for the notification that is calling it. The endpoints and connections created after the
 * call are not closed by it. The function is a function of the SDK: there is no QD method that it calls.
 *
 * @deprecated Use dxfg_v4_system_close_all_and_await_termination().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also if
 * the waiting thread is interrupted). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_system_close_all_and_await_termination()")
int32_t dxfg_system_close_all_and_await_termination(graal_isolatethread_t *thread);

/** @} */ // end of System

#ifdef __cplusplus
}
#endif

#endif // DXFG_SYSTEM_H
