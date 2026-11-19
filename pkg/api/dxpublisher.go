// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"sync"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// DXPublisher publishes events; get it with DXEndpoint.GetPublisher.
type DXPublisher struct {
	publisher *native.DXPublisherHandle

	mu            sync.Mutex
	subscriptions map[eventcodes.EventCode]*ObservableSubscription
}

// Publish publishes events, such as *quote.Quote, to the subscribers of the endpoint.
func (p *DXPublisher) Publish(events []interface{}) error {
	return p.publisher.Publish(events)
}

// Subscription returns the symbols that the feeds are subscribed to for the event type, so that a data provider
// publishes only the events that someone needs; later calls return the same subscription. It is valid until the
// endpoint is closed.
func (p *DXPublisher) Subscription(eventType eventcodes.EventCode) (*ObservableSubscription, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if sub := p.subscriptions[eventType]; sub != nil {
		return sub, nil
	}
	handle, err := p.publisher.Subscription(eventType.NativeCode())
	if err != nil {
		return nil, err
	}
	sub := &ObservableSubscription{
		sub:       handle,
		listeners: newChangeListenerRegistry(),
	}
	if p.subscriptions == nil {
		p.subscriptions = make(map[eventcodes.EventCode]*ObservableSubscription)
	}
	p.subscriptions[eventType] = sub
	return sub, nil
}

// ObservableSubscription is the set of the symbols that the feeds are subscribed to for an event type; get it with
// DXPublisher.Subscription and follow it with a change listener. The symbols include *osub.WildcardSymbol for a
// subscription to all symbols and *osub.TimeSeriesSubscriptionSymbol for a time series subscription.
type ObservableSubscription struct {
	sub       *native.ObservableSubscriptionHandle
	listeners changeListenerRegistry
}

// IsClosed reports whether the subscription is closed: it is closed with the endpoint, and a closed subscription
// notifies no listener.
func (s *ObservableSubscription) IsClosed() (bool, error) {
	return s.sub.IsClosed()
}

// EventTypes returns the event types of the subscription.
func (s *ObservableSubscription) EventTypes() ([]eventcodes.EventCode, error) {
	classes, err := s.sub.EventTypes()
	if err != nil {
		return nil, err
	}
	return eventCodes(classes), nil
}

// ContainsEventType reports whether the subscription has the event type.
func (s *ObservableSubscription) ContainsEventType(eventType eventcodes.EventCode) (bool, error) {
	return s.sub.ContainsEventType(eventType.NativeCode())
}

// AddChangeListener adds a listener of the changes of the subscribed symbols; see
// osub.ObservableSubscriptionChangeListener. If symbols are subscribed, the listener gets them in SymbolsAdded before
// AddChangeListener returns. Adding a listener to a closed subscription does nothing. Remove the listener with
// RemoveChangeListener to release it, also after the subscription is closed. After the endpoint is closed,
// AddChangeListener, EventTypes and ContainsEventType return an error and IsClosed returns true.
func (s *ObservableSubscription) AddChangeListener(listener osub.ObservableSubscriptionChangeListener) error {
	return s.listeners.add(listener, func() (*native.ChangeListener, error) {
		return attachChangeListener(listener, s.sub.AddChangeListener)
	}, s.detach)
}

// RemoveChangeListener removes a listener added with AddChangeListener; its SubscriptionClosed is called before
// RemoveChangeListener returns, and nothing after that (if AddChangeListener of that listener is still running, after
// AddChangeListener returns). It may be called from the listener. A listener added several times is removed once. A
// failure of the SDK is logged.
func (s *ObservableSubscription) RemoveChangeListener(listener osub.ObservableSubscriptionChangeListener) {
	s.listeners.remove(listener, s.detach)
}

func (s *ObservableSubscription) detach(cl *native.ChangeListener) error {
	return detachChangeListener(cl, s.sub.RemoveChangeListener)
}
