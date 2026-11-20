// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"encoding/json"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
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
		{"nettest p :7001 -S 100 --foo", 2, "nettest: unknown flag --foo\n"},
		{"nettest p :7001", 2, "nettest: -S is missing\n"},
		{"nettest c 127.0.0.1:7001 -S 100 --batch 10", 2, "nettest: --batch is an option of the mode p\n"},
		{"nettest p :7001 -S 100 -Xmx512m", 2, "nettest: unsupported option: -Xmx512m\n"},
		{"nettest p :7001 -S 100 --legacy", 2, "nettest: unsupported option: --legacy\n"},
		{"nettest c 127.0.0.1:7001 -S 100 --summary s.json", 2,
			"nettest: unsupported option: --summary (not yet implemented)\n"},
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
}

// TestNetTestLoopback runs a publisher and a consumer of the tool for some seconds: both print their lines, the
// consumer's lag is positive and below 1 s, and --duration ends both with 0.
func TestNetTestLoopback(t *testing.T) {
	re := netTestRegexes(t)
	port := freePort(t)
	producer := startTools(t, "nettest p :"+port+" -S 100 -r */5k -s 1s --duration 14s -c stream -n prod")
	time.Sleep(time.Second)
	status, out := runTools(t, "nettest c 127.0.0.1:"+port+" -S 100 -C 2 -s 1s --duration 8s -c stream -n cons")
	if status != 0 {
		t.Fatalf("consumer: status %d:\n%s", status, out)
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
}

// countLines counts the statistics blocks, the instance lines and the runtime lines of the output of a side and checks
// the side, the name and the instances of each block, and on the consumer a lag from 0 to 1 s.
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
			if m[1] != side || m[2] != name || m[5] != instances || side == "Consumer" && m[3] != "0" && (lag <= 0 || lag >= 1000000) {
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
// end with -Ddxscheme.nanoTime=false.
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
			status, out := runTools(t, "nettest c 127.0.0.1:"+port+" -S 10 -s 1s --duration 20s")
			if status != 1 || !strings.Contains(out, "Error during nettest: nano_missing: ") {
				t.Errorf("consumer: status %d, want 1 and nano_missing:\n%s", status, out)
			}
			status, out = runTools(t, "nettest c 127.0.0.1:"+port+" -S 10 -s 1s --duration 2s -Ddxscheme.nanoTime=false")
			if status != 0 {
				t.Errorf("consumer with -Ddxscheme.nanoTime=false: status %d, want 0:\n%s", status, out)
			}
		}
		_ = producer // the cleanup stops it
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
