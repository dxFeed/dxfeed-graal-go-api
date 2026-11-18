// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The filter of the subscriptions and the events by named categories of QD filter expressions, its builder and
 * the tracker of its updates (see the group @ref V4_FeedFilter). Experimental.
 */

#ifndef DXFG_V4_FEED_FILTER_H
#define DXFG_V4_FEED_FILTER_H

#include "dxfg_v4_common.h"
#include "dxfg_v4_endpoint.h"
#include "dxfg_v4_events.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_FeedFilter Feed filter (v4, experimental)
 *  @brief The filter of the subscriptions and the events by named categories of QD filter expressions.
 *
 *  A builder (dxfg_v4_DXFeedFilter_newBuilder()) takes the categories of a filter, each one a name and a QD filter
 *  expression (e.g. `IBM,AAPL`, `A*`, `range-A-M-`, `ipf[file.ipf,update=10m]`), and builds the filter
 *  (dxfg_v4_DXFeedFilter_Builder_build()). The filter tells whether a subscription (an event type and a symbol) or an
 *  event passes any of its categories, and which category is the first one that it passes, in the order the
 *  categories were added.
 *
 *  A filter does not change. A dynamic filter (dxfg_v4_DXFeedFilter_isDynamic()), e.g. one with an `ipf[...]`
 *  expression with an `update` period, gets a new version on each update: its tracker
 *  (dxfg_v4_DXFeedFilter_getTracker()) returns the latest version (dxfg_v4_DXFeedFilterTracker_getCurrentFilter()) and
 *  calls its listeners (dxfg_v4_DXFeedFilterListener_new()) after each update.
 *
 *  Each function that returns a handle returns a new one, which the caller releases with
 *  dxfg_v4_JavaObjectHandler_release(). The builder, the filter, the tracker and the listener can be released in any
 *  order.
 *
 *  All the functions and types of the group are experimental (see @ref V4_Conventions): the Java API they call is
 *  `@Experimental` in QD and may change. #DXFG_V4_FEATURE_FEED_FILTER reports them.
 *
 *  Java: `com.dxfeed.api.DXFeedFilter`, `DXFeedFilter.Builder`, `DXFeedFilterTracker`, `DXFeedFilterListener`.
 *  @{
 */

/**
 * @brief The handle of a builder of a feed filter (dxfg_v4_DXFeedFilter_newBuilder()).
 * @experimental
 * @since 3.13.0
 */
typedef struct dxfg_v4_feed_filter_builder_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_v4_feed_filter_builder_t;

/**
 * @brief The handle of a feed filter (dxfg_v4_DXFeedFilter_Builder_build()).
 * @experimental
 * @since 3.13.0
 */
typedef struct dxfg_v4_feed_filter_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_v4_feed_filter_t;

/**
 * @brief The handle of the tracker of the updates of a feed filter (dxfg_v4_DXFeedFilter_getTracker()).
 * @experimental
 * @since 3.13.0
 */
typedef struct dxfg_v4_feed_filter_tracker_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_v4_feed_filter_tracker_t;

/**
 * @brief The handle of a listener of the updates of a feed filter (dxfg_v4_DXFeedFilterListener_new()).
 * @experimental
 * @since 3.13.0
 */
typedef struct dxfg_v4_feed_filter_listener_t {
    /** The handle of the Java object. */
    dxfg_java_object_handler handler;
} dxfg_v4_feed_filter_listener_t;

/**
 * @brief The function of a listener of the updates of a feed filter (dxfg_v4_DXFeedFilterListener_new()).
 *
 * It is called after the filter of a tracker to which the listener is added has been updated; the tracker returns the
 * new version with dxfg_v4_DXFeedFilterTracker_getCurrentFilter(). Fast updates may come as one call. It is called on
 * a thread of the SDK that updates the filter (e.g. the thread that reloads the IPF of an `ipf[...]` expression), never
 * on the thread that adds the listener. The function may call any function of the SDK on that thread, including
 * dxfg_v4_DXFeedFilterTracker_addListener() and dxfg_v4_DXFeedFilterTracker_removeListener() of its own tracker and
 * the release of its own listener; it should return quickly, since the next update waits for it.
 *
 * @param[in] thread The current GraalVM Isolate's thread (the thread of the SDK that calls the function).
 * @param[in] user_data The user data given to dxfg_v4_DXFeedFilterListener_new().
 * @experimental
 * @since 3.13.0
 */
typedef void (*dxfg_v4_feed_filter_listener_func)(graal_isolatethread_t *thread, void *user_data);

/**
 * @brief Creates a new builder of a feed filter for the endpoint, without categories.
 *
 * The expressions of the categories are parsed with the data scheme of the endpoint at
 * dxfg_v4_DXFeedFilter_Builder_build(). A builder is not thread-safe: use it on one thread at a time.
 *
 * Java: `DXFeedFilter.newBuilder(DXEndpoint)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] endpoint The endpoint; must not be NULL.
 * @param[out] builder The pointer to the variable that receives the builder (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_newBuilder(graal_isolatethread_t *thread, dxfg_v4_endpoint_t *endpoint,
                                        DXFG_OUT dxfg_v4_feed_filter_builder_t **builder);

/**
 * @brief Adds a category to the builder: a name and its QD filter expression.
 *
 * The categories are evaluated in the order they are added; the first one that matches gives the category. The
 * expression is checked at dxfg_v4_DXFeedFilter_Builder_build(), not here.
 *
 * Java: `DXFeedFilter.Builder.withCategory(String, String)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder; must not be NULL.
 * @param[in] categoryName The name of the category, a non-empty UTF-8 string that the builder does not have yet.
 * @param[in] filterExpression The QD filter expression of the category, a non-empty UTF-8 string.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error: a NULL
 * or empty name or expression, or a name that the builder already has (`IllegalArgumentException`); the builder is
 * unchanged then. Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_Builder_withCategory(graal_isolatethread_t *thread,
                                                  dxfg_v4_feed_filter_builder_t *builder, const char *categoryName,
                                                  const char *filterExpression);

/**
 * @brief Adds categories to the builder: the names and the QD filter expressions of two lists of the same size, in
 * the order of the lists.
 *
 * The lists are checked as a whole before the builder changes: on an error no category is added. Empty lists add
 * nothing. The lists stay the caller's.
 *
 * Java: `DXFeedFilter.Builder.withCategories(Map)` with a `LinkedHashMap`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder; must not be NULL.
 * @param[in] categoryNames The names of the categories: non-empty UTF-8 strings, different from each other and from
 * the names that the builder already has; must not be NULL.
 * @param[in] filterExpressions The QD filter expressions of the categories, one for each name: non-empty UTF-8
 * strings; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error: a NULL
 * list, lists of different sizes, a NULL or empty name or expression, or a repeated name (`IllegalArgumentException`).
 * Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_Builder_withCategories(graal_isolatethread_t *thread,
                                                    dxfg_v4_feed_filter_builder_t *builder,
                                                    const dxfg_v4_string_list *categoryNames,
                                                    const dxfg_v4_string_list *filterExpressions);

/**
 * @brief Builds a filter of the categories of the builder.
 *
 * The builder stays usable: more categories can be added and another filter built.
 *
 * Java: `DXFeedFilter.Builder.build()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] builder The builder; must not be NULL.
 * @param[out] filter The pointer to the variable that receives the filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error: no
 * categories, or an expression that is not a valid QD filter (`IllegalArgumentException`, e.g. QD's
 * `FilterSyntaxException`). Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_Builder_build(graal_isolatethread_t *thread, dxfg_v4_feed_filter_builder_t *builder,
                                           DXFG_OUT dxfg_v4_feed_filter_t **filter);

/**
 * @brief Checks whether a subscription (an event type and a symbol) passes any category of the filter.
 *
 * Java: `DXFeedFilter.accept(Class, Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The filter; must not be NULL.
 * @param[in] eventType The event type of the subscription.
 * @param[in] symbol The subscription symbol (a string, a candle, an indexed event subscription or a time series
 * subscription symbol); must not be NULL.
 * @param[out] accepted The pointer to the variable that receives 1 if the subscription passes, 0 otherwise; must not be
 * NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_accept(graal_isolatethread_t *thread, dxfg_v4_feed_filter_t *filter,
                                    dxfg_v4_event_clazz_t eventType, dxfg_v4_symbol_t *symbol,
                                    DXFG_OUT int32_t *accepted);

/**
 * @brief Checks whether an event passes any category of the filter, by the type and the symbol of the event.
 *
 * Java: `DXFeedFilter.acceptEvent(EventType)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The filter; must not be NULL.
 * @param[in] event The event; must not be NULL.
 * @param[out] accepted The pointer to the variable that receives 1 if the event passes, 0 otherwise; must not be NULL.
 * It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_acceptEvent(graal_isolatethread_t *thread, dxfg_v4_feed_filter_t *filter,
                                         dxfg_v4_event_type_t *event, DXFG_OUT int32_t *accepted);

/**
 * @brief Returns the name of the first category of the filter that a subscription (an event type and a symbol)
 * passes, in the order the categories were added.
 *
 * Java: `DXFeedFilter.getCategory(Class, Object)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The filter; must not be NULL.
 * @param[in] eventType The event type of the subscription.
 * @param[in] symbol The subscription symbol (see dxfg_v4_DXFeedFilter_accept()); must not be NULL.
 * @param[out] category The pointer to the variable that receives the name of the category (release it with
 * dxfg_v4_String_release()), or NULL if the subscription passes no category (it is not an error); must not be NULL.
 * It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_getCategory(graal_isolatethread_t *thread, dxfg_v4_feed_filter_t *filter,
                                         dxfg_v4_event_clazz_t eventType, dxfg_v4_symbol_t *symbol,
                                         DXFG_OUT const char **category);

/**
 * @brief Returns the name of the first category of the filter that an event passes, in the order the categories were
 * added.
 *
 * Java: `DXFeedFilter.getEventCategory(EventType)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The filter; must not be NULL.
 * @param[in] event The event; must not be NULL.
 * @param[out] category The pointer to the variable that receives the name of the category (release it with
 * dxfg_v4_String_release()), or NULL if the event passes no category (it is not an error); must not be NULL. It is
 * written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_getEventCategory(graal_isolatethread_t *thread, dxfg_v4_feed_filter_t *filter,
                                              dxfg_v4_event_type_t *event, DXFG_OUT const char **category);

/**
 * @brief Checks whether the filter is dynamic: whether it gets new versions, which its tracker reports
 * (dxfg_v4_DXFeedFilter_getTracker()).
 *
 * Java: `DXFeedFilter.isDynamic()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The filter; must not be NULL.
 * @param[out] dynamic The pointer to the variable that receives 1 if the filter is dynamic, 0 otherwise; must not be
 * NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_isDynamic(graal_isolatethread_t *thread, dxfg_v4_feed_filter_t *filter,
                                       DXFG_OUT int32_t *dynamic);

/**
 * @brief Returns the tracker of the filter: the holder of its latest version, the same for all the versions.
 *
 * Each call returns a new handle of the same tracker.
 *
 * Java: `DXFeedFilter.getTracker()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] filter The filter; must not be NULL.
 * @param[out] tracker The pointer to the variable that receives the tracker (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilter_getTracker(graal_isolatethread_t *thread, dxfg_v4_feed_filter_t *filter,
                                        DXFG_OUT dxfg_v4_feed_filter_tracker_t **tracker);

/**
 * @brief Returns the latest version of the filter of the tracker.
 *
 * The tracker follows the updates only while it has a listener (dxfg_v4_DXFeedFilterTracker_addListener()); without
 * one it may return an older version. Each call returns a new handle.
 *
 * Java: `DXFeedFilterTracker.getCurrentFilter()`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] tracker The tracker; must not be NULL.
 * @param[out] filter The pointer to the variable that receives the filter (release it with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilterTracker_getCurrentFilter(graal_isolatethread_t *thread,
                                                     dxfg_v4_feed_filter_tracker_t *tracker,
                                                     DXFG_OUT dxfg_v4_feed_filter_t **filter);

/**
 * @brief Adds a listener to the tracker: it is called after each later update of the filter.
 *
 * The tracker of a filter that is not dynamic ignores the listener. The listener stays added until it is removed; a
 * listener whose last handle is released is not called anymore and removes itself at the next update.
 *
 * Java: `DXFeedFilterTracker.addListener(DXFeedFilterListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] tracker The tracker; must not be NULL.
 * @param[in] listener The listener; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilterTracker_addListener(graal_isolatethread_t *thread, dxfg_v4_feed_filter_tracker_t *tracker,
                                                dxfg_v4_feed_filter_listener_t *listener);

/**
 * @brief Removes a listener from the tracker.
 *
 * A call of the listener that has already started on another thread may still run after the function returns;
 * releasing the last handle of the listener waits for it (see dxfg_v4_DXFeedFilterListener_new()).
 *
 * Java: `DXFeedFilterTracker.removeListener(DXFeedFilterListener)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] tracker The tracker; must not be NULL.
 * @param[in] listener The listener; must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilterTracker_removeListener(graal_isolatethread_t *thread,
                                                   dxfg_v4_feed_filter_tracker_t *tracker,
                                                   dxfg_v4_feed_filter_listener_t *listener);

/**
 * @brief Creates a new listener of the updates of a feed filter.
 *
 * The listener does nothing until it is added to a tracker with dxfg_v4_DXFeedFilterTracker_addListener(). The
 * function is called as described at dxfg_v4_feed_filter_listener_func, with the user data only: it gets no handles.
 *
 * Releasing the last handle of the listener with dxfg_v4_JavaObjectHandler_release() (or
 * dxfg_v4_NativeListener_deactivate(), or dxfg_v4_NativeListener_deactivateAll()) stops the calls of the function:
 * after the release returns the function is not running (except on the releasing thread itself, when the function
 * releases its own listener) and is not called anymore, so the user data can be freed. A released listener that is
 * still added removes itself from the tracker at its next update; until then the tracker keeps following the
 * updates (e.g. reloads the IPF file), so remove the listener before the release to stop at once.
 *
 * Java: `DXFeedFilterListener.filterUpdated(DXFeedFilterTracker)`.
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] user_func The function to call after the updates; must not be NULL.
 * @param[in] user_data The user data to pass to the function; the SDK does not use it. It must stay valid until the
 * listener is deactivated.
 * @param[out] listener The pointer to the variable that receives the listener (release its handle with
 * dxfg_v4_JavaObjectHandler_release()); must not be NULL. It is written on success only.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 * @experimental
 * @since 3.13.0
 */
DXFG_EXPERIMENTAL
int32_t dxfg_v4_DXFeedFilterListener_new(graal_isolatethread_t *thread, dxfg_v4_feed_filter_listener_func user_func,
                                         void *user_data, DXFG_OUT dxfg_v4_feed_filter_listener_t **listener);

/** @} */

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_FEED_FILTER_H
