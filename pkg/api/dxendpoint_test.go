package api

import (
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
