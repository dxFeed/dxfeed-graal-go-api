// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// TimeSeriesTxModelOptions are the settings of a new TimeSeriesTxModel.
type TimeSeriesTxModelOptions struct {
	TxModelOptions
	// FromTime is the time in milliseconds since the epoch the model is subscribed from; math.MaxInt64 means not
	// subscribed.
	FromTime int64
}

// TimeSeriesTxModelHandle is a native TimeSeriesTxModel with its listener.
type TimeSeriesTxModelHandle struct {
	*txModelHandle
}

type timeSeriesBuilder = C.dxfg_v4_time_series_tx_model_builder_t

// NewTimeSeriesTxModel builds a native TimeSeriesTxModel with the options.
func NewTimeSeriesTxModel(o TimeSeriesTxModelOptions) (*TimeSeriesTxModelHandle, error) {
	h, err := newTxModel(o.TxModelOptions, txModelSteps[timeSeriesBuilder]{
		newBuilder: func(thread *isolateThread, eventType C.dxfg_v4_event_clazz_t, out **timeSeriesBuilder) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_newBuilder(thread.ptr, eventType, out)
		},
		withSymbol: func(thread *isolateThread, cur *timeSeriesBuilder, symbol *C.dxfg_v4_symbol_t,
			out **timeSeriesBuilder) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_Builder_withSymbol(thread.ptr, cur, symbol, out)
		},
		withBatch: func(thread *isolateThread, cur *timeSeriesBuilder, value C.int32_t,
			out **timeSeriesBuilder) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_Builder_withBatchProcessing(thread.ptr, cur, value, out)
		},
		withSnapshot: func(thread *isolateThread, cur *timeSeriesBuilder, value C.int32_t,
			out **timeSeriesBuilder) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_Builder_withSnapshotProcessing(thread.ptr, cur, value, out)
		},
		withListener: func(thread *isolateThread, cur *timeSeriesBuilder, listener *C.dxfg_v4_tx_model_listener_t,
			out **timeSeriesBuilder) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_Builder_withListener(thread.ptr, cur, listener, out)
		},
		withExecutor: func(thread *isolateThread, cur *timeSeriesBuilder, executor *C.dxfg_v4_executor_t,
			out **timeSeriesBuilder) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_Builder_withExecutor(thread.ptr, cur, executor, out)
		},
		withFeed: func(thread *isolateThread, cur *timeSeriesBuilder, feed *C.dxfg_v4_feed_t,
			out **timeSeriesBuilder) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_Builder_withFeed(thread.ptr, cur, feed, out)
		},
		build: func(thread *isolateThread, cur *timeSeriesBuilder) (unsafe.Pointer, C.int32_t) {
			var model *C.dxfg_v4_time_series_tx_model_t
			status := C.dxfg_v4_TimeSeriesTxModel_Builder_build(thread.ptr, cur, &model)
			return unsafe.Pointer(model), status
		},
		configure: func(b *txModelBuilder[timeSeriesBuilder]) error {
			return b.step(func(cur *timeSeriesBuilder, out **timeSeriesBuilder) C.int32_t {
				return C.dxfg_v4_TimeSeriesTxModel_Builder_withFromTime(b.thread.ptr, cur, C.int64_t(o.FromTime), out)
			})
		},
		close: func(thread *isolateThread, model unsafe.Pointer) C.int32_t {
			return C.dxfg_v4_TimeSeriesTxModel_close(thread.ptr, (*C.dxfg_v4_time_series_tx_model_t)(model))
		},
	})
	if err != nil {
		return nil, err
	}
	return &TimeSeriesTxModelHandle{h}, nil
}

// Attach attaches the model to the feed.
func (h *TimeSeriesTxModelHandle) Attach(feed *DXFeedHandle) error {
	return h.withFeed(feed, func(thread *isolateThread, model unsafe.Pointer, f *C.dxfg_v4_feed_t) C.int32_t {
		return C.dxfg_v4_TimeSeriesTxModel_attach(thread.ptr, (*C.dxfg_v4_time_series_tx_model_t)(model), f)
	})
}

// Detach detaches the model from the feed.
func (h *TimeSeriesTxModelHandle) Detach(feed *DXFeedHandle) error {
	return h.withFeed(feed, func(thread *isolateThread, model unsafe.Pointer, f *C.dxfg_v4_feed_t) C.int32_t {
		return C.dxfg_v4_TimeSeriesTxModel_detach(thread.ptr, (*C.dxfg_v4_time_series_tx_model_t)(model), f)
	})
}

// FromTime returns the time in milliseconds since the epoch the model is subscribed from; math.MaxInt64 means not
// subscribed.
func (h *TimeSeriesTxModelHandle) FromTime() (int64, error) {
	var fromTime C.int64_t
	err := h.call(func(thread *isolateThread, model unsafe.Pointer) error {
		return checkStatus(thread, C.dxfg_v4_TimeSeriesTxModel_getFromTime(thread.ptr,
			(*C.dxfg_v4_time_series_tx_model_t)(model), &fromTime))
	})
	return int64(fromTime), err
}

// SetFromTime sets the time in milliseconds since the epoch the model is subscribed from; math.MaxInt64 unsubscribes.
func (h *TimeSeriesTxModelHandle) SetFromTime(fromTime int64) error {
	return h.call(func(thread *isolateThread, model unsafe.Pointer) error {
		return checkStatus(thread, C.dxfg_v4_TimeSeriesTxModel_setFromTime(thread.ptr,
			(*C.dxfg_v4_time_series_tx_model_t)(model), C.int64_t(fromTime)))
	})
}
