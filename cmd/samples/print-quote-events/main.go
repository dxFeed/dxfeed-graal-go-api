package main

import (
	"dxfeed-graal-go-api/pkg/api"
	"dxfeed-graal-go-api/pkg/events"
	"dxfeed-graal-go-api/pkg/events/market"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	symbol := os.Args[1]
	// The experimental property must be enabled.
	api.SetSystemProperty("dxfeed.experimental.dxlink.enable", "true")
	// Set scheme for dxLink.
	api.SetSystemProperty("scheme", "ext:resource:dxlink.xml")

	endpoint := api.NewEndpoint(api.Feed).Connect("demo.dxfeed.com:7300")
	defer endpoint.Close()

	subscription := endpoint.GetFeed().CreateSubscription(events.Quote)
	subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			switch v := event.(type) {
			case market.Quote:
				fmt.Printf("%s\n", v.String())
			default:
				fmt.Printf("Unsupported event %T!\n", v)
			}
		}
	}))
	subscription.AddSymbol(symbol)
	time.Sleep(time.Duration(math.MaxInt64))
}

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []interface{}) {
	pr(events)
}
