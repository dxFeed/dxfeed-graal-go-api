// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command file-parser reads the events of some types and symbols from a tape file and prints them.
//
//	go run ./cmd/samples/file-parser -file WriteTapeFile.out.txt -types Quote -symbols TEST1,TEST2
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// eventListener adapts a function to common.EventListener.
type eventListener func(events []interface{})

func (f eventListener) Update(events []interface{}) { f(events) }

func main() {
	file := flag.String("file", "", "the tape file to read")
	types := flag.String("types", "", "comma-separated event types: "+typeNames())
	symbols := flag.String("symbols", "", `comma-separated symbols, for example "IBM,AAPL,MSFT"`)
	flag.Parse()
	if *file == "" || *types == "" || *symbols == "" {
		flag.Usage()
		os.Exit(2)
	}
	var codes []eventcodes.EventCode
	for _, name := range strings.Split(*types, ",") {
		code, ok := eventcodes.ByName(strings.TrimSpace(name))
		if !ok || code == eventcodes.OrderBase {
			fmt.Fprintf(os.Stderr, "unknown event type %q\n", name)
			os.Exit(2)
		}
		codes = append(codes, code)
	}

	// A stream feed endpoint neither conflates nor drops events.
	endpoint, err := api.NewEndpoint(api.StreamFeed)
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()
	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}
	subscription, err := feed.CreateSubscription(codes...)
	if err != nil {
		panic(err)
	}
	defer func() { _ = subscription.Close() }()
	var mu sync.Mutex
	count := 0
	err = subscription.AddListener(eventListener(func(events []interface{}) {
		mu.Lock()
		defer mu.Unlock()
		for _, event := range events {
			count++
			fmt.Printf("%d: %v\n", count, event)
		}
	}))
	if err != nil {
		panic(err)
	}
	var names []string
	for _, symbol := range strings.Split(*symbols, ",") {
		if symbol = strings.TrimSpace(symbol); symbol != "" {
			names = append(names, symbol)
		}
	}
	if err := subscription.AddSymbols(osub.Strings(names...)...); err != nil {
		panic(err)
	}

	if err := endpoint.Connect("file:" + *file + "[speed=max]"); err != nil {
		panic(err)
	}
	// Wait until the file is read, then close the endpoint, which waits until the events are processed.
	ctx := context.Background()
	if err := endpoint.AwaitNotConnected(ctx); err != nil {
		panic(err)
	}
	if err := endpoint.CloseAndAwaitTermination(ctx); err != nil {
		panic(err)
	}
}

// typeNames returns the names of the event types, sorted and separated by commas. OrderBase is left out: it has no
// events of its own.
func typeNames() string {
	var names []string
	for _, code := range eventcodes.All() {
		if code != eventcodes.OrderBase {
			names = append(names, code.String())
		}
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}
