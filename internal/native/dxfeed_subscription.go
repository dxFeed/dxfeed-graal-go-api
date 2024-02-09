package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnEventReceived(graal_isolatethread_t *thread, dxfg_event_type_list *events, void *user_data);
*/
import "C"
import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"unsafe"
)

type DXFeedSubscription struct {
	ptr *C.dxfg_subscription_t
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
func OnEventReceived(thread *C.graal_isolatethread_t, eventsList *C.dxfg_event_type_list, userData unsafe.Pointer) {
	Restore(userData).(common.EventListener).Update(newEventMapper().goEvents(eventsList))
}

func (s DXFeedSubscription) AttachListener(listener common.EventListener) error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		l := C.dxfg_DXFeedEventListener_new(thread.ptr, (*[0]byte)(C.OnEventReceived), Save(listener))
		C.dxfg_DXFeedSubscription_addEventListener(thread.ptr, s.ptr, l)
		return nil
	})
	return err
}

func (s DXFeedSubscription) AddSymbol(symbol any) error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		cSymbol := s.convertSymbol(symbol)
		if cSymbol != nil {
			C.dxfg_DXFeedSubscription_addSymbol(thread.ptr, s.ptr, cSymbol)
			return nil
		} else {
			return fmt.Errorf("Unsupported symbol %T!\n", symbol)
		}
	})
	return err
}

func (s DXFeedSubscription) AddSymbols(symbols ...any) error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_symbol_list, interface{}](symbols)
		C.dxfg_DXFeedSubscription_addSymbols(thread.ptr, s.ptr, (*C.dxfg_symbol_list)(unsafe.Pointer(l)))
		return nil
	})
	return err
}

func (s DXFeedSubscription) convertSymbol(symbol any) *C.dxfg_symbol_t {
	value := newEventMapper().cSymbol(symbol)
	return (*C.dxfg_symbol_t)(value)
}

func (s DXFeedSubscription) RemoveSymbol(symbol any) error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		cSymbol := s.convertSymbol(symbol)
		if cSymbol != nil {
			C.dxfg_DXFeedSubscription_removeSymbol(thread.ptr, s.ptr, cSymbol)
			return nil
		} else {
			return fmt.Errorf("Unsupported symbol %T!\n", symbol)
		}
	})
	return err
}

func (s DXFeedSubscription) RemoveSymbols(symbols ...any) error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_symbol_list, interface{}](symbols)
		C.dxfg_DXFeedSubscription_removeSymbols(thread.ptr, s.ptr, (*C.dxfg_symbol_list)(unsafe.Pointer(l)))
		return nil
	})
	return err
}

func (s DXFeedSubscription) Clear() {
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		C.dxfg_DXFeedSubscription_clear(thread.ptr, s.ptr)
		return nil
	})
}

func (s DXFeedSubscription) Close() {
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		C.dxfg_DXFeedSubscription_close(thread.ptr, s.ptr)
		return nil
	})
}
