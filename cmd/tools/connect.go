// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"math"
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

func (c Connect) Run(args []string) int {
	return runTool("connect", args, connectHelp,
		[]option{fromTimeOption, propertiesOption, quietOption, forceStreamOption}, func(a toolArguments) error {
			properties, err := a.properties()
			if err != nil {
				return usageError{err}
			}
			address, types, symbols, err := parseCommon(a)
			if err != nil {
				return err
			}
			return connect(address, types, symbols, properties, a.has(forceStreamOption), a.has(quietOption),
				a.value(fromTimeOption))
		})
}

func connect(
	address string,
	types []eventcodes.EventCode,
	symbols []osub.Symbol,
	properties map[string]string,
	forceStream bool,
	quiet bool,
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
	if !quiet {
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
