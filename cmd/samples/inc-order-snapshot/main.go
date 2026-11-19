// Command inc-order-snapshot prints the orders of a symbol and an order source with model.IndexedTxModel: first the
// snapshot of the order book, then its incremental updates, as complete transactions.
//
//	go run ./cmd/samples/inc-order-snapshot -address demo.dxfeed.com:7300 -symbol AAPL -source NTV
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/model"
)

func main() {
	// For token-based authorization, use the address format "demo.dxfeed.com:7300[login=entitle:<token>]".
	address := flag.String("address", "demo.dxfeed.com:7300", "the address of the feed")
	symbol := flag.String("symbol", "AAPL", "the symbol")
	sourceName := flag.String("source", "NTV", "the order source")
	limit := flag.Int("limit", 10, "the number of orders to print per transaction; 0 prints all")
	flag.Parse()
	if err := run(*address, *symbol, *sourceName, *limit); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address, symbol, sourceName string, limit int) error {
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

	m, err := model.NewIndexedTxModelBuilder[*order.Order]().
		WithBatchProcessing(true).    // all the pending transactions in one call
		WithSnapshotProcessing(true). // a snapshot as its resulting state
		WithFeed(feed).
		WithSymbol(symbol).
		WithSources(source).
		WithListener(func(source events.IndexedEventSourceInterface, orders []*order.Order, isSnapshot bool) {
			printOrders(symbol, source, orders, isSnapshot, limit)
		}).
		Build()
	if err != nil {
		return err
	}
	defer func() { _ = m.Close() }()

	if err := endpoint.Connect(address); err != nil {
		return err
	}
	fmt.Printf("Subscribed to the orders of %s. Press [ENTER] to exit.\n", symbol)
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	return nil
}

func printOrders(symbol string, source events.IndexedEventSourceInterface, orders []*order.Order, isSnapshot bool,
	limit int) {
	kind := "update"
	if isSnapshot {
		kind = "full"
	}
	name := ""
	if source.Name() != nil {
		name = *source.Name()
	}
	fmt.Printf("Snapshot { Symbol: '%s', Source: '%s', Type: %s }\n", symbol, name, kind)
	for i, o := range orders {
		if limit > 0 && i == limit {
			fmt.Printf("   { ... %d records left ... }\n", len(orders)-i)
			break
		}
		if events.IsRemove(o) || !o.HasSize() {
			fmt.Printf("    REMOVAL { Index: %#x }\n", o.Index())
		} else {
			fmt.Printf("    %s\n", o.String())
		}
	}
}
