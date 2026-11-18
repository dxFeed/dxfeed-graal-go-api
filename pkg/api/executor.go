package api

import (
	"io"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// Executor runs the notifications of an endpoint: event listeners, state changes and the like; see
// EndpointBuilder.WithExecutor. The SDK's default is a thread pool with one thread per processor. Create one with
// NewFixedThreadPoolExecutor or NewInPlaceExecutor and close it when no endpoint is built with it anymore: closing
// releases the Go handle, while the endpoints already built keep using the executor.
type Executor interface {
	io.Closer
	handle() *native.ExecutorHandle
}

var (
	_ Executor = (*ThreadPoolExecutor)(nil)
	_ Executor = (*InPlaceExecutor)(nil)
)

// ThreadPoolExecutor runs the notifications on a fixed number of SDK threads.
type ThreadPoolExecutor struct {
	h *native.ExecutorHandle
}

// NewFixedThreadPoolExecutor creates an executor with the given number of threads, named with the given prefix.
func NewFixedThreadPoolExecutor(threads int, name string) (*ThreadPoolExecutor, error) {
	h, err := native.NewFixedThreadPoolExecutor(threads, name)
	if err != nil {
		return nil, err
	}
	return &ThreadPoolExecutor{h: h}, nil
}

// Close releases the executor's handle.
func (e *ThreadPoolExecutor) Close() error {
	return e.h.Free()
}

func (e *ThreadPoolExecutor) handle() *native.ExecutorHandle {
	if e == nil {
		return nil
	}
	return e.h
}

// InPlaceExecutor queues the notifications and runs them on the goroutine that calls ProcessAllPendingTasks, as the
// C++ API InPlaceExecutor does. It lets an application receive events in its own loop: call ProcessAllPendingTasks
// from one goroutine, since calls from several goroutines run the listeners concurrently, and do not Close the executor
// while a call runs.
//
// The waits of an endpoint built with this executor (AwaitProcessed, AwaitNotConnected, CloseAndAwaitTermination)
// complete only while ProcessAllPendingTasks keeps being called, so do not wait on the goroutine that calls it.
type InPlaceExecutor struct {
	h *native.ExecutorHandle
}

// NewInPlaceExecutor creates an executor whose notifications run in ProcessAllPendingTasks.
func NewInPlaceExecutor() (*InPlaceExecutor, error) {
	h, err := native.NewInPlaceExecutor()
	if err != nil {
		return nil, err
	}
	return &InPlaceExecutor{h: h}, nil
}

// ProcessAllPendingTasks runs the queued notifications on the calling goroutine and returns when the queue is empty.
func (e *InPlaceExecutor) ProcessAllPendingTasks() error {
	return e.h.ProcessAllPendingTasks()
}

// Close releases the executor's handle; ProcessAllPendingTasks fails after it.
func (e *InPlaceExecutor) Close() error {
	return e.h.Free()
}

func (e *InPlaceExecutor) handle() *native.ExecutorHandle {
	if e == nil {
		return nil
	}
	return e.h
}
