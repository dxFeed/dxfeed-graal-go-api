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
	"sort"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

// PriceLevelServiceHandle holds a native client of the ORCS price level service.
type PriceLevelServiceHandle struct {
	self Handler
}

// OrderSourceSymbols are the symbols available for an order source.
type OrderSourceSymbols struct {
	Source  *order.Source
	Symbols []string
}

// AuthOrderSourceData is the content of an AuthOrderSource: the symbols by the order source identifier and by the
// order source. The symbols are sorted.
type AuthOrderSourceData struct {
	ByIDs     map[int32][]string
	BySources []OrderSourceSymbols
}

// NewPriceLevelServiceHandle creates a client of the price level service connected to the address; close it with
// Close and release it with Free.
func NewPriceLevelServiceHandle(address string) (*PriceLevelServiceHandle, error) {
	var service *C.dxfg_v4_price_level_service_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		cAddress := C.CString(address)
		defer C.free(unsafe.Pointer(cAddress))
		return checkStatus(thread, C.dxfg_v4_PriceLevelService_new(thread.ptr, cAddress, &service))
	})
	if err != nil {
		return nil, err
	}
	return &PriceLevelServiceHandle{self: NewJavaHandle(unsafe.Pointer(service))}, nil
}

var errPriceLevelServiceClosed = errors.New("native: the price level service is closed")

func (h *PriceLevelServiceHandle) ptr() *C.dxfg_v4_price_level_service_t {
	return (*C.dxfg_v4_price_level_service_t)(h.self.Ptr())
}

// Orders requests the orders of the candle symbol and the order source from from to to; a nil caller is the default
// caller of the SDK. It blocks while the request runs.
func (h *PriceLevelServiceHandle) Orders(symbol *candle.CandleSymbol, source *order.Source, from, to int64,
	caller *string) ([]interface{}, error) {
	ptr := h.ptr()
	if ptr == nil {
		return nil, errPriceLevelServiceClosed
	}
	if source == nil {
		return nil, errors.New("native: the order source is nil")
	}
	cSource, err := newCIndexedEventSource(source)
	if err != nil {
		return nil, err
	}
	defer freeCIndexedEventSource(cSource)
	return h.request(symbol, caller, func(thread *isolateThread, cSymbol *C.dxfg_v4_symbol_t, cCaller *C.char,
		list **C.dxfg_v4_event_type_list) C.int32_t {
		if cCaller == nil {
			return C.dxfg_v4_PriceLevelService_getOrders(thread.ptr, ptr, cSymbol, cSource, C.int64_t(from),
				C.int64_t(to), list)
		}
		return C.dxfg_v4_PriceLevelService_getOrdersWithCaller(thread.ptr, ptr, cSymbol, cSource, C.int64_t(from),
			C.int64_t(to), cCaller, list)
	})
}

// Quotes requests the quotes of the candle symbol from from to to; a nil caller is the default caller of the SDK. It
// blocks while the request runs.
func (h *PriceLevelServiceHandle) Quotes(symbol *candle.CandleSymbol, from, to int64, caller *string) ([]interface{},
	error) {
	ptr := h.ptr()
	if ptr == nil {
		return nil, errPriceLevelServiceClosed
	}
	return h.request(symbol, caller, func(thread *isolateThread, cSymbol *C.dxfg_v4_symbol_t, cCaller *C.char,
		list **C.dxfg_v4_event_type_list) C.int32_t {
		if cCaller == nil {
			return C.dxfg_v4_PriceLevelService_getQuotes(thread.ptr, ptr, cSymbol, C.int64_t(from), C.int64_t(to),
				list)
		}
		return C.dxfg_v4_PriceLevelService_getQuotesWithCaller(thread.ptr, ptr, cSymbol, C.int64_t(from),
			C.int64_t(to), cCaller, list)
	})
}

// request runs a request that returns a list of events with the symbol and the caller in C memory.
func (h *PriceLevelServiceHandle) request(symbol *candle.CandleSymbol, caller *string,
	call func(thread *isolateThread, cSymbol *C.dxfg_v4_symbol_t, cCaller *C.char,
		list **C.dxfg_v4_event_type_list) C.int32_t) ([]interface{}, error) {
	cSymbol, err := newCSymbol(symbol)
	if err != nil {
		return nil, err
	}
	defer freeCSymbol(cSymbol)
	var cCaller *C.char
	if caller != nil {
		cCaller = C.CString(*caller)
		defer C.free(unsafe.Pointer(cCaller))
	}
	var result []interface{}
	err = executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_event_type_list
		if err := checkStatus(thread, call(thread, cSymbol, cCaller, &list)); err != nil {
			return err
		}
		var err error
		result, err = takeEvents(thread, list)
		return err
	})
	return result, err
}

// AuthOrderSource requests the order sources and their symbols that are available to the client.
func (h *PriceLevelServiceHandle) AuthOrderSource() (AuthOrderSourceData, error) {
	ptr := h.ptr()
	if ptr == nil {
		return AuthOrderSourceData{}, errPriceLevelServiceClosed
	}
	var data AuthOrderSourceData
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var aos *C.dxfg_v4_auth_order_source_t
		if err := checkStatus(thread, C.dxfg_v4_PriceLevelService_getAuthOrderSource(thread.ptr, ptr, &aos)); err != nil {
			return err
		}
		var err error
		data, err = takeAuthOrderSource(thread, aos)
		return err
	})
	return data, err
}

// Close closes the connection of the client; the handle stays valid until Free. It does nothing after Free.
func (h *PriceLevelServiceHandle) Close() error {
	ptr := h.ptr()
	if ptr == nil {
		return nil
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_PriceLevelService_close(thread.ptr, ptr))
	})
}

// Free releases the handle once; later calls do nothing.
func (h *PriceLevelServiceHandle) Free() error {
	return h.self.Free()
}

// NewAuthOrderSource creates an AuthOrderSource from the symbols by the order source identifier and returns its
// content as the SDK reads it back.
func NewAuthOrderSource(byIDs map[int32][]string) (AuthOrderSourceData, error) {
	ids := make([]int32, 0, len(byIDs))
	for id := range byIDs {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	list := (*C.dxfg_v4_symbols_by_order_source_id_map_entry_list)(C.calloc(1,
		C.size_t(unsafe.Sizeof(C.dxfg_v4_symbols_by_order_source_id_map_entry_list{}))))
	defer freeIDEntryList(list)
	if len(ids) > 0 {
		var entry *C.dxfg_v4_symbols_by_order_source_id_map_entry_t
		list.elements = (**C.dxfg_v4_symbols_by_order_source_id_map_entry_t)(C.calloc(C.size_t(len(ids)),
			C.size_t(unsafe.Sizeof(entry))))
		list.size = C.int32_t(len(ids))
		for i, id := range ids {
			e := (*C.dxfg_v4_symbols_by_order_source_id_map_entry_t)(C.calloc(1,
				C.size_t(unsafe.Sizeof(C.dxfg_v4_symbols_by_order_source_id_map_entry_t{}))))
			e.order_source_id = C.int32_t(id)
			e.symbols = newStringList(byIDs[id])
			unsafe.Slice(list.elements, len(ids))[i] = e
		}
	}
	var data AuthOrderSourceData
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var aos *C.dxfg_v4_auth_order_source_t
		if err := checkStatus(thread, C.dxfg_v4_AuthOrderSource_create(thread.ptr, list, &aos)); err != nil {
			return err
		}
		var err error
		data, err = takeAuthOrderSource(thread, aos)
		return err
	})
	return data, err
}

// freeIDEntryList frees a list of entries built by NewAuthOrderSource, its entries and their symbols included.
func freeIDEntryList(list *C.dxfg_v4_symbols_by_order_source_id_map_entry_list) {
	if list.elements != nil {
		for _, e := range unsafe.Slice(list.elements, int(list.size)) {
			if e != nil {
				freeStringList(e.symbols)
				C.free(unsafe.Pointer(e))
			}
		}
		C.free(unsafe.Pointer(list.elements))
	}
	C.free(unsafe.Pointer(list))
}

// takeAuthOrderSource reads the content of an AuthOrderSource and releases its handle.
func takeAuthOrderSource(thread *isolateThread, aos *C.dxfg_v4_auth_order_source_t) (AuthOrderSourceData, error) {
	defer C.dxfg_v4_JavaObjectHandler_release(thread.ptr, (*C.dxfg_v4_java_object_handler)(unsafe.Pointer(aos)))
	data := AuthOrderSourceData{ByIDs: map[int32][]string{}}

	var byIDs *C.dxfg_v4_symbols_by_order_source_id_map_entry_list
	if err := checkStatus(thread, C.dxfg_v4_AuthOrderSource_getByIds(thread.ptr, aos, &byIDs)); err != nil {
		return AuthOrderSourceData{}, err
	}
	if byIDs != nil {
		defer C.dxfg_v4_symbols_by_order_source_id_map_entry_list_release(thread.ptr, byIDs)
		if byIDs.size > 0 && byIDs.elements != nil {
			for _, e := range unsafe.Slice(byIDs.elements, int(byIDs.size)) {
				if e != nil {
					data.ByIDs[int32(e.order_source_id)] = sortedStrings(e.symbols)
				}
			}
		}
	}

	var bySources *C.dxfg_v4_symbols_by_order_source_map_entry_list
	if err := checkStatus(thread, C.dxfg_v4_AuthOrderSource_getByOrderSources(thread.ptr, aos, &bySources)); err != nil {
		return AuthOrderSourceData{}, err
	}
	if bySources != nil {
		defer C.dxfg_v4_symbols_by_order_source_map_entry_list_release(thread.ptr, bySources)
		if bySources.size > 0 && bySources.elements != nil {
			for _, e := range unsafe.Slice(bySources.elements, int(bySources.size)) {
				if e == nil {
					continue
				}
				source, err := goIndexedEventSource((*C.dxfg_v4_indexed_event_source_t)(unsafe.Pointer(e.order_source)))
				if err != nil {
					return AuthOrderSourceData{}, err
				}
				orderSource, ok := source.(*order.Source)
				if !ok {
					return AuthOrderSourceData{}, fmt.Errorf("native: %v is not an order source", source)
				}
				data.BySources = append(data.BySources, OrderSourceSymbols{Source: orderSource,
					Symbols: sortedStrings(e.symbols)})
			}
		}
	}
	sort.Slice(data.BySources, func(i, j int) bool {
		return data.BySources[i].Source.Id() < data.BySources[j].Source.Id()
	})
	return data, nil
}

// sortedStrings copies the strings of a list that the SDK returned, sorted; never nil.
func sortedStrings(list *C.dxfg_string_list) []string {
	result := goStrings((*C.dxfg_v4_string_list)(unsafe.Pointer(list)))
	if result == nil {
		result = []string{}
	}
	sort.Strings(result)
	return result
}

// ValidatePriceLevels checks the orders with PriceLevelChecker: the bid is below the ask at each moment, and no gap
// between subsequent orders exceeds timeGapBound milliseconds.
func ValidatePriceLevels(orders []interface{}, timeGapBound int64, printQuotes bool) (bool, error) {
	if err := checkEvents("validate", orders); err != nil {
		return false, err
	}
	var valid C.int32_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_v4_event_type_list, interface{}](orders)
		// The SDK only reads the list; the caller keeps its ownership.
		defer l.Free()
		print := C.int32_t(0)
		if printQuotes {
			print = 1
		}
		return checkStatus(thread, C.dxfg_v4_PriceLevelChecker_validate(thread.ptr,
			(*C.dxfg_v4_event_type_list)(l.Ptr()), C.int64_t(timeGapBound), print, &valid))
	})
	return valid != 0, err
}
