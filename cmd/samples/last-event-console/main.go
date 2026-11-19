// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command last-event-console subscribes to the last events of a large set of symbols, which the process keeps in
// memory, and prints the last Quote, Trade, Summary and Profile of each symbol typed in the console. An empty line,
// "exit", "quit" or "q" ends it.
//
//	go run ./cmd/samples/last-event-console -address demo.dxfeed.com:7300
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func main() {
	address := flag.String("address", "demo.dxfeed.com:7300", "the address to connect to")
	// The records of the permanent subscription: the composite (NBBO) events. The records of the regional events add
	// "&" and the exchange code; a regular expression may list them, for example
	// "Quote,Trade,Summary,Profile,Quote&[A-Z]".
	records := flag.String("records", "Quote,Trade,Summary,Profile", "comma-separated records to keep in memory")
	// The symbols are usually an instrument profile file, given by a URL that contains ".ipf"; this one has all the
	// symbols of the demo feed.
	symbols := flag.String("symbols", "https://dxfeed.s3.amazonaws.com/masterdata/ipf/demo/mux-demo.ipf.zip",
		"comma-separated symbols, or the URL of an instrument profile file")
	flag.Parse()

	// The property "dxfeed.qd.subscribe.ticker" subscribes the endpoint permanently to the records of the symbols:
	// "<records> <symbols>". All the stocks and options of all the exchanges take gigabytes of memory and much
	// network and CPU to keep up to date: subscribe permanently only to what is often needed. The promises below
	// request an event from the data provider when it is not in memory.
	endpoint, err := api.NewEndpointBuilder().
		WithProperty("dxfeed.qd.subscribe.ticker", *records+" "+*symbols).
		Build()
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

	fmt.Println("Type symbols to get their quote, trade, summary, and profile event snapshots")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// A symbol like "IBM" is the composite of the stock; "IBM&N" is the information of one exchange.
		symbol := strings.TrimSpace(scanner.Text())
		if symbol == "" || symbol == "exit" || symbol == "quit" || symbol == "q" {
			return
		}
		types := []eventcodes.EventCode{eventcodes.Quote, eventcodes.Trade, eventcodes.Summary}
		// Profile events are composite only: a regional symbol ("&" and an exchange code at its end) never gets one, so
		// do not even wait for it.
		if !isRegional(symbol) {
			types = append(types, eventcodes.Profile)
		}
		// A promise completes at once when the event is in memory; otherwise the feed requests it. Wait at most 1
		// second for all of them together: a promise that is not done by then prints "<null>".
		var promises []*api.Promise[interface{}]
		for _, eventType := range types {
			promise, err := feed.GetLastEventPromise(eventType, osub.String(symbol))
			if err != nil {
				panic(err)
			}
			promises = append(promises, promise)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		timedOut := false
		for _, promise := range promises {
			select {
			case <-promise.Done():
			case <-ctx.Done():
			}
			// Both cases may be ready at the deadline: the promise decides.
			select {
			case <-promise.Done():
			default:
				timedOut = true
			}
		}
		cancel()
		if timedOut {
			fmt.Println("Request timed out")
		}
		for _, promise := range promises {
			if event, err := promise.Result(); err == nil && event != nil {
				fmt.Println(event)
			} else {
				fmt.Println("<null>")
			}
			_ = promise.Cancel()
		}
	}
}

// isRegional tells whether a symbol is regional: "&" and an exchange code end it, before the attributes "{...}" if
// any, as in "IBM&N".
func isRegional(symbol string) bool {
	if i := strings.IndexByte(symbol, '{'); i >= 0 && strings.HasSuffix(symbol, "}") {
		symbol = symbol[:i]
	}
	return len(symbol) >= 2 && symbol[len(symbol)-2] == '&'
}
