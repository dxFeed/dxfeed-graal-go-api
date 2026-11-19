// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command connect connects to an address, subscribes to events of some types and symbols, optionally from a time, and
// prints them until Ctrl+C.
//
//	go run ./cmd/samples/connect -address demo.dxfeed.com:7300 -types Quote,Trade -symbols MSFT,IBM
//	go run ./cmd/samples/connect -address demo.dxfeed.com:7300 -types TimeAndSale -symbols AAPL
//	go run ./cmd/samples/connect -address demo.dxfeed.com:7300 -types Candle -symbols "AAPL{=d}" -from 20230901Z
//
// The address is a remote host or a tape file ("file:<path>"). To pass an authorization token, add it to the address:
// "demo.dxfeed.com:7300[login=entitle:<token>]".
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/dxtime"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// eventTypes maps the names of the event types to their codes.
var eventTypes = map[string]eventcodes.EventCode{
	"Quote": eventcodes.Quote, "Profile": eventcodes.Profile, "Summary": eventcodes.Summary,
	"Greeks": eventcodes.Greeks, "Candle": eventcodes.Candle, "Underlying": eventcodes.Underlying,
	"TheoPrice": eventcodes.TheoPrice, "Trade": eventcodes.Trade, "TradeETH": eventcodes.TradeETH,
	"Configuration": eventcodes.Configuration, "Message": eventcodes.Message, "TimeAndSale": eventcodes.TimeAndSale,
	"Order": eventcodes.Order, "AnalyticOrder": eventcodes.AnalyticOrder, "OtcMarketsOrder": eventcodes.OtcMarketsOrder,
	"SpreadOrder": eventcodes.SpreadOrder, "Series": eventcodes.Series, "OptionSale": eventcodes.OptionSale,
	"TextMessage": eventcodes.TextMessage, "MarketMaker": eventcodes.MarketMaker,
	"TextConfiguration": eventcodes.TextConfiguration, "NuamOrder": eventcodes.NuamOrder,
	"NuamTimeAndSale": eventcodes.NuamTimeAndSale, "NuamTrade": eventcodes.NuamTrade,
	"OrderImbalance": eventcodes.OrderImbalance,
}

// eventListener adapts a function to common.EventListener.
type eventListener func(events []interface{})

func (f eventListener) Update(events []interface{}) { f(events) }

func main() {
	address := flag.String("address", "demo.dxfeed.com:7300", "the address: a remote host or a tape file")
	types := flag.String("types", "Quote", "comma-separated event types: "+typeNames())
	symbols := flag.String("symbols", "AAPL", `comma-separated symbols, for example "IBM,AAPL,MSFT" or "AAPL{=d}"`)
	from := flag.String("from", "", `the time to subscribe from, for example "20070101-123456", "2007-11-02Z" `+
		`or milliseconds since the epoch; empty to subscribe without history`)
	flag.Parse()

	var codes []eventcodes.EventCode
	for _, name := range strings.Split(*types, ",") {
		code, ok := eventTypes[strings.TrimSpace(name)]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown event type %q\n", name)
			os.Exit(2)
		}
		codes = append(codes, code)
	}
	var subscriptionSymbols []osub.Symbol
	for _, symbol := range strings.Split(*symbols, ",") {
		if symbol = strings.TrimSpace(symbol); symbol != "" {
			subscriptionSymbols = append(subscriptionSymbols, osub.String(symbol))
		}
	}
	if *from != "" {
		fromTime, err := dxtime.DefaultFormat().Parse(*from)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad time %q: %v\n", *from, err)
			os.Exit(2)
		}
		for i, symbol := range subscriptionSymbols {
			subscriptionSymbols[i] = osub.NewTimeSeriesSubscriptionSymbol(symbol, fromTime.UnixMilli())
		}
	}

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
	subscription, err := feed.CreateSubscription(codes...)
	if err != nil {
		panic(err)
	}
	// Close the subscription before the endpoint: Close waits for the listener calls in progress.
	defer func() { _ = subscription.Close() }()
	var mu sync.Mutex
	err = subscription.AddListener(eventListener(func(events []interface{}) {
		mu.Lock()
		defer mu.Unlock()
		for _, event := range events {
			fmt.Println(event)
		}
	}))
	if err != nil {
		panic(err)
	}
	if err := subscription.AddSymbols(subscriptionSymbols...); err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}

// typeNames returns the names of the event types, sorted and separated by commas.
func typeNames() string {
	names := make([]string, 0, len(eventTypes))
	for name := range eventTypes {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}
