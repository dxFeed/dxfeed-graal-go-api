// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/toolargs"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
)

// perfTestInterval is the measurement interval of PerfTest, as in the C++ tool.
const perfTestInterval = 2 * time.Second

type PerfTest struct{}

func (c PerfTest) ShortDescription() string {
	return "Connects to specified address and calculates performance counters."
}

func (c PerfTest) Run(args []string) int {
	return runTool("perftest", args, perfTestHelp, []toolargs.Option{toolargs.Properties, toolargs.ForceStream,
		toolargs.CPUUsageByCore, toolargs.DetachListener}, func(a toolargs.Args) error {
		common, err := parseCommon(a, "", "")
		if err != nil {
			return err
		}
		return perf(common, a.Has(toolargs.ForceStream), a.Has(toolargs.CPUUsageByCore),
			a.Has(toolargs.DetachListener))
	})
}

func perf(common commonArgs, forceStream, cpuUsageByCore, detachListener bool) error {
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
	sub, err := feed.CreateSubscription(common.types...)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	defer func() { _ = sub.Close() }()
	d := &diag{cpuUsageByCore: cpuUsageByCore}
	if !detachListener {
		err = sub.AddListener(PrintEvents(func(eventsList []interface{}) {
			d.mu.Lock()
			d.listenerCounter++
			d.eventCounter += len(eventsList)
			d.mu.Unlock()
		}))
		if err != nil {
			return fmt.Errorf("AddListener: %w", err)
		}
	}
	for _, symbol := range common.symbols {
		if err := sub.AddSymbol(symbol); err != nil {
			return fmt.Errorf("AddSymbol: %w", err)
		}
	}
	if err := endpoint.Connect(common.address); err != nil {
		return fmt.Errorf("Connect: %w", err)
	}
	d.start()
	for range time.Tick(perfTestInterval) {
		d.PrintDiag()
	}
	return nil
}

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

type diag struct {
	listenerCounter int
	eventCounter    int
	cpuUsageByCore  bool
	lastTime        time.Time
	lastCPUTimeNs   int64
	mu              sync.Mutex
}

// start takes the time and the CPU time from which the first interval counts.
func (d *diag) start() {
	d.lastTime = time.Now()
	d.lastCPUTimeNs = -1
	if s, err := api.GetRuntimeStatistics(); err == nil {
		d.lastCPUTimeNs = s.ProcessCPUTimeNs
	}
}

// cpuUsage returns the CPU usage of the process since the last call in percent: of all the processors, or of one
// core with --cpu-usage-by-core (a fully loaded core is 100%); NaN when the CPU time of this interval or of the last
// one is not available.
func (d *diag) cpuUsage(now time.Time) float64 {
	last := d.lastCPUTimeNs
	d.lastCPUTimeNs = -1
	s, err := api.GetRuntimeStatistics()
	if err != nil || s.ProcessCPUTimeNs < 0 {
		return math.NaN()
	}
	d.lastCPUTimeNs = s.ProcessCPUTimeNs
	if last < 0 {
		return math.NaN()
	}
	elapsed := now.Sub(d.lastTime)
	usage := float64(s.ProcessCPUTimeNs-last) / float64(elapsed.Nanoseconds()) * 100
	if !d.cpuUsageByCore && s.AvailableProcessors > 0 {
		usage /= float64(s.AvailableProcessors)
	}
	return usage
}

func (d *diag) PrintDiag() {
	now := time.Now()
	cpu := d.cpuUsage(now)
	d.mu.Lock()
	seconds := now.Sub(d.lastTime).Seconds()
	eventPerSec := float64(d.eventCounter) / seconds
	listenerCallsPerSec := float64(d.listenerCounter) / seconds
	d.eventCounter = 0
	d.listenerCounter = 0
	d.mu.Unlock()
	d.lastTime = now
	fmt.Println("----------------------------------------------")
	fmt.Printf("Rate of events (avg)           : %f (events/s)\n", eventPerSec)
	fmt.Printf("Rate of listener calls         : %f (calls/s)\n", listenerCallsPerSec)
	fmt.Printf("Number of events in call (avg) : %f (events)\n", eventPerSec/listenerCallsPerSec)
	fmt.Printf("CPU usage                      : %s (%%)\n", format(cpu))
}
