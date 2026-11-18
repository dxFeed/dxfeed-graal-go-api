package api

import (
	"errors"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// The names of the properties of OnDemandService that OnDemandPropertyChangeListener is told about.
const (
	// OnDemandPropertyReplaySupported is the property of IsReplaySupported; its value is a bool.
	OnDemandPropertyReplaySupported = "replaySupported"
	// OnDemandPropertyReplay is the property of IsReplay; its value is a bool.
	OnDemandPropertyReplay = "replay"
	// OnDemandPropertyClear is the property of IsClear; its value is a bool.
	OnDemandPropertyClear = "clear"
	// OnDemandPropertyTime is the property of Time; its value is a time.Time.
	OnDemandPropertyTime = "time"
	// OnDemandPropertySpeed is the property of Speed; its value is a float64.
	OnDemandPropertySpeed = "speed"
)

// OnDemandPropertyChangeListener is told about the changes of the properties of an OnDemandService: see
// OnDemandService.AddPropertyChangeListener.
type OnDemandPropertyChangeListener interface {
	// PropertyChanged is called after a property changed. name is one of the OnDemandProperty constants (a later SDK
	// may add others); a value is a bool, a float64, a time.Time, or nil for a value of an unknown type.
	PropertyChanged(name string, oldValue, newValue any)
}

// OnDemandService replays the historical market data of an on-demand data provider through the feed of an endpoint:
// the events come to the subscriptions of the feed as if they came from the market now. Get it with
// DXEndpoint.OnDemandService.
//
// The endpoint is connected to an on-demand data provider with an address that has an "(ondemand:<address>)"
// component, such as "(ondemand:demo.dxfeed.com:7680)"; IsReplaySupported tells whether it is. Replay switches the
// feed to the replay of the data from a time at a speed; Pause and SetSpeed control it; StopAndResume returns to the
// ordinary data and StopAndClear stops the replay and clears the data. An endpoint with the role OnDemandFeed has no
// ordinary data: it is not connected while the replay is stopped.
//
// The methods may be called from any goroutine. The service lives as long as its endpoint: after Close of the endpoint
// its methods return errors. See the Java API OnDemandService for details.
type OnDemandService struct {
	endpoint *DXEndpoint
	handle   *native.OnDemandServiceHandle

	// listeners is replaced, never changed in place; one native listener, attached on the first add and removed when
	// the endpoint closes, passes the changes to them.
	mu             sync.Mutex
	listeners      []OnDemandPropertyChangeListener
	nativeListener *native.OnDemandPropertyListener
	closed         bool
}

// OnDemandService returns the on-demand service of the endpoint, the same one on each call. The endpoint must have the
// role Feed or OnDemandFeed.
func (e *DXEndpoint) OnDemandService() (*OnDemandService, error) {
	e.onDemandMu.Lock()
	defer e.onDemandMu.Unlock()
	if e.onDemandClosed {
		return nil, errors.New("dxfeed: the endpoint is closed")
	}
	if e.onDemand != nil {
		return e.onDemand, nil
	}
	if e.role != Feed && e.role != OnDemandFeed {
		return nil, errors.New("dxfeed: the on-demand service needs an endpoint with the role Feed or OnDemandFeed")
	}
	handle, err := e.endpointHandle.OnDemandService()
	if err != nil {
		return nil, err
	}
	e.onDemand = &OnDemandService{endpoint: e, handle: handle}
	return e.onDemand, nil
}

// closeOnDemandService releases the on-demand service before the endpoint closes.
func (e *DXEndpoint) closeOnDemandService() error {
	e.onDemandMu.Lock()
	s := e.onDemand
	e.onDemand = nil
	e.onDemandClosed = true
	e.onDemandMu.Unlock()
	if s == nil {
		return nil
	}
	return s.close()
}

// Endpoint returns the endpoint of the service.
func (s *OnDemandService) Endpoint() *DXEndpoint {
	return s.endpoint
}

// IsReplaySupported reports whether the endpoint is connected to an on-demand data provider, so that Replay works.
func (s *OnDemandService) IsReplaySupported() (bool, error) {
	return s.handle.IsReplaySupported()
}

// IsReplay reports whether the service is in the replay mode: after Replay and before StopAndResume or StopAndClear.
// A paused replay is still the replay mode.
func (s *OnDemandService) IsReplay() (bool, error) {
	return s.handle.IsReplay()
}

// IsClear reports whether the service is in the clear mode: after StopAndClear and before Replay or StopAndResume.
func (s *OnDemandService) IsClear() (bool, error) {
	return s.handle.IsClear()
}

// Time returns the time that is being replayed in the replay mode, otherwise the last replayed time; the Unix epoch
// before the first replay.
func (s *OnDemandService) Time() (time.Time, error) {
	millis, err := s.handle.Time()
	return time.UnixMilli(millis), err
}

// Speed returns the replay speed: 1 is the real time, 0 is a paused replay. It is 0 outside the replay mode.
func (s *OnDemandService) Speed() (float64, error) {
	return s.handle.Speed()
}

// Replay starts the replay of the data from the time at the real speed; it is ReplayWithSpeed with the speed 1.
func (s *OnDemandService) Replay(from time.Time) error {
	return s.ReplayWithSpeed(from, 1)
}

// ReplayWithSpeed starts the replay of the data from the time at the speed, which must not be negative: 1 is the real
// time, 0 pauses the replay at the time. It fails when the replay is not supported (see IsReplaySupported).
func (s *OnDemandService) ReplayWithSpeed(from time.Time, speed float64) error {
	return s.handle.Replay(from.UnixMilli(), speed)
}

// Pause pauses the replay at the current time; it is SetSpeed with the speed 0. Outside the replay mode it does
// nothing.
func (s *OnDemandService) Pause() error {
	return s.handle.Pause()
}

// StopAndResume stops the replay and resumes the ordinary data. An endpoint with the role OnDemandFeed has no ordinary
// data, so for it this is the same as StopAndClear, but the data is not cleared.
func (s *OnDemandService) StopAndResume() error {
	return s.handle.StopAndResume()
}

// StopAndClear stops the replay and clears the data; Replay continues the replay and StopAndResume resumes the
// ordinary data.
func (s *OnDemandService) StopAndClear() error {
	return s.handle.StopAndClear()
}

// SetSpeed changes the replay speed, which must not be negative; 0 pauses the replay.
func (s *OnDemandService) SetSpeed(speed float64) error {
	return s.handle.SetSpeed(speed)
}

// AddPropertyChangeListener adds a listener of the changes of the properties of the service. A listener may be called
// in the goroutine of the method that changed the property, before the method returns, or on an SDK thread; it must
// not block. A listener that panics is logged and the other listeners miss that change.
func (s *OnDemandService) AddPropertyChangeListener(listener OnDemandPropertyChangeListener) error {
	if listener == nil {
		return errors.New("dxfeed: the on-demand property change listener is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errors.New("dxfeed: the endpoint is closed")
	}
	if s.nativeListener == nil {
		l, err := s.handle.AttachPropertyListener(s.propertyChanged)
		if err != nil {
			return err
		}
		s.nativeListener = l
	}
	listeners := make([]OnDemandPropertyChangeListener, 0, len(s.listeners)+1)
	s.listeners = append(append(listeners, s.listeners...), listener)
	return nil
}

// RemovePropertyChangeListener removes a listener added with AddPropertyChangeListener; a change that is being
// reported may still reach it. The listener is found by ==, so its dynamic type must be comparable.
func (s *OnDemandService) RemovePropertyChangeListener(listener OnDemandPropertyChangeListener) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, l := range s.listeners {
		if l == listener {
			listeners := make([]OnDemandPropertyChangeListener, 0, len(s.listeners)-1)
			s.listeners = append(append(listeners, s.listeners[:i]...), s.listeners[i+1:]...)
			return
		}
	}
}

// propertyChanged is the callback of the native listener: it passes the change to a copy of the listeners.
func (s *OnDemandService) propertyChanged(name string, oldValue, newValue any) {
	s.mu.Lock()
	listeners := s.listeners
	s.mu.Unlock()
	if len(listeners) == 0 {
		return
	}
	oldValue, newValue = onDemandPropertyValue(oldValue), onDemandPropertyValue(newValue)
	for _, listener := range listeners {
		listener.PropertyChanged(name, oldValue, newValue)
	}
}

// onDemandPropertyValue turns a native time, an int64 of milliseconds, into a time.Time.
func onDemandPropertyValue(value any) any {
	if millis, ok := value.(int64); ok {
		return time.UnixMilli(millis)
	}
	return value
}

// close removes the native listener and releases the service. It does not hold the lock while it releases the
// listener: the release waits for a running callback, which takes the lock.
func (s *OnDemandService) close() error {
	s.mu.Lock()
	l := s.nativeListener
	s.nativeListener = nil
	s.listeners = nil
	s.closed = true
	s.mu.Unlock()
	var err error
	if l != nil {
		err = s.handle.DetachPropertyListener(l)
	}
	return errors.Join(err, s.handle.Free())
}
