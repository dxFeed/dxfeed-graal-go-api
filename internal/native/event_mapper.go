// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"unsafe"
)

type eventMapper struct {
	mappers map[int32]mappers.MapperInterface
}

func newEventMapper() *eventMapper {
	return &eventMapper{mappers: mappers.AvailableMappers()}
}

func (m *eventMapper) goEvents(eventsList *C.dxfg_v4_event_type_list) []interface{} {
	if eventsList == nil || eventsList.elements == nil || int(eventsList.size) == 0 {
		return nil
	}

	list := make([]interface{}, 0, int(eventsList.size))
	elementsSlice := unsafe.Slice(eventsList.elements, C.size_t(eventsList.size))

	for _, event := range elementsSlice {
		if goEvent := m.goEvent(event); goEvent != nil {
			list = append(list, goEvent)
		}
	}

	return list
}

// goEvent returns nil for an event class that has no Go mapper yet; such events are skipped and their class is
// reported once.
func (m *eventMapper) goEvent(event *C.dxfg_v4_event_type_t) interface{} {
	if event == nil {
		return nil
	}
	mapper := m.mappers[int32(event.clazz)]
	if mapper == nil {
		reportUnmappedClazz(int32(event.clazz))
		return nil
	}
	return mapper.GoEvent(unsafe.Pointer(event))
}
