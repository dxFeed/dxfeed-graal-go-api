package api

import "dxfeed-graal-go-api/pkg/native"

const (
	Feed native.Role = iota
	OnDemandFeed
	StreamFeed
	Publisher
	StreamPublisher
	LocalHub
)

type DXEndpoint struct {
	role     native.Role
	endpoint *native.DXEndpointHandle
	feed     *DXFeed
}

func NewEndpoint(role native.Role) *DXEndpoint {
	return &DXEndpoint{role: role, endpoint: native.NewDXEndpointHandle(role)}
}

func CreateEndpoint(role native.Role) *DXEndpoint {
	return NewEndpoint(role)
}

func (e *DXEndpoint) Connect(address string) *DXEndpoint {
	e.endpoint.Connect(address)
	return e
}

func (e *DXEndpoint) Close() {
}

func (e *DXEndpoint) GetFeed() *DXFeed {
	e.feed = &DXFeed{feed: e.endpoint.GetFeed()}
	return e.feed
}
