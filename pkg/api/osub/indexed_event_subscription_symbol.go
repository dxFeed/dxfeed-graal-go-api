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
