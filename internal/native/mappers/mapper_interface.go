package mappers

import "C"
import (
	"unsafe"
)

type MapperInterface interface {
	GoEvent(native unsafe.Pointer) interface{}
	// CEvent maps the event to C memory; free it with FreeCEvent.
	CEvent(event interface{}) unsafe.Pointer
	// FreeCEvent frees an event made by CEvent together with its strings.
	FreeCEvent(event unsafe.Pointer)
}
