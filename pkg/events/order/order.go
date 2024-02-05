package order

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
)

type Order struct {
	Base
	marketMaker *string
}

func NewOrder(eventSymbol string) *Order {
	return &Order{Base: *NewBase(eventSymbol)}
}

func (o *Order) MarketMaker() *string {
	return o.marketMaker
}

func (o *Order) SetMarketMaker(marketMaker *string) {
	o.marketMaker = marketMaker
}

func (o *Order) Type() eventcodes.EventCode {
	return eventcodes.Order
}

func (o *Order) String() string {
	return "Order{" + o.baseFieldsToString() +
		", marketMaker=" + formatutil.FormatString(o.MarketMaker()) +
		"}"
}
