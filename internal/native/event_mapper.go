package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"unsafe"
)

type eventMapper struct {
	mappers map[int32]mappers.MapperInterface
}

func newEventMapper() *eventMapper {
	eventMappers := make(map[int32]mappers.MapperInterface)
	eventMappers[C.DXFG_EVENT_QUOTE] = mappers.QuoteMapper{}
	eventMappers[C.DXFG_EVENT_TIME_AND_SALE] = mappers.TimeAndSaleMapper{}
	eventMappers[C.DXFG_EVENT_PROFILE] = mappers.ProfileMapper{}
	eventMappers[C.DXFG_EVENT_ORDER] = mappers.OrderMapper{}
	eventMappers[C.DXFG_EVENT_SPREAD_ORDER] = mappers.SpreadOrderMapper{}

	return &eventMapper{mappers: eventMappers}
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
