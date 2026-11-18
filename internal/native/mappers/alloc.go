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

// cStringOrNil returns nil for a nil string and a C copy of the string otherwise.
func cStringOrNil(s *string) *C.char {
	if s == nil {
		return nil
	}
	return cString(*s)
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
