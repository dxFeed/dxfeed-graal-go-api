// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// ExecutorHandle holds a Java executor that runs the notifications of endpoints and subscriptions.
type ExecutorHandle struct {
	self Handler
}

// NewFixedThreadPoolExecutor creates a Java fixed thread pool with the given number of threads and thread name prefix.
func NewFixedThreadPoolExecutor(threads int, name string) (*ExecutorHandle, error) {
	var ptr *C.dxfg_v4_executor_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		namePtr := C.CString(name)
		defer C.free(unsafe.Pointer(namePtr))
		return checkStatus(thread, C.dxfg_v4_Executors_newFixedThreadPool(thread.ptr, C.int(threads), namePtr, &ptr))
	})
	if err != nil {
		return nil, err
	}
	return &ExecutorHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

// NewInPlaceExecutor creates an executor that only queues its tasks; ProcessAllPendingTasks runs them.
func NewInPlaceExecutor() (*ExecutorHandle, error) {
	var ptr *C.dxfg_v4_executor_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_new(thread.ptr, &ptr))
	})
	if err != nil {
		return nil, err
	}
	return &ExecutorHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

// ProcessAllPendingTasks runs the queued tasks of an executor made by NewInPlaceExecutor on the calling goroutine.
func (e *ExecutorHandle) ProcessAllPendingTasks() error {
	ptr := e.ptr()
	if ptr == nil {
		return errors.New("native: the executor is closed")
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_ExecutorBaseOnConcurrentLinkedQueue_processAllPendingTasks(thread.ptr, ptr))
	})
}

// Free releases the handle; the endpoints and subscriptions that use the executor keep it.
func (e *ExecutorHandle) Free() error {
	if e == nil {
		return nil
	}
	return e.self.Free()
}

// IsFreed reports whether Free was called.
func (e *ExecutorHandle) IsFreed() bool {
	return e.ptr() == nil
}

func (e *ExecutorHandle) ptr() *C.dxfg_v4_executor_t {
	if e == nil {
		return nil
	}
	return (*C.dxfg_v4_executor_t)(e.self.Ptr())
}
