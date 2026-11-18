package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func SetSystemProperty(key, value string) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		keyPtr := C.CString(key)
		defer C.free(unsafe.Pointer(keyPtr))

		valuePtr := C.CString(value)
		defer C.free(unsafe.Pointer(valuePtr))

		return checkStatus(thread, C.dxfg_v4_system_set_property(thread.ptr, keyPtr, valuePtr))
	})
}

func GetSystemProperty(key string) (string, error) {
	value := ""
	err := executeInIsolateThread(func(thread *isolateThread) error {
		keyPtr := C.CString(key)
		defer C.free(unsafe.Pointer(keyPtr))

		var valuePtr *C.char
		if err := checkStatus(thread, C.dxfg_v4_system_get_property(thread.ptr, keyPtr, &valuePtr)); err != nil {
			return err
		}
		if valuePtr != nil {
			defer C.dxfg_v4_system_release_property(thread.ptr, valuePtr)
			value = C.GoString(valuePtr)
		}
		return nil
	})
	return value, err
}
