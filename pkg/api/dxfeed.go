package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

type DXFeed struct {
	feed *native.DXFeedHandle
}

func (f *DXFeed) CreateSubscription(eventType ...eventcodes.EventCode) (*DXFeedSubscription, error) {
	data := make([]int32, len(eventType))
	for i := range data {
		data[i] = eventType[i].NativeCode()
	}
	sub, err := f.feed.CreateSubscription(data...)
	if err != nil {
		return nil, err
	}
	return &DXFeedSubscription{sub: sub}, nil
}
