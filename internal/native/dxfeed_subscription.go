package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnEventReceived(graal_isolatethread_t *thread, dxfg_v4_event_type_list *events, void *user_data);
*/
import "C"
import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"unsafe"
)

type DXFeedSubscription struct {
	ptr *C.dxfg_v4_subscription_t
}

type dxfg_symbol_t struct {
	t      C.int32_t
	symbol *C.char
}

type dxfg_time_series_subscription_symbol_t struct {
	t         C.int32_t
	symbol    *dxfg_symbol_t
	from_time C.int64_t
}

type dxfg_indexed_event_source_t struct {
	t    C.int32_t
	id   C.int32_t
	name *C.char
}

type dxfg_indexed_event_subscription_symbol_t struct {
	t      C.int32_t
	symbol *dxfg_symbol_t
	source *dxfg_indexed_event_source_t
}

func ConvertString(value *C.char) *string {
	if value == nil {
		return nil
	} else {
		result := C.GoString(value)
		return &result
	}
}

//export OnEventReceived
func OnEventReceived(thread *C.graal_isolatethread_t, eventsList *C.dxfg_v4_event_type_list, userData unsafe.Pointer) {
	events := newEventMapper().goEvents(eventsList)
	if len(events) == 0 {
		return
	}
	Restore(userData).(common.EventListener).Update(events)
}

func (s DXFeedSubscription) AttachListener(listener common.EventListener) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		var l *C.dxfg_v4_feed_event_listener_t
		status := C.dxfg_v4_DXFeedEventListener_new(thread.ptr, (*[0]byte)(C.OnEventReceived), Save(listener), &l)
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_addEventListener(thread.ptr, s.ptr, l))
	})
}

func (s DXFeedSubscription) AddSymbol(symbol any) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		cSymbol := s.convertSymbol(symbol)
		if cSymbol == nil {
			return fmt.Errorf("unsupported symbol %T", symbol)
		}
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_addSymbol(thread.ptr, s.ptr, cSymbol))
	})
}

func (s DXFeedSubscription) AddSymbols(symbols ...any) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_v4_symbol_list, interface{}](symbols)
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_addSymbols(thread.ptr, s.ptr,
			(*C.dxfg_v4_symbol_list)(unsafe.Pointer(l))))
	})
}

func (s DXFeedSubscription) convertSymbol(symbol any) *C.dxfg_v4_symbol_t {
	value := newEventMapper().cSymbol(symbol)
	return (*C.dxfg_v4_symbol_t)(value)
}

func (s DXFeedSubscription) RemoveSymbol(symbol any) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		cSymbol := s.convertSymbol(symbol)
		if cSymbol == nil {
			return fmt.Errorf("unsupported symbol %T", symbol)
		}
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_removeSymbol(thread.ptr, s.ptr, cSymbol))
	})
}

func (s DXFeedSubscription) RemoveSymbols(symbols ...any) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_v4_symbol_list, interface{}](symbols)
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_removeSymbols(thread.ptr, s.ptr,
			(*C.dxfg_v4_symbol_list)(unsafe.Pointer(l))))
	})
}

// Clear and Close have no error result yet: a failure is taken from the thread and dropped.
func (s DXFeedSubscription) Clear() {
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_clear(thread.ptr, s.ptr))
	})
}

func (s DXFeedSubscription) Close() {
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_close(thread.ptr, s.ptr))
	})
}
