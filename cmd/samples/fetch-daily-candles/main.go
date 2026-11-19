// Command fetch-daily-candles fetches the daily candles of a symbol for the last 20 days, prints them and exits.
//
//	go run ./cmd/samples/fetch-daily-candles -address demo.dxfeed.com:7300 -symbol AAPL
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

const days = 20

func main() {
	address := flag.String("address", "demo.dxfeed.com:7300", "the address to connect to")
	baseSymbol := flag.String("symbol", "AAPL", "the symbol of the instrument, for example IBM, AAPL or SPX")
	timeout := flag.Duration("timeout", 5*time.Second, "how long to wait for the candles")
	flag.Parse()

	symbol, err := candle.NewCandleSymbolWithAttributes(*baseSymbol, candle.NewCandlePeriod(1, candle.CandleTypeDay))
	if err != nil {
		panic(err)
	}
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()
	if err := endpoint.Connect(*address); err != nil {
		panic(err)
	}
	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fetching the last %d days of candles of %s\n", days, *baseSymbol)
	to := time.Now()
	from := to.AddDate(0, 0, -days)
	promise, err := feed.GetTimeSeriesPromise(eventcodes.Candle, osub.Candle(symbol), from.UnixMilli(), to.UnixMilli())
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	candles, err := promise.Await(ctx)
	if err != nil {
		panic(err)
	}
	for _, c := range candles {
		fmt.Println(c)
	}
}
