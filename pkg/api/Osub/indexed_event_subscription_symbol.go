package Osub

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

type IndexedEventSubscriptionSymbol struct {
	symbol any
	source events.IndexedEventSourceInterface
}

func NewIndexedEventSubscriptionSymbol(symbol any, source events.IndexedEventSourceInterface) *IndexedEventSubscriptionSymbol {
	return &IndexedEventSubscriptionSymbol{symbol, source}
}

func (i IndexedEventSubscriptionSymbol) Symbol() any {
	return i.symbol
}

func (i IndexedEventSubscriptionSymbol) Source() events.IndexedEventSourceInterface {
	return i.source
}
