// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package events

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
)

// InstrumentProfile is the description of a financial instrument, including its identifiers, classification and trading data.
type InstrumentProfile struct {
	instrumentType        *string
	symbol                *string
	description           *string
	localSymbol           *string
	localDescription      *string
	country               *string
	opol                  *string
	exchangeData          *string
	exchanges             *string
	currency              *string
	baseCurrency          *string
	cfi                   *string
	isin                  *string
	sedol                 *string
	cusip                 *string
	icb                   int64
	sic                   int64
	multiplier            float64
	product               *string
	underlying            *string
	spc                   float64
	additionalUnderlyings *string
	mmy                   *string
	expiration            int64
	lastTrade             int64
	strike                float64
	optionType            *string
	expirationStyle       *string
	settlementStyle       *string
	priceIncrements       *string
	tradingHours          *string
	customFields          map[string]string
}

// InstrumentType returns the type of the instrument (for example STOCK, OPTION or FUTURE).
func (p *InstrumentProfile) InstrumentType() *string {
	return p.instrumentType
}

// SetInstrumentType sets the type of the instrument (for example STOCK, OPTION or FUTURE).
func (p *InstrumentProfile) SetInstrumentType(instrumentType *string) {
	p.instrumentType = instrumentType
}

// Symbol returns the identifier of the instrument, unique within its type.
func (p *InstrumentProfile) Symbol() *string {
	return p.symbol
}

// SetSymbol sets the identifier of the instrument, unique within its type.
func (p *InstrumentProfile) SetSymbol(symbol *string) {
	p.symbol = symbol
}

// Description returns the human-readable description of the instrument.
func (p *InstrumentProfile) Description() *string {
	return p.description
}

// SetDescription sets the human-readable description of the instrument.
func (p *InstrumentProfile) SetDescription(description *string) {
	p.description = description
}

// LocalSymbol returns the identifier of the instrument on its local exchange.
func (p *InstrumentProfile) LocalSymbol() *string {
	return p.localSymbol
}

// SetLocalSymbol sets the identifier of the instrument on its local exchange.
func (p *InstrumentProfile) SetLocalSymbol(localSymbol *string) {
	p.localSymbol = localSymbol
}

// LocalDescription returns the description of the instrument in the local language.
func (p *InstrumentProfile) LocalDescription() *string {
	return p.localDescription
}

// SetLocalDescription sets the description of the instrument in the local language.
func (p *InstrumentProfile) SetLocalDescription(localDescription *string) {
	p.localDescription = localDescription
}

// Country returns the ISO 3166-1 alpha-2 country code of the instrument.
func (p *InstrumentProfile) Country() *string {
	return p.country
}

// SetCountry sets the ISO 3166-1 alpha-2 country code of the instrument.
func (p *InstrumentProfile) SetCountry(country *string) {
	p.country = country
}

// Opol returns the ISO 10383 market identification code (MIC) of the primary listing exchange.
func (p *InstrumentProfile) Opol() *string {
	return p.opol
}

// SetOpol sets the ISO 10383 market identification code (MIC) of the primary listing exchange.
func (p *InstrumentProfile) SetOpol(opol *string) {
	p.opol = opol
}

// ExchangeData returns the exchange-specific data of the instrument.
func (p *InstrumentProfile) ExchangeData() *string {
	return p.exchangeData
}

// SetExchangeData sets the exchange-specific data of the instrument.
func (p *InstrumentProfile) SetExchangeData(exchangeData *string) {
	p.exchangeData = exchangeData
}

// Exchanges returns the list of exchanges where the instrument is traded.
func (p *InstrumentProfile) Exchanges() *string {
	return p.exchanges
}

// SetExchanges sets the list of exchanges where the instrument is traded.
func (p *InstrumentProfile) SetExchanges(exchanges *string) {
	p.exchanges = exchanges
}

// Currency returns the trading currency of the instrument.
func (p *InstrumentProfile) Currency() *string {
	return p.currency
}

// SetCurrency sets the trading currency of the instrument.
func (p *InstrumentProfile) SetCurrency(currency *string) {
	p.currency = currency
}

// BaseCurrency returns the base currency of the instrument.
func (p *InstrumentProfile) BaseCurrency() *string {
	return p.baseCurrency
}

// SetBaseCurrency sets the base currency of the instrument.
func (p *InstrumentProfile) SetBaseCurrency(baseCurrency *string) {
	p.baseCurrency = baseCurrency
}

// Cfi returns the ISO 10962 Classification of Financial Instruments (CFI) code.
func (p *InstrumentProfile) Cfi() *string {
	return p.cfi
}

// SetCfi sets the ISO 10962 Classification of Financial Instruments (CFI) code.
func (p *InstrumentProfile) SetCfi(cfi *string) {
	p.cfi = cfi
}

// Isin returns the International Securities Identification Number (ISIN).
func (p *InstrumentProfile) Isin() *string {
	return p.isin
}

// SetIsin sets the International Securities Identification Number (ISIN).
func (p *InstrumentProfile) SetIsin(isin *string) {
	p.isin = isin
}

// Sedol returns the Stock Exchange Daily Official List (SEDOL) number.
func (p *InstrumentProfile) Sedol() *string {
	return p.sedol
}

// SetSedol sets the Stock Exchange Daily Official List (SEDOL) number.
func (p *InstrumentProfile) SetSedol(sedol *string) {
	p.sedol = sedol
}

// Cusip returns the Committee on Uniform Security Identification Procedures (CUSIP) number.
func (p *InstrumentProfile) Cusip() *string {
	return p.cusip
}

// SetCusip sets the Committee on Uniform Security Identification Procedures (CUSIP) number.
func (p *InstrumentProfile) SetCusip(cusip *string) {
	p.cusip = cusip
}

// Icb returns the Industry Classification Benchmark (ICB) code.
func (p *InstrumentProfile) Icb() int64 {
	return p.icb
}

// SetIcb sets the Industry Classification Benchmark (ICB) code.
func (p *InstrumentProfile) SetIcb(icb int64) {
	p.icb = icb
}

// Sic returns the Standard Industrial Classification (SIC) code.
func (p *InstrumentProfile) Sic() int64 {
	return p.sic
}

// SetSic sets the Standard Industrial Classification (SIC) code.
func (p *InstrumentProfile) SetSic(sic int64) {
	p.sic = sic
}

// Multiplier returns the market value multiplier of the instrument.
func (p *InstrumentProfile) Multiplier() float64 {
	return p.multiplier
}

// SetMultiplier sets the market value multiplier of the instrument.
func (p *InstrumentProfile) SetMultiplier(multiplier float64) {
	p.multiplier = multiplier
}

// Product returns the product for futures and options on futures.
func (p *InstrumentProfile) Product() *string {
	return p.product
}

// SetProduct sets the product for futures and options on futures.
func (p *InstrumentProfile) SetProduct(product *string) {
	p.product = product
}

// Underlying returns the primary underlying symbol for options.
func (p *InstrumentProfile) Underlying() *string {
	return p.underlying
}

// SetUnderlying sets the primary underlying symbol for options.
func (p *InstrumentProfile) SetUnderlying(underlying *string) {
	p.underlying = underlying
}

// Spc returns the shares per contract for options.
func (p *InstrumentProfile) Spc() float64 {
	return p.spc
}

// SetSpc sets the shares per contract for options.
func (p *InstrumentProfile) SetSpc(spc float64) {
	p.spc = spc
}

// AdditionalUnderlyings returns the additional underlyings for options, including those with non-standard deliverables.
func (p *InstrumentProfile) AdditionalUnderlyings() *string {
	return p.additionalUnderlyings
}

// SetAdditionalUnderlyings sets the additional underlyings for options, including those with non-standard deliverables.
func (p *InstrumentProfile) SetAdditionalUnderlyings(additionalUnderlyings *string) {
	p.additionalUnderlyings = additionalUnderlyings
}

// Mmy returns the maturity month-year of the instrument.
func (p *InstrumentProfile) Mmy() *string {
	return p.mmy
}

// SetMmy sets the maturity month-year of the instrument.
func (p *InstrumentProfile) SetMmy(mmy *string) {
	p.mmy = mmy
}

// Expiration returns the day of expiration, as a day number since the Unix epoch.
func (p *InstrumentProfile) Expiration() int64 {
	return p.expiration
}

// SetExpiration sets the day of expiration, as a day number since the Unix epoch.
func (p *InstrumentProfile) SetExpiration(expiration int64) {
	p.expiration = expiration
}

// LastTrade returns the last trading day, as a day number since the Unix epoch.
func (p *InstrumentProfile) LastTrade() int64 {
	return p.lastTrade
}

// SetLastTrade sets the last trading day, as a day number since the Unix epoch.
func (p *InstrumentProfile) SetLastTrade(lastTrade int64) {
	p.lastTrade = lastTrade
}

// Strike returns the strike price of the option.
func (p *InstrumentProfile) Strike() float64 {
	return p.strike
}

// SetStrike sets the strike price of the option.
func (p *InstrumentProfile) SetStrike(strike float64) {
	p.strike = strike
}

// OptionType returns the type of the option.
func (p *InstrumentProfile) OptionType() *string {
	return p.optionType
}

// SetOptionType sets the type of the option.
func (p *InstrumentProfile) SetOptionType(optionType *string) {
	p.optionType = optionType
}

// ExpirationStyle returns the expiration cycle style of the option.
func (p *InstrumentProfile) ExpirationStyle() *string {
	return p.expirationStyle
}

// SetExpirationStyle sets the expiration cycle style of the option.
func (p *InstrumentProfile) SetExpirationStyle(expirationStyle *string) {
	p.expirationStyle = expirationStyle
}

// SettlementStyle returns the settlement style of the option.
func (p *InstrumentProfile) SettlementStyle() *string {
	return p.settlementStyle
}

// SetSettlementStyle sets the settlement style of the option.
func (p *InstrumentProfile) SetSettlementStyle(settlementStyle *string) {
	p.settlementStyle = settlementStyle
}

// PriceIncrements returns the price increments (tick sizes) of the instrument.
func (p *InstrumentProfile) PriceIncrements() *string {
	return p.priceIncrements
}

// SetPriceIncrements sets the price increments (tick sizes) of the instrument.
func (p *InstrumentProfile) SetPriceIncrements(priceIncrements *string) {
	p.priceIncrements = priceIncrements
}

// TradingHours returns the trading hours specification of the instrument.
func (p *InstrumentProfile) TradingHours() *string {
	return p.tradingHours
}

// SetTradingHours sets the trading hours specification of the instrument.
func (p *InstrumentProfile) SetTradingHours(tradingHours *string) {
	p.tradingHours = tradingHours
}

// CustomFields returns the custom (non-standard) fields of the instrument profile.
func (p *InstrumentProfile) CustomFields() map[string]string {
	return p.customFields
}

// SetCustomFields sets the custom (non-standard) fields of the instrument profile.
func (p *InstrumentProfile) SetCustomFields(customFields map[string]string) {
	p.customFields = customFields
}

// NewInstrumentProfile creates an empty InstrumentProfile.
func NewInstrumentProfile() *InstrumentProfile {
	emptyValue := ""
	emptyString := &emptyValue

	return &InstrumentProfile{
		instrumentType:        emptyString,
		symbol:                emptyString,
		description:           emptyString,
		localSymbol:           emptyString,
		localDescription:      emptyString,
		country:               emptyString,
		opol:                  emptyString,
		exchangeData:          emptyString,
		exchanges:             emptyString,
		currency:              emptyString,
		baseCurrency:          emptyString,
		cfi:                   emptyString,
		isin:                  emptyString,
		sedol:                 emptyString,
		cusip:                 emptyString,
		icb:                   0,
		sic:                   0,
		multiplier:            0,
		product:               emptyString,
		underlying:            emptyString,
		spc:                   0,
		additionalUnderlyings: emptyString,
		mmy:                   emptyString,
		expiration:            0,
		lastTrade:             0,
		strike:                0,
		optionType:            emptyString,
		expirationStyle:       emptyString,
		settlementStyle:       emptyString,
		priceIncrements:       emptyString,
		tradingHours:          emptyString,
		customFields:          make(map[string]string),
	}
}

// String returns a textual representation of the instrument profile.
func (q *InstrumentProfile) String() string {
	return "InstrumentProfile{" + formatutil.FormatString(q.Symbol()) +
		", description=" + formatutil.FormatString(q.Description()) +
		", type=" + formatutil.FormatString(q.InstrumentType()) +
		"}"
}
