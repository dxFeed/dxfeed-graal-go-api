// Package series provides the Series event: the properties of an option series, such as its implied volatility and
// put/call volumes.
package series

import (
	"math"
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// Series is an indexed event with the properties of an option series of an underlying symbol, for one expiration:
// implied volatility, call and put volumes, forward price, dividend and interest. The index is unique per expiration.
type Series struct {
	eventSymbol  *string
	eventTime    int64
	eventFlags   int32
	index        int64
	timeSequence int64
	expiration   int32

	volatility   float64
	callVolume   float64
	putVolume    float64
	putCallRatio float64
	forwardPrice float64
	dividend     float64
	interest     float64
}

// NewSeries creates a Series event for the given event symbol, the symbol of the underlying.
func NewSeries(eventSymbol string) *Series {
	return &Series{
		eventSymbol:  &eventSymbol,
		volatility:   math.NaN(),
		callVolume:   math.NaN(),
		putVolume:    math.NaN(),
		putCallRatio: math.NaN(),
		forwardPrice: math.NaN(),
		dividend:     math.NaN(),
		interest:     math.NaN(),
	}
}

// Type returns the event code of this event type.
func (s *Series) Type() eventcodes.EventCode {
	return eventcodes.Series
}

// Source is always DEFAULT for Series.
func (s *Series) Source() *events.IndexedEventSource {
	return events.DefaultIndexedEventSource()
}

// EventSymbol returns the symbol of this event.
func (s *Series) EventSymbol() *string {
	return s.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (s *Series) SetEventSymbol(eventSymbol string) {
	s.eventSymbol = &eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (s *Series) EventTime() int64 {
	return s.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (s *Series) SetEventTime(eventTime int64) {
	s.eventTime = eventTime
}

// EventFlags returns the transactional event flags; see events.EventFlag.
func (s *Series) EventFlags() int32 {
	return s.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (s *Series) SetEventFlags(eventFlags int32) {
	s.eventFlags = eventFlags
}

// Index returns the unique per-symbol index of this series.
func (s *Series) Index() int64 {
	return s.index
}

// SetIndex sets the unique per-symbol index of this series.
func (s *Series) SetIndex(index int64) {
	s.index = index
}

// TimeSequence returns the time and the sequence of this event packed in one value.
func (s *Series) TimeSequence() int64 {
	return s.timeSequence
}

// SetTimeSequence sets the time and the sequence of this event packed in one value.
func (s *Series) SetTimeSequence(timeSequence int64) {
	s.timeSequence = timeSequence
}

// Time returns the time of this event in milliseconds.
func (s *Series) Time() int64 {
	return timeutil.TimeOfTimeSequence(s.timeSequence)
}

// SetTime sets the time of this event in milliseconds; the sequence is kept.
func (s *Series) SetTime(timeMillis int64) {
	s.timeSequence = timeutil.WithTime(s.timeSequence, timeMillis)
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (s *Series) Sequence() int32 {
	return timeutil.SequenceOfTimeSequence(s.timeSequence)
}

// SetSequence sets the sequence number of this event and returns an error if it is out of range.
func (s *Series) SetSequence(sequence int32) error {
	timeSequence, err := timeutil.WithSequence(s.timeSequence, sequence)
	s.timeSequence = timeSequence
	return err
}

// Expiration returns the day id (days since 1970-01-01) of the expiration of this series.
func (s *Series) Expiration() int32 {
	return s.expiration
}

// SetExpiration sets the day id of the expiration of this series.
func (s *Series) SetExpiration(expiration int32) {
	s.expiration = expiration
}

// Volatility returns the implied volatility of this series.
func (s *Series) Volatility() float64 {
	return s.volatility
}

// SetVolatility sets the implied volatility of this series.
func (s *Series) SetVolatility(volatility float64) {
	s.volatility = volatility
}

// CallVolume returns the traded volume of the call options of this series for the day.
func (s *Series) CallVolume() float64 {
	return s.callVolume
}

// SetCallVolume sets the traded volume of the call options.
func (s *Series) SetCallVolume(callVolume float64) {
	s.callVolume = callVolume
}

// PutVolume returns the traded volume of the put options of this series for the day.
func (s *Series) PutVolume() float64 {
	return s.putVolume
}

// SetPutVolume sets the traded volume of the put options.
func (s *Series) SetPutVolume(putVolume float64) {
	s.putVolume = putVolume
}

// OptionVolume returns the traded volume of the options of this series: the sum of the put and call volumes, or the
// one of them that is known.
func (s *Series) OptionVolume() float64 {
	if math.IsNaN(s.putVolume) {
		return s.callVolume
	}
	if math.IsNaN(s.callVolume) {
		return s.putVolume
	}
	return s.putVolume + s.callVolume
}

// PutCallRatio returns the ratio of the put volume to the call volume.
func (s *Series) PutCallRatio() float64 {
	return s.putCallRatio
}

// SetPutCallRatio sets the ratio of the put volume to the call volume.
func (s *Series) SetPutCallRatio(putCallRatio float64) {
	s.putCallRatio = putCallRatio
}

// ForwardPrice returns the implied forward price of this series.
func (s *Series) ForwardPrice() float64 {
	return s.forwardPrice
}

// SetForwardPrice sets the implied forward price.
func (s *Series) SetForwardPrice(forwardPrice float64) {
	s.forwardPrice = forwardPrice
}

// Dividend returns the implied simple dividend return of the underlying.
func (s *Series) Dividend() float64 {
	return s.dividend
}

// SetDividend sets the implied simple dividend return.
func (s *Series) SetDividend(dividend float64) {
	s.dividend = dividend
}

// Interest returns the implied simple interest return of the underlying.
func (s *Series) Interest() float64 {
	return s.interest
}

// SetInterest sets the implied simple interest return.
func (s *Series) SetInterest(interest float64) {
	s.interest = interest
}

// String returns a text representation of the event.
func (s *Series) String() string {
	return "Series{" + formatutil.FormatString(s.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(s.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(s.EventFlags())) +
		", index=" + formatutil.HexFormat(s.Index()) +
		", time=" + formatutil.FormatTime(s.Time()) +
		", sequence=" + strconv.FormatInt(int64(s.Sequence()), 10) +
		", expiration=" + strconv.FormatInt(int64(timeutil.GetYearMonthDayByDayId(s.Expiration())), 10) +
		", volatility=" + formatutil.FormatFloat64(s.volatility) +
		", callVolume=" + formatutil.FormatFloat64(s.callVolume) +
		", putVolume=" + formatutil.FormatFloat64(s.putVolume) +
		", putCallRatio=" + formatutil.FormatFloat64(s.putCallRatio) +
		", forwardPrice=" + formatutil.FormatFloat64(s.forwardPrice) +
		", dividend=" + formatutil.FormatFloat64(s.dividend) +
		", interest=" + formatutil.FormatFloat64(s.interest) +
		"}"
}
