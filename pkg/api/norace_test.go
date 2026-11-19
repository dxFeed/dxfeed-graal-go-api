// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build !race

package api

// raceEnabled reports whether the tests run with the race detector, which allows at most 8128 live goroutines.
const raceEnabled = false
