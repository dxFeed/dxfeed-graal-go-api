// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import "testing"

func TestHistogramVectors(t *testing.T) {
	var kit struct {
		Histogram struct {
			Lowest  int64 `json:"lowest"`
			Highest int64 `json:"highest"`
			Digits  int   `json:"digits"`
		} `json:"histogram"`
		Vectors []struct {
			About         string     `json:"about"`
			Recorded      [][2]int64 `json:"recorded"`
			RecordedRange *struct {
				From      int64 `json:"from"`
				To        int64 `json:"to"`
				CountEach int64 `json:"count_each"`
			} `json:"recorded_range"`
			Count   int64 `json:"count"`
			P50Ns   int64 `json:"p50_ns"`
			P90Ns   int64 `json:"p90_ns"`
			P99Ns   int64 `json:"p99_ns"`
			P999Ns  int64 `json:"p999_ns"`
			P9999Ns int64 `json:"p9999_ns"`
		} `json:"vectors"`
	}
	readKit(t, "histogram-vectors.json", &kit)
	if kit.Histogram.Lowest != HistogramLowest || kit.Histogram.Highest != HistogramHighest ||
		kit.Histogram.Digits != HistogramDigits {
		t.Fatalf("the kit's histogram is %+v", kit.Histogram)
	}
	for _, v := range kit.Vectors {
		h := NewHistogram()
		for _, r := range v.Recorded {
			h.Record(r[0], r[1])
		}
		if r := v.RecordedRange; r != nil {
			for value := r.From; value <= r.To; value++ {
				h.Record(value, r.CountEach)
			}
		}
		got := [5]int64{h.ValueAtPercentile(50), h.ValueAtPercentile(90), h.ValueAtPercentile(99),
			h.ValueAtPercentile(99.9), h.ValueAtPercentile(99.99)}
		want := [5]int64{v.P50Ns, v.P90Ns, v.P99Ns, v.P999Ns, v.P9999Ns}
		if h.Count() != v.Count || got != want {
			t.Errorf("%s: count %d, percentiles %v; want %d, %v", v.About, h.Count(), got, v.Count, want)
		}
	}
}

func TestHistogramLimits(t *testing.T) {
	h := NewHistogram()
	h.Record(-5, 1)
	h.Record(HistogramHighest*2, 1)
	if h.ValueAtPercentile(0) != 0 || h.ValueAtPercentile(100) < HistogramHighest || h.Count() != 2 {
		t.Errorf("%d %d %d", h.ValueAtPercentile(0), h.ValueAtPercentile(100), h.Count())
	}
	other := NewHistogram()
	other.Add(h)
	h.Reset()
	if h.Count() != 0 || other.Count() != 2 || h.ValueAtPercentile(50) != 0 {
		t.Errorf("add/reset: %d %d", h.Count(), other.Count())
	}
}

func TestHistogramRoundedCount(t *testing.T) {
	// 7 values: the count at p90 is 6.3, which HdrHistogram 2.1.9 rounds to 6 (later versions take 7).
	h := NewHistogram()
	for i := int64(1); i <= 7; i++ {
		h.Record(i*1000, 1)
	}
	if got := h.ValueAtPercentile(90); got != 6003 {
		t.Errorf("p90 of 1000..7000 = %d, want 6003", got)
	}
}
