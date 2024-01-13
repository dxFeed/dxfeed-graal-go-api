// SPDX-License-Identifier: MPL-2.0

#ifndef DXFEED_GRAAL_NATIVE_SDK_CATCH_EXCEPTION_H_
#define DXFEED_GRAAL_NATIVE_SDK_CATCH_EXCEPTION_H_

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "graal_isolate.h"


typedef struct dxfg_stack_trace_element_t {
    const char *class_name;
    const char *class_loader_name;
    const char *file_name;
    const char *method_name;
    int32_t line_number;
    int32_t is_native_method;
    const char *module_name;
    const char *module_version;
} dxfg_stack_trace_element_t;

typedef struct dxfg_stack_trace_element_list {
    int32_t size;
    dxfg_stack_trace_element_t **elements;
} dxfg_stack_trace_element_list;

/**
 * @brief The Exception.
 * <a href="https://docs.oracle.com/javase/8/docs/api/java/lang/Exception.html">Javadoc</a>
 */
typedef struct dxfg_exception_t {
    const char *class_name;
    const char *message;
    const char *print_stack_trace;
    dxfg_stack_trace_element_list *stack_trace;
    struct dxfg_exception_t *cause;
} dxfg_exception_t;

dxfg_exception_t* dxfg_get_and_clear_thread_exception_t(graal_isolatethread_t *thread);
void dxfg_Exception_release(graal_isolatethread_t *thread, dxfg_exception_t *exception);

#ifdef __cplusplus
}
#endif

#endif // DXFEED_GRAAL_NATIVE_SDK_CATCH_EXCEPTION_H_
