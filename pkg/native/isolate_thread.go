package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
#include <time.h>
#include <memory.h>
typedef struct tm tm_t;
*/
import "C"

type IsolateThread struct {
	ptr *C.graal_isolatethread_t
}

func attachIsolateThread() *IsolateThread {
	isolate := getOrCreateIsolate()
	thread := &IsolateThread{}
	C.graal_attach_thread(isolate.ptr, &thread.ptr)
	return thread
}

func (t *IsolateThread) detachIsolateThread() {
	C.graal_detach_thread(t.ptr)
}
