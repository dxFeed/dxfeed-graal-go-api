package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
	"unsafe"
)

type UnderlyingMapper struct{}

func (m UnderlyingMapper) CEvent(event interface{}) unsafe.Pointer {
	u := event.(*underlying.Underlying)
	n := (*C.dxfg_underlying_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_underlying_t{}))))
	n.market_event.event_type.clazz = C.DXFG_EVENT_UNDERLYING
	n.market_event.event_symbol = C.CString(*u.EventSymbol())
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
	n := (*C.dxfg_underlying_t)(native)
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
