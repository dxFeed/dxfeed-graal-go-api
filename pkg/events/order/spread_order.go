package order

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// SpreadOrder is an Order that represents a multi-leg spread.
type SpreadOrder struct {
	Base
	spreadSymbol *string
}

// NewSpreadOrder creates a SpreadOrder for the given event symbol.
func NewSpreadOrder(eventSymbol string) *SpreadOrder {
	return &SpreadOrder{Base: *NewBase(eventSymbol)}
}

// SpreadSymbol returns the symbol of the spread this order belongs to.
func (o *SpreadOrder) SpreadSymbol() *string {
	return o.spreadSymbol
}

// SetSpreadSymbol sets the symbol of the spread this order belongs to.
func (o *SpreadOrder) SetSpreadSymbol(spreadSymbol *string) {
	o.spreadSymbol = spreadSymbol
}

// Type returns the event code of this event type.
func (o *SpreadOrder) Type() eventcodes.EventCode {
	return eventcodes.SpreadOrder
}

// String returns a text representation of the event.
func (o *SpreadOrder) String() string {
	return "SpreadOrder{" + o.baseFieldsToString() +
		", spreadSymbol=" + formatutil.FormatString(o.SpreadSymbol()) +
		"}"
}
