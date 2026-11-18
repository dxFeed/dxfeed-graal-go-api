package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnStateChanged(graal_isolatethread_t *thread, dxfg_v4_endpoint_state_t old_state, dxfg_v4_endpoint_state_t new_state, void *user_data);
*/
import "C"
import (
	"errors"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"sync"
	"unsafe"
)

type DXEndpointHandle struct {
	self     Handler
	feedOnce sync.Once
	feed     *DXFeedHandle
	feedErr  error

	publisherOnce sync.Once
	publisher     *DXPublisherHandle
	publisherErr  error
}

func NewDXEndpointHandle(role common.Role) (*DXEndpointHandle, error) {
	var ptr *C.dxfg_v4_endpoint_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_createWithRole(thread.ptr, (C.dxfg_v4_endpoint_role_t)(role), &ptr))
	})
	if err != nil {
		return nil, err
	}

	return &DXEndpointHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

func NewDXEndpointHandleWithProperties(role common.Role, properties map[string]string) (*DXEndpointHandle, error) {
	var ptr *C.dxfg_v4_endpoint_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var builder *C.dxfg_v4_endpoint_builder_t
		if err := checkStatus(thread, C.dxfg_v4_DXEndpoint_newBuilder(thread.ptr, &builder)); err != nil {
			return err
		}
		defer C.dxfg_v4_JavaObjectHandler_release(thread.ptr, (*C.dxfg_v4_java_object_handler)(unsafe.Pointer(builder)))
		if err := checkStatus(thread, C.dxfg_v4_DXEndpoint_Builder_withRole(thread.ptr, builder, (C.dxfg_v4_endpoint_role_t)(role))); err != nil {
			return err
		}
		for key, value := range properties {
			// The SDK copies both strings during the call, so they are freed right after it.
			keyPtr := C.CString(key)
			valuePtr := C.CString(value)
			status := C.dxfg_v4_DXEndpoint_Builder_withProperty(thread.ptr, builder, keyPtr, valuePtr)
			C.free(unsafe.Pointer(keyPtr))
			C.free(unsafe.Pointer(valuePtr))
			if err := checkStatus(thread, status); err != nil {
				return err
			}
		}
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_Builder_build(thread.ptr, builder, &ptr))
	})
	if err != nil {
		return nil, err
	}

	return &DXEndpointHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

func (e *DXEndpointHandle) Close() error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_close(thread.ptr, e.ptr()))
	})
	return errors.Join(err, e.Free())
}

func (e *DXEndpointHandle) CloseAndAwaitTermination() error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_closeAndAwaitTermination(thread.ptr, e.ptr()))
	})
	return errors.Join(err, e.Free())
}

func (e *DXEndpointHandle) Connect(address string) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		addressPtr := C.CString(address)
		defer C.free(unsafe.Pointer(addressPtr))

		return checkStatus(thread, C.dxfg_v4_DXEndpoint_connect(thread.ptr, e.ptr(), addressPtr))
	})
}

func (e *DXEndpointHandle) Reconnect() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_reconnect(thread.ptr, e.ptr()))
	})
}

func (e *DXEndpointHandle) Disconnect() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_disconnect(thread.ptr, e.ptr()))
	})
}

func (e *DXEndpointHandle) DisconnectAndClear() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_disconnectAndClear(thread.ptr, e.ptr()))
	})
}

func (e *DXEndpointHandle) AwaitProcessed() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_awaitProcessed(thread.ptr, e.ptr()))
	})
}

func (e *DXEndpointHandle) AwaitNotConnected() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_awaitNotConnected(thread.ptr, e.ptr()))
	})
}

func (e *DXEndpointHandle) GetFeed() (*DXFeedHandle, error) {
	// The result, an error included, is kept: a later call returns the same one.
	e.feedOnce.Do(func() {
		var ptr *C.dxfg_v4_feed_t
		e.feedErr = executeInIsolateThread(func(thread *isolateThread) error {
			return checkStatus(thread, C.dxfg_v4_DXEndpoint_getFeed(thread.ptr, e.ptr(), &ptr))
		})
		e.feed = NewDXFeedHandle(ptr)
	})

	return e.feed, e.feedErr
}

func (e *DXEndpointHandle) GetPublisher() (*DXPublisherHandle, error) {
	// The result, an error included, is kept: a later call returns the same one.
	e.publisherOnce.Do(func() {
		var ptr *C.dxfg_v4_publisher_t
		e.publisherErr = executeInIsolateThread(func(thread *isolateThread) error {
			return checkStatus(thread, C.dxfg_v4_DXEndpoint_getPublisher(thread.ptr, e.ptr(), &ptr))
		})
		e.publisher = NewDXPublisherHandle(ptr)
	})

	return e.publisher, e.publisherErr
}

//export OnStateChanged
func OnStateChanged(thread *C.graal_isolatethread_t, old C.dxfg_v4_endpoint_state_t, new C.dxfg_v4_endpoint_state_t, userData unsafe.Pointer) {
	Restore(userData).(common.ConnectionStateListener).UpdateState(common.ConnectionState(old), common.ConnectionState(new))
}

func (e *DXEndpointHandle) AttachListener(listener common.ConnectionStateListener) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		var l *C.dxfg_v4_endpoint_state_change_listener_t
		status := C.dxfg_v4_PropertyChangeListener_new(thread.ptr, (*[0]byte)(C.OnStateChanged), Save(listener), &l)
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_addStateChangeListener(thread.ptr, e.ptr(), l))
	})
}

func (e *DXEndpointHandle) Free() error {
	return errors.Join(e.feed.Free(), e.publisher.Free(), e.self.Free())
}

func (e *DXEndpointHandle) ptr() *C.dxfg_v4_endpoint_t {
	return (*C.dxfg_v4_endpoint_t)(e.self.Ptr())
}
