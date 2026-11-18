package api

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// DXFeedSubscription is a subscription to events of some types for a set of symbols; create it with
// DXFeed.CreateSubscription.
//
// A symbol is an osub.Symbol: osub.String for a symbol such as "AAPL" (osub.Strings converts several strings at once),
// *osub.WildcardSymbol for all symbols, *osub.TimeSeriesSubscriptionSymbol for time series events from a time, and
// *osub.IndexedEventSubscriptionSymbol for indexed events of one source.
type DXFeedSubscription struct {
	sub       *native.DXFeedSubscription
	mu        sync.Mutex
	listeners []attachedListener
	closed    atomic.Bool
}

// attachedListener pairs a listener added with AddListener with its native listener.
type attachedListener struct {
	listener common.EventListener
	native   *native.EventListener
}

// IsClosed reports whether Close succeeded.
func (s *DXFeedSubscription) IsClosed() bool {
	return s.closed.Load()
}

// AddListener adds a listener that receives the events of the subscription. The listener is called from an SDK thread;
// a panic in it is logged and does not end the process.
func (s *DXFeedSubscription) AddListener(listener common.EventListener) error {
	nativeListener, err := s.sub.AttachListener(listener)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.listeners = append(s.listeners, attachedListener{listener: listener, native: nativeListener})
	s.mu.Unlock()
	return nil
}

// RemoveListener removes a listener added with AddListener; the listener is not called after RemoveListener returns.
// A listener added several times is removed once. A failure of the SDK is logged.
func (s *DXFeedSubscription) RemoveListener(listener common.EventListener) {
	s.mu.Lock()
	var removed *native.EventListener
	for i, l := range s.listeners {
		if l.listener == listener {
			removed = l.native
			s.listeners = append(s.listeners[:i:i], s.listeners[i+1:]...)
			break
		}
	}
	s.mu.Unlock()
	if removed == nil {
		return
	}
	if err := s.sub.DetachListener(removed); err != nil {
		log.Printf("dxfeed: remove event listener: %v", err)
	}
}

// AddSymbol adds a symbol to the subscription; see DXFeedSubscription for the symbol types.
func (s *DXFeedSubscription) AddSymbol(symbol osub.Symbol) error {
	return s.sub.AddSymbol(symbol)
}

// AddSymbols adds symbols to the subscription; see DXFeedSubscription for the symbol types.
func (s *DXFeedSubscription) AddSymbols(symbols ...osub.Symbol) error {
	return s.sub.AddSymbols(symbols...)
}

// RemoveSymbol removes a symbol from the subscription.
func (s *DXFeedSubscription) RemoveSymbol(symbol osub.Symbol) error {
	return s.sub.RemoveSymbol(symbol)
}

// RemoveSymbols removes symbols from the subscription.
func (s *DXFeedSubscription) RemoveSymbols(symbols ...osub.Symbol) error {
	return s.sub.RemoveSymbols(symbols...)
}

// Clear removes all symbols from the subscription.
func (s *DXFeedSubscription) Clear() error {
	return s.sub.Clear()
}

// Close closes the subscription; it stops receiving events.
func (s *DXFeedSubscription) Close() error {
	if err := s.sub.Close(); err != nil {
		return err
	}
	s.closed.Store(true)
	s.mu.Lock()
	listeners := s.listeners
	s.listeners = nil
	s.mu.Unlock()
	var err error
	for _, l := range listeners {
		err = errors.Join(err, l.native.Release())
	}
	return err
}
