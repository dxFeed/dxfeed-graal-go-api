// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// startQuotePublisher starts a publisher on a free local port that publishes a quote of the symbol every 50 ms until
// the test ends, and returns its address.
func startQuotePublisher(t *testing.T, symbol string) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	_ = listener.Close()
	server, err := api.NewEndpoint(api.Publisher)
	if err != nil {
		t.Fatal(err)
	}
	if err := server.Connect(":" + port); err != nil {
		t.Fatal(err)
	}
	publisher, err := server.GetPublisher()
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		for price := 1.0; ; price++ {
			q := quote.NewQuote(symbol)
			q.SetBidPrice(price)
			q.SetBidTime(time.Now().UnixMilli())
			_ = publisher.Publish([]any{q})
			select {
			case <-done:
				return
			case <-time.After(50 * time.Millisecond):
			}
		}
	}()
	t.Cleanup(func() {
		close(done)
		<-stopped
		_ = server.Close()
	})
	return "127.0.0.1:" + port
}

func TestDiagnoseData(t *testing.T) {
	address := startQuotePublisher(t, "DIAG")
	report := filepath.Join(t.TempDir(), "report.txt")
	status, out := runTools(t, "diagnose "+address+"[user=me,password=s3cr3t] Quote DIAG --duration 2 -o "+report)
	for _, want := range []string{
		"=== Address: " + address + "[user=***,password=***]\n",
		"  " + address + ": resolved to 1 IP address(es)\n",
		"    127.0.0.1: TCP connected in ",
		"Connecting -> Connected\n",
		"  Quote: ",
		"    flags: ",
		"  Quote DIAG: ",
		"Verdict: DATA\n",
		"The report is written to " + report + "\n",
	} {
		if want == "    flags: " {
			if strings.Contains(out, want) {
				t.Errorf("flags of quotes are printed:\n%s", out)
			}
			continue
		}
		if !strings.Contains(out, want) {
			t.Errorf("no %q in the output:\n%s", want, out)
		}
	}
	if status != 0 || strings.Contains(out, "s3cr3t") {
		t.Errorf("status %d, want 0, and no password:\n%s", status, out)
	}
	text, err := os.ReadFile(report)
	if err != nil || !strings.Contains(string(text), "Verdict: DATA\n") ||
		strings.Contains(string(text), "The report is written") || strings.Contains(string(text), "s3cr3t") {
		t.Errorf("report %v:\n%s", err, text)
	}
}

func TestDiagnosePartialDataAndCompare(t *testing.T) {
	address := startQuotePublisher(t, "DIAG")
	status, out := runTools(t, "diagnose "+address+" Quote DIAG,NONE --duration 2 --compare 127.0.0.1:1")
	for _, want := range []string{
		"Verdict: PARTIAL DATA\n", "  no events: Quote NONE\n", "Verdict: CANNOT CONNECT\n", "=== Comparison\n",
		"  only from " + address + ": Quote DIAG\n", "  The verdicts differ.\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in the output:\n%s", want, out)
		}
	}
	if status != 4 {
		t.Errorf("status %d, want 4, the worst verdict:\n%s", status, out)
	}
}

func TestDiagnoseArgumentErrors(t *testing.T) {
	for _, test := range []struct{ args, want string }{
		{"diagnose a Quote AAPL --duration 0", "diagnose: the duration must be a whole number of seconds from 1 to 3600: 0"},
		{"diagnose a Quote AAPL --duration 3601", "diagnose: the duration must be a whole number of seconds from 1 to 3600: 3601"},
		{"diagnose a Quote", "diagnose: the symbols argument is missing"},
		{"diagnose a Quote AAPL -o", "diagnose: option -o needs a value"},
	} {
		status, out := runTools(t, test.args)
		if status != 2 || !strings.HasPrefix(out, test.want+"\n\n") {
			t.Errorf("tools %s: status %d, want 2 and %q:\n%s", test.args, status, test.want, out)
		}
	}
}
