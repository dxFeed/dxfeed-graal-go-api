package native

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include "dxfg_api.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/sdk"
)

const (
	dxfgExecuteSuccessfully = 0
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
		var thread *C.graal_isolatethread_t
		err := checkIsolateCall(func() C.int {
			return C.graal_create_isolate(nil, &isolateInstance.ptr, &thread)
		})
		if err != nil {
			panic(err)
		}
		if err := checkSdkVersion(thread); err != nil {
			panic(err)
		}
	})
	return isolateInstance
}

// checkSdkVersion compares the release of the loaded native library with the one whose headers the API is built with.
// It runs on the thread that created the isolate, before getOrCreateIsolate returns, so it calls the SDK directly.
func checkSdkVersion(thread *C.graal_isolatethread_t) error {
	var version *C.char
	if C.dxfg_v4_system_get_version(thread, &version) != dxfgExecuteSuccessfully {
		if ex := C.dxfg_get_and_clear_thread_exception_t(thread); ex != nil {
			C.dxfg_Exception_release(thread, ex)
		}
		return fmt.Errorf("native: cannot read the version of the dxFeed Graal Native SDK library; this API needs %s",
			sdk.Version)
	}
	defer C.dxfg_v4_String_release(thread, version)
	return sdk.CheckVersion(C.GoString(version))
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
