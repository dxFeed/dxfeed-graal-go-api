// Package api connects to dxFeed: an endpoint (DXEndpoint) gives a feed to subscribe to market events and a publisher
// to publish them. It wraps the dxFeed Graal Native SDK.
package api

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
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

	// stateListeners is replaced, never changed in place, so UpdateState can call a snapshot without the lock.
	stateMu        sync.Mutex
	stateListeners []common.ConnectionStateListener
}

// UpdateState is called by the SDK when the connection state changes and notifies the state listeners; applications
// do not call it.
func (e *DXEndpoint) UpdateState(old common.ConnectionState, new common.ConnectionState) {
	e.stateMu.Lock()
	listeners := e.stateListeners
	e.stateMu.Unlock()
	for _, listener := range listeners {
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
	return NewEndpointBuilder().WithRole(role).WithProperties(properties).Build()
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

// Role returns the role of the endpoint.
func (e *DXEndpoint) Role() common.Role {
	return e.role
}

// State returns the state of the connection of the endpoint; Closed after Close.
func (e *DXEndpoint) State() (common.ConnectionState, error) {
	return e.endpointHandle.State()
}

// IsClosed reports whether the endpoint is closed.
func (e *DXEndpoint) IsClosed() bool {
	state, err := e.endpointHandle.State()
	return err == nil && state == common.Closed
}

// User sets the user name for the services that require credentials; call it before Connect, together with Password.
func (e *DXEndpoint) User(user string) error {
	return e.endpointHandle.User(user)
}

// Password sets the password for the services that require credentials; call it before Connect, together with User.
func (e *DXEndpoint) Password(password string) error {
	return e.endpointHandle.Password(password)
}

// Reconnect drops the connections of the endpoint and connects again to the same address, switching to another
// address of a multi-address connection. It does nothing if the endpoint is not connected with Connect or is
// disconnected with Disconnect.
func (e *DXEndpoint) Reconnect() error {
	return e.endpointHandle.Reconnect()
}

// Disconnect drops the connections of the endpoint; the state becomes NotConnected at once. The endpoint stays usable:
// Connect connects it again.
func (e *DXEndpoint) Disconnect() error {
	return e.endpointHandle.Disconnect()
}

// DisconnectAndClear drops the connections of the endpoint, as Disconnect, and clears the data that it stores, such
// as the last events of the feed.
func (e *DXEndpoint) DisconnectAndClear() error {
	return e.endpointHandle.DisconnectAndClear()
}

// EventTypes returns the event types that the endpoint supports. The list may hold types that the Go API cannot map to
// events yet: subscribing to them delivers nothing.
func (e *DXEndpoint) EventTypes() ([]eventcodes.EventCode, error) {
	clazzes, err := e.endpointHandle.EventTypes()
	if err != nil {
		return nil, err
	}
	codes := make([]eventcodes.EventCode, 0, len(clazzes))
	for _, clazz := range clazzes {
		codes = append(codes, eventcodes.EventCode(clazz))
	}
	return codes, nil
}

// maxAggregationPeriod is the longest aggregation period the SDK accepts.
const maxAggregationPeriod = 24 * time.Hour

// SetRequestedAggregationPeriod asks the servers to deliver the events of the endpoint conflated over the given period;
// zero asks for real-time delivery. The SDK keeps whole milliseconds: a shorter period is rounded up to 1 ms. The
// servers keep the period within their own bounds: see AggregationPeriodInfo. The period must be between 0 and 24
// hours. The property "dxfeed.requestedAggregationPeriod" and the address option "[requestedAggregationPeriod=5s]" do
// the same at creation.
func (e *DXEndpoint) SetRequestedAggregationPeriod(period time.Duration) error {
	if period < 0 || period > maxAggregationPeriod {
		return fmt.Errorf("dxfeed: aggregation period %v is out of range [0, %v]", period, maxAggregationPeriod)
	}
	millis := period.Milliseconds()
	if millis == 0 && period > 0 {
		millis = 1
	}
	return e.endpointHandle.SetRequestedAggregationPeriodMillis(millis)
}

// ClearRequestedAggregationPeriod withdraws the requested aggregation period: the servers use their default one.
func (e *DXEndpoint) ClearRequestedAggregationPeriod() error {
	return e.endpointHandle.SetRequestedAggregationPeriodMillis(-1)
}

// RequestedAggregationPeriod returns the aggregation period that the endpoint requests from the servers; ok is false
// when it requests none.
func (e *DXEndpoint) RequestedAggregationPeriod() (period time.Duration, ok bool, err error) {
	millis, err := e.endpointHandle.RequestedAggregationPeriodMillis()
	if err != nil || millis < 0 {
		return 0, false, err
	}
	return time.Duration(millis) * time.Millisecond, true, nil
}

// AggregationPeriodInfo is the aggregation period that the servers apply over all connections of an endpoint.
type AggregationPeriodInfo struct {
	// Known is false until a server reports the period on connection, and always for an endpoint without such a
	// server, such as a local hub or a publisher; Min and Max are zero then.
	Known bool
	// Min is the shortest period that the servers apply.
	Min time.Duration
	// Max is the longest period that the servers apply.
	Max time.Duration
}

// AggregationPeriodInfo returns the aggregation period that the servers apply.
func (e *DXEndpoint) AggregationPeriodInfo() (AggregationPeriodInfo, error) {
	lo, hi, err := e.endpointHandle.AggregationPeriodInfo()
	if err != nil || lo < 0 || hi < 0 {
		return AggregationPeriodInfo{}, err
	}
	return AggregationPeriodInfo{
		Known: true,
		Min:   time.Duration(lo) * time.Millisecond,
		Max:   time.Duration(hi) * time.Millisecond,
	}, nil
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

// AddListener adds a listener of connection state changes. It may be called from any goroutine, also from a listener;
// a state change that is being reported may or may not reach the new listener.
func (e *DXEndpoint) AddListener(listener common.ConnectionStateListener) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()
	listeners := make([]common.ConnectionStateListener, 0, len(e.stateListeners)+1)
	e.stateListeners = append(append(listeners, e.stateListeners...), listener)
}

// RemoveListener removes a listener added with AddListener. It may be called from any goroutine, also from a listener;
// a state change that is being reported may still reach the removed listener.
func (e *DXEndpoint) RemoveListener(listener common.ConnectionStateListener) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()
	for i, l := range e.stateListeners {
		if l == listener {
			listeners := make([]common.ConnectionStateListener, 0, len(e.stateListeners)-1)
			e.stateListeners = append(append(listeners, e.stateListeners[:i]...), e.stateListeners[i+1:]...)
			return
		}
	}
}
