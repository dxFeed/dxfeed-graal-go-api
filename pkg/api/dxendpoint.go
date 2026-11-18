// Package api connects to dxFeed: an endpoint (DXEndpoint) gives a feed to subscribe to market events and a publisher
// to publish them. It wraps the dxFeed Graal Native SDK.
package api

import (
	"context"
	"io"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// Roles of an endpoint, as in the dxFeed Java API DXEndpoint.Role.
const (
	// Feed receives events from a remote address with conflation of fast updates.
	Feed common.Role = iota
	// OnDemandFeed is a feed for historical data replay on demand.
	OnDemandFeed
	// StreamFeed receives every event, without conflation.
	StreamFeed
	// Publisher publishes events to a remote address.
	Publisher
	// StreamPublisher publishes every event, without conflation.
	StreamPublisher
	// LocalHub connects a feed and a publisher in the same process, without a network connection.
	LocalHub
)

// The endpoint and the subscription own native resources: close them.
var (
	_ io.Closer = (*DXEndpoint)(nil)
	_ io.Closer = (*DXFeedSubscription)(nil)
)

// DXEndpoint is a connection to dxFeed: create it with NewEndpoint, connect it with Connect and close it with Close.
type DXEndpoint struct {
	role            common.Role
	endpointHandle  *native.DXEndpointHandle
	feedHandle      *DXFeed
	publisherHandle *DXPublisher

	stateListenerList []common.ConnectionStateListener
}

// UpdateState is called by the SDK when the connection state changes and notifies the state listeners; applications
// do not call it.
func (e *DXEndpoint) UpdateState(old common.ConnectionState, new common.ConnectionState) {
	for _, listener := range e.stateListenerList {
		listener.UpdateState(old, new)
	}
}

// NewEndpoint creates an endpoint with the given role and the default properties.
func NewEndpoint(role common.Role) (*DXEndpoint, error) {
	handle, err := native.NewDXEndpointHandle(role)
	return newEndpoint(role, handle, err)
}

// newEndpoint wraps a native endpoint and attaches the state listener; every constructor goes through it.
func newEndpoint(role common.Role, handle *native.DXEndpointHandle, err error) (*DXEndpoint, error) {
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

// NewEndpointWithProperties creates an endpoint with the given role and properties, such as "dxfeed.address".
func NewEndpointWithProperties(role common.Role, properties map[string]string) (*DXEndpoint, error) {
	handle, err := native.NewDXEndpointHandleWithProperties(role, properties)
	return newEndpoint(role, handle, err)
}

// Connect connects the endpoint to the address, for example "demo.dxfeed.com:7300", and returns at once; the
// connection is made in the background.
func (e *DXEndpoint) Connect(address string) error {
	return e.endpointHandle.Connect(address)
}

// Close closes the endpoint and its connections and releases its native resources at once, without waiting for
// pending data.
func (e *DXEndpoint) Close() error {
	return e.endpointHandle.Close()
}

// GetFeed returns the feed of the endpoint to create subscriptions.
func (e *DXEndpoint) GetFeed() (*DXFeed, error) {
	handle, err := e.endpointHandle.GetFeed()
	if err != nil {
		return nil, err
	}

	e.feedHandle = &DXFeed{feed: handle}
	return e.feedHandle, nil
}

// GetPublisher returns the publisher of the endpoint to publish events.
func (e *DXEndpoint) GetPublisher() (*DXPublisher, error) {
	handle, err := e.endpointHandle.GetPublisher()
	if err != nil {
		return nil, err
	}
	e.publisherHandle = &DXPublisher{publisher: handle}
	return e.publisherHandle, nil
}

// AwaitNotConnected waits until the endpoint is not connected. If ctx ends first, it returns ctx.Err(); the wait in the
// SDK cannot be interrupted and goes on in the background.
func (e *DXEndpoint) AwaitNotConnected(ctx context.Context) error {
	return awaitWithContext(ctx, e.endpointHandle.AwaitNotConnected)
}

// CloseAndAwaitTermination closes the endpoint and waits until its connections are closed and its pending data is
// processed. The endpoint is closed even if ctx is already done; if ctx ends first, the call returns ctx.Err() and the
// endpoint goes on closing in the background.
func (e *DXEndpoint) CloseAndAwaitTermination(ctx context.Context) error {
	return runWithContext(ctx, e.endpointHandle.CloseAndAwaitTermination)
}

// AwaitProcessed waits until all data that is published or subscribed before the call is processed. If ctx ends first,
// it returns ctx.Err(); the wait in the SDK cannot be interrupted and goes on in the background.
func (e *DXEndpoint) AwaitProcessed(ctx context.Context) error {
	return awaitWithContext(ctx, e.endpointHandle.AwaitProcessed)
}

// awaitWithContext runs a blocking SDK wait unless ctx is already done; see runWithContext.
func awaitWithContext(ctx context.Context, wait func() error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return runWithContext(ctx, wait)
}

// runWithContext runs a blocking SDK call and returns its error, or ctx.Err() if ctx ends first; when both happen at
// once, either may be returned. The call keeps running after that: the SDK offers no way to cancel it.
func runWithContext(ctx context.Context, wait func() error) error {
	done := make(chan error, 1)
	go func() { done <- wait() }()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// AddListener adds a listener of connection state changes. It is not safe to call concurrently with the other listener
// methods.
func (e *DXEndpoint) AddListener(listener common.ConnectionStateListener) {
	e.stateListenerList = append(e.stateListenerList, listener)
}

// RemoveListener removes a listener added with AddListener.
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
