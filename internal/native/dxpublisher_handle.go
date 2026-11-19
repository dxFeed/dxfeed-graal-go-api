package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	goevents "github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
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
	available := mappers.AvailableMappers()
	for _, event := range events {
		e, ok := event.(goevents.EventType)
		if !ok {
			return fmt.Errorf("cannot publish %T: not an event type", event)
		}
		if available[e.Type().NativeCode()] == nil {
			return fmt.Errorf("cannot publish %T: event code %d is not supported yet", event, e.Type())
		}
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_v4_event_type_list, interface{}](events)
		// The SDK copies the events; the caller keeps ownership of the list.
		defer l.Free()
		return checkStatus(thread, C.dxfg_v4_DXPublisher_publishEvents(thread.ptr, p.ptr(),
			(*C.dxfg_v4_event_type_list)(unsafe.Pointer(l))))
	})
}

func (p *DXPublisherHandle) ptr() *C.dxfg_v4_publisher_t {
	return (*C.dxfg_v4_publisher_t)(p.handle.Ptr())
}
