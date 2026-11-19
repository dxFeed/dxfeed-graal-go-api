// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package quote provides the Quote event: the best bid and ask prices and sizes of a symbol.
package quote

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"math"
	"strconv"
)

// Quote is a quote event with the current best bid and ask prices and sizes for a symbol.
type Quote struct {
	eventSymbol        *string
	eventTime          int64
	timeMillisSequence int32
	timeNanoPart       int32
	bidTime            int64
	bidExchangeCode    rune
	bidPrice           float64
	bidSize            float64
	askTime            int64
	askExchangeCode    rune
	askPrice           float64
	askSize            float64
}

const maxSequence = (1 << 22) - 1

// NewQuote creates a Quote event for the given event symbol.
func NewQuote(eventSymbol string) *Quote {
	return &Quote{
		eventSymbol: &eventSymbol,
		bidPrice:    math.NaN(),
		bidSize:     math.NaN(),
		askPrice:    math.NaN(),
		askSize:     math.NaN(),
	}
}

// Type returns the event code identifying this event type.
func (q *Quote) Type() eventcodes.EventCode {
	return eventcodes.Quote
}

// EventSymbol returns the symbol of this event, or nil if it is not set.
func (q *Quote) EventSymbol() *string {
	return q.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (q *Quote) SetEventSymbol(eventSymbol string) {
	*q.eventSymbol = eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (q *Quote) EventTime() int64 {
	return q.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (q *Quote) SetEventTime(eventTime int64) {
	q.eventTime = eventTime
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (q *Quote) Sequence() int32 {
	return q.timeMillisSequence & maxSequence
}

// SetSequence sets the sequence number of this event; the value is not validated.
func (q *Quote) SetSequence(sequence int32) {
	q.timeMillisSequence = (q.timeMillisSequence & ^maxSequence) | sequence
}

// Time returns the time of the last bid or ask change in milliseconds, the later of BidTime and AskTime.
func (q *Quote) Time() int64 {
	return mathutil.FloorDivInt(mathutil.MaxInt(q.bidTime, q.askTime), 1000)*1000 + int64(uint32(q.timeMillisSequence)>>22)
}

// TimeNanos returns the time of the last bid or ask change in nanoseconds.
func (q *Quote) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(q.Time(), q.timeNanoPart)
}

// TimeNanoPart returns the microsecond and nanosecond part of the time of the last bid or ask change.
func (q *Quote) TimeNanoPart() int32 {
	return q.timeNanoPart
}

// SetTimeNanoPart sets the microsecond and nanosecond part of the time of the last bid or ask change.
func (q *Quote) SetTimeNanoPart(timeNanoPart int32) {
	q.timeNanoPart = timeNanoPart
}

// BidTime returns the time of the last bid change in milliseconds.
func (q *Quote) BidTime() int64 {
	return q.bidTime
}

// SetBidTime sets the time of the last bid change in milliseconds.
func (q *Quote) SetBidTime(bidTime int64) {
	q.bidTime = bidTime
	q.recomputeTimeMillisPart()
}

// BidExchangeCode returns the exchange code of the bid.
func (q *Quote) BidExchangeCode() rune {
	return q.bidExchangeCode
}

// SetBidExchangeCode sets the exchange code of the bid.
func (q *Quote) SetBidExchangeCode(bidExchangeCode rune) {
	q.bidExchangeCode = bidExchangeCode
}

// BidPrice returns the bid price.
func (q *Quote) BidPrice() float64 {
	return q.bidPrice
}

// SetBidPrice sets the bid price.
func (q *Quote) SetBidPrice(bidPrice float64) {
	q.bidPrice = bidPrice
}

// BidSize returns the bid size.
func (q *Quote) BidSize() float64 {
	return q.bidSize
}

// SetBidSize sets the bid size.
func (q *Quote) SetBidSize(bidSize float64) {
	q.bidSize = bidSize
}

// AskTime returns the time of the last ask change in milliseconds.
func (q *Quote) AskTime() int64 {
	return q.askTime
}

// SetAskTime sets the time of the last ask change in milliseconds.
func (q *Quote) SetAskTime(askTime int64) {
	q.askTime = askTime
	q.recomputeTimeMillisPart()
}

// AskExchangeCode returns the exchange code of the ask.
func (q *Quote) AskExchangeCode() rune {
	return q.askExchangeCode
}

// SetAskExchangeCode sets the exchange code of the ask.
func (q *Quote) SetAskExchangeCode(askExchangeCode rune) {
	q.askExchangeCode = askExchangeCode
}

// AskPrice returns the ask price.
func (q *Quote) AskPrice() float64 {
	return q.askPrice
}

// SetAskPrice sets the ask price.
func (q *Quote) SetAskPrice(askPrice float64) {
	q.askPrice = askPrice
}

// AskSize returns the ask size.
func (q *Quote) AskSize() float64 {
	return q.askSize
}

// SetAskSize sets the ask size.
func (q *Quote) SetAskSize(askSize float64) {
	q.askSize = askSize
}

// String returns a string representation of the event.
func (q *Quote) String() string {
	return "Quote{" + formatutil.FormatString(q.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(q.EventTime()) +
		", time=" + formatutil.FormatTime(q.Time()) +
		", timeNanoPart=" + strconv.FormatInt(int64(q.timeNanoPart), 10) +
		", sequence=" + strconv.FormatInt(int64(q.Sequence()), 10) +
		", bidTime=" + formatutil.FormatTime(q.bidTime) +
		", bidExchange=" + formatutil.FormatChar(q.bidExchangeCode) +
		", bidPrice=" + formatutil.FormatFloat64(q.bidPrice) +
		", bidSize=" + formatutil.FormatFloat64(q.bidSize) +
		", askTime=" + formatutil.FormatTime(q.askTime) +
		", askExchange=" + formatutil.FormatChar(q.askExchangeCode) +
		", askPrice=" + formatutil.FormatFloat64(q.askPrice) +
		", askSize=" + formatutil.FormatFloat64(q.askSize) +
		"}"
}

// TimeMillisSequence returns the time in milliseconds part combined with the sequence number, as a packed value.
func (q *Quote) TimeMillisSequence() int32 {
	return q.timeMillisSequence
}

// SetTimeMillisSequence sets the packed time in milliseconds part and sequence number.
func (q *Quote) SetTimeMillisSequence(timeMillisSequence int32) {
	q.timeMillisSequence = timeMillisSequence
}

func (q *Quote) recomputeTimeMillisPart() {
	q.timeMillisSequence = timeutil.GetMillisFromTime(mathutil.MaxInt(q.askTime, q.bidTime))<<22 | q.Sequence()
}
