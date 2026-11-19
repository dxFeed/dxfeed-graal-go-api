// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command ondemand-replay replays the quotes of Accenture (ACN), which dropped under $1 in the "Flash Crash" of
// May 6, 2010 from 14:47:48 to 14:48:02 EST, from the dxFeed on-demand history data replay service. The service needs
// the credentials of a dxFeed account with access to it; Ctrl+C stops the replay.
//
//	go run ./cmd/samples/ondemand-replay -address "(ondemand:demo.dxfeed.com:7680)" -user <user> -password <password>
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
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

const timeLayout = "2006-01-02 15:04:05.000 MST"

func main() {
	// The address of the on-demand service: the "(ondemand:...)" component makes the replay available.
	address := flag.String("address", "(ondemand:demo.dxfeed.com:7680)", "the address of the on-demand service")
	user := flag.String("user", "demo", "the user name of the on-demand service")
	password := flag.String("password", "demo", "the password of the on-demand service")
	symbol := flag.String("symbol", "ACN", "the symbol to replay the quotes of")
	flag.Parse()
	// Ctrl+C from now on ends the sample through ctx, so that the deferred calls close the endpoint.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	endpoint, err := api.NewEndpointWithProperties(api.OnDemandFeed, map[string]string{
		// The replayed events keep the time when they happened.
		"dxendpoint.eventTime": "true",
		"dxfeed.user":          *user,
		"dxfeed.password":      *password,
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()

	onDemand, err := endpoint.OnDemandService()
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
	est := time.FixedZone("EST", -5*60*60)
	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			if q, ok := event.(*quote.Quote); ok {
				fmt.Printf("%s : %s bid %v / ask %v\n", time.UnixMilli(q.EventTime()).In(est).Format(timeLayout),
					*q.EventSymbol(), q.BidPrice(), q.AskPrice())
			}
		}
	}))
	if err != nil {
		panic(err)
	}
	if err = subscription.AddSymbol(osub.String(*symbol)); err != nil {
		panic(err)
	}

	if err = endpoint.Connect(*address); err != nil {
		panic(err)
	}
	from := time.Date(2010, 5, 6, 14, 47, 48, 0, est)
	to := time.Date(2010, 5, 6, 14, 48, 2, 0, est)
	if err = onDemand.Replay(from); err != nil {
		panic(err)
	}

	// Print the state every second until the replay reaches the end time or Ctrl+C; returning from main closes the
	// subscription and the endpoint.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		now, err := onDemand.Time()
		if err != nil {
			panic(err)
		}
		if !now.Before(to) {
			return
		}
		state, _ := endpoint.State()
		fmt.Printf("Current state is %s, on-demand time is %s\n", state, now.In(est).Format(timeLayout))
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
