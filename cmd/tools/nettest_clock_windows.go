// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

var getSystemTimePreciseAsFileTime = syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemTimePreciseAsFileTime")

// wallNanos returns the system time in nanoseconds since the Unix epoch. On Windows it reads
// GetSystemTimePreciseAsFileTime, as the C++ and Rust tools do: time.Now() there steps by up to 1 ms, and the
// monotonic clock of Go drifts from the system time.
func wallNanos() int64 {
	var ft syscall.Filetime
	syscall.SyscallN(getSystemTimePreciseAsFileTime.Addr(), uintptr(unsafe.Pointer(&ft)))
	return ft.Nanoseconds()
}
