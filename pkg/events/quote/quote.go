package quote

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
	"math"
	"strconv"
)

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

func NewQuote(eventSymbol string) *Quote {
	return &Quote{
		eventSymbol: &eventSymbol,
		bidPrice:    math.NaN(),
		bidSize:     math.NaN(),
		askPrice:    math.NaN(),
		askSize:     math.NaN(),
	}
}

func (q *Quote) Type() eventcodes.EventCode {
	return eventcodes.Quote
}

func (q *Quote) EventSymbol() *string {
	return q.eventSymbol
}

func (q *Quote) SetEventSymbol(eventSymbol string) {
	*q.eventSymbol = eventSymbol
}

func (q *Quote) EventTime() int64 {
	return q.eventTime
}

func (q *Quote) SetEventTime(eventTime int64) {
	q.eventTime = eventTime
}

func (q *Quote) Sequence() int32 {
	return q.timeMillisSequence & maxSequence
}

func (q *Quote) SetSequence(sequence int32) {
	q.timeMillisSequence = (q.timeMillisSequence & ^maxSequence) | sequence
}

func (q *Quote) Time() int64 {
	return mathutil.FloorDivInt(mathutil.MaxInt(q.bidTime, q.askTime), 1000)*1000 + int64(uint32(q.timeMillisSequence)>>22)
}

func (q *Quote) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(q.Time(), q.timeNanoPart)
}

func (q *Quote) TimeNanoPart() int32 {
	return q.timeNanoPart
}

func (q *Quote) SetTimeNanoPart(timeNanoPart int32) {
	q.timeNanoPart = timeNanoPart
}

func (q *Quote) BidTime() int64 {
	return q.bidTime
}

func (q *Quote) SetBidTime(bidTime int64) {
	q.bidTime = bidTime
	q.recomputeTimeMillisPart()
}

func (q *Quote) BidExchangeCode() rune {
	return q.bidExchangeCode
}

func (q *Quote) SetBidExchangeCode(bidExchangeCode rune) {
	q.bidExchangeCode = bidExchangeCode
}

func (q *Quote) BidPrice() float64 {
	return q.bidPrice
}

func (q *Quote) SetBidPrice(bidPrice float64) {
	q.bidPrice = bidPrice
}

func (q *Quote) BidSize() float64 {
	return q.bidSize
}

func (q *Quote) SetBidSize(bidSize float64) {
	q.bidSize = bidSize
}

func (q *Quote) AskTime() int64 {
	return q.askTime
}

func (q *Quote) SetAskTime(askTime int64) {
	q.askTime = askTime
	q.recomputeTimeMillisPart()
}

func (q *Quote) AskExchangeCode() rune {
	return q.askExchangeCode
}

func (q *Quote) SetAskExchangeCode(askExchangeCode rune) {
	q.askExchangeCode = askExchangeCode
}

func (q *Quote) AskPrice() float64 {
	return q.askPrice
}

func (q *Quote) SetAskPrice(askPrice float64) {
	q.askPrice = askPrice
}

func (q *Quote) AskSize() float64 {
	return q.askSize
}

func (q *Quote) SetAskSize(askSize float64) {
	q.askSize = askSize
}

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

func (q *Quote) GetTimeMillisSequence() int32 {
	return q.timeMillisSequence
}

func (q *Quote) SetTimeMillisSequence(timeMillisSequence int32) {
	q.timeMillisSequence = timeMillisSequence
}

func (q *Quote) recomputeTimeMillisPart() {
	q.timeMillisSequence = timeutil.GetMillisFromTime(mathutil.MaxInt(q.askTime, q.bidTime))<<22 | q.Sequence()
}
