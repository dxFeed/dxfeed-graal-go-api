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

// publisher is the publish-profiles sample running on a free local port.
type publisher struct {
	port   string
	server *exec.Cmd
	out    *syncBuffer
	exited chan struct{}
	err    error
}

// startPublisher builds publish-profiles, starts it on a free port and waits until it listens; the test stops it.
func startPublisher(t *testing.T, dir string) *publisher {
	t.Helper()
	program := build(t, dir, "publish-profiles")

	// A free port: listen on it, then free it for the publisher.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	p := &publisher{port: strconv.Itoa(listener.Addr().(*net.TCPAddr).Port), out: &syncBuffer{},
		exited: make(chan struct{})}
	_ = listener.Close()

	p.server = exec.Command(program, "-address", ":"+p.port)
	p.server.Dir = dir
	p.server.Stdout, p.server.Stderr = p.out, p.out
	if err := p.server.Start(); err != nil {
		t.Fatal(err)
	}
	go func() {
		p.err = p.server.Wait()
		close(p.exited)
	}()
	t.Cleanup(func() { p.stop() })

	// Wait until the server listens, so that a client does not wait for a reconnect. The log tells it, not a probe
	// connection, which TestPublisherSurvivesDisconnects tests on its own.
	deadline := time.Now().Add(30 * time.Second)
	for !strings.Contains(p.out.String(), "Listening at") {
		if p.hasExited() {
			t.Fatalf("publish-profiles: %s", p.stop())
		}
		if time.Now().After(deadline) {
			t.Fatalf("publish-profiles does not listen on %s:\n%s", p.port, p.stop())
		}
		time.Sleep(100 * time.Millisecond)
	}
	return p
}

func (p *publisher) hasExited() bool {
	select {
	case <-p.exited:
		return true
	default:
		return false
	}
}

// stop stops the server and describes how it ended and its output, which is complete only after Wait.
func (p *publisher) stop() string {
	if p.hasExited() {
		return fmt.Sprintf("exited by itself: %v\n%s", p.err, p.out)
	}
	_ = p.server.Process.Kill()
	<-p.exited
	return p.out.String()
}

func TestPublishAndRequestProfile(t *testing.T) {
	dir := t.TempDir()
	requester := build(t, dir, "request-profile")
	p := startPublisher(t, dir)

	out, err := tryRun(dir, requester, "-address", "127.0.0.1:"+p.port, "-symbol", "A:TEST", "-timeout", "20s")
	if err != nil || !strings.Contains(out, "A:TEST") || !strings.Contains(out, "Test symbol") {
		t.Errorf("request-profile: %v\n%s\npublish-profiles: %s", err, out, p.stop())
	}
}

// TestPublisherSurvivesDisconnects connects to a publisher and disconnects at once, so that the SDK writes to a closed
// connection. On Linux, the resulting SIGPIPE ended the process unless the API handled it.
func TestPublisherSurvivesDisconnects(t *testing.T) {
	p := startPublisher(t, t.TempDir())
	for i := 0; i < 30; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:"+p.port)
		if err != nil {
			t.Fatalf("connection %d: %v\npublish-profiles: %s", i, err, p.stop())
		}
		if i%2 == 1 {
			_ = conn.(*net.TCPConn).SetLinger(0) // a reset instead of a close
		}
		_ = conn.Close()
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)
	if p.hasExited() {
		t.Fatalf("publish-profiles ended after the disconnects: %s", p.stop())
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
