// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// toolsArgsEnv makes the test binary run the tools with the arguments it holds (separated by spaces) instead of the
// tests, as main does but without the system properties, so a test can check the exit status and output of a run.
const toolsArgsEnv = "DXFEED_TOOLS_TEST_ARGS"

func TestMain(m *testing.M) {
	if args, ok := os.LookupEnv(toolsArgsEnv); ok {
		args := append([]string{"tools"}, strings.Fields(args)...)
		if tool, status := selectTool(args, os.Stdout, os.Stderr); tool == nil {
			os.Exit(status)
		} else {
			os.Exit(tool.Run(args[2:]))
		}
	}
	os.Exit(m.Run())
}

// runTools runs the tools with args in a child process and returns its exit status and output.
func runTools(t *testing.T, args string) (int, string) {
	t.Helper()
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), toolsArgsEnv+"="+args)
	out, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	switch {
	case err == nil:
		return 0, string(out)
	case errors.As(err, &exitErr):
		return exitErr.ExitCode(), string(out)
	default:
		t.Fatalf("tools %s: %v", args, err)
		return 0, ""
	}
}

func TestSelectToolUsage(t *testing.T) {
	for _, args := range [][]string{{"tools"}, {"tools", "--help"}, {"tools", "-help"}, {"tools", "-h"},
		{"tools", "help"}, {"tools", "HELP"}} {
		var stdout, stderr bytes.Buffer
		if tool, status := selectTool(args, &stdout, &stderr); tool != nil || status != 0 {
			t.Errorf("%q: tool %v, status %d; want none, 0", args, tool, status)
		}
		if stderr.Len() != 0 {
			t.Errorf("%q: stderr %q, want none", args, stderr.String())
		}
		checkUsage(t, args, stdout.String())
	}
}

func TestSelectToolUnknown(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if tool, status := selectTool([]string{"tools", "Unknown"}, &stdout, &stderr); tool != nil || status != 2 {
		t.Errorf("tool %v, status %d; want none, 2", tool, status)
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout %q, want none", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Unknown tool: Unknown") {
		t.Errorf("stderr %q does not name the tool", stderr.String())
	}
	checkUsage(t, []string{"tools", "Unknown"}, stderr.String())
}

// checkUsage checks that the usage lists every tool once, in the order of the names.
func checkUsage(t *testing.T, args []string, usage string) {
	t.Helper()
	last := -1
	for _, name := range []string{"connect", "dump", "latencytest", "perftest"} {
		line := "\t" + name + " "
		if strings.Count(usage, line) != 1 {
			t.Errorf("%q: usage lists %s %d times, want once:\n%s", args, name, strings.Count(usage, line), usage)
			continue
		}
		index := strings.Index(usage, line)
		if index < last {
			t.Errorf("%q: %s is out of order:\n%s", args, name, usage)
		}
		last = index
	}
}

func TestRunExitStatus(t *testing.T) {
	for _, test := range []struct {
		args   string
		status int
	}{{"--help", 0}, {"-h", 0}, {"help", 0}, {"unknown", 2}} {
		status, out := runTools(t, test.args)
		if status != test.status {
			t.Errorf("tools %s: status %d, want %d:\n%s", test.args, status, test.status, out)
		}
		if !strings.Contains(out, "Usage: tools <tool>") {
			t.Errorf("tools %s: no usage:\n%s", test.args, out)
		}
	}
}

// TestToolArgumentErrors checks that a wrong argument of a tool prints the error and the help of the tool and exits
// with the status 2, without a panic and before the tool connects anywhere.
func TestToolArgumentErrors(t *testing.T) {
	for _, test := range []struct{ args, want string }{
		{"connect a Quote AAPL -p", "connect: option -p needs a value"},
		{"connect a Quote AAPL --bogus", "connect: unknown option --bogus"},
		{"connect a Quote AAPL -p x", `connect: option --properties: "x" is not name=value`},
		{"dump a Quote AAPL -t", "dump: option -t needs a value"},
		{"dump a Quote AAPL extra", "dump: unexpected argument extra"},
		{"perftest a Quote AAPL -q", "perftest: unknown option -q"},
		{"latencytest a TimeAndSale AAPL --ignore-exchanges", "latency test: option --ignore-exchanges needs a value"},
	} {
		status, out := runTools(t, test.args)
		if status != 2 {
			t.Errorf("tools %s: status %d, want 2:\n%s", test.args, status, out)
		}
		if !strings.HasPrefix(out, test.want+"\n\n") || !strings.Contains(out, "\nUsage:\n") || strings.Contains(out, "panic") {
			t.Errorf("tools %s: output does not start with %q and the help:\n%s", test.args, test.want, out)
		}
	}
}

// TestToolRunError checks that an error of a run is printed with a newline and gives the exit status 1.
func TestToolRunError(t *testing.T) {
	status, out := runTools(t, "connect localhost:notaport Quote AAPL -q")
	if status != 1 {
		t.Errorf("status %d, want 1:\n%s", status, out)
	}
	// The output starts with the log of the SDK.
	if !strings.Contains(out, "\nError during connect: Connect to localhost:notaport: ") || !strings.HasSuffix(out, "\n") {
		t.Errorf("output does not end with the error line:\n%s", out)
	}
}

func TestLatencyTestHelp(t *testing.T) {
	status, out := runTools(t, "latencytest")
	if status != 0 {
		t.Errorf("status %d, want 0:\n%s", status, out)
	}
	if !strings.Contains(out, "calculates their latency") || strings.Contains(out, "performance counters") {
		t.Errorf("the help is not of the latency test:\n%s", out)
	}
}
