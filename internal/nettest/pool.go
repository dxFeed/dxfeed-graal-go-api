// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

// The constants of the random source of the pools.
const (
	golden     = 0x9E3779B97F4A7C15
	multiplier = 0x2545F4914F6CDD1D
)

// TableRows is the number of the rows of the price and size table of an instance.
const TableRows = 1024

// SequenceMask is the mask of the 22-bit sequence of an instance.
const SequenceMask = 1<<22 - 1

// Random is the xorshift64* source of an instance.
type Random struct{ state uint64 }

// Seed returns the seed of the instance (from 1).
func Seed(instance int) uint64 {
	return golden * uint64(instance)
}

// NewRandom returns the random source of the instance (from 1).
func NewRandom(instance int) *Random {
	state := Seed(instance)*golden + multiplier
	if state == 0 {
		state = multiplier
	}
	return &Random{state}
}

// State returns the current state.
func (r *Random) State() uint64 { return r.state }

// Next returns the next draw.
func (r *Random) Next() uint64 {
	s := r.state
	s ^= s >> 12
	s ^= s << 25
	s ^= s >> 27
	r.state = s
	return s * multiplier
}

// Row is a row of the table: a price and a size.
type Row struct{ Price, Size float64 }

// NewTable builds the table of an instance from its random source, two draws per row.
func NewTable(r *Random) []Row {
	table := make([]Row, TableRows)
	for i := range table {
		table[i].Price = float64(100+r.Next()%100000) / 100.0
		table[i].Size = float64(1 + r.Next()%1000)
	}
	return table
}

// Choose returns the symbol index and the row of an event from its draw.
func Choose(draw uint64, poolSize int) (symbol, row int) {
	return int((draw >> 32) % uint64(poolSize)), int((draw >> 8) % TableRows)
}

// PackTimeSequence packs a time in milliseconds and a sequence as the time-and-sequence of Trade and the index of
// TimeAndSale: the seconds in bits 32 to 63, the milliseconds in bits 22 to 31 and the sequence in bits 0 to 21.
func PackTimeSequence(timeMillis int64, sequence int) int64 {
	seconds := timeMillis / 1000
	millis := timeMillis % 1000
	if millis < 0 {
		seconds--
		millis += 1000
	}
	return seconds<<32 | millis<<22 | int64(sequence&SequenceMask)
}

// SplitTime splits a time in nanoseconds since the epoch into the milliseconds and the nanosecond part.
func SplitTime(epochNs int64) (timeMillis int64, nanoPart int32) {
	return epochNs / 1000000, int32(epochNs % 1000000)
}
