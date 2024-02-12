package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/Osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"math"
	"time"
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
	// For token-based authorization, use the following address format:
	// "demo.dxfeed.com:7300[login=entitle:token]"
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

	err = endpoint.Connect("demo.dxfeed.com:7300")
	if err != nil {
		panic(err)
	}

	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	subscription, err := feed.CreateSubscription(eventcodes.Order)
	if err != nil {
		panic(err)
	}
	defer subscription.Close()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			switch v := event.(type) {
			case *order.Order:
				fmt.Printf("%s\n", v.String())
			}
		}
	}))
	indexedSymbol := Osub.NewIndexedEventSubscriptionSymbol("AAPL", order.AgregateAsk())
	err = subscription.AddSymbol(indexedSymbol)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(math.MaxInt64))
}
