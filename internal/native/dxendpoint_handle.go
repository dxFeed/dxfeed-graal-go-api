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
	"sync/atomic"
	"unsafe"
)

type DXEndpointHandle struct {
	self     Handler
	feedOnce sync.Once
	feed     *DXFeedHandle
	feedErr  error
	// feedMade is the feed once it is made, for Close, which must not wait for or start GetFeed.
	feedMade atomic.Pointer[DXFeedHandle]

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

// BuilderStep is one setting of an endpoint builder; the steps are applied in order, so a later one wins.
type BuilderStep interface {
	apply(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t
}

type builderStepFunc func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t

func (f builderStepFunc) apply(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t {
	return f(thread, builder)
}

// WithRole sets the role of the endpoint.
func WithRole(role common.Role) BuilderStep {
	return builderStepFunc(func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t {
		return C.dxfg_v4_DXEndpoint_Builder_withRole(thread.ptr, builder, (C.dxfg_v4_endpoint_role_t)(role))
	})
}

// WithName sets the name of the endpoint.
func WithName(name string) BuilderStep {
	return builderStepFunc(func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t {
		// The SDK copies the string during the call.
		namePtr := C.CString(name)
		defer C.free(unsafe.Pointer(namePtr))
		return C.dxfg_v4_DXEndpoint_Builder_withName(thread.ptr, builder, namePtr)
	})
}

// WithProperty sets a property of the endpoint.
func WithProperty(key, value string) BuilderStep {
	return builderStepFunc(func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t {
		// The SDK copies both strings during the call.
		keyPtr := C.CString(key)
		defer C.free(unsafe.Pointer(keyPtr))
		valuePtr := C.CString(value)
		defer C.free(unsafe.Pointer(valuePtr))
		return C.dxfg_v4_DXEndpoint_Builder_withProperty(thread.ptr, builder, keyPtr, valuePtr)
	})
}

// WithPropertiesFile sets the properties of the endpoint from a Java properties file.
func WithPropertiesFile(path string) BuilderStep {
	return builderStepFunc(func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t {
		pathPtr := C.CString(path)
		defer C.free(unsafe.Pointer(pathPtr))
		return C.dxfg_v4_DXEndpoint_Builder_withProperties(thread.ptr, builder, pathPtr)
	})
}

// WithExecutor sets the executor of the notifications of the endpoint; nil sets the default one.
func WithExecutor(executor *ExecutorHandle) BuilderStep {
	return builderStepFunc(func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) C.int32_t {
		return C.dxfg_v4_DXEndpoint_Builder_withExecutor(thread.ptr, builder, executor.ptr())
	})
}

// withBuilder runs f with a new native builder that has the steps applied, and releases the builder afterwards.
func withBuilder(steps []BuilderStep, f func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) error) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		var builder *C.dxfg_v4_endpoint_builder_t
		if err := checkStatus(thread, C.dxfg_v4_DXEndpoint_newBuilder(thread.ptr, &builder)); err != nil {
			return err
		}
		defer C.dxfg_v4_JavaObjectHandler_release(thread.ptr, (*C.dxfg_v4_java_object_handler)(unsafe.Pointer(builder)))
		for _, step := range steps {
			if err := checkStatus(thread, step.apply(thread, builder)); err != nil {
				return err
			}
		}
		return f(thread, builder)
	})
}

// BuildDXEndpointHandle builds an endpoint with a builder that has the steps applied.
func BuildDXEndpointHandle(steps ...BuilderStep) (*DXEndpointHandle, error) {
	var ptr *C.dxfg_v4_endpoint_t
	err := withBuilder(steps, func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_Builder_build(thread.ptr, builder, &ptr))
	})
	if err != nil {
		return nil, err
	}

	return &DXEndpointHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

// BuilderSupportsProperty reports whether a builder that has the steps applied supports the property.
func BuilderSupportsProperty(key string, steps ...BuilderStep) (bool, error) {
	var supported C.int32_t
	err := withBuilder(steps, func(thread *isolateThread, builder *C.dxfg_v4_endpoint_builder_t) error {
		keyPtr := C.CString(key)
		defer C.free(unsafe.Pointer(keyPtr))
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_Builder_supportsProperty(thread.ptr, builder, keyPtr, &supported))
	})
	return supported != 0, err
}

// Role returns the role of the endpoint.
func (e *DXEndpointHandle) Role() (common.Role, error) {
	var role C.dxfg_v4_endpoint_role_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_getRole(thread.ptr, e.ptr(), &role))
	})
	return common.Role(role), err
}

// State returns the state of the endpoint; Closed after Close, whose Free drops the handle.
func (e *DXEndpointHandle) State() (common.ConnectionState, error) {
	ptr := e.ptr()
	if ptr == nil {
		return common.Closed, nil
	}
	var state C.dxfg_v4_endpoint_state_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_getState(thread.ptr, ptr, &state))
	})
	return common.ConnectionState(state), err
}

// User sets the user name of the endpoint.
func (e *DXEndpointHandle) User(user string) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		userPtr := C.CString(user)
		defer C.free(unsafe.Pointer(userPtr))
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_user(thread.ptr, e.ptr(), userPtr))
	})
}

// Password sets the password of the endpoint.
func (e *DXEndpointHandle) Password(password string) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		passwordPtr := C.CString(password)
		defer C.free(unsafe.Pointer(passwordPtr))
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_password(thread.ptr, e.ptr(), passwordPtr))
	})
}

// EventTypes returns the event classes that the endpoint supports.
func (e *DXEndpointHandle) EventTypes() ([]int32, error) {
	var result []int32
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_event_clazz_list
		if err := checkStatus(thread, C.dxfg_v4_DXEndpoint_getEventTypes(thread.ptr, e.ptr(), &list)); err != nil {
			return err
		}
		result = takeEventClazzList(thread, list)
		return nil
	})
	return result, err
}

// RequestedAggregationPeriodMillis returns the requested aggregation period in milliseconds, negative if none.
func (e *DXEndpointHandle) RequestedAggregationPeriodMillis() (int64, error) {
	var period C.int64_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_getRequestedAggregationPeriodMillis(thread.ptr, e.ptr(), &period))
	})
	return int64(period), err
}

// SetRequestedAggregationPeriodMillis sets the requested aggregation period in milliseconds; negative clears it.
func (e *DXEndpointHandle) SetRequestedAggregationPeriodMillis(period int64) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_setRequestedAggregationPeriodMillis(thread.ptr, e.ptr(), C.int64_t(period)))
	})
}

// AggregationPeriodInfo returns the bounds of the aggregation period that the servers apply, in milliseconds;
// negative if unknown.
func (e *DXEndpointHandle) AggregationPeriodInfo() (minMillis, maxMillis int64, err error) {
	var lo, hi C.int64_t
	err = executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_getAggregationPeriodInfo(thread.ptr, e.ptr(), &lo, &hi))
	})
	return int64(lo), int64(hi), err
}

func (e *DXEndpointHandle) Close() error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_close(thread.ptr, e.ptr()))
	})
	e.cancelPromises()
	return errors.Join(err, e.Free())
}

func (e *DXEndpointHandle) CloseAndAwaitTermination() error {
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXEndpoint_closeAndAwaitTermination(thread.ptr, e.ptr()))
	})
	e.cancelPromises()
	return errors.Join(err, e.Free())
}

// cancelPromises cancels the pending promises of the feed of the closed endpoint.
func (e *DXEndpointHandle) cancelPromises() {
	if feed := e.feedMade.Load(); feed != nil {
		feed.promises.close()
	}
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
		e.feedMade.Store(e.feed)
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
	defer recoverCallback("connection state listener")
	if listener, ok := Restore(userData).(common.ConnectionStateListener); ok {
		listener.UpdateState(common.ConnectionState(old), common.ConnectionState(new))
	}
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
