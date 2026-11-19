// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package underlying provides the Underlying event: volatility and put/call volume data of an underlying symbol.
package underlying

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"math"
	"strconv"
)

const maxSequence = (1 << 22) - 1

// Underlying is a time series event with the volatility and the option volume data of an underlying symbol.
type Underlying struct {
	eventSymbol *string
	eventTime   int64
	eventFlags  int32
	index       int64

	volatility      float64
	frontVolatility float64
	backVolatility  float64
	callVolume      float64
	putVolume       float64
	putCallRatio    float64
}

// NewUnderlying creates an Underlying event for the given event symbol.
func NewUnderlying(eventSymbol string) *Underlying {
	return &Underlying{
		eventSymbol:     &eventSymbol,
		volatility:      math.NaN(),
		frontVolatility: math.NaN(),
		backVolatility:  math.NaN(),
		callVolume:      math.NaN(),
		putVolume:       math.NaN(),
		putCallRatio:    math.NaN(),
	}
}

// Type returns the event code identifying this event type.
func (u *Underlying) Type() eventcodes.EventCode {
	return eventcodes.Underlying
}

// Source is always DEFAULT for Underlying (multiple sources per symbol are not supported).
func (u *Underlying) Source() *events.IndexedEventSource {
	return events.DefaultIndexedEventSource()
}

// EventSymbol returns the symbol of this event, or nil if it is not set.
func (u *Underlying) EventSymbol() *string {
	return u.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (u *Underlying) SetEventSymbol(eventSymbol string) {
	*u.eventSymbol = eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (u *Underlying) EventTime() int64 {
	return u.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (u *Underlying) SetEventTime(eventTime int64) {
	u.eventTime = eventTime
}

// EventFlags returns the transactional event flags of this event.
func (u *Underlying) EventFlags() int32 {
	return u.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (u *Underlying) SetEventFlags(eventFlags int32) {
	u.eventFlags = eventFlags
}

// Index returns the unique per-symbol index of this event.
func (u *Underlying) Index() int64 {
	return u.index
}

// SetIndex sets the unique per-symbol index of this event.
func (u *Underlying) SetIndex(index int64) {
	u.index = index
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (u *Underlying) Sequence() int32 {
	return int32(u.index & maxSequence)
}

// SetSequence sets the sequence number of this event and returns an error if the value is out of range.
func (u *Underlying) SetSequence(sequence int32) error {
	if sequence < 0 || int64(sequence) > maxSequence {
		return fmt.Errorf("Sequence(%d) is < 0 or > MaxSequence(%d)", sequence, maxSequence)
	}
	u.index = (u.index & ^maxSequence) | int64(sequence)
	return nil
}

// Time returns the time of this event in milliseconds.
func (u *Underlying) Time() int64 {
	return ((u.index >> 32) * 1000) + ((u.index >> 22) & 0x3ff)
}

// SetTime sets the time of this event in milliseconds.
func (u *Underlying) SetTime(value int64) {
	u.index = (timeutil.GetSecondsFromTime(value) << 32) |
		(int64(timeutil.GetMillisFromTime(value)) << 22) |
		int64(u.Sequence())
}

// Volatility returns the 30-day implied volatility for options of the underlying symbol.
func (u *Underlying) Volatility() float64 {
	return u.volatility
}

// SetVolatility sets the 30-day implied volatility.
func (u *Underlying) SetVolatility(volatility float64) {
	u.volatility = volatility
}

// FrontVolatility returns the front month implied volatility for options of the underlying symbol.
func (u *Underlying) FrontVolatility() float64 {
	return u.frontVolatility
}

// SetFrontVolatility sets the front month implied volatility.
func (u *Underlying) SetFrontVolatility(frontVolatility float64) {
	u.frontVolatility = frontVolatility
}

// BackVolatility returns the back month implied volatility for options of the underlying symbol.
func (u *Underlying) BackVolatility() float64 {
	return u.backVolatility
}

// SetBackVolatility sets the back month implied volatility.
func (u *Underlying) SetBackVolatility(backVolatility float64) {
	u.backVolatility = backVolatility
}

// CallVolume returns the call options traded volume for the day.
func (u *Underlying) CallVolume() float64 {
	return u.callVolume
}

// SetCallVolume sets the call options traded volume.
func (u *Underlying) SetCallVolume(callVolume float64) {
	u.callVolume = callVolume
}

// PutVolume returns the put options traded volume for the day.
func (u *Underlying) PutVolume() float64 {
	return u.putVolume
}

// SetPutVolume sets the put options traded volume.
func (u *Underlying) SetPutVolume(putVolume float64) {
	u.putVolume = putVolume
}

// OptionVolume returns options traded volume for a day.
func (u *Underlying) OptionVolume() float64 {
	if math.IsNaN(u.putVolume) {
		return u.callVolume
	}
	if math.IsNaN(u.callVolume) {
		return u.putVolume
	}
	return u.putVolume + u.callVolume
}

// PutCallRatio returns the ratio of put options volume to call options volume.
func (u *Underlying) PutCallRatio() float64 {
	return u.putCallRatio
}

// SetPutCallRatio sets the put-call ratio.
func (u *Underlying) SetPutCallRatio(putCallRatio float64) {
	u.putCallRatio = putCallRatio
}

// String returns a string representation of the event.
func (u *Underlying) String() string {
	return "Underlying{" + formatutil.FormatString(u.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(u.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(u.EventFlags())) +
		", time=" + formatutil.FormatTime(u.Time()) +
		", sequence=" + strconv.FormatInt(int64(u.Sequence()), 10) +
		", volatility=" + formatutil.FormatFloat64(u.volatility) +
		", frontVolatility=" + formatutil.FormatFloat64(u.frontVolatility) +
		", backVolatility=" + formatutil.FormatFloat64(u.backVolatility) +
		", callVolume=" + formatutil.FormatFloat64(u.callVolume) +
		", putVolume=" + formatutil.FormatFloat64(u.putVolume) +
		", putCallRatio=" + formatutil.FormatFloat64(u.putCallRatio) +
		"}"
}
