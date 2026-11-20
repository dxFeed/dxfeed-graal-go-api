// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/diagnose"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/sdk"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/toolargs"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

const (
	// diagnoseDefaultDuration and diagnoseMaxDuration are the default and the greatest time to wait for the events
	// of an address, in seconds.
	diagnoseDefaultDuration = 10
	diagnoseMaxDuration     = 3600
	// diagnoseFirstEvents is the number of the first events of a type that the report shows.
	diagnoseFirstEvents = 5
	// diagnoseOnlyFromLimit is the number of the pairs received from one address only that a comparison lists.
	diagnoseOnlyFromLimit = 20
	// diagnoseTimeFormat is the format of the times of the report: local time with milliseconds.
	diagnoseTimeFormat = "20060102-150405.000"
)

type Diagnose struct{}

func (c Diagnose) ShortDescription() string {
	return "Diagnoses the connection to the address: the network, the endpoint and the events of the symbols."
}

func (c Diagnose) Run(args []string) int {
	return runTool("diagnose", args, diagnoseHelp, []toolargs.Option{toolargs.Compare, toolargs.Duration,
		toolargs.Report, toolargs.Properties}, func(a toolargs.Args) error {
		duration := diagnoseDefaultDuration
		if value := a.Value(toolargs.Duration); value != nil {
			var err error
			if duration, err = parseDiagnoseDuration(*value); err != nil {
				return usageError{err}
			}
		}
		common, err := parseCommon(a, "", "")
		if err != nil {
			return err
		}
		_, typesArg, symbolsArg, _ := a.Common("", "")
		d := &diagnosis{
			printer:    &reportPrinter{},
			common:     common,
			duration:   time.Duration(duration) * time.Second,
			properties: a.Value(toolargs.Properties),
		}
		return d.run(typesArg, symbolsArg, a.Value(toolargs.Compare), a.Value(toolargs.Report))
	})
}

// parseDiagnoseDuration parses the --duration value: a whole number of seconds from 1 to diagnoseMaxDuration.
func parseDiagnoseDuration(value string) (int, error) {
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds < 1 || seconds > diagnoseMaxDuration {
		return 0, fmt.Errorf("the duration must be a whole number of seconds from 1 to %d: %s", diagnoseMaxDuration,
			value)
	}
	return seconds, nil
}

// reportPrinter prints the lines of the report to stdout and keeps them for the report file, with the secrets
// masked. Listeners print from SDK threads, so it is safe for concurrent use.
type reportPrinter struct {
	mu     sync.Mutex
	report strings.Builder
}

func (p *reportPrinter) printf(format string, args ...any) {
	line := diagnose.Redact(fmt.Sprintf(format, args...)) + "\n"
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Print(line)
	p.report.WriteString(line)
}

type diagnosis struct {
	printer    *reportPrinter
	common     commonArgs
	duration   time.Duration
	properties *string
}

// addressResult is what the diagnosis of an address found.
type addressResult struct {
	address      string
	verdict      diagnose.Verdict
	wasConnected bool
	total        int
	counts       map[string]int
	received     map[diagnose.Pair]bool
}

func (d *diagnosis) run(typesArg, symbolsArg string, compare, report *string) error {
	p := d.printer
	if err := setToolProperties(); err != nil {
		return err
	}
	for key, value := range d.common.properties {
		if err := api.SetSystemProperty(key, value); err != nil {
			return fmt.Errorf("SetSystemProperty %s: %w", key, err)
		}
	}
	p.printf("Diagnose %s %s %s", d.common.address, typesArg, symbolsArg)
	p.printf("Duration: %d s per address", int(d.duration.Seconds()))
	if d.properties != nil {
		p.printf("Properties: %s", *d.properties)
	}
	p.printf("dxFeed Graal Go API: %s", moduleVersion())
	p.printf("Graal Native SDK: %s", sdk.Version)
	p.printf("Platform: %s/%s, %s, %d processors", runtime.GOOS, runtime.GOARCH, runtime.Version(),
		runtime.NumCPU())

	results := []*addressResult{}
	for _, address := range append([]string{d.common.address}, valueList(compare)...) {
		result, err := d.diagnoseAddress(address)
		if err != nil {
			return err
		}
		results = append(results, result)
	}
	if len(results) == 2 {
		d.printComparison(results[0], results[1])
	}

	worst := diagnose.Data
	for _, result := range results {
		worst = max(worst, result.verdict.Kind)
	}
	if report != nil {
		if err := os.WriteFile(*report, []byte(p.report.String()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "The report is not written to %s: %v\n", *report, err)
			return exitStatus(1)
		}
		fmt.Printf("The report is written to %s\n", *report)
	}
	if code := worst.ExitCode(); code != 0 {
		return exitStatus(code)
	}
	return nil
}

func valueList(value *string) []string {
	if value == nil {
		return nil
	}
	return []string{*value}
}

// typeStats are the statistics of the events of a type.
type typeStats struct {
	count      int
	firstAfter time.Duration
	first      []string
	// flags counts the events with SNAPSHOT_BEGIN, SNAPSHOT_END, TX_PENDING and REMOVE_EVENT.
	flags                  [4]int
	lagMin, lagMax, lagSum int64
	lagCount               int
	lastBySymbol           map[string]time.Time
	countBySymbol          map[string]int
	connectTime            time.Time
	receivedAny            bool
	// mu is shared by the statistics of all the types of an address.
	mu *sync.Mutex
}

// typeListener collects the statistics of the events of a type.
type typeListener struct{ stats *typeStats }

func (l typeListener) Update(eventsList []any) {
	s := l.stats
	receive := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(eventsList) > 0 && !s.receivedAny {
		s.receivedAny = true
		s.firstAfter = receive.Sub(s.connectTime)
	}
	for _, event := range eventsList {
		s.count++
		if len(s.first) < diagnoseFirstEvents {
			s.first = append(s.first, eventText(event))
		}
		if indexed, ok := event.(events.IndexedEvent); ok {
			for i, set := range []bool{events.IsSnapshotBegin(indexed), events.IsSnapshotEnd(indexed),
				events.IsPending(indexed), events.IsRemove(indexed)} {
				if set {
					s.flags[i]++
				}
			}
		}
		if eventTime := eventTimeMillis(event); eventTime > 0 {
			lag := receive.UnixMilli() - eventTime
			if s.lagCount == 0 || lag < s.lagMin {
				s.lagMin = lag
			}
			if s.lagCount == 0 || lag > s.lagMax {
				s.lagMax = lag
			}
			s.lagSum += lag
			s.lagCount++
		}
		symbol := eventSymbol(event)
		s.countBySymbol[symbol]++
		s.lastBySymbol[symbol] = receive
	}
}

func eventText(event any) string {
	if s, ok := event.(events.StringConverter); ok {
		return s.String()
	}
	return fmt.Sprint(event)
}

// eventTimeMillis returns the time of an event: its event time, else the time of a quote, a trade or a time series
// event; 0 when it has none.
func eventTimeMillis(event any) int64 {
	if e, ok := event.(interface{ EventTime() int64 }); ok && e.EventTime() != 0 {
		return e.EventTime()
	}
	if e, ok := event.(interface{ Time() int64 }); ok {
		return e.Time()
	}
	return 0
}

func eventSymbol(event any) string {
	switch e := event.(type) {
	case interface{ EventSymbol() *string }:
		if symbol := e.EventSymbol(); symbol != nil {
			return *symbol
		}
	case interface{ EventSymbol() *candle.CandleSymbol }:
		if symbol := e.EventSymbol(); symbol != nil {
			return symbol.String()
		}
	}
	return "?"
}

// stateListener prints the state changes of the endpoint.
type stateListener struct {
	printer   *reportPrinter
	mu        sync.Mutex
	connected bool
}

func (l *stateListener) UpdateState(old, new common.ConnectionState) {
	if new == common.Connected {
		l.mu.Lock()
		l.connected = true
		l.mu.Unlock()
	}
	l.printer.printf("  %s %s -> %s", time.Now().Format(diagnoseTimeFormat), old, new)
}

func (l *stateListener) wasConnected() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.connected
}

// connectorListener prints each new connection error of a connector.
type connectorListener struct {
	printer *reportPrinter
	mu      sync.Mutex
	counts  map[string]int64
}

func (l *connectorListener) ConnectorChanged(connector api.MessageConnectorInfo) {
	if connector.LastError == nil {
		return
	}
	key := connector.Name + " " + connector.Address
	l.mu.Lock()
	isNew := connector.LastError.Count > l.counts[key]
	if isNew {
		l.counts[key] = connector.LastError.Count
	}
	l.mu.Unlock()
	if isNew {
		l.printer.printf("  %s connection error: %s", time.Now().Format(diagnoseTimeFormat),
			connectionErrorText(connector.LastError))
	}
}

// connectionErrorText returns the text of a connection error without the prefix of the Go errors.
func connectionErrorText(err *api.ConnectionError) string {
	return strings.TrimPrefix(err.Error(), "dxfeed: ")
}

func (d *diagnosis) diagnoseAddress(address string) (*addressResult, error) {
	p := d.printer
	p.printf("")
	p.printf("=== Address: %s", address)
	noHostResolved, anyTCPConnected := d.probe(address)

	p.printf("Endpoint:")
	endpoint, err := api.NewEndpointBuilder().WithRole(api.Feed).WithProperties(d.common.properties).
		WithName("DiagnoseTool").Build()
	if err != nil {
		return nil, fmt.Errorf("NewEndpoint: %w", err)
	}
	defer func() { _ = endpoint.Close() }()
	states := &stateListener{printer: p}
	endpoint.AddListener(states)
	if err := endpoint.AddMessageConnectorListener(&connectorListener{printer: p,
		counts: map[string]int64{}}); err != nil {
		p.printf("  the connection errors are not printed: %v", err)
	}
	feed, err := endpoint.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("GetFeed: %w", err)
	}
	types := sortedTypes(d.common.types)
	mu := &sync.Mutex{}
	stats := map[string]*typeStats{}
	var subscriptions []*api.DXFeedSubscription
	defer func() {
		for _, subscription := range subscriptions {
			_ = subscription.Close()
		}
	}()
	for _, eventType := range types {
		s := &typeStats{mu: mu, lastBySymbol: map[string]time.Time{}, countBySymbol: map[string]int{}}
		stats[eventType.String()] = s
		subscription, err := feed.CreateSubscription(eventType)
		if err != nil {
			return nil, fmt.Errorf("CreateSubscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
		if err := subscription.AddListener(typeListener{stats: s}); err != nil {
			return nil, fmt.Errorf("AddListener: %w", err)
		}
		for _, symbol := range d.common.symbols {
			if err := subscription.AddSymbol(symbol); err != nil {
				return nil, fmt.Errorf("AddSymbol: %w", err)
			}
		}
	}

	connectTime := time.Now()
	mu.Lock()
	for _, s := range stats {
		s.connectTime = connectTime
	}
	mu.Unlock()
	var notAccepted error
	if notAccepted = endpoint.Connect(address); notAccepted != nil {
		p.printf("  the address is not accepted: %v", notAccepted)
	} else {
		time.Sleep(d.duration)
	}

	var lastError *api.ConnectionError
	if connectors, err := endpoint.Connectors(); err != nil {
		p.printf("  the connectors are not available: %v", err)
	} else {
		for _, connector := range connectors {
			p.printf("  connector: %s", connectorText(connector))
		}
		if lastError, err = endpoint.LastConnectionError(); err != nil {
			p.printf("  the connectors are not available: %v", err)
		}
	}
	if err := endpoint.CloseAndAwaitTermination(context.Background()); err != nil {
		p.printf("  close failed: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	result := &addressResult{address: address, wasConnected: states.wasConnected(), counts: map[string]int{},
		received: map[diagnose.Pair]bool{}}
	typeNames := make([]string, len(types))
	for i, eventType := range types {
		typeNames[i] = eventType.String()
	}
	d.printEvents(typeNames, stats, result)

	in := diagnose.Input{Types: typeNames, Symbols: symbolNames(d.common.symbols), Received: result.received,
		WasConnected: result.wasConnected, NoHostResolved: noHostResolved}
	if lastError != nil {
		in.HasLastError, in.LastError = true, connectionErrorText(lastError)
		in.LastErrorIsDNS = lastError.Category == api.ConnectionErrorDNS
	}
	result.verdict = diagnose.ComputeVerdict(in)
	if notAccepted != nil && result.verdict.Kind == diagnose.CannotConnect {
		result.verdict.Hint = "The endpoint did not accept the address: " + notAccepted.Error()
	}
	p.printf("Verdict: %s", result.verdict.Kind)
	if result.verdict.Hint != "" {
		p.printf("Hint: %s", result.verdict.Hint)
	}
	for _, missing := range result.verdict.Missing {
		p.printf("  no events: %s", missing)
	}
	if anyTCPConnected && !result.wasConnected {
		p.printf("Note: the TCP probe connected but the endpoint did not: often TLS, a login or a proxy.")
	}
	return result, nil
}

// probe prints the network probe of the address. It returns whether there were hosts and none of them resolved, and
// whether a TCP connect succeeded.
func (d *diagnosis) probe(address string) (noHostResolved, anyTCPConnected bool) {
	p := d.printer
	p.printf("Network probe (DNS, then a TCP connect to each IP address, %d ms each):",
		diagnose.ProbeTimeout.Milliseconds())
	anyHost, anyResolved := false, false
	for _, target := range diagnose.SplitAddress(address) {
		if target.Skip != "" {
			p.printf("  %s: skipped, %s", target.Part, target.Skip)
			continue
		}
		anyHost = true
		name := target.Host + ":" + strconv.Itoa(target.Port)
		probe := diagnose.ProbeHost(context.Background(), target.Host, target.Port, diagnose.ProbeTimeout)
		if probe.Err != "" {
			p.printf("  %s: DNS failed: %s", name, probe.Err)
			continue
		}
		anyResolved = true
		p.printf("  %s: resolved to %d IP address(es)", name, len(probe.IPs))
		for _, ip := range probe.IPs {
			if ip.Err == "" {
				anyTCPConnected = true
				p.printf("    %s: TCP connected in %d ms", ip.IP, ip.Elapsed.Milliseconds())
			} else {
				p.printf("    %s: TCP failed after %d ms: %s", ip.IP, ip.Elapsed.Milliseconds(), ip.Err)
			}
		}
	}
	if !anyHost {
		p.printf("  (no host:port part to probe)")
	}
	return anyHost && !anyResolved, anyTCPConnected
}

// printEvents prints the events of each type and the last event per symbol, and fills the counts of the result.
func (d *diagnosis) printEvents(typeNames []string, stats map[string]*typeStats, result *addressResult) {
	p := d.printer
	p.printf("Events:")
	for _, name := range typeNames {
		s := stats[name]
		if s.count == 0 {
			p.printf("  %s: no events", name)
			continue
		}
		result.total += s.count
		result.counts[name] = s.count
		p.printf("  %s: %d events, the first one %d ms after the connect", name, s.count,
			s.firstAfter.Milliseconds())
		for _, text := range s.first {
			p.printf("    %s", text)
		}
		if s.flags[0]+s.flags[1]+s.flags[2]+s.flags[3] > 0 {
			p.printf("    flags: SNAPSHOT_BEGIN %d, SNAPSHOT_END %d, TX_PENDING %d, REMOVE_EVENT %d", s.flags[0],
				s.flags[1], s.flags[2], s.flags[3])
		}
		if s.lagCount == 0 {
			p.printf("    receive time - event time: n/a (the events have no time)")
		} else {
			p.printf("    receive time - event time: min %d ms, avg %d ms, max %d ms", s.lagMin,
				s.lagSum/int64(s.lagCount), s.lagMax)
		}
	}
	var lines []string
	for _, name := range typeNames {
		s := stats[name]
		symbols := make([]string, 0, len(s.countBySymbol))
		for symbol := range s.countBySymbol {
			symbols = append(symbols, symbol)
			result.received[diagnose.Pair{Type: name, Symbol: symbol}] = true
		}
		sort.Strings(symbols)
		for _, symbol := range symbols {
			lines = append(lines, fmt.Sprintf("  %s %s: %s (%d events)", name, symbol,
				s.lastBySymbol[symbol].Format(diagnoseTimeFormat), s.countBySymbol[symbol]))
		}
	}
	if len(lines) > 0 {
		p.printf("Last event per symbol:")
		for _, line := range lines {
			p.printf("%s", line)
		}
	}
}

func (d *diagnosis) printComparison(first, second *addressResult) {
	p := d.printer
	p.printf("")
	p.printf("=== Comparison")
	for _, result := range []*addressResult{first, second} {
		connected := "not connected"
		if result.wasConnected {
			connected = "connected"
		}
		p.printf("  %s: %s, %d events, %s", result.address, result.verdict.Kind, result.total, connected)
	}
	names := map[string]bool{}
	for _, result := range []*addressResult{first, second} {
		for name := range result.counts {
			names[name] = true
		}
	}
	for _, name := range sortedKeys(names) {
		p.printf("  %s: %d vs %d", name, first.counts[name], second.counts[name])
	}
	for _, pair := range [][2]*addressResult{{first, second}, {second, first}} {
		var only []string
		for received := range pair[0].received {
			if !pair[1].received[received] {
				only = append(only, received.Type+" "+received.Symbol)
			}
		}
		sort.Strings(only)
		for i, text := range only {
			if i == diagnoseOnlyFromLimit {
				p.printf("  only from %s: %d more", pair[0].address, len(only)-diagnoseOnlyFromLimit)
				break
			}
			p.printf("  only from %s: %s", pair[0].address, text)
		}
	}
	if first.verdict.Kind == second.verdict.Kind {
		p.printf("  The verdicts are the same.")
	} else {
		p.printf("  The verdicts differ.")
	}
}

func connectorText(c api.MessageConnectorInfo) string {
	text := fmt.Sprintf("%s %s, %v, %v", c.Name, c.Address, c.Kind, c.State)
	if len(c.CurrentAddresses) > 0 {
		text += ", connected to " + strings.Join(c.CurrentAddresses, ",")
	}
	return text + fmt.Sprintf(", %d connections (%d closed), read %d bytes, written %d bytes", c.ConnectionCount,
		c.ClosedConnectionCount, c.Statistics.ReadBytes, c.Statistics.WrittenBytes)
}

// sortedTypes returns the event types sorted by name, without repetitions.
func sortedTypes(types []eventcodes.EventCode) []eventcodes.EventCode {
	seen := map[eventcodes.EventCode]bool{}
	var sorted []eventcodes.EventCode
	for _, eventType := range types {
		if !seen[eventType] {
			seen[eventType] = true
			sorted = append(sorted, eventType)
		}
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].String() < sorted[j].String() })
	return sorted
}

// symbolNames returns the names of the symbols for the verdict; the wildcard symbol is diagnose.Wildcard.
func symbolNames(symbols []osub.Symbol) []string {
	names := make([]string, len(symbols))
	for i, symbol := range symbols {
		if _, ok := symbol.(*osub.WildcardSymbol); ok {
			names[i] = diagnose.Wildcard
		} else {
			names[i] = fmt.Sprint(symbol)
		}
	}
	return names
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
