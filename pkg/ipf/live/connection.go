package live

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// State is the state of an InstrumentProfileConnection.
type State int32

// The states of a connection.
const (
	// NotConnected is the state of a new connection, before Start.
	NotConnected State = iota
	// Connecting is the state while the connection connects to its source.
	Connecting
	// Connected is the state when the connection is established and the profiles are being read.
	Connected
	// Completed is the state when the initial profiles have been read; live updates may follow.
	Completed
	// Closed is the state after Close.
	Closed
)

// String returns the name of the state, such as "COMPLETED".
func (s State) String() string {
	switch s {
	case NotConnected:
		return "NOT_CONNECTED"
	case Connecting:
		return "CONNECTING"
	case Connected:
		return "CONNECTED"
	case Completed:
		return "COMPLETED"
	case Closed:
		return "CLOSED"
	}
	return fmt.Sprintf("State(%d)", int32(s))
}

// StateChangeListener is notified about the state changes of a connection.
type StateChangeListener interface {
	// StateChanged receives the previous and the new state.
	StateChanged(oldState, newState State)
}

type stateListenerEntry struct {
	listener StateChangeListener
	native   *native.IpfStateListener
}

// InstrumentProfileConnection reads instrument profiles from a file, a URL or an instrument profile server into a
// collector and keeps them up to date: a server pushes live updates, a file or a URL is checked again every update
// period. Create it with CreateConnection, then Start it; Close it when it is no longer needed.
type InstrumentProfileConnection struct {
	// mu guards the handle: the calls hold it for reading, the release in Close for writing.
	mu        sync.RWMutex
	handle    *native.InstrumentProfileConnectionHandle
	collector *InstrumentProfileCollector
	closed    bool
	listeners []stateListenerEntry
}

var errConnectionClosed = errors.New("live: the instrument profile connection is closed")

// CreateConnection creates a connection that reads the profiles from the address into the collector. The address is
// host:port of an instrument profile server, a URL or a file path; a .zip or .gz file is unpacked. The suffix
// [update=<period>], such as [update=30s], sets the update period (1 minute by default). The connection is not
// started.
func CreateConnection(address string, collector *InstrumentProfileCollector) (*InstrumentProfileConnection, error) {
	if collector == nil {
		return nil, errors.New("live: nil collector")
	}
	var handle *native.InstrumentProfileConnectionHandle
	err := collector.read(func(h *native.InstrumentProfileCollectorHandle) (err error) {
		handle, err = native.NewInstrumentProfileConnection(address, h)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &InstrumentProfileConnection{handle: handle, collector: collector}, nil
}

func (c *InstrumentProfileConnection) read(f func(h *native.InstrumentProfileConnectionHandle) error) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return errConnectionClosed
	}
	return f(c.handle)
}

// Collector returns the collector of the connection.
func (c *InstrumentProfileConnection) Collector() *InstrumentProfileCollector {
	return c.collector
}

// Address returns the address of the connection without its options, such as [update=30s].
func (c *InstrumentProfileConnection) Address() (string, error) {
	var address string
	err := c.read(func(h *native.InstrumentProfileConnectionHandle) (err error) {
		address, err = h.Address()
		return err
	})
	return address, err
}

// UpdatePeriod returns the period of the update checks of a source without live updates, and of the reconnections.
func (c *InstrumentProfileConnection) UpdatePeriod() (time.Duration, error) {
	var millis int64
	err := c.read(func(h *native.InstrumentProfileConnectionHandle) (err error) {
		millis, err = h.UpdatePeriod()
		return err
	})
	return time.Duration(millis) * time.Millisecond, err
}

// SetUpdatePeriod changes the update period, which must not be negative; it is rounded down to milliseconds.
func (c *InstrumentProfileConnection) SetUpdatePeriod(period time.Duration) error {
	if period < 0 {
		return fmt.Errorf("live: negative update period %v", period)
	}
	return c.read(func(h *native.InstrumentProfileConnectionHandle) error {
		return h.SetUpdatePeriod(period.Milliseconds())
	})
}

// State returns the state of the connection; it is Closed after Close.
func (c *InstrumentProfileConnection) State() (State, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return Closed, nil
	}
	state, err := c.handle.State()
	return State(state), err
}

// LastModified returns the modification time of the last data read, or the zero time if it is unknown.
func (c *InstrumentProfileConnection) LastModified() (time.Time, error) {
	var millis int64
	err := c.read(func(h *native.InstrumentProfileConnectionHandle) (err error) {
		millis, err = h.LastModified()
		return err
	})
	if millis == 0 {
		return time.Time{}, err
	}
	return time.UnixMilli(millis), err
}

// Start starts the connection, which connects in the background; it fails unless the state is NotConnected. The
// listeners receive the change to Connecting before Start returns.
func (c *InstrumentProfileConnection) Start() error {
	return c.read(func(h *native.InstrumentProfileConnectionHandle) error {
		return h.Start()
	})
}

// WaitUntilCompleted waits until the initial profiles are read (the state Completed), the connection is closed or the
// timeout elapses; it reports whether the profiles were read.
func (c *InstrumentProfileConnection) WaitUntilCompleted(timeout time.Duration) (bool, error) {
	var completed bool
	err := c.read(func(h *native.InstrumentProfileConnectionHandle) (err error) {
		completed, err = h.WaitUntilCompleted(timeout)
		return err
	})
	return completed, err
}

// AddStateChangeListener adds a listener of the state changes. The listeners are called on the thread that changes
// the state: the goroutine of Start and Close for Connecting and Closed, a thread of the connection for the others. A
// listener must not block and must not call Close; a listener that panics is logged. Adding a listener that is already
// added does nothing.
func (c *InstrumentProfileConnection) AddStateChangeListener(listener StateChangeListener) error {
	if listener == nil {
		return errors.New("live: nil state change listener")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errConnectionClosed
	}
	for _, e := range c.listeners {
		if e.listener == listener {
			return nil
		}
	}
	l, err := c.handle.AttachStateListener(func(oldState, newState int32) {
		listener.StateChanged(State(oldState), State(newState))
	})
	if err != nil {
		return err
	}
	listeners := make([]stateListenerEntry, 0, len(c.listeners)+1)
	c.listeners = append(append(listeners, c.listeners...), stateListenerEntry{listener, l})
	return nil
}

// RemoveStateChangeListener removes a listener; a change that is being delivered may still reach it. The listener is
// found by ==, so its dynamic type must be comparable.
func (c *InstrumentProfileConnection) RemoveStateChangeListener(listener StateChangeListener) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, e := range c.listeners {
		if e.listener == listener {
			listeners := make([]stateListenerEntry, 0, len(c.listeners)-1)
			c.listeners = append(append(listeners, c.listeners[:i]...), c.listeners[i+1:]...)
			if c.closed {
				return nil
			}
			return c.handle.DetachStateListener(e.native)
		}
	}
	return nil
}

// Close closes the connection: the listeners receive the change to Closed, nothing more reaches the collector, and
// the native connection and listeners are released. A WaitUntilCompleted in progress returns false. The collector is
// not closed. Later calls do nothing.
func (c *InstrumentProfileConnection) Close() error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil
	}
	// Close without the write lock: it wakes up a WaitUntilCompleted that holds the read lock.
	err := c.handle.Close()
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return err
	}
	c.closed = true
	errs := []error{err}
	for _, e := range c.listeners {
		errs = append(errs, c.handle.DetachStateListener(e.native))
	}
	c.listeners = nil
	errs = append(errs, c.handle.Free())
	return errors.Join(errs...)
}
