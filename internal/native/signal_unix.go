// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build cgo && unix

package native

import (
	"os"
	"os/signal"
	"syscall"
)

// unhandledSignals receives SIGPIPE and SIGXFSZ; nothing reads it, and the signals that do not fit are dropped.
var unhandledSignals = make(chan os.Signal, 1)

// handleWriteSignals keeps SIGPIPE and SIGXFSZ from ending the process, as a JVM does: a write to a closed socket or
// pipe then fails with EPIPE, and a write over the file size limit with EFBIG. The SDK handles these errors, but it
// writes on its own threads, which Go did not create: Go passes such a signal of a thread like that to the default
// action, which ends the process, unless the program is notified of the signal (os/signal, "SIGPIPE").
//
// signal.Notify is used rather than signal.Ignore: an ignored signal stays ignored in the programs that the
// application starts. A side effect: a write of the program to a closed standard output or error fails with EPIPE
// instead of ending the program.
func handleWriteSignals() {
	signal.Notify(unhandledSignals, syscall.SIGPIPE, syscall.SIGXFSZ)
}
