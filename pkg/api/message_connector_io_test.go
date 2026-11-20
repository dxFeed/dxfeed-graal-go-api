// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// TestConnectorIOStatistics connects a feed to a local publisher, receives quotes and reads the I/O statistics of the
// connectors of both sides.
func TestConnectorIOStatistics(t *testing.T) {
	port := freePort()
	server, err := api.NewEndpoint(api.Publisher)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Close() }()
	if err := server.Connect(":" + port); err != nil {
		t.Fatal(err)
	}
	waitListening(port)
	publisher, err := server.GetPublisher()
	if err != nil {
		t.Fatal(err)
	}

	client, err := api.NewEndpoint(api.Feed)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()
	if err := client.Connect("127.0.0.1:" + port); err != nil {
		t.Fatal(err)
	}
	feed, err := client.GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = subscription.Close() }()
	received := make(eventChannel, 1000)
	if err := subscription.AddListener(received); err != nil {
		t.Fatal(err)
	}
	if err := subscription.AddSymbol(osub.String("IO-STATS")); err != nil {
		t.Fatal(err)
	}

	// Publish until the feed receives a quote: the subscription reaches the publisher after the connect.
	const want = 10
	got := 0
	for deadline := time.Now().Add(30 * time.Second); got < want; {
		if time.Now().After(deadline) {
			t.Fatalf("received %d quotes within 30 s; want %d", got, want)
		}
		q := quote.NewQuote("IO-STATS")
		q.SetBidPrice(float64(got + 1))
		if err := publisher.Publish([]interface{}{q}); err != nil {
			t.Fatal(err)
		}
		select {
		case <-received:
			got++
		case <-time.After(100 * time.Millisecond):
		}
	}

	stats, err := client.ConnectorIOStatistics()
	if err != nil {
		t.Fatalf("ConnectorIOStatistics of the feed: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("ConnectorIOStatistics of the feed = %+v; want one connector", stats)
	}
	s := stats[0]
	if s.Address != "127.0.0.1:"+port || s.ReadDataRecords < want || s.WrittenSubRecords < 1 || s.ReadBytes <= 0 ||
		s.WrittenBytes <= 0 || s.ActiveConnectionCount != 1 || s.ReadDataLagSumUs < 0 || s.ReadRTTSumUs < 0 {
		t.Errorf("ConnectorIOStatistics of the feed = %+v; want %d data records read, a subscription written and "+
			"one open connection", s, want)
	}
	connectors, err := client.Connectors()
	if err != nil || len(connectors) != 1 || connectors[0].Name != s.Name {
		t.Errorf("Connectors() = %+v, %v; want the connector %q", connectors, err, s.Name)
	}

	stats, err = server.ConnectorIOStatistics()
	if err != nil || len(stats) != 1 {
		t.Fatalf("ConnectorIOStatistics of the publisher = %+v, %v; want one connector", stats, err)
	}
	s = stats[0]
	if s.WrittenDataRecords < want || s.ReadSubRecords < 1 || s.ActiveConnectionCount != 1 {
		t.Errorf("ConnectorIOStatistics of the publisher = %+v; want %d data records written, a subscription read "+
			"and one open connection", s, want)
	}

	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ConnectorIOStatistics(); err == nil || errors.Is(err, api.ErrIOStatisticsUnsupported) {
		t.Errorf("ConnectorIOStatistics after Close = %v; want the error of a closed endpoint", err)
	}
}

func TestConnectorIOStatisticsOfLocalHub(t *testing.T) {
	endpoint, err := api.NewEndpoint(api.LocalHub)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()
	stats, err := endpoint.ConnectorIOStatistics()
	if err != nil || len(stats) != 0 {
		t.Errorf("ConnectorIOStatistics of a local hub = %+v, %v; want none", stats, err)
	}
}
