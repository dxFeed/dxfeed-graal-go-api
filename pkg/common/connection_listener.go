// Package common holds the types shared by the other packages: roles, connection states, listeners and JavaError.
package common

import "fmt"

// ConnectionState is the state of the connection of an endpoint.
type ConnectionState int32

// String returns the name of the state.
func (c ConnectionState) String() string {
	switch c {
	case NotConnected:
		return "Not connected"
	case Connecting:
		return "Connecting"
	case Connected:
		return "Connected"
	case Closed:
		return "Closed"
	default:
		return fmt.Sprintf("Unsupported connection state %d", int(c))
	}
}

// The connection states.
const (
	// NotConnected is the state of an endpoint that is not connected or has disconnected.
	NotConnected ConnectionState = iota
	// Connecting is the state while the endpoint connects.
	Connecting
	// Connected is the state of a connected endpoint.
	Connected
	// Closed is the state of a closed endpoint.
	Closed
)

// ConnectionStateListener is notified when the connection state of an endpoint changes; see DXEndpoint.AddListener in
// package api. It is called from an SDK thread.
type ConnectionStateListener interface {
	UpdateState(old ConnectionState, new ConnectionState)
}
