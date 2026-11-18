package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
)

type DXFeedHandle struct {
	handle Handler
}

func NewDXFeedHandle(ptr *C.dxfg_v4_feed_t) *DXFeedHandle {
	return &DXFeedHandle{handle: NewJavaHandle(unsafe.Pointer(ptr))}
}

func (f *DXFeedHandle) CreateSubscription(eventTypes ...int32) (*DXFeedSubscription, error) {
	available := mappers.AvailableMappers()
	for _, eventType := range eventTypes {
		if available[eventType] == nil {
			return nil, fmt.Errorf("cannot subscribe: event code %d is not supported yet", eventType)
		}
	}
	var ptr *C.dxfg_v4_subscription_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		list := createEventClazzList(eventTypes...)
		defer destroyEventClazzList(list)
		return checkStatus(thread, C.dxfg_v4_DXFeed_createSubscriptionWithEventTypes(thread.ptr, f.ptr(),
			(*C.dxfg_v4_event_clazz_list)(unsafe.Pointer(list)), &ptr))
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

func (f *DXFeedHandle) ptr() *C.dxfg_v4_feed_t {
	return (*C.dxfg_v4_feed_t)(f.handle.Ptr())
}
