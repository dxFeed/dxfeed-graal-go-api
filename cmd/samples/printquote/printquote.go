package main

import (
	"dxfeed-graal-go-api/pkg/api"
	"dxfeed-graal-go-api/pkg/events"
	"math"
	"os"
	"time"
)

func main() {
	symbol := os.Args[1]
	endpoint := api.CreateEndpoint(api.Feed)
	defer endpoint.Close()

	subscription := endpoint.GetFeed().CreateSubscription(events.Quote)
	subscription.AddSymbol(symbol)

	time.Sleep(time.Duration(math.MaxInt64))
}
