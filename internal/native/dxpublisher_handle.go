package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type DXPublisherHandle struct {
	handle Handler
}

func NewDXPublisherHandle(ptr *C.dxfg_publisher_t) *DXPublisherHandle {
	return &DXPublisherHandle{handle: NewJavaHandle(unsafe.Pointer(ptr))}
}

func (p *DXPublisherHandle) Free() error {
	if p != nil {
		return p.handle.Free()
	}
	return nil
}

func (p *DXPublisherHandle) Publish(events []interface{}) error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_event_type_list, interface{}](events)
		_ = C.dxfg_DXPublisher_publishEvents(thread.ptr, p.ptr(), (*C.dxfg_event_type_list)(unsafe.Pointer(l)))
		return nil
	})
	return err
}

func (p *DXPublisherHandle) ptr() *C.dxfg_publisher_t {
	return (*C.dxfg_publisher_t)(p.handle.Ptr())
}
