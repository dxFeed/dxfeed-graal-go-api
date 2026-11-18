// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK CandleWebService, HistoryEndpoint functions and types declarations
 */

#ifndef DXFG_CANDLEWEBSERVICE_H
#define DXFG_CANDLEWEBSERVICE_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "dxfg_javac.h"
#include "graal_isolate.h"

#include "dxfg_events.h"

/** @defgroup CandleWebService CandleWebService
 *  @brief The client of the history (candle web) service: a request of the time series events of a symbol within a
 *  time range over HTTP, configured with a builder.
 *
 *  dxfg_HistoryEndpoint_newBuilder() creates a builder, the functions dxfg_HistoryEndpoint_Builder_withAddress(),
 *  dxfg_HistoryEndpoint_Builder_withUserName() and the other `dxfg_HistoryEndpoint_Builder_with*` functions set its
 *  parameters and dxfg_HistoryEndpoint_Builder_build() creates the endpoint, which dxfg_HistoryEndpoint_getTimeSeries()
 *  uses to request the events.
 *
 *  A handle that a function returns is owned by the caller: release it with dxfg_JavaObjectHandler_release() (see the
 *  group Javac for the ownership rules). All the functions of this group return #DXFG_EXECUTE_SUCCESSFULLY (0) or
 *  #DXFG_EXECUTE_FAIL (-1) on error, and the result is passed through an `[out]` parameter; it is not changed on
 *  error. A NULL `[out]` pointer is an error too. dxfg_get_and_clear_thread_exception_t() gives the exception.
 *  @{
 */

/**
 * @brief A handle of a history endpoint, the client of the history service; release it with
 * dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_HistoryEndpoint_Builder_build(). The endpoint holds only the settings of the requests: it does
 * not keep a connection open, there is nothing to close. Releasing the handle drops the reference to the settings.
 */
typedef struct dxfg_history_endpoint_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_history_endpoint_t;

/**
 * @brief A handle of a builder of history endpoints (Java `HistoryEndpoint.Builder`); release it with
 * dxfg_JavaObjectHandler_release().
 *
 * Create it with dxfg_HistoryEndpoint_newBuilder(), set the parameters with the dxfg_HistoryEndpoint_Builder_with*()
 * functions (they change this builder) and create an endpoint with dxfg_HistoryEndpoint_Builder_build(). The builder
 * stays valid after that: it can be changed and used to build more endpoints.
 */
typedef struct dxfg_history_endpoint_builder_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_history_endpoint_builder_t;

/**
 * @brief The compression of the data that the history service sends.
 *
 * Java: `HistoryEndpoint.Compression`. The compression is requested from the service and the response is
 * decompressed accordingly.
 */
typedef enum dxfg_history_endpoint_compression_t {
    /** The data is not compressed. */
    DXFG_HISTORY_ENDPOINT_COMPRESSION_NONE = 0,
    /** The data is compressed with ZIP. */
    DXFG_HISTORY_ENDPOINT_COMPRESSION_ZIP,
    /** The data is compressed with GZIP. */
    DXFG_HISTORY_ENDPOINT_COMPRESSION_GZIP,
} dxfg_history_endpoint_compression_t;

/**
 * @brief The format of the data that the history service sends.
 *
 * Java: `HistoryEndpoint.Format`. The format is requested from the service and the response is parsed accordingly.
 */
typedef enum dxfg_history_endpoint_format_t {
    /** The data is plain text. */
    DXFG_HISTORY_ENDPOINT_FORMAT_TEXT = 0,
    /** The data is comma-separated values. */
    DXFG_HISTORY_ENDPOINT_FORMAT_CSV,
    /** The data is in the binary format. */
    DXFG_HISTORY_ENDPOINT_FORMAT_BINARY,
} dxfg_history_endpoint_format_t;

/**
 * @brief Creates a builder of history endpoints with the default settings.
 *
 * The default settings are:
 * - #DXFG_HISTORY_ENDPOINT_COMPRESSION_ZIP as the compression;
 * - #DXFG_HISTORY_ENDPOINT_FORMAT_BINARY as the format.
 *
 * The address, the user name, the password and the authentication token are not set.
 *
 * Java: `HistoryEndpoint.newBuilder()`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_newBuilder().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] builder The new builder (release the handle with dxfg_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_newBuilder()")
int32_t dxfg_HistoryEndpoint_newBuilder(graal_isolatethread_t *thread,
                                        DXFG_OUT dxfg_history_endpoint_builder_t **builder);

/**
 * @brief Requests the time series events of a symbol within a time range from the history service.
 *
 * The function blocks the calling thread while the request is executed and the response is read: it sends a request
 * to the address of the endpoint with the query parameters `records` (the simple name of the event class), `symbols`,
 * `start` and `stop` (the range), `format` and `compression` (the settings of the endpoint), then decompresses and
 * parses the response. It fails if the request, the reading or the parsing fails.
 *
 * `eventClazz` must be the class of a time series event (for example #DXFG_EVENT_CANDLE or
 * #DXFG_EVENT_TIME_AND_SALE); the sources do not say what happens for another class.
 *
 * Java: `HistoryEndpoint.getTimeSeries(Class, Object, Long, Long)`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_getTimeSeries().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The history endpoint.
 * @param[in] eventClazz The class of the events to retrieve.
 * @param[in] symbol The symbol whose time series is requested; its string form is sent to the service. The caller
 * keeps its ownership.
 * @param[in] fromTime The start of the time range, in milliseconds since the Unix epoch.
 * @param[in] toTime The end of the time range, in milliseconds since the Unix epoch.
 * @param[out] events The list of the events of the specified class and symbol within the time range; it can be empty.
 * Release it with dxfg_CList_EventType_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_getTimeSeries()")
int32_t dxfg_HistoryEndpoint_getTimeSeries(graal_isolatethread_t *thread, dxfg_history_endpoint_t *endpoint,
                                           dxfg_event_clazz_t eventClazz, dxfg_symbol_t *symbol, int64_t fromTime,
                                           int64_t toTime, DXFG_OUT dxfg_event_type_list **events);

/**
 * @brief Sets the address of the history service.
 *
 * The address is the URL of the service without the query: the SDK appends the query of each request to it. The
 * string is copied: the caller keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withAddress(String)`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_Builder_withAddress().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] address The address of the service: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_Builder_withAddress()")
int32_t dxfg_HistoryEndpoint_Builder_withAddress(graal_isolatethread_t *thread,
                                                 dxfg_history_endpoint_builder_t *builder, const char *address);

/**
 * @brief Sets the user name that the requests to the history service use.
 *
 * The user name and the password are used to open the connection to the service. The string is copied: the caller
 * keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withUserName(String)`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_Builder_withUserName().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] userName The user name: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_Builder_withUserName()")
int32_t dxfg_HistoryEndpoint_Builder_withUserName(graal_isolatethread_t *thread,
                                                  dxfg_history_endpoint_builder_t *builder, const char *userName);

/**
 * @brief Sets the password that the requests to the history service use.
 *
 * The user name and the password are used to open the connection to the service. The string is copied: the caller
 * keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withPassword(String)`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_Builder_withPassword().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] password The password: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_Builder_withPassword()")
int32_t dxfg_HistoryEndpoint_Builder_withPassword(graal_isolatethread_t *thread,
                                                  dxfg_history_endpoint_builder_t *builder, const char *password);

/**
 * @brief Sets the authentication token that the requests to the history service use.
 *
 * If the token is not NULL and not empty, each request carries it in the `Authorization: Bearer <token>` header. The
 * string is copied: the caller keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withAuthToken(String)`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_Builder_withAuthToken().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] authToken The authentication token: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_Builder_withAuthToken()")
int32_t dxfg_HistoryEndpoint_Builder_withAuthToken(graal_isolatethread_t *thread,
                                                   dxfg_history_endpoint_builder_t *builder, const char *authToken);

/**
 * @brief Sets the compression of the data that the history service sends.
 *
 * Java: `HistoryEndpoint.Builder.withCompression(HistoryEndpoint.Compression)`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_Builder_withCompression().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] compression The compression, see #dxfg_history_endpoint_compression_t.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_Builder_withCompression()")
int32_t dxfg_HistoryEndpoint_Builder_withCompression(graal_isolatethread_t *thread,
                                                     dxfg_history_endpoint_builder_t *builder,
                                                     dxfg_history_endpoint_compression_t compression);

/**
 * @brief Sets the format of the data that the history service sends.
 *
 * Java: `HistoryEndpoint.Builder.withFormat(HistoryEndpoint.Format)`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_Builder_withFormat().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] format The format, see #dxfg_history_endpoint_format_t.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_Builder_withFormat()")
int32_t dxfg_HistoryEndpoint_Builder_withFormat(graal_isolatethread_t *thread, dxfg_history_endpoint_builder_t *builder,
                                                dxfg_history_endpoint_format_t format);

/**
 * @brief Creates a history endpoint with the settings of the builder.
 *
 * The endpoint takes the address, the user name, the password, the authentication token, the compression and the
 * format that are set in the builder at the moment of the call; later changes of the builder do not affect it. The
 * function does not check that the address is set.
 *
 * Java: `HistoryEndpoint.Builder.build()`.
 *
 * @deprecated Use dxfg_v4_HistoryEndpoint_Builder_build().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[out] endpoint The new endpoint (release the handle with dxfg_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_HistoryEndpoint_Builder_build()")
int32_t dxfg_HistoryEndpoint_Builder_build(graal_isolatethread_t *thread, dxfg_history_endpoint_builder_t *builder,
                                           DXFG_OUT dxfg_history_endpoint_t **endpoint);

/** @} */ // end of CandleWebService

#ifdef __cplusplus
}
#endif

#endif // DXFG_CANDLEWEBSERVICE_H
