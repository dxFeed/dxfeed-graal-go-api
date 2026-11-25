// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

// #include <stdlib.h>
import "C"
import (
	"sync/atomic"
	"unsafe"
)

// liveAllocations counts the C allocations of CEvent that FreeCEvent has not freed yet.
var liveAllocations atomic.Int64

// LiveAllocations returns the number of C allocations made by CEvent and not yet freed by FreeCEvent; tests use it to
// find leaks.
func LiveAllocations() int64 {
	return liveAllocations.Load()
}

func cString(s string) *C.char {
	liveAllocations.Add(1)
	return C.CString(s)
}

func cCalloc(size uintptr) unsafe.Pointer {
	liveAllocations.Add(1)
	return C.calloc(1, C.size_t(size))
}

// cFree frees memory allocated by cString or cCalloc; nil does nothing.
func cFree(p unsafe.Pointer) {
	if p == nil {
		return
	}
	liveAllocations.Add(-1)
	C.free(p)
}

// arenaChunk is the largest size of the C memory blocks of an Allocator: a few hundred events with their strings.
const arenaChunk = 64 << 10

// arenaMinChunk is the smallest size of a block that Reserve sets.
const arenaMinChunk = 1 << 10

// arenaAlign aligns the allocations of an Allocator for any C type of the events.
const arenaAlign = 16

// Allocator allocates the C memory of the events of one batch from a few large blocks, so that the batch costs a
// constant number of cgo calls; Free releases all of it at once. A nil Allocator allocates and frees every event
// separately, with CEvent and FreeCEvent. An Allocator is not safe for concurrent use.
type Allocator struct {
	blocks []unsafe.Pointer
	last   unsafe.Pointer // the last block
	used   uintptr        // the used bytes of the last block
	size   uintptr        // the size of the last block
	next   uintptr        // the size of the next block set by Reserve; 0 for arenaChunk
}

// Reserve sets the size of the first block for a batch of a known size, so that a small batch does not take a whole
// arenaChunk; the blocks after it take arenaChunk.
func (a *Allocator) Reserve(size uintptr) {
	a.next = min(max(size, arenaMinChunk), arenaChunk)
}

// calloc returns zeroed memory of the size: from the blocks of the allocator, or from C.calloc for a nil allocator.
func (a *Allocator) calloc(size uintptr) unsafe.Pointer {
	if a == nil {
		return cCalloc(size)
	}
	size = (size + arenaAlign - 1) &^ (arenaAlign - 1)
	if a.last == nil || a.used+size > a.size {
		block := uintptr(arenaChunk)
		if a.next != 0 {
			block, a.next = a.next, 0
		}
		block = max(size, block)
		// The blocks come from calloc and are never reused, so the memory is zeroed.
		p := cCalloc(block)
		if p == nil {
			panic("mappers: out of C memory")
		}
		a.blocks = append(a.blocks, p)
		a.last, a.used, a.size = p, 0, block
	}
	p := unsafe.Add(a.last, a.used)
	a.used += size
	return p
}

// Calloc returns zeroed C memory of the size from the blocks of the allocator, valid until Free.
func (a *Allocator) Calloc(size uintptr) unsafe.Pointer { return a.calloc(size) }

// cString returns a C copy of the string.
func (a *Allocator) cString(s string) *C.char {
	if a == nil {
		return cString(s)
	}
	p := a.calloc(uintptr(len(s)) + 1)
	copy(unsafe.Slice((*byte)(p), len(s)), s)
	return (*C.char)(p)
}

// cStringOrNil returns nil for a nil string and a C copy of the string otherwise.
func (a *Allocator) cStringOrNil(s *string) *C.char {
	if s == nil {
		return nil
	}
	return a.cString(*s)
}

// Free releases the memory of the allocator; the events allocated from it must not be used after that. The allocator
// can be used again.
func (a *Allocator) Free() {
	if a == nil {
		return
	}
	for _, block := range a.blocks {
		cFree(block)
	}
	clear(a.blocks)
	a.blocks = a.blocks[:0]
	a.last, a.used, a.size, a.next = nil, 0, 0, 0
}
