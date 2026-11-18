package api

import (
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

type stateRecorder struct {
	mu     sync.Mutex
	states []string
	seen   chan struct{}
	once   sync.Once
}

func newStateRecorder() *stateRecorder {
	return &stateRecorder{seen: make(chan struct{})}
}

func (r *stateRecorder) UpdateState(old common.ConnectionState, new common.ConnectionState) {
	r.mu.Lock()
	r.states = append(r.states, old.String()+" -> "+new.String())
	r.mu.Unlock()
	if new.String() == "Connecting" {
		r.once.Do(func() { close(r.seen) })
	}
}

func (r *stateRecorder) recorded() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.states...)
}

func TestEndpointStateListener(t *testing.T) {
	constructors := map[string]func() (*DXEndpoint, error){
		"NewEndpoint": func() (*DXEndpoint, error) { return NewEndpoint(Feed) },
		"NewEndpointWithProperties": func() (*DXEndpoint, error) {
			return NewEndpointWithProperties(Feed, map[string]string{"dxfeed.aggregationPeriod": "1s"})
		},
	}
	for name, create := range constructors {
		t.Run(name, func(t *testing.T) {
			endpoint, err := create()
			if err != nil {
				t.Fatalf("create endpoint: %v", err)
			}
			defer func() { _ = endpoint.Close() }()

			recorder := newStateRecorder()
			endpoint.AddListener(recorder)
			// Nothing listens on port 1, so the endpoint stays in Connecting.
			if err := endpoint.Connect("127.0.0.1:1"); err != nil {
				t.Fatalf("connect: %v", err)
			}

			select {
			case <-recorder.seen:
			case <-time.After(5 * time.Second):
				t.Fatalf("no transition to Connecting within 5 s, recorded: %v", recorder.recorded())
			}
		})
	}
}

type panickingStateListener struct{}

func (panickingStateListener) UpdateState(common.ConnectionState, common.ConnectionState) {
	panic("state listener failure")
}

// TestStateListenerPanicIsRecovered checks that a panic of a connection state listener is logged and does not end the
// process.
func TestStateListenerPanicIsRecovered(t *testing.T) {
	logged := &syncBuffer{}
	defer log.SetOutput(log.Writer())
	log.SetOutput(logged)

	endpoint, err := NewEndpoint(Feed)
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	defer func() { _ = endpoint.Close() }()
	recorder := newStateRecorder()
	endpoint.AddListener(recorder)
	// The panicking listener comes after the recorder, so the recorder still sees the transition.
	endpoint.AddListener(panickingStateListener{})
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("connect: %v", err)
	}
	select {
	case <-recorder.seen:
	case <-time.After(5 * time.Second):
		t.Fatalf("no transition to Connecting within 5 s, recorded: %v", recorder.recorded())
	}
	deadline := time.Now().Add(5 * time.Second)
	for !strings.Contains(logged.String(), "state listener failure") {
		if time.Now().After(deadline) {
			t.Fatalf("the panic was not logged; log: %q", logged.String())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestStateListenersChangeWhileStateChanges adds and removes listeners while the SDK reports state changes; under
// -race it fails if the listener list is not synchronized.
func TestStateListenersChangeWhileStateChanges(t *testing.T) {
	endpoint, err := NewEndpoint(Feed)
	if err != nil {
		t.Fatalf("create endpoint: %v", err)
	}
	recorder := newStateRecorder()
	endpoint.AddListener(recorder)
	removed := newStateRecorder()
	endpoint.AddListener(removed)
	endpoint.RemoveListener(removed)

	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				l := newStateRecorder()
				endpoint.AddListener(l)
				endpoint.RemoveListener(l)
			}
		}
	}()
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("connect: %v", err)
	}
	select {
	case <-recorder.seen:
	case <-time.After(5 * time.Second):
		t.Fatalf("no transition to Connecting within 5 s, recorded: %v", recorder.recorded())
	}
	if err := endpoint.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	close(done)
	wg.Wait()
	if states := removed.recorded(); len(states) != 0 {
		t.Errorf("a removed listener was called: %v", states)
	}
}
