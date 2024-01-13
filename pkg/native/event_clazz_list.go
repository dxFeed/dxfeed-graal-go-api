package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type eventClazzList struct {
	size     C.int32_t
	elements **C.dxfg_event_clazz_t
}

func createEventClazzList(eventTypes ...int32) *eventClazzList {
	e := &eventClazzList{
		size: C.int32_t(len(eventTypes)),
	}

	ptrSize := unsafe.Sizeof(uintptr(0))
	e.elements = (**C.dxfg_event_clazz_t)(C.malloc(C.size_t(e.size) * C.size_t(ptrSize)))

	for i, eventType := range eventTypes {
		element := C.dxfg_event_clazz_t(eventType)
		ptr := (**C.dxfg_event_clazz_t)(unsafe.Pointer(uintptr(unsafe.Pointer(e.elements)) + uintptr(i)*ptrSize))
		*ptr = &element
	}

	return e
}
