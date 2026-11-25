// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/series"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/theoprice"
)

// SeriesMapper maps the Series event.
type SeriesMapper struct{}

func (SeriesMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*series.Series)
	n := (*C.dxfg_v4_series_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_series_t{})))
	n.market_event.event_type.clazz = C.DXFG_EVENT_SERIES
	n.market_event.event_symbol = a.cString(*e.EventSymbol())
	n.market_event.event_time = C.int64_t(e.EventTime())
	n.event_flags = C.int32_t(e.EventFlags())
	n.index = C.int64_t(e.Index())
	n.time_sequence = C.int64_t(e.TimeSequence())
	n.expiration = C.int32_t(e.Expiration())
	n.volatility = C.double(e.Volatility())
	n.call_volume = C.double(e.CallVolume())
	n.put_volume = C.double(e.PutVolume())
	n.put_call_ratio = C.double(e.PutCallRatio())
	n.forward_price = C.double(e.ForwardPrice())
	n.dividend = C.double(e.Dividend())
	n.interest = C.double(e.Interest())
	return unsafe.Pointer(n)
}

func (SeriesMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_series_t)(native)
	e := series.NewSeries(C.GoString(n.market_event.event_symbol))
	e.SetEventTime(int64(n.market_event.event_time))
	e.SetEventFlags(int32(n.event_flags))
	e.SetIndex(int64(n.index))
	e.SetTimeSequence(int64(n.time_sequence))
	e.SetExpiration(int32(n.expiration))
	e.SetVolatility(float64(n.volatility))
	e.SetCallVolume(float64(n.call_volume))
	e.SetPutVolume(float64(n.put_volume))
	e.SetPutCallRatio(float64(n.put_call_ratio))
	e.SetForwardPrice(float64(n.forward_price))
	e.SetDividend(float64(n.dividend))
	e.SetInterest(float64(n.interest))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (SeriesMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_series_t)(event)
	cFree(unsafe.Pointer(n.market_event.event_symbol))
	cFree(event)
}

// TheoPriceMapper maps the TheoPrice event.
type TheoPriceMapper struct{}

func (TheoPriceMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*theoprice.TheoPrice)
	n := (*C.dxfg_v4_theo_price_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_theo_price_t{})))
	n.market_event.event_type.clazz = C.DXFG_EVENT_THEO_PRICE
	n.market_event.event_symbol = a.cString(*e.EventSymbol())
	n.market_event.event_time = C.int64_t(e.EventTime())
	n.event_flags = C.int32_t(e.EventFlags())
	n.index = C.int64_t(e.Index())
	n.price = C.double(e.Price())
	n.underlying_price = C.double(e.UnderlyingPrice())
	n.delta = C.double(e.Delta())
	n.gamma = C.double(e.Gamma())
	n.dividend = C.double(e.Dividend())
	n.interest = C.double(e.Interest())
	return unsafe.Pointer(n)
}

func (TheoPriceMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_theo_price_t)(native)
	e := theoprice.NewTheoPrice(C.GoString(n.market_event.event_symbol))
	e.SetEventTime(int64(n.market_event.event_time))
	e.SetEventFlags(int32(n.event_flags))
	e.SetIndex(int64(n.index))
	e.SetPrice(float64(n.price))
	e.SetUnderlyingPrice(float64(n.underlying_price))
	e.SetDelta(float64(n.delta))
	e.SetGamma(float64(n.gamma))
	e.SetDividend(float64(n.dividend))
	e.SetInterest(float64(n.interest))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (TheoPriceMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_theo_price_t)(event)
	cFree(unsafe.Pointer(n.market_event.event_symbol))
	cFree(event)
}
