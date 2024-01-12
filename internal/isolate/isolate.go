package isolate

/*
#cgo CFLAGS: -I/${SRCDIR}/libs
#cgo LDFLAGS: -L/${SRCDIR}/bin -lDxFeedGraalNativeSdk
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "sync"

type Isolate struct {
	isolate *C.graal_isolate_t
}

var (
	instance *Isolate
	once     sync.Once
)

func GetInstance() *Isolate {
	once.Do(func() {
		instance = &Isolate{}
		instance.isolate = (*C.graal_isolate_t)(C.malloc(C.sizeof_uintptr_t))
		C.graal_create_isolate(nil, &instance.isolate, nil)
	})
	return instance
}

func (i *Isolate) CurrentThread() *C.graal_isolatethread_t {
	thread := (*C.graal_isolatethread_t)(C.malloc(C.sizeof_uintptr_t))
	C.graal_attach_thread(i.isolate, &thread)
	return thread
}
