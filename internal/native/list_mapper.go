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

func allocElement[T CMapper, U comparable](element U, mappers map[int32]mappers.MapperInterface) *T {
	switch t := any(element).(type) {
	case int32:
		return (*T)(C.malloc(C.size_t(unsafe.Sizeof(element))))
	case events.EventType:
		// all market events have to implement this interface
		mapper := mappers[int32(t.Type())]
		return (*T)(mapper.CEvent(t))
	case string:
		return (*T)(unsafe.Pointer(newEventMapper().cStringSymbol(t)))
	case Osub.WildcardSymbol:
		return (*T)(unsafe.Pointer(newEventMapper().cWildCardSymbol()))
	case Osub.TimeSeriesSubscriptionSymbol:
		return (*T)(unsafe.Pointer(newEventMapper().cTimeSeriesSymbol(t.GetSymbol(), t.GetFromTime())))
	default:
		fmt.Printf("Couldn't alloc element for %T\n", element)
		return nil
	}
}
