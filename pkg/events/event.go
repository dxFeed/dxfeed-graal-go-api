package events

type EventType int32

const (
	Quote EventType = iota
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
