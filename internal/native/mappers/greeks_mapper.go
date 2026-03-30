package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
	"unsafe"
)

type GreeksMapper struct{}

func (m GreeksMapper) CEvent(event interface{}) unsafe.Pointer {
	g := event.(*greeks.Greeks)
	n := (*C.dxfg_greeks_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_greeks_t{}))))
	n.market_event.event_type.clazz = C.DXFG_EVENT_GREEKS
	n.market_event.event_symbol = C.CString(*g.EventSymbol())
	n.market_event.event_time = C.int64_t(g.EventTime())
	n.event_flags = C.int32_t(g.EventFlags())
	n.index = C.int64_t(g.Index())
	n.price = C.double(g.Price())
	n.volatility = C.double(g.Volatility())
	n.delta = C.double(g.Delta())
	n.gamma = C.double(g.Gamma())
	n.theta = C.double(g.Theta())
	n.rho = C.double(g.Rho())
	n.vega = C.double(g.Vega())
	return unsafe.Pointer(n)
}

func (m GreeksMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_greeks_t)(native)
	g := greeks.NewGreeks(C.GoString(n.market_event.event_symbol))
	g.SetEventTime(int64(n.market_event.event_time))
	g.SetEventFlags(int32(n.event_flags))
	g.SetIndex(int64(n.index))
	g.SetPrice(float64(n.price))
	g.SetVolatility(float64(n.volatility))
	g.SetDelta(float64(n.delta))
	g.SetGamma(float64(n.gamma))
	g.SetTheta(float64(n.theta))
	g.SetRho(float64(n.rho))
	g.SetVega(float64(n.vega))
	return g
}
