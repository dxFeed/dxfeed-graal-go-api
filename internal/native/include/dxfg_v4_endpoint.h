// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_endpoint.h and the v4 names of its types (see the group @ref V4_Endpoint).
 */

#ifndef DXFG_V4_ENDPOINT_H
#define DXFG_V4_ENDPOINT_H

#include "dxfg_endpoint.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Endpoint Endpoint (v4)
 *  @brief The v4 functions and type names of the group @ref Endpoint. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_feed_t */
typedef dxfg_feed_t dxfg_v4_feed_t;

/** @copydoc dxfg_publisher_t */
typedef dxfg_publisher_t dxfg_v4_publisher_t;

/** @copydoc dxfg_event_clazz_list_t */
typedef dxfg_event_clazz_list_t dxfg_v4_event_clazz_list;

/** @copydoc dxfg_endpoint_t */
typedef dxfg_endpoint_t dxfg_v4_endpoint_t;

/** @copydoc dxfg_endpoint_builder_t */
typedef dxfg_endpoint_builder_t dxfg_v4_endpoint_builder_t;

/** @copydoc dxfg_endpoint_state_change_listener_t */
typedef dxfg_endpoint_state_change_listener_t dxfg_v4_endpoint_state_change_listener_t;

/** @copydoc dxfg_endpoint_role_t */
typedef dxfg_endpoint_role_t dxfg_v4_endpoint_role_t;

/** @copydoc dxfg_endpoint_state_t */
typedef dxfg_endpoint_state_t dxfg_v4_endpoint_state_t;

/** @copydoc dxfg_endpoint_state_change_listener_func */
typedef dxfg_endpoint_state_change_listener_func dxfg_v4_endpoint_state_change_listener_func;

/**
 * @brief Creates a new endpoint builder.
 *
 * The builder has the #DXFG_ENDPOINT_ROLE_FEED role and the default executor until they are changed.
 *
 * Java: `DXEndpoint.newBuilder()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] builder The pointer to the variable that receives the builder (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_newBuilder(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_endpoint_builder_t **builder);

/**
 * @brief Sets the role of the endpoints built by the builder.
 *
 * The default role is #DXFG_ENDPOINT_ROLE_FEED.
 *
 * Java: `DXEndpoint.Builder.withRole(DXEndpoint.Role)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] role The role of the endpoints.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_Builder_withRole(graal_isolatethread_t *thread, dxfg_v4_endpoint_builder_t *builder,
                                            dxfg_v4_endpoint_role_t role);

/**
 * @brief Sets the name of the endpoints built by the builder.
 *
 * The name distinguishes several endpoints in the same process in the logs and other diagnostics. It is the same as
 * dxfg_v4_DXEndpoint_Builder_withProperty() with the key `name`. A short name that describes the role of the endpoint
 * in the application is recommended.
 *
 * Java: `DXEndpoint.Builder.withName(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] name The name, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_Builder_withName(graal_isolatethread_t *thread, dxfg_v4_endpoint_builder_t *builder,
                                            const char *name);

/**
 * @brief Sets a property of the endpoints built by the builder.
 *
 * A property that the builder does not support is ignored: check it with dxfg_v4_DXEndpoint_Builder_supportsProperty().
 * A property that is set here takes precedence over the same property from the properties file and from the system
 * properties (see dxfg_system_set_property()).
 *
 * Java: `DXEndpoint.Builder.withProperty(String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] key The name of the property, a NUL-terminated UTF-8 string; the SDK copies it.
 * @param[in] value The value of the property, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_Builder_withProperty(graal_isolatethread_t *thread, dxfg_v4_endpoint_builder_t *builder,
                                                const char *key, const char *value);

/**
 * @brief Sets the properties of the endpoints built by the builder from a properties file.
 *
 * The file is read in the format of a Java properties file (`java.util.Properties.load`) and every property that the
 * builder supports is set as by dxfg_v4_DXEndpoint_Builder_withProperty(). The unsupported properties are ignored.
 * An error is reported if the file cannot be read. The file is closed before the function returns, also on an error
 * (since 3.16.0; before, it stayed open until the garbage collection, and on Windows it could not be deleted meanwhile).
 *
 * Java: `DXEndpoint.Builder.withProperties(Properties)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] file_path The path of the properties file, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_Builder_withProperties(graal_isolatethread_t *thread, dxfg_v4_endpoint_builder_t *builder,
                                                  const char *file_path);

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
 * The executor of a shared event queue (dxfg_v4_EventQueue_getExecutor()) makes the polls of the queue call the
 * listeners of the endpoint and of the subscriptions that inherit it; those subscriptions are not attached to the queue
 * (see @ref V4_EventQueue for the limits). The executor of a subscription attached to a queue is refused with
 * `IllegalStateException`.
 *
 * Java: `DXEndpoint.Builder.withExecutor(Executor)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] executor The executor (e.g. from dxfg_v4_Executors_newFixedThreadPool()) or NULL to use the default one.
 * The builder does not take the handle: free it with dxfg_v4_JavaObjectHandler_release() when it is no longer needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_Builder_withExecutor(graal_isolatethread_t *thread, dxfg_v4_endpoint_builder_t *builder,
                                                dxfg_v4_executor_t *executor);

/**
 * @brief Checks whether the builder supports a property.
 *
 * The properties that are not supported are ignored by dxfg_v4_DXEndpoint_Builder_withProperty().
 *
 * Java: `DXEndpoint.Builder.supportsProperty(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[in] key The name of the property, a NUL-terminated UTF-8 string.
 * @param[out] supported The pointer to the variable that receives 1 if the property is supported, 0 if it is not
 * supported; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_Builder_supportsProperty(graal_isolatethread_t *thread, dxfg_v4_endpoint_builder_t *builder,
                                                    const char *key, DXFG_OUT int32_t *supported);

/**
 * @brief Builds a new endpoint with the builder's settings.
 *
 * Every call builds a new endpoint. An endpoint with the role #DXFG_ENDPOINT_ROLE_FEED or #DXFG_ENDPOINT_ROLE_PUBLISHER
 * connects automatically on creation if `dxfeed.address` (feed) or `dxpublisher.address` (publisher) is set; the
 * endpoints of the other roles do not.
 *
 * Java: `DXEndpoint.Builder.build()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder.
 * @param[out] endpoint The pointer to the variable that receives the new endpoint. The caller owns the endpoint: close
 * it with dxfg_v4_DXEndpoint_close() or dxfg_v4_DXEndpoint_closeAndAwaitTermination() and free its handle with
 * dxfg_v4_JavaObjectHandler_release() (releasing the handle does not close the endpoint); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_Builder_build(graal_isolatethread_t *thread, dxfg_v4_endpoint_builder_t *builder,
                                         DXFG_OUT dxfg_v4_endpoint_t **endpoint);

/**
 * @brief Returns the default application-wide singleton endpoint with the #DXFG_ENDPOINT_ROLE_FEED role, the same one
 * that dxfg_v4_DXFeed_getInstance() uses.
 *
 * The endpoint is created on the first call and is configured by the properties of the endpoint and the system
 * properties (see dxfg_system_set_property()) as of that moment.
 *
 * All callers share this endpoint: closing it closes it for everyone in the process.
 *
 * Java: `DXEndpoint.getInstance()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] endpoint The pointer to the variable that receives the endpoint (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getInstance(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_endpoint_t **endpoint);

/**
 * @brief Creates a <b>new</b> endpoint with the given role on every call, the same as
 * dxfg_v4_DXEndpoint_createWithRole().
 *
 * Unlike the Java `DXEndpoint.getInstance(Role)`, this function does not return a singleton: the SDK calls
 * `DXEndpoint.create(Role)`. Prefer dxfg_v4_DXEndpoint_newBuilder() (or dxfg_v4_DXEndpoint_createWithRole()) in the new
 * code.
 *
 * Java: `DXEndpoint.create(DXEndpoint.Role)`.
 *
 * @deprecated Use dxfg_v4_DXEndpoint_createWithRole(), which does the same: the name of this function promises the
 * singleton of the Java `DXEndpoint.getInstance(Role)`, which it does not return. Deprecated since 3.16.0, removed in
 * 4.0.0.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] role The role of the endpoint.
 * @param[out] endpoint The pointer to the variable that receives the new endpoint. The caller owns the endpoint: close
 * it with dxfg_v4_DXEndpoint_close() or dxfg_v4_DXEndpoint_closeAndAwaitTermination() and free its handle with
 * dxfg_v4_JavaObjectHandler_release(); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_DXEndpoint_createWithRole()")
int32_t dxfg_v4_DXEndpoint_getInstanceWithRole(graal_isolatethread_t *thread, dxfg_v4_endpoint_role_t role,
                                               DXFG_OUT dxfg_v4_endpoint_t **endpoint);

/**
 * @brief Creates a new endpoint with the #DXFG_ENDPOINT_ROLE_FEED role.
 *
 * Java: `DXEndpoint.create()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] endpoint The pointer to the variable that receives the new endpoint. The caller owns the endpoint: close
 * it with dxfg_v4_DXEndpoint_close() or dxfg_v4_DXEndpoint_closeAndAwaitTermination() and free its handle with
 * dxfg_v4_JavaObjectHandler_release(); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_create(graal_isolatethread_t *thread, DXFG_OUT dxfg_v4_endpoint_t **endpoint);

/**
 * @brief Creates a new endpoint with the given role.
 *
 * Java: `DXEndpoint.create(DXEndpoint.Role)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] role The role of the endpoint.
 * @param[out] endpoint The pointer to the variable that receives the new endpoint. The caller owns the endpoint: close
 * it with dxfg_v4_DXEndpoint_close() or dxfg_v4_DXEndpoint_closeAndAwaitTermination() and free its handle with
 * dxfg_v4_JavaObjectHandler_release(); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_createWithRole(graal_isolatethread_t *thread, dxfg_v4_endpoint_role_t role,
                                          DXFG_OUT dxfg_v4_endpoint_t **endpoint);

/**
 * @brief Closes the endpoint: disconnects it and releases its resources. The endpoint cannot be used after that.
 *
 * The state of the endpoint becomes #DXFG_ENDPOINT_STATE_CLOSED immediately.
 *
 * The function does not wait for the notifications that are already scheduled: a state change listener
 * (including the change to #DXFG_ENDPOINT_STATE_CLOSED) or an event listener may still be called after it returns.
 * Use dxfg_v4_DXEndpoint_closeAndAwaitTermination() to wait for them, or dxfg_v4_JavaObjectHandler_release() of the
 * listener to stop its calls (see dxfg_v4_JavaObjectHandler_release()).
 *
 * Java: `DXEndpoint.close()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_close(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint);

/**
 * @brief Closes the endpoint and waits until all the data processing and the scheduled notifications are completed.
 *
 * It does the same as dxfg_v4_DXEndpoint_close() and blocks until the outstanding tasks of the endpoint are finished.
 * It is designed for the #DXFG_ENDPOINT_ROLE_STREAM_FEED role, after dxfg_v4_DXEndpoint_awaitNotConnected() returns, to
 * make sure that a file was completely processed.
 *
 * Do not call it from a listener of this endpoint: it would wait for the notification that is calling it.
 *
 * Java: `DXEndpoint.closeAndAwaitTermination()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also if
 * the waiting thread is interrupted). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 */
int32_t dxfg_v4_DXEndpoint_closeAndAwaitTermination(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint);

/**
 * @brief Returns the role of the endpoint.
 *
 * Java: `DXEndpoint.getRole()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] role The pointer to the variable that receives the role of the endpoint; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getRole(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                   DXFG_OUT dxfg_v4_endpoint_role_t *role);

/**
 * @brief Sets the user name of the endpoint.
 *
 * Call it before dxfg_v4_DXEndpoint_connect(), together with dxfg_v4_DXEndpoint_password(), to configure the
 * credentials of a service that requires them. A NULL user name is an error.
 *
 * Java: `DXEndpoint.user(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] user The user name, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_user(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint, const char *user);

/**
 * @brief Sets the password of the endpoint.
 *
 * Call it before dxfg_v4_DXEndpoint_connect(), together with dxfg_v4_DXEndpoint_user(), to configure the credentials of
 * a service that requires them. A NULL password is an error.
 *
 * Java: `DXEndpoint.password(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] password The password, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_password(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint, const char *password);

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
 * services that require credentials call dxfg_v4_DXEndpoint_user() and dxfg_v4_DXEndpoint_password() first.
 *
 * Java: `DXEndpoint.connect(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] address The address of the data source, a NUL-terminated UTF-8 string; the SDK copies it.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_connect(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint, const char *address);

/**
 * @brief Terminates all the network connections of the endpoint and initiates connecting again with the same address.
 *
 * It works like dxfg_v4_DXEndpoint_disconnect() followed by dxfg_v4_DXEndpoint_connect() with the current address, but
 * the internal resources may be reused. A connection with several target addresses tries to switch to an alternative
 * address; the configured reconnect timeouts apply.
 *
 * The function does not connect an endpoint that was not connected with dxfg_v4_DXEndpoint_connect() or was
 * disconnected with dxfg_v4_DXEndpoint_disconnect(). Whether a state change listener sees the intermediate
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED state is not specified.
 *
 * Java: `DXEndpoint.reconnect()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_reconnect(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint);

/**
 * @brief Terminates all the remote network connections of the endpoint.
 *
 * The function does nothing if the endpoint is closed. Otherwise, the state of the endpoint immediately becomes
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED. The resources of the endpoint are not released: use dxfg_v4_DXEndpoint_close()
 * for that.
 *
 * Java: `DXEndpoint.disconnect()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_disconnect(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint);

/**
 * @brief Terminates all the remote network connections of the endpoint and clears the stored data.
 *
 * The function does nothing if the endpoint is closed. Otherwise, the state of the endpoint immediately becomes
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED. The resources of the endpoint are not released: use dxfg_v4_DXEndpoint_close()
 * for that.
 *
 * Java: `DXEndpoint.disconnectAndClear()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_disconnectAndClear(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also if
 * the waiting thread is interrupted). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 */
int32_t dxfg_v4_DXEndpoint_awaitProcessed(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint);

/**
 * @brief Waits until the state of the endpoint becomes #DXFG_ENDPOINT_STATE_NOT_CONNECTED or
 * #DXFG_ENDPOINT_STATE_CLOSED.
 *
 * It is the signal that the files opened with dxfg_v4_DXEndpoint_connect() (`file:...` addresses) were finished
 * reading, but not necessarily completely processed by the subscription listeners. Call
 * dxfg_v4_DXEndpoint_closeAndAwaitTermination() after it returns to make sure that all the processing has completed.
 *
 * The function blocks the calling thread. With a single-threaded executor of the endpoint, do not call it from the
 * executor's thread (e.g. from a listener): it would wait forever.
 *
 * Java: `DXEndpoint.awaitNotConnected()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also if
 * the waiting thread is interrupted). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 */
int32_t dxfg_v4_DXEndpoint_awaitNotConnected(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint);

/**
 * @brief Returns the state of the endpoint.
 *
 * Java: `DXEndpoint.getState()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] state The pointer to the variable that receives the state of the endpoint; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getState(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                    DXFG_OUT dxfg_v4_endpoint_state_t *state);

/**
 * @brief Adds the state change listener (see dxfg_v4_PropertyChangeListener_new()).
 *
 * The state changes are delivered by a single task on the endpoint's executor that reports the change from
 * the last reported state to the current one, so the fast changes are coalesced (e.g.
 * #DXFG_ENDPOINT_STATE_NOT_CONNECTED -> #DXFG_ENDPOINT_STATE_CONNECTED without #DXFG_ENDPOINT_STATE_CONNECTING).
 * A notification may come after dxfg_v4_DXEndpoint_close() or dxfg_v4_DXEndpoint_removeStateChangeListener() returns.
 *
 * The listener stays added until it is removed with dxfg_v4_DXEndpoint_removeStateChangeListener(); releasing the last
 * handle of the listener only deactivates it.
 *
 * Java: `DXEndpoint.addStateChangeListener(PropertyChangeListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] listener The listener. The caller keeps its handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_addStateChangeListener(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                                  dxfg_v4_endpoint_state_change_listener_t *listener);

/**
 * @brief Removes the state change listener that was added with dxfg_v4_DXEndpoint_addStateChangeListener().
 *
 * A notification that is already scheduled may still reach the listener after the function returns: release the
 * listener's handle with dxfg_v4_JavaObjectHandler_release() (or deactivate it with
 * dxfg_v4_NativeListener_deactivate()) to be sure that its function is not called anymore. The function does not
 * release the handle.
 *
 * Java: `DXEndpoint.removeStateChangeListener(PropertyChangeListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] listener The listener to remove. The caller keeps its handle.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_removeStateChangeListener(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                                     dxfg_v4_endpoint_state_change_listener_t *listener);

/**
 * @brief Returns the feed of the endpoint.
 *
 * The feed is associated with the endpoint: it receives the events from the endpoint's connection. Releasing the
 * handle of the feed does not affect the endpoint.
 *
 * Java: `DXEndpoint.getFeed()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] feed The pointer to the variable that receives the feed (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getFeed(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                   DXFG_OUT dxfg_v4_feed_t **feed);

/**
 * @brief Returns the publisher of the endpoint.
 *
 * The publisher is associated with the endpoint: it publishes the events to the endpoint's connection. Releasing the
 * handle of the publisher does not affect the endpoint.
 *
 * Java: `DXEndpoint.getPublisher()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] publisher The pointer to the variable that receives the publisher (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getPublisher(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                        DXFG_OUT dxfg_v4_publisher_t **publisher);

/**
 * @brief Returns the event types that the endpoint supports.
 *
 * The result is the set of the event types supported by the endpoint and by the current version of the dxFeed Graal
 * Native SDK; the event types that the SDK does not know are left out.
 *
 * Java: `DXEndpoint.getEventTypes()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] eventTypes The pointer to the variable that receives the list of the supported event types (release it
 * with dxfg_v4_event_clazz_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getEventTypes(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                         DXFG_OUT dxfg_v4_event_clazz_list **eventTypes);

/**
 * @brief Returns the aggregation period that the endpoint requests from the server.
 *
 * Java: `DXEndpoint.getRequestedAggregationPeriod()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] period The requested aggregation period (release the handle with dxfg_v4_JavaObjectHandler_release()), or
 * NULL if no period is requested. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `period` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getRequestedAggregationPeriod(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                                         DXFG_OUT dxfg_v4_time_period_t **period);

/**
 * @brief Returns the aggregation period (in millis) that the endpoint requests from the server.
 *
 * Java: `DXEndpoint.getRequestedAggregationPeriod()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] period The requested aggregation period in millis, or -1 if no period is requested. It is written on
 * success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `period` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getRequestedAggregationPeriodMillis(graal_isolatethread_t *thread,
                                                               dxfg_v4_endpoint_t *endpoint, DXFG_OUT int64_t *period);

/**
 * @brief Sets the aggregation period that the endpoint requests from the server.
 *
 * The server keeps the period within its own bounds and reports the period it applies: see
 * dxfg_v4_DXEndpoint_getAggregationPeriodInfo(). The period can also be requested with the
 * `dxfeed.requestedAggregationPeriod` property or the `[requestedAggregationPeriod=5s]` address option.
 *
 * Java: `DXEndpoint.setRequestedAggregationPeriod(TimePeriod)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] period The requested aggregation period: zero (dxfg_v4_TimePeriod_ZERO()) requests real-time delivery,
 * NULL clears the request (the server uses its default period). It must not exceed 24 hours. The endpoint does not take
 * the handle: release it with dxfg_v4_JavaObjectHandler_release() when it is no longer needed.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_v4_get_and_clear_thread_exception() to determine if an
 * exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_setRequestedAggregationPeriod(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                                         dxfg_v4_time_period_t *period);

/**
 * @brief Sets the aggregation period (in millis) that the endpoint requests from the server.
 *
 * The same as dxfg_v4_DXEndpoint_setRequestedAggregationPeriod(), without a time period handle.
 *
 * Java: `DXEndpoint.setRequestedAggregationPeriod(TimePeriod)`, called with `TimePeriod.valueOf(long)`, or with `null`
 * for a negative period.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[in] period The requested aggregation period in millis: 0 requests real-time delivery, a negative value clears
 * the request (the server uses its default period). It must not exceed 24 hours (86400000).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (for
 * example, when the period exceeds 24 hours). Use dxfg_v4_get_and_clear_thread_exception() to determine if an
 * exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_setRequestedAggregationPeriodMillis(graal_isolatethread_t *thread,
                                                               dxfg_v4_endpoint_t *endpoint, int64_t period);

/**
 * @brief Returns the aggregation period that the servers apply, over all connections of the endpoint.
 *
 * The servers report the period when the connection is established. Both bounds are negative (unknown) until then,
 * and when no server reports the period (e.g. a local hub, a publisher or an older server).
 *
 * Java: `DXEndpoint.getAggregationPeriodInfo()`, `TimePeriodInfo.getMin()`, `TimePeriodInfo.getMax()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint.
 * @param[out] min The minimum aggregation period in millis, negative if unknown. It is written on success only.
 * @param[out] max The maximum aggregation period in millis, negative if unknown. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error (also
 * when `min` or `max` is NULL). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_DXEndpoint_getAggregationPeriodInfo(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                                    DXFG_OUT int64_t *min, DXFG_OUT int64_t *max);

/**
 * @brief Creates a new endpoint state change listener.
 *
 * The listener does nothing until it is added to an endpoint with dxfg_v4_DXEndpoint_addStateChangeListener(). The
 * function is called as described at dxfg_v4_endpoint_state_change_listener_func.
 *
 * Releasing the last handle of the listener with dxfg_v4_JavaObjectHandler_release() (or
 * dxfg_v4_NativeListener_deactivate()) stops the calls of the function: after the release returns the function is not
 * running and is not called anymore, so the user data can be freed. Remove the listener from the endpoint as well,
 * otherwise it stays added and does nothing.
 *
 * Java: `PropertyChangeListener.propertyChange(PropertyChangeEvent)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call on the state changes (see dxfg_v4_DXEndpoint_addStateChangeListener()).
 * @param[in] user_data The user data to pass to the function; the SDK does not use it. It must stay valid until the
 * listener is deactivated.
 * @param[out] listener The pointer to the variable that receives the listener (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_PropertyChangeListener_new(graal_isolatethread_t *thread,
                                           dxfg_v4_endpoint_state_change_listener_func user_func, void *user_data,
                                           DXFG_OUT dxfg_v4_endpoint_state_change_listener_t **listener);

/** @} */ // end of V4_Endpoint

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_ENDPOINT_H
