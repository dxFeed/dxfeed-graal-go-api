// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package osub

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// IndexedEventSubscriptionSymbol subscribes to indexed events, such as orders, of a symbol from one source.
type IndexedEventSubscriptionSymbol struct {
	symbol any
	source events.IndexedEventSourceInterface
}

// NewIndexedEventSubscriptionSymbol returns a symbol for the events of symbol from source, for example an order source.
// The symbol is a string or an osub.String, or a *TimeSeriesSubscriptionSymbol for time series indexed events.
func NewIndexedEventSubscriptionSymbol(symbol any, source events.IndexedEventSourceInterface) *IndexedEventSubscriptionSymbol {
	return &IndexedEventSubscriptionSymbol{symbol, source}
}

// Symbol returns the underlying symbol.
func (i IndexedEventSubscriptionSymbol) Symbol() any {
	return i.symbol
}

// Source returns the source of the events.
func (i IndexedEventSubscriptionSymbol) Source() events.IndexedEventSourceInterface {
	return i.source
}
