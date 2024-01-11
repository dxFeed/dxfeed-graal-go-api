package native

/*
#cgo CFLAGS: -I${SRCDIR}/graal
#include "dxfg_api.h"
*/
import "C"
import (
	"runtime"
	"sync"
)

type isolate struct {
	ptr *C.graal_isolate_t
}

var (
	isolateOnce     sync.Once
	isolateInstance *isolate
)

func getOrCreateIsolate() *isolate {
	isolateOnce.Do(func() {
		isolateInstance = &isolate{}
		err := checkIsolateCall(func() C.int {
			return C.graal_create_isolate(nil, &isolateInstance.ptr, nil)
		})
		if err != nil {
			panic(err)
		}
	})
	return isolateInstance
}

type isolateThread struct {
	ptr          *C.graal_isolatethread_t
	shouldDetach bool
}

func executeInIsolateThread(call func(thread *isolateThread) error) error {
	thread := attachCurrentThread()
	defer thread.detach()
	return call(thread)
}

func attachCurrentThread() *isolateThread {
	runtime.LockOSThread()
	isolate := getOrCreateIsolate()
	thread := &isolateThread{ptr: C.graal_get_current_thread(isolate.ptr), shouldDetach: false}
	if thread.ptr == nil {
		err := checkIsolateCall(func() C.int {
			return C.graal_attach_thread(isolate.ptr, &thread.ptr)
		})
		if err != nil {
			panic(err)
		}
		thread.shouldDetach = true
	}
	return thread
}

func (t *isolateThread) detach() {
	defer runtime.UnlockOSThread()
	if t.ptr != nil && t.shouldDetach {
		err := checkIsolateCall(func() C.int {
			return C.graal_detach_thread(t.ptr)
		})
		if err != nil {
			panic(err)
		}
	}
	t.ptr = nil
}
