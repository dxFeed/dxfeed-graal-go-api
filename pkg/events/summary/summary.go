package summary

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/market"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
	"math"
	"strconv"
)

const (
	dayClosePriceTypeMask      = 3
	dayClosePriceTypeShift     = 2
	prevDayClosePriceTypeMask  = 3
	prevDayClosePriceTypeShift = 0
)

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

func (s *Summary) Type() eventcodes.EventCode {
	return eventcodes.Summary
}

func (s *Summary) DayClosePriceType() market.PriceType {
	return market.PriceTypeValueOf(mathutil.GetBits(int64(s.flags), dayClosePriceTypeMask, dayClosePriceTypeShift))
}

func (s *Summary) SetDayClosePriceType(t market.PriceType) {
	s.SetFlags(int32(mathutil.SetBits(int64(s.flags), dayClosePriceTypeMask, dayClosePriceTypeShift, int64(t.Code()))))
}

func (s *Summary) PrevDayClosePriceType() market.PriceType {
	return market.PriceTypeValueOf(mathutil.GetBits(int64(s.flags), prevDayClosePriceTypeMask, prevDayClosePriceTypeShift))
}

func (s *Summary) SetPrevDayClosePriceType(t market.PriceType) {
	s.SetFlags(int32(mathutil.SetBits(int64(s.flags), prevDayClosePriceTypeMask, prevDayClosePriceTypeShift, int64(t.Code()))))
}

func (s *Summary) EventSymbol() *string {
	return s.eventSymbol
}

func (s *Summary) SetEventSymbol(eventSymbol string) {
	*s.eventSymbol = eventSymbol
}

func (s *Summary) EventTime() int64 {
	return s.eventTime
}

func (s *Summary) SetEventTime(eventTime int64) {
	s.eventTime = eventTime
}

func (s *Summary) DayId() int32 {
	return s.dayId
}

func (s *Summary) SetDayId(dayId int32) {
	s.dayId = dayId
}

func (s *Summary) DayOpenPrice() float64 {
	return s.dayOpenPrice
}

func (s *Summary) SetDayOpenPrice(v float64) {
	s.dayOpenPrice = v
}

func (s *Summary) DayHighPrice() float64 {
	return s.dayHighPrice
}

func (s *Summary) SetDayHighPrice(v float64) {
	s.dayHighPrice = v
}

func (s *Summary) DayLowPrice() float64 {
	return s.dayLowPrice
}

func (s *Summary) SetDayLowPrice(v float64) {
	s.dayLowPrice = v
}

func (s *Summary) DayClosePrice() float64 {
	return s.dayClosePrice
}

func (s *Summary) SetDayClosePrice(v float64) {
	s.dayClosePrice = v
}

func (s *Summary) PrevDayId() int32 {
	return s.prevDayId
}

func (s *Summary) SetPrevDayId(prevDayId int32) {
	s.prevDayId = prevDayId
}

func (s *Summary) PrevDayClosePrice() float64 {
	return s.prevDayClosePrice
}

func (s *Summary) SetPrevDayClosePrice(v float64) {
	s.prevDayClosePrice = v
}

func (s *Summary) PrevDayVolume() float64 {
	return s.prevDayVolume
}

func (s *Summary) SetPrevDayVolume(v float64) {
	s.prevDayVolume = v
}

func (s *Summary) OpenInterest() int64 {
	return s.openInterest
}

func (s *Summary) SetOpenInterest(openInterest int64) {
	s.openInterest = openInterest
}

func (s *Summary) Flags() int32 {
	return s.flags
}

func (s *Summary) SetFlags(flags int32) {
	s.flags = flags
}

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
