// Command multiple-market-depth keeps the order books of several symbols, one model.MarketDepthModel per symbol, and
// prints each book when it changes, at most once per aggregation period. [ENTER] prints the book of the first symbol;
// the end of the input (Ctrl+D, Ctrl+Z on Windows) exits.
//
//	go run ./cmd/samples/multiple-market-depth -address demo.dxfeed.com:7300 -symbols AAPL,TSLA -source NTV
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

// book is the last state of the order book of a symbol, as the listener of its model got it.
type book struct {
	buy, sell []*order.Order
}

func main() {
	// For token-based authorization, use the address format "demo.dxfeed.com:7300[login=entitle:<token>]".
	address := flag.String("address", "demo.dxfeed.com:7300", "the address of the feed")
	symbols := flag.String("symbols", "AAPL,TSLA", "the symbols, separated by commas")
	sourceName := flag.String("source", "NTV", "the order source")
	depth := flag.Int("depth", 10, "the number of the orders of each side; 0 means all")
	period := flag.Duration("period", 10*time.Second, "the aggregation period of the incremental changes")
	flag.Parse()
	// The symbols without blanks and repetitions.
	var symbolList []string
	seen := map[string]bool{}
	for _, symbol := range strings.Split(*symbols, ",") {
		if symbol = strings.TrimSpace(symbol); symbol != "" && !seen[symbol] {
			seen[symbol] = true
			symbolList = append(symbolList, symbol)
		}
	}
	if len(symbolList) == 0 {
		fmt.Fprintln(os.Stderr, "no symbols")
		os.Exit(2)
	}
	if err := run(*address, symbolList, *sourceName, *depth, *period); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address string, symbols []string, sourceName string, depth int, period time.Duration) error {
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

	// mu guards books and the output, so that the books of different symbols are not printed in between.
	var mu sync.Mutex
	books := map[string]book{}
	for _, symbol := range symbols {
		symbol := symbol
		m, err := model.NewMarketDepthModelBuilder[*order.Order]().
			WithFeed(feed).
			WithSymbol(symbol).
			WithSources(source).
			WithDepthLimit(depth).
			WithAggregationPeriod(period).
			WithListener(func(buy, sell []*order.Order) {
				mu.Lock()
				defer mu.Unlock()
				books[symbol] = book{buy, sell}
				printBook(symbol, sourceName, buy, sell)
			}).
			Build()
		if err != nil {
			return err
		}
		defer func() { _ = m.Close() }()
	}

	if err := endpoint.Connect(address); err != nil {
		return err
	}
	first := symbols[0]
	fmt.Printf("Press [ENTER] to print the order book of %s.\n", first)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		mu.Lock()
		if b, ok := books[first]; ok {
			fmt.Println("=============================Print manually=============================")
			printBook(first, sourceName, b.buy, b.sell)
			fmt.Println("========================================================================")
		}
		mu.Unlock()
	}
	return nil
}

func printBook(symbol string, source string, buy, sell []*order.Order) {
	fmt.Printf("%s#%s:\n", symbol, source)
	for i := 0; i < len(buy) || i < len(sell); i++ {
		fmt.Printf("%s \t %s\n", orderRow("Buy ", buy, i), orderRow("Sell ", sell, i))
	}
	fmt.Println()
}

func orderRow(side string, orders []*order.Order, i int) string {
	if i >= len(orders) {
		return side + " [None]"
	}
	o := orders[i]
	return fmt.Sprintf("%s [Size: %8.2f, Price: %8.4f]", side, o.Size(), o.Price())
}
