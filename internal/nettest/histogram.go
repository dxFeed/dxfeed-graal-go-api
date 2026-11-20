// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"math"
	"math/bits"
)

// The parameters of the e2e lag histogram of the specification: the lowest value 1 ns, the highest 3600 s and 3
// significant digits.
const (
	HistogramLowest  = 1
	HistogramHighest = 3600000000000
	HistogramDigits  = 3
)

// The layout of the histogram for these parameters, as HdrHistogram computes it: sub-buckets of 2048 values (the
// smallest power of two above 2 * 10^3), a unit of 1 and as many buckets as cover the highest value.
const (
	subBucketCountMagnitude     = 11
	subBucketHalfCountMagnitude = subBucketCountMagnitude - 1
	subBucketCount              = 1 << subBucketCountMagnitude
	subBucketHalfCount          = subBucketCount / 2
	subBucketMask               = subBucketCount - 1
	leadingZeroCountBase        = 64 - subBucketCountMagnitude
)

// Histogram is the e2e lag histogram: an HdrHistogram whose percentiles are those of Java's HdrHistogram 2.1
// (getValueAtPercentile, the highest value equivalent to the bucket of the percentile). It is not safe for concurrent
// use.
type Histogram struct {
	counts []int64
	total  int64
}

// NewHistogram returns an empty histogram.
func NewHistogram() *Histogram {
	buckets := 1
	for smallestUntrackable := int64(subBucketCount); smallestUntrackable <= HistogramHighest; {
		if smallestUntrackable > math.MaxInt64/2 {
			buckets++
			break
		}
		smallestUntrackable <<= 1
		buckets++
	}
	return &Histogram{counts: make([]int64, (buckets+1)*subBucketHalfCount)}
}

func bucketIndex(value int64) int {
	return leadingZeroCountBase - bits.LeadingZeros64(uint64(value|subBucketMask))
}

func countsIndex(value int64) int {
	bucket := bucketIndex(value)
	subBucket := int(value >> bucket)
	return (bucket+1)<<subBucketHalfCountMagnitude + subBucket - subBucketHalfCount
}

func valueFromIndex(index int) int64 {
	bucket := index>>subBucketHalfCountMagnitude - 1
	subBucket := index&(subBucketHalfCount-1) + subBucketHalfCount
	if bucket < 0 {
		subBucket -= subBucketHalfCount
		bucket = 0
	}
	return int64(subBucket) << bucket
}

// highestEquivalentValue returns the highest value that falls into the bucket of the value.
func highestEquivalentValue(value int64) int64 {
	bucket := bucketIndex(value)
	subBucket := int(value >> bucket)
	lowest := int64(subBucket) << bucket
	sizeShift := bucket
	if subBucket >= subBucketCount {
		sizeShift++
	}
	return lowest + int64(1)<<sizeShift - 1
}

// Record records a value count times: a negative value as 0, a value above the highest as the highest.
func (h *Histogram) Record(value, count int64) {
	if value < 0 {
		value = 0
	}
	if value > HistogramHighest {
		value = HistogramHighest
	}
	h.counts[countsIndex(value)] += count
	h.total += count
}

// Add adds the counts of another histogram.
func (h *Histogram) Add(other *Histogram) {
	for i, count := range other.counts {
		h.counts[i] += count
	}
	h.total += other.total
}

// Reset empties the histogram.
func (h *Histogram) Reset() {
	clear(h.counts)
	h.total = 0
}

// Count returns the number of the recorded values.
func (h *Histogram) Count() int64 { return h.total }

// ValueAtPercentile returns the value at the percentile (0 to 100), 0 when the histogram is empty.
func (h *Histogram) ValueAtPercentile(percentile float64) int64 {
	// As HdrHistogram 2.1.9: the count at the percentile rounded to the nearest, at least 1 (later versions take the
	// ceiling one ulp below the percentile, which differs when the fraction of the count is below 0.5).
	requested := math.Max(math.Min(percentile, 100), 0)
	countAtPercentile := max(int64(requested/100*float64(h.total)+0.5), 1)
	var total int64
	for i, count := range h.counts {
		total += count
		if total >= countAtPercentile {
			value := valueFromIndex(i)
			if percentile == 0 {
				return value
			}
			return highestEquivalentValue(value)
		}
	}
	return 0
}
