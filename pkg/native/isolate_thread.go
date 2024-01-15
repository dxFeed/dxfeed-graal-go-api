package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

type isolateThread struct {
	ptr *C.graal_isolatethread_t
}

type isolateThreadOperation func(thread *isolateThread) error

func executeInIsolateThread(operation isolateThreadOperation) error {
	thread, detach := attachIsolateThread()
	defer detach()
	return operation(thread)
}

func attachIsolateThread() (*isolateThread, func()) {
	isolate := getOrCreateIsolate()
	thread := &isolateThread{ptr: C.graal_get_current_thread(isolate.ptr)}
	detach := func() {}
	if thread.ptr == nil {
		C.graal_attach_thread(isolate.ptr, &thread.ptr)
		detach = func() { thread.detachIsolateThread() }
	}
	return thread, detach
}

func (t *isolateThread) detachIsolateThread() {
	C.graal_detach_thread(t.ptr)
}
