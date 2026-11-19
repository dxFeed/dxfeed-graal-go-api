// Command quote-and-trade shows several ways to get the quotes and trades of a symbol: a Quote listener that prints the
// mid price, a promise of the last quote, a listener of quotes and trades, and snapshots of the last trade every
// second until Ctrl+C.
//
//	go run ./cmd/samples/quote-and-trade -address demo.dxfeed.com:7300 -symbol AAPL
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
)

// eventListener adapts a function to common.EventListener.
type eventListener func(events []interface{})

func (f eventListener) Update(events []interface{}) { f(events) }

func main() {
	address := flag.String("address", "demo.dxfeed.com:7300", "the address to connect to")
	symbol := flag.String("symbol", "AAPL", "the symbol of the instrument, for example IBM, AAPL or SPX")
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

	quoteListener(feed, *symbol)
	quotePromise(feed, *symbol)
	quoteAndTradeListener(feed, *symbol)
	tradeSnapshots(feed, *symbol)
}

// quoteListener prints the mid price of each quote. The subscriptions of this sample live until the endpoint closes.
func quoteListener(feed *api.DXFeed, symbol string) {
	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		panic(err)
	}
	// Add the listener before the symbols, so that it gets the first events.
	err = subscription.AddListener(eventListener(func(events []interface{}) {
		for _, event := range events {
			if q, ok := event.(*quote.Quote); ok {
				fmt.Printf("Mid = %v\n", (q.BidPrice()+q.AskPrice())/2)
			}
		}
	}))
	if err != nil {
		panic(err)
	}
	if err := subscription.AddSymbol(osub.String(symbol)); err != nil {
		panic(err)
	}
}

// quotePromise waits up to 5 seconds for the last quote and prints its mid price.
func quotePromise(feed *api.DXFeed, symbol string) {
	promise, err := feed.GetLastEventPromise(eventcodes.Quote, osub.String(symbol))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	event, err := promise.Await(ctx)
	if err != nil {
		fmt.Printf("No quote for %s: %v\n", symbol, err)
	} else if q, ok := event.(*quote.Quote); ok && q != nil {
		fmt.Printf("Mid = %v\n", (q.BidPrice()+q.AskPrice())/2)
	} else {
		fmt.Printf("No quote for %s\n", symbol)
	}
}

// quoteAndTradeListener prints the quotes and the trades of the symbol.
func quoteAndTradeListener(feed *api.DXFeed, symbol string) {
	subscription, err := feed.CreateSubscription(eventcodes.Quote, eventcodes.Trade)
	if err != nil {
		panic(err)
	}
	err = subscription.AddListener(eventListener(func(events []interface{}) {
		for _, event := range events {
			switch e := event.(type) {
			case *quote.Quote:
				fmt.Println(e)
			case *trade.Trade:
				fmt.Println(e)
			}
		}
	}))
	if err != nil {
		panic(err)
	}
	if err := subscription.AddSymbol(osub.String(symbol)); err != nil {
		panic(err)
	}
}

// tradeSnapshots prints the last trade from the cache of the feed every second until Ctrl+C. GetLastEvent returns the
// last event only for a subscribed symbol: the subscription keeps it subscribed.
func tradeSnapshots(feed *api.DXFeed, symbol string) {
	subscription, err := feed.CreateSubscription(eventcodes.Trade)
	if err != nil {
		panic(err)
	}
	if err := subscription.AddSymbol(osub.String(symbol)); err != nil {
		panic(err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-interrupt:
			return
		case <-ticker.C:
			last, err := feed.GetLastEvent(trade.NewTrade(symbol))
			if err != nil {
				panic(err)
			}
			fmt.Printf("LAST: %v\n", last)
		}
	}
}
