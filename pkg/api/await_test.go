// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func TestAwaitWithContextReturnsWaitResult(t *testing.T) {
	want := errors.New("wait failed")
	if err := awaitWithContext(context.Background(), func() error { return want }); !errors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
}

func TestAwaitWithContextCancelled(t *testing.T) {
	release := make(chan struct{})
	defer close(release)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := awaitWithContext(ctx, func() error { <-release; return nil })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestAwaitWithContextDoneBeforeCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	if err := awaitWithContext(ctx, func() error { called = true; return nil }); !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want %v", err, context.Canceled)
	}
	if called {
		t.Fatal("the wait was started on a done context")
	}
}

func TestEndpointWaitsWithContext(t *testing.T) {
	endpoint, err := NewEndpoint(LocalHub)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := endpoint.AwaitProcessed(ctx); err != nil {
		t.Fatalf("AwaitProcessed: %v", err)
	}
	if err := endpoint.AwaitNotConnected(ctx); err != nil {
		t.Fatalf("AwaitNotConnected: %v", err)
	}
	if err := endpoint.CloseAndAwaitTermination(ctx); err != nil {
		t.Fatalf("CloseAndAwaitTermination: %v", err)
	}
}

func TestCloseAndAwaitTerminationClosesOnDoneContext(t *testing.T) {
	endpoint, err := NewEndpoint(LocalHub)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = endpoint.CloseAndAwaitTermination(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("CloseAndAwaitTermination: %v", err)
	}
	// The close goes on in the background; a Close at the same time must not race with it (the race detector checks)
	// or release the handle twice. Its error, if the handle is already released, does not matter here.
	_ = endpoint.Close()
}

func TestSubscriptionIsClosed(t *testing.T) {
	endpoint, err := NewEndpoint(LocalHub)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()
	feed, err := endpoint.GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatal(err)
	}
	if subscription.IsClosed() {
		t.Fatal("a new subscription is closed")
	}
	if err := subscription.Close(); err != nil {
		t.Fatal(err)
	}
	if !subscription.IsClosed() {
		t.Fatal("the subscription is not closed after Close")
	}
}
