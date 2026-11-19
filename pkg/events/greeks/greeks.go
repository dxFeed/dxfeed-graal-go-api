// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package greeks provides the Greeks event: option price, volatility and sensitivities (delta, gamma, theta, rho, vega).
package greeks

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

// Greeks is a time series event with the option price, Black-Scholes volatility and the greeks (sensitivities) of an option.
type Greeks struct {
	eventSymbol *string
	eventTime   int64
	eventFlags  int32
	index       int64

	price      float64
	volatility float64
	delta      float64
	gamma      float64
	theta      float64
	rho        float64
	vega       float64
}

// NewGreeks creates a Greeks event for the given event symbol.
func NewGreeks(eventSymbol string) *Greeks {
	return &Greeks{
		eventSymbol: &eventSymbol,
		price:       math.NaN(),
		volatility:  math.NaN(),
		delta:       math.NaN(),
		gamma:       math.NaN(),
		theta:       math.NaN(),
		rho:         math.NaN(),
		vega:        math.NaN(),
	}
}

// Type returns the event code identifying this event type.
func (g *Greeks) Type() eventcodes.EventCode {
	return eventcodes.Greeks
}

// Source is always DEFAULT for Greeks (multiple sources per symbol are not supported).
func (g *Greeks) Source() *events.IndexedEventSource {
	return events.DefaultIndexedEventSource()
}

// EventSymbol returns the symbol of this event, or nil if it is not set.
func (g *Greeks) EventSymbol() *string {
	return g.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (g *Greeks) SetEventSymbol(eventSymbol string) {
	*g.eventSymbol = eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (g *Greeks) EventTime() int64 {
	return g.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (g *Greeks) SetEventTime(eventTime int64) {
	g.eventTime = eventTime
}

// EventFlags returns the transactional event flags of this event.
func (g *Greeks) EventFlags() int32 {
	return g.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (g *Greeks) SetEventFlags(eventFlags int32) {
	g.eventFlags = eventFlags
}

// Index returns the unique per-symbol index of this event.
func (g *Greeks) Index() int64 {
	return g.index
}

// SetIndex sets the unique per-symbol index of this event.
func (g *Greeks) SetIndex(index int64) {
	g.index = index
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (g *Greeks) Sequence() int32 {
	return int32(g.index & maxSequence)
}

// SetSequence sets the sequence number of this event and returns an error if the value is out of range.
func (g *Greeks) SetSequence(sequence int32) error {
	if sequence < 0 || int64(sequence) > maxSequence {
		return fmt.Errorf("Sequence(%d) is < 0 or > MaxSequence(%d)", sequence, maxSequence)
	}
	g.index = (g.index & ^maxSequence) | int64(sequence)
	return nil
}

// Time returns the time of this event in milliseconds.
func (g *Greeks) Time() int64 {
	return ((g.index >> 32) * 1000) + ((g.index >> 22) & 0x3ff)
}

// SetTime sets the time of this event in milliseconds.
func (g *Greeks) SetTime(value int64) {
	g.index = (timeutil.GetSecondsFromTime(value) << 32) |
		(int64(timeutil.GetMillisFromTime(value)) << 22) |
		int64(g.Sequence())
}

// Price returns the option market price.
func (g *Greeks) Price() float64 {
	return g.price
}

// SetPrice sets the option market price.
func (g *Greeks) SetPrice(price float64) {
	g.price = price
}

// Volatility returns the Black-Scholes implied volatility of the option.
func (g *Greeks) Volatility() float64 {
	return g.volatility
}

// SetVolatility sets the Black-Scholes implied volatility of the option.
func (g *Greeks) SetVolatility(volatility float64) {
	g.volatility = volatility
}

// Delta returns the option delta: the change of the option price per unit change of the underlying price.
func (g *Greeks) Delta() float64 {
	return g.delta
}

// SetDelta sets the option delta.
func (g *Greeks) SetDelta(delta float64) {
	g.delta = delta
}

// Gamma returns the option gamma: the change of delta per unit change of the underlying price.
func (g *Greeks) Gamma() float64 {
	return g.gamma
}

// SetGamma sets the option gamma.
func (g *Greeks) SetGamma(gamma float64) {
	g.gamma = gamma
}

// Theta returns the option theta: the change of the option price per unit of time decay.
func (g *Greeks) Theta() float64 {
	return g.theta
}

// SetTheta sets the option theta.
func (g *Greeks) SetTheta(theta float64) {
	g.theta = theta
}

// Rho returns the option rho: the change of the option price per unit change of the interest rate.
func (g *Greeks) Rho() float64 {
	return g.rho
}

// SetRho sets the option rho.
func (g *Greeks) SetRho(rho float64) {
	g.rho = rho
}

// Vega returns the option vega: the change of the option price per unit change of the volatility.
func (g *Greeks) Vega() float64 {
	return g.vega
}

// SetVega sets the option vega.
func (g *Greeks) SetVega(vega float64) {
	g.vega = vega
}

// String returns a string representation of the event.
func (g *Greeks) String() string {
	return "Greeks{" + formatutil.FormatString(g.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(g.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(g.EventFlags())) +
		", time=" + formatutil.FormatTime(g.Time()) +
		", sequence=" + strconv.FormatInt(int64(g.Sequence()), 10) +
		", price=" + formatutil.FormatFloat64(g.price) +
		", volatility=" + formatutil.FormatFloat64(g.volatility) +
		", delta=" + formatutil.FormatFloat64(g.delta) +
		", gamma=" + formatutil.FormatFloat64(g.gamma) +
		", theta=" + formatutil.FormatFloat64(g.theta) +
		", rho=" + formatutil.FormatFloat64(g.rho) +
		", vega=" + formatutil.FormatFloat64(g.vega) +
		"}"
}
