// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// TxModelListener receives the events of one or more complete transactions of a model, all of one source. When
// isSnapshot is set, the events are a snapshot: the state kept for the source has to be cleared first. The events
// belong to the listener. The listener runs in the executor of the model, one call at a time.
type TxModelListener[E events.IndexedEvent] func(source events.IndexedEventSourceInterface, events []E,
	isSnapshot bool)

// IndexedTxModelBuilder sets up an IndexedTxModel; create it with NewIndexedTxModelBuilder. The With methods change
// the builder and return it, so they can be chained.
type IndexedTxModelBuilder[E events.IndexedEvent] struct {
	txModelBuilder
	sources  []events.IndexedEventSourceInterface
	listener TxModelListener[E]
}

// NewIndexedTxModelBuilder creates a builder of a model of the indexed events of type E. E must be a pointer to a
// concrete indexed event type, such as *order.Order; an interface type, such as events.IndexedEvent, makes Build
// return an error. By default the model processes the transactions in batches and passes a snapshot as it comes.
func NewIndexedTxModelBuilder[E events.IndexedEvent]() *IndexedTxModelBuilder[E] {
	return &IndexedTxModelBuilder[E]{txModelBuilder: newTxModelBuilder[E]()}
}

// WithFeed sets the feed to attach the model to when it is built.
func (b *IndexedTxModelBuilder[E]) WithFeed(feed *api.DXFeed) *IndexedTxModelBuilder[E] {
	b.feed = feed
	return b
}

// WithSymbol sets the symbol of the model: a string, an osub.String or a *candle.CandleSymbol. It is required.
func (b *IndexedTxModelBuilder[E]) WithSymbol(symbol any) *IndexedTxModelBuilder[E] {
	b.options.Symbol = symbol
	return b
}

// WithSources sets the sources of the events, such as order sources; none means all the sources.
func (b *IndexedTxModelBuilder[E]) WithSources(sources ...events.IndexedEventSourceInterface) *IndexedTxModelBuilder[E] {
	b.sources = append([]events.IndexedEventSourceInterface(nil), sources...)
	return b
}

// WithListener sets the listener of the transactions.
func (b *IndexedTxModelBuilder[E]) WithListener(listener TxModelListener[E]) *IndexedTxModelBuilder[E] {
	b.listener = listener
	return b
}

// WithBatchProcessing sets whether the listener gets all the pending transactions in one call (true, the default) or
// one transaction per call (false).
func (b *IndexedTxModelBuilder[E]) WithBatchProcessing(batchProcessing bool) *IndexedTxModelBuilder[E] {
	b.options.BatchProcessing = batchProcessing
	return b
}

// WithSnapshotProcessing sets whether a snapshot is passed as its resulting state, without the removed events (true),
// or as it comes (false, the default).
func (b *IndexedTxModelBuilder[E]) WithSnapshotProcessing(snapshotProcessing bool) *IndexedTxModelBuilder[E] {
	b.options.SnapshotProcessing = snapshotProcessing
	return b
}

// WithExecutor sets the executor that runs the listener; by default it is the executor of the feed. The model keeps
// using the executor after the executor is closed.
func (b *IndexedTxModelBuilder[E]) WithExecutor(executor api.Executor) *IndexedTxModelBuilder[E] {
	b.executor = executor
	return b
}

// Build builds the model and attaches it to the feed, if one is set. Close the model when it is no longer needed.
func (b *IndexedTxModelBuilder[E]) Build() (*IndexedTxModel[E], error) {
	options, err := b.nativeOptions()
	if err != nil {
		return nil, err
	}
	if listener := b.listener; listener != nil {
		options.Listener = func(source events.IndexedEventSourceInterface, received []interface{}, isSnapshot bool) {
			listener(source, typedEvents[E](received), isSnapshot)
		}
	}
	handle, err := native.NewIndexedTxModel(native.IndexedTxModelOptions{TxModelOptions: options, Sources: b.sources})
	if err != nil {
		return nil, err
	}
	m := &IndexedTxModel[E]{handle: handle}
	m.model.init(handle, options)
	return m, nil
}

// IndexedTxModel is a model of the indexed events of type E, such as *order.Order, of one symbol: it passes the
// complete transactions of the events to its listener and skips the incomplete ones; create it with
// NewIndexedTxModelBuilder. Its methods may be called from several goroutines.
type IndexedTxModel[E events.IndexedEvent] struct {
	model  txModel
	handle *native.IndexedTxModelHandle
}

// IsBatchProcessing returns whether the listener gets all the pending transactions in one call.
func (m *IndexedTxModel[E]) IsBatchProcessing() bool {
	return m.model.batchProcessing
}

// IsSnapshotProcessing returns whether a snapshot is passed as its resulting state.
func (m *IndexedTxModel[E]) IsSnapshotProcessing() bool {
	return m.model.snapshotProcessing
}

// Attach attaches the model to the feed: it subscribes to the events of the feed.
func (m *IndexedTxModel[E]) Attach(feed *api.DXFeed) error {
	return m.model.withFeed(feed, m.handle.Attach)
}

// Detach detaches the model from the feed.
func (m *IndexedTxModel[E]) Detach(feed *api.DXFeed) error {
	return m.model.withFeed(feed, m.handle.Detach)
}

// Sources returns the sources of the events; empty means all the sources.
func (m *IndexedTxModel[E]) Sources() ([]events.IndexedEventSourceInterface, error) {
	var sources []events.IndexedEventSourceInterface
	err := m.model.call(func() error {
		var err error
		sources, err = m.handle.Sources()
		return err
	})
	return sources, err
}

// SetSources sets the sources of the events: the model subscribes to the new ones and unsubscribes from the others;
// none means all the sources.
func (m *IndexedTxModel[E]) SetSources(sources ...events.IndexedEventSourceInterface) error {
	return m.model.call(func() error { return m.handle.SetSources(sources) })
}

// Close detaches the model from all the feeds for good and releases it: the listener is not called after Close
// returns, unless Close is called from the listener itself. Later calls do nothing: a concurrent second Close
// returns at once, so the guarantee holds for the Close that performs the release.
func (m *IndexedTxModel[E]) Close() error {
	return m.model.close()
}
