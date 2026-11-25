// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/nettest"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/sdk"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
)

// The limits of the consumer: the events checked for the nanosecond part, the times without data and how often the
// watchdog checks them.
const (
	nanoCheckEvents       = 1000
	noDataAfterRecord     = 30 * time.Second
	noDataAfterConnection = 60 * time.Second
	watchdogPeriod        = 200 * time.Millisecond
	nanoTimeProperty      = "dxscheme.nanoTime"
	netTestImplementation = "go-api"
)

type NetTest struct{}

func (c NetTest) ShortDescription() string {
	return "Publishes or receives generated events at a controlled rate and prints the rate and the lag."
}

// Run runs NetTest as the NetTest specification nettest/1 defines it: an invalid command line or none exits with 2,
// an error of the run with 1, the end of --duration or a signal with 0.
func (c NetTest) Run(args []string) int {
	config, err := nettest.Parse(args)
	if err == nil && config.Help {
		fmt.Print(nettest.HelpText)
		return 0
	}
	if err == nettest.ErrNoArguments {
		fmt.Fprint(os.Stderr, nettest.HelpText)
		return 2
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "nettest: %v\n\n%s", err, nettest.HelpText)
		return 2
	}
	return runNetTest(&config, args)
}

// netTestStatusError is an error of the run with its status in the summary file.
type netTestStatusError struct {
	status string
	err    error
}

func (e *netTestStatusError) Error() string { return e.status + ": " + e.err.Error() }

func (e *netTestStatusError) Unwrap() error { return e.err }

// netTestInstance is an instance of the tool: an endpoint, its counts of the current interval and its records since
// the start.
type netTestInstance struct {
	endpoint *api.DXEndpoint
	records  atomic.Int64
	lagSum   atomic.Int64
	total    atomic.Int64
}

// netTestSnapshot holds what the tool reads from the native runtime once per interval, for the metrics and the
// summary.
type netTestSnapshot struct {
	connections []int
	connectors  []nettest.SummaryConnector
	connector   nettest.ConnectorMetrics
	runtime     nettest.RuntimeSample
	// taken tells whether a tick took the snapshot; all are the instances that the scrapes read.
	taken bool
	all   []*netTestInstance
}

// netTestRun is a run of the tool.
type netTestRun struct {
	config   *nettest.Config
	nanoTime bool
	// start is the start of the run and the base of now.
	start time.Time
	// instances are the endpoints, built after the checks of the start.
	instances []*netTestInstance
	// stop ends the run; failure holds the error that ended it, if any; signaled tells that a signal ended it.
	stop     chan struct{}
	stopOnce sync.Once
	failure  error
	signaled atomic.Bool
	// workers are the goroutines that use the endpoints: the publishers and the watchdog of the consumer.
	workers sync.WaitGroup
	// The consumer: the events received, those with a nanosecond part, the time of the last record and the lags.
	received     atomic.Int64
	withNano     atomic.Int64
	lastRecordNs atomic.Int64
	lag          *nettest.LagRecorder
	// watchdog stops the consumer without data; connections returns the open connections of the instances.
	watchdog    netTestWatchdog
	connections func() int
	// The publisher: the fills and the publishing calls, and the schedule slip of the last batch and of the interval.
	schedule      *nettest.Schedule
	publishCalls  atomic.Int64
	publishEvents atomic.Int64
	publishNs     atomic.Int64
	slipLast      atomic.Int64
	slipMax       atomic.Int64
	// warmupStartNs is the start of the warm-up: the first record on c, the first grant on p; 0 before it.
	warmupStartNs atomic.Int64
	// The intervals: the start of the current one, the last tick and the runtime at it.
	counting      bool
	countStart    time.Time
	intervalStart time.Time
	lastTick      time.Time
	lastRuntime   nettest.RuntimeSample
	// startRuntime is the runtime once the instances have started, and startRuntimeAt its time: the start of the
	// growths of the runtime in the summary, as the first sample of the SDK C tool.
	startRuntime   nettest.RuntimeSample
	startRuntimeAt time.Time
	summary        *nettest.SummaryBuilder
	// snapshot is read by the scrapes of the metrics.
	snapshotLock sync.Mutex
	snapshot     netTestSnapshot
	metrics      *nettest.MetricsServer
}

func (r *netTestRun) fail(err error) {
	r.stopOnce.Do(func() {
		r.failure = err
		close(r.stop)
	})
}

func (r *netTestRun) end() { r.fail(nil) }

// stopped tells whether the run has been stopped.
func (r *netTestRun) stopped() bool {
	select {
	case <-r.stop:
		return true
	default:
		return false
	}
}

// watchSignals ends the run at the first SIGINT or SIGTERM, as the end of --duration does, and the process at once at
// the second one, without the summary and with 128 plus the signal number: the tool may wait for something that does
// not check the stop, such as the reading of a large symbol file.
func (r *netTestRun) watchSignals(signals <-chan os.Signal, done <-chan struct{}) {
	select {
	case <-signals:
	case <-done:
		return
	}
	r.signaled.Store(true)
	r.end()
	select {
	case s := <-signals:
		code := 1
		if number, ok := s.(syscall.Signal); ok {
			code = 128 + int(number)
		}
		os.Exit(code)
	case <-done:
	}
}

// now returns the wall-clock time of the events, see wallNanos.
func (r *netTestRun) now() time.Time { return time.Unix(0, wallNanos()) }

// runNetTest runs the tool and returns its exit code: it stops the workers, writes the summary file and only then
// closes the endpoints.
func runNetTest(config *nettest.Config, args []string) int {
	nanoTime := "true"
	if value, ok := config.Property(nanoTimeProperty); ok {
		nanoTime = value
	}
	// Any case of true is on, as Boolean.parseBoolean of Java reads the property.
	r := &netTestRun{config: config, nanoTime: strings.EqualFold(nanoTime, "true"), start: time.Now(),
		stop: make(chan struct{}), summary: nettest.NewSummaryBuilder(config.Mode),
		watchdog: netTestWatchdog{afterRecord: noDataAfterRecord, afterConnection: noDataAfterConnection}}
	r.connections = r.openConnections
	// The signals are taken from here on: a signal while the symbols are read ends the run with the summary too.
	signals, done := make(chan os.Signal, 2), make(chan struct{})
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go r.watchSignals(signals, done)
	defer func() {
		signal.Stop(signals)
		close(done)
	}()
	if config.Mode == nettest.Consumer {
		r.lag = nettest.NewLagRecorder()
	}
	err := r.run(nanoTime)
	// Stop the goroutines that use the endpoints before closing them.
	r.end()
	r.workers.Wait()
	code := 0
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during nettest: %v\n", err)
		code = 1
	}
	if config.Summary != "" {
		if err := nettest.WriteSummary(config.Summary, r.buildSummary(args, err, code)); err != nil {
			fmt.Fprintf(os.Stderr, "Error during nettest: the summary file %s: %v\n", config.Summary, err)
			code = 1
		}
	}
	r.close()
	return code
}

func (r *netTestRun) run(nanoTime string) error {
	config := r.config
	symbols, err := netTestSymbols(config)
	if err != nil {
		return err
	}
	// A signal before the start of the endpoints: the summary of the empty run, without the endpoints.
	if r.stopped() {
		return nil
	}
	if config.MetricsAddress != "" {
		r.metrics, err = nettest.ListenMetrics(config.MetricsAddress)
		if err != nil {
			return fmt.Errorf("the metrics address %s cannot be listened on: %w", config.MetricsAddress, err)
		}
		r.metrics.Serve(r.scrape)
		fmt.Printf("{*NetTestMetrics-%s*} %s\n", config.Name, r.metrics.URL())
	}
	for _, p := range config.Properties {
		if err := api.SetSystemProperty(p.Key, p.Value); err != nil {
			return fmt.Errorf("SetSystemProperty %s: %w", p.Key, err)
		}
	}
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
	r.counting, r.intervalStart = true, time.Now()
	r.countStart = r.intervalStart
	r.snapshotLock.Lock()
	r.snapshot.all = r.instances
	r.snapshotLock.Unlock()
	if config.Mode == nettest.Consumer {
		err = r.startConsumers(symbols)
	} else {
		err = r.startProducers(symbols)
	}
	if err != nil {
		return err
	}
	r.startRuntime, r.startRuntimeAt = netTestRuntimeSample(), time.Now()
	r.waitAndPrint()
	return r.failure
}

// close stops the metrics and closes the endpoints.
func (r *netTestRun) close() {
	if r.metrics != nil {
		r.metrics.Close()
	}
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
		if err := r.connect(instance); err != nil {
			return err
		}
	}
	r.workers.Add(1)
	go r.watchConsumer()
	return nil
}

// netTestListener counts the events of a consumer instance and records their e2e lags.
type netTestListener struct {
	run      *netTestRun
	instance *netTestInstance
}

func (l *netTestListener) Update(events []any) {
	r := l.run
	now := wallNanos()
	r.warmupStartNs.CompareAndSwap(0, now)
	var lagSum, withNano int64
	lags := make([]int64, 0, len(events))
	symbols := make([]string, 0, len(events))
	for _, event := range events {
		var timeMillis int64
		var nano int32
		var symbol *string
		switch e := event.(type) {
		case *trade.Trade:
			timeMillis, nano, symbol = e.Time(), e.TimeNanoPart(), e.EventSymbol()
		case *timeandsale.TimeAndSale:
			timeMillis, nano, symbol = e.Time(), e.TimeNanoPart(), e.EventSymbol()
		}
		if nano != 0 {
			withNano++
		}
		lag := now - (timeMillis*1000000 + int64(nano))
		lagSum += lag
		lags = append(lags, lag)
		if symbol != nil {
			symbols = append(symbols, *symbol)
		}
	}
	r.lag.Record(lags, symbols)
	l.instance.records.Add(int64(len(events)))
	l.instance.lagSum.Add(lagSum)
	l.instance.total.Add(int64(len(events)))
	r.lastRecordNs.Store(now)
	r.withNano.Add(withNano)
	if received := r.received.Add(int64(len(events))); r.nanoTime && received >= nanoCheckEvents &&
		received-int64(len(events)) < nanoCheckEvents && r.withNano.Load() == 0 {
		r.fail(&netTestStatusError{nettest.StatusNanoMissing, fmt.Errorf("the first %d events have no nanosecond "+
			"part: a hop of the path lacks %s=true (or run with -D%s=false)", received, nanoTimeProperty,
			nanoTimeProperty)})
	}
}

// watchConsumer stops the consumer without data, as the watchdog tells, until the run stops.
func (r *netTestRun) watchConsumer() {
	defer r.workers.Done()
	ticker := time.NewTicker(watchdogPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-r.stop:
			return
		case <-ticker.C:
			var lastRecord time.Time
			if last := r.lastRecordNs.Load(); last != 0 {
				lastRecord = time.Unix(0, last)
			}
			if err := r.watchdog.check(r.now(), r.connections(), lastRecord); err != nil {
				r.fail(&netTestStatusError{nettest.StatusNoData, err})
			}
		}
	}
}

// netTestWatchdog tells when the consumer has no data, as the specification says: with open connections, no records
// for 30 s after the last record, or for 60 s after the connection when no record came. The open connections are
// those of the connector statistics, as in the statistics lines: a listening endpoint in the state CONNECTED without
// peers has none. After the connections were all closed, the time counts from the next connection again.
type netTestWatchdog struct {
	afterRecord, afterConnection time.Duration
	// connectedSince is the time since which connections are open, zero while there are none.
	connectedSince time.Time
}

// check returns the error of no data at the time, or nil; lastRecord is zero before the first record.
func (w *netTestWatchdog) check(now time.Time, connections int, lastRecord time.Time) error {
	if connections <= 0 {
		w.connectedSince = time.Time{}
		return nil
	}
	if w.connectedSince.IsZero() {
		w.connectedSince = now
	}
	if lastRecord.IsZero() {
		if now.Sub(w.connectedSince) >= w.afterConnection {
			return fmt.Errorf("no records for %d ms after the connection, with %d open connections",
				w.afterConnection.Milliseconds(), connections)
		}
		return nil
	}
	since := lastRecord
	if since.Before(w.connectedSince) {
		since = w.connectedSince
	}
	if now.Sub(since) >= w.afterRecord {
		return fmt.Errorf("no records for %d ms after the last record, with %d open connections",
			w.afterRecord.Milliseconds(), connections)
	}
	return nil
}

// openConnections returns the open connections of the connectors of all instances.
func (r *netTestRun) openConnections() int {
	count := 0
	for _, instance := range r.instances {
		if connectors, err := instance.endpoint.Connectors(); err == nil {
			for _, connector := range connectors {
				count += connector.ConnectionCount
			}
		}
	}
	return count
}

// serverStartTimeout is how long the start waits for a server socket to bind, as in the NetTest tool of the SDK.
const serverStartTimeout = 3 * time.Second

// connect connects an instance to the address. A listen address may carry connector options, such as
// ":7400[bindAddr=127.0.0.1]", and the SDK binds its port in the background and only logs a failure: the connector of a
// bound server socket reports itself as connected within milliseconds, one that cannot bind (the port is taken) stays
// connecting. So for a listen address it waits until every server socket connector is connected, at most
// serverStartTimeout. A client address never waits: with no peer its connector stays connecting, which is normal.
func (r *netTestRun) connect(instance *netTestInstance) error {
	address := r.config.Address
	if err := instance.endpoint.Connect(address); err != nil {
		return fmt.Errorf("connect to %s: %w", address, err)
	}
	if !r.config.ServerAddress() {
		return nil
	}
	deadline := time.Now().Add(serverStartTimeout)
	for {
		connectors, err := instance.endpoint.Connectors()
		if err != nil {
			return fmt.Errorf("the connectors of %s: %w", address, err)
		}
		pending := false
		for _, connector := range connectors {
			if connector.Kind == api.ServerSocketConnector && connector.State != api.ConnectorConnected {
				pending = true
			}
		}
		if !pending {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("the address %s cannot be listened on: the server socket did not start in %v "+
				"(is the port already in use?)", address, serverStartTimeout)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// startProducers connects each instance and starts its publishing loop.
func (r *netTestRun) startProducers(symbols []string) error {
	if r.config.Rate > 0 {
		r.schedule = nettest.NewSchedule(r.config.Rate, r.config.Batch)
	}
	for i, instance := range r.instances {
		if err := r.connect(instance); err != nil {
			return err
		}
		publisher, err := instance.endpoint.GetPublisher()
		if err != nil {
			return err
		}
		r.workers.Add(1)
		go r.publish(i+1, instance, publisher, symbols)
	}
	return nil
}

// publish publishes the batches of an instance until the run stops: wait for the schedule, read the clock, fill the
// batch from the pool and publish it.
func (r *netTestRun) publish(number int, instance *netTestInstance, publisher *api.DXPublisher, symbols []string) {
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
		var due time.Time
		if r.schedule != nil {
			due = r.schedule.Grant()
			r.warmupStartNs.CompareAndSwap(0, due.UnixNano())
			if wait := time.Until(due); wait > 0 {
				timer.Reset(wait)
				select {
				case <-r.stop:
					return
				case <-timer.C:
				}
			}
		}
		// The monotonic clock for the slip and the elapsed time, the wall clock for the events.
		start, wall := time.Now(), wallNanos()
		if r.schedule != nil {
			slip := max(start.Sub(due).Nanoseconds(), 0)
			r.slipLast.Store(slip)
			for current := r.slipMax.Load(); slip > current && !r.slipMax.CompareAndSwap(current, slip); {
				current = r.slipMax.Load()
			}
		} else {
			r.warmupStartNs.CompareAndSwap(0, wall)
		}
		timeMillis, nano := nettest.SplitTime(wall)
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
		records := int64(len(batch))
		instance.records.Add(records)
		instance.lagSum.Add(records * elapsed)
		instance.total.Add(records)
		r.publishCalls.Add(1)
		r.publishEvents.Add(records)
		r.publishNs.Add(elapsed)
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

// readConnectors reads the connectors of all instances: the open connections of each instance, each connector with
// its statistics v1 and v2, and their sums; without v2 on any instance the sums have no v2.
func (r *netTestRun) readConnectors() netTestSnapshot {
	s := netTestSnapshot{connections: make([]int, len(r.instances)), connectors: []nettest.SummaryConnector{}}
	v2 := &nettest.ConnectorMetricsV2{}
	for i, instance := range r.instances {
		connectors, err := instance.endpoint.Connectors()
		if err != nil {
			continue
		}
		io, ioErr := instance.endpoint.ConnectorIOStatistics()
		if ioErr != nil || len(io) != len(connectors) {
			v2 = nil
		}
		for j, connector := range connectors {
			s.connections[i] += connector.ConnectionCount
			st := connector.Statistics
			entry := nettest.SummaryConnector{Address: connector.Address, V1: &nettest.SummaryConnectorV1{
				ReadBytes: st.ReadBytes, WrittenBytes: st.WrittenBytes, ReadDataRecords: st.ReadDataRecords,
				WrittenDataRecords: st.WrittenDataRecords, ActiveConnections: st.ActiveConnectionCount,
				ClosedConnections: st.ClosedConnectionCount}}
			s.connector.ReadBytes += st.ReadBytes
			s.connector.WrittenBytes += st.WrittenBytes
			s.connector.ReadRecords += st.ReadDataRecords
			s.connector.WrittenRecords += st.WrittenDataRecords
			if ioErr == nil && j < len(io) {
				c := io[j]
				entry.V2 = &nettest.SummaryConnectorV2{ReadSubRecords: c.ReadSubRecords,
					WrittenSubRecords: c.WrittenSubRecords, DataReadLagUs: c.ReadDataLagSumUs,
					DataWriteLagUs: c.WrittenDataLagSumUs, ReadRTTUs: c.ReadRTTSumUs, WriteRTTUs: c.WrittenRTTSumUs}
				if v2 != nil {
					v2.ReadLagUs += c.ReadDataLagSumUs
					v2.RTTUs += c.ReadRTTSumUs + c.WrittenRTTSumUs
					v2.SubRecords += c.ReadSubRecords + c.WrittenSubRecords
				}
			}
			s.connectors = append(s.connectors, entry)
		}
	}
	if len(r.instances) > 0 {
		s.connector.V2 = v2
	}
	return s
}

// scrape returns the metrics from the counters and the snapshot of the last tick, without native calls.
func (r *netTestRun) scrape() *nettest.MetricsSnapshot {
	config := r.config
	side := "producer"
	if config.Mode == nettest.Consumer {
		side = "consumer"
	}
	s := &nettest.MetricsSnapshot{
		Info: nettest.MetricsInfo{Tool: netTestImplementation, Side: side, Name: config.Name,
			Version: moduleVersion(), SDKVersion: sdk.Version, Role: config.Role, Event: config.Event,
			Collector: config.Collector},
		Records:   make([]int64, config.Instances),
		PublishNs: r.publishNs.Load(), PublishCalls: r.publishCalls.Load(), SlipNs: r.slipLast.Load(),
		Scheduled: config.Mode == nettest.Producer && config.Rate > 0,
	}
	if r.lag != nil {
		lag := r.lag.Metrics()
		s.Lag = &lag
	}
	r.snapshotLock.Lock()
	snapshot := r.snapshot
	r.snapshotLock.Unlock()
	s.Connections = make([]int, config.Instances)
	copy(s.Connections, snapshot.connections)
	for i, instance := range snapshot.all {
		s.Records[i] = instance.total.Load()
	}
	s.Connector = snapshot.connector
	s.Runtime = snapshot.runtime
	if !snapshot.taken {
		s.Runtime = nettest.UnknownRuntimeSample()
	}
	return s
}

// waitAndPrint ends an interval at each multiple of the period on the wall clock until the run stops, and the last
// one at the stop.
func (r *netTestRun) waitAndPrint() {
	var deadline <-chan time.Time
	if r.config.Duration > 0 {
		deadline = time.After(r.config.Duration)
	}
	period := r.config.Period
	timer := time.NewTimer(time.Until(nextTick(time.Now(), period)))
	defer timer.Stop()
	for {
		select {
		case <-r.stop:
			// The publishers end their batch first, so that the last interval has all of them.
			r.workers.Wait()
			r.tick(time.Now(), true)
			return
		case <-deadline:
			r.end()
		case now := <-timer.C:
			r.tick(now, false)
			timer.Reset(time.Until(nextTick(now, period)))
		}
	}
}

// nextTick returns the next multiple of the period since the Unix epoch after now.
func nextTick(now time.Time, period time.Duration) time.Time {
	p := period.Nanoseconds()
	return time.Unix(0, now.UnixNano()/p*p+p)
}

// warmup tells whether an interval that starts at the time belongs to the warm-up: before the warm-up starts, or when
// it starts before the warm-up ends.
func (r *netTestRun) warmup(start time.Time) bool {
	if r.config.Warmup <= 0 {
		return false
	}
	warmupStart := r.warmupStartNs.Load()
	return warmupStart == 0 || start.UnixNano() < warmupStart+r.config.Warmup.Nanoseconds()
}

// tick ends an interval: it takes the counts, reads the connectors and the runtime once for the metrics, adds the
// interval to the summary and prints the statistics block with the runtime line in one write. The first tick prints
// nothing; the last block, since the last tick or since the start of the instances when no tick came, is printed only
// when it is at least one second long and the run did not fail.
func (r *netTestRun) tick(now time.Time, final bool) {
	statistics := make([]nettest.InstanceStatistics, len(r.instances))
	counts := nettest.IntervalCounts{Start: r.intervalStart, End: now, Warmup: r.warmup(r.intervalStart)}
	for i, instance := range r.instances {
		statistics[i] = nettest.InstanceStatistics{Records: instance.records.Swap(0),
			LagSum: instance.lagSum.Swap(0)}
		counts.Records += statistics[i].Records
		counts.PublishNs += statistics[i].LagSum
	}
	snapshot := r.readConnectors()
	snapshot.runtime = netTestRuntimeSample()
	snapshot.taken, snapshot.all = true, r.instances
	r.snapshotLock.Lock()
	r.snapshot = snapshot
	r.snapshotLock.Unlock()
	for _, connections := range snapshot.connections {
		counts.Connections += connections
	}
	if r.lag != nil {
		lag := r.lag.TakeInterval()
		counts.Lag = &lag
	}
	if r.schedule != nil {
		counts.Slip, counts.SlipMax, counts.SlipEnd = true, r.slipMax.Swap(0), r.slipLast.Load()
	}
	if r.counting && (now.After(r.intervalStart) || counts.Records > 0) {
		r.summary.Add(counts)
	}
	r.intervalStart = now
	previous, previousRuntime := r.lastTick, r.lastRuntime
	r.lastTick, r.lastRuntime = now, snapshot.runtime
	if previous.IsZero() && final {
		// A run shorter than the period: the last block counts from the start of the instances.
		previous, previousRuntime = r.startRuntimeAt, r.startRuntime
	}
	if previous.IsZero() || final && (now.Sub(previous) < time.Second || r.failure != nil) {
		return
	}
	interval := now.Sub(previous).Nanoseconds()
	text := nettest.StatisticsBlock(r.config.Mode, r.config.Name, interval, counts.Connections, statistics) + "\n" +
		nettest.RuntimeLine(r.config.Name, interval, previousRuntime, snapshot.runtime) + "\n"
	_, _ = os.Stdout.WriteString(text)
}

// buildSummary returns the summary of the run that ended with the error and the exit code.
func (r *netTestRun) buildSummary(args []string, err error, code int) *nettest.Summary {
	config := r.config
	end := time.Now()
	status := nettest.StatusOK
	var statusErr *netTestStatusError
	switch {
	case errors.As(err, &statusErr):
		status = statusErr.status
	case err != nil:
		status = nettest.StatusError
	case r.signaled.Load():
		status = nettest.StatusStoppedBySignal
	}
	var errorText *string
	if err != nil {
		errorText = nettest.String(err.Error())
	}
	side := "producer"
	if config.Mode == nettest.Consumer {
		side = "consumer"
	}
	resolution := "ms"
	if config.Mode == nettest.Producer && r.nanoTime || config.Mode == nettest.Consumer && r.withNano.Load() > 0 {
		resolution = "ns"
	}
	s := &nettest.Summary{
		Schema: "nettest-summary/1", SpecVersion: nettest.SpecVersion,
		Tool: nettest.SummaryTool{Implementation: netTestImplementation, Name: "tools nettest",
			Version: nettest.String(moduleVersion()), SDKVersion: nettest.String(sdk.Version)},
		Side: side, Name: config.Name, Args: append([]string{}, args...), Config: config.SummaryConfig(r.nanoTime),
		StartEpochNs: r.start.UnixNano(), EndEpochNs: end.UnixNano(), Status: status, ExitCode: code, Error: errorText,
		Records: r.summary.Records, Rate: nettest.Rate(r.summary.Records, end.Sub(r.countStart).Nanoseconds()),
		TimeSource: "actual", TimeResolution: resolution, Intervals: r.summary.Intervals,
		Connectors: r.snapshot.connectors,
	}
	s.Tool.QDVersion, s.Tool.Build = netTestBuildInfo()
	if s.Connectors == nil {
		s.Connectors = []nettest.SummaryConnector{}
	}
	if config.Mode == nettest.Consumer {
		s.E2ELag = r.summary.Lag()
		s.SymbolsSeen = nettest.Int(r.lag.SymbolsSeen())
	} else {
		calls, events, ns := r.publishCalls.Load(), r.publishEvents.Load(), r.publishNs.Load()
		s.PublishCall = &nettest.SummaryPublish{Calls: calls, Events: events, NsTotal: ns,
			NsPerEventMean: nettest.IntIf(events > 0, ns/max(events, 1))}
		s.ScheduleSlip = r.summary.Slip()
	}
	if !r.lastTick.IsZero() && !r.startRuntimeAt.IsZero() {
		s.Runtime = nettest.NewSummaryRuntime(r.startRuntime, r.snapshot.runtime,
			r.lastTick.Sub(r.startRuntimeAt).Nanoseconds())
	}
	return s
}

// netTestBuildInfo reads the QDS version and the build of the library for the summary; a value the library does not
// give is null.
func netTestBuildInfo() (*string, nettest.SummaryBuild) {
	info, err := api.GetBuildInfo()
	if err != nil {
		return nil, nettest.SummaryBuild{}
	}
	value := func(v string) *string {
		if v == "" {
			return nil
		}
		return nettest.String(v)
	}
	return value(info.QDVersion), nettest.SummaryBuild{Vendor: value(info.Vendor), March: value(info.March),
		GC: value(info.GC)}
}

// netTestRuntimeSample reads the heap and the collections of the isolate and the CPU time of the process; a value
// that is not given is -1.
func netTestRuntimeSample() nettest.RuntimeSample {
	s, err := api.GetRuntimeStatistics()
	if err != nil {
		return nettest.UnknownRuntimeSample()
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
		GCCounts:           [2]int64{s.YoungGCCount, s.FullGCCount},
		GCTimesMs:          [2]int64{s.YoungGCTimeMs, s.FullGCTimeMs},
		CPUTimeNs:          s.ProcessCPUTimeNs,
	}
}
