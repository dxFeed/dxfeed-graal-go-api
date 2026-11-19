package model_test

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/model"
)

// The scenarios follow the C++ API tests of MarketDepthModel.

type depthFixture struct {
	*orderFixture

	mu          sync.Mutex
	calls       int
	buy, sell   []*order.Order
	changesBuy  int
	changesSell int
}

func newDepthFixture(t *testing.T) *depthFixture {
	return &depthFixture{orderFixture: newFixture[*order.Order](t)}
}

// listenBook records a notification of the listener and whether each side changed.
func (f *depthFixture) listenBook(buy, sell []*order.Order) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if !sameOrders(f.buy, buy) {
		f.changesBuy++
	}
	if !sameOrders(f.sell, sell) {
		f.changesSell++
	}
	f.buy, f.sell = buy, sell
}

func sameOrders(o1, o2 []*order.Order) bool {
	if len(o1) != len(o2) {
		return false
	}
	for i := range o1 {
		if o1[i].String() != o2[i].String() {
			return false
		}
	}
	return true
}

// depthBuilder returns a builder of a model of the orders of testSymbol of the default source.
func depthBuilder(f *depthFixture) *model.MarketDepthModelBuilder[*order.Order] {
	return model.NewMarketDepthModelBuilder[*order.Order]().WithFeed(f.feed).WithSymbol(testSymbol).
		WithSources(order.Default()).WithListener(f.listenBook)
}

func buildDepth(f *depthFixture, b *model.MarketDepthModelBuilder[*order.Order]) *model.MarketDepthModel[*order.Order] {
	f.t.Helper()
	m, err := b.Build()
	if err != nil {
		f.t.Fatalf("Build: %v", err)
	}
	f.t.Cleanup(func() { _ = m.Close() })
	return m
}

func newBookOrder(t *testing.T, index int64, s side.Side, price, size float64, flags int32) *order.Order {
	t.Helper()
	o := order.NewOrder(testSymbol)
	if err := o.SetIndex(index); err != nil {
		t.Fatalf("SetIndex: %v", err)
	}
	o.SetSide(s)
	o.SetPrice(price)
	o.SetSize(size)
	o.SetEventFlags(flags)
	return o
}

// publishAndCheck publishes the orders and checks whether the listener was notified. A notification is awaited for
// up to 10 s; its absence is checked for 300 ms.
func (f *depthFixture) publishAndCheck(changed bool, orders ...*order.Order) {
	f.t.Helper()
	f.publish(orders...)
	f.checkChanged(changed)
}

func (f *depthFixture) checkChanged(changed bool) {
	f.t.Helper()
	wait := 300 * time.Millisecond
	if changed {
		wait = 10 * time.Second
	}
	for deadline := time.Now().Add(wait); time.Now().Before(deadline); time.Sleep(5 * time.Millisecond) {
		f.process()
		f.mu.Lock()
		called := f.calls > 0
		f.mu.Unlock()
		if called && changed {
			break
		}
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if (f.calls > 0) != changed {
		f.t.Fatalf("listener calls: %d, want changed %v", f.calls, changed)
	}
	f.calls = 0
}

func (f *depthFixture) checkSizes(buy, sell int) {
	f.t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.buy) != buy || len(f.sell) != sell {
		f.t.Fatalf("book: %d buy, %d sell orders, want %d and %d", len(f.buy), len(f.sell), buy, sell)
	}
}

// checkOrder checks that the order is at the position of its side of the book.
func (f *depthFixture) checkOrder(s side.Side, o *order.Order, pos int) {
	f.t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	orders := f.sell
	if s == side.Buy {
		orders = f.buy
	}
	if pos >= len(orders) {
		f.t.Fatalf("%v orders: %d, want more than %d", s, len(orders), pos)
	}
	got := orders[pos]
	if got.Index() != o.Index() || got.Side() != o.Side() || got.Size() != o.Size() || got.Price() != o.Price() {
		f.t.Fatalf("%v order %d: %v, want %v", s, pos, got, o)
	}
}

func TestMarketDepthModelRemoveNonExistent(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f))
	f.publishAndCheck(false, newBookOrder(t, 0, side.Buy, 1, 1, events.RemoveEvent.Flag()))
}

func TestMarketDepthModelRemoveBySizeAndByFlags(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f))
	o1 := newBookOrder(t, 2, side.Buy, 3, 1, 0)
	o2 := newBookOrder(t, 1, side.Buy, 2, 1, 0)
	o3 := newBookOrder(t, 0, side.Buy, 1, 1, snapshotBeginEnd)
	f.publishAndCheck(true, o3, o2, o1)
	f.checkSizes(3, 0)
	f.checkOrder(side.Buy, o1, 0)
	f.checkOrder(side.Buy, o2, 1)
	f.checkOrder(side.Buy, o3, 2)

	f.publishAndCheck(true, newBookOrder(t, 2, side.Buy, 1, math.NaN(), 0))
	f.checkSizes(2, 0)
	f.checkOrder(side.Buy, o2, 0)

	f.publishAndCheck(true, newBookOrder(t, 1, side.Buy, 1, 1, events.RemoveEvent.Flag()))
	f.checkSizes(1, 0)
	f.checkOrder(side.Buy, o3, 0)

	f.publishAndCheck(true, newBookOrder(t, 0, side.Buy, 1, 1, events.SnapshotEnd.Flag()|events.RemoveEvent.Flag()))
	f.checkSizes(0, 0)
}

func TestMarketDepthModelOrderChangeSide(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f))
	buy := newBookOrder(t, 0, side.Buy, 1, 1, snapshotBeginEnd)
	f.publishAndCheck(true, buy)
	f.checkSizes(1, 0)
	f.checkOrder(side.Buy, buy, 0)

	sell := newBookOrder(t, 0, side.Sell, 1, 1, 0)
	f.publishAndCheck(true, sell)
	f.checkSizes(0, 1)
	f.checkOrder(side.Sell, sell, 0)
}

func TestMarketDepthModelOrderPriorityAfterUpdate(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f))
	b1 := newBookOrder(t, 0, side.Buy, 100, 1, snapshotBeginEnd)
	b2 := newBookOrder(t, 1, side.Buy, 150, 1, 0)
	s1 := newBookOrder(t, 3, side.Sell, 150, 1, 0)
	s2 := newBookOrder(t, 2, side.Sell, 100, 1, 0)

	f.publishAndCheck(true, b1)
	f.checkOrder(side.Buy, b1, 0)
	f.publishAndCheck(true, b2)
	f.checkOrder(side.Buy, b2, 0)
	f.checkOrder(side.Buy, b1, 1)

	f.publishAndCheck(true, s1)
	f.checkOrder(side.Sell, s1, 0)
	f.publishAndCheck(true, s2)
	f.checkOrder(side.Sell, s2, 0)
	f.checkOrder(side.Sell, s1, 1)
}

func TestMarketDepthModelMultipleUpdatesWithMixedSides(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f))
	buyLow := newBookOrder(t, 0, side.Buy, 100, 1, snapshotBeginEnd)
	buyHigh := newBookOrder(t, 1, side.Buy, 200, 1, 0)
	sellLow := newBookOrder(t, 2, side.Sell, 150, 1, 0)
	sellHigh := newBookOrder(t, 3, side.Sell, 250, 1, 0)

	f.publishAndCheck(true, buyLow, sellHigh, buyHigh, sellLow)
	f.checkSizes(2, 2)
	f.checkOrder(side.Buy, buyHigh, 0)
	f.checkOrder(side.Sell, sellLow, 0)
}

func TestMarketDepthModelDuplicateIndexUpdatesOrder(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f))
	original := newBookOrder(t, 0, side.Buy, 100, 1, snapshotBeginEnd)
	duplicate := newBookOrder(t, 0, side.Buy, 150, 1, 0)

	f.publishAndCheck(true, original, duplicate)
	f.checkSizes(1, 0)
	f.checkOrder(side.Buy, duplicate, 0)
}

// TestMarketDepthModelDepthLimit checks that the changes below the depth limit are not notified of.
func TestMarketDepthModelDepthLimit(t *testing.T) {
	f := newDepthFixture(t)
	m := buildDepth(f, depthBuilder(f).WithDepthLimit(3))
	if got := m.DepthLimit(); got != 3 {
		t.Fatalf("DepthLimit: %d, want 3", got)
	}

	f.publishAndCheck(true, newBookOrder(t, 0, side.Buy, 5, 1, snapshotBeginEnd))
	f.publishAndCheck(true, newBookOrder(t, 1, side.Buy, 4, 1, 0))
	f.publishAndCheck(true, newBookOrder(t, 2, side.Buy, 3, 1, 0))

	f.publishAndCheck(false, newBookOrder(t, 3, side.Buy, 2, 1, 0))          // outside the limit
	f.publishAndCheck(false, newBookOrder(t, 4, side.Buy, 1, 1, 0))          // outside the limit
	f.publishAndCheck(false, newBookOrder(t, 4, side.Buy, 1, 2, 0))          // changed outside the limit
	f.publishAndCheck(false, newBookOrder(t, 3, side.Buy, 1, math.NaN(), 0)) // removed outside the limit
	f.publishAndCheck(true, newBookOrder(t, 2, side.Buy, 3, 2, 0))           // changed within the limit
	f.publishAndCheck(true, newBookOrder(t, 1, side.Buy, 3, math.NaN(), 0))  // removed within the limit
	f.checkSizes(3, 0)

	f.publishAndCheck(true, newBookOrder(t, 4, side.Sell, 1, 1, 0))
	f.publishAndCheck(true, newBookOrder(t, 5, side.Sell, 2, 1, 0))
	f.publishAndCheck(true, newBookOrder(t, 6, side.Sell, 3, 1, 0))

	f.publishAndCheck(false, newBookOrder(t, 7, side.Sell, 4, 1, 0))          // outside the limit
	f.publishAndCheck(false, newBookOrder(t, 8, side.Sell, 5, 1, 0))          // outside the limit
	f.publishAndCheck(false, newBookOrder(t, 8, side.Sell, 5, 2, 0))          // changed outside the limit
	f.publishAndCheck(false, newBookOrder(t, 8, side.Sell, 5, math.NaN(), 0)) // removed outside the limit
	f.publishAndCheck(true, newBookOrder(t, 6, side.Sell, 4, 2, 0))           // changed within the limit
	f.publishAndCheck(true, newBookOrder(t, 5, side.Sell, 2, math.NaN(), 0))  // removed within the limit
	f.checkSizes(2, 3)

	if err := m.SetDepthLimit(0); err != nil { // no limit
		t.Fatalf("SetDepthLimit(0): %v", err)
	}
	f.checkChanged(true)
	f.checkSizes(2, 3)
	f.publishAndCheck(true, newBookOrder(t, 4, side.Buy, 1, 3, 0))
	f.publishAndCheck(true, newBookOrder(t, 8, side.Sell, 1, 3, 0))

	if err := m.SetDepthLimit(1); err != nil {
		t.Fatalf("SetDepthLimit(1): %v", err)
	}
	f.checkChanged(true)
	f.checkSizes(1, 1)
	f.publishAndCheck(true, newBookOrder(t, 0, side.Buy, 2, 1, 0))
	f.publishAndCheck(false, newBookOrder(t, 1, side.Buy, 2, 1, 0))
	if err := m.SetDepthLimit(-1); err == nil {
		t.Error("SetDepthLimit(-1): no error")
	}
}

// TestMarketDepthModelRandomBook applies random orders and checks the book after each of them.
func TestMarketDepthModelRandomBook(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f))
	f.publishAndCheck(false, newBookOrder(t, 0, side.Buy, math.NaN(), math.NaN(),
		snapshotBeginEnd|events.RemoveEvent.Flag()))

	const bookSize = 50
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	book := make([]*order.Order, bookSize)
	expectedBuy, expectedSell := 0, 0
	count := func(o *order.Order, s side.Side) int {
		if o != nil && o.Side() == s && o.Size() != 0 {
			return 1
		}
		return 0
	}
	for i := 0; i < 500; i++ {
		index := random.Intn(bookSize)
		var s side.Side = side.Sell
		if random.Intn(2) != 0 {
			s = side.Buy
		}
		value := float64(random.Intn(10))
		o := newBookOrder(t, int64(index), s, value*10, value, 0)
		o.SetScope(order.ScopeOrder)
		old := book[index]
		book[index] = o
		expectedBuy += count(o, side.Buy) - count(old, side.Buy)
		expectedSell += count(o, side.Sell) - count(old, side.Sell)
		f.publish(o)
		f.waitSizes(expectedBuy, expectedSell)
	}
	f.publish(newBookOrder(t, 0, side.Undefined, 0, 0, snapshotBeginEnd|events.RemoveEvent.Flag()))
	f.waitSizes(0, 0)
}

// waitSizes waits for up to 10 s for the book to have the numbers of orders.
func (f *depthFixture) waitSizes(buy, sell int) {
	f.t.Helper()
	for deadline := time.Now().Add(10 * time.Second); ; time.Sleep(time.Millisecond) {
		f.process()
		f.mu.Lock()
		done := len(f.buy) == buy && len(f.sell) == sell
		f.mu.Unlock()
		if done {
			return
		}
		if time.Now().After(deadline) {
			f.checkSizes(buy, sell)
		}
	}
}

func TestMarketDepthModelAggregation(t *testing.T) {
	f := newDepthFixture(t)
	var calls atomic.Int32
	m := buildDepth(f, depthBuilder(f).WithAggregationPeriod(300*time.Millisecond).
		WithListener(func(_, _ []*order.Order) { calls.Add(1) }))
	if got := m.AggregationPeriod(); got != 300*time.Millisecond {
		t.Fatalf("AggregationPeriod: %v, want 300ms", got)
	}
	f.publish(newBookOrder(t, 0, side.Buy, 1, 1, snapshotBeginEnd))
	waitFor(t, func() bool { return calls.Load() == 1 }) // a snapshot at once

	// The incremental changes within the period come in one notification.
	f.publish(newBookOrder(t, 1, side.Buy, 2, 1, 0))
	f.publish(newBookOrder(t, 2, side.Buy, 3, 1, 0))
	waitFor(t, func() bool { return calls.Load() == 2 })
	time.Sleep(500 * time.Millisecond)
	if got := calls.Load(); got != 2 {
		t.Fatalf("listener calls: %d, want 2", got)
	}

	if err := m.SetAggregationPeriod(time.Hour); err != nil {
		t.Fatalf("SetAggregationPeriod: %v", err)
	}
	f.publish(newBookOrder(t, 3, side.Buy, 4, 1, 0))
	time.Sleep(100 * time.Millisecond)
	if got := calls.Load(); got != 2 {
		t.Fatalf("listener calls: %d, want 2", got)
	}
	// Without a period the pending notification is made at once.
	if err := m.SetAggregationPeriod(0); err != nil {
		t.Fatalf("SetAggregationPeriod: %v", err)
	}
	if got := calls.Load(); got != 3 {
		t.Fatalf("listener calls: %d, want 3", got)
	}
	if err := m.SetAggregationPeriod(-time.Second); err == nil {
		t.Error("SetAggregationPeriod(-1s): no error")
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	for deadline := time.Now().Add(10 * time.Second); !condition(); time.Sleep(5 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatal("the condition is not met in 10 s")
		}
	}
}

func TestMarketDepthModelNoListenerCallAfterClose(t *testing.T) {
	f := newDepthFixture(t)
	var calls atomic.Int32
	m := buildDepth(f, depthBuilder(f).WithAggregationPeriod(300*time.Millisecond).
		WithListener(func(_, _ []*order.Order) { calls.Add(1) }))
	f.publish(newBookOrder(t, 0, side.Buy, 1, 1, 0))
	// A slow run may let the timer fire before Close: only the calls after it count.
	before := calls.Load()
	if err := m.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	time.Sleep(600 * time.Millisecond)
	if got := calls.Load(); got != before {
		t.Fatalf("listener calls after Close: %d", got-before)
	}
	if err := m.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}

func TestMarketDepthModelCloseFromListener(t *testing.T) {
	for _, period := range []time.Duration{0, 50 * time.Millisecond} {
		t.Run(period.String(), func(t *testing.T) {
			f := newDepthFixture(t)
			var calls atomic.Int32
			var m *model.MarketDepthModel[*order.Order]
			closed := make(chan error, 1)
			m = buildDepth(f, depthBuilder(f).WithAggregationPeriod(period).
				WithListener(func(_, _ []*order.Order) {
					calls.Add(1)
					select {
					case closed <- m.Close():
					default:
					}
				}))
			f.publish(newBookOrder(t, 0, side.Buy, 1, 1, 0))
			waitFor(t, func() bool { f.process(); return calls.Load() > 0 })
			if err := <-closed; err != nil {
				t.Fatalf("Close in the listener: %v", err)
			}
			f.publish(newBookOrder(t, 1, side.Buy, 2, 1, 0))
			time.Sleep(200 * time.Millisecond)
			if got := calls.Load(); got != 1 {
				t.Fatalf("listener calls: %d, want 1", got)
			}
		})
	}
}

// TestMarketDepthModelCloseWaitsForListener checks that Close on another goroutine returns after the listener call
// in progress.
func TestMarketDepthModelCloseWaitsForListener(t *testing.T) {
	f := newDepthFixture(t)
	entered, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	var callEnded atomic.Bool
	m := buildDepth(f, depthBuilder(f).WithAggregationPeriod(50*time.Millisecond).
		WithListener(func(_, _ []*order.Order) {
			once.Do(func() {
				close(entered)
				<-release
				callEnded.Store(true)
			})
		}))
	f.publish(newBookOrder(t, 0, side.Buy, 1, 1, 0))
	select {
	case <-entered:
	case <-time.After(10 * time.Second):
		t.Fatal("the listener was not called")
	}
	closeReturned := make(chan bool, 1)
	go func() {
		_ = m.Close()
		closeReturned <- callEnded.Load()
	}()
	select {
	case <-closeReturned:
		t.Fatal("Close returned during the listener call")
	case <-time.After(100 * time.Millisecond):
	}
	close(release)
	if ended := <-closeReturned; !ended {
		t.Fatal("Close returned before the listener call ended")
	}
}

func TestMarketDepthModelSnapshotKeepsOtherSources(t *testing.T) {
	f := newDepthFixture(t)
	buildDepth(f, depthBuilder(f).WithSources(order.Default(), order.NtvL3()))
	withSource := func(o *order.Order, source *order.Source) *order.Order {
		o.SetOrderSource(source)
		return o
	}
	ntv := withSource(newBookOrder(t, 0, side.Buy, 1, 1, snapshotBeginEnd), order.NtvL3())
	f.publishAndCheck(true, ntv)
	f.publishAndCheck(true, newBookOrder(t, 1, side.Buy, 3, 1, events.SnapshotBegin.Flag()),
		newBookOrder(t, 0, side.Buy, 2, 1, events.SnapshotEnd.Flag()))
	f.checkSizes(3, 0)

	o := newBookOrder(t, 0, side.Buy, 10, 1, snapshotBeginEnd)
	f.publishAndCheck(true, o)
	f.checkSizes(2, 0)
	f.checkOrder(side.Buy, o, 0)
	f.checkOrder(side.Buy, ntv, 1)
}

func TestMarketDepthModelWithoutListener(t *testing.T) {
	f := newDepthFixture(t)
	m, err := model.NewMarketDepthModelBuilder[*order.Order]().WithFeed(f.feed).WithSymbol(testSymbol).
		WithSources(order.Default()).Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	f.publish(newBookOrder(t, 0, side.Buy, 1, 1, snapshotBeginEnd))
	f.publish(newBookOrder(t, 1, side.Sell, 2, 1, 0))
	if err := m.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestMarketDepthModelBuildErrors(t *testing.T) {
	if _, err := model.NewMarketDepthModelBuilder[*order.Order]().Build(); err == nil {
		t.Error("Build without a symbol: no error")
	}
	if _, err := model.NewMarketDepthModelBuilder[*order.Order]().WithSymbol(testSymbol).WithDepthLimit(-1).
		Build(); err == nil {
		t.Error("Build with a negative depth limit: no error")
	}
	if _, err := model.NewMarketDepthModelBuilder[*order.Order]().WithSymbol(testSymbol).
		WithAggregationPeriod(-time.Millisecond).Build(); err == nil {
		t.Error("Build with a negative aggregation period: no error")
	}
	if _, err := model.NewMarketDepthModelBuilder[model.OrderEvent]().WithSymbol(testSymbol).Build(); err == nil {
		t.Error("Build with an interface order type: no error")
	}
}
