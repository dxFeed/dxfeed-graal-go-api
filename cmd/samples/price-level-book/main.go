// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command price-level-book prints the price levels (market by price) of a symbol, aggregated from the orders (market by
// order) of a source by model.MarketDepthModel: the size of all the orders at each price, up to a depth limit, at
// most once per aggregation period. Use a market by order source, such as NTV or GLBX.
//
//	go run ./cmd/samples/price-level-book -address demo.dxfeed.com:7300 -symbol MSFT -source NTV -depth 10
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
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/model"
)

// priceLevel is the total size of the orders of one side at one price.
type priceLevel struct {
	price, size float64
}

// aggregate merges the orders of one side, sorted from the best price, into at most depth price levels.
func aggregate(orders []*order.Order, depth int) []priceLevel {
	var levels []priceLevel
	for _, o := range orders {
		if n := len(levels); n > 0 && levels[n-1].price == o.Price() {
			levels[n-1].size += o.Size()
			continue
		}
		if len(levels) == depth {
			break
		}
		levels = append(levels, priceLevel{o.Price(), o.Size()})
	}
	return levels
}

func main() {
	// For token-based authorization, use the address format "demo.dxfeed.com:7300[login=entitle:<token>]".
	address := flag.String("address", "demo.dxfeed.com:7300", "the address of the feed")
	symbol := flag.String("symbol", "MSFT", "the symbol")
	sourceName := flag.String("source", "NTV", "the market by order source")
	depth := flag.Int("depth", 10, "the number of the price levels of each side")
	period := flag.Duration("period", 10*time.Second, "the aggregation period of the incremental changes")
	flag.Parse()
	if *depth <= 0 {
		fmt.Fprintln(os.Stderr, "the depth must be positive")
		os.Exit(2)
	}
	if err := run(*address, *symbol, *sourceName, *depth, *period); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address, symbol, sourceName string, depth int, period time.Duration) error {
	source, err := order.ValueOfName(sourceName)
	if err != nil {
		return err
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

	// The model keeps all the orders (no depth limit), so that the price levels have the full size at each price.
	var printMu sync.Mutex
	m, err := model.NewMarketDepthModelBuilder[*order.Order]().
		WithFeed(feed).
		WithSymbol(symbol).
		WithSources(source).
		WithAggregationPeriod(period).
		WithListener(func(buy, sell []*order.Order) {
			printMu.Lock()
			defer printMu.Unlock()
			printLevels(aggregate(buy, depth), aggregate(sell, depth))
		}).
		Build()
	if err != nil {
		return err
	}
	defer func() { _ = m.Close() }()

	if err := endpoint.Connect(address); err != nil {
		return err
	}
	fmt.Printf("Subscribed to the price levels of %s. Press [ENTER] to exit.\n", symbol)
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	return nil
}

func printLevels(buy, sell []priceLevel) {
	if len(buy) == 0 && len(sell) == 0 {
		return
	}
	fmt.Println(strings.Repeat("=", 66))
	fmt.Printf("%-31s || %-31s\n", center("BID", 31), center("ASK", 31))
	fmt.Printf("%-15s|%-15s || %-15s|%-15s\n", center("Price", 15), center("Size", 15), center("Price", 15),
		center("Size", 15))
	fmt.Println(strings.Repeat("-", 66))
	for i := 0; i < len(buy) || i < len(sell); i++ {
		fmt.Printf("%s || %s\n", levelRow(buy, i), levelRow(sell, i))
	}
	fmt.Println(strings.Repeat("=", 66))
}

func levelRow(levels []priceLevel, i int) string {
	if i >= len(levels) {
		return fmt.Sprintf("%14s | %-14s", "", "")
	}
	return fmt.Sprintf("%14.4f | %-14.2f", levels[i].price, levels[i].size)
}

// center centers s in a field of the given width.
func center(s string, width int) string {
	left := (width - len(s)) / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", width-len(s)-left)
}
