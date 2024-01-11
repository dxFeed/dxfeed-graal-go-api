package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/Osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
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
		slice[i] = allocElement[T, U](element)
	}

	return &ListMapper[T]{
		elements: e,
		size:     C.int32_t(size),
	}
}

func allocElement[T CMapper, U comparable](element U) *T {
	switch t := any(element).(type) {
	case int32:
		return (*T)(C.malloc(C.size_t(unsafe.Sizeof(element))))
	case *quote.Quote:
		return (*T)(unsafe.Pointer(newEventMapper().cQuote(t)))
	case *timeandsale.TimeAndSale:
		return (*T)(unsafe.Pointer(newEventMapper().cTimeAndSale(t)))
	case *profile.Profile:
		return (*T)(unsafe.Pointer(newEventMapper().cProfile(t)))
	case string:
		return (*T)(unsafe.Pointer(newEventMapper().cStringSymbol(t)))
	case Osub.WildcardSymbol:
		return (*T)(unsafe.Pointer(newEventMapper().cWildCardSymbol()))
	}
	return nil
}
