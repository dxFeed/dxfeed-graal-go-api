// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func readKitBytes(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// families returns the HELP and TYPE lines of an exposition.
func families(text string) []string {
	var result []string
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "# ") {
			result = append(result, line)
		}
	}
	sort.Strings(result)
	return result
}

func TestMetricsFamilies(t *testing.T) {
	for _, side := range []string{"consumer", "producer"} {
		sample, err := os.ReadFile(filepath.Join("testdata", "metrics-sample-"+side+".txt"))
		if err != nil {
			t.Fatal(err)
		}
		s := &MetricsSnapshot{Info: MetricsInfo{Tool: "go-api", Side: side, Name: "bench"}, Records: []int64{5, 6},
			Connections: []int{1, 1}, Connector: ConnectorMetrics{V2: &ConnectorMetricsV2{}},
			Runtime: RuntimeSample{HeapUsedBytes: 1, HeapCommittedBytes: 2, GCCount: 3, GCTimeMs: 4, CPUTimeNs: 5}}
		if side == "consumer" {
			s.Lag = &LagMetrics{Buckets: make([]int64, len(LagBuckets)+1)}
		}
		got := FormatMetrics(s)
		if want := families(string(sample)); !reflect.DeepEqual(families(got), want) {
			t.Errorf("%s: the families are\n%s\nwant\n%s", side, strings.Join(families(got), "\n"),
				strings.Join(want, "\n"))
		}
		if !strings.Contains(got, `nettest_records_total{tool="go-api",side="`+side+`",name="bench",connection="2"} 6`) {
			t.Errorf("%s: no series of the second instance:\n%s", side, got)
		}
	}
}

func TestMetricsWithoutV2AndRuntime(t *testing.T) {
	got := FormatMetrics(&MetricsSnapshot{Info: MetricsInfo{Tool: "go-api", Side: "producer", Name: "n"},
		Records: []int64{1}, Connections: []int{0},
		Runtime: RuntimeSample{HeapUsedBytes: -1, HeapCommittedBytes: -1, GCCount: -1, GCTimeMs: -1, CPUTimeNs: -1}})
	if strings.Contains(got, "nettest_connector_rtt_seconds_total") || strings.Contains(got, "nettest_runtime_gc_total{") ||
		!strings.Contains(got, "# TYPE nettest_runtime_gc_total counter") {
		t.Errorf("without v2 and runtime:\n%s", got)
	}
}

func TestLagRecorder(t *testing.T) {
	l := NewLagRecorder()
	l.Record([]int64{-5, 0, 50000, 60000}, []string{"A", "B", "A", "C"})
	l.Record([]int64{2000000000}, []string{"A"})
	m := l.Metrics()
	want := []int64{3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 5}
	if !reflect.DeepEqual(m.Buckets, want) || m.Count != 5 || m.Negative != 1 || m.SumNs != 2000109995 {
		t.Errorf("metrics %+v, want the buckets %v", m, want)
	}
	if l.SymbolsSeen() != 3 {
		t.Errorf("symbols seen %d", l.SymbolsSeen())
	}
	in := l.TakeInterval()
	if in.Count != 5 || in.Max != 2000000000 || in.Negative != 1 || in.Histogram.Count() != 5 {
		t.Errorf("interval %+v", in)
	}
	if next := l.TakeInterval(); next.Count != 0 || next.Histogram.Count() != 0 {
		t.Errorf("next interval %+v", next)
	}
}

// consumerSummary returns the summary of a consumer with a warm-up interval of lags of 1 s and a measured one of lags
// from 1 to 1000 us.
func consumerSummary(t *testing.T) *Summary {
	t.Helper()
	config, err := Parse([]string{"c", "127.0.0.1:7001", "-S", "100", "--warmup", "10s"})
	if err != nil {
		t.Fatal(err)
	}
	b := NewSummaryBuilder(Consumer)
	start := time.Unix(1700000000, 0)
	warm := NewLagRecorder()
	warm.Record([]int64{1000000000, 1000000000}, nil)
	warmLag := warm.TakeInterval()
	b.Add(IntervalCounts{Start: start, End: start.Add(10 * time.Second), Warmup: true, Records: 2, Lag: &warmLag})
	measured := NewLagRecorder()
	var lags []int64
	for i := int64(1); i <= 1000; i++ {
		lags = append(lags, i*1000)
	}
	lags = append(lags, -7)
	measured.Record(lags, nil)
	lag := measured.TakeInterval()
	b.Add(IntervalCounts{Start: start.Add(10 * time.Second), End: start.Add(20 * time.Second), Records: 1001,
		Lag: &lag, Connections: 1})
	return &Summary{Schema: "nettest-summary/1", SpecVersion: SpecVersion,
		Tool: SummaryTool{Implementation: "go-api", Name: "tools nettest"}, Side: "consumer", Name: config.Name,
		Args: []string{}, Config: config.SummaryConfig(true), StartEpochNs: start.UnixNano(),
		EndEpochNs: start.Add(20 * time.Second).UnixNano(), Status: StatusOK, Records: b.Records,
		Rate: Rate(b.Records, 20e9), TimeSource: "actual", TimeResolution: "ns", Intervals: b.Intervals,
		E2ELag: b.Lag(), SymbolsSeen: Int(100), Connectors: []SummaryConnector{{Address: "127.0.0.1:7001",
			V1: &SummaryConnectorV1{ActiveConnections: 1}}}}
}

func TestSummaryBuilder(t *testing.T) {
	s := consumerSummary(t)
	lag := s.E2ELag
	if s.Records != 1003 || len(s.Intervals) != 2 || s.Intervals[0].Phase != "warmup" ||
		s.Intervals[1].Phase != "measure" {
		t.Fatalf("records %d, intervals %+v", s.Records, s.Intervals)
	}
	// The warm-up lags of 1 s are left out of the totals; the negative lag enters the mean and the count.
	if lag.Count != 1001 || lag.Negative != 1 || *lag.MaxNs != 1000000 || *lag.MeanNs != (500500000-7)/1001 ||
		*lag.P50Ns > 501000 || *lag.P50Ns < 499000 || *lag.P9999Ns > 1001000 {
		t.Errorf("e2e lag %+v", lag)
	}
	if *s.Intervals[0].LagMaxNs != 1000000000 || s.Intervals[0].PublishNsPerEvent != nil {
		t.Errorf("warm-up interval %+v", s.Intervals[0])
	}
	b := NewSummaryBuilder(Producer)
	start := time.Unix(1700000000, 0)
	for i, slip := range []int64{0, 1000, 3000} {
		b.Add(IntervalCounts{Start: start.Add(time.Duration(i) * time.Second),
			End: start.Add(time.Duration(i+1) * time.Second), Warmup: i == 0, Records: 100, PublishNs: 5000,
			Slip: true, SlipMax: slip + 10, SlipEnd: slip})
	}
	if slip := b.Slip(); slip.MaxNs != 3010 || slip.EndNs != 3000 || *slip.GrowthNsPerS != 2000 ||
		*b.Intervals[0].PublishNsPerEvent != 50 || b.Lag().Count != 0 || b.Lag().P50Ns != nil {
		t.Errorf("slip %+v, intervals %+v", slip, b.Intervals)
	}
	if NewSummaryBuilder(Producer).Slip() != nil {
		t.Error("a slip without a schedule")
	}
}

func TestSummarySchema(t *testing.T) {
	schema := readKitBytes(t, "nettest-summary.schema.json")
	path := filepath.Join(t.TempDir(), "summary.json")
	if err := CheckSummaryPath(path); err != nil {
		t.Fatal(err)
	}
	if err := WriteSummary(path, consumerSummary(t)); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("the temporary file is left: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateJSON(schema, data); err != nil {
		t.Errorf("%v\n%s", err, data)
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	// The validator rejects a missing key, an extra key and a wrong value.
	for name, change := range map[string]func(map[string]any){
		"missing": func(d map[string]any) { delete(d, "e2e_lag") },
		"extra":   func(d map[string]any) { d["extra"] = 1 },
		"status":  func(d map[string]any) { d["status"] = "fine" },
		"nested":  func(d map[string]any) { d["config"].(map[string]any)["symbols"] = 0 },
	} {
		var copied map[string]any
		_ = json.Unmarshal(data, &copied)
		change(copied)
		changed, _ := json.Marshal(copied)
		if ValidateJSON(schema, changed) == nil {
			t.Errorf("%s: the validator accepts it", name)
		}
	}
}

func TestSummaryPath(t *testing.T) {
	dir := t.TempDir()
	for _, path := range []string{"", filepath.Join(dir, "no-such-dir", "s.json"), dir, dir + string(filepath.Separator)} {
		if err := CheckSummaryPath(path); err == nil || !strings.HasPrefix(err.Error(), "--summary cannot write") {
			t.Errorf("%q: %v", path, err)
		}
	}
}

func TestParseMetricsAddress(t *testing.T) {
	for value, want := range map[string]string{"0": "127.0.0.1:0", "9100": "127.0.0.1:9100",
		"0.0.0.0:9100": "0.0.0.0:9100", ":9100": "127.0.0.1:9100", "[::1]:9100": "[::1]:9100"} {
		if got, err := ParseMetricsAddress(value); err != nil || got != want {
			t.Errorf("%q: %q %v, want %q", value, got, err, want)
		}
	}
	for _, value := range []string{"65536", "abc", "-1", "+1", "host:", ""} {
		if _, err := ParseMetricsAddress(value); err == nil {
			t.Errorf("%q: no error", value)
		}
	}
}
