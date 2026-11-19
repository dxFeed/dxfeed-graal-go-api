package model

import (
	"math"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// TimeSeriesTxModelListener receives the events of one or more complete transactions of a TimeSeriesTxModel. When
// isSnapshot is set, the events are a snapshot: the state kept for the symbol has to be cleared first. The events
// belong to the listener. The listener runs in the executor of the model, one call at a time.
type TimeSeriesTxModelListener[E events.TimeSeriesEvent] func(events []E, isSnapshot bool)

// TimeSeriesTxModelBuilder sets up a TimeSeriesTxModel; create it with NewTimeSeriesTxModelBuilder. The With methods
// change the builder and return it, so they can be chained.
type TimeSeriesTxModelBuilder[E events.TimeSeriesEvent] struct {
	txModelBuilder
	fromTime int64
	listener TimeSeriesTxModelListener[E]
}

// NewTimeSeriesTxModelBuilder creates a builder of a model of the time series events of type E. E must be a pointer
// to a concrete time series event type, such as *candle.Candle or *timeandsale.TimeAndSale; an interface type makes
// Build return an error. By default the model processes the transactions in batches, passes a snapshot as it comes
// and is not subscribed: set the time to subscribe from with WithFromTime or SetFromTime.
func NewTimeSeriesTxModelBuilder[E events.TimeSeriesEvent]() *TimeSeriesTxModelBuilder[E] {
	return &TimeSeriesTxModelBuilder[E]{txModelBuilder: newTxModelBuilder[E](), fromTime: math.MaxInt64}
}

// WithFeed sets the feed to attach the model to when it is built.
func (b *TimeSeriesTxModelBuilder[E]) WithFeed(feed *api.DXFeed) *TimeSeriesTxModelBuilder[E] {
	b.feed = feed
	return b
}

// WithSymbol sets the symbol of the model: a string, an osub.String or a *candle.CandleSymbol. It is required.
func (b *TimeSeriesTxModelBuilder[E]) WithSymbol(symbol any) *TimeSeriesTxModelBuilder[E] {
	b.options.Symbol = symbol
	return b
}

// WithFromTime sets the time in milliseconds since the epoch to subscribe from; math.MaxInt64, the default, leaves the
// model not subscribed.
func (b *TimeSeriesTxModelBuilder[E]) WithFromTime(fromTime int64) *TimeSeriesTxModelBuilder[E] {
	b.fromTime = fromTime
	return b
}

// WithListener sets the listener of the transactions.
func (b *TimeSeriesTxModelBuilder[E]) WithListener(listener TimeSeriesTxModelListener[E]) *TimeSeriesTxModelBuilder[E] {
	b.listener = listener
	return b
}

// WithBatchProcessing sets whether the listener gets all the pending transactions in one call (true, the default) or
// one transaction per call (false).
func (b *TimeSeriesTxModelBuilder[E]) WithBatchProcessing(batchProcessing bool) *TimeSeriesTxModelBuilder[E] {
	b.options.BatchProcessing = batchProcessing
	return b
}

// WithSnapshotProcessing sets whether a snapshot is passed as its resulting state, without the removed events (true),
// or as it comes (false, the default).
func (b *TimeSeriesTxModelBuilder[E]) WithSnapshotProcessing(snapshotProcessing bool) *TimeSeriesTxModelBuilder[E] {
	b.options.SnapshotProcessing = snapshotProcessing
	return b
}

// WithExecutor sets the executor that runs the listener; by default it is the executor of the feed. The model keeps
// using the executor after the executor is closed.
func (b *TimeSeriesTxModelBuilder[E]) WithExecutor(executor api.Executor) *TimeSeriesTxModelBuilder[E] {
	b.executor = executor
	return b
}

// Build builds the model and attaches it to the feed, if one is set. Close the model when it is no longer needed.
func (b *TimeSeriesTxModelBuilder[E]) Build() (*TimeSeriesTxModel[E], error) {
	options, err := b.nativeOptions()
	if err != nil {
		return nil, err
	}
	if listener := b.listener; listener != nil {
		options.Listener = func(_ events.IndexedEventSourceInterface, received []interface{}, isSnapshot bool) {
			listener(typedEvents[E](received), isSnapshot)
		}
	}
	handle, err := native.NewTimeSeriesTxModel(native.TimeSeriesTxModelOptions{TxModelOptions: options,
		FromTime: b.fromTime})
	if err != nil {
		return nil, err
	}
	m := &TimeSeriesTxModel[E]{handle: handle}
	m.model.init(handle, options)
	return m, nil
}

// TimeSeriesTxModel is a model of the time series events of type E, such as *candle.Candle, of one symbol from a
// time: it passes the complete transactions of the events to its listener and skips the incomplete ones; create it
// with NewTimeSeriesTxModelBuilder. Its methods may be called from several goroutines.
type TimeSeriesTxModel[E events.TimeSeriesEvent] struct {
	model  txModel
	handle *native.TimeSeriesTxModelHandle
}

// IsBatchProcessing returns whether the listener gets all the pending transactions in one call.
func (m *TimeSeriesTxModel[E]) IsBatchProcessing() bool {
	return m.model.batchProcessing
}

// IsSnapshotProcessing returns whether a snapshot is passed as its resulting state.
func (m *TimeSeriesTxModel[E]) IsSnapshotProcessing() bool {
	return m.model.snapshotProcessing
}

// Attach attaches the model to the feed: it subscribes to the events of the feed.
func (m *TimeSeriesTxModel[E]) Attach(feed *api.DXFeed) error {
	return m.model.withFeed(feed, m.handle.Attach)
}

// Detach detaches the model from the feed.
func (m *TimeSeriesTxModel[E]) Detach(feed *api.DXFeed) error {
	return m.model.withFeed(feed, m.handle.Detach)
}

// FromTime returns the time in milliseconds since the epoch the model is subscribed from; math.MaxInt64 means not
// subscribed.
func (m *TimeSeriesTxModel[E]) FromTime() (int64, error) {
	var fromTime int64
	err := m.model.call(func() error {
		var err error
		fromTime, err = m.handle.FromTime()
		return err
	})
	return fromTime, err
}

// SetFromTime sets the time in milliseconds since the epoch to subscribe from: the model gets a new snapshot from
// that time; math.MaxInt64 unsubscribes.
func (m *TimeSeriesTxModel[E]) SetFromTime(fromTime int64) error {
	return m.model.call(func() error { return m.handle.SetFromTime(fromTime) })
}

// Close detaches the model from all the feeds for good and releases it: the listener is not called after Close
// returns, unless Close is called from the listener itself. Later calls do nothing: a concurrent second Close
// returns at once, so the guarantee holds for the Close that performs the release.
func (m *TimeSeriesTxModel[E]) Close() error {
	return m.model.close()
}
