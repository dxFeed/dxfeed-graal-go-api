package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"unsafe"
)

type OrderMapper struct {
}

func (m OrderMapper) CEvent(event interface{}) unsafe.Pointer {
	orderEvent := event.(*order.Order)

	q := (*C.dxfg_order_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_order_t{}))))
	q.order_base.market_event.event_type.clazz = C.DXFG_EVENT_ORDER

	q.order_base.market_event.event_symbol = C.CString(*orderEvent.EventSymbol())
	q.order_base.market_event.event_time = C.int64_t(orderEvent.EventTime())

	q.order_base.event_flags = C.int32_t(orderEvent.EventFlags())
	q.order_base.index = C.int64_t(orderEvent.Index())
	q.order_base.time_sequence = C.int64_t(orderEvent.TimeSequence())
	q.order_base.time_nano_part = C.int32_t(orderEvent.TimeNanoPart())
	q.order_base.action_time = C.int64_t(orderEvent.ActionTime())
	q.order_base.order_id = C.int64_t(orderEvent.OrderId())
	q.order_base.aux_order_id = C.int64_t(orderEvent.AuxOrderId())
	q.order_base.price = C.double(orderEvent.Price())
	q.order_base.size = C.double(orderEvent.Size())
	q.order_base.executed_size = C.double(orderEvent.ExecutedSize())
	q.order_base.count = C.int64_t(orderEvent.Count())
	q.order_base.flags = C.int32_t(orderEvent.Flags())
	q.order_base.trade_id = C.int64_t(orderEvent.TradeId())
	q.order_base.trade_price = C.double(orderEvent.TradePrice())
	q.order_base.trade_size = C.double(orderEvent.TradeSize())

	if orderEvent.MarketMaker() != nil {
		q.market_maker = C.CString(*orderEvent.MarketMaker())
	}

	return unsafe.Pointer(q)
}

func (m OrderMapper) GoEvent(native unsafe.Pointer) interface{} {
	orderNative := (*C.dxfg_order_t)(native)
	o := order.NewOrder(C.GoString(orderNative.order_base.market_event.event_symbol))
	o.SetEventTime(int64(orderNative.order_base.market_event.event_time))

	o.SetEventFlags(int32(orderNative.order_base.event_flags))
	_ = o.SetIndex(int64(orderNative.order_base.index))
	o.SetTimeSequence(int64(orderNative.order_base.time_sequence))
	o.SetTimeNanoPart(int32(orderNative.order_base.time_nano_part))
	o.SetActionTime(int64(orderNative.order_base.action_time))
	o.SetOrderId(int64(orderNative.order_base.order_id))
	o.SetAuxOrderId(int64(orderNative.order_base.aux_order_id))
	o.SetPrice(float64(orderNative.order_base.price))
	o.SetSize(float64(orderNative.order_base.size))
	o.SetExecutedSize(float64(orderNative.order_base.executed_size))
	o.SetCount(int64(orderNative.order_base.count))
	o.SetFlags(int32(orderNative.order_base.flags))
	o.SetTradeId(int64(orderNative.order_base.trade_id))
	o.SetTradePrice(float64(orderNative.order_base.trade_price))
	o.SetTradeSize(float64(orderNative.order_base.trade_size))

	o.SetMarketMaker(convertString(orderNative.market_maker))

	return o
}
