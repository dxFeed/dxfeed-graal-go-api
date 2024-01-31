package mappers

import "C"
import (
	"unsafe"
)

type MapperInterface interface {
	GoEvent(native unsafe.Pointer) interface{}
	CEvent(event interface{}) unsafe.Pointer
}
