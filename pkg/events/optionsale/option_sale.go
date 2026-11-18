// Package optionsale provides the OptionSale event: a trade of an option, with the option analytics at its time.
package optionsale

import (
	"math"
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
)

// The flags have the same layout as the flags of TimeAndSale.
const (
	tteMask   = 0xff
	tteShift  = 8
	sideMask  = 3
	sideShift = 5
	spreadLeg = 1 << 4
	eth       = 1 << 3
	validTick = 1 << 2
	typeMask  = 3
	typeShift = 0
)

// OptionSale is a trade of an option, or another sale event, with the price, the underlying price, the volatility and
// the delta of the option at the time of the sale. It is an indexed event with the default source: the index is
// unique per symbol and identifies the sale for a later correction or cancellation.
type OptionSale struct {
	eventSymbol            *string
	eventTime              int64
	eventFlags             int32
	index                  int64
	timeSequence           int64
	timeNanoPart           int32
	exchangeCode           int16
	price                  float64
	size                   float64
	bidPrice               float64
	askPrice               float64
	exchangeSaleConditions *string
	flags                  int32
	underlyingPrice        float64
	volatility             float64
	delta                  float64
	optionSymbol           *string
}

// NewOptionSale creates an OptionSale for the given event symbol.
func NewOptionSale(eventSymbol string) *OptionSale {
	return &OptionSale{
		eventSymbol:     &eventSymbol,
		price:           math.NaN(),
		size:            math.NaN(),
		bidPrice:        math.NaN(),
		askPrice:        math.NaN(),
		underlyingPrice: math.NaN(),
		volatility:      math.NaN(),
		delta:           math.NaN(),
	}
}

// Type returns the event code of this event type.
func (o *OptionSale) Type() eventcodes.EventCode {
	return eventcodes.OptionSale
}

// Source is always DEFAULT for OptionSale.
func (o *OptionSale) Source() *events.IndexedEventSource {
	return events.DefaultIndexedEventSource()
}

// EventSymbol returns the symbol of this event.
func (o *OptionSale) EventSymbol() *string {
	return o.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (o *OptionSale) SetEventSymbol(eventSymbol string) {
	o.eventSymbol = &eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (o *OptionSale) EventTime() int64 {
	return o.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (o *OptionSale) SetEventTime(eventTime int64) {
	o.eventTime = eventTime
}

// EventFlags returns the transactional event flags; see events.EventFlag.
func (o *OptionSale) EventFlags() int32 {
	return o.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (o *OptionSale) SetEventFlags(eventFlags int32) {
	o.eventFlags = eventFlags
}

// Index returns the unique per-symbol index of this sale.
func (o *OptionSale) Index() int64 {
	return o.index
}

// SetIndex sets the unique per-symbol index of this sale.
func (o *OptionSale) SetIndex(index int64) {
	o.index = index
}

// TimeSequence returns the time and the sequence of this event packed in one value.
func (o *OptionSale) TimeSequence() int64 {
	return o.timeSequence
}

// SetTimeSequence sets the time and the sequence of this event packed in one value.
func (o *OptionSale) SetTimeSequence(timeSequence int64) {
	o.timeSequence = timeSequence
}

// Time returns the time of the sale in milliseconds since the Unix epoch.
func (o *OptionSale) Time() int64 {
	return timeutil.TimeOfTimeSequence(o.timeSequence)
}

// SetTime sets the time of the sale in milliseconds since the Unix epoch; the sequence is kept.
func (o *OptionSale) SetTime(timeMillis int64) {
	o.timeSequence = timeutil.WithTime(o.timeSequence, timeMillis)
}

// TimeNanos returns the time of the sale in nanoseconds since the Unix epoch.
func (o *OptionSale) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(o.Time(), o.timeNanoPart)
}

// SetTimeNanos sets the time of the sale in nanoseconds since the Unix epoch.
func (o *OptionSale) SetTimeNanos(timeNanos int64) {
	o.SetTime(timeutil.GetMillisFromNanos(timeNanos))
	o.timeNanoPart = int32(timeutil.GetNanoPartFromNanos(timeNanos))
}

// TimeNanoPart returns the microseconds and nanoseconds part of the time of the sale.
func (o *OptionSale) TimeNanoPart() int32 {
	return o.timeNanoPart
}

// SetTimeNanoPart sets the microseconds and nanoseconds part of the time of the sale.
func (o *OptionSale) SetTimeNanoPart(timeNanoPart int32) {
	o.timeNanoPart = timeNanoPart
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (o *OptionSale) Sequence() int32 {
	return timeutil.SequenceOfTimeSequence(o.timeSequence)
}

// SetSequence sets the sequence number of this event; it returns an error if the sequence is out of range.
func (o *OptionSale) SetSequence(sequence int32) error {
	timeSequence, err := timeutil.WithSequence(o.timeSequence, sequence)
	if err != nil {
		return err
	}
	o.timeSequence = timeSequence
	return nil
}

// ExchangeCode returns the exchange code of the sale.
func (o *OptionSale) ExchangeCode() int16 {
	return o.exchangeCode
}

// SetExchangeCode sets the exchange code of the sale.
func (o *OptionSale) SetExchangeCode(exchangeCode int16) {
	o.exchangeCode = exchangeCode
}

// Price returns the price of the sale.
func (o *OptionSale) Price() float64 {
	return o.price
}

// SetPrice sets the price of the sale.
func (o *OptionSale) SetPrice(price float64) {
	o.price = price
}

// Size returns the size of the sale.
func (o *OptionSale) Size() float64 {
	return o.size
}

// SetSize sets the size of the sale.
func (o *OptionSale) SetSize(size float64) {
	o.size = size
}

// BidPrice returns the bid price on the market at the time of the sale.
func (o *OptionSale) BidPrice() float64 {
	return o.bidPrice
}

// SetBidPrice sets the bid price on the market at the time of the sale.
func (o *OptionSale) SetBidPrice(bidPrice float64) {
	o.bidPrice = bidPrice
}

// AskPrice returns the ask price on the market at the time of the sale.
func (o *OptionSale) AskPrice() float64 {
	return o.askPrice
}

// SetAskPrice sets the ask price on the market at the time of the sale.
func (o *OptionSale) SetAskPrice(askPrice float64) {
	o.askPrice = askPrice
}

// ExchangeSaleConditions returns the sale conditions provided for this event by the data feed, or nil.
func (o *OptionSale) ExchangeSaleConditions() *string {
	return o.exchangeSaleConditions
}

// SetExchangeSaleConditions sets the sale conditions of this event.
func (o *OptionSale) SetExchangeSaleConditions(exchangeSaleConditions *string) {
	o.exchangeSaleConditions = exchangeSaleConditions
}

// Flags returns the packed flags of this event; they have the same layout as the flags of TimeAndSale.
func (o *OptionSale) Flags() int32 {
	return o.flags
}

// SetFlags sets the packed flags of this event.
func (o *OptionSale) SetFlags(flags int32) {
	o.flags = flags
}

// TradeThroughExempt returns the trade-through-exempt flag as an ASCII character code.
func (o *OptionSale) TradeThroughExempt() rune {
	return rune(mathutil.GetBits(int64(o.flags), tteMask, tteShift))
}

// SetTradeThroughExempt sets the trade-through-exempt flag as an ASCII character code.
func (o *OptionSale) SetTradeThroughExempt(value rune) {
	o.flags = int32(mathutil.SetBits(int64(o.flags), tteMask, tteShift, int64(value)))
}

// AggressorSide returns the side of the aggressor that initiated the trade.
func (o *OptionSale) AggressorSide() side.Side {
	return side.SideValueOf(mathutil.GetBits(int64(o.flags), sideMask, sideShift))
}

// SetAggressorSide sets the side of the aggressor that initiated the trade.
func (o *OptionSale) SetAggressorSide(value side.Side) {
	o.flags = int32(mathutil.SetBits(int64(o.flags), sideMask, sideShift, int64(value)))
}

// IsSpreadLeg reports whether this sale is a spread leg.
func (o *OptionSale) IsSpreadLeg() bool {
	return o.flags&spreadLeg != 0
}

// SetIsSpreadLeg sets the spread leg flag.
func (o *OptionSale) SetIsSpreadLeg(value bool) {
	o.setFlag(spreadLeg, value)
}

// IsExtendedTradingHours reports whether the sale happened during extended trading hours (ETH).
func (o *OptionSale) IsExtendedTradingHours() bool {
	return o.flags&eth != 0
}

// SetIsExtendedTradingHours sets the extended trading hours (ETH) flag.
func (o *OptionSale) SetIsExtendedTradingHours(value bool) {
	o.setFlag(eth, value)
}

// IsValidTick reports whether this event is a valid tick for the intraday price chart.
func (o *OptionSale) IsValidTick() bool {
	return o.flags&validTick != 0
}

// SetIsValidTick sets the valid tick flag.
func (o *OptionSale) SetIsValidTick(value bool) {
	o.setFlag(validTick, value)
}

// SaleType returns the type of this event: new, correction or cancel.
func (o *OptionSale) SaleType() timeandsale.Type {
	return timeandsale.TypeValueOf(mathutil.GetBits(int64(o.flags), typeMask, typeShift))
}

// SetSaleType sets the type of this event: new, correction or cancel.
func (o *OptionSale) SetSaleType(value timeandsale.Type) {
	o.flags = int32(mathutil.SetBits(int64(o.flags), typeMask, typeShift, int64(value)))
}

// IsNew reports whether this is a new event.
func (o *OptionSale) IsNew() bool {
	return o.SaleType() == timeandsale.TypeNew
}

// IsCorrection reports whether this event corrects a previous one.
func (o *OptionSale) IsCorrection() bool {
	return o.SaleType() == timeandsale.TypeCorrection
}

// IsCancel reports whether this event cancels a previous one.
func (o *OptionSale) IsCancel() bool {
	return o.SaleType() == timeandsale.TypeCancel
}

// UnderlyingPrice returns the underlying price at the time of the sale.
func (o *OptionSale) UnderlyingPrice() float64 {
	return o.underlyingPrice
}

// SetUnderlyingPrice sets the underlying price at the time of the sale.
func (o *OptionSale) SetUnderlyingPrice(underlyingPrice float64) {
	o.underlyingPrice = underlyingPrice
}

// Volatility returns the Black-Scholes implied volatility of the option at the time of the sale.
func (o *OptionSale) Volatility() float64 {
	return o.volatility
}

// SetVolatility sets the Black-Scholes implied volatility of the option at the time of the sale.
func (o *OptionSale) SetVolatility(volatility float64) {
	o.volatility = volatility
}

// Delta returns the delta of the option at the time of the sale.
func (o *OptionSale) Delta() float64 {
	return o.delta
}

// SetDelta sets the delta of the option at the time of the sale.
func (o *OptionSale) SetDelta(delta float64) {
	o.delta = delta
}

// OptionSymbol returns the option symbol of the sale, or nil.
func (o *OptionSale) OptionSymbol() *string {
	return o.optionSymbol
}

// SetOptionSymbol sets the option symbol of the sale.
func (o *OptionSale) SetOptionSymbol(optionSymbol *string) {
	o.optionSymbol = optionSymbol
}

func (o *OptionSale) setFlag(flag int32, value bool) {
	if value {
		o.flags |= flag
	} else {
		o.flags &^= flag
	}
}

// String returns a text representation of the event.
func (o *OptionSale) String() string {
	return "OptionSale{" + formatutil.FormatString(o.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(o.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(o.EventFlags())) +
		", index=" + formatutil.HexFormat(o.Index()) +
		", time=" + formatutil.FormatTime(o.Time()) +
		", timeNanoPart=" + strconv.FormatInt(int64(o.timeNanoPart), 10) +
		", sequence=" + strconv.FormatInt(int64(o.Sequence()), 10) +
		", exchange=" + formatutil.FormatChar(rune(o.exchangeCode)) +
		", price=" + formatutil.FormatFloat64(o.price) +
		", size=" + formatutil.FormatFloat64(o.size) +
		", bid=" + formatutil.FormatFloat64(o.bidPrice) +
		", ask=" + formatutil.FormatFloat64(o.askPrice) +
		", ESC=" + formatutil.FormatString(o.exchangeSaleConditions) +
		", TTE=" + formatutil.FormatChar(o.TradeThroughExempt()) +
		", side=" + o.AggressorSide().String() +
		", spread=" + formatutil.FormatBool(o.IsSpreadLeg()) +
		", ETH=" + formatutil.FormatBool(o.IsExtendedTradingHours()) +
		", validTick=" + formatutil.FormatBool(o.IsValidTick()) +
		", type=" + formatutil.FormatInt64(int64(o.SaleType())) +
		", underlyingPrice=" + formatutil.FormatFloat64(o.underlyingPrice) +
		", volatility=" + formatutil.FormatFloat64(o.volatility) +
		", delta=" + formatutil.FormatFloat64(o.delta) +
		", optionSymbol=" + formatutil.FormatString(o.optionSymbol) +
		"}"
}
