package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

type DXPublisher struct {
	publisher *native.DXPublisherHandle
}

func (p *DXPublisher) Publish(events []interface{}) error {
	return p.publisher.Publish(events)
}
