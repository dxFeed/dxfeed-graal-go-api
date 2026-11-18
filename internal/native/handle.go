package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"sync/atomic"
	"unsafe"
)

type Handler interface {
	Ptr() unsafe.Pointer
	Free() error
}

// JavaHandle holds a handle of a Java object. It may be used and freed from several goroutines: a wait abandoned on
// a done context still runs while the owner closes the object.
type JavaHandle struct {
	ptr atomic.Pointer[byte]
}

func NewJavaHandle(ptr unsafe.Pointer) Handler {
	j := &JavaHandle{}
	j.ptr.Store((*byte)(ptr))
	return j
}

// Ptr returns the handle, or nil after Free.
func (j *JavaHandle) Ptr() unsafe.Pointer {
	return unsafe.Pointer(j.ptr.Load())
}

// Free releases the handle once; later calls do nothing.
func (j *JavaHandle) Free() error {
	ptr := j.ptr.Swap(nil)
	if ptr == nil {
		return nil
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_JavaObjectHandler_release(thread.ptr, (*C.dxfg_v4_java_object_handler)(unsafe.Pointer(ptr))))
	})
}
