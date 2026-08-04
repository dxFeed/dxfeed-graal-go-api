package underlying

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
	"math"
	"strconv"
)

const maxSequence = (1 << 22) - 1

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

func (u *Underlying) Type() eventcodes.EventCode {
	return eventcodes.Underlying
}

// Source is always DEFAULT for Underlying (multiple sources per symbol are not supported).
func (u *Underlying) Source() *events.IndexedEventSource {
	return events.DefaultIndexedEventSource()
}

func (u *Underlying) EventSymbol() *string {
	return u.eventSymbol
}

func (u *Underlying) SetEventSymbol(eventSymbol string) {
	*u.eventSymbol = eventSymbol
}

func (u *Underlying) EventTime() int64 {
	return u.eventTime
}

func (u *Underlying) SetEventTime(eventTime int64) {
	u.eventTime = eventTime
}

func (u *Underlying) EventFlags() int32 {
	return u.eventFlags
}

func (u *Underlying) SetEventFlags(eventFlags int32) {
	u.eventFlags = eventFlags
}

func (u *Underlying) Index() int64 {
	return u.index
}

func (u *Underlying) SetIndex(index int64) {
	u.index = index
}

func (u *Underlying) Sequence() int32 {
	return int32(u.index & maxSequence)
}

func (u *Underlying) SetSequence(sequence int32) error {
	if sequence < 0 || int64(sequence) > maxSequence {
		return fmt.Errorf("Sequence(%d) is < 0 or > MaxSequence(%d)", sequence, maxSequence)
	}
	u.index = (u.index & ^maxSequence) | int64(sequence)
	return nil
}

func (u *Underlying) Time() int64 {
	return ((u.index >> 32) * 1000) + ((u.index >> 22) & 0x3ff)
}

func (u *Underlying) SetTime(value int64) {
	u.index = (timeutil.GetSecondsFromTime(value) << 32) |
		(int64(timeutil.GetMillisFromTime(value)) << 22) |
		int64(u.Sequence())
}

func (u *Underlying) Volatility() float64 {
	return u.volatility
}

func (u *Underlying) SetVolatility(volatility float64) {
	u.volatility = volatility
}

func (u *Underlying) FrontVolatility() float64 {
	return u.frontVolatility
}

func (u *Underlying) SetFrontVolatility(frontVolatility float64) {
	u.frontVolatility = frontVolatility
}

func (u *Underlying) BackVolatility() float64 {
	return u.backVolatility
}

func (u *Underlying) SetBackVolatility(backVolatility float64) {
	u.backVolatility = backVolatility
}

func (u *Underlying) CallVolume() float64 {
	return u.callVolume
}

func (u *Underlying) SetCallVolume(callVolume float64) {
	u.callVolume = callVolume
}

func (u *Underlying) PutVolume() float64 {
	return u.putVolume
}

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

func (u *Underlying) PutCallRatio() float64 {
	return u.putCallRatio
}

func (u *Underlying) SetPutCallRatio(putCallRatio float64) {
	u.putCallRatio = putCallRatio
}

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
