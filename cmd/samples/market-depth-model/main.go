// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command market-depth-model prints the order book of a symbol with model.MarketDepthModel: the best buy and sell
// orders of the sources, up to the depth limit, at most once per aggregation period.
//
//	go run ./cmd/samples/market-depth-model -address demo.dxfeed.com:7300 -symbol AAPL -sources NTV -depth 10
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/model"
)

func main() {
	// For token-based authorization, use the address format "demo.dxfeed.com:7300[login=entitle:<token>]".
	address := flag.String("address", "demo.dxfeed.com:7300", "the address of the feed")
	symbol := flag.String("symbol", "AAPL", "the symbol")
	sources := flag.String("sources", "NTV", "the order sources, separated by commas; empty means all")
	depth := flag.Int("depth", 10, "the number of the orders of each side; 0 means all")
	period := flag.Duration("period", time.Second, "the aggregation period of the incremental changes")
	flag.Parse()
	if err := run(*address, *symbol, *sources, *depth, *period); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address, symbol, sourceNames string, depth int, period time.Duration) error {
	var sources []events.IndexedEventSourceInterface
	for _, name := range strings.Split(sourceNames, ",") {
		if name = strings.TrimSpace(name); name == "" {
			continue
		}
		source, err := order.ValueOfName(name)
		if err != nil {
			return err
		}
		sources = append(sources, source)
	}
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		return err
	}
	defer func() { _ = endpoint.Close() }()
	feed, err := endpoint.GetFeed()
	if err != nil {
		return err
	}

	var printMu sync.Mutex
	m, err := model.NewMarketDepthModelBuilder[*order.Order]().
		WithFeed(feed).
		WithSymbol(symbol).
		WithSources(sources...).
		WithDepthLimit(depth).
		WithAggregationPeriod(period).
		WithListener(func(buy, sell []*order.Order) {
			printMu.Lock()
			defer printMu.Unlock()
			printBook(symbol, buy, sell)
		}).
		Build()
	if err != nil {
		return err
	}
	defer func() { _ = m.Close() }()

	if err := endpoint.Connect(address); err != nil {
		return err
	}
	fmt.Printf("Subscribed to the order book of %s. Press [ENTER] to exit.\n", symbol)
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	return nil
}

func printBook(symbol string, buy, sell []*order.Order) {
	if len(buy) == 0 && len(sell) == 0 {
		return
	}
	line := strings.Repeat("=", 81)
	fmt.Println(line)
	fmt.Printf("%-39s || %-39s\n", symbol+" BID", "ASK")
	fmt.Printf("%2s|%-8s|%-6s|%10s|%9s || %2s|%-8s|%-6s|%10s|%9s\n",
		"EX", "Source", "MM", "Price", "Size", "EX", "Source", "MM", "Price", "Size")
	fmt.Println(strings.Repeat("-", 81))
	for i := 0; i < len(buy) || i < len(sell); i++ {
		fmt.Printf("%-39s || %-39s\n", orderRow(buy, i), orderRow(sell, i))
	}
}

func orderRow(orders []*order.Order, i int) string {
	if i >= len(orders) {
		return ""
	}
	o := orders[i]
	source := ""
	if s, err := o.OrderSource(); err == nil && s.Name() != nil {
		source = *s.Name()
	}
	marketMaker := ""
	if o.MarketMaker() != nil {
		marketMaker = *o.MarketMaker()
	}
	return fmt.Sprintf("%2s|%-8.8s|%-6.6s|%10.4f|%9.2f", exchange(o.ExchangeCode()), source, marketMaker, o.Price(),
		o.Size())
}

// exchange returns the exchange code as printable text.
func exchange(code rune) string {
	switch {
	case code == 0:
		return `\0`
	case code >= ' ' && code <= '~':
		return string(code)
	}
	return fmt.Sprintf(`\u%04x`, code)
}
