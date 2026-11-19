// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK Endpoint functions and types declarations
 */

#ifndef DXFG_ENDPOINT_H
#define DXFG_ENDPOINT_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"
#include "dxfg_javac.h"

/** @defgroup Endpoint Endpoint
 *  @brief The endpoint (Java `DXEndpoint`): a connection to a data feed or to a publisher, its role and state, and
 *  the access to its feed and publisher.
 *
 *  The ownership rules of the handles, the strings and the lists are described in the group @ref Javac.
 *
 *  A typical life cycle of an endpoint:
 *  1. Create it with dxfg_DXEndpoint_Builder_build() (see the group @ref Builder), dxfg_DXEndpoint_create() or
 *     dxfg_DXEndpoint_create2(). The default endpoint of the process is returned by dxfg_DXEndpoint_getInstance().
 *  2. Optionally add a state change listener (dxfg_PropertyChangeListener_new(),
 *     dxfg_DXEndpoint_addStateChangeListener()) and set the credentials (dxfg_DXEndpoint_user(),
 *     dxfg_DXEndpoint_password()).
 *  3. Connect with dxfg_DXEndpoint_connect() and take the feed (dxfg_DXEndpoint_getFeed()) or the publisher
 *     (dxfg_DXEndpoint_getPublisher()).
 *  4. Close it with dxfg_DXEndpoint_close() or dxfg_DXEndpoint_closeAndAwaitTermination() and release all the
 *     handles with dxfg_JavaObjectHandler_release(). Releasing a handle does not close the endpoint. An endpoint that
 *     stays open keeps its threads running: close all of them with dxfg_system_close_all_and_await_termination()
 *     before graal_tear_down_isolate().
 *  @{
 */

/**
 * @brief Forward declaration of the handle of a Java DXFeed object (see dxfg_feed.h): the result of
 * dxfg_DXEndpoint_getFeed().
 */
typedef struct dxfg_feed_t dxfg_feed_t;

/** @brief Forward declaration of the handle of a Java DXPublisher object (see dxfg_publisher.h): the result of
 * dxfg_DXEndpoint_getPublisher(). */
typedef struct dxfg_publisher_t dxfg_publisher_t;

/** @brief Forward declaration of the list of event classes (see dxfg_events.h): the result of
 * dxfg_DXEndpoint_getEventTypes(). */
typedef struct dxfg_event_clazz_list_t dxfg_event_clazz_list_t;

/**
 * @brief A handle of a Java DXEndpoint object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_DXEndpoint_Builder_build(), dxfg_DXEndpoint_create(), dxfg_DXEndpoint_create2() or get the
 * default one with dxfg_DXEndpoint_getInstance(). Releasing the handle does not close the endpoint: call
 * dxfg_DXEndpoint_close() or dxfg_DXEndpoint_closeAndAwaitTermination().
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.html">Javadoc</a>
 */
typedef struct dxfg_endpoint_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_endpoint_t;

/**
 * @brief A handle of a Java DXEndpoint.Builder object; release it with dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_DXEndpoint_newBuilder(). The endpoints that the builder has built are independent of it:
 * releasing the builder's handle does not affect them.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Builder.html">Javadoc</a>
 */
typedef struct dxfg_endpoint_builder_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_endpoint_builder_t;

/**
 * @brief A handle of a Java PropertyChangeListener object that reports the state changes of an endpoint; release it
 * with dxfg_JavaObjectHandler_release().
 *
 * The listener calls a dxfg_endpoint_state_change_listener_func when the state of an endpoint to which it is added
 * changes. Create it with dxfg_PropertyChangeListener_new() and add it with dxfg_DXEndpoint_addStateChangeListener().
 * The last handle of the listener that is released deactivates the listener: its function is not called anymore.
 *
 * <a href="https://docs.oracle.com/javase/8/docs/api/java/beans/PropertyChangeListener.html">Javadoc</a>
 */
typedef struct dxfg_endpoint_state_change_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_endpoint_state_change_listener_t;

/**
 * @brief The role of an endpoint: it is specified when the endpoint is created and cannot be changed afterwards.
 *
 * The role defines what dxfg_DXEndpoint_connect() connects to and which of dxfg_DXEndpoint_getFeed() and
 * dxfg_DXEndpoint_getPublisher() is meaningful. The default role is #DXFG_ENDPOINT_ROLE_FEED.
 *
 * Java: `DXEndpoint.Role`.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Role.html">Javadoc</a>
 */
typedef enum dxfg_endpoint_role_t {
    /** Connects to a remote data feed provider; optimized for real-time or delayed data. The default role. The
     * events are conflated when the event processing cannot keep up. */
    DXFG_ENDPOINT_ROLE_FEED = 0,
    /** Similar to #DXFG_ENDPOINT_ROLE_FEED, but only for the replay of historical data with the on-demand service:
     * it is not connected to the data provider until the replay is started. It cannot be connected to an
     * ordinary data feed. */
    DXFG_ENDPOINT_ROLE_ON_DEMAND_FEED,
    /** Similar to #DXFG_ENDPOINT_ROLE_FEED, but designed for bulk parsing of data from files: the events are not
     * conflated and not skipped and are processed as fast as possible. The last event requests do not work in this
     * role. */
    DXFG_ENDPOINT_ROLE_STREAM_FEED,
    /** Connects to a remote publisher hub or creates a publisher on the local host; the publisher of the endpoint
     * publishes the events to all the connected feeds. The last event requests and the time series subscription do
     * not work in this role. */
    DXFG_ENDPOINT_ROLE_PUBLISHER,
    /** Similar to #DXFG_ENDPOINT_ROLE_PUBLISHER, but designed for bulk publishing: the events are not conflated and
     * not skipped and are processed as fast as possible. */
    DXFG_ENDPOINT_ROLE_STREAM_PUBLISHER,
    /** A local hub without the network connections: the events published with the publisher of the endpoint are
     * delivered to the feed of the same endpoint only. */
    DXFG_ENDPOINT_ROLE_LOCAL_HUB,
} dxfg_endpoint_role_t;

/**
 * @brief The state of an endpoint, see dxfg_DXEndpoint_getState().
 *
 * An endpoint is created #DXFG_ENDPOINT_STATE_NOT_CONNECTED, unless it connects on creation (see
 * dxfg_DXEndpoint_Builder_build()). dxfg_DXEndpoint_connect() makes it #DXFG_ENDPOINT_STATE_CONNECTING, the established
 * connection makes it #DXFG_ENDPOINT_STATE_CONNECTED, a lost network connection returns it to
 * #DXFG_ENDPOINT_STATE_CONNECTING. #DXFG_ENDPOINT_STATE_CLOSED is the final state.
 *
 * Java: `DXEndpoint.State`.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.State.html">Javadoc</a>
 */
typedef enum dxfg_endpoint_state_t {
    /** The endpoint is not connected to a remote endpoint: it is just created, or was disconnected, or a file that
     * it read is finished. */
    DXFG_ENDPOINT_STATE_NOT_CONNECTED = 0,
    /** dxfg_DXEndpoint_connect() was called, but the connection is not established yet or was lost. */
    DXFG_ENDPOINT_STATE_CONNECTING,
    /** The connection to a remote endpoint is established. */
    DXFG_ENDPOINT_STATE_CONNECTED,
    /** The endpoint was closed by dxfg_DXEndpoint_close(); no connection can be made anymore. */
    DXFG_ENDPOINT_STATE_CLOSED,
} dxfg_endpoint_state_t;

/**
 * @brief The function of a dxfg_endpoint_state_change_listener_t: called when the state of the endpoint changes.
 *
 * The function is called by a thread of the executor of the endpoint (the default thread pool or the executor that
 * was set with dxfg_DXEndpoint_Builder_withExecutor()), not by the thread that added the listener or changed the
 * state. A change is reported from the last reported state to the current one, so the fast changes are coalesced
 * (see dxfg_DXEndpoint_addStateChangeListener()). The states are passed by value, there is nothing to release.
 * Do not call dxfg_DXEndpoint_closeAndAwaitTermination() of the same endpoint from the function.
 *
 * Java: `PropertyChangeListener.propertyChange(PropertyChangeEvent)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] old_state The previous state of the endpoint.
 * @param[in] new_state The new state of the endpoint.
 * @param[in] user_data The pointer that was passed to dxfg_PropertyChangeListener_new(); the SDK does not use it.
 */
typedef void (*dxfg_endpoint_state_change_listener_func)(graal_isolatethread_t *thread, dxfg_endpoint_state_t old_state,
                                                         dxfg_endpoint_state_t new_state, void *user_data);

/** @defgroup Builder Endpoint builder
 *  @brief The builder of endpoints (Java `DXEndpoint.Builder`).
 *
 * The recommended way to create an endpoint: the builder allows to set the role, the name, the executor and the
 * properties
 * (see <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Builder.html">DXEndpoint.Builder</a>).
 * Every dxfg_DXEndpoint_Builder_build() call creates a new endpoint.
 *
 *  @{
 */

/**
 * @brief Creates a new endpoint builder.
 *
 * The builder has the #DXFG_ENDPOINT_ROLE_FEED role and the default executor until they are changed.
 *
 * Java: `DXEndpoint.newBuilder()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_newBuilder().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The builder (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_newBuilder()")
dxfg_endpoint_builder_t*    dxfg_DXEndpoint_newBuilder(graal_isolatethread_t *thread);

/**
 * @brief Sets the role of the endpoints built by the builder.
 *
 * The default role is #DXFG_ENDPOINT_ROLE_FEED.
 *
 * Java: `DXEndpoint.Builder.withRole(DXEndpoint.Role)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_withRole().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] role The role of the endpoints.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_withRole()")
int32_t                     dxfg_DXEndpoint_Builder_withRole(graal_isolatethread_t *thread, dxfg_endpoint_builder_t *builder, dxfg_endpoint_role_t role);

/**
 * @brief Sets the name of the endpoints built by the builder.
 *
 * The name distinguishes several endpoints in the same process in the logs and other diagnostics. It is the same as
 * dxfg_DXEndpoint_Builder_withProperty() with the key `name`. A short name that describes the role of the endpoint in
 * the application is recommended.
 *
 * Java: `DXEndpoint.Builder.withName(String)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_withName().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] name The name, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_withName()")
int32_t                     dxfg_DXEndpoint_Builder_withName(graal_isolatethread_t *thread, dxfg_endpoint_builder_t *builder, const char *name);

/**
 * @brief Sets a property of the endpoints built by the builder.
 *
 * A property that the builder does not support is ignored: check it with dxfg_DXEndpoint_Builder_supportsProperty().
 * A property that is set here takes precedence over the same property from the properties file and from the system
 * properties (see dxfg_system_set_property()).
 *
 * Java: `DXEndpoint.Builder.withProperty(String, String)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_withProperty().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] key The name of the property, a NUL-terminated UTF-8 string; the SDK copies it.
 * @param[in] value The value of the property, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_withProperty()")
int32_t                     dxfg_DXEndpoint_Builder_withProperty(graal_isolatethread_t *thread, dxfg_endpoint_builder_t *builder, const char *key, const char *value);

/**
 * @brief Sets the properties of the endpoints built by the builder from a properties file.
 *
 * The file is read in the format of a Java properties file (`java.util.Properties.load`) and every property that the
 * builder supports is set as by dxfg_DXEndpoint_Builder_withProperty(). The unsupported properties are ignored.
 * An error is reported if the file cannot be read.
 *
 * Java: `DXEndpoint.Builder.withProperties(Properties)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_withProperties().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] file_path The path of the properties file, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_withProperties()")
int32_t                     dxfg_DXEndpoint_Builder_withProperties(graal_isolatethread_t *thread, dxfg_endpoint_builder_t *builder, const char *file_path);

/**
 * @brief Sets the executor of the notifications of the endpoints built by this builder: event processing, state
 * changes, etc.
 *
 * Use it instead of the deprecated dxfg_DXEndpoint_executor().
 *
 * By default, the endpoint uses a thread pool with the size equal to the number of available processors
 * (the `dxfeed.threadPoolSize` property for the feed roles and `dxpublisher.threadPoolSize` for the publisher).
 * The executor must not reject the submitted tasks.
 *
 * Java: `DXEndpoint.Builder.withExecutor(Executor)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_withExecutor().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] executor The executor (e.g. from dxfg_Executors_newFixedThreadPool()) or NULL to use the default one.
 * The builder does not take the handle: free it with dxfg_JavaObjectHandler_release() when it is no longer needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_withExecutor()")
int32_t                     dxfg_DXEndpoint_Builder_withExecutor(graal_isolatethread_t *thread, dxfg_endpoint_builder_t *builder, dxfg_executor_t *executor);

/**
 * @brief Checks whether the builder supports a property.
 *
 * The properties that are not supported are ignored by dxfg_DXEndpoint_Builder_withProperty().
 *
 * Java: `DXEndpoint.Builder.supportsProperty(String)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_supportsProperty().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] key The name of the property, a NUL-terminated UTF-8 string.
 * @return 1 if the property is supported, 0 if it is not supported, or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_supportsProperty()")
int32_t                     dxfg_DXEndpoint_Builder_supportsProperty(graal_isolatethread_t *thread, dxfg_endpoint_builder_t *builder, const char *key);

/**
 * @brief Builds a new endpoint with the builder's settings.
 *
 * Every call builds a new endpoint. An endpoint with the role #DXFG_ENDPOINT_ROLE_FEED or #DXFG_ENDPOINT_ROLE_PUBLISHER
 * connects automatically on creation if `dxfeed.address` (feed) or `dxpublisher.address` (publisher) is set; the
 * endpoints of the other roles do not.
 *
 * Java: `DXEndpoint.Builder.build()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_build().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @return The new endpoint, or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 * The caller owns the endpoint: close it with dxfg_DXEndpoint_close() or dxfg_DXEndpoint_closeAndAwaitTermination()
 * and free its handle with dxfg_JavaObjectHandler_release() (releasing the handle does not close the endpoint).
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_build()")
dxfg_endpoint_t*            dxfg_DXEndpoint_Builder_build(graal_isolatethread_t *thread, dxfg_endpoint_builder_t *builder);

/** @} */ // end of Builder

/**
 * @brief Returns the default application-wide singleton endpoint with the #DXFG_ENDPOINT_ROLE_FEED role, the same one
 * that dxfg_DXFeed_getInstance() uses.
 *
 * The endpoint is created on the first call and is configured by the properties of the endpoint and the system
 * properties (see dxfg_system_set_property()) as of that moment.
 *
 * All callers share this endpoint: closing it closes it for everyone in the process.
 *
 * Java: `DXEndpoint.getInstance()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The endpoint (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getInstance()")
dxfg_endpoint_t*                dxfg_DXEndpoint_getInstance(graal_isolatethread_t *thread);

/**
 * @brief Creates a <b>new</b> endpoint with the given role on every call, the same as dxfg_DXEndpoint_create2().
 *
 * Unlike the Java `DXEndpoint.getInstance(Role)`, this function does not return a singleton: the SDK calls
 * `DXEndpoint.create(Role)`. Prefer dxfg_DXEndpoint_newBuilder() (or dxfg_DXEndpoint_create2()) in the new code.
 *
 * Java: `DXEndpoint.create(DXEndpoint.Role)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_createWithRole().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] role The role of the endpoint.
 * @return The new endpoint, or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 * The caller owns the endpoint: close it with dxfg_DXEndpoint_close() or dxfg_DXEndpoint_closeAndAwaitTermination()
 * and free its handle with dxfg_JavaObjectHandler_release().
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_createWithRole()")
dxfg_endpoint_t*                dxfg_DXEndpoint_getInstance2(graal_isolatethread_t *thread, dxfg_endpoint_role_t role);

/**
 * @brief Creates a new endpoint with the #DXFG_ENDPOINT_ROLE_FEED role.
 *
 * Java: `DXEndpoint.create()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_create().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The new endpoint, or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 * The caller owns the endpoint: close it with dxfg_DXEndpoint_close() or dxfg_DXEndpoint_closeAndAwaitTermination()
 * and free its handle with dxfg_JavaObjectHandler_release().
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_create()")
dxfg_endpoint_t*                dxfg_DXEndpoint_create(graal_isolatethread_t *thread);

/**
 * @brief Creates a new endpoint with the given role.
 *
 * Java: `DXEndpoint.create(DXEndpoint.Role)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_createWithRole().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] role The role of the endpoint.
 * @return The new endpoint, or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 * The caller owns the endpoint: close it with dxfg_DXEndpoint_close() or dxfg_DXEndpoint_closeAndAwaitTermination()
 * and free its handle with dxfg_JavaObjectHandler_release().
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_createWithRole()")
dxfg_endpoint_t*                dxfg_DXEndpoint_create2(graal_isolatethread_t *thread, dxfg_endpoint_role_t role);

/**
 * @brief Closes the endpoint: disconnects it and releases its resources. The endpoint cannot be used after that.
 *
 * The state of the endpoint becomes #DXFG_ENDPOINT_STATE_CLOSED immediately.
 *
 * The function does not wait for the notifications that are already scheduled: a state change listener
 * (including the change to #DXFG_ENDPOINT_STATE_CLOSED) or an event listener may still be called after it returns.
 * Use dxfg_DXEndpoint_closeAndAwaitTermination() to wait for them, or dxfg_JavaObjectHandler_release() of the listener
 * to stop its calls (see dxfg_JavaObjectHandler_release()).
 *
 * Java: `DXEndpoint.close()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_close().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_close()")
int32_t                         dxfg_DXEndpoint_close(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Closes the endpoint and waits until all the data processing and the scheduled notifications are completed.
 *
 * It does the same as dxfg_DXEndpoint_close() and blocks until the outstanding tasks of the endpoint are finished. It
 * is designed for the #DXFG_ENDPOINT_ROLE_STREAM_FEED role, after dxfg_DXEndpoint_awaitNotConnected() returns, to make
 * sure that a file was completely processed.
 *
 * Do not call it from a listener of this endpoint: it would wait for the notification that is calling it.
 *
 * Java: `DXEndpoint.closeAndAwaitTermination()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_closeAndAwaitTermination().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also if
 * the waiting thread is interrupted). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_closeAndAwaitTermination()")
int32_t                         dxfg_DXEndpoint_closeAndAwaitTermination(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Returns the role of the endpoint.
 *
 * Java: `DXEndpoint.getRole()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getRole().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return The role of the endpoint, or -1 on error: it is not a value of #dxfg_endpoint_role_t, so compare the result
 * with -1 before using it. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getRole()")
dxfg_endpoint_role_t            dxfg_DXEndpoint_getRole(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Sets the user name of the endpoint.
 *
 * Call it before dxfg_DXEndpoint_connect(), together with dxfg_DXEndpoint_password(), to configure the credentials of
 * a service that requires them. A NULL user name is an error.
 *
 * Java: `DXEndpoint.user(String)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_user().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] user The user name, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_user()")
int32_t                         dxfg_DXEndpoint_user(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint, const char *user);

/**
 * @brief Sets the password of the endpoint.
 *
 * Call it before dxfg_DXEndpoint_connect(), together with dxfg_DXEndpoint_user(), to configure the credentials of a
 * service that requires them. A NULL password is an error.
 *
 * Java: `DXEndpoint.password(String)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_password().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] password The password, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_password()")
int32_t                         dxfg_DXEndpoint_password(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint, const char *password);

/**
 * @brief Connects the endpoint to the remote address.
 *
 * The previously established connections are closed if the new address is different from the old one. The function
 * does nothing if the address does not change or if the endpoint is closed. Otherwise, the state of the endpoint
 * immediately becomes #DXFG_ENDPOINT_STATE_CONNECTING.
 *
 * The function does not wait until the connection is established: it is established asynchronously. It waits only
 * until the state change listeners have processed the change to #DXFG_ENDPOINT_STATE_CONNECTING.
 *
 * The address string comes with the market data vendor agreement, for example `host:port` to connect to a server
 * with TCP/IP, or `:port` to listen for a connection. A malformed address or a NULL address is an error. For the
 * services that require credentials call dxfg_DXEndpoint_user() and dxfg_DXEndpoint_password() first.
 *
 * Java: `DXEndpoint.connect(String)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_connect().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] address The address of the data source, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_connect()")
int32_t                         dxfg_DXEndpoint_connect(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint, const char *address);

/**
 * @brief Terminates all the network connections of the endpoint and initiates connecting again with the same address.
 *
 * It works like dxfg_DXEndpoint_disconnect() followed by dxfg_DXEndpoint_connect() with the current address, but the
 * internal resources may be reused. A connection with several target addresses tries to switch to an alternative
 * address; the configured reconnect timeouts apply.
 *
 * The function does not connect an endpoint that was not connected with dxfg_DXEndpoint_connect() or was
 * disconnected with dxfg_DXEndpoint_disconnect(). Whether a state change listener sees the intermediate
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED state is not specified.
 *
 * Java: `DXEndpoint.reconnect()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_reconnect().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_reconnect()")
int32_t                         dxfg_DXEndpoint_reconnect(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Terminates all the remote network connections of the endpoint.
 *
 * The function does nothing if the endpoint is closed. Otherwise, the state of the endpoint immediately becomes
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED. The resources of the endpoint are not released: use dxfg_DXEndpoint_close() for
 * that.
 *
 * Java: `DXEndpoint.disconnect()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_disconnect().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_disconnect()")
int32_t                         dxfg_DXEndpoint_disconnect(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Terminates all the remote network connections of the endpoint and clears the stored data.
 *
 * The function does nothing if the endpoint is closed. Otherwise, the state of the endpoint immediately becomes
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED. The resources of the endpoint are not released: use dxfg_DXEndpoint_close() for
 * that.
 *
 * Java: `DXEndpoint.disconnectAndClear()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_disconnectAndClear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_disconnectAndClear()")
int32_t                         dxfg_DXEndpoint_disconnectAndClear(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Waits until the endpoint stops processing data (becomes quiescent).
 *
 * It is important when data is written to a file with a `tape:...` address: wait for it to make sure that all the
 * published data was written before closing the endpoint.
 *
 * The function blocks the calling thread. With a single-threaded executor of the endpoint, do not call it from the
 * executor's thread (e.g. from a listener): it would wait forever.
 *
 * Java: `DXEndpoint.awaitProcessed()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_awaitProcessed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also if
 * the waiting thread is interrupted). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_awaitProcessed()")
int32_t                         dxfg_DXEndpoint_awaitProcessed(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Waits until the state of the endpoint becomes #DXFG_ENDPOINT_STATE_NOT_CONNECTED or
 * #DXFG_ENDPOINT_STATE_CLOSED.
 *
 * It is the signal that the files opened with dxfg_DXEndpoint_connect() (`file:...` addresses) were finished reading,
 * but not necessarily completely processed by the subscription listeners. Call
 * dxfg_DXEndpoint_closeAndAwaitTermination() after it returns to make sure that all the processing has completed.
 *
 * The function blocks the calling thread. With a single-threaded executor of the endpoint, do not call it from the
 * executor's thread (e.g. from a listener): it would wait forever.
 *
 * Java: `DXEndpoint.awaitNotConnected()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_awaitNotConnected().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also if
 * the waiting thread is interrupted). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_awaitNotConnected()")
int32_t                         dxfg_DXEndpoint_awaitNotConnected(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Returns the state of the endpoint.
 *
 * Java: `DXEndpoint.getState()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getState().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return The state of the endpoint, or -1 on error: it is not a value of #dxfg_endpoint_state_t, so compare the
 * result with -1 before using it. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getState()")
dxfg_endpoint_state_t           dxfg_DXEndpoint_getState(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Adds the state change listener (see dxfg_PropertyChangeListener_new()).
 *
 * The state changes are delivered by a single task on the endpoint's executor that reports the change from
 * the last reported state to the current one, so the fast changes are coalesced (e.g.
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED -> #DXFG_ENDPOINT_STATE_CONNECTED without #DXFG_ENDPOINT_STATE_CONNECTING).
 * A notification may come after dxfg_DXEndpoint_close() or dxfg_DXEndpoint_removeStateChangeListener() returns.
 *
 * The listener stays added until it is removed with dxfg_DXEndpoint_removeStateChangeListener(); releasing the last
 * handle of the listener only deactivates it.
 *
 * Java: `DXEndpoint.addStateChangeListener(PropertyChangeListener)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_addStateChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] listener The listener. The caller keeps its handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_addStateChangeListener()")
int32_t                         dxfg_DXEndpoint_addStateChangeListener(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint, dxfg_endpoint_state_change_listener_t *listener);

/**
 * @brief Removes the state change listener that was added with dxfg_DXEndpoint_addStateChangeListener().
 *
 * A notification that is already scheduled may still reach the listener after the function returns: release the
 * listener's handle with dxfg_JavaObjectHandler_release() (or deactivate it with dxfg_NativeListener_deactivate()) to
 * be sure that its function is not called anymore. The function does not release the handle.
 *
 * Java: `DXEndpoint.removeStateChangeListener(PropertyChangeListener)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_removeStateChangeListener().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] listener The listener to remove. The caller keeps its handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_removeStateChangeListener()")
int32_t                         dxfg_DXEndpoint_removeStateChangeListener(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint, dxfg_endpoint_state_change_listener_t *listener);

/**
 * @brief Returns the feed of the endpoint.
 *
 * The feed is associated with the endpoint: it receives the events from the endpoint's connection. Releasing the
 * handle of the feed does not affect the endpoint.
 *
 * Java: `DXEndpoint.getFeed()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getFeed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return The feed (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getFeed()")
dxfg_feed_t*                    dxfg_DXEndpoint_getFeed(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Returns the publisher of the endpoint.
 *
 * The publisher is associated with the endpoint: it publishes the events to the endpoint's connection. Releasing the
 * handle of the publisher does not affect the endpoint.
 *
 * Java: `DXEndpoint.getPublisher()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getPublisher().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return The publisher (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getPublisher()")
dxfg_publisher_t*               dxfg_DXEndpoint_getPublisher(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Changes the executor of the endpoint's notifications.
 *
 * The executor runs all the asynchronous interaction with the user code: event processing, state changes, etc.
 *
 * Call it before dxfg_DXEndpoint_getFeed() (and before the other functions that use the endpoint's feed,
 * e.g. dxfg_OnDemandService_getInstance2()). Otherwise, the feed keeps using the default executor, which is shut down,
 * and the last events requests (dxfg_DXFeed_getLastEventPromise(), etc.) never complete
 * (a `RejectedExecutionException` is logged).
 *
 * The executor must not reject the submitted tasks. A NULL executor is an error.
 *
 * Java: `DXEndpoint.executor(Executor)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_Builder_withExecutor(). The first call logs a deprecation warning.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] executor The executor. The endpoint does not take the handle: free it with
 * dxfg_JavaObjectHandler_release() when it is no longer needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_Builder_withExecutor()")
int32_t                         dxfg_DXEndpoint_executor(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint, dxfg_executor_t *executor);

/**
 * @brief Returns the event types that the endpoint supports.
 *
 * The result is the set of the event types supported by the endpoint and by the current version of the dxFeed Graal
 * Native SDK; the event types that the SDK does not know are left out.
 *
 * Java: `DXEndpoint.getEventTypes()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getEventTypes().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return The list of the supported event types (release it with dxfg_CList_EventClazz_release()), or NULL on error.
 * Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getEventTypes()")
dxfg_event_clazz_list_t *dxfg_DXEndpoint_getEventTypes(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Returns the aggregation period that the endpoint requests from the server.
 *
 * Java: `DXEndpoint.getRequestedAggregationPeriod()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getRequestedAggregationPeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] period The requested aggregation period (release the handle with dxfg_JavaObjectHandler_release()), or
 * NULL if no period is requested. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `period` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getRequestedAggregationPeriod()")
int32_t dxfg_DXEndpoint_getRequestedAggregationPeriod(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint,
                                                      DXFG_OUT dxfg_time_period_t **period);

/**
 * @brief Returns the aggregation period (in millis) that the endpoint requests from the server.
 *
 * Java: `DXEndpoint.getRequestedAggregationPeriod()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getRequestedAggregationPeriodMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] period The requested aggregation period in millis, or -1 if no period is requested. It is written on
 * success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `period` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getRequestedAggregationPeriodMillis()")
int32_t dxfg_DXEndpoint_getRequestedAggregationPeriodMillis(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint,
                                                            DXFG_OUT int64_t *period);

/**
 * @brief Sets the aggregation period that the endpoint requests from the server.
 *
 * The server keeps the period within its own bounds and reports the period it applies: see
 * dxfg_DXEndpoint_getAggregationPeriodInfo(). The period can also be requested with the
 * `dxfeed.requestedAggregationPeriod` property or the `[requestedAggregationPeriod=5s]` address option.
 *
 * Java: `DXEndpoint.setRequestedAggregationPeriod(TimePeriod)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_setRequestedAggregationPeriod().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] period The requested aggregation period: zero (dxfg_TimePeriod_ZERO()) requests real-time delivery, NULL
 * clears the request (the server uses its default period). It must not exceed 24 hours. The endpoint does not take the
 * handle: release it with dxfg_JavaObjectHandler_release() when it is no longer needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_setRequestedAggregationPeriod()")
int32_t dxfg_DXEndpoint_setRequestedAggregationPeriod(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint,
                                                      dxfg_time_period_t *period);

/**
 * @brief Sets the aggregation period (in millis) that the endpoint requests from the server.
 *
 * The same as dxfg_DXEndpoint_setRequestedAggregationPeriod(), without a time period handle.
 *
 * Java: `DXEndpoint.setRequestedAggregationPeriod(TimePeriod)`, called with `TimePeriod.valueOf(long)`, or with `null`
 * for a negative period.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_setRequestedAggregationPeriodMillis().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] period The requested aggregation period in millis: 0 requests real-time delivery, a negative value clears
 * the request (the server uses its default period). It must not exceed 24 hours (86400000).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_get_and_clear_thread_exception_t() to determine if an
 * exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_setRequestedAggregationPeriodMillis()")
int32_t dxfg_DXEndpoint_setRequestedAggregationPeriodMillis(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint,
                                                            int64_t period);

/**
 * @brief Returns the aggregation period that the servers apply, over all connections of the endpoint.
 *
 * The servers report the period when the connection is established. Both bounds are negative (unknown) until then,
 * and when no server reports the period (e.g. a local hub, a publisher or an older server).
 *
 * Java: `DXEndpoint.getAggregationPeriodInfo()`, `TimePeriodInfo.getMin()`, `TimePeriodInfo.getMax()`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_getAggregationPeriodInfo().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] min The minimum aggregation period in millis, negative if unknown. It is written on success only.
 * @param[out] max The maximum aggregation period in millis, negative if unknown. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `min` or `max` is NULL). Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_getAggregationPeriodInfo()")
int32_t dxfg_DXEndpoint_getAggregationPeriodInfo(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint,
                                                 DXFG_OUT int64_t *min, DXFG_OUT int64_t *max);

/**
 * @brief Creates a new endpoint state change listener.
 *
 * The listener does nothing until it is added to an endpoint with dxfg_DXEndpoint_addStateChangeListener(). The
 * function is called as described at dxfg_endpoint_state_change_listener_func.
 *
 * Releasing the last handle of the listener with dxfg_JavaObjectHandler_release() (or dxfg_NativeListener_deactivate())
 * stops the calls of the function: after the release returns the function is not running and is not called anymore,
 * so the user data can be freed. Remove the listener from the endpoint as well, otherwise it stays added and does
 * nothing.
 *
 * Java: `PropertyChangeListener.propertyChange(PropertyChangeEvent)`.
 *
 * @deprecated Use dxfg_v4_PropertyChangeListener_new().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call on the state changes (see dxfg_DXEndpoint_addStateChangeListener()).
 * @param[in] user_data The user data to pass to the function; the SDK does not use it. It must stay valid until the
 * listener is deactivated.
 * @return The listener (release its handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_PropertyChangeListener_new()")
dxfg_endpoint_state_change_listener_t* dxfg_PropertyChangeListener_new(graal_isolatethread_t *thread, dxfg_endpoint_state_change_listener_func user_func, void *user_data);

/** @} */ // end of Endpoint

#ifdef __cplusplus
}
#endif

#endif // DXFG_ENDPOINT_H
