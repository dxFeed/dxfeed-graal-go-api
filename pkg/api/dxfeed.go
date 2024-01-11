package api

type DXFeed struct {
}

func (f *DXFeed) CreateSubscription() DXFeedSubscription {
	return DXFeedSubscription{}
}
