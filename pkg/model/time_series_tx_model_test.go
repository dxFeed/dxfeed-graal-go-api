package model_test

import (
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/model"
)

// The scenarios follow the C++ API tests of TimeSeriesTxModel.

type tnsFixture = fixture[*timeandsale.TimeAndSale]

// tnsBuilder returns a builder of a model of the trades of testSymbol from time 0 whose listener counts the
// notifications.
func tnsBuilder(f *tnsFixture, withFeed bool) *model.TimeSeriesTxModelBuilder[*timeandsale.TimeAndSale] {
	b := model.NewTimeSeriesTxModelBuilder[*timeandsale.TimeAndSale]().WithSymbol(testSymbol).WithFromTime(0).
		WithListener(f.listen)
	if withFeed {
		b = b.WithFeed(f.feed)
	}
	return b
}

func buildTns(f *tnsFixture,
	b *model.TimeSeriesTxModelBuilder[*timeandsale.TimeAndSale]) *model.TimeSeriesTxModel[*timeandsale.TimeAndSale] {
	f.t.Helper()
	m, err := b.Build()
	if err != nil {
		f.t.Fatalf("Build: %v", err)
	}
	f.t.Cleanup(func() { _ = m.Close() })
	return m
}

func newTns(index int64, size float64, flags int32) *timeandsale.TimeAndSale {
	tns := timeandsale.NewTimeAndSale(testSymbol)
	tns.SetIndex(index)
	tns.SetSize(size)
	tns.SetEventFlags(flags)
	return tns
}

func newTnsAt(time int64, size float64, flags int32) *timeandsale.TimeAndSale {
	tns := timeandsale.NewTimeAndSale(testSymbol)
	tns.SetTime(time)
	tns.SetSize(size)
	tns.SetEventFlags(flags)
	return tns
}

func fromTime(t *testing.T, m *model.TimeSeriesTxModel[*timeandsale.TimeAndSale]) int64 {
	t.Helper()
	value, err := m.FromTime()
	if err != nil {
		t.Fatalf("FromTime: %v", err)
	}
	return value
}

func setFromTime(t *testing.T, m *model.TimeSeriesTxModel[*timeandsale.TimeAndSale], value int64) {
	t.Helper()
	if err := m.SetFromTime(value); err != nil {
		t.Fatalf("SetFromTime(%d): %v", value, err)
	}
}

func TestTimeSeriesTxModelInitialState(t *testing.T) {
	if _, err := model.NewTimeSeriesTxModelBuilder[*timeandsale.TimeAndSale]().Build(); err == nil {
		t.Error("Build without a symbol: no error")
	}
	m, err := model.NewTimeSeriesTxModelBuilder[*timeandsale.TimeAndSale]().WithSymbol(testSymbol).Build()
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
	if got := fromTime(t, m); got != math.MaxInt64 {
		t.Errorf("FromTime: %d, want math.MaxInt64", got)
	}
}

// TestTimeSeriesTxModelNegativeFromTime checks that -1 is a time, not a failure.
func TestTimeSeriesTxModelNegativeFromTime(t *testing.T) {
	f := newFixture[*timeandsale.TimeAndSale](t)
	m := buildTns(f, tnsBuilder(f, true).WithFromTime(-1))
	if got := fromTime(t, m); got != -1 {
		t.Fatalf("FromTime: %d, want -1", got)
	}
	setFromTime(t, m, 0)
	if got := fromTime(t, m); got != 0 {
		t.Fatalf("FromTime: %d, want 0", got)
	}
	setFromTime(t, m, -1)
	if got := fromTime(t, m); got != -1 {
		t.Fatalf("FromTime: %d, want -1", got)
	}
}

func TestTimeSeriesTxModelChangeFromTime(t *testing.T) {
	f := newFixture[*timeandsale.TimeAndSale](t)
	m := buildTns(f, tnsBuilder(f, true).WithFromTime(math.MaxInt64)) // not subscribed

	f.publish(newTnsAt(0, 1, snapshotBeginEnd))
	f.assertChanged(false) // not subscribed

	setFromTime(t, m, 0)
	f.publish(newTnsAt(0, 2, snapshotBeginEnd))
	f.assertChanged(true)

	setFromTime(t, m, 1000)
	if got := fromTime(t, m); got != 1000 {
		t.Fatalf("FromTime: %d, want 1000", got)
	}
	f.publish(newTnsAt(1000, 3, 0))
	f.assertChanged(true)

	setFromTime(t, m, 1000)
	f.publish(newTnsAt(0, 3, 0))
	f.assertChanged(false) // an event before the time of the subscription

	setFromTime(t, m, math.MaxInt64)
	f.publish(newTnsAt(0, 4, snapshotBeginEnd))
	f.assertChanged(false) // not subscribed
}

func TestTimeSeriesTxModelAttachDetach(t *testing.T) {
	forEachProcessingMode(t, func(t *testing.T, f *tnsFixture, batch, snapshot bool) {
		expectedFlags := snapshotBeginEnd
		if snapshot {
			expectedFlags = 0
		}
		m := buildTns(f, tnsBuilder(f, false).WithBatchProcessing(batch).WithSnapshotProcessing(snapshot))
		if m.IsBatchProcessing() != batch || m.IsSnapshotProcessing() != snapshot {
			t.Fatalf("processing: batch %v, snapshot %v", m.IsBatchProcessing(), m.IsSnapshotProcessing())
		}
		if got := fromTime(t, m); got != 0 {
			t.Fatalf("FromTime: %d, want 0", got)
		}
		f.publish(newTns(0, 1, snapshotBeginEnd))
		f.assertChanged(false) // not attached

		if err := m.Attach(f.feed); err != nil {
			t.Fatalf("Attach: %v", err)
		}
		f.publish(newTns(0, 1, snapshotBeginEnd))
		f.assertChanged(true)
		f.assertReceived(1)

		if err := m.Detach(f.feed); err != nil {
			t.Fatalf("Detach: %v", err)
		}
		f.publish(newTns(0, 2, snapshotBeginEnd))
		f.assertChanged(false) // detached

		if err := m.Attach(f.feed); err != nil {
			t.Fatalf("Attach: %v", err)
		}
		f.publish(newTns(0, 3, snapshotBeginEnd))
		f.assertChanged(true)
		f.assertReceived(2)
		f.assertEvent(0, 1, expectedFlags)
		f.assertEvent(0, 3, expectedFlags)
	})
}

func TestTimeSeriesTxModelNoNotificationsAfterClose(t *testing.T) {
	forEachProcessingMode(t, func(t *testing.T, f *tnsFixture, batch, snapshot bool) {
		m := buildTns(f, tnsBuilder(f, true).WithBatchProcessing(batch).WithSnapshotProcessing(snapshot))
		f.publish(newTns(0, 12.34, snapshotBeginEnd))
		f.assertChanged(true)
		f.assertSnapshots(1)

		if err := m.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
		if err := m.Close(); err != nil {
			t.Fatalf("second Close: %v", err)
		}
		f.publish(newTns(2, 56.78, 0))
		f.assertChanged(false)
		if _, err := m.FromTime(); err == nil {
			t.Error("FromTime after Close: no error")
		}
		if err := m.SetFromTime(0); err == nil {
			t.Error("SetFromTime after Close: no error")
		}
	})
}

func TestTimeSeriesTxModelCloseInListener(t *testing.T) {
	f := newFixture[*timeandsale.TimeAndSale](t)
	var m *model.TimeSeriesTxModel[*timeandsale.TimeAndSale]
	closed := make(chan error, 1)
	m, err := model.NewTimeSeriesTxModelBuilder[*timeandsale.TimeAndSale]().WithFeed(f.feed).WithSymbol(testSymbol).
		WithFromTime(0).WithListener(func([]*timeandsale.TimeAndSale, bool) {
		select {
		case closed <- m.Close():
		default:
		}
	}).Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	defer m.Close()
	f.publish(newTns(0, 1, snapshotBeginEnd))
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

// TestTimeSeriesTxModelSnapshotThenIncremental checks that a snapshot comes first and then the incremental updates.
func TestTimeSeriesTxModelSnapshotThenIncremental(t *testing.T) {
	f := newFixture[*timeandsale.TimeAndSale](t)
	buildTns(f, tnsBuilder(f, true).WithSnapshotProcessing(true))
	f.publish(newTns(2, 20, events.SnapshotBegin.Flag()),
		newTns(1, 10, 0),
		newTns(0, 0, events.SnapshotEnd.Flag()|events.RemoveEvent.Flag()))
	f.assertChanged(true)
	f.assertSnapshots(1)
	f.assertReceived(2) // snapshot processing drops the removed event
	f.assertEvent(2, 20, 0)
	f.assertEvent(1, 10, 0)

	f.publish(newTns(3, 30, 0))
	f.assertChanged(true)
	f.assertSnapshots(0)
	f.assertReceived(1)
	f.assertEvent(3, 30, 0)
}

func TestTimeSeriesTxModelBuildErrors(t *testing.T) {
	if _, err := model.NewTimeSeriesTxModelBuilder[events.TimeSeriesEvent]().WithSymbol(testSymbol).
		Build(); err == nil {
		t.Error("Build with an interface event type: no error")
	}
	if _, err := model.NewTimeSeriesTxModelBuilder[*timeandsale.TimeAndSale]().WithSymbol(42).Build(); err == nil {
		t.Error("Build with an unsupported symbol: no error")
	}
}
