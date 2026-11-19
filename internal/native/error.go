// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// checkStatus turns the status of a dxfg_v4 call into an error. A failed call leaves a Java exception on the isolate
// thread; it is taken right away, so that a later failure does not report it again.
func checkStatus(thread *isolateThread, status C.int32_t) error {
	if status == dxfgExecuteSuccessfully {
		return nil
	}
	if err := takeJavaException(thread.ptr); err != nil {
		return err
	}
	return fmt.Errorf("native: call failed with code %d without a pending java exception", int32(status))
}

// takeJavaException takes and releases the pending Java exception of the thread; nil if there is none.
func takeJavaException(thread *C.graal_isolatethread_t) error {
	var ptr *C.dxfg_v4_exception_t
	if C.dxfg_v4_get_and_clear_thread_exception(thread, &ptr) != dxfgExecuteSuccessfully || ptr == nil {
		return nil
	}
	defer C.dxfg_v4_Exception_release(thread, ptr)
	return javaError(ptr)
}

// javaError converts an exception of the SDK; the caller keeps and releases it.
func javaError(ptr *C.dxfg_v4_exception_t) error {
	return common.JavaError{
		ClassName:  C.GoString(ptr.class_name),
		Message:    C.GoString(ptr.message),
		StackTrace: C.GoString(ptr.print_stack_trace),
	}
}

func checkIsolateCall(call func() C.int) error {
	e := call()
	if !errors.Is((IsolateError)(e), NoError) {
		return (IsolateError)(e)
	}
	return nil
}

type IsolateError int32

const (
	NoError                                                 IsolateError = 0
	Unspecified                                             IsolateError = 1
	NullArgument                                            IsolateError = 2
	AllocationFailed                                        IsolateError = 3
	UnattachedThread                                        IsolateError = 4
	UninitializedIsolate                                    IsolateError = 5
	LocateImageFailed                                       IsolateError = 6
	OpenImageFailed                                         IsolateError = 7
	MapHeapFailed                                           IsolateError = 8
	ReserveAddressSpaceFailed                               IsolateError = 801
	InsufficientAddressSpace                                IsolateError = 802
	ProtectHeapFailed                                       IsolateError = 9
	UnsupportedIsolateParametersVersion                     IsolateError = 10
	ThreadingInitializationFailed                           IsolateError = 11
	UncaughtException                                       IsolateError = 12
	IsolateInitializationFailed                             IsolateError = 13
	OpenAuxImageFailed                                      IsolateError = 14
	ReadAuxImageMetaFailed                                  IsolateError = 15
	MapAuxImageFailed                                       IsolateError = 16
	InsufficientAuxImageMemory                              IsolateError = 17
	AuxImageUnsupported                                     IsolateError = 18
	FreeAddressSpaceFailed                                  IsolateError = 19
	FreeImageHeapFailed                                     IsolateError = 20
	AuxImagePrimaryImageMismatch                            IsolateError = 21
	ArgumentParsingFailed                                   IsolateError = 22
	CpuFeatureCheckFailed                                   IsolateError = 23
	PageSizeCheckFailed                                     IsolateError = 24
	DynamicMethodAddressResolutionGotFdCreateFailed         IsolateError = 25
	DynamicMethodAddressResolutionGotFdResizeFailed         IsolateError = 26
	DynamicMethodAddressResolutionGotFdMapFailed            IsolateError = 27
	DynamicMethodAddressResolutionGotMmapFailed             IsolateError = 28
	DynamicMethodAddressResolutionGotWrongMmap              IsolateError = 29
	DynamicMethodAddressResolutionGotFdInvalid              IsolateError = 30
	DynamicMethodAddressResolutionGotUniqueFileCreateFailed IsolateError = 31
	UnknownStackBoundaries                                  IsolateError = 32
)

func (e IsolateError) Error() string {
	return fmt.Sprintf("isolate: %s", e.String())
}

func (e IsolateError) String() string {
	switch {
	case errors.Is(e, NoError):
		return "No error occurred."
	case errors.Is(e, Unspecified):
		return "An unspecified error occurred."
	case errors.Is(e, NullArgument):
		return "An argument was NULL."
	case errors.Is(e, AllocationFailed):
		return "Memory allocation failed, the OS is probably out of memory."
	case errors.Is(e, UnattachedThread):
		return "The specified thread is not attached to the isolate."
	case errors.Is(e, UninitializedIsolate):
		return "The specified isolate is unknown."
	case errors.Is(e, LocateImageFailed):
		return "Locating the image file failed."
	case errors.Is(e, OpenImageFailed):
		return "Opening the located image file failed."
	case errors.Is(e, MapHeapFailed):
		return "Mapping the heap from the image file into memory failed."
	case errors.Is(e, ReserveAddressSpaceFailed):
		return "Reserving address space for the new isolate failed."
	case errors.Is(e, InsufficientAddressSpace):
		return "The image heap does not fit in the available address space."
	case errors.Is(e, ProtectHeapFailed):
		return "Setting the protection of the heap memory failed."
	case errors.Is(e, UnsupportedIsolateParametersVersion):
		return "The version of the specified isolate parameters is unsupported."
	case errors.Is(e, ThreadingInitializationFailed):
		return "Initialization of threading in the isolate failed."
	case errors.Is(e, UncaughtException):
		return "Some exception is not caught."
	case errors.Is(e, IsolateInitializationFailed):
		return "Initialization the isolate failed."
	case errors.Is(e, OpenAuxImageFailed):
		return "Opening the located auxiliary image file failed."
	case errors.Is(e, ReadAuxImageMetaFailed):
		return "Reading the opened auxiliary image file failed."
	case errors.Is(e, MapAuxImageFailed):
		return "Mapping the auxiliary image file into memory failed."
	case errors.Is(e, InsufficientAuxImageMemory):
		return "Insufficient memory for the auxiliary image."
	case errors.Is(e, AuxImageUnsupported):
		return "Auxiliary images are not supported on this platform or edition."
	case errors.Is(e, FreeAddressSpaceFailed):
		return "Releasing the isolate's address space failed."
	case errors.Is(e, FreeImageHeapFailed):
		return "Releasing the isolate's image heap memory failed."
	case errors.Is(e, AuxImagePrimaryImageMismatch):
		return "The auxiliary image was built from a different primary image."
	case errors.Is(e, ArgumentParsingFailed):
		return "The isolate arguments could not be parsed."
	case errors.Is(e, CpuFeatureCheckFailed):
		return "Current target does not support the CPU features that are required by the image."
	case errors.Is(e, PageSizeCheckFailed):
		return "Image page size is incompatible with run-time page size. " +
			"Rebuild image with -H:PageSize=[pagesize] to set appropriately."
	case errors.Is(e, DynamicMethodAddressResolutionGotFdCreateFailed):
		return "Creating an in-memory file for the GOT failed."
	case errors.Is(e, DynamicMethodAddressResolutionGotFdResizeFailed):
		return "Resizing the in-memory file for the GOT failed."
	case errors.Is(e, DynamicMethodAddressResolutionGotFdMapFailed):
		return "Mapping and populating the in-memory file for the GOT failed."
	case errors.Is(e, DynamicMethodAddressResolutionGotMmapFailed):
		return "Mapping the GOT before an isolate's heap failed (no mapping)."
	case errors.Is(e, DynamicMethodAddressResolutionGotWrongMmap):
		return "Mapping the GOT before an isolate's heap failed (wrong mapping)."
	case errors.Is(e, DynamicMethodAddressResolutionGotFdInvalid):
		return "Mapping the GOT before an isolate's heap failed (invalid file)."
	case errors.Is(e, DynamicMethodAddressResolutionGotUniqueFileCreateFailed):
		return "Could not create unique GOT file even after retrying."
	case errors.Is(e, UnknownStackBoundaries):
		return "Could not determine the stack boundaries."
	default:
		return "Unknown error."
	}
}
