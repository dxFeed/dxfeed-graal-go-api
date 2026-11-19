// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

type Dump struct{}

func (c Dump) ShortDescription() string {
	return "Dumps all data and subscription information received from address."
}

func (c Dump) Run(args []string) int {
	return runTool("dump", args, dumpHelp, []option{tapeOption, propertiesOption, quietOption},
		func(a toolArguments) error {
			properties, err := a.properties()
			if err != nil {
				return usageError{err}
			}
			address, types, symbols, err := parseCommon(a)
			if err != nil {
				return err
			}
			return dump(address, a.value(tapeOption), symbols, types, properties, a.has(quietOption))
		})
}

func dump(
	inputFile string,
	outputFile *string,
	symbols []osub.Symbol,
	types []eventcodes.EventCode,
	properties map[string]string,
	quiet bool,
) error {
	for key, value := range properties {
		if err := api.SetSystemProperty(key, value); err != nil {
			return fmt.Errorf("SetSystemProperty %s: %w", key, err)
		}
	}
	var listeners []common.EventListener

	if !quiet {
		listeners = append(listeners, DumpEvents(func(eventsList []interface{}) {
			for _, event := range eventsList {
				switch v := event.(type) {
				case events.StringConverter:
					fmt.Printf("%s\n", v.String())
				default:
					fmt.Printf("Unsupported event %T!\n", v)
				}
			}
		}))
	}

	inputEndpoint, err := api.NewEndpointWithProperties(api.StreamFeed, properties)
	if err != nil {
		return fmt.Errorf("NewEndpoint: %w", err)
	}

	if err != nil {
		return fmt.Errorf("Connect to %s: %w", inputFile, err)
	}
	feed, err := inputEndpoint.GetFeed()
	if err != nil {
		return fmt.Errorf("GetFeed: %w", err)
	}

	subscription, err := feed.CreateSubscription(types...)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	defer func() { _ = subscription.Close() }()
	count := 0
	var outputEndpoint *api.DXEndpoint

	if outputFile != nil {
		outputEndpoint, err = api.NewEndpointWithProperties(api.StreamPublisher, properties)
		if err != nil {
			return fmt.Errorf("NewEndpoint Publisher: %w", err)
		}

		publisher, err := outputEndpoint.GetPublisher()
		if err != nil {
			return fmt.Errorf("GetPublisher: %w", err)
		}

		listeners = append(listeners, DumpEvents(func(eventsList []interface{}) {
			count = count + len(eventsList)
			err := publisher.Publish(eventsList)
			if err != nil {
				fmt.Printf("Publish error %v", err)
			}
		}))
		err = outputEndpoint.Connect(fmt.Sprintf("tape:%s", *outputFile))
		if err != nil {
			return fmt.Errorf("Connect to %s: %w", *outputFile, err)
		}
	}
	for _, listener := range listeners {
		err = subscription.AddListener(listener)
		if err != nil {
			return fmt.Errorf("AddListener: %w", err)
		}
	}

	for _, symbol := range symbols {
		err = subscription.AddSymbol(symbol)
		if err != nil {
			return fmt.Errorf("AddSymbol: %w", err)
		}
	}

	err = inputEndpoint.Connect(inputFile)
	if err != nil {
		return fmt.Errorf("Connect to %s: %w", inputFile, err)
	}

	err = inputEndpoint.AwaitNotConnected(context.Background())
	if err != nil {
		return fmt.Errorf("AwaitNotConnected: %w", err)
	}
	err = inputEndpoint.CloseAndAwaitTermination(context.Background())
	if err != nil {
		return fmt.Errorf("CloseAndAwaitTermination: %w", err)
	}
	if outputEndpoint != nil {
		err = outputEndpoint.AwaitProcessed(context.Background())
		if err != nil {
			return fmt.Errorf("AwaitProcessed: %w", err)
		}
		err = outputEndpoint.CloseAndAwaitTermination(context.Background())
		if err != nil {
			return fmt.Errorf("CloseAndAwaitTermination: %w", err)
		}
	}
	fmt.Printf("Published %d events\n", count)
	return nil
}

type DumpEvents func(events []interface{})

func (dm DumpEvents) Update(events []any) {
	dm(events)
}
