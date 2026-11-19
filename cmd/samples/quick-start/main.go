// Command quick-start prints quotes of AAPL and IBM until Ctrl+C, from the dxFeed demo feed or the -address.
// README.md shows this file as it is; a test keeps the two identical.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// listener turns a function into an event listener (common.EventListener).
type listener func(events []any)

func (l listener) Update(events []any) { l(events) }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// For token-based authorization, use the address "demo.dxfeed.com:7300[login=entitle:token]".
	address := flag.String("address", "demo.dxfeed.com:7300", "the address to connect to")
	flag.Parse()
	// Ctrl+C from now on ends the sample through ctx, so that the deferred calls close the endpoint.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		return err
	}
	defer func() { _ = endpoint.Close() }()

	if err := endpoint.Connect(*address); err != nil {
		return err
	}

	feed, err := endpoint.GetFeed()
	if err != nil {
		return err
	}
	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		return err
	}
	defer func() { _ = subscription.Close() }()

	err = subscription.AddListener(listener(func(events []any) {
		for _, event := range events {
			if q, ok := event.(*quote.Quote); ok {
				fmt.Println(q)
			}
		}
	}))
	if err != nil {
		return err
	}
	if err := subscription.AddSymbols(osub.Strings("AAPL", "IBM")...); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
