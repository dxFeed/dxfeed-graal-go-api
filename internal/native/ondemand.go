// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>

// cgo exports drop const: the SDK passes const pointers, which the callback only reads.
extern void OnOnDemandPropertyChanged(graal_isolatethread_t *thread, char *property_name,
	dxfg_v4_on_demand_property_value_t *old_value, dxfg_v4_on_demand_property_value_t *new_value, void *user_data);
*/
import "C"
import (
	"errors"
	"unsafe"
)

// OnDemandServiceHandle holds the native on-demand service of an endpoint.
type OnDemandServiceHandle struct {
	self Handler
}

// OnDemandService returns the on-demand service of the endpoint, which must have the Feed or the OnDemandFeed role.
// Release the handle with Free; it does not stop the replay and does not close the endpoint.
func (e *DXEndpointHandle) OnDemandService() (*OnDemandServiceHandle, error) {
	ptr := e.ptr()
	if ptr == nil {
		return nil, errors.New("native: the endpoint is closed")
	}
	var service *C.dxfg_v4_on_demand_service_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_OnDemandService_getInstanceByEndpoint(thread.ptr, ptr, &service))
	})
	if err != nil {
		return nil, err
	}
	return &OnDemandServiceHandle{self: NewJavaHandle(unsafe.Pointer(service))}, nil
}

func (s *OnDemandServiceHandle) ptr() *C.dxfg_v4_on_demand_service_t {
	return (*C.dxfg_v4_on_demand_service_t)(s.self.Ptr())
}

// call runs f with the service in an isolate thread; it fails when the handle is released.
func (s *OnDemandServiceHandle) call(f func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t) error {
	ptr := s.ptr()
	if ptr == nil {
		return errors.New("native: the on-demand service is released")
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, f(thread, ptr))
	})
}

func (s *OnDemandServiceHandle) flag(get func(*C.graal_isolatethread_t, *C.dxfg_v4_on_demand_service_t,
	*C.int32_t) C.int32_t) (bool, error) {
	var value C.int32_t
	err := s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return get(thread.ptr, service, &value)
	})
	return value != 0, err
}

// IsReplaySupported reports whether the endpoint is connected to an on-demand data provider.
func (s *OnDemandServiceHandle) IsReplaySupported() (bool, error) {
	return s.flag(func(t *C.graal_isolatethread_t, p *C.dxfg_v4_on_demand_service_t, v *C.int32_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_isReplaySupported(t, p, v)
	})
}

// IsReplay reports whether the service is in the replay mode.
func (s *OnDemandServiceHandle) IsReplay() (bool, error) {
	return s.flag(func(t *C.graal_isolatethread_t, p *C.dxfg_v4_on_demand_service_t, v *C.int32_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_isReplay(t, p, v)
	})
}

// IsClear reports whether the service is in the clear mode.
func (s *OnDemandServiceHandle) IsClear() (bool, error) {
	return s.flag(func(t *C.graal_isolatethread_t, p *C.dxfg_v4_on_demand_service_t, v *C.int32_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_isClear(t, p, v)
	})
}

// Time returns the current or the last replayed time, in milliseconds since the Unix epoch.
func (s *OnDemandServiceHandle) Time() (int64, error) {
	var value C.int64_t
	err := s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_getTime(thread.ptr, service, &value)
	})
	return int64(value), err
}

// Speed returns the replay speed.
func (s *OnDemandServiceHandle) Speed() (float64, error) {
	var value C.double
	err := s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_getSpeed(thread.ptr, service, &value)
	})
	return float64(value), err
}

// Replay starts the replay from time, in milliseconds since the Unix epoch, at speed.
func (s *OnDemandServiceHandle) Replay(time int64, speed float64) error {
	return s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_replayWithSpeed(thread.ptr, service, C.int64_t(time), C.double(speed))
	})
}

// Pause pauses the replay.
func (s *OnDemandServiceHandle) Pause() error {
	return s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_pause(thread.ptr, service)
	})
}

// StopAndResume stops the replay and resumes the ordinary data feed.
func (s *OnDemandServiceHandle) StopAndResume() error {
	return s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_stopAndResume(thread.ptr, service)
	})
}

// StopAndClear stops the replay and clears the data.
func (s *OnDemandServiceHandle) StopAndClear() error {
	return s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_stopAndClear(thread.ptr, service)
	})
}

// SetSpeed changes the replay speed.
func (s *OnDemandServiceHandle) SetSpeed(speed float64) error {
	return s.call(func(thread *isolateThread, service *C.dxfg_v4_on_demand_service_t) C.int32_t {
		return C.dxfg_v4_OnDemandService_setSpeed(thread.ptr, service, C.double(speed))
	})
}

// Free releases the handle once; later calls do nothing.
func (s *OnDemandServiceHandle) Free() error {
	return s.self.Free()
}

// OnDemandPropertyCallback receives a change of a property of the on-demand service. A value is nil, a bool, a float64
// or an int64 time in milliseconds since the Unix epoch.
type OnDemandPropertyCallback func(name string, oldValue, newValue any)

func onDemandPropertyValue(v *C.dxfg_v4_on_demand_property_value_t) any {
	if v == nil {
		return nil
	}
	switch v.kind {
	case C.DXFG_V4_ON_DEMAND_PROPERTY_VALUE_BOOLEAN:
		return v.boolean_value != 0
	case C.DXFG_V4_ON_DEMAND_PROPERTY_VALUE_DOUBLE:
		return float64(v.double_value)
	case C.DXFG_V4_ON_DEMAND_PROPERTY_VALUE_TIME:
		return int64(v.time_value)
	default:
		return nil
	}
}

//export OnOnDemandPropertyChanged
func OnOnDemandPropertyChanged(thread *C.graal_isolatethread_t, name *C.char,
	oldValue, newValue *C.dxfg_v4_on_demand_property_value_t, userData unsafe.Pointer) {
	defer recoverCallback("on-demand property change listener")
	callback, ok := Restore(userData).(OnDemandPropertyCallback)
	if !ok {
		return
	}
	callback(C.GoString(name), onDemandPropertyValue(oldValue), onDemandPropertyValue(newValue))
}

// OnDemandPropertyListener is a native property change listener of an on-demand service.
type OnDemandPropertyListener struct {
	handle   Handler
	userData unsafe.Pointer
}

// AttachPropertyListener creates a native listener that passes the property changes to callback and adds it to the
// service.
func (s *OnDemandServiceHandle) AttachPropertyListener(callback OnDemandPropertyCallback) (*OnDemandPropertyListener,
	error) {
	ptr := s.ptr()
	if ptr == nil {
		return nil, errors.New("native: the on-demand service is released")
	}
	userData := Save(callback)
	var l *C.dxfg_v4_on_demand_property_change_listener_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		status := C.dxfg_v4_OnDemandPropertyChangeListener_new(thread.ptr,
			(*[0]byte)(C.OnOnDemandPropertyChanged), userData, &l)
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		return checkStatus(thread, C.dxfg_v4_OnDemandService_addPropertyChangeListener(thread.ptr, ptr, l))
	})
	if l == nil {
		Unref(userData)
		return nil, err
	}
	pl := &OnDemandPropertyListener{handle: NewJavaHandle(unsafe.Pointer(l)), userData: userData}
	if err != nil {
		return nil, errors.Join(err, pl.Release())
	}
	return pl, nil
}

// DetachPropertyListener removes the listener from the service, if it is not released, and releases the listener.
func (s *OnDemandServiceHandle) DetachPropertyListener(pl *OnDemandPropertyListener) error {
	var err error
	if ptr, l := s.ptr(), pl.handle.Ptr(); ptr != nil && l != nil {
		err = executeInIsolateThread(func(thread *isolateThread) error {
			return checkStatus(thread, C.dxfg_v4_OnDemandService_removePropertyChangeListener(thread.ptr, ptr,
				(*C.dxfg_v4_on_demand_property_change_listener_t)(l)))
		})
	}
	return errors.Join(err, pl.Release())
}

// Release releases the native listener and forgets the Go callback; later calls do nothing.
func (pl *OnDemandPropertyListener) Release() error {
	if pl.handle.Ptr() == nil {
		return nil
	}
	if err := pl.handle.Free(); err != nil {
		// The SDK may still call the listener: keep the Go callback for it.
		return err
	}
	Unref(pl.userData)
	return nil
}
