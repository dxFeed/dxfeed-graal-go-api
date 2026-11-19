// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Command auth connects to an address that takes an authorization token, renews the token every 10 seconds and
// prints the quotes of a symbol until Ctrl+C. The demo feed accepts any token; this sample makes random ones.
//
//	go run ./cmd/samples/auth -address demo.dxfeed.com:7300 -symbol AAPL
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

var (
	address = flag.String("address", "demo.dxfeed.com:7300", "the address to connect to, without the token")
	symbol  = flag.String("symbol", "AAPL", "the symbol to print the quotes of")
)

type listener func(events []any)

func (l listener) Update(events []any) { l(events) }

type stateListener func(old, new common.ConnectionState)

func (l stateListener) UpdateState(old, new common.ConnectionState) { l(old, new) }

// newToken stands for a request of a fresh token from an authorization service.
func newToken() (*auth.Token, error) {
	value := make([]byte, 4+rand.Intn(6))
	for i := range value {
		value[i] = byte('a' + rand.Intn(26))
	}
	return auth.NewBearer(string(value))
}

// connect connects the endpoint with a new token. Connecting a connected endpoint reconnects it; the subscriptions
// are subscribed again.
func connect(endpoint *api.DXEndpoint) error {
	token, err := newToken()
	if err != nil {
		return err
	}
	return endpoint.Connect(*address + "[login=entitle:" + token.Value() + "]")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		return err
	}
	defer func() { _ = endpoint.Close() }()
	endpoint.AddListener(stateListener(func(old, new common.ConnectionState) {
		fmt.Printf("Connection state changed: %s->%s\n", old, new)
	}))

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
	if err := subscription.AddSymbols(osub.String(*symbol)); err != nil {
		return err
	}

	if err := connect(endpoint); err != nil {
		return err
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := connect(endpoint); err != nil {
				return err
			}
		}
	}
}
