package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type Handler interface {
	Ptr() unsafe.Pointer
	Free() error
}

type JavaHandle struct {
	ptr unsafe.Pointer
}

func NewJavaHandle(ptr unsafe.Pointer) Handler {
	return &JavaHandle{ptr: ptr}
}

func (j *JavaHandle) Ptr() unsafe.Pointer {
	return j.ptr
}

func (j *JavaHandle) Free() error {
	if j.ptr == nil {
		return nil
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_JavaObjectHandler_release(thread.ptr, (*C.dxfg_java_object_handler)(j.ptr))
			j.ptr = nil
		})
	})
}
