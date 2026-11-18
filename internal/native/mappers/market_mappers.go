package mappers

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/marketmaker"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/optionsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/orderimbalance"
)

// OptionSaleMapper maps the OptionSale event.
type OptionSaleMapper struct{}

// CEvent makes a native event from the Go event.
func (OptionSaleMapper) CEvent(event interface{}) unsafe.Pointer {
	e := event.(*optionsale.OptionSale)
	n := (*C.dxfg_v4_option_sale_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_option_sale_t{})))
	n.market_event.event_type.clazz = C.DXFG_EVENT_OPTION_SALE
	n.market_event.event_symbol = cString(*e.EventSymbol())
	n.market_event.event_time = C.int64_t(e.EventTime())
	n.event_flags = C.int32_t(e.EventFlags())
	n.index = C.int64_t(e.Index())
	n.time_sequence = C.int64_t(e.TimeSequence())
	n.time_nano_part = C.int32_t(e.TimeNanoPart())
	n.exchange_code = C.int16_t(e.ExchangeCode())
	n.price = C.double(e.Price())
	n.size = C.double(e.Size())
	n.bid_price = C.double(e.BidPrice())
	n.ask_price = C.double(e.AskPrice())
	n.exchange_sale_conditions = cStringOrNil(e.ExchangeSaleConditions())
	n.flags = C.int32_t(e.Flags())
	n.underlying_price = C.double(e.UnderlyingPrice())
	n.volatility = C.double(e.Volatility())
	n.delta = C.double(e.Delta())
	n.option_symbol = cStringOrNil(e.OptionSymbol())
	return unsafe.Pointer(n)
}

// GoEvent makes a Go event from the native event.
func (OptionSaleMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_option_sale_t)(native)
	e := optionsale.NewOptionSale(C.GoString(n.market_event.event_symbol))
	e.SetEventTime(int64(n.market_event.event_time))
	e.SetEventFlags(int32(n.event_flags))
	e.SetIndex(int64(n.index))
	e.SetTimeSequence(int64(n.time_sequence))
	e.SetTimeNanoPart(int32(n.time_nano_part))
	e.SetExchangeCode(int16(n.exchange_code))
	e.SetPrice(float64(n.price))
	e.SetSize(float64(n.size))
	e.SetBidPrice(float64(n.bid_price))
	e.SetAskPrice(float64(n.ask_price))
	e.SetExchangeSaleConditions(convertString(n.exchange_sale_conditions))
	e.SetFlags(int32(n.flags))
	e.SetUnderlyingPrice(float64(n.underlying_price))
	e.SetVolatility(float64(n.volatility))
	e.SetDelta(float64(n.delta))
	e.SetOptionSymbol(convertString(n.option_symbol))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (OptionSaleMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_option_sale_t)(event)
	cFree(unsafe.Pointer(n.market_event.event_symbol))
	cFree(unsafe.Pointer(n.exchange_sale_conditions))
	cFree(unsafe.Pointer(n.option_symbol))
	cFree(event)
}

// MarketMakerMapper maps the MarketMaker event.
type MarketMakerMapper struct{}

// CEvent makes a native event from the Go event.
func (MarketMakerMapper) CEvent(event interface{}) unsafe.Pointer {
	e := event.(*marketmaker.MarketMaker)
	n := (*C.dxfg_v4_market_maker_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_market_maker_t{})))
	n.market_event.event_type.clazz = C.DXFG_EVENT_MARKET_MAKER
	n.market_event.event_symbol = cString(*e.EventSymbol())
	n.market_event.event_time = C.int64_t(e.EventTime())
	n.event_flags = C.int32_t(e.EventFlags())
	n.index = C.int64_t(e.Index())
	n.bid_time = C.int64_t(e.BidTime())
	n.bid_price = C.double(e.BidPrice())
	n.bid_size = C.double(e.BidSize())
	n.bid_count = C.int64_t(e.BidCount())
	n.ask_time = C.int64_t(e.AskTime())
	n.ask_price = C.double(e.AskPrice())
	n.ask_size = C.double(e.AskSize())
	n.ask_count = C.int64_t(e.AskCount())
	return unsafe.Pointer(n)
}

// GoEvent makes a Go event from the native event.
func (MarketMakerMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_market_maker_t)(native)
	e := marketmaker.NewMarketMaker(C.GoString(n.market_event.event_symbol))
	e.SetEventTime(int64(n.market_event.event_time))
	e.SetEventFlags(int32(n.event_flags))
	_ = e.SetIndex(int64(n.index))
	e.SetBidTime(int64(n.bid_time))
	e.SetBidPrice(float64(n.bid_price))
	e.SetBidSize(float64(n.bid_size))
	e.SetBidCount(int64(n.bid_count))
	e.SetAskTime(int64(n.ask_time))
	e.SetAskPrice(float64(n.ask_price))
	e.SetAskSize(float64(n.ask_size))
	e.SetAskCount(int64(n.ask_count))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (MarketMakerMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_market_maker_t)(event)
	cFree(unsafe.Pointer(n.market_event.event_symbol))
	cFree(event)
}

// OtcMarketsOrderMapper maps the OtcMarketsOrder event.
type OtcMarketsOrderMapper struct{}

// CEvent makes a native event from the Go event.
func (OtcMarketsOrderMapper) CEvent(event interface{}) unsafe.Pointer {
	e := event.(*order.OtcMarketsOrder)
	n := (*C.dxfg_v4_otc_markets_order_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_otc_markets_order_t{})))
	orderToC(&n.order_base, &e.Order, C.DXFG_EVENT_OTC_MARKETS_ORDER)
	n.quote_access_payment = C.int32_t(e.QuoteAccessPayment())
	n.otc_markets_flags = C.int32_t(e.OtcMarketsFlags())
	return unsafe.Pointer(n)
}

// GoEvent makes a Go event from the native event.
func (OtcMarketsOrderMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_otc_markets_order_t)(native)
	e := order.NewOtcMarketsOrder(C.GoString(n.order_base.order_base.market_event.event_symbol))
	orderFromC(&e.Order, &n.order_base)
	e.SetQuoteAccessPayment(int32(n.quote_access_payment))
	e.SetOtcMarketsFlags(int32(n.otc_markets_flags))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (OtcMarketsOrderMapper) FreeCEvent(event unsafe.Pointer) {
	freeCOrder(&(*C.dxfg_v4_otc_markets_order_t)(event).order_base)
	cFree(event)
}

// OrderImbalanceMapper maps the OrderImbalance event.
type OrderImbalanceMapper struct{}

// CEvent makes a native event from the Go event.
func (OrderImbalanceMapper) CEvent(event interface{}) unsafe.Pointer {
	e := event.(*orderimbalance.OrderImbalance)
	n := (*C.dxfg_v4_order_imbalance_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_order_imbalance_t{})))
	n.market_event.event_type.clazz = C.DXFG_EVENT_ORDER_IMBALANCE
	n.market_event.event_symbol = cString(*e.EventSymbol())
	n.market_event.event_time = C.int64_t(e.EventTime())
	n.source_id = C.int32_t(e.SourceId())
	n.time_sequence = C.int64_t(e.TimeSequence())
	n.ref_price = C.double(e.RefPrice())
	n.paired_size = C.double(e.PairedSize())
	n.imbalance_size = C.double(e.ImbalanceSize())
	n.near_price = C.double(e.NearPrice())
	n.far_price = C.double(e.FarPrice())
	n.flags = C.int32_t(e.Flags())
	return unsafe.Pointer(n)
}

// GoEvent makes a Go event from the native event.
func (OrderImbalanceMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_order_imbalance_t)(native)
	e := orderimbalance.NewOrderImbalance(C.GoString(n.market_event.event_symbol))
	e.SetEventTime(int64(n.market_event.event_time))
	e.SetSourceId(int32(n.source_id))
	e.SetTimeSequence(int64(n.time_sequence))
	e.SetRefPrice(float64(n.ref_price))
	e.SetPairedSize(float64(n.paired_size))
	e.SetImbalanceSize(float64(n.imbalance_size))
	e.SetNearPrice(float64(n.near_price))
	e.SetFarPrice(float64(n.far_price))
	e.SetFlags(int32(n.flags))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (OrderImbalanceMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_order_imbalance_t)(event)
	cFree(unsafe.Pointer(n.market_event.event_symbol))
	cFree(event)
}
