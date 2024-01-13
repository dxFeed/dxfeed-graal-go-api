package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DXEndpoint struct {
	ptr *C.dxfg_endpoint_t
}

type Role int32

func NewEndpoint(role Role) *DXEndpoint {
	thread := attachIsolateThread()
	defer thread.detachIsolateThread()
	endpoint := &DXEndpoint{}
	endpoint.ptr = C.dxfg_DXEndpoint_create2(thread.ptr, (C.dxfg_endpoint_role_t)(role))
	return endpoint
}

func (e *DXEndpoint) Connect(address string) {
	thread := attachIsolateThread()
	defer thread.detachIsolateThread()
	addressPtr := C.CString(address)
	defer C.free(unsafe.Pointer(addressPtr))
	C.dxfg_DXEndpoint_connect(thread.ptr, e.ptr, addressPtr)
}

func (e *DXEndpoint) GetFeed() *DXFeed {
	thread := attachIsolateThread()
	defer thread.detachIsolateThread()
	feed := &DXFeed{}
	feed.ptr = C.dxfg_DXEndpoint_getFeed(thread.ptr, e.ptr)
	return feed
}
