package api

import (
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
)

// expectEvent waits for an event from the collector.
func expectEvent(t *testing.T, collector *eventCollector, what string) interface{} {
	t.Helper()
	select {
	case event := <-collector.events:
		return event
	case <-time.After(10 * time.Second):
		t.Fatalf("%s: no event within 10 s", what)
		return nil
	}
}

// expectNoEvent checks that the collector receives no event within a second.
func expectNoEvent(t *testing.T, collector *eventCollector, what string) {
	t.Helper()
	select {
	case event := <-collector.events:
		t.Fatalf("%s: unexpected event %v", what, event)
	case <-time.After(time.Second):
	}
}

func TestDetachedSubscriptionAttachAndDetach(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "ATTACH/A"
	sub, err := NewDXFeedSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("NewDXFeedSubscription: %v", err)
	}
	defer sub.Close()
	collector := &eventCollector{events: make(chan interface{}, 16)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	if err := sub.AddSymbol(osub.String(symbol)); err != nil {
		t.Fatalf("add symbol: %v", err)
	}
	bid := 0.0
	publish := func() {
		t.Helper()
		bid++
		q := quote.NewQuote(symbol)
		q.SetBidPrice(bid)
		if err := publisher.Publish([]interface{}{q}); err != nil {
			t.Fatalf("publish: %v", err)
		}
	}

	publish()
	expectNoEvent(t, collector, "a detached subscription")
	// Detaching a subscription that is not attached does nothing.
	if err := feed.DetachSubscription(sub); err != nil {
		t.Fatalf("DetachSubscription of a detached subscription: %v", err)
	}
	if err := sub.Detach(feed); err != nil {
		t.Fatalf("Detach of a detached subscription: %v", err)
	}

	steps := []struct {
		name   string
		attach func() error
		detach func() error
	}{
		{"DXFeed.AttachSubscription", func() error { return feed.AttachSubscription(sub) },
			func() error { return feed.DetachSubscription(sub) }},
		{"DXFeedSubscription.Attach", func() error { return sub.Attach(feed) },
			func() error { return sub.Detach(feed) }},
	}
	for _, step := range steps {
		if err := step.attach(); err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
		// Attaching twice does nothing.
		if err := step.attach(); err != nil {
			t.Fatalf("%s again: %v", step.name, err)
		}
		publish()
		if q, ok := expectEvent(t, collector, step.name).(*quote.Quote); !ok || q.BidPrice() != bid {
			t.Fatalf("%s: got %v, want the quote with bid %v", step.name, q, bid)
		}
		if err := step.detach(); err != nil {
			t.Fatalf("detach after %s: %v", step.name, err)
		}
		publish()
		expectNoEvent(t, collector, "detached after "+step.name)
	}
}

func TestDetachSubscriptionAndClear(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "ATTACH/CLEAR"
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer sub.Close()
	collector := &eventCollector{events: make(chan interface{}, 16)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	if err := sub.AddSymbol(osub.String(symbol)); err != nil {
		t.Fatalf("add symbol: %v", err)
	}
	q := quote.NewQuote(symbol)
	q.SetBidPrice(5)
	if err := publisher.Publish([]interface{}{q}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	expectEvent(t, collector, "the published quote")

	if err := feed.DetachSubscriptionAndClear(sub); err != nil {
		t.Fatalf("DetachSubscriptionAndClear: %v", err)
	}
	// The clearing may deliver an empty quote of the symbol; nothing else arrives after the detach.
	deadline := time.After(time.Second)
	for drained := false; !drained; {
		select {
		case event := <-collector.events:
			if cleared, ok := event.(*quote.Quote); !ok || *cleared.EventSymbol() != symbol ||
				!math.IsNaN(cleared.BidPrice()) {
				t.Fatalf("got %v after DetachSubscriptionAndClear, want at most an empty quote of %s", event, symbol)
			}
		case <-deadline:
			drained = true
		}
	}
	q.SetBidPrice(6)
	if err := publisher.Publish([]interface{}{q}); err != nil {
		t.Fatalf("publish: %v", err)
	}
	expectNoEvent(t, collector, "a detached subscription")
}

func TestAttachNilArguments(t *testing.T) {
	feed, _ := newLocalHub(t)
	sub, err := NewDXFeedSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("NewDXFeedSubscription: %v", err)
	}
	defer sub.Close()
	calls := map[string]func() error{
		"AttachSubscription(nil)":         func() error { return feed.AttachSubscription(nil) },
		"DetachSubscription(nil)":         func() error { return feed.DetachSubscription(nil) },
		"DetachSubscriptionAndClear(nil)": func() error { return feed.DetachSubscriptionAndClear(nil) },
		"Attach(nil)":                     func() error { return sub.Attach(nil) },
		"Detach(nil)":                     func() error { return sub.Detach(nil) },
	}
	for name, call := range calls {
		if err := call(); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if _, err := NewDXFeedSubscription(); err == nil {
		t.Error("NewDXFeedSubscription without event types: no error")
	}
}

func TestTimeSeriesSubscriptionFromTime(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "TIMESERIES/A"
	sub, err := feed.CreateTimeSeriesSubscription(eventcodes.TimeAndSale)
	if err != nil {
		t.Fatalf("CreateTimeSeriesSubscription: %v", err)
	}
	defer sub.Close()
	if from, err := sub.FromTime(); err != nil || from != math.MaxInt64 {
		t.Fatalf("FromTime of a new subscription = %d, %v; want math.MaxInt64", from, err)
	}
	base := time.Now().Add(-time.Hour).UnixMilli() / 1000 * 1000
	if err := sub.SetFromTime(base + 1000); err != nil {
		t.Fatalf("SetFromTime: %v", err)
	}
	if from, err := sub.FromTime(); err != nil || from != base+1000 {
		t.Fatalf("FromTime = %d, %v; want %d", from, err, base+1000)
	}
	collector := &eventCollector{events: make(chan interface{}, 16)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	if err := sub.AddSymbol(osub.String(symbol)); err != nil {
		t.Fatalf("add symbol: %v", err)
	}

	// The event before the from-time is not delivered.
	var published []interface{}
	for i := 2; i >= 0; i-- {
		e := timeandsale.NewTimeAndSale(symbol)
		e.SetTime(base + int64(i)*1000)
		published = append(published, e)
	}
	if err := publisher.Publish(published); err != nil {
		t.Fatalf("publish: %v", err)
	}
	for _, want := range []int64{base + 2000, base + 1000} {
		if e, ok := expectEvent(t, collector, "time and sale").(*timeandsale.TimeAndSale); !ok || e.Time() != want {
			t.Fatalf("got %v, want the time and sale at %d", e, want)
		}
	}
	expectNoEvent(t, collector, "the time and sale before the from-time")
}
