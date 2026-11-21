// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"math/big"
	"strconv"
	"strings"
)

// Group prints an integer grouped by thousands with ",", as Java's NumberFormat.getIntegerInstance(Locale.US).
func Group(value int64) string {
	digits := strconv.FormatInt(value, 10)
	sign := ""
	if digits[0] == '-' {
		sign, digits = "-", digits[1:]
	}
	var b strings.Builder
	b.WriteString(sign)
	for i, digit := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(digit)
	}
	return b.String()
}

// InstanceStatistics are the counts of an instance in an interval: the records and, on p, the sum over the batches
// of their records times their fill and publishing time in nanoseconds, on c the sum of the e2e lags in nanoseconds.
type InstanceStatistics struct {
	Records int64
	LagSum  int64
}

// muldiv returns a * b / c with integer division toward zero, without overflow in between; c is not 0.
func muldiv(a, b, c int64) int64 {
	var x big.Int
	x.Mul(big.NewInt(a), big.NewInt(b))
	return x.Quo(&x, big.NewInt(c)).Int64()
}

// Rate returns the records per second of an interval of intervalNs, 0 for an empty interval.
func Rate(records, intervalNs int64) int64 {
	if records == 0 || intervalNs <= 0 {
		return 0
	}
	return muldiv(records, 1000000000, intervalNs)
}

// Lag returns the mean lag in microseconds, lagSum / records / 1000 toward zero, 0 without records.
func Lag(lagSum, records int64) int64 {
	if records == 0 {
		return 0
	}
	return lagSum / records / 1000
}

// StatisticsBlock returns the statistics block of an interval without the last line break: the summary line of the
// side (Producer or Consumer) and a line per instance.
func StatisticsBlock(mode, name string, intervalNs int64, connections int, instances []InstanceStatistics) string {
	side := "Producer"
	if mode == Consumer {
		side = "Consumer"
	}
	var total InstanceStatistics
	for _, s := range instances {
		total.Records += s.Records
		total.LagSum += s.LagSum
	}
	var b strings.Builder
	b.WriteString("{*NetTest" + side + "Side-" + name + "*} data " + Group(Rate(total.Records, intervalNs)) +
		" rps, lag " + Group(Lag(total.LagSum, total.Records)) + " us, instances " + strconv.Itoa(len(instances)) +
		", connections " + strconv.Itoa(connections))
	for i, s := range instances {
		b.WriteString("\n    " + name + "." + strconv.Itoa(i+1) + " data " + Group(Rate(s.Records, intervalNs)) +
			" rps, lag " + Group(Lag(s.LagSum, s.Records)) + " us")
	}
	return b.String()
}

// RuntimeSample is a sample of the runtime; a value that the runtime does not give is negative.
type RuntimeSample struct {
	HeapUsedBytes      int64
	HeapCommittedBytes int64
	// GCCount and GCTimeMs are the sums of the collections of both kinds and of their time.
	GCCount   int64
	GCTimeMs  int64
	CPUTimeNs int64
	// GCCounts and GCTimesMs are the collections and their time by kind, young and full, for the summary.
	GCCounts, GCTimesMs [2]int64
}

// UnknownRuntimeSample returns a sample with no value given: every field, the sums and the values by kind, is -1.
func UnknownRuntimeSample() RuntimeSample {
	return RuntimeSample{HeapUsedBytes: -1, HeapCommittedBytes: -1, GCCount: -1, GCTimeMs: -1,
		GCCounts: [2]int64{-1, -1}, GCTimesMs: [2]int64{-1, -1}, CPUTimeNs: -1}
}

// NewSummaryRuntime returns the runtime of the summary from the samples at the start and at the end of the run, as
// the SDK C tool: the heap at the end, the growth of the collections and of their time over the run (null when a
// value is not given or a counter went back), and the mean CPU of the run.
func NewSummaryRuntime(first, last RuntimeSample, periodNs int64) SummaryRuntime {
	megabytes := func(bytes int64) *int64 { return IntIf(bytes >= 0, bytes>>20) }
	growth := func(from, to int64) (int64, bool) { return to - from, from >= 0 && to >= from }
	both := func(from, to [2]int64) *int64 {
		young, youngOK := growth(from[0], to[0])
		full, fullOK := growth(from[1], to[1])
		return IntIf(youngOK && fullOK, young+full)
	}
	r := SummaryRuntime{HeapUsedMb: megabytes(last.HeapUsedBytes), HeapCommittedMb: megabytes(last.HeapCommittedBytes),
		GCCount: both(first.GCCounts, last.GCCounts), GCTimeMs: both(first.GCTimesMs, last.GCTimesMs)}
	if cpu, ok := growth(first.CPUTimeNs, last.CPUTimeNs); ok && periodNs > 0 {
		r.CPUPercentMean = Int(cpu * 100 / periodNs)
	}
	return r
}

// RuntimeLine returns the runtime line of an interval from the samples at its start and its end.
func RuntimeLine(name string, intervalNs int64, previous, current RuntimeSample) string {
	mb := func(bytes int64) string {
		if bytes < 0 {
			return "n/a"
		}
		return Group(bytes>>20) + " MB"
	}
	delta := func(previous, current int64) (int64, bool) {
		return current - previous, previous >= 0 && current >= previous
	}
	gc := "n/a"
	count, countOK := delta(previous.GCCount, current.GCCount)
	time, timeOK := delta(previous.GCTimeMs, current.GCTimeMs)
	if countOK && timeOK {
		gc = Group(count) + " (" + Group(time) + " ms)"
	}
	cpu := "n/a"
	if ns, ok := delta(previous.CPUTimeNs, current.CPUTimeNs); ok && intervalNs > 0 {
		cpu = Group(muldiv(ns, 100, intervalNs)) + "%"
	}
	return "{*NetTestRuntime-" + name + "*} heap " + mb(current.HeapUsedBytes) + " of " +
		mb(current.HeapCommittedBytes) + ", gc " + gc + ", cpu " + cpu
}
