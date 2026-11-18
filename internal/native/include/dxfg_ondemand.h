// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK OnDemand functions and types declarations
 */

#ifndef DXFG_ONDEMAND_H
#define DXFG_ONDEMAND_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"
#include "dxfg_endpoint.h"
#include "dxfg_javac.h"

/** @defgroup OnDemand OnDemand
 *  @brief The control of the on-demand historical tick data replay of an endpoint: start, pause, change the speed and
 *  stop the replay, and query its state.
 *
 *  The service is used to switch a feed from the ordinary real-time or delayed data to the replay of the tick-by-tick
 *  history of the market data without any change in the code that consumes the events. The replay works only when
 *  the endpoint is connected to an on-demand data provider (an address with the "(ondemand:<address>)" component, see
 *  dxfg_DXEndpoint_connect()) and the appropriate credentials are provided.
 *
 *  A handle that a function returns is owned by the caller: release it with dxfg_JavaObjectHandler_release() (see the
 *  group Javac for the ownership rules). A function that returns `int32_t` status returns #DXFG_EXECUTE_SUCCESSFULLY
 *  (0) or #DXFG_EXECUTE_FAIL (-1) on error, a function that returns a pointer returns NULL on error, and a function
 *  that returns a value gives the value that is stated in its `@return` on error. The exception is available from
 *  dxfg_get_and_clear_thread_exception_t().
 *  @{
 */

/**
 * @brief Forward declarations.
 */
typedef struct dxfg_java_object_handler dxfg_java_object_handler;

/**
 * @brief A handle of a Java OnDemandService object; release it with dxfg_JavaObjectHandler_release().
 *
 * One service is associated with each endpoint that has the role #DXFG_ENDPOINT_ROLE_FEED or
 * #DXFG_ENDPOINT_ROLE_ON_DEMAND_FEED: get it with dxfg_OnDemandService_getInstance() or
 * dxfg_OnDemandService_getInstance2(). Releasing the handle does not stop the replay and does not close the endpoint.
 *
 * <a href="https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ondemand/OnDemandService.html">Javadoc</a>
 */
typedef struct dxfg_on_demand_service_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_on_demand_service_t;

/**
 * @brief Returns the on-demand service of the default endpoint with the #DXFG_ENDPOINT_ROLE_ON_DEMAND_FEED role.
 *
 * The endpoint is the application-wide default endpoint of that role (dxfg_OnDemandService_getEndpoint() returns it;
 * dxfg_DXEndpoint_getInstance2() creates a new endpoint instead), which is not connected to any other real-time or
 * delayed data feed. It is created on the first call, is configured by the system properties and is closed by
 * dxfg_system_close_all_and_await_termination(). To switch another feed to the replay and back, use
 * dxfg_OnDemandService_getInstance2() with the endpoint of that feed.
 *
 * The call fails if the on-demand implementation is not available in the library.
 *
 * Java: `OnDemandService.getInstance()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_getInstance().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @return The service (release the handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_getInstance()")
dxfg_on_demand_service_t*         dxfg_OnDemandService_getInstance(graal_isolatethread_t *thread);

/**
 * @brief Returns the on-demand service of the specified endpoint.
 *
 * Each endpoint is associated with one service, so the calls for the same endpoint return handles of the same
 * service. The endpoint must have the role #DXFG_ENDPOINT_ROLE_FEED or #DXFG_ENDPOINT_ROLE_ON_DEMAND_FEED, otherwise
 * the call fails. It also fails if the on-demand implementation is not available in the library.
 *
 * Java: `OnDemandService.getInstance(DXEndpoint)`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_getInstanceByEndpoint().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint; the caller keeps its ownership.
 * @return The service (release the handle with dxfg_JavaObjectHandler_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_getInstanceByEndpoint()")
dxfg_on_demand_service_t*         dxfg_OnDemandService_getInstance2(graal_isolatethread_t *thread, dxfg_endpoint_t *endpoint);

/**
 * @brief Returns the endpoint that is associated with the service.
 *
 * Java: `OnDemandService.getEndpoint()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_getEndpoint().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return The endpoint (release the handle with dxfg_JavaObjectHandler_release(); this does not close the endpoint),
 * or NULL on error. Use dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_getEndpoint()")
dxfg_endpoint_t*                  dxfg_OnDemandService_getEndpoint(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Tells whether the on-demand historical data replay mode is supported.
 *
 * It is supported when the endpoint is connected to an address with the "(ondemand:<address>)" component that
 * specifies an on-demand historical data provider. While it is not supported, the replay functions fail.
 *
 * Java: `OnDemandService.isReplaySupported()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_isReplaySupported().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return 1 if the replay mode is supported, 0 if not, or -1 on error. Use dxfg_get_and_clear_thread_exception_t() to
 * get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_isReplaySupported()")
int32_t                           dxfg_OnDemandService_isReplaySupported(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Tells whether the service is in the replay mode.
 *
 * The replay mode is in effect after dxfg_OnDemandService_replay() or dxfg_OnDemandService_replay2() and before
 * dxfg_OnDemandService_stopAndResume() or dxfg_OnDemandService_stopAndClear().
 *
 * Java: `OnDemandService.isReplay()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_isReplay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return 1 if the service is in the replay mode, 0 if not, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_isReplay()")
int32_t                           dxfg_OnDemandService_isReplay(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Tells whether the service is in the clear mode.
 *
 * The clear mode is in effect after dxfg_OnDemandService_stopAndClear() and before
 * dxfg_OnDemandService_stopAndResume(), dxfg_OnDemandService_replay() or dxfg_OnDemandService_replay2().
 *
 * Java: `OnDemandService.isClear()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_isClear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return 1 if the service is in the clear mode, 0 if not, or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_isClear()")
int32_t                           dxfg_OnDemandService_isClear(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Returns the current or the last replay time.
 *
 * In the replay mode it is the time that is being replayed, otherwise the last time that was replayed. If the replay
 * was never started, it is 0 (the Unix epoch).
 *
 * Java: `OnDemandService.getTime()`, the result of `Date.getTime()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_getTime().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return The time in milliseconds since the Unix epoch (UTC), or -1 on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_getTime()")
int64_t                           dxfg_OnDemandService_getTime(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Returns the replay speed.
 *
 * The speed is measured with respect to the real-time playback speed: 1 is the real-time speed. It is 0 while the
 * service is not in the replay mode.
 *
 * Java: `OnDemandService.getSpeed()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_getSpeed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return The replay speed, or negative infinity (-Inf) on error. Use dxfg_get_and_clear_thread_exception_t() to get
 * the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_getSpeed()")
double                            dxfg_OnDemandService_getSpeed(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Turns the replay mode on from the specified time with the real-time speed (1).
 *
 * The call is the same as dxfg_OnDemandService_replay2() with the speed 1. It fails if the replay mode is not
 * supported (see dxfg_OnDemandService_isReplaySupported()).
 *
 * Java: `OnDemandService.replay(Date)`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_replay().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @param[in] time The time to start the replay from, in milliseconds since the Unix epoch (UTC).
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_replay()")
int32_t                           dxfg_OnDemandService_replay(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service, int64_t time);

/**
 * @brief Turns the replay mode on from the specified time with the specified speed.
 *
 * The feed is disconnected from the other data providers, pre-buffers the historical tick data and replays it as if it
 * came from the data feeds now. After the call dxfg_OnDemandService_isReplay() returns 1,
 * dxfg_OnDemandService_isClear() returns 0, dxfg_OnDemandService_getTime() returns the time that is being replayed and
 * dxfg_OnDemandService_getSpeed() returns `speed`. The call fails if the replay mode is not supported (see
 * dxfg_OnDemandService_isReplaySupported()) or if the speed is negative.
 *
 * Java: `OnDemandService.replay(Date, double)`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_replayWithSpeed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @param[in] time The time to start the replay from, in milliseconds since the Unix epoch (UTC).
 * @param[in] speed The speed to start the replay with: 1 is the real-time speed, more than 1 is faster, less than 1 is
 * slower and 0 is a pause.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_replayWithSpeed()")
int32_t                           dxfg_OnDemandService_replay2(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service, int64_t time, double speed);

/**
 * @brief Pauses the replay and keeps the data snapshot.
 *
 * It is the same as dxfg_OnDemandService_setSpeed() with the speed 0: the current replay time is captured and the
 * replay stops at it. Outside the replay mode the call has no effect (and succeeds).
 *
 * Java: `OnDemandService.pause()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_pause().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_pause()")
int32_t                           dxfg_OnDemandService_pause(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Stops the replay and resumes the ordinary data feed.
 *
 * The call has no effect if the service is neither in the replay mode nor in the clear mode. After the call
 * dxfg_OnDemandService_isReplay() and dxfg_OnDemandService_isClear() return 0, dxfg_OnDemandService_getTime() returns
 * the last replayed time and dxfg_OnDemandService_getSpeed() returns 0. An endpoint with the role
 * #DXFG_ENDPOINT_ROLE_ON_DEMAND_FEED has no ordinary feed, so for it the call works like
 * dxfg_OnDemandService_stopAndClear().
 *
 * Java: `OnDemandService.stopAndResume()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_stopAndResume().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_stopAndResume()")
int32_t                           dxfg_OnDemandService_stopAndResume(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Stops the incoming data and clears it without resuming the data updates.
 *
 * The call works both in the replay mode and for the ordinary data feed: all the incoming data updates are suspended
 * and the current data is cleared. After the call dxfg_OnDemandService_isReplay() returns 0,
 * dxfg_OnDemandService_isClear() returns 1, dxfg_OnDemandService_getTime() returns the last replayed time and
 * dxfg_OnDemandService_getSpeed() returns 0. The ordinary feed is resumed by dxfg_OnDemandService_stopAndResume(),
 * the replay is continued by dxfg_OnDemandService_replay() or dxfg_OnDemandService_replay2().
 *
 * Java: `OnDemandService.stopAndClear()`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_stopAndClear().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_stopAndClear()")
int32_t                           dxfg_OnDemandService_stopAndClear(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service);

/**
 * @brief Changes the replay speed and continues the replay at the current time.
 *
 * The speed is measured with respect to the real-time playback speed. The call fails if the speed is negative, or if
 * the service is not in the replay mode and the speed is not 0.
 *
 * Java: `OnDemandService.setSpeed(double)`.
 *
 * @deprecated Use dxfg_v4_OnDemandService_setSpeed().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] service The service.
 * @param[in] speed The new replay speed: 1 is the real-time speed, more than 1 is faster, less than 1 is slower and 0
 * is a pause.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_get_and_clear_thread_exception_t() to determine if an exception was thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_OnDemandService_setSpeed()")
int32_t                           dxfg_OnDemandService_setSpeed(graal_isolatethread_t *thread, dxfg_on_demand_service_t *service, double speed);

/** @} */ // end of OnDemand

#ifdef __cplusplus
}
#endif

#endif // DXFG_ONDEMAND_H
