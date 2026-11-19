// Command print-quote-events prints the quotes of a symbol from a dxLink address until Ctrl+C.
//
//	go run ./cmd/samples/print-quote-events -address dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws -symbol AAPL
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

type PrintState func(old common.ConnectionState, new common.ConnectionState)

func (pr PrintState) UpdateState(old common.ConnectionState, new common.ConnectionState) {
	pr(old, new)
}

func main() {
	// For token-based authorization, use the address format
	// "dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws[login=dxlink:token]".
	address := flag.String("address", "dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws",
		"the address to connect to; dxLink and QD addresses both work")
	symbol := flag.String("symbol", "AAPL", "the symbol to print the quotes of")
	flag.Parse()
	// Ctrl+C from now on ends the sample through ctx, so that the deferred calls close the endpoint.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// The experimental property must be enabled.
	if err := api.SetSystemProperty("dxfeed.experimental.dxlink.enable", "true"); err != nil {
		panic(err)
	}
	// Set scheme for dxLink.
	if err := api.SetSystemProperty("scheme", "ext:opt:sysprops,resource:dxlink.xml"); err != nil {
		panic(err)
	}

	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func(endpoint *api.DXEndpoint) {
		_ = endpoint.Close()
	}(endpoint)
	endpoint.AddListener(PrintState(func(old common.ConnectionState, new common.ConnectionState) {
		fmt.Printf("Connection state changed from %s to %s\n", old, new)
	}))

	err = endpoint.Connect(*address)
	if err != nil {
		panic(err)
	}

	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		panic(err)
	}
	defer func() { _ = subscription.Close() }()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			if q, ok := event.(*quote.Quote); ok {
				fmt.Printf("%s\n", q.String())
			}
		}
	}))
	if err != nil {
		panic(err)
	}

	err = subscription.AddSymbol(osub.String(*symbol))
	if err != nil {
		panic(err)
	}

	// Print the quotes until Ctrl+C; returning from main closes the subscription and the endpoint.
	<-ctx.Done()
}
