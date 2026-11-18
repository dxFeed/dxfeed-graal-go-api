package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// ConnectionErrorCategory tells what kind of failure a ConnectionError is.
type ConnectionErrorCategory int32

// The categories of connection errors, as ConnectionError::Category of the C++ API.
const (
	// ConnectionErrorOther is an error of no other category.
	ConnectionErrorOther ConnectionErrorCategory = iota
	// ConnectionErrorDNS is a host name that could not be resolved.
	ConnectionErrorDNS
	// ConnectionErrorConnectFailed is a connect attempt that failed, for example refused by the host.
	ConnectionErrorConnectFailed
	// ConnectionErrorReadTimeout is a read that timed out.
	ConnectionErrorReadTimeout
	// ConnectionErrorHeartbeatTimeout is a connection dropped because the other side sent no heartbeat in time.
	ConnectionErrorHeartbeatTimeout
	// ConnectionErrorClosedByRemote is a connection closed by the other side.
	ConnectionErrorClosedByRemote
	// ConnectionErrorTLS is a failure of TLS, such as a handshake or certificate error.
	ConnectionErrorTLS
	// ConnectionErrorIO is another I/O error, for example a connection reset.
	ConnectionErrorIO
)

// String returns the name of the category, for example "CONNECT_FAILED".
func (c ConnectionErrorCategory) String() string {
	switch c {
	case ConnectionErrorOther:
		return "OTHER"
	case ConnectionErrorDNS:
		return "DNS"
	case ConnectionErrorConnectFailed:
		return "CONNECT_FAILED"
	case ConnectionErrorReadTimeout:
		return "READ_TIMEOUT"
	case ConnectionErrorHeartbeatTimeout:
		return "HEARTBEAT_TIMEOUT"
	case ConnectionErrorClosedByRemote:
		return "CLOSED_BY_REMOTE"
	case ConnectionErrorTLS:
		return "TLS"
	case ConnectionErrorIO:
		return "IO"
	}
	return fmt.Sprintf("ConnectionErrorCategory(%d)", int32(c))
}

// ConnectionError is a failed connect attempt or a dropped connection of an endpoint: see
// DXEndpoint.LastConnectionError and MessageConnectorInfo.LastError. It is a copy that does not change later.
type ConnectionError struct {
	// Time is when the error happened.
	Time time.Time
	// Address is the address of the connection, for example "127.0.0.1:7300".
	Address string
	// Category tells what kind of failure it is.
	Category ConnectionErrorCategory
	// Dropped is true if an established connection was dropped, false if a connect attempt failed.
	Dropped bool
	// ExceptionClassName is the class name of the Java exception, for example "java.net.ConnectException"; empty when
	// there is none.
	ExceptionClassName string
	// Message is the message of the exception; empty when there is none.
	Message string
	// Count is the number of errors the connector has recorded, this one included: it tells a new error from one seen
	// before.
	Count int64
}

var _ error = (*ConnectionError)(nil)

// Error describes the error, for example "dxfeed: connect to 127.0.0.1:1 failed (CONNECT_FAILED): Connection refused".
func (e *ConnectionError) Error() string {
	what := "connect to " + e.Address + " failed"
	if e.Dropped {
		what = "connection to " + e.Address + " dropped"
	}
	if e.Message == "" {
		return fmt.Sprintf("dxfeed: %s (%v)", what, e.Category)
	}
	return fmt.Sprintf("dxfeed: %s (%v): %s", what, e.Category, e.Message)
}

func newConnectionError(d *native.ConnectionErrorData) *ConnectionError {
	if d == nil {
		return nil
	}
	return &ConnectionError{
		Time:               time.UnixMilli(d.TimeMillis),
		Address:            d.Address,
		Category:           ConnectionErrorCategory(d.Category),
		Dropped:            d.Dropped,
		ExceptionClassName: d.ExceptionClassName,
		Message:            d.Message,
		Count:              d.Count,
	}
}

// MessageConnectorKind is the kind of a message connector.
type MessageConnectorKind int32

// The kinds of message connectors, as MessageConnectorInfo::Kind of the C++ API.
const (
	// ClientSocketConnector connects to an address, for example "demo.dxfeed.com:7300".
	ClientSocketConnector MessageConnectorKind = iota
	// ServerSocketConnector accepts connections on a port, for example ":7400".
	ServerSocketConnector
	// OtherConnector is any other connector, such as one that reads a file.
	OtherConnector
)

// String returns the name of the kind, for example "CLIENT_SOCKET".
func (k MessageConnectorKind) String() string {
	switch k {
	case ClientSocketConnector:
		return "CLIENT_SOCKET"
	case ServerSocketConnector:
		return "SERVER_SOCKET"
	case OtherConnector:
		return "OTHER"
	}
	return fmt.Sprintf("MessageConnectorKind(%d)", int32(k))
}

// MessageConnectorState is the state of a message connector.
type MessageConnectorState int32

// The states of message connectors, as MessageConnectorInfo::State of the C++ API.
const (
	// ConnectorConnecting is a connector that connects or waits to reconnect.
	ConnectorConnecting MessageConnectorState = iota
	// ConnectorConnected is a connector with an open connection.
	ConnectorConnected
	// ConnectorDisconnected is a stopped connector.
	ConnectorDisconnected
)

// String returns the name of the state, for example "CONNECTED".
func (s MessageConnectorState) String() string {
	switch s {
	case ConnectorConnecting:
		return "CONNECTING"
	case ConnectorConnected:
		return "CONNECTED"
	case ConnectorDisconnected:
		return "DISCONNECTED"
	}
	return fmt.Sprintf("MessageConnectorState(%d)", int32(s))
}

// MessageConnectorStatistics are the counters of a message connector since it was created.
type MessageConnectorStatistics struct {
	// ReadBytes is the number of bytes read.
	ReadBytes int64
	// WrittenBytes is the number of bytes written.
	WrittenBytes int64
	// ReadDataRecords is the number of data records read.
	ReadDataRecords int64
	// WrittenDataRecords is the number of data records written.
	WrittenDataRecords int64
	// ActiveConnectionCount is the number of the open connections.
	ActiveConnectionCount int64
	// ClosedConnectionCount is the number of the closed connections.
	ClosedConnectionCount int64
}

// MessageConnectorInfo is a snapshot of a message connector of an endpoint: one connection setting, such as an address
// given to Connect. Each value is current when it is read; the values are not read together as one consistent set.
type MessageConnectorInfo struct {
	// Name is the name of the connector, as in the log, for example "ClientSocket-Distributor".
	Name string
	// Kind is the kind of the connector.
	Kind MessageConnectorKind
	// Address is the address the connector was configured with, credentials masked; a server socket connector of
	// ":7400" has "*:7400".
	Address string
	// CurrentAddresses are the "host:port" addresses a client socket connector is connected to now, one per connection;
	// empty when it is not connected and for the other kinds.
	CurrentAddresses []string
	// State is the state of the connector.
	State MessageConnectorState
	// Active is true if the connector is started (it connects or reconnects), false if it is stopped.
	Active bool
	// ConnectionCount is the number of the open connections of the connector.
	ConnectionCount int
	// ClosedConnectionCount is the number of the connections closed since the connector was created; it grows on each
	// reconnect.
	ClosedConnectionCount int64
	// Statistics are the counters of the connector.
	Statistics MessageConnectorStatistics
	// LastError is the last connection error of the connector, or nil when it has had none.
	LastError *ConnectionError
}

func newMessageConnectorInfo(d native.MessageConnectorData) MessageConnectorInfo {
	return MessageConnectorInfo{
		Name:                  d.Name,
		Kind:                  MessageConnectorKind(d.Kind),
		Address:               d.Address,
		CurrentAddresses:      d.CurrentAddresses,
		State:                 MessageConnectorState(d.State),
		Active:                d.Active,
		ConnectionCount:       int(d.ConnectionCount),
		ClosedConnectionCount: d.ClosedConnectionCount,
		Statistics: MessageConnectorStatistics{
			ReadBytes:             d.ReadBytes,
			WrittenBytes:          d.WrittenBytes,
			ReadDataRecords:       d.ReadDataRecords,
			WrittenDataRecords:    d.WrittenDataRecords,
			ActiveConnectionCount: d.ActiveConnections,
			ClosedConnectionCount: d.ClosedConnections,
		},
		LastError: newConnectionError(d.LastError),
	}
}

// MessageConnectorListener is told about the state changes of the message connectors of an endpoint: see
// DXEndpoint.AddMessageConnectorListener.
type MessageConnectorListener interface {
	// ConnectorChanged is called after a connector changed its state, with a snapshot of the connector read in the
	// call; fast changes of one connector may come as one call.
	ConnectorChanged(connector MessageConnectorInfo)
}

// Connectors returns a snapshot of each message connector of the endpoint: a client socket connector for each address
// given to Connect, a server socket connector for each port it listens on. A local hub has none.
func (e *DXEndpoint) Connectors() ([]MessageConnectorInfo, error) {
	data, err := e.endpointHandle.Connectors()
	if err != nil {
		return nil, err
	}
	result := make([]MessageConnectorInfo, len(data))
	for i, d := range data {
		result[i] = newMessageConnectorInfo(d)
	}
	return result, nil
}

// LastConnectionError returns the newest connection error of all the connectors of the endpoint, or nil when there is
// none. A successful connection does not clear it: compare the Time or the Count with those seen before. Read it in a
// MessageConnectorListener, which is called after the error is recorded; a ConnectionStateListener may be called
// before that.
func (e *DXEndpoint) LastConnectionError() (*ConnectionError, error) {
	d, err := e.endpointHandle.LastConnectionError()
	if err != nil {
		return nil, err
	}
	return newConnectionError(d), nil
}

// AddMessageConnectorListener adds a listener of the state changes of the message connectors of the endpoint, also of
// those that a later Connect creates. The listeners are called on one SDK thread, one call at a time, so a listener
// should return quickly; it may call the methods of the endpoint, Connect included. The current states are not passed
// to a new listener: call Connectors after adding it. It may be called from any goroutine, also from a listener, which
// may also close the endpoint. A listener that panics is logged and the other listeners miss that change.
func (e *DXEndpoint) AddMessageConnectorListener(listener MessageConnectorListener) error {
	if listener == nil {
		return errors.New("dxfeed: the message connector listener is nil")
	}
	e.connectorMu.Lock()
	defer e.connectorMu.Unlock()
	if e.connectorsClosed {
		return errors.New("dxfeed: the endpoint is closed")
	}
	if e.nativeConnectorListener == nil {
		// One native listener per endpoint passes the changes to all the Go listeners.
		l, err := e.endpointHandle.AttachMessageConnectorListener(e.connectorChanged)
		if err != nil {
			return err
		}
		e.nativeConnectorListener = l
	}
	listeners := make([]MessageConnectorListener, 0, len(e.connectorListeners)+1)
	e.connectorListeners = append(append(listeners, e.connectorListeners...), listener)
	return nil
}

// RemoveMessageConnectorListener removes a listener added with AddMessageConnectorListener. It may be called from any
// goroutine, also from a listener; a change that is being reported may still reach the removed listener.
func (e *DXEndpoint) RemoveMessageConnectorListener(listener MessageConnectorListener) {
	e.connectorMu.Lock()
	defer e.connectorMu.Unlock()
	for i, l := range e.connectorListeners {
		if l == listener {
			listeners := make([]MessageConnectorListener, 0, len(e.connectorListeners)-1)
			e.connectorListeners = append(append(listeners, e.connectorListeners[:i]...), e.connectorListeners[i+1:]...)
			return
		}
	}
}

// connectorChanged is the callback of the native listener: it passes the snapshot to a copy of the listeners.
func (e *DXEndpoint) connectorChanged(d native.MessageConnectorData, err error) {
	if err != nil {
		// The connector could not be read, for example while the endpoint closes: there is nothing to report.
		return
	}
	e.connectorMu.Lock()
	listeners := e.connectorListeners
	e.connectorMu.Unlock()
	if len(listeners) == 0 {
		return
	}
	info := newMessageConnectorInfo(d)
	for _, listener := range listeners {
		listener.ConnectorChanged(info)
	}
}

// closeConnectorListener removes the native listener before the endpoint closes. It does not hold the lock while it
// releases the listener: the release waits for a running callback, which takes the lock.
func (e *DXEndpoint) closeConnectorListener() error {
	e.connectorMu.Lock()
	l := e.nativeConnectorListener
	e.nativeConnectorListener = nil
	e.connectorListeners = nil
	e.connectorsClosed = true
	e.connectorMu.Unlock()
	if l == nil {
		return nil
	}
	return e.endpointHandle.DetachMessageConnectorListener(l)
}
