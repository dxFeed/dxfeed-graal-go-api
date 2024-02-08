package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"unsafe"
)

type TradeETHMapper struct {
}

func (t TradeETHMapper) GoEvent(native unsafe.Pointer) interface{} {
	tradeNative := (*C.dxfg_trade_eth_t)(native)
	tradeEvent := trade.NewTradeETH(C.GoString(tradeNative.trade_base.market_event.event_symbol))
	tradeEvent.SetEventTime(int64(tradeNative.trade_base.market_event.event_time))
	tradeEvent.SetTimeSequence(int64(tradeNative.trade_base.time_sequence))
	tradeEvent.SetTimeNanoPart(int32(tradeNative.trade_base.time_nano_part))
	tradeEvent.SetExchangeCode(int16(tradeNative.trade_base.exchange_code))
	tradeEvent.SetPrice(float64(tradeNative.trade_base.price))
	tradeEvent.SetSize(float64(tradeNative.trade_base.size))
	tradeEvent.SetChange(float64(tradeNative.trade_base.change))
	tradeEvent.SetDayId(int32(tradeNative.trade_base.day_id))
	tradeEvent.SetDayVolume(float64(tradeNative.trade_base.day_volume))
	tradeEvent.SetDayTurnover(float64(tradeNative.trade_base.day_turnover))
	tradeEvent.SetFlags(int32(tradeNative.trade_base.flags))
	return tradeEvent
}

func (t TradeETHMapper) CEvent(event interface{}) unsafe.Pointer {
	tradeEvent := event.(*trade.TradeETH)
	q := (*C.dxfg_trade_eth_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_trade_eth_t{}))))
	q.trade_base.market_event.event_type.clazz = C.DXFG_EVENT_TRADE_ETH
	q.trade_base.market_event.event_symbol = C.CString(*tradeEvent.EventSymbol())
	q.trade_base.market_event.event_time = C.int64_t(tradeEvent.EventTime())
	q.trade_base.time_sequence = C.int64_t(tradeEvent.TimeSequence())
	q.trade_base.time_nano_part = C.int32_t(tradeEvent.TimeNanoPart())
	q.trade_base.exchange_code = C.int16_t(tradeEvent.ExchangeCode())
	q.trade_base.price = C.double(tradeEvent.Price())
	q.trade_base.change = C.double(tradeEvent.Change())
	q.trade_base.size = C.double(tradeEvent.Size())
	q.trade_base.day_id = C.int32_t(tradeEvent.DayId())
	q.trade_base.day_volume = C.double(tradeEvent.DayVolume())
	q.trade_base.day_turnover = C.double(tradeEvent.DayTurnover())
	q.trade_base.flags = C.int32_t(tradeEvent.Flags())
	return unsafe.Pointer(q)
}
