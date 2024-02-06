package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func ParseSymbols(symbols string) ([]any, error) {
	var result []any
	err := executeInIsolateThread(func(thread *isolateThread) error {
		symbolsPtr := C.CString(symbols)
		defer C.free(unsafe.Pointer(symbolsPtr))
		return checkCall(func() {
			resultPtr := C.dxfg_Tools_parseSymbols(thread.ptr, symbolsPtr)
			defer C.dxfg_CList_String_release(thread.ptr, resultPtr)

			if resultPtr == nil || resultPtr.elements == nil || int(resultPtr.size) == 0 {
				return
			}

			size := int(resultPtr.size)
			list := make([]any, size)
			elementsSlice := unsafe.Slice(resultPtr.elements, C.size_t(resultPtr.size))

			for i, event := range elementsSlice {
				symbol := C.GoString(event)
				list[i] = symbol
			}
			result = list
		})
	})

	return result, err
}

func ParseTime(time string) (int64, error) {
	var result int64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		timePtr := C.CString(time)
		defer C.free(unsafe.Pointer(timePtr))
		return checkCall(func() {
			defaultTimeFormat := C.dxfg_TimeFormat_DEFAULT(thread.ptr)
			defer C.dxfg_JavaObjectHandler_release(thread.ptr, (*C.dxfg_java_object_handler)(unsafe.Pointer(defaultTimeFormat)))
			result = int64(C.dxfg_TimeFormat_parse(thread.ptr, defaultTimeFormat, timePtr))
		})
	})

	return result, err
}
