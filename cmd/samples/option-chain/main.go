// Command option-chain reads instrument profiles, builds the option chain of a symbol and prints its nearest series
// with the strikes around a price, and the call and put symbols at each strike. The dxFeed demo service is the default
// source of the profiles:
//
//	go run ./cmd/samples/option-chain -symbol AAPL -price 200
//	go run ./cmd/samples/option-chain -ipf pkg/ipf/testdata/profiles.ipf -symbol TEST -strikes 2
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf/option"
)

func main() {
	address := flag.String("ipf", "https://demo:demo@tools.dxfeed.com/ipf", "the instrument profile file or URL")
	symbol := flag.String("symbol", "AAPL", "the product or underlying symbol of the chain")
	strikes := flag.Int("strikes", 10, "the number of strikes around the price")
	series := flag.Int("series", 3, "the number of the nearest series")
	price := flag.Float64("price", 0, "the price to choose the strikes around; 0 takes the middle strike of each series")
	flag.Parse()

	if err := run(*address, *symbol, *strikes, *series, *price); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address, symbol string, nStrikes, nSeries int, price float64) error {
	reader, err := ipf.NewInstrumentProfileReader()
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()
	profiles, err := reader.ReadFromFile(address)
	if err != nil {
		return err
	}
	chain := option.BuildChains(profiles)[symbol]
	if chain == nil {
		return fmt.Errorf("no options of %s in %d profiles", symbol, len(profiles))
	}
	allSeries := chain.Series()
	fmt.Printf("%s: %d series\n", symbol, len(allSeries))
	for _, s := range allSeries[:max(0, min(nSeries, len(allSeries)))] {
		fmt.Printf("\n%s\n", s)
		around := price
		if around == 0 {
			strikes := s.Strikes()
			around = strikes[len(strikes)/2]
		}
		calls, puts := s.Calls(), s.Puts()
		for _, strike := range s.NStrikesAround(nStrikes, around) {
			fmt.Printf("%12g  %-24s %s\n", strike, symbolOf(calls[strike]), symbolOf(puts[strike]))
		}
	}
	return nil
}

func symbolOf(profile *events.InstrumentProfile) string {
	if profile == nil {
		return "-"
	}
	return *profile.Symbol()
}
