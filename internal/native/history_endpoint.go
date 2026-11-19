package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// HistoryEndpointSettings are the settings of a history endpoint; a nil string is not set.
type HistoryEndpointSettings struct {
	Address, UserName, Password, AuthToken *string
	// Compression and Format are the values of dxfg_history_endpoint_compression_t and dxfg_history_endpoint_format_t.
	Compression, Format int32
}

// HistoryEndpointHandle holds a native history endpoint.
type HistoryEndpointHandle struct {
	self Handler
}

// NewHistoryEndpointHandle builds a history endpoint with the settings; release it with Free.
func NewHistoryEndpointHandle(settings HistoryEndpointSettings) (*HistoryEndpointHandle, error) {
	var endpoint *C.dxfg_v4_history_endpoint_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var builder *C.dxfg_v4_history_endpoint_builder_t
		if err := checkStatus(thread, C.dxfg_v4_HistoryEndpoint_newBuilder(thread.ptr, &builder)); err != nil {
			return err
		}
		defer C.dxfg_v4_JavaObjectHandler_release(thread.ptr, (*C.dxfg_v4_java_object_handler)(unsafe.Pointer(builder)))
		strings := []struct {
			value *string
			set   func(*C.char) C.int32_t
		}{
			{settings.Address, func(s *C.char) C.int32_t {
				return C.dxfg_v4_HistoryEndpoint_Builder_withAddress(thread.ptr, builder, s)
			}},
			{settings.UserName, func(s *C.char) C.int32_t {
				return C.dxfg_v4_HistoryEndpoint_Builder_withUserName(thread.ptr, builder, s)
			}},
			{settings.Password, func(s *C.char) C.int32_t {
				return C.dxfg_v4_HistoryEndpoint_Builder_withPassword(thread.ptr, builder, s)
			}},
			{settings.AuthToken, func(s *C.char) C.int32_t {
				return C.dxfg_v4_HistoryEndpoint_Builder_withAuthToken(thread.ptr, builder, s)
			}},
		}
		for _, s := range strings {
			if s.value == nil {
				continue
			}
			cs := C.CString(*s.value)
			status := s.set(cs)
			C.free(unsafe.Pointer(cs))
			if err := checkStatus(thread, status); err != nil {
				return err
			}
		}
		if err := checkStatus(thread, C.dxfg_v4_HistoryEndpoint_Builder_withCompression(thread.ptr, builder,
			C.dxfg_v4_history_endpoint_compression_t(settings.Compression))); err != nil {
			return err
		}
		if err := checkStatus(thread, C.dxfg_v4_HistoryEndpoint_Builder_withFormat(thread.ptr, builder,
			C.dxfg_v4_history_endpoint_format_t(settings.Format))); err != nil {
			return err
		}
		return checkStatus(thread, C.dxfg_v4_HistoryEndpoint_Builder_build(thread.ptr, builder, &endpoint))
	})
	if err != nil {
		return nil, err
	}
	return &HistoryEndpointHandle{self: NewJavaHandle(unsafe.Pointer(endpoint))}, nil
}

// TimeSeries requests the events of the class and the symbol within [from, to], in milliseconds since the Unix epoch.
// It blocks while the request runs. The symbol is a subscription symbol as newCSymbol takes it.
func (h *HistoryEndpointHandle) TimeSeries(clazz int32, symbol any, from, to int64) ([]interface{}, error) {
	ptr := (*C.dxfg_v4_history_endpoint_t)(h.self.Ptr())
	if ptr == nil {
		return nil, errors.New("native: the history endpoint is closed")
	}
	cSymbol, err := newCSymbol(symbol)
	if err != nil {
		return nil, err
	}
	defer freeCSymbol(cSymbol)
	var result []interface{}
	err = executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_event_type_list
		status := C.dxfg_v4_HistoryEndpoint_getTimeSeries(thread.ptr, ptr, C.dxfg_v4_event_clazz_t(clazz), cSymbol,
			C.int64_t(from), C.int64_t(to), &list)
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		if list == nil {
			return nil
		}
		defer C.dxfg_v4_event_type_list_release(thread.ptr, list)
		result = newEventMapper().goEvents(list)
		return nil
	})
	return result, err
}

// Free releases the endpoint once; later calls do nothing.
func (h *HistoryEndpointHandle) Free() error {
	return h.self.Free()
}
