package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnEventReceived(graal_isolatethread_t *thread, dxfg_v4_event_type_list *events, void *user_data);
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

type DXFeedSubscription struct {
	ptr *C.dxfg_v4_subscription_t
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
	defer recoverCallback("event listener")
	listener, ok := Restore(userData).(common.EventListener)
	if !ok {
		return
	}
	events := newEventMapper().goEvents(eventsList)
	if len(events) == 0 {
		return
	}
	listener.Update(events)
}

// EventListener is a native listener attached to a subscription with AttachListener.
type EventListener struct {
	handle   Handler
	userData unsafe.Pointer
}

// AttachListener creates a native listener that passes the events to listener and adds it to the subscription.
func (s DXFeedSubscription) AttachListener(listener common.EventListener) (*EventListener, error) {
	userData := Save(listener)
	var l *C.dxfg_v4_feed_event_listener_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		status := C.dxfg_v4_DXFeedEventListener_new(thread.ptr, (*[0]byte)(C.OnEventReceived), userData, &l)
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_addEventListener(thread.ptr, s.ptr, l))
	})
	if l == nil {
		Unref(userData)
		return nil, err
	}
	el := &EventListener{handle: NewJavaHandle(unsafe.Pointer(l)), userData: userData}
	if err != nil {
		return nil, errors.Join(err, el.Release())
	}
	return el, nil
}

// DetachListener removes the listener from the subscription and releases it.
func (s DXFeedSubscription) DetachListener(el *EventListener) error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_removeEventListener(thread.ptr, s.ptr,
			(*C.dxfg_v4_feed_event_listener_t)(el.handle.Ptr())))
	})
	return errors.Join(err, el.Release())
}

// Release releases the native listener: the SDK does not call it after Release returns, so the Go listener is
// forgotten too. Later calls do nothing.
func (el *EventListener) Release() error {
	if el.handle.Ptr() == nil {
		return nil
	}
	if err := el.handle.Free(); err != nil {
		// The SDK may still call the listener: keep the Go listener for it.
		return err
	}
	Unref(el.userData)
	return nil
}

func (s DXFeedSubscription) AddSymbol(symbol osub.Symbol) error {
	return s.withSymbol(symbol, func(thread *isolateThread, cSymbol *C.dxfg_v4_symbol_t) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_addSymbol(thread.ptr, s.ptr, cSymbol)
	})
}

func (s DXFeedSubscription) AddSymbols(symbols ...osub.Symbol) error {
	return s.withSymbols(symbols, func(thread *isolateThread, list *C.dxfg_v4_symbol_list) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_addSymbols(thread.ptr, s.ptr, list)
	})
}

func (s DXFeedSubscription) RemoveSymbol(symbol osub.Symbol) error {
	return s.withSymbol(symbol, func(thread *isolateThread, cSymbol *C.dxfg_v4_symbol_t) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_removeSymbol(thread.ptr, s.ptr, cSymbol)
	})
}

func (s DXFeedSubscription) RemoveSymbols(symbols ...osub.Symbol) error {
	return s.withSymbols(symbols, func(thread *isolateThread, list *C.dxfg_v4_symbol_list) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_removeSymbols(thread.ptr, s.ptr, list)
	})
}

// withSymbol passes the symbol, built in C memory, to call and frees it afterwards; the SDK only reads it.
func (s DXFeedSubscription) withSymbol(symbol osub.Symbol,
	call func(thread *isolateThread, cSymbol *C.dxfg_v4_symbol_t) C.int32_t) error {
	cSymbol, err := newCSymbol(symbol)
	if err != nil {
		return err
	}
	defer freeCSymbol(cSymbol)
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, call(thread, cSymbol))
	})
}

// withSymbols passes the symbols, built in C memory, to call and frees them afterwards; the SDK only reads them.
func (s DXFeedSubscription) withSymbols(symbols []osub.Symbol,
	call func(thread *isolateThread, list *C.dxfg_v4_symbol_list) C.int32_t) error {
	list, err := newCSymbolList(symbols)
	if err != nil {
		return err
	}
	defer freeCSymbolList(list)
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, call(thread, list))
	})
}

func (s DXFeedSubscription) Clear() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_clear(thread.ptr, s.ptr))
	})
}

func (s DXFeedSubscription) Close() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_close(thread.ptr, s.ptr))
	})
}
