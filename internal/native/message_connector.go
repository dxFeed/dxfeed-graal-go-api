package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnMessageConnectorChanged(graal_isolatethread_t *thread, dxfg_v4_message_connector_t *connector, void *user_data);
*/
import "C"
import (
	"errors"
	"unsafe"
)

// ConnectionErrorData is a copy of a connection error of the SDK.
type ConnectionErrorData struct {
	TimeMillis         int64
	Address            string
	Category           int32
	Dropped            bool
	ExceptionClassName string
	Message            string
	Count              int64
}

// MessageConnectorData is a snapshot of a message connector of the SDK.
type MessageConnectorData struct {
	Name                  string
	Kind                  int32
	Address               string
	CurrentAddresses      []string
	State                 int32
	Active                bool
	ConnectionCount       int32
	ClosedConnectionCount int64
	ReadBytes             int64
	WrittenBytes          int64
	ReadDataRecords       int64
	WrittenDataRecords    int64
	ActiveConnections     int64
	ClosedConnections     int64
	LastError             *ConnectionErrorData
}

// takeConnectionError copies an error that the SDK returned and releases it; nil for NULL.
func takeConnectionError(thread *isolateThread, e *C.dxfg_v4_connection_error_t) *ConnectionErrorData {
	if e == nil {
		return nil
	}
	defer C.dxfg_v4_ConnectionError_release(thread.ptr, e)
	data := &ConnectionErrorData{
		TimeMillis: int64(e.time),
		Category:   int32(e.category),
		Dropped:    e.dropped != 0,
		Count:      int64(e.count),
	}
	if e.address != nil {
		data.Address = C.GoString(e.address)
	}
	if e.exception_class_name != nil {
		data.ExceptionClassName = C.GoString(e.exception_class_name)
	}
	if e.message != nil {
		data.Message = C.GoString(e.message)
	}
	return data
}

// takeString copies a string that the SDK returned and releases it; "" for NULL.
func takeString(thread *isolateThread, s *C.char) string {
	if s == nil {
		return ""
	}
	defer C.dxfg_v4_String_release(thread.ptr, s)
	return C.GoString(s)
}

// readMessageConnector reads the values of a connector one by one, as the SDK offers them.
func readMessageConnector(thread *isolateThread, c *C.dxfg_v4_message_connector_t) (MessageConnectorData, error) {
	var d MessageConnectorData
	var name, address *C.char
	var kind, state, active, count C.int32_t
	var closed C.int64_t
	var addresses *C.dxfg_v4_string_list
	var stats C.dxfg_v4_message_connector_stats_t
	var lastError *C.dxfg_v4_connection_error_t
	err := checkStatus(thread, C.dxfg_v4_MessageConnector_getName(thread.ptr, c, &name))
	d.Name = takeString(thread, name)
	if err != nil {
		return d, err
	}
	err = checkStatus(thread, C.dxfg_v4_MessageConnector_getAddress(thread.ptr, c, &address))
	d.Address = takeString(thread, address)
	if err != nil {
		return d, err
	}
	if err := checkStatus(thread, C.dxfg_v4_MessageConnector_getCurrentAddresses(thread.ptr, c, &addresses)); err != nil {
		return d, err
	}
	d.CurrentAddresses = goStrings(addresses)
	releaseStringList(thread, addresses)
	for _, call := range []func() C.int32_t{
		func() C.int32_t { return C.dxfg_v4_MessageConnector_getKind(thread.ptr, c, &kind) },
		func() C.int32_t { return C.dxfg_v4_MessageConnector_getState(thread.ptr, c, &state) },
		func() C.int32_t { return C.dxfg_v4_MessageConnector_isActive(thread.ptr, c, &active) },
		func() C.int32_t { return C.dxfg_v4_MessageConnector_getConnectionCount(thread.ptr, c, &count) },
		func() C.int32_t { return C.dxfg_v4_MessageConnector_getClosedConnectionCount(thread.ptr, c, &closed) },
		func() C.int32_t { return C.dxfg_v4_MessageConnector_getStatistics(thread.ptr, c, &stats) },
		func() C.int32_t { return C.dxfg_v4_MessageConnector_getLastError(thread.ptr, c, &lastError) },
	} {
		if err := checkStatus(thread, call()); err != nil {
			return d, err
		}
	}
	d.Kind = int32(kind)
	d.State = int32(state)
	d.Active = active != 0
	d.ConnectionCount = int32(count)
	d.ClosedConnectionCount = int64(closed)
	d.ReadBytes = int64(stats.read_bytes)
	d.WrittenBytes = int64(stats.written_bytes)
	d.ReadDataRecords = int64(stats.read_data_records)
	d.WrittenDataRecords = int64(stats.written_data_records)
	d.ActiveConnections = int64(stats.active_connection_count)
	d.ClosedConnections = int64(stats.closed_connection_count)
	d.LastError = takeConnectionError(thread, lastError)
	return d, nil
}

// Connectors returns a snapshot of each message connector of the endpoint.
func (e *DXEndpointHandle) Connectors() ([]MessageConnectorData, error) {
	ptr := e.ptr()
	if ptr == nil {
		return nil, errors.New("native: the endpoint is closed")
	}
	var result []MessageConnectorData
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_message_connector_list
		if err := checkStatus(thread, C.dxfg_v4_DXEndpoint_getConnectors(thread.ptr, ptr, &list)); err != nil {
			return err
		}
		defer C.dxfg_v4_message_connector_list_release(thread.ptr, list)
		if list == nil || list.elements == nil || list.size <= 0 {
			return nil
		}
		connectors := unsafe.Slice(list.elements, int(list.size))
		result = make([]MessageConnectorData, 0, len(connectors))
		for _, c := range connectors {
			if c == nil {
				continue
			}
			d, err := readMessageConnector(thread, c)
			if err != nil {
				return err
			}
			result = append(result, d)
		}
		return nil
	})
	return result, err
}

// LastConnectionError returns the newest connection error of all the connectors of the endpoint; nil when there is
// none.
func (e *DXEndpointHandle) LastConnectionError() (*ConnectionErrorData, error) {
	ptr := e.ptr()
	if ptr == nil {
		return nil, errors.New("native: the endpoint is closed")
	}
	var result *ConnectionErrorData
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var ce *C.dxfg_v4_connection_error_t
		if err := checkStatus(thread, C.dxfg_v4_DXEndpoint_getLastConnectionError(thread.ptr, ptr, &ce)); err != nil {
			return err
		}
		result = takeConnectionError(thread, ce)
		return nil
	})
	return result, err
}

// MessageConnectorCallback receives a snapshot of a connector whose state changed, or the error of reading it.
type MessageConnectorCallback func(connector MessageConnectorData, err error)

//export OnMessageConnectorChanged
func OnMessageConnectorChanged(thread *C.graal_isolatethread_t, connector *C.dxfg_v4_message_connector_t,
	userData unsafe.Pointer) {
	defer recoverCallback("message connector listener")
	callback, ok := Restore(userData).(MessageConnectorCallback)
	if !ok || connector == nil {
		return
	}
	// The connector handle is valid during this call only: read it now.
	callback(readMessageConnector(&isolateThread{ptr: thread}, connector))
}

// MessageConnectorListener is a native listener of the connectors of an endpoint (AttachMessageConnectorListener).
type MessageConnectorListener struct {
	handle   Handler
	userData unsafe.Pointer
}

// AttachMessageConnectorListener creates a native listener that passes a snapshot of each changed connector to
// callback, on one SDK thread, and adds it to the endpoint.
func (e *DXEndpointHandle) AttachMessageConnectorListener(callback MessageConnectorCallback) (*MessageConnectorListener,
	error) {
	ptr := e.ptr()
	if ptr == nil {
		return nil, errors.New("native: the endpoint is closed")
	}
	userData := Save(callback)
	var l *C.dxfg_v4_message_connector_listener_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		status := C.dxfg_v4_MessageConnectorListener_new(thread.ptr, (*[0]byte)(C.OnMessageConnectorChanged), userData,
			&l)
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_addMessageConnectorListener(thread.ptr, ptr, l))
	})
	if l == nil {
		Unref(userData)
		return nil, err
	}
	ml := &MessageConnectorListener{handle: NewJavaHandle(unsafe.Pointer(l)), userData: userData}
	if err != nil {
		return nil, errors.Join(err, ml.Release())
	}
	return ml, nil
}

// DetachMessageConnectorListener removes the listener from the endpoint, if it is not closed, and releases it.
func (e *DXEndpointHandle) DetachMessageConnectorListener(ml *MessageConnectorListener) error {
	var err error
	if ptr, l := e.ptr(), ml.handle.Ptr(); ptr != nil && l != nil {
		err = executeInIsolateThread(func(thread *isolateThread) error {
			return checkStatus(thread, C.dxfg_v4_DXEndpoint_removeMessageConnectorListener(thread.ptr, ptr,
				(*C.dxfg_v4_message_connector_listener_t)(l)))
		})
	}
	return errors.Join(err, ml.Release())
}

// Release releases the native listener: the SDK does not call it after Release returns, so the Go callback is
// forgotten too. Later calls do nothing.
func (ml *MessageConnectorListener) Release() error {
	if ml.handle.Ptr() == nil {
		return nil
	}
	if err := ml.handle.Free(); err != nil {
		// The SDK may still call the listener: keep the Go callback for it.
		return err
	}
	Unref(ml.userData)
	return nil
}
