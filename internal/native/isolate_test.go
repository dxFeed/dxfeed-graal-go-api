// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

import (
	"errors"
	"runtime"
	"sync"
	"testing"
	"unsafe"
)

func TestIsolateCreation(t *testing.T) {
	first, err := getOrCreateIsolate()
	if err != nil {
		t.Fatal(err)
	}
	second, err := getOrCreateIsolate()
	if err != nil {
		t.Fatal(err)
	}
	if first.ptr != second.ptr {
		t.Errorf("Multiple calls to getOrCreateIsolate returned different isolates")
	}
}

func TestMultipleAttachIsolateThreadInSameThread(t *testing.T) {
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		return executeInIsolateThread(func(nestedThread *isolateThread) error {
			if thread.ptr != nestedThread.ptr {
				t.Errorf("Nested call to executeInIsolateThread returned a different thread instance")
			}
			return nil
		})
	})
}

// TestThreadStaysAttached checks that an OS thread stays attached after a call: the next call on the same thread gets
// the same IsolateThread without attaching it again.
func TestThreadStaysAttached(t *testing.T) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		// Both calls run on the same OS thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		var first, second *isolateThread
		_ = executeInIsolateThread(func(thread *isolateThread) error { first = thread; return nil })
		_ = executeInIsolateThread(func(thread *isolateThread) error { second = thread; return nil })
		if first == nil || second == nil || first.ptr != second.ptr {
			t.Errorf("the second call got another IsolateThread: %v, %v", first, second)
		}
	}()
	<-done
}

// TestIsolateErrorIsReturned checks that a failure of the isolate creation is returned by every call instead of a
// panic. The real creation cannot be made to fail in a test, so the stored failure is replaced for the test.
func TestIsolateErrorIsReturned(t *testing.T) {
	if _, err := getOrCreateIsolate(); err != nil {
		t.Fatal(err)
	}
	failure := IsolateError(UninitializedIsolate)
	isolateErr = failure
	defer func() { isolateErr = nil }()

	for i := 0; i < 2; i++ {
		called := false
		err := executeInIsolateThread(func(*isolateThread) error {
			called = true
			return nil
		})
		if !errors.Is(err, failure) {
			t.Fatalf("call %d: error = %v, want %v", i, err, failure)
		}
		if called {
			t.Fatalf("call %d: the function ran without an isolate", i)
		}
	}
	if _, err := ParseSymbols("AAPL"); !errors.Is(err, failure) {
		t.Fatalf("ParseSymbols: error = %v, want %v", err, failure)
	}
}

// executeInIsolateThread pins its goroutine to an OS thread for the duration of the call, so every call in progress
// costs one real thread. Keep the concurrency modest: with too many goroutines at once the runtime hits the process
// thread limit and cgo aborts with "pthread_create failed: Resource temporarily unavailable".
func TestConcurrentAttachIsolateThread(t *testing.T) {
	const goroutines = 1000

	var wg sync.WaitGroup
	attached := make(chan uintptr, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := executeInIsolateThread(func(thread *isolateThread) error {
				if thread.ptr == nil {
					return errors.New("attached isolate thread is nil")
				}
				attached <- uintptr(unsafe.Pointer(thread.ptr))
				return nil
			})
			if err != nil {
				t.Errorf("executeInIsolateThread: %v", err)
			}
		}()
	}
	wg.Wait()
	close(attached)

	if len(attached) != goroutines {
		t.Errorf("attached %d threads, want %d", len(attached), goroutines)
	}
	// All attaches must share the single process-wide isolate.
	if isolate, err := getOrCreateIsolate(); err != nil || isolate.ptr == nil {
		t.Error("isolate pointer is nil after concurrent attaches")
	}
}
