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

// IndexedTxModelOptions are the settings of a new IndexedTxModel.
type IndexedTxModelOptions struct {
	// EventType is the native code of the event type; required.
	EventType int32
	// Symbol is the symbol of the model, as newCSymbol accepts it; required.
	Symbol any
	// Feed is the feed to attach the model to; nil leaves the model detached.
	Feed *DXFeedHandle
	// Sources are the sources of the events; empty means all the sources.
	Sources []events.IndexedEventSourceInterface
	// Listener receives the transactions; nil means none.
	Listener TxModelListenerFunc
	// BatchProcessing and SnapshotProcessing are passed to the builder as they are.
	BatchProcessing    bool
	SnapshotProcessing bool
	// Executor runs the listener; nil means the executor of the feed.
	Executor *ExecutorHandle
}

// IndexedTxModelHandle is a native IndexedTxModel with its listener.
type IndexedTxModelHandle struct {
	model    Handler
	listener *nativeListener
}

// NewIndexedTxModel builds a native IndexedTxModel with the options.
func NewIndexedTxModel(o IndexedTxModelOptions) (*IndexedTxModelHandle, error) {
	h := &IndexedTxModelHandle{}
	if o.Listener != nil {
		listener, err := newTxModelListener(o.Listener)
		if err != nil {
			return nil, err
		}
		h.listener = listener
	}
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var builder *C.dxfg_v4_indexed_tx_model_builder_t
		if err := checkStatus(thread, C.dxfg_v4_IndexedTxModel_newBuilder(thread.ptr,
			C.dxfg_v4_event_clazz_t(o.EventType), &builder)); err != nil {
			return err
		}
		release := func() error {
			if builder == nil {
				return nil
			}
			err := checkStatus(thread, C.dxfg_v4_JavaObjectHandler_release(thread.ptr,
				(*C.dxfg_v4_java_object_handler)(unsafe.Pointer(builder))))
			builder = nil
			return err
		}
		// step passes the current builder to f, replaces it with the one that f returns and releases the previous one.
		step := func(f func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t) error {
			var next *C.dxfg_v4_indexed_tx_model_builder_t
			if err := checkStatus(thread, f(builder, &next)); err != nil {
				return err
			}
			err := release()
			builder = next
			return err
		}
		err := configureIndexedTxModel(thread, step, o, h)
		if err == nil {
			var model *C.dxfg_v4_indexed_tx_model_t
			err = checkStatus(thread, C.dxfg_v4_IndexedTxModel_Builder_build(thread.ptr, builder, &model))
			if model != nil {
				h.model = NewJavaHandle(unsafe.Pointer(model))
			}
		}
		return errors.Join(err, release())
	})
	if err != nil || h.model == nil {
		if err == nil {
			err = errors.New("native: the SDK returned no transaction model")
		}
		return nil, errors.Join(err, h.Free())
	}
	return h, nil
}

// configureIndexedTxModel applies the options to the builder, one step at a time.
func configureIndexedTxModel(thread *isolateThread,
	step func(f func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t) error,
	o IndexedTxModelOptions, h *IndexedTxModelHandle) error {
	symbol, err := newCSymbol(o.Symbol)
	if err != nil {
		return err
	}
	defer freeCSymbol(symbol)
	if err := step(func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t {
		return C.dxfg_v4_IndexedTxModel_Builder_withSymbol(thread.ptr, cur, symbol, out)
	}); err != nil {
		return err
	}
	if err := step(func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t {
		return C.dxfg_v4_IndexedTxModel_Builder_withBatchProcessing(thread.ptr, cur, cBool(o.BatchProcessing), out)
	}); err != nil {
		return err
	}
	if err := step(func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t {
		return C.dxfg_v4_IndexedTxModel_Builder_withSnapshotProcessing(thread.ptr, cur,
			cBool(o.SnapshotProcessing), out)
	}); err != nil {
		return err
	}
	if len(o.Sources) > 0 {
		sources, err := newCIndexedEventSourceList(o.Sources)
		if err != nil {
			return err
		}
		err = step(func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withSources(thread.ptr, cur, sources, out)
		})
		freeCIndexedEventSourceList(sources)
		if err != nil {
			return err
		}
	}
	if h.listener != nil {
		listener := (*C.dxfg_v4_tx_model_listener_t)(h.listener.handle.Ptr())
		if err := step(func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withListener(thread.ptr, cur, listener, out)
		}); err != nil {
			return err
		}
	}
	if o.Executor != nil {
		executor := o.Executor.ptr()
		if executor == nil {
			return errors.New("native: the executor of the model is closed")
		}
		if err := step(func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withExecutor(thread.ptr, cur, executor, out)
		}); err != nil {
			return err
		}
	}
	if o.Feed != nil {
		feed := o.Feed.ptr()
		if feed == nil {
			return errors.New("native: the feed of the model is closed")
		}
		if err := step(func(cur *C.dxfg_v4_indexed_tx_model_builder_t, out **C.dxfg_v4_indexed_tx_model_builder_t) C.int32_t {
			return C.dxfg_v4_IndexedTxModel_Builder_withFeed(thread.ptr, cur, feed, out)
		}); err != nil {
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

func (h *IndexedTxModelHandle) ptr() (*C.dxfg_v4_indexed_tx_model_t, error) {
	if h.model == nil || h.model.Ptr() == nil {
		return nil, errors.New("native: the transaction model is closed")
	}
	return (*C.dxfg_v4_indexed_tx_model_t)(h.model.Ptr()), nil
}

// call runs f with the model, or fails if the model is released.
func (h *IndexedTxModelHandle) call(f func(thread *isolateThread, model *C.dxfg_v4_indexed_tx_model_t) error) error {
	model, err := h.ptr()
	if err != nil {
		return err
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return f(thread, model)
	})
}

// Attach attaches the model to the feed.
func (h *IndexedTxModelHandle) Attach(feed *DXFeedHandle) error {
	return h.withFeed(feed, func(thread *isolateThread, model *C.dxfg_v4_indexed_tx_model_t, f *C.dxfg_v4_feed_t) C.int32_t {
		return C.dxfg_v4_IndexedTxModel_attach(thread.ptr, model, f)
	})
}

// Detach detaches the model from the feed.
func (h *IndexedTxModelHandle) Detach(feed *DXFeedHandle) error {
	return h.withFeed(feed, func(thread *isolateThread, model *C.dxfg_v4_indexed_tx_model_t, f *C.dxfg_v4_feed_t) C.int32_t {
		return C.dxfg_v4_IndexedTxModel_detach(thread.ptr, model, f)
	})
}

func (h *IndexedTxModelHandle) withFeed(feed *DXFeedHandle,
	f func(thread *isolateThread, model *C.dxfg_v4_indexed_tx_model_t, feed *C.dxfg_v4_feed_t) C.int32_t) error {
	if feed == nil || feed.ptr() == nil {
		return errors.New("native: the feed is nil or closed")
	}
	return h.call(func(thread *isolateThread, model *C.dxfg_v4_indexed_tx_model_t) error {
		return checkStatus(thread, f(thread, model, feed.ptr()))
	})
}

// Sources returns the sources of the model; empty means all the sources.
func (h *IndexedTxModelHandle) Sources() ([]events.IndexedEventSourceInterface, error) {
	var result []events.IndexedEventSourceInterface
	err := h.call(func(thread *isolateThread, model *C.dxfg_v4_indexed_tx_model_t) error {
		var list *C.dxfg_v4_indexed_event_source_list
		if err := checkStatus(thread, C.dxfg_v4_IndexedTxModel_getSources(thread.ptr, model, &list)); err != nil {
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
	return h.call(func(thread *isolateThread, model *C.dxfg_v4_indexed_tx_model_t) error {
		return checkStatus(thread, C.dxfg_v4_IndexedTxModel_setSources(thread.ptr, model, list))
	})
}

// Close closes the model: it is detached from all the feeds for good. Later calls do nothing.
func (h *IndexedTxModelHandle) Close() error {
	model, err := h.ptr()
	if err != nil {
		return nil
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_IndexedTxModel_close(thread.ptr, model))
	})
}

// Free closes the model, releases it and then its listener, so the listener is not called after Free returns
// (except when Free runs in the listener itself). Later calls do nothing.
func (h *IndexedTxModelHandle) Free() error {
	return errors.Join(h.FreeModel(), h.FreeListener())
}

// FreeModel closes and releases the model. Later calls do nothing.
func (h *IndexedTxModelHandle) FreeModel() error {
	if h.model == nil {
		return nil
	}
	return errors.Join(h.Close(), h.model.Free())
}

// FreeListener releases the listener: it is not called after FreeListener returns, except when FreeListener runs in
// the listener itself. Later calls do nothing.
func (h *IndexedTxModelHandle) FreeListener() error {
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
