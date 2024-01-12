package api

import "dxfeed-graal-go-api/pkg/events"

type DXFeed struct {
}

func (f *DXFeed) CreateSubscription(eventType events.EventType) *DXFeedSubscription {
	return &DXFeedSubscription{}
}
