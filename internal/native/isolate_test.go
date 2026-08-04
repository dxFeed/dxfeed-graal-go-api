package native

import (
	"errors"
	"sync"
	"testing"
	"unsafe"
)

func TestIsolateCreation(t *testing.T) {
	if getOrCreateIsolate().ptr != getOrCreateIsolate().ptr {
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

func TestDoubleDetachIsolateThread(t *testing.T) {
	thread := attachCurrentThread()
	thread.detach()
	thread.detach()
}

// attachCurrentThread pins its goroutine to an OS thread for the duration of
// the call, so every in-flight iteration costs one real thread. Keep the
// concurrency modest: with too many goroutines at once the runtime hits the
// process thread limit and cgo aborts with
// "pthread_create failed: Resource temporarily unavailable".
func TestConcurrentAttachDetachIsolateThread(t *testing.T) {
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
	if getOrCreateIsolate().ptr == nil {
		t.Error("isolate pointer is nil after concurrent attach/detach")
	}
}
