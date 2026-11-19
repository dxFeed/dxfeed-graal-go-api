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
		code, ok := eventTypes[strings.TrimSpace(name)]
		if !ok {
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

// typeNames returns the names of the event types, sorted and separated by commas.
func typeNames() string {
	names := make([]string, 0, len(eventTypes))
	for name := range eventTypes {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}
