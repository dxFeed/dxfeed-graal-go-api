package native

/*
#include "dxfg_api.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
)

// DXFeedTimeSeriesSubscription is a subscription created with DXFeedHandle.CreateTimeSeriesSubscription; its
// DXFeedSubscription is the base structure of the time series subscription.
type DXFeedTimeSeriesSubscription struct {
	DXFeedSubscription
	ts *C.dxfg_v4_time_series_subscription_t
}

// checkEventTypes fails for an event code that has no mapper: the events of such a type could not be read.
func checkEventTypes(eventTypes []int32) error {
	available := mappers.AvailableMappers()
	for _, eventType := range eventTypes {
		if available[eventType] == nil {
			return fmt.Errorf("cannot subscribe: event code %d is not supported yet", eventType)
		}
	}
	return nil
}

// withEventClazzList calls call with the event types as a C list.
func withEventClazzList(eventTypes []int32, call func(thread *isolateThread, list *C.dxfg_v4_event_clazz_list) C.int32_t) error {
	if err := checkEventTypes(eventTypes); err != nil {
		return err
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		list := createEventClazzList(eventTypes...)
		defer destroyEventClazzList(list)
		return checkStatus(thread, call(thread, (*C.dxfg_v4_event_clazz_list)(unsafe.Pointer(list))))
	})
}

// NewDetachedSubscription creates a subscription to the event types that is not attached to any feed.
func NewDetachedSubscription(eventTypes ...int32) (*DXFeedSubscription, error) {
	var ptr *C.dxfg_v4_subscription_t
	err := withEventClazzList(eventTypes, func(thread *isolateThread, list *C.dxfg_v4_event_clazz_list) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_newWithEventTypes(thread.ptr, list, &ptr)
	})
	if err != nil {
		return nil, err
	}
	return &DXFeedSubscription{ptr: ptr}, nil
}

// CreateTimeSeriesSubscription creates a time series subscription to the event types and attaches it to the feed.
func (f *DXFeedHandle) CreateTimeSeriesSubscription(eventTypes ...int32) (*DXFeedTimeSeriesSubscription, error) {
	var ptr *C.dxfg_v4_time_series_subscription_t
	err := withEventClazzList(eventTypes, func(thread *isolateThread, list *C.dxfg_v4_event_clazz_list) C.int32_t {
		return C.dxfg_v4_DXFeed_createTimeSeriesSubscriptionWithEventTypes(thread.ptr, f.ptr(), list, &ptr)
	})
	if err != nil {
		return nil, err
	}
	return &DXFeedTimeSeriesSubscription{DXFeedSubscription: DXFeedSubscription{
		ptr: (*C.dxfg_v4_subscription_t)(unsafe.Pointer(&ptr.sub))}, ts: ptr}, nil
}

// SetFromTime sets the time in milliseconds from which the subscription receives the events.
func (s *DXFeedTimeSeriesSubscription) SetFromTime(fromTime int64) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedTimeSeriesSubscription_setFromTime(thread.ptr, s.ts,
			C.int64_t(fromTime)))
	})
}

// FromTime returns the time in milliseconds from which the subscription receives the events.
func (s *DXFeedTimeSeriesSubscription) FromTime() (int64, error) {
	var fromTime C.int64_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedTimeSeriesSubscription_getFromTime(thread.ptr, s.ts, &fromTime))
	})
	return int64(fromTime), err
}

// AttachSubscription attaches the subscription to the feed.
func (f *DXFeedHandle) AttachSubscription(sub *DXFeedSubscription) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeed_attachSubscription(thread.ptr, f.ptr(), sub.ptr))
	})
}

// DetachSubscription detaches the subscription from the feed.
func (f *DXFeedHandle) DetachSubscription(sub *DXFeedSubscription) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeed_detachSubscription(thread.ptr, f.ptr(), sub.ptr))
	})
}

// DetachSubscriptionAndClear detaches the subscription from the feed and clears the data delivered to it.
func (f *DXFeedHandle) DetachSubscriptionAndClear(sub *DXFeedSubscription) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeed_detachSubscriptionAndClear(thread.ptr, f.ptr(), sub.ptr))
	})
}

// Attach attaches the subscription to the feed.
func (s DXFeedSubscription) Attach(feed *DXFeedHandle) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_attach(thread.ptr, s.ptr, feed.ptr()))
	})
}

// Detach detaches the subscription from the feed.
func (s DXFeedSubscription) Detach(feed *DXFeedHandle) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_detach(thread.ptr, s.ptr, feed.ptr()))
	})
}
