// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command feed-ipf-connect reads instrument profiles from an IPF file or service, subscribes to the quotes of their
// symbols and prints them until Ctrl+C.
//
//	go run ./cmd/samples/feed-ipf-connect -ipf "https://demo:demo@tools.dxfeed.com/ipf?TYPE=STOCK&compression=zip"
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
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
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
	ipfFile := flag.String("ipf", "https://demo:demo@tools.dxfeed.com/ipf?TYPE=STOCK&compression=zip",
		"the IPF file or service URL to read the instrument profiles from; user:password@ in a URL are its credentials")
	// For token-based authorization, use the address format
	// "dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws[login=dxlink:token]".
	address := flag.String("address", "dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws",
		"the address to connect to; dxLink and QD addresses both work")
	flag.Parse()
	// Ctrl+C from now on ends the sample through ctx, so that the deferred calls close the endpoint.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// The experimental property must be enabled.
	if err := api.SetSystemProperty("dxfeed.experimental.dxlink.enable", "true"); err != nil {
		panic(err)
	}
	// Set scheme for dxLink.
	if err := api.SetSystemProperty("scheme", "ext:opt:sysprops,resource:dxlink.xml"); err != nil {
		panic(err)
	}

	fmt.Printf("Reading instruments from %s\n", *ipfFile)
	ipfReader, err := ipf.NewInstrumentProfileReader()
	if err != nil {
		panic(err)
	}
	defer func(reader *ipf.InstrumentProfileReader) {
		_ = reader.Close()
	}(ipfReader)

	instrumentProfiles, err := ipfReader.ReadFromFile(*ipfFile)
	if err != nil {
		panic(err)
	}
	var symbols []osub.Symbol
	for _, instrumentProfile := range instrumentProfiles {
		symbols = append(symbols, osub.String(*instrumentProfile.Symbol()))
		fmt.Printf("%s\n", instrumentProfile.String())
	}
	completed, _ := ipfReader.WasComplete()
	fmt.Printf("Was Completed %t\n", completed)
	lastModified, _ := ipfReader.LastModified()
	fmt.Printf("LastModified %d\n", lastModified)

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

	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		panic(err)
	}
	defer func() { _ = subscription.Close() }()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			if q, ok := event.(*quote.Quote); ok {
				fmt.Printf("%s\n", q.String())
			}
		}
	}))
	if err != nil {
		panic(err)
	}

	err = subscription.AddSymbols(symbols...)
	if err != nil {
		panic(err)
	}

	// Print the quotes until Ctrl+C; returning from main closes the subscription, the endpoint and the reader.
	<-ctx.Done()
}
