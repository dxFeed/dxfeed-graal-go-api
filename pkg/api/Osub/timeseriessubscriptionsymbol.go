package Osub

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events"

type TimeSeriesSubscriptionSymbol struct {
	symbol   any
	fromTime int64
	source   events.IndexedEventSource
}

func NewTimeSeriesSubscriptionSymbol(symbol any, fromTime int64) *TimeSeriesSubscriptionSymbol {
	return &TimeSeriesSubscriptionSymbol{symbol, fromTime, *events.DefaultIndexedEventSource()}
}

func (symbol TimeSeriesSubscriptionSymbol) FromTime() int64 {
	return symbol.fromTime
}

func (symbol TimeSeriesSubscriptionSymbol) Symbol() any {
	return symbol.symbol
}
