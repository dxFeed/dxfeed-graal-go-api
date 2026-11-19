package model_test

import (
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/model"
)

// The scenarios follow the C++ and Java API tests of IndexedTxModel: the events are published through a local hub
// whose executor runs the notifications in process().

const testSymbol = "TEST-SYMBOL"

var snapshotBeginEnd = events.SnapshotBegin.Flag() | events.SnapshotEnd.Flag()

type fixture struct {
	t         *testing.T
	executor  *api.InPlaceExecutor
	feed      *api.DXFeed
	publisher *api.DXPublisher

	mu            sync.Mutex
	notifications int
	snapshots     int
	received      []*order.Order
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	executor, err := api.NewInPlaceExecutor()
	if err != nil {
		t.Fatalf("NewInPlaceExecutor: %v", err)
	}
	t.Cleanup(func() { _ = executor.Close() })
	endpoint, err := api.NewEndpointBuilder().WithRole(api.LocalHub).WithExecutor(executor).Build()
	if err != nil {
		t.Fatalf("build endpoint: %v", err)
	}
	t.Cleanup(func() { _ = endpoint.Close() })
	feed, err := endpoint.GetFeed()
	if err != nil {
		t.Fatalf("GetFeed: %v", err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		t.Fatalf("GetPublisher: %v", err)
	}
	return &fixture{t: t, executor: executor, feed: feed, publisher: publisher}
}

// builder returns a builder of a model of the orders of testSymbol whose listener counts the notifications.
func (f *fixture) builder(withFeed bool) *model.IndexedTxModelBuilder[*order.Order] {
	b := model.NewIndexedTxModelBuilder[*order.Order]().WithSymbol(testSymbol).
		WithListener(func(_ events.IndexedEventSourceInterface, orders []*order.Order, isSnapshot bool) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.notifications++
			if isSnapshot {
				f.snapshots++
			}
			f.received = append(f.received, orders...)
		})
	if withFeed {
		b = b.WithFeed(f.feed)
	}
	return b
}

func (f *fixture) build(b *model.IndexedTxModelBuilder[*order.Order]) *model.IndexedTxModel[*order.Order] {
	f.t.Helper()
	m, err := b.Build()
	if err != nil {
		f.t.Fatalf("Build: %v", err)
	}
	f.t.Cleanup(func() { _ = m.Close() })
	return m
}

func newOrder(t *testing.T, index int64, source *order.Source, size float64, flags int32) *order.Order {
	t.Helper()
	o := order.NewOrder(testSymbol)
	// The source is kept in the high bits of the index, so it is set after the index.
	if err := o.SetIndex(index); err != nil {
		t.Fatalf("SetIndex: %v", err)
	}
	o.SetOrderSource(source)
	o.SetSize(size)
	o.SetEventFlags(flags)
	o.SetSide(side.Buy)
	return o
}

func (f *fixture) publish(orders ...*order.Order) {
	f.t.Helper()
	list := make([]interface{}, len(orders))
	for i, o := range orders {
		list[i] = o
	}
	if err := f.publisher.Publish(list); err != nil {
		f.t.Fatalf("Publish: %v", err)
	}
	f.process()
}

func (f *fixture) process() {
	f.t.Helper()
	if err := f.executor.ProcessAllPendingTasks(); err != nil {
		f.t.Fatalf("ProcessAllPendingTasks: %v", err)
	}
}

// assertChanged checks whether the listener was notified since the last check. A notification is awaited for up to
// 10 s and the rest of them for 100 ms more; their absence is checked for 400 ms.
func (f *fixture) assertChanged(changed bool) {
	f.t.Helper()
	wait := 300 * time.Millisecond
	if changed {
		wait = 10 * time.Second
	}
	for deadline := time.Now().Add(wait); ; time.Sleep(10 * time.Millisecond) {
		f.process()
		f.mu.Lock()
		notified := f.notifications > 0
		f.mu.Unlock()
		if notified || time.Now().After(deadline) {
			break
		}
	}
	// The notifications of one publication may take a few runs of the executor.
	for i := 0; i < 5; i++ {
		time.Sleep(20 * time.Millisecond)
		f.process()
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if (f.notifications > 0) != changed {
		f.t.Fatalf("notifications: %d, want changed %v", f.notifications, changed)
	}
	f.notifications = 0
}

func (f *fixture) assertSnapshots(count int) {
	f.t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.snapshots != count {
		f.t.Fatalf("snapshot notifications: %d, want %d", f.snapshots, count)
	}
	f.snapshots = 0
}

func (f *fixture) assertReceived(count int) {
	f.t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.received) != count {
		f.t.Fatalf("received %d events, want %d", len(f.received), count)
	}
}

// assertEvent checks the oldest received event that was not checked yet.
func (f *fixture) assertEvent(index int64, size float64, flags int32) {
	f.t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.received) == 0 {
		f.t.Fatal("no event received")
	}
	o := f.received[0]
	f.received = f.received[1:]
	if o.EventSymbol() == nil || *o.EventSymbol() != testSymbol || o.Index() != index || o.Size() != size ||
		o.EventFlags() != flags {
		f.t.Fatalf("event %v, want index %d, size %v, flags %#x", o, index, size, flags)
	}
}

func sourceIDs(t *testing.T, m *model.IndexedTxModel[*order.Order]) []int64 {
	t.Helper()
	sources, err := m.Sources()
	if err != nil {
		t.Fatalf("Sources: %v", err)
	}
	ids := make([]int64, 0, len(sources))
	for _, s := range sources {
		ids = append(ids, s.Id())
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func sortedIDs(sources ...*order.Source) []int64 {
	ids := make([]int64, 0, len(sources))
	for _, s := range sources {
		ids = append(ids, s.Id())
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func TestIndexedTxModelInitialState(t *testing.T) {
	if _, err := model.NewIndexedTxModelBuilder[*order.Order]().Build(); err == nil {
		t.Error("Build without a symbol: no error")
	}
	m, err := model.NewIndexedTxModelBuilder[*order.Order]().WithSymbol(testSymbol).Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	defer m.Close()
	if !m.IsBatchProcessing() {
		t.Error("IsBatchProcessing: false, want true")
	}
	if m.IsSnapshotProcessing() {
		t.Error("IsSnapshotProcessing: true, want false")
	}
}

func TestIndexedTxModelChangeSources(t *testing.T) {
	f := newFixture(t)
	dex := newOrder(t, 0, order.DexL3(), 1, snapshotBeginEnd)
	ntvLower := newOrder(t, 0, order.NtvL2(), 2, snapshotBeginEnd)
	ntv := newOrder(t, 0, order.NtvL3(), 3, snapshotBeginEnd)
	m := f.build(f.builder(true).
		WithSources(order.AgregateAsk(), order.AgregateBid()).
		WithSources(order.NtvL2(), order.NtvL3())) // the later sources replace the earlier ones
	if got, want := sourceIDs(t, m), sortedIDs(order.NtvL2(), order.NtvL3()); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Sources: %v, want %v", got, want)
	}

	f.publish(dex) // a source the model is not subscribed to
	f.assertChanged(false)
	f.publish(ntvLower, ntv, dex) // two subscribed sources and one not
	f.assertChanged(true)
	f.assertSnapshots(2)
	f.assertReceived(2)

	if err := m.SetSources(order.DexL3()); err != nil {
		t.Fatalf("SetSources: %v", err)
	}
	if got, want := sourceIDs(t, m), sortedIDs(order.DexL3()); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Sources: %v, want %v", got, want)
	}
	f.publish(dex)
	f.assertChanged(true)
	f.assertSnapshots(1)
	f.assertReceived(3)
}

func TestIndexedTxModelEmptySourcesSubscribeToAll(t *testing.T) {
	f := newFixture(t)
	m := f.build(f.builder(true))
	if got := sourceIDs(t, m); len(got) != 0 {
		t.Fatalf("Sources: %v, want none", got)
	}
	f.publish(newOrder(t, 0, order.NtvL2(), 2, snapshotBeginEnd), newOrder(t, 0, order.NtvL3(), 3, snapshotBeginEnd),
		newOrder(t, 0, order.DexL3(), 1, snapshotBeginEnd))
	f.assertChanged(true)
	f.assertSnapshots(3)
	f.assertReceived(3)
}

// forEachProcessingMode runs the scenario for each combination of the batch and snapshot processing modes.
func forEachProcessingMode(t *testing.T, scenario func(t *testing.T, f *fixture, batch, snapshot bool)) {
	for _, batch := range []bool{false, true} {
		for _, snapshot := range []bool{false, true} {
			batch, snapshot := batch, snapshot
			t.Run(fmt.Sprintf("batch=%v,snapshot=%v", batch, snapshot), func(t *testing.T) {
				scenario(t, newFixture(t), batch, snapshot)
			})
		}
	}
}

func TestIndexedTxModelAttachDetach(t *testing.T) {
	forEachProcessingMode(t, func(t *testing.T, f *fixture, batch, snapshot bool) {
		expectedFlags := snapshotBeginEnd
		if snapshot {
			expectedFlags = 0
		}
		m := f.build(f.builder(false).WithBatchProcessing(batch).WithSnapshotProcessing(snapshot))
		f.publish(newOrder(t, 0, order.Default(), 1, snapshotBeginEnd))
		f.assertChanged(false) // not attached

		if err := m.Attach(f.feed); err != nil {
			t.Fatalf("Attach: %v", err)
		}
		f.publish(newOrder(t, 0, order.Default(), 1, snapshotBeginEnd))
		f.assertChanged(true)
		f.assertReceived(1)

		if err := m.Detach(f.feed); err != nil {
			t.Fatalf("Detach: %v", err)
		}
		f.publish(newOrder(t, 0, order.Default(), 2, snapshotBeginEnd))
		f.assertChanged(false) // detached

		if err := m.Attach(f.feed); err != nil {
			t.Fatalf("Attach: %v", err)
		}
		f.publish(newOrder(t, 0, order.Default(), 3, snapshotBeginEnd))
		f.assertChanged(true)
		f.assertReceived(2)
		f.assertEvent(0, 1, expectedFlags)
		f.assertEvent(0, 3, expectedFlags)
	})
}

func TestIndexedTxModelNoNotificationsAfterClose(t *testing.T) {
	forEachProcessingMode(t, func(t *testing.T, f *fixture, batch, snapshot bool) {
		m := f.build(f.builder(true).WithBatchProcessing(batch).WithSnapshotProcessing(snapshot))
		f.publish(newOrder(t, 0, order.Default(), 12.34, snapshotBeginEnd))
		f.assertChanged(true)
		f.assertSnapshots(1)

		if err := m.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
		if err := m.Close(); err != nil {
			t.Fatalf("second Close: %v", err)
		}
		f.publish(newOrder(t, 2, order.Default(), 56.78, 0))
		f.assertChanged(false)
		if _, err := m.Sources(); err == nil {
			t.Error("Sources after Close: no error")
		}
		if err := m.Attach(f.feed); err == nil {
			t.Error("Attach after Close: no error")
		}
	})
}

// TestIndexedTxModelSnapshotThenIncremental checks that a snapshot comes first and then the incremental updates, and
// that a pending transaction is passed only when it completes.
func TestIndexedTxModelSnapshotThenIncremental(t *testing.T) {
	f := newFixture(t)
	f.build(f.builder(true).WithSnapshotProcessing(true))
	f.publish(newOrder(t, 2, order.Default(), 20, events.SnapshotBegin.Flag()),
		newOrder(t, 1, order.Default(), 10, 0),
		newOrder(t, 0, order.Default(), 0, events.SnapshotEnd.Flag()|events.RemoveEvent.Flag()))
	f.assertChanged(true)
	f.assertSnapshots(1)
	f.assertReceived(2) // snapshot processing drops the removed event
	f.assertEvent(2, 20, 0)
	f.assertEvent(1, 10, 0)

	f.publish(newOrder(t, 1, order.Default(), 11, events.TxPending.Flag()))
	f.assertChanged(false) // the transaction is not complete
	f.publish(newOrder(t, 3, order.Default(), 30, 0))
	f.assertChanged(true)
	f.assertSnapshots(0)
	f.assertReceived(2)
	f.assertEvent(1, 11, events.TxPending.Flag())
	f.assertEvent(3, 30, 0)
}

// TestIndexedTxModelCloseInListener checks that the model can be closed from its listener.
func TestIndexedTxModelCloseInListener(t *testing.T) {
	f := newFixture(t)
	var m *model.IndexedTxModel[*order.Order]
	closed := make(chan error, 1)
	m, err := model.NewIndexedTxModelBuilder[*order.Order]().WithFeed(f.feed).WithSymbol(testSymbol).
		WithListener(func(events.IndexedEventSourceInterface, []*order.Order, bool) {
			select {
			case closed <- m.Close():
			default:
			}
		}).Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	defer m.Close()
	f.publish(newOrder(t, 0, order.Default(), 1, snapshotBeginEnd))
	for deadline := time.Now().Add(10 * time.Second); len(closed) == 0 && time.Now().Before(deadline); {
		f.process()
		time.Sleep(10 * time.Millisecond)
	}
	select {
	case err := <-closed:
		if err != nil {
			t.Fatalf("Close in the listener: %v", err)
		}
	default:
		t.Fatal("the listener was not called")
	}
}

func TestIndexedTxModelBuildErrors(t *testing.T) {
	if _, err := model.NewIndexedTxModelBuilder[*order.Order]().WithSymbol(nil).Build(); err == nil {
		t.Error("Build with a nil symbol: no error")
	}
	if _, err := model.NewIndexedTxModelBuilder[*order.Order]().WithSymbol(42).Build(); err == nil {
		t.Error("Build with an unsupported symbol: no error")
	}
	executor, err := api.NewInPlaceExecutor()
	if err != nil {
		t.Fatalf("NewInPlaceExecutor: %v", err)
	}
	_ = executor.Close()
	if _, err := model.NewIndexedTxModelBuilder[*order.Order]().WithSymbol(testSymbol).WithExecutor(executor).
		Build(); err == nil {
		t.Error("Build with a closed executor: no error")
	}
}
