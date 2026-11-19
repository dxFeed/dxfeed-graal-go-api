// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// propertyRecorder records the property changes of an on-demand service.
type propertyRecorder struct {
	mu      sync.Mutex
	changes []string
}

func (r *propertyRecorder) PropertyChanged(name string, oldValue, newValue any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.changes = append(r.changes, fmt.Sprintf("%s %v -> %v", name, format(oldValue), format(newValue)))
}

func format(value any) any {
	if t, ok := value.(time.Time); ok {
		return t.UnixMilli()
	}
	return value
}

func (r *propertyRecorder) has(change string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.changes {
		if c == change {
			return true
		}
	}
	return false
}

// waitFor waits until the recorder has the change: a listener may be called on an SDK thread.
func (r *propertyRecorder) waitFor(t *testing.T, change string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !r.has(change) {
		if time.Now().After(deadline) {
			r.mu.Lock()
			defer r.mu.Unlock()
			t.Fatalf("no property change %q within 5 s; got %q", change, r.changes)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type onDemandState struct {
	supported, replay, clear bool
	time                     int64
	speed                    float64
}

func readOnDemandState(t *testing.T, s *OnDemandService) onDemandState {
	t.Helper()
	var st onDemandState
	var err error
	check := func(e error) {
		if e != nil && err == nil {
			err = e
		}
	}
	var e error
	st.supported, e = s.IsReplaySupported()
	check(e)
	st.replay, e = s.IsReplay()
	check(e)
	st.clear, e = s.IsClear()
	check(e)
	tm, e := s.Time()
	check(e)
	st.time = tm.UnixMilli()
	st.speed, e = s.Speed()
	check(e)
	if err != nil {
		t.Fatalf("read the on-demand state: %v", err)
	}
	return st
}

// TestOnDemandServiceStateTransitions drives the on-demand service of an endpoint that is connected to an unreachable
// on-demand provider, so no data comes: the state follows the calls alone.
func TestOnDemandServiceStateTransitions(t *testing.T) {
	endpoint, err := NewEndpoint(OnDemandFeed)
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	defer func() { _ = endpoint.Close() }()
	service, err := endpoint.OnDemandService()
	if err != nil {
		t.Fatalf("on-demand service: %v", err)
	}
	if again, err := endpoint.OnDemandService(); err != nil || again != service {
		t.Fatalf("second OnDemandService call returned %p, %v; want the same service %p", again, err, service)
	}
	if service.Endpoint() != endpoint {
		t.Fatal("Endpoint does not return the endpoint of the service")
	}
	recorder := &propertyRecorder{}
	if err := service.AddPropertyChangeListener(recorder); err != nil {
		t.Fatalf("add listener: %v", err)
	}

	if got, want := readOnDemandState(t, service), (onDemandState{}); got != want {
		t.Fatalf("initial state %+v, want %+v", got, want)
	}
	from := time.UnixMilli(1700000000000)
	if err := service.Replay(from); err == nil {
		t.Fatal("Replay without an on-demand address returned no error")
	}

	if err := endpoint.Connect("(ondemand:127.0.0.1:1)"); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if got := readOnDemandState(t, service); !got.supported || got.replay {
		t.Fatalf("state after connect %+v, want the replay supported and not started", got)
	}
	recorder.waitFor(t, "replaySupported false -> true")

	if err := service.ReplayWithSpeed(from, 2); err != nil {
		t.Fatalf("replay: %v", err)
	}
	// The provider is unreachable, so the replay waits for data; the time does not go back from the start.
	if got := readOnDemandState(t, service); !got.supported || !got.replay || got.clear || got.speed != 2 ||
		got.time < from.UnixMilli() {
		t.Fatalf("state after replay %+v, want the replay mode from %d at the speed 2", got, from.UnixMilli())
	}
	recorder.waitFor(t, "replay false -> true")
	recorder.waitFor(t, "speed 0 -> 2")
	recorder.waitFor(t, fmt.Sprintf("time 0 -> %d", from.UnixMilli()))

	if err := service.SetSpeed(5); err != nil {
		t.Fatalf("set speed: %v", err)
	}
	if speed, _ := service.Speed(); speed != 5 {
		t.Fatalf("speed %v after SetSpeed(5)", speed)
	}
	recorder.waitFor(t, "speed 2 -> 5")
	if err := service.SetSpeed(-1); err == nil {
		t.Fatal("SetSpeed(-1) returned no error")
	}

	if err := service.Pause(); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if got := readOnDemandState(t, service); !got.replay || got.speed != 0 {
		t.Fatalf("state after pause %+v, want the replay mode at the speed 0", got)
	}

	if err := service.StopAndResume(); err != nil {
		t.Fatalf("stop and resume: %v", err)
	}
	if got := readOnDemandState(t, service); got.replay || got.clear {
		t.Fatalf("state after StopAndResume %+v, want neither the replay nor the clear mode", got)
	}
	recorder.waitFor(t, "replay true -> false")

	if err := service.StopAndClear(); err != nil {
		t.Fatalf("stop and clear: %v", err)
	}
	if got := readOnDemandState(t, service); got.replay || !got.clear || got.time < from.UnixMilli() {
		t.Fatalf("state after StopAndClear %+v, want the clear mode at the last replayed time", got)
	}

	// The changes of one service are reported one at a time, in order: once the remaining listener has a change, a
	// removed listener would have had it too.
	remaining := &propertyRecorder{}
	if err := service.AddPropertyChangeListener(remaining); err != nil {
		t.Fatalf("add the second listener: %v", err)
	}
	service.RemovePropertyChangeListener(recorder)
	recorder.mu.Lock()
	recorder.changes = nil
	recorder.mu.Unlock()
	if err := service.ReplayWithSpeed(from, 3); err != nil {
		t.Fatalf("replay after removing the listener: %v", err)
	}
	remaining.waitFor(t, "speed 0 -> 3")
	if recorder.has("speed 0 -> 3") {
		t.Fatal("a removed listener got a property change")
	}

	if err := endpoint.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := service.IsReplay(); err == nil {
		t.Fatal("IsReplay after Close of the endpoint returned no error")
	}
	if _, err := endpoint.OnDemandService(); err == nil {
		t.Fatal("OnDemandService after Close of the endpoint returned no error")
	}
	if err := service.AddPropertyChangeListener(recorder); err == nil {
		t.Fatal("AddPropertyChangeListener after Close of the endpoint returned no error")
	}
}

// TestOnDemandServiceRejectsOtherRoles checks that only the Feed and OnDemandFeed endpoints have the service.
func TestOnDemandServiceRejectsOtherRoles(t *testing.T) {
	endpoint, err := NewEndpoint(Publisher)
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	defer func() { _ = endpoint.Close() }()
	if _, err := endpoint.OnDemandService(); err == nil {
		t.Fatal("OnDemandService of a Publisher endpoint returned no error")
	}
	if err := (*OnDemandService)(nil).AddPropertyChangeListener(nil); err == nil {
		t.Fatal("AddPropertyChangeListener(nil) returned no error")
	}
}
