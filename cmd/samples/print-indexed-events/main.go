// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command print-indexed-events prints the orders of a symbol from one order source (the aggregated asks) and the
// connection state until Ctrl+C.
//
//	go run ./cmd/samples/print-indexed-events -address demo.dxfeed.com:7300 -symbol AAPL
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

type PrintState func(old common.ConnectionState, new common.ConnectionState)

func (pr PrintState) UpdateState(old common.ConnectionState, new common.ConnectionState) {
	pr(old, new)
}

func main() {
	// For token-based authorization, use the address format "demo.dxfeed.com:7300[login=entitle:token]".
	address := flag.String("address", "demo.dxfeed.com:7300", "the address to connect to")
	symbol := flag.String("symbol", "AAPL", "the symbol to print the orders of")
	flag.Parse()
	// Ctrl+C from now on ends the sample through ctx, so that the deferred calls close the endpoint.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func(endpoint *api.DXEndpoint) {
		_ = endpoint.Close()
	}(endpoint)
	endpoint.AddListener(PrintState(func(old common.ConnectionState, new common.ConnectionState) {
		fmt.Printf("Connection state changed from %s to %s\n", old, new)
	}))

	err = endpoint.Connect(*address)
	if err != nil {
		panic(err)
	}

	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	subscription, err := feed.CreateSubscription(eventcodes.Order)
	if err != nil {
		panic(err)
	}
	defer func() { _ = subscription.Close() }()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			if o, ok := event.(*order.Order); ok {
				fmt.Printf("%s\n", o.String())
			}
		}
	}))
	if err != nil {
		panic(err)
	}
	indexedSymbol := osub.NewIndexedEventSubscriptionSymbol(*symbol, order.AgregateAsk())
	err = subscription.AddSymbol(indexedSymbol)
	if err != nil {
		panic(err)
	}

	// Print the orders until Ctrl+C; returning from main closes the subscription and the endpoint.
	<-ctx.Done()
}
