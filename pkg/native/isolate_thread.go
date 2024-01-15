package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

type isolateThread struct {
	ptr *C.graal_isolatethread_t
}

type isolateThreadOperation func(thread *isolateThread)

func executeWithIsolateThread(operation isolateThreadOperation) {
	thread := attachIsolateThread()
	defer thread.detachIsolateThread()
	operation(thread)
}

func attachIsolateThread() *isolateThread {
	isolate := getOrCreateIsolate()
	thread := &isolateThread{}
	C.graal_attach_thread(isolate.ptr, &thread.ptr)
	return thread
}

func (t *isolateThread) detachIsolateThread() {
	C.graal_detach_thread(t.ptr)
}
