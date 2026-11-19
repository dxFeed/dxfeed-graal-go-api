// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

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

// Default returns the default order source.
func Default() *Source {
	return defaultSource
}

// CompsoiteBid returns the order source for the bid side of a composite quote.
func CompsoiteBid() *Source {
	return compsoiteBid
}

// CompsoiteAsk returns the order source for the ask side of a composite quote.
func CompsoiteAsk() *Source {
	return compsoiteAsk
}

// RegionalBid returns the order source for the bid side of a regional quote.
func RegionalBid() *Source {
	return regionalBid
}

// RegionalAsk returns the order source for the ask side of a regional quote.
func RegionalAsk() *Source {
	return regionalAsk
}

// AgregateBid returns the order source for the aggregate bid side.
func AgregateBid() *Source {
	return agregateBid
}

// AgregateAsk returns the order source for the aggregate ask side.
func AgregateAsk() *Source {
	return agregateAsk
}

// NtvL2 returns the NTV order source (price level book).
func NtvL2() *Source {
	return ntvL2
}

// NtvL3 returns the NTV order source (individual orders).
func NtvL3() *Source {
	return ntvL3
}

// NfxL3 returns the NFX order source (individual orders).
func NfxL3() *Source {
	return nfxL3
}

// EspdL3 returns the ESPD order source (individual orders).
func EspdL3() *Source {
	return espdL3
}

// XnfiL3 returns the XNFI order source (individual orders).
func XnfiL3() *Source {
	return xnfiL3
}

// IceL3 returns the ICE order source (individual orders).
func IceL3() *Source {
	return iceL3
}

// IseL3 returns the ISE order source (individual orders).
func IseL3() *Source {
	return iseL3
}

// DeaL3 returns the DEA order source (individual orders).
func DeaL3() *Source {
	return deaL3
}

// DexL3 returns the DEX order source (individual orders).
func DexL3() *Source {
	return dexL3
}

// ByxL3 returns the BYX order source (individual orders).
func ByxL3() *Source {
	return byxL3
}

// BzxL3 returns the BZX order source (individual orders).
func BzxL3() *Source {
	return bzxL3
}

// BateL3 returns the BATE order source (individual orders).
func BateL3() *Source {
	return bateL3
}

// ChixL3 returns the CHIX order source (individual orders).
func ChixL3() *Source {
	return chixL3
}

// CeuxL3 returns the CEUX order source (individual orders).
func CeuxL3() *Source {
	return ceuxL3
}

// BxtrL3 returns the BXTR order source (individual orders).
func BxtrL3() *Source {
	return bxtrL3
}

// IstL3 returns the IST order source (individual orders).
func IstL3() *Source {
	return istL3
}

// Bi20L3 returns the BI20 order source (individual orders).
func Bi20L3() *Source {
	return bi20L3
}

// AbeL3 returns the ABE order source (individual orders).
func AbeL3() *Source {
	return abeL3
}

// FairL3 returns the FAIR order source (individual orders).
func FairL3() *Source {
	return fairL3
}

// GlbxL3 returns the GLBX order source (individual orders).
func GlbxL3() *Source {
	return glbxL3
}

// GlbxL2 returns the GLBX order source (price level book).
func GlbxL2() *Source {
	return glbxL2
}

// ErisL3 returns the ERIS order source (individual orders).
func ErisL3() *Source {
	return erisL3
}

// XeurL2 returns the XEUR order source (price level book).
func XeurL2() *Source {
	return xeurL2
}

// XeurL3 returns the XEUR order source (individual orders).
func XeurL3() *Source {
	return xeurL3
}

// CfeL3 returns the CFE order source (individual orders).
func CfeL3() *Source {
	return cfeL3
}

// C20xL3 returns the C2OX order source (individual orders).
func C20xL3() *Source {
	return c20xL3
}

// SmfeL2 returns the SMFE order source (price level book).
func SmfeL2() *Source {
	return smfeL2
}

// SmfeL3 returns the SMFE order source (individual orders).
func SmfeL3() *Source {
	return smfeL3
}

// IexL2 returns the IEX order source (price level book).
func IexL2() *Source {
	return iexL2
}

// MemxL2 returns the MEMX order source (price level book).
func MemxL2() *Source {
	return memxL2
}

// MemxL3 returns the MEMX order source (individual orders).
func MemxL3() *Source {
	return memxL3
}
