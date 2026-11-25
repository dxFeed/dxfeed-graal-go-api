// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package main

import "time"

// wallNanos returns the system time in nanoseconds since the Unix epoch.
func wallNanos() int64 { return time.Now().UnixNano() }
