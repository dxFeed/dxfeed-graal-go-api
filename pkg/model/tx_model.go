// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/apibridge"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

var errClosed = errors.New("model: the model is closed")

// txModelBuilder holds the settings common to the builders of the transaction models.
type txModelBuilder struct {
	options  native.TxModelOptions
	feed     *api.DXFeed
	executor api.Executor
	err      error
}

// newTxModelBuilder returns the common settings of a model of the events of type E: batch processing, no snapshot
// processing.
func newTxModelBuilder[E events.IndexedEvent]() txModelBuilder {
	b := txModelBuilder{options: native.TxModelOptions{BatchProcessing: true}}
	var zero E
	if any(zero) == nil {
		b.err = errors.New("model: the event type must be a concrete event type, such as *order.Order")
		return b
	}
	b.options.EventType = zero.Type().NativeCode()
	return b
}

// nativeOptions checks the settings and returns them with the native handles of the feed and the executor.
func (b *txModelBuilder) nativeOptions() (native.TxModelOptions, error) {
	if b.err != nil {
		return native.TxModelOptions{}, b.err
	}
	if b.options.Symbol == nil {
		return native.TxModelOptions{}, errors.New("model: the symbol of the model is not set")
	}
	options := b.options
	if b.feed != nil {
		options.Feed = apibridge.FeedHandle(b.feed)
		if options.Feed == nil {
			return native.TxModelOptions{}, errors.New("model: the feed is closed")
		}
	}
	if b.executor != nil {
		options.Executor = apibridge.ExecutorHandle(b.executor)
		if options.Executor == nil {
			return native.TxModelOptions{}, errors.New("model: the executor is closed")
		}
	}
	return options, nil
}

// typedEvents returns the events of type E and skips the others.
func typedEvents[E events.IndexedEvent](received []interface{}) []E {
	typed := make([]E, 0, len(received))
	for _, e := range received {
		if t, ok := e.(E); ok {
			typed = append(typed, t)
		}
	}
	return typed
}

// txModelHandle is the native model behind a transaction model.
type txModelHandle interface {
	Close() error
	FreeModel() error
	FreeListener() error
}

// txModel is the state common to the transaction models: the calls hold mu for reading, the release holds it for
// writing.
type txModel struct {
	mu                 sync.RWMutex
	handle             txModelHandle
	closed             atomic.Bool
	batchProcessing    bool
	snapshotProcessing bool
}

func (m *txModel) init(handle txModelHandle, options native.TxModelOptions) {
	m.handle = handle
	m.batchProcessing = options.BatchProcessing
	m.snapshotProcessing = options.SnapshotProcessing
}

func (m *txModel) call(f func() error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.closed.Load() {
		return errClosed
	}
	return f()
}

// withFeed runs f with the native handle of the feed.
func (m *txModel) withFeed(feed *api.DXFeed, f func(*native.DXFeedHandle) error) error {
	h := apibridge.FeedHandle(feed)
	if h == nil {
		return errors.New("model: the feed is nil")
	}
	return m.call(func() error { return f(h) })
}

func (m *txModel) close() error {
	if m.closed.Swap(true) {
		return nil
	}
	m.mu.RLock()
	err := m.handle.Close()
	m.mu.RUnlock()
	// The write lock waits for the running calls; the listener is released after it, since a listener that waits
	// for the read lock would never return to the release.
	m.mu.Lock()
	err = errors.Join(err, m.handle.FreeModel())
	m.mu.Unlock()
	return errors.Join(err, m.handle.FreeListener())
}
