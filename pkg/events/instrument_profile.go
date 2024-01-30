package events

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
)

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
}

func (p *InstrumentProfile) InstrumentType() *string {
	return p.instrumentType
}

func (p *InstrumentProfile) SetInstrumentType(instrumentType *string) {
	p.instrumentType = instrumentType
}

func (p *InstrumentProfile) Symbol() *string {
	return p.symbol
}

func (p *InstrumentProfile) SetSymbol(symbol *string) {
	p.symbol = symbol
}

func (p *InstrumentProfile) Description() *string {
	return p.description
}

func (p *InstrumentProfile) SetDescription(description *string) {
	p.description = description
}

func (p *InstrumentProfile) LocalSymbol() *string {
	return p.localSymbol
}

func (p *InstrumentProfile) SetLocalSymbol(localSymbol *string) {
	p.localSymbol = localSymbol
}

func (p *InstrumentProfile) LocalDescription() *string {
	return p.localDescription
}

func (p *InstrumentProfile) SetLocalDescription(localDescription *string) {
	p.localDescription = localDescription
}

func (p *InstrumentProfile) Country() *string {
	return p.country
}

func (p *InstrumentProfile) SetCountry(country *string) {
	p.country = country
}

func (p *InstrumentProfile) Opol() *string {
	return p.opol
}

func (p *InstrumentProfile) SetOpol(opol *string) {
	p.opol = opol
}

func (p *InstrumentProfile) ExchangeData() *string {
	return p.exchangeData
}

func (p *InstrumentProfile) SetExchangeData(exchangeData *string) {
	p.exchangeData = exchangeData
}

func (p *InstrumentProfile) Exchanges() *string {
	return p.exchanges
}

func (p *InstrumentProfile) SetExchanges(exchanges *string) {
	p.exchanges = exchanges
}

func (p *InstrumentProfile) Currency() *string {
	return p.currency
}

func (p *InstrumentProfile) SetCurrency(currency *string) {
	p.currency = currency
}

func (p *InstrumentProfile) BaseCurrency() *string {
	return p.baseCurrency
}

func (p *InstrumentProfile) SetBaseCurrency(baseCurrency *string) {
	p.baseCurrency = baseCurrency
}

func (p *InstrumentProfile) Cfi() *string {
	return p.cfi
}

func (p *InstrumentProfile) SetCfi(cfi *string) {
	p.cfi = cfi
}

func (p *InstrumentProfile) Isin() *string {
	return p.isin
}

func (p *InstrumentProfile) SetIsin(isin *string) {
	p.isin = isin
}

func (p *InstrumentProfile) Sedol() *string {
	return p.sedol
}

func (p *InstrumentProfile) SetSedol(sedol *string) {
	p.sedol = sedol
}

func (p *InstrumentProfile) Cusip() *string {
	return p.cusip
}

func (p *InstrumentProfile) SetCusip(cusip *string) {
	p.cusip = cusip
}

func (p *InstrumentProfile) Icb() int64 {
	return p.icb
}

func (p *InstrumentProfile) SetIcb(icb int64) {
	p.icb = icb
}

func (p *InstrumentProfile) Sic() int64 {
	return p.sic
}

func (p *InstrumentProfile) SetSic(sic int64) {
	p.sic = sic
}

func (p *InstrumentProfile) Multiplier() float64 {
	return p.multiplier
}

func (p *InstrumentProfile) SetMultiplier(multiplier float64) {
	p.multiplier = multiplier
}

func (p *InstrumentProfile) Product() *string {
	return p.product
}

func (p *InstrumentProfile) SetProduct(product *string) {
	p.product = product
}

func (p *InstrumentProfile) Underlying() *string {
	return p.underlying
}

func (p *InstrumentProfile) SetUnderlying(underlying *string) {
	p.underlying = underlying
}

func (p *InstrumentProfile) Spc() float64 {
	return p.spc
}

func (p *InstrumentProfile) SetSpc(spc float64) {
	p.spc = spc
}

func (p *InstrumentProfile) AdditionalUnderlyings() *string {
	return p.additionalUnderlyings
}

func (p *InstrumentProfile) SetAdditionalUnderlyings(additionalUnderlyings *string) {
	p.additionalUnderlyings = additionalUnderlyings
}

func (p *InstrumentProfile) Mmy() *string {
	return p.mmy
}

func (p *InstrumentProfile) SetMmy(mmy *string) {
	p.mmy = mmy
}

func (p *InstrumentProfile) Expiration() int64 {
	return p.expiration
}

func (p *InstrumentProfile) SetExpiration(expiration int64) {
	p.expiration = expiration
}

func (p *InstrumentProfile) LastTrade() int64 {
	return p.lastTrade
}

func (p *InstrumentProfile) SetLastTrade(lastTrade int64) {
	p.lastTrade = lastTrade
}

func (p *InstrumentProfile) Strike() float64 {
	return p.strike
}

func (p *InstrumentProfile) SetStrike(strike float64) {
	p.strike = strike
}

func (p *InstrumentProfile) OptionType() *string {
	return p.optionType
}

func (p *InstrumentProfile) SetOptionType(optionType *string) {
	p.optionType = optionType
}

func (p *InstrumentProfile) ExpirationStyle() *string {
	return p.expirationStyle
}

func (p *InstrumentProfile) SetExpirationStyle(expirationStyle *string) {
	p.expirationStyle = expirationStyle
}

func (p *InstrumentProfile) SettlementStyle() *string {
	return p.settlementStyle
}

func (p *InstrumentProfile) SetSettlementStyle(settlementStyle *string) {
	p.settlementStyle = settlementStyle
}

func (p *InstrumentProfile) PriceIncrements() *string {
	return p.priceIncrements
}

func (p *InstrumentProfile) SetPriceIncrements(priceIncrements *string) {
	p.priceIncrements = priceIncrements
}

func (p *InstrumentProfile) TradingHours() *string {
	return p.tradingHours
}

func (p *InstrumentProfile) SetTradingHours(tradingHours *string) {
	p.tradingHours = tradingHours
}

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
	}
}

func (q *InstrumentProfile) String() string {
	return "InstrumentProfile{" + formatutil.FormatString(q.Symbol()) +
		", description=" + formatutil.FormatString(q.Description()) +
		", type=" + formatutil.FormatString(q.InstrumentType()) +
		"}"
}
