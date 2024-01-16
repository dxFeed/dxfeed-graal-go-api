package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
)

type JavaError struct {
	ClassName  string
	Message    string
	StackTrace string
}

func (j JavaError) Error() string {
	return fmt.Sprintf("java: exception of type '%s'. %s \n %s", j.ClassName, j.Message, j.StackTrace)
}

func checkCall(call func()) error {
	call()
	return getJavaThreadErrorIfExist()
}

func getJavaThreadErrorIfExist() error {
	thread, detach := attachIsolateThread()
	defer detach()

	ptr := C.dxfg_get_and_clear_thread_exception_t(thread.ptr)
	if ptr == nil {
		return nil
	}
	defer C.dxfg_Exception_release(thread.ptr, ptr)

	return JavaError{
		ClassName:  C.GoString(ptr.class_name),
		Message:    C.GoString(ptr.message),
		StackTrace: C.GoString(ptr.print_stack_trace),
	}
}
