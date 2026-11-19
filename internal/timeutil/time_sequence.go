// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package timeutil

import "fmt"

// MaxSequence is the largest sequence number that fits in the low 22 bits of a time-sequence value.
const MaxSequence = (1 << 22) - 1

// TimeOfTimeSequence returns the time in milliseconds packed in a time-sequence value: seconds in the high 32 bits,
// milliseconds in bits 22-31.
func TimeOfTimeSequence(timeSequence int64) int64 {
	return (timeSequence>>32)*1000 + (timeSequence>>22)&0x3ff
}

// SequenceOfTimeSequence returns the sequence number packed in the low 22 bits of a time-sequence value.
func SequenceOfTimeSequence(timeSequence int64) int32 {
	return int32(timeSequence & MaxSequence)
}

// WithTime returns the time-sequence value with the time replaced and the sequence kept.
func WithTime(timeSequence int64, timeMillis int64) int64 {
	return GetSecondsFromTime(timeMillis)<<32 | int64(GetMillisFromTime(timeMillis))<<22 |
		int64(SequenceOfTimeSequence(timeSequence))
}

// WithSequence returns the time-sequence value with the sequence replaced, or an error if the sequence is out of
// range.
func WithSequence(timeSequence int64, sequence int32) (int64, error) {
	if sequence < 0 || sequence > MaxSequence {
		return timeSequence, fmt.Errorf("sequence %d is out of range [0, %d]", sequence, MaxSequence)
	}
	return timeSequence&^MaxSequence | int64(sequence), nil
}
