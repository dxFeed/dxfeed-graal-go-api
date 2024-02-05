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
	eventMappers := make(map[int32]mappers.MapperInterface)

	eventMappers[C.DXFG_EVENT_QUOTE] = mappers.QuoteMapper{}
	eventMappers[C.DXFG_EVENT_TIME_AND_SALE] = mappers.TimeAndSaleMapper{}
	eventMappers[C.DXFG_EVENT_PROFILE] = mappers.ProfileMapper{}
	eventMappers[C.DXFG_EVENT_ORDER] = mappers.OrderMapper{}
	eventMappers[C.DXFG_EVENT_SPREAD_ORDER] = mappers.SpreadOrderMapper{}

	size := len(elements)
	e := (**T)(C.malloc(C.size_t(size) * C.size_t(unsafe.Sizeof((*int)(nil)))))
	slice := unsafe.Slice(e, C.size_t(size))
	for i, element := range elements {
		slice[i] = allocElement[T, U](element, eventMappers)
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
	default:
		fmt.Printf("Couldn't alloc element for %T\n", element)
		return nil
	}
}
