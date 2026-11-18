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
