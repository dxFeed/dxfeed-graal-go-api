// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// readKit reads a JSON file of the kit in testdata.
func readKit(t *testing.T, name string, v any) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
}

// features are the optional features of the specification that the tool has.
var features = map[string]bool{"isolate": false, "legacy": false}

func TestCLIVectors(t *testing.T) {
	var kit struct {
		Vectors []struct {
			Section        string         `json:"section"`
			Args           []string       `json:"args"`
			Exit           int            `json:"exit"`
			Error          string         `json:"error"`
			ErrorPrefix    string         `json:"error_prefix"`
			Requires       []string       `json:"requires"`
			Settings       map[string]any `json:"settings"`
			StdoutContains string         `json:"stdout_contains"`
		} `json:"vectors"`
	}
	readKit(t, "cli-vectors.json", &kit)
	skipped := map[string]int{}
	for _, v := range kit.Vectors {
		if !hasFeatures(v.Requires) {
			skipped[v.Section+" "+strings.Join(v.Requires, ",")]++
			continue
		}
		name := fmt.Sprintf("%q", v.Args)
		c, err := Parse(v.Args)
		if err == nil && c.Help {
			if v.Exit != 0 || !strings.Contains(HelpText, v.StdoutContains) {
				t.Errorf("%s: the help has no %q", name, v.StdoutContains)
			}
			continue
		}
		if v.Exit == 0 {
			if err != nil {
				t.Errorf("%s: %v", name, err)
				continue
			}
			settings := c.Settings()
			for key, want := range v.Settings {
				got, ok := settings[key]
				if !ok || fmt.Sprint(normalize(got)) != fmt.Sprint(normalize(want)) {
					t.Errorf("%s: %s is %v, want %v", name, key, got, want)
				}
			}
			continue
		}
		if err == nil {
			t.Errorf("%s: no error, want exit %d %q", name, v.Exit, v.Error)
			continue
		}
		if v.Error != "" && err.Error() != v.Error ||
			v.ErrorPrefix != "" && !strings.HasPrefix(err.Error(), v.ErrorPrefix) {
			t.Errorf("%s: error %q, want %q%q", name, err, v.Error, v.ErrorPrefix)
		}
	}
	t.Logf("skipped: %v", skipped)
}

func hasFeatures(requires []string) bool {
	for _, r := range requires {
		if name, not := strings.CutPrefix(r, "!"); not == features[name] {
			return false
		}
	}
	return true
}

// normalize makes the numbers of JSON and of Settings comparable.
func normalize(v any) any {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int64:
		return float64(n)
	}
	return v
}

func TestLines(t *testing.T) {
	var kit struct {
		Grouping []struct {
			Value int64  `json:"value"`
			Text  string `json:"text"`
		} `json:"grouping"`
		Statistics []struct {
			Side        string `json:"side"`
			Name        string `json:"name"`
			IntervalNs  int64  `json:"interval_ns"`
			Connections int    `json:"connections"`
			Instances   []struct {
				Records    int64 `json:"records"`
				WeightedNs int64 `json:"weighted_ns"`
				LagSumNs   int64 `json:"lag_sum_ns"`
			} `json:"instances"`
			Text string `json:"text"`
		} `json:"statistics"`
		Runtime []struct {
			Name       string            `json:"name"`
			IntervalNs int64             `json:"interval_ns"`
			Previous   map[string]*int64 `json:"previous"`
			Current    map[string]*int64 `json:"current"`
			Text       string            `json:"text"`
		} `json:"runtime"`
		Captured []string `json:"captured"`
	}
	readKit(t, "lines.json", &kit)
	for _, g := range kit.Grouping {
		if got := Group(g.Value); got != g.Text {
			t.Errorf("Group(%d) = %q, want %q", g.Value, got, g.Text)
		}
	}
	var regexes struct {
		Regexes map[string]string `json:"regexes"`
	}
	readKit(t, "regexes.json", &regexes)
	re := map[string]*regexp.Regexp{}
	for name, expression := range regexes.Regexes {
		re[name] = regexp.MustCompile(expression)
	}
	for _, s := range kit.Statistics {
		mode := Producer
		if s.Side == "consumer" {
			mode = Consumer
		}
		var instances []InstanceStatistics
		for _, i := range s.Instances {
			instances = append(instances, InstanceStatistics{Records: i.Records, LagSum: i.WeightedNs + i.LagSumNs})
		}
		got := StatisticsBlock(mode, s.Name, s.IntervalNs, s.Connections, instances)
		if got != s.Text {
			t.Errorf("statistics:\n%s\nwant\n%s", got, s.Text)
		}
		for i, line := range strings.Split(got, "\n") {
			if expression := map[bool]string{true: "summary", false: "instance"}[i == 0]; !re[expression].MatchString(line) {
				t.Errorf("%q does not match %s", line, expression)
			}
		}
	}
	sample := func(m map[string]*int64) RuntimeSample {
		value := func(key string) int64 {
			if v := m[key]; v != nil {
				return *v
			}
			return -1
		}
		return RuntimeSample{HeapUsedBytes: value("heap_used_bytes"), HeapCommittedBytes: value("heap_committed_bytes"),
			GCCount: value("gc_count"), GCTimeMs: value("gc_time_ms"), CPUTimeNs: value("cpu_ns")}
	}
	for _, r := range kit.Runtime {
		got := RuntimeLine(r.Name, r.IntervalNs, sample(r.Previous), sample(r.Current))
		if got != r.Text || !re["runtime"].MatchString(got) {
			t.Errorf("runtime %q, want %q", got, r.Text)
		}
	}
	for _, line := range kit.Captured {
		matched := false
		for _, expression := range re {
			matched = matched || expression.MatchString(line)
		}
		if !matched {
			t.Errorf("captured %q matches no expression", line)
		}
	}
}

func hex64(t *testing.T, text string) uint64 {
	t.Helper()
	v, err := strconv.ParseUint(strings.TrimPrefix(text, "0x"), 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestPools(t *testing.T) {
	var kit struct {
		Instances []struct {
			Instance     int      `json:"instance"`
			Seed         string   `json:"seed"`
			InitialState string   `json:"initial_state"`
			FirstDraws   []string `json:"first_draws"`
			TableRows    []struct {
				Row   int     `json:"row"`
				Price float64 `json:"price"`
				Size  float64 `json:"size"`
			} `json:"table_rows"`
			Events struct {
				PoolSize    int `json:"pool_size"`
				FirstEvents []struct {
					Draw        string  `json:"draw"`
					SymbolIndex int     `json:"symbol_index"`
					Row         int     `json:"row"`
					Price       float64 `json:"price"`
					Size        float64 `json:"size"`
				} `json:"first_events"`
			} `json:"events"`
		} `json:"instances"`
		PackTimeSequence []struct {
			EpochNs  int64  `json:"epoch_ns"`
			Sequence int    `json:"sequence"`
			TimeMs   int64  `json:"time_ms"`
			NanoPart int32  `json:"nano_part"`
			Packed   string `json:"packed"`
		} `json:"pack_time_sequence"`
	}
	readKit(t, "pool-vectors.json", &kit)
	for _, i := range kit.Instances {
		if Seed(i.Instance) != hex64(t, i.Seed) {
			t.Errorf("instance %d: seed %#x", i.Instance, Seed(i.Instance))
		}
		r := NewRandom(i.Instance)
		if r.State() != hex64(t, i.InitialState) {
			t.Errorf("instance %d: state %#x", i.Instance, r.State())
		}
		for k, draw := range i.FirstDraws {
			if got := r.Next(); got != hex64(t, draw) {
				t.Errorf("instance %d: draw %d is %#x, want %s", i.Instance, k, got, draw)
			}
		}
		r = NewRandom(i.Instance)
		table := NewTable(r)
		for _, row := range i.TableRows {
			if table[row.Row] != (Row{row.Price, row.Size}) {
				t.Errorf("instance %d: row %d is %v, want %v", i.Instance, row.Row, table[row.Row], row)
			}
		}
		for k, e := range i.Events.FirstEvents {
			draw := r.Next()
			symbol, row := Choose(draw, i.Events.PoolSize)
			if draw != hex64(t, e.Draw) || symbol != e.SymbolIndex || row != e.Row ||
				table[row] != (Row{e.Price, e.Size}) {
				t.Errorf("instance %d: event %d is %#x %d %d %v, want %+v", i.Instance, k, draw, symbol, row,
					table[row], e)
			}
		}
	}
	for _, p := range kit.PackTimeSequence {
		timeMs, nano := SplitTime(p.EpochNs)
		if timeMs != p.TimeMs || nano != p.NanoPart ||
			uint64(PackTimeSequence(timeMs, p.Sequence)) != hex64(t, p.Packed) {
			t.Errorf("%d: %d %d %#x, want %+v", p.EpochNs, timeMs, nano, PackTimeSequence(timeMs, p.Sequence), p)
		}
	}
}

func TestGeneratedSymbols(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "symbols-generated-10.txt"))
	if err != nil {
		t.Fatal(err)
	}
	first, _ := GenerateSymbols(10)
	if want := strings.Fields(string(data)); strings.Join(first, " ") != strings.Join(want, " ") {
		t.Errorf("the first ten are %q, want %q", first, want)
	}
	var kit struct {
		Cases []struct {
			Total   int    `json:"total"`
			Draws   int    `json:"draws"`
			FNV1a64 string `json:"fnv1a64"`
			Samples []struct {
				Index  int    `json:"index"`
				Symbol string `json:"symbol"`
			} `json:"samples"`
		} `json:"cases"`
	}
	readKit(t, "symbols-generated-vectors.json", &kit)
	for _, c := range kit.Cases {
		symbols, draws := GenerateSymbols(c.Total)
		hash := fnv.New64a()
		for _, symbol := range symbols {
			hash.Write([]byte(symbol + "\n"))
		}
		if draws != c.Draws || hash.Sum64() != hex64(t, c.FNV1a64) {
			t.Errorf("%d: draws %d, hash %#x; want %d, %s", c.Total, draws, hash.Sum64(), c.Draws, c.FNV1a64)
		}
		for _, s := range c.Samples {
			if symbols[s.Index] != s.Symbol {
				t.Errorf("%d: symbol %d is %s, want %s", c.Total, s.Index, symbols[s.Index], s.Symbol)
			}
		}
	}
}

func TestSlice(t *testing.T) {
	symbols := []string{"A", "B", "C", "D", "E"}
	for _, test := range []struct {
		instance, instances int
		want                string
	}{{1, 1, "A B C D E"}, {1, 2, "A C E"}, {2, 2, "B D"}, {3, 3, "C"}, {4, 8, "D"}, {6, 8, ""}} {
		if got := strings.Join(Slice(symbols, test.instance, test.instances), " "); got != test.want {
			t.Errorf("Slice(%d, %d) = %q, want %q", test.instance, test.instances, got, test.want)
		}
	}
}

// TestRateBounds checks the bounds of -r: below 1 and above 1000G are bad rates, with the text of the SDK tools.
func TestRateBounds(t *testing.T) {
	for _, test := range []struct {
		rate string
		ok   bool
	}{{"*/1", true}, {"*/0.001k", true}, {"*/1000G", true}, {"*/0.5", false}, {"*/0.0009k", false},
		{"*/1000.1G", false}, {"*/2000G", false}} {
		_, err := Parse([]string{"p", ":7001", "-S", "10", "-r", test.rate})
		want := "-r has a bad rate (a positive number with an optional k, M or G): '" + test.rate + "'"
		if test.ok && err != nil || !test.ok && (err == nil || err.Error() != want) {
			t.Errorf("-r %s: %v", test.rate, err)
		}
	}
}
