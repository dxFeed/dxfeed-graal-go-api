// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/nettest"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
)

// The limits of the consumer: the events checked for the nanosecond part and the times without data.
const (
	nanoCheckEvents        = 1000
	noDataAfterFirstRecord = 30 * time.Second
	noDataAfterConnection  = 60 * time.Second
	nanoTimeProperty       = "dxscheme.nanoTime"
)

type NetTest struct{}

func (c NetTest) ShortDescription() string {
	return "Publishes or receives generated events at a controlled rate and prints the rate and the lag."
}

// Run runs NetTest as the NetTest specification nettest/1 defines it: an invalid command line or none exits with 2,
// an error of the run with 1, the end of --duration or a signal with 0.
func (c NetTest) Run(args []string) int {
	for _, arg := range args {
		if nettest.IsHelp(arg) {
			fmt.Print(nettest.HelpText)
			return 0
		}
	}
	config, err := nettest.Parse(args)
	if err == nettest.ErrNoArguments {
		fmt.Fprint(os.Stderr, nettest.HelpText)
		return 2
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "nettest: %v\n\n%s", err, nettest.HelpText)
		return 2
	}
	if err := runNetTest(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Error during nettest: %v\n", err)
		return 1
	}
	return 0
}

// netTestInstance is an instance of the tool: an endpoint and its counts of the current interval.
type netTestInstance struct {
	endpoint *api.DXEndpoint
	records  atomic.Int64
	lagSum   atomic.Int64
}

// netTestRun is a run of the tool.
type netTestRun struct {
	config    *nettest.Config
	instances []*netTestInstance
	// stop ends the run; failure holds the error that ended it, if any.
	stop     chan struct{}
	stopOnce sync.Once
	failure  error
	// workers are the goroutines that use the endpoints: the publishers and the watchdog of the consumer.
	workers sync.WaitGroup
	// The consumer: the events received, those with a nanosecond part, and the time of the last record.
	received     atomic.Int64
	withNano     atomic.Int64
	lastRecordNs atomic.Int64
	nanoCheck    bool
}

func (r *netTestRun) fail(err error) {
	r.stopOnce.Do(func() {
		r.failure = err
		close(r.stop)
	})
}

func (r *netTestRun) end() { r.fail(nil) }

func runNetTest(config *nettest.Config) error {
	symbols, err := netTestSymbols(config)
	if err != nil {
		return err
	}
	if config.ServerAddress() {
		// The SDK binds the port in the background and only logs a failure: check that the port is free first.
		listener, err := net.Listen("tcp", config.Address)
		if err != nil {
			return fmt.Errorf("the address %s cannot be listened on: %w", config.Address, err)
		}
		_ = listener.Close()
	}
	for _, p := range config.Properties {
		if err := api.SetSystemProperty(p.Key, p.Value); err != nil {
			return fmt.Errorf("SetSystemProperty %s: %w", p.Key, err)
		}
	}
	nanoTime := "true"
	if value, ok := config.Property(nanoTimeProperty); ok {
		nanoTime = value
	}
	r := &netTestRun{config: config, stop: make(chan struct{}), nanoCheck: nanoTime == "true"}
	defer func() {
		// Stop the goroutines that use the endpoints before closing them.
		r.end()
		r.workers.Wait()
		r.close()
	}()
	for i := 1; i <= config.Instances; i++ {
		builder := api.NewEndpointBuilder().WithRole(netTestRole(config.Role)).
			WithName(config.Name+"."+strconv.Itoa(i)).WithProperty(nanoTimeProperty, nanoTime)
		for _, p := range config.Properties {
			if p.Key != nanoTimeProperty {
				builder.WithProperty(p.Key, p.Value)
			}
		}
		endpoint, err := builder.Build()
		if err != nil {
			return fmt.Errorf("the endpoint %s.%d: %w", config.Name, i, err)
		}
		r.instances = append(r.instances, &netTestInstance{endpoint: endpoint})
	}
	if config.Mode == nettest.Consumer {
		err = r.startConsumers(symbols)
	} else {
		err = r.startProducers(symbols)
	}
	if err != nil {
		return err
	}
	r.waitAndPrint()
	return r.failure
}

// close closes the endpoints.
func (r *netTestRun) close() {
	for _, instance := range r.instances {
		_ = instance.endpoint.CloseAndAwaitTermination(context.Background())
	}
}

func netTestRole(role string) common.Role {
	switch role {
	case nettest.RoleStreamPublisher:
		return api.StreamPublisher
	case nettest.RoleFeed:
		return api.Feed
	case nettest.RoleStreamFeed:
		return api.StreamFeed
	}
	return api.Publisher
}

// netTestSymbols returns the generated symbols or those of the IPF file.
func netTestSymbols(config *nettest.Config) ([]string, error) {
	if config.SymbolFile == "" {
		symbols, _ := nettest.GenerateSymbols(config.Symbols)
		return symbols, nil
	}
	reader, err := ipf.NewInstrumentProfileReader()
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()
	profiles, err := reader.ReadFromFile(config.SymbolFile)
	if err != nil {
		return nil, fmt.Errorf("the symbol file %s: %w", config.SymbolFile, err)
	}
	var symbols []string
	for _, profile := range profiles {
		if len(symbols) == config.Symbols {
			break
		}
		if symbol := profile.Symbol(); symbol != nil && *symbol != "" {
			symbols = append(symbols, *symbol)
		}
	}
	switch {
	case len(symbols) == 0:
		return nil, fmt.Errorf("the symbol file %s has no symbols", config.SymbolFile)
	case len(symbols) < config.Symbols:
		fmt.Fprintf(os.Stderr, "Warning: the symbol file %s has %d symbols, fewer than %d\n", config.SymbolFile,
			len(symbols), config.Symbols)
	}
	return symbols, nil
}

func (r *netTestRun) eventCode() eventcodes.EventCode {
	if r.config.Event == nettest.EventTimeAndSale {
		return eventcodes.TimeAndSale
	}
	return eventcodes.Trade
}

// startConsumers subscribes each instance to its slice of the symbols, the listener before the symbols, and connects
// it last.
func (r *netTestRun) startConsumers(symbols []string) error {
	for i, instance := range r.instances {
		feed, err := instance.endpoint.GetFeed()
		if err != nil {
			return err
		}
		subscription, err := feed.CreateSubscription(r.eventCode())
		if err != nil {
			return err
		}
		if err := subscription.AddListener(&netTestListener{run: r, instance: instance}); err != nil {
			return err
		}
		slice := nettest.Slice(symbols, i+1, len(r.instances))
		subscriptionSymbols := make([]osub.Symbol, len(slice))
		for j, symbol := range slice {
			subscriptionSymbols[j] = osub.String(symbol)
		}
		if err := subscription.AddSymbols(subscriptionSymbols...); err != nil {
			return err
		}
	}
	for _, instance := range r.instances {
		if err := instance.endpoint.Connect(r.config.Address); err != nil {
			return fmt.Errorf("connect to %s: %w", r.config.Address, err)
		}
	}
	r.workers.Add(1)
	go r.watchConsumer()
	return nil
}

// netTestListener counts the events of a consumer instance and the sum of their e2e lags.
type netTestListener struct {
	run      *netTestRun
	instance *netTestInstance
}

func (l *netTestListener) Update(events []any) {
	r := l.run
	now := time.Now().UnixNano()
	var lagSum, withNano int64
	for _, event := range events {
		var timeMillis int64
		var nano int32
		switch e := event.(type) {
		case *trade.Trade:
			timeMillis, nano = e.Time(), e.TimeNanoPart()
		case *timeandsale.TimeAndSale:
			timeMillis, nano = e.Time(), e.TimeNanoPart()
		}
		if nano != 0 {
			withNano++
		}
		lagSum += now - (timeMillis*1000000 + int64(nano))
	}
	l.instance.records.Add(int64(len(events)))
	l.instance.lagSum.Add(lagSum)
	r.lastRecordNs.Store(now)
	r.withNano.Add(withNano)
	if received := r.received.Add(int64(len(events))); r.nanoCheck && received >= nanoCheckEvents &&
		received-int64(len(events)) < nanoCheckEvents && r.withNano.Load() == 0 {
		r.fail(fmt.Errorf("nano_missing: the first %d events have no nanosecond part: a hop of the path lacks "+
			"%s=true (or run with -D%s=false)", received, nanoTimeProperty, nanoTimeProperty))
	}
}

// watchConsumer stops the consumer without data: with open connections and no records for 30 s after the last
// record, or for 60 s after the connection when no record came.
func (r *netTestRun) watchConsumer() {
	defer r.workers.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var connectedSince time.Time
	for {
		select {
		case <-r.stop:
			return
		case now := <-ticker.C:
			if r.connections() == 0 {
				connectedSince = time.Time{}
				continue
			}
			if connectedSince.IsZero() {
				connectedSince = now
			}
			last := r.lastRecordNs.Load()
			// After a reconnect the time without records counts from the connection.
			since := time.Unix(0, last)
			if since.Before(connectedSince) {
				since = connectedSince
			}
			switch {
			case last == 0 && now.Sub(connectedSince) >= noDataAfterConnection:
				r.fail(fmt.Errorf("no_data: no records for %v after the connection", noDataAfterConnection))
			case last != 0 && now.Sub(since) >= noDataAfterFirstRecord:
				r.fail(fmt.Errorf("no_data: no records for %v", noDataAfterFirstRecord))
			}
		}
	}
}

// startProducers connects each instance and starts its publishing loop.
func (r *netTestRun) startProducers(symbols []string) error {
	var schedule *nettest.Schedule
	if r.config.Rate > 0 {
		schedule = nettest.NewSchedule(r.config.Rate, r.config.Batch)
	}
	for i, instance := range r.instances {
		if err := instance.endpoint.Connect(r.config.Address); err != nil {
			return fmt.Errorf("connect to %s: %w", r.config.Address, err)
		}
		publisher, err := instance.endpoint.GetPublisher()
		if err != nil {
			return err
		}
		r.workers.Add(1)
		go r.publish(i+1, instance, publisher, symbols, schedule)
	}
	return nil
}

// publish publishes the batches of an instance until the run stops: wait for the schedule, read the clock, fill the
// batch from the pool and publish it.
func (r *netTestRun) publish(number int, instance *netTestInstance, publisher *api.DXPublisher, symbols []string,
	schedule *nettest.Schedule) {
	defer r.workers.Done()
	random := nettest.NewRandom(number)
	table := nettest.NewTable(random)
	timeAndSale := r.config.Event == nettest.EventTimeAndSale
	var trades []*trade.Trade
	var sales []*timeandsale.TimeAndSale
	for _, symbol := range symbols {
		if timeAndSale {
			sales = append(sales, timeandsale.NewTimeAndSale(symbol))
		} else {
			trades = append(trades, newPoolTrade(symbol))
		}
	}
	batch := make([]any, r.config.Batch)
	sequence := 0
	timer := time.NewTimer(0)
	<-timer.C
	for {
		select {
		case <-r.stop:
			return
		default:
		}
		if schedule != nil {
			if wait := time.Until(schedule.Grant()); wait > 0 {
				timer.Reset(wait)
				select {
				case <-r.stop:
					return
				case <-timer.C:
				}
			}
		}
		start := time.Now()
		timeMillis, nano := nettest.SplitTime(start.UnixNano())
		for j := range batch {
			index, row := nettest.Choose(random.Next(), len(symbols))
			packed := nettest.PackTimeSequence(timeMillis, sequence)
			sequence = (sequence + 1) & nettest.SequenceMask
			if timeAndSale {
				event := sales[index]
				if !r.config.Reuse {
					copied := *event
					event = &copied
				}
				event.SetIndex(packed)
				event.SetTimeNanoPart(nano)
				event.SetPrice(table[row].Price)
				event.SetSize(table[row].Size)
				batch[j] = event
			} else {
				event := trades[index]
				if !r.config.Reuse {
					copied := *event
					event = &copied
				}
				event.SetTimeSequence(packed)
				event.SetTimeNanoPart(nano)
				event.SetPrice(table[row].Price)
				event.SetSize(table[row].Size)
				batch[j] = event
			}
		}
		if err := publisher.Publish(batch); err != nil {
			r.fail(fmt.Errorf("publish: %w", err))
			return
		}
		elapsed := time.Since(start).Nanoseconds()
		instance.records.Add(int64(len(batch)))
		instance.lagSum.Add(int64(len(batch)) * elapsed)
	}
}

// newPoolTrade returns the pool trade of a symbol: NaN in the price, the change, the size, the day volume and the day
// turnover, as the Trade of the other APIs.
func newPoolTrade(symbol string) *trade.Trade {
	t := trade.NewTrade(symbol)
	t.SetPrice(math.NaN())
	t.SetChange(math.NaN())
	t.SetSize(math.NaN())
	t.SetDayVolume(math.NaN())
	t.SetDayTurnover(math.NaN())
	return t
}

// connections returns the open connections of the connectors of all instances.
func (r *netTestRun) connections() int {
	total := 0
	for _, instance := range r.instances {
		connectors, err := instance.endpoint.Connectors()
		if err != nil {
			continue
		}
		for _, connector := range connectors {
			total += connector.ConnectionCount
		}
	}
	return total
}

// waitAndPrint prints a statistics block at the end of each period on the wall clock until the run stops; the first
// tick only starts the count. At the stop it prints the last block when it is at least one second long.
func (r *netTestRun) waitAndPrint() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	var deadline <-chan time.Time
	if r.config.Duration > 0 {
		deadline = time.After(r.config.Duration)
	}
	period := r.config.Period
	var last time.Time
	var lastRuntime nettest.RuntimeSample
	timer := time.NewTimer(time.Until(nextTick(time.Now(), period)))
	defer timer.Stop()
	for {
		select {
		case <-r.stop:
			if r.failure == nil {
				r.printBlock(&last, &lastRuntime, time.Now(), true)
			}
			return
		case <-signals:
			r.end()
		case <-deadline:
			r.end()
		case now := <-timer.C:
			r.printBlock(&last, &lastRuntime, now, false)
			timer.Reset(time.Until(nextTick(now, period)))
		}
	}
}

// nextTick returns the next multiple of the period since the Unix epoch after now.
func nextTick(now time.Time, period time.Duration) time.Time {
	p := period.Nanoseconds()
	return time.Unix(0, now.UnixNano()/p*p+p)
}

// printBlock takes the counts of the interval since the last tick and prints them with the runtime line in one write;
// the first tick only starts the count, and the last block is printed only when it is at least one second long.
func (r *netTestRun) printBlock(last *time.Time, lastRuntime *nettest.RuntimeSample, now time.Time, final bool) {
	statistics := make([]nettest.InstanceStatistics, len(r.instances))
	for i, instance := range r.instances {
		statistics[i] = nettest.InstanceStatistics{Records: instance.records.Swap(0),
			LagSum: instance.lagSum.Swap(0)}
	}
	runtime := netTestRuntimeSample()
	previous, previousRuntime := *last, *lastRuntime
	*last, *lastRuntime = now, runtime
	if previous.IsZero() || final && now.Sub(previous) < time.Second {
		return
	}
	interval := now.Sub(previous).Nanoseconds()
	text := nettest.StatisticsBlock(r.config.Mode, r.config.Name, interval, r.connections(), statistics) + "\n" +
		nettest.RuntimeLine(r.config.Name, interval, previousRuntime, runtime) + "\n"
	_, _ = os.Stdout.WriteString(text)
}

// netTestRuntimeSample reads the heap and the collections of the isolate and the CPU time of the process; a value
// that is not given is -1.
func netTestRuntimeSample() nettest.RuntimeSample {
	s, err := api.GetRuntimeStatistics()
	if err != nil {
		return nettest.RuntimeSample{HeapUsedBytes: -1, HeapCommittedBytes: -1, GCCount: -1, GCTimeMs: -1,
			CPUTimeNs: -1}
	}
	sum := func(a, b int64) int64 {
		if a < 0 || b < 0 {
			return -1
		}
		return a + b
	}
	return nettest.RuntimeSample{
		HeapUsedBytes:      s.HeapUsedBytes,
		HeapCommittedBytes: s.HeapCommittedBytes,
		GCCount:            sum(s.YoungGCCount, s.FullGCCount),
		GCTimeMs:           sum(s.YoungGCTimeMs, s.FullGCTimeMs),
		CPUTimeNs:          s.ProcessCPUTimeNs,
	}
}
