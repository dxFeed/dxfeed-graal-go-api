// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnTxModelEventsReceived(graal_isolatethread_t *thread, dxfg_indexed_event_source_t *source,
	dxfg_event_type_list *events, int32_t isSnapshot, void *user_data);
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// TxModelListenerFunc receives the events of one or more complete transactions of a model, all of one source.
type TxModelListenerFunc func(source events.IndexedEventSourceInterface, events []interface{}, isSnapshot bool)

//export OnTxModelEventsReceived
func OnTxModelEventsReceived(thread *C.graal_isolatethread_t, source *C.dxfg_indexed_event_source_t,
	eventsList *C.dxfg_event_type_list, isSnapshot C.int32_t, userData unsafe.Pointer) {
	defer recoverCallback("transaction model listener")
	listener, ok := Restore(userData).(TxModelListenerFunc)
	if !ok {
		return
	}
	goSource, err := goIndexedEventSource((*C.dxfg_v4_indexed_event_source_t)(unsafe.Pointer(source)))
	if err != nil {
		logGo(LogLevelError, "transaction model listener: %v", err)
		return
	}
	goEvents := newEventMapper().goEvents((*C.dxfg_v4_event_type_list)(unsafe.Pointer(eventsList)))
	if goEvents == nil {
		goEvents = []interface{}{}
	}
	listener(goSource, goEvents, isSnapshot != 0)
}

// TxModelOptions are the settings common to the transaction models.
type TxModelOptions struct {
	// EventType is the native code of the event type; required.
	EventType int32
	// Symbol is the symbol of the model, as newCSymbol accepts it; required.
	Symbol any
	// Feed is the feed to attach the model to; nil leaves the model detached.
	Feed *DXFeedHandle
	// Listener receives the transactions; nil means none.
	Listener TxModelListenerFunc
	// BatchProcessing and SnapshotProcessing are passed to the builder as they are.
	BatchProcessing    bool
	SnapshotProcessing bool
	// Executor runs the listener; nil means the executor of the feed.
	Executor *ExecutorHandle
}

// txModelBuilder is a native builder of a transaction model. Each step replaces the builder with the one that the SDK
// returns and releases the previous one.
type txModelBuilder[B any] struct {
	thread *isolateThread
	cur    *B
}

// step passes the current builder to f and replaces it with the one that f returns.
func (b *txModelBuilder[B]) step(f func(cur *B, out **B) C.int32_t) error {
	var next *B
	if err := checkStatus(b.thread, f(b.cur, &next)); err != nil {
		return err
	}
	err := b.release()
	b.cur = next
	return err
}

func (b *txModelBuilder[B]) release() error {
	if b.cur == nil {
		return nil
	}
	err := checkStatus(b.thread, C.dxfg_v4_JavaObjectHandler_release(b.thread.ptr,
		(*C.dxfg_v4_java_object_handler)(unsafe.Pointer(b.cur))))
	b.cur = nil
	return err
}

// txModelSteps are the native calls of the builder of one kind of model. C functions are not values in Go, so each
// field wraps one.
type txModelSteps[B any] struct {
	newBuilder   func(thread *isolateThread, eventType C.dxfg_v4_event_clazz_t, out **B) C.int32_t
	withSymbol   func(thread *isolateThread, cur *B, symbol *C.dxfg_v4_symbol_t, out **B) C.int32_t
	withBatch    func(thread *isolateThread, cur *B, value C.int32_t, out **B) C.int32_t
	withSnapshot func(thread *isolateThread, cur *B, value C.int32_t, out **B) C.int32_t
	withListener func(thread *isolateThread, cur *B, listener *C.dxfg_v4_tx_model_listener_t, out **B) C.int32_t
	withExecutor func(thread *isolateThread, cur *B, executor *C.dxfg_v4_executor_t, out **B) C.int32_t
	withFeed     func(thread *isolateThread, cur *B, feed *C.dxfg_v4_feed_t, out **B) C.int32_t
	build        func(thread *isolateThread, cur *B) (unsafe.Pointer, C.int32_t)
	// configure applies the options of the kind of model after the common ones; nil means none.
	configure func(b *txModelBuilder[B]) error
	// close closes a model of the kind.
	close func(thread *isolateThread, model unsafe.Pointer) C.int32_t
}

// newTxModel builds a native model with the options and its listener.
func newTxModel[B any](o TxModelOptions, s txModelSteps[B]) (*txModelHandle, error) {
	h := &txModelHandle{close: s.close}
	if o.Listener != nil {
		listener, err := newTxModelListener(o.Listener)
		if err != nil {
			return nil, err
		}
		h.listener = listener
	}
	err := executeInIsolateThread(func(thread *isolateThread) error {
		// The out parameter is a variable of its own: cgo refuses a pointer into a struct that holds Go pointers.
		var cur *B
		if err := checkStatus(thread, s.newBuilder(thread, C.dxfg_v4_event_clazz_t(o.EventType), &cur)); err != nil {
			return err
		}
		b := &txModelBuilder[B]{thread: thread, cur: cur}
		err := configureTxModel(b, s, o, h)
		if err == nil {
			model, status := s.build(thread, b.cur)
			err = checkStatus(thread, status)
			if model != nil {
				h.model = NewJavaHandle(model)
			}
		}
		return errors.Join(err, b.release())
	})
	if err != nil || h.model == nil {
		if err == nil {
			err = errors.New("native: the SDK returned no transaction model")
		}
		return nil, errors.Join(err, h.Free())
	}
	return h, nil
}

// configureTxModel applies the options to the builder, one step at a time.
func configureTxModel[B any](b *txModelBuilder[B], s txModelSteps[B], o TxModelOptions, h *txModelHandle) error {
	thread := b.thread
	symbol, err := newCSymbol(o.Symbol)
	if err != nil {
		return err
	}
	defer freeCSymbol(symbol)
	if err := b.step(func(cur *B, out **B) C.int32_t { return s.withSymbol(thread, cur, symbol, out) }); err != nil {
		return err
	}
	batch, snapshot := cBool(o.BatchProcessing), cBool(o.SnapshotProcessing)
	if err := b.step(func(cur *B, out **B) C.int32_t { return s.withBatch(thread, cur, batch, out) }); err != nil {
		return err
	}
	if err := b.step(func(cur *B, out **B) C.int32_t { return s.withSnapshot(thread, cur, snapshot, out) }); err != nil {
		return err
	}
	if s.configure != nil {
		if err := s.configure(b); err != nil {
			return err
		}
	}
	if h.listener != nil {
		listener := (*C.dxfg_v4_tx_model_listener_t)(h.listener.handle.Ptr())
		if err := b.step(func(cur *B, out **B) C.int32_t { return s.withListener(thread, cur, listener, out) }); err != nil {
			return err
		}
	}
	if o.Executor != nil {
		executor := o.Executor.ptr()
		if executor == nil {
			return errors.New("native: the executor of the model is closed")
		}
		if err := b.step(func(cur *B, out **B) C.int32_t { return s.withExecutor(thread, cur, executor, out) }); err != nil {
			return err
		}
	}
	if o.Feed != nil {
		feed := o.Feed.ptr()
		if feed == nil {
			return errors.New("native: the feed of the model is closed")
		}
		if err := b.step(func(cur *B, out **B) C.int32_t { return s.withFeed(thread, cur, feed, out) }); err != nil {
			return err
		}
	}
	return nil
}

func newTxModelListener(listener TxModelListenerFunc) (*nativeListener, error) {
	userData := Save(listener)
	var l *C.dxfg_v4_tx_model_listener_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_TxModelListener_new(thread.ptr,
			(*[0]byte)(C.OnTxModelEventsReceived), userData, &l))
	})
	if l == nil {
		Unref(userData)
		if err == nil {
			err = errors.New("native: the SDK returned no transaction model listener")
		}
		return nil, err
	}
	nl := &nativeListener{handle: NewJavaHandle(unsafe.Pointer(l)), userData: userData}
	if err != nil {
		return nil, errors.Join(err, nl.Release())
	}
	return nl, nil
}

// txModelHandle is a native transaction model with its listener.
type txModelHandle struct {
	model    Handler
	listener *nativeListener
	close    func(thread *isolateThread, model unsafe.Pointer) C.int32_t
}

func (h *txModelHandle) ptr() (unsafe.Pointer, error) {
	if h.model == nil || h.model.Ptr() == nil {
		return nil, errors.New("native: the transaction model is closed")
	}
	return h.model.Ptr(), nil
}

// call runs f with the model, or fails if the model is released.
func (h *txModelHandle) call(f func(thread *isolateThread, model unsafe.Pointer) error) error {
	model, err := h.ptr()
	if err != nil {
		return err
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return f(thread, model)
	})
}

// withFeed runs f with the model and the feed.
func (h *txModelHandle) withFeed(feed *DXFeedHandle,
	f func(thread *isolateThread, model unsafe.Pointer, feed *C.dxfg_v4_feed_t) C.int32_t) error {
	if feed == nil || feed.ptr() == nil {
		return errors.New("native: the feed is nil or closed")
	}
	return h.call(func(thread *isolateThread, model unsafe.Pointer) error {
		return checkStatus(thread, f(thread, model, feed.ptr()))
	})
}

// Close closes the model: it is detached from all the feeds for good. Later calls do nothing.
func (h *txModelHandle) Close() error {
	model, err := h.ptr()
	if err != nil {
		return nil
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, h.close(thread, model))
	})
}

// Free closes the model, releases it and then its listener, so the listener is not called after Free returns
// (except when Free runs in the listener itself). Later calls do nothing.
func (h *txModelHandle) Free() error {
	return errors.Join(h.FreeModel(), h.FreeListener())
}

// FreeModel closes and releases the model. Later calls do nothing.
func (h *txModelHandle) FreeModel() error {
	if h.model == nil {
		return nil
	}
	return errors.Join(h.Close(), h.model.Free())
}

// FreeListener releases the listener: it is not called after FreeListener returns, except when FreeListener runs in
// the listener itself. Later calls do nothing.
func (h *txModelHandle) FreeListener() error {
	if h.listener == nil {
		return nil
	}
	return h.listener.Release()
}

func cBool(value bool) C.int32_t {
	if value {
		return 1
	}
	return 0
}
