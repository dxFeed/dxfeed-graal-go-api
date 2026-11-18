// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK IPF functions and types declarations
 */

#ifndef DXFG_IPF_H
#define DXFG_IPF_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"
#include "dxfg_javac.h"

/** @defgroup InstrumentProfile InstrumentProfile
 *  @brief Instrument profiles: reading of instrument profile files, live connections to an instrument profile
 *  source, the collector of the profiles, and the custom fields.
 *  @{
 */

/**
 * @brief A handle of a Java InstrumentProfileConnection object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_InstrumentProfileConnection_createConnection(). Releasing the handle does not close the
 * connection: call dxfg_InstrumentProfileConnection_close().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileConnection.html">Javadoc</a>
 */
typedef struct dxfg_ipf_connection_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_ipf_connection_t;

/**
 * @brief A handle of a Java InstrumentProfileReader object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_InstrumentProfileReader_new().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/InstrumentProfileReader.html">Javadoc</a>
 */
typedef struct dxfg_instrument_profile_reader_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_instrument_profile_reader_t;

/**
 * @brief A handle of a Java PropertyChangeListener object; release it with dxfg_JavaObjectHandler_release().
 *
 * The listener calls a dxfg_ipf_connection_state_change_listener_func when the state of an instrument profile
 * connection changes. Create it with dxfg_IpfPropertyChangeListener_new(). The last handle of the listener that is
 * released deactivates the listener: its function is not called anymore.
 *
 * <a href="https://docs.oracle.com/javase/8/docs/api/java/beans/PropertyChangeListener.html">Javadoc</a>
 */
typedef struct dxfg_ipf_connection_state_change_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_ipf_connection_state_change_listener_t;

/**
 * @brief The state of an instrument profile connection.
 *
 * Java: `InstrumentProfileConnection.State`. The state changes from NOT_CONNECTED to CONNECTING, CONNECTED and
 * COMPLETED; CLOSED can follow any state.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileConnection.State.html">Javadoc</a>
 */
typedef enum dxfg_ipf_connection_state_t {
    /** The connection is not started yet: dxfg_InstrumentProfileConnection_start() was not called. */
    DXFG_IPF_CONNECTION_STATE_NOT_CONNECTED = 0,
    /** The connection is being established. */
    DXFG_IPF_CONNECTION_STATE_CONNECTING,
    /** The connection was established. */
    DXFG_IPF_CONNECTION_STATE_CONNECTED,
    /** The initial snapshot of the instrument profiles was fully read (this state is set only once). */
    DXFG_IPF_CONNECTION_STATE_COMPLETED,
    /** The connection was closed by dxfg_InstrumentProfileConnection_close(). */
    DXFG_IPF_CONNECTION_STATE_CLOSED,
} dxfg_ipf_connection_state_t;

/**
 * @brief The function of a dxfg_ipf_connection_state_change_listener_t: called when the state of the connection
 * changes.
 *
 * The function is called by the Java thread that changes the state, not by the thread that added the listener:
 * CONNECTING and CLOSED by the thread that called dxfg_InstrumentProfileConnection_start() and
 * dxfg_InstrumentProfileConnection_close() respectively, CONNECTED and COMPLETED by the background thread of the
 * connection. The states are passed by value, there is nothing to release.
 *
 * Java: `PropertyChangeListener.propertyChange(PropertyChangeEvent)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] old_state The previous state of the connection.
 * @param[in] new_state The new state of the connection.
 * @param[in] user_data The pointer that was passed to dxfg_IpfPropertyChangeListener_new(); the SDK does not use it.
 */
typedef void (*dxfg_ipf_connection_state_change_listener_func)(graal_isolatethread_t *thread, dxfg_ipf_connection_state_t old_state, dxfg_ipf_connection_state_t new_state, void *user_data);

/**
 * @brief A handle of a Java InstrumentProfileCollector object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_InstrumentProfileCollector_new().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileCollector.html">Javadoc</a>
 */
typedef struct dxfg_ipf_collector_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_ipf_collector_t;


/**
 * @brief A handle of a Java InstrumentProfile object; release it with dxfg_JavaObjectHandler_release().
 *
 * dxfg_InstrumentProfile_release() releases the handle as well. Create a profile with dxfg_InstrumentProfile_new(). A
 * profile that is stored in a collector is the object itself, not a copy: do not modify it.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/InstrumentProfile.html">Javadoc</a>
 */
typedef struct dxfg_instrument_profile_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_instrument_profile_t;

/**
 * @brief A list of instrument profile handles.
 *
 * Release the list and all the profile handles in it with dxfg_CList_InstrumentProfile_release(), or only the list with
 * dxfg_CList_InstrumentProfile_wrapper_release().
 */
typedef struct dxfg_instrument_profile_list {
    /** The number of the elements. */
    int32_t size;
    /** The array of size pointers to the instrument profile handles; NULL if the list is empty. */
    dxfg_instrument_profile_t **elements;
} dxfg_instrument_profile_list;

/**
 * @brief A handle of a Java InstrumentProfileUpdateListener object; release it with dxfg_JavaObjectHandler_release().
 *
 * The listener calls a dxfg_ipf_update_listener_function when the set of the instrument profiles of a collector
 * changes. Create it with dxfg_InstrumentProfileUpdateListener_new(). The last handle of the listener that is released
 * deactivates the listener: its function is not called anymore.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileUpdateListener.html">Javadoc</a>
 */
typedef struct dxfg_ipf_update_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_ipf_update_listener_t;

/**
 * @brief A handle of a Java Iterator object over instrument profiles; release it with dxfg_JavaObjectHandler_release().
 *
 * It is returned by dxfg_InstrumentProfileCollector_view() and passed to a dxfg_ipf_update_listener_function. Read it
 * with dxfg_Iterable_InstrumentProfile_hasNext() and dxfg_Iterable_InstrumentProfile_next() or
 * dxfg_Iterable_InstrumentProfile_next2().
 */
typedef struct dxfg_iterable_ip_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_iterable_ip_t;

/**
 * @brief The function of a dxfg_ipf_update_listener_t: called when the set of the instrument profiles of the collector
 * changes.
 *
 * The function is called by the executor of the collector (a single-threaded executor by default, see
 * dxfg_InstrumentProfileCollector_setExecutor()), not by the thread that added the listener. A listener that is added
 * to a collector that is not empty is called for the existing profiles as well. A removed instrument profile is passed
 * as a profile whose type is `REMOVED`.
 *
 * The SDK releases the passed iterator right after the function returns: the callee must not release it and must not
 * use it after the return. The profiles that the callee gets from it with dxfg_Iterable_InstrumentProfile_next() belong
 * to the callee and are released with dxfg_JavaObjectHandler_release(); the structures that
 * dxfg_Iterable_InstrumentProfile_next2() returns belong to the callee too and are released with
 * dxfg_instrument_profile_free().
 *
 * Java: `InstrumentProfileUpdateListener.instrumentProfilesUpdated(Iterator)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] profiles The iterator over the changed instrument profiles. It belongs to the SDK and is valid during the
 * call only.
 * @param[in] user_data The pointer that was passed to dxfg_InstrumentProfileUpdateListener_new(); the SDK does not use
 * it.
 */
typedef void (*dxfg_ipf_update_listener_function)(graal_isolatethread_t *thread, dxfg_iterable_ip_t *profiles, void *user_data);

/**
 * @brief Creates an instrument profile connection to an address that pushes the profiles into a collector.
 *
 * The address is `<host>:<port>` of a server, a URL, or a file path. The clause `[update=<period>]` can be added at the
 * end of the address to set the update period (see dxfg_InstrumentProfileConnection_getUpdatePeriod()); the default
 * period is 1 minute. The connection is not started: call dxfg_InstrumentProfileConnection_start().
 *
 * Java: `InstrumentProfileConnection.createConnection(String, InstrumentProfileCollector)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_createConnection().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] collector The collector handle that the connection pushes the updates into.
 * @return The connection (release it with dxfg_JavaObjectHandler_release(); that does not close the connection), or
 * NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_createConnection()")
dxfg_ipf_connection_t*                        dxfg_InstrumentProfileConnection_createConnection(graal_isolatethread_t *thread, const char *address, dxfg_ipf_collector_t *collector);
/**
 * @brief Returns the address of the connection, without the additional options given as a part of it.
 *
 * Java: `InstrumentProfileConnection.getAddress()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_getAddress().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return The address (a UTF-8 string; release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_getAddress()")
const char*                                   dxfg_InstrumentProfileConnection_getAddress(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
/**
 * @brief Returns the update period of the connection in milliseconds.
 *
 * It is the period of an update check when the instrument profile source does not support live updates and when the
 * connection is dropped. The default period is 1 minute unless the address of the connection sets another one.
 *
 * Java: `InstrumentProfileConnection.getUpdatePeriod()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_getUpdatePeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return The update period in milliseconds, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to determine
 * if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_getUpdatePeriod()")
int64_t                                       dxfg_InstrumentProfileConnection_getUpdatePeriod(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
/**
 * @brief Changes the update period of the connection.
 *
 * Java: `InstrumentProfileConnection.setUpdatePeriod(long)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_setUpdatePeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] newValue The update period in milliseconds; it must not be negative (the function fails otherwise).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_setUpdatePeriod()")
int32_t                                       dxfg_InstrumentProfileConnection_setUpdatePeriod(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, int64_t newValue);
/**
 * @brief Returns the state of the connection.
 *
 * Java: `InstrumentProfileConnection.getState()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_getState().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return The state, or -1 (it is not a dxfg_ipf_connection_state_t value) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_getState()")
dxfg_ipf_connection_state_t                   dxfg_InstrumentProfileConnection_getState(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
/**
 * @brief Returns the last modification time of the instrument profiles in milliseconds, or 0 if it is unknown.
 *
 * The time is in milliseconds, but the actual granularity of the time is a second.
 *
 * Java: `InstrumentProfileConnection.getLastModified()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_getLastModified().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return The last modification time in milliseconds, or 0 if it is unknown, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_getLastModified()")
int64_t                                       dxfg_InstrumentProfileConnection_getLastModified(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
/**
 * @brief Starts the connection.
 *
 * The state of the connection immediately changes to #DXFG_IPF_CONNECTION_STATE_CONNECTING and the actual connection is
 * established in the background. The function fails if the state is not #DXFG_IPF_CONNECTION_STATE_NOT_CONNECTED (the
 * connection is already started or closed).
 *
 * Java: `InstrumentProfileConnection.start()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_start().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_start()")
int32_t                                       dxfg_InstrumentProfileConnection_start(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
/**
 * @brief Closes the connection.
 *
 * The state of the connection immediately changes to #DXFG_IPF_CONNECTION_STATE_CLOSED and the background update
 * stops: nothing that the connection receives after the call reaches the collector. Its thread and its connection to
 * the server end when the next data arrives (a live server sends a heartbeat every 10 seconds) or at the 60-second read
 * timeout. The collector may still notify its listeners about the updates received before the call: release a listener
 * with dxfg_JavaObjectHandler_release() before freeing its user_data. Closing a closed connection does nothing.
 *
 * Java: `InstrumentProfileConnection.close()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_close()")
int32_t                                       dxfg_InstrumentProfileConnection_close(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
/**
 * @brief Creates a listener of the state changes of a connection that calls a C function.
 *
 * Add the listener with dxfg_InstrumentProfileConnection_addStateChangeListener(). The function is called by a Java
 * thread, see dxfg_ipf_connection_state_change_listener_func. Releasing the last handle of the listener deactivates it:
 * after dxfg_JavaObjectHandler_release() returns the function is not called anymore, so user_data can be freed.
 *
 * Java: `PropertyChangeListener.propertyChange(PropertyChangeEvent)`.
 *
 * @deprecated Use dxfg_v4_IpfPropertyChangeListener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call when the state of the connection changes.
 * @param[in] user_data The pointer that is passed to user_func as is; the SDK does not use it.
 * @return The listener (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_IpfPropertyChangeListener_new()")
dxfg_ipf_connection_state_change_listener_t*  dxfg_IpfPropertyChangeListener_new(graal_isolatethread_t *thread, dxfg_ipf_connection_state_change_listener_func user_func, void *user_data);
/**
 * @brief Adds a listener of the state changes of the connection.
 *
 * The listener is called for the following state changes only, not for the current state. Remove it with
 * dxfg_InstrumentProfileConnection_removeStateChangeListener().
 *
 * Java: `InstrumentProfileConnection.addStateChangeListener(PropertyChangeListener)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_addStateChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] listener The listener handle, see dxfg_IpfPropertyChangeListener_new().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_addStateChangeListener()")
int32_t                                       dxfg_InstrumentProfileConnection_addStateChangeListener(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, dxfg_ipf_connection_state_change_listener_t *listener);
/**
 * @brief Removes a listener of the state changes of the connection.
 *
 * The listener handle is not released.
 *
 * Java: `InstrumentProfileConnection.removeStateChangeListener(PropertyChangeListener)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_removeStateChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] listener The listener handle that was added with
 * dxfg_InstrumentProfileConnection_addStateChangeListener().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_removeStateChangeListener()")
int32_t                                       dxfg_InstrumentProfileConnection_removeStateChangeListener(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, dxfg_ipf_connection_state_change_listener_t *listener);
/**
 * @brief Waits until the initial snapshot of the instrument profiles is fully read, for at most the timeout.
 *
 * The calling thread is blocked until the state #DXFG_IPF_CONNECTION_STATE_COMPLETED is reached, the connection is
 * closed, or the timeout elapses.
 *
 * Java: `InstrumentProfileConnection.waitUntilCompleted(long, TimeUnit)`. The timeout is in milliseconds.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileConnection_waitUntilCompleted().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] connection The connection handle.
 * @param[in] timeoutInMs The maximum time to wait in milliseconds.
 * @return 1 if the state #DXFG_IPF_CONNECTION_STATE_COMPLETED was reached, 0 if the time elapsed or the connection was
 * closed before the snapshot was fully read, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to determine
 * if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileConnection_waitUntilCompleted()")
int32_t                                       dxfg_InstrumentProfileConnection_waitUntilCompleted(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, int64_t timeoutInMs);

/**
 * @brief Releases the list structure only and keeps the instrument profile handles in it.
 *
 * The array of the elements and the dxfg_instrument_profile_list structure are freed. The profile handles stay valid:
 * release each of them with dxfg_JavaObjectHandler_release().
 *
 * @deprecated Use dxfg_v4_instrument_profile_list_wrapper_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ips The list of the instrument profile handles.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_instrument_profile_list_wrapper_release()")
int32_t                                       dxfg_CList_InstrumentProfile_wrapper_release(graal_isolatethread_t *thread, dxfg_instrument_profile_list *ips);
/**
 * @brief Releases the list and all the instrument profile handles in it.
 *
 * @deprecated Use dxfg_v4_instrument_profile_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ips The list of the instrument profile handles.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_instrument_profile_list_release()")
int32_t                                       dxfg_CList_InstrumentProfile_release(graal_isolatethread_t *thread, dxfg_instrument_profile_list *ips);

/**
 * @brief Creates an instrument profile collector.
 *
 * The collector keeps one instrument profile per symbol and notifies its listeners about the updates.
 *
 * Java: `new InstrumentProfileCollector()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The collector (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_new()")
dxfg_ipf_collector_t*                         dxfg_InstrumentProfileCollector_new(graal_isolatethread_t *thread);
/**
 * @brief Returns the last modification time of the instrument profiles of the collector in milliseconds.
 *
 * The time is in milliseconds, but the actual granularity of the time is a second.
 *
 * Java: `InstrumentProfileCollector.getLastUpdateTime()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_getLastUpdateTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @return The last modification time in milliseconds, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to
 * determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_getLastUpdateTime()")
int64_t                                       dxfg_InstrumentProfileCollector_getLastUpdateTime(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector);
/**
 * @brief Updates one instrument profile in the collector.
 *
 * The collector stores the Java instrument profile object behind the handle, not a copy: do not modify the profile
 * after the call. It is the same as dxfg_InstrumentProfileCollector_updateInstrumentProfiles() with one profile and no
 * generation.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfile(InstrumentProfile)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] ip The instrument profile handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandle()")
int32_t                                       dxfg_InstrumentProfileCollector_updateInstrumentProfile(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_instrument_profile_t *ip);
/**
 * @brief Updates a list of instrument profiles in the collector and assigns them a generation tag.
 *
 * The removal of an instrument profile is represented by a profile whose type is `REMOVED`. The collector stores the
 * Java instrument profile objects behind the handles, not copies: do not modify the profiles after the call. The list
 * and its handles still belong to the caller.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfiles(List, Object)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] ips The list of the instrument profile handles.
 * @param[in] generation The generation tag, the handle of a Java object, or NULL for none. The tag is used by
 * dxfg_InstrumentProfileCollector_removeGenerations().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfileHandles()")
int32_t                                       dxfg_InstrumentProfileCollector_updateInstrumentProfiles(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_instrument_profile_list *ips, dxfg_java_object_handler *generation);
/**
 * @brief Removes the instrument profiles whose generation tags are in a list.
 *
 * The profiles without a generation tag are not removed. The time of the call grows with the number of the profiles in
 * the collector: use it sparingly. The list and its handles still belong to the caller.
 *
 * Java: `InstrumentProfileCollector.removeGenerations(Set)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_removeGenerations().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] generations The list of the generation tags, the handles of Java objects.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_removeGenerations()")
int32_t                                       dxfg_InstrumentProfileCollector_removeGenerations(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_java_object_handler_list *generations);
/**
 * @brief Returns an iterator over the instrument profiles of the collector.
 *
 * The iterator is a concurrent view of the set of the profiles. It exposes the profiles that are not removed; however,
 * if the iteration is concurrent with a removal, a removed profile (its type is `REMOVED`) can be exposed.
 *
 * Java: `InstrumentProfileCollector.view()`. The SDK returns the iterator of the view.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_view().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @return The iterator (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_view()")
dxfg_iterable_ip_t*                           dxfg_InstrumentProfileCollector_view(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector);
/**
 * @brief Returns the executor that calls the update listeners of the collector.
 *
 * By default a single-threaded executor is used to process the updates asynchronously.
 *
 * Java: `InstrumentProfileCollector.getExecutor()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_getExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @return The executor (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_getExecutor()")
dxfg_executor_t*                              dxfg_InstrumentProfileCollector_getExecutor(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector);
/**
 * @brief Changes the executor that calls the update listeners of the collector.
 *
 * Java: `InstrumentProfileCollector.setExecutor(Executor)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_setExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] executor The executor handle; it must not be NULL (the function fails otherwise).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_setExecutor()")
int32_t                                       dxfg_InstrumentProfileCollector_setExecutor(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_executor_t *executor);
/**
 * @brief Adds a listener that is notified about any updates in the set of the instrument profiles.
 *
 * If the set of the profiles is not empty, the listener is notified about it as well (by the executor of the
 * collector). Remove the listener with dxfg_InstrumentProfileCollector_removeUpdateListener().
 *
 * Java: `InstrumentProfileCollector.addUpdateListener(InstrumentProfileUpdateListener)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_addUpdateListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] listener The listener handle, see dxfg_InstrumentProfileUpdateListener_new().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_addUpdateListener()")
int32_t                                       dxfg_InstrumentProfileCollector_addUpdateListener(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_ipf_update_listener_t* listener);
/**
 * @brief Removes a listener that was added with dxfg_InstrumentProfileCollector_addUpdateListener().
 *
 * The listener handle is not released.
 *
 * Java: `InstrumentProfileCollector.removeUpdateListener(InstrumentProfileUpdateListener)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_removeUpdateListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] listener The listener handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_removeUpdateListener()")
int32_t                                       dxfg_InstrumentProfileCollector_removeUpdateListener(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_ipf_update_listener_t* listener);

/**
 * @brief Checks whether the iterator has a next instrument profile.
 *
 * Java: `Iterator.hasNext()`.
 *
 * @deprecated Use dxfg_v4_Iterable_InstrumentProfile_hasNext().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] iterable_ip The iterator handle.
 * @return 1 if the iterator has a next instrument profile, 0 if it has not, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Iterable_InstrumentProfile_hasNext()")
int32_t                                       dxfg_Iterable_InstrumentProfile_hasNext(graal_isolatethread_t *thread, dxfg_iterable_ip_t *iterable_ip);
/**
 * @brief Returns the next instrument profile of the iterator as a handle.
 *
 * Check dxfg_Iterable_InstrumentProfile_hasNext() first: the function fails if there is no next profile. The handle
 * refers to the profile object itself (the collector stores the profiles, not copies): do not modify it. Use
 * dxfg_Iterable_InstrumentProfile_next2() to get a copy.
 *
 * Java: `Iterator.next()`.
 *
 * @deprecated Use dxfg_v4_Iterable_InstrumentProfile_nextHandle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] iterable_ip The iterator handle.
 * @return The instrument profile (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Iterable_InstrumentProfile_nextHandle()")
dxfg_instrument_profile_t*                    dxfg_Iterable_InstrumentProfile_next(graal_isolatethread_t *thread, dxfg_iterable_ip_t *iterable_ip);

/**
 * @brief Creates a listener of the updates of a collector that calls a C function.
 *
 * Add the listener with dxfg_InstrumentProfileCollector_addUpdateListener(). The function is called by the executor of
 * the collector, see dxfg_ipf_update_listener_function. Releasing the last handle of the listener deactivates it: after
 * dxfg_JavaObjectHandler_release() returns the function is not called anymore, so user_data can be freed.
 *
 * Java: `InstrumentProfileUpdateListener.instrumentProfilesUpdated(Iterator)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileUpdateListener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call when the set of the instrument profiles changes.
 * @param[in] user_data The pointer that is passed to user_func as is; the SDK does not use it.
 * @return The listener (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileUpdateListener_new()")
dxfg_ipf_update_listener_t*                   dxfg_InstrumentProfileUpdateListener_new(graal_isolatethread_t *thread, dxfg_ipf_update_listener_function user_func, void *user_data);

/**
 * @brief Creates an instrument profile reader.
 *
 * Java: `new InstrumentProfileReader()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The reader (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_new()")
dxfg_instrument_profile_reader_t*             dxfg_InstrumentProfileReader_new(graal_isolatethread_t *thread);
/**
 * @brief Returns the last modification time in milliseconds from the last read from an address, or 0 if it is unknown.
 *
 * Java: `InstrumentProfileReader.getLastModified()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_getLastModified().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @return The last modification time in milliseconds, or 0 if it is unknown, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_getLastModified()")
int64_t                                       dxfg_InstrumentProfileReader_getLastModified(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader);
/**
 * @brief Checks whether the instrument profile file was fully read by the last read operation.
 *
 * Java: `InstrumentProfileReader.wasComplete()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_wasComplete().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @return 1 if the file was fully read, 0 if it was not, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to
 * determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_wasComplete()")
int32_t                                       dxfg_InstrumentProfileReader_wasComplete(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader);
/**
 * @brief Reads and returns the instrument profiles from an address.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The operation updates the values of
 * dxfg_InstrumentProfileReader_getLastModified() and dxfg_InstrumentProfileReader_wasComplete(). The function fails if
 * the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFileAsHandles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @return The list of the instrument profiles (release it and the profiles in it with
 * dxfg_CList_InstrumentProfile_release()), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFileAsHandles()")
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_readFromFile(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address);
/**
 * @brief Reads and returns the instrument profiles from an address with basic authentication credentials.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The user and the password take precedence over the
 * authentication information of the URL. The operation updates the values of
 * dxfg_InstrumentProfileReader_getLastModified() and dxfg_InstrumentProfileReader_wasComplete(). The function fails if
 * the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, String, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPasswordAsHandles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] user The user name of the basic authentication, a UTF-8 string, or NULL.
 * @param[in] password The password of the basic authentication, a UTF-8 string, or NULL.
 * @return The list of the instrument profiles (release it and the profiles in it with
 * dxfg_CList_InstrumentProfile_release()), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPasswordAsHandles()")
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_readFromFile2(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, const char *user, const char *password);
/**
 * @brief Reads and returns the instrument profiles from an address with token credentials.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The token takes precedence over the authentication
 * information of the URL. The operation updates the values of dxfg_InstrumentProfileReader_getLastModified() and
 * dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument Profile
 * Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, AuthToken)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFileWithTokenAsHandles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] token The token handle, or NULL.
 * @return The list of the instrument profiles (release it and the profiles in it with
 * dxfg_CList_InstrumentProfile_release()), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFileWithTokenAsHandles()")
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_readFromFile3(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, dxfg_auth_token_t *token);
/**
 * @brief Converts an address specification into the URL that is read by the
 * dxfg_InstrumentProfileReader_readFromFile*() functions.
 *
 * An address of the form `<host>:<port>` (it contains a colon, no slash and ends with a valid port number) becomes
 * `http://<host>:<port>/ipf/all.ipf.gz`, the query part of the address is kept; any other address is returned as is.
 *
 * Java: `InstrumentProfileReader.resolveSourceURL(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_resolveSourceURL().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] address The address, a UTF-8 string.
 * @return The URL (a UTF-8 string; release it with dxfg_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_resolveSourceURL()")
const char*                                   dxfg_InstrumentProfileReader_resolveSourceURL(graal_isolatethread_t *thread, const char *address);
/**
 * @brief Reads and returns the instrument profiles from an input stream.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The stream is not closed. The operation updates the value of
 * dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument Profile
 * Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readWithAddressAsHandles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[in] address The origin of the stream for debugging purposes, a UTF-8 string.
 * @return The list of the instrument profiles (release it and the profiles in it with
 * dxfg_CList_InstrumentProfile_release()), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readWithAddressAsHandles()")
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_read2(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, const char* address);
/**
 * @brief Reads and returns the instrument profiles from a compressed input stream.
 *
 * It is deprecated in Java; use dxfg_InstrumentProfileReader_read2(). The data compression formats "zip" and "gzip" are
 * recognized automatically: for zip the first file entry is read as a plain data stream, other data is parsed as is.
 * The stream is not closed. The operation updates the value of dxfg_InstrumentProfileReader_wasComplete(). The function
 * fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readCompressed(InputStream)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readCompressedAsHandles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @return The list of the instrument profiles (release it and the profiles in it with
 * dxfg_CList_InstrumentProfile_release()), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readCompressedAsHandles()")
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_readCompressed(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is);
/**
 * @brief Reads and returns the instrument profiles from an input stream with plain (not compressed) data.
 *
 * It is deprecated in Java; use dxfg_InstrumentProfileReader_read2(). Unlike it, the function does not decompress the
 * stream and does not report an incomplete one. The function fails if the data does not conform to the Instrument
 * Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readAsHandles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @return The list of the instrument profiles (release it and the profiles in it with
 * dxfg_CList_InstrumentProfile_release()), or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the
 * exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readAsHandles()")
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_read(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is);
/**
 * @brief Creates an instrument profile with the default values: empty strings and zero numbers.
 *
 * Java: `new InstrumentProfile()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The instrument profile (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_new()")
dxfg_instrument_profile_t*                    dxfg_InstrumentProfile_new(graal_isolatethread_t *thread);
/**
 * @brief Creates an instrument profile as a copy of another instrument profile.
 *
 * Java: `new InstrumentProfile(InstrumentProfile)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_newCopy().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The handle of the instrument profile to copy.
 * @return The instrument profile (release it with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_newCopy()")
dxfg_instrument_profile_t*                    dxfg_InstrumentProfile_new2(graal_isolatethread_t *thread, dxfg_instrument_profile_t* ip);
/**
 * @brief Returns the type of the instrument of the instrument profile.
 *
 * It takes precedence over the other fields in conflict cases. It is a mandatory field and must not be empty. Examples:
 * `STOCK`, `FUTURE`, `OPTION`.
 *
 * Java: `InstrumentProfile.getType()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getType().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The type (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getType()")
const char*                                   dxfg_InstrumentProfile_getType(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the type of the instrument of the instrument profile.
 *
 * It takes precedence over the other fields in conflict cases. It is a mandatory field and must not be empty. Examples:
 * `STOCK`, `FUTURE`, `OPTION`.
 *
 * Java: `InstrumentProfile.setType(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setType().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The type, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setType()")
int32_t                                       dxfg_InstrumentProfile_setType(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the identifier of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. It is a mandatory field and must not be empty. Examples:
 * `GOOG`, `/YGM9`, `.ZYEAD`.
 *
 * Java: `InstrumentProfile.getSymbol()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The symbol (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getSymbol()")
const char*                                   dxfg_InstrumentProfile_getSymbol(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the identifier of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. It is a mandatory field and must not be empty. Examples:
 * `GOOG`, `/YGM9`, `.ZYEAD`.
 *
 * Java: `InstrumentProfile.setSymbol(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The symbol, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setSymbol()")
int32_t                                       dxfg_InstrumentProfile_setSymbol(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the description of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. Examples: `Google Inc.`, `Mini Gold
 * Futures,Jun-2009,ETH`.
 *
 * Java: `InstrumentProfile.getDescription()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getDescription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The description (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getDescription()")
const char*                                   dxfg_InstrumentProfile_getDescription(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the description of the instrument of the instrument profile.
 *
 * It is preferably an international one in the Latin alphabet. Examples: `Google Inc.`, `Mini Gold
 * Futures,Jun-2009,ETH`.
 *
 * Java: `InstrumentProfile.setDescription(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setDescription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The description, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setDescription()")
int32_t                                       dxfg_InstrumentProfile_setDescription(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the identifier of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the symbol.
 *
 * Java: `InstrumentProfile.getLocalSymbol()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getLocalSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The local symbol (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getLocalSymbol()")
const char*                                   dxfg_InstrumentProfile_getLocalSymbol(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the identifier of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the symbol.
 *
 * Java: `InstrumentProfile.setLocalSymbol(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setLocalSymbol().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The local symbol, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setLocalSymbol()")
int32_t                                       dxfg_InstrumentProfile_setLocalSymbol(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the description of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the description.
 *
 * Java: `InstrumentProfile.getLocalDescription()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getLocalDescription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The local description (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or
 * NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getLocalDescription()")
const char*                                   dxfg_InstrumentProfile_getLocalDescription(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the description of the instrument in the national language of the instrument profile.
 *
 * It is empty if it is the same as the description.
 *
 * Java: `InstrumentProfile.setLocalDescription(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setLocalDescription().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The local description, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setLocalDescription()")
int32_t                                       dxfg_InstrumentProfile_setLocalDescription(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the country of origin (incorporation) of the corresponding company or parent entity of the instrument
 * profile.
 *
 * It uses the two-letter country code of ISO 3166-1. Examples: `US`, `RU`.
 *
 * Java: `InstrumentProfile.getCountry()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getCountry().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The country (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getCountry()")
const char*                                   dxfg_InstrumentProfile_getCountry(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the country of origin (incorporation) of the corresponding company or parent entity of the instrument
 * profile.
 *
 * It uses the two-letter country code of ISO 3166-1. Examples: `US`, `RU`.
 *
 * Java: `InstrumentProfile.setCountry(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setCountry().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The country, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setCountry()")
int32_t                                       dxfg_InstrumentProfile_setCountry(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the official Place Of Listing, the organization that has listed the instrument of the instrument
 * profile.
 *
 * An instrument with multiple listings uses a separate profile for each listing. It uses the Market Identifier Code
 * (MIC) of ISO 10383. Examples: `XNAS`, `RTSX`.
 *
 * Java: `InstrumentProfile.getOPOL()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getOPOL().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The OPOL (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getOPOL()")
const char*                                   dxfg_InstrumentProfile_getOPOL(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the official Place Of Listing, the organization that has listed the instrument of the instrument profile.
 *
 * An instrument with multiple listings uses a separate profile for each listing. It uses the Market Identifier Code
 * (MIC) of ISO 10383. Examples: `XNAS`, `RTSX`.
 *
 * Java: `InstrumentProfile.setOPOL(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setOPOL().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The OPOL, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setOPOL()")
int32_t                                       dxfg_InstrumentProfile_setOPOL(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the exchange-specific data required to identify the instrument when communicating with the exchange of
 * the instrument profile.
 *
 * It uses an exchange-specific format.
 *
 * Java: `InstrumentProfile.getExchangeData()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getExchangeData().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The exchange data (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getExchangeData()")
const char*                                   dxfg_InstrumentProfile_getExchangeData(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the exchange-specific data required to identify the instrument when communicating with the exchange of
 * the instrument profile.
 *
 * It uses an exchange-specific format.
 *
 * Java: `InstrumentProfile.setExchangeData(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setExchangeData().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The exchange data, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setExchangeData()")
int32_t                                       dxfg_InstrumentProfile_setExchangeData(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the list of the exchanges where the instrument is quoted or traded of the instrument profile.
 *
 * The list is a sequence of MICs separated by semicolons and sorted by MIC. Example: `ARCX;CBSX;XNAS;XNYS`.
 *
 * Java: `InstrumentProfile.getExchanges()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getExchanges().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The exchanges (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getExchanges()")
const char*                                   dxfg_InstrumentProfile_getExchanges(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the list of the exchanges where the instrument is quoted or traded of the instrument profile.
 *
 * The list is a sequence of MICs separated by semicolons and sorted by MIC. Example: `ARCX;CBSX;XNAS;XNYS`.
 *
 * Java: `InstrumentProfile.setExchanges(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setExchanges().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The exchanges, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setExchanges()")
int32_t                                       dxfg_InstrumentProfile_setExchanges(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the currency of quotation, pricing and trading of the instrument profile.
 *
 * It uses the three-letter currency code of ISO 4217. Examples: `USD`, `RUB`.
 *
 * Java: `InstrumentProfile.getCurrency()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getCurrency().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The currency (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getCurrency()")
const char*                                   dxfg_InstrumentProfile_getCurrency(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the currency of quotation, pricing and trading of the instrument profile.
 *
 * It uses the three-letter currency code of ISO 4217. Examples: `USD`, `RUB`.
 *
 * Java: `InstrumentProfile.setCurrency(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setCurrency().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The currency, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setCurrency()")
int32_t                                       dxfg_InstrumentProfile_setCurrency(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the base currency of a currency pair (FOREX instruments) of the instrument profile.
 *
 * It uses the three-letter currency code like the currency.
 *
 * Java: `InstrumentProfile.getBaseCurrency()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getBaseCurrency().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The base currency (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getBaseCurrency()")
const char*                                   dxfg_InstrumentProfile_getBaseCurrency(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the base currency of a currency pair (FOREX instruments) of the instrument profile.
 *
 * It uses the three-letter currency code like the currency.
 *
 * Java: `InstrumentProfile.setBaseCurrency(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setBaseCurrency().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The base currency, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setBaseCurrency()")
int32_t                                       dxfg_InstrumentProfile_setBaseCurrency(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the Classification of Financial Instruments code of the instrument profile.
 *
 * It is a mandatory field for OPTION instruments: it is the only way to distinguish Call/Put, American/European
 * exercise and Cash/Physical delivery. It uses the six-letter code of ISO 10962; `X` can be used extensively and
 * trailing letters can be omitted (they are assumed to be `X`). Examples: `ESNTPB`, `ESXXXX`, `ES`, `OPASPS`.
 *
 * Java: `InstrumentProfile.getCFI()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getCFI().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The CFI (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getCFI()")
const char*                                   dxfg_InstrumentProfile_getCFI(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the Classification of Financial Instruments code of the instrument profile.
 *
 * It is a mandatory field for OPTION instruments: it is the only way to distinguish Call/Put, American/European
 * exercise and Cash/Physical delivery. It uses the six-letter code of ISO 10962; `X` can be used extensively and
 * trailing letters can be omitted (they are assumed to be `X`). Examples: `ESNTPB`, `ESXXXX`, `ES`, `OPASPS`.
 *
 * Java: `InstrumentProfile.setCFI(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setCFI().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The CFI, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setCFI()")
int32_t                                       dxfg_InstrumentProfile_setCFI(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the International Securities Identifying Number of the instrument profile.
 *
 * It uses the twelve-letter code of ISO 6166. Examples: `DE0007100000`, `US38259P5089`.
 *
 * Java: `InstrumentProfile.getISIN()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getISIN().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The ISIN (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getISIN()")
const char*                                   dxfg_InstrumentProfile_getISIN(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the International Securities Identifying Number of the instrument profile.
 *
 * It uses the twelve-letter code of ISO 6166. Examples: `DE0007100000`, `US38259P5089`.
 *
 * Java: `InstrumentProfile.setISIN(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setISIN().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The ISIN, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setISIN()")
int32_t                                       dxfg_InstrumentProfile_setISIN(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the Stock Exchange Daily Official List code of the instrument profile.
 *
 * It uses the seven-letter code assigned by the London Stock Exchange. Examples: `2310967`, `5766857`.
 *
 * Java: `InstrumentProfile.getSEDOL()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getSEDOL().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The SEDOL (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getSEDOL()")
const char*                                   dxfg_InstrumentProfile_getSEDOL(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the Stock Exchange Daily Official List code of the instrument profile.
 *
 * It uses the seven-letter code assigned by the London Stock Exchange. Examples: `2310967`, `5766857`.
 *
 * Java: `InstrumentProfile.setSEDOL(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setSEDOL().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The SEDOL, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setSEDOL()")
int32_t                                       dxfg_InstrumentProfile_setSEDOL(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the Committee on Uniform Security Identification Procedures code of the instrument profile.
 *
 * It uses the nine-letter code assigned by the CUSIP Services Bureau. Example: `38259P508`.
 *
 * Java: `InstrumentProfile.getCUSIP()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getCUSIP().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The CUSIP (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getCUSIP()")
const char*                                   dxfg_InstrumentProfile_getCUSIP(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the Committee on Uniform Security Identification Procedures code of the instrument profile.
 *
 * It uses the nine-letter code assigned by the CUSIP Services Bureau. Example: `38259P508`.
 *
 * Java: `InstrumentProfile.setCUSIP(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setCUSIP().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The CUSIP, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setCUSIP()")
int32_t                                       dxfg_InstrumentProfile_setCUSIP(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the Industry Classification Benchmark of the instrument profile.
 *
 * It uses a four-digit number from the ICB catalog. Example: 9535.
 *
 * Java: `InstrumentProfile.getICB()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getICB().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The ICB number, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getICB()")
int32_t                                       dxfg_InstrumentProfile_getICB(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the Industry Classification Benchmark of the instrument profile.
 *
 * It uses a four-digit number from the ICB catalog. Example: 9535.
 *
 * Java: `InstrumentProfile.setICB(int)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setICB().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The ICB.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setICB()")
int32_t                                       dxfg_InstrumentProfile_setICB(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, int32_t value);
/**
 * @brief Returns the Standard Industrial Classification of the instrument profile.
 *
 * It uses a four-digit number from the SIC catalog. Example: 7371.
 *
 * Java: `InstrumentProfile.getSIC()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getSIC().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The SIC number, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getSIC()")
int32_t                                       dxfg_InstrumentProfile_getSIC(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the Standard Industrial Classification of the instrument profile.
 *
 * It uses a four-digit number from the SIC catalog. Example: 7371.
 *
 * Java: `InstrumentProfile.setSIC(int)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setSIC().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The SIC.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setSIC()")
int32_t                                       dxfg_InstrumentProfile_setSIC(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, int32_t value);
/**
 * @brief Returns the market value multiplier of the instrument profile.
 *
 * Examples: 100, 33.2.
 *
 * Java: `InstrumentProfile.getMultiplier()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getMultiplier().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The multiplier, or negative infinity on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getMultiplier()")
double                                        dxfg_InstrumentProfile_getMultiplier(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the market value multiplier of the instrument profile.
 *
 * Examples: 100, 33.2.
 *
 * Java: `InstrumentProfile.setMultiplier(double)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setMultiplier().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The multiplier.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setMultiplier()")
int32_t                                       dxfg_InstrumentProfile_setMultiplier(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, double value);
/**
 * @brief Returns the product for futures and options on futures (the name of the underlying asset) of the instrument
 * profile.
 *
 * Example: `/YG`.
 *
 * Java: `InstrumentProfile.getProduct()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getProduct().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The product (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getProduct()")
const char*                                   dxfg_InstrumentProfile_getProduct(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the product for futures and options on futures (the name of the underlying asset) of the instrument
 * profile.
 *
 * Example: `/YG`.
 *
 * Java: `InstrumentProfile.setProduct(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setProduct().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The product, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setProduct()")
int32_t                                       dxfg_InstrumentProfile_setProduct(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the primary underlying symbol for options of the instrument profile.
 *
 * Examples: `C`, `/YGM9`.
 *
 * Java: `InstrumentProfile.getUnderlying()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getUnderlying().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The underlying (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getUnderlying()")
const char*                                   dxfg_InstrumentProfile_getUnderlying(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the primary underlying symbol for options of the instrument profile.
 *
 * Examples: `C`, `/YGM9`.
 *
 * Java: `InstrumentProfile.setUnderlying(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setUnderlying().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The underlying, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setUnderlying()")
int32_t                                       dxfg_InstrumentProfile_setUnderlying(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the shares per contract for options of the instrument profile.
 *
 * Examples: 1, 100.
 *
 * Java: `InstrumentProfile.getSPC()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getSPC().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The SPC, or negative infinity on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getSPC()")
double                                        dxfg_InstrumentProfile_getSPC(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the shares per contract for options of the instrument profile.
 *
 * Examples: 1, 100.
 *
 * Java: `InstrumentProfile.setSPC(double)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setSPC().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The SPC.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setSPC()")
int32_t                                       dxfg_InstrumentProfile_setSPC(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, double value);
/**
 * @brief Returns the additional underlyings for options, including additional cash of the instrument profile.
 *
 * The format is a list of `<UNDERLYING> <SPC>` items separated by `; ` and sorted by underlying. Examples: `SE 50`,
 * `FIS 53; US$ 45.46`.
 *
 * Java: `InstrumentProfile.getAdditionalUnderlyings()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getAdditionalUnderlyings().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The additional underlyings (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()),
 * or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getAdditionalUnderlyings()")
const char*                                   dxfg_InstrumentProfile_getAdditionalUnderlyings(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the additional underlyings for options, including additional cash of the instrument profile.
 *
 * The format is a list of `<UNDERLYING> <SPC>` items separated by `; ` and sorted by underlying. Examples: `SE 50`,
 * `FIS 53; US$ 45.46`.
 *
 * Java: `InstrumentProfile.setAdditionalUnderlyings(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setAdditionalUnderlyings().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The additional underlyings, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setAdditionalUnderlyings()")
int32_t                                       dxfg_InstrumentProfile_setAdditionalUnderlyings(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the maturity month-year as provided for the corresponding FIX tag (200) of the instrument profile.
 *
 * Depending on the data source it uses `YYYYMM` (year and month), `YYYYMMDD` (full date) or `YYYYMMwN` (week number
 * within the month).
 *
 * Java: `InstrumentProfile.getMMY()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getMMY().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The MMY (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getMMY()")
const char*                                   dxfg_InstrumentProfile_getMMY(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the maturity month-year as provided for the corresponding FIX tag (200) of the instrument profile.
 *
 * Depending on the data source it uses `YYYYMM` (year and month), `YYYYMMDD` (full date) or `YYYYMMwN` (week number
 * within the month).
 *
 * Java: `InstrumentProfile.setMMY(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setMMY().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The MMY, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setMMY()")
int32_t                                       dxfg_InstrumentProfile_setMMY(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the day id of the expiration of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.getExpiration()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getExpiration().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The day id of the expiration (the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set), or
 * -1 on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getExpiration()")
int32_t                                       dxfg_InstrumentProfile_getExpiration(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the day id of the expiration of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.setExpiration(int)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setExpiration().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The expiration.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setExpiration()")
int32_t                                       dxfg_InstrumentProfile_setExpiration(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, int32_t value);
/**
 * @brief Returns the day id of the last trading day of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.getLastTrade()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getLastTrade().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The day id of the last trading day (the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not
 * set), or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getLastTrade()")
int32_t                                       dxfg_InstrumentProfile_getLastTrade(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the day id of the last trading day of the instrument profile.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set.
 *
 * Java: `InstrumentProfile.setLastTrade(int)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setLastTrade().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The last trade.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setLastTrade()")
int32_t                                       dxfg_InstrumentProfile_setLastTrade(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, int32_t value);
/**
 * @brief Returns the strike price for options of the instrument profile.
 *
 * Examples: 80, 22.5.
 *
 * Java: `InstrumentProfile.getStrike()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getStrike().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The strike, or negative infinity on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getStrike()")
double                                        dxfg_InstrumentProfile_getStrike(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the strike price for options of the instrument profile.
 *
 * Examples: 80, 22.5.
 *
 * Java: `InstrumentProfile.setStrike(double)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setStrike().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The strike.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setStrike()")
int32_t                                       dxfg_InstrumentProfile_setStrike(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, double value);
/**
 * @brief Returns the type of the option of the instrument profile.
 *
 * It is one of `STAN` (Standard Options), `LEAP` (Long-term Equity AnticiPation Securities), `SDO` (Special Dated
 * Options), `BINY` (Binary Options), `FLEX` (FLexible EXchange Options), `VSO` (Variable Start Options), `RNGE`
 * (Range).
 *
 * Java: `InstrumentProfile.getOptionType()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getOptionType().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The option type (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getOptionType()")
const char*                                   dxfg_InstrumentProfile_getOptionType(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the type of the option of the instrument profile.
 *
 * It is one of `STAN` (Standard Options), `LEAP` (Long-term Equity AnticiPation Securities), `SDO` (Special Dated
 * Options), `BINY` (Binary Options), `FLEX` (FLexible EXchange Options), `VSO` (Variable Start Options), `RNGE`
 * (Range).
 *
 * Java: `InstrumentProfile.setOptionType(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setOptionType().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The option type, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setOptionType()")
int32_t                                       dxfg_InstrumentProfile_setOptionType(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the expiration cycle style of the instrument profile.
 *
 * Examples: `Weeklys`, `Quarterlys`.
 *
 * Java: `InstrumentProfile.getExpirationStyle()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getExpirationStyle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The expiration style (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL
 * on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getExpirationStyle()")
const char*                                   dxfg_InstrumentProfile_getExpirationStyle(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the expiration cycle style of the instrument profile.
 *
 * Examples: `Weeklys`, `Quarterlys`.
 *
 * Java: `InstrumentProfile.setExpirationStyle(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setExpirationStyle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The expiration style, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setExpirationStyle()")
int32_t                                       dxfg_InstrumentProfile_setExpirationStyle(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the settlement price determination style of the instrument profile.
 *
 * Examples: `Open`, `Close`.
 *
 * Java: `InstrumentProfile.getSettlementStyle()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getSettlementStyle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The settlement style (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL
 * on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getSettlementStyle()")
const char*                                   dxfg_InstrumentProfile_getSettlementStyle(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the settlement price determination style of the instrument profile.
 *
 * Examples: `Open`, `Close`.
 *
 * Java: `InstrumentProfile.setSettlementStyle(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setSettlementStyle().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The settlement style, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setSettlementStyle()")
int32_t                                       dxfg_InstrumentProfile_setSettlementStyle(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the minimum allowed price increments with the corresponding price ranges of the instrument profile.
 *
 * The format is an increment, or a list of `<INCREMENT> <UPPER_LIMIT>` ranges separated by `; ` and sorted by upper
 * limit. Examples: `0.25`, `0.01 3; 0.05`.
 *
 * Java: `InstrumentProfile.getPriceIncrements()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getPriceIncrements().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The price increments (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL
 * on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getPriceIncrements()")
const char*                                   dxfg_InstrumentProfile_getPriceIncrements(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the minimum allowed price increments with the corresponding price ranges of the instrument profile.
 *
 * The format is an increment, or a list of `<INCREMENT> <UPPER_LIMIT>` ranges separated by `; ` and sorted by upper
 * limit. Examples: `0.25`, `0.01 3; 0.05`.
 *
 * Java: `InstrumentProfile.setPriceIncrements(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setPriceIncrements().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The price increments, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setPriceIncrements()")
int32_t                                       dxfg_InstrumentProfile_setPriceIncrements(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the trading hours specification of the instrument profile.
 *
 * The format is the one of `Schedule.getInstance(String)` of the dxFeed API.
 *
 * Java: `InstrumentProfile.getTradingHours()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getTradingHours().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The trading hours (a UTF-8 string, empty if it is not set; release it with dxfg_String_release()), or NULL on
 * error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getTradingHours()")
const char*                                   dxfg_InstrumentProfile_getTradingHours(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Sets the trading hours specification of the instrument profile.
 *
 * The format is the one of `Schedule.getInstance(String)` of the dxFeed API.
 *
 * Java: `InstrumentProfile.setTradingHours(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setTradingHours().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] value The trading hours, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setTradingHours()")
int32_t                                       dxfg_InstrumentProfile_setTradingHours(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *value);
/**
 * @brief Returns the value of a field of the instrument profile by the name of the field, in the text form.
 *
 * The field is a standard field of the profile (TYPE, SYMBOL, DESCRIPTION, LOCAL_SYMBOL, LOCAL_DESCRIPTION, COUNTRY,
 * OPOL, EXCHANGE_DATA, EXCHANGES, CURRENCY, BASE_CURRENCY, CFI, ISIN, SEDOL, CUSIP, ICB, SIC, MULTIPLIER, PRODUCT,
 * UNDERLYING, SPC, ADDITIONAL_UNDERLYINGS, MMY, EXPIRATION, LAST_TRADE, STRIKE, OPTION_TYPE, EXPIRATION_STYLE,
 * SETTLEMENT_STYLE, PRICE_INCREMENTS, TRADING_HOURS) or, for any other name, a custom field. The numeric fields are
 * formatted like dxfg_InstrumentProfileField_formatNumber() does, the date fields like
 * dxfg_InstrumentProfileField_formatDate() does (a zero value is an empty string).
 *
 * Java: `InstrumentProfile.getField(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string; the standard field names are case sensitive.
 * @return The value (a UTF-8 string, empty if a custom field is not set; release it with dxfg_String_release()), or
 * NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getField()")
const char*                                   dxfg_InstrumentProfile_getField(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *name);
/**
 * @brief Sets the value of a field of the instrument profile by the name of the field, in the text form.
 *
 * The field is a standard field of the profile (TYPE, SYMBOL, DESCRIPTION, LOCAL_SYMBOL, LOCAL_DESCRIPTION, COUNTRY,
 * OPOL, EXCHANGE_DATA, EXCHANGES, CURRENCY, BASE_CURRENCY, CFI, ISIN, SEDOL, CUSIP, ICB, SIC, MULTIPLIER, PRODUCT,
 * UNDERLYING, SPC, ADDITIONAL_UNDERLYINGS, MMY, EXPIRATION, LAST_TRADE, STRIKE, OPTION_TYPE, EXPIRATION_STYLE,
 * SETTLEMENT_STYLE, PRICE_INCREMENTS, TRADING_HOURS) or, for any other name, a custom field. The text of a numeric or a
 * date field is parsed like dxfg_InstrumentProfileField_parseNumber() and dxfg_InstrumentProfileField_parseDate() do;
 * the function fails if the text has a wrong format.
 *
 * Java: `InstrumentProfile.setField(String, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string; the standard field names are case sensitive.
 * @param[in] value The value, a UTF-8 string; NULL is treated as an empty string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setField()")
int32_t                                       dxfg_InstrumentProfile_setField(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *name, const char *value);
/**
 * @brief Returns the value of a field of the instrument profile by the name of the field, in the numeric form.
 *
 * The standard numeric fields are ICB, SIC, MULTIPLIER, SPC, EXPIRATION, LAST_TRADE and STRIKE; the function fails for
 * the other standard fields. A custom field that is not set gives 0, a custom field whose text looks like a date
 * (`yyyy-MM-dd`) gives the day id, any other text is parsed as a number.
 *
 * Java: `InstrumentProfile.getNumericField(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getNumericField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @return The value, or negative infinity on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getNumericField()")
double                                        dxfg_InstrumentProfile_getNumericField(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *name);
/**
 * @brief Sets the value of a field of the instrument profile by the name of the field, in the numeric form.
 *
 * The function fails for the standard fields that are not numeric (see dxfg_InstrumentProfile_getNumericField()). A
 * custom field is stored as the text of the number.
 *
 * Java: `InstrumentProfile.setNumericField(String, double)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setNumericField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The value.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setNumericField()")
int32_t                                       dxfg_InstrumentProfile_setNumericField(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *name, double value);
/**
 * @brief Returns the day id of a date field of the instrument profile by the name of the field.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT. A custom field that is not set gives 0.
 *
 * Java: `InstrumentProfile.getDateField(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getDateField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @return The day id (the number of days since January 1, 1970, 00:00:00 GMT; 0 if it is not set), or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getDateField()")
int32_t                                       dxfg_InstrumentProfile_getDateField(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *name);
/**
 * @brief Sets the day id of a date field of the instrument profile by the name of the field.
 *
 * The day id is the number of days since January 1, 1970, 00:00:00 GMT. A custom field is stored as the text
 * `yyyy-MM-dd` (an empty string for 0).
 *
 * Java: `InstrumentProfile.setDateField(String, int)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_setDateField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The day id (the number of days since January 1, 1970, 00:00:00 GMT).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_setDateField()")
int32_t                                       dxfg_InstrumentProfile_setDateField(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip, const char *name, int32_t value);
/**
 * @brief Returns the names of the non-empty custom fields of the instrument profile.
 *
 * Java: `InstrumentProfile.addNonEmptyCustomFieldNames(Collection)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_getNonEmptyCustomFieldNames().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return The list of the names (release it with dxfg_CList_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_getNonEmptyCustomFieldNames()")
dxfg_string_list*                             dxfg_InstrumentProfile_getNonEmptyCustomFieldNames(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);
/**
 * @brief Releases an instrument profile handle.
 *
 * It is the same as dxfg_JavaObjectHandler_release() for the handle.
 *
 * @deprecated Use dxfg_v4_InstrumentProfile_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] ip The instrument profile handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfile_release()")
int32_t                                       dxfg_InstrumentProfile_release(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);

/**
 * @brief A handle of a Java InstrumentProfileCustomFields object; release it with dxfg_JavaObjectHandler_release().
 *
 * The custom fields are the fields of an instrument profile that are not standard fields, as the name-value pairs of
 * the text values. InstrumentProfileCustomFields is a class of the SDK, it does not belong to the dxFeed API. Use
 * dxfg_JavaObjectHandler_clone() to extend the life of the handle (to clone the object reference).
 */
typedef struct dxfg_instrument_profile_custom_fields_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_instrument_profile_custom_fields_t;

/**
 * @brief The fields of an instrument profile in a plain C structure, with the handle of its custom fields.
 *
 * The strings are UTF-8 null-terminated strings that belong to the structure (an empty string if the field is not set).
 * A structure that the SDK fills is released with dxfg_instrument_profile_free() (a single profile, see
 * dxfg_Iterable_InstrumentProfile_next2()), dxfg_instrument_profiles_array_free() (an array) or
 * dxfg_instrument_profile2_list_free() (a list): that frees the strings and releases the custom fields handle.
 *
 * A structure that is passed to the SDK, e.g. to dxfg_InstrumentProfileCollector_updateInstrumentProfile2(), is only
 * read: the SDK copies the strings, does not keep or free the pointers, treats a NULL string as an empty one and
 * ignores instrument_profile_custom_fields_hash.
 *
 * The custom fields are not copied: the Java instrument profile and the custom fields handle of the structure start
 * with the same storage of the fields, both in a structure that the SDK fills and in a structure passed to the SDK
 * (the profile gets the storage the handle has at the call). The sharing lasts only while the storage has room: a
 * change made with dxfg_InstrumentProfileCustomFields_setField() and the like through that handle may change the
 * profile (e.g. the profile in a collector), but a change that adds a field may move the fields of the handle to a
 * new, larger storage, and from then on the profile no longer sees the changes of the handle, nor the handle those of
 * the profile. So do not rely on the sharing in either direction: to change the custom fields of a profile, pass the
 * structure to the SDK again, and to change them independently, make a copy with
 * dxfg_InstrumentProfileCustomFields_new3() first.
 *
 * Every field is described with the corresponding dxfg_InstrumentProfile_get*() function.
 */
typedef struct dxfg_instrument_profile2_t {
  /** The type of the instrument, see dxfg_InstrumentProfile_getType(). */
  char* type;
  /** The identifier of the instrument, see dxfg_InstrumentProfile_getSymbol(). */
  char* symbol;
  /** The description of the instrument, see dxfg_InstrumentProfile_getDescription(). */
  char* description;
  /** The identifier of the instrument in the national language, see dxfg_InstrumentProfile_getLocalSymbol(). */
  char* local_symbol;
  /** The description in the national language, see dxfg_InstrumentProfile_getLocalDescription(). */
  char* local_description;
  /** The country of origin (ISO 3166-1 code), see dxfg_InstrumentProfile_getCountry(). */
  char* country;
  /** The official Place Of Listing (MIC), see dxfg_InstrumentProfile_getOPOL(). */
  char* opol;
  /** The exchange-specific data, see dxfg_InstrumentProfile_getExchangeData(). */
  char* exchange_data;
  /** The list of the exchanges where the instrument is traded, see dxfg_InstrumentProfile_getExchanges(). */
  char* exchanges;
  /** The currency of quotation, pricing and trading (ISO 4217 code), see dxfg_InstrumentProfile_getCurrency(). */
  char* currency;
  /** The base currency of a currency pair, see dxfg_InstrumentProfile_getBaseCurrency(). */
  char* base_currency;
  /** The Classification of Financial Instruments code, see dxfg_InstrumentProfile_getCFI(). */
  char* cfi;
  /** The International Securities Identifying Number, see dxfg_InstrumentProfile_getISIN(). */
  char* isin;
  /** The Stock Exchange Daily Official List code, see dxfg_InstrumentProfile_getSEDOL(). */
  char* sedol;
  /** The CUSIP code, see dxfg_InstrumentProfile_getCUSIP(). */
  char* cusip;
  /** The Industry Classification Benchmark number, see dxfg_InstrumentProfile_getICB(). */
  int32_t icb;
  /** The Standard Industrial Classification number, see dxfg_InstrumentProfile_getSIC(). */
  int32_t sic;
  /** The market value multiplier, see dxfg_InstrumentProfile_getMultiplier(). */
  double multiplier;
  /** The product for futures and options on futures, see dxfg_InstrumentProfile_getProduct(). */
  char* product;
  /** The primary underlying symbol for options, see dxfg_InstrumentProfile_getUnderlying(). */
  char* underlying;
  /** The shares per contract for options, see dxfg_InstrumentProfile_getSPC(). */
  double spc;
  /** The additional underlyings for options, see dxfg_InstrumentProfile_getAdditionalUnderlyings(). */
  char* additional_underlyings;
  /** The maturity month-year, see dxfg_InstrumentProfile_getMMY(). */
  char* mmy;
  /** The day id of the expiration (days since January 1, 1970, 00:00:00 GMT; 0 if not set). */
  int32_t expiration;
  /** The day id of the last trading day (days since January 1, 1970, 00:00:00 GMT; 0 if not set). */
  int32_t last_trade;
  /** The strike price for options, see dxfg_InstrumentProfile_getStrike(). */
  double strike;
  /** The type of the option, see dxfg_InstrumentProfile_getOptionType(). */
  char* option_type;
  /** The expiration cycle style, see dxfg_InstrumentProfile_getExpirationStyle(). */
  char* expiration_style;
  /** The settlement price determination style, see dxfg_InstrumentProfile_getSettlementStyle(). */
  char* settlement_style;
  /** The minimum allowed price increments, see dxfg_InstrumentProfile_getPriceIncrements(). */
  char* price_increments;
  /** The trading hours specification, see dxfg_InstrumentProfile_getTradingHours(). */
  char* trading_hours;
  /**
   * The handle of the custom fields (the fields that are not standard fields), or NULL if there are none. It
   * belongs to the structure and is released with it: clone it with dxfg_JavaObjectHandler_clone() to keep it.
   * If the structure is passed to the SDK, the custom fields of the handle are assigned to the profile.
   */
  dxfg_instrument_profile_custom_fields_t* instrument_profile_custom_fields;
  /**
   * The hash code of the custom fields, 0 if they are empty (the handle is then NULL). Equal custom fields give
   * equal hash codes.
   */
  int32_t instrument_profile_custom_fields_hash;
} dxfg_instrument_profile2_t;

/**
 * @brief A list of dxfg_instrument_profile2_t structures.
 *
 * Release the list and the structures in it with dxfg_instrument_profile2_list_free().
 */
typedef struct dxfg_instrument_profile2_list_t {
    /** The number of the elements. */
    int32_t size;
    /** The array of size pointers to the instrument profile structures; NULL if the list is empty. */
    dxfg_instrument_profile2_t **elements;
} dxfg_instrument_profile2_list_t;

/**
 * @brief Formats a number to the text form that is used in the instrument profile fields.
 *
 * The text has no digit grouping; 0 is formatted as an empty string.
 *
 * Java: `InstrumentProfileField.formatNumber(double)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileField_formatNumber().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] number The number to format.
 * @param[out] result The formatted number, a UTF-8 string; release it with dxfg_String_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileField_formatNumber()")
int32_t dxfg_InstrumentProfileField_formatNumber(graal_isolatethread_t *thread, double number, DXFG_OUT char** result);

/**
 * @brief Parses a number from the text form that is used in the instrument profile fields.
 *
 * A NULL or an empty string gives 0. The function fails if the text is not a number.
 *
 * Java: `InstrumentProfileField.parseNumber(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileField_parseNumber().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The text to parse, a UTF-8 string.
 * @param[out] number The parsed number.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileField_parseNumber()")
int32_t dxfg_InstrumentProfileField_parseNumber(graal_isolatethread_t *thread, const char* string, DXFG_OUT double* number);

/**
 * @brief Formats a day id to the text form that is used in the instrument profile fields.
 *
 * The text is `yyyy-MM-dd` for the GMT time zone; 0 is formatted as an empty string.
 *
 * Java: `InstrumentProfileField.formatDate(int)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileField_formatDate().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] date The day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @param[out] result The formatted date, a UTF-8 string; release it with dxfg_String_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileField_formatDate()")
int32_t dxfg_InstrumentProfileField_formatDate(graal_isolatethread_t *thread, int32_t date, DXFG_OUT char** result);

/**
 * @brief Parses a day id from the text form that is used in the instrument profile fields.
 *
 * The text is `yyyy-MM-dd` for the GMT time zone. A NULL or an empty string gives 0. The function fails if the text is
 * not a date of this format.
 *
 * Java: `InstrumentProfileField.parseDate(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileField_parseDate().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] string The text to parse, a UTF-8 string.
 * @param[out] date The parsed day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileField_parseDate()")
int32_t dxfg_InstrumentProfileField_parseDate(graal_isolatethread_t *thread, const char* string, DXFG_OUT int32_t* date);

/**
 * @brief Creates empty custom fields.
 *
 * Java: `new InstrumentProfileCustomFields()`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] custom_fields The handle of the new custom fields; release it with dxfg_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_new()")
int32_t dxfg_InstrumentProfileCustomFields_new(graal_isolatethread_t *thread, DXFG_OUT dxfg_instrument_profile_custom_fields_t **custom_fields);

/**
 * @brief Creates custom fields from an array of strings: the names and the values, consecutive.
 *
 * The array is name 0, value 0, name 1, value 1 and so on. If the array is NULL or the size is not positive, empty
 * custom fields are created.
 *
 * Java: `new InstrumentProfileCustomFields(String[])`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_newWithNamesAndValues().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields_array The array of the UTF-8 strings: names and values, consecutive.
 * @param[in] size The number of the strings in the array.
 * @param[out] custom_fields The handle of the new custom fields; release it with dxfg_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_newWithNamesAndValues()")
int32_t dxfg_InstrumentProfileCustomFields_new2(graal_isolatethread_t *thread, const char** custom_fields_array, int32_t size, DXFG_OUT dxfg_instrument_profile_custom_fields_t **custom_fields);

/**
 * @brief Creates custom fields as a copy of other custom fields.
 *
 * Java: `new InstrumentProfileCustomFields(InstrumentProfileCustomFields)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_newCopy().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] source The handle of the custom fields to copy.
 * @param[out] custom_fields The handle of the new custom fields; release it with dxfg_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_newCopy()")
int32_t dxfg_InstrumentProfileCustomFields_new3(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t * source, DXFG_OUT dxfg_instrument_profile_custom_fields_t **custom_fields);

/**
 * @brief Creates custom fields from a list of strings: the names and the values, consecutive.
 *
 * The list is name 0, value 0, name 1, value 1 and so on. If the list is NULL or empty, empty custom fields are
 * created.
 *
 * Java: `new InstrumentProfileCustomFields(String[])`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_newWithNamesAndValues().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields_string_list The list of the UTF-8 strings: names and values, consecutive.
 * @param[out] custom_fields The handle of the new custom fields; release it with dxfg_JavaObjectHandler_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_newWithNamesAndValues()")
int32_t dxfg_InstrumentProfileCustomFields_new4(graal_isolatethread_t *thread, const dxfg_string_list* custom_fields_string_list, DXFG_OUT dxfg_instrument_profile_custom_fields_t **custom_fields);


/**
 * @brief Returns the value of a custom field by the name.
 *
 * Java: `InstrumentProfileCustomFields.getField(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_getField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] result The value, a UTF-8 string; release it with dxfg_String_release(). It is NULL if there is no such
 * field.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_getField()")
int32_t dxfg_InstrumentProfileCustomFields_getField(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, const char *name, DXFG_OUT char** result);

/**
 * @brief Sets the value of a custom field by the name.
 *
 * An empty value clears the field: it is not listed by dxfg_InstrumentProfileCustomFields_getNonEmptyFieldNames()
 * anymore.
 *
 * Java: `InstrumentProfileCustomFields.setField(String, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_setField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The value, a UTF-8 string; it must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_setField()")
int32_t dxfg_InstrumentProfileCustomFields_setField(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, const char *name, const char* value);

/**
 * @brief Returns the numeric value of a custom field by the name.
 *
 * A field that is not set or is empty gives 0. A field whose text looks like a date (`yyyy-MM-dd`) gives the day id
 * (the number of days since January 1, 1970, 00:00:00 GMT); the text of any other field is parsed as a number.
 *
 * Java: `InstrumentProfileCustomFields.getNumericField(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_getNumericField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] result The numeric value of the field.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_getNumericField()")
int32_t dxfg_InstrumentProfileCustomFields_getNumericField(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, const char *name, DXFG_OUT double* result);

/**
 * @brief Sets the numeric value of a custom field by the name.
 *
 * The field is stored as the text of the number, see dxfg_InstrumentProfileField_formatNumber().
 *
 * Java: `InstrumentProfileCustomFields.setNumericField(String, double)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_setNumericField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The value.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_setNumericField()")
int32_t dxfg_InstrumentProfileCustomFields_setNumericField(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, const char *name, double value);

/**
 * @brief Returns the day id of a date custom field by the name.
 *
 * A field that is not set or is empty gives 0. The function fails if the text of the field is not a date, see
 * dxfg_InstrumentProfileField_parseDate().
 *
 * Java: `InstrumentProfileCustomFields.getDateField(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_getDateField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[out] result The day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_getDateField()")
int32_t dxfg_InstrumentProfileCustomFields_getDateField(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, const char *name, DXFG_OUT int32_t* result);

/**
 * @brief Sets the day id of a date custom field by the name.
 *
 * The field is stored as the text of the date, see dxfg_InstrumentProfileField_formatDate().
 *
 * Java: `InstrumentProfileCustomFields.setDateField(String, int)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_setDateField().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[in] name The name of the field, a UTF-8 string.
 * @param[in] value The day id: the number of days since January 1, 1970, 00:00:00 GMT.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_setDateField()")
int32_t dxfg_InstrumentProfileCustomFields_setDateField(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, const char *name, int32_t value);

/**
 * @brief Returns the names of the non-empty custom fields.
 *
 * @deprecated The updated parameter is useless. Use dxfg_v4_InstrumentProfileCustomFields_getNonEmptyFieldNames()
 * instead.
 *
 * Java: `InstrumentProfileCustomFields.addNonEmptyFieldNames(Collection)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[out] target_field_names The list of the names (empty if there are none); release it with
 * dxfg_CList_String_release().
 * @param[out] updated 1 if any name was written to the list, 0 otherwise. The pointer may be NULL if the flag is not
 * needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_getNonEmptyFieldNames(), the updated parameter is useless")
int32_t dxfg_InstrumentProfileCustomFields_addNonEmptyFieldNames(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, DXFG_OUT dxfg_string_list** target_field_names, DXFG_OUT int32_t* updated);

/**
 * @brief Returns the names of the non-empty custom fields.
 *
 * Java: `InstrumentProfileCustomFields.addNonEmptyFieldNames(Collection)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCustomFields_getNonEmptyFieldNames().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] custom_fields The custom fields handle.
 * @param[out] target_field_names The list of the names (empty if there are none); release it with
 * dxfg_CList_String_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCustomFields_getNonEmptyFieldNames()")
int32_t dxfg_InstrumentProfileCustomFields_getNonEmptyFieldNames(graal_isolatethread_t *thread, dxfg_instrument_profile_custom_fields_t *custom_fields, DXFG_OUT dxfg_string_list** target_field_names);

/**
 * @brief Updates one instrument profile in the collector from a structure.
 *
 * The collector gets a new Java instrument profile built from the structure: the structure is only read and can be
 * released after the call. The custom fields are shared, not copied (see dxfg_instrument_profile2_t): do not change
 * them through instrument_profile_custom_fields after the call.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfile(InstrumentProfile)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfile().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] instrument_profile The structure with the fields of the instrument profile.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfile()")
int32_t dxfg_InstrumentProfileCollector_updateInstrumentProfile2(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_instrument_profile2_t *instrument_profile);

/**
 * @brief Updates an array of instrument profiles in the collector from structures and assigns them a generation tag.
 *
 * The collector gets new Java instrument profiles built from the structures: the structures are only read and can be
 * released after the call; the custom fields are shared, not copied (see dxfg_instrument_profile2_t). The call does
 * nothing if the array is NULL or the size is not positive. The removal of an instrument profile is represented by a
 * profile whose type is `REMOVED`.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfiles(List, Object)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] instrument_profiles The array of the structures with the fields of the instrument profiles.
 * @param[in] size The number of the structures in the array.
 * @param[in] generation The generation tag, the handle of a Java object, or NULL for none. The tag is used by
 * dxfg_InstrumentProfileCollector_removeGenerations().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles()")
int32_t dxfg_InstrumentProfileCollector_updateInstrumentProfiles2(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, const dxfg_instrument_profile2_t* instrument_profiles, int32_t size, dxfg_java_object_handler *generation);

/**
 * @brief Updates a list of instrument profiles in the collector from structures and assigns them a generation tag.
 *
 * The collector gets new Java instrument profiles built from the structures: the structures are only read and can be
 * released after the call; the custom fields are shared, not copied (see dxfg_instrument_profile2_t). The removal of
 * an instrument profile is represented by a profile whose type is `REMOVED`.
 *
 * Java: `InstrumentProfileCollector.updateInstrumentProfiles(List, Object)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] collector The collector handle.
 * @param[in] instrument_profiles The list of the structures with the fields of the instrument profiles.
 * @param[in] generation The generation tag, the handle of a Java object, or NULL for none. The tag is used by
 * dxfg_InstrumentProfileCollector_removeGenerations().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles()")
int32_t dxfg_InstrumentProfileCollector_updateInstrumentProfiles3(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, const dxfg_instrument_profile2_list_t* instrument_profiles, dxfg_java_object_handler *generation);

/**
 * @brief Returns the next instrument profile of the iterator as a structure.
 *
 * Check dxfg_Iterable_InstrumentProfile_hasNext() first: the function fails if there is no next profile. The strings
 * of the structure are copies; its custom fields handle shares the custom fields of the profile in the collector (see
 * dxfg_instrument_profile2_t).
 *
 * Java: `Iterator.next()`.
 *
 * @deprecated Use dxfg_v4_Iterable_InstrumentProfile_next().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] iterable_ip The iterator handle.
 * @param[out] instrument_profile The next instrument profile; release it with dxfg_instrument_profile_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Iterable_InstrumentProfile_next()")
int32_t dxfg_Iterable_InstrumentProfile_next2(graal_isolatethread_t *thread, dxfg_iterable_ip_t *iterable_ip, DXFG_OUT dxfg_instrument_profile2_t **instrument_profile);

/**
 * @brief Reads and returns the instrument profiles from an address as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The operation updates the values of
 * dxfg_InstrumentProfileReader_getLastModified() and dxfg_InstrumentProfileReader_wasComplete(). The function fails if
 * the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFile().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[out] instrument_profiles The array of the structures with the instrument profiles, NULL if there are none;
 * release it with dxfg_instrument_profiles_array_free().
 * @param[out] size The number of the structures in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFile()")
int32_t dxfg_InstrumentProfileReader_readFromFile4(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, DXFG_OUT dxfg_instrument_profile2_t **instrument_profiles, DXFG_OUT int32_t* size);

/**
 * @brief Reads and returns the instrument profiles from an address with basic authentication credentials as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The user and the password take precedence over the
 * authentication information of the URL. The operation updates the values of
 * dxfg_InstrumentProfileReader_getLastModified() and dxfg_InstrumentProfileReader_wasComplete(). The function fails if
 * the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, String, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPassword().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] user The user name of the basic authentication, a UTF-8 string, or NULL.
 * @param[in] password The password of the basic authentication, a UTF-8 string, or NULL.
 * @param[out] instrument_profiles The array of the structures with the instrument profiles, NULL if there are none;
 * release it with dxfg_instrument_profiles_array_free().
 * @param[out] size The number of the structures in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPassword()")
int32_t dxfg_InstrumentProfileReader_readFromFile5(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, const char *user, const char *password, DXFG_OUT dxfg_instrument_profile2_t **instrument_profiles, DXFG_OUT int32_t* size);

/**
 * @brief Reads and returns the instrument profiles from an address with token credentials as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The token takes precedence over the authentication
 * information of the URL. The operation updates the values of dxfg_InstrumentProfileReader_getLastModified() and
 * dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument Profile
 * Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, AuthToken)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFileWithToken().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] token The token handle, or NULL.
 * @param[out] instrument_profiles The array of the structures with the instrument profiles, NULL if there are none;
 * release it with dxfg_instrument_profiles_array_free().
 * @param[out] size The number of the structures in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFileWithToken()")
int32_t dxfg_InstrumentProfileReader_readFromFile6(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, dxfg_auth_token_t *token, DXFG_OUT dxfg_instrument_profile2_t **instrument_profiles, DXFG_OUT int32_t* size);

/**
 * @brief Reads and returns the instrument profiles from an address as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The operation updates the values of
 * dxfg_InstrumentProfileReader_getLastModified() and dxfg_InstrumentProfileReader_wasComplete(). The function fails if
 * the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFile().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_instrument_profile2_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFile()")
int32_t dxfg_InstrumentProfileReader_readFromFile7(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, DXFG_OUT dxfg_instrument_profile2_list_t **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an address with basic authentication credentials as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The user and the password take precedence over the
 * authentication information of the URL. The operation updates the values of
 * dxfg_InstrumentProfileReader_getLastModified() and dxfg_InstrumentProfileReader_wasComplete(). The function fails if
 * the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, String, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPassword().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] user The user name of the basic authentication, a UTF-8 string, or NULL.
 * @param[in] password The password of the basic authentication, a UTF-8 string, or NULL.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_instrument_profile2_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPassword()")
int32_t dxfg_InstrumentProfileReader_readFromFile8(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, const char *user, const char *password, DXFG_OUT dxfg_instrument_profile2_list_t **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an address with token credentials as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The address is a URL, a file path or `<host>:<port>` (see
 * dxfg_InstrumentProfileReader_resolveSourceURL()). Authentication information can be passed as a part of the URL, for
 * example `http://user:password@host:port/path/file.ipf`. The token takes precedence over the authentication
 * information of the URL. The operation updates the values of dxfg_InstrumentProfileReader_getLastModified() and
 * dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument Profile
 * Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readFromFile(String, AuthToken)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readFromFileWithToken().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] address The address, a UTF-8 string.
 * @param[in] token The token handle, or NULL.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_instrument_profile2_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readFromFileWithToken()")
int32_t dxfg_InstrumentProfileReader_readFromFile9(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, dxfg_auth_token_t *token, DXFG_OUT dxfg_instrument_profile2_list_t **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an input stream with plain (not compressed) data as structures.
 *
 * It is deprecated in Java. Use dxfg_InstrumentProfileReader_read4() or dxfg_InstrumentProfileReader_read6(). Unlike
 * dxfg_InstrumentProfileReader_read4(), the function does not decompress the stream and does not report an incomplete
 * one. The function fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_read().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrument_profiles The array of the structures with the instrument profiles, NULL if there are none;
 * release it with dxfg_instrument_profiles_array_free().
 * @param[out] size The number of the structures in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_read()")
int32_t dxfg_InstrumentProfileReader_read3(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, DXFG_OUT dxfg_instrument_profile2_t **instrument_profiles, DXFG_OUT int32_t* size);

/**
 * @brief Reads and returns the instrument profiles from an input stream as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The stream is not closed. The operation updates the value of
 * dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument Profile
 * Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readWithAddress().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[in] address The origin of the stream for debugging purposes, a UTF-8 string.
 * @param[out] instrument_profiles The array of the structures with the instrument profiles, NULL if there are none;
 * release it with dxfg_instrument_profiles_array_free().
 * @param[out] size The number of the structures in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readWithAddress()")
int32_t dxfg_InstrumentProfileReader_read4(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, const char *address, DXFG_OUT dxfg_instrument_profile2_t **instrument_profiles, DXFG_OUT int32_t* size);

/**
 * @brief Reads and returns the instrument profiles from an input stream with plain (not compressed) data as structures.
 *
 * It is deprecated in Java. Use dxfg_InstrumentProfileReader_read4() or dxfg_InstrumentProfileReader_read6(). Unlike
 * dxfg_InstrumentProfileReader_read4(), the function does not decompress the stream and does not report an incomplete
 * one. The function fails if the data does not conform to the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_read().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_instrument_profile2_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_read()")
int32_t dxfg_InstrumentProfileReader_read5(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, DXFG_OUT dxfg_instrument_profile2_list_t **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from an input stream as structures.
 *
 * The data compression formats "zip" and "gzip" are recognized automatically: for zip the first file entry is read as a
 * plain data stream, other data is parsed as is. The stream is not closed. The operation updates the value of
 * dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to the Instrument Profile
 * Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.read(InputStream, String)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readWithAddress().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[in] address The origin of the stream for debugging purposes, a UTF-8 string.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_instrument_profile2_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readWithAddress()")
int32_t dxfg_InstrumentProfileReader_read6(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, const char *address, DXFG_OUT dxfg_instrument_profile2_list_t **instrument_profiles);

/**
 * @brief Reads and returns the instrument profiles from a compressed input stream as structures.
 *
 * It is deprecated in Java. The data compression formats "zip" and "gzip" are recognized automatically: for zip the
 * first file entry is read as a plain data stream, other data is parsed as is. The stream is not closed. The operation
 * updates the value of dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to
 * the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readCompressed(InputStream)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readCompressed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrument_profiles The array of the structures with the instrument profiles, NULL if there are none;
 * release it with dxfg_instrument_profiles_array_free().
 * @param[out] size The number of the structures in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readCompressed()")
int32_t dxfg_InstrumentProfileReader_readCompressed2(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, DXFG_OUT dxfg_instrument_profile2_t **instrument_profiles, DXFG_OUT int32_t* size);

/**
 * @brief Reads and returns the instrument profiles from a compressed input stream as structures.
 *
 * It is deprecated in Java. The data compression formats "zip" and "gzip" are recognized automatically: for zip the
 * first file entry is read as a plain data stream, other data is parsed as is. The stream is not closed. The operation
 * updates the value of dxfg_InstrumentProfileReader_wasComplete(). The function fails if the data does not conform to
 * the Instrument Profile Format or on an I/O error.
 *
 * Java: `InstrumentProfileReader.readCompressed(InputStream)`.
 *
 * @deprecated Use dxfg_v4_InstrumentProfileReader_readCompressed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] reader The reader handle.
 * @param[in] is The input stream handle.
 * @param[out] instrument_profiles The list of the structures with the instrument profiles; release it with
 * dxfg_instrument_profile2_list_free().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_InstrumentProfileReader_readCompressed()")
int32_t dxfg_InstrumentProfileReader_readCompressed3(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, DXFG_OUT dxfg_instrument_profile2_list_t **instrument_profiles);

/**
 * @brief Frees the structure of an instrument profile that the SDK has filled, and releases its custom fields handle.
 *
 * The structure itself is freed as well: use it for a structure that was returned by
 * dxfg_Iterable_InstrumentProfile_next2(), not for an element of an array or a list.
 *
 * @deprecated Use dxfg_v4_instrument_profile_data_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrument_profile The structure of the instrument profile; NULL does nothing.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_instrument_profile_data_release()")
int32_t dxfg_instrument_profile_free(graal_isolatethread_t *thread, dxfg_instrument_profile2_t *instrument_profile);

/**
 * @brief Frees an array of the structures of instrument profiles that the SDK has filled, with the strings of every
 * structure, and releases their custom fields handles.
 *
 * @deprecated Use dxfg_v4_instrument_profile_data_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrument_profiles The array of the structures; NULL does nothing.
 * @param[in] size The number of the structures in the array.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_instrument_profile_data_list_release()")
int32_t dxfg_instrument_profiles_array_free(graal_isolatethread_t *thread, dxfg_instrument_profile2_t *instrument_profiles, int32_t size);

/**
 * @brief Frees a list of the structures of instrument profiles that the SDK has filled, with the structures in it, and
 * releases their custom fields handles.
 *
 * @deprecated Use dxfg_v4_instrument_profile_data_list_release().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] instrument_profiles The list of the structures.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_instrument_profile_data_list_release()")
int32_t dxfg_instrument_profile2_list_free(graal_isolatethread_t *thread, dxfg_instrument_profile2_list_t* instrument_profiles);

/** @} */ // end of InstrumentProfile

#ifdef __cplusplus
}
#endif

#endif // DXFG_IPF_H
