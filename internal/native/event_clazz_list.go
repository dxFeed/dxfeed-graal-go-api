// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type eventClazzList struct {
	size     C.int32_t
	elements **C.dxfg_v4_event_clazz_t
}

// createEventClazzList builds the list in C memory. cgo's C.malloc never returns NULL: it ends the process when C is
// out of memory.
func createEventClazzList(eventTypes ...int32) *eventClazzList {
	l := &eventClazzList{
		size: C.int32_t(len(eventTypes)),
	}

	if l.size == 0 {
		return l
	}

	l.elements = (**C.dxfg_v4_event_clazz_t)(C.malloc(C.size_t(l.size) * C.size_t(unsafe.Sizeof((*C.dxfg_v4_event_clazz_t)(nil)))))

	// Convert the C pointer to a Go slice for easier indexing
	slice := unsafe.Slice(l.elements, C.size_t(l.size))

	for i, eventType := range eventTypes {
		slice[i] = (*C.dxfg_v4_event_clazz_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_v4_event_clazz_t(0)))))

		*slice[i] = C.dxfg_v4_event_clazz_t(eventType)
	}

	return l
}

func destroyEventClazzList(l *eventClazzList) {
	if l == nil || l.size <= 0 || l.elements == nil {
		return
	}

	// Convert the C pointer to a Go slice for easier iterating
	for _, elem := range unsafe.Slice(l.elements, C.size_t(l.size)) {
		C.free(unsafe.Pointer(elem))
	}

	C.free(unsafe.Pointer(l.elements))
}

// takeEventClazzList copies a list of event classes that the SDK returned and releases the list.
func takeEventClazzList(thread *isolateThread, list *C.dxfg_v4_event_clazz_list) []int32 {
	if list == nil {
		return nil
	}
	defer C.dxfg_v4_event_clazz_list_release(thread.ptr, list)
	if list.size <= 0 || list.elements == nil {
		return nil
	}
	result := make([]int32, 0, int(list.size))
	for _, element := range unsafe.Slice(list.elements, int(list.size)) {
		if element != nil {
			result = append(result, int32(*element))
		}
	}
	return result
}
