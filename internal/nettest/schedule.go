// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"sync"
	"sync/atomic"
	"time"
)

// Schedule is the rate limit of all instances together: batch k, counted over all instances, is due at
// T(k) = t0 + k * batch / rate, where t0 is the time of the first grant.
type Schedule struct {
	nsPerBatch float64
	next       atomic.Int64
	once       sync.Once
	t0         time.Time
}

// NewSchedule returns the schedule of batches of batch events at rate events per second.
func NewSchedule(rate float64, batch int) *Schedule {
	return &Schedule{nsPerBatch: float64(batch) * 1e9 / rate}
}

// Grant returns the time when the next batch is due; the first call takes the current time as t0.
func (s *Schedule) Grant() time.Time {
	s.once.Do(func() { s.t0 = time.Now() })
	k := s.next.Add(1) - 1
	return s.t0.Add(time.Duration(float64(k) * s.nsPerBatch))
}
