package trade

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
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

// TradeBase is the base for the Trade and TradeETH events: the price, size and daily statistics of the latest trade.
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

// NewTradeBase creates a new TradeBase for the given event symbol.
func NewTradeBase(eventSymbol string) *TradeBase {
	return &TradeBase{eventSymbol: &eventSymbol}
}

// EventSymbol returns the symbol of this event.
func (t *TradeBase) EventSymbol() *string {
	return t.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (t *TradeBase) SetEventSymbol(eventSymbol *string) {
	t.eventSymbol = eventSymbol
}

// TimeSequence returns the time and sequence of the last trade packed into a single value.
func (t *TradeBase) TimeSequence() int64 {
	return t.timeSequence
}

// SetTimeSequence sets the time and sequence of the last trade packed into a single value.
func (t *TradeBase) SetTimeSequence(timeSequence int64) {
	t.timeSequence = timeSequence
}

// EventTime returns the source-provided time of this event in milliseconds since the Unix epoch.
func (t *TradeBase) EventTime() int64 {
	return t.eventTime
}

// SetEventTime sets the source-provided time of this event in milliseconds since the Unix epoch.
func (t *TradeBase) SetEventTime(eventTime int64) {
	t.eventTime = eventTime
}

// TimeNanoPart returns the microseconds and nanoseconds part of the event time.
func (t *TradeBase) TimeNanoPart() int32 {
	return t.timeNanoPart
}

// SetTimeNanoPart sets the microseconds and nanoseconds part of the event time.
func (t *TradeBase) SetTimeNanoPart(timeNanoPart int32) {
	t.timeNanoPart = timeNanoPart
}

// ExchangeCode returns the exchange code of the event.
func (t *TradeBase) ExchangeCode() int16 {
	return t.exchangeCode
}

// SetExchangeCode sets the exchange code of the event.
func (t *TradeBase) SetExchangeCode(exchangeCode int16) {
	t.exchangeCode = exchangeCode
}

// Price returns the price of the last trade.
func (t *TradeBase) Price() float64 {
	return t.price
}

// SetPrice sets the price of the last trade.
func (t *TradeBase) SetPrice(price float64) {
	t.price = price
}

// Change returns the price change of the last trade relative to the previous day close.
func (t *TradeBase) Change() float64 {
	return t.change
}

// SetChange sets the price change of the last trade relative to the previous day close.
func (t *TradeBase) SetChange(change float64) {
	t.change = change
}

// Size returns the size of the last trade.
func (t *TradeBase) Size() float64 {
	return t.size
}

// SetSize sets the size of the last trade.
func (t *TradeBase) SetSize(size float64) {
	t.size = size
}

// DayId returns the identifier of the current trading day, the number of days since the Unix epoch.
func (t *TradeBase) DayId() int32 {
	return t.dayId
}

// SetDayId sets the identifier of the current trading day, the number of days since the Unix epoch.
func (t *TradeBase) SetDayId(dayId int32) {
	t.dayId = dayId
}

// DayVolume returns the total volume traded for the day.
func (t *TradeBase) DayVolume() float64 {
	return t.dayVolume
}

// SetDayVolume sets the total volume traded for the day.
func (t *TradeBase) SetDayVolume(dayVolume float64) {
	t.dayVolume = dayVolume
}

// DayTurnover returns the total turnover traded for the day.
func (t *TradeBase) DayTurnover() float64 {
	return t.dayTurnover
}

// SetDayTurnover sets the total turnover traded for the day.
func (t *TradeBase) SetDayTurnover(dayTurnover float64) {
	t.dayTurnover = dayTurnover
}

// Flags returns the packed flags of this event.
func (t *TradeBase) Flags() int32 {
	return t.flags
}

// SetFlags sets the packed flags of this event.
func (t *TradeBase) SetFlags(flags int32) {
	t.flags = flags
}

// TickDirection returns the tick direction of the last trade relative to the previous one.
func (t *TradeBase) TickDirection() Direction {
	return ValueOf(mathutil.GetBits(int64(t.Flags()), directionMask, directionShift))
}

// SetTickDirection sets the tick direction of the last trade relative to the previous one.
func (t *TradeBase) SetTickDirection(value Direction) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.Flags()), directionMask, directionShift, int64(value))))
}

// IsExtendedTradingHours reports whether the trade occurred during extended trading hours (ETH).
func (t *TradeBase) IsExtendedTradingHours() bool {
	return (t.Flags() & eth) != 0
}

// SetIsExtendedTradingHours sets the extended trading hours (ETH) flag.
func (t *TradeBase) SetIsExtendedTradingHours(value bool) {
	if value {
		t.SetFlags(t.Flags() | eth)
	} else {
		t.SetFlags(t.Flags() & ^eth)
	}
}

// Sequence returns the sequence number of this event within the same millisecond.
func (t *TradeBase) Sequence() int64 {
	return t.TimeSequence() & maxSequence
}

// SetSequence sets the sequence number of this event within the same millisecond; it returns an error if out of range.
func (t *TradeBase) SetSequence(value int64) error {
	if value < 0 || value > maxSequence {
		return fmt.Errorf("sequence(%d) is < 0 or > MaxSequence(%d)", value, maxSequence)
	}
	t.timeSequence = (t.TimeSequence() & ^maxSequence) | value
	return nil
}

// Time returns the time of this event in milliseconds since the Unix epoch.
func (t *TradeBase) Time() int64 {
	return ((t.TimeSequence() >> 32) * 1000) + ((t.TimeSequence() >> 22) & 0x3ff)
}

// SetTime sets the time of this event in milliseconds since the Unix epoch.
func (t *TradeBase) SetTime(value int64) {
	t.timeSequence = timeutil.GetSecondsFromTime(value)<<32 |
		(int64(timeutil.GetMillisFromTime(value)) << 22) |
		t.Sequence()
}

// TimeNanos returns the time of this event in nanoseconds since the Unix epoch.
func (t *TradeBase) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(t.Time(), t.TimeNanoPart())
}

// Type returns the event code of this event.
func (t *TradeBase) Type() eventcodes.EventCode {
	return eventcodes.Trade
}

// SetTimeNanos sets the time of this event in nanoseconds since the Unix epoch.
func (t *TradeBase) SetTimeNanos(value int64) {
	t.SetTime(timeutil.GetMillisFromNanos(value))
	t.SetTimeNanoPart(int32(timeutil.GetNanoPartFromNanos(value)))
}

// String returns a string representation of this event.
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
