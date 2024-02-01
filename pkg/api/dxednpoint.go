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

	stateListenerList []common.ConnectionStateListener
}

func (e *DXEndpoint) UpdateState(old common.ConnectionState, new common.ConnectionState) {
	for _, listener := range e.stateListenerList {
		listener.UpdateState(old, new)
	}
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
	err = handle.AttachListener(e)
	if err != nil {
		_ = handle.Close()
		return nil, err
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

func (e *DXEndpoint) AddListener(listener common.ConnectionStateListener) {
	e.stateListenerList = append(e.stateListenerList, listener)
}

func (e *DXEndpoint) RemoveListener(listener common.ConnectionStateListener) {
	e.stateListenerList = removeStateListenerFromSlice(e.stateListenerList, listener)
}

func removeStateListenerFromSlice(list []common.ConnectionStateListener, observerToRemove common.ConnectionStateListener) []common.ConnectionStateListener {
	listLength := len(list)
	for i, observer := range list {
		if observerToRemove == observer {
			list[listLength-1], list[i] = list[i], list[listLength-1]
			return list[:listLength-1]
		}
	}
	return list
}
