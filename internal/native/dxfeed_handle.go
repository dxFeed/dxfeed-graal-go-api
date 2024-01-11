package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DXFeedHandle struct {
	handle Handler
}

func NewDXFeedHandle(ptr *C.dxfg_feed_t) *DXFeedHandle {
	return &DXFeedHandle{handle: NewJavaHandle(unsafe.Pointer(ptr))}
}

func (f *DXFeedHandle) CreateSubscription(eventTypes ...int32) (*DXFeedSubscription, error) {
	var ptr *C.dxfg_subscription_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		list := createEventClazzList(eventTypes...)
		defer destroyEventClazzList(list)
		return checkCall(func() {
			ptr = C.dxfg_DXFeed_createSubscription2(thread.ptr, f.ptr(), (*C.dxfg_event_clazz_list_t)(unsafe.Pointer(list)))
		})
	})
	if err != nil {
		return nil, err
	}

	return &DXFeedSubscription{ptr: ptr}, nil
}

func (f *DXFeedHandle) Free() error {
	if f != nil {
		return f.handle.Free()
	}
	return nil
}

func (f *DXFeedHandle) ptr() *C.dxfg_feed_t {
	return (*C.dxfg_feed_t)(f.handle.Ptr())
}
