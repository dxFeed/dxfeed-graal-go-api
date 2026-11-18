package api

import (
	"errors"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// DXFeed is the feed of an endpoint; get it with DXEndpoint.GetFeed.
type DXFeed struct {
	feed *native.DXFeedHandle
}

// CreateSubscription creates a subscription to events of the given types. Add symbols and a listener to it, and close
// it when it is no longer needed.
func (f *DXFeed) CreateSubscription(eventType ...eventcodes.EventCode) (*DXFeedSubscription, error) {
	data := make([]int32, len(eventType))
	for i := range data {
		data[i] = eventType[i].NativeCode()
	}
	sub, err := f.feed.CreateSubscription(data...)
	if err != nil {
		return nil, err
	}
	return newDXFeedSubscription(sub), nil
}

// GetLastEvent returns the last event for the type and symbol of the given lasting event, such as a *quote.Quote with
// its event symbol set. When the last event is not available (there is no subscription to the type and symbol, or the
// event has not arrived yet), the result is a copy of the given event. The call reads the local cache of the feed and
// makes no remote calls; the cache has the event only while an attached subscription is subscribed to it.
func (f *DXFeed) GetLastEvent(event interface{}) (interface{}, error) {
	return f.feed.GetLastEvent(event)
}

// GetLastEvents is the bulk version of GetLastEvent: it returns the last event for every given event, in the same
// order.
func (f *DXFeed) GetLastEvents(events []interface{}) ([]interface{}, error) {
	return f.feed.GetLastEvents(events)
}

// GetLastEventIfSubscribed returns the last event of the lasting type for the symbol, or nil if there is no
// subscription to them or the event has not arrived yet. It reads the local cache of the feed.
func (f *DXFeed) GetLastEventIfSubscribed(eventType eventcodes.EventCode, symbol osub.Symbol) (interface{}, error) {
	return f.feed.GetLastEventIfSubscribed(eventType.NativeCode(), symbol)
}

// GetIndexedEventsIfSubscribed returns the indexed events of the type for the symbol and source, ordered by index, if
// there is a subscription to them. The list is empty if there is no subscription or the events have not arrived yet;
// it is not a consistent snapshot. It reads the local cache of the feed.
func (f *DXFeed) GetIndexedEventsIfSubscribed(eventType eventcodes.EventCode, symbol osub.Symbol,
	source events.IndexedEventSourceInterface) ([]interface{}, error) {
	if source == nil || source.Name() == nil {
		return nil, errors.New("cannot get the indexed events: no source")
	}
	return f.feed.GetIndexedEventsIfSubscribed(eventType.NativeCode(), symbol, *source.Name())
}

// GetTimeSeriesIfSubscribed returns the time series events of the type for the symbol from fromTime to toTime
// (inclusive, in milliseconds since the epoch; math.MaxInt64 is no upper limit), ordered by time, if there is a
// subscription to them. The list is empty if there is no subscription or the events have not arrived yet; it is not a
// consistent snapshot. The symbol must not be an *osub.TimeSeriesSubscriptionSymbol. It reads the local cache of the
// feed.
func (f *DXFeed) GetTimeSeriesIfSubscribed(eventType eventcodes.EventCode, symbol osub.Symbol, fromTime,
	toTime int64) ([]interface{}, error) {
	return f.feed.GetTimeSeriesIfSubscribed(eventType.NativeCode(), symbol, fromTime, toTime)
}
