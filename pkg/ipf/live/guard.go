package live

import (
	"errors"
	"reflect"
	"sync"
)

// guard keeps a native handle valid while calls use it: the calls do not hold a lock during the native calls, which
// may block (WaitUntilCompleted) or call listeners that call back, and close waits for the calls in progress before
// the handle is released.
type guard struct {
	mu     sync.Mutex
	closed bool
	calls  sync.WaitGroup
}

// enter starts a call; it returns false after close.
func (g *guard) enter() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return false
	}
	g.calls.Add(1)
	return true
}

// exit ends a call started by enter.
func (g *guard) exit() {
	g.calls.Done()
}

// close marks the guard closed and runs f under its lock; it returns false if the guard was already closed.
func (g *guard) close(f func()) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return false
	}
	g.closed = true
	f()
	return true
}

// wait waits for the calls in progress; call it after close.
func (g *guard) wait() {
	g.calls.Wait()
}

// checkListener returns an error for a nil listener or one that == cannot compare.
func checkListener(listener any) error {
	if listener == nil {
		return errors.New("live: nil listener")
	}
	if !reflect.TypeOf(listener).Comparable() {
		return errors.New("live: the listener type is not comparable; use a pointer")
	}
	return nil
}
