package native

import (
	"bytes"
	"log"
	"strings"
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
