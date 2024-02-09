package Osub

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

type IndexedEventSubscriptionSymbol struct {
	symbol any
	source events.IndexedEventSourceInterface
}

func NewIndexedEventSubscriptionSymbolWithOrderSource(symbol any, source *order.OrderSource) *IndexedEventSubscriptionSymbol {
	return &IndexedEventSubscriptionSymbol{symbol, source}
}

func NewIndexedEventSubscriptionSymbolWithIndexedSymbol(symbol any, source *events.IndexedEventSource) *IndexedEventSubscriptionSymbol {
	return &IndexedEventSubscriptionSymbol{symbol, source}
}

func (i IndexedEventSubscriptionSymbol) Symbol() any {
	return i.symbol
}

func (i IndexedEventSubscriptionSymbol) Source() events.IndexedEventSourceInterface {
	return i.source
}
