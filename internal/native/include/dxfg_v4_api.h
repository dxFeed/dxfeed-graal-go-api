// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The header to include for the v4 functions of the C API (see @ref V4_Conventions): it includes every
 * `dxfg_v4_<area>.h` header, and through them the older headers with the types. dxfg_api.h includes it.
 */

#ifndef DXFG_V4_API_H
#define DXFG_V4_API_H

#include "dxfg_v4_common.h"

#include "dxfg_v4_candlewebservice.h"
#include "dxfg_v4_catch_exception.h"
#include "dxfg_v4_endpoint.h"
#include "dxfg_v4_event_model.h"
#include "dxfg_v4_event_queue.h"
#include "dxfg_v4_events.h"
#include "dxfg_v4_feed.h"
#include "dxfg_v4_feed_filter.h"
#include "dxfg_v4_glossary.h"
#include "dxfg_v4_indexed_event_tx_model.h"
#include "dxfg_v4_ipf.h"
#include "dxfg_v4_javac.h"
#include "dxfg_v4_logging.h"
#include "dxfg_v4_message_connector.h"
#include "dxfg_v4_ondemand.h"
#include "dxfg_v4_orcs.h"
#include "dxfg_v4_publisher.h"
#include "dxfg_v4_qds.h"
#include "dxfg_v4_schedule.h"
#include "dxfg_v4_subscription.h"
#include "dxfg_v4_subscription_controller.h"
#include "dxfg_v4_system.h"

#endif // DXFG_V4_API_H
