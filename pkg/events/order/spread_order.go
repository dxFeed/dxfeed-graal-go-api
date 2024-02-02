package order

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
)

type SpreadOrder struct {
	Base
	spreadSymbol *string
}

func NewSpreadOrder(eventSymbol string) *SpreadOrder {
	return &SpreadOrder{Base: Base{eventSymbol: &eventSymbol}}
}

func (o *SpreadOrder) SpreadSymbol() *string {
	return o.spreadSymbol
}

func (o *SpreadOrder) SetSpreadSymbol(spreadSymbol *string) {
	o.spreadSymbol = spreadSymbol
}

func (o *SpreadOrder) Type() eventcodes.EventCode {
	return eventcodes.SpreadOrder
}

func (o *SpreadOrder) String() string {
	return "SpreadOrder{" + o.baseFieldsToString() +
		", spreadSymbol=" + formatutil.FormatString(o.SpreadSymbol()) +
		"}"
}
