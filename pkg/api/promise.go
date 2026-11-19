// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// ErrPromiseCancelled is the error of a promise that was cancelled, by Cancel or by the context of Await.
var ErrPromiseCancelled = native.ErrPromiseCancelled

// ErrPromiseNotDone is returned by Promise.Result for a promise that is not done yet.
var ErrPromiseNotDone = errors.New("dxfeed: the promise is not done")

// Promise is the result of a feed request that arrives later: an event (Promise[interface{}], nil if there is none)
// or a list of events (Promise[[]interface{}]). It completes once, normally, exceptionally or by cancellation, and
// releases its native memory then. Cancel a promise that is no longer needed, so that the request stops.
type Promise[T any] struct {
	p *native.PromiseHandle
}

// Done returns a channel that is closed when the promise is done.
func (p *Promise[T]) Done() <-chan struct{} {
	return p.p.Done()
}

// IsDone reports whether the promise is done.
func (p *Promise[T]) IsDone() bool {
	select {
	case <-p.p.Done():
		return true
	default:
		return false
	}
}

// IsCancelled reports whether the promise is done by cancellation.
func (p *Promise[T]) IsCancelled() bool {
	return p.IsDone() && p.p.IsCancelled()
}

// Await waits until the promise is done and returns its result or its error. When the context is done first, Await
// cancels the promise and returns the error of the context.
func (p *Promise[T]) Await(ctx context.Context) (T, error) {
	select {
	case <-p.p.Done():
		return p.Result()
	case <-ctx.Done():
		// The promise may be done at the same moment; its outcome wins then.
		select {
		case <-p.p.Done():
			return p.Result()
		default:
		}
		_ = p.p.Cancel()
		var zero T
		return zero, ctx.Err()
	}
}

// Result returns the result or the error of a done promise; for a promise that is not done yet it returns
// ErrPromiseNotDone.
func (p *Promise[T]) Result() (T, error) {
	var zero T
	if !p.IsDone() {
		return zero, ErrPromiseNotDone
	}
	value, err := p.p.Result()
	if err != nil {
		return zero, err
	}
	result, _ := value.(T)
	return result, nil
}

// Cancel cancels the promise; it does nothing if the promise is done.
func (p *Promise[T]) Cancel() error {
	return p.p.Cancel()
}

// GetLastEventPromise requests the last event of the lasting type for the symbol, from the cache of the feed or from
// the uplink. The result is nil when the event is not available.
func (f *DXFeed) GetLastEventPromise(eventType eventcodes.EventCode, symbol osub.Symbol) (*Promise[interface{}],
	error) {
	p, err := f.feed.GetLastEventPromise(eventType.NativeCode(), symbol)
	if err != nil {
		return nil, err
	}
	return &Promise[interface{}]{p: p}, nil
}

// GetLastEventsPromises requests the last events of the lasting type for the symbols: one promise per symbol, in the
// order of the symbols.
func (f *DXFeed) GetLastEventsPromises(eventType eventcodes.EventCode,
	symbols ...osub.Symbol) ([]*Promise[interface{}], error) {
	handles, err := f.feed.GetLastEventsPromises(eventType.NativeCode(), symbols)
	if err != nil {
		return nil, err
	}
	promises := make([]*Promise[interface{}], len(handles))
	for i, h := range handles {
		promises[i] = &Promise[interface{}]{p: h}
	}
	return promises, nil
}

// GetIndexedEventsPromise requests a consistent snapshot of the indexed events of the type for the symbol and source,
// ordered by index. The symbol must not be an *osub.IndexedEventSubscriptionSymbol.
func (f *DXFeed) GetIndexedEventsPromise(eventType eventcodes.EventCode, symbol osub.Symbol,
	source events.IndexedEventSourceInterface) (*Promise[[]interface{}], error) {
	if !hasSource(source) {
		return nil, errNoSource
	}
	p, err := f.feed.GetIndexedEventsPromise(eventType.NativeCode(), symbol, source)
	if err != nil {
		return nil, err
	}
	return &Promise[[]interface{}]{p: p}, nil
}

// GetTimeSeriesPromise requests the time series events of the type for the symbol from fromTime to toTime
// (inclusive, in milliseconds since the epoch; math.MaxInt64 is no upper limit), ordered by time. The symbol must not
// be an *osub.TimeSeriesSubscriptionSymbol.
func (f *DXFeed) GetTimeSeriesPromise(eventType eventcodes.EventCode, symbol osub.Symbol, fromTime,
	toTime int64) (*Promise[[]interface{}], error) {
	p, err := f.feed.GetTimeSeriesPromise(eventType.NativeCode(), symbol, fromTime, toTime)
	if err != nil {
		return nil, err
	}
	return &Promise[[]interface{}]{p: p}, nil
}
