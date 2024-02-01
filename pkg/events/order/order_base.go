package order

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
	"math"
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
	b.index = (int64(timeutil.GetSecondsFromTime(value)) << 32) |
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
	return ActionValueOf(mathutil.GetBits(int64(b.Flags()), actionMask, actionShift))
}

func (b *Base) SetAction(value Action) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), actionMask, actionShift, int64(value))))
}

func (b *Base) Scope() Scope {
	return ScopeValueOf(mathutil.GetBits(int64(b.Flags()), scopeMask, scopeShift))
}

func (b *Base) SetScope(value Scope) {
	b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), scopeMask, scopeShift, int64(value))))
}

func (b *Base) OrderSource() (*OrderSource, error) {
	sourceId := b.Index() >> 48
	if !IsSpecialSourceId(sourceId) {
		sourceId = b.Index() >> 32
	}
	value, err := OrderSourceValueOfIdentifier(sourceId)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *Base) SetOrderSource(value *OrderSource) {

}

func (b *Base) SetIndexedEventSource(value OrderSource) {
	//b.SetFlags(int32(mathutil.SetBits(int64(b.Flags()), scopeMask, scopeShift, int64(value))))
}

//
//func (b *Base) baseFieldsToString() string {
//	return formatutil.FormatString(b.EventSymbol()) +
//		", eventTime=" + formatutil.FormatTime(b.EventTime()) +
//		", source=" + formatutil.HexFormat(int64(t.EventFlags())) +
//		", eventFlags=" + formatutil.HexFormat(int64(t.EventFlags())) +
//		", time=" + formatutil.FormatTime(t.Time()) +
//		", timeNanoPart=" + strconv.FormatInt(int64(t.timeNanoPart), 10) +
//		", sequence=" + strconv.FormatInt(int64(t.Sequence()), 10) +
//		", exchange=" + formatutil.FormatChar(rune(t.exchangeCode)) +
//		", price=" + formatutil.FormatFloat64(t.Price()) +
//		", size=" + formatutil.FormatFloat64(t.Size()) +
//		", bid=" + formatutil.FormatFloat64(t.BidPrice()) +
//		", ask=" + formatutil.FormatFloat64(t.AskPrice()) +
//		", ESC=" + formatutil.FormatString(t.ExchangeSaleConditions()) +
//		", TTE=" + formatutil.FormatChar(t.TradeThroughExempt()) +
//		", side=" + formatutil.FormatInt64(int64(t.AggressorSide())) +
//		", spread=" + formatutil.FormatBool(t.IsSpreadLeg()) +
//		", ETH=" + formatutil.FormatBool(t.IsExtendedTradingHours()) +
//		", validTick=" + formatutil.FormatBool(t.IsValidTick()) +
//		", type=" + formatutil.FormatInt64(int64(t.TimeAndSaleType())) +
//		formatutil.FormatNullableString(", buyer=%s", t.Buyer(), "") +
//		formatutil.FormatNullableString(", seller=%s", t.Seller(), "") +
//		"}"
//}
