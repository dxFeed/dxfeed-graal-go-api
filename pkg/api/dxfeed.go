package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
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
