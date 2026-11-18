// Package summary provides the Summary event: a snapshot of the trading session summary of a symbol (day and previous day prices, open interest).
package summary

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/market"
	"math"
	"strconv"
)

const (
	dayClosePriceTypeMask      = 3
	dayClosePriceTypeShift     = 2
	prevDayClosePriceTypeMask  = 3
	prevDayClosePriceTypeShift = 0
)

// Summary is a summary information snapshot about the trading session of a symbol, including day and previous day prices and open interest.
type Summary struct {
	eventSymbol *string
	eventTime   int64

	dayId             int32
	dayOpenPrice      float64
	dayHighPrice      float64
	dayLowPrice       float64
	dayClosePrice     float64
	prevDayId         int32
	prevDayClosePrice float64
	prevDayVolume     float64
	openInterest      int64
	flags             int32
}

// NewSummary creates a Summary event for the given event symbol.
func NewSummary(eventSymbol string) *Summary {
	return &Summary{
		eventSymbol:       &eventSymbol,
		dayOpenPrice:      math.NaN(),
		dayHighPrice:      math.NaN(),
		dayLowPrice:       math.NaN(),
		dayClosePrice:     math.NaN(),
		prevDayClosePrice: math.NaN(),
		prevDayVolume:     math.NaN(),
	}
}

// Type returns the event code identifying this event type.
func (s *Summary) Type() eventcodes.EventCode {
	return eventcodes.Summary
}

// DayClosePriceType returns the price type of the day close price.
func (s *Summary) DayClosePriceType() market.PriceType {
	return market.PriceTypeValueOf(mathutil.GetBits(int64(s.flags), dayClosePriceTypeMask, dayClosePriceTypeShift))
}

// SetDayClosePriceType sets the price type of the day close price.
func (s *Summary) SetDayClosePriceType(t market.PriceType) {
	s.SetFlags(int32(mathutil.SetBits(int64(s.flags), dayClosePriceTypeMask, dayClosePriceTypeShift, int64(t.Code()))))
}

// PrevDayClosePriceType returns the price type of the previous day close price.
func (s *Summary) PrevDayClosePriceType() market.PriceType {
	return market.PriceTypeValueOf(mathutil.GetBits(int64(s.flags), prevDayClosePriceTypeMask, prevDayClosePriceTypeShift))
}

// SetPrevDayClosePriceType sets the price type of the previous day close price.
func (s *Summary) SetPrevDayClosePriceType(t market.PriceType) {
	s.SetFlags(int32(mathutil.SetBits(int64(s.flags), prevDayClosePriceTypeMask, prevDayClosePriceTypeShift, int64(t.Code()))))
}

// EventSymbol returns the symbol of this event, or nil if it is not set.
func (s *Summary) EventSymbol() *string {
	return s.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (s *Summary) SetEventSymbol(eventSymbol string) {
	*s.eventSymbol = eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (s *Summary) EventTime() int64 {
	return s.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (s *Summary) SetEventTime(eventTime int64) {
	s.eventTime = eventTime
}

// DayId returns the identifier of the day that this summary represents.
func (s *Summary) DayId() int32 {
	return s.dayId
}

// SetDayId sets the identifier of the day that this summary represents.
func (s *Summary) SetDayId(dayId int32) {
	s.dayId = dayId
}

// DayOpenPrice returns the first (open) price of the day.
func (s *Summary) DayOpenPrice() float64 {
	return s.dayOpenPrice
}

// SetDayOpenPrice sets the first (open) price of the day.
func (s *Summary) SetDayOpenPrice(v float64) {
	s.dayOpenPrice = v
}

// DayHighPrice returns the maximal (high) price of the day.
func (s *Summary) DayHighPrice() float64 {
	return s.dayHighPrice
}

// SetDayHighPrice sets the maximal (high) price of the day.
func (s *Summary) SetDayHighPrice(v float64) {
	s.dayHighPrice = v
}

// DayLowPrice returns the minimal (low) price of the day.
func (s *Summary) DayLowPrice() float64 {
	return s.dayLowPrice
}

// SetDayLowPrice sets the minimal (low) price of the day.
func (s *Summary) SetDayLowPrice(v float64) {
	s.dayLowPrice = v
}

// DayClosePrice returns the last (close) price of the day.
func (s *Summary) DayClosePrice() float64 {
	return s.dayClosePrice
}

// SetDayClosePrice sets the last (close) price of the day.
func (s *Summary) SetDayClosePrice(v float64) {
	s.dayClosePrice = v
}

// PrevDayId returns the identifier of the previous day that this summary represents.
func (s *Summary) PrevDayId() int32 {
	return s.prevDayId
}

// SetPrevDayId sets the identifier of the previous day that this summary represents.
func (s *Summary) SetPrevDayId(prevDayId int32) {
	s.prevDayId = prevDayId
}

// PrevDayClosePrice returns the last (close) price of the previous day.
func (s *Summary) PrevDayClosePrice() float64 {
	return s.prevDayClosePrice
}

// SetPrevDayClosePrice sets the last (close) price of the previous day.
func (s *Summary) SetPrevDayClosePrice(v float64) {
	s.prevDayClosePrice = v
}

// PrevDayVolume returns the total volume traded on the previous day.
func (s *Summary) PrevDayVolume() float64 {
	return s.prevDayVolume
}

// SetPrevDayVolume sets the total volume traded on the previous day.
func (s *Summary) SetPrevDayVolume(v float64) {
	s.prevDayVolume = v
}

// OpenInterest returns the open interest of the symbol as the number of open contracts.
func (s *Summary) OpenInterest() int64 {
	return s.openInterest
}

// SetOpenInterest sets the open interest of the symbol.
func (s *Summary) SetOpenInterest(openInterest int64) {
	s.openInterest = openInterest
}

// Flags returns the raw flags bit mask holding the price types and other flags.
func (s *Summary) Flags() int32 {
	return s.flags
}

// SetFlags sets the raw flags bit mask.
func (s *Summary) SetFlags(flags int32) {
	s.flags = flags
}

// String returns a string representation of the event.
func (s *Summary) String() string {
	return "Summary{" + formatutil.FormatString(s.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(s.EventTime()) +
		", day=" + strconv.FormatInt(int64(timeutil.GetYearMonthDayByDayId(s.DayId())), 10) +
		", dayOpen=" + formatutil.FormatFloat64(s.dayOpenPrice) +
		", dayHigh=" + formatutil.FormatFloat64(s.dayHighPrice) +
		", dayLow=" + formatutil.FormatFloat64(s.dayLowPrice) +
		", dayClose=" + formatutil.FormatFloat64(s.dayClosePrice) +
		", dayCloseType=" + s.DayClosePriceType().String() +
		", prevDay=" + strconv.FormatInt(int64(timeutil.GetYearMonthDayByDayId(s.PrevDayId())), 10) +
		", prevDayClose=" + formatutil.FormatFloat64(s.prevDayClosePrice) +
		", prevDayCloseType=" + s.PrevDayClosePriceType().String() +
		", prevDayVolume=" + formatutil.FormatFloat64(s.prevDayVolume) +
		", openInterest=" + strconv.FormatInt(s.openInterest, 10) +
		"}"
}
