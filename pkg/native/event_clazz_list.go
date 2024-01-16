package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type eventClazzList struct {
	size     C.int32_t
	elements **C.dxfg_event_clazz_t
}

func createEventClazzList(eventTypes ...int32) *eventClazzList {
	l := &eventClazzList{
		size: C.int32_t(len(eventTypes)),
	}

	if l.size == 0 {
		return l
	}

	l.elements = (**C.dxfg_event_clazz_t)(C.malloc(C.size_t(l.size) * C.size_t(unsafe.Sizeof((*C.dxfg_event_clazz_t)(nil)))))

	if l.elements == nil {
		panic("Memory allocation for eventClazzList.elements failed")
	}

	// Convert the C pointer to a Go slice for easier indexing
	slice := unsafe.Slice(l.elements, C.size_t(l.size))

	for i, eventType := range eventTypes {
		slice[i] = (*C.dxfg_event_clazz_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_event_clazz_t(0)))))

		if slice[i] == nil {
			// Free previously allocated elements before panicking
			for j := 0; j < i; j++ {
				C.free(unsafe.Pointer(slice[j]))
			}

			C.free(unsafe.Pointer(l.elements))

			panic("Memory allocation for eventClazzList's element failed")
		}

		*slice[i] = C.dxfg_event_clazz_t(eventType)
	}

	return l
}

func destroyEventClazzList(l *eventClazzList) {
	if l == nil || l.size <= 0 || l.elements == nil {
		return
	}

	// Convert the C pointer to a Go slice for easier indexing
	for _, elem := range unsafe.Slice(l.elements, C.size_t(l.size)) {
		C.free(unsafe.Pointer(elem))
	}

	C.free(unsafe.Pointer(l.elements))
}
