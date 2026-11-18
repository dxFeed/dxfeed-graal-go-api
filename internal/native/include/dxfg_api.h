// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The header to include for the whole C API of the SDK: it includes every header of the API, graal_isolate.h
 * among them.
 *
 * A program that loads the library at run time includes graal_isolate_dynamic.h before it (see there).
 */

#ifndef DXFG_API_H
#define DXFG_API_H

#ifdef __cplusplus
extern "C" {
#endif

#include "dxfg_common.h"

#include "dxfg_catch_exception.h"
#include "dxfg_endpoint.h"
#include "dxfg_events.h"
#include "dxfg_feed.h"
#include "dxfg_publisher.h"
#include "dxfg_subscription_controller.h"
#include "dxfg_subscription.h"
#include "dxfg_system.h"
#include "dxfg_javac.h"
#include "dxfg_qds.h"
#include "graal_isolate.h"
#include "graal_isolate_dynamic.h"
#include "dxfg_event_model.h"
#include "dxfg_ipf.h"
#include "dxfg_ondemand.h"
#include "dxfg_orcs.h"
#include "dxfg_schedule.h"
#include "dxfg_glossary.h"
#include "dxfg_logging.h"
#include "dxfg_candlewebservice.h"
#include "dxfg_indexed_event_tx_model.h"

#include "dxfg_v4_api.h"

#ifdef __cplusplus
}
#endif

#endif // DXFG_API_H
