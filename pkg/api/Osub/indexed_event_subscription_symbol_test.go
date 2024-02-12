package Osub

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"reflect"
	"testing"
)

func TestIndexedEventSource(t *testing.T) {
	test_name := "test_key"
	eventSource := events.NewIndexedEventSource(1, test_name)
	symbol := NewIndexedEventSubscriptionSymbol("TEST", eventSource)

	if *symbol.source.Name() != test_name || reflect.TypeOf(symbol.source) != reflect.TypeOf(eventSource) ||
		symbol.source.Type() != events.IndexedEventSourceType {
		t.Fatalf(`Wrong values in subscription symbol`)
	}
}

func TestOrderSourceType(t *testing.T) {
	eventSource := order.NtvL2()
	symbol := NewIndexedEventSubscriptionSymbol("TEST", eventSource)

	if *symbol.source.Name() != *eventSource.Name() || reflect.TypeOf(symbol.source) != reflect.TypeOf(eventSource) ||
		symbol.source.Type() != events.OrderSourceType {
		t.Fatalf(`Wrong values in subscription symbol`)
	}
}
