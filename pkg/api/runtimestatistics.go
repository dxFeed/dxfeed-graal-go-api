// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// RuntimeStatistics holds the runtime statistics of the Graal isolate of the SDK and of its process: the heap, the
// garbage collections, the threads, the CPU time of the process and the uptime (see [GetRuntimeStatistics]).
//
// Every value is -1 when the platform or the garbage collector does not give it ("not given"), otherwise it is 0 or
// more. The heap pool and GC values follow the Serial GC, the GC of the SDK builds; a build with another GC (G1) has
// no such pools and collectors, and these values are -1 there. AllocatedBytes is -1 with GraalVM 25.5.
//
// Each value is read on its own, so the statistics are not one atomic snapshot: a heap value and a GC count may come
// from moments around one collection. The counters are cumulative since the isolate started: compute the change
// between two calls. ProcessCPUTimeNs and AvailableProcessors belong to the whole process, the other values to the
// isolate.
type RuntimeStatistics struct {
	// HeapUsedBytes is the bytes of the heap of the isolate in use, garbage not yet collected included.
	HeapUsedBytes int64
	// HeapCommittedBytes is the bytes of the heap of the isolate taken from the system, HeapUsedBytes or more.
	HeapCommittedBytes int64
	// HeapMaxBytes is the bytes the heap of the isolate may grow to (-Xmx, by default a part of the physical memory).
	HeapMaxBytes int64
	// YoungGCCount is the incremental (young generation) collections of the isolate since it started; -1 on G1.
	YoungGCCount int64
	// YoungGCTimeMs is the total time of the incremental collections, in milliseconds; -1 on G1.
	YoungGCTimeMs int64
	// FullGCCount is the collections of the whole heap of the isolate since it started; -1 on G1.
	FullGCCount int64
	// FullGCTimeMs is the total time of the collections of the whole heap, in milliseconds; -1 on G1.
	FullGCTimeMs int64
	// ThreadCount is the live Java threads of the isolate: its own threads and the attached native threads.
	ThreadCount int64
	// PeakThreadCount is the peak of ThreadCount since the isolate started.
	PeakThreadCount int64
	// ProcessCPUTimeNs is the CPU time of the whole process, user plus system, in nanoseconds. The step is the one of
	// the system (on Windows 15.625 ms).
	ProcessCPUTimeNs int64
	// UptimeMs is the time since the isolate started, in milliseconds.
	UptimeMs int64
	// AllocatedBytes is the bytes allocated by all threads of the isolate since it started; -1 with GraalVM 25.5.
	AllocatedBytes int64
	// EdenUsedBytes is the bytes in use in the eden space of the young generation; -1 on G1.
	EdenUsedBytes int64
	// SurvivorUsedBytes is the bytes in use in the survivor space of the young generation; -1 on G1.
	SurvivorUsedBytes int64
	// OldUsedBytes is the bytes in use in the old generation: its growth across collections shows objects that stay
	// alive; -1 on G1.
	OldUsedBytes int64
	// DaemonThreadCount is the live daemon threads of the isolate; the attached native threads count as daemon
	// threads.
	DaemonThreadCount int64
	// TotalStartedThreadCount is the Java threads started or attached in the isolate since it started.
	TotalStartedThreadCount int64
	// AvailableProcessors is the processors available to the process, container limits considered: divide a CPU load
	// by it per core.
	AvailableProcessors int64
}

// GetRuntimeStatistics reads the runtime statistics of the Graal isolate of the SDK and of its process. A value that
// the platform or the garbage collector does not give is -1, and the call still succeeds (see [RuntimeStatistics]).
//
// Reading changes nothing. The first call initializes the sources and may take some milliseconds: call it once at
// start-up to take the cost out of a measurement. The call creates the isolate if it is not created yet.
func GetRuntimeStatistics() (RuntimeStatistics, error) {
	s, err := native.GetRuntimeStatistics()
	return RuntimeStatistics(s), err
}

// String returns the statistics in one line: every value with its name, -1 as is.
func (s RuntimeStatistics) String() string {
	return fmt.Sprintf("RuntimeStatistics{heapUsedBytes = %d, heapCommittedBytes = %d, heapMaxBytes = %d, "+
		"youngGcCount = %d, youngGcTimeMs = %d, fullGcCount = %d, fullGcTimeMs = %d, threadCount = %d, "+
		"peakThreadCount = %d, processCpuTimeNs = %d, uptimeMs = %d, allocatedBytes = %d, edenUsedBytes = %d, "+
		"survivorUsedBytes = %d, oldUsedBytes = %d, daemonThreadCount = %d, totalStartedThreadCount = %d, "+
		"availableProcessors = %d}",
		s.HeapUsedBytes, s.HeapCommittedBytes, s.HeapMaxBytes, s.YoungGCCount, s.YoungGCTimeMs, s.FullGCCount,
		s.FullGCTimeMs, s.ThreadCount, s.PeakThreadCount, s.ProcessCPUTimeNs, s.UptimeMs, s.AllocatedBytes,
		s.EdenUsedBytes, s.SurvivorUsedBytes, s.OldUsedBytes, s.DaemonThreadCount, s.TotalStartedThreadCount,
		s.AvailableProcessors)
}
