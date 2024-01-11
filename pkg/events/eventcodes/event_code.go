package eventcodes

type EventCode int32

const (
	Quote EventCode = iota
	Profile
	Summary
	Greeks
	Candle
	DailyCandle
	Underlying
	TheoPrice
	Trade
	TradeETH
	Configuration
	Message
	TimeAndSale
	OrderBase
	Order
	AnalyticOrder
	SpreadOrder
	Series
	OptionSale
)

func (e EventCode) NativeCode() int32 {
	return int32(e)
}
