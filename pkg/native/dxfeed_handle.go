package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DXFeedHandle struct {
	ptr *C.dxfg_feed_t
}

func (f DXFeedHandle) CreateSubscription(eventTypes ...int32) *DXFeedSubscription {
	sub := &DXFeedSubscription{}
	_ = executeInIsolateThread(func(thread *isolateThread) error {
		list := createEventClazzList(eventTypes...)
		sub.ptr = C.dxfg_DXFeed_createSubscription2(thread.ptr, f.ptr, (*C.dxfg_event_clazz_list_t)(unsafe.Pointer(list)))
		destroyEventClazzList(list)

		return nil
	})
	return sub
}
