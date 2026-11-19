// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build cgo && !unix

package native

// handleWriteSignals does nothing: the system has no SIGPIPE or SIGXFSZ.
func handleWriteSignals() {}
