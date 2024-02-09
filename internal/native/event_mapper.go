package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/Osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"unsafe"
)

type eventMapper struct {
	mappers map[int32]mappers.MapperInterface
}

func newEventMapper() *eventMapper {
	return &eventMapper{mappers: mappers.AvailableMappers()}
}

func (m *eventMapper) goEvents(eventsList *C.dxfg_event_type_list) []interface{} {
	if eventsList == nil || eventsList.elements == nil || int(eventsList.size) == 0 {
		return nil
	}

	size := int(eventsList.size)
	list := make([]interface{}, size)
	elementsSlice := unsafe.Slice(eventsList.elements, C.size_t(eventsList.size))

	for i, event := range elementsSlice {
		list[i] = m.goEvent(event)
	}

	return list
}

func (m *eventMapper) goEvent(event *C.dxfg_event_type_t) interface{} {
	mapper := m.mappers[int32(event.clazz)]
	if mapper != nil {
		return mapper.GoEvent(unsafe.Pointer(event))
	} else {
		panic(fmt.Sprintf("unknown event eventcodes %v", event.clazz))
	}
}

// TODO add recursive release for symbols
func (m *eventMapper) cSymbol(symbol any) unsafe.Pointer {
	switch value := symbol.(type) {
	case string:
		return unsafe.Pointer(m.cStringSymbol(value))
	case *Osub.WildcardSymbol:
		return unsafe.Pointer(m.cWildCardSymbol())
	case *Osub.IndexedEventSubscriptionSymbol:
		return unsafe.Pointer(m.cIndexedEventSubscriptionSymbol(value.Symbol(), value.Source()))
	case *Osub.TimeSeriesSubscriptionSymbol:
		return unsafe.Pointer(m.cTimeSeriesSymbol(value.Symbol(), value.FromTime()))
	default:
		return nil
	}
}

func (m *eventMapper) cStringSymbol(str string) *dxfg_symbol_t {
	ss := &dxfg_symbol_t{}
	ss.t = 0
	ss.symbol = C.CString(str)
	return ss
}

func (m *eventMapper) cWildCardSymbol() *dxfg_symbol_t {
	ss := &dxfg_symbol_t{}
	ss.t = 2
	return ss
}

func (m *eventMapper) cTimeSeriesSymbol(str any, fromTime int64) *dxfg_time_series_subscription_symbol_t {
	ss := &dxfg_time_series_subscription_symbol_t{}
	ss.t = 4
	ss.symbol = (*dxfg_symbol_t)(m.cSymbol(str))
	ss.from_time = C.int64_t(fromTime)
	return ss
}

func (m *eventMapper) cIndexedEventSubscriptionSymbol(str any, source events.IndexedEventSourceInterface) *dxfg_indexed_event_subscription_symbol_t {
	ss := &dxfg_indexed_event_subscription_symbol_t{}
	ss.t = 3
	ss.symbol = (*dxfg_symbol_t)(m.cSymbol(str))
	nativeSource := &dxfg_indexed_event_source_t{}
	nativeSource.id = C.int32_t(source.Id())
	nativeSource.name = C.CString(*source.Name())
	switch value := source.(type) {
	case *events.IndexedEventSource:
		nativeSource.t = 0
	case *order.OrderSource:
		nativeSource.t = 1
	default:
		panic(fmt.Sprintf("Undefined source %T", value))
	}
	ss.source = nativeSource
	return ss
}
