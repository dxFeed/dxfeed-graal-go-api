package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DXFeed struct {
	ptr *C.dxfg_feed_t
}

func (f DXFeed) CreateSubscription(eventTypes ...int32) *DXFeedSubscription {
	thread := attachIsolateThread()
	defer thread.detachIsolateThread()
	sub := &DXFeedSubscription{}
	list := createEventClazzList(eventTypes...)
	sub.ptr = C.dxfg_DXFeed_createSubscription2(thread.ptr, f.ptr, (*C.dxfg_event_clazz_list_t)(unsafe.Pointer(list)))
	return sub
}
