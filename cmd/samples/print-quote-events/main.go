package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"math"
	"time"
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

func main() {
	// The experimental property must be enabled.
	api.SetSystemProperty("dxfeed.experimental.dxlink.enable", "true")
	// Set scheme for dxLink.
	api.SetSystemProperty("scheme", "ext:opt:sysprops,resource:dxlink.xml")

	// For token-based authorization, use the following address format:
	// "dxlink:wss://demo.dxfeed.com/dxlink-ws[login=dxlink:token]"
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func(endpoint *api.DXEndpoint) {
		_ = endpoint.Close()
	}(endpoint)

	err = endpoint.Connect("dxlink:wss://demo.dxfeed.com/dxlink-ws")
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
	defer subscription.Close()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			switch v := event.(type) {
			case *quote.Quote:
				fmt.Printf("%s\n", v.String())
			}
		}
	}))

	err = subscription.AddSymbol("AAPL")
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(math.MaxInt64))
}
