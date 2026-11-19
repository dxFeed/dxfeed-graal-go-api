package api

import (
	"errors"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

var (
	errNilFeed         = errors.New("dxfeed: the feed is nil")
	errNilSubscription = errors.New("dxfeed: the subscription is nil")
)

// DXFeedTimeSeriesSubscription is a subscription to time series events, such as candle.Candle or
// timeandsale.TimeAndSale, from a time; create it with DXFeed.CreateTimeSeriesSubscription. Its plain symbols are
// subscribed from the time set with SetFromTime; until then it receives no events.
type DXFeedTimeSeriesSubscription struct {
	*DXFeedSubscription
	ts *native.DXFeedTimeSeriesSubscription
}

// SetFromTime sets the time in milliseconds since the epoch from which the subscription receives the events.
func (s *DXFeedTimeSeriesSubscription) SetFromTime(fromTime int64) error {
	return s.ts.SetFromTime(fromTime)
}

// FromTime returns the time set with SetFromTime, math.MaxInt64 for a new subscription.
func (s *DXFeedTimeSeriesSubscription) FromTime() (int64, error) {
	return s.ts.FromTime()
}

// NewDXFeedSubscription creates a subscription to events of the given types that is not attached to any feed: it
// receives no events until DXFeed.AttachSubscription or DXFeedSubscription.Attach attaches it to a feed.
func NewDXFeedSubscription(eventType ...eventcodes.EventCode) (*DXFeedSubscription, error) {
	sub, err := native.NewDetachedSubscription(nativeCodes(eventType)...)
	if err != nil {
		return nil, err
	}
	return newDXFeedSubscription(sub), nil
}

// CreateTimeSeriesSubscription creates a time series subscription to events of the given types and attaches it to the
// feed. The types are expected to be time series event types, such as eventcodes.Candle or eventcodes.TimeAndSale;
// the SDK does not check them.
func (f *DXFeed) CreateTimeSeriesSubscription(eventType ...eventcodes.EventCode) (*DXFeedTimeSeriesSubscription,
	error) {
	ts, err := f.feed.CreateTimeSeriesSubscription(nativeCodes(eventType)...)
	if err != nil {
		return nil, err
	}
	return &DXFeedTimeSeriesSubscription{DXFeedSubscription: newDXFeedSubscription(&ts.DXFeedSubscription), ts: ts},
		nil
}

// AttachSubscription attaches the subscription to the feed: the feed delivers its events to the subscription. It does
// nothing if the subscription is already attached to the feed. A subscription can be attached to several feeds.
func (f *DXFeed) AttachSubscription(sub *DXFeedSubscription) error {
	if sub == nil {
		return errNilSubscription
	}
	return f.feed.AttachSubscription(sub.sub)
}

// DetachSubscription detaches the subscription from the feed; it does nothing if the subscription is not attached to
// the feed. The subscription keeps its symbols and listeners.
func (f *DXFeed) DetachSubscription(sub *DXFeedSubscription) error {
	if sub == nil {
		return errNilSubscription
	}
	return f.feed.DetachSubscription(sub.sub)
}

// DetachSubscriptionAndClear detaches the subscription from the feed and clears the data delivered to it by
// publishing empty events to it. It does nothing if the subscription is not attached to the feed.
func (f *DXFeed) DetachSubscriptionAndClear(sub *DXFeedSubscription) error {
	if sub == nil {
		return errNilSubscription
	}
	return f.feed.DetachSubscriptionAndClear(sub.sub)
}

// Attach attaches the subscription to the feed; see DXFeed.AttachSubscription.
func (s *DXFeedSubscription) Attach(feed *DXFeed) error {
	if feed == nil {
		return errNilFeed
	}
	return s.sub.Attach(feed.feed)
}

// Detach detaches the subscription from the feed; see DXFeed.DetachSubscription.
func (s *DXFeedSubscription) Detach(feed *DXFeed) error {
	if feed == nil {
		return errNilFeed
	}
	return s.sub.Detach(feed.feed)
}

func nativeCodes(eventType []eventcodes.EventCode) []int32 {
	codes := make([]int32, len(eventType))
	for i := range codes {
		codes[i] = eventType[i].NativeCode()
	}
	return codes
}
