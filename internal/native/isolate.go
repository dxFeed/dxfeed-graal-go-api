// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#cgo CFLAGS: -I${SRCDIR}/include -Werror=deprecated-declarations
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"unsafe"

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

// ErrIsolateCreated is returned by SetIsolateArguments after the creation of the isolate has started.
var ErrIsolateCreated = errors.New("native: the isolate is already created, its arguments cannot be changed")

var (
	isolateArgsMu sync.Mutex
	isolateArgs   []string
	// isolateCreating is set under isolateArgsMu when the creation starts; the arguments are fixed from then on.
	isolateCreating bool
)

// SetIsolateArguments stores the arguments that the isolate is created with. It does not create the isolate.
func SetIsolateArguments(args []string) error {
	for i, arg := range args {
		if strings.IndexByte(arg, 0) >= 0 {
			return fmt.Errorf("native: the isolate argument #%d contains a NUL character", i)
		}
	}
	isolateArgsMu.Lock()
	defer isolateArgsMu.Unlock()
	if isolateCreating {
		return ErrIsolateCreated
	}
	isolateArgs = append([]string(nil), args...)
	return nil
}

// takeIsolateArguments fixes the arguments and returns them.
func takeIsolateArguments() []string {
	isolateArgsMu.Lock()
	defer isolateArgsMu.Unlock()
	isolateCreating = true
	return isolateArgs
}

// createIsolate calls graal_create_isolate with the stored arguments. Without arguments it passes no parameters, so
// the SDK uses its defaults.
func createIsolate(isolatePtr **C.graal_isolate_t, thread **C.graal_isolatethread_t) C.int {
	args := takeIsolateArguments()
	if len(args) == 0 {
		return C.graal_create_isolate(nil, isolatePtr, thread)
	}
	// The SDK skips the first argument: it is the name of the program.
	argv := make([]*C.char, 0, len(args)+1)
	argv = append(argv, C.CString("dxfeed-graal-go-api"))
	for _, arg := range args {
		argv = append(argv, C.CString(arg))
	}
	defer func() {
		for _, ptr := range argv {
			C.free(unsafe.Pointer(ptr))
		}
	}()
	// The array of the C strings is in C memory as well; the isolate reads it only during the call.
	cArgv := (**C.char)(C.malloc(C.size_t(len(argv)) * C.size_t(unsafe.Sizeof(argv[0]))))
	defer C.free(unsafe.Pointer(cArgv))
	copy(unsafe.Slice(cArgv, len(argv)), argv)
	var params C.graal_create_isolate_params_t
	params.version = C.__graal_create_isolate_params_version
	params.argc = C.int(len(argv))
	params.argv = cArgv
	// An unknown option fails the creation instead of being silently dropped.
	params.ignore_unrecognized_args = 0
	return C.graal_create_isolate(&params, isolatePtr, thread)
}

// getOrCreateIsolate creates the isolate on the first call; a failure is returned by this call and every later one.
func getOrCreateIsolate() (*isolate, error) {
	isolateOnce.Do(func() {
		handleWriteSignals()
		// The thread that creates the isolate is attached to it; keep the goroutine on it for the version check.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		instance := &isolate{}
		var thread *C.graal_isolatethread_t
		isolateErr = checkIsolateCall(func() C.int {
			return createIsolate(&instance.ptr, &thread)
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
	ptr *C.graal_isolatethread_t
}

// executeInIsolateThread runs call on the OS thread of the goroutine, attached to the isolate. An error of creating
// the isolate or of attaching the thread is returned, otherwise the error of call.
//
// An OS thread is attached on its first call and stays attached until the process exits: Go keeps its OS threads and
// runs other goroutines on them, and the next call on the thread finds it attached. GraalVM keeps the attached threads
// in one list under a global lock, so graal_get_current_thread of a thread that is not attached and
// graal_detach_thread take time that grows with the number of attached threads; attaching and detaching on every call
// of many goroutines at once made the calls hang on Linux for minutes.
//
// Go ends an OS thread only when a goroutine locked to it with runtime.LockOSThread exits without unlocking; such a
// thread, if it called the API, leaves its IsolateThread in the isolate until the process exits.
func executeInIsolateThread(call func(thread *isolateThread) error) error {
	// The goroutine stays on the attached thread for the call.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	thread, err := currentThread()
	if err != nil {
		return err
	}
	return call(thread)
}

// isAttached reports whether the current OS thread is attached to the isolate. The goroutine must be locked to its OS
// thread.
func isAttached() bool {
	isolate, err := getOrCreateIsolate()
	return err == nil && C.graal_get_current_thread(isolate.ptr) != nil
}

// currentThread returns the IsolateThread of the current OS thread and attaches the thread if it is not attached. The
// goroutine must be locked to its OS thread.
func currentThread() (*isolateThread, error) {
	isolate, err := getOrCreateIsolate()
	if err != nil {
		return nil, err
	}
	thread := &isolateThread{ptr: C.graal_get_current_thread(isolate.ptr)}
	if thread.ptr == nil {
		if err := checkIsolateCall(func() C.int {
			return C.graal_attach_thread(isolate.ptr, &thread.ptr)
		}); err != nil {
			return nil, err
		}
	}
	return thread, nil
}
