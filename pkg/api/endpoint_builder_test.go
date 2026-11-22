// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

func buildEndpoint(t *testing.T, b *EndpointBuilder) *DXEndpoint {
	t.Helper()
	endpoint, err := b.Build()
	if err != nil {
		t.Fatalf("build endpoint: %v", err)
	}
	t.Cleanup(func() { _ = endpoint.Close() })
	return endpoint
}

func wantPeriod(t *testing.T, endpoint *DXEndpoint, want time.Duration) {
	t.Helper()
	period, ok, err := endpoint.RequestedAggregationPeriod()
	if err != nil || !ok || period != want {
		t.Errorf("RequestedAggregationPeriod() = %v, %v, %v; want %v", period, ok, err, want)
	}
}

func TestEndpointBuilder(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder().
		WithRole(LocalHub).
		WithName("builder-test").
		WithProperties(map[string]string{"dxfeed.requestedAggregationPeriod": "1s"}))
	wantPeriod(t, endpoint, time.Second)
	if endpoint.Role() != LocalHub {
		t.Errorf("Role() = %v, want LocalHub", endpoint.Role())
	}
	types, err := endpoint.EventTypes()
	if err != nil {
		t.Fatalf("EventTypes: %v", err)
	}
	found := false
	for _, code := range types {
		found = found || code == eventcodes.Quote
	}
	if !found {
		t.Errorf("EventTypes() = %v, want Quote among them", types)
	}
}

func TestEndpointBuilderSupportsProperty(t *testing.T) {
	b := NewEndpointBuilder()
	supported, err := b.SupportsProperty("dxfeed.address")
	if err != nil || !supported {
		t.Errorf(`SupportsProperty("dxfeed.address") = %v, %v; want true`, supported, err)
	}
	supported, err = b.SupportsProperty("no.such.property")
	if err != nil || supported {
		t.Errorf(`SupportsProperty("no.such.property") = %v, %v; want false`, supported, err)
	}
}

func TestEndpointBuilderPropertiesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dxfeed.properties")
	if err := os.WriteFile(path, []byte("dxfeed.requestedAggregationPeriod=2s\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	wantPeriod(t, buildEndpoint(t, NewEndpointBuilder().WithPropertiesFile(path)), 2*time.Second)
	// The steps apply in the order of the calls: a later property wins over the file and over an earlier property.
	wantPeriod(t, buildEndpoint(t, NewEndpointBuilder().
		WithPropertiesFile(path).
		WithProperty("dxfeed.requestedAggregationPeriod", "3s")), 3*time.Second)
	wantPeriod(t, buildEndpoint(t, NewEndpointBuilder().
		WithProperty("dxfeed.requestedAggregationPeriod", "3s").
		WithPropertiesFile(path)), 2*time.Second)

	_, err := NewEndpointBuilder().WithPropertiesFile(filepath.Join(t.TempDir(), "missing.properties")).Build()
	if err == nil {
		t.Error("Build with a missing properties file: no error")
	}
}

// TestEndpointStateTransitions goes through the states with Connect, Disconnect and Close.
func TestEndpointStateTransitions(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	wantState := func(want common.ConnectionState) {
		t.Helper()
		state, err := endpoint.State()
		if err != nil || state != want {
			t.Fatalf("State() = %v, %v; want %v", state, err, want)
		}
	}
	wantState(common.NotConnected)
	if err := endpoint.User("demo"); err != nil {
		t.Fatalf("User: %v", err)
	}
	if err := endpoint.Password("demo"); err != nil {
		t.Fatalf("Password: %v", err)
	}
	// Nothing listens on port 1, so the endpoint stays in Connecting.
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	wantState(common.Connecting)
	if err := endpoint.Reconnect(); err != nil {
		t.Fatalf("Reconnect: %v", err)
	}
	if err := endpoint.Disconnect(); err != nil {
		t.Fatalf("Disconnect: %v", err)
	}
	wantState(common.NotConnected)
	if err := endpoint.DisconnectAndClear(); err != nil {
		t.Fatalf("DisconnectAndClear: %v", err)
	}
	if endpoint.IsClosed() {
		t.Fatal("IsClosed() = true before Close")
	}
	if err := endpoint.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	wantState(common.Closed)
	if !endpoint.IsClosed() {
		t.Fatal("IsClosed() = false after Close")
	}
}

func TestEndpointAggregationPeriod(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	if _, ok, err := endpoint.RequestedAggregationPeriod(); err != nil || ok {
		t.Fatalf("RequestedAggregationPeriod() of a new endpoint: ok %v, err %v; want none", ok, err)
	}
	if err := endpoint.SetRequestedAggregationPeriod(1500 * time.Millisecond); err != nil {
		t.Fatalf("SetRequestedAggregationPeriod: %v", err)
	}
	period, ok, err := endpoint.RequestedAggregationPeriod()
	if err != nil || !ok || period != 1500*time.Millisecond {
		t.Fatalf("RequestedAggregationPeriod() = %v, %v, %v; want 1.5s", period, ok, err)
	}
	if err := endpoint.SetRequestedAggregationPeriod(0); err != nil {
		t.Fatalf("SetRequestedAggregationPeriod(0): %v", err)
	}
	if period, ok, _ := endpoint.RequestedAggregationPeriod(); !ok || period != 0 {
		t.Fatalf("RequestedAggregationPeriod() = %v, %v; want 0 (real time)", period, ok)
	}
	if err := endpoint.SetRequestedAggregationPeriod(500 * time.Microsecond); err != nil {
		t.Fatalf("SetRequestedAggregationPeriod(500us): %v", err)
	}
	wantPeriod(t, endpoint, time.Millisecond)
	for _, bad := range []time.Duration{-time.Second, 25 * time.Hour} {
		if err := endpoint.SetRequestedAggregationPeriod(bad); err == nil {
			t.Errorf("SetRequestedAggregationPeriod(%v): no error", bad)
		}
	}
	if err := endpoint.ClearRequestedAggregationPeriod(); err != nil {
		t.Fatalf("ClearRequestedAggregationPeriod: %v", err)
	}
	if _, ok, err := endpoint.RequestedAggregationPeriod(); err != nil || ok {
		t.Fatalf("RequestedAggregationPeriod() after clear: ok %v, err %v; want none", ok, err)
	}
	info, err := endpoint.AggregationPeriodInfo()
	if err != nil || info.Known {
		t.Fatalf("AggregationPeriodInfo() of an unconnected endpoint = %+v, %v; want unknown", info, err)
	}
}

// TestInPlaceExecutor checks that the events reach the listener only in ProcessAllPendingTasks.
func TestInPlaceExecutor(t *testing.T) {
	executor, err := NewInPlaceExecutor()
	if err != nil {
		t.Fatalf("NewInPlaceExecutor: %v", err)
	}
	defer func() { _ = executor.Close() }()
	endpoint := buildEndpoint(t, NewEndpointBuilder().WithRole(LocalHub).WithExecutor(executor))
	feed, err := endpoint.GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		t.Fatal(err)
	}
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sub.Close() }()
	collector := &eventCollector{events: make(chan interface{}, 10)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatal(err)
	}
	if err := sub.AddSymbol(osub.String("IBM")); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for published := false; ; {
		if err := executor.ProcessAllPendingTasks(); err != nil {
			t.Fatalf("ProcessAllPendingTasks: %v", err)
		}
		select {
		case <-collector.events:
			return
		default:
		}
		if !published {
			// The subscription reaches the publisher through the executor too, so publish after the first run.
			if err := publisher.Publish([]interface{}{quote.NewQuote("IBM")}); err != nil {
				t.Fatal(err)
			}
			// AwaitProcessed would wait for the tasks that only ProcessAllPendingTasks runs: give the SDK time instead.
			time.Sleep(200 * time.Millisecond)
			select {
			case <-collector.events:
				t.Fatal("an event arrived outside ProcessAllPendingTasks")
			default:
			}
			published = true
		}
		if time.Now().After(deadline) {
			t.Fatal("no event within 5 s")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestInPlaceExecutorCloseAndAwaitTermination checks that closing an endpoint with an in-place executor ends although
// nobody runs its pending tasks: since SDK 3.19.0 the closing call runs them itself while it waits, so the listener
// gets the pending event on the closing goroutine.
func TestInPlaceExecutorCloseAndAwaitTermination(t *testing.T) {
	executor, err := NewInPlaceExecutor()
	if err != nil {
		t.Fatalf("NewInPlaceExecutor: %v", err)
	}
	defer func() { _ = executor.Close() }()
	endpoint := buildEndpoint(t, NewEndpointBuilder().WithRole(LocalHub).WithExecutor(executor))
	feed, err := endpoint.GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sub.Close() }()
	collector := &eventCollector{events: make(chan interface{}, 10)}
	if err := sub.AddListener(collector); err != nil {
		t.Fatal(err)
	}
	if err := sub.AddSymbol(osub.String("IBM")); err != nil {
		t.Fatal(err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		t.Fatal(err)
	}
	// The subscription and the events reach the listener only through the executor: run the subscription once, then
	// leave the notification of the published event pending.
	if err := executor.ProcessAllPendingTasks(); err != nil {
		t.Fatalf("ProcessAllPendingTasks: %v", err)
	}
	if err := publisher.Publish([]interface{}{quote.NewQuote("IBM")}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := endpoint.CloseAndAwaitTermination(ctx); err != nil {
		t.Fatalf("CloseAndAwaitTermination: %v", err)
	}
	select {
	case <-collector.events:
	default:
		t.Error("the pending event did not reach the listener while the endpoint was closed")
	}
}

func TestClosedExecutor(t *testing.T) {
	executor, err := NewFixedThreadPoolExecutor(1, "closed-test")
	if err != nil {
		t.Fatalf("NewFixedThreadPoolExecutor: %v", err)
	}
	endpoint := buildEndpoint(t, NewEndpointBuilder().WithExecutor(executor))
	if endpoint.Role() != Feed {
		t.Errorf("Role() = %v, want Feed", endpoint.Role())
	}
	if err := executor.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := NewEndpointBuilder().WithExecutor(executor).Build(); err == nil {
		t.Error("Build with a closed executor: no error")
	}
	var nilExecutor *InPlaceExecutor
	if _, err := NewEndpointBuilder().WithExecutor(nilExecutor).Build(); err == nil {
		t.Error("Build with a nil *InPlaceExecutor: no error")
	}
	inPlace, err := NewInPlaceExecutor()
	if err != nil {
		t.Fatal(err)
	}
	_ = inPlace.Close()
	if err := inPlace.ProcessAllPendingTasks(); err == nil {
		t.Error("ProcessAllPendingTasks after Close: no error")
	}
}
