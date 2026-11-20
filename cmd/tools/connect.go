// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/parser"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/toolargs"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

type Connect struct{}

func (c Connect) ShortDescription() string {
	return "Connects to the address and prints the events."
}

func (c Connect) Run(args []string) int {
	return runTool("connect", args, connectHelp, []toolargs.Option{toolargs.FromTime, toolargs.Source,
		toolargs.Properties, toolargs.Tape, toolargs.Quiet, toolargs.ForceStream}, func(a toolargs.Args) error {
		fromTime, source, err := parseSubscriptionOptions(a)
		if err != nil {
			return err
		}
		common, err := parseCommon(a, "", "")
		if err != nil {
			return err
		}
		return connect(common, fromTime, source, a.Value(toolargs.Tape), a.Has(toolargs.ForceStream),
			a.Has(toolargs.Quiet))
	})
}

// parseSubscriptionOptions parses -f and -s: the time of a time series subscription and the order source of an
// indexed one. They cannot be used together.
func parseSubscriptionOptions(a toolargs.Args) (fromTime *int64, source *order.Source, err error) {
	from, src := a.Value(toolargs.FromTime), a.Value(toolargs.Source)
	if from != nil && src != nil {
		return nil, nil, usageError{errors.New("options --from-time and --source cannot be used together")}
	}
	if from != nil {
		t, err := parser.ParseTime(*from)
		if err != nil {
			return nil, nil, usageError{fmt.Errorf("option --from-time: %w", err)}
		}
		fromTime = &t
	}
	if src != nil {
		source, err = order.ValueOfName(*src)
		if err != nil {
			return nil, nil, usageError{fmt.Errorf("option --source: %w", err)}
		}
	}
	return fromTime, source, nil
}

func connect(common commonArgs, fromTime *int64, source *order.Source, tape *string, forceStream, quiet bool) error {
	role := api.Feed
	if forceStream {
		role = api.StreamFeed
	}
	endpoint, err := newEndpoint(role, common.properties)
	if err != nil {
		return err
	}
	defer func() { _ = endpoint.Close() }()

	feed, err := endpoint.GetFeed()
	if err != nil {
		return fmt.Errorf("GetFeed: %w", err)
	}
	subscription, err := feed.CreateSubscription(common.types...)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	defer func() { _ = subscription.Close() }()
	if !quiet {
		if err := subscription.AddListener(printEvents()); err != nil {
			return fmt.Errorf("AddListener: %w", err)
		}
	}
	if tape != nil {
		writer, err := newTapeWriter(*tape, common.properties)
		if err != nil {
			return err
		}
		defer writer.close()
		if err := subscription.AddListener(writer); err != nil {
			return fmt.Errorf("AddListener: %w", err)
		}
	}
	for _, symbol := range common.symbols {
		var s osub.Symbol = symbol
		switch {
		case fromTime != nil:
			s = osub.NewTimeSeriesSubscriptionSymbol(symbol, *fromTime)
		case source != nil:
			s = osub.NewIndexedEventSubscriptionSymbol(symbol, source)
		}
		if err := subscription.AddSymbol(s); err != nil {
			return fmt.Errorf("AddSymbol: %w", err)
		}
	}
	if err := endpoint.Connect(common.address); err != nil {
		return fmt.Errorf("Connect to %s: %w", common.address, err)
	}
	time.Sleep(time.Duration(math.MaxInt64))
	return nil
}

// printEvents is a listener that prints each event on a line.
func printEvents() PrintEvents {
	return func(eventsList []interface{}) {
		for _, event := range eventsList {
			switch v := event.(type) {
			case events.StringConverter:
				fmt.Printf("%s\n", v.String())
			default:
				fmt.Printf("Unsupported event %T!\n", v)
			}
		}
	}
}
