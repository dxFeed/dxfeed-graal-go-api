package main

import (
	"dxfeed-graal-go-api/pkg/api"
	"dxfeed-graal-go-api/pkg/native"
	"math"
	"time"
)

func main() {
	api.SetSystemProperty("monitoring.stat", "1s")
	e := native.NewDXEndpointHandle(api.Feed)
	e.Connect("demo.dxfeed.com:7300")
	e.GetFeed().CreateSubscription(0).AddSymbol("ETH/USD:GDAX")
	time.Sleep(time.Duration(math.MaxInt64))
}
