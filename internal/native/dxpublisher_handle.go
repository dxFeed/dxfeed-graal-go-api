// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

type DXPublisherHandle struct {
	handle Handler

	mu            sync.Mutex
	subscriptions map[int32]*ObservableSubscriptionHandle
}

func NewDXPublisherHandle(ptr *C.dxfg_v4_publisher_t) *DXPublisherHandle {
	return &DXPublisherHandle{handle: NewJavaHandle(unsafe.Pointer(ptr))}
}

// Free releases the publisher and its observable subscriptions.
func (p *DXPublisherHandle) Free() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	subscriptions := p.subscriptions
	p.subscriptions = nil
	p.mu.Unlock()
	var err error
	for _, sub := range subscriptions {
		err = errors.Join(err, sub.handle.Free())
	}
	return errors.Join(err, p.handle.Free())
}

func (p *DXPublisherHandle) Publish(events []interface{}) error {
	if err := checkEvents("publish", events); err != nil {
		return err
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_v4_event_type_list, interface{}](events)
		// The SDK copies the events; the caller keeps ownership of the list.
		defer l.Free()
		return checkStatus(thread, C.dxfg_v4_DXPublisher_publishEvents(thread.ptr, p.ptr(),
			(*C.dxfg_v4_event_type_list)(l.Ptr())))
	})
}

func (p *DXPublisherHandle) ptr() *C.dxfg_v4_publisher_t {
	return (*C.dxfg_v4_publisher_t)(p.handle.Ptr())
}
