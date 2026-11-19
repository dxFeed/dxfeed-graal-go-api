// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command print-last-events requests the last Quote of several symbols and the time and sales of the last minute of
// one of them with promises, waits for the results and prints them.
//
//	go run ./cmd/samples/print-last-events -address demo.dxfeed.com:7300 -symbols AAPL,IBM,MSFT
package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func main() {
	address := flag.String("address", "demo.dxfeed.com:7300", "the address to connect to")
	symbols := flag.String("symbols", "AAPL,IBM,MSFT", "comma-separated symbols")
	timeout := flag.Duration("timeout", 10*time.Second, "how long to wait for the results")
	flag.Parse()

	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()
	if err := endpoint.Connect(*address); err != nil {
		panic(err)
	}
	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// One promise per symbol; Await cancels a promise that is not done when the context ends.
	names := strings.Split(*symbols, ",")
	quotes, err := feed.GetLastEventsPromises(eventcodes.Quote, osub.Strings(names...)...)
	if err != nil {
		panic(err)
	}
	for i, promise := range quotes {
		quote, err := promise.Await(ctx)
		switch {
		case err != nil:
			fmt.Printf("%s: %v\n", names[i], err)
		case quote == nil:
			fmt.Printf("%s: no quote\n", names[i])
		default:
			fmt.Println(quote)
		}
	}

	from := time.Now().Add(-time.Minute).UnixMilli()
	sales, err := feed.GetTimeSeriesPromise(eventcodes.TimeAndSale, osub.String(names[0]), from, math.MaxInt64)
	if err != nil {
		panic(err)
	}
	events, err := sales.Await(ctx)
	if err != nil {
		fmt.Printf("%s time and sales: %v\n", names[0], err)
		return
	}
	fmt.Printf("%s: %d time and sales in the last minute\n", names[0], len(events))
	for _, event := range events {
		fmt.Println(event)
	}
}
