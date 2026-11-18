package osub

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events"

// TimeSeriesSubscriptionSymbol subscribes to time series events, such as candles, of a symbol from a time.
type TimeSeriesSubscriptionSymbol struct {
	symbol   any
	fromTime int64
	source   events.IndexedEventSource
}

// NewTimeSeriesSubscriptionSymbol returns a symbol for the events of symbol from fromTime, in milliseconds since the
// Unix epoch.
func NewTimeSeriesSubscriptionSymbol(symbol any, fromTime int64) *TimeSeriesSubscriptionSymbol {
	return &TimeSeriesSubscriptionSymbol{symbol, fromTime, *events.DefaultIndexedEventSource()}
}

// FromTime returns the time, in milliseconds since the Unix epoch, from which events are requested.
func (symbol TimeSeriesSubscriptionSymbol) FromTime() int64 {
	return symbol.fromTime
}

// Symbol returns the underlying symbol.
func (symbol TimeSeriesSubscriptionSymbol) Symbol() any {
	return symbol.symbol
}
