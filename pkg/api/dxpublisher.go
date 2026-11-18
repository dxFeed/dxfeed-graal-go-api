package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// DXPublisher publishes events; get it with DXEndpoint.GetPublisher.
type DXPublisher struct {
	publisher *native.DXPublisherHandle
}

// Publish publishes events, such as *quote.Quote, to the subscribers of the endpoint.
func (p *DXPublisher) Publish(events []interface{}) error {
	return p.publisher.Publish(events)
}
