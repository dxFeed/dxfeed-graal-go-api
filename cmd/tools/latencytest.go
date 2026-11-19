// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/montanaflynn/stats"
	"math"
	"slices"
	"strings"
	"sync"
	"time"
)

type LatencyTest struct{}

func (c LatencyTest) ShortDescription() string {
	return "Connects to the address and calculates the latency of TimeAndSale events."
}

func (c LatencyTest) Run(args []string) int {
	return runTool("latency test", args, latencyTestHelp, []option{ignoreExchangesOption, forceStreamOption},
		func(a toolArguments) error {
			address, types, symbols, err := parseCommon(a)
			if err != nil {
				return err
			}
			return latency(address, types, symbols, a.has(forceStreamOption), a.value(ignoreExchangesOption))
		})
}

func latency(address string, types []eventcodes.EventCode, symbols []osub.Symbol, forceStream bool, ignoreExchanges *string) error {
	role := api.Feed
	if forceStream {
		role = api.StreamFeed
	}
	var ignoredExchanges []string
	if ignoreExchanges != nil {
		ignoredExchanges = strings.Split(*ignoreExchanges, ",")
	}

	endpoint, err := api.NewEndpoint(role)
	if err != nil {
		return fmt.Errorf("NewEndpoint: %w", err)
	}
	err = endpoint.Connect(address)
	if err != nil {
		return fmt.Errorf("Connect: %w", err)
	}
	feed, err := endpoint.GetFeed()
	if err != nil {
		return fmt.Errorf("GetFeed: %w", err)
	}
	sub, err := feed.CreateSubscription(types...)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	d := createLatency()

	err = sub.AddListener(PrintEvents(func(eventsList []interface{}) {
		currentTime := time.Now().UnixMilli()
		d.mu.Lock()
		d.addListenerCounter(1)
		d.addEventCounter(len(eventsList))
		for _, event := range eventsList {
			switch v := event.(type) {
			case *timeandsale.TimeAndSale:
				{
					if len(ignoredExchanges) > 0 && slices.Contains(ignoredExchanges, formatutil.FormatChar(rune(v.ExchangeCode()))) {
						continue
					}
					if v.IsNew() {
						d.addSymbols(v.EventSymbol())
						delta := float64(currentTime - v.Time())
						d.addDeltas(delta)
					}
				}

			default:
				fmt.Printf("Unsupported event %T!\n", v)
			}
		}
		d.mu.Unlock()
	}))
	if err != nil {
		return fmt.Errorf("AddListener: %w", err)
	}
	for _, symbol := range symbols {
		err = sub.AddSymbol(symbol)
		if err != nil {
			return fmt.Errorf("AddSymbol: %w", err)

		}
	}

	go func() {
		for {
			interval := 2
			time.Sleep(time.Duration(interval) * time.Second)
			d.PrintDiag(interval)
		}
	}()

	time.Sleep(time.Duration(math.MaxInt64))
	return nil
}

func createLatency() *latencyDiag {
	return &latencyDiag{
		deltas:    make([]float64, 0),
		startTime: time.Now(),
		symbols:   make(map[string]string),
		minTotal:  math.MaxFloat64,
		maxTotal:  math.SmallestNonzeroFloat64,
	}
}

type latencyDiag struct {
	listenerCounter int
	eventCounter    int
	startTime       time.Time
	mu              sync.Mutex
	minTotal        float64
	maxTotal        float64
	symbols         map[string]string
	deltas          []float64
}

func (d *latencyDiag) addListenerCounter(i int) {
	d.listenerCounter += i
}

func (d *latencyDiag) addEventCounter(i int) {
	d.eventCounter += i
}

func (d *latencyDiag) addSymbols(symbol *string) {
	d.symbols[*symbol] = *symbol
}

func (d *latencyDiag) addDeltas(delta float64) {
	d.deltas = append(d.deltas, delta)
}

func (d *latencyDiag) PrintDiag(interval int) {
	d.mu.Lock()

	t := time.Now()

	average, _ := stats.Mean(d.deltas)
	minValue, _ := stats.Min(d.deltas)
	maxValue, _ := stats.Max(d.deltas)
	percentile, _ := stats.Percentile(d.deltas, 99)
	stdDev, _ := stats.StdDevP(d.deltas)
	stdErr := stdDev / math.Sqrt(float64(len(d.deltas)))

	eventPerSec := float64(d.eventCounter) / 2.0
	listenerCallsPerSec := float64(d.listenerCounter) / 2.0
	if !math.IsNaN(minValue) {
		d.minTotal = math.Min(d.minTotal, minValue)
	}
	if !math.IsNaN(maxValue) {
		d.maxTotal = math.Max(d.maxTotal, maxValue)
	}

	fmt.Println("----------------------------------------------")
	fmt.Printf("Rate of events (avg)           : %s (events/s)\n", format(eventPerSec))
	fmt.Printf("Rate of listener calls         : %s (calls/s)\n", format(listenerCallsPerSec))
	fmt.Printf("Number of events in call (avg) : %s (events)\n", format(eventPerSec/listenerCallsPerSec))

	fmt.Printf("Rate of unique symbols         : %d symbols/interval\n", len(d.symbols))
	fmt.Printf("Min current                    : %s (ms)\n", format(minValue))
	fmt.Printf("Max current                    : %s (ms)\n", format(maxValue))
	fmt.Printf("Min total                      : %s (ms)\n", format(d.minTotal))
	fmt.Printf("Max total                      : %s (ms)\n", format(d.maxTotal))
	fmt.Printf("99th percentile                : %s (ms)\n", format(percentile))
	fmt.Printf("Mean                           : %s (ms)\n", format(average))
	fmt.Printf("StdDev                         : %s (ms)\n", format(stdDev))
	fmt.Printf("Error                          : %s (ms)\n", format(stdErr))
	fmt.Printf("Sample size (N)                : %d (events)\n", len(d.deltas))

	fmt.Printf("Measurement interval           : %d (s)\n", interval)
	fmt.Printf("Running time                   : %s \n", time.Since(d.startTime))
	fmt.Printf("Timestamp                      : %s \n", t.Format("20060102-150405.000000"))

	d.eventCounter = 0
	d.listenerCounter = 0
	d.symbols = map[string]string{}
	d.deltas = []float64{}
	d.mu.Unlock()
}

func format(value float64) string {
	if math.IsNaN(value) {
		return "---"
	} else if math.IsInf(value, 1) {
		return "---"
	} else if value == math.MaxFloat64 {
		return "---"
	}
	return fmt.Sprintf("%.2f", value)
}
