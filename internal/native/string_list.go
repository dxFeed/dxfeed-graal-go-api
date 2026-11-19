// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// goStrings copies the strings of a list that the SDK returned; nil for a nil or empty list. A NULL element becomes
// an empty string.
func goStrings(list *C.dxfg_v4_string_list) []string {
	if list == nil || list.elements == nil || list.size <= 0 {
		return nil
	}
	elements := unsafe.Slice(list.elements, int(list.size))
	result := make([]string, len(elements))
	for i, element := range elements {
		if element != nil {
			result[i] = C.GoString(element)
		}
	}
	return result
}

// releaseStringList releases a list that the SDK returned. A list built by newStringList is freed with
// freeStringList instead: the SDK never takes the ownership of a list passed to it.
func releaseStringList(thread *isolateThread, list *C.dxfg_v4_string_list) {
	if list != nil {
		C.dxfg_v4_string_list_release(thread.ptr, list)
	}
}

// newStringList builds a list in C memory to pass to the SDK; the caller frees it with freeStringList.
func newStringList(values []string) *C.dxfg_v4_string_list {
	list := (*C.dxfg_v4_string_list)(C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_string_list{}))))
	if len(values) == 0 {
		return list
	}
	var element *C.char
	elements := (**C.char)(C.calloc(C.size_t(len(values)), C.size_t(unsafe.Sizeof(element))))
	slice := unsafe.Slice(elements, len(values))
	for i, value := range values {
		slice[i] = C.CString(value)
	}
	list.size = C.int32_t(len(values))
	list.elements = elements
	return list
}

// freeStringList frees a list built by newStringList, its strings included.
func freeStringList(list *C.dxfg_v4_string_list) {
	if list == nil {
		return
	}
	if list.elements != nil {
		for _, element := range unsafe.Slice(list.elements, int(list.size)) {
			C.free(unsafe.Pointer(element))
		}
		C.free(unsafe.Pointer(list.elements))
	}
	C.free(unsafe.Pointer(list))
}
