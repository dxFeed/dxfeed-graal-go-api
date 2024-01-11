package api

type DXFeed struct {
}

func (f *DXFeed) CreateSubscription(events ...int32) *DXFeedSubscription {
	return &DXFeedSubscription{}
}
