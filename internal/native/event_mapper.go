package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"unsafe"
)

type eventMapper struct {
}

func newEventMapper() *eventMapper {
	return &eventMapper{}
}

func (m *eventMapper) goEvents(eventsList *C.dxfg_event_type_list) []interface{} {
	if eventsList == nil || eventsList.elements == nil || int(eventsList.size) == 0 {
		return nil
	}

	size := int(eventsList.size)
	list := make([]interface{}, size)
	elementsSlice := unsafe.Slice(eventsList.elements, C.size_t(eventsList.size))

	for i, event := range elementsSlice {
		list[i] = m.goEvent(event)
	}

	return list
}

func (m *eventMapper) goEvent(event *C.dxfg_event_type_t) interface{} {
	switch event.clazz {
	case C.DXFG_EVENT_QUOTE:
		return m.goQuote((*C.dxfg_quote_t)(unsafe.Pointer(event)))
	case C.DXFG_EVENT_TIME_AND_SALE:
		return m.goTimeAndSale((*C.dxfg_time_and_sale_t)(unsafe.Pointer(event)))
	case C.DXFG_EVENT_PROFILE:
		return m.goProfile((*C.dxfg_profile_t)(unsafe.Pointer(event)))
	default:
		panic("unknown event eventcodes")
	}
}

func (m *eventMapper) goQuote(quoteNative *C.dxfg_quote_t) *quote.Quote {
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

func (m *eventMapper) cQuote(quote *quote.Quote) *C.dxfg_quote_t {
	q := (*C.dxfg_quote_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_quote_t{}))))
	q.market_event.event_type.clazz = C.DXFG_EVENT_QUOTE
	q.market_event.event_symbol = C.CString(*quote.EventSymbol())
	q.market_event.event_time = C.int64_t(quote.EventTime())
	q.time_millis_sequence = C.int32_t(quote.GetTimeMillisSequence())
	q.time_nano_part = C.int32_t(quote.TimeNanoPart())
	q.bid_time = C.int64_t(quote.BidTime())
	q.bid_exchange_code = C.int16_t(quote.BidExchangeCode())
	q.bid_price = C.double(quote.BidPrice())
	q.bid_size = C.double(quote.BidSize())
	q.ask_time = C.int64_t(quote.AskTime())
	q.ask_exchange_code = C.int16_t(quote.AskExchangeCode())
	q.ask_price = C.double(quote.AskPrice())
	q.ask_size = C.double(quote.AskSize())
	return q
}

func (m *eventMapper) goTimeAndSale(timeAndSale *C.dxfg_time_and_sale_t) *timeandsale.TimeAndSale {
	t := timeandsale.NewTimeAndSale(C.GoString(timeAndSale.market_event.event_symbol))
	t.SetEventTime(int64(timeAndSale.market_event.event_time))
	t.SetTimeNanoPart(int32(timeAndSale.time_nano_part))
	t.SetExchangeCode(int16(timeAndSale.exchange_code))
	t.SetPrice(float64(timeAndSale.price))
	t.SetSize(float64(timeAndSale.size))
	t.SetBidPrice(float64(timeAndSale.bid_price))
	t.SetAskPrice(float64(timeAndSale.ask_price))
	t.SetExchangeSaleConditions(ConvertString(timeAndSale.exchange_sale_conditions))
	t.SetBuyer(ConvertString(timeAndSale.buyer))
	t.SetSeller(ConvertString(timeAndSale.seller))
	t.SetEventFlags(int32(timeAndSale.event_flags))
	t.SetIndex(int64(timeAndSale.index))
	t.SetFlags(int32(timeAndSale.flags))
	return t
}

func (m *eventMapper) cTimeAndSale(timeAndSale *timeandsale.TimeAndSale) *C.dxfg_time_and_sale_t {
	t := (*C.dxfg_time_and_sale_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_time_and_sale_t{}))))
	t.market_event.event_type.clazz = C.DXFG_EVENT_TIME_AND_SALE
	t.market_event.event_symbol = C.CString(*timeAndSale.EventSymbol())
	t.market_event.event_time = C.int64_t(timeAndSale.EventTime())
	t.time_nano_part = C.int32_t(timeAndSale.TimeNanoPart())
	t.exchange_code = C.int16_t(timeAndSale.ExchangeCode())
	t.price = C.double(timeAndSale.Price())
	t.size = C.double(timeAndSale.Size())
	t.bid_price = C.double(timeAndSale.BidPrice())
	t.ask_price = C.double(timeAndSale.AskPrice())
	t.exchange_sale_conditions = CString(timeAndSale.ExchangeSaleConditions())
	t.buyer = CString(timeAndSale.Buyer())
	t.seller = CString(timeAndSale.Seller())
	t.event_flags = C.int32_t(timeAndSale.EventFlags())
	t.index = C.int64_t(timeAndSale.Index())
	t.flags = C.int32_t(timeAndSale.Flags())
	return t
}

func (m *eventMapper) goProfile(native *C.dxfg_profile_t) *profile.Profile {
	p := profile.NewProfile(C.GoString(native.market_event.event_symbol))
	p.SetEventTime(int64(native.market_event.event_time))
	p.SetDescription(ConvertString(native.description))
	p.SetStatusReason(ConvertString(native.status_reason))
	p.SetHaltStartTime(int64(native.halt_start_time))
	p.SetHaltEndTime(int64(native.halt_end_time))
	p.SetHighLimitPrice(float64(native.high_limit_price))
	p.SetLowLimitPrice(float64(native.low_limit_price))
	p.SetHigh52WeekPrice(float64(native.high_52_week_price))
	p.SetLow52WeekPrice(float64(native.low_52_week_price))
	p.SetBeta(float64(native.beta))
	p.SetEarningsPerShare(float64(native.earnings_per_share))
	p.SetDividendFrequency(float64(native.dividend_frequency))
	p.SetExDividendAmount(float64(native.ex_dividend_amount))
	p.SetExDividendDayId(int32(native.ex_dividend_day_id))
	p.SetShares(float64(native.shares))
	p.SetFreeFloat(float64(native.free_float))
	p.SetFlags(int32(native.flags))
	return p
}

func (m *eventMapper) cProfile(profile *profile.Profile) *C.dxfg_profile_t {
	p := (*C.dxfg_profile_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_profile_t{}))))
	p.market_event.event_type.clazz = C.DXFG_EVENT_PROFILE
	p.market_event.event_symbol = C.CString(*profile.EventSymbol())
	p.market_event.event_time = C.int64_t(profile.EventTime())
	if profile.Description() != nil {
		p.description = C.CString(*profile.Description())
	}
	if profile.StatusReason() != nil {
		p.status_reason = C.CString(*profile.StatusReason())
	}
	p.halt_end_time = C.int64_t(profile.HaltEndTime())
	p.halt_start_time = C.int64_t(profile.HaltStartTime())
	p.halt_end_time = C.int64_t(profile.HaltEndTime())
	p.high_limit_price = C.double(profile.HighLimitPrice())
	p.low_limit_price = C.double(profile.LowLimitPrice())
	p.high_52_week_price = C.double(profile.High52WeekPrice())
	p.low_52_week_price = C.double(profile.Low52WeekPrice())
	p.beta = C.double(profile.Beta())
	p.earnings_per_share = C.double(profile.EarningsPerShare())
	p.dividend_frequency = C.double(profile.DividendFrequency())
	p.ex_dividend_amount = C.double(profile.ExDividendAmount())
	p.ex_dividend_day_id = C.int32_t(profile.ExDividendDayId())
	p.shares = C.double(profile.Shares())
	p.free_float = C.double(profile.FreeFloat())
	p.flags = C.int32_t(profile.Flags())
	return p
}

func CString(str *string) *C.char {
	if str == nil {
		return nil
	}
	return C.CString(*str)
}

func (m *eventMapper) cStringSymbol(str string) *dxfg_symbol_t {
	ss := &dxfg_symbol_t{}
	ss.t = 0
	ss.symbol = C.CString(str)
	return ss
}

func (m *eventMapper) cWildCardSymbol() *dxfg_symbol_t {
	ss := &dxfg_symbol_t{}
	ss.t = 2
	return ss
}
