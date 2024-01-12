// SPDX-License-Identifier: MPL-2.0

#ifndef DXFEED_GRAAL_NATIVE_SDK_IPF_H_
#define DXFEED_GRAAL_NATIVE_SDK_IPF_H_

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "graal_isolate.h"
#include "dxfg_javac.h"

/** @defgroup InstrumentProfile
 *  @{
 */

/**
 * @brief The InstrumentProfileConnection.
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileConnection.html">Javadoc</a>
 */
typedef struct dxfg_ipf_connection_t {
    dxfg_java_object_handler handler;
} dxfg_ipf_connection_t;

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/InstrumentProfileReader.html">Javadoc</a>
 */
typedef struct dxfg_instrument_profile_reader_t {
    dxfg_java_object_handler handler;
} dxfg_instrument_profile_reader_t;

/**
 * @brief The PropertyChangeListener.
 * <a href="https://docs.oracle.com/javase/8/docs/api/java/beans/PropertyChangeListener.html">Javadoc</a>
 */
typedef struct dxfg_ipf_connection_state_change_listener_t {
    dxfg_java_object_handler handler;
} dxfg_ipf_connection_state_change_listener_t;

/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileConnection.State.html">Javadoc</a>
 */
typedef enum dxfg_ipf_connection_state_t {
    DXFG_IPF_CONNECTION_STATE_NOT_CONNECTED = 0,
    DXFG_IPF_CONNECTION_STATE_CONNECTING,
    DXFG_IPF_CONNECTION_STATE_CONNECTED,
    DXFG_IPF_CONNECTION_STATE_COMPLETED,
    DXFG_IPF_CONNECTION_STATE_CLOSED,
} dxfg_ipf_connection_state_t;

typedef void (*dxfg_ipf_connection_state_change_listener_func)(graal_isolatethread_t *thread, dxfg_ipf_connection_state_t old_state, dxfg_ipf_connection_state_t new_state, void *user_data);

/**
 * @brief The InstrumentProfileCollector.
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileCollector.html">Javadoc</a>
 */
typedef struct dxfg_ipf_collector_t {
    dxfg_java_object_handler handler;
} dxfg_ipf_collector_t;


/**
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/InstrumentProfile.html">Javadoc</a>
 */
typedef struct dxfg_instrument_profile_t {
    const char *type;
    const char *symbol;
    const char *description;
    const char *local_symbol;
    const char *local_description;
    const char *country;
    const char *opol;
    const char *exchange_data;
    const char *exchanges;
    const char *currency;
    const char *base_currency;
    const char *cfi;
    const char *isin;
    const char *sedol;
    const char *cusip;
    int32_t icb;
    int32_t sic;
    double multiplier;
    const char *product;
    const char *underlying;
    double spc;
    const char *additional_underlyings;
    const char *mmy;
    int32_t expiration;
    int32_t last_trade;
    double strike;
    const char *option_type;
    const char *expiration_style;
    const char *settlement_style;
    const char *price_increments;
    const char *trading_hours;
    dxfg_string_list *custom_fields;
} dxfg_instrument_profile_t;

typedef struct dxfg_instrument_profile_list {
    int32_t size;
    dxfg_instrument_profile_t **elements;
} dxfg_instrument_profile_list;

/**
 * @brief The InstrumentProfileUpdateListener.
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileUpdateListener.html">Javadoc</a>
 */
typedef struct dxfg_ipf_update_listener_t {
    dxfg_java_object_handler handler;
} dxfg_ipf_update_listener_t;

typedef struct dxfg_iterable_ip_t {
    dxfg_java_object_handler handler;
} dxfg_iterable_ip_t;

typedef void (*dxfg_ipf_update_listener_function)(graal_isolatethread_t *thread, dxfg_iterable_ip_t *profiles, void *user_data);

dxfg_ipf_connection_t*                        dxfg_InstrumentProfileConnection_createConnection(graal_isolatethread_t *thread, const char *address, dxfg_ipf_collector_t *collector);
const char*                                   dxfg_InstrumentProfileConnection_getAddress(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
int64_t                                       dxfg_InstrumentProfileConnection_getUpdatePeriod(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
int32_t                                       dxfg_InstrumentProfileConnection_setUpdatePeriod(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, int64_t newValue);
dxfg_ipf_connection_state_t                   dxfg_InstrumentProfileConnection_getState(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
int64_t                                       dxfg_InstrumentProfileConnection_getLastModified(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
int32_t                                       dxfg_InstrumentProfileConnection_start(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
int32_t                                       dxfg_InstrumentProfileConnection_close(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection);
dxfg_ipf_connection_state_change_listener_t*  dxfg_IpfPropertyChangeListener_new(graal_isolatethread_t *thread, dxfg_ipf_connection_state_change_listener_func user_func, void *user_data);
int32_t                                       dxfg_InstrumentProfileConnection_addStateChangeListener(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, dxfg_ipf_connection_state_change_listener_t *listener);
int32_t                                       dxfg_InstrumentProfileConnection_removeStateChangeListener(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, dxfg_ipf_connection_state_change_listener_t *listener);
int32_t                                       dxfg_InstrumentProfileConnection_waitUntilCompleted(graal_isolatethread_t *thread, dxfg_ipf_connection_t *connection, int64_t timeoutInMs);

int32_t                                       dxfg_CList_InstrumentProfile_release(graal_isolatethread_t *thread, dxfg_instrument_profile_list *ips);

dxfg_ipf_collector_t*                         dxfg_InstrumentProfileCollector_new(graal_isolatethread_t *thread);
int64_t                                       dxfg_InstrumentProfileCollector_getLastUpdateTime(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector);
int32_t                                       dxfg_InstrumentProfileCollector_updateInstrumentProfile(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_instrument_profile_t *ip);
int32_t                                       dxfg_InstrumentProfileCollector_updateInstrumentProfiles(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_instrument_profile_list *ips, dxfg_java_object_handler *generation);
int32_t                                       dxfg_InstrumentProfileCollector_removeGenerations(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_java_object_handler_list *generations);
dxfg_iterable_ip_t*                           dxfg_InstrumentProfileCollector_view(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector);
dxfg_executor_t*                              dxfg_InstrumentProfileCollector_getExecutor(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector);
int32_t                                       dxfg_InstrumentProfileCollector_setExecutor(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_executor_t *executor);
int32_t                                       dxfg_InstrumentProfileCollector_addUpdateListener(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_ipf_update_listener_t* listener);
int32_t                                       dxfg_InstrumentProfileCollector_removeUpdateListener(graal_isolatethread_t *thread, dxfg_ipf_collector_t *collector, dxfg_ipf_update_listener_t* listener);

int32_t                                       dxfg_Iterable_InstrumentProfile_hasNext(graal_isolatethread_t *thread, dxfg_iterable_ip_t *iterable_ip);
dxfg_instrument_profile_t*                    dxfg_Iterable_InstrumentProfile_next(graal_isolatethread_t *thread, dxfg_iterable_ip_t *iterable_ip);

dxfg_ipf_update_listener_t*                   dxfg_InstrumentProfileUpdateListener_new(graal_isolatethread_t *thread, dxfg_ipf_update_listener_function user_func, void *user_data);

dxfg_instrument_profile_reader_t*             dxfg_InstrumentProfileReader_new(graal_isolatethread_t *thread);
int64_t                                       dxfg_InstrumentProfileReader_getLastModified(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader);
int32_t                                       dxfg_InstrumentProfileReader_wasComplete(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader);
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_readFromFile(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address);
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_readFromFile2(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, const char *address, const char *user, const char *password);
const char*                                   dxfg_InstrumentProfileReader_resolveSourceURL(graal_isolatethread_t *thread, const char *address);
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_read2(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is, const char* address);
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_readCompressed(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is);
dxfg_instrument_profile_list*                 dxfg_InstrumentProfileReader_read(graal_isolatethread_t *thread, dxfg_instrument_profile_reader_t *reader, dxfg_input_stream_t* is);

int32_t                                       dxfg_InstrumentProfile_release(graal_isolatethread_t *thread, dxfg_instrument_profile_t *ip);

/** @} */ // end of InstrumentProfile

#ifdef __cplusplus
}
#endif

#endif // DXFEED_GRAAL_NATIVE_SDK_IPF_H_
