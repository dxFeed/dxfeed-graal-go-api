package order

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"math"
	"strconv"
)

const (
	// ACTION values are taken from OrderAction enum.
	actionMask  = 0x0f
	actionShift = 11

	// EXCHANGE values are ASCII chars in [0, 127].
	exchangeMask  = 0x7f
	exchangeShift = 4

	// SIDE values are taken from Side enum.
	sideMask  = 3
	sideShift = 2

	// SCOPE values are taken from Scope enum.
	scopeMask  = 3
	scopeShift = 0
)

const maxSequence = (1 << 22) - 1

// Base holds the fields shared by all order events (Order, AnalyticOrder, SpreadOrder).
type Base struct {
	eventSymbol  *string
	eventTime    int64
	eventFlags   int32
	index        int64
	timeSequence int64
	timeNanoPart int32
	actionTime   int64
	orderId      int64
	auxOrderId   int64
	price        float64
	size         float64
	executedSize float64
	count        int64
	flags        int32
	tradeId      int64
	tradePrice   float64
	tradeSize    float64
}

// NewBase creates a Base for the given event symbol.
func NewBase(eventSymbol string) *Base {
	return &Base{eventSymbol: &eventSymbol,
		price:        math.NaN(),
		size:         math.NaN(),
		executedSize: math.NaN(),
		tradePrice:   math.NaN(),
		tradeSize:    math.NaN(),
	}
}

// EventSymbol returns the symbol of this event.
func (b *Base) EventSymbol() *string {
	return b.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (b *Base) SetEventSymbol(eventSymbol *string) {
	b.eventSymbol = eventSymbol
}

// EventTime returns the time of the event, used for transactional events.
func (b *Base) EventTime() int64 {
	return b.eventTime
}

// SetEventTime sets the time of the event.
func (b *Base) SetEventTime(eventTime int64) {
	b.eventTime = eventTime
}

// EventFlags returns the transactional event flags.
func (b *Base) EventFlags() int32 {
	return b.eventFlags
}

// SetEventFlags sets the transactional event flags.
func (b *Base) SetEventFlags(eventFlags int32) {
	b.eventFlags = eventFlags
}

// Index returns the unique per-symbol index of this order, which also encodes the order source.
func (b *Base) Index() int64 {
	return b.index
}

// SetIndex sets the order index; it returns an error if the index is negative.
func (b *Base) SetIndex(index int64) error {
	if index < 0 {
		return fmt.Errorf("negative index: %d", index)
	}
	b.index = index
	return nil
}

// TimeSequence returns the time and sequence of the order packed into one number.
func (b *Base) TimeSequence() int64 {
	return b.timeSequence
}

// SetTimeSequence sets the time and sequence of the order packed into one number.
func (b *Base) SetTimeSequence(timeSequence int64) {
	b.timeSequence = timeSequence
}

// TimeNanoPart returns the microseconds and nanoseconds part of the order time.
func (b *Base) TimeNanoPart() int32 {
	return b.timeNanoPart
}

// SetTimeNanoPart sets the microseconds and nanoseconds part of the order time.
func (b *Base) SetTimeNanoPart(timeNanoPart int32) {
	b.timeNanoPart = timeNanoPart
}

// ActionTime returns the time of the last action (in milliseconds since the epoch).
func (b *Base) ActionTime() int64 {
	return b.actionTime
}

// SetActionTime sets the time of the last action.
func (b *Base) SetActionTime(actionTime int64) {
	b.actionTime = actionTime
}

// OrderId returns the order ID.
func (b *Base) OrderId() int64 {
	return b.orderId
}

// SetOrderId sets the order ID.
func (b *Base) SetOrderId(orderId int64) {
	b.orderId = orderId
}

// AuxOrderId returns the auxiliary order ID (for example, the ID of the matching order for a trade).
func (b *Base) AuxOrderId() int64 {
	return b.auxOrderId
}

// SetAuxOrderId sets the auxiliary order ID.
func (b *Base) SetAuxOrderId(auxOrderId int64) {
	b.auxOrderId = auxOrderId
}

// Price returns the price of the order.
func (b *Base) Price() float64 {
	return b.price
}

// SetPrice sets the price of the order.
func (b *Base) SetPrice(price float64) {
	b.price = price
}

// Size returns the size of the order.
func (b *Base) Size() float64 {
	return b.size
}

// SetSize sets the size of the order.
func (b *Base) SetSize(size float64) {
	b.size = size
}

// ExecutedSize returns the executed size of the order.
func (b *Base) ExecutedSize() float64 {
	return b.executedSize
}

// SetExecutedSize sets the executed size of the order.
func (b *Base) SetExecutedSize(executedSize float64) {
	b.executedSize = executedSize
}

// Count returns the number of individual orders at this price level.
func (b *Base) Count() int64 {
	return b.count
}

// SetCount sets the number of individual orders at this price level.
func (b *Base) SetCount(count int64) {
	b.count = count
}

// Flags returns the raw flags bit field (action, exchange, side and scope).
func (b *Base) Flags() int32 {
	return b.flags
}

// SetFlags sets the raw flags bit field.
func (b *Base) SetFlags(flags int32) {
	b.flags = flags
}

// TradeId returns the trade ID.
func (b *Base) TradeId() int64 {
	return b.tradeId
}

// SetTradeId sets the trade ID.
func (b *Base) SetTradeId(tradeId int64) {
	b.tradeId = tradeId
}

// TradePrice returns the trade price.
func (b *Base) TradePrice() float64 {
	return b.tradePrice
}

// SetTradePrice sets the trade price.
func (b *Base) SetTradePrice(tradePrice float64) {
	b.tradePrice = tradePrice
}

// TradeSize returns the trade size.
func (b *Base) TradeSize() float64 {
	return b.tradeSize
}

// SetTradeSize sets the trade size.
func (b *Base) SetTradeSize(tradeSize float64) {
	b.tradeSize = tradeSize
}

// HasSize reports whether the order has a non-zero size.
func (b *Base) HasSize() bool {
	return b.size != 0 && !math.IsNaN(b.size)
}

// ExchangeCode returns the exchange code of the order as a character.
func (b *Base) ExchangeCode() rune {
	return rune(mathutil.GetBits(int64(b.flags), exchangeMask, exchangeShift))
}

// SetExchangeCode sets the exchange code; it returns an error if the code is not an ASCII character.
func (b *Base) SetExchangeCode(value rune) error {
	err := mathutil.CheckChar(int64(value), exchangeMask, "exchangeCode")
	if err != nil {
		return err
	}
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), exchangeMask, exchangeShift, int64(value))))
	return nil
}

// Side returns the side of the order.
func (b *Base) Side() side.Side {
	bits := mathutil.GetBits(int64(b.flags), sideMask, sideShift)
	return side.SideValueOf(bits)
}

// SetSide sets the side of the order.
func (b *Base) SetSide(value side.Side) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.flags), sideMask, sideShift, int64(value))))
}

// Time returns the order time in milliseconds since the epoch.
func (b *Base) Time() int64 {
	return ((b.timeSequence >> 32) * 1000) + ((b.timeSequence >> 22) & 0x3ff)
}

// SetTime sets the order time in milliseconds since the epoch.
func (b *Base) SetTime(value int64) {
	b.timeSequence = (timeutil.GetSecondsFromTime(value) << 32) |
		(int64(timeutil.GetMillisFromTime(value)) << 22) |
		b.Sequence()
}

// Sequence returns the sequence number of the order within the same millisecond.
func (b *Base) Sequence() int64 {
	return b.timeSequence & maxSequence
}

// SetSequence sets the sequence number; it returns an error if the value is out of range.
func (b *Base) SetSequence(value int64) error {
	if value < 0 || value > maxSequence {
		return fmt.Errorf("Sequence(%d) is < 0 or > MaxSequence(%d)", value, maxSequence)
	}

	b.timeSequence = (b.timeSequence & ^maxSequence) | value
	return nil
}

// TimeNanos returns the order time in nanoseconds since the epoch.
func (b *Base) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(b.Time(), b.TimeNanoPart())
}

// SetTimeNanos sets the order time in nanoseconds since the epoch.
func (b *Base) SetTimeNanos(value int64) {
	b.SetTime(timeutil.GetMillisFromNanos(value))
	b.SetTimeNanoPart(int32(timeutil.GetNanoPartFromNanos(value)))
}

// Action returns the action of the order.
func (b *Base) Action() Action {
	return ActionValueOf(int32(mathutil.GetBits(int64(b.Flags()), actionMask, actionShift)))
}

// SetAction sets the action of the order.
func (b *Base) SetAction(value Action) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), actionMask, actionShift, int64(value))))
}

// Scope returns the scope of the order.
func (b *Base) Scope() Scope {
	return ScopeValueOf(int32(mathutil.GetBits(int64(b.Flags()), scopeMask, scopeShift)))
}

// SetScope sets the scope of the order.
func (b *Base) SetScope(value Scope) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), scopeMask, scopeShift, int64(value))))
}

// OrderSourceName returns the name of the order source, or an empty string if it cannot be determined.
func (b *Base) OrderSourceName() *string {
	source, err := b.OrderSource()
	if err != nil {
		emptyStr := ""
		return &emptyStr
	}
	return source.Name()
}

// OrderSource returns the source of the order decoded from its index.
func (b *Base) OrderSource() (*Source, error) {
	sourceId := b.Index() >> 48
	if !IsSpecialSourceId(sourceId) {
		sourceId = b.Index() >> 32
	}
	value, err := ValueOfIdentifier(sourceId)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("incorrect value for eventsource %d", b.Index())
	}
	return value, nil
}

// SetOrderSource sets the source of the order in its index.
func (b *Base) SetOrderSource(value *Source) {
	var shift int64
	if IsSpecialSourceId(int64(value.Id())) {
		shift = 48
	} else {
		shift = 32
	}
	var mask int64
	if IsSpecialSourceId(b.index >> 48) {
		mask = ^(int64(-1) << 48)
	} else {
		mask = ^(int64(-1) << 32)
	}
	_ = b.SetIndex((int64(value.Id()) << shift) | (b.Index() & mask))
}

func (b *Base) baseFieldsToString() string {
	return formatutil.FormatString(b.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(b.EventTime()) +
		", source=" + formatutil.FormatString(b.OrderSourceName()) +
		", eventFlags=" + formatutil.HexFormat(int64(b.EventFlags())) +
		", index=" + formatutil.HexFormat(int64(b.Index())) +
		", time=" + formatutil.FormatTime(b.Time()) +
		", sequence=" + strconv.FormatInt(int64(b.Sequence()), 10) +
		", timeNanoPart=" + strconv.FormatInt(int64(b.TimeNanoPart()), 10) +
		", action=" + formatutil.FormatInt64(int64(b.Action())) +
		", actionTime=" + formatutil.FormatTime(b.ActionTime()) +
		", orderId=" + formatutil.FormatInt64(b.OrderId()) +
		", auxOrderId=" + formatutil.FormatInt64(b.AuxOrderId()) +
		", price=" + formatutil.FormatFloat64(b.Price()) +
		", size=" + formatutil.FormatFloat64(b.Size()) +
		", executedSize=" + formatutil.FormatFloat64(b.ExecutedSize()) +
		", count=" + formatutil.FormatInt64(b.Count()) +
		", exchange=" + formatutil.FormatChar(b.ExchangeCode()) +
		", side=" + b.Side().String() +
		", scope=" + b.Scope().String() +
		", tradeId=" + formatutil.FormatInt64(b.TradeId()) +
		", tradePrice=" + formatutil.FormatFloat64(b.TradePrice()) +
		", tradeSize=" + formatutil.FormatFloat64(b.TradeSize())
}
