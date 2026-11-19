// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command publish-profiles starts a server that publishes a Profile event for every requested symbol that ends with
// ":TEST", until Ctrl+C. The request-profile sample requests such profiles.
//
//	go run ./cmd/samples/publish-profiles -address :7700
package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
)

// profilePublisher publishes the profiles of the ":TEST" symbols that the clients subscribe to.
type profilePublisher struct {
	publisher *api.DXPublisher
}

func (p profilePublisher) SymbolsAdded(symbols []osub.Symbol) {
	var events []interface{}
	for _, symbol := range symbols {
		if s, ok := symbol.(osub.String); ok && len(s) > len(":TEST") && strings.HasSuffix(string(s), ":TEST") {
			event := profile.NewProfile(string(s))
			description := "Test symbol"
			event.SetDescription(&description)
			events = append(events, event)
		}
	}
	if len(events) > 0 {
		_ = p.publisher.Publish(events)
	}
}

func (p profilePublisher) SymbolsRemoved([]osub.Symbol) {}

func (p profilePublisher) SubscriptionClosed() {}

func main() {
	address := flag.String("address", ":7700", "the address to listen on")
	flag.Parse()

	endpoint, err := api.NewEndpoint(api.Publisher)
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()
	if err := endpoint.Connect(*address); err != nil {
		panic(err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		panic(err)
	}
	subscription, err := publisher.Subscription(eventcodes.Profile)
	if err != nil {
		panic(err)
	}
	if err := subscription.AddChangeListener(profilePublisher{publisher}); err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
