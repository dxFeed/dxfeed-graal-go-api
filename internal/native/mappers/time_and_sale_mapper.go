// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"unsafe"
)

type TimeAndSaleMapper struct {
}

func convertString(value *C.char) *string {
	if value == nil {
		return nil
	} else {
		result := C.GoString(value)
		return &result
	}
}

func CString(str *string) *C.char {
	if str == nil {
		return nil
	}
	return cString(*str)
}

func (ts TimeAndSaleMapper) CEvent(event interface{}) unsafe.Pointer {
	t := (*C.dxfg_v4_time_and_sale_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_time_and_sale_t{})))
	timeAndSaleToC(t, event.(*timeandsale.TimeAndSale), C.DXFG_EVENT_TIME_AND_SALE)
	return unsafe.Pointer(t)
}

func (ts TimeAndSaleMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_time_and_sale_t)(native)
	t := timeandsale.NewTimeAndSale(C.GoString(n.market_event.event_symbol))
	timeAndSaleFromC(t, n)
	return t
}

// FreeCEvent frees an event made by CEvent.
func (TimeAndSaleMapper) FreeCEvent(event unsafe.Pointer) {
	freeCTimeAndSale((*C.dxfg_v4_time_and_sale_t)(event))
	cFree(event)
}

// timeAndSaleToC fills the native time and sale with the fields of the Go event and sets its clazz.
func timeAndSaleToC(t *C.dxfg_v4_time_and_sale_t, timeAndSale *timeandsale.TimeAndSale, clazz C.dxfg_event_clazz_t) {
	t.market_event.event_type.clazz = clazz
	t.market_event.event_symbol = cString(*timeAndSale.EventSymbol())
	t.market_event.event_time = C.int64_t(timeAndSale.EventTime())
	t.time_nano_part = C.int32_t(timeAndSale.TimeNanoPart())
	t.exchange_code = C.int16_t(timeAndSale.ExchangeCode())
	t.price = C.double(timeAndSale.Price())
	t.size = C.double(timeAndSale.Size())
	t.bid_price = C.double(timeAndSale.BidPrice())
	t.ask_price = C.double(timeAndSale.AskPrice())
	t.exchange_sale_conditions = CString(timeAndSale.ExchangeSaleConditions())
	t.buyer = CString(timeAndSale.Buyer())
	t.seller = CString(timeAndSale.Seller())
	t.event_flags = C.int32_t(timeAndSale.EventFlags())
	t.index = C.int64_t(timeAndSale.Index())
	t.flags = C.int32_t(timeAndSale.Flags())
}

// timeAndSaleFromC fills the Go event with the fields of the native time and sale, except the event symbol.
func timeAndSaleFromC(t *timeandsale.TimeAndSale, n *C.dxfg_v4_time_and_sale_t) {
	t.SetEventTime(int64(n.market_event.event_time))
	t.SetTimeNanoPart(int32(n.time_nano_part))
	t.SetExchangeCode(int16(n.exchange_code))
	t.SetPrice(float64(n.price))
	t.SetSize(float64(n.size))
	t.SetBidPrice(float64(n.bid_price))
	t.SetAskPrice(float64(n.ask_price))
	t.SetExchangeSaleConditions(convertString(n.exchange_sale_conditions))
	t.SetBuyer(convertString(n.buyer))
	t.SetSeller(convertString(n.seller))
	t.SetEventFlags(int32(n.event_flags))
	t.SetIndex(int64(n.index))
	t.SetFlags(int32(n.flags))
}

// freeCTimeAndSale frees the strings of a native time and sale made by timeAndSaleToC, but not the event itself.
func freeCTimeAndSale(n *C.dxfg_v4_time_and_sale_t) {
	cFree(unsafe.Pointer(n.market_event.event_symbol))
	cFree(unsafe.Pointer(n.exchange_sale_conditions))
	cFree(unsafe.Pointer(n.buyer))
	cFree(unsafe.Pointer(n.seller))
}
