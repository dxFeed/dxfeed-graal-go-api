// Command price-level-service-client requests the price levels (orders) and the quotes of a candle symbol from the
// ORCS price level service, prints the order sources available to the client and, with -validate, checks the price
// levels with orcs.ValidatePriceLevels. The service needs an authorization token in the address:
//
//	go run ./cmd/samples/price-level-service-client -address "orcs.example.com:7777[login=entitle:<token>]"
//	go run ./cmd/samples/price-level-service-client -address "<address>" -symbol "AAPL{=d,gr=s}" -source NTV \
//		-from 2026-01-10T10:00:00Z -to 2026-01-10T23:00:00Z -validate -verbose
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/orcs"
)

func main() {
	address := flag.String("address", "", "the address of the price level service, with the token: "+
		"host:port[login=entitle:<token>]")
	symbol := flag.String("symbol", "AAPL{=d,gr=s}", "the candle symbol with the granularity attribute (gr)")
	sourceName := flag.String("source", "NTV", "the order source")
	now := time.Now().UTC()
	from := flag.String("from", now.Add(-time.Hour).Format(time.RFC3339), "the start time (RFC 3339)")
	to := flag.String("to", now.Format(time.RFC3339), "the end time (RFC 3339)")
	validate := flag.Bool("validate", false, "check the price levels (the details go to the log of the SDK)")
	verbose := flag.Bool("verbose", false, "print the symbols of each order source")
	flag.Parse()
	if *address == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*address, *symbol, *sourceName, *from, *to, *validate, *verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address, symbolText, sourceName, fromText, toText string, validate, verbose bool) error {
	source, err := order.ValueOfName(sourceName)
	if err != nil {
		return err
	}
	from, err := time.Parse(time.RFC3339, fromText)
	if err != nil {
		return err
	}
	to, err := time.Parse(time.RFC3339, toText)
	if err != nil {
		return err
	}
	symbol := candle.NewCandleSymbol(symbolText)

	service, err := orcs.NewPriceLevelService(address)
	if err != nil {
		return err
	}
	defer func() { _ = service.Close() }()

	authOrderSource, err := service.GetAuthOrderSource()
	if err != nil {
		return err
	}
	var names []string
	for _, e := range authOrderSource.ByOrderSources() {
		names = append(names, *e.Source.Name())
	}
	fmt.Printf("Authorized sources: [%s]\n", strings.Join(names, ", "))
	if verbose {
		for _, e := range authOrderSource.ByOrderSources() {
			fmt.Printf("%s: [%s]\n", *e.Source.Name(), strings.Join(e.Symbols, ", "))
		}
	}

	start := time.Now()
	orders, err := service.GetOrdersWithCaller(symbol, source, from, to, "sample")
	if err != nil {
		return err
	}
	fmt.Printf("Received %d orders in %d ms\n", len(orders), time.Since(start).Milliseconds())
	if validate {
		valid, err := orcs.ValidatePriceLevels(orders, time.Minute, true)
		if err != nil {
			return err
		}
		fmt.Printf("The orders are valid: %v\n", valid)
	}

	start = time.Now()
	quotes, err := service.GetQuotesWithCaller(candle.NewCandleSymbol(symbol.BaseSymbol()+"{gr=s}"), from, to,
		"sample")
	if err != nil {
		return err
	}
	fmt.Printf("Received %d quotes in %d ms\n", len(quotes), time.Since(start).Milliseconds())
	return nil
}
