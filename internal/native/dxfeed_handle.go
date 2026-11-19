// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DXFeedHandle struct {
	handle   Handler
	promises promiseSet
}

func NewDXFeedHandle(ptr *C.dxfg_v4_feed_t) *DXFeedHandle {
	return &DXFeedHandle{handle: NewJavaHandle(unsafe.Pointer(ptr))}
}

func (f *DXFeedHandle) CreateSubscription(eventTypes ...int32) (*DXFeedSubscription, error) {
	var ptr *C.dxfg_v4_subscription_t
	err := withEventClazzList(eventTypes, func(thread *isolateThread, list *C.dxfg_v4_event_clazz_list) C.int32_t {
		return C.dxfg_v4_DXFeed_createSubscriptionWithEventTypes(thread.ptr, f.ptr(), list, &ptr)
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
