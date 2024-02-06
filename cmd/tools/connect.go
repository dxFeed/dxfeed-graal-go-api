package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/Osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/parser"
	"math"
	"os"
	"time"
)

type Connect struct{}

func (c Connect) ShortDescription() string {
	return "Connects to specified address(es)."
}

func (c Connect) Run(args []string) {
	dxarguments := DXArguments{args}

	arguments := dxarguments.arguments()
	if len(arguments) < 3 {
		fmt.Println(`
		Connect
		=======
		
		Usage:
		  Connect <address> <types> <symbols> [-f <time>] [<options>]
		
		Where:
			address - The address to connect to retrieve data (remote host or local tape file).
					  To pass an authorization token, add to the address: ""[login=entitle:<token>]"",
					  e.g.: demo.dxfeed.com:7300[login=entitle:<token>]
			types   - Is comma-separated list of dxfeed event types ({eventTypeNames}).
			symbol  - Is comma-separated list of symbol names to get events for (e.g. ""IBM,AAPL,MSFT"").
					  for Candle event specify symbol with aggregation like in ""AAPL{{=d}}""
			--force-stream    Enforces a streaming contract for subscription. The StreamFeed role is used instead of Feed.
		
		Sample: connect "dxlink:wss://demo.dxfeed.com/dxlink-ws" Quote AAPL -p dxfeed.experimental.dxlink.enable=true
		Sample: connect demo.dxfeed.com:7300 Quote AAPL`)
		os.Exit(0)
	}

	address := arguments[0]
	symbols := parser.ParseSymbols(arguments[2])
	types := parser.ParseEventTypes(arguments[1])

	err := connect(address, types, symbols, dxarguments.properties(), dxarguments.forceStream(), dxarguments.isQuite(), dxarguments.time())
	if err != nil {
		fmt.Printf("Error during connect: %v", err)
	}
}

func connect(
	address string,
	types []eventcodes.EventCode,
	symbols []any,
	properties map[string]string,
	forceStream bool,
	isQuite bool,
	fromTime *string,
) error {
	for key, value := range properties {
		api.SetSystemProperty(key, value)
	}
	role := api.Feed
	if forceStream {
		role = api.StreamFeed
	}
	endpoint, err := api.NewEndpointWithProperties(role, properties)

	if err != nil {
		return fmt.Errorf("CreateEndpoint: %we", err)
	}

	defer func(endpoint *api.DXEndpoint) {
		_ = endpoint.Close()
	}(endpoint)

	err = endpoint.Connect(address)

	if err != nil {
		return fmt.Errorf("Connect to %s: %we", address, err)
	}

	feed, err := endpoint.GetFeed()

	if err != nil {
		return fmt.Errorf("GetFeed: %we", err)
	}

	subscription, err := feed.CreateSubscription(types...)

	defer subscription.Close()
	if !isQuite {
		err = subscription.AddListener(PrintEvents(func(eventsList []interface{}) {
			for _, event := range eventsList {
				switch v := event.(type) {
				case events.StringConverter:
					fmt.Printf("%s\n", v.String())
				default:
					fmt.Printf("Unsupported event %T!\n", v)
				}
			}
		}))
		if err != nil {
			return fmt.Errorf("AddListener: %we", err)
		}
	}

	if err != nil {
		return fmt.Errorf("CreateSubscription: %we", err)
	}
	for _, symbol := range symbols {
		if fromTime != nil {
			parseTime, err := parser.ParseTime(*fromTime)
			if err != nil {
				return err
			}
			timeSeriesSubscriptionSymbol := Osub.NewTimeSeriesSubscriptionSymbol(symbol, parseTime)
			subscription.AddSymbol(timeSeriesSubscriptionSymbol)
		} else {
			err = subscription.AddSymbol(symbol)
		}
		if err != nil {
			return fmt.Errorf("AddSymbol: %we", err)
		}

	}
	time.Sleep(time.Duration(math.MaxInt64))
	return nil
}
