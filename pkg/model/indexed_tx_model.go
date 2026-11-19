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

// TxModelListener receives the events of one or more complete transactions of a model, all of one source. When
// isSnapshot is set, the events are a snapshot: the state kept for the source has to be cleared first. The events
// belong to the listener. The listener runs in the executor of the model, one call at a time.
type TxModelListener[E events.IndexedEvent] func(source events.IndexedEventSourceInterface, events []E,
	isSnapshot bool)

// IndexedTxModelBuilder sets up an IndexedTxModel; create it with NewIndexedTxModelBuilder. The With methods change
// the builder and return it, so they can be chained.
type IndexedTxModelBuilder[E events.IndexedEvent] struct {
	options  native.IndexedTxModelOptions
	feed     *api.DXFeed
	listener TxModelListener[E]
	executor api.Executor
	err      error
}

// NewIndexedTxModelBuilder creates a builder of a model of the indexed events of type E. E must be a pointer to a
// concrete indexed event type, such as *order.Order; an interface type, such as events.IndexedEvent, makes Build
// return an error. By default the model processes the transactions in batches and passes a snapshot as it comes.
func NewIndexedTxModelBuilder[E events.IndexedEvent]() *IndexedTxModelBuilder[E] {
	b := &IndexedTxModelBuilder[E]{options: native.IndexedTxModelOptions{BatchProcessing: true}}
	var zero E
	if any(zero) == nil {
		b.err = errors.New("model: the event type must be a concrete event type, such as *order.Order")
		return b
	}
	b.options.EventType = zero.Type().NativeCode()
	return b
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
	b.options.Sources = append([]events.IndexedEventSourceInterface(nil), sources...)
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
	if b.err != nil {
		return nil, b.err
	}
	if b.options.Symbol == nil {
		return nil, errors.New("model: the symbol of the model is not set")
	}
	options := b.options
	if b.feed != nil {
		options.Feed = apibridge.FeedHandle(b.feed)
		if options.Feed == nil {
			return nil, errors.New("model: the feed is closed")
		}
	}
	if b.executor != nil {
		options.Executor = apibridge.ExecutorHandle(b.executor)
		if options.Executor == nil {
			return nil, errors.New("model: the executor is closed")
		}
	}
	if b.listener != nil {
		options.Listener = typedListener(b.listener)
	}
	handle, err := native.NewIndexedTxModel(options)
	if err != nil {
		return nil, err
	}
	return &IndexedTxModel[E]{
		handle:             handle,
		batchProcessing:    options.BatchProcessing,
		snapshotProcessing: options.SnapshotProcessing,
	}, nil
}

// typedListener passes the events of type E to the listener and skips the others.
func typedListener[E events.IndexedEvent](listener TxModelListener[E]) native.TxModelListenerFunc {
	return func(source events.IndexedEventSourceInterface, received []interface{}, isSnapshot bool) {
		typed := make([]E, 0, len(received))
		for _, e := range received {
			if t, ok := e.(E); ok {
				typed = append(typed, t)
			}
		}
		listener(source, typed, isSnapshot)
	}
}

// IndexedTxModel is a model of the indexed events of type E, such as *order.Order, of one symbol: it passes the
// complete transactions of the events to its listener and skips the incomplete ones; create it with
// NewIndexedTxModelBuilder. Its methods may be called from several goroutines.
type IndexedTxModel[E events.IndexedEvent] struct {
	// mu is held for reading by the calls and for writing when the model is released.
	mu                 sync.RWMutex
	handle             *native.IndexedTxModelHandle
	closed             atomic.Bool
	batchProcessing    bool
	snapshotProcessing bool
}

// IsBatchProcessing returns whether the listener gets all the pending transactions in one call.
func (m *IndexedTxModel[E]) IsBatchProcessing() bool {
	return m.batchProcessing
}

// IsSnapshotProcessing returns whether a snapshot is passed as its resulting state.
func (m *IndexedTxModel[E]) IsSnapshotProcessing() bool {
	return m.snapshotProcessing
}

// Attach attaches the model to the feed: it subscribes to the events of the feed.
func (m *IndexedTxModel[E]) Attach(feed *api.DXFeed) error {
	return m.withFeed(feed, (*native.IndexedTxModelHandle).Attach)
}

// Detach detaches the model from the feed.
func (m *IndexedTxModel[E]) Detach(feed *api.DXFeed) error {
	return m.withFeed(feed, (*native.IndexedTxModelHandle).Detach)
}

func (m *IndexedTxModel[E]) withFeed(feed *api.DXFeed,
	f func(*native.IndexedTxModelHandle, *native.DXFeedHandle) error) error {
	h := apibridge.FeedHandle(feed)
	if h == nil {
		return errors.New("model: the feed is nil")
	}
	return m.call(func() error { return f(m.handle, h) })
}

// Sources returns the sources of the events; empty means all the sources.
func (m *IndexedTxModel[E]) Sources() ([]events.IndexedEventSourceInterface, error) {
	var sources []events.IndexedEventSourceInterface
	err := m.call(func() error {
		var err error
		sources, err = m.handle.Sources()
		return err
	})
	return sources, err
}

// SetSources sets the sources of the events: the model subscribes to the new ones and unsubscribes from the others;
// none means all the sources.
func (m *IndexedTxModel[E]) SetSources(sources ...events.IndexedEventSourceInterface) error {
	return m.call(func() error { return m.handle.SetSources(sources) })
}

func (m *IndexedTxModel[E]) call(f func() error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.closed.Load() {
		return errClosed
	}
	return f()
}

// Close detaches the model from all the feeds for good and releases it: the listener is not called after Close
// returns, unless Close is called from the listener itself. Later calls do nothing: a concurrent second Close
// returns at once, so the guarantee holds for the Close that performs the release.
func (m *IndexedTxModel[E]) Close() error {
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

var errClosed = errors.New("model: the model is closed")
