package order

var (
	sourceDefault = newOrderSourceWithIDNameFlagsNoError(0, "DEFAULT", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook)

	/// Bid side of a composite ``Quote``
	///
	/// It is a synthetic source.
	/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
	sourceCompsoiteBid = newOrderSourceWithIDNameFlagsNoError(1, "COMPOSITE_BID", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook)
	/// Ask side of a composite ``Quote``.
	/// It is a synthetic source.
	/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
	sourceCompsoiteAsk = newOrderSourceWithIDNameFlagsNoError(2, "COMPOSITE_ASK", 0)
	/// Bid side of a regional ``Quote``.
	/// It is a synthetic source.
	/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

	sourceRegionalBid = newOrderSourceWithIDNameFlagsNoError(3, "REGIONAL_BID", 0)
	/// Ask side of a regional ``Quote``.
	/// It is a synthetic source.
	/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

	sourceRegionalAsk = newOrderSourceWithIDNameFlagsNoError(4, "REGIONAL_ASK", 0)
	/// Bid side of an aggregate order book (futures depth and NASDAQ Level II).
	/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
	sourceAgregateBid = newOrderSourceWithIDNameFlagsNoError(5, "AGGREGATE_BID", 0)

	/// Ask side of an aggregate order book (futures depth and NASDAQ Level II).
	/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
	sourceAgregateAsk = newOrderSourceWithIDNameFlagsNoError(6, "AGGREGATE_ASK", 0)

	/// NASDAQ Total View.
	sourceNTV = newOrderSourceName("NTV", pubOrder|fullOrderBook)

	/// NASDAQ Total View. Record for price level book.
	sourcentv = newOrderSourceName("ntv", pubOrder)

	/// NASDAQ Futures Exchange.
	sourceNFX = newOrderSourceName("NFX", pubOrder)

	/// NASDAQ eSpeed.
	sourceESPD = newOrderSourceName("ESPD", pubOrder)

	/// NASDAQ Fixed Income.
	sourceXNFI = newOrderSourceName("XNFI", pubOrder)

	/// Intercontinental Exchange.
	sourceICE = newOrderSourceName("ICE", pubOrder)

	/// International Securities Exchange.
	sourceISE = newOrderSourceName("ISE", pubOrder|pubSpreadOrder)

	/// Direct-Edge EDGA Exchange.
	sourceDEA = newOrderSourceName("DEA", pubOrder)

	/// Direct-Edge EDGX Exchange.
	sourceDEX = newOrderSourceName("DEX", pubOrder)

	/// Bats BYX Exchange.
	sourceBYX = newOrderSourceName("BYX", pubOrder)

	/// Bats BZX Exchange.
	sourceBZX = newOrderSourceName("BZX", pubOrder)

	/// Bats Europe BXE Exchange.
	sourceBATE = newOrderSourceName("BATE", pubOrder)

	/// Bats Europe CXE Exchange.
	sourceCHIX = newOrderSourceName("CHIX", pubOrder)

	/// Bats Europe DXE Exchange.
	sourceCEUX = newOrderSourceName("CEUX", pubOrder)

	/// Bats Europe TRF.
	sourceBXTR = newOrderSourceName("BXTR", pubOrder)

	/// Borsa Istanbul Exchange.
	sourceIST = newOrderSourceName("IST", pubOrder)

	/// Borsa Istanbul Exchange. Record for particular top 20 order book.
	sourceBI20 = newOrderSourceName("BI20", pubOrder)

	/// ABE (abe.io) exchange.
	sourceABE = newOrderSourceName("ABE", pubOrder)

	/// FAIR (FairX) exchange.
	sourceFAIR = newOrderSourceName("FAIR", pubOrder)

	/// CME Globex.
	sourceGLBX = newOrderSourceName("GLBX", pubOrder|pubAnalyticOrder)

	/// CME Globex. Record for price level book.
	sourceglbx = newOrderSourceName("glbx", pubOrder)

	/// Eris Exchange group of companies.
	sourceERIS = newOrderSourceName("ERIS", pubOrder)

	/// Eurex Exchange.
	sourceXEUR = newOrderSourceName("XEUR", pubOrder)

	/// Eurex Exchange. Record for price level book.
	sourcexeur = newOrderSourceName("xeur", pubOrder)

	/// CBOE Futures Exchange.
	sourceCFE = newOrderSourceName("CFE", pubOrder)

	/// CBOE Options C2 Exchange.
	sourceC2OX = newOrderSourceName("C2OX", pubOrder)

	/// Small Exchange.
	sourceSMFE = newOrderSourceName("SMFE", pubOrder)

	/// Small Exchange. Record for price level book.
	sourcesmfe = newOrderSourceName("smfe", pubOrder)

	/// Investors exchange. Record for price level book.
	sourceiex = newOrderSourceName("iex", pubOrder)

	/// Members Exchange.
	sourceMEMX = newOrderSourceName("MEMX", pubOrder)

	/// Members Exchange. Record for price level book.
	sourcememx = newOrderSourceName("memx", pubOrder)
)

func SourceDefault() *OrderSource {
	return sourceDefault
}

func SourceCompsoiteBid() *OrderSource {
	return sourceCompsoiteBid
}

func SourceCompsoiteAsk() *OrderSource {
	return sourceCompsoiteAsk
}

func SourceRegionalBid() *OrderSource {
	return sourceRegionalBid
}

func SourceRegionalAsk() *OrderSource {
	return sourceRegionalAsk
}

func SourceAgregateBid() *OrderSource {
	return sourceAgregateBid
}

func SourceAgregateAsk() *OrderSource {
	return sourceAgregateAsk
}

func SourceNTV() *OrderSource {
	return sourceNTV
}

func Sourcentv() *OrderSource {
	return sourcentv
}

func SourceNFX() *OrderSource {
	return sourceNFX
}

func SourceESPD() *OrderSource {
	return sourceESPD
}

func SourceXNFI() *OrderSource {
	return sourceXNFI
}

func SourceICE() *OrderSource {
	return sourceICE
}

func SourceISE() *OrderSource {
	return sourceISE
}

func SourceDEA() *OrderSource {
	return sourceDEA
}

func SourceDEX() *OrderSource {
	return sourceDEX
}

func SourceBYX() *OrderSource {
	return sourceBYX
}

func SourceBZX() *OrderSource {
	return sourceBZX
}

func SourceBATE() *OrderSource {
	return sourceBATE
}

func SourceCHIX() *OrderSource {
	return sourceCHIX
}

func SourceCEUX() *OrderSource {
	return sourceCEUX
}

func SourceBXTR() *OrderSource {
	return sourceBXTR
}

func SourceIST() *OrderSource {
	return sourceIST
}

func SourceBI20() *OrderSource {
	return sourceBI20
}

func SourceABE() *OrderSource {
	return sourceABE
}

func SourceFAIR() *OrderSource {
	return sourceFAIR
}

func SourceGLBX() *OrderSource {
	return sourceGLBX
}

func Sourceglbx() *OrderSource {
	return sourceglbx
}

func SourceERIS() *OrderSource {
	return sourceERIS
}

func SourceXEUR() *OrderSource {
	return sourceXEUR
}

func Sourcexeur() *OrderSource {
	return sourcexeur
}

func SourceCFE() *OrderSource {
	return sourceCFE
}

func SourceC2OX() *OrderSource {
	return sourceC2OX
}

func SourceSMFE() *OrderSource {
	return sourceSMFE
}

func Sourcesmfe() *OrderSource {
	return sourcesmfe
}

func Sourceiex() *OrderSource {
	return sourceiex
}

func SourceMEMX() *OrderSource {
	return sourceMEMX
}

func Sourcememx() *OrderSource {
	return sourcememx
}
