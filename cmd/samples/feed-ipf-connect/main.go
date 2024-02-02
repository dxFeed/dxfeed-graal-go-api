package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
	"math"
	"time"
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

type PrintState func(old common.ConnectionState, new common.ConnectionState)

func (pr PrintState) UpdateState(old common.ConnectionState, new common.ConnectionState) {
	pr(old, new)
}

// An sample that demonstrates a subscription using InstrumentProfile.
func main() {
	// The experimental property must be enabled.
	api.SetSystemProperty("dxfeed.experimental.dxlink.enable", "true")
	// Set scheme for dxLink.
	api.SetSystemProperty("scheme", "ext:opt:sysprops,resource:dxlink.xml")

	ipfFile := "https://demo:demo@tools.dxfeed.com/ipf?TYPE=STOCK&compression=zip"
	fmt.Printf("Reading instruments from %s\n", ipfFile)
	ipfReader, err := ipf.NewInstrumentProfileReader()
	if err != nil {
		panic(err)
	}
	defer func(reader *ipf.InstrumentProfileReader) {
		_ = reader.Close()
	}(ipfReader)

	profiles, err := ipfReader.ReadFromFile(ipfFile)
	if err != nil {
		panic(err)
	}
	var symbols []any
	for _, profile := range profiles {
		symbols = append(symbols, *profile.Symbol())
		fmt.Printf("%s\n", profile.String())
	}
	completed, _ := ipfReader.WasComplete()
	fmt.Printf("Was Completed %t\n", completed)
	lastModified, _ := ipfReader.GetLastModified()
	fmt.Printf("LastModified %d\n", lastModified)

	// For token-based authorization, use the following address format:
	// "dxlink:wss://demo.dxfeed.com/dxlink-ws[login=dxlink:token]"
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
	err = endpoint.Connect("dxlink:wss://demo.dxfeed.com/dxlink-ws")
	if err != nil {
		panic(err)
	}

	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	subscription, err := feed.CreateSubscription(eventcodes.Order)
	if err != nil {
		panic(err)
	}
	defer subscription.Close()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			switch v := event.(type) {
			case *quote.Quote:
				fmt.Printf("%s\n", v.String())
			case *profile.Profile:
				fmt.Printf("%s\n", v.String())
			case *timeandsale.TimeAndSale:
				fmt.Printf("%s\n", v.String())
			case *order.Order:
				fmt.Printf("%s\n", v.String())
			}
		}
	}))

	err = subscription.AddSymbols(symbols...)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(math.MaxInt64))
}
