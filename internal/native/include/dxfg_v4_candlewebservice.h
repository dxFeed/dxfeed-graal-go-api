// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_candlewebservice.h and the v4 names of its types (see the group @ref V4_CandleWebService).
 */

#ifndef DXFG_V4_CANDLEWEBSERVICE_H
#define DXFG_V4_CANDLEWEBSERVICE_H

#include "dxfg_candlewebservice.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"
#include "dxfg_v4_events.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_CandleWebService CandleWebService (v4)
 *  @brief The v4 functions and type names of the group @ref CandleWebService. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
 *  @{
 */

/** @copydoc dxfg_history_endpoint_t */
typedef dxfg_history_endpoint_t dxfg_v4_history_endpoint_t;

/** @copydoc dxfg_history_endpoint_builder_t */
typedef dxfg_history_endpoint_builder_t dxfg_v4_history_endpoint_builder_t;

/** @copydoc dxfg_history_endpoint_compression_t */
typedef dxfg_history_endpoint_compression_t dxfg_v4_history_endpoint_compression_t;

/** @copydoc dxfg_history_endpoint_format_t */
typedef dxfg_history_endpoint_format_t dxfg_v4_history_endpoint_format_t;

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[out] builder The new builder (release the handle with dxfg_v4_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_newBuilder(graal_isolatethread_t *thread,
                                           DXFG_OUT dxfg_v4_history_endpoint_builder_t **builder);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The history endpoint.
 * @param[in] eventClazz The class of the events to retrieve.
 * @param[in] symbol The symbol whose time series is requested; its string form is sent to the service. The caller
 * keeps its ownership.
 * @param[in] fromTime The start of the time range, in milliseconds since the Unix epoch.
 * @param[in] toTime The end of the time range, in milliseconds since the Unix epoch.
 * @param[out] events The list of the events of the specified class and symbol within the time range; it can be empty.
 * Release it with dxfg_v4_event_type_list_release().
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_getTimeSeries(graal_isolatethread_t *thread, dxfg_v4_history_endpoint_t *endpoint,
                                              dxfg_v4_event_clazz_t eventClazz, dxfg_v4_symbol_t *symbol,
                                              int64_t fromTime, int64_t toTime,
                                              DXFG_OUT dxfg_v4_event_type_list **events);

/**
 * @brief Sets the address of the history service.
 *
 * The address is the URL of the service without the query: the SDK appends the query of each request to it. The
 * string is copied: the caller keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withAddress(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] address The address of the service: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_Builder_withAddress(graal_isolatethread_t *thread,
                                                    dxfg_v4_history_endpoint_builder_t *builder, const char *address);

/**
 * @brief Sets the user name that the requests to the history service use.
 *
 * The user name and the password are used to open the connection to the service. The string is copied: the caller
 * keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withUserName(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] userName The user name: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_Builder_withUserName(graal_isolatethread_t *thread,
                                                     dxfg_v4_history_endpoint_builder_t *builder, const char *userName);

/**
 * @brief Sets the password that the requests to the history service use.
 *
 * The user name and the password are used to open the connection to the service. The string is copied: the caller
 * keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withPassword(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] password The password: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_Builder_withPassword(graal_isolatethread_t *thread,
                                                     dxfg_v4_history_endpoint_builder_t *builder, const char *password);

/**
 * @brief Sets the authentication token that the requests to the history service use.
 *
 * If the token is not NULL and not empty, each request carries it in the `Authorization: Bearer <token>` header. The
 * string is copied: the caller keeps its ownership and can free it after the call.
 *
 * Java: `HistoryEndpoint.Builder.withAuthToken(String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] authToken The authentication token: a UTF-8 null-terminated string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_Builder_withAuthToken(graal_isolatethread_t *thread,
                                                      dxfg_v4_history_endpoint_builder_t *builder,
                                                      const char *authToken);

/**
 * @brief Sets the compression of the data that the history service sends.
 *
 * Java: `HistoryEndpoint.Builder.withCompression(HistoryEndpoint.Compression)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] compression The compression, see #dxfg_v4_history_endpoint_compression_t.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_Builder_withCompression(graal_isolatethread_t *thread,
                                                        dxfg_v4_history_endpoint_builder_t *builder,
                                                        dxfg_v4_history_endpoint_compression_t compression);

/**
 * @brief Sets the format of the data that the history service sends.
 *
 * Java: `HistoryEndpoint.Builder.withFormat(HistoryEndpoint.Format)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[in] format The format, see #dxfg_v4_history_endpoint_format_t.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_Builder_withFormat(graal_isolatethread_t *thread,
                                                   dxfg_v4_history_endpoint_builder_t *builder,
                                                   dxfg_v4_history_endpoint_format_t format);

/**
 * @brief Creates a history endpoint with the settings of the builder.
 *
 * The endpoint takes the address, the user name, the password, the authentication token, the compression and the
 * format that are set in the builder at the moment of the call; later changes of the builder do not affect it. The
 * function does not check that the address is set.
 *
 * Java: `HistoryEndpoint.Builder.build()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The history endpoint's builder.
 * @param[out] endpoint The new endpoint (release the handle with dxfg_v4_JavaObjectHandler_release()).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_HistoryEndpoint_Builder_build(graal_isolatethread_t *thread,
                                              dxfg_v4_history_endpoint_builder_t *builder,
                                              DXFG_OUT dxfg_v4_history_endpoint_t **endpoint);

/** @} */ // end of V4_CandleWebService

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_CANDLEWEBSERVICE_H
