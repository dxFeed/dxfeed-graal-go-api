package api

import (
	"strings"
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
