package isolate

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "dxfeed-graal-go-api/pkg/api"

type Endpoint struct {
	endpoint *C.dxfg_endpoint_t
}

func NewEndpoint(role api.Role) *Endpoint {
	thread := GetInstance().CurrentThread()
	e := &Endpoint{}
	e.endpoint = (*C.dxfg_endpoint_t)(C.malloc(8))
	e.endpoint = C.dxfg_DXEndpoint_create2(thread, (C.dxfg_endpoint_role_t)(role))
	return e
}

func (e *Endpoint) Connect(address string) {
	thread := GetInstance().CurrentThread()
	C.dxfg_DXEndpoint_connect(thread, e.endpoint, C.CString(address))
}
