// Package theoprice provides the TheoPrice event: the theoretical price of an option computed by a pricing model.
package theoprice

import (
	"math"
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// TheoPrice is a time series event with the theoretical price of an option and the values of the pricing model it
// was computed with: underlying price, delta, gamma, dividend and interest. The time and the sequence are packed in
// the index.
type TheoPrice struct {
	eventSymbol *string
	eventTime   int64
	eventFlags  int32
	index       int64

	price           float64
	underlyingPrice float64
	delta           float64
	gamma           float64
	dividend        float64
	interest        float64
}

// NewTheoPrice creates a TheoPrice event for the given event symbol.
func NewTheoPrice(eventSymbol string) *TheoPrice {
	return &TheoPrice{
		eventSymbol:     &eventSymbol,
		price:           math.NaN(),
		underlyingPrice: math.NaN(),
		delta:           math.NaN(),
		gamma:           math.NaN(),
		dividend:        math.NaN(),
		interest:        math.NaN(),
	}
}

// Type returns the event code of this event type.
func (p *TheoPrice) Type() eventcodes.EventCode {
	return eventcodes.TheoPrice
}

// Source is always DEFAULT for TheoPrice.
func (p *TheoPrice) Source() *events.IndexedEventSource {
	return events.DefaultIndexedEventSource()
}

// EventSymbol returns the symbol of this event.
func (p *TheoPrice) EventSymbol() *string {
	return p.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (p *TheoPrice) SetEventSymbol(eventSymbol string) {
	p.eventSymbol = &eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (p *TheoPrice) EventTime() int64 {
	return p.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (p *TheoPrice) SetEventTime(eventTime int64) {
	p.eventTime = eventTime
}

// EventFlags returns the transactional event flags; see events.EventFlag.
func (p *TheoPrice) EventFlags() int32 {
	return p.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (p *TheoPrice) SetEventFlags(eventFlags int32) {
	p.eventFlags = eventFlags
}

// Index returns the index of this event: its time and sequence packed in one value.
func (p *TheoPrice) Index() int64 {
	return p.index
}

// SetIndex sets the index of this event.
func (p *TheoPrice) SetIndex(index int64) {
	p.index = index
}

// Time returns the time of this event in milliseconds.
func (p *TheoPrice) Time() int64 {
	return timeutil.TimeOfTimeSequence(p.index)
}

// SetTime sets the time of this event in milliseconds; the sequence is kept.
func (p *TheoPrice) SetTime(timeMillis int64) {
	p.index = timeutil.WithTime(p.index, timeMillis)
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (p *TheoPrice) Sequence() int32 {
	return timeutil.SequenceOfTimeSequence(p.index)
}

// SetSequence sets the sequence number of this event and returns an error if it is out of range.
func (p *TheoPrice) SetSequence(sequence int32) error {
	index, err := timeutil.WithSequence(p.index, sequence)
	p.index = index
	return err
}

// Price returns the theoretical option price.
func (p *TheoPrice) Price() float64 {
	return p.price
}

// SetPrice sets the theoretical option price.
func (p *TheoPrice) SetPrice(price float64) {
	p.price = price
}

// UnderlyingPrice returns the price of the underlying the theoretical price was computed with.
func (p *TheoPrice) UnderlyingPrice() float64 {
	return p.underlyingPrice
}

// SetUnderlyingPrice sets the price of the underlying.
func (p *TheoPrice) SetUnderlyingPrice(underlyingPrice float64) {
	p.underlyingPrice = underlyingPrice
}

// Delta returns the delta of the theoretical price.
func (p *TheoPrice) Delta() float64 {
	return p.delta
}

// SetDelta sets the delta of the theoretical price.
func (p *TheoPrice) SetDelta(delta float64) {
	p.delta = delta
}

// Gamma returns the gamma of the theoretical price.
func (p *TheoPrice) Gamma() float64 {
	return p.gamma
}

// SetGamma sets the gamma of the theoretical price.
func (p *TheoPrice) SetGamma(gamma float64) {
	p.gamma = gamma
}

// Dividend returns the implied simple dividend return of the underlying.
func (p *TheoPrice) Dividend() float64 {
	return p.dividend
}

// SetDividend sets the implied simple dividend return.
func (p *TheoPrice) SetDividend(dividend float64) {
	p.dividend = dividend
}

// Interest returns the implied simple interest return of the underlying.
func (p *TheoPrice) Interest() float64 {
	return p.interest
}

// SetInterest sets the implied simple interest return.
func (p *TheoPrice) SetInterest(interest float64) {
	p.interest = interest
}

// String returns a text representation of the event.
func (p *TheoPrice) String() string {
	return "TheoPrice{" + formatutil.FormatString(p.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(p.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(p.EventFlags())) +
		", time=" + formatutil.FormatTime(p.Time()) +
		", sequence=" + strconv.FormatInt(int64(p.Sequence()), 10) +
		", price=" + formatutil.FormatFloat64(p.price) +
		", underlyingPrice=" + formatutil.FormatFloat64(p.underlyingPrice) +
		", delta=" + formatutil.FormatFloat64(p.delta) +
		", gamma=" + formatutil.FormatFloat64(p.gamma) +
		", dividend=" + formatutil.FormatFloat64(p.dividend) +
		", interest=" + formatutil.FormatFloat64(p.interest) +
		"}"
}
