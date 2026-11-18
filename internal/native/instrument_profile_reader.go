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
	var ptr *C.dxfg_v4_instrument_profile_reader_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileReader_new(thread.ptr, &ptr))
	})
	if err != nil {
		return nil, err
	}
	return &InstrumentProfileReader{handle: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

func ResolveSourceURL(address string) (*string, error) {
	var result *string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		addressPtr := C.CString(address)
		defer C.free(unsafe.Pointer(addressPtr))
		var value *C.char
		if err := checkStatus(thread, C.dxfg_v4_InstrumentProfileReader_resolveSourceURL(thread.ptr, addressPtr, &value)); err != nil {
			return err
		}
		if value != nil {
			defer C.dxfg_v4_String_release(thread.ptr, value)
			temp := C.GoString(value)
			result = &temp
		}
		return nil
	})
	return result, err
}

func (r *InstrumentProfileReader) Close() error {
	return r.handle.Free()
}

func (r *InstrumentProfileReader) GetLastModified() (int64, error) {
	var result C.int64_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileReader_getLastModified(thread.ptr, r.ptr(), &result))
	})
	return int64(result), err
}

func (r *InstrumentProfileReader) WasComplete() (bool, error) {
	var result C.int32_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileReader_wasComplete(thread.ptr, r.ptr(), &result))
	})
	return result == 1, err
}

func (r *InstrumentProfileReader) ReadFromFile(address string) ([]*events.InstrumentProfile, error) {
	var resultList []*events.InstrumentProfile

	err := executeInIsolateThread(func(thread *isolateThread) error {
		addressPtr := C.CString(address)
		defer C.free(unsafe.Pointer(addressPtr))

		var nativeProfiles *C.dxfg_v4_instrument_profile_data_list

		err := checkStatus(thread, C.dxfg_v4_InstrumentProfileReader_readFromFile(thread.ptr, r.ptr(), addressPtr,
			&nativeProfiles))

		if nativeProfiles != nil {
			defer C.dxfg_v4_instrument_profile_data_list_release(thread.ptr, nativeProfiles)
		}

		if err != nil {
			return err
		}

		resultList, err = newProfileMapper().goProfiles(
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

		var nativeProfiles *C.dxfg_v4_instrument_profile_data_list

		err := checkStatus(thread, C.dxfg_v4_InstrumentProfileReader_readFromFileWithUserAndPassword(thread.ptr, r.ptr(),
			addressPtr, userPtr, passwordPtr, &nativeProfiles))

		if nativeProfiles != nil {
			defer C.dxfg_v4_instrument_profile_data_list_release(thread.ptr, nativeProfiles)
		}

		if err != nil {
			return err
		}

		resultList, err = newProfileMapper().goProfiles(
			thread,
			nativeProfiles,
		)
		return err
	})

	return resultList, err
}

func (r *InstrumentProfileReader) ptr() *C.dxfg_v4_instrument_profile_reader_t {
	return (*C.dxfg_v4_instrument_profile_reader_t)(r.handle.Ptr())
}
