package mappers

/*
#include "../graal/dxfg_api.h"
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
	return C.CString(*str)
}

func (ts TimeAndSaleMapper) CEvent(event interface{}) unsafe.Pointer {
	timeAndSale := event.(*timeandsale.TimeAndSale)
	t := (*C.dxfg_time_and_sale_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_time_and_sale_t{}))))
	t.market_event.event_type.clazz = C.DXFG_EVENT_TIME_AND_SALE
	t.market_event.event_symbol = C.CString(*timeAndSale.EventSymbol())
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
	return unsafe.Pointer(t)
}

func (ts TimeAndSaleMapper) GoEvent(native unsafe.Pointer) interface{} {
	timeAndSale1 := (*C.dxfg_time_and_sale_t)(native)
	t := timeandsale.NewTimeAndSale(C.GoString(timeAndSale1.market_event.event_symbol))
	t.SetEventTime(int64(timeAndSale1.market_event.event_time))
	t.SetTimeNanoPart(int32(timeAndSale1.time_nano_part))
	t.SetExchangeCode(int16(timeAndSale1.exchange_code))
	t.SetPrice(float64(timeAndSale1.price))
	t.SetSize(float64(timeAndSale1.size))
	t.SetBidPrice(float64(timeAndSale1.bid_price))
	t.SetAskPrice(float64(timeAndSale1.ask_price))
	t.SetExchangeSaleConditions(convertString(timeAndSale1.exchange_sale_conditions))
	t.SetBuyer(convertString(timeAndSale1.buyer))
	t.SetSeller(convertString(timeAndSale1.seller))
	t.SetEventFlags(int32(timeAndSale1.event_flags))
	t.SetIndex(int64(timeAndSale1.index))
	t.SetFlags(int32(timeAndSale1.flags))
	return t
}
