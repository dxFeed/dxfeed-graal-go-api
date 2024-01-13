package api

import (
	"dxfeed-graal-go-api/pkg/events"
	"dxfeed-graal-go-api/pkg/native"
)

type DXFeed struct {
	feed *native.DXFeed
}

func (f *DXFeed) CreateSubscription(eventType ...events.EventType) *DXFeedSubscription {
	data := make([]int32, len(eventType))
	for i := range data {
		data[i] = int32(eventType[i])
	}
	return &DXFeedSubscription{sub: f.feed.CreateSubscription(data...)}
}
