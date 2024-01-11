package profile

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
	"math"
)

const (
	ssrMask     = 3
	ssrShift    = 2
	statusMask  = 3
	statusShift = 0
)

type Profile struct {
	eventSymbol       *string
	eventTime         int64
	description       *string
	statusReason      *string
	haltStartTime     int64
	haltEndTime       int64
	highLimitPrice    float64
	lowLimitPrice     float64
	high52WeekPrice   float64
	low52WeekPrice    float64
	beta              float64
	earningsPerShare  float64
	dividendFrequency float64
	exDividendAmount  float64
	exDividendDayId   int32
	shares            float64
	freeFloat         float64

	flags int32
}

func (p *Profile) ShortSaleRestriction() ShortSaleRestriction {
	return ShortSaleRestrictionValueOf(mathutil.GetBits(int64(p.flags), ssrMask, ssrShift))
}

func (p *Profile) SetShortSaleRestriction(shortSaleRestriction ShortSaleRestriction) {
	p.SetFlags(int32(mathutil.SetBits(int64(p.flags), ssrMask, ssrShift, int64(shortSaleRestriction))))
}

func (p *Profile) IsShortSaleRestricted() bool {
	return p.TradingStatus() == ActiveShortSaleRestriction
}

func (p *Profile) TradingStatus() TradingStatus {
	return TradingStatusValueOf(mathutil.GetBits(int64(p.flags), statusMask, statusShift))
}

func (p *Profile) SetTradingStatus(tradingStatus TradingStatus) {
	p.SetFlags(int32(mathutil.SetBits(int64(p.flags), statusMask, statusShift, int64(tradingStatus))))
}

func (p *Profile) IsTradingHalted() bool {
	return p.TradingStatus() == HaltedTradingStatus
}

func (p *Profile) EventSymbol() *string {
	return p.eventSymbol
}

func (p *Profile) SetEventSymbol(eventSymbol *string) {
	p.eventSymbol = eventSymbol
}

func (p *Profile) EventTime() int64 {
	return p.eventTime
}

func (p *Profile) SetEventTime(eventTime int64) {
	p.eventTime = eventTime
}

func (p *Profile) Description() *string {
	return p.description
}

func (p *Profile) SetDescription(description *string) {
	p.description = description
}

func (p *Profile) StatusReason() *string {
	return p.statusReason
}

func (p *Profile) SetStatusReason(statusReason *string) {
	p.statusReason = statusReason
}

func (p *Profile) HaltStartTime() int64 {
	return p.haltStartTime
}

func (p *Profile) SetHaltStartTime(haltStartTime int64) {
	p.haltStartTime = haltStartTime
}

func (p *Profile) HaltEndTime() int64 {
	return p.haltEndTime
}

func (p *Profile) SetHaltEndTime(haltEndTime int64) {
	p.haltEndTime = haltEndTime
}

func (p *Profile) HighLimitPrice() float64 {
	return p.highLimitPrice
}

func (p *Profile) SetHighLimitPrice(highLimitPrice float64) {
	p.highLimitPrice = highLimitPrice
}

func (p *Profile) LowLimitPrice() float64 {
	return p.lowLimitPrice
}

func (p *Profile) SetLowLimitPrice(lowLimitPrice float64) {
	p.lowLimitPrice = lowLimitPrice
}

func (p *Profile) High52WeekPrice() float64 {
	return p.high52WeekPrice
}

func (p *Profile) SetHigh52WeekPrice(high52WeekPrice float64) {
	p.high52WeekPrice = high52WeekPrice
}

func (p *Profile) Low52WeekPrice() float64 {
	return p.low52WeekPrice
}

func (p *Profile) SetLow52WeekPrice(low52WeekPrice float64) {
	p.low52WeekPrice = low52WeekPrice
}

func (p *Profile) Beta() float64 {
	return p.beta
}

func (p *Profile) SetBeta(beta float64) {
	p.beta = beta
}

func (p *Profile) EarningsPerShare() float64 {
	return p.earningsPerShare
}

func (p *Profile) SetEarningsPerShare(earningsPerShare float64) {
	p.earningsPerShare = earningsPerShare
}

func (p *Profile) DividendFrequency() float64 {
	return p.dividendFrequency
}

func (p *Profile) SetDividendFrequency(dividendFrequency float64) {
	p.dividendFrequency = dividendFrequency
}

func (p *Profile) ExDividendAmount() float64 {
	return p.exDividendAmount
}

func (p *Profile) SetExDividendAmount(exDividendAmount float64) {
	p.exDividendAmount = exDividendAmount
}

func (p *Profile) ExDividendDayId() int32 {
	return p.exDividendDayId
}

func (p *Profile) SetExDividendDayId(exDividendDayId int32) {
	p.exDividendDayId = exDividendDayId
}

func (p *Profile) Shares() float64 {
	return p.shares
}

func (p *Profile) SetShares(shares float64) {
	p.shares = shares
}

func (p *Profile) FreeFloat() float64 {
	return p.freeFloat
}

func (p *Profile) SetFreeFloat(freeFloat float64) {
	p.freeFloat = freeFloat
}

func (p *Profile) SetFlags(flags int32) {
	p.flags = flags
}

func (p *Profile) Flags() int32 {
	return p.flags
}

func NewProfile(eventSymbol string) *Profile {
	return &Profile{eventSymbol: &eventSymbol,
		highLimitPrice:    math.NaN(),
		lowLimitPrice:     math.NaN(),
		high52WeekPrice:   math.NaN(),
		low52WeekPrice:    math.NaN(),
		beta:              math.NaN(),
		earningsPerShare:  math.NaN(),
		dividendFrequency: math.NaN(),
		exDividendAmount:  math.NaN(),
		shares:            math.NaN(),
		freeFloat:         math.NaN(),
	}
}

func (p *Profile) String() string {
	return "Profile{" + formatutil.FormatString(p.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(p.EventTime()) +
		", description='" + formatutil.FormatString(p.Description()) + "'" +
		", SSR=" + formatutil.FormatInt64(int64(p.ShortSaleRestriction())) +
		", status=" + formatutil.FormatInt64(int64(p.TradingStatus())) +
		", statusReason='" + formatutil.FormatString(p.StatusReason()) + "'" +
		", haltStartTime=" + formatutil.FormatInt64(p.HaltStartTime()) +
		", haltEndTime=" + formatutil.FormatInt64(p.HaltEndTime()) +
		", highLimitPrice=" + formatutil.FormatFloat64(p.HighLimitPrice()) +
		", lowLimitPrice=" + formatutil.FormatFloat64(p.LowLimitPrice()) +
		", high52WeekPrice=" + formatutil.FormatFloat64(p.High52WeekPrice()) +
		", low52WeekPrice=" + formatutil.FormatFloat64(p.Low52WeekPrice()) +
		", beta=" + formatutil.FormatFloat64(p.Beta()) +
		", earningsPerShare=" + formatutil.FormatFloat64(p.EarningsPerShare()) +
		", dividendFrequency=" + formatutil.FormatFloat64(p.DividendFrequency()) +
		", exDividendAmount=" + formatutil.FormatFloat64(p.ExDividendAmount()) +
		", exDividendDay=" + formatutil.FormatInt64(int64(p.ExDividendDayId())) +
		", shares=" + formatutil.FormatFloat64(p.Shares()) +
		", freeFloat=" + formatutil.FormatFloat64(p.FreeFloat()) +
		"}"
}
