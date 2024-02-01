package common

import "fmt"

type ConnectionState int32

func (c ConnectionState) String() string {
	switch c {
	case notConnected:
		return "Not connected"
	case connecting:
		return "Connecting"
	case connected:
		return "Connected"
	case closed:
		return "Closed"
	default:
		return fmt.Sprintf("Unsupproted connection state %d", int(c))
	}
}

const (
	notConnected ConnectionState = iota
	connecting
	connected
	closed
)

type ConnectionStateListener interface {
	UpdateState(old ConnectionState, new ConnectionState)
}
