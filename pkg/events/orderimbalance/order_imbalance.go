// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package orderimbalance provides the OrderImbalance event: the imbalance of the orders of an auction.
package orderimbalance

import (
	"fmt"
	"math"
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

// The flags hold the imbalance side in bits 0 to 2 and the auction type in bits 3 to 5.
const (
	sideMask         = 7
	sideShift        = 0
	auctionTypeMask  = 7
	auctionTypeShift = 3
)

// ImbalanceSide is the side of an order imbalance.
type ImbalanceSide int32

const (
	// ImbalanceSideUndefined is an undefined side.
	ImbalanceSideUndefined ImbalanceSide = iota
	// ImbalanceSideBuy is a buy side imbalance.
	ImbalanceSideBuy
	// ImbalanceSideSell is a sell side imbalance.
	ImbalanceSideSell
	// ImbalanceSideNoImbalance means that there is no imbalance.
	ImbalanceSideNoImbalance
	// ImbalanceSideInsufficientOrders means that there are not enough orders to calculate the imbalance.
	ImbalanceSideInsufficientOrders
	// ImbalanceSidePaused means that the auction is paused.
	ImbalanceSidePaused
)

// String returns the name of the side, as in the Java API.
func (s ImbalanceSide) String() string {
	switch s {
	case ImbalanceSideUndefined:
		return "UNDEFINED"
	case ImbalanceSideBuy:
		return "BUY"
	case ImbalanceSideSell:
		return "SELL"
	case ImbalanceSideNoImbalance:
		return "NO_IMBALANCE"
	case ImbalanceSideInsufficientOrders:
		return "INSUFFICIENT_ORDERS"
	case ImbalanceSidePaused:
		return "PAUSED"
	}
	return fmt.Sprintf("ImbalanceSide(%d)", int32(s))
}

// AuctionType is the type of an auction.
type AuctionType int32

const (
	// AuctionTypeUndefined is an undefined auction type.
	AuctionTypeUndefined AuctionType = iota
	// AuctionTypeOpening is the opening auction.
	AuctionTypeOpening
	// AuctionTypeClosing is the closing auction.
	AuctionTypeClosing
	// AuctionTypeReopening is an auction that reopens the trading after a halt.
	AuctionTypeReopening
	// AuctionTypeOther is another auction.
	AuctionTypeOther
)

// String returns the name of the auction type, as in the Java API.
func (t AuctionType) String() string {
	switch t {
	case AuctionTypeUndefined:
		return "UNDEFINED"
	case AuctionTypeOpening:
		return "OPENING"
	case AuctionTypeClosing:
		return "CLOSING"
	case AuctionTypeReopening:
		return "REOPENING"
	case AuctionTypeOther:
		return "OTHER"
	}
	return fmt.Sprintf("AuctionType(%d)", int32(t))
}

// OrderImbalance is the imbalance of the buy and sell orders of an auction: the reference price, the paired and the
// imbalance sizes and the indicative prices. It is a market event, neither lasting nor indexed; its source is an
// order source.
type OrderImbalance struct {
	eventSymbol   *string
	eventTime     int64
	sourceId      int32
	timeSequence  int64
	refPrice      float64
	pairedSize    float64
	imbalanceSize float64
	nearPrice     float64
	farPrice      float64
	flags         int32
}

// NewOrderImbalance creates an OrderImbalance for the given event symbol.
func NewOrderImbalance(eventSymbol string) *OrderImbalance {
	return &OrderImbalance{
		eventSymbol:   &eventSymbol,
		refPrice:      math.NaN(),
		pairedSize:    math.NaN(),
		imbalanceSize: math.NaN(),
		nearPrice:     math.NaN(),
		farPrice:      math.NaN(),
	}
}

// Type returns the event code of this event type.
func (o *OrderImbalance) Type() eventcodes.EventCode {
	return eventcodes.OrderImbalance
}

// EventSymbol returns the symbol of this event.
func (o *OrderImbalance) EventSymbol() *string {
	return o.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (o *OrderImbalance) SetEventSymbol(eventSymbol string) {
	o.eventSymbol = &eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (o *OrderImbalance) EventTime() int64 {
	return o.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (o *OrderImbalance) SetEventTime(eventTime int64) {
	o.eventTime = eventTime
}

// SourceId returns the identifier of the order source of this event.
func (o *OrderImbalance) SourceId() int32 {
	return o.sourceId
}

// SetSourceId sets the identifier of the order source of this event.
func (o *OrderImbalance) SetSourceId(sourceId int32) {
	o.sourceId = sourceId
}

// Source returns the order source of this event, or an error for an unknown identifier.
func (o *OrderImbalance) Source() (*order.Source, error) {
	return order.ValueOfIdentifier(int64(o.sourceId))
}

// SetSource sets the order source of this event.
func (o *OrderImbalance) SetSource(source *order.Source) {
	o.sourceId = int32(source.Id())
}

// TimeSequence returns the time and the sequence of this event packed in one value.
func (o *OrderImbalance) TimeSequence() int64 {
	return o.timeSequence
}

// SetTimeSequence sets the time and the sequence of this event packed in one value.
func (o *OrderImbalance) SetTimeSequence(timeSequence int64) {
	o.timeSequence = timeSequence
}

// Time returns the time of the imbalance in milliseconds since the Unix epoch.
func (o *OrderImbalance) Time() int64 {
	return timeutil.TimeOfTimeSequence(o.timeSequence)
}

// SetTime sets the time of the imbalance in milliseconds since the Unix epoch; the sequence is kept.
func (o *OrderImbalance) SetTime(timeMillis int64) {
	o.timeSequence = timeutil.WithTime(o.timeSequence, timeMillis)
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (o *OrderImbalance) Sequence() int32 {
	return timeutil.SequenceOfTimeSequence(o.timeSequence)
}

// SetSequence sets the sequence number of this event; it returns an error if the sequence is out of range.
func (o *OrderImbalance) SetSequence(sequence int32) error {
	timeSequence, err := timeutil.WithSequence(o.timeSequence, sequence)
	if err != nil {
		return err
	}
	o.timeSequence = timeSequence
	return nil
}

// RefPrice returns the reference price: the price for which the paired and the imbalance sizes are calculated.
func (o *OrderImbalance) RefPrice() float64 {
	return o.refPrice
}

// SetRefPrice sets the reference price.
func (o *OrderImbalance) SetRefPrice(refPrice float64) {
	o.refPrice = refPrice
}

// PairedSize returns the number of shares that can be matched at the reference price.
func (o *OrderImbalance) PairedSize() float64 {
	return o.pairedSize
}

// SetPairedSize sets the number of shares that can be matched at the reference price.
func (o *OrderImbalance) SetPairedSize(pairedSize float64) {
	o.pairedSize = pairedSize
}

// ImbalanceSize returns the number of shares that are not matched at the reference price.
func (o *OrderImbalance) ImbalanceSize() float64 {
	return o.imbalanceSize
}

// SetImbalanceSize sets the number of shares that are not matched at the reference price.
func (o *OrderImbalance) SetImbalanceSize(imbalanceSize float64) {
	o.imbalanceSize = imbalanceSize
}

// NearPrice returns the hypothetical price at which all the interest on the book could trade now.
func (o *OrderImbalance) NearPrice() float64 {
	return o.nearPrice
}

// SetNearPrice sets the near price.
func (o *OrderImbalance) SetNearPrice(nearPrice float64) {
	o.nearPrice = nearPrice
}

// FarPrice returns the hypothetical price at which only the auction interest could trade now.
func (o *OrderImbalance) FarPrice() float64 {
	return o.farPrice
}

// SetFarPrice sets the far price.
func (o *OrderImbalance) SetFarPrice(farPrice float64) {
	o.farPrice = farPrice
}

// Flags returns the packed flags: the imbalance side and the auction type.
func (o *OrderImbalance) Flags() int32 {
	return o.flags
}

// SetFlags sets the packed flags.
func (o *OrderImbalance) SetFlags(flags int32) {
	o.flags = flags
}

// ImbalanceSide returns the side of the imbalance.
func (o *OrderImbalance) ImbalanceSide() ImbalanceSide {
	return ImbalanceSide(mathutil.GetBits(int64(o.flags), sideMask, sideShift))
}

// SetImbalanceSide sets the side of the imbalance.
func (o *OrderImbalance) SetImbalanceSide(value ImbalanceSide) {
	o.flags = int32(mathutil.SetBits(int64(o.flags), sideMask, sideShift, int64(value)))
}

// AuctionType returns the type of the auction.
func (o *OrderImbalance) AuctionType() AuctionType {
	return AuctionType(mathutil.GetBits(int64(o.flags), auctionTypeMask, auctionTypeShift))
}

// SetAuctionType sets the type of the auction.
func (o *OrderImbalance) SetAuctionType(value AuctionType) {
	o.flags = int32(mathutil.SetBits(int64(o.flags), auctionTypeMask, auctionTypeShift, int64(value)))
}

// String returns a text representation of the event.
func (o *OrderImbalance) String() string {
	source := strconv.FormatInt(int64(o.sourceId), 10)
	if s, err := o.Source(); err == nil && s != nil && s.Name() != nil {
		source = *s.Name()
	}
	return "OrderImbalance{" + formatutil.FormatString(o.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(o.EventTime()) +
		", source=" + source +
		", time=" + formatutil.FormatTime(o.Time()) +
		", sequence=" + strconv.FormatInt(int64(o.Sequence()), 10) +
		", refPrice=" + formatutil.FormatFloat64(o.refPrice) +
		", pairedSize=" + formatutil.FormatFloat64(o.pairedSize) +
		", imbalanceSize=" + formatutil.FormatFloat64(o.imbalanceSize) +
		", nearPrice=" + formatutil.FormatFloat64(o.nearPrice) +
		", farPrice=" + formatutil.FormatFloat64(o.farPrice) +
		", imbalanceSide=" + o.ImbalanceSide().String() +
		", auctionType=" + o.AuctionType().String() +
		"}"
}
