// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
)

// benchBatch is the number of events in a published batch, as in NetTest with --batch 100.
const benchBatch = 100

// benchTrades returns a batch of trades, one for each of benchBatch symbols, and the symbols.
func benchTrades() ([]interface{}, []osub.Symbol) {
	batch := make([]interface{}, benchBatch)
	symbols := make([]osub.Symbol, benchBatch)
	for i := range batch {
		symbol := "SYM" + strconv.Itoa(i)
		event := trade.NewTrade(symbol)
		event.SetPrice(100 + float64(i))
		event.SetSize(10)
		batch[i] = event
		symbols[i] = osub.String(symbol)
	}
	return batch, symbols
}

// BenchmarkPublishTrade publishes batches of trades without subscribers: the cost of the Go conversion and of the
// native call, per event.
func BenchmarkPublishTrade(b *testing.B) {
	_, publisher := newLocalHub(b)
	batch, _ := benchTrades()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := publisher.Publish(batch); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*benchBatch), "ns/event")
}

// countingListener counts the received events and signals when a batch is complete.
type countingListener struct {
	received atomic.Int64
	batch    chan struct{}
}

func (l *countingListener) Update(events []interface{}) {
	if l.received.Add(int64(len(events)))%benchBatch == 0 {
		l.batch <- struct{}{}
	}
}

// BenchmarkReceiveTrade publishes batches of trades through a local hub to a subscription and waits for each batch:
// the cost of both conversions and of the delivery, per event.
func BenchmarkReceiveTrade(b *testing.B) {
	feed, publisher := newLocalHub(b)
	batch, symbols := benchTrades()
	subscription, err := feed.CreateSubscription(eventcodes.Trade)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { _ = subscription.Close() })
	listener := &countingListener{batch: make(chan struct{}, 1)}
	if err := subscription.AddListener(listener); err != nil {
		b.Fatal(err)
	}
	if err := subscription.AddSymbols(symbols...); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// The ticker delivers changed records only.
		for _, event := range batch {
			event.(*trade.Trade).SetSize(float64(i))
		}
		if err := publisher.Publish(batch); err != nil {
			b.Fatal(err)
		}
		select {
		case <-listener.batch:
		case <-time.After(5 * time.Second):
			b.Fatalf("batch %d: received %d events, want %d", i, listener.received.Load(), (i+1)*benchBatch)
		}
	}
	b.StopTimer()
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*benchBatch), "ns/event")
}

// BenchmarkPublishOneTrade publishes one trade at a time: the fixed cost of a batch.
func BenchmarkPublishOneTrade(b *testing.B) {
	_, publisher := newLocalHub(b)
	batch, _ := benchTrades()
	batch = batch[:1]
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := publisher.Publish(batch); err != nil {
			b.Fatal(err)
		}
	}
}
