// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// ObservableSubscriptionHandle is the subscription of a publisher to the events of one type; get it with
// DXPublisherHandle.Subscription. The publisher releases it in Free.
type ObservableSubscriptionHandle struct {
	handle Handler
}

// Subscription returns the observable subscription of the publisher to the event class; later calls return the same
// handle.
func (p *DXPublisherHandle) Subscription(eventType int32) (*ObservableSubscriptionHandle, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if sub := p.subscriptions[eventType]; sub != nil {
		return sub, nil
	}
	if p.handle.Ptr() == nil {
		return nil, errors.New("native: the publisher is released")
	}
	var ptr *C.dxfg_v4_observable_subscription_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXPublisher_getSubscription(thread.ptr, p.ptr(),
			C.dxfg_v4_event_clazz_t(eventType), &ptr))
	})
	if ptr == nil {
		if err == nil {
			err = errors.New("native: the SDK returned no subscription")
		}
		return nil, err
	}
	sub := &ObservableSubscriptionHandle{handle: NewJavaHandle(unsafe.Pointer(ptr))}
	if err != nil {
		return nil, errors.Join(err, sub.handle.Free())
	}
	if p.subscriptions == nil {
		p.subscriptions = make(map[int32]*ObservableSubscriptionHandle)
	}
	p.subscriptions[eventType] = sub
	return sub, nil
}

// errSubscriptionReleased is returned for a subscription whose native handle the publisher has released with the
// endpoint.
var errSubscriptionReleased = errors.New("native: the subscription is released with the endpoint")

func (s *ObservableSubscriptionHandle) ptr() *C.dxfg_v4_observable_subscription_t {
	return (*C.dxfg_v4_observable_subscription_t)(s.handle.Ptr())
}

// IsClosed reports whether the subscription is closed; a subscription released with the endpoint is closed.
func (s *ObservableSubscriptionHandle) IsClosed() (bool, error) {
	if s.handle.Ptr() == nil {
		return true, nil
	}
	var closed C.int32_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_ObservableSubscription_isClosed(thread.ptr, s.ptr(), &closed))
	})
	return closed != 0, err
}

// EventTypes returns the event classes of the subscription.
func (s *ObservableSubscriptionHandle) EventTypes() ([]int32, error) {
	if s.handle.Ptr() == nil {
		return nil, errSubscriptionReleased
	}
	var result []int32
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_event_clazz_list
		if err := checkStatus(thread, C.dxfg_v4_ObservableSubscription_getEventTypes(thread.ptr, s.ptr(), &list)); err != nil {
			return err
		}
		result = takeEventClazzList(thread, list)
		return nil
	})
	return result, err
}

// ContainsEventType reports whether the subscription has the event class.
func (s *ObservableSubscriptionHandle) ContainsEventType(eventType int32) (bool, error) {
	if s.handle.Ptr() == nil {
		return false, errSubscriptionReleased
	}
	var contains C.int32_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_ObservableSubscription_containsEventType(thread.ptr, s.ptr(),
			C.dxfg_v4_event_clazz_t(eventType), &contains))
	})
	return contains != 0, err
}

// AddChangeListener adds the change listener to the subscription.
func (s *ObservableSubscriptionHandle) AddChangeListener(cl *ChangeListener) error {
	if s.handle.Ptr() == nil {
		return errSubscriptionReleased
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_ObservableSubscription_addChangeListener(thread.ptr, s.ptr(), cl.ptr()))
	})
}

// RemoveChangeListener removes the change listener from the subscription; it does not release it. A subscription
// released with the endpoint holds no listeners, so there is nothing to remove.
func (s *ObservableSubscriptionHandle) RemoveChangeListener(cl *ChangeListener) error {
	if s.handle.Ptr() == nil {
		return nil
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_ObservableSubscription_removeChangeListener(thread.ptr, s.ptr(), cl.ptr()))
	})
}
