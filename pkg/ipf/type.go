package ipf

// The types of instruments, the values of the TYPE field of an instrument profile
// (events.InstrumentProfile.InstrumentType).
const (
	TypeCurrency        = "CURRENCY"
	TypeForex           = "FOREX"
	TypeBond            = "BOND"
	TypeIndex           = "INDEX"
	TypeStock           = "STOCK"
	TypeETF             = "ETF"
	TypeMutualFund      = "MUTUAL_FUND"
	TypeMoneyMarketFund = "MONEY_MARKET_FUND"
	TypeProduct         = "PRODUCT"
	TypeFuture          = "FUTURE"
	TypeOption          = "OPTION"
	TypeWarrant         = "WARRANT"
	TypeCertificate     = "CERTIFICATE"
	TypeCFD             = "CFD"
	TypeSpread          = "SPREAD"
	TypeOther           = "OTHER"
	// TypeRemoved marks a profile that is removed: a live connection or a collector update with this type removes
	// the profile of the symbol.
	TypeRemoved = "REMOVED"
)
