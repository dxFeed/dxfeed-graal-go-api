package order

var (
	defaultSource = newOrderSourceWithIDNameFlags(0, "DEFAULT", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook)

	/// Bid side of a composite ``Quote``
	///
	/// It is a synthetic source.
	/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
	compsoiteBid = newOrderSourceWithIDNameFlags(1, "COMPOSITE_BID", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook)
	/// Ask side of a composite ``Quote``.
	/// It is a synthetic source.
	/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
	compsoiteAsk = newOrderSourceWithIDNameFlags(2, "COMPOSITE_ASK", 0)
	/// Bid side of a regional ``Quote``.
	/// It is a synthetic source.
	/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

	regionalBid = newOrderSourceWithIDNameFlags(3, "REGIONAL_BID", 0)
	/// Ask side of a regional ``Quote``.
	/// It is a synthetic source.
	/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

	regionalAsk = newOrderSourceWithIDNameFlags(4, "REGIONAL_ASK", 0)
	/// Bid side of an aggregate order book (futures depth and NASDAQ Level II).
	/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
	agregateBid = newOrderSourceWithIDNameFlags(5, "AGGREGATE_BID", 0)

	/// Ask side of an aggregate order book (futures depth and NASDAQ Level II).
	/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
	agregateAsk = newOrderSourceWithIDNameFlags(6, "AGGREGATE_ASK", 0)

	/// NASDAQ Total View.
	ntvL3 = newOrderSourceName("NTV", pubOrder|fullOrderBook)

	/// NASDAQ Total View. Record for price level book.
	ntvL2 = newOrderSourceName("ntv", pubOrder)

	/// NASDAQ Futures Exchange.
	nfxL3 = newOrderSourceName("NFX", pubOrder)

	/// NASDAQ eSpeed.
	espdL3 = newOrderSourceName("ESPD", pubOrder)

	/// NASDAQ Fixed Income.
	xnfiL3 = newOrderSourceName("XNFI", pubOrder)

	/// Intercontinental Exchange.
	iceL3 = newOrderSourceName("ICE", pubOrder)

	/// International Securities Exchange.
	iseL3 = newOrderSourceName("ISE", pubOrder|pubSpreadOrder)

	/// Direct-Edge EDGA Exchange.
	deaL3 = newOrderSourceName("DEA", pubOrder)

	/// Direct-Edge EDGX Exchange.
	dexL3 = newOrderSourceName("DEX", pubOrder)

	/// Bats BYX Exchange.
	byxL3 = newOrderSourceName("BYX", pubOrder)

	/// Bats BZX Exchange.
	bzxL3 = newOrderSourceName("BZX", pubOrder)

	/// Bats Europe BXE Exchange.
	bateL3 = newOrderSourceName("BATE", pubOrder)

	/// Bats Europe CXE Exchange.
	chixL3 = newOrderSourceName("CHIX", pubOrder)

	/// Bats Europe DXE Exchange.
	ceuxL3 = newOrderSourceName("CEUX", pubOrder)

	/// Bats Europe TRF.
	bxtrL3 = newOrderSourceName("BXTR", pubOrder)

	/// Borsa Istanbul Exchange.
	istL3 = newOrderSourceName("IST", pubOrder)

	/// Borsa Istanbul Exchange. Record for particular top 20 order book.
	bi20L3 = newOrderSourceName("BI20", pubOrder)

	/// ABE (abe.io) exchange.
	abeL3 = newOrderSourceName("ABE", pubOrder)

	/// FAIR (FairX) exchange.
	fairL3 = newOrderSourceName("FAIR", pubOrder)

	/// CME Globex.
	glbxL3 = newOrderSourceName("GLBX", pubOrder|pubAnalyticOrder)

	/// CME Globex. Record for price level book.
	glbxL2 = newOrderSourceName("glbx", pubOrder)

	/// Eris Exchange group of companies.
	erisL3 = newOrderSourceName("ERIS", pubOrder)

	/// Eurex Exchange.
	xeurL3 = newOrderSourceName("XEUR", pubOrder)

	/// Eurex Exchange. Record for price level book.
	xeurL2 = newOrderSourceName("xeur", pubOrder)

	/// CBOE Futures Exchange.
	cfeL3 = newOrderSourceName("CFE", pubOrder)

	/// CBOE Options C2 Exchange.
	c20xL3 = newOrderSourceName("C2OX", pubOrder)

	/// Small Exchange.
	smfeL3 = newOrderSourceName("SMFE", pubOrder)

	/// Small Exchange. Record for price level book.
	smfeL2 = newOrderSourceName("smfe", pubOrder)

	/// Investors exchange. Record for price level book.
	iexL2 = newOrderSourceName("iex", pubOrder)

	/// Members Exchange.
	memxL3 = newOrderSourceName("MEMX", pubOrder)

	/// Members Exchange. Record for price level book.
	memxL2 = newOrderSourceName("memx", pubOrder)
)

func Default() *Source {
	return defaultSource
}

func CompsoiteBid() *Source {
	return compsoiteBid
}

func CompsoiteAsk() *Source {
	return compsoiteAsk
}

func RegionalBid() *Source {
	return regionalBid
}

func RegionalAsk() *Source {
	return regionalAsk
}

func AgregateBid() *Source {
	return agregateBid
}

func AgregateAsk() *Source {
	return agregateAsk
}

func NtvL2() *Source {
	return ntvL2
}

func NtvL3() *Source {
	return ntvL3
}

func NfxL3() *Source {
	return nfxL3
}

func EspdL3() *Source {
	return espdL3
}

func XnfiL3() *Source {
	return xnfiL3
}

func IceL3() *Source {
	return iceL3
}

func IseL3() *Source {
	return iseL3
}

func DeaL3() *Source {
	return deaL3
}

func DexL3() *Source {
	return dexL3
}

func ByxL3() *Source {
	return byxL3
}

func BzxL3() *Source {
	return bzxL3
}

func BateL3() *Source {
	return bateL3
}

func ChixL3() *Source {
	return chixL3
}

func CeuxL3() *Source {
	return ceuxL3
}

func BxtrL3() *Source {
	return bxtrL3
}

func IstL3() *Source {
	return istL3
}

func Bi20L3() *Source {
	return bi20L3
}

func AbeL3() *Source {
	return abeL3
}

func FairL3() *Source {
	return fairL3
}

func GlbxL3() *Source {
	return glbxL3
}

func GlbxL2() *Source {
	return glbxL2
}

func ErisL3() *Source {
	return erisL3
}

func XeurL2() *Source {
	return xeurL2
}

func XeurL3() *Source {
	return xeurL3
}

func CfeL3() *Source {
	return cfeL3
}

func C20xL3() *Source {
	return c20xL3
}

func SmfeL2() *Source {
	return smfeL2
}

func SmfeL3() *Source {
	return smfeL3
}

func IexL2() *Source {
	return iexL2
}

func MemxL2() *Source {
	return memxL2
}

func MemxL3() *Source {
	return memxL3
}
