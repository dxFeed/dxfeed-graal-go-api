package native

/*
#cgo CFLAGS: -I${SRCDIR}/graal
#cgo LDFLAGS: -L${SRCDIR}/graal -lDxFeedGraalNativeSdk
#include "dxfg_api.h"
*/
import "C"
import "sync"

type isolate struct {
	ptr *C.graal_isolate_t
}

var (
	isolateInstance *isolate
	isolateOnce     sync.Once
)

func getOrCreateIsolate() *isolate {
	isolateOnce.Do(func() {
		isolateInstance = &isolate{}
		C.graal_create_isolate(nil, &isolateInstance.ptr, nil)
	})
	return isolateInstance
}
