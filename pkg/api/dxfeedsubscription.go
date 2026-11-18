package api

import (
	"sync/atomic"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// DXFeedSubscription is a subscription to events of some types for a set of symbols; create it with
// DXFeed.CreateSubscription.
//
// A symbol is a string such as "AAPL", or one of the osub types: *osub.WildcardSymbol for all symbols,
// *osub.TimeSeriesSubscriptionSymbol for time series events from a time, *osub.IndexedEventSubscriptionSymbol for
// indexed events of one source.
type DXFeedSubscription struct {
	sub               *native.DXFeedSubscription
	eventListenerList []common.EventListener
	closed            atomic.Bool
}

// IsClosed reports whether Close succeeded.
func (s *DXFeedSubscription) IsClosed() bool {
	return s.closed.Load()
}

// AddListener adds a listener that receives the events of the subscription. The listener is called from an SDK thread;
// a panic in it is logged and does not end the process.
func (s *DXFeedSubscription) AddListener(listener common.EventListener) error {
	s.eventListenerList = append(s.eventListenerList, listener)
	return s.sub.AttachListener(listener)
}

// RemoveListener is meant to remove a listener added with AddListener. Known limitation: the listener still receives
// events; close the subscription to stop them.
func (s *DXFeedSubscription) RemoveListener(listener common.EventListener) {
	s.eventListenerList = removeFromSlice(s.eventListenerList, listener)
}

// AddSymbol adds a symbol to the subscription; see DXFeedSubscription for the symbol types.
func (s *DXFeedSubscription) AddSymbol(symbol any) error {
	return s.sub.AddSymbol(symbol)
}

// AddSymbols adds symbols to the subscription; see DXFeedSubscription for the symbol types.
func (s *DXFeedSubscription) AddSymbols(symbols ...any) error {
	return s.sub.AddSymbols(symbols...)
}

// RemoveSymbol removes a symbol from the subscription.
func (s *DXFeedSubscription) RemoveSymbol(symbol any) error {
	return s.sub.RemoveSymbol(symbol)
}

// RemoveSymbols removes symbols from the subscription.
func (s *DXFeedSubscription) RemoveSymbols(symbols ...any) error {
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
	return nil
}

func removeFromSlice(list []common.EventListener, observerToRemove common.EventListener) []common.EventListener {
	listLength := len(list)
	for i, observer := range list {
		if observerToRemove == observer {
			list[listLength-1], list[i] = list[i], list[listLength-1]
			return list[:listLength-1]
		}
	}
	return list
}
