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
	"fmt"
	"reflect"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	goevents "github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// checkEvents checks that every value is a non-nil event of a type with a Go mapper and with an event symbol; verb
// names the operation in the error.
func checkEvents(verb string, events []interface{}) error {
	available := mappers.AvailableMappers()
	for _, event := range events {
		e, ok := event.(goevents.EventType)
		if !ok {
			return fmt.Errorf("cannot %s %T: not an event type", verb, event)
		}
		if rv := reflect.ValueOf(event); rv.Kind() == reflect.Ptr && rv.IsNil() {
			return fmt.Errorf("cannot %s %T: nil event", verb, event)
		}
		if available[e.Type().NativeCode()] == nil {
			return fmt.Errorf("cannot %s %T: event code %d is not supported yet", verb, event, e.Type())
		}
		if s, ok := event.(interface{ EventSymbol() *string }); ok && s.EventSymbol() == nil {
			return fmt.Errorf("cannot %s %T: event symbol is not set", verb, event)
		}
	}
	return nil
}

// checkEventType checks that the event class has a Go mapper; verb names the operation in the error.
func checkEventType(verb string, eventType int32) error {
	if mappers.AvailableMappers()[eventType] == nil {
		return fmt.Errorf("cannot %s: event code %d is not supported yet", verb, eventType)
	}
	return nil
}

// takeEvent converts an event that the SDK returned and releases it; a nil event gives nil.
func takeEvent(thread *isolateThread, event *C.dxfg_v4_event_type_t) (interface{}, error) {
	if event == nil {
		return nil, nil
	}
	result := newEventMapper().goEvent(event)
	return result, checkStatus(thread, C.dxfg_v4_EventType_release(thread.ptr, event))
}

// takeEvents converts a list of events that the SDK returned and releases it; the result is never nil.
func takeEvents(thread *isolateThread, list *C.dxfg_v4_event_type_list) ([]interface{}, error) {
	if list == nil {
		return []interface{}{}, nil
	}
	result := newEventMapper().goEvents(list)
	if result == nil {
		result = []interface{}{}
	}
	return result, checkStatus(thread, C.dxfg_v4_event_type_list_release(thread.ptr, list))
}

// withCSymbol runs f with the symbol built in C memory and frees it afterwards.
func withCSymbol(symbol osub.Symbol, f func(cSymbol *C.dxfg_v4_symbol_t) error) error {
	cSymbol, err := newCSymbol(symbol)
	if err != nil {
		return err
	}
	defer freeCSymbol(cSymbol)
	return f(cSymbol)
}

// GetLastEvent returns the last event for the type and symbol of the given event, or a copy of the event when the last
// event is not available.
func (f *DXFeedHandle) GetLastEvent(event interface{}) (interface{}, error) {
	if err := checkEvents("get the last event of", []interface{}{event}); err != nil {
		return nil, err
	}
	e := event.(goevents.EventType)
	mapper := mappers.AvailableMappers()[e.Type().NativeCode()]
	var result interface{}
	err := executeInIsolateThread(func(thread *isolateThread) error {
		cEvent := mapper.CEvent(e)
		// The SDK copies the event; the caller keeps its ownership.
		defer mapper.FreeCEvent(cEvent)
		var last *C.dxfg_v4_event_type_t
		if err := checkStatus(thread, C.dxfg_v4_DXFeed_getLastEvent(thread.ptr, f.ptr(),
			(*C.dxfg_v4_event_type_t)(cEvent), &last)); err != nil {
			return err
		}
		var err error
		result, err = takeEvent(thread, last)
		return err
	})
	return result, err
}

// GetLastEvents returns the last event for every given event, in the same order (see GetLastEvent).
func (f *DXFeedHandle) GetLastEvents(events []interface{}) ([]interface{}, error) {
	if err := checkEvents("get the last event of", events); err != nil {
		return nil, err
	}
	var result []interface{}
	err := executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_v4_event_type_list, interface{}](events)
		// The SDK copies the events; the caller keeps ownership of the list.
		defer l.Free()
		var last *C.dxfg_v4_event_type_list
		if err := checkStatus(thread, C.dxfg_v4_DXFeed_getLastEvents(thread.ptr, f.ptr(),
			(*C.dxfg_v4_event_type_list)(unsafe.Pointer(l)), &last)); err != nil {
			return err
		}
		var err error
		result, err = takeEvents(thread, last)
		return err
	})
	return result, err
}

// GetLastEventIfSubscribed returns the last event of the class for the symbol, or nil if there is no subscription for
// them or the event has not arrived yet.
func (f *DXFeedHandle) GetLastEventIfSubscribed(eventType int32, symbol osub.Symbol) (interface{}, error) {
	if err := checkEventType("get the last event", eventType); err != nil {
		return nil, err
	}
	var result interface{}
	err := withCSymbol(symbol, func(cSymbol *C.dxfg_v4_symbol_t) error {
		return executeInIsolateThread(func(thread *isolateThread) error {
			var event *C.dxfg_v4_event_type_t
			if err := checkStatus(thread, C.dxfg_v4_DXFeed_getLastEventIfSubscribed(thread.ptr, f.ptr(),
				C.dxfg_v4_event_clazz_t(eventType), cSymbol, &event)); err != nil {
				return err
			}
			var err error
			result, err = takeEvent(thread, event)
			return err
		})
	})
	return result, err
}

// GetIndexedEventsIfSubscribed returns the indexed events of the class for the symbol and the source name; the list is
// empty if there is no subscription for them or the events have not arrived yet.
func (f *DXFeedHandle) GetIndexedEventsIfSubscribed(eventType int32, symbol osub.Symbol,
	source string) ([]interface{}, error) {
	if err := checkEventType("get the indexed events", eventType); err != nil {
		return nil, err
	}
	var result []interface{}
	err := withCSymbol(symbol, func(cSymbol *C.dxfg_v4_symbol_t) error {
		cSource := C.CString(source)
		defer C.free(unsafe.Pointer(cSource))
		return executeInIsolateThread(func(thread *isolateThread) error {
			var list *C.dxfg_v4_event_type_list
			if err := checkStatus(thread, C.dxfg_v4_DXFeed_getIndexedEventsIfSubscribed(thread.ptr, f.ptr(),
				C.dxfg_v4_event_clazz_t(eventType), cSymbol, cSource, &list)); err != nil {
				return err
			}
			var err error
			result, err = takeEvents(thread, list)
			return err
		})
	})
	return result, err
}

// GetTimeSeriesIfSubscribed returns the time series events of the class for the symbol from fromTime to toTime
// (inclusive, milliseconds since the epoch); the list is empty if there is no subscription for them or the events have
// not arrived yet.
func (f *DXFeedHandle) GetTimeSeriesIfSubscribed(eventType int32, symbol osub.Symbol, fromTime,
	toTime int64) ([]interface{}, error) {
	if err := checkEventType("get the time series", eventType); err != nil {
		return nil, err
	}
	if _, ok := symbol.(*osub.TimeSeriesSubscriptionSymbol); ok {
		return nil, errors.New("cannot get the time series: the symbol must not be a time series subscription symbol")
	}
	var result []interface{}
	err := withCSymbol(symbol, func(cSymbol *C.dxfg_v4_symbol_t) error {
		return executeInIsolateThread(func(thread *isolateThread) error {
			var list *C.dxfg_v4_event_type_list
			if err := checkStatus(thread, C.dxfg_v4_DXFeed_getTimeSeriesIfSubscribed(thread.ptr, f.ptr(),
				C.dxfg_v4_event_clazz_t(eventType), cSymbol, C.int64_t(fromTime), C.int64_t(toTime),
				&list)); err != nil {
				return err
			}
			var err error
			result, err = takeEvents(thread, list)
			return err
		})
	})
	return result, err
}
