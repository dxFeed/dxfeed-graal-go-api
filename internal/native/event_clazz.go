// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
*/
import "C"

// sdkEventClazzes maps every dxfg_event_clazz_t constant name to its value; a test checks the Go event codes against it.
func sdkEventClazzes() map[string]int32 {
	return map[string]int32{
		"QUOTE":              C.DXFG_EVENT_QUOTE,
		"PROFILE":            C.DXFG_EVENT_PROFILE,
		"SUMMARY":            C.DXFG_EVENT_SUMMARY,
		"GREEKS":             C.DXFG_EVENT_GREEKS,
		"CANDLE":             C.DXFG_EVENT_CANDLE,
		"DAILY_CANDLE":       C.DXFG_EVENT_DAILY_CANDLE,
		"UNDERLYING":         C.DXFG_EVENT_UNDERLYING,
		"THEO_PRICE":         C.DXFG_EVENT_THEO_PRICE,
		"TRADE":              C.DXFG_EVENT_TRADE,
		"TRADE_ETH":          C.DXFG_EVENT_TRADE_ETH,
		"CONFIGURATION":      C.DXFG_EVENT_CONFIGURATION,
		"MESSAGE":            C.DXFG_EVENT_MESSAGE,
		"TIME_AND_SALE":      C.DXFG_EVENT_TIME_AND_SALE,
		"ORDER_BASE":         C.DXFG_EVENT_ORDER_BASE,
		"ORDER":              C.DXFG_EVENT_ORDER,
		"ANALYTIC_ORDER":     C.DXFG_EVENT_ANALYTIC_ORDER,
		"OTC_MARKETS_ORDER":  C.DXFG_EVENT_OTC_MARKETS_ORDER,
		"SPREAD_ORDER":       C.DXFG_EVENT_SPREAD_ORDER,
		"SERIES":             C.DXFG_EVENT_SERIES,
		"OPTION_SALE":        C.DXFG_EVENT_OPTION_SALE,
		"TEXT_MESSAGE":       C.DXFG_EVENT_TEXT_MESSAGE,
		"MARKET_MAKER":       C.DXFG_EVENT_MARKET_MAKER,
		"TEXT_CONFIGURATION": C.DXFG_EVENT_TEXT_CONFIGURATION,
		"NUAM_ORDER":         C.DXFG_EVENT_NUAM_ORDER,
		"NUAM_TIME_AND_SALE": C.DXFG_EVENT_NUAM_TIME_AND_SALE,
		"NUAM_TRADE":         C.DXFG_EVENT_NUAM_TRADE,
		"ORDER_IMBALANCE":    C.DXFG_EVENT_ORDER_IMBALANCE,
	}
}
