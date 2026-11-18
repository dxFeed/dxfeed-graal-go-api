package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// ParseSymbols splits a list of symbols as the SDK tools do; the order of the result is not defined.
func ParseSymbols(symbols string) ([]string, error) {
	var result []string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		symbolsPtr := C.CString(symbols)
		defer C.free(unsafe.Pointer(symbolsPtr))

		var list *C.dxfg_v4_string_list
		if err := checkStatus(thread, C.dxfg_v4_Tools_parseSymbols(thread.ptr, symbolsPtr, &list)); err != nil {
			return err
		}
		defer releaseStringList(thread, list)

		result = goStrings(list)
		return nil
	})

	return result, err
}

func ParseTime(time string) (int64, error) {
	var result C.int64_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		timePtr := C.CString(time)
		defer C.free(unsafe.Pointer(timePtr))

		var format *C.dxfg_v4_time_format_t
		if err := checkStatus(thread, C.dxfg_v4_TimeFormat_DEFAULT(thread.ptr, &format)); err != nil {
			return err
		}
		defer C.dxfg_v4_JavaObjectHandler_release(thread.ptr, (*C.dxfg_v4_java_object_handler)(unsafe.Pointer(format)))

		return checkStatus(thread, C.dxfg_v4_TimeFormat_parse(thread.ptr, format, timePtr, &result))
	})

	return int64(result), err
}
