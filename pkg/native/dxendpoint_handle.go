package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DXEndpointHandle struct {
	ptr *C.dxfg_endpoint_t
}

type Role int32

func NewDXEndpointHandle(role Role) *DXEndpointHandle {
	endpoint := &DXEndpointHandle{}
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		ptr, _ := checkCall(C.dxfg_DXEndpoint_create2(thread.ptr, (C.dxfg_endpoint_role_t)(role)))
		endpoint.ptr = ptr
		return nil
	})
	return endpoint
}

func (e *DXEndpointHandle) Connect(address string) {
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		addressPtr := C.CString(address)
		defer C.free(unsafe.Pointer(addressPtr))
		C.dxfg_DXEndpoint_connect(thread.ptr, e.ptr, addressPtr)
		return nil
	})
}

func (e *DXEndpointHandle) GetFeed() *DXFeedHandle {
	feed := &DXFeedHandle{}
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		feed.ptr = C.dxfg_DXEndpoint_getFeed(thread.ptr, e.ptr)
		return nil
	})
	return feed
}
