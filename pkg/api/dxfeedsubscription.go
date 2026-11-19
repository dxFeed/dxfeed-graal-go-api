package api

import (
	"errors"
	"sync/atomic"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// DXFeedSubscription is a subscription to events of some types for a set of symbols; create it with
// DXFeed.CreateSubscription.
//
// A symbol is an osub.Symbol: osub.String for a symbol such as "AAPL" (osub.Strings converts several strings at once),
// *osub.CandleSymbol for the candles of a candle.CandleSymbol, *osub.WildcardSymbol for all symbols,
// *osub.TimeSeriesSubscriptionSymbol for time series events from a time, and *osub.IndexedEventSubscriptionSymbol for
// indexed events of one source.
type DXFeedSubscription struct {
	sub             *native.DXFeedSubscription
	listeners       listenerRegistry[common.EventListener, *native.EventListener]
	changeListeners changeListenerRegistry
	closed          atomic.Bool
}

func newDXFeedSubscription(sub *native.DXFeedSubscription) *DXFeedSubscription {
	return &DXFeedSubscription{
		sub:             sub,
		listeners:       listenerRegistry[common.EventListener, *native.EventListener]{kind: "event listener"},
		changeListeners: newChangeListenerRegistry(),
	}
}

// IsClosed reports whether Close succeeded.
func (s *DXFeedSubscription) IsClosed() bool {
	return s.closed.Load()
}

// AddListener adds a listener that receives the events of the subscription. The listener is called from an SDK thread;
// a panic in it is logged and does not end the process.
func (s *DXFeedSubscription) AddListener(listener common.EventListener) error {
	return s.listeners.add(listener, func() (*native.EventListener, error) {
		return s.sub.AttachListener(listener)
	}, s.sub.DetachListener)
}

// RemoveListener removes a listener added with AddListener; the listener is not called after RemoveListener returns
// (if AddListener of that listener is still running, after AddListener returns). It may be called from the listener.
// A listener added several times is removed once. A failure of the SDK is logged.
func (s *DXFeedSubscription) RemoveListener(listener common.EventListener) {
	s.listeners.remove(listener, s.sub.DetachListener)
}

// AddChangeListener adds a listener of the changes of the symbols of the subscription; see
// osub.ObservableSubscriptionChangeListener. If the subscription has symbols, the listener gets them in SymbolsAdded
// before AddChangeListener returns. Adding a listener to a closed subscription does nothing.
func (s *DXFeedSubscription) AddChangeListener(listener osub.ObservableSubscriptionChangeListener) error {
	return s.changeListeners.add(listener, func() (*native.ChangeListener, error) {
		return attachChangeListener(listener, s.sub.AddChangeListener)
	}, func(cl *native.ChangeListener) error {
		return detachChangeListener(cl, s.sub.RemoveChangeListener)
	})
}

// RemoveChangeListener removes a listener added with AddChangeListener; its SubscriptionClosed is called before
// RemoveChangeListener returns, and nothing after that (if AddChangeListener of that listener is still running, after
// AddChangeListener returns). It may be called from the listener. A listener added several times is removed once. A
// failure of the SDK is logged.
func (s *DXFeedSubscription) RemoveChangeListener(listener osub.ObservableSubscriptionChangeListener) {
	s.changeListeners.remove(listener, func(cl *native.ChangeListener) error {
		return detachChangeListener(cl, s.sub.RemoveChangeListener)
	})
}

// attachChangeListener creates a native change listener and adds it with add.
func attachChangeListener(listener osub.ObservableSubscriptionChangeListener,
	add func(*native.ChangeListener) error) (*native.ChangeListener, error) {
	cl, err := native.NewChangeListener(listener)
	if err != nil {
		return nil, err
	}
	if err := add(cl); err != nil {
		return nil, errors.Join(err, cl.Release())
	}
	return cl, nil
}

// detachChangeListener removes the native change listener with remove and releases it.
func detachChangeListener(cl *native.ChangeListener, remove func(*native.ChangeListener) error) error {
	return errors.Join(remove(cl), cl.Release())
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

// Close closes the subscription; it stops receiving events, and its change listeners get SubscriptionClosed.
func (s *DXFeedSubscription) Close() error {
	if err := s.sub.Close(); err != nil {
		return err
	}
	s.closed.Store(true)
	var err error
	for _, l := range s.listeners.close() {
		err = errors.Join(err, l.Release())
	}
	for _, l := range s.changeListeners.close() {
		err = errors.Join(err, l.Release())
	}
	return err
}
