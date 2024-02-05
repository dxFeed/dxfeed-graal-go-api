package order

var (
	defaultSource = newOrderSourceWithIDNameFlagsNoError(0, "DEFAULT", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook)

	/// Bid side of a composite ``Quote``
	///
	/// It is a synthetic source.
	/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
	compsoiteBid = newOrderSourceWithIDNameFlagsNoError(1, "COMPOSITE_BID", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook)
	/// Ask side of a composite ``Quote``.
	/// It is a synthetic source.
	/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
	compsoiteAsk = newOrderSourceWithIDNameFlagsNoError(2, "COMPOSITE_ASK", 0)
	/// Bid side of a regional ``Quote``.
	/// It is a synthetic source.
	/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

	regionalBid = newOrderSourceWithIDNameFlagsNoError(3, "REGIONAL_BID", 0)
	/// Ask side of a regional ``Quote``.
	/// It is a synthetic source.
	/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

	regionalAsk = newOrderSourceWithIDNameFlagsNoError(4, "REGIONAL_ASK", 0)
	/// Bid side of an aggregate order book (futures depth and NASDAQ Level II).
	/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
	agregateBid = newOrderSourceWithIDNameFlagsNoError(5, "AGGREGATE_BID", 0)

	/// Ask side of an aggregate order book (futures depth and NASDAQ Level II).
	/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
	agregateAsk = newOrderSourceWithIDNameFlagsNoError(6, "AGGREGATE_ASK", 0)

	/// NASDAQ Total View.
	ntvL2 = newOrderSourceName("NTV", pubOrder|fullOrderBook)

	/// NASDAQ Total View. Record for price level book.
	ntvL3 = newOrderSourceName("ntv", pubOrder)

	/// NASDAQ Futures Exchange.
	nfx = newOrderSourceName("NFX", pubOrder)

	/// NASDAQ eSpeed.
	espd = newOrderSourceName("ESPD", pubOrder)

	/// NASDAQ Fixed Income.
	xnfi = newOrderSourceName("XNFI", pubOrder)

	/// Intercontinental Exchange.
	ice = newOrderSourceName("ICE", pubOrder)

	/// International Securities Exchange.
	ise = newOrderSourceName("ISE", pubOrder|pubSpreadOrder)

	/// Direct-Edge EDGA Exchange.
	dea = newOrderSourceName("DEA", pubOrder)

	/// Direct-Edge EDGX Exchange.
	dex = newOrderSourceName("DEX", pubOrder)

	/// Bats BYX Exchange.
	byx = newOrderSourceName("BYX", pubOrder)

	/// Bats BZX Exchange.
	bzx = newOrderSourceName("BZX", pubOrder)

	/// Bats Europe BXE Exchange.
	bate = newOrderSourceName("BATE", pubOrder)

	/// Bats Europe CXE Exchange.
	chix = newOrderSourceName("CHIX", pubOrder)

	/// Bats Europe DXE Exchange.
	ceux = newOrderSourceName("CEUX", pubOrder)

	/// Bats Europe TRF.
	bxtr = newOrderSourceName("BXTR", pubOrder)

	/// Borsa Istanbul Exchange.
	ist = newOrderSourceName("IST", pubOrder)

	/// Borsa Istanbul Exchange. Record for particular top 20 order book.
	bi20 = newOrderSourceName("BI20", pubOrder)

	/// ABE (abe.io) exchange.
	abe = newOrderSourceName("ABE", pubOrder)

	/// FAIR (FairX) exchange.
	fair = newOrderSourceName("FAIR", pubOrder)

	/// CME Globex.
	glbxL2 = newOrderSourceName("GLBX", pubOrder|pubAnalyticOrder)

	/// CME Globex. Record for price level book.
	glbxL3 = newOrderSourceName("glbx", pubOrder)

	/// Eris Exchange group of companies.
	eris = newOrderSourceName("ERIS", pubOrder)

	/// Eurex Exchange.
	xeurL2 = newOrderSourceName("XEUR", pubOrder)

	/// Eurex Exchange. Record for price level book.
	xeurL3 = newOrderSourceName("xeur", pubOrder)

	/// CBOE Futures Exchange.
	cfe = newOrderSourceName("CFE", pubOrder)

	/// CBOE Options C2 Exchange.
	c20x = newOrderSourceName("C2OX", pubOrder)

	/// Small Exchange.
	smfeL2 = newOrderSourceName("SMFE", pubOrder)

	/// Small Exchange. Record for price level book.
	smfeL3 = newOrderSourceName("smfe", pubOrder)

	/// Investors exchange. Record for price level book.
	iex = newOrderSourceName("iex", pubOrder)

	/// Members Exchange.
	memxL2 = newOrderSourceName("MEMX", pubOrder)

	/// Members Exchange. Record for price level book.
	memxL3 = newOrderSourceName("memx", pubOrder)
)

func Default() *OrderSource {
	return defaultSource
}

func CompsoiteBid() *OrderSource {
	return compsoiteBid
}

func CompsoiteAsk() *OrderSource {
	return compsoiteAsk
}

func RegionalBid() *OrderSource {
	return regionalBid
}

func RegionalAsk() *OrderSource {
	return regionalAsk
}

func AgregateBid() *OrderSource {
	return agregateBid
}

func AgregateAsk() *OrderSource {
	return agregateAsk
}

func NtvL2() *OrderSource {
	return ntvL2
}

func NtvL3() *OrderSource {
	return ntvL3
}

func Nfx() *OrderSource {
	return nfx
}

func Espd() *OrderSource {
	return espd
}

func Xnfi() *OrderSource {
	return xnfi
}

func Ice() *OrderSource {
	return ice
}

func Ise() *OrderSource {
	return ise
}

func Dea() *OrderSource {
	return dea
}

func Dex() *OrderSource {
	return dex
}

func Byx() *OrderSource {
	return byx
}

func Bzx() *OrderSource {
	return bzx
}

func Bate() *OrderSource {
	return bate
}

func CHIX() *OrderSource {
	return chix
}

func Ceux() *OrderSource {
	return ceux
}

func Bxtr() *OrderSource {
	return bxtr
}

func Ist() *OrderSource {
	return ist
}

func Bi20() *OrderSource {
	return bi20
}

func Abe() *OrderSource {
	return abe
}

func Fair() *OrderSource {
	return fair
}

func GlbxL3() *OrderSource {
	return glbxL3
}

func GlbxL2() *OrderSource {
	return glbxL2
}

func Eris() *OrderSource {
	return eris
}

func XeurL2() *OrderSource {
	return xeurL2
}

func XeurL3() *OrderSource {
	return xeurL3
}

func Cfe() *OrderSource {
	return cfe
}

func C20x() *OrderSource {
	return c20x
}

func SmfeL2() *OrderSource {
	return smfeL2
}

func SmfeL3() *OrderSource {
	return smfeL3
}

func Iex() *OrderSource {
	return iex
}

func MemxL2() *OrderSource {
	return memxL2
}

func MemxL3() *OrderSource {
	return memxL3
}
