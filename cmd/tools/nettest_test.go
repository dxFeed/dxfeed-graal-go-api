// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/nettest"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
)

// startTools starts the tools with args in a child process; wait waits for it to exit and returns its exit status
// and output.
func startTools(t *testing.T, args string) (wait func() (int, string)) {
	t.Helper()
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), toolsArgsEnv+"="+args)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	var once sync.Once
	status := 0
	wait = func() (int, string) {
		once.Do(func() {
			if err := cmd.Wait(); err != nil {
				status = cmd.ProcessState.ExitCode()
			}
		})
		return status, out.String()
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		wait()
	})
	return wait
}

// freePort returns a free local port.
func freePort(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	return strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
}

// netTestRegexes reads the expressions of the output lines from the kit.
func netTestRegexes(t *testing.T) map[string]*regexp.Regexp {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "internal", "nettest", "testdata", "regexes.json"))
	if err != nil {
		t.Fatal(err)
	}
	var kit struct {
		Regexes map[string]string `json:"regexes"`
	}
	if err := json.Unmarshal(data, &kit); err != nil {
		t.Fatal(err)
	}
	result := map[string]*regexp.Regexp{}
	for name, expression := range kit.Regexes {
		result[name] = regexp.MustCompile(expression)
	}
	return result
}

func TestNetTestCommandLine(t *testing.T) {
	for _, test := range []struct {
		args   string
		status int
		want   string
	}{
		{"nettest", 2, "Usage:\n  tools nettest "},
		{"nettest --help", 0, "spec_version nettest/1\n"},
		{"nettest p :7001 -S 100 -h", 0, "spec_version nettest/1\n"},
		{"nettest p :7001 -S 100 -h --foo", 0, "spec_version nettest/1\n"},
		{"nettest p :7001 --foo --help", 0, "spec_version nettest/1\n"},
		{"nettest x -help", 0, "spec_version nettest/1\n"},
		{"nettest p :7001 --name=-h -S 0", 2, "nettest: -S "},
		{"nettest p :7001 -n --help -S 0", 2, "nettest: -S "},
		{"nettest p :7001 -S 100 --help=x", 2, "nettest: --help takes no value: '--help=x'\n"},
		{"nettest p :7001 -S 100 --foo", 2, "nettest: unknown flag --foo\n"},
		{"nettest p :7001", 2, "nettest: -S is missing\n"},
		{"nettest c 127.0.0.1:7001 -S 100 --batch 10", 2, "nettest: --batch is an option of the mode p\n"},
		{"nettest p :7001 -S 100 -Xmx512m", 2, "nettest: unsupported option: -Xmx512m\n"},
		{"nettest p :7001 -S 100 --legacy", 2, "nettest: unsupported option: --legacy\n"},
		{"nettest c 127.0.0.1:7001 -S 100 --summary no-such-dir/s.json", 2,
			"nettest: --summary cannot write to the directory of 'no-such-dir/s.json'\n"},
		{"nettest c 127.0.0.1:7001 -S 100 --metrics-port 65536", 2, "nettest: --metrics-port needs [<host>:]<port>"},
	} {
		status, out := runTools(t, test.args)
		if status != test.status || !strings.Contains(out, test.want) {
			t.Errorf("tools %s: status %d, want %d and %q:\n%s", test.args, status, test.status, test.want, out)
		}
		if test.status == 2 && test.want[0] == 'n' && !strings.HasPrefix(out, test.want) {
			t.Errorf("tools %s: the error is not the first line:\n%s", test.args, out)
		}
	}
}

func TestNetTestPortInUse(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	status, out := runTools(t, "nettest p :"+port+" -S 10 --duration 2s")
	if status != 1 || !strings.Contains(out, "Error during nettest: the address :"+port+" cannot be listened on") {
		t.Errorf("status %d, want 1:\n%s", status, out)
	}
	// Windows lets a socket bind 127.0.0.1 next to one on all interfaces: hold the metrics port on 127.0.0.1.
	metricsListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = metricsListener.Close() }()
	metricsPort := strconv.Itoa(metricsListener.Addr().(*net.TCPAddr).Port)
	status, out = runTools(t, "nettest c 127.0.0.1:"+port+" -S 10 --duration 2s --metrics-port 127.0.0.1:"+metricsPort)
	if status != 1 || !strings.Contains(out, "Error during nettest: the metrics address 127.0.0.1:"+metricsPort+" cannot be "+
		"listened on") {
		t.Errorf("metrics: status %d, want 1:\n%s", status, out)
	}
}

// TestNetTestListenOptions checks that a listen address with connector options starts and serves a consumer.
func TestNetTestListenOptions(t *testing.T) {
	port := freePort(t)
	summary := filepath.Join(t.TempDir(), "c.json")
	producer := startTools(t, "nettest p :"+port+"[bindAddr=127.0.0.1] -S 10 -r */1k --duration 6s")
	time.Sleep(time.Second)
	status, out := runTools(t, "nettest c 127.0.0.1:"+port+" -S 10 --duration 3s --summary "+summary)
	if status != 0 {
		t.Fatalf("consumer: status %d:\n%s", status, out)
	}
	if c := readSummary(t, summary); c["status"] != "ok" || c["symbols_seen"] != float64(10) {
		t.Errorf("consumer summary: status %v, symbols_seen %v:\n%s", c["status"], c["symbols_seen"], out)
	}
	if status, out := producer(); status != 0 {
		t.Errorf("producer: status %d:\n%s", status, out)
	}
	// A taken port still fails the start.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()
	taken := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	status, out = runTools(t, "nettest p :"+taken+"[bindAddr=127.0.0.1] -S 10 --duration 2s")
	if status != 1 || !strings.Contains(out, "Error during nettest: the address :"+taken+
		"[bindAddr=127.0.0.1] cannot be listened on") {
		t.Errorf("taken port: status %d, want 1:\n%s", status, out)
	}
}

// TestNetTestShortRun checks that a run shorter than the period prints its last block, counted from the start, when
// it lasts at least one second.
func TestNetTestShortRun(t *testing.T) {
	status, out := runTools(t, "nettest p :"+freePort(t)+" -S 10 -s 60m --duration 2s")
	if status != 0 || strings.Count(out, "{*NetTestProducerSide-NetTest*} data ") != 1 {
		t.Errorf("status %d, want 0 and one statistics block:\n%s", status, out)
	}
}

// TestNetTestLoopback runs a publisher and a consumer of the tool for some seconds: both print their lines, the
// consumer's lag is positive and below 1 s, --duration ends both with 0, a scrape of the consumer's metrics has its
// families, and both summary files validate against the schema of the kit.
func TestNetTestLoopback(t *testing.T) {
	re := netTestRegexes(t)
	port, metricsPort := freePort(t), freePort(t)
	dir := t.TempDir()
	producerSummary, consumerSummary := filepath.Join(dir, "p.json"), filepath.Join(dir, "c.json")
	producer := startTools(t, "nettest p :"+port+" -S 100 -r */5k -s 1s --duration 14s -c stream -n prod --summary "+
		producerSummary)
	time.Sleep(time.Second)
	consumer := startTools(t, "nettest c 127.0.0.1:"+port+" -S 100 -C 2 -s 1s --duration 8s -c stream -n cons "+
		"--warmup 2s --metrics-port 127.0.0.1:"+metricsPort+" --summary "+consumerSummary)
	time.Sleep(4 * time.Second)
	scrape := scrapeMetrics(t, "http://127.0.0.1:"+metricsPort+"/metrics")
	status, out := consumer()
	if status != 0 {
		t.Fatalf("consumer: status %d:\n%s", status, out)
	}
	metricsLine := "{*NetTestMetrics-cons*} http://127.0.0.1:" + metricsPort + "/metrics"
	if !re["metrics"].MatchString(metricsLine) || !strings.Contains(out, "{*NetTestMetrics-cons*} http://127.0.0.1:"+
		metricsPort+"/metrics") {
		t.Errorf("consumer: no metrics line:\n%s", out)
	}
	for _, family := range []string{"nettest_info gauge", "nettest_records_total counter",
		"nettest_e2e_lag_seconds histogram", "nettest_e2e_lag_negative_total counter", "nettest_connections gauge",
		"nettest_connector_read_records_total counter", "nettest_runtime_cpu_seconds_total counter"} {
		if !strings.Contains(scrape, "# TYPE "+family+"\n") {
			t.Errorf("the scrape has no %s:\n%s", family, scrape)
		}
	}
	if !strings.Contains(scrape, `nettest_records_total{tool="go-api",side="consumer",name="cons",connection="2"} `) {
		t.Errorf("the scrape has no records of the second instance:\n%s", scrape)
	}
	blocks, instances, runtimes := countLines(t, re, out, "Consumer", "cons", "2")
	if blocks < 3 || instances != 2*blocks || runtimes != blocks {
		t.Errorf("consumer: %d blocks, %d instance lines, %d runtime lines; want 3 or more blocks, each with 2 "+
			"instance lines and a runtime line:\n%s", blocks, instances, runtimes, out)
	}
	status, out = producer()
	blocks, instances, runtimes = countLines(t, re, out, "Producer", "prod", "1")
	if status != 0 || blocks < 3 || instances != blocks || runtimes != blocks {
		t.Errorf("producer: status %d, %d blocks, %d instance lines, %d runtime lines; want 0 and 3 or more blocks:\n%s",
			status, blocks, instances, runtimes, out)
	}
	c := readSummary(t, consumerSummary)
	if c["status"] != "ok" || c["symbols_seen"] != float64(100) || c["time_resolution"] != "ns" {
		t.Errorf("consumer summary: status %v, symbols_seen %v, time_resolution %v", c["status"], c["symbols_seen"],
			c["time_resolution"])
	}
	phases := map[any]int{}
	for _, interval := range c["intervals"].([]any) {
		phases[interval.(map[string]any)["phase"]]++
	}
	lag := c["e2e_lag"].(map[string]any)
	if phases["warmup"] == 0 || phases["measure"] == 0 || lag["count"].(float64) <= 0 ||
		lag["p50_ns"].(float64) <= 0 || lag["p50_ns"].(float64) >= 1e9 {
		t.Errorf("consumer summary: phases %v, e2e_lag %v", phases, lag)
	}
	p := readSummary(t, producerSummary)
	if p["status"] != "ok" || p["schedule_slip"] == nil || p["publish_call"].(map[string]any)["calls"].(float64) <= 0 ||
		p["e2e_lag"] != nil {
		t.Errorf("producer summary: status %v, schedule_slip %v, publish_call %v", p["status"], p["schedule_slip"],
			p["publish_call"])
	}
}

// scrapeMetrics returns a scrape of the metrics and checks its content type.
func scrapeMetrics(t *testing.T, url string) string {
	t.Helper()
	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got := response.Header.Get("Content-Type"); got != "text/plain; version=0.0.4; charset=utf-8" {
		t.Errorf("the content type is %q", got)
	}
	return string(body)
}

// readSummary validates a summary file against the schema of the kit, checks that the records of its intervals add
// up and returns it.
func readSummary(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	schema, err := os.ReadFile(filepath.Join("..", "..", "internal", "nettest", "testdata",
		"nettest-summary.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := nettest.ValidateJSON(schema, data); err != nil {
		t.Errorf("%s: %v\n%s", path, err, data)
	}
	var summary map[string]any
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatal(err)
	}
	// The library of SDK 3.19.0 and newer gives its build information.
	tool := summary["tool"].(map[string]any)
	build := tool["build"].(map[string]any)
	if tool["qd_version"] == nil || build["vendor"] == nil || build["gc"] == nil {
		t.Errorf("%s: tool %v", path, tool)
	}
	var records float64
	for _, interval := range summary["intervals"].([]any) {
		records += interval.(map[string]any)["records"].(float64)
	}
	if records != summary["records"] {
		t.Errorf("%s: the intervals have %v records, the summary %v", path, records, summary["records"])
	}
	return summary
}

// countLines counts the statistics blocks, the instance lines and the runtime lines of the output of a side and checks
// the side, the name and the instances of each block, and on the consumer a lag from 0 to 1 s (a block of a burst on
// a busy machine may have a mean below 1 us; the summary checks that the lag is positive).
func countLines(t *testing.T, re map[string]*regexp.Regexp, out, side, name, instances string) (int, int, int) {
	t.Helper()
	blocks, instanceLines, runtimes := 0, 0, 0
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimRight(line, "\r")
		switch {
		case re["summary"].MatchString(line):
			blocks++
			m := re["summary"].FindStringSubmatch(line)
			lag, _ := strconv.Atoi(strings.ReplaceAll(m[4], ",", ""))
			if m[1] != side || m[2] != name || m[5] != instances || side == "Consumer" && m[3] != "0" && (lag < 0 || lag >= 1000000) {
				t.Errorf("line %q: want %s, %s, %s instances and a lag from 0 to 1 s on the consumer", line, side,
					name, instances)
			}
		case re["instance"].MatchString(line):
			instanceLines++
		case re["runtime"].MatchString(line):
			runtimes++
		}
	}
	return blocks, instanceLines, runtimes
}

// receiveTrades connects a STREAM_FEED endpoint with dxscheme.nanoTime=true to the address and returns the nanosecond
// parts of the trades it receives in the time.
func receiveTrades(t *testing.T, address string, symbols []string, wait time.Duration) []int32 {
	t.Helper()
	endpoint, err := api.NewEndpointBuilder().WithRole(api.StreamFeed).WithProperty("dxscheme.nanoTime", "true").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()
	feed, err := endpoint.GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	subscription, err := feed.CreateSubscription(eventcodes.Trade)
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var nanos []int32
	if err := subscription.AddListener(listenerFunc(func(events []any) {
		mu.Lock()
		defer mu.Unlock()
		for _, event := range events {
			if e, ok := event.(*trade.Trade); ok {
				nanos = append(nanos, e.TimeNanoPart())
			}
		}
	})); err != nil {
		t.Fatal(err)
	}
	for _, symbol := range symbols {
		if err := subscription.AddSymbols(osub.String(symbol)); err != nil {
			t.Fatal(err)
		}
	}
	if err := endpoint.Connect(address); err != nil {
		t.Fatal(err)
	}
	time.Sleep(wait)
	mu.Lock()
	defer mu.Unlock()
	return append([]int32(nil), nanos...)
}

type listenerFunc func(events []any)

func (f listenerFunc) Update(events []any) { f(events) }

// TestNetTestNanoTime checks the nanosecond part: the tool's publisher sends it, and does not with
// -Ddxscheme.nanoTime=false; the tool's consumer stops with nano_missing against such a publisher, and runs to the
// end with -Ddxscheme.nanoTime=false; the property is case-insensitive.
func TestNetTestNanoTime(t *testing.T) {
	symbols := []string{"YQKNT", "XRR", "T", "NDFY", "EUKK", "UXCZWE", "JBZY", "QJD", "LENRR", "IC"}
	for _, nanoTime := range []bool{true, false} {
		port := freePort(t)
		args := "nettest p :" + port + " -S 10 -r */2k -s 1s --duration 60s"
		if !nanoTime {
			args += " -Ddxscheme.nanoTime=false"
		}
		producer := startTools(t, args)
		time.Sleep(time.Second)
		nanos := receiveTrades(t, "127.0.0.1:"+port, symbols, 2*time.Second)
		nonzero := 0
		for _, nano := range nanos {
			if nano != 0 {
				nonzero++
			}
		}
		if len(nanos) < 100 || nanoTime != (nonzero > len(nanos)/2) || !nanoTime && nonzero != 0 {
			t.Errorf("nanoTime %v: %d of %d trades have a nanosecond part", nanoTime, nonzero, len(nanos))
		}
		if !nanoTime {
			summary := filepath.Join(t.TempDir(), "c.json")
			status, out := runTools(t, "nettest c 127.0.0.1:"+port+" -S 10 -s 1s --duration 20s --summary "+summary)
			if status != 1 || !strings.Contains(out, "Error during nettest: nano_missing: ") {
				t.Errorf("consumer: status %d, want 1 and nano_missing:\n%s", status, out)
			}
			if s := readSummary(t, summary); s["status"] != "nano_missing" || s["exit_code"] != float64(1) {
				t.Errorf("consumer summary: status %v, exit_code %v", s["status"], s["exit_code"])
			}
			status, out = runTools(t, "nettest c 127.0.0.1:"+port+" -S 10 -s 1s --duration 20s -Ddxscheme.nanoTime=TRUE")
			if status != 1 || !strings.Contains(out, "Error during nettest: nano_missing: ") {
				t.Errorf("consumer with -Ddxscheme.nanoTime=TRUE: status %d, want 1 and nano_missing:\n%s",
					status, out)
			}
			status, out = runTools(t, "nettest c 127.0.0.1:"+port+" -S 10 -s 1s --duration 2s -Ddxscheme.nanoTime=false")
			if status != 0 {
				t.Errorf("consumer with -Ddxscheme.nanoTime=false: status %d, want 0:\n%s", status, out)
			}
		}
		_ = producer // the cleanup stops it
	}
}

// TestWallNanosResolution checks that the event clock has the sub-millisecond resolution: time.Now() on Windows
// steps by up to 1 ms, so 100 steps would take at least 50 ms and millions of reads.
func TestWallNanosResolution(t *testing.T) {
	const steps, maxReads = 100, 100_000
	last, changes := wallNanos(), 0
	for reads := 0; changes < steps; reads++ {
		if reads == maxReads {
			t.Fatalf("the event clock changed %d times in %d reads, want %d", changes, maxReads, steps)
		}
		if now := wallNanos(); now != last {
			last, changes = now, changes+1
		}
	}
}

func TestNextTick(t *testing.T) {
	for _, test := range []struct {
		now    int64
		period time.Duration
		want   int64
	}{
		{1700000000123456789, 10 * time.Second, 1700000010000000000},
		{1700000003000000000, 7 * time.Second, 1700000008000000000},
		{1700000008000000000, 7 * time.Second, 1700000015000000000},
		{1700000000123456789, 500 * time.Millisecond, 1700000000500000000},
	} {
		if got := nextTick(time.Unix(0, test.now), test.period).UnixNano(); got != test.want {
			t.Errorf("nextTick(%d, %v) = %d, want %d", test.now, test.period, got, test.want)
		}
	}
}

func TestPoolTrade(t *testing.T) {
	p := newPoolTrade("A")
	for name, value := range map[string]float64{"price": p.Price(), "change": p.Change(), "size": p.Size(),
		"dayVolume": p.DayVolume(), "dayTurnover": p.DayTurnover()} {
		if !math.IsNaN(value) {
			t.Errorf("%s is %v, want NaN", name, value)
		}
	}
}

// TestNetTestWatchdog checks the times without data: they count only with open connections, from the connection
// again after the connections were all closed.
func TestNetTestWatchdog(t *testing.T) {
	w := netTestWatchdog{afterRecord: 30 * time.Second, afterConnection: 60 * time.Second}
	start := time.Unix(1700000000, 0)
	at := func(seconds int) time.Time { return start.Add(time.Duration(seconds) * time.Second) }
	for _, step := range []struct {
		now, connections, lastRecord int
		want                         string
	}{
		{0, 0, -1, ""}, {100, 0, -1, ""}, // listening without peers: no limit
		{101, 1, -1, ""}, {160, 1, -1, ""}, {161, 1, -1, "no records for 60000 ms after the connection, with 1 open " +
			"connections"},
		{162, 0, -1, ""}, {170, 2, -1, ""}, {229, 2, -1, ""}, // the connections closed and opened again
		{230, 2, 201, ""}, {259, 2, 230, ""}, {260, 2, 230, "no records for 30000 ms after the last record, with 2 open " +
			"connections"},
		{261, 0, 230, ""}, {300, 1, 230, ""}, {329, 1, 230, ""}, // the time counts from the new connection
		{330, 1, 230, "no records for 30000 ms after the last record, with 1 open connections"},
	} {
		var lastRecord time.Time
		if step.lastRecord >= 0 {
			lastRecord = at(step.lastRecord)
		}
		err := w.check(at(step.now), step.connections, lastRecord)
		if got := fmt.Sprint(err); err == nil && step.want != "" || err != nil && got != step.want {
			t.Errorf("at %d s with %d connections: %v, want %q", step.now, step.connections, err, step.want)
		}
	}
}

// TestNetTestWatchdogListening runs the watchdog on a consumer endpoint that listens on a port: it does not stop
// without a peer, although the endpoint is in the state CONNECTED, and stops with no_data when a peer connects and
// sends nothing.
func TestNetTestWatchdogListening(t *testing.T) {
	port := freePort(t)
	endpoint, err := api.NewEndpointBuilder().WithRole(api.StreamFeed).Build()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = endpoint.CloseAndAwaitTermination(context.Background()) }()
	if err := endpoint.Connect(":" + port); err != nil {
		t.Fatal(err)
	}
	r := &netTestRun{start: time.Now(), stop: make(chan struct{}),
		instances: []*netTestInstance{{endpoint: endpoint}},
		watchdog:  netTestWatchdog{afterRecord: time.Second, afterConnection: time.Second}}
	r.connections = r.openConnections
	r.workers.Add(1)
	go r.watchConsumer()
	defer func() {
		r.end()
		r.workers.Wait()
	}()
	time.Sleep(2500 * time.Millisecond)
	if state, _ := endpoint.State(); state != common.Connected || r.openConnections() != 0 {
		t.Errorf("the listening endpoint: state %v, %d open connections", state, r.openConnections())
	}
	select {
	case <-r.stop:
		t.Fatalf("stopped without a peer: %v", r.failure)
	default:
	}
	peer, err := api.NewEndpointBuilder().WithRole(api.Publisher).Build()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = peer.CloseAndAwaitTermination(context.Background()) }()
	if err := peer.Connect("127.0.0.1:" + port); err != nil {
		t.Fatal(err)
	}
	select {
	case <-r.stop:
	case <-time.After(15 * time.Second):
		t.Fatal("no stop with a silent peer")
	}
	var statusErr *netTestStatusError
	if !errors.As(r.failure, &statusErr) || statusErr.status != nettest.StatusNoData ||
		!strings.Contains(r.failure.Error(), "after the connection, with 1 open connections") {
		t.Errorf("the failure %v, want no_data after the connection", r.failure)
	}
}
