// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// changeRecorder records the calls of a change listener as "added:A,B", "removed:A" and "closed".
type changeRecorder struct {
	mu    sync.Mutex
	calls []string
	// onAdded, if set, is called from SymbolsAdded after the call is recorded.
	onAdded func()
}

func (r *changeRecorder) record(call string) {
	r.mu.Lock()
	r.calls = append(r.calls, call)
	r.mu.Unlock()
}

func symbolsText(symbols []osub.Symbol) string {
	texts := make([]string, 0, len(symbols))
	for _, s := range symbols {
		texts = append(texts, fmt.Sprint(s))
	}
	sort.Strings(texts)
	return strings.Join(texts, ",")
}

func (r *changeRecorder) SymbolsAdded(symbols []osub.Symbol) {
	r.record("added:" + symbolsText(symbols))
	if r.onAdded != nil {
		r.onAdded()
	}
}

func (r *changeRecorder) SymbolsRemoved(symbols []osub.Symbol) {
	r.record("removed:" + symbolsText(symbols))
}

func (r *changeRecorder) SubscriptionClosed() {
	r.record("closed")
}

// take returns the recorded calls and forgets them.
func (r *changeRecorder) take() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	calls := r.calls
	r.calls = nil
	return calls
}

// waitFor waits until a recorded call is want and returns the calls up to it.
func (r *changeRecorder) waitFor(t *testing.T, want string) []string {
	t.Helper()
	var seen []string
	for deadline := time.Now().Add(10 * time.Second); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		for _, call := range r.take() {
			seen = append(seen, call)
			if call == want {
				return seen
			}
		}
	}
	t.Fatalf("no %q within 10 s; calls %q", want, seen)
	return nil
}

// configRecorder is a changeRecorder that also records ConfigurationChanged as "config".
type configRecorder struct {
	*changeRecorder
}

func (r configRecorder) ConfigurationChanged() {
	r.record("config")
}

func expectCalls(t *testing.T, what string, got []string, want ...string) {
	t.Helper()
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("%s: calls %q, want %q", what, got, want)
	}
}

// TestChangeListenerSeesSymbolChanges checks the calls of a change listener of a feed subscription: they come on the
// goroutine that changes the subscription, before the change returns.
func TestChangeListenerSeesSymbolChanges(t *testing.T) {
	sub := newSubscription(t, eventcodes.Quote)
	if err := sub.AddSymbols(osub.Strings("AAPL", "IBM")...); err != nil {
		t.Fatalf("AddSymbols: %v", err)
	}
	withConfig := configRecorder{&changeRecorder{}}
	plain := &changeRecorder{}
	for _, l := range []osub.ObservableSubscriptionChangeListener{withConfig, plain} {
		if err := sub.AddChangeListener(l); err != nil {
			t.Fatalf("AddChangeListener: %v", err)
		}
	}
	expectCalls(t, "AddChangeListener", withConfig.take(), "added:AAPL,IBM")
	expectCalls(t, "AddChangeListener", plain.take(), "added:AAPL,IBM")

	if err := sub.AddSymbol(osub.String("MSFT")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	if err := sub.RemoveSymbols(osub.Strings("AAPL", "IBM")...); err != nil {
		t.Fatalf("RemoveSymbols: %v", err)
	}
	if err := sub.SetAggregationPeriod(time.Second); err != nil {
		t.Fatalf("SetAggregationPeriod: %v", err)
	}
	expectCalls(t, "changes", withConfig.take(), "added:MSFT", "removed:AAPL,IBM", "config")
	expectCalls(t, "changes", plain.take(), "added:MSFT", "removed:AAPL,IBM")

	sub.RemoveChangeListener(withConfig)
	expectCalls(t, "RemoveChangeListener", withConfig.take(), "closed")
	if err := sub.AddSymbol(osub.String("TSLA")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	expectCalls(t, "after RemoveChangeListener", withConfig.take())
	expectCalls(t, "the other listener", plain.take(), "added:TSLA")

	if err := sub.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	expectCalls(t, "Close", plain.take(), "closed")
	if err := sub.AddChangeListener(plain); err != nil {
		t.Errorf("AddChangeListener of a closed subscription: %v", err)
	}
	expectCalls(t, "AddChangeListener of a closed subscription", plain.take())
}

// TestChangeListenerRemovesItself checks that a change listener may remove itself from its calls, both from the call
// of AddChangeListener and from a later one.
func TestChangeListenerRemovesItself(t *testing.T) {
	sub := newSubscription(t, eventcodes.Quote)
	if err := sub.AddSymbol(osub.String("AAPL")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	first := &changeRecorder{}
	first.onAdded = func() { sub.RemoveChangeListener(first) }
	if err := sub.AddChangeListener(first); err != nil {
		t.Fatalf("AddChangeListener: %v", err)
	}
	expectCalls(t, "removed in the call of AddChangeListener", first.take(), "added:AAPL", "closed")

	later := &changeRecorder{}
	if err := sub.AddChangeListener(later); err != nil {
		t.Fatalf("AddChangeListener: %v", err)
	}
	later.onAdded = func() { sub.RemoveChangeListener(later) }
	if err := sub.AddSymbol(osub.String("IBM")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	if err := sub.AddSymbol(osub.String("MSFT")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	expectCalls(t, "removed in a later call", later.take(), "added:AAPL", "added:IBM", "closed")
	expectCalls(t, "after the removal", first.take())
}

// panickingChangeListener panics on its first SymbolsAdded.
type panickingChangeListener struct {
	*changeRecorder
	once sync.Once
}

func (l *panickingChangeListener) SymbolsAdded(symbols []osub.Symbol) {
	panicked := false
	l.once.Do(func() { panicked = true })
	if panicked {
		panic("change listener failure")
	}
	l.changeRecorder.SymbolsAdded(symbols)
}

// TestChangeListenerPanicIsRecovered checks that a panic of a change listener is logged, does not end the process and
// does not stop the later calls.
func TestChangeListenerPanicIsRecovered(t *testing.T) {
	logged := &syncBuffer{}
	defer log.SetOutput(log.Writer())
	log.SetOutput(logged)

	sub := newSubscription(t, eventcodes.Quote)
	listener := &panickingChangeListener{changeRecorder: &changeRecorder{}}
	if err := sub.AddChangeListener(listener); err != nil {
		t.Fatalf("AddChangeListener: %v", err)
	}
	for _, symbol := range []string{"AAPL", "IBM"} {
		if err := sub.AddSymbol(osub.String(symbol)); err != nil {
			t.Fatalf("AddSymbol: %v", err)
		}
	}
	expectCalls(t, "after the panic", listener.take(), "added:IBM")
	if !strings.Contains(logged.String(), "change listener failure") {
		t.Errorf("the panic is not logged; log:\n%s", logged.String())
	}
}

// TestPublisherSubscription checks that the subscription of a publisher reports the symbols that a feed subscribes
// to.
func TestPublisherSubscription(t *testing.T) {
	feed, publisher := newLocalHub(t)
	observed, err := publisher.Subscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	if again, err := publisher.Subscription(eventcodes.Quote); err != nil || again != observed {
		t.Errorf("a second Subscription(Quote) = %p, %v; want the same %p", again, err, observed)
	}
	if closed, err := observed.IsClosed(); err != nil || closed {
		t.Errorf("IsClosed() = %v, %v; want false", closed, err)
	}
	if types, err := observed.EventTypes(); err != nil || len(types) != 1 || types[0] != eventcodes.Quote {
		t.Errorf("EventTypes() = %v, %v; want [Quote]", types, err)
	}
	for code, want := range map[eventcodes.EventCode]bool{eventcodes.Quote: true, eventcodes.Trade: false} {
		if got, err := observed.ContainsEventType(code); err != nil || got != want {
			t.Errorf("ContainsEventType(%v) = %v, %v; want %v", code, got, err, want)
		}
	}

	listener := configRecorder{&changeRecorder{}}
	if err := observed.AddChangeListener(listener); err != nil {
		t.Fatalf("AddChangeListener: %v", err)
	}
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("CreateSubscription: %v", err)
	}
	defer func() { _ = sub.Close() }()
	if err := sub.AddSymbol(osub.String("PUBLISHED")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	listener.waitFor(t, "added:PUBLISHED")
	if err := sub.RemoveSymbol(osub.String("PUBLISHED")); err != nil {
		t.Fatalf("RemoveSymbol: %v", err)
	}
	listener.waitFor(t, "removed:PUBLISHED")

	observed.RemoveChangeListener(listener)
	listener.waitFor(t, "closed")
	if err := sub.AddSymbol(osub.String("LATER")); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	expectCalls(t, "after RemoveChangeListener", listener.take())
}

// TestPublisherTimeSeriesSubscription checks that the subscription of a publisher reports a time series subscription
// with its time.
func TestPublisherTimeSeriesSubscription(t *testing.T) {
	feed, publisher := newLocalHub(t)
	observed, err := publisher.Subscription(eventcodes.TimeAndSale)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	symbols := make(chan osub.Symbol, 16)
	listener := &symbolsListener{added: symbols}
	if err := observed.AddChangeListener(listener); err != nil {
		t.Fatalf("AddChangeListener: %v", err)
	}
	defer observed.RemoveChangeListener(listener)
	sub, err := feed.CreateSubscription(eventcodes.TimeAndSale)
	if err != nil {
		t.Fatalf("CreateSubscription: %v", err)
	}
	defer func() { _ = sub.Close() }()
	if err := sub.AddSymbol(osub.NewTimeSeriesSubscriptionSymbol("IBM", 1000)); err != nil {
		t.Fatalf("AddSymbol: %v", err)
	}
	timeout := time.After(10 * time.Second)
	for {
		select {
		case s := <-symbols:
			if ts, ok := s.(*osub.TimeSeriesSubscriptionSymbol); ok && ts.Symbol() == "IBM" {
				if ts.FromTime() != 1000 {
					t.Errorf("FromTime() = %d, want 1000", ts.FromTime())
				}
				return
			}
		case <-timeout:
			t.Fatal("no time series symbol within 10 s")
		}
	}
}

// symbolsListener passes the added symbols to a channel.
type symbolsListener struct {
	added chan osub.Symbol
}

func (l *symbolsListener) SymbolsAdded(symbols []osub.Symbol) {
	for _, s := range symbols {
		l.added <- s
	}
}

func (l *symbolsListener) SymbolsRemoved([]osub.Symbol) {}

func (l *symbolsListener) SubscriptionClosed() {}

// TestPublisherSubscriptionAfterEndpointClose checks that the subscription of a publisher reports itself closed after
// the endpoint is closed, and that its listener can still be removed, and that GetPublisher returns one publisher.
func TestPublisherSubscriptionAfterEndpointClose(t *testing.T) {
	endpoint, err := NewEndpoint(LocalHub)
	if err != nil {
		t.Fatalf("NewEndpoint: %v", err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		t.Fatalf("GetPublisher: %v", err)
	}
	if again, err := endpoint.GetPublisher(); err != nil || again != publisher {
		t.Errorf("a second GetPublisher() = %p, %v; want the same %p", again, err, publisher)
	}
	observed, err := publisher.Subscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	listener := &changeRecorder{}
	if err := observed.AddChangeListener(listener); err != nil {
		t.Fatalf("AddChangeListener: %v", err)
	}
	if err := endpoint.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if closed, err := observed.IsClosed(); err != nil || !closed {
		t.Errorf("IsClosed() after Close = %v, %v; want true", closed, err)
	}
	if _, err := observed.EventTypes(); err == nil {
		t.Error("EventTypes() after Close: no error")
	}
	observed.RemoveChangeListener(listener)
}
