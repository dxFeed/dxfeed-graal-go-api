package live

import (
	"errors"
	"fmt"
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
	guard     guard
	handle    *native.InstrumentProfileConnectionHandle
	collector *InstrumentProfileCollector
	listeners []stateListenerEntry // guarded by guard.mu
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
	err := collector.call(func(h *native.InstrumentProfileCollectorHandle) (err error) {
		handle, err = native.NewInstrumentProfileConnection(address, h)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &InstrumentProfileConnection{handle: handle, collector: collector}, nil
}

func (c *InstrumentProfileConnection) call(f func(h *native.InstrumentProfileConnectionHandle) error) error {
	if !c.guard.enter() {
		return errConnectionClosed
	}
	defer c.guard.exit()
	return f(c.handle)
}

// Collector returns the collector of the connection.
func (c *InstrumentProfileConnection) Collector() *InstrumentProfileCollector {
	return c.collector
}

// Address returns the address of the connection without its options, such as [update=30s].
func (c *InstrumentProfileConnection) Address() (string, error) {
	var address string
	err := c.call(func(h *native.InstrumentProfileConnectionHandle) (err error) {
		address, err = h.Address()
		return err
	})
	return address, err
}

// UpdatePeriod returns the period of the update checks of a source without live updates, and of the reconnections.
func (c *InstrumentProfileConnection) UpdatePeriod() (time.Duration, error) {
	var millis int64
	err := c.call(func(h *native.InstrumentProfileConnectionHandle) (err error) {
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
	return c.call(func(h *native.InstrumentProfileConnectionHandle) error {
		return h.SetUpdatePeriod(period.Milliseconds())
	})
}

// State returns the state of the connection; it is Closed after Close.
func (c *InstrumentProfileConnection) State() (State, error) {
	state := Closed
	err := c.call(func(h *native.InstrumentProfileConnectionHandle) error {
		s, err := h.State()
		state = State(s)
		return err
	})
	if errors.Is(err, errConnectionClosed) {
		return Closed, nil
	}
	return state, err
}

// LastModified returns the modification time of the last data read, or the zero time if it is unknown. The SDK sets
// it only after it has read all the data of a download: for a file or a plain HTTP source a moment after the state
// becomes Completed (right after WaitUntilCompleted returns or the listeners get Completed, it may still be the zero
// time), for a live streaming source only when the stream ends.
func (c *InstrumentProfileConnection) LastModified() (time.Time, error) {
	var millis int64
	err := c.call(func(h *native.InstrumentProfileConnectionHandle) (err error) {
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
	return c.call(func(h *native.InstrumentProfileConnectionHandle) error {
		return h.Start()
	})
}

// WaitUntilCompleted waits until the initial profiles are read (the state Completed), the connection is closed or the
// timeout elapses; it reports whether the profiles were read. LastModified may lag behind: see there.
func (c *InstrumentProfileConnection) WaitUntilCompleted(timeout time.Duration) (bool, error) {
	var completed bool
	err := c.call(func(h *native.InstrumentProfileConnectionHandle) (err error) {
		completed, err = h.WaitUntilCompleted(timeout)
		return err
	})
	return completed, err
}

// AddStateChangeListener adds a listener of the state changes. The listeners are called on the thread that changes
// the state: the goroutine of Start and Close for Connecting and Closed, a thread of the connection for the others. A
// listener must not block; a listener that panics is logged. The listener is found by ==, so its type must be
// comparable, such as a pointer; adding a listener that is already added does nothing.
func (c *InstrumentProfileConnection) AddStateChangeListener(listener StateChangeListener) error {
	if err := checkListener(listener); err != nil {
		return err
	}
	return c.call(func(h *native.InstrumentProfileConnectionHandle) error {
		c.guard.mu.Lock()
		for _, e := range c.listeners {
			if e.listener == listener {
				c.guard.mu.Unlock()
				return nil
			}
		}
		c.guard.mu.Unlock()
		l, err := h.AttachStateListener(func(oldState, newState int32) {
			listener.StateChanged(State(oldState), State(newState))
		})
		if err != nil {
			return err
		}
		// The listener is detached without the lock: a callback in progress may need it.
		c.guard.mu.Lock()
		if c.guard.closed {
			c.guard.mu.Unlock()
			return errors.Join(errConnectionClosed, h.DetachStateListener(l))
		}
		for _, e := range c.listeners {
			if e.listener == listener {
				c.guard.mu.Unlock()
				return h.DetachStateListener(l)
			}
		}
		listeners := make([]stateListenerEntry, 0, len(c.listeners)+1)
		c.listeners = append(append(listeners, c.listeners...), stateListenerEntry{listener, l})
		c.guard.mu.Unlock()
		return nil
	})
}

// RemoveStateChangeListener removes a listener; a change that is being delivered may still reach it.
func (c *InstrumentProfileConnection) RemoveStateChangeListener(listener StateChangeListener) error {
	if checkListener(listener) != nil {
		return nil
	}
	return c.call(func(h *native.InstrumentProfileConnectionHandle) error {
		c.guard.mu.Lock()
		var removed *native.IpfStateListener
		for i, e := range c.listeners {
			if e.listener == listener {
				removed = e.native
				listeners := make([]stateListenerEntry, 0, len(c.listeners)-1)
				c.listeners = append(append(listeners, c.listeners[:i]...), c.listeners[i+1:]...)
				break
			}
		}
		c.guard.mu.Unlock()
		if removed == nil {
			return nil
		}
		return h.DetachStateListener(removed)
	})
}

// Close closes the connection: the listeners receive the change to Closed, nothing more reaches the collector, and a
// WaitUntilCompleted in progress returns false. The native connection and listeners are released after the calls in
// progress end. The collector is not closed. Later calls do nothing; Close must not be called from a state listener.
func (c *InstrumentProfileConnection) Close() error {
	var listeners []stateListenerEntry
	if !c.guard.close(func() {
		listeners, c.listeners = c.listeners, nil
	}) {
		return nil
	}
	// The calls in progress still hold the handle, so it is valid; the native close wakes up WaitUntilCompleted.
	errs := []error{c.handle.Close()}
	for _, e := range listeners {
		errs = append(errs, c.handle.DetachStateListener(e.native))
	}
	c.guard.wait()
	errs = append(errs, c.handle.Free())
	return errors.Join(errs...)
}
