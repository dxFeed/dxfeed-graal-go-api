package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"unsafe"
)

type CMapper interface {
	// Define methods that your C types must satisfy
}

type ListMapper[T CMapper] struct {
	size     C.int32_t
	elements **T
}

func NewListMapper[T CMapper, U comparable](elements []U) *ListMapper[T] {
	size := len(elements)
	e := (**T)(C.malloc(C.size_t(size) * C.size_t(unsafe.Sizeof((*int)(nil)))))
	slice := unsafe.Slice(e, C.size_t(size))
	for i, element := range elements {
		slice[i] = allocElement[T, U](element, mappers.AvailableMappers())
	}

	return &ListMapper[T]{
		elements: e,
		size:     C.int32_t(size),
	}
}

// allocElement maps an event to C memory; the callers check beforehand that every element is an event with a mapper.
func allocElement[T CMapper, U comparable](element U, mappers map[int32]mappers.MapperInterface) *T {
	event, ok := any(element).(events.EventType)
	if !ok {
		return nil
	}
	mapper := mappers[event.Type().NativeCode()]
	if mapper == nil {
		return nil
	}
	return (*T)(mapper.CEvent(event))
}
