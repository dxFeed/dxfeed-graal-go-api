// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"runtime"
	"sync/atomic"
)

// listenerCall serializes the calls of a listener that come from several goroutines, such as the goroutine of the
// feed and the goroutine of a timer, and lets Close wait for the call in progress. A goroutine that is in the call
// may enter it again: the listener may call the model.
type listenerCall struct {
	sem   chan struct{}
	owner atomic.Int64
	// depth is the number of the nested calls of the owner; only the owner changes it.
	depth int
}

func newListenerCall() *listenerCall {
	return &listenerCall{sem: make(chan struct{}, 1)}
}

// enter starts a call; it returns false, without a call, if done is closed while it waits for another one.
func (c *listenerCall) enter(done <-chan struct{}) bool {
	id := goroutineID()
	if c.owner.Load() == id {
		c.depth++
		return true
	}
	select {
	case c.sem <- struct{}{}:
	case <-done:
		return false
	}
	c.owner.Store(id)
	c.depth = 1
	return true
}

// exit ends a call that enter started.
func (c *listenerCall) exit() {
	c.depth--
	if c.depth == 0 {
		c.owner.Store(0)
		<-c.sem
	}
}

// wait waits for the call in progress, unless it runs on this goroutine.
func (c *listenerCall) wait() {
	if c.owner.Load() == goroutineID() {
		return
	}
	c.sem <- struct{}{}
	<-c.sem
}

// goroutineID returns the id of the current goroutine, from the first line of its stack: "goroutine 18 [running]:".
// The ids start from 1.
func goroutineID() int64 {
	var buf [32]byte
	n := runtime.Stack(buf[:], false)
	const prefix = "goroutine "
	var id int64
	for _, c := range buf[min(len(prefix), n):n] {
		if c < '0' || c > '9' {
			break
		}
		id = id*10 + int64(c-'0')
	}
	return id
}
