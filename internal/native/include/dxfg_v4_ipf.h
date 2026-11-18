// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_ipf.h and the v4 names of its types (see the group @ref V4_InstrumentProfile).
 */

#ifndef DXFG_V4_IPF_H
#define DXFG_V4_IPF_H

#include "dxfg_ipf.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_InstrumentProfile InstrumentProfile (v4)
 *  @brief The v4 functions and type names of the group @ref InstrumentProfile. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_ipf_connection_t */
typedef dxfg_ipf_connection_t dxfg_v4_ipf_connection_t;

/** @copydoc dxfg_instrument_profile_reader_t */
typedef dxfg_instrument_profile_reader_t dxfg_v4_instrument_profile_reader_t;

/** @copydoc dxfg_ipf_connection_state_change_listener_t */
typedef dxfg_ipf_connection_state_change_listener_t dxfg_v4_ipf_connection_state_change_listener_t;

/** @copydoc dxfg_ipf_connection_state_t */
typedef dxfg_ipf_connection_state_t dxfg_v4_ipf_connection_state_t;

/** @copydoc dxfg_ipf_connection_state_change_listener_func */
typedef dxfg_ipf_connection_state_change_listener_func dxfg_v4_ipf_connection_state_change_listener_func;

/** @copydoc dxfg_ipf_collector_t */
typedef dxfg_ipf_collector_t dxfg_v4_ipf_collector_t;

/** @copydoc dxfg_instrument_profile_t */
typedef dxfg_instrument_profile_t dxfg_v4_instrument_profile_t;

/** @copydoc dxfg_instrument_profile_list */
typedef dxfg_instrument_profile_list dxfg_v4_instrument_profile_list;

/** @copydoc dxfg_ipf_update_listener_t */
typedef dxfg_ipf_update_listener_t dxfg_v4_ipf_update_listener_t;

/** @copydoc dxfg_iterable_ip_t */
typedef dxfg_iterable_ip_t dxfg_v4_iterable_ip_t;

/** @copydoc dxfg_ipf_update_listener_function */
typedef dxfg_ipf_update_listener_function dxfg_v4_ipf_update_listener_function;

/** @copydoc dxfg_instrument_profile_custom_fields_t */
typedef dxfg_instrument_profile_custom_fields_t dxfg_v4_instrument_profile_custom_fields_t;

/** @copydoc dxfg_instrument_profile2_t */
typedef dxfg_instrument_profile2_t dxfg_v4_instrument_profile_data_t;

/** @copydoc dxfg_instrument_profile2_list_t */
typedef dxfg_instrument_profile2_list_t dxfg_v4_instrument_profile_data_list;

/**
 * @brief Creates an instrument profile connection to an address that pushes the profiles into a collector.
 *
 * The address is `<host>:<port>` of a server, a URL, or a file path. The clause `[update=<period>]` can be added at the
 * end of the address to set the update period (see dxfg_v4_InstrumentProfileConnection_getUpdatePeriod()); the default
 * period is 1 minute. The connection is not started: call dxfg_v4_InstrumentProfileConnection_start().
 *
 * Java: `InstrumentProfileConnection.createConnection(String, InstrumentProfileCollector)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] collector The collector handle that the connection pushes the updates into.
 * @param[out] ipfConnection The pointer to the variable that receives the connection (release it with
 * dxfg_v4_JavaObjectHandler_release(); that does not close the connection); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_createConnection(graal_isolatethread_t *thread, const char *address,
                                                             dxfg_v4_ipf_collector_t *collector,
                                                             DXFG_OUT dxfg_v4_ipf_connection_t **ipfConnection);

/**
 * @brief Returns the address of the connection, without the additional options given as a part of it.
 *
 * Java: `InstrumentProfileConnection.getAddress()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[out] address The pointer to the variable that receives the address (a UTF-8 string; release it with
 * dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_getAddress(graal_isolatethread_t *thread,
                                                       dxfg_v4_ipf_connection_t *connection,
                                                       DXFG_OUT const char **address);

/**
 * @brief Returns the update period of the connection in milliseconds.
 *
 * It is the period of an update check when the instrument profile source does not support live updates and when the
 * connection is dropped. The default period is 1 minute unless the address of the connection sets another one.
 *
 * Java: `InstrumentProfileConnection.getUpdatePeriod()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[out] updatePeriod The pointer to the variable that receives the update period in milliseconds; must not be
 * NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_getUpdatePeriod(graal_isolatethread_t *thread,
                                                            dxfg_v4_ipf_connection_t *connection,
                                                            DXFG_OUT int64_t *updatePeriod);

/**
 * @brief Changes the update period of the connection.
 *
 * Java: `InstrumentProfileConnection.setUpdatePeriod(long)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] newValue The update period in milliseconds; it must not be negative (the function fails otherwise).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_setUpdatePeriod(graal_isolatethread_t *thread,
                                                            dxfg_v4_ipf_connection_t *connection, int64_t newValue);

/**
 * @brief Returns the state of the connection.
 *
 * Java: `InstrumentProfileConnection.getState()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[out] state The pointer to the variable that receives the state; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_getState(graal_isolatethread_t *thread,
                                                     dxfg_v4_ipf_connection_t *connection,
                                                     DXFG_OUT dxfg_v4_ipf_connection_state_t *state);

/**
 * @brief Returns the last modification time of the instrument profiles in milliseconds, or 0 if it is unknown.
 *
 * The time is in milliseconds, but the actual granularity of the time is a second.
 *
 * Java: `InstrumentProfileConnection.getLastModified()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[out] lastModified The pointer to the variable that receives the last modification time in milliseconds, or 0
 * if it is unknown; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_getLastModified(graal_isolatethread_t *thread,
                                                            dxfg_v4_ipf_connection_t *connection,
                                                            DXFG_OUT int64_t *lastModified);

/**
 * @brief Starts the connection.
 *
 * The state of the connection immediately changes to #DXFG_IPF_CONNECTION_STATE_CONNECTING and the actual connection is
 * established in the background. The function fails if the state is not #DXFG_IPF_CONNECTION_STATE_NOT_CONNECTED (the
 * connection is already started or closed).
 *
 * Java: `InstrumentProfileConnection.start()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_start(graal_isolatethread_t *thread, dxfg_v4_ipf_connection_t *connection);

/**
 * @brief Closes the connection.
 *
 * The state of the connection immediately changes to #DXFG_IPF_CONNECTION_STATE_CLOSED and the background update
 * stops: nothing that the connection receives after the call reaches the collector. Its thread and its connection to
 * the server end when the next data arrives (a live server sends a heartbeat every 10 seconds) or at the 60-second read
 * timeout. The collector may still notify its listeners about the updates received before the call: release a listener
 * with dxfg_v4_JavaObjectHandler_release() before freeing its user_data. Closing a closed connection does nothing.
 *
 * Java: `InstrumentProfileConnection.close()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_close(graal_isolatethread_t *thread, dxfg_v4_ipf_connection_t *connection);

/**
 * @brief Creates a listener of the state changes of a connection that calls a C function.
 *
 * Add the listener with dxfg_v4_InstrumentProfileConnection_addStateChangeListener(). The function is called by a Java
 * thread, see dxfg_v4_ipf_connection_state_change_listener_func. Releasing the last handle of the listener deactivates
 * it: after dxfg_v4_JavaObjectHandler_release() returns the function is not called anymore, so user_data can be freed.
 *
 * Java: `PropertyChangeListener.propertyChange(PropertyChangeEvent)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call when the state of the connection changes.
 * @param[in] user_data The pointer that is passed to user_func as is; the SDK does not use it.
 * @param[out] ipfConnectionStateChangeListener The pointer to the variable that receives the listener (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_IpfPropertyChangeListener_new(
    graal_isolatethread_t *thread, dxfg_v4_ipf_connection_state_change_listener_func user_func, void *user_data,
    DXFG_OUT dxfg_v4_ipf_connection_state_change_listener_t **ipfConnectionStateChangeListener);

/**
 * @brief Adds a listener of the state changes of the connection.
 *
 * The listener is called for the following state changes only, not for the current state. Remove it with
 * dxfg_v4_InstrumentProfileConnection_removeStateChangeListener().
 *
 * Java: `InstrumentProfileConnection.addStateChangeListener(PropertyChangeListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] listener The listener handle, see dxfg_v4_IpfPropertyChangeListener_new().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_addStateChangeListener(
    graal_isolatethread_t *thread, dxfg_v4_ipf_connection_t *connection,
    dxfg_v4_ipf_connection_state_change_listener_t *listener);

/**
 * @brief Removes a listener of the state changes of the connection.
 *
 * The listener handle is not released.
 *
 * Java: `InstrumentProfileConnection.removeStateChangeListener(PropertyChangeListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] listener The listener handle that was added with
 * dxfg_v4_InstrumentProfileConnection_addStateChangeListener().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_removeStateChangeListener(
    graal_isolatethread_t *thread, dxfg_v4_ipf_connection_t *connection,
    dxfg_v4_ipf_connection_state_change_listener_t *listener);

/**
 * @brief Waits until the initial snapshot of the instrument profiles is fully read, for at most the timeout.
 *
 * The calling thread is blocked until the state #DXFG_IPF_CONNECTION_STATE_COMPLETED is reached, the connection is
 * closed, or the timeout elapses.
 *
 * Java: `InstrumentProfileConnection.waitUntilCompleted(long, TimeUnit)`. The timeout is in milliseconds.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] timeoutInMs The maximum time to wait in milliseconds.
 * @param[out] completed The pointer to the variable that receives 1 if the state #DXFG_IPF_CONNECTION_STATE_COMPLETED
 * was reached, 0 if the time elapsed or the connection was closed before the snapshot was fully read; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileConnection_waitUntilCompleted(graal_isolatethread_t *thread,
                                                               dxfg_v4_ipf_connection_t *connection,
                                                               int64_t timeoutInMs, DXFG_OUT int32_t *completed);

/**
 * @brief Releases the list structure only and keeps the instrument profile handles in it.
 *
 * The array of the elements and the dxfg_v4_instrument_profile_list structure are freed. The profile handles stay
 * valid: release each of them with dxfg_v4_JavaObjectHandler_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ips The list of the instrument profile handles.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_instrument_profile_list_wrapper_release(graal_isolatethread_t *thread,
                                                        dxfg_v4_instrument_profile_list *ips);

/**
 * @brief Releases the list and all the instrument profile handles in it.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ips The list of the instrument profile handles.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_instrument_profile_list_release(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_list *ips);

/**
 * @brief Creates an instrument profile collector.
 *
 * The collector keeps one instrument profile per symbol and notifies its listeners about the updates.
 *
 * Java: `new InstrumentProfileCollector()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] ipfCollector The pointer to the variable that receives the collector (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_new(graal_isolatethread_t *thread,
                                               DXFG_OUT dxfg_v4_ipf_collector_t **ipfCollector);

/**
 * @brief Returns the last modification time of the instrument profiles of the collector in milliseconds.
 *
 * The time is in milliseconds, but the actual granularity of the time is a second.
 *
 * Java: `InstrumentProfileCollector.getLastUpdateTime()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[out] lastUpdateTime The pointer to the variable that receives the last modification time in milliseconds; must
 * not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_getLastUpdateTime(graal_isolatethread_t *thread,
                                                             dxfg_v4_ipf_collector_t *collector,
                                                             DXFG_OUT int64_t *lastUpdateTime);

/**
 * @brief Updates one instrument profile in the collector.
 *
 * The collector stores the Java instrument profile object behind the handle, not a copy: do not modify the profile
 * after the call. It is the same as dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandles() with one
 * profile and no generation.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfile(InstrumentProfile)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] ip The instrument profile handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandle(graal_isolatethread_t *thread,
                                                                         dxfg_v4_ipf_collector_t *collector,
                                                                         dxfg_v4_instrument_profile_t *ip);

/**
 * @brief Updates a list of instrument profiles in the collector and assigns them a generation tag.
 *
 * The removal of an instrument profile is represented by a profile whose type is `REMOVED`. The collector stores the
 * Java instrument profile objects behind the handles, not copies: do not modify the profiles after the call. The list
 * and its handles still belong to the caller.
 *
 * A list with a NULL element, or with a positive size and NULL `elements`, is an error: the function fails before it
 * updates any profile. A NULL list or a list of size 0 updates nothing.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfiles(List, Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] ips The list of the instrument profile handles.
 * @param[in] generation The generation tag, the handle of a Java object, or NULL for none. The tag is used by
 * dxfg_v4_InstrumentProfileCollector_removeGenerations().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandles(graal_isolatethread_t *thread,
                                                                          dxfg_v4_ipf_collector_t *collector,
                                                                          dxfg_v4_instrument_profile_list *ips,
                                                                          dxfg_v4_java_object_handler *generation);

/**
 * @brief Removes the instrument profiles whose generation tags are in a list.
 *
 * The profiles without a generation tag are not removed. The time of the call grows with the number of the profiles in
 * the collector: use it sparingly. The list and its handles still belong to the caller.
 *
 * Java: `InstrumentProfileCollector.removeGenerations(Set)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] generations The list of the generation tags, the handles of Java objects.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_removeGenerations(graal_isolatethread_t *thread,
                                                             dxfg_v4_ipf_collector_t *collector,
                                                             dxfg_v4_java_object_handler_list *generations);

/**
 * @brief Returns an iterator over the instrument profiles of the collector.
 *
 * The iterator is a concurrent view of the set of the profiles. It exposes the profiles that are not removed; however,
 * if the iteration is concurrent with a removal, a removed profile (its type is `REMOVED`) can be exposed.
 *
 * Java: `InstrumentProfileCollector.view()`. The SDK returns the iterator of the view.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[out] iterableIp The pointer to the variable that receives the iterator (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_view(graal_isolatethread_t *thread, dxfg_v4_ipf_collector_t *collector,
                                                DXFG_OUT dxfg_v4_iterable_ip_t **iterableIp);

/**
 * @brief Returns the executor that calls the update listeners of the collector.
 *
 * By default a single-threaded executor is used to process the updates asynchronously.
 *
 * Java: `InstrumentProfileCollector.getExecutor()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[out] executor The pointer to the variable that receives the executor (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_getExecutor(graal_isolatethread_t *thread,
                                                       dxfg_v4_ipf_collector_t *collector,
                                                       DXFG_OUT dxfg_v4_executor_t **executor);

/**
 * @brief Changes the executor that calls the update listeners of the collector.
 *
 * The executor of a shared event queue (dxfg_v4_EventQueue_getExecutor()) makes the polls of the queue call the
 * update listeners; the executor of a subscription attached to a queue is refused with `IllegalStateException` (see
 * @ref V4_EventQueue).
 *
 * Java: `InstrumentProfileCollector.setExecutor(Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] executor The executor handle; it must not be NULL (the function fails otherwise).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_setExecutor(graal_isolatethread_t *thread,
                                                       dxfg_v4_ipf_collector_t *collector,
                                                       dxfg_v4_executor_t *executor);

/**
 * @brief Adds a listener that is notified about any updates in the set of the instrument profiles.
 *
 * If the set of the profiles is not empty, the listener is notified about it as well (by the executor of the
 * collector). Remove the listener with dxfg_v4_InstrumentProfileCollector_removeUpdateListener().
 *
 * Java: `InstrumentProfileCollector.addUpdateListener(InstrumentProfileUpdateListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] listener The listener handle, see dxfg_v4_InstrumentProfileUpdateListener_new().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_addUpdateListener(graal_isolatethread_t *thread,
                                                             dxfg_v4_ipf_collector_t *collector,
                                                             dxfg_v4_ipf_update_listener_t* listener);

/**
 * @brief Removes a listener that was added with dxfg_v4_InstrumentProfileCollector_addUpdateListener().
 *
 * The listener handle is not released.
 *
 * Java: `InstrumentProfileCollector.removeUpdateListener(InstrumentProfileUpdateListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_removeUpdateListener(graal_isolatethread_t *thread,
                                                                dxfg_v4_ipf_collector_t *collector,
                                                                dxfg_v4_ipf_update_listener_t* listener);

/**
 * @brief Checks whether the iterator has a next instrument profile.
 *
 * Java: `Iterator.hasNext()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] iterable_ip The iterator handle.
 * @param[out] hasNext The pointer to the variable that receives 1 if the iterator has a next instrument profile, 0 if
 * it has not; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Iterable_InstrumentProfile_hasNext(graal_isolatethread_t *thread, dxfg_v4_iterable_ip_t *iterable_ip,
                                                   DXFG_OUT int32_t *hasNext);

/**
 * @brief Returns the next instrument profile of the iterator as a handle.
 *
 * Check dxfg_v4_Iterable_InstrumentProfile_hasNext() first: the function fails if there is no next profile. The handle
 * refers to the profile object itself (the collector stores the profiles, not copies): do not modify it. Use
 * dxfg_v4_Iterable_InstrumentProfile_next() to get a copy.
 *
 * Java: `Iterator.next()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] iterable_ip The iterator handle.
 * @param[out] instrumentProfile The pointer to the variable that receives the instrument profile (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Iterable_InstrumentProfile_nextHandle(graal_isolatethread_t *thread, dxfg_v4_iterable_ip_t *iterable_ip,
                                                      DXFG_OUT dxfg_v4_instrument_profile_t **instrumentProfile);

/**
 * @brief Creates a listener of the updates of a collector that calls a C function.
 *
 * Add the listener with dxfg_v4_InstrumentProfileCollector_addUpdateListener(). The function is called by the executor
 * of the collector, see dxfg_v4_ipf_update_listener_function. Releasing the last handle of the listener deactivates it:
 * after dxfg_v4_JavaObjectHandler_release() returns the function is not called anymore, so user_data can be freed.
 *
 * Java: `InstrumentProfileUpdateListener.instrumentProfilesUpdated(Iterator)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call when the set of the instrument profiles changes.
 * @param[in] user_data The pointer that is passed to user_func as is; the SDK does not use it.
 * @param[out] ipfUpdateListener The pointer to the variable that receives the listener (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileUpdateListener_new(graal_isolatethread_t *thread,
                                                    dxfg_v4_ipf_update_listener_function user_func, void *user_data,
                                                    DXFG_OUT dxfg_v4_ipf_update_listener_t **ipfUpdateListener);

/**
 * @brief Creates an instrument profile reader.
 *
 * Java: `new InstrumentProfileReader()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] instrumentProfileReader The pointer to the variable that receives the reader (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_new(graal_isolatethread_t *thread,
                                            DXFG_OUT dxfg_v4_instrument_profile_reader_t **instrumentProfileReader);

/**
 * @brief Returns the last modification time in milliseconds from the last read from an address, or 0 if it is unknown.
 *
 * Java: `InstrumentProfileReader.getLastModified()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[out] lastModified The pointer to the variable that receives the last modification time in milliseconds, or 0
 * if it is unknown; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_getLastModified(graal_isolatethread_t *thread,
                                                        dxfg_v4_instrument_profile_reader_t *reader,
                                                        DXFG_OUT int64_t *lastModified);

/**
 * @brief Checks whether the instrument profile file was fully read by the last read operation.
 *
 * Java: `InstrumentProfileReader.wasComplete()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[out] complete The pointer to the variable that receives 1 if the file was fully read, 0 if it was not; must
 * not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_wasComplete(graal_isolatethread_t *thread,
                                                    dxfg_v4_instrument_profile_reader_t *reader,
                                                    DXFG_OUT int32_t *complete);

/**
 * @brief Reads and returns the instrument profiles from an address.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_v4_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL,
 * for example `http://user:password@host:port/path/file.ipf`. The operation updates the values of
 * dxfg_v4_InstrumentProfileReader_getLastModified() and dxfg_v4_InstrumentProfileReader_wasComplete(). The function
 * fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[out] instrumentProfileList The pointer to the variable that receives the list of the instrument profiles
 * (release it and the profiles in it with dxfg_v4_instrument_profile_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readFromFileAsHandles(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, const char *address,
    DXFG_OUT dxfg_v4_instrument_profile_list **instrumentProfileList);

/**
 * @brief Reads and returns the instrument profiles from an address with basic authentication credentials.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_v4_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL,
 * for example `http://user:password@host:port/path/file.ipf`. The user and the password take precedence over the
 * authentication information of the URL. The operation updates the values of
 * dxfg_v4_InstrumentProfileReader_getLastModified() and dxfg_v4_InstrumentProfileReader_wasComplete(). The function
 * fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] user The user name of the basic authentication, a UTF-8 string, or NULL.
 * @param[in] password The password of the basic authentication, a UTF-8 string, or NULL.
 * @param[out] instrumentProfileList The pointer to the variable that receives the list of the instrument profiles
 * (release it and the profiles in it with dxfg_v4_instrument_profile_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPasswordAsHandles(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, const char *address, const char *user,
    const char *password, DXFG_OUT dxfg_v4_instrument_profile_list **instrumentProfileList);

/**
 * @brief Reads and returns the instrument profiles from an address with token credentials.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_v4_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL,
 * for example `http://user:password@host:port/path/file.ipf`. The token takes precedence over the authentication
 * information of the URL. The operation updates the values of dxfg_v4_InstrumentProfileReader_getLastModified() and
 * dxfg_v4_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument
 * Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, AuthToken)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] token The token handle, or NULL.
 * @param[out] instrumentProfileList The pointer to the variable that receives the list of the instrument profiles
 * (release it and the profiles in it with dxfg_v4_instrument_profile_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readFromFileWithTokenAsHandles(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, const char *address,
    dxfg_v4_auth_token_t *token, DXFG_OUT dxfg_v4_instrument_profile_list **instrumentProfileList);

/**
 * @brief Converts an address specification into the URL that is read by the
 * dxfg_v4_InstrumentProfileReader_readFromFile*() functions.
 *
 * An address of the form `<host>:<port>` (it contains a colon, no slash and ends with a valid port number) becomes
 * `http://<host>:<port>/ipf/all.ipf.gz`, the query part of the address is kept; any other address is returned as is.
 *
 * Java: `InstrumentProfileReader.resolveSourceURL(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] address The address, a UTF-8 string.
 * @param[out] url The pointer to the variable that receives the URL (a UTF-8 string; release it with
 * dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_resolveSourceURL(graal_isolatethread_t *thread, const char *address,
                                                         DXFG_OUT const char **url);

/**
 * @brief Reads and returns the instrument profiles from an input stream.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The stream is not closed. The operation updates the value of
 * dxfg_v4_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument
 * Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[in] address The origin of the stream for debugging purposes, a UTF-8 string.
 * @param[out] instrumentProfileList The pointer to the variable that receives the list of the instrument profiles
 * (release it and the profiles in it with dxfg_v4_instrument_profile_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readWithAddressAsHandles(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, dxfg_v4_input_stream_t* is,
    const char* address, DXFG_OUT dxfg_v4_instrument_profile_list **instrumentProfileList);

/**
 * @brief Reads and returns the instrument profiles from a compressed input stream.
 *
 * It is deprecated in Java; use dxfg_v4_InstrumentProfileReader_readWithAddressAsHandles(). The data compression
 * formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a plain data stream,
 * other data is parsed as is. The stream is not closed. The operation updates the value of
 * dxfg_v4_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument
 * Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readCompressed(InputStream)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrumentProfileList The pointer to the variable that receives the list of the instrument profiles
 * (release it and the profiles in it with dxfg_v4_instrument_profile_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readCompressedAsHandles(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, dxfg_v4_input_stream_t* is,
    DXFG_OUT dxfg_v4_instrument_profile_list **instrumentProfileList);

/**
 * @brief Reads and returns the instrument profiles from an input stream with plain (not compressed) data.
 *
 * It is deprecated in Java; use dxfg_v4_InstrumentProfileReader_readWithAddressAsHandles(). Unlike it, the function
 * does not decompress the stream and does not report an incomplete one. The function fails if the data does not conform
 * to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrumentProfileList The pointer to the variable that receives the list of the instrument profiles
 * (release it and the profiles in it with dxfg_v4_instrument_profile_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readAsHandles(graal_isolatethread_t *thread,
                                                      dxfg_v4_instrument_profile_reader_t *reader,
                                                      dxfg_v4_input_stream_t* is,
                                                      DXFG_OUT dxfg_v4_instrument_profile_list **instrumentProfileList);

/**
 * @brief Creates an instrument profile with the default values: empty strings and zero numbers.
 *
 * Java: `new InstrumentProfile()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] instrumentProfile The pointer to the variable that receives the instrument profile (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_new(graal_isolatethread_t *thread,
                                      DXFG_OUT dxfg_v4_instrument_profile_t **instrumentProfile);

/**
 * @brief Creates an instrument profile as a copy of another instrument profile.
 *
 * Java: `new InstrumentProfile(InstrumentProfile)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The handle of the instrument profile to copy.
 * @param[out] instrumentProfile The pointer to the variable that receives the instrument profile (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_newCopy(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t* ip,
                                          DXFG_OUT dxfg_v4_instrument_profile_t **instrumentProfile);

/**
 * @brief Returns the type of the instrument of the instrument profile.
 *
 * It takes precedence over the other fields in conflict cases. It is a mandatory field and must not be empty. Examples:
 * `STOCK`, `FUTURE`, `OPTION`.
 *
 * Java: `InstrumentProfile.getType()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] type The pointer to the variable that receives the type (a UTF-8 string, empty if it is not set; release
 * it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getType(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                          DXFG_OUT const char **type);

/**
 * @brief Sets the type of the instrument of the instrument profile.
 *
 * It takes precedence over the other fields in conflict cases. It is a mandatory field and must not be empty. Examples:
 * `STOCK`, `FUTURE`, `OPTION`.
 *
 * Java: `InstrumentProfile.setType(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The type, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setType(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                          const char *value);

/**
 * @brief Returns the identifier of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. It is a mandatory field and must not be empty. Examples:
 * `GOOG`, `/YGM9`, `.ZYEAD`.
 *
 * Java: `InstrumentProfile.getSymbol()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] symbol The pointer to the variable that receives the symbol (a UTF-8 string, empty if it is not set;
 * release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getSymbol(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                            DXFG_OUT const char **symbol);

/**
 * @brief Sets the identifier of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. It is a mandatory field and must not be empty. Examples:
 * `GOOG`, `/YGM9`, `.ZYEAD`.
 *
 * Java: `InstrumentProfile.setSymbol(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The symbol, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setSymbol(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                            const char *value);

/**
 * @brief Returns the description of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. Examples: `Google Inc.`, `Mini Gold
 * Futures,Jun-2009,ETH`.
 *
 * Java: `InstrumentProfile.getDescription()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] description The pointer to the variable that receives the description (a UTF-8 string, empty if it is not
 * set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getDescription(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                 DXFG_OUT const char **description);

/**
 * @brief Sets the description of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. Examples: `Google Inc.`, `Mini Gold
 * Futures,Jun-2009,ETH`.
 *
 * Java: `InstrumentProfile.setDescription(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The description, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setDescription(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                 const char *value);

/**
 * @brief Returns the identifier of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the symbol.
 *
 * Java: `InstrumentProfile.getLocalSymbol()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] localSymbol The pointer to the variable that receives the local symbol (a UTF-8 string, empty if it is
 * not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getLocalSymbol(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                 DXFG_OUT const char **localSymbol);

/**
 * @brief Sets the identifier of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the symbol.
 *
 * Java: `InstrumentProfile.setLocalSymbol(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The local symbol, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setLocalSymbol(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                 const char *value);

/**
 * @brief Returns the description of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the description.
 *
 * Java: `InstrumentProfile.getLocalDescription()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] localDescription The pointer to the variable that receives the local description (a UTF-8 string, empty
 * if it is not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getLocalDescription(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                      DXFG_OUT const char **localDescription);

/**
 * @brief Sets the description of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the description.
 *
 * Java: `InstrumentProfile.setLocalDescription(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The local description, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setLocalDescription(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                      const char *value);

/**
 * @brief Returns the country of origin (incorporation) of the corresponding company or parent entity of the instrument
 * profile.
 *
 * It uses the two-letter country code of ISO 3166-1. Examples: `US`, `RU`.
 *
 * Java: `InstrumentProfile.getCountry()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] country The pointer to the variable that receives the country (a UTF-8 string, empty if it is not set;
 * release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getCountry(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                             DXFG_OUT const char **country);

/**
 * @brief Sets the country of origin (incorporation) of the corresponding company or parent entity of the instrument
 * profile.
 *
 * It uses the two-letter country code of ISO 3166-1. Examples: `US`, `RU`.
 *
 * Java: `InstrumentProfile.setCountry(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The country, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setCountry(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                             const char *value);

/**
 * @brief Returns the official Place Of Listing, the organization that has listed the instrument of the instrument
 * profile.
 *
 * An instrument with multiple listings uses a separate profile for each listing. It uses the Market Identifier Code
 * (MIC) of ISO 10383. Examples: `XNAS`, `RTSX`.
 *
 * Java: `InstrumentProfile.getOPOL()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] opol The pointer to the variable that receives the OPOL (a UTF-8 string, empty if it is not set; release
 * it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getOPOL(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                          DXFG_OUT const char **opol);

/**
 * @brief Sets the official Place Of Listing, the organization that has listed the instrument of the instrument profile.
 *
 * An instrument with multiple listings uses a separate profile for each listing. It uses the Market Identifier Code
 * (MIC) of ISO 10383. Examples: `XNAS`, `RTSX`.
 *
 * Java: `InstrumentProfile.setOPOL(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The OPOL, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setOPOL(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                          const char *value);

/**
 * @brief Returns the exchange-specific data required to identify the instrument when communicating with the exchange of
 * the instrument profile.
 *
 * It uses an exchange-specific format.
 *
 * Java: `InstrumentProfile.getExchangeData()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] exchangeData The pointer to the variable that receives the exchange data (a UTF-8 string, empty if it is
 * not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getExchangeData(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  DXFG_OUT const char **exchangeData);

/**
 * @brief Sets the exchange-specific data required to identify the instrument when communicating with the exchange of
 * the instrument profile.
 *
 * It uses an exchange-specific format.
 *
 * Java: `InstrumentProfile.setExchangeData(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The exchange data, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setExchangeData(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  const char *value);

/**
 * @brief Returns the list of the exchanges where the instrument is quoted or traded of the instrument profile.
 *
 * The list is a sequence of MICs separated by semicolons and sorted by MIC. Example: `ARCX;CBSX;XNAS;XNYS`.
 *
 * Java: `InstrumentProfile.getExchanges()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] exchanges The pointer to the variable that receives the exchanges (a UTF-8 string, empty if it is not
 * set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getExchanges(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                               DXFG_OUT const char **exchanges);

/**
 * @brief Sets the list of the exchanges where the instrument is quoted or traded of the instrument profile.
 *
 * The list is a sequence of MICs separated by semicolons and sorted by MIC. Example: `ARCX;CBSX;XNAS;XNYS`.
 *
 * Java: `InstrumentProfile.setExchanges(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The exchanges, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setExchanges(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                               const char *value);

/**
 * @brief Returns the currency of quotation, pricing and trading of the instrument profile.
 *
 * It uses the three-letter currency code of ISO 4217. Examples: `USD`, `RUB`.
 *
 * Java: `InstrumentProfile.getCurrency()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] currency The pointer to the variable that receives the currency (a UTF-8 string, empty if it is not set;
 * release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getCurrency(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                              DXFG_OUT const char **currency);

/**
 * @brief Sets the currency of quotation, pricing and trading of the instrument profile.
 *
 * It uses the three-letter currency code of ISO 4217. Examples: `USD`, `RUB`.
 *
 * Java: `InstrumentProfile.setCurrency(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The currency, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setCurrency(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                              const char *value);

/**
 * @brief Returns the base currency of a currency pair (FOREX instruments) of the instrument profile.
 *
 * It uses the three-letter currency code like the currency.
 *
 * Java: `InstrumentProfile.getBaseCurrency()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] baseCurrency The pointer to the variable that receives the base currency (a UTF-8 string, empty if it is
 * not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getBaseCurrency(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  DXFG_OUT const char **baseCurrency);

/**
 * @brief Sets the base currency of a currency pair (FOREX instruments) of the instrument profile.
 *
 * It uses the three-letter currency code like the currency.
 *
 * Java: `InstrumentProfile.setBaseCurrency(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The base currency, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setBaseCurrency(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  const char *value);

/**
 * @brief Returns the Classification of Financial Instruments code of the instrument profile.
 *
 * It is a mandatory field for OPTION instruments: it is the only way to distinguish Call/Put, American/European
 * exercise and Cash/Physical delivery. It uses the six-letter code of ISO 10962; `X` can be used extensively and
 * trailing letters can be omitted (they are assumed to be `X`). Examples: `ESNTPB`, `ESXXXX`, `ES`, `OPASPS`.
 *
 * Java: `InstrumentProfile.getCFI()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] cfi The pointer to the variable that receives the CFI (a UTF-8 string, empty if it is not set; release it
 * with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getCFI(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         DXFG_OUT const char **cfi);

/**
 * @brief Sets the Classification of Financial Instruments code of the instrument profile.
 *
 * It is a mandatory field for OPTION instruments: it is the only way to distinguish Call/Put, American/European
 * exercise and Cash/Physical delivery. It uses the six-letter code of ISO 10962; `X` can be used extensively and
 * trailing letters can be omitted (they are assumed to be `X`). Examples: `ESNTPB`, `ESXXXX`, `ES`, `OPASPS`.
 *
 * Java: `InstrumentProfile.setCFI(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The CFI, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setCFI(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         const char *value);

/**
 * @brief Returns the International Securities Identifying Number of the instrument profile.
 *
 * It uses the twelve-letter code of ISO 6166. Examples: `DE0007100000`, `US38259P5089`.
 *
 * Java: `InstrumentProfile.getISIN()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] isin The pointer to the variable that receives the ISIN (a UTF-8 string, empty if it is not set; release
 * it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getISIN(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                          DXFG_OUT const char **isin);

/**
 * @brief Sets the International Securities Identifying Number of the instrument profile.
 *
 * It uses the twelve-letter code of ISO 6166. Examples: `DE0007100000`, `US38259P5089`.
 *
 * Java: `InstrumentProfile.setISIN(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The ISIN, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setISIN(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                          const char *value);

/**
 * @brief Returns the Stock Exchange Daily Official List code of the instrument profile.
 *
 * It uses the seven-letter code assigned by the London Stock Exchange. Examples: `2310967`, `5766857`.
 *
 * Java: `InstrumentProfile.getSEDOL()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] sedol The pointer to the variable that receives the SEDOL (a UTF-8 string, empty if it is not set;
 * release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getSEDOL(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                           DXFG_OUT const char **sedol);

/**
 * @brief Sets the Stock Exchange Daily Official List code of the instrument profile.
 *
 * It uses the seven-letter code assigned by the London Stock Exchange. Examples: `2310967`, `5766857`.
 *
 * Java: `InstrumentProfile.setSEDOL(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The SEDOL, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setSEDOL(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                           const char *value);

/**
 * @brief Returns the Committee on Uniform Security Identification Procedures code of the instrument profile.
 *
 * It uses the nine-letter code assigned by the CUSIP Services Bureau. Example: `38259P508`.
 *
 * Java: `InstrumentProfile.getCUSIP()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] cusip The pointer to the variable that receives the CUSIP (a UTF-8 string, empty if it is not set;
 * release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getCUSIP(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                           DXFG_OUT const char **cusip);

/**
 * @brief Sets the Committee on Uniform Security Identification Procedures code of the instrument profile.
 *
 * It uses the nine-letter code assigned by the CUSIP Services Bureau. Example: `38259P508`.
 *
 * Java: `InstrumentProfile.setCUSIP(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The CUSIP, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setCUSIP(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                           const char *value);

/**
 * @brief Returns the Industry Classification Benchmark of the instrument profile.
 *
 * It uses a four-digit number from the ICB catalog. Example: 9535.
 *
 * Java: `InstrumentProfile.getICB()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] icb The pointer to the variable that receives the ICB number; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getICB(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         DXFG_OUT int32_t *icb);

/**
 * @brief Sets the Industry Classification Benchmark of the instrument profile.
 *
 * It uses a four-digit number from the ICB catalog. Example: 9535.
 *
 * Java: `InstrumentProfile.setICB(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The ICB.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setICB(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         int32_t value);

/**
 * @brief Returns the Standard Industrial Classification of the instrument profile.
 *
 * It uses a four-digit number from the SIC catalog. Example: 7371.
 *
 * Java: `InstrumentProfile.getSIC()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] sic The pointer to the variable that receives the SIC number; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getSIC(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         DXFG_OUT int32_t *sic);

/**
 * @brief Sets the Standard Industrial Classification of the instrument profile.
 *
 * It uses a four-digit number from the SIC catalog. Example: 7371.
 *
 * Java: `InstrumentProfile.setSIC(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The SIC.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setSIC(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         int32_t value);

/**
 * @brief Returns the market value multiplier of the instrument profile.
 *
 * Examples: 100, 33.2.
 *
 * Java: `InstrumentProfile.getMultiplier()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] multiplier The pointer to the variable that receives the multiplier; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getMultiplier(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                DXFG_OUT double *multiplier);

/**
 * @brief Sets the market value multiplier of the instrument profile.
 *
 * Examples: 100, 33.2.
 *
 * Java: `InstrumentProfile.setMultiplier(double)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The multiplier.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setMultiplier(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                double value);

/**
 * @brief Returns the product for futures and options on futures (the name of the underlying asset) of the instrument
 * profile.
 *
 * Example: `/YG`.
 *
 * Java: `InstrumentProfile.getProduct()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] product The pointer to the variable that receives the product (a UTF-8 string, empty if it is not set;
 * release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getProduct(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                             DXFG_OUT const char **product);

/**
 * @brief Sets the product for futures and options on futures (the name of the underlying asset) of the instrument
 * profile.
 *
 * Example: `/YG`.
 *
 * Java: `InstrumentProfile.setProduct(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The product, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setProduct(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                             const char *value);

/**
 * @brief Returns the primary underlying symbol for options of the instrument profile.
 *
 * Examples: `C`, `/YGM9`.
 *
 * Java: `InstrumentProfile.getUnderlying()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] underlying The pointer to the variable that receives the underlying (a UTF-8 string, empty if it is not
 * set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getUnderlying(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                DXFG_OUT const char **underlying);

/**
 * @brief Sets the primary underlying symbol for options of the instrument profile.
 *
 * Examples: `C`, `/YGM9`.
 *
 * Java: `InstrumentProfile.setUnderlying(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The underlying, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setUnderlying(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                const char *value);

/**
 * @brief Returns the shares per contract for options of the instrument profile.
 *
 * Examples: 1, 100.
 *
 * Java: `InstrumentProfile.getSPC()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] spc The pointer to the variable that receives the SPC; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getSPC(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         DXFG_OUT double *spc);

/**
 * @brief Sets the shares per contract for options of the instrument profile.
 *
 * Examples: 1, 100.
 *
 * Java: `InstrumentProfile.setSPC(double)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The SPC.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setSPC(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip, double value);

/**
 * @brief Returns the additional underlyings for options, including additional cash of the instrument profile.
 *
 * The format is a list of `<UNDERLYING> <SPC>` items separated by `; ` and sorted by underlying. Examples: `SE 50`,
 * `FIS 53; US$ 45.46`.
 *
 * Java: `InstrumentProfile.getAdditionalUnderlyings()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] additionalUnderlyings The pointer to the variable that receives the additional underlyings (a UTF-8
 * string, empty if it is not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getAdditionalUnderlyings(graal_isolatethread_t *thread,
                                                           dxfg_v4_instrument_profile_t *ip,
                                                           DXFG_OUT const char **additionalUnderlyings);

/**
 * @brief Sets the additional underlyings for options, including additional cash of the instrument profile.
 *
 * The format is a list of `<UNDERLYING> <SPC>` items separated by `; ` and sorted by underlying. Examples: `SE 50`,
 * `FIS 53; US$ 45.46`.
 *
 * Java: `InstrumentProfile.setAdditionalUnderlyings(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The additional underlyings, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setAdditionalUnderlyings(graal_isolatethread_t *thread,
                                                           dxfg_v4_instrument_profile_t *ip, const char *value);

/**
 * @brief Returns the maturity month-year as provided for the corresponding FIX tag (200) of the instrument profile.
 *
 * Depending on the data source it uses `YYYYMM` (year and month), `YYYYMMDD` (full date) or `YYYYMMwN` (week number
 * within the month).
 *
 * Java: `InstrumentProfile.getMMY()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] mmy The pointer to the variable that receives the MMY (a UTF-8 string, empty if it is not set; release it
 * with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getMMY(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         DXFG_OUT const char **mmy);

/**
 * @brief Sets the maturity month-year as provided for the corresponding FIX tag (200) of the instrument profile.
 *
 * Depending on the data source it uses `YYYYMM` (year and month), `YYYYMMDD` (full date) or `YYYYMMwN` (week number
 * within the month).
 *
 * Java: `InstrumentProfile.setMMY(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The MMY, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setMMY(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                         const char *value);

/**
 * @brief Returns the day id of the expiration of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.getExpiration()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] expiration The pointer to the variable that receives the day id of the expiration (the number of days
 * since January 1, 1970, 00:00:00 GMT; 0 if it is not set); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getExpiration(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                DXFG_OUT int32_t *expiration);

/**
 * @brief Sets the day id of the expiration of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.setExpiration(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The expiration.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setExpiration(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                int32_t value);

/**
 * @brief Returns the day id of the last trading day of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.getLastTrade()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] lastTrade The pointer to the variable that receives the day id of the last trading day (the number of
 * days since January 1, 1970, 00:00:00 GMT; 0 if it is not set); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getLastTrade(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                               DXFG_OUT int32_t *lastTrade);

/**
 * @brief Sets the day id of the last trading day of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.setLastTrade(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The last trade.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setLastTrade(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                               int32_t value);

/**
 * @brief Returns the strike price for options of the instrument profile.
 *
 * Examples: 80, 22.5.
 *
 * Java: `InstrumentProfile.getStrike()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] strike The pointer to the variable that receives the strike; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getStrike(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                            DXFG_OUT double *strike);

/**
 * @brief Sets the strike price for options of the instrument profile.
 *
 * Examples: 80, 22.5.
 *
 * Java: `InstrumentProfile.setStrike(double)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The strike.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setStrike(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                            double value);

/**
 * @brief Returns the type of the option of the instrument profile.
 *
 * It is one of `STAN` (Standard Options), `LEAP` (Long-term Equity AnticiPation Securities), `SDO` (Special Dated
 * Options), `BINY` (Binary Options), `FLEX` (FLexible EXchange Options), `VSO` (Variable Start Options), `RNGE`
 * (Range).
 *
 * Java: `InstrumentProfile.getOptionType()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] optionType The pointer to the variable that receives the option type (a UTF-8 string, empty if it is not
 * set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getOptionType(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                DXFG_OUT const char **optionType);

/**
 * @brief Sets the type of the option of the instrument profile.
 *
 * It is one of `STAN` (Standard Options), `LEAP` (Long-term Equity AnticiPation Securities), `SDO` (Special Dated
 * Options), `BINY` (Binary Options), `FLEX` (FLexible EXchange Options), `VSO` (Variable Start Options), `RNGE`
 * (Range).
 *
 * Java: `InstrumentProfile.setOptionType(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The option type, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setOptionType(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                const char *value);

/**
 * @brief Returns the expiration cycle style of the instrument profile.
 *
 * Examples: `Weeklys`, `Quarterlys`.
 *
 * Java: `InstrumentProfile.getExpirationStyle()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] expirationStyle The pointer to the variable that receives the expiration style (a UTF-8 string, empty if
 * it is not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getExpirationStyle(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                     DXFG_OUT const char **expirationStyle);

/**
 * @brief Sets the expiration cycle style of the instrument profile.
 *
 * Examples: `Weeklys`, `Quarterlys`.
 *
 * Java: `InstrumentProfile.setExpirationStyle(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The expiration style, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setExpirationStyle(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                     const char *value);

/**
 * @brief Returns the settlement price determination style of the instrument profile.
 *
 * Examples: `Open`, `Close`.
 *
 * Java: `InstrumentProfile.getSettlementStyle()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] settlementStyle The pointer to the variable that receives the settlement style (a UTF-8 string, empty if
 * it is not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getSettlementStyle(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                     DXFG_OUT const char **settlementStyle);

/**
 * @brief Sets the settlement price determination style of the instrument profile.
 *
 * Examples: `Open`, `Close`.
 *
 * Java: `InstrumentProfile.setSettlementStyle(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The settlement style, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setSettlementStyle(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                     const char *value);

/**
 * @brief Returns the minimum allowed price increments with the corresponding price ranges of the instrument profile.
 *
 * The format is an increment, or a list of `<INCREMENT> <UPPER_LIMIT>` ranges separated by `; ` and sorted by upper
 * limit. Examples: `0.25`, `0.01 3; 0.05`.
 *
 * Java: `InstrumentProfile.getPriceIncrements()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] priceIncrements The pointer to the variable that receives the price increments (a UTF-8 string, empty if
 * it is not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getPriceIncrements(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                     DXFG_OUT const char **priceIncrements);

/**
 * @brief Sets the minimum allowed price increments with the corresponding price ranges of the instrument profile.
 *
 * The format is an increment, or a list of `<INCREMENT> <UPPER_LIMIT>` ranges separated by `; ` and sorted by upper
 * limit. Examples: `0.25`, `0.01 3; 0.05`.
 *
 * Java: `InstrumentProfile.setPriceIncrements(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The price increments, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setPriceIncrements(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                     const char *value);

/**
 * @brief Returns the trading hours specification of the instrument profile.
 *
 * The format is the one of `Schedule.getInstance(String)` of the dxFeed API.
 *
 * Java: `InstrumentProfile.getTradingHours()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] tradingHours The pointer to the variable that receives the trading hours (a UTF-8 string, empty if it is
 * not set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getTradingHours(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  DXFG_OUT const char **tradingHours);

/**
 * @brief Sets the trading hours specification of the instrument profile.
 *
 * The format is the one of `Schedule.getInstance(String)` of the dxFeed API.
 *
 * Java: `InstrumentProfile.setTradingHours(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The trading hours, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setTradingHours(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  const char *value);

/**
 * @brief Returns the value of a field of the instrument profile by the name of the field, in the text form.
 *
 * The field is a standard field of the profile (TYPE, SYMBOL, DESCRIPTION, LOCAL_SYMBOL, LOCAL_DESCRIPTION, COUNTRY,
 * OPOL, EXCHANGE_DATA, EXCHANGES, CURRENCY, BASE_CURRENCY, CFI, ISIN, SEDOL, CUSIP, ICB, SIC, MULTIPLIER, PRODUCT,
 * UNDERLYING, SPC, ADDITIONAL_UNDERLYINGS, MMY, EXPIRATION, LAST_TRADE, STRIKE, OPTION_TYPE, EXPIRATION_STYLE,
 * SETTLEMENT_STYLE, PRICE_INCREMENTS, TRADING_HOURS) or, for any other name, a custom field. The numeric fields are
 * formatted like dxfg_v4_InstrumentProfileField_formatNumber() does, the date fields like
 * dxfg_v4_InstrumentProfileField_formatDate() does (a zero value is an empty string).
 *
 * Java: `InstrumentProfile.getField(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string; the standard field names are case sensitive.
 * @param[out] field The pointer to the variable that receives the value (a UTF-8 string, empty if a custom field is not
 * set; release it with dxfg_v4_String_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getField(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                           const char *name, DXFG_OUT const char **field);

/**
 * @brief Sets the value of a field of the instrument profile by the name of the field, in the text form.
 *
 * The field is a standard field of the profile (TYPE, SYMBOL, DESCRIPTION, LOCAL_SYMBOL, LOCAL_DESCRIPTION, COUNTRY,
 * OPOL, EXCHANGE_DATA, EXCHANGES, CURRENCY, BASE_CURRENCY, CFI, ISIN, SEDOL, CUSIP, ICB, SIC, MULTIPLIER, PRODUCT,
 * UNDERLYING, SPC, ADDITIONAL_UNDERLYINGS, MMY, EXPIRATION, LAST_TRADE, STRIKE, OPTION_TYPE, EXPIRATION_STYLE,
 * SETTLEMENT_STYLE, PRICE_INCREMENTS, TRADING_HOURS) or, for any other name, a custom field. The text of a numeric or a
 * date field is parsed like dxfg_v4_InstrumentProfileField_parseNumber() and dxfg_v4_InstrumentProfileField_parseDate()
 * do; the function fails if the text has a wrong format.
 *
 * Java: `InstrumentProfile.setField(String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string; the standard field names are case sensitive.
 * @param[in] value The value, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setField(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                           const char *name, const char *value);

/**
 * @brief Returns the value of a field of the instrument profile by the name of the field, in the numeric form.
 *
 * The standard numeric fields are ICB, SIC, MULTIPLIER, SPC, EXPIRATION, LAST_TRADE and STRIKE; the function fails for
 * the other standard fields. A custom field that is not set gives 0, a custom field whose text looks like a date
 * (`yyyy-MM-dd`) gives the day id, any other text is parsed as a number.
 *
 * Java: `InstrumentProfile.getNumericField(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] numericField The pointer to the variable that receives the value; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getNumericField(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  const char *name, DXFG_OUT double *numericField);

/**
 * @brief Sets the value of a field of the instrument profile by the name of the field, in the numeric form.
 *
 * The function fails for the standard fields that are not numeric (see dxfg_v4_InstrumentProfile_getNumericField()). A
 * custom field is stored as the text of the number.
 *
 * Java: `InstrumentProfile.setNumericField(String, double)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The value.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setNumericField(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                                  const char *name, double value);

/**
 * @brief Returns the day id of a date field of the instrument profile by the name of the field.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT. A custom field that is not set gives 0.
 *
 * Java: `InstrumentProfile.getDateField(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] dateField The pointer to the variable that receives the day id (the number of days since January 1, 1970,
 * 00:00:00 GMT; 0 if it is not set); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getDateField(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                               const char *name, DXFG_OUT int32_t *dateField);

/**
 * @brief Sets the day id of a date field of the instrument profile by the name of the field.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT. A custom field is stored as the text
 * `yyyy-MM-dd` (an empty string for 0).
 *
 * Java: `InstrumentProfile.setDateField(String, int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The day id (the number of days since January 1, 1970, 00:00:00 GMT).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_setDateField(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip,
                                               const char *name, int32_t value);

/**
 * @brief Returns the names of the non-empty custom fields of the instrument profile.
 *
 * Java: `InstrumentProfile.addNonEmptyCustomFieldNames(Collection)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[out] stringList The pointer to the variable that receives the list of the names (release it with
 * dxfg_v4_string_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_getNonEmptyCustomFieldNames(graal_isolatethread_t *thread,
                                                              dxfg_v4_instrument_profile_t *ip,
                                                              DXFG_OUT dxfg_v4_string_list **stringList);

/**
 * @brief Releases an instrument profile handle.
 *
 * It is the same as dxfg_v4_JavaObjectHandler_release() for the handle.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfile_release(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_t *ip);

/**
 * @brief Formats a number to the text form that is used in the instrument profile fields.
 *
 * The text has no digit grouping; 0 is formatted as an empty string.
 *
 * Java: `InstrumentProfileField.formatNumber(double)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] number The number to format.
 * @param[out] result The formatted number, a UTF-8 string; release it with dxfg_v4_String_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileField_formatNumber(graal_isolatethread_t *thread, double number,
                                                    DXFG_OUT char** result);

/**
 * @brief Parses a number from the text form that is used in the instrument profile fields.
 *
 * A NULL or an empty string gives 0. The function fails if the text is not a number.
 *
 * Java: `InstrumentProfileField.parseNumber(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The text to parse, a UTF-8 string.
 * @param[out] number The parsed number.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileField_parseNumber(graal_isolatethread_t *thread, const char* string,
                                                   DXFG_OUT double* number);

/**
 * @brief Formats a day id to the text form that is used in the instrument profile fields.
 *
 * The text is `yyyy-MM-dd` for the GMT time zone; 0 is formatted as an empty string.
 *
 * Java: `InstrumentProfileField.formatDate(int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] date The day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @param[out] result The formatted date, a UTF-8 string; release it with dxfg_v4_String_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileField_formatDate(graal_isolatethread_t *thread, int32_t date, DXFG_OUT char** result);

/**
 * @brief Parses a day id from the text form that is used in the instrument profile fields.
 *
 * The text is `yyyy-MM-dd` for the GMT time zone. A NULL or an empty string gives 0. The function fails if the text is
 * not a date of this format.
 *
 * Java: `InstrumentProfileField.parseDate(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The text to parse, a UTF-8 string.
 * @param[out] date The parsed day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileField_parseDate(graal_isolatethread_t *thread, const char* string,
                                                 DXFG_OUT int32_t* date);

/**
 * @brief Creates empty custom fields.
 *
 * Java: `new InstrumentProfileCustomFields()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] custom_fields The handle of the new custom fields; release it with dxfg_v4_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_new(graal_isolatethread_t *thread,
                                                  DXFG_OUT dxfg_v4_instrument_profile_custom_fields_t **custom_fields);

/**
 * @brief Creates custom fields as a copy of other custom fields.
 *
 * Java: `new InstrumentProfileCustomFields(InstrumentProfileCustomFields)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The handle of the custom fields to copy.
 * @param[out] custom_fields The handle of the new custom fields; release it with dxfg_v4_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_newCopy(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_custom_fields_t * source,
    DXFG_OUT dxfg_v4_instrument_profile_custom_fields_t **custom_fields);

/**
 * @brief Creates custom fields from a list of strings: the names and the values, consecutive.
 *
 * The list is name 0, value 0, name 1, value 1 and so on. If the list is NULL or empty, empty custom fields are
 * created.
 *
 * Java: `new InstrumentProfileCustomFields(String[])`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields_string_list The list of the UTF-8 strings: names and values, consecutive.
 * @param[out] custom_fields The handle of the new custom fields; release it with dxfg_v4_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_newWithNamesAndValues(
    graal_isolatethread_t *thread, const dxfg_v4_string_list* custom_fields_string_list,
    DXFG_OUT dxfg_v4_instrument_profile_custom_fields_t **custom_fields);

/**
 * @brief Returns the value of a custom field by the name.
 *
 * Java: `InstrumentProfileCustomFields.getField(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] result The value, a UTF-8 string; release it with dxfg_v4_String_release(). It is NULL if there is no
 * such field.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_getField(graal_isolatethread_t *thread,
                                                       dxfg_v4_instrument_profile_custom_fields_t *custom_fields,
                                                       const char *name, DXFG_OUT char** result);

/**
 * @brief Sets the value of a custom field by the name.
 *
 * An empty value clears the field: it is not listed by dxfg_v4_InstrumentProfileCustomFields_getNonEmptyFieldNames()
 * anymore.
 *
 * Java: `InstrumentProfileCustomFields.setField(String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The value, a UTF-8 string; it must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_setField(graal_isolatethread_t *thread,
                                                       dxfg_v4_instrument_profile_custom_fields_t *custom_fields,
                                                       const char *name, const char* value);

/**
 * @brief Returns the numeric value of a custom field by the name.
 *
 * A field that is not set or is empty gives 0. A field whose text looks like a date (`yyyy-MM-dd`) gives the day id
 * (the number of days since January 1, 1970, 00:00:00 GMT); the text of any other field is parsed as a number.
 *
 * Java: `InstrumentProfileCustomFields.getNumericField(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] result The numeric value of the field.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_getNumericField(graal_isolatethread_t *thread,
                                                              dxfg_v4_instrument_profile_custom_fields_t *custom_fields,
                                                              const char *name, DXFG_OUT double* result);

/**
 * @brief Sets the numeric value of a custom field by the name.
 *
 * The field is stored as the text of the number, see dxfg_v4_InstrumentProfileField_formatNumber().
 *
 * Java: `InstrumentProfileCustomFields.setNumericField(String, double)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The value.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_setNumericField(graal_isolatethread_t *thread,
                                                              dxfg_v4_instrument_profile_custom_fields_t *custom_fields,
                                                              const char *name, double value);

/**
 * @brief Returns the day id of a date custom field by the name.
 *
 * A field that is not set or is empty gives 0. The function fails if the text of the field is not a date, see
 * dxfg_v4_InstrumentProfileField_parseDate().
 *
 * Java: `InstrumentProfileCustomFields.getDateField(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] result The day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_getDateField(graal_isolatethread_t *thread,
                                                           dxfg_v4_instrument_profile_custom_fields_t *custom_fields,
                                                           const char *name, DXFG_OUT int32_t* result);

/**
 * @brief Sets the day id of a date custom field by the name.
 *
 * The field is stored as the text of the date, see dxfg_v4_InstrumentProfileField_formatDate().
 *
 * Java: `InstrumentProfileCustomFields.setDateField(String, int)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_setDateField(graal_isolatethread_t *thread,
                                                           dxfg_v4_instrument_profile_custom_fields_t *custom_fields,
                                                           const char *name, int32_t value);

/**
 * @brief Returns the names of the non-empty custom fields.
 *
 * Java: `InstrumentProfileCustomFields.addNonEmptyFieldNames(Collection)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[out] target_field_names The list of the names (empty if there are none); release it with
 * dxfg_v4_string_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCustomFields_getNonEmptyFieldNames(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_custom_fields_t *custom_fields,
    DXFG_OUT dxfg_v4_string_list** target_field_names);

/**
 * @brief Updates one instrument profile in the collector from a structure.
 *
 * The collector gets a new Java instrument profile built from the structure: the structure is only read and can be
 * released after the call. The custom fields are shared, not copied (see dxfg_v4_instrument_profile_data_t): do not
 * change them through instrument_profile_custom_fields after the call.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfile(InstrumentProfile)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] instrument_profile The structure with the fields of the instrument profile.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_updateInstrumentProfile(
    graal_isolatethread_t *thread, dxfg_v4_ipf_collector_t *collector,
    dxfg_v4_instrument_profile_data_t *instrument_profile);

/**
 * @brief Updates a list of instrument profiles in the collector from structures and assigns them a generation tag.
 *
 * The collector gets new Java instrument profiles built from the structures: the structures are only read and can be
 * released after the call; the custom fields are shared, not copied (see dxfg_v4_instrument_profile_data_t). The
 * removal of an instrument profile is represented by a profile whose type is `REMOVED`.
 *
 * A list with a NULL element, or with a positive size and NULL `elements`, is an error: the function fails before it
 * updates any profile. A NULL list or a list of size 0 updates nothing.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfiles(List, Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] instrument_profiles The list of the structures with the fields of the instrument profiles.
 * @param[in] generation The generation tag, the handle of a Java object, or NULL for none. The tag is used by
 * dxfg_v4_InstrumentProfileCollector_removeGenerations().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles(
    graal_isolatethread_t *thread, dxfg_v4_ipf_collector_t *collector,
    const dxfg_v4_instrument_profile_data_list* instrument_profiles, dxfg_v4_java_object_handler *generation);

/**
 * @brief Returns the next instrument profile of the iterator as a structure.
 *
 * Check dxfg_v4_Iterable_InstrumentProfile_hasNext() first: the function fails if there is no next profile. The strings
 * of the structure are copies; its custom fields handle shares the custom fields of the profile in the collector (see
 * dxfg_v4_instrument_profile_data_t).
 *
 * Java: `Iterator.next()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] iterable_ip The iterator handle.
 * @param[out] instrument_profile The next instrument profile; release it with
 * dxfg_v4_instrument_profile_data_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Iterable_InstrumentProfile_next(graal_isolatethread_t *thread, dxfg_v4_iterable_ip_t *iterable_ip,
                                                DXFG_OUT dxfg_v4_instrument_profile_data_t **instrument_profile);

/**
 * @brief Reads and returns the instrument profiles from an address as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_v4_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL,
 * for example `http://user:password@host:port/path/file.ipf`. The operation updates the values of
 * dxfg_v4_InstrumentProfileReader_getLastModified() and dxfg_v4_InstrumentProfileReader_wasComplete(). The function
 * fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_v4_instrument_profile_data_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readFromFile(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, const char *address,
    DXFG_OUT dxfg_v4_instrument_profile_data_list **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an address with basic authentication credentials as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_v4_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL,
 * for example `http://user:password@host:port/path/file.ipf`. The user and the password take precedence over the
 * authentication information of the URL. The operation updates the values of
 * dxfg_v4_InstrumentProfileReader_getLastModified() and dxfg_v4_InstrumentProfileReader_wasComplete(). The function
 * fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] user The user name of the basic authentication, a UTF-8 string, or NULL.
 * @param[in] password The password of the basic authentication, a UTF-8 string, or NULL.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_v4_instrument_profile_data_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPassword(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, const char *address, const char *user,
    const char *password, DXFG_OUT dxfg_v4_instrument_profile_data_list **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an address with token credentials as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_v4_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL,
 * for example `http://user:password@host:port/path/file.ipf`. The token takes precedence over the authentication
 * information of the URL. The operation updates the values of dxfg_v4_InstrumentProfileReader_getLastModified() and
 * dxfg_v4_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument
 * Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, AuthToken)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] token The token handle, or NULL.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_v4_instrument_profile_data_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readFromFileWithToken(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, const char *address,
    dxfg_v4_auth_token_t *token, DXFG_OUT dxfg_v4_instrument_profile_data_list **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an input stream with plain (not compressed) data as structures.
 *
 * It is deprecated in Java. Use dxfg_v4_InstrumentProfileReader_readWithAddress(). Unlike
 * dxfg_v4_InstrumentProfileReader_readWithAddress(), the function does not decompress the stream and does not report an
 * incomplete one. The function fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_v4_instrument_profile_data_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_read(graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader,
                                             dxfg_v4_input_stream_t* is,
                                             DXFG_OUT dxfg_v4_instrument_profile_data_list **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an input stream as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The stream is not closed. The operation updates the value of
 * dxfg_v4_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument
 * Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[in] address The origin of the stream for debugging purposes, a UTF-8 string.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_v4_instrument_profile_data_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readWithAddress(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, dxfg_v4_input_stream_t* is,
    const char *address, DXFG_OUT dxfg_v4_instrument_profile_data_list **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from a compressed input stream as structures.
 *
 * It is deprecated in Java. The data compression formats "zip" and "gzip" are recognized automatically: for zip the
 * first file entry is read as a plain data stream, other data is parsed as is. The stream is not closed. The operation
 * updates the value of dxfg_v4_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform
 * to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readCompressed(InputStream)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_v4_instrument_profile_data_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_InstrumentProfileReader_readCompressed(
    graal_isolatethread_t *thread, dxfg_v4_instrument_profile_reader_t *reader, dxfg_v4_input_stream_t* is,
    DXFG_OUT dxfg_v4_instrument_profile_data_list **instrument_profiles);

/**
 * @brief Writes the instrument profiles from structures to a file.
 *
 * The name of the file selects the format, as in Java: a name that ends with `.zip` (in any case) gives a zip archive
 * with one entry, named as the file without its directory and the `.zip` suffix; a name that ends with `.gz` gives gzip
 * data; any other name gives plain text. The data is in the Instrument Profile Format: for each type of the profiles a
 * header lists the fields that have values, the custom fields included, and the data ends with the `##COMPLETE` tag.
 * dxfg_v4_InstrumentProfileReader_readFromFile() reads the file back. An existing file is overwritten.
 *
 * The structures are only read, as in dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles(): the custom fields
 * of each structure are the ones its handle holds at the call. All the arguments are checked before the file is opened:
 * a NULL argument or a NULL element of the list fails the call with `java.lang.IllegalArgumentException` and leaves the
 * file as it was. An I/O error, e.g. a missing directory, fails the call with `java.io.IOException` (for a file that
 * cannot be opened, its subclass `java.io.FileNotFoundException`); a file that was opened may be left incomplete.
 *
 * Java: `InstrumentProfileWriter.writeToFile(String, List)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] file The path of the file, a UTF-8 string; must not be NULL.
 * @param[in] instrument_profiles The list of the structures of the instrument profiles; must not be NULL, may be empty.
 * The caller keeps it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_InstrumentProfileWriter_writeToFile(graal_isolatethread_t *thread, const char *file,
                                                    const dxfg_v4_instrument_profile_data_list *instrument_profiles);

/**
 * @brief Writes the instrument profiles from handles to a file.
 *
 * The same as dxfg_v4_InstrumentProfileWriter_writeToFile(), for a list of handles: the format, the errors and the
 * checks are the same. The handles are only read.
 *
 * Java: `InstrumentProfileWriter.writeToFile(String, List)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] file The path of the file, a UTF-8 string; must not be NULL.
 * @param[in] instrument_profiles The list of the handles of the instrument profiles; must not be NULL, may be empty.
 * The caller keeps it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_InstrumentProfileWriter_writeHandlesToFile(graal_isolatethread_t *thread, const char *file,
                                                           const dxfg_v4_instrument_profile_list *instrument_profiles);

/**
 * @brief Writes the instrument profiles from structures to memory, as a list of bytes.
 *
 * The bytes are the content of a file with the given name written by dxfg_v4_InstrumentProfileWriter_writeToFile():
 * the name selects plain text, `.gz` or `.zip`, and gives the zip entry its name (without the directory and `.zip`);
 * nothing is written to a file.
 * dxfg_v4_ByteArrayInputStream_new() and dxfg_v4_InstrumentProfileReader_readWithAddress() read the bytes back. The
 * list is never NULL on success, also for an empty list of profiles (the data still has the `##COMPLETE` tag). Data of
 * more than 2147483639 bytes (`Integer.MAX_VALUE - 8`) does not fit a list: the call fails with
 * `java.lang.IllegalStateException`.
 *
 * The structures are only read, and a NULL argument or a NULL element of the list fails the call with
 * `java.lang.IllegalArgumentException`, as in dxfg_v4_InstrumentProfileWriter_writeToFile().
 *
 * Java: `InstrumentProfileWriter.write(OutputStream, String, List)` with a `ByteArrayOutputStream`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] name The name that selects the format, a UTF-8 string, e.g. `profiles.ipf.gz`; must not be NULL.
 * @param[in] instrument_profiles The list of the structures of the instrument profiles; must not be NULL, may be empty.
 * The caller keeps it.
 * @param[out] bytes The pointer to the variable that receives the list of the bytes (release it with
 * dxfg_v4_byte_list_release()); must not be NULL. On failure the variable is left unchanged.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_InstrumentProfileWriter_writeToBytes(graal_isolatethread_t *thread, const char *name,
                                                     const dxfg_v4_instrument_profile_data_list *instrument_profiles,
                                                     DXFG_OUT dxfg_v4_byte_list **bytes);

/**
 * @brief Writes the instrument profiles from handles to memory, as a list of bytes.
 *
 * The same as dxfg_v4_InstrumentProfileWriter_writeToBytes(), for a list of handles: the format, the errors and the
 * checks are the same. The handles are only read.
 *
 * Java: `InstrumentProfileWriter.write(OutputStream, String, List)` with a `ByteArrayOutputStream`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] name The name that selects the format, a UTF-8 string, e.g. `profiles.ipf.gz`; must not be NULL.
 * @param[in] instrument_profiles The list of the handles of the instrument profiles; must not be NULL, may be empty.
 * The caller keeps it.
 * @param[out] bytes The pointer to the variable that receives the list of the bytes (release it with
 * dxfg_v4_byte_list_release()); must not be NULL. On failure the variable is left unchanged.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @since 3.13.0
 */
int32_t dxfg_v4_InstrumentProfileWriter_writeHandlesToBytes(graal_isolatethread_t *thread, const char *name,
                                                            const dxfg_v4_instrument_profile_list *instrument_profiles,
                                                            DXFG_OUT dxfg_v4_byte_list **bytes);

/**
 * @brief Frees the structure of an instrument profile that the SDK has filled, and releases its custom fields handle.
 *
 * The structure itself is freed as well: use it for a structure that was returned by
 * dxfg_v4_Iterable_InstrumentProfile_next(), not for an element of an array or a list.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrument_profile The structure of the instrument profile; NULL does nothing.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_instrument_profile_data_release(graal_isolatethread_t *thread,
                                                dxfg_v4_instrument_profile_data_t *instrument_profile);

/**
 * @brief Frees a list of the structures of instrument profiles that the SDK has filled, with the structures in it, and
 * releases their custom fields handles.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrument_profiles The list of the structures.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_instrument_profile_data_list_release(graal_isolatethread_t *thread,
                                                     dxfg_v4_instrument_profile_data_list* instrument_profiles);

/** @} */ // end of V4_InstrumentProfile

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_IPF_H
