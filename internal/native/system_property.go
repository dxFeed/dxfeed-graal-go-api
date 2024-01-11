package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func SetSystemProperty(key, value string) {
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		keyPtr := C.CString(key)
		defer C.free(unsafe.Pointer(keyPtr))

		valuePtr := C.CString(value)
		defer C.free(unsafe.Pointer(valuePtr))

		C.dxfg_system_set_property(thread.ptr, keyPtr, valuePtr)
		return nil
	})
}

func GetSystemProperty(key string) string {
	value := ""
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		keyPtr := C.CString(key)
		defer C.free(unsafe.Pointer(keyPtr))

		valuePtr := C.dxfg_system_get_property(thread.ptr, keyPtr)
		defer C.dxfg_String_release(thread.ptr, valuePtr)

		value = C.GoString(valuePtr)
		return nil
	})
	return value
}
