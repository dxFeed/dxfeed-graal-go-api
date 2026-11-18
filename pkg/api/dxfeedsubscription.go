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
// *osub.CandleSymbol for the candles of a candle.CandleSymbol, *osub.WildcardSymbol for all symbols,
// *osub.TimeSeriesSubscriptionSymbol for time series events from a time, and *osub.IndexedEventSubscriptionSymbol for
// indexed events of one source.
type DXFeedSubscription struct {
	sub       *native.DXFeedSubscription
	mu        sync.Mutex
	listeners []*attachedListener
	closed    atomic.Bool
}

// attachedListener pairs a listener added with AddListener with its native listener. The entry is in the list before
// the native listener is attached, so a listener may remove itself from its first Update; native is nil until then.
type attachedListener struct {
	listener common.EventListener
	native   *native.EventListener
	removed  bool
}

// IsClosed reports whether Close succeeded.
func (s *DXFeedSubscription) IsClosed() bool {
	return s.closed.Load()
}

// AddListener adds a listener that receives the events of the subscription. The listener is called from an SDK thread;
// a panic in it is logged and does not end the process.
func (s *DXFeedSubscription) AddListener(listener common.EventListener) error {
	entry := &attachedListener{listener: listener}
	s.mu.Lock()
	if s.closed.Load() {
		// A closed subscription delivers no events.
		s.mu.Unlock()
		return nil
	}
	s.listeners = append(s.listeners, entry)
	s.mu.Unlock()

	nativeListener, err := s.sub.AttachListener(listener)
	s.mu.Lock()
	if err != nil {
		s.removeEntry(entry)
		s.mu.Unlock()
		return err
	}
	if !entry.removed {
		entry.native = nativeListener
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()
	// RemoveListener or Close came while the native listener was being attached.
	if err := s.sub.DetachListener(nativeListener); err != nil {
		log.Printf("dxfeed: remove event listener: %v", err)
	}
	return nil
}

// RemoveListener removes a listener added with AddListener; the listener is not called after RemoveListener returns
// (if AddListener of that listener is still running, after AddListener returns). It may be called from the listener.
// A listener added several times is removed once. A failure of the SDK is logged.
func (s *DXFeedSubscription) RemoveListener(listener common.EventListener) {
	s.mu.Lock()
	var removed *native.EventListener
	for _, l := range s.listeners {
		if l.listener == listener {
			removed = l.native
			s.removeEntry(l)
			break
		}
	}
	s.mu.Unlock()
	if removed == nil {
		// Not added, or AddListener is still attaching it and detaches it when done.
		return
	}
	if err := s.sub.DetachListener(removed); err != nil {
		log.Printf("dxfeed: remove event listener: %v", err)
	}
}

// removeEntry marks the entry removed and takes it out of the list; s.mu must be held.
func (s *DXFeedSubscription) removeEntry(entry *attachedListener) {
	entry.removed = true
	for i, l := range s.listeners {
		if l == entry {
			s.listeners = append(s.listeners[:i:i], s.listeners[i+1:]...)
			return
		}
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
	for _, l := range listeners {
		l.removed = true
	}
	s.mu.Unlock()
	var err error
	for _, l := range listeners {
		if l.native != nil {
			err = errors.Join(err, l.native.Release())
		}
	}
	return err
}
