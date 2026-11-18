package api

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type connectorRecorder struct {
	changes chan MessageConnectorInfo
}

func (r *connectorRecorder) ConnectorChanged(connector MessageConnectorInfo) {
	select {
	case r.changes <- connector:
	default:
	}
}

func TestConnectorsOfLocalHub(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder().WithRole(LocalHub))
	connectors, err := endpoint.Connectors()
	if err != nil || len(connectors) != 0 {
		t.Fatalf("Connectors() of a local hub = %v, %v; want none", connectors, err)
	}
	lastError, err := endpoint.LastConnectionError()
	if err != nil || lastError != nil {
		t.Fatalf("LastConnectionError() of a local hub = %v, %v; want nil", lastError, err)
	}
}

// TestConnectorListenerAndLastError connects to a port where nothing listens and waits for the failed connect.
func TestConnectorListenerAndLastError(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	recorder := &connectorRecorder{changes: make(chan MessageConnectorInfo, 100)}
	if err := endpoint.AddMessageConnectorListener(recorder); err != nil {
		t.Fatalf("AddMessageConnectorListener: %v", err)
	}
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	var seen *ConnectionError
	for timeout := time.After(30 * time.Second); seen == nil; {
		select {
		case info := <-recorder.changes:
			if info.Kind != ClientSocketConnector || info.Address != "127.0.0.1:1" {
				t.Fatalf("ConnectorChanged(%+v): want the client socket connector of 127.0.0.1:1", info)
			}
			seen = info.LastError
		case <-timeout:
			t.Fatal("no connector change with an error within 30 s")
		}
	}
	if seen.Category != ConnectionErrorConnectFailed || seen.Dropped || seen.Count < 1 ||
		!strings.Contains(seen.Address, "127.0.0.1") {
		t.Errorf("the error of the listener = %+v; want a failed connect to 127.0.0.1", seen)
	}

	lastError, err := endpoint.LastConnectionError()
	if err != nil || lastError == nil {
		t.Fatalf("LastConnectionError() = %v, %v; want an error", lastError, err)
	}
	if lastError.Category != ConnectionErrorConnectFailed || lastError.Count < seen.Count ||
		time.Since(lastError.Time) > time.Minute {
		t.Errorf("LastConnectionError() = %+v; want a recent CONNECT_FAILED", lastError)
	}

	connectors, err := endpoint.Connectors()
	if err != nil || len(connectors) != 1 {
		t.Fatalf("Connectors() = %+v, %v; want one", connectors, err)
	}
	c := connectors[0]
	if c.Kind != ClientSocketConnector || c.Address != "127.0.0.1:1" || !c.Active || c.State == ConnectorConnected ||
		c.ConnectionCount != 0 || c.LastError == nil {
		t.Errorf("Connectors()[0] = %+v; want an active, not connected client socket connector with an error", c)
	}

	endpoint.RemoveMessageConnectorListener(recorder)
	if err := endpoint.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := endpoint.AddMessageConnectorListener(recorder); err == nil {
		t.Error("AddMessageConnectorListener after Close: no error")
	}
	if _, err := endpoint.Connectors(); err == nil {
		t.Error("Connectors after Close: no error")
	}
}

func TestConnectionErrorText(t *testing.T) {
	e := &ConnectionError{Address: "127.0.0.1:1", Category: ConnectionErrorConnectFailed, Message: "Connection refused"}
	if got, want := e.Error(), "dxfeed: connect to 127.0.0.1:1 failed (CONNECT_FAILED): Connection refused"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	e = &ConnectionError{Address: "host:7300", Category: ConnectionErrorHeartbeatTimeout, Dropped: true}
	if got, want := e.Error(), "dxfeed: connection to host:7300 dropped (HEARTBEAT_TIMEOUT)"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	for _, s := range []interface{ String() string }{ConnectionErrorCategory(99), MessageConnectorKind(99),
		MessageConnectorState(99)} {
		if !strings.HasSuffix(s.String(), "(99)") {
			t.Errorf("String() of an unknown value = %q", s.String())
		}
	}
}

func TestAddMessageConnectorListenerNil(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	if err := endpoint.AddMessageConnectorListener(nil); err == nil {
		t.Error("AddMessageConnectorListener(nil): no error")
	}
}

type panickingConnectorListener struct {
	calls atomic.Int32
}

func (l *panickingConnectorListener) ConnectorChanged(MessageConnectorInfo) {
	l.calls.Add(1)
	panic("connector listener failure")
}

// TestConnectorListenerPanicIsRecovered checks that a panic of a connector listener is logged and does not end the
// process or break the endpoint. A listener after the panicking one misses that change, so the recorder goes first.
func TestConnectorListenerPanicIsRecovered(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	recorder := &connectorRecorder{changes: make(chan MessageConnectorInfo, 100)}
	panicking := &panickingConnectorListener{}
	for _, l := range []MessageConnectorListener{recorder, panicking} {
		if err := endpoint.AddMessageConnectorListener(l); err != nil {
			t.Fatalf("AddMessageConnectorListener: %v", err)
		}
	}
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	select {
	case <-recorder.changes:
	case <-time.After(10 * time.Second):
		t.Fatal("no connector change within 10 s")
	}
	for deadline := time.Now().Add(5 * time.Second); panicking.calls.Load() == 0; time.Sleep(10 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatal("the panicking listener was not called within 5 s")
		}
	}
	if _, err := endpoint.Connectors(); err != nil {
		t.Errorf("Connectors after a listener panic: %v", err)
	}
}

// selfRemovingConnectorListener removes itself from ConnectorChanged.
type selfRemovingConnectorListener struct {
	endpoint *DXEndpoint
	calls    atomic.Int32
}

func (l *selfRemovingConnectorListener) ConnectorChanged(MessageConnectorInfo) {
	l.calls.Add(1)
	l.endpoint.RemoveMessageConnectorListener(l)
}

// TestConnectorListenerRemovesItselfFromCallback checks that a listener may remove itself without a deadlock.
func TestConnectorListenerRemovesItselfFromCallback(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	listener := &selfRemovingConnectorListener{endpoint: endpoint}
	if err := endpoint.AddMessageConnectorListener(listener); err != nil {
		t.Fatalf("AddMessageConnectorListener: %v", err)
	}
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	for deadline := time.Now().Add(10 * time.Second); listener.calls.Load() == 0; time.Sleep(10 * time.Millisecond) {
		if time.Now().After(deadline) {
			t.Fatal("the listener was not called within 10 s")
		}
	}
	// The endpoint still works after the removal from the callback.
	if _, err := endpoint.Connectors(); err != nil {
		t.Errorf("Connectors: %v", err)
	}
}

// TestRemoveMessageConnectorListenerKeepsOthers checks that after a removal only the other listener gets the changes.
func TestRemoveMessageConnectorListenerKeepsOthers(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	removed := &connectorRecorder{changes: make(chan MessageConnectorInfo, 100)}
	kept := &connectorRecorder{changes: make(chan MessageConnectorInfo, 100)}
	for _, l := range []MessageConnectorListener{removed, kept} {
		if err := endpoint.AddMessageConnectorListener(l); err != nil {
			t.Fatalf("AddMessageConnectorListener: %v", err)
		}
	}
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	for _, r := range []*connectorRecorder{removed, kept} {
		select {
		case <-r.changes:
		case <-time.After(10 * time.Second):
			t.Fatal("no connector change within 10 s")
		}
	}

	endpoint.RemoveMessageConnectorListener(removed)
	if err := endpoint.Connect("127.0.0.1:2"); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	for seen := false; !seen; {
		select {
		case info := <-kept.changes:
			seen = info.Address == "127.0.0.1:2"
		case <-time.After(10 * time.Second):
			t.Fatal("the kept listener got no change of 127.0.0.1:2 within 10 s")
		}
	}
	for {
		select {
		case info := <-removed.changes:
			if info.Address == "127.0.0.1:2" {
				t.Fatalf("the removed listener got a change after the removal: %+v", info)
			}
		default:
			return
		}
	}
}

// closingConnectorListener closes the endpoint with CloseAndAwaitTermination from the first ConnectorChanged.
type closingConnectorListener struct {
	endpoint *DXEndpoint
	done     chan error
	once     atomic.Bool
}

func (l *closingConnectorListener) ConnectorChanged(MessageConnectorInfo) {
	if !l.once.CompareAndSwap(false, true) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	l.done <- l.endpoint.CloseAndAwaitTermination(ctx)
}

// TestConnectorListenerClosesEndpointFromCallback checks that a listener may call CloseAndAwaitTermination without a
// deadlock between the release of the native listener and the running callback.
func TestConnectorListenerClosesEndpointFromCallback(t *testing.T) {
	endpoint := buildEndpoint(t, NewEndpointBuilder())
	listener := &closingConnectorListener{endpoint: endpoint, done: make(chan error, 1)}
	if err := endpoint.AddMessageConnectorListener(listener); err != nil {
		t.Fatalf("AddMessageConnectorListener: %v", err)
	}
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	select {
	case err := <-listener.done:
		if err != nil {
			t.Errorf("CloseAndAwaitTermination from a listener: %v", err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("CloseAndAwaitTermination from a listener did not return within 15 s")
	}
	if !endpoint.IsClosed() {
		t.Error("the endpoint is not closed")
	}
}
