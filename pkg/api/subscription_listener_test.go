// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// TestRemoveListenerStopsDelivery checks that a removed listener gets no events and the other listener still does.
func TestRemoveListenerStopsDelivery(t *testing.T) {
	feed, publisher := newLocalHub(t)
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer sub.Close()
	kept := &eventCollector{events: make(chan interface{}, 16)}
	removed := &eventCollector{events: make(chan interface{}, 16)}
	for _, l := range []*eventCollector{kept, removed} {
		if err := sub.AddListener(l); err != nil {
			t.Fatalf("add listener: %v", err)
		}
	}
	if err := sub.AddSymbol(osub.String("REMOVE/LISTENER")); err != nil {
		t.Fatalf("add symbol: %v", err)
	}
	publish := func(bid float64) {
		t.Helper()
		q := quote.NewQuote("REMOVE/LISTENER")
		q.SetBidPrice(bid)
		if err := publisher.Publish([]interface{}{q}); err != nil {
			t.Fatalf("publish: %v", err)
		}
	}
	receive := func(c *eventCollector, name string) {
		t.Helper()
		select {
		case <-c.events:
		case <-time.After(10 * time.Second):
			t.Fatalf("%s listener got no event within 10 s", name)
		}
	}

	publish(1)
	receive(kept, "kept")
	receive(removed, "removed")
	sub.RemoveListener(removed)
	publish(2)
	receive(kept, "kept")
	select {
	case got := <-removed.events:
		t.Fatalf("event delivered to a removed listener: %v", got)
	case <-time.After(time.Second):
	}
}

// selfRemovingEventListener removes itself from the subscription on the first events.
type selfRemovingEventListener struct {
	sub   *DXFeedSubscription
	calls chan struct{}
}

func (l *selfRemovingEventListener) Update([]interface{}) {
	l.calls <- struct{}{}
	l.sub.RemoveListener(l)
}

// TestListenerRemovesItself checks that a listener can remove itself from Update without a deadlock.
func TestListenerRemovesItself(t *testing.T) {
	feed, publisher := newLocalHub(t)
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer sub.Close()
	l := &selfRemovingEventListener{sub: sub, calls: make(chan struct{}, 16)}
	if err := sub.AddListener(l); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	if err := sub.AddSymbol(osub.String("REMOVE/SELF")); err != nil {
		t.Fatalf("add symbol: %v", err)
	}
	for i := 0; i < 2; i++ {
		q := quote.NewQuote("REMOVE/SELF")
		q.SetBidPrice(float64(i + 1)) // a changed quote, so it is delivered again
		if err := publisher.Publish([]interface{}{q}); err != nil {
			t.Fatalf("publish: %v", err)
		}
		if i == 0 {
			select {
			case <-l.calls:
			case <-time.After(10 * time.Second):
				t.Fatal("no event within 10 s")
			}
		}
	}
	select {
	case <-l.calls:
		t.Fatal("event delivered after the listener removed itself")
	case <-time.After(time.Second):
	}
}
