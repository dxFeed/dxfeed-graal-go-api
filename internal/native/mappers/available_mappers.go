// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

/*
#cgo CFLAGS: -I${SRCDIR}/../include -Werror=deprecated-declarations
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

func AvailableMappers() map[int32]MapperInterface {
	eventMappers := map[int32]MapperInterface{
		C.DXFG_EVENT_QUOTE:          QuoteMapper{},
		C.DXFG_EVENT_SUMMARY:        SummaryMapper{},
		C.DXFG_EVENT_GREEKS:         GreeksMapper{},
		C.DXFG_EVENT_UNDERLYING:     UnderlyingMapper{},
		C.DXFG_EVENT_TIME_AND_SALE:  TimeAndSaleMapper{},
		C.DXFG_EVENT_PROFILE:        ProfileMapper{},
		C.DXFG_EVENT_ORDER:          OrderMapper{},
		C.DXFG_EVENT_SPREAD_ORDER:   SpreadOrderMapper{},
		C.DXFG_EVENT_CANDLE:         CandleMapper{},
		C.DXFG_EVENT_TRADE:          TradeMapper{},
		C.DXFG_EVENT_TRADE_ETH:      TradeETHMapper{},
		C.DXFG_EVENT_ANALYTIC_ORDER: AnalyticOrderMapper{},
		C.DXFG_EVENT_DAILY_CANDLE:   DailyCandleMapper{},
		C.DXFG_EVENT_SERIES:         SeriesMapper{},
		C.DXFG_EVENT_THEO_PRICE:     TheoPriceMapper{},
		C.DXFG_EVENT_MESSAGE:        MessageMapper{},
		C.DXFG_EVENT_TEXT_MESSAGE:   TextMessageMapper{},
		C.DXFG_EVENT_CONFIGURATION:  ConfigurationMapper{},

		C.DXFG_EVENT_TEXT_CONFIGURATION: TextConfigurationMapper{},
		C.DXFG_EVENT_OPTION_SALE:        OptionSaleMapper{},
		C.DXFG_EVENT_MARKET_MAKER:       MarketMakerMapper{},
		C.DXFG_EVENT_OTC_MARKETS_ORDER:  OtcMarketsOrderMapper{},
		C.DXFG_EVENT_ORDER_IMBALANCE:    OrderImbalanceMapper{},
		C.DXFG_EVENT_NUAM_ORDER:         NuamOrderMapper{},
		C.DXFG_EVENT_NUAM_TIME_AND_SALE: NuamTimeAndSaleMapper{},
		C.DXFG_EVENT_NUAM_TRADE:         NuamTradeMapper{},
	}
	return eventMappers
}
