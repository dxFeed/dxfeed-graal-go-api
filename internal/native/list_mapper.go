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
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"unsafe"
)

type CMapper interface {
	// Define methods that your C types must satisfy
}

// cList has the layout of the C list types: the size, then the elements.
type cList[T CMapper] struct {
	size     C.int32_t
	elements **T
}

// ListMapper is a C list of events. The list, its elements and the events with their strings are in the C memory of
// one allocator, so a list costs a few cgo calls whatever its size.
type ListMapper[T CMapper] struct {
	list  *cList[T]
	alloc mappers.Allocator
}

// eventSizeEstimate is the C memory that an event with its strings usually takes, for the first block of a list.
const eventSizeEstimate = 256

func NewListMapper[T CMapper, U comparable](elements []U) *ListMapper[T] {
	l := &ListMapper[T]{}
	l.alloc.Reserve(unsafe.Sizeof(cList[T]{}) + uintptr(len(elements))*(unsafe.Sizeof((*T)(nil))+eventSizeEstimate))
	l.list = (*cList[T])(l.alloc.Calloc(unsafe.Sizeof(cList[T]{})))
	l.list.size = C.int32_t(len(elements))
	l.list.elements = (**T)(l.alloc.Calloc(uintptr(len(elements)) * unsafe.Sizeof((*T)(nil))))
	slice := unsafe.Slice(l.list.elements, len(elements))
	available := mappers.AvailableMappers()
	for i, element := range elements {
		slice[i] = allocElement[T, U](&l.alloc, element, available)
	}
	return l
}

// Ptr returns the C list; it is valid until Free.
func (l *ListMapper[T]) Ptr() unsafe.Pointer { return unsafe.Pointer(l.list) }

// allocElement maps an event to C memory of the allocator; the callers check beforehand that every element is an
// event with a mapper.
func allocElement[T CMapper, U comparable](a *mappers.Allocator, element U, mappers map[int32]mappers.MapperInterface) *T {
	event, ok := any(element).(events.EventType)
	if !ok {
		return nil
	}
	mapper := mappers[event.Type().NativeCode()]
	if mapper == nil {
		return nil
	}
	return (*T)(mapper.CEvent(a, event))
}

// Free frees the list with its events; the list must not be used after that.
func (l *ListMapper[T]) Free() {
	if l == nil || l.list == nil {
		return
	}
	l.list = nil
	l.alloc.Free()
}
