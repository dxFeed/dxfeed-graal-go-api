package Osub

type TimeSeriesSubscriptionSymbol struct {
	symbol   any
	fromTime int64
	source   IndexedEventSource
}

func NewTimeSeriesSubscriptionSymbol(symbol any, fromTime int64) TimeSeriesSubscriptionSymbol {
	return TimeSeriesSubscriptionSymbol{symbol, fromTime, DefaultIndexedEventSource}
}

func (symbol TimeSeriesSubscriptionSymbol) GetFromTime() int64 {
	return symbol.fromTime
}

func (symbol TimeSeriesSubscriptionSymbol) GetSymbol() any {
	return symbol.symbol
}
