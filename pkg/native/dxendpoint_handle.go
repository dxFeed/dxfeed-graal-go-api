package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type DXEndpointHandle struct {
	ptr *C.dxfg_endpoint_t
}

type Role int32

func NewDXEndpointHandle(role Role) (*DXEndpointHandle, error) {
	var ptr *C.dxfg_endpoint_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			ptr = C.dxfg_DXEndpoint_create2(thread.ptr, (C.dxfg_endpoint_role_t)(role))
		})
	})
	if err != nil {
		return nil, err
	}

	return &DXEndpointHandle{ptr: ptr}, nil
}

func (e *DXEndpointHandle) Connect(address string) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		addressPtr := C.CString(address)
		defer C.free(unsafe.Pointer(addressPtr))
		return checkCall(func() {
			C.dxfg_DXEndpoint_connect(thread.ptr, e.ptr, addressPtr)
		})
	})
}

func (e *DXEndpointHandle) GetFeed() (*DXFeedHandle, error) {
	var ptr *C.dxfg_feed_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		ptr = C.dxfg_DXEndpoint_getFeed(thread.ptr, e.ptr)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &DXFeedHandle{ptr: ptr}, nil
}
