// Command print-candles subscribes to the hourly candles of AAPL for the last day and prints them as they arrive.
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
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
	// The candle symbol of hourly candles: the same as candle.NewCandleSymbol("AAPL{=h}").
	symbol, err := candle.NewCandleSymbolWithAttributes("AAPL", candle.NewCandlePeriod(1, candle.CandleTypeHour))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Candle symbol %s: period %s, price %s\n", symbol, symbol.Period(), symbol.Price())

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

	subscription, err := feed.CreateSubscription(eventcodes.Candle)
	if err != nil {
		panic(err)
	}
	defer func() { _ = subscription.Close() }()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			if c, ok := event.(*candle.Candle); ok {
				fmt.Println(c)
			}
		}
	}))
	if err != nil {
		panic(err)
	}

	// Candles are time series events: request them from a time. osub.Candle(symbol) alone subscribes to the last
	// candle only.
	fromTime := time.Now().Add(-24 * time.Hour).UnixMilli()
	err = subscription.AddSymbol(osub.NewTimeSeriesSubscriptionSymbol(osub.Candle(symbol), fromTime))
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(math.MaxInt64))
}
