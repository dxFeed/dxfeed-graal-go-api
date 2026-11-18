// Package candle provides the Candle event and candle symbols (time-series OHLC aggregates).
package candle

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"math"
	"strconv"
)

const maxSequence = (1 << 22) - 1

// Candle is a candle event: an OHLC aggregate of market data for one candle symbol (a time-series event).
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

// NewCandle creates a Candle for the given event symbol with price fields set to NaN.
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

// EventSymbol returns the candle symbol of this event.
func (c *Candle) EventSymbol() *CandleSymbol {
	return c.eventSymbol
}

// SetEventSymbol sets the candle symbol of this event.
func (c *Candle) SetEventSymbol(eventSymbol *CandleSymbol) {
	c.eventSymbol = eventSymbol
}

// EventTime returns the time when the event was created or modified, in milliseconds since the Unix epoch.
func (c *Candle) EventTime() int64 {
	return c.eventTime
}

// SetEventTime sets the time when the event was created or modified, in milliseconds since the Unix epoch.
func (c *Candle) SetEventTime(eventTime int64) {
	c.eventTime = eventTime
}

// Count returns the total number of original trades (or ticks) in this candle.
func (c *Candle) Count() int64 {
	return c.count
}

// SetCount sets the total number of original trades (or ticks) in this candle.
func (c *Candle) SetCount(count int64) {
	c.count = count
}

// Open returns the first (open) price of this candle.
func (c *Candle) Open() float64 {
	return c.open
}

// SetOpen sets the first (open) price of this candle.
func (c *Candle) SetOpen(open float64) {
	c.open = open
}

// High returns the maximal (high) price of this candle.
func (c *Candle) High() float64 {
	return c.high
}

// SetHigh sets the maximal (high) price of this candle.
func (c *Candle) SetHigh(high float64) {
	c.high = high
}

// Low returns the minimal (low) price of this candle.
func (c *Candle) Low() float64 {
	return c.low
}

// SetLow sets the minimal (low) price of this candle.
func (c *Candle) SetLow(low float64) {
	c.low = low
}

// Close returns the last (close) price of this candle.
func (c *Candle) Close() float64 {
	return c.close
}

// SetClose sets the last (close) price of this candle.
func (c *Candle) SetClose(close float64) {
	c.close = close
}

// Volume returns the total volume of this candle.
func (c *Candle) Volume() float64 {
	return c.volume
}

// SetVolume sets the total volume of this candle.
func (c *Candle) SetVolume(volume float64) {
	c.volume = volume
}

// Vwap returns the volume-weighted average price (VWAP) of this candle.
func (c *Candle) Vwap() float64 {
	return c.vwap
}

// SetVwap sets the volume-weighted average price (VWAP) of this candle.
func (c *Candle) SetVwap(vwap float64) {
	c.vwap = vwap
}

// BidVolume returns the bid volume in this candle.
func (c *Candle) BidVolume() float64 {
	return c.bidVolume
}

// SetBidVolume sets the bid volume in this candle.
func (c *Candle) SetBidVolume(bidVolume float64) {
	c.bidVolume = bidVolume
}

// AskVolume returns the ask volume in this candle.
func (c *Candle) AskVolume() float64 {
	return c.askVolume
}

// SetAskVolume sets the ask volume in this candle.
func (c *Candle) SetAskVolume(askVolume float64) {
	c.askVolume = askVolume
}

// ImpVolatility returns the implied volatility of this candle.
func (c *Candle) ImpVolatility() float64 {
	return c.impVolatility
}

// SetImpVolatility sets the implied volatility of this candle.
func (c *Candle) SetImpVolatility(impVolatility float64) {
	c.impVolatility = impVolatility
}

// OpenInterest returns the open interest of this candle.
func (c *Candle) OpenInterest() float64 {
	return c.openInterest
}

// SetOpenInterest sets the open interest of this candle.
func (c *Candle) SetOpenInterest(openInterest float64) {
	c.openInterest = openInterest
}

// EventFlags returns the transactional event flags.
func (c *Candle) EventFlags() int32 {
	return c.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (c *Candle) SetEventFlags(eventFlags int32) {
	c.eventFlags = eventFlags
}

// Index returns the unique per-symbol index of this event.
func (c *Candle) Index() int64 {
	return c.index
}

// SetIndex sets the unique per-symbol index of this event.
func (c *Candle) SetIndex(index int64) {
	c.index = index
}

// Sequence returns the sequence number of this event within its millisecond.
func (c *Candle) Sequence() int64 {
	return c.index & maxSequence
}

// SetSequence sets the sequence number of this event within its millisecond.
func (q *Candle) SetSequence(sequence int64) error {
	if sequence < 0 || sequence > maxSequence {
		return fmt.Errorf("sequence(%d) is < 0 or > MaxSequence(%d)", sequence, maxSequence)
	}
	q.index = (q.index & ^maxSequence) | sequence
	return nil
}

// Time returns the time of this candle in milliseconds since the Unix epoch, encoded in the index.
func (q *Candle) Time() int64 {
	return ((q.index >> 32) * 1000) + ((q.index >> 22) & 0x3ff)
}

// SetTime sets the time of this candle in milliseconds since the Unix epoch, encoded in the index.
func (q *Candle) SetTime(value int64) {
	q.index = (timeutil.GetSecondsFromTime(value) << 32) |
		(int64(timeutil.GetMillisFromTime(value)) << 22) |
		q.Sequence()
}

// String returns a textual representation of the candle.
func (c *Candle) String() string {
	return "Candle{" + c.fieldsString() + "}"
}

// fieldsString returns the fields of the candle for String, without the type name and braces.
func (c *Candle) fieldsString() string {
	return c.EventSymbol().String() +
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
		", openInterest=" + formatutil.FormatFloat64(c.OpenInterest())
}

// Type returns the event code of the candle event.
func (c *Candle) Type() eventcodes.EventCode {
	return eventcodes.Candle
}
