package order

import (
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// The OTC Markets flags.
const (
	otcOpen           = 1
	otcUnsolicited    = 1 << 1
	otcPriceTypeMask  = 3
	otcPriceTypeShift = 2
	otcSaturated      = 1 << 4
	otcAutoExecution  = 1 << 5
	otcNmsConditional = 1 << 6
)

// OtcMarketsPriceType is the price type of an OTC Markets quote.
type OtcMarketsPriceType int32

const (
	// OtcMarketsPriceTypeUnpriced is a quote without a price: an indication of interest.
	OtcMarketsPriceTypeUnpriced OtcMarketsPriceType = iota
	// OtcMarketsPriceTypeActual is a quote with an actual price.
	OtcMarketsPriceTypeActual
	// OtcMarketsPriceTypeWanted is a quote that wants an offer (OW) or a bid (BW) and has no price.
	OtcMarketsPriceTypeWanted
)

// String returns the name of the price type, as in the Java API.
func (t OtcMarketsPriceType) String() string {
	switch t {
	case OtcMarketsPriceTypeUnpriced:
		return "UNPRICED"
	case OtcMarketsPriceTypeActual:
		return "ACTUAL"
	case OtcMarketsPriceTypeWanted:
		return "WANTED"
	}
	return fmt.Sprintf("OtcMarketsPriceType(%d)", int32(t))
}

// OtcMarketsOrder is an Order with the OTC Markets fields: a bid or an ask of an OTC Markets quote. The original
// quote is presented as two orders, one for each side.
type OtcMarketsOrder struct {
	Order
	quoteAccessPayment int32
	otcMarketsFlags    int32
}

// NewOtcMarketsOrder creates an OtcMarketsOrder for the given event symbol.
func NewOtcMarketsOrder(eventSymbol string) *OtcMarketsOrder {
	return &OtcMarketsOrder{Order: *NewOrder(eventSymbol)}
}

// Type returns the event code of this event type.
func (o *OtcMarketsOrder) Type() eventcodes.EventCode {
	return eventcodes.OtcMarketsOrder
}

// QuoteAccessPayment returns the Quote Access Payment (QAP): from 1 to 30 is a rebate, from -1 to -30 is an access
// fee, 0 is neither.
func (o *OtcMarketsOrder) QuoteAccessPayment() int32 {
	return o.quoteAccessPayment
}

// SetQuoteAccessPayment sets the Quote Access Payment (QAP).
func (o *OtcMarketsOrder) SetQuoteAccessPayment(quoteAccessPayment int32) {
	o.quoteAccessPayment = quoteAccessPayment
}

// OtcMarketsFlags returns the packed OTC Markets flags.
func (o *OtcMarketsOrder) OtcMarketsFlags() int32 {
	return o.otcMarketsFlags
}

// SetOtcMarketsFlags sets the packed OTC Markets flags.
func (o *OtcMarketsOrder) SetOtcMarketsFlags(otcMarketsFlags int32) {
	o.otcMarketsFlags = otcMarketsFlags
}

// IsOpen reports whether the quote is open: it can be executed.
func (o *OtcMarketsOrder) IsOpen() bool {
	return o.otcMarketsFlags&otcOpen != 0
}

// SetOpen sets whether the quote is open.
func (o *OtcMarketsOrder) SetOpen(value bool) {
	o.setOtcFlag(otcOpen, value)
}

// IsUnsolicited reports whether the quote is unsolicited.
func (o *OtcMarketsOrder) IsUnsolicited() bool {
	return o.otcMarketsFlags&otcUnsolicited != 0
}

// SetUnsolicited sets whether the quote is unsolicited.
func (o *OtcMarketsOrder) SetUnsolicited(value bool) {
	o.setOtcFlag(otcUnsolicited, value)
}

// OtcMarketsPriceType returns the price type of the quote.
func (o *OtcMarketsOrder) OtcMarketsPriceType() OtcMarketsPriceType {
	return OtcMarketsPriceType(mathutil.GetBits(int64(o.otcMarketsFlags), otcPriceTypeMask, otcPriceTypeShift))
}

// SetOtcMarketsPriceType sets the price type of the quote.
func (o *OtcMarketsOrder) SetOtcMarketsPriceType(value OtcMarketsPriceType) {
	o.otcMarketsFlags = int32(mathutil.SetBits(int64(o.otcMarketsFlags), otcPriceTypeMask, otcPriceTypeShift,
		int64(value)))
}

// IsSaturated reports whether the quote is saturated: it was replaced by a quote of the same side.
func (o *OtcMarketsOrder) IsSaturated() bool {
	return o.otcMarketsFlags&otcSaturated != 0
}

// SetSaturated sets whether the quote is saturated.
func (o *OtcMarketsOrder) SetSaturated(value bool) {
	o.setOtcFlag(otcSaturated, value)
}

// IsAutoExecution reports whether the quote can be executed automatically.
func (o *OtcMarketsOrder) IsAutoExecution() bool {
	return o.otcMarketsFlags&otcAutoExecution != 0
}

// SetAutoExecution sets whether the quote can be executed automatically.
func (o *OtcMarketsOrder) SetAutoExecution(value bool) {
	o.setOtcFlag(otcAutoExecution, value)
}

// IsNmsConditional reports whether the quote is an NMS conditional quote.
func (o *OtcMarketsOrder) IsNmsConditional() bool {
	return o.otcMarketsFlags&otcNmsConditional != 0
}

// SetNmsConditional sets whether the quote is an NMS conditional quote.
func (o *OtcMarketsOrder) SetNmsConditional(value bool) {
	o.setOtcFlag(otcNmsConditional, value)
}

func (o *OtcMarketsOrder) setOtcFlag(flag int32, value bool) {
	if value {
		o.otcMarketsFlags |= flag
	} else {
		o.otcMarketsFlags &^= flag
	}
}

// String returns a text representation of the event.
func (o *OtcMarketsOrder) String() string {
	return "OtcMarketsOrder{" + o.baseFieldsToString() +
		", marketMaker=" + formatutil.FormatString(o.MarketMaker()) +
		", QAP=" + formatutil.FormatInt64(int64(o.quoteAccessPayment)) +
		", open=" + formatutil.FormatBool(o.IsOpen()) +
		", unsolicited=" + formatutil.FormatBool(o.IsUnsolicited()) +
		", priceType=" + o.OtcMarketsPriceType().String() +
		", saturated=" + formatutil.FormatBool(o.IsSaturated()) +
		", autoEx=" + formatutil.FormatBool(o.IsAutoExecution()) +
		", NMS=" + formatutil.FormatBool(o.IsNmsConditional()) +
		"}"
}
