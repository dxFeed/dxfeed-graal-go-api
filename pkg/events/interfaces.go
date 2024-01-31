package events

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"

type StringConverter interface {
	String() string
}

type EventType interface {
	Type() eventcodes.EventCode
}
