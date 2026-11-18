package api

import (
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
)

// TestSnapshotFlagsRoundTrip publishes a snapshot of three time series events, newest first as the snapshot order
// requires, through a local hub and checks that a time series subscription receives the snapshot flags: the first
// event begins the snapshot, the last one ends it, the middle one has no flags.
func TestSnapshotFlagsRoundTrip(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "SNAPSHOT/FLAGS"
	sub, err := feed.CreateSubscription(eventcodes.TimeAndSale)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer sub.Close()
	collector := &eventCollector{events: make(chan interface{}, 16)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	// The snapshot covers the whole subscribed time range, so the hub passes its end on.
	base := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC).UnixMilli()
	from := base - 2000
	if err := sub.AddSymbol(osub.NewTimeSeriesSubscriptionSymbol(symbol, from)); err != nil {
		t.Fatalf("add symbol: %v", err)
	}

	flags := []int32{
		events.SnapshotBegin.Flag(),
		0,
		events.SnapshotEnd.Flag(),
	}
	published := make([]interface{}, len(flags))
	for i, f := range flags {
		e := timeandsale.NewTimeAndSale(symbol)
		e.SetTime(base - int64(i)*1000)
		e.SetPrice(float64(100 + i))
		e.SetEventFlags(f)
		published[i] = e
	}
	if err := publisher.Publish(published); err != nil {
		t.Fatalf("publish: %v", err)
	}

	var received []events.IndexedEvent
	deadline := time.After(10 * time.Second)
	for len(received) == 0 || !events.IsSnapshotEndOrSnip(received[len(received)-1]) {
		select {
		case got := <-collector.events:
			e, ok := got.(events.IndexedEvent)
			if !ok {
				t.Fatalf("received %T, not an indexed event", got)
			}
			received = append(received, e)
		case <-deadline:
			t.Fatalf("snapshot end not received within 10 s; received %v", received)
		}
	}
	if len(received) != len(published) {
		t.Fatalf("received %d events, want %d: %v", len(received), len(published), received)
	}
	for i, e := range received {
		if e.EventFlags() != flags[i] {
			t.Errorf("event %d: flags %#x, want %#x: %v", i, e.EventFlags(), flags[i], e)
		}
	}
	if !events.IsSnapshotBegin(received[0]) || events.IsSnapshotEnd(received[0]) {
		t.Errorf("first event: IsSnapshotBegin %v, IsSnapshotEnd %v", events.IsSnapshotBegin(received[0]),
			events.IsSnapshotEnd(received[0]))
	}
	if !events.IsSnapshotEndOrSnip(received[2]) || events.IsPending(received[2]) || events.IsRemove(received[2]) {
		t.Errorf("last event: flags %#x", received[2].EventFlags())
	}
}
