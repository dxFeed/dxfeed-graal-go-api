// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package order

import (
	"fmt"
	"math"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// The Nuam flags hold the order type in bits 0 to 3 and the time-in-force type in bits 4 to 7.
const (
	nuamOrderTypeMask        = 15
	nuamOrderTypeShift       = 0
	nuamTimeInForceTypeMask  = 15
	nuamTimeInForceTypeShift = 4
)

// NuamOrderType is the type of a Nuam Exchange order.
type NuamOrderType int32

const (
	// NuamOrderTypeUndefined is an undefined order type.
	NuamOrderTypeUndefined NuamOrderType = iota
	// NuamOrderTypeLimit is a limit order.
	NuamOrderTypeLimit
	// NuamOrderTypeMarket is a market order.
	NuamOrderTypeMarket
	// NuamOrderTypeMarketToLimit is a market-to-limit order.
	NuamOrderTypeMarketToLimit
	// NuamOrderTypeBestOrder is a best order.
	NuamOrderTypeBestOrder
	// NuamOrderTypeImbalance is an imbalance order.
	NuamOrderTypeImbalance
)

// String returns the name of the order type, as in the Java API.
func (t NuamOrderType) String() string {
	switch t {
	case NuamOrderTypeUndefined:
		return "UNDEFINED"
	case NuamOrderTypeLimit:
		return "LIMIT"
	case NuamOrderTypeMarket:
		return "MARKET"
	case NuamOrderTypeMarketToLimit:
		return "MARKET_TO_LIMIT"
	case NuamOrderTypeBestOrder:
		return "BEST_ORDER"
	case NuamOrderTypeImbalance:
		return "IMBALANCE"
	}
	return fmt.Sprintf("NuamOrderType(%d)", int32(t))
}

// NuamTimeInForceType is the time in force of a Nuam Exchange order.
type NuamTimeInForceType int32

const (
	// NuamTimeInForceTypeUndefined is an undefined time in force.
	NuamTimeInForceTypeUndefined NuamTimeInForceType = iota
	// NuamTimeInForceTypeRestOfDay is valid for the rest of the day.
	NuamTimeInForceTypeRestOfDay
	// NuamTimeInForceTypeGtc is good till cancelled.
	NuamTimeInForceTypeGtc
	// NuamTimeInForceTypeImmediateOrCancel is immediate or cancel.
	NuamTimeInForceTypeImmediateOrCancel
	// NuamTimeInForceTypeFillOrKill is fill or kill.
	NuamTimeInForceTypeFillOrKill
	// NuamTimeInForceTypeGoodTilSession is good till the end of the session.
	NuamTimeInForceTypeGoodTilSession
	// NuamTimeInForceTypeNumberOfDays is valid for a number of days, see NuamOrder.TimeInForceData.
	NuamTimeInForceTypeNumberOfDays
)

// String returns the name of the time in force, as in the Java API.
func (t NuamTimeInForceType) String() string {
	switch t {
	case NuamTimeInForceTypeUndefined:
		return "UNDEFINED"
	case NuamTimeInForceTypeRestOfDay:
		return "REST_OF_DAY"
	case NuamTimeInForceTypeGtc:
		return "GTC"
	case NuamTimeInForceTypeImmediateOrCancel:
		return "IMMEDIATE_OR_CANCEL"
	case NuamTimeInForceTypeFillOrKill:
		return "FILL_OR_KILL"
	case NuamTimeInForceTypeGoodTilSession:
		return "GOOD_TIL_SESSION"
	case NuamTimeInForceTypeNumberOfDays:
		return "NUMBER_OF_DAYS"
	}
	return fmt.Sprintf("NuamTimeInForceType(%d)", int32(t))
}

// NuamOrder is an Order of the Nuam Exchange with the Nuam-specific order fields.
type NuamOrder struct {
	Order
	actorId               int32
	participantId         int32
	submitterId           int32
	onBehalfOfSubmitterId int32
	clientOrderId         *string
	customerAccount       *string
	customerInfo          *string
	exchangeInfo          *string
	timeInForceData       int32
	triggerOrderBookId    int32
	triggerPrice          float64
	triggerSessionType    int32
	orderQuantity         float64
	displayQuantity       float64
	refreshQuantity       float64
	leavesQuantity        float64
	matchedQuantity       float64
	nuamFlags             int32
}

// NewNuamOrder creates a NuamOrder for the given event symbol.
func NewNuamOrder(eventSymbol string) *NuamOrder {
	return &NuamOrder{
		Order:           *NewOrder(eventSymbol),
		triggerPrice:    math.NaN(),
		orderQuantity:   math.NaN(),
		displayQuantity: math.NaN(),
		refreshQuantity: math.NaN(),
		leavesQuantity:  math.NaN(),
		matchedQuantity: math.NaN(),
	}
}

// Type returns the event code of this event type.
func (o *NuamOrder) Type() eventcodes.EventCode {
	return eventcodes.NuamOrder
}

// ActorId returns the actor ID.
func (o *NuamOrder) ActorId() int32 { return o.actorId }

// SetActorId sets the actor ID.
func (o *NuamOrder) SetActorId(value int32) { o.actorId = value }

// ParticipantId returns the participant ID.
func (o *NuamOrder) ParticipantId() int32 { return o.participantId }

// SetParticipantId sets the participant ID.
func (o *NuamOrder) SetParticipantId(value int32) { o.participantId = value }

// SubmitterId returns the submitter ID.
func (o *NuamOrder) SubmitterId() int32 { return o.submitterId }

// SetSubmitterId sets the submitter ID.
func (o *NuamOrder) SetSubmitterId(value int32) { o.submitterId = value }

// OnBehalfOfSubmitterId returns the ID of the submitter on whose behalf the order was submitted.
func (o *NuamOrder) OnBehalfOfSubmitterId() int32 { return o.onBehalfOfSubmitterId }

// SetOnBehalfOfSubmitterId sets the ID of the submitter on whose behalf the order was submitted.
func (o *NuamOrder) SetOnBehalfOfSubmitterId(value int32) { o.onBehalfOfSubmitterId = value }

// ClientOrderId returns the client order ID, or nil.
func (o *NuamOrder) ClientOrderId() *string { return o.clientOrderId }

// SetClientOrderId sets the client order ID.
func (o *NuamOrder) SetClientOrderId(value *string) { o.clientOrderId = value }

// CustomerAccount returns the customer account, or nil.
func (o *NuamOrder) CustomerAccount() *string { return o.customerAccount }

// SetCustomerAccount sets the customer account.
func (o *NuamOrder) SetCustomerAccount(value *string) { o.customerAccount = value }

// CustomerInfo returns the customer information, or nil.
func (o *NuamOrder) CustomerInfo() *string { return o.customerInfo }

// SetCustomerInfo sets the customer information.
func (o *NuamOrder) SetCustomerInfo(value *string) { o.customerInfo = value }

// ExchangeInfo returns the exchange information, or nil.
func (o *NuamOrder) ExchangeInfo() *string { return o.exchangeInfo }

// SetExchangeInfo sets the exchange information.
func (o *NuamOrder) SetExchangeInfo(value *string) { o.exchangeInfo = value }

// TimeInForceData returns the data of the time in force, such as the number of days for NumberOfDays.
func (o *NuamOrder) TimeInForceData() int32 { return o.timeInForceData }

// SetTimeInForceData sets the data of the time in force.
func (o *NuamOrder) SetTimeInForceData(value int32) { o.timeInForceData = value }

// TriggerOrderBookId returns the ID of the order book that triggers the order.
func (o *NuamOrder) TriggerOrderBookId() int32 { return o.triggerOrderBookId }

// SetTriggerOrderBookId sets the ID of the order book that triggers the order.
func (o *NuamOrder) SetTriggerOrderBookId(value int32) { o.triggerOrderBookId = value }

// TriggerPrice returns the trigger price.
func (o *NuamOrder) TriggerPrice() float64 { return o.triggerPrice }

// SetTriggerPrice sets the trigger price.
func (o *NuamOrder) SetTriggerPrice(value float64) { o.triggerPrice = value }

// TriggerSessionType returns the trigger session type.
func (o *NuamOrder) TriggerSessionType() int32 { return o.triggerSessionType }

// SetTriggerSessionType sets the trigger session type.
func (o *NuamOrder) SetTriggerSessionType(value int32) { o.triggerSessionType = value }

// OrderQuantity returns the order quantity.
func (o *NuamOrder) OrderQuantity() float64 { return o.orderQuantity }

// SetOrderQuantity sets the order quantity.
func (o *NuamOrder) SetOrderQuantity(value float64) { o.orderQuantity = value }

// DisplayQuantity returns the displayed quantity.
func (o *NuamOrder) DisplayQuantity() float64 { return o.displayQuantity }

// SetDisplayQuantity sets the displayed quantity.
func (o *NuamOrder) SetDisplayQuantity(value float64) { o.displayQuantity = value }

// RefreshQuantity returns the refresh quantity.
func (o *NuamOrder) RefreshQuantity() float64 { return o.refreshQuantity }

// SetRefreshQuantity sets the refresh quantity.
func (o *NuamOrder) SetRefreshQuantity(value float64) { o.refreshQuantity = value }

// LeavesQuantity returns the remaining (leaves) quantity.
func (o *NuamOrder) LeavesQuantity() float64 { return o.leavesQuantity }

// SetLeavesQuantity sets the remaining (leaves) quantity.
func (o *NuamOrder) SetLeavesQuantity(value float64) { o.leavesQuantity = value }

// MatchedQuantity returns the matched quantity.
func (o *NuamOrder) MatchedQuantity() float64 { return o.matchedQuantity }

// SetMatchedQuantity sets the matched quantity.
func (o *NuamOrder) SetMatchedQuantity(value float64) { o.matchedQuantity = value }

// NuamFlags returns the packed Nuam flags: the order type and the time-in-force type.
func (o *NuamOrder) NuamFlags() int32 { return o.nuamFlags }

// SetNuamFlags sets the packed Nuam flags.
func (o *NuamOrder) SetNuamFlags(value int32) { o.nuamFlags = value }

// NuamOrderType returns the type of the order.
func (o *NuamOrder) NuamOrderType() NuamOrderType {
	return NuamOrderType(mathutil.GetBits(int64(o.nuamFlags), nuamOrderTypeMask, nuamOrderTypeShift))
}

// SetNuamOrderType sets the type of the order.
func (o *NuamOrder) SetNuamOrderType(value NuamOrderType) {
	o.nuamFlags = int32(mathutil.SetBits(int64(o.nuamFlags), nuamOrderTypeMask, nuamOrderTypeShift, int64(value)))
}

// TimeInForce returns the time in force of the order.
func (o *NuamOrder) TimeInForce() NuamTimeInForceType {
	return NuamTimeInForceType(mathutil.GetBits(int64(o.nuamFlags), nuamTimeInForceTypeMask,
		nuamTimeInForceTypeShift))
}

// SetTimeInForce sets the time in force of the order.
func (o *NuamOrder) SetTimeInForce(value NuamTimeInForceType) {
	o.nuamFlags = int32(mathutil.SetBits(int64(o.nuamFlags), nuamTimeInForceTypeMask, nuamTimeInForceTypeShift,
		int64(value)))
}

// String returns a text representation of the event.
func (o *NuamOrder) String() string {
	return "NuamOrder{" + o.baseFieldsToString() +
		", marketMaker=" + formatutil.FormatString(o.MarketMaker()) +
		", actorId=" + formatutil.FormatInt64(int64(o.actorId)) +
		", participantId=" + formatutil.FormatInt64(int64(o.participantId)) +
		", submitterId=" + formatutil.FormatInt64(int64(o.submitterId)) +
		", onBehalfOfSubmitterId=" + formatutil.FormatInt64(int64(o.onBehalfOfSubmitterId)) +
		", clientOrderId=" + formatutil.FormatString(o.clientOrderId) +
		", customerAccount=" + formatutil.FormatString(o.customerAccount) +
		", customerInfo=" + formatutil.FormatString(o.customerInfo) +
		", exchangeInfo=" + formatutil.FormatString(o.exchangeInfo) +
		", timeInForce=" + o.TimeInForce().String() +
		", timeInForceData=" + formatutil.FormatInt64(int64(o.timeInForceData)) +
		", triggerOrderBookId=" + formatutil.FormatInt64(int64(o.triggerOrderBookId)) +
		", triggerPrice=" + formatutil.FormatFloat64(o.triggerPrice) +
		", triggerSessionType=" + formatutil.FormatInt64(int64(o.triggerSessionType)) +
		", orderType=" + o.NuamOrderType().String() +
		", orderQuantity=" + formatutil.FormatFloat64(o.orderQuantity) +
		", displayQuantity=" + formatutil.FormatFloat64(o.displayQuantity) +
		", refreshQuantity=" + formatutil.FormatFloat64(o.refreshQuantity) +
		", leavesQuantity=" + formatutil.FormatFloat64(o.leavesQuantity) +
		", matchedQuantity=" + formatutil.FormatFloat64(o.matchedQuantity) +
		"}"
}
