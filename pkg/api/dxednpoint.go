package api

import (
	"dxfeed-graal-go-api/pkg/native"
)

const (
	Feed native.Role = iota
	OnDemandFeed
	StreamFeed
	Publisher
	StreamPublisher
	LocalHub
)

type DXEndpoint struct {
	role           native.Role
	endpointHandle *native.DXEndpointHandle
	feedHandle     *DXFeed
}

func NewEndpoint(role native.Role) (*DXEndpoint, error) {
	handle, err := native.NewDXEndpointHandle(role)
	if err != nil {
		return nil, err
	}

	e := &DXEndpoint{
		role:           role,
		endpointHandle: handle,
	}
	return e, nil
}

func CreateEndpoint(role native.Role) (*DXEndpoint, error) {
	return NewEndpoint(role)
}

func (e *DXEndpoint) Connect(address string) error {
	return e.endpointHandle.Connect(address)
}

func (e *DXEndpoint) Close() {
}

func (e *DXEndpoint) GetFeed() (*DXFeed, error) {
	handle, err := e.endpointHandle.GetFeed()
	if err != nil {
		return nil, err
	}

	e.feedHandle = &DXFeed{feed: handle}
	return e.feedHandle, nil
}
