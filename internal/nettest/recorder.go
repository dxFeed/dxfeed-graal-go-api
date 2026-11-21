// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"sync"
	"time"
)

// LagBuckets are the upper bounds of the buckets of nettest_e2e_lag_seconds in nanoseconds, without +Inf.
var LagBuckets = []int64{50000, 100000, 250000, 500000, 1000000, 2500000, 5000000, 10000000, 25000000, 50000000,
	100000000, 250000000, 1000000000}

// LagInterval is the e2e lag of the consumer in an interval.
type LagInterval struct {
	Histogram *Histogram
	// Count, Sum and Max are those of the lags with their sign; Max is meaningful only when Count is not 0.
	Count, Sum, Max, Negative int64
}

// LagRecorder records the e2e lags of the consumer: those of the current interval for the summary, the cumulative
// buckets for the metrics and the symbols seen. It is safe for concurrent use.
type LagRecorder struct {
	mu       sync.Mutex
	interval LagInterval
	// buckets counts the lags of each bucket of LagBuckets and, last, those above them; sum, count and negative are
	// cumulative.
	buckets              []int64
	sum, count, negative int64
	symbols              map[string]struct{}
}

// NewLagRecorder returns an empty recorder.
func NewLagRecorder() *LagRecorder {
	return &LagRecorder{interval: LagInterval{Histogram: NewHistogram()}, buckets: make([]int64, len(LagBuckets)+1),
		symbols: map[string]struct{}{}}
}

// Record records the lags of the events of a call of the listener and their symbols.
func (l *LagRecorder) Record(lags []int64, symbols []string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, lag := range lags {
		in := &l.interval
		if in.Count == 0 || lag > in.Max {
			in.Max = lag
		}
		in.Count++
		in.Sum += lag
		in.Histogram.Record(lag, 1)
		l.count++
		l.sum += lag
		if lag < 0 {
			in.Negative++
			l.negative++
		}
		bucket := len(LagBuckets)
		for i, bound := range LagBuckets {
			if lag <= bound {
				bucket = i
				break
			}
		}
		l.buckets[bucket]++
	}
	for _, symbol := range symbols {
		if _, ok := l.symbols[symbol]; !ok {
			l.symbols[symbol] = struct{}{}
		}
	}
}

// TakeInterval returns the lags of the interval and starts a new one.
func (l *LagRecorder) TakeInterval() LagInterval {
	l.mu.Lock()
	defer l.mu.Unlock()
	result := l.interval
	l.interval = LagInterval{Histogram: NewHistogram()}
	return result
}

// LagMetrics are the cumulative e2e lags for the metrics.
type LagMetrics struct {
	// Buckets are the cumulative counts of the buckets of LagBuckets and +Inf.
	Buckets                []int64
	SumNs, Count, Negative int64
}

// Metrics returns the cumulative lags.
func (l *LagRecorder) Metrics() LagMetrics {
	l.mu.Lock()
	defer l.mu.Unlock()
	m := LagMetrics{Buckets: make([]int64, len(l.buckets)), SumNs: l.sum, Count: l.count, Negative: l.negative}
	var total int64
	for i, count := range l.buckets {
		total += count
		m.Buckets[i] = total
	}
	return m
}

// SymbolsSeen returns the number of the distinct symbols received.
func (l *LagRecorder) SymbolsSeen() int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return int64(len(l.symbols))
}

// IntervalCounts are the counts of an interval of the run.
type IntervalCounts struct {
	Start, End time.Time
	// Warmup tells whether the interval belongs to the warm-up.
	Warmup      bool
	Records     int64
	Connections int
	// Lag is the e2e lag of the consumer, nil on the publisher.
	Lag *LagInterval
	// PublishNs is the time of the fills and the publishing calls of the publisher, weighted by their records.
	PublishNs int64
	// Slip tells whether the publisher has a schedule; SlipMax and SlipEnd are its slip in the interval and at its end.
	Slip             bool
	SlipMax, SlipEnd int64
}

// SummaryBuilder collects the intervals of a run and their totals for the summary file. It is not safe for concurrent
// use.
type SummaryBuilder struct {
	Mode      string
	Intervals []SummaryInterval
	Records   int64
	// The e2e lag after the warm-up.
	lag                                   *Histogram
	lagCount, lagSum, lagMax, lagNegative int64
	// The schedule slip: its maximum over all intervals and over the measured ones, its value at the end and the slip
	// at the ends of the measured intervals.
	slip                              bool
	slipMax, slipMaxMeasured, slipEnd int64
	slipEnds                          []slipPoint
}

type slipPoint struct{ timeMs, slipNs int64 }

// NewSummaryBuilder returns a builder of the mode.
func NewSummaryBuilder(mode string) *SummaryBuilder {
	return &SummaryBuilder{Mode: mode, lag: NewHistogram(), Intervals: []SummaryInterval{}}
}

// Add adds an interval.
func (b *SummaryBuilder) Add(c IntervalCounts) {
	length := c.End.Sub(c.Start).Nanoseconds()
	phase := "measure"
	if c.Warmup {
		phase = "warmup"
	}
	in := SummaryInterval{StartEpochMs: c.Start.UnixMilli(), EndEpochMs: c.End.UnixMilli(), Phase: phase,
		Records: c.Records, Rate: Rate(c.Records, length), Connections: c.Connections}
	b.Records += c.Records
	if lag := c.Lag; lag != nil {
		if lag.Count > 0 {
			in.LagMeanNs = Int(lag.Sum / lag.Count)
			in.LagP50Ns = Int(lag.Histogram.ValueAtPercentile(50))
			in.LagP99Ns = Int(lag.Histogram.ValueAtPercentile(99))
			in.LagMaxNs = Int(lag.Max)
		}
		if !c.Warmup && lag.Count > 0 {
			if b.lagCount == 0 || lag.Max > b.lagMax {
				b.lagMax = lag.Max
			}
			b.lag.Add(lag.Histogram)
			b.lagCount += lag.Count
			b.lagSum += lag.Sum
			b.lagNegative += lag.Negative
		}
	} else if c.Records > 0 {
		in.PublishNsPerEvent = Int(c.PublishNs / c.Records)
	}
	if c.Slip {
		b.slip = true
		in.SlipMaxNs, in.SlipEndNs = Int(c.SlipMax), Int(c.SlipEnd)
		b.slipMax = max(b.slipMax, c.SlipMax)
		b.slipEnd = c.SlipEnd
		if !c.Warmup {
			b.slipMaxMeasured = max(b.slipMaxMeasured, c.SlipMax)
			b.slipEnds = append(b.slipEnds, slipPoint{c.End.UnixMilli(), c.SlipEnd})
		}
	}
	b.Intervals = append(b.Intervals, in)
}

// Lag returns the e2e lag after the warm-up.
func (b *SummaryBuilder) Lag() *SummaryLag {
	return NewLag(b.lag, b.lagNegative, b.lagSum, b.lagMax)
}

// Slip returns the schedule slip, nil without a schedule, as the SDK C tool computes it: the maximum of the measured
// intervals (of all intervals when none is measured), and the growth as the slope of the slip at the ends of the
// measured intervals over their end times in milliseconds, truncated, null with fewer than two of them.
func (b *SummaryBuilder) Slip() *SummarySlip {
	if !b.slip {
		return nil
	}
	s := &SummarySlip{MaxNs: b.slipMax, EndNs: b.slipEnd}
	if len(b.slipEnds) > 0 {
		s.MaxNs = b.slipMaxMeasured
	}
	if n := len(b.slipEnds); n >= 2 {
		first, last := b.slipEnds[0], b.slipEnds[n-1]
		if millis := last.timeMs - first.timeMs; millis > 0 {
			s.GrowthNsPerS = Int((last.slipNs - first.slipNs) * 1000 / millis)
		}
	}
	return s
}
