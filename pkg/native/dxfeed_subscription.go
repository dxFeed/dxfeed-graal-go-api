package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
extern void OnEventReceived(graal_isolatethread_t *thread, dxfg_event_type_list *events, void *user_data);
*/
import "C"
import (
	"dxfeed-graal-go-api/pkg/events/market"
	"unsafe"
)
import gopointer "github.com/mattn/go-pointer"

type DXFeedSubscription struct {
	ptr *C.dxfg_subscription_t
}

type dxfg_symbol_t struct {
	t      C.int32_t
	symbol *C.char
}

//export OnEventReceived
func OnEventReceived(thread *C.graal_isolatethread_t, events *C.dxfg_event_type_list, userData unsafe.Pointer) {
	size := int(events.size)
	list := make([]interface{}, size)
	for i := 0; i < size; i++ {
		eventPtr := *(**C.dxfg_event_type_t)(unsafe.Pointer(uintptr(unsafe.Pointer(events.elements)) + uintptr(i)*unsafe.Sizeof(*events.elements)))
		quote := (*C.dxfg_quote_t)(unsafe.Pointer(eventPtr))

		q := market.NewQuote(C.GoString(quote.market_event.event_symbol))
		q.SetEventTime(int64(quote.market_event.event_time))
		q.SetTimeMillisSequence(int32(quote.time_millis_sequence))
		q.SetTimeNanoPart(int32(quote.time_nano_part))
		q.SetBidTime(int64(quote.bid_time))
		q.SetBidExchangeCode(rune(quote.bid_exchange_code))
		q.SetBidPrice(float64(quote.bid_price))
		q.SetBidSize(float64(quote.bid_size))
		q.SetAskTime(int64(quote.ask_time))
		q.SetAskExchangeCode(rune(quote.ask_exchange_code))
		q.SetAskPrice(float64(quote.ask_price))
		q.SetAskSize(float64(quote.ask_size))
		list[i] = q
	}
	gopointer.Restore(userData).(EventListener).Update(list)
}

func (s DXFeedSubscription) AttachListener(listener EventListener) {
	executeInIsolateThread(func(thread *isolateThread) {
		l := C.dxfg_DXFeedEventListener_new(thread.ptr, (*[0]byte)(C.OnEventReceived), gopointer.Save(listener))
		C.dxfg_DXFeedSubscription_addEventListener(thread.ptr, s.ptr, l)
	})
}

func (s DXFeedSubscription) AddSymbol(symbol string) {
	executeInIsolateThread(func(thread *isolateThread) {
		ss := &dxfg_symbol_t{}
		ss.t = 0
		ss.symbol = C.CString(symbol)
		C.dxfg_DXFeedSubscription_addSymbol(thread.ptr, s.ptr, (*C.dxfg_symbol_t)(unsafe.Pointer(ss)))
	})
}
