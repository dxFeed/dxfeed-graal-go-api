// Command live-ipf keeps instrument profiles up to date from an instrument profile source and prints the updates: the
// number of the new and changed profiles and the removed ones. The source is a URL, a file or host:port of an
// instrument profile server; the dxFeed demo service is the default:
//
//	go run ./cmd/samples/live-ipf
//	go run ./cmd/samples/live-ipf -address "profiles.ipf[update=10s]" -duration 1m
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf/live"
)

type printUpdates struct{}

func (printUpdates) InstrumentProfilesUpdated(profiles []*events.InstrumentProfile) {
	removed := 0
	for _, p := range profiles {
		if *p.InstrumentType() == ipf.TypeRemoved {
			removed++
			fmt.Printf("removed %s\n", *p.Symbol())
		}
	}
	fmt.Printf("%s: %d profiles updated, %d removed\n", time.Now().Format(time.TimeOnly), len(profiles)-removed, removed)
}

type printStates struct{}

func (printStates) StateChanged(oldState, newState live.State) {
	fmt.Printf("connection: %v -> %v\n", oldState, newState)
}

func main() {
	address := flag.String("address", "https://demo:demo@tools.dxfeed.com/ipf", "the instrument profile source")
	duration := flag.Duration("duration", 0, "how long to run; 0 runs until interrupted")
	flag.Parse()

	if err := run(*address, *duration); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address string, duration time.Duration) error {
	collector, err := live.NewInstrumentProfileCollector()
	if err != nil {
		return err
	}
	defer func() { _ = collector.Close() }()
	connection, err := live.CreateConnection(address, collector)
	if err != nil {
		return err
	}
	defer func() { _ = connection.Close() }()

	if err := connection.AddStateChangeListener(printStates{}); err != nil {
		return err
	}
	if err := collector.AddUpdateListener(printUpdates{}); err != nil {
		return err
	}
	if err := connection.Start(); err != nil {
		return err
	}
	if completed, err := connection.WaitUntilCompleted(time.Minute); err != nil {
		return err
	} else if !completed {
		return fmt.Errorf("the profiles from %s were not read in a minute", address)
	}
	profiles, err := collector.View()
	if err != nil {
		return err
	}
	fmt.Printf("%d profiles read; waiting for updates\n", len(profiles))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	var timeout <-chan time.Time
	if duration > 0 {
		timeout = time.After(duration)
	}
	select {
	case <-interrupt:
	case <-timeout:
	}
	return nil
}
