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
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// IndexedTxModelOptions are the settings of a new IndexedTxModel.
type IndexedTxModelOptions struct {
	TxModelOptions
	// Sources are the sources of the events; empty means all the sources.
	Sources []events.IndexedEventSourceInterface
}

// IndexedTxModelHandle is a native IndexedTxModel with its listener.
type IndexedTxModelHandle struct {
	*txModelHandle
}

type indexedBuilder = C.dxfg_v4_indexed_tx_model_builder_t

// NewIndexedTxModel builds a native IndexedTxModel with the options.
func NewIndexedTxModel(o IndexedTxModelOptions) (*IndexedTxModelHandle, error) {
	h, err := newTxModel(o.TxModelOptions, txModelSteps[indexedBuilder]{
		newBuilder: func(thread *isolateThread, eventType C.dxfg_v4_event_clazz_t, out **indexedBuilder) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_newBuilder(thread.ptr, eventType, out)
		},
		withSymbol: func(thread *isolateThread, cur *indexedBuilder, symbol *C.dxfg_v4_symbol_t,
			out **indexedBuilder) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withSymbol(thread.ptr, cur, symbol, out)
		},
		withBatch: func(thread *isolateThread, cur *indexedBuilder, value C.int32_t, out **indexedBuilder) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withBatchProcessing(thread.ptr, cur, value, out)
		},
		withSnapshot: func(thread *isolateThread, cur *indexedBuilder, value C.int32_t, out **indexedBuilder) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withSnapshotProcessing(thread.ptr, cur, value, out)
		},
		withListener: func(thread *isolateThread, cur *indexedBuilder, listener *C.dxfg_v4_tx_model_listener_t,
			out **indexedBuilder) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withListener(thread.ptr, cur, listener, out)
		},
		withExecutor: func(thread *isolateThread, cur *indexedBuilder, executor *C.dxfg_v4_executor_t,
			out **indexedBuilder) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withExecutor(thread.ptr, cur, executor, out)
		},
		withFeed: func(thread *isolateThread, cur *indexedBuilder, feed *C.dxfg_v4_feed_t,
			out **indexedBuilder) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withFeed(thread.ptr, cur, feed, out)
		},
		build: func(thread *isolateThread, cur *indexedBuilder) (unsafe.Pointer, C.int32_t) {
			var model *C.dxfg_v4_indexed_tx_model_t
			status := C.dxfg_v4_IndexedTxModel_Builder_build(thread.ptr, cur, &model)
			return unsafe.Pointer(model), status
		},
		configure: func(b *txModelBuilder[indexedBuilder]) error {
			if len(o.Sources) == 0 {
				return nil
			}
			sources, err := newCIndexedEventSourceList(o.Sources)
			if err != nil {
				return err
			}
			defer freeCIndexedEventSourceList(sources)
			return b.step(func(cur *indexedBuilder, out **indexedBuilder) C.int32_t {
				return C.dxfg_v4_IndexedTxModel_Builder_withSources(b.thread.ptr, cur, sources, out)
			})
		},
		close: func(thread *isolateThread, model unsafe.Pointer) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_close(thread.ptr, (*C.dxfg_v4_indexed_tx_model_t)(model))
		},
	})
	if err != nil {
		return nil, err
	}
	return &IndexedTxModelHandle{h}, nil
}

// Attach attaches the model to the feed.
func (h *IndexedTxModelHandle) Attach(feed *DXFeedHandle) error {
	return h.withFeed(feed, func(thread *isolateThread, model unsafe.Pointer, f *C.dxfg_v4_feed_t) C.int32_t {
		return C.dxfg_v4_IndexedTxModel_attach(thread.ptr, (*C.dxfg_v4_indexed_tx_model_t)(model), f)
	})
}

// Detach detaches the model from the feed.
func (h *IndexedTxModelHandle) Detach(feed *DXFeedHandle) error {
	return h.withFeed(feed, func(thread *isolateThread, model unsafe.Pointer, f *C.dxfg_v4_feed_t) C.int32_t {
		return C.dxfg_v4_IndexedTxModel_detach(thread.ptr, (*C.dxfg_v4_indexed_tx_model_t)(model), f)
	})
}

// Sources returns the sources of the model; empty means all the sources.
func (h *IndexedTxModelHandle) Sources() ([]events.IndexedEventSourceInterface, error) {
	var result []events.IndexedEventSourceInterface
	err := h.call(func(thread *isolateThread, model unsafe.Pointer) error {
		var list *C.dxfg_v4_indexed_event_source_list
		if err := checkStatus(thread, C.dxfg_v4_IndexedTxModel_getSources(thread.ptr,
			(*C.dxfg_v4_indexed_tx_model_t)(model), &list)); err != nil {
			return err
		}
		var err error
		result, err = takeIndexedEventSourceList(thread, list)
		return err
	})
	return result, err
}

// SetSources sets the sources of the model; empty means all the sources.
func (h *IndexedTxModelHandle) SetSources(sources []events.IndexedEventSourceInterface) error {
	list, err := newCIndexedEventSourceList(sources)
	if err != nil {
		return err
	}
	defer freeCIndexedEventSourceList(list)
	return h.call(func(thread *isolateThread, model unsafe.Pointer) error {
		return checkStatus(thread, C.dxfg_v4_IndexedTxModel_setSources(thread.ptr,
			(*C.dxfg_v4_indexed_tx_model_t)(model), list))
	})
}

// newCIndexedEventSourceList builds a list of the sources in C memory; free it with freeCIndexedEventSourceList.
func newCIndexedEventSourceList(sources []events.IndexedEventSourceInterface) (*C.dxfg_v4_indexed_event_source_list,
	error) {
	list := (*C.dxfg_v4_indexed_event_source_list)(C.calloc(1,
		C.size_t(unsafe.Sizeof(C.dxfg_v4_indexed_event_source_list{}))))
	if len(sources) == 0 {
		return list, nil
	}
	list.elements = (**C.dxfg_indexed_event_source_t)(C.calloc(C.size_t(len(sources)),
		C.size_t(unsafe.Sizeof((*C.dxfg_indexed_event_source_t)(nil)))))
	elements := unsafe.Slice(list.elements, len(sources))
	for i, source := range sources {
		s, err := newCIndexedEventSource(source)
		if err != nil {
			freeCIndexedEventSourceList(list)
			return nil, err
		}
		elements[i] = (*C.dxfg_indexed_event_source_t)(unsafe.Pointer(s))
		list.size = C.int32_t(i + 1)
	}
	return list, nil
}

func freeCIndexedEventSourceList(list *C.dxfg_v4_indexed_event_source_list) {
	if list == nil {
		return
	}
	if list.elements != nil {
		for _, s := range unsafe.Slice(list.elements, int(list.size)) {
			freeCIndexedEventSource((*C.dxfg_v4_indexed_event_source_t)(unsafe.Pointer(s)))
		}
		C.free(unsafe.Pointer(list.elements))
	}
	C.free(unsafe.Pointer(list))
}

// takeIndexedEventSourceList copies the sources of a list that the SDK returned and releases the list.
func takeIndexedEventSourceList(thread *isolateThread,
	list *C.dxfg_v4_indexed_event_source_list) ([]events.IndexedEventSourceInterface, error) {
	if list == nil {
		return nil, nil
	}
	var result []events.IndexedEventSourceInterface
	var err error
	if list.size > 0 && list.elements != nil {
		for _, s := range unsafe.Slice(list.elements, int(list.size)) {
			source, sourceErr := goIndexedEventSource((*C.dxfg_v4_indexed_event_source_t)(unsafe.Pointer(s)))
			if sourceErr != nil {
				err = sourceErr
				break
			}
			result = append(result, source)
		}
	}
	err = errors.Join(err, checkStatus(thread, C.dxfg_v4_indexed_event_source_list_release(thread.ptr, list)))
	if err != nil {
		return nil, err
	}
	return result, nil
}
