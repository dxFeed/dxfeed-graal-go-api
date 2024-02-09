package trade

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
	"strconv"
)

const (
	// DIRECTION values are taken from Direction enum.
	directionMask  = 7
	directionShift = 1
	// ETH mask.
	eth = 1

	maxSequence = (1 << 22) - 1
)

type TradeBase struct {
	eventSymbol  *string
	eventTime    int64
	timeSequence int64
	timeNanoPart int32
	exchangeCode int16
	price        float64
	change       float64
	size         float64
	dayId        int32
	dayVolume    float64
	dayTurnover  float64
	flags        int32
}

func NewTradeBase(eventSymbol string) *TradeBase {
	return &TradeBase{eventSymbol: &eventSymbol}
}

func (t *TradeBase) EventSymbol() *string {
	return t.eventSymbol
}

func (t *TradeBase) SetEventSymbol(eventSymbol *string) {
	t.eventSymbol = eventSymbol
}

func (t *TradeBase) TimeSequence() int64 {
	return t.timeSequence
}

func (t *TradeBase) SetTimeSequence(timeSequence int64) {
	t.timeSequence = timeSequence
}

func (t *TradeBase) EventTime() int64 {
	return t.eventTime
}

func (t *TradeBase) SetEventTime(eventTime int64) {
	t.eventTime = eventTime
}

func (t *TradeBase) TimeNanoPart() int32 {
	return t.timeNanoPart
}

func (t *TradeBase) SetTimeNanoPart(timeNanoPart int32) {
	t.timeNanoPart = timeNanoPart
}

func (t *TradeBase) ExchangeCode() int16 {
	return t.exchangeCode
}

func (t *TradeBase) SetExchangeCode(exchangeCode int16) {
	t.exchangeCode = exchangeCode
}

func (t *TradeBase) Price() float64 {
	return t.price
}

func (t *TradeBase) SetPrice(price float64) {
	t.price = price
}

func (t *TradeBase) Change() float64 {
	return t.change
}

func (t *TradeBase) SetChange(change float64) {
	t.change = change
}

func (t *TradeBase) Size() float64 {
	return t.size
}

func (t *TradeBase) SetSize(size float64) {
	t.size = size
}

func (t *TradeBase) DayId() int32 {
	return t.dayId
}

func (t *TradeBase) SetDayId(dayId int32) {
	t.dayId = dayId
}

func (t *TradeBase) DayVolume() float64 {
	return t.dayVolume
}

func (t *TradeBase) SetDayVolume(dayVolume float64) {
	t.dayVolume = dayVolume
}

func (t *TradeBase) DayTurnover() float64 {
	return t.dayTurnover
}

func (t *TradeBase) SetDayTurnover(dayTurnover float64) {
	t.dayTurnover = dayTurnover
}

func (t *TradeBase) Flags() int32 {
	return t.flags
}

func (t *TradeBase) SetFlags(flags int32) {
	t.flags = flags
}

func (t *TradeBase) TickDirection() Direction {
	return ValueOf(mathutil.GetBits(int64(t.Flags()), directionMask, directionShift))
}

func (t *TradeBase) SetTickDirection(value Direction) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.Flags()), directionMask, directionShift, int64(value))))
}

func (t *TradeBase) IsExtendedTradingHours() bool {
	return (t.Flags() & eth) != 0
}

func (t *TradeBase) SetIsExtendedTradingHours(value bool) {
	if value {
		t.SetFlags(t.Flags() | eth)
	} else {
		t.SetFlags(t.Flags() & ^eth)
	}
}

func (t *TradeBase) Sequence() int64 {
	return t.TimeSequence() & maxSequence
}

func (t *TradeBase) SetSequence(value int64) error {
	if value < 0 || value > maxSequence {
		return fmt.Errorf("sequence(%d) is < 0 or > MaxSequence(%d)", value, maxSequence)
	}
	t.timeSequence = (t.TimeSequence() & ^maxSequence) | value
	return nil
}

func (t *TradeBase) Time() int64 {
	return ((t.TimeSequence() >> 32) * 1000) + ((t.TimeSequence() >> 22) & 0x3ff)
}

func (t *TradeBase) SetTime(value int64) {
	t.timeSequence = timeutil.GetSecondsFromTime(value)<<32 |
		int64(timeutil.GetMillisFromTime(value)<<22) |
		value
}

func (t *TradeBase) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(t.Time(), t.TimeNanoPart())
}

func (t *TradeBase) Type() eventcodes.EventCode {
	return eventcodes.Trade
}

func (t *TradeBase) SetTimeNanos(value int64) {
	t.SetTime(timeutil.GetMillisFromNanos(value))
	t.SetTimeNanoPart(int32(timeutil.GetNanoPartFromNanos(value)))
}

func (t *TradeBase) String() string {
	return formatutil.FormatString(t.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(t.EventTime()) +
		", time=" + formatutil.FormatTime(t.Time()) +
		", sequence=" + strconv.FormatInt(t.Sequence(), 10) +
		", timeNanoPart=" + strconv.FormatInt(int64(t.TimeNanoPart()), 10) +
		", exchange=" + formatutil.FormatChar(rune(t.ExchangeCode())) +
		", price=" + formatutil.FormatFloat64(t.Price()) +
		", change=" + formatutil.FormatFloat64(t.Change()) +
		", size=" + formatutil.FormatFloat64(t.Size()) +
		", day=" + strconv.FormatInt(int64(timeutil.GetYearMonthDayByDayId(t.DayId())), 10) +
		", dayVolume=" + formatutil.FormatFloat64(t.DayVolume()) +
		", dayTurnover=" + formatutil.FormatFloat64(t.DayTurnover()) +
		", direction=" + formatutil.FormatInt64(int64(t.TickDirection())) +
		", ETH=" + formatutil.FormatBool(t.IsExtendedTradingHours())
}
