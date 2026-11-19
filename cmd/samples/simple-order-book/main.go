// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command simple-order-book builds the order book of a symbol and an order source from the transactions of
// model.IndexedTxModel and prints its top levels with the scope of each order every second.
//
//	go run ./cmd/samples/simple-order-book -address demo.dxfeed.com:7300 -symbol AAPL -source NTV -scope aggregate
//
// With -duration it runs for the given time and exits with code 0 if the book had orders, 1 otherwise.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/model"
)

func main() {
	// For token-based authorization, use the address format "demo.dxfeed.com:7300[login=entitle:<token>]".
	address := flag.String("address", "demo.dxfeed.com:7300", "the address of the feed")
	symbol := flag.String("symbol", "AAPL", "the symbol")
	sourceName := flag.String("source", "NTV", "the order source")
	scopeName := flag.String("scope", "", `shows only the orders of the scope: "composite", "regional", "aggregate" `+
		`or "order"; empty means all`)
	depth := flag.Int("depth", 10, "the number of the levels to print")
	duration := flag.Duration("duration", 0, "runs for the given time and exits; 0 runs until interrupted")
	flag.Parse()
	hadOrders, err := run(*address, *symbol, *sourceName, *scopeName, *depth, *duration)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *duration > 0 && !hadOrders {
		os.Exit(1)
	}
}

func run(address, symbol, sourceName, scopeName string, depth int, duration time.Duration) (bool, error) {
	source, err := order.ValueOfName(sourceName)
	if err != nil {
		return false, err
	}
	scope, err := findScope(scopeName)
	if err != nil {
		return false, err
	}
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		return false, err
	}
	defer func() { _ = endpoint.Close() }()
	feed, err := endpoint.GetFeed()
	if err != nil {
		return false, err
	}

	book := newOrderBook(scope)
	m, err := model.NewIndexedTxModelBuilder[*order.Order]().
		WithFeed(feed).
		WithSymbol(symbol).
		WithSources(source).
		WithListener(func(_ events.IndexedEventSourceInterface, orders []*order.Order, isSnapshot bool) {
			book.update(orders, isSnapshot)
		}).
		Build()
	if err != nil {
		return false, err
	}
	defer func() { _ = m.Close() }()
	if err := endpoint.Connect(address); err != nil {
		return false, err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	var timeout <-chan time.Time
	if duration > 0 {
		timeout = time.After(duration)
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	hadOrders := false
	for {
		select {
		case <-ticker.C:
			buy, sell := book.orders()
			hadOrders = hadOrders || len(buy) > 0 || len(sell) > 0
			printBook(buy, sell, depth)
		case <-timeout:
			return hadOrders, nil
		case <-interrupt:
			return hadOrders, nil
		}
	}
}

// orderBook keeps the orders of one source by their index.
type orderBook struct {
	scope *order.Scope // shows only the orders of this scope; nil shows all of them

	mu            sync.Mutex
	ordersByIndex map[int64]*order.Order
}

func newOrderBook(scope *order.Scope) *orderBook {
	return &orderBook{scope: scope, ordersByIndex: make(map[int64]*order.Order)}
}

func (b *orderBook) update(orders []*order.Order, isSnapshot bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if isSnapshot { // a snapshot replaces the book
		clear(b.ordersByIndex)
	}
	for _, o := range orders {
		// A removed order and an order without a size leave the book.
		if events.IsRemove(o) || !o.HasSize() || (b.scope != nil && o.Scope() != *b.scope) {
			delete(b.ordersByIndex, o.Index())
		} else {
			b.ordersByIndex[o.Index()] = o
		}
	}
}

// orders returns the buy orders from the highest price and the sell orders from the lowest one.
func (b *orderBook) orders() (buy, sell []*order.Order) {
	b.mu.Lock()
	for _, o := range b.ordersByIndex {
		if o.Side() == side.Buy {
			buy = append(buy, o)
		} else {
			sell = append(sell, o)
		}
	}
	b.mu.Unlock()
	sortOrders(buy, func(p1, p2 float64) bool { return p1 > p2 })
	sortOrders(sell, func(p1, p2 float64) bool { return p1 < p2 })
	return buy, sell
}

// sortOrders sorts the orders by price, and by index if the prices are equal.
func sortOrders(orders []*order.Order, better func(p1, p2 float64) bool) {
	sort.Slice(orders, func(i, j int) bool {
		if orders[i].Price() != orders[j].Price() {
			return better(orders[i].Price(), orders[j].Price())
		}
		return orders[i].Index() < orders[j].Index()
	})
}

func printBook(buy, sell []*order.Order, depth int) {
	for i := 0; i < depth && (i < len(buy) || i < len(sell)); i++ {
		fmt.Printf("%s\t%s\n", orderInfo("Buy", buy, i), orderInfo("Sell", sell, i))
	}
	if len(buy) > 0 || len(sell) > 0 {
		fmt.Println()
	}
}

func orderInfo(name string, orders []*order.Order, i int) string {
	if i >= len(orders) {
		return name + " [None]"
	}
	o := orders[i]
	source := ""
	if s, err := o.OrderSource(); err == nil && s.Name() != nil {
		source = *s.Name()
	}
	return fmt.Sprintf("%s [Source: %s, Scope: %s, Size: %8.4f, Price: %8.2f]", name, source, o.Scope(), o.Size(),
		o.Price())
}

// findScope finds a scope by its name, ignoring the case; an empty name means no scope.
func findScope(name string) (*order.Scope, error) {
	if name == "" {
		return nil, nil
	}
	for _, scope := range []order.Scope{order.ScopeComposite, order.ScopeRegional, order.ScopeAggregate,
		order.ScopeOrder} {
		if strings.EqualFold(scope.String(), name) {
			return &scope, nil
		}
	}
	return nil, fmt.Errorf("unknown scope: %s", name)
}
