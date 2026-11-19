// Command candle-data-download downloads the daily candles of a symbol for the last month from the dxFeed candle web
// service and prints them. The service needs credentials:
//
//	go run ./cmd/samples/candle-data-download -user <user> -password <password>
//	go run ./cmd/samples/candle-data-download -token <token> -symbol "IBM{=h}" -days 3
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/candlewebservice"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func main() {
	address := flag.String("address", "https://tools.dxfeed.com/candledata", "the address of the candle web service")
	user := flag.String("user", "", "the user name")
	password := flag.String("password", "", "the password")
	token := flag.String("token", "", "the authentication token, instead of the user name and the password")
	symbol := flag.String("symbol", "AAPL{=d}", "the candle symbol")
	days := flag.Int("days", 30, "the number of the last days to download")
	flag.Parse()

	builder := candlewebservice.NewHistoryEndpointBuilder().WithAddress(*address)
	if *token != "" {
		builder.WithAuthToken(*token)
	} else {
		builder.WithUserName(*user).WithPassword(*password)
	}
	endpoint, err := builder.Build()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = endpoint.Close() }()

	to := time.Now()
	from := to.AddDate(0, 0, -*days)
	candles, err := endpoint.GetTimeSeries(eventcodes.Candle, osub.String(*symbol), from, to)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%d candles of %s from %s to %s\n", len(candles), *symbol, from.Format(time.DateOnly),
		to.Format(time.DateOnly))
	for _, c := range candles {
		fmt.Println(c)
	}
}
