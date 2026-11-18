package profile

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"math"
)

const (
	ssrMask     = 3
	ssrShift    = 2
	statusMask  = 3
	statusShift = 0
)

// Profile is a snapshot of security information such as the description, trading status and price limits.
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

// Type returns the event code of this event.
func (p *Profile) Type() eventcodes.EventCode {
	return eventcodes.Profile
}

// ShortSaleRestriction returns the short sale restriction status of the security instrument.
func (p *Profile) ShortSaleRestriction() ShortSaleRestriction {
	return ShortSaleRestrictionValueOf(mathutil.GetBits(int64(p.flags), ssrMask, ssrShift))
}

// SetShortSaleRestriction sets the short sale restriction status of the security instrument.
func (p *Profile) SetShortSaleRestriction(shortSaleRestriction ShortSaleRestriction) {
	p.SetFlags(int32(mathutil.SetBits(int64(p.flags), ssrMask, ssrShift, int64(shortSaleRestriction))))
}

// IsShortSaleRestricted reports whether short sale of the security instrument is restricted.
func (p *Profile) IsShortSaleRestricted() bool {
	return p.ShortSaleRestriction() == ActiveShortSaleRestriction
}

// TradingStatus returns the current trading status of the security instrument.
func (p *Profile) TradingStatus() TradingStatus {
	return TradingStatusValueOf(mathutil.GetBits(int64(p.flags), statusMask, statusShift))
}

// SetTradingStatus sets the current trading status of the security instrument.
func (p *Profile) SetTradingStatus(tradingStatus TradingStatus) {
	p.SetFlags(int32(mathutil.SetBits(int64(p.flags), statusMask, statusShift, int64(tradingStatus))))
}

// IsTradingHalted reports whether trading of the security instrument is halted.
func (p *Profile) IsTradingHalted() bool {
	return p.TradingStatus() == HaltedTradingStatus
}

// EventSymbol returns the symbol of this event.
func (p *Profile) EventSymbol() *string {
	return p.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (p *Profile) SetEventSymbol(eventSymbol *string) {
	p.eventSymbol = eventSymbol
}

// EventTime returns the source-provided time of this event in milliseconds since the Unix epoch.
func (p *Profile) EventTime() int64 {
	return p.eventTime
}

// SetEventTime sets the source-provided time of this event in milliseconds since the Unix epoch.
func (p *Profile) SetEventTime(eventTime int64) {
	p.eventTime = eventTime
}

// Description returns the description of the security instrument.
func (p *Profile) Description() *string {
	return p.description
}

// SetDescription sets the description of the security instrument.
func (p *Profile) SetDescription(description *string) {
	p.description = description
}

// StatusReason returns the description of the reason that trading is halted at the moment.
func (p *Profile) StatusReason() *string {
	return p.statusReason
}

// SetStatusReason sets the description of the reason that trading is halted at the moment.
func (p *Profile) SetStatusReason(statusReason *string) {
	p.statusReason = statusReason
}

// HaltStartTime returns the starting time of the trading halt in milliseconds since the Unix epoch.
func (p *Profile) HaltStartTime() int64 {
	return p.haltStartTime
}

// SetHaltStartTime sets the starting time of the trading halt in milliseconds since the Unix epoch.
func (p *Profile) SetHaltStartTime(haltStartTime int64) {
	p.haltStartTime = haltStartTime
}

// HaltEndTime returns the ending time of the trading halt in milliseconds since the Unix epoch.
func (p *Profile) HaltEndTime() int64 {
	return p.haltEndTime
}

// SetHaltEndTime sets the ending time of the trading halt in milliseconds since the Unix epoch.
func (p *Profile) SetHaltEndTime(haltEndTime int64) {
	p.haltEndTime = haltEndTime
}

// HighLimitPrice returns the maximal (high) allowed price.
func (p *Profile) HighLimitPrice() float64 {
	return p.highLimitPrice
}

// SetHighLimitPrice sets the maximal (high) allowed price.
func (p *Profile) SetHighLimitPrice(highLimitPrice float64) {
	p.highLimitPrice = highLimitPrice
}

// LowLimitPrice returns the minimal (low) allowed price.
func (p *Profile) LowLimitPrice() float64 {
	return p.lowLimitPrice
}

// SetLowLimitPrice sets the minimal (low) allowed price.
func (p *Profile) SetLowLimitPrice(lowLimitPrice float64) {
	p.lowLimitPrice = lowLimitPrice
}

// High52WeekPrice returns the maximal (high) price in the last 52 weeks.
func (p *Profile) High52WeekPrice() float64 {
	return p.high52WeekPrice
}

// SetHigh52WeekPrice sets the maximal (high) price in the last 52 weeks.
func (p *Profile) SetHigh52WeekPrice(high52WeekPrice float64) {
	p.high52WeekPrice = high52WeekPrice
}

// Low52WeekPrice returns the minimal (low) price in the last 52 weeks.
func (p *Profile) Low52WeekPrice() float64 {
	return p.low52WeekPrice
}

// SetLow52WeekPrice sets the minimal (low) price in the last 52 weeks.
func (p *Profile) SetLow52WeekPrice(low52WeekPrice float64) {
	p.low52WeekPrice = low52WeekPrice
}

// Beta returns the correlation coefficient of the instrument to the S&P500 index.
func (p *Profile) Beta() float64 {
	return p.beta
}

// SetBeta sets the correlation coefficient of the instrument to the S&P500 index.
func (p *Profile) SetBeta(beta float64) {
	p.beta = beta
}

// EarningsPerShare returns the earnings per share (the company earnings divided by the number of shares).
func (p *Profile) EarningsPerShare() float64 {
	return p.earningsPerShare
}

// SetEarningsPerShare sets the earnings per share.
func (p *Profile) SetEarningsPerShare(earningsPerShare float64) {
	p.earningsPerShare = earningsPerShare
}

// DividendFrequency returns the frequency of cash dividends payments per year (calculated).
func (p *Profile) DividendFrequency() float64 {
	return p.dividendFrequency
}

// SetDividendFrequency sets the frequency of cash dividends payments per year.
func (p *Profile) SetDividendFrequency(dividendFrequency float64) {
	p.dividendFrequency = dividendFrequency
}

// ExDividendAmount returns the amount of the most recent regular cash dividend paid.
func (p *Profile) ExDividendAmount() float64 {
	return p.exDividendAmount
}

// SetExDividendAmount sets the amount of the most recent regular cash dividend paid.
func (p *Profile) SetExDividendAmount(exDividendAmount float64) {
	p.exDividendAmount = exDividendAmount
}

// ExDividendDayId returns the identifier of the day of the most recent ex-dividend date.
func (p *Profile) ExDividendDayId() int32 {
	return p.exDividendDayId
}

// SetExDividendDayId sets the identifier of the day of the most recent ex-dividend date.
func (p *Profile) SetExDividendDayId(exDividendDayId int32) {
	p.exDividendDayId = exDividendDayId
}

// Shares returns the number of shares outstanding.
func (p *Profile) Shares() float64 {
	return p.shares
}

// SetShares sets the number of shares outstanding.
func (p *Profile) SetShares(shares float64) {
	p.shares = shares
}

// FreeFloat returns the number of shares available for trading by the public.
func (p *Profile) FreeFloat() float64 {
	return p.freeFloat
}

// SetFreeFloat sets the number of shares available for trading by the public.
func (p *Profile) SetFreeFloat(freeFloat float64) {
	p.freeFloat = freeFloat
}

// SetFlags sets the packed flags of this event.
func (p *Profile) SetFlags(flags int32) {
	p.flags = flags
}

// Flags returns the packed flags of this event.
func (p *Profile) Flags() int32 {
	return p.flags
}

// NewProfile creates a new Profile for the given event symbol.
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

// String returns a string representation of this event.
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
