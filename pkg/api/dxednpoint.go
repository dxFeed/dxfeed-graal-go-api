package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

const (
	Feed common.Role = iota
	OnDemandFeed
	StreamFeed
	Publisher
	StreamPublisher
	LocalHub
)

type DXEndpoint struct {
	role            common.Role
	endpointHandle  *native.DXEndpointHandle
	feedHandle      *DXFeed
	publisherHandle *DXPublisher
}

func NewEndpoint(role common.Role) (*DXEndpoint, error) {
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

func CreateEndpoint(role common.Role) (*DXEndpoint, error) {
	return NewEndpoint(role)
}

func NewEndpointWithProperties(role common.Role, properties map[string]string) (*DXEndpoint, error) {
	handle, err := native.NewDXEndpointHandleWithProperties(role, properties)
	if err != nil {
		return nil, err
	}

	e := &DXEndpoint{
		role:           role,
		endpointHandle: handle,
	}
	return e, nil
}

func (e *DXEndpoint) Connect(address string) error {
	return e.endpointHandle.Connect(address)
}

func (e *DXEndpoint) Close() error {
	return e.endpointHandle.Close()
}

func (e *DXEndpoint) GetFeed() (*DXFeed, error) {
	handle, err := e.endpointHandle.GetFeed()
	if err != nil {
		return nil, err
	}

	e.feedHandle = &DXFeed{feed: handle}
	return e.feedHandle, nil
}

func (e *DXEndpoint) GetPublisher() (*DXPublisher, error) {
	handle, err := e.endpointHandle.GetPublisher()
	if err != nil {
		return nil, err
	}
	e.publisherHandle = &DXPublisher{publisher: handle}
	return e.publisherHandle, nil
}

func (e *DXEndpoint) AwaitNotConnected() error {
	return e.endpointHandle.AwaitNotConnected()
}

func (e *DXEndpoint) CloseAndAwaitTermination() error {
	return e.endpointHandle.CloseAndAwaitTermination()
}

func (e *DXEndpoint) AwaitProcessed() error {
	return e.endpointHandle.AwaitProcessed()

}
