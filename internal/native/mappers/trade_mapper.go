// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"unsafe"
)

type TradeMapper struct {
}

func (t TradeMapper) GoEvent(native unsafe.Pointer) interface{} {
	tradeNative := (*C.dxfg_v4_trade_t)(native)
	tradeEvent := trade.NewTrade(C.GoString(tradeNative.trade_base.market_event.event_symbol))
	tradeFromC(&tradeEvent.TradeBase, &tradeNative.trade_base)
	return tradeEvent
}

func (t TradeMapper) CEvent(event interface{}) unsafe.Pointer {
	q := (*C.dxfg_v4_trade_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_trade_t{})))
	tradeToC(&q.trade_base, &event.(*trade.Trade).TradeBase, C.DXFG_EVENT_TRADE)
	return unsafe.Pointer(q)
}

// FreeCEvent frees an event made by CEvent.
func (TradeMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_trade_t)(event)
	cFree(unsafe.Pointer(n.trade_base.market_event.event_symbol))
	cFree(event)
}

// tradeToC fills the native trade base with the fields of the Go trade base and sets its clazz.
func tradeToC(q *C.dxfg_v4_trade_base_t, tradeEvent *trade.TradeBase, clazz C.dxfg_event_clazz_t) {
	q.market_event.event_type.clazz = clazz
	q.market_event.event_symbol = cString(*tradeEvent.EventSymbol())
	q.market_event.event_time = C.int64_t(tradeEvent.EventTime())
	q.time_sequence = C.int64_t(tradeEvent.TimeSequence())
	q.time_nano_part = C.int32_t(tradeEvent.TimeNanoPart())
	q.exchange_code = C.int16_t(tradeEvent.ExchangeCode())
	q.price = C.double(tradeEvent.Price())
	q.change = C.double(tradeEvent.Change())
	q.size = C.double(tradeEvent.Size())
	q.day_id = C.int32_t(tradeEvent.DayId())
	q.day_volume = C.double(tradeEvent.DayVolume())
	q.day_turnover = C.double(tradeEvent.DayTurnover())
	q.flags = C.int32_t(tradeEvent.Flags())
}

// tradeFromC fills the Go trade base with the fields of the native trade base, except the event symbol.
func tradeFromC(tradeEvent *trade.TradeBase, tradeNative *C.dxfg_v4_trade_base_t) {
	tradeEvent.SetEventTime(int64(tradeNative.market_event.event_time))
	tradeEvent.SetTimeSequence(int64(tradeNative.time_sequence))
	tradeEvent.SetTimeNanoPart(int32(tradeNative.time_nano_part))
	tradeEvent.SetExchangeCode(int16(tradeNative.exchange_code))
	tradeEvent.SetPrice(float64(tradeNative.price))
	tradeEvent.SetSize(float64(tradeNative.size))
	tradeEvent.SetChange(float64(tradeNative.change))
	tradeEvent.SetDayId(int32(tradeNative.day_id))
	tradeEvent.SetDayVolume(float64(tradeNative.day_volume))
	tradeEvent.SetDayTurnover(float64(tradeNative.day_turnover))
	tradeEvent.SetFlags(int32(tradeNative.flags))
}
