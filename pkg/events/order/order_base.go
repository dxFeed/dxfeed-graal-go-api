package order

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
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

func NewBase(eventSymbol string) *Base {
	return &Base{eventSymbol: &eventSymbol,
		price:        math.NaN(),
		size:         math.NaN(),
		executedSize: math.NaN(),
		tradePrice:   math.NaN(),
		tradeSize:    math.NaN(),
	}
}

func (b *Base) EventSymbol() *string {
	return b.eventSymbol
}

func (b *Base) SetEventSymbol(eventSymbol *string) {
	b.eventSymbol = eventSymbol
}

func (b *Base) EventTime() int64 {
	return b.eventTime
}

func (b *Base) SetEventTime(eventTime int64) {
	b.eventTime = eventTime
}

func (b *Base) EventFlags() int32 {
	return b.eventFlags
}

func (b *Base) SetEventFlags(eventFlags int32) {
	b.eventFlags = eventFlags
}

func (b *Base) Index() int64 {
	return b.index
}

func (b *Base) SetIndex(index int64) error {
	if index < 0 {
		return fmt.Errorf("negative index: %d", index)
	}
	b.index = index
	return nil
}

func (b *Base) TimeSequence() int64 {
	return b.timeSequence
}

func (b *Base) SetTimeSequence(timeSequence int64) {
	b.timeSequence = timeSequence
}

func (b *Base) TimeNanoPart() int32 {
	return b.timeNanoPart
}

func (b *Base) SetTimeNanoPart(timeNanoPart int32) {
	b.timeNanoPart = timeNanoPart
}

func (b *Base) ActionTime() int64 {
	return b.actionTime
}

func (b *Base) SetActionTime(actionTime int64) {
	b.actionTime = actionTime
}

func (b *Base) OrderId() int64 {
	return b.orderId
}

func (b *Base) SetOrderId(orderId int64) {
	b.orderId = orderId
}

func (b *Base) AuxOrderId() int64 {
	return b.auxOrderId
}

func (b *Base) SetAuxOrderId(auxOrderId int64) {
	b.auxOrderId = auxOrderId
}

func (b *Base) Price() float64 {
	return b.price
}

func (b *Base) SetPrice(price float64) {
	b.price = price
}

func (b *Base) Size() float64 {
	return b.size
}

func (b *Base) SetSize(size float64) {
	b.size = size
}

func (b *Base) ExecutedSize() float64 {
	return b.executedSize
}

func (b *Base) SetExecutedSize(executedSize float64) {
	b.executedSize = executedSize
}

func (b *Base) Count() int64 {
	return b.count
}

func (b *Base) SetCount(count int64) {
	b.count = count
}

func (b *Base) Flags() int32 {
	return b.flags
}

func (b *Base) SetFlags(flags int32) {
	b.flags = flags
}

func (b *Base) TradeId() int64 {
	return b.tradeId
}

func (b *Base) SetTradeId(tradeId int64) {
	b.tradeId = tradeId
}

func (b *Base) TradePrice() float64 {
	return b.tradePrice
}

func (b *Base) SetTradePrice(tradePrice float64) {
	b.tradePrice = tradePrice
}

func (b *Base) TradeSize() float64 {
	return b.tradeSize
}

func (b *Base) SetTradeSize(tradeSize float64) {
	b.tradeSize = tradeSize
}

func (b *Base) HasSize() bool {
	return b.size != 0 && !math.IsNaN(b.size)
}

func (b *Base) ExchangeCode() rune {
	return rune(mathutil.GetBits(int64(b.flags), exchangeMask, exchangeShift))
}

func (b *Base) SetExchangeCode(value rune) error {
	err := mathutil.CheckChar(int64(value), exchangeMask, "exchangeCode")
	if err != nil {
		return err
	}
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), exchangeMask, exchangeShift, int64(value))))
	return nil
}

func (b *Base) Side() side.Side {
	bits := mathutil.GetBits(int64(b.flags), sideMask, sideShift)
	return side.SideValueOf(bits)
}

func (b *Base) SetSide(value side.Side) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.flags), sideMask, sideShift, int64(value))))
}

func (b *Base) Time() int64 {
	return ((b.index >> 32) * 1000) + ((b.index >> 22) & 0x3ff)
}

func (b *Base) SetTime(value int64) {
	b.index = (timeutil.GetSecondsFromTime(value) << 32) |
		int64(timeutil.GetMillisFromTime(value)<<22) |
		b.Sequence()
}

func (b *Base) Sequence() int64 {
	return b.index & maxSequence
}

func (b *Base) SetSequence(value int64) error {
	if value < 0 || value > maxSequence {
		return fmt.Errorf("Sequence(%d) is < 0 or > MaxSequence(%d)", value, maxSequence)
	}

	b.index = (b.index & ^maxSequence) | value
	return nil
}

func (b *Base) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(b.Time(), b.TimeNanoPart())
}

func (b *Base) SetTimeNanos(value int64) {
	b.SetTime(timeutil.GetMillisFromNanos(value))
	b.SetTimeNanoPart(int32(timeutil.GetNanoPartFromNanos(value)))
}

func (b *Base) Action() Action {
	return ActionValueOf(int32(mathutil.GetBits(int64(b.Flags()), actionMask, actionShift)))
}

func (b *Base) SetAction(value Action) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), actionMask, actionShift, int64(value))))
}

func (b *Base) Scope() Scope {
	return ScopeValueOf(int32(mathutil.GetBits(int64(b.Flags()), scopeMask, scopeShift)))
}

func (b *Base) SetScope(value Scope) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), scopeMask, scopeShift, int64(value))))
}

func (b *Base) OrderSourceName() *string {
	source, err := b.OrderSource()
	if err != nil {
		emptyStr := ""
		return &emptyStr
	}
	return source.Name()
}

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
