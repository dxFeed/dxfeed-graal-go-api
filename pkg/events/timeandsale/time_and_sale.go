package timeandsale

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"math"
	"strconv"
)

const (
	// TTE (TradeThroughExempt) values are ASCII chars in [0, 255].
	tteMask  = 0xff
	tteShift = 8

	// SIDE values are taken from Side enum.
	sideMask  = 3
	sideShift = 5

	spreadLeg = 1 << 4
	eth       = 1 << 3
	validTick = 1 << 2

	// TYPE values are taken from TimeAndSaleType enum.
	typeMask  = 3
	typeShift = 0
)

// TimeAndSale represents a trade or other market event with a price, such as a trade or a correction, in the time and sale feed.
type TimeAndSale struct {
	eventSymbol            *string
	eventTime              int64
	timeNanoPart           int32
	exchangeCode           int16
	price                  float64
	size                   float64
	bidPrice               float64
	askPrice               float64
	exchangeSaleConditions *string
	flags                  int32
	buyer                  *string
	seller                 *string
	eventFlags             int32
	index                  int64
}

const maxSequence = (1 << 22) - 1

// NewTimeAndSale creates a new TimeAndSale for the given event symbol.
func NewTimeAndSale(eventSymbol string) *TimeAndSale {
	return &TimeAndSale{
		eventSymbol: &eventSymbol,
		price:       math.NaN(),
		size:        math.NaN(),
		bidPrice:    math.NaN(),
		askPrice:    math.NaN(),
	}
}

// Type returns the event code of this event.
func (t *TimeAndSale) Type() eventcodes.EventCode {
	return eventcodes.TimeAndSale
}

// EventSymbol returns the symbol of this event.
func (t *TimeAndSale) EventSymbol() *string {
	return t.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (t *TimeAndSale) SetEventSymbol(value string) {
	*t.eventSymbol = value
}

// EventTime returns the source-provided time of this event in milliseconds since the Unix epoch.
func (t *TimeAndSale) EventTime() int64 {
	return t.eventTime
}

// SetEventTime sets the source-provided time of this event in milliseconds since the Unix epoch.
func (t *TimeAndSale) SetEventTime(value int64) {
	t.eventTime = value
}

// TimeNanoPart returns the microseconds and nanoseconds part of the event time.
func (t *TimeAndSale) TimeNanoPart() int32 {
	return t.timeNanoPart
}

// SetTimeNanoPart sets the microseconds and nanoseconds part of the event time.
func (t *TimeAndSale) SetTimeNanoPart(value int32) {
	t.timeNanoPart = value
}

// ExchangeCode returns the exchange code of the event.
func (t *TimeAndSale) ExchangeCode() int16 {
	return t.exchangeCode
}

// SetExchangeCode sets the exchange code of the event.
func (t *TimeAndSale) SetExchangeCode(value int16) {
	t.exchangeCode = value
}

// Price returns the price of the last trade.
func (t *TimeAndSale) Price() float64 {
	return t.price
}

// SetPrice sets the price of the last trade.
func (t *TimeAndSale) SetPrice(value float64) {
	t.price = value
}

// Size returns the size of the last trade.
func (t *TimeAndSale) Size() float64 {
	return t.size
}

// SetSize sets the size of the last trade.
func (t *TimeAndSale) SetSize(value float64) {
	t.size = value
}

// BidPrice returns the current bid price on the market when this event occurred.
func (t *TimeAndSale) BidPrice() float64 {
	return t.bidPrice
}

// SetBidPrice sets the current bid price on the market when this event occurred.
func (t *TimeAndSale) SetBidPrice(value float64) {
	t.bidPrice = value
}

// AskPrice returns the current ask price on the market when this event occurred.
func (t *TimeAndSale) AskPrice() float64 {
	return t.askPrice
}

// SetAskPrice sets the current ask price on the market when this event occurred.
func (t *TimeAndSale) SetAskPrice(value float64) {
	t.askPrice = value
}

// ExchangeSaleConditions returns the sale conditions provided by the exchange.
func (t *TimeAndSale) ExchangeSaleConditions() *string {
	return t.exchangeSaleConditions
}

// SetExchangeSaleConditions sets the sale conditions provided by the exchange.
func (t *TimeAndSale) SetExchangeSaleConditions(value *string) {
	t.exchangeSaleConditions = value
}

// Buyer returns the buyer of this trade, if known.
func (t *TimeAndSale) Buyer() *string {
	return t.buyer
}

// SetBuyer sets the buyer of this trade.
func (t *TimeAndSale) SetBuyer(value *string) {
	t.buyer = value
}

// Seller returns the seller of this trade, if known.
func (t *TimeAndSale) Seller() *string {
	return t.seller
}

// SetSeller sets the seller of this trade.
func (t *TimeAndSale) SetSeller(value *string) {
	t.seller = value
}

// EventFlags returns the transactional event flags.
func (t *TimeAndSale) EventFlags() int32 {
	return t.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (t *TimeAndSale) SetEventFlags(value int32) {
	t.eventFlags = value
}

// Index returns the unique per-symbol index of this event, combining time and sequence.
func (t *TimeAndSale) Index() int64 {
	return t.index
}

// SetIndex sets the unique per-symbol index of this event.
func (t *TimeAndSale) SetIndex(value int64) {
	t.index = value
}

// Sequence returns the sequence number of this event within the same millisecond.
func (t *TimeAndSale) Sequence() int64 {
	return t.index & maxSequence
}

// SetSequence sets the sequence number of this event within the same millisecond; it returns an error if out of range.
func (t *TimeAndSale) SetSequence(value int64) error {
	if value < 0 || value > maxSequence {
		return fmt.Errorf("Sequence(%d) is < 0 or > MaxSequence(%d)", value, maxSequence)
	}

	t.index = (t.index & ^maxSequence) | value
	return nil
}

// Time returns the time of this event in milliseconds since the Unix epoch.
func (t *TimeAndSale) Time() int64 {
	return ((t.index >> 32) * 1000) + ((t.index >> 22) & 0x3ff)
}

// SetTime sets the time of this event in milliseconds since the Unix epoch.
func (t *TimeAndSale) SetTime(value int64) {
	t.index = (timeutil.GetSecondsFromTime(value) << 32) |
		(int64(timeutil.GetMillisFromTime(value)) << 22) |
		t.Sequence()
}

// TimeNanos returns the time of this event in nanoseconds since the Unix epoch.
func (t *TimeAndSale) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(t.Time(), t.TimeNanoPart())
}

// SetTimeNanos sets the time of this event in nanoseconds since the Unix epoch.
func (t *TimeAndSale) SetTimeNanos(value int64) {
	t.SetTime(timeutil.GetMillisFromNanos(value))
	t.SetTimeNanoPart(int32(timeutil.GetNanoPartFromNanos(value)))
}

// TradeThroughExempt returns the trade-through-exempt flag as an ASCII character code.
func (t *TimeAndSale) TradeThroughExempt() rune {
	return rune(mathutil.GetBits(int64(t.flags), tteMask, tteShift))
}

// SetTradeThroughExempt sets the trade-through-exempt flag as an ASCII character code.
func (t *TimeAndSale) SetTradeThroughExempt(value rune) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.flags), tteMask, tteShift, int64(value))))
}

// AggressorSide returns the side of the aggressor that initiated the trade.
func (t *TimeAndSale) AggressorSide() side.Side {
	bits := mathutil.GetBits(int64(t.flags), sideMask, sideShift)
	return side.SideValueOf(bits)
}

// SetAggressorSide sets the side of the aggressor that initiated the trade.
func (t *TimeAndSale) SetAggressorSide(value side.Side) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.flags), sideMask, sideShift, int64(value))))
}

// IsSpreadLeg reports whether this trade is a spread leg.
func (t *TimeAndSale) IsSpreadLeg() bool {
	return (t.flags & spreadLeg) != 0
}

// SetIsSpreadLeg sets the spread leg flag.
func (t *TimeAndSale) SetIsSpreadLeg(value bool) {
	if value {
		t.SetFlags(t.flags | spreadLeg)
	} else {
		t.SetFlags(t.flags & ^spreadLeg)
	}
}

// IsExtendedTradingHours reports whether the trade occurred during extended trading hours (ETH).
func (t *TimeAndSale) IsExtendedTradingHours() bool {
	return (t.flags & eth) != 0
}

// SetIsExtendedTradingHours sets the extended trading hours (ETH) flag.
func (t *TimeAndSale) SetIsExtendedTradingHours(value bool) {
	if value {
		t.SetFlags(t.flags | eth)
	} else {
		t.SetFlags(t.flags & ^eth)
	}
}

// IsValidTick reports whether this event represents a valid tick for the intraday price chart.
func (t *TimeAndSale) IsValidTick() bool {
	return (t.flags & validTick) != 0
}

// SetIsValidTick sets the valid tick flag.
func (t *TimeAndSale) SetIsValidTick(value bool) {
	if value {
		t.SetFlags(t.flags | validTick)
	} else {
		t.SetFlags(t.flags & ^validTick)
	}
}

// TimeAndSaleType returns the type of this event: new, correction or cancel.
func (t *TimeAndSale) TimeAndSaleType() Type {
	return TypeValueOf(mathutil.GetBits(int64(t.flags), typeMask, typeShift))
}

// SetTimeAndSaleType sets the type of this event: new, correction or cancel.
func (t *TimeAndSale) SetTimeAndSaleType(value Type) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.flags), typeMask, typeShift, int64(value))))
}

// IsNew reports whether this is a new event.
func (t *TimeAndSale) IsNew() bool {
	return t.TimeAndSaleType() == TypeNew
}

// IsCorrection reports whether this event corrects a previous one.
func (t *TimeAndSale) IsCorrection() bool {
	return t.TimeAndSaleType() == TypeCorrection
}

// IsCancel reports whether this event cancels a previous one.
func (t *TimeAndSale) IsCancel() bool {
	return t.TimeAndSaleType() == TypeCancel
}

// SetFlags sets the packed flags of this event.
func (t *TimeAndSale) SetFlags(value int32) {
	t.flags = value
}

// Flags returns the packed flags of this event.
func (t *TimeAndSale) Flags() int32 {
	return t.flags
}

// String returns a string representation of this event.
func (t *TimeAndSale) String() string {
	return "TimeAndSale{" + formatutil.FormatString(t.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(t.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(t.EventFlags())) +
		", time=" + formatutil.FormatTime(t.Time()) +
		", timeNanoPart=" + strconv.FormatInt(int64(t.timeNanoPart), 10) +
		", sequence=" + strconv.FormatInt(int64(t.Sequence()), 10) +
		", exchange=" + formatutil.FormatChar(rune(t.exchangeCode)) +
		", price=" + formatutil.FormatFloat64(t.Price()) +
		", size=" + formatutil.FormatFloat64(t.Size()) +
		", bid=" + formatutil.FormatFloat64(t.BidPrice()) +
		", ask=" + formatutil.FormatFloat64(t.AskPrice()) +
		", ESC=" + formatutil.FormatString(t.ExchangeSaleConditions()) +
		", TTE=" + formatutil.FormatChar(t.TradeThroughExempt()) +
		", side=" + t.AggressorSide().String() +
		", spread=" + formatutil.FormatBool(t.IsSpreadLeg()) +
		", ETH=" + formatutil.FormatBool(t.IsExtendedTradingHours()) +
		", validTick=" + formatutil.FormatBool(t.IsValidTick()) +
		", type=" + formatutil.FormatInt64(int64(t.TimeAndSaleType())) +
		formatutil.FormatNullableString(", buyer=%s", t.Buyer(), "") +
		formatutil.FormatNullableString(", seller=%s", t.Seller(), "") +
		"}"
}
