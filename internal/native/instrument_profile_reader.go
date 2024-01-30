package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"unsafe"
)

type InstrumentProfileReader struct {
	handle Handler
}

func NewInstrumentProfileReader() (*InstrumentProfileReader, error) {
	var ptr *C.dxfg_instrument_profile_reader_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			ptr = C.dxfg_InstrumentProfileReader_new(thread.ptr)
		})
	})
	if err != nil {
		return nil, err
	}
	return &InstrumentProfileReader{handle: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

func (r *InstrumentProfileReader) GetLastModified() (int64, error) {
	var result int64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			result = int64(C.dxfg_InstrumentProfileReader_getLastModified(thread.ptr, r.ptr()))
		})
	})
	return result, err
}

func (r *InstrumentProfileReader) WasComplete() (bool, error) {
	var result bool
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			result = int32(C.dxfg_InstrumentProfileReader_wasComplete(thread.ptr, r.ptr())) == 1
		})
	})
	return result, err
}

func (r *InstrumentProfileReader) ReadFromFile(address string) ([]*events.InstrumentProfile, error) {
	var resultList []*events.InstrumentProfile

	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			ptr := C.dxfg_InstrumentProfileReader_readFromFile(thread.ptr,
				r.ptr(),
				C.CString(address))
			resultList = newProfileMapper().goProfiles(ptr)
			C.dxfg_CList_InstrumentProfile_release(thread.ptr, ptr)
		})
	})
	return resultList, err
}

func (r *InstrumentProfileReader) ReadFromFileWithPassword(address string, user string, password string) ([]*events.InstrumentProfile, error) {
	var resultList []*events.InstrumentProfile

	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			ptr := C.dxfg_InstrumentProfileReader_readFromFile2(thread.ptr,
				r.ptr(),
				C.CString(address),
				C.CString(user),
				C.CString(password))
			resultList = newProfileMapper().goProfiles(ptr)
			C.dxfg_CList_InstrumentProfile_release(thread.ptr, ptr)
		})
	})
	return resultList, err
}

func (r *InstrumentProfileReader) ReadDataWithAddress(data []byte, address string) ([]*events.InstrumentProfile, error) {
	var resultList []*events.InstrumentProfile

	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			inputStream := C.dxfg_ByteArrayInputStream_new(thread.ptr, (*C.char)(unsafe.Pointer(&data[0])), C.int32_t(len(data)))
			ptr := C.dxfg_InstrumentProfileReader_read2(thread.ptr,
				r.ptr(),
				inputStream,
				C.CString(address))
			resultList = newProfileMapper().goProfiles(ptr)
			C.dxfg_CList_InstrumentProfile_release(thread.ptr, ptr)
			C.dxfg_JavaObjectHandler_release(thread.ptr, (*C.dxfg_java_object_handler)(unsafe.Pointer(inputStream)))
		})
	})
	return resultList, err
}

func (r *InstrumentProfileReader) ReadCompressedData(data []byte) ([]*events.InstrumentProfile, error) {
	var resultList []*events.InstrumentProfile

	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			inputStream := C.dxfg_ByteArrayInputStream_new(thread.ptr, (*C.char)(unsafe.Pointer(&data[0])), C.int32_t(len(data)))
			ptr := C.dxfg_InstrumentProfileReader_readCompressed(thread.ptr,
				r.ptr(),
				inputStream)
			resultList = newProfileMapper().goProfiles(ptr)
			C.dxfg_CList_InstrumentProfile_release(thread.ptr, ptr)
			C.dxfg_JavaObjectHandler_release(thread.ptr, (*C.dxfg_java_object_handler)(unsafe.Pointer(inputStream)))
		})
	})
	return resultList, err
}

func ResolveSourceURL(address string) (*string, error) {
	var result *string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			value := C.dxfg_InstrumentProfileReader_resolveSourceURL(thread.ptr, C.CString(address))
			if value == nil {
				result = nil
			} else {
				temp := C.GoString(value)
				result = &temp
			}
		})
	})
	return result, err
}

func (r *InstrumentProfileReader) ReadData(data []byte) ([]*events.InstrumentProfile, error) {
	var resultList []*events.InstrumentProfile

	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			inputStream := C.dxfg_ByteArrayInputStream_new(thread.ptr, (*C.char)(unsafe.Pointer(&data[0])), C.int32_t(len(data)))
			ptr := C.dxfg_InstrumentProfileReader_read(thread.ptr,
				r.ptr(),
				inputStream)
			resultList = newProfileMapper().goProfiles(ptr)
			C.dxfg_CList_InstrumentProfile_release(thread.ptr, ptr)
			C.dxfg_JavaObjectHandler_release(thread.ptr, (*C.dxfg_java_object_handler)(unsafe.Pointer(inputStream)))
		})
	})
	return resultList, err
}

func (r *InstrumentProfileReader) ptr() *C.dxfg_instrument_profile_reader_t {
	return (*C.dxfg_instrument_profile_reader_t)(r.handle.Ptr())
}
