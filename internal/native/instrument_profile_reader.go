package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
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

func ResolveSourceURL(address string) (*string, error) {
	var result *string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			addressPtr := C.CString(address)
			defer C.free(unsafe.Pointer(addressPtr))
			value := C.dxfg_InstrumentProfileReader_resolveSourceURL(thread.ptr, addressPtr)
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

func (r *InstrumentProfileReader) Close() error {
	return r.handle.Free()
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
		addressPtr := C.CString(address)
		defer C.free(unsafe.Pointer(addressPtr))

		var nativeProfiles *C.dxfg_instrument_profile2_list_t

		err := checkResultCall(func() C.int32_t {
			return C.dxfg_InstrumentProfileReader_readFromFile7(thread.ptr,
				r.ptr(),
				addressPtr,
				&nativeProfiles)
		})

		if nativeProfiles != nil {
			defer C.dxfg_instrument_profile2_list_free(thread.ptr, nativeProfiles)
		}

		if err != nil {
			return err
		}

		resultList, err = newProfileMapper().goProfiles2(
			thread,
			nativeProfiles,
		)
		return err
	})

	return resultList, err
}

func (r *InstrumentProfileReader) ReadFromFileWithPassword(
	address string,
	user string,
	password string,
) ([]*events.InstrumentProfile, error) {
	var resultList []*events.InstrumentProfile

	err := executeInIsolateThread(func(thread *isolateThread) error {
		addressPtr := C.CString(address)
		userPtr := C.CString(user)
		passwordPtr := C.CString(password)

		defer C.free(unsafe.Pointer(addressPtr))
		defer C.free(unsafe.Pointer(userPtr))
		defer C.free(unsafe.Pointer(passwordPtr))

		var nativeProfiles *C.dxfg_instrument_profile2_list_t

		err := checkResultCall(func() C.int32_t {
			return C.dxfg_InstrumentProfileReader_readFromFile8(thread.ptr,
				r.ptr(),
				addressPtr,
				userPtr,
				passwordPtr,
				&nativeProfiles)
		})

		if nativeProfiles != nil {
			defer C.dxfg_instrument_profile2_list_free(thread.ptr, nativeProfiles)
		}

		if err != nil {
			return err
		}

		resultList, err = newProfileMapper().goProfiles2(
			thread,
			nativeProfiles,
		)
		return err
	})

	return resultList, err
}

func (r *InstrumentProfileReader) ptr() *C.dxfg_instrument_profile_reader_t {
	return (*C.dxfg_instrument_profile_reader_t)(r.handle.Ptr())
}
