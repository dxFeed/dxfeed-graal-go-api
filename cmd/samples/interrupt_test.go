// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package samples_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// untilInterrupt are the samples that run until Ctrl+C.
var untilInterrupt = []string{"quick-start", "auth", "print-quote-events", "print-indexed-events", "print-candles",
	"feed-ipf-connect", "ondemand-replay"}

// buildAll builds the samples into dir with one go build and returns the paths of the programs by name.
func buildAll(t *testing.T, dir string, names ...string) map[string]string {
	t.Helper()
	args := []string{"build", "-o", dir + string(filepath.Separator)}
	programs := make(map[string]string, len(names))
	for _, name := range names {
		args = append(args, "./"+name)
		programs[name] = filepath.Join(dir, name)
		if runtime.GOOS == "windows" {
			programs[name] += ".exe"
		}
	}
	if out, err := exec.Command("go", args...).CombinedOutput(); err != nil {
		t.Fatalf("go %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return programs
}

// TestSamplesHelp checks that -h prints the usage and exits with status 0 instead of starting the sample.
func TestSamplesHelp(t *testing.T) {
	programs := buildAll(t, t.TempDir(), untilInterrupt...)
	for _, name := range untilInterrupt {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		out, err := exec.CommandContext(ctx, programs[name], "-h").CombinedOutput()
		cancel()
		if err != nil || !strings.Contains(string(out), "-address") {
			t.Errorf("%s -h: %v\n%s", name, err, out)
		}
	}
}

// TestSamplesExitOnInterrupt connects the samples to a local publisher, interrupts them and checks that they exit
// with status 0. Windows has no interrupt signal to send to a process.
func TestSamplesExitOnInterrupt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.Interrupt cannot be sent to a process on Windows")
	}
	dir := t.TempDir()
	programs := buildAll(t, dir, untilInterrupt...)
	p := startPublisher(t, dir)
	address := "127.0.0.1:" + p.port
	ipf := filepath.Join(dir, "profiles.ipf")
	if err := os.WriteFile(ipf, []byte("#STOCK::=TYPE,SYMBOL,DESCRIPTION\nSTOCK,A:TEST,Test symbol\n##COMPLETE\n"),
		0o600); err != nil {
		t.Fatal(err)
	}
	args := map[string][]string{
		"quick-start":      {"-address", address},
		"feed-ipf-connect": {"-ipf", ipf, "-address", address},
		"ondemand-replay":  {"-address", "(ondemand:" + address + ")", "-symbol", "A:TEST"},
	}

	for _, name := range untilInterrupt {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			arguments, ok := args[name]
			if !ok {
				arguments = []string{"-address", address, "-symbol", "A:TEST"}
			}
			cmd := exec.Command(programs[name], arguments...)
			cmd.Dir = dir
			out := &syncBuffer{}
			cmd.Stdout, cmd.Stderr = out, out
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			exited := make(chan error, 1)
			go func() { exited <- cmd.Wait() }()

			// Let the sample connect and subscribe; it must still run then.
			select {
			case err := <-exited:
				t.Fatalf("exited before the interrupt: %v\n%s", err, out)
			case <-time.After(3 * time.Second):
			}
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				t.Fatal(err)
			}
			select {
			case err := <-exited:
				if err != nil {
					t.Errorf("exit after the interrupt: %v\n%s", err, out)
				}
			case <-time.After(30 * time.Second):
				_ = cmd.Process.Kill()
				<-exited
				t.Errorf("did not exit in 30 s after the interrupt:\n%s", out)
			}
		})
	}
}
