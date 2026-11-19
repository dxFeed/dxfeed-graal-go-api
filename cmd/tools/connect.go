package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/parser"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

type Connect struct{}

func (c Connect) ShortDescription() string {
	return "Connects to the address and prints the events."
}

func (c Connect) Run(args []string) {
	dxarguments := DXArguments{args}

	arguments := dxarguments.arguments()
	if len(arguments) < 3 {
		fmt.Print(connectHelp)
		os.Exit(0)
	}

	address := arguments[0]
	symbols, err := parser.ParseSymbols(arguments[2])
	if err != nil {
		fmt.Printf("Error parsing symbols: %v\n", err)
		os.Exit(1)
	}
	types := parser.ParseEventTypes(arguments[1])

	err = connect(address, types, symbols, dxarguments.properties(), dxarguments.forceStream(), dxarguments.isQuite(), dxarguments.time())
	if err != nil {
		fmt.Printf("Error during connect: %v", err)
	}
}

func connect(
	address string,
	types []eventcodes.EventCode,
	symbols []osub.Symbol,
	properties map[string]string,
	forceStream bool,
	isQuite bool,
	fromTime *string,
) error {
	for key, value := range properties {
		if err := api.SetSystemProperty(key, value); err != nil {
			return fmt.Errorf("SetSystemProperty %s: %w", key, err)
		}
	}
	role := api.Feed
	if forceStream {
		role = api.StreamFeed
	}
	endpoint, err := api.NewEndpointWithProperties(role, properties)

	if err != nil {
		return fmt.Errorf("NewEndpointWithProperties: %w", err)
	}

	defer func(endpoint *api.DXEndpoint) {
		_ = endpoint.Close()
	}(endpoint)

	err = endpoint.Connect(address)

	if err != nil {
		return fmt.Errorf("Connect to %s: %w", address, err)
	}

	feed, err := endpoint.GetFeed()

	if err != nil {
		return fmt.Errorf("GetFeed: %w", err)
	}

	subscription, err := feed.CreateSubscription(types...)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	defer func() { _ = subscription.Close() }()
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
			return fmt.Errorf("AddListener: %w", err)
		}
	}

	for _, symbol := range symbols {
		if fromTime != nil {
			parseTime, parseErr := parser.ParseTime(*fromTime)
			if parseErr != nil {
				return parseErr
			}
			timeSeriesSubscriptionSymbol := osub.NewTimeSeriesSubscriptionSymbol(symbol, parseTime)
			err = subscription.AddSymbol(timeSeriesSubscriptionSymbol)
		} else {
			err = subscription.AddSymbol(symbol)
		}
		if err != nil {
			return fmt.Errorf("AddSymbol: %w", err)
		}

	}
	time.Sleep(time.Duration(math.MaxInt64))
	return nil
}
