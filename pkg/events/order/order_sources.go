package order

import "sync"

type OrderSources interface {
	OsDefault() *OrderSource
	OsCompsoiteBid() *OrderSource
	OsCompsoiteAsk() *OrderSource
	OsRegionalBid() *OrderSource
	OsRegionalAsk() *OrderSource
	OsAgregateBid() *OrderSource
	OsAgregateAsk() *OrderSource
	OsNTV() *OrderSource
	Osntv() *OrderSource
	OsNFX() *OrderSource
	OsESPD() *OrderSource
	OsXNFI() *OrderSource
	OsICE() *OrderSource
	OsISE() *OrderSource
	OsDEA() *OrderSource
	OsDEX() *OrderSource
	OsBYX() *OrderSource
	OsBZX() *OrderSource
	OsBATE() *OrderSource
	OsCHIX() *OrderSource
	OsCEUX() *OrderSource
	OsBXTR() *OrderSource
	OsIST() *OrderSource
	OsBI20() *OrderSource
	OsABE() *OrderSource
	OsFAIR() *OrderSource
	OsGLBX() *OrderSource
	Osglbx() *OrderSource
	OsERIS() *OrderSource
	OsXEUR() *OrderSource
	Osxeur() *OrderSource
	OsCFE() *OrderSource
	OsC2OX() *OrderSource
	OsSMFE() *OrderSource
	Ossmfe() *OrderSource
	Osiex() *OrderSource
	OsMEMX() *OrderSource
	Osmemx() *OrderSource
}

type orderSourceConsts struct {
	osDefault      *OrderSource
	osCompsoiteBid *OrderSource
	osCompsoiteAsk *OrderSource
	osRegionalBid  *OrderSource
	osRegionalAsk  *OrderSource
	osAgregateBid  *OrderSource
	osAgregateAsk  *OrderSource
	osNTV          *OrderSource
	osntv          *OrderSource
	osNFX          *OrderSource
	osESPD         *OrderSource
	osXNFI         *OrderSource
	osICE          *OrderSource
	osISE          *OrderSource
	osDEA          *OrderSource
	osDEX          *OrderSource
	osBYX          *OrderSource
	osBZX          *OrderSource
	osBATE         *OrderSource
	osCHIX         *OrderSource
	osCEUX         *OrderSource
	osBXTR         *OrderSource
	osIST          *OrderSource
	osBI20         *OrderSource
	osABE          *OrderSource
	osFAIR         *OrderSource
	osGLBX         *OrderSource
	osglbx         *OrderSource
	osERIS         *OrderSource
	osXEUR         *OrderSource
	osxeur         *OrderSource
	osCFE          *OrderSource
	osC2OX         *OrderSource
	osSMFE         *OrderSource
	ossmfe         *OrderSource
	osiex          *OrderSource
	osMEMX         *OrderSource
	osmemx         *OrderSource
}

func (o orderSourceConsts) OsDefault() *OrderSource {
	return o.osDefault
}

func (o orderSourceConsts) OsCompsoiteBid() *OrderSource {
	return o.osCompsoiteBid
}

func (o orderSourceConsts) OsCompsoiteAsk() *OrderSource {
	return o.osCompsoiteAsk
}

func (o orderSourceConsts) OsRegionalBid() *OrderSource {
	return o.osRegionalBid
}

func (o orderSourceConsts) OsRegionalAsk() *OrderSource {
	return o.osRegionalAsk
}

func (o orderSourceConsts) OsAgregateBid() *OrderSource {
	return o.osAgregateBid
}

func (o orderSourceConsts) OsAgregateAsk() *OrderSource {
	return o.osAgregateAsk
}

func (o orderSourceConsts) OsNTV() *OrderSource {
	return o.osNTV
}

func (o orderSourceConsts) Osntv() *OrderSource {
	return o.osntv
}

func (o orderSourceConsts) OsNFX() *OrderSource {
	return o.osNFX
}

func (o orderSourceConsts) OsESPD() *OrderSource {
	return o.osESPD
}

func (o orderSourceConsts) OsXNFI() *OrderSource {
	return o.osXNFI
}

func (o orderSourceConsts) OsICE() *OrderSource {
	return o.osICE
}

func (o orderSourceConsts) OsISE() *OrderSource {
	return o.osISE
}

func (o orderSourceConsts) OsDEA() *OrderSource {
	return o.osDEA
}

func (o orderSourceConsts) OsDEX() *OrderSource {
	return o.osDEX
}

func (o orderSourceConsts) OsBYX() *OrderSource {
	return o.osBYX
}

func (o orderSourceConsts) OsBZX() *OrderSource {
	return o.osBZX
}

func (o orderSourceConsts) OsBATE() *OrderSource {
	return o.osBATE
}

func (o orderSourceConsts) OsCHIX() *OrderSource {
	return o.osCHIX
}

func (o orderSourceConsts) OsCEUX() *OrderSource {
	return o.osCEUX
}

func (o orderSourceConsts) OsBXTR() *OrderSource {
	return o.osBXTR
}

func (o orderSourceConsts) OsIST() *OrderSource {
	return o.osIST
}

func (o orderSourceConsts) OsBI20() *OrderSource {
	return o.osBI20
}

func (o orderSourceConsts) OsABE() *OrderSource {
	return o.osABE
}

func (o orderSourceConsts) OsFAIR() *OrderSource {
	return o.osFAIR
}

func (o orderSourceConsts) OsGLBX() *OrderSource {
	return o.osGLBX
}

func (o orderSourceConsts) Osglbx() *OrderSource {
	return o.osglbx
}

func (o orderSourceConsts) OsERIS() *OrderSource {
	return o.osERIS
}

func (o orderSourceConsts) OsXEUR() *OrderSource {
	return o.osXEUR
}

func (o orderSourceConsts) Osxeur() *OrderSource {
	return o.osxeur
}

func (o orderSourceConsts) OsCFE() *OrderSource {
	return o.osCFE
}

func (o orderSourceConsts) OsC2OX() *OrderSource {
	return o.osC2OX
}

func (o orderSourceConsts) OsSMFE() *OrderSource {
	return o.osSMFE
}

func (o orderSourceConsts) Ossmfe() *OrderSource {
	return o.ossmfe
}

func (o orderSourceConsts) Osiex() *OrderSource {
	return o.osiex
}

func (o orderSourceConsts) OsMEMX() *OrderSource {
	return o.osMEMX
}

func (o orderSourceConsts) Osmemx() *OrderSource {
	return o.osmemx
}

var instance OrderSources = nil
var once sync.Once

func GetConsts() OrderSources {
	once.Do(func() {
		instance = &orderSourceConsts{ /// The default source with zero identifier for all events that do not support multiple sources.
			osDefault: newOrderSourceWithIDNameFlagsNoError(0, "DEFAULT", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook),

			/// Bid side of a composite ``Quote``
			///
			/// It is a synthetic source.
			/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
			osCompsoiteBid: newOrderSourceWithIDNameFlagsNoError(1, "COMPOSITE_BID", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook),
			/// Ask side of a composite ``Quote``.
			/// It is a synthetic source.
			/// The subscription on composite ``Quote`` event is observed when this source is subscribed to.
			osCompsoiteAsk: newOrderSourceWithIDNameFlagsNoError(2, "COMPOSITE_ASK", 0),
			/// Bid side of a regional ``Quote``.
			/// It is a synthetic source.
			/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

			osRegionalBid: newOrderSourceWithIDNameFlagsNoError(3, "REGIONAL_BID", 0),
			/// Ask side of a regional ``Quote``.
			/// It is a synthetic source.
			/// The subscription on regional ``Quote`` event is observed when this source is subscribed to.

			osRegionalAsk: newOrderSourceWithIDNameFlagsNoError(4, "REGIONAL_ASK", 0),
			/// Bid side of an aggregate order book (futures depth and NASDAQ Level II).
			/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
			osAgregateBid: newOrderSourceWithIDNameFlagsNoError(5, "AGGREGATE_BID", 0),

			/// Ask side of an aggregate order book (futures depth and NASDAQ Level II).
			/// This source cannot be directly published via dxFeed API, but otherwise it is fully operational.
			osAgregateAsk: newOrderSourceWithIDNameFlagsNoError(6, "AGGREGATE_ASK", 0),

			/// NASDAQ Total View.
			osNTV: newOrderSourceName("NTV", pubOrder|fullOrderBook),

			/// NASDAQ Total View. Record for price level book.
			osntv: newOrderSourceName("ntv", pubOrder),

			/// NASDAQ Futures Exchange.
			osNFX: newOrderSourceName("NFX", pubOrder),

			/// NASDAQ eSpeed.
			osESPD: newOrderSourceName("ESPD", pubOrder),

			/// NASDAQ Fixed Income.
			osXNFI: newOrderSourceName("XNFI", pubOrder),

			/// Intercontinental Exchange.
			osICE: newOrderSourceName("ICE", pubOrder),

			/// International Securities Exchange.
			osISE: newOrderSourceName("ISE", pubOrder|pubSpreadOrder),

			/// Direct-Edge EDGA Exchange.
			osDEA: newOrderSourceName("DEA", pubOrder),

			/// Direct-Edge EDGX Exchange.
			osDEX: newOrderSourceName("DEX", pubOrder),

			/// Bats BYX Exchange.
			osBYX: newOrderSourceName("BYX", pubOrder),

			/// Bats BZX Exchange.
			osBZX: newOrderSourceName("BZX", pubOrder),

			/// Bats Europe BXE Exchange.
			osBATE: newOrderSourceName("BATE", pubOrder),

			/// Bats Europe CXE Exchange.
			osCHIX: newOrderSourceName("CHIX", pubOrder),

			/// Bats Europe DXE Exchange.
			osCEUX: newOrderSourceName("CEUX", pubOrder),

			/// Bats Europe TRF.
			osBXTR: newOrderSourceName("BXTR", pubOrder),

			/// Borsa Istanbul Exchange.
			osIST: newOrderSourceName("IST", pubOrder),

			/// Borsa Istanbul Exchange. Record for particular top 20 order book.
			osBI20: newOrderSourceName("BI20", pubOrder),

			/// ABE (abe.io) exchange.
			osABE: newOrderSourceName("ABE", pubOrder),

			/// FAIR (FairX) exchange.
			osFAIR: newOrderSourceName("FAIR", pubOrder),

			/// CME Globex.
			osGLBX: newOrderSourceName("GLBX", pubOrder|pubAnalyticOrder),

			/// CME Globex. Record for price level book.
			osglbx: newOrderSourceName("glbx", pubOrder),

			/// Eris Exchange group of companies.
			osERIS: newOrderSourceName("ERIS", pubOrder),

			/// Eurex Exchange.
			osXEUR: newOrderSourceName("XEUR", pubOrder),

			/// Eurex Exchange. Record for price level book.
			osxeur: newOrderSourceName("xeur", pubOrder),

			/// CBOE Futures Exchange.
			osCFE: newOrderSourceName("CFE", pubOrder),

			/// CBOE Options C2 Exchange.
			osC2OX: newOrderSourceName("C2OX", pubOrder),

			/// Small Exchange.
			osSMFE: newOrderSourceName("SMFE", pubOrder),

			/// Small Exchange. Record for price level book.
			ossmfe: newOrderSourceName("smfe", pubOrder),

			/// Investors exchange. Record for price level book.
			osiex: newOrderSourceName("iex", pubOrder),

			/// Members Exchange.
			osMEMX: newOrderSourceName("MEMX", pubOrder),

			/// Members Exchange. Record for price level book.
			osmemx: newOrderSourceName("memx", pubOrder)}
	})
	return instance
}
