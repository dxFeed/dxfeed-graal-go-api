package native

import (
	"sync"
	"testing"
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

func TestConcurrentAttachDetachIsolateThread(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			_ = executeInIsolateThread(func(thread *isolateThread) error {
				defer wg.Done()
				return nil
			})
		}()
	}
	wg.Wait()
}
