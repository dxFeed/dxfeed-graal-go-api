package model

import (
	"errors"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
)

// MarketDepthModelListener receives the buy and the sell orders of the book, each side sorted from the best price and
// cut to the depth limit. The slices belong to the listener; the orders are shared with the model and must not be
// changed.
type MarketDepthModelListener[E OrderEvent] func(buy, sell []E)

// MarketDepthModelBuilder sets up a MarketDepthModel; create it with NewMarketDepthModelBuilder. The With methods
// change the builder and return it, so they can be chained.
type MarketDepthModelBuilder[E OrderEvent] struct {
	builder           *IndexedTxModelBuilder[E]
	listener          MarketDepthModelListener[E]
	depthLimit        int
	aggregationPeriod time.Duration
}

// NewMarketDepthModelBuilder creates a builder of a model of the order book of the orders of type E, such as
// *order.Order. By default the book has no depth limit and the listener is notified of each change at once.
func NewMarketDepthModelBuilder[E OrderEvent]() *MarketDepthModelBuilder[E] {
	return &MarketDepthModelBuilder[E]{builder: NewIndexedTxModelBuilder[E]()}
}

// WithFeed sets the feed to attach the model to when it is built; the model receives no orders without it.
func (b *MarketDepthModelBuilder[E]) WithFeed(feed *api.DXFeed) *MarketDepthModelBuilder[E] {
	b.builder.WithFeed(feed)
	return b
}

// WithSymbol sets the symbol of the model: a string or an osub.String. It is required.
func (b *MarketDepthModelBuilder[E]) WithSymbol(symbol any) *MarketDepthModelBuilder[E] {
	b.builder.WithSymbol(symbol)
	return b
}

// WithSources sets the order sources of the book; none means all the sources.
func (b *MarketDepthModelBuilder[E]) WithSources(sources ...events.IndexedEventSourceInterface) *MarketDepthModelBuilder[E] {
	b.builder.WithSources(sources...)
	return b
}

// WithListener sets the listener of the changes of the book.
func (b *MarketDepthModelBuilder[E]) WithListener(listener MarketDepthModelListener[E]) *MarketDepthModelBuilder[E] {
	b.listener = listener
	return b
}

// WithDepthLimit sets the maximum number of the orders of each side that the listener gets, and so the changes it is
// notified of; 0, the default, means no limit. A negative limit makes Build return an error.
func (b *MarketDepthModelBuilder[E]) WithDepthLimit(depthLimit int) *MarketDepthModelBuilder[E] {
	b.depthLimit = depthLimit
	return b
}

// WithAggregationPeriod sets the period over which the incremental changes are gathered into one notification; 0,
// the default, notifies of each change at once. A snapshot is always notified of at once. A negative period makes
// Build return an error.
func (b *MarketDepthModelBuilder[E]) WithAggregationPeriod(period time.Duration) *MarketDepthModelBuilder[E] {
	b.aggregationPeriod = period
	return b
}

// Build builds the model and attaches it to the feed, if one is set. Close the model when it is no longer needed.
func (b *MarketDepthModelBuilder[E]) Build() (*MarketDepthModel[E], error) {
	if err := checkDepthLimit(b.depthLimit); err != nil {
		return nil, err
	}
	if err := checkAggregationPeriod(b.aggregationPeriod); err != nil {
		return nil, err
	}
	m := &MarketDepthModel[E]{
		call:              newListenerCall(),
		done:              make(chan struct{}),
		ordersByIndex:     make(map[int64]*bookEntry[E]),
		buy:               newOrderBook(compareBuyOrders[E], b.depthLimit),
		sell:              newOrderBook(compareSellOrders[E], b.depthLimit),
		listener:          b.listener,
		aggregationPeriod: b.aggregationPeriod,
	}
	indexed, err := b.builder.WithListener(m.eventsReceived).Build()
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	m.indexed = indexed
	closed := m.closed
	m.mu.Unlock()
	if closed {
		// The listener closed the model before Build returned.
		return m, indexed.Close()
	}
	return m, nil
}

// MarketDepthModel is a model of the order book of a symbol, built on IndexedTxModel: it keeps the buy and the sell
// orders of its sources, sorted from the best price, and passes them to its listener when the part of the book within
// the depth limit changes; create it with NewMarketDepthModelBuilder. A snapshot of a source replaces its orders in
// the book. The orders without a size and the removed ones leave the book. Its methods may be called from several
// goroutines; the listener is called one call at a time and may call the methods of the model.
type MarketDepthModel[E OrderEvent] struct {
	mu   sync.Mutex
	call *listenerCall
	// done is closed by Close: a notification that waits for the call in progress gives up.
	done              chan struct{}
	ordersByIndex     map[int64]*bookEntry[E]
	buy               orderBook[E]
	sell              orderBook[E]
	indexed           *IndexedTxModel[E]
	listener          MarketDepthModelListener[E]
	depthLimit        int
	aggregationPeriod time.Duration
	// timer runs the aggregated notification; taskScheduled reports whether one is pending.
	timer         *time.Timer
	taskScheduled bool
	closed        bool
}

func (m *MarketDepthModel[E]) eventsReceived(source events.IndexedEventSourceInterface, orders []E, isSnapshot bool) {
	m.mu.Lock()
	if m.closed || !m.update(source.Id(), orders, isSnapshot) {
		m.mu.Unlock()
		return
	}
	notifyNow := false
	if !isSnapshot && m.aggregationPeriod != 0 {
		m.scheduleTask()
	} else {
		m.cancelTask()
		notifyNow = true
	}
	m.mu.Unlock()
	if notifyNow {
		m.notifyListener()
	}
}

// update applies the orders to the book; it reports whether the part of the book within the depth limit changed.
func (m *MarketDepthModel[E]) update(source int64, orders []E, isSnapshot bool) bool {
	if isSnapshot {
		m.clearBySource(source)
	}
	for _, o := range orders {
		if old, ok := m.ordersByIndex[o.Index()]; ok {
			delete(m.ordersByIndex, o.Index())
			m.side(old.order).erase(old)
		}
		if o.HasSize() && !events.IsRemove(o) {
			e := &bookEntry[E]{order: o, source: source}
			m.ordersByIndex[o.Index()] = e
			m.side(o).insert(e)
		}
	}
	return m.buy.changed || m.sell.changed
}

// side returns the side of the book of the order: a buy order goes to the buy side, the others to the sell side.
func (m *MarketDepthModel[E]) side(o E) *orderBook[E] {
	if o.Side() == side.Buy {
		return &m.buy
	}
	return &m.sell
}

func (m *MarketDepthModel[E]) clearBySource(source int64) {
	for index, e := range m.ordersByIndex {
		if e.source == source {
			delete(m.ordersByIndex, index)
		}
	}
	m.buy.clearBySource(source)
	m.sell.clearBySource(source)
}

// notifyListener passes the book to the listener; it must be called without m.mu.
func (m *MarketDepthModel[E]) notifyListener() {
	if !m.call.enter(m.done) {
		return
	}
	defer m.call.exit()
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return
	}
	buy, sell := m.buy.toSlice(), m.sell.toSlice()
	m.taskScheduled = false
	m.mu.Unlock()
	if m.listener != nil {
		m.listener(buy, sell)
	}
}

// scheduleTask schedules the aggregated notification, unless one is pending; it must be called with m.mu.
func (m *MarketDepthModel[E]) scheduleTask() {
	if !m.taskScheduled {
		m.taskScheduled = true
		m.timer = time.AfterFunc(m.aggregationPeriod, m.notifyListener)
	}
}

// cancelTask cancels the pending aggregated notification and reports whether there was one; it must be called with
// m.mu.
func (m *MarketDepthModel[E]) cancelTask() bool {
	if !m.taskScheduled {
		return false
	}
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
	m.taskScheduled = false
	return true
}

// DepthLimit returns the maximum number of the orders of each side that the listener gets; 0 means no limit.
func (m *MarketDepthModel[E]) DepthLimit() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.buy.depthLimit
}

// SetDepthLimit sets the maximum number of the orders of each side that the listener gets; 0 means no limit. The
// listener is notified of the book with the new limit at once.
func (m *MarketDepthModel[E]) SetDepthLimit(depthLimit int) error {
	if err := checkDepthLimit(depthLimit); err != nil {
		return err
	}
	m.mu.Lock()
	if depthLimit == m.buy.depthLimit {
		m.mu.Unlock()
		return nil
	}
	m.buy.setDepthLimit(depthLimit)
	m.sell.setDepthLimit(depthLimit)
	m.cancelTask()
	m.mu.Unlock()
	m.notifyListener()
	return nil
}

// AggregationPeriod returns the period over which the incremental changes are gathered into one notification.
func (m *MarketDepthModel[E]) AggregationPeriod() time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.aggregationPeriod
}

// SetAggregationPeriod sets the period over which the incremental changes are gathered into one notification; 0
// notifies of each change at once. A pending notification is rescheduled with the new period, or made at once with 0.
func (m *MarketDepthModel[E]) SetAggregationPeriod(period time.Duration) error {
	if err := checkAggregationPeriod(period); err != nil {
		return err
	}
	m.mu.Lock()
	if period == m.aggregationPeriod {
		m.mu.Unlock()
		return nil
	}
	m.aggregationPeriod = period
	notifyNow := false
	if m.cancelTask() {
		if period != 0 {
			m.scheduleTask()
		} else {
			notifyNow = true
		}
	}
	m.mu.Unlock()
	if notifyNow {
		m.notifyListener()
	}
	return nil
}

// Close detaches the model from all the feeds for good and releases it. A pending aggregated notification is
// cancelled, and the listener is not called after Close returns: Close waits for the call in progress, unless it is
// called from the listener itself. Later calls do nothing but wait the same way.
func (m *MarketDepthModel[E]) Close() error {
	m.mu.Lock()
	closing := !m.closed
	if closing {
		m.closed = true
		close(m.done)
		m.cancelTask()
	}
	indexed := m.indexed
	m.mu.Unlock()
	var err error
	if closing && indexed != nil {
		err = indexed.Close()
	}
	m.call.wait()
	return err
}

func checkDepthLimit(depthLimit int) error {
	if depthLimit < 0 {
		return errors.New("model: the depth limit is negative")
	}
	return nil
}

func checkAggregationPeriod(period time.Duration) error {
	if period < 0 {
		return errors.New("model: the aggregation period is negative")
	}
	return nil
}
