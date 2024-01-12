// SPDX-License-Identifier: MPL-2.0

#ifndef DXFEED_GRAAL_NATIVE_SDK_SYSTEM_H_
#define DXFEED_GRAAL_NATIVE_SDK_SYSTEM_H_

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "graal_isolate.h"

/** @defgroup System
 *  @{
 */

/**
 * @brief Sets the system property indicated by the specified key.
 */
int32_t dxfg_system_set_property(graal_isolatethread_t *thread, const char *key, const char *value);

/**
 * @brief Gets the system property indicated by the specified key.
 */
const char *dxfg_system_get_property(graal_isolatethread_t *thread, const char *key);

/**
 * @brief Frees pointer that was previously allocated in Java method (by UnmanagedMemory).
 */
int32_t dxfg_system_release_property(graal_isolatethread_t *thread, const char *value);

/** @} */ // end of System

#ifdef __cplusplus
}
#endif

#endif // DXFEED_GRAAL_NATIVE_SDK_SYSTEM_H_
