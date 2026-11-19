// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package order

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// Order is a snapshot of a limit order, quote or price level with an optional market maker.
type Order struct {
	Base
	marketMaker *string
}

// NewOrder creates an Order for the given event symbol.
func NewOrder(eventSymbol string) *Order {
	return &Order{Base: *NewBase(eventSymbol)}
}

// MarketMaker returns the market maker or other aggregate identifier of the order.
func (o *Order) MarketMaker() *string {
	return o.marketMaker
}

// SetMarketMaker sets the market maker or other aggregate identifier of the order.
func (o *Order) SetMarketMaker(marketMaker *string) {
	o.marketMaker = marketMaker
}

// Type returns the event code of this event type.
func (o *Order) Type() eventcodes.EventCode {
	return eventcodes.Order
}

// String returns a text representation of the event.
func (o *Order) String() string {
	return "Order{" + o.baseFieldsToString() +
		", marketMaker=" + formatutil.FormatString(o.MarketMaker()) +
		"}"
}
