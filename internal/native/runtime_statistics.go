// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
*/
import "C"

// RuntimeStatistics holds the values of dxfg_v4_runtime_stats_t; pkg/api converts it to its own struct of the same
// fields.
type RuntimeStatistics struct {
	HeapUsedBytes           int64
	HeapCommittedBytes      int64
	HeapMaxBytes            int64
	YoungGCCount            int64
	YoungGCTimeMs           int64
	FullGCCount             int64
	FullGCTimeMs            int64
	ThreadCount             int64
	PeakThreadCount         int64
	ProcessCPUTimeNs        int64
	UptimeMs                int64
	AllocatedBytes          int64
	EdenUsedBytes           int64
	SurvivorUsedBytes       int64
	OldUsedBytes            int64
	DaemonThreadCount       int64
	TotalStartedThreadCount int64
	AvailableProcessors     int64
}

// GetRuntimeStatistics reads the runtime statistics of the isolate and its process.
func GetRuntimeStatistics() (RuntimeStatistics, error) {
	var s C.dxfg_v4_runtime_stats_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_system_get_runtime_statistics(thread.ptr, &s))
	})
	if err != nil {
		return RuntimeStatistics{}, err
	}
	return RuntimeStatistics{
		HeapUsedBytes:           int64(s.heap_used_bytes),
		HeapCommittedBytes:      int64(s.heap_committed_bytes),
		HeapMaxBytes:            int64(s.heap_max_bytes),
		YoungGCCount:            int64(s.young_gc_count),
		YoungGCTimeMs:           int64(s.young_gc_time_ms),
		FullGCCount:             int64(s.full_gc_count),
		FullGCTimeMs:            int64(s.full_gc_time_ms),
		ThreadCount:             int64(s.thread_count),
		PeakThreadCount:         int64(s.peak_thread_count),
		ProcessCPUTimeNs:        int64(s.process_cpu_time_ns),
		UptimeMs:                int64(s.uptime_ms),
		AllocatedBytes:          int64(s.allocated_bytes),
		EdenUsedBytes:           int64(s.eden_used_bytes),
		SurvivorUsedBytes:       int64(s.survivor_used_bytes),
		OldUsedBytes:            int64(s.old_used_bytes),
		DaemonThreadCount:       int64(s.daemon_thread_count),
		TotalStartedThreadCount: int64(s.total_started_thread_count),
		AvailableProcessors:     int64(s.available_processors),
	}, nil
}
