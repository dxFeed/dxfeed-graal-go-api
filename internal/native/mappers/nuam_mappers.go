// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
)

// NuamOrderMapper maps the NuamOrder event.
type NuamOrderMapper struct{}

// CEvent makes a native event from the Go event.
func (NuamOrderMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*order.NuamOrder)
	n := (*C.dxfg_v4_nuam_order_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_nuam_order_t{})))
	orderToC(a, &n.order_base, &e.Order, C.DXFG_EVENT_NUAM_ORDER)
	n.actor_id = C.int32_t(e.ActorId())
	n.participant_id = C.int32_t(e.ParticipantId())
	n.submitter_id = C.int32_t(e.SubmitterId())
	n.on_behalf_of_submitter_id = C.int32_t(e.OnBehalfOfSubmitterId())
	n.client_order_id = a.cStringOrNil(e.ClientOrderId())
	n.customer_account = a.cStringOrNil(e.CustomerAccount())
	n.customer_info = a.cStringOrNil(e.CustomerInfo())
	n.exchange_info = a.cStringOrNil(e.ExchangeInfo())
	n.time_in_force_data = C.int32_t(e.TimeInForceData())
	n.trigger_order_book_id = C.int32_t(e.TriggerOrderBookId())
	n.trigger_price = C.double(e.TriggerPrice())
	n.trigger_session_type = C.int32_t(e.TriggerSessionType())
	n.order_quantity = C.double(e.OrderQuantity())
	n.display_quantity = C.double(e.DisplayQuantity())
	n.refresh_quantity = C.double(e.RefreshQuantity())
	n.leaves_quantity = C.double(e.LeavesQuantity())
	n.matched_quantity = C.double(e.MatchedQuantity())
	n.nuam_flags = C.int32_t(e.NuamFlags())
	return unsafe.Pointer(n)
}

// GoEvent makes a Go event from the native event.
func (NuamOrderMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_nuam_order_t)(native)
	e := order.NewNuamOrder(C.GoString(n.order_base.order_base.market_event.event_symbol))
	orderFromC(&e.Order, &n.order_base)
	e.SetActorId(int32(n.actor_id))
	e.SetParticipantId(int32(n.participant_id))
	e.SetSubmitterId(int32(n.submitter_id))
	e.SetOnBehalfOfSubmitterId(int32(n.on_behalf_of_submitter_id))
	e.SetClientOrderId(convertString(n.client_order_id))
	e.SetCustomerAccount(convertString(n.customer_account))
	e.SetCustomerInfo(convertString(n.customer_info))
	e.SetExchangeInfo(convertString(n.exchange_info))
	e.SetTimeInForceData(int32(n.time_in_force_data))
	e.SetTriggerOrderBookId(int32(n.trigger_order_book_id))
	e.SetTriggerPrice(float64(n.trigger_price))
	e.SetTriggerSessionType(int32(n.trigger_session_type))
	e.SetOrderQuantity(float64(n.order_quantity))
	e.SetDisplayQuantity(float64(n.display_quantity))
	e.SetRefreshQuantity(float64(n.refresh_quantity))
	e.SetLeavesQuantity(float64(n.leaves_quantity))
	e.SetMatchedQuantity(float64(n.matched_quantity))
	e.SetNuamFlags(int32(n.nuam_flags))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (NuamOrderMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_nuam_order_t)(event)
	freeCOrder(&n.order_base)
	cFree(unsafe.Pointer(n.client_order_id))
	cFree(unsafe.Pointer(n.customer_account))
	cFree(unsafe.Pointer(n.customer_info))
	cFree(unsafe.Pointer(n.exchange_info))
	cFree(event)
}

// NuamTimeAndSaleMapper maps the NuamTimeAndSale event.
type NuamTimeAndSaleMapper struct{}

// CEvent makes a native event from the Go event.
func (NuamTimeAndSaleMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*timeandsale.NuamTimeAndSale)
	n := (*C.dxfg_v4_nuam_time_and_sale_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_nuam_time_and_sale_t{})))
	timeAndSaleToC(a, &n.time_and_sale_base, &e.TimeAndSale, C.DXFG_EVENT_NUAM_TIME_AND_SALE)
	n.match_id = C.int64_t(e.MatchId())
	return unsafe.Pointer(n)
}

// GoEvent makes a Go event from the native event.
func (NuamTimeAndSaleMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_nuam_time_and_sale_t)(native)
	e := timeandsale.NewNuamTimeAndSale(C.GoString(n.time_and_sale_base.market_event.event_symbol))
	timeAndSaleFromC(&e.TimeAndSale, &n.time_and_sale_base)
	e.SetMatchId(int64(n.match_id))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (NuamTimeAndSaleMapper) FreeCEvent(event unsafe.Pointer) {
	freeCTimeAndSale(&(*C.dxfg_v4_nuam_time_and_sale_t)(event).time_and_sale_base)
	cFree(event)
}

// NuamTradeMapper maps the NuamTrade event.
type NuamTradeMapper struct{}

// CEvent makes a native event from the Go event.
func (NuamTradeMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*trade.NuamTrade)
	n := (*C.dxfg_v4_nuam_trade_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_nuam_trade_t{})))
	tradeToC(a, &n.trade_base.trade_base, &e.TradeBase, C.DXFG_EVENT_NUAM_TRADE)
	n.trade_stat_time = C.int64_t(e.TradeStatTime())
	n.last_significant_price = C.double(e.LastSignificantPrice())
	n.last_price_for_all = C.double(e.LastPriceForAll())
	n.number_of_trades = C.int32_t(e.NumberOfTrades())
	n.vwap = C.double(e.Vwap())
	return unsafe.Pointer(n)
}

// GoEvent makes a Go event from the native event.
func (NuamTradeMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_nuam_trade_t)(native)
	e := trade.NewNuamTrade(C.GoString(n.trade_base.trade_base.market_event.event_symbol))
	tradeFromC(&e.TradeBase, &n.trade_base.trade_base)
	e.SetTradeStatTime(int64(n.trade_stat_time))
	e.SetLastSignificantPrice(float64(n.last_significant_price))
	e.SetLastPriceForAll(float64(n.last_price_for_all))
	e.SetNumberOfTrades(int32(n.number_of_trades))
	e.SetVwap(float64(n.vwap))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (NuamTradeMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_nuam_trade_t)(event)
	cFree(unsafe.Pointer(n.trade_base.trade_base.market_event.event_symbol))
	cFree(event)
}
