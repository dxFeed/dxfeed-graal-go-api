package eventcodes

// EventCode identifies an event type; each value equals the matching dxfg_event_clazz_t constant of the native SDK.
type EventCode int32

const (
	Quote             EventCode = 0
	Profile           EventCode = 1
	Summary           EventCode = 2
	Greeks            EventCode = 3
	Candle            EventCode = 4
	DailyCandle       EventCode = 5
	Underlying        EventCode = 6
	TheoPrice         EventCode = 7
	Trade             EventCode = 8
	TradeETH          EventCode = 9
	Configuration     EventCode = 10
	Message           EventCode = 11
	TimeAndSale       EventCode = 12
	OrderBase         EventCode = 13
	Order             EventCode = 14
	AnalyticOrder     EventCode = 15
	OtcMarketsOrder   EventCode = 16
	SpreadOrder       EventCode = 17
	Series            EventCode = 18
	OptionSale        EventCode = 19
	TextMessage       EventCode = 20
	MarketMaker       EventCode = 21
	TextConfiguration EventCode = 22
	NuamOrder         EventCode = 23
	NuamTimeAndSale   EventCode = 24
	NuamTrade         EventCode = 25
	OrderImbalance    EventCode = 26
)

func (e EventCode) NativeCode() int32 {
	return int32(e)
}
