package native

/*
#cgo CFLAGS: -I${SRCDIR}/include -Werror=deprecated-declarations
#include "dxfg_api.h"
*/
import "C"
import (
	"errors"
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
	// isolateErr is the failure of the isolate creation or of the SDK version check; every later call reports it again.
	isolateErr error
)

// getOrCreateIsolate creates the isolate on the first call; a failure is returned by this call and every later one.
func getOrCreateIsolate() (*isolate, error) {
	isolateOnce.Do(func() {
		// The thread that creates the isolate is attached to it; keep the goroutine on it for the version check.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		instance := &isolate{}
		var thread *C.graal_isolatethread_t
		isolateErr = checkIsolateCall(func() C.int {
			return C.graal_create_isolate(nil, &instance.ptr, &thread)
		})
		if isolateErr == nil {
			isolateErr = checkSdkVersion(thread)
		}
		isolateInstance = instance
	})
	if isolateErr != nil {
		return nil, isolateErr
	}
	return isolateInstance, nil
}

// checkSdkVersion compares the release of the loaded native library with the one whose headers the API is built with.
// It runs on the thread that created the isolate, before getOrCreateIsolate returns, so it calls the SDK directly.
// A library older than 3.13.0 lacks dxfg_v4_system_get_version, so linking or loading it fails before this check.
func checkSdkVersion(thread *C.graal_isolatethread_t) error {
	var version *C.char
	if C.dxfg_v4_system_get_version(thread, &version) != dxfgExecuteSuccessfully {
		_ = takeJavaException(thread)
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

// executeInIsolateThread runs call on an OS thread attached to the isolate. An error of creating the isolate or of
// attaching or detaching the thread is returned, the error of call joined with a detach failure.
func executeInIsolateThread(call func(thread *isolateThread) error) (err error) {
	thread, err := attachCurrentThread()
	if err != nil {
		return err
	}
	defer func() {
		// The error of call is kept as it is unless detaching fails too.
		if detachErr := thread.detach(); detachErr != nil {
			err = errors.Join(err, detachErr)
		}
	}()
	return call(thread)
}

func attachCurrentThread() (*isolateThread, error) {
	runtime.LockOSThread()
	isolate, err := getOrCreateIsolate()
	if err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}
	thread := &isolateThread{ptr: C.graal_get_current_thread(isolate.ptr), shouldDetach: false}
	if thread.ptr == nil {
		err := checkIsolateCall(func() C.int {
			return C.graal_attach_thread(isolate.ptr, &thread.ptr)
		})
		if err != nil {
			runtime.UnlockOSThread()
			return nil, err
		}
		thread.shouldDetach = true
	}
	return thread, nil
}

func (t *isolateThread) detach() error {
	defer runtime.UnlockOSThread()
	var err error
	if t.ptr != nil && t.shouldDetach {
		err = checkIsolateCall(func() C.int {
			return C.graal_detach_thread(t.ptr)
		})
	}
	t.ptr = nil
	return err
}
