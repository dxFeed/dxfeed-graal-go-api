package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"

func AvailableMappers() map[int32]MapperInterface {
	eventMappers := map[int32]MapperInterface{
		C.DXFG_EVENT_QUOTE:          QuoteMapper{},
		C.DXFG_EVENT_TIME_AND_SALE:  TimeAndSaleMapper{},
		C.DXFG_EVENT_PROFILE:        ProfileMapper{},
		C.DXFG_EVENT_ORDER:          OrderMapper{},
		C.DXFG_EVENT_SPREAD_ORDER:   SpreadOrderMapper{},
		C.DXFG_EVENT_CANDLE:         CandleMapper{},
		C.DXFG_EVENT_TRADE:          TradeMapper{},
		C.DXFG_EVENT_TRADE_ETH:      TradeETHMapper{},
		C.DXFG_EVENT_ANALYTIC_ORDER: AnalyticOrderMapper{},
	}
	return eventMappers
}
