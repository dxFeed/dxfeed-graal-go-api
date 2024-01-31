package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"unsafe"
)

type QuoteMapper struct {
}

func (m QuoteMapper) CEvent(event interface{}) unsafe.Pointer {
	quoteEvent := event.(*quote.Quote)

	q := (*C.dxfg_quote_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_quote_t{}))))
	q.market_event.event_type.clazz = C.DXFG_EVENT_QUOTE
	q.market_event.event_symbol = C.CString(*quoteEvent.EventSymbol())
	q.market_event.event_time = C.int64_t(quoteEvent.EventTime())
	q.time_millis_sequence = C.int32_t(quoteEvent.GetTimeMillisSequence())
	q.time_nano_part = C.int32_t(quoteEvent.TimeNanoPart())
	q.bid_time = C.int64_t(quoteEvent.BidTime())
	q.bid_exchange_code = C.int16_t(quoteEvent.BidExchangeCode())
	q.bid_price = C.double(quoteEvent.BidPrice())
	q.bid_size = C.double(quoteEvent.BidSize())
	q.ask_time = C.int64_t(quoteEvent.AskTime())
	q.ask_exchange_code = C.int16_t(quoteEvent.AskExchangeCode())
	q.ask_price = C.double(quoteEvent.AskPrice())
	q.ask_size = C.double(quoteEvent.AskSize())
	return unsafe.Pointer(q)
}

func (m QuoteMapper) GoEvent(native unsafe.Pointer) interface{} {
	quoteNative := (*C.dxfg_quote_t)(native)
	q := quote.NewQuote(C.GoString(quoteNative.market_event.event_symbol))
	q.SetEventSymbol(C.GoString(quoteNative.market_event.event_symbol))
	q.SetEventTime(int64(quoteNative.market_event.event_time))
	q.SetTimeMillisSequence(int32(quoteNative.time_millis_sequence))
	q.SetTimeNanoPart(int32(quoteNative.time_nano_part))
	q.SetBidTime(int64(quoteNative.bid_time))
	q.SetBidExchangeCode(rune(quoteNative.bid_exchange_code))
	q.SetBidPrice(float64(quoteNative.bid_price))
	q.SetBidSize(float64(quoteNative.bid_size))
	q.SetAskTime(int64(quoteNative.ask_time))
	q.SetAskExchangeCode(rune(quoteNative.ask_exchange_code))
	q.SetAskPrice(float64(quoteNative.ask_price))
	q.SetAskSize(float64(quoteNative.ask_size))
	return q
}
