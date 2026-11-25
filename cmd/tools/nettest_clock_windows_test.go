// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package main

import (
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// preciseNanos reads GetSystemTimePreciseAsFileTime, the clock of the C++ and Rust tools, on its own.
func preciseNanos() int64 {
	var ft syscall.Filetime
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemTimePreciseAsFileTime")
	_, _, _ = proc.Call(uintptr(unsafe.Pointer(&ft)))
	return ft.Nanoseconds()
}

// TestNetTestClockWindows checks that the event clock stays within 50 µs of the system clock over 20 s.
func TestNetTestClockWindows(t *testing.T) {
	if testing.Short() {
		t.Skip("20 s")
	}
	t.Parallel()
	const tolerance = 50 * time.Microsecond
	r := &netTestRun{start: time.Now()}
	// time.Now() steps by up to 1 ms on Windows: a coarse check of the epoch and the units of both precise reads.
	if diff := time.Since(r.now()); diff < -5*time.Millisecond || diff > 5*time.Millisecond {
		t.Fatalf("event clock %v off time.Now()", -diff)
	}
	for end := time.Now().Add(20 * time.Second); time.Now().Before(end); time.Sleep(100 * time.Millisecond) {
		before := preciseNanos()
		now := r.now().UnixNano()
		after := preciseNanos()
		if now < before-tolerance.Nanoseconds() || now > after+tolerance.Nanoseconds() {
			t.Fatalf("event clock %d us off the system clock", (now-(before+after)/2)/1000)
		}
	}
}
