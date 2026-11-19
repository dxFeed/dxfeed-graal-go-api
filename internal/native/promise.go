package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnPromiseDone(graal_isolatethread_t *thread, dxfg_promise_t *promise, void *user_data);
*/
import "C"
import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	goevents "github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// ErrPromiseCancelled is the error of a promise that was cancelled.
var ErrPromiseCancelled = errors.New("dxfeed: the promise is cancelled")

// promiseKind tells how the result of a promise is read.
type promiseKind int

const (
	// promiseOfEvent is a promise of one event (dxfg_v4_promise_event_t).
	promiseOfEvent promiseKind = iota
	// promiseOfEvents is a promise of a list of events (dxfg_v4_promise_events_t).
	promiseOfEvents
)

// promiseOwner releases the native memory of one or more promises once all of them are done: a handle of its own, or
// the list that holds the handles of several promises.
type promiseOwner struct {
	pending  atomic.Int32
	promises []*PromiseHandle
	release  func(thread *isolateThread) error
}

// done counts a done promise and releases the memory after the last one.
func (o *promiseOwner) done() {
	if o.pending.Add(-1) != 0 {
		return
	}
	// The SDK needs the promise until the handler returns, so the release runs after it, on another goroutine. It
	// waits for the running Cancel calls, and later ones do nothing.
	go func() {
		for _, p := range o.promises {
			p.mu.Lock()
			p.released = true
			p.mu.Unlock()
		}
		if err := executeInIsolateThread(o.release); err != nil {
			log.Printf("dxfeed: release a promise: %v", err)
		}
	}()
}

// PromiseHandle is a native promise of the result of a feed request. The result is read once, when the promise is
// done; after that the native promise is released.
type PromiseHandle struct {
	ptr   *C.dxfg_v4_promise_t
	kind  promiseKind
	owner *promiseOwner

	// mu guards the calls to the native promise against its release; released is set under it. The handler never
	// takes it, so Cancel may hold it while the SDK runs the handler.
	mu       sync.Mutex
	released bool

	done      chan struct{}
	result    interface{}
	err       error
	cancelled bool
}

//export OnPromiseDone
func OnPromiseDone(thread *C.graal_isolatethread_t, _ *C.dxfg_promise_t, userData unsafe.Pointer) {
	defer recoverCallback("promise handler")
	p, ok := Restore(userData).(*PromiseHandle)
	Unref(userData)
	if !ok {
		return
	}
	p.complete(&isolateThread{ptr: thread})
}

// complete reads the outcome of the promise, marks it done and lets the owner release the memory.
func (p *PromiseHandle) complete(thread *isolateThread) {
	defer p.owner.done()
	defer close(p.done)
	var hasResult C.int32_t
	if err := checkStatus(thread, C.dxfg_v4_Promise_hasResult(thread.ptr, p.ptr, &hasResult)); err != nil {
		p.err = err
		return
	}
	if hasResult != 0 {
		p.result, p.err = p.takeResult(thread)
		return
	}
	var cancelled C.int32_t
	if err := checkStatus(thread, C.dxfg_v4_Promise_isCancelled(thread.ptr, p.ptr, &cancelled)); err != nil {
		p.err = err
		return
	}
	if cancelled != 0 {
		p.cancelled = true
		p.err = ErrPromiseCancelled
		return
	}
	p.err = p.takeException(thread)
}

func (p *PromiseHandle) takeResult(thread *isolateThread) (interface{}, error) {
	switch p.kind {
	case promiseOfEvent:
		var event *C.dxfg_v4_event_type_t
		if err := checkStatus(thread, C.dxfg_v4_Promise_EventType_getResult(thread.ptr,
			(*C.dxfg_v4_promise_event_t)(unsafe.Pointer(p.ptr)), &event)); err != nil {
			return nil, err
		}
		return takeEvent(thread, event)
	default:
		var list *C.dxfg_v4_event_type_list
		if err := checkStatus(thread, C.dxfg_v4_Promise_List_EventType_getResult(thread.ptr,
			(*C.dxfg_v4_promise_events_t)(unsafe.Pointer(p.ptr)), &list)); err != nil {
			return nil, err
		}
		return takeEvents(thread, list)
	}
}

func (p *PromiseHandle) takeException(thread *isolateThread) error {
	var exception *C.dxfg_v4_exception_t
	if err := checkStatus(thread, C.dxfg_v4_Promise_getException(thread.ptr, p.ptr, &exception)); err != nil {
		return err
	}
	if exception == nil {
		return errors.New("native: the promise is done without a result or an exception")
	}
	defer C.dxfg_v4_Exception_release(thread.ptr, exception)
	return javaError(exception)
}

// Done returns a channel that is closed when the promise is done.
func (p *PromiseHandle) Done() <-chan struct{} {
	return p.done
}

// Result returns the result of a done promise: an event (or nil) for a promise of an event, a non-nil list of events
// for a promise of a list. It must be called after Done is closed.
func (p *PromiseHandle) Result() (interface{}, error) {
	return p.result, p.err
}

// IsCancelled reports whether a done promise was cancelled. It must be called after Done is closed.
func (p *PromiseHandle) IsCancelled() bool {
	return p.cancelled
}

// Cancel cancels the promise; it does nothing if the promise is done.
func (p *PromiseHandle) Cancel() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.released {
		return nil
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_Promise_cancel(thread.ptr, p.ptr))
	})
}

// newPromises wraps the native promises that owner releases and adds the handler to each of them. On an error the
// promises that got no handler are cancelled with it, so that the owner still releases the memory.
func newPromises(thread *isolateThread, ptrs []*C.dxfg_v4_promise_t, kind promiseKind,
	release func(thread *isolateThread) error) ([]*PromiseHandle, error) {
	owner := &promiseOwner{release: release}
	if len(ptrs) == 0 {
		owner.pending.Store(1)
		owner.done()
		return []*PromiseHandle{}, nil
	}
	owner.pending.Store(int32(len(ptrs)))
	promises := make([]*PromiseHandle, len(ptrs))
	for i, ptr := range ptrs {
		promises[i] = &PromiseHandle{ptr: ptr, kind: kind, owner: owner, done: make(chan struct{})}
	}
	owner.promises = promises
	for i, p := range promises {
		userData := Save(p)
		status := C.dxfg_v4_Promise_whenDone(thread.ptr, p.ptr, (*[0]byte)(C.OnPromiseDone), userData)
		if err := checkStatus(thread, status); err != nil {
			Unref(userData)
			// The promises without a handler are done here, failed with the error; the ones with a handler are
			// cancelled, so that their handlers run and the owner releases the memory.
			for _, rest := range promises[i:] {
				rest.err = err
				close(rest.done)
				owner.done()
			}
			for _, started := range promises[:i] {
				_ = checkStatus(thread, C.dxfg_v4_Promise_cancel(thread.ptr, started.ptr))
			}
			return nil, err
		}
	}
	return promises, nil
}

// newPromise wraps one native promise whose handle is released on its own.
func newPromise(thread *isolateThread, ptr *C.dxfg_v4_promise_t, kind promiseKind) (*PromiseHandle, error) {
	release := func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_JavaObjectHandler_release(thread.ptr,
			(*C.dxfg_v4_java_object_handler)(unsafe.Pointer(ptr))))
	}
	promises, err := newPromises(thread, []*C.dxfg_v4_promise_t{ptr}, kind, release)
	if err != nil {
		return nil, err
	}
	return promises[0], nil
}

// GetLastEventPromise requests the last event of the class for the symbol.
func (f *DXFeedHandle) GetLastEventPromise(eventType int32, symbol osub.Symbol) (*PromiseHandle, error) {
	if err := checkEventType("request the last event", eventType); err != nil {
		return nil, err
	}
	var result *PromiseHandle
	err := withCSymbol(symbol, func(cSymbol *C.dxfg_v4_symbol_t) error {
		return executeInIsolateThread(func(thread *isolateThread) error {
			var ptr *C.dxfg_v4_promise_event_t
			if err := checkStatus(thread, C.dxfg_v4_DXFeed_getLastEventPromise(thread.ptr, f.ptr(),
				C.dxfg_v4_event_clazz_t(eventType), cSymbol, &ptr)); err != nil {
				return err
			}
			var err error
			result, err = newPromise(thread, &ptr.handler, promiseOfEvent)
			return err
		})
	})
	return result, err
}

// GetLastEventsPromises requests the last events of the class for the symbols, one promise per symbol in the order
// of the symbols.
func (f *DXFeedHandle) GetLastEventsPromises(eventType int32, symbols []osub.Symbol) ([]*PromiseHandle, error) {
	if err := checkEventType("request the last events", eventType); err != nil {
		return nil, err
	}
	cSymbols, err := newCSymbolList(symbols)
	if err != nil {
		return nil, err
	}
	defer freeCSymbolList(cSymbols)
	var result []*PromiseHandle
	err = executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_promise_list
		if err := checkStatus(thread, C.dxfg_v4_DXFeed_getLastEventsPromises(thread.ptr, f.ptr(),
			C.dxfg_v4_event_clazz_t(eventType), cSymbols, &list)); err != nil {
			return err
		}
		release := func(thread *isolateThread) error {
			return checkStatus(thread, C.dxfg_v4_promise_list_release(thread.ptr, list))
		}
		var ptrs []*C.dxfg_v4_promise_t
		if list != nil && list.list.size > 0 {
			for _, h := range unsafe.Slice(list.list.elements, int(list.list.size)) {
				ptrs = append(ptrs, (*C.dxfg_v4_promise_t)(unsafe.Pointer(h)))
			}
		}
		var err error
		result, err = newPromises(thread, ptrs, promiseOfEvent, release)
		return err
	})
	return result, err
}

// GetIndexedEventsPromise requests the indexed events of the class for the symbol and the source.
func (f *DXFeedHandle) GetIndexedEventsPromise(eventType int32, symbol osub.Symbol,
	source goevents.IndexedEventSourceInterface) (*PromiseHandle, error) {
	if err := checkEventType("request the indexed events", eventType); err != nil {
		return nil, err
	}
	if _, ok := symbol.(*osub.IndexedEventSubscriptionSymbol); ok {
		return nil, errors.New("cannot request the indexed events: the symbol must not be an indexed event " +
			"subscription symbol")
	}
	cSource, err := newCIndexedEventSource(source)
	if err != nil {
		return nil, err
	}
	defer freeCIndexedEventSource(cSource)
	var result *PromiseHandle
	err = withCSymbol(symbol, func(cSymbol *C.dxfg_v4_symbol_t) error {
		return executeInIsolateThread(func(thread *isolateThread) error {
			var ptr *C.dxfg_v4_promise_events_t
			if err := checkStatus(thread, C.dxfg_v4_DXFeed_getIndexedEventsPromise(thread.ptr, f.ptr(),
				C.dxfg_v4_event_clazz_t(eventType), cSymbol, cSource, &ptr)); err != nil {
				return err
			}
			var err error
			result, err = newPromise(thread, &ptr.base, promiseOfEvents)
			return err
		})
	})
	return result, err
}

// GetTimeSeriesPromise requests the time series events of the class for the symbol from fromTime to toTime
// (inclusive, milliseconds since the epoch).
func (f *DXFeedHandle) GetTimeSeriesPromise(eventType int32, symbol osub.Symbol, fromTime,
	toTime int64) (*PromiseHandle, error) {
	if err := checkEventType("request the time series", eventType); err != nil {
		return nil, err
	}
	if _, ok := symbol.(*osub.TimeSeriesSubscriptionSymbol); ok {
		return nil, errors.New("cannot request the time series: the symbol must not be a time series subscription " +
			"symbol")
	}
	var result *PromiseHandle
	err := withCSymbol(symbol, func(cSymbol *C.dxfg_v4_symbol_t) error {
		return executeInIsolateThread(func(thread *isolateThread) error {
			var ptr *C.dxfg_v4_promise_events_t
			if err := checkStatus(thread, C.dxfg_v4_DXFeed_getTimeSeriesPromise(thread.ptr, f.ptr(),
				C.dxfg_v4_event_clazz_t(eventType), cSymbol, C.int64_t(fromTime), C.int64_t(toTime),
				&ptr)); err != nil {
				return err
			}
			var err error
			result, err = newPromise(thread, &ptr.base, promiseOfEvents)
			return err
		})
	})
	return result, err
}
