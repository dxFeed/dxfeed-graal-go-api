// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/toolargs"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
)

type Dump struct{}

func (c Dump) ShortDescription() string {
	return "Dumps all data and subscription information received from address."
}

func (c Dump) Run(args []string) int {
	return runTool("dump", args, dumpHelp, []toolargs.Option{toolargs.Tape, toolargs.Properties, toolargs.Quiet},
		func(a toolargs.Args) error {
			common, err := parseCommon(a, "all", "all")
			if err != nil {
				return err
			}
			return dump(common, a.Value(toolargs.Tape), a.Has(toolargs.Quiet))
		})
}

func dump(common commonArgs, tape *string, quiet bool) error {
	inputEndpoint, err := newEndpoint(api.StreamFeed, common.properties)
	if err != nil {
		return err
	}
	defer func() { _ = inputEndpoint.Close() }()
	feed, err := inputEndpoint.GetFeed()
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
	var writer *tapeWriter
	if tape != nil {
		if writer, err = newTapeWriter(*tape, common.properties); err != nil {
			return err
		}
		defer writer.close()
		if err := subscription.AddListener(writer); err != nil {
			return fmt.Errorf("AddListener: %w", err)
		}
	}
	for _, symbol := range common.symbols {
		if err := subscription.AddSymbol(symbol); err != nil {
			return fmt.Errorf("AddSymbol: %w", err)
		}
	}
	if err := inputEndpoint.Connect(common.address); err != nil {
		return fmt.Errorf("Connect to %s: %w", common.address, err)
	}
	if err := inputEndpoint.AwaitNotConnected(context.Background()); err != nil {
		return fmt.Errorf("AwaitNotConnected: %w", err)
	}
	if err := inputEndpoint.CloseAndAwaitTermination(context.Background()); err != nil {
		return fmt.Errorf("CloseAndAwaitTermination: %w", err)
	}
	if writer != nil {
		if err := writer.finish(); err != nil {
			return err
		}
		fmt.Printf("Published %d events\n", writer.count.Load())
	}
	return nil
}

// tapeWriter is a listener that publishes the events it gets to a tape file through a StreamPublisher endpoint, as
// the -t option of the tools does.
type tapeWriter struct {
	endpoint  *api.DXEndpoint
	publisher *api.DXPublisher
	count     atomic.Int64
}

func newTapeWriter(file string, properties map[string]string) (*tapeWriter, error) {
	endpoint, err := newEndpoint(api.StreamPublisher, properties)
	if err != nil {
		return nil, err
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		_ = endpoint.Close()
		return nil, fmt.Errorf("GetPublisher: %w", err)
	}
	if err := endpoint.Connect("tape:" + file); err != nil {
		_ = endpoint.Close()
		return nil, fmt.Errorf("Connect to tape:%s: %w", file, err)
	}
	return &tapeWriter{endpoint: endpoint, publisher: publisher}, nil
}

// Update publishes the events to the tape.
func (w *tapeWriter) Update(events []any) {
	w.count.Add(int64(len(events)))
	if err := w.publisher.Publish(events); err != nil {
		fmt.Fprintf(os.Stderr, "Publish error: %v\n", err)
	}
}

// finish waits until the published events are written and closes the endpoint.
func (w *tapeWriter) finish() error {
	if err := w.endpoint.AwaitProcessed(context.Background()); err != nil {
		return fmt.Errorf("AwaitProcessed: %w", err)
	}
	if err := w.endpoint.CloseAndAwaitTermination(context.Background()); err != nil {
		return fmt.Errorf("CloseAndAwaitTermination: %w", err)
	}
	return nil
}

// close closes the endpoint without waiting; it does nothing after finish.
func (w *tapeWriter) close() {
	_ = w.endpoint.Close()
}
