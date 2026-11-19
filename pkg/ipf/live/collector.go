package live

import (
	"errors"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// InstrumentProfileUpdateListener is notified about the updates of the profiles of a collector.
type InstrumentProfileUpdateListener interface {
	// InstrumentProfilesUpdated receives the changed profiles: new, changed and removed ones (a removed profile has
	// the type ipf.TypeRemoved). The profiles are copies that the listener may keep.
	InstrumentProfilesUpdated(profiles []*events.InstrumentProfile)
}

type updateListenerEntry struct {
	listener InstrumentProfileUpdateListener
	native   *native.InstrumentProfileUpdateListener
}

// InstrumentProfileCollector collects instrument profiles, one per symbol, and notifies its listeners about the
// updates. An InstrumentProfileConnection fills it from a file, a URL or a server; UpdateInstrumentProfiles fills it
// by hand. Close it when it is no longer needed.
type InstrumentProfileCollector struct {
	// mu guards the handle: the calls hold it for reading, Close for writing.
	mu        sync.RWMutex
	handle    *native.InstrumentProfileCollectorHandle
	closed    bool
	listeners []updateListenerEntry
}

var errCollectorClosed = errors.New("live: the instrument profile collector is closed")

// NewInstrumentProfileCollector creates an empty collector.
func NewInstrumentProfileCollector() (*InstrumentProfileCollector, error) {
	handle, err := native.NewInstrumentProfileCollector()
	if err != nil {
		return nil, err
	}
	return &InstrumentProfileCollector{handle: handle}, nil
}

func (c *InstrumentProfileCollector) read(f func(h *native.InstrumentProfileCollectorHandle) error) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed {
		return errCollectorClosed
	}
	return f(c.handle)
}

// LastUpdateTime returns the time of the last update of the profiles.
func (c *InstrumentProfileCollector) LastUpdateTime() (time.Time, error) {
	var millis int64
	err := c.read(func(h *native.InstrumentProfileCollectorHandle) (err error) {
		millis, err = h.LastUpdateTime()
		return err
	})
	return time.UnixMilli(millis), err
}

// UpdateInstrumentProfile adds or replaces the profile of its symbol; a profile of the type ipf.TypeRemoved removes
// it. The collector keeps a copy: later changes of the profile do not reach it.
func (c *InstrumentProfileCollector) UpdateInstrumentProfile(profile *events.InstrumentProfile) error {
	return c.UpdateInstrumentProfiles([]*events.InstrumentProfile{profile})
}

// UpdateInstrumentProfiles updates the profiles like UpdateInstrumentProfile does, all at once.
func (c *InstrumentProfileCollector) UpdateInstrumentProfiles(profiles []*events.InstrumentProfile) error {
	return c.read(func(h *native.InstrumentProfileCollectorHandle) error {
		return h.UpdateInstrumentProfiles(profiles)
	})
}

// View returns copies of the profiles of the collector that are not removed.
func (c *InstrumentProfileCollector) View() ([]*events.InstrumentProfile, error) {
	var profiles []*events.InstrumentProfile
	err := c.read(func(h *native.InstrumentProfileCollectorHandle) (err error) {
		profiles, err = h.View()
		return err
	})
	return profiles, err
}

// AddUpdateListener adds a listener of the updates. If the collector has profiles, the listener first receives all
// of them. The listeners are called on a thread of the collector, one update at a time; a listener must not block,
// and a listener that panics is logged. Adding a listener that is already added does nothing.
func (c *InstrumentProfileCollector) AddUpdateListener(listener InstrumentProfileUpdateListener) error {
	if listener == nil {
		return errors.New("live: nil update listener")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errCollectorClosed
	}
	for _, e := range c.listeners {
		if e.listener == listener {
			return nil
		}
	}
	l, err := c.handle.AttachUpdateListener(listener.InstrumentProfilesUpdated)
	if err != nil {
		return err
	}
	listeners := make([]updateListenerEntry, 0, len(c.listeners)+1)
	c.listeners = append(append(listeners, c.listeners...), updateListenerEntry{listener, l})
	return nil
}

// RemoveUpdateListener removes a listener; an update that is being delivered may still reach it. The listener is
// found by ==, so its dynamic type must be comparable.
func (c *InstrumentProfileCollector) RemoveUpdateListener(listener InstrumentProfileUpdateListener) error {
	c.mu.Lock()
	var entry *updateListenerEntry
	for i, e := range c.listeners {
		if e.listener == listener {
			entry = &e
			listeners := make([]updateListenerEntry, 0, len(c.listeners)-1)
			c.listeners = append(append(listeners, c.listeners[:i]...), c.listeners[i+1:]...)
			break
		}
	}
	if entry == nil || c.closed {
		c.mu.Unlock()
		return nil
	}
	defer c.mu.Unlock()
	return c.handle.DetachUpdateListener(entry.native)
}

// Close removes the listeners and releases the native collector; a connection that fills the collector should be
// closed first. Later calls do nothing.
func (c *InstrumentProfileCollector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	var errs []error
	for _, e := range c.listeners {
		errs = append(errs, c.handle.DetachUpdateListener(e.native))
	}
	c.listeners = nil
	errs = append(errs, c.handle.Free())
	return errors.Join(errs...)
}
