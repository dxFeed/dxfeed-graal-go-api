// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
	"unsafe"
)

type UnderlyingMapper struct{}

func (m UnderlyingMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	u := event.(*underlying.Underlying)
	n := (*C.dxfg_v4_underlying_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_underlying_t{})))
	n.market_event.event_type.clazz = C.DXFG_EVENT_UNDERLYING
	n.market_event.event_symbol = a.cString(*u.EventSymbol())
	n.market_event.event_time = C.int64_t(u.EventTime())
	n.event_flags = C.int32_t(u.EventFlags())
	n.index = C.int64_t(u.Index())
	n.volatility = C.double(u.Volatility())
	n.front_volatility = C.double(u.FrontVolatility())
	n.back_volatility = C.double(u.BackVolatility())
	n.call_volume = C.double(u.CallVolume())
	n.put_volume = C.double(u.PutVolume())
	n.put_call_ratio = C.double(u.PutCallRatio())
	return unsafe.Pointer(n)
}

func (m UnderlyingMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_underlying_t)(native)
	u := underlying.NewUnderlying(C.GoString(n.market_event.event_symbol))
	u.SetEventTime(int64(n.market_event.event_time))
	u.SetEventFlags(int32(n.event_flags))
	u.SetIndex(int64(n.index))
	u.SetVolatility(float64(n.volatility))
	u.SetFrontVolatility(float64(n.front_volatility))
	u.SetBackVolatility(float64(n.back_volatility))
	u.SetCallVolume(float64(n.call_volume))
	u.SetPutVolume(float64(n.put_volume))
	u.SetPutCallRatio(float64(n.put_call_ratio))
	return u
}

// FreeCEvent frees an event made by CEvent.
func (UnderlyingMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_underlying_t)(event)
	cFree(unsafe.Pointer(n.market_event.event_symbol))
	cFree(event)
}
