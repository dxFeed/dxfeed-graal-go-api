package order

import (
	"fmt"
	"regexp"
)

type OrderSource struct {
	IndexedEventSource
	pubFlags  int64
	isBuiltin bool
}

const (
	pubOrder           = 0x0001
	pubAnalyticOrder   = 0x0002
	pubOtcMarketsOrder = 0x0004
	pubSpreadOrder     = 0x0008
	fullOrderBook      = 0x0010
	flagsSize          = 5
)

var (
	/// The default source with zero identifier for all events that do not support multiple sources.
	defaultOrderSource = newOrderSourceWithIDNameFlags(0, "DEFAULT", pubOrder|pubAnalyticOrder|pubSpreadOrder|fullOrderBook)

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
	NTV, _ = newOrderSourceName("NTV", pubOrder|fullOrderBook)

	/// NASDAQ Total View. Record for price level book.
	ntv, _ = newOrderSourceName("ntv", pubOrder)

	/// NASDAQ Futures Exchange.
	NFX, _ = newOrderSourceName("NFX", pubOrder)

	/// NASDAQ eSpeed.
	ESPD, _ = newOrderSourceName("ESPD", pubOrder)

	/// NASDAQ Fixed Income.
	XNFI, _ = newOrderSourceName("XNFI", pubOrder)

	/// Intercontinental Exchange.
	ICE, _ = newOrderSourceName("ICE", pubOrder)

	/// International Securities Exchange.
	ISE, _ = newOrderSourceName("ISE", pubOrder|pubSpreadOrder)

	/// Direct-Edge EDGA Exchange.
	DEA, _ = newOrderSourceName("DEA", pubOrder)

	/// Direct-Edge EDGX Exchange.
	DEX, _ = newOrderSourceName("DEX", pubOrder)

	/// Bats BYX Exchange.
	BYX, _ = newOrderSourceName("BYX", pubOrder)

	/// Bats BZX Exchange.
	BZX, _ = newOrderSourceName("BZX", pubOrder)

	/// Bats Europe BXE Exchange.
	BATE, _ = newOrderSourceName("BATE", pubOrder)

	/// Bats Europe CXE Exchange.
	CHIX, _ = newOrderSourceName("CHIX", pubOrder)

	/// Bats Europe DXE Exchange.
	CEUX, _ = newOrderSourceName("CEUX", pubOrder)

	/// Bats Europe TRF.
	BXTR, _ = newOrderSourceName("BXTR", pubOrder)

	/// Borsa Istanbul Exchange.
	IST, _ = newOrderSourceName("IST", pubOrder)

	/// Borsa Istanbul Exchange. Record for particular top 20 order book.
	BI20, _ = newOrderSourceName("BI20", pubOrder)

	/// ABE (abe.io) exchange.
	ABE, _ = newOrderSourceName("ABE", pubOrder)

	/// FAIR (FairX) exchange.
	FAIR, _ = newOrderSourceName("FAIR", pubOrder)

	/// CME Globex.
	GLBX, _ = newOrderSourceName("GLBX", pubOrder|pubAnalyticOrder)

	/// CME Globex. Record for price level book.
	glbx, _ = newOrderSourceName("glbx", pubOrder)

	/// Eris Exchange group of companies.
	ERIS, _ = newOrderSourceName("ERIS", pubOrder)

	/// Eurex Exchange.
	XEUR, _ = newOrderSourceName("XEUR", pubOrder)

	/// Eurex Exchange. Record for price level book.
	xeur, _ = newOrderSourceName("xeur", pubOrder)

	/// CBOE Futures Exchange.
	CFE, _ = newOrderSourceName("CFE", pubOrder)

	/// CBOE Options C2 Exchange.
	C2OX, _ = newOrderSourceName("C2OX", pubOrder)

	/// Small Exchange.
	SMFE, _ = newOrderSourceName("SMFE", pubOrder)

	/// Small Exchange. Record for price level book.
	smfe, _ = newOrderSourceName("smfe", pubOrder)

	/// Investors exchange. Record for price level book.
	iex, _ = newOrderSourceName("iex", pubOrder)

	/// Members Exchange.
	MEMX, _ = newOrderSourceName("MEMX", pubOrder)

	/// Members Exchange. Record for price level book.
	memx, _              = newOrderSourceName("memx", pubOrder)
	allOrderSourceValues = []*OrderSource{
		defaultOrderSource,
		compsoiteBid,
		compsoiteAsk,
		regionalBid,
		regionalAsk,
		agregateBid,
		agregateAsk,
		NTV,
		ntv,
		NFX,
		ESPD,
		XNFI,
		ICE,
		ISE,
		DEA,
		DEX,
		BYX,
		BZX,
		BATE,
		CHIX,
		CEUX,
		BXTR,
		IST,
		BI20,
		ABE,
		FAIR,
		GLBX,
		glbx,
		ERIS,
		XEUR,
		xeur,
		CFE,
		C2OX,
		SMFE,
		smfe,
		iex,
		MEMX,
		memx}
)

func newOrderSourceWithIDName(identifier int64, name *string) *OrderSource {
	return &OrderSource{IndexedEventSource: IndexedEventSource{identifier: identifier, name: name}, pubFlags: 0, isBuiltin: false}
}

func newOrderSourceWithIDNameFlags(identifier int64, name string, pubflags int64) *OrderSource {

	return &OrderSource{IndexedEventSource: IndexedEventSource{identifier: identifier, name: &name}, pubFlags: pubflags, isBuiltin: true}
}

func newOrderSourceName(name string, pubflags int64) (*OrderSource, error) {
	id, err := OrderSourceComposeId(name)
	if err != nil {
		return nil, err
	}
	return newOrderSourceWithIDNameFlags(id, name, pubflags), nil
}

func OrderSourceComposeId(name string) (int64, error) {
	var sourceId int64
	count := len(name)
	if count == 0 || count > 4 {
		return 0, nil
	}
	for _, ch := range name {
		str := fmt.Sprintf("%c", ch)
		notAlpha := OrderSourceCheck(str)
		if !notAlpha {
			return 0, fmt.Errorf("source name must contain only alphanumeric characters. Current %s", str)
		}
		sourceId = (sourceId << 8) | int64(ch)
	}

	return sourceId, nil
}

func OrderSourceCheck(char string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(char)
}

func IsSpecialSourceId(value int64) bool {
	return value >= 1 && value <= 6
}

func OrderSourceValueOfIdentifier(id int64) (*OrderSource, error) {
	panic("impl it")
	return nil, nil
}
