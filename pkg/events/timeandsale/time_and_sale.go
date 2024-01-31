package timeandsale

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/timeutil"
	"math"
	"strconv"
)

const (
	// TTE (TradeThroughExempt) values are ASCII chars in [0, 255].
	tteMask  = 0xff
	tteShift = 8

	// SIDE values are taken from Side enum.
	sideMask  = 3
	sideShift = 5

	spreadLeg = 1 << 4
	eth       = 1 << 3
	validTick = 1 << 2

	// TYPE values are taken from TimeAndSaleType enum.
	typeMask  = 3
	typeShift = 0
)

type TimeAndSale struct {
	eventSymbol            *string
	eventTime              int64
	timeNanoPart           int32
	exchangeCode           int16
	price                  float64
	size                   float64
	bidPrice               float64
	askPrice               float64
	exchangeSaleConditions *string
	flags                  int32
	buyer                  *string
	seller                 *string
	eventFlags             int32
	index                  int64
}

const maxSequence = (1 << 22) - 1

func NewTimeAndSale(eventSymbol string) *TimeAndSale {
	return &TimeAndSale{
		eventSymbol: &eventSymbol,
		price:       math.NaN(),
		size:        math.NaN(),
		bidPrice:    math.NaN(),
		askPrice:    math.NaN(),
	}
}

func (t *TimeAndSale) Type() eventcodes.EventCode {
	return eventcodes.TimeAndSale
}

func (t *TimeAndSale) EventSymbol() *string {
	return t.eventSymbol
}

func (t *TimeAndSale) SetEventSymbol(value string) {
	*t.eventSymbol = value
}

func (t *TimeAndSale) EventTime() int64 {
	return t.eventTime
}

func (t *TimeAndSale) SetEventTime(value int64) {
	t.eventTime = value
}

func (t *TimeAndSale) TimeNanoPart() int32 {
	return t.timeNanoPart
}

func (t *TimeAndSale) SetTimeNanoPart(value int32) {
	t.timeNanoPart = value
}

func (t *TimeAndSale) ExchangeCode() int16 {
	return t.exchangeCode
}

func (t *TimeAndSale) SetExchangeCode(value int16) {
	t.exchangeCode = value
}

func (t *TimeAndSale) Price() float64 {
	return t.price
}

func (t *TimeAndSale) SetPrice(value float64) {
	t.price = value
}

func (t *TimeAndSale) Size() float64 {
	return t.size
}

func (t *TimeAndSale) SetSize(value float64) {
	t.size = value
}

func (t *TimeAndSale) BidPrice() float64 {
	return t.bidPrice
}

func (t *TimeAndSale) SetBidPrice(value float64) {
	t.bidPrice = value
}

func (t *TimeAndSale) AskPrice() float64 {
	return t.askPrice
}

func (t *TimeAndSale) SetAskPrice(value float64) {
	t.askPrice = value
}

func (t *TimeAndSale) ExchangeSaleConditions() *string {
	return t.exchangeSaleConditions
}

func (t *TimeAndSale) SetExchangeSaleConditions(value *string) {
	t.exchangeSaleConditions = value
}

func (t *TimeAndSale) Buyer() *string {
	return t.buyer
}

func (t *TimeAndSale) SetBuyer(value *string) {
	t.buyer = value
}

func (t *TimeAndSale) Seller() *string {
	return t.seller
}

func (t *TimeAndSale) SetSeller(value *string) {
	t.seller = value
}

func (t *TimeAndSale) EventFlags() int32 {
	return t.eventFlags
}

func (t *TimeAndSale) SetEventFlags(value int32) {
	t.eventFlags = value
}

func (t *TimeAndSale) Index() int64 {
	return t.index
}

func (t *TimeAndSale) SetIndex(value int64) {
	t.index = value
}

func (t *TimeAndSale) Sequence() int64 {
	return t.index & maxSequence
}

func (t *TimeAndSale) SetSequence(value int64) error {
	if value < 0 || value > maxSequence {
		return fmt.Errorf("Sequence(%d) is < 0 or > MaxSequence(%d)", value, maxSequence)
	}

	t.index = (t.index & ^maxSequence) | value
	return nil
}

func (t *TimeAndSale) Time() int64 {
	return ((t.index >> 32) * 1000) + ((t.index >> 22) & 0x3ff)
}

func (t *TimeAndSale) SetTime(value int64) {
	t.index = (int64(timeutil.GetSecondsFromTime(value)) << 32) |
		int64(timeutil.GetMillisFromTime(value)<<22) |
		t.Sequence()
}

func (t *TimeAndSale) TimeNanos() int64 {
	return timeutil.GetNanosFromMillisAndNanoPart(t.Time(), t.TimeNanoPart())
}

func (t *TimeAndSale) SetTimeNanos(value int64) {
	t.SetTime(timeutil.GetMillisFromNanos(value))
	t.SetTimeNanoPart(int32(timeutil.GetNanoPartFromNanos(value)))
}

func (t *TimeAndSale) TradeThroughExempt() rune {
	return rune(mathutil.GetBits(int64(t.flags), tteMask, tteShift))
}

func (t *TimeAndSale) SetTradeThroughExempt(value rune) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.flags), tteMask, tteShift, int64(value))))
}

func (t *TimeAndSale) AggressorSide() side.Side {
	bits := mathutil.GetBits(int64(t.flags), sideMask, sideShift)
	return side.SideValueOf(bits)
}

func (t *TimeAndSale) SetAggressorSide(value side.Side) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.flags), sideMask, sideShift, int64(value))))
}

func (t *TimeAndSale) IsSpreadLeg() bool {
	return (t.flags & spreadLeg) != 0
}

func (t *TimeAndSale) SetIsSpreadLeg(value bool) {
	if value {
		t.SetFlags(t.flags | spreadLeg)
	} else {
		t.SetFlags(t.flags & ^spreadLeg)
	}
}

func (t *TimeAndSale) IsExtendedTradingHours() bool {
	return (t.flags & eth) != 0
}

func (t *TimeAndSale) SetIsExtendedTradingHours(value bool) {
	if value {
		t.SetFlags(t.flags | eth)
	} else {
		t.SetFlags(t.flags & ^eth)
	}
}

func (t *TimeAndSale) IsValidTick() bool {
	return (t.flags & validTick) != 0
}

func (t *TimeAndSale) SetIsValidTick(value bool) {
	if value {
		t.SetFlags(t.flags | validTick)
	} else {
		t.SetFlags(t.flags & ^validTick)
	}
}

func (t *TimeAndSale) TimeAndSaleType() Type {
	return TypeValueOf(mathutil.GetBits(int64(t.flags), typeMask, typeShift))
}

func (t *TimeAndSale) SetTimeAndSaleType(value Type) {
	t.SetFlags(int32(mathutil.SetBits(int64(t.flags), typeMask, typeShift, int64(value))))
}

func (t *TimeAndSale) IsNew() bool {
	return t.TimeAndSaleType() == TypeNew
}

func (t *TimeAndSale) IsCorrection() bool {
	return t.TimeAndSaleType() == TypeCorrection
}

func (t *TimeAndSale) IsCancel() bool {
	return t.TimeAndSaleType() == TypeCancel
}

func (t *TimeAndSale) SetFlags(value int32) {
	t.flags = value
}

func (t *TimeAndSale) Flags() int32 {
	return t.flags
}

func (t *TimeAndSale) String() string {
	return "TimeAndSale{" + formatutil.FormatString(t.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(t.EventTime()) +
		", eventFlags=" + formatutil.HexFormat(int64(t.EventFlags())) +
		", time=" + formatutil.FormatTime(t.Time()) +
		", timeNanoPart=" + strconv.FormatInt(int64(t.timeNanoPart), 10) +
		", sequence=" + strconv.FormatInt(int64(t.Sequence()), 10) +
		", exchange=" + formatutil.FormatChar(rune(t.exchangeCode)) +
		", price=" + formatutil.FormatFloat64(t.Price()) +
		", size=" + formatutil.FormatFloat64(t.Size()) +
		", bid=" + formatutil.FormatFloat64(t.BidPrice()) +
		", ask=" + formatutil.FormatFloat64(t.AskPrice()) +
		", ESC=" + formatutil.FormatString(t.ExchangeSaleConditions()) +
		", TTE=" + formatutil.FormatChar(t.TradeThroughExempt()) +
		", side=" + formatutil.FormatInt64(int64(t.AggressorSide())) +
		", spread=" + formatutil.FormatBool(t.IsSpreadLeg()) +
		", ETH=" + formatutil.FormatBool(t.IsExtendedTradingHours()) +
		", validTick=" + formatutil.FormatBool(t.IsValidTick()) +
		", type=" + formatutil.FormatInt64(int64(t.TimeAndSaleType())) +
		formatutil.FormatNullableString(", buyer=%s", t.Buyer(), "") +
		formatutil.FormatNullableString(", seller=%s", t.Seller(), "") +
		"}"
}
