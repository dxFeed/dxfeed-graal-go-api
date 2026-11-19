// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package live

import (
	"errors"
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
	guard     guard
	handle    *native.InstrumentProfileCollectorHandle
	listeners []updateListenerEntry // guarded by guard.mu
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

func (c *InstrumentProfileCollector) call(f func(h *native.InstrumentProfileCollectorHandle) error) error {
	if !c.guard.enter() {
		return errCollectorClosed
	}
	defer c.guard.exit()
	return f(c.handle)
}

// LastUpdateTime returns the time of the last update of the profiles.
func (c *InstrumentProfileCollector) LastUpdateTime() (time.Time, error) {
	var millis int64
	err := c.call(func(h *native.InstrumentProfileCollectorHandle) (err error) {
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
	return c.call(func(h *native.InstrumentProfileCollectorHandle) error {
		return h.UpdateInstrumentProfiles(profiles)
	})
}

// View returns copies of the profiles of the collector that are not removed.
func (c *InstrumentProfileCollector) View() ([]*events.InstrumentProfile, error) {
	var profiles []*events.InstrumentProfile
	err := c.call(func(h *native.InstrumentProfileCollectorHandle) (err error) {
		profiles, err = h.View()
		return err
	})
	return profiles, err
}

// AddUpdateListener adds a listener of the updates. If the collector has profiles, the listener first receives all
// of them. The listeners are called on a thread of the collector, one update at a time; a listener must not block,
// and a listener that panics is logged. The listener is found by ==, so its type must be comparable, such as a
// pointer; adding a listener that is already added does nothing.
func (c *InstrumentProfileCollector) AddUpdateListener(listener InstrumentProfileUpdateListener) error {
	if err := checkListener(listener); err != nil {
		return err
	}
	return c.call(func(h *native.InstrumentProfileCollectorHandle) error {
		c.guard.mu.Lock()
		for _, e := range c.listeners {
			if e.listener == listener {
				c.guard.mu.Unlock()
				return nil
			}
		}
		c.guard.mu.Unlock()
		l, err := h.AttachUpdateListener(listener.InstrumentProfilesUpdated)
		if err != nil {
			return err
		}
		// The listener is detached without the lock: a callback in progress may need it.
		c.guard.mu.Lock()
		if c.guard.closed {
			// Closed concurrently: Close has taken the listeners already.
			c.guard.mu.Unlock()
			return errors.Join(errCollectorClosed, h.DetachUpdateListener(l))
		}
		for _, e := range c.listeners {
			if e.listener == listener {
				// Added concurrently: keep the other one.
				c.guard.mu.Unlock()
				return h.DetachUpdateListener(l)
			}
		}
		listeners := make([]updateListenerEntry, 0, len(c.listeners)+1)
		c.listeners = append(append(listeners, c.listeners...), updateListenerEntry{listener, l})
		c.guard.mu.Unlock()
		return nil
	})
}

// RemoveUpdateListener removes a listener; an update that is being delivered may still reach it.
func (c *InstrumentProfileCollector) RemoveUpdateListener(listener InstrumentProfileUpdateListener) error {
	if checkListener(listener) != nil {
		return nil
	}
	return c.call(func(h *native.InstrumentProfileCollectorHandle) error {
		c.guard.mu.Lock()
		var removed *native.InstrumentProfileUpdateListener
		for i, e := range c.listeners {
			if e.listener == listener {
				removed = e.native
				listeners := make([]updateListenerEntry, 0, len(c.listeners)-1)
				c.listeners = append(append(listeners, c.listeners[:i]...), c.listeners[i+1:]...)
				break
			}
		}
		c.guard.mu.Unlock()
		if removed == nil {
			return nil
		}
		return h.DetachUpdateListener(removed)
	})
}

// Close removes the listeners and releases the native collector after the calls in progress end; a connection that
// fills the collector should be closed first. Later calls do nothing.
func (c *InstrumentProfileCollector) Close() error {
	var listeners []updateListenerEntry
	if !c.guard.close(func() {
		listeners, c.listeners = c.listeners, nil
	}) {
		return nil
	}
	var errs []error
	for _, e := range listeners {
		errs = append(errs, c.handle.DetachUpdateListener(e.native))
	}
	c.guard.wait()
	errs = append(errs, c.handle.Free())
	return errors.Join(errs...)
}
