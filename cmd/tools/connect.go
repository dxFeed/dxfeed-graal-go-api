package main

import (
	"dxfeed-graal-go-api/pkg/api"
)

type Connect struct{}

func (c Connect) Run(args []string) {
	address := args[2]
	endpoint := api.CreateEndpoint(api.Feed)
	endpoint.Connect(address)
	subscription := endpoint.GetFeed().CreateSubscription()
	subscription.AddSymbol("AAPL")
}
