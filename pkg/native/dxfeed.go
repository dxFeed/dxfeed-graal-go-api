package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

type DXFeed struct {
	ptr *C.dxfg_feed_t
}

func (f DXFeed) CreateSubscription(eventTypes ...int32) *DXFeedSubscription {
	return &DXFeedSubscription{ptr: nil}
}
