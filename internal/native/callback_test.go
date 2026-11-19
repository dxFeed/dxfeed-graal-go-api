// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"testing"
)

func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(previous) })
	return &buf
}

func TestReportUnmappedClazzOncePerClass(t *testing.T) {
	logged := captureLog(t)
	clazzes := sdkEventClazzes()
	// The reported classes are global: forget the ones of this test, so it does not depend on what ran before.
	forget := func() {
		for _, clazz := range []int32{clazzes["MESSAGE"], clazzes["SERIES"], -7} {
			reportedClazzes.Delete(clazz)
		}
	}
	forget()
	t.Cleanup(forget)
	reportUnmappedClazz(clazzes["MESSAGE"])
	reportUnmappedClazz(clazzes["MESSAGE"])
	reportUnmappedClazz(clazzes["SERIES"])
	reportUnmappedClazz(-7)

	out := logged.String()
	if n := strings.Count(out, "MESSAGE"); n != 1 {
		t.Errorf("MESSAGE reported %d times, want 1; log: %q", n, out)
	}
	if !strings.Contains(out, "SERIES") || !strings.Contains(out, "unknown (-7)") {
		t.Errorf("SERIES or the unknown class is not reported; log: %q", out)
	}
}

func TestRecoverCallback(t *testing.T) {
	logged := captureLog(t)
	func() {
		defer recoverCallback("test listener")
		panic("boom")
	}()
	if out := logged.String(); !strings.Contains(out, "test listener panicked: boom") {
		t.Fatalf("the panic is not logged; log: %q", out)
	}
}

func TestReportGoesToLogListener(t *testing.T) {
	logged := captureLog(t)
	clazz := sdkEventClazzes()["MESSAGE"]
	reportedClazzes.Delete(clazz)
	t.Cleanup(func() { reportedClazzes.Delete(clazz) })

	// The SDK listener is process-wide: keep only the records of the Go API, the SDK may log from its threads.
	var mu sync.Mutex
	var first, second []LogRecord
	collect := func(to *[]LogRecord) func(LogRecord) {
		return func(r LogRecord) {
			if r.LoggerName == GoLoggerName {
				mu.Lock()
				defer mu.Unlock()
				*to = append(*to, r)
			}
		}
	}
	if err := SetLogListener(collect(&first)); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = SetLogListener(nil) })
	reportUnmappedClazz(clazz)
	if len(first) != 1 || first[0].LoggerName != GoLoggerName || first[0].Level != LogLevelWarn ||
		!strings.Contains(first[0].Message, "MESSAGE") {
		t.Fatalf("records = %+v, want one warning about MESSAGE", first)
	}
	// A replaced listener gets nothing more.
	if err := SetLogListener(collect(&second)); err != nil {
		t.Fatal(err)
	}
	logGo(LogLevelError, "second %d", 2)
	if len(first) != 1 || len(second) != 1 || second[0].Message != "second 2" {
		t.Errorf("first = %+v, second = %+v", first, second)
	}
	// Without a listener the standard log gets the message again.
	if err := SetLogListener(nil); err != nil {
		t.Fatal(err)
	}
	logGo(LogLevelInfo, "to the log")
	if !strings.Contains(logged.String(), "dxfeed: to the log") {
		t.Errorf("log = %q", logged.String())
	}
	if strings.Contains(logged.String(), "MESSAGE") {
		t.Errorf("the reported record also reached the log: %q", logged.String())
	}
}

func TestPanickingLogListenerFallsBackToLog(t *testing.T) {
	logged := captureLog(t)
	if err := SetLogListener(func(r LogRecord) {
		if r.LoggerName == GoLoggerName {
			panic("boom")
		}
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = SetLogListener(nil) })
	logGo(LogLevelWarn, "kept")
	if out := logged.String(); !strings.Contains(out, "boom") || !strings.Contains(out, "dxfeed: kept") {
		t.Errorf("log = %q", out)
	}
}
