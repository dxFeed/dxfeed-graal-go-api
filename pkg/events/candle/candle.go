package candle

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
	"math"
	"strconv"
)

const maxSequence = (1 << 22) - 1

type Candle struct {
	eventSymbol   *CandleSymbol
	eventTime     int64
	count         int64
	open          float64
	high          float64
	low           float64
	close         float64
	volume        float64
	vwap          float64
	bidVolume     float64
	askVolume     float64
	impVolatility float64
	openInterest  float64
	eventFlags    int32
	index         int64
}

func NewCandle(eventSymbol string) *Candle {
	return &Candle{
		eventSymbol:   NewCandleSymbol(eventSymbol),
		open:          math.NaN(),
		high:          math.NaN(),
		low:           math.NaN(),
		close:         math.NaN(),
		volume:        math.NaN(),
		vwap:          math.NaN(),
		bidVolume:     math.NaN(),
		askVolume:     math.NaN(),
		impVolatility: math.NaN(),
		openInterest:  math.NaN(),
	}
}

func (c *Candle) EventSymbol() *CandleSymbol {
	return c.eventSymbol
}

func (c *Candle) SetEventSymbol(eventSymbol *CandleSymbol) {
	c.eventSymbol = eventSymbol
}

func (c *Candle) EventTime() int64 {
	return c.eventTime
}

func (c *Candle) SetEventTime(eventTime int64) {
	c.eventTime = eventTime
}

func (c *Candle) Count() int64 {
	return c.count
}

func (c *Candle) SetCount(count int64) {
	c.count = count
}

func (c *Candle) Open() float64 {
	return c.open
}

func (c *Candle) SetOpen(open float64) {
	c.open = open
}

func (c *Candle) High() float64 {
	return c.high
}

func (c *Candle) SetHigh(high float64) {
	c.high = high
}

func (c *Candle) Low() float64 {
	return c.low
}

func (c *Candle) SetLow(low float64) {
	c.low = low
}

func (c *Candle) Close() float64 {
	return c.close
}

func (c *Candle) SetClose(close float64) {
	c.close = close
}

func (c *Candle) Volume() float64 {
	return c.volume
}

func (c *Candle) SetVolume(volume float64) {
	c.volume = volume
}

func (c *Candle) Vwap() float64 {
	return c.vwap
}

func (c *Candle) SetVwap(vwap float64) {
	c.vwap = vwap
}

func (c *Candle) BidVolume() float64 {
	return c.bidVolume
}

func (c *Candle) SetBidVolume(bidVolume float64) {
	c.bidVolume = bidVolume
}

func (c *Candle) AskVolume() float64 {
	return c.askVolume
}

func (c *Candle) SetAskVolume(askVolume float64) {
	c.askVolume = askVolume
}

func (c *Candle) ImpVolatility() float64 {
	return c.impVolatility
}

func (c *Candle) SetImpVolatility(impVolatility float64) {
	c.impVolatility = impVolatility
}

func (c *Candle) OpenInterest() float64 {
	return c.openInterest
}

func (c *Candle) SetOpenInterest(openInterest float64) {
	c.openInterest = openInterest
}

func (c *Candle) EventFlags() int32 {
	return c.eventFlags
}

func (c *Candle) SetEventFlags(eventFlags int32) {
	c.eventFlags = eventFlags
}

func (c *Candle) Index() int64 {
	return c.index
}

func (c *Candle) SetIndex(index int64) {
	c.index = index
}

func (c *Candle) Sequence() int64 {
	return c.index & maxSequence
}

func (q *Candle) SetSequence(sequence int64) error {
	if sequence < 0 || sequence > maxSequence {
		return fmt.Errorf("sequence(%d) is < 0 or > MaxSequence(%d)", sequence, maxSequence)
	}
	q.index = (q.index & ^maxSequence) | sequence
	return nil
}

func (q *Candle) Time() int64 {
	return ((q.index >> 32) * 1000) + ((q.index >> 22) & 0x3ff)
}

func (q *Candle) SetTime(value int64) {
	q.index = (int64(timeutil.GetSecondsFromTime(value)) << 32) |
		int64(timeutil.GetMillisFromTime(value)<<22) |
		q.Sequence()
}

func (c *Candle) String() string {
	return "Candle{" + c.EventSymbol().String() +
		", eventTime=" + formatutil.FormatTime(c.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(c.EventFlags())) +
		", time=" + formatutil.FormatTime(c.Time()) +
		", sequence=" + strconv.FormatInt(c.Sequence(), 10) +
		", count=" + strconv.FormatInt(c.Count(), 10) +
		", open=" + formatutil.FormatFloat64(c.Open()) +
		", high=" + formatutil.FormatFloat64(c.High()) +
		", low=" + formatutil.FormatFloat64(c.Low()) +
		", close=" + formatutil.FormatFloat64(c.Close()) +
		", volume=" + formatutil.FormatFloat64(c.Volume()) +
		", vwap=" + formatutil.FormatFloat64(c.Vwap()) +
		", bidVolume=" + formatutil.FormatFloat64(c.BidVolume()) +
		", askVolume=" + formatutil.FormatFloat64(c.AskVolume()) +
		", impVolatility=" + formatutil.FormatFloat64(c.ImpVolatility()) +
		", openInterest=" + formatutil.FormatFloat64(c.OpenInterest()) +
		"}"
}

func (c *Candle) Type() eventcodes.EventCode {
	return eventcodes.Candle
}
