package api

type DXEndpoint struct {
	role Role
	feed *DXFeed
}

type Role int32

const (
	Feed Role = iota
	OnDemandFeed
	StreamFeed
	Publisher
	StreamPublisher
	LocalHub
)

func NewEndpoint(role Role) *DXEndpoint {
	return &DXEndpoint{role: role}
}

func CreateEndpoint(role Role) *DXEndpoint {
	return NewEndpoint(role)
}

func (e *DXEndpoint) Connect(address string) *DXEndpoint {
	return e
}

func (e *DXEndpoint) Close() {
}

func (e *DXEndpoint) GetFeed() *DXFeed {
	e.feed = &DXFeed{}
	return e.feed
}
