package native

/*
#include "dxfg_api.h"
#include <stdlib.h>

extern void OnInstrumentProfilesUpdated(graal_isolatethread_t *thread, dxfg_v4_iterable_ip_t *profiles, void *user_data);
extern void OnIpfConnectionStateChanged(graal_isolatethread_t *thread, dxfg_v4_ipf_connection_state_t old_state,
	dxfg_v4_ipf_connection_state_t new_state, void *user_data);
*/
import "C"
import (
	"errors"
	"log"
	"time"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// InstrumentProfileCollectorHandle holds a native instrument profile collector.
type InstrumentProfileCollectorHandle struct {
	self Handler
}

// NewInstrumentProfileCollector creates a collector; release it with Free.
func NewInstrumentProfileCollector() (*InstrumentProfileCollectorHandle, error) {
	var ptr *C.dxfg_v4_ipf_collector_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileCollector_new(thread.ptr, &ptr))
	})
	if ptr == nil {
		if err == nil {
			err = errors.New("native: the SDK returned no collector")
		}
		return nil, err
	}
	c := &InstrumentProfileCollectorHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}
	if err != nil {
		return nil, errors.Join(err, c.Free())
	}
	return c, nil
}

func (c *InstrumentProfileCollectorHandle) ptr() *C.dxfg_v4_ipf_collector_t {
	return (*C.dxfg_v4_ipf_collector_t)(c.self.Ptr())
}

func (c *InstrumentProfileCollectorHandle) call(f func(thread *isolateThread, collector *C.dxfg_v4_ipf_collector_t) error) error {
	ptr := c.ptr()
	if ptr == nil {
		return errors.New("native: the instrument profile collector is released")
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return f(thread, ptr)
	})
}

// Free releases the native collector; later calls do nothing.
func (c *InstrumentProfileCollectorHandle) Free() error {
	return c.self.Free()
}

// LastUpdateTime returns the time of the last update, in milliseconds since the Unix epoch.
func (c *InstrumentProfileCollectorHandle) LastUpdateTime() (int64, error) {
	var value C.int64_t
	err := c.call(func(thread *isolateThread, collector *C.dxfg_v4_ipf_collector_t) error {
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileCollector_getLastUpdateTime(thread.ptr, collector, &value))
	})
	return int64(value), err
}

// UpdateInstrumentProfiles passes the profiles to the collector; the collector keeps copies.
func (c *InstrumentProfileCollectorHandle) UpdateInstrumentProfiles(profiles []*events.InstrumentProfile) error {
	if err := checkProfiles(profiles); err != nil {
		return err
	}
	return c.call(func(thread *isolateThread, collector *C.dxfg_v4_ipf_collector_t) error {
		list, err := newCProfileList(thread, profiles)
		if err != nil {
			return err
		}
		defer freeCProfileList(thread, list)
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileCollector_updateInstrumentProfiles(thread.ptr, collector,
			list, nil))
	})
}

// View returns copies of the profiles of the collector that are not removed.
func (c *InstrumentProfileCollectorHandle) View() ([]*events.InstrumentProfile, error) {
	var result []*events.InstrumentProfile
	err := c.call(func(thread *isolateThread, collector *C.dxfg_v4_ipf_collector_t) error {
		var iterator *C.dxfg_v4_iterable_ip_t
		if err := checkStatus(thread, C.dxfg_v4_InstrumentProfileCollector_view(thread.ptr, collector,
			&iterator)); err != nil {
			return err
		}
		defer releaseJavaObject(thread, unsafe.Pointer(iterator))
		var err error
		result, err = goIterableProfiles(thread, iterator)
		return err
	})
	return result, err
}

// goIterableProfiles copies the profiles of an iterator to Go profiles; the iterator is not released.
func goIterableProfiles(thread *isolateThread, iterator *C.dxfg_v4_iterable_ip_t) ([]*events.InstrumentProfile,
	error) {
	var result []*events.InstrumentProfile
	for {
		var hasNext C.int32_t
		if err := checkStatus(thread, C.dxfg_v4_Iterable_InstrumentProfile_hasNext(thread.ptr, iterator,
			&hasNext)); err != nil {
			return result, err
		}
		if hasNext == 0 {
			return result, nil
		}
		var data *C.dxfg_v4_instrument_profile_data_t
		if err := checkStatus(thread, C.dxfg_v4_Iterable_InstrumentProfile_next(thread.ptr, iterator,
			&data)); err != nil {
			return result, err
		}
		profile, err := goProfile(thread, data)
		C.dxfg_v4_instrument_profile_data_release(thread.ptr, data)
		if err != nil {
			return result, err
		}
		result = append(result, profile)
	}
}

// InstrumentProfilesCallback receives the changed profiles of a collector, removed ones included (their type is
// REMOVED).
type InstrumentProfilesCallback func(profiles []*events.InstrumentProfile)

//export OnInstrumentProfilesUpdated
func OnInstrumentProfilesUpdated(thread *C.graal_isolatethread_t, profiles *C.dxfg_v4_iterable_ip_t,
	userData unsafe.Pointer) {
	defer recoverCallback("instrument profile update listener")
	callback, ok := Restore(userData).(InstrumentProfilesCallback)
	if !ok {
		return
	}
	changed, err := goIterableProfiles(&isolateThread{ptr: thread}, profiles)
	if err != nil {
		log.Printf("dxfeed: instrument profile update listener: %v", err)
	}
	callback(changed)
}

// InstrumentProfileUpdateListener is a native update listener of a collector.
type InstrumentProfileUpdateListener struct {
	nativeListener
}

// AttachUpdateListener creates a native listener that passes the updates to callback and adds it to the collector.
func (c *InstrumentProfileCollectorHandle) AttachUpdateListener(callback InstrumentProfilesCallback) (
	*InstrumentProfileUpdateListener, error,
) {
	userData := Save(callback)
	var l *C.dxfg_v4_ipf_update_listener_t
	err := c.call(func(thread *isolateThread, collector *C.dxfg_v4_ipf_collector_t) error {
		if err := checkStatus(thread, C.dxfg_v4_InstrumentProfileUpdateListener_new(thread.ptr,
			(*[0]byte)(C.OnInstrumentProfilesUpdated), userData, &l)); err != nil {
			return err
		}
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileCollector_addUpdateListener(thread.ptr, collector, l))
	})
	if l == nil {
		Unref(userData)
		if err == nil {
			err = errors.New("native: the SDK returned no update listener")
		}
		return nil, err
	}
	ul := &InstrumentProfileUpdateListener{nativeListener{handle: NewJavaHandle(unsafe.Pointer(l)), userData: userData}}
	if err != nil {
		return nil, errors.Join(err, ul.Release())
	}
	return ul, nil
}

// DetachUpdateListener removes the listener from the collector, if both are not released, and releases the listener.
func (c *InstrumentProfileCollectorHandle) DetachUpdateListener(ul *InstrumentProfileUpdateListener) error {
	var err error
	if ptr, l := c.ptr(), ul.handle.Ptr(); ptr != nil && l != nil {
		err = executeInIsolateThread(func(thread *isolateThread) error {
			return checkStatus(thread, C.dxfg_v4_InstrumentProfileCollector_removeUpdateListener(thread.ptr, ptr,
				(*C.dxfg_v4_ipf_update_listener_t)(l)))
		})
	}
	return errors.Join(err, ul.Release())
}

// InstrumentProfileConnectionHandle holds a native instrument profile connection.
type InstrumentProfileConnectionHandle struct {
	self Handler
}

// NewInstrumentProfileConnection creates a connection that pushes the profiles from the address into the collector;
// release it with Free, which does not close it.
func NewInstrumentProfileConnection(address string, collector *InstrumentProfileCollectorHandle) (
	*InstrumentProfileConnectionHandle, error,
) {
	if err := checkNoNul(address); err != nil {
		return nil, err
	}
	var ptr *C.dxfg_v4_ipf_connection_t
	err := collector.call(func(thread *isolateThread, c *C.dxfg_v4_ipf_collector_t) error {
		addressPtr, free := cString(address)
		defer free()
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileConnection_createConnection(thread.ptr, addressPtr, c,
			&ptr))
	})
	if ptr == nil {
		if err == nil {
			err = errors.New("native: the SDK returned no connection")
		}
		return nil, err
	}
	conn := &InstrumentProfileConnectionHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}
	if err != nil {
		return nil, errors.Join(err, conn.Free())
	}
	return conn, nil
}

func (c *InstrumentProfileConnectionHandle) call(f func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t) error {
	ptr := (*C.dxfg_v4_ipf_connection_t)(c.self.Ptr())
	if ptr == nil {
		return errors.New("native: the instrument profile connection is released")
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, f(thread, ptr))
	})
}

// Free releases the native connection; it does not close it. Later calls do nothing.
func (c *InstrumentProfileConnectionHandle) Free() error {
	return c.self.Free()
}

// Address returns the address of the connection without its options.
func (c *InstrumentProfileConnectionHandle) Address() (string, error) {
	var result string
	err := c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		var address *C.char
		status := C.dxfg_v4_InstrumentProfileConnection_getAddress(thread.ptr, conn, &address)
		if address != nil {
			result = C.GoString(address)
			C.dxfg_v4_String_release(thread.ptr, address)
		}
		return status
	})
	return result, err
}

// UpdatePeriod returns the update period in milliseconds.
func (c *InstrumentProfileConnectionHandle) UpdatePeriod() (int64, error) {
	var value C.int64_t
	err := c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		return C.dxfg_v4_InstrumentProfileConnection_getUpdatePeriod(thread.ptr, conn, &value)
	})
	return int64(value), err
}

// SetUpdatePeriod changes the update period in milliseconds, which must not be negative.
func (c *InstrumentProfileConnectionHandle) SetUpdatePeriod(millis int64) error {
	return c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		return C.dxfg_v4_InstrumentProfileConnection_setUpdatePeriod(thread.ptr, conn, C.int64_t(millis))
	})
}

// State returns the state of the connection: 0 not connected, 1 connecting, 2 connected, 3 completed, 4 closed.
func (c *InstrumentProfileConnectionHandle) State() (int32, error) {
	var value C.dxfg_v4_ipf_connection_state_t
	err := c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		return C.dxfg_v4_InstrumentProfileConnection_getState(thread.ptr, conn, &value)
	})
	return int32(value), err
}

// LastModified returns the modification time of the last received data, in milliseconds since the Unix epoch, or 0.
func (c *InstrumentProfileConnectionHandle) LastModified() (int64, error) {
	var value C.int64_t
	err := c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		return C.dxfg_v4_InstrumentProfileConnection_getLastModified(thread.ptr, conn, &value)
	})
	return int64(value), err
}

// Start starts the connection; it fails if the connection is not in the not connected state.
func (c *InstrumentProfileConnectionHandle) Start() error {
	return c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		return C.dxfg_v4_InstrumentProfileConnection_start(thread.ptr, conn)
	})
}

// Close closes the connection; closing a closed connection does nothing.
func (c *InstrumentProfileConnectionHandle) Close() error {
	return c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		return C.dxfg_v4_InstrumentProfileConnection_close(thread.ptr, conn)
	})
}

// WaitUntilCompleted waits until the initial profiles are fully read, the connection is closed or the timeout
// elapses; it reports whether the profiles were read.
func (c *InstrumentProfileConnectionHandle) WaitUntilCompleted(timeout time.Duration) (bool, error) {
	var completed C.int32_t
	err := c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		return C.dxfg_v4_InstrumentProfileConnection_waitUntilCompleted(thread.ptr, conn,
			C.int64_t(timeout.Milliseconds()), &completed)
	})
	return completed != 0, err
}

// IpfStateCallback receives a change of the state of a connection (see State).
type IpfStateCallback func(oldState, newState int32)

//export OnIpfConnectionStateChanged
func OnIpfConnectionStateChanged(_ *C.graal_isolatethread_t, oldState, newState C.dxfg_v4_ipf_connection_state_t,
	userData unsafe.Pointer) {
	defer recoverCallback("instrument profile connection state listener")
	callback, ok := Restore(userData).(IpfStateCallback)
	if !ok {
		return
	}
	callback(int32(oldState), int32(newState))
}

// IpfStateListener is a native state change listener of a connection.
type IpfStateListener struct {
	nativeListener
}

// AttachStateListener creates a native listener that passes the state changes to callback and adds it to the
// connection.
func (c *InstrumentProfileConnectionHandle) AttachStateListener(callback IpfStateCallback) (*IpfStateListener, error) {
	userData := Save(callback)
	var l *C.dxfg_v4_ipf_connection_state_change_listener_t
	err := c.call(func(thread *isolateThread, conn *C.dxfg_v4_ipf_connection_t) C.int32_t {
		if status := C.dxfg_v4_IpfPropertyChangeListener_new(thread.ptr,
			(*[0]byte)(C.OnIpfConnectionStateChanged), userData, &l); status != 0 {
			return status
		}
		return C.dxfg_v4_InstrumentProfileConnection_addStateChangeListener(thread.ptr, conn, l)
	})
	if l == nil {
		Unref(userData)
		if err == nil {
			err = errors.New("native: the SDK returned no state listener")
		}
		return nil, err
	}
	sl := &IpfStateListener{nativeListener{handle: NewJavaHandle(unsafe.Pointer(l)), userData: userData}}
	if err != nil {
		return nil, errors.Join(err, sl.Release())
	}
	return sl, nil
}

// DetachStateListener removes the listener from the connection, if both are not released, and releases the listener.
func (c *InstrumentProfileConnectionHandle) DetachStateListener(sl *IpfStateListener) error {
	var err error
	if ptr, l := c.self.Ptr(), sl.handle.Ptr(); ptr != nil && l != nil {
		err = executeInIsolateThread(func(thread *isolateThread) error {
			return checkStatus(thread, C.dxfg_v4_InstrumentProfileConnection_removeStateChangeListener(thread.ptr,
				(*C.dxfg_v4_ipf_connection_t)(ptr), (*C.dxfg_v4_ipf_connection_state_change_listener_t)(l)))
		})
	}
	return errors.Join(err, sl.Release())
}
