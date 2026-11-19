package events

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"

// StringConverter is implemented by values that can be rendered as a string.
type StringConverter interface {
	String() string
}

// EventType is implemented by all event types and reports the event code of the event.
type EventType interface {
	Type() eventcodes.EventCode
}
