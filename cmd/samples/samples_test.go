// Package samples_test builds the samples that need no network and runs them, so that every CI platform runs them.
package samples_test

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// build builds the sample in cmd/samples/<name> into dir and returns the path of the program.
func build(t *testing.T, dir, name string) string {
	t.Helper()
	program := filepath.Join(dir, name)
	if runtime.GOOS == "windows" {
		program += ".exe"
	}
	out, err := exec.Command("go", "build", "-o", program, "./"+name).CombinedOutput()
	if err != nil {
		t.Fatalf("go build ./%s: %v\n%s", name, err, out)
	}
	return program
}

// run runs a program in dir and returns its standard output.
func run(t *testing.T, dir, program string, args ...string) string {
	t.Helper()
	out, err := tryRun(dir, program, args...)
	if err != nil {
		t.Fatalf("%s %v: %v\n%s", filepath.Base(program), args, err, out)
	}
	return out
}

// tryRun runs a program in dir for at most a minute and returns its standard output, or, if it fails, its standard
// output and standard error.
func tryRun(dir, program string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, program, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return "stdout:\n" + stdout.String() + "\nstderr:\n" + stderr.String(), err
	}
	return stdout.String(), nil
}

// lines returns the lines of the output that contain s.
func lines(output, s string) []string {
	var result []string
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, s) {
			result = append(result, line)
		}
	}
	return result
}

func TestTapeFiles(t *testing.T) {
	dir := t.TempDir()
	writer := build(t, dir, "write-tape-file")
	converter := build(t, dir, "convert-tape-file")
	parser := build(t, dir, "file-parser")

	run(t, dir, writer)
	out := run(t, dir, parser, "-file", "WriteTapeFile.out.txt", "-types", "Quote", "-symbols", "TEST1,TEST2")
	if got := lines(out, "Quote{"); len(got) != 2 || !strings.Contains(out, "TEST1") || !strings.Contains(out, "TEST2") {
		t.Errorf("the written tape file has %d quotes:\n%s", len(got), out)
	}

	input, err := os.ReadFile(filepath.Join("convert-tape-file", "ConvertTapeFile.in"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ConvertTapeFile.in"), input, 0o644); err != nil {
		t.Fatal(err)
	}
	run(t, dir, converter)
	out = run(t, dir, parser, "-file", "ConvertTapeFile.out", "-types", "Quote", "-symbols", "AAPL")
	if got := lines(out, "Quote{"); len(got) != 20 {
		t.Errorf("the converted tape file has %d quotes, want 20:\n%s", len(got), out)
	}
}

func TestPublishAndRequestProfile(t *testing.T) {
	dir := t.TempDir()
	publisher := build(t, dir, "publish-profiles")
	requester := build(t, dir, "request-profile")

	// A free port: listen on it, then free it for the publisher.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	_ = listener.Close()

	server := exec.Command(publisher, "-address", ":"+port)
	server.Dir = dir
	serverOut := &syncBuffer{}
	server.Stdout, server.Stderr = serverOut, serverOut
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	exited := make(chan struct{})
	var exitErr error
	go func() {
		exitErr = server.Wait()
		close(exited)
	}()
	// stop stops the server and describes how it ended and its output, which is complete only after Wait.
	stop := func() string {
		select {
		case <-exited:
			return fmt.Sprintf("exited by itself: %v\n%s", exitErr, serverOut)
		default:
		}
		_ = server.Process.Kill()
		<-exited
		return serverOut.String()
	}
	defer stop()

	// Wait until the server listens, so that the client does not wait for a reconnect. The log tells it: a probe
	// connection is not used, since the server writes to the connection that the probe has closed.
	deadline := time.Now().Add(30 * time.Second)
	for !strings.Contains(serverOut.String(), "Listening at") {
		select {
		case <-exited:
			t.Fatalf("publish-profiles: %s", stop())
		default:
		}
		if time.Now().After(deadline) {
			t.Fatalf("publish-profiles does not listen on %s:\n%s", port, stop())
		}
		time.Sleep(100 * time.Millisecond)
	}

	out, err := tryRun(dir, requester, "-address", "127.0.0.1:"+port, "-symbol", "A:TEST", "-timeout", "20s")
	if err != nil || !strings.Contains(out, "A:TEST") || !strings.Contains(out, "Test symbol") {
		t.Errorf("request-profile: %v\n%s\npublish-profiles: %s", err, out, stop())
	}
}

// syncBuffer is a bytes.Buffer that a process writes while the test reads it.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
