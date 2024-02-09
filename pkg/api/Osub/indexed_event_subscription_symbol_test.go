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
	symbol := NewIndexedEventSubscriptionSymbolWithIndexedSymbol("TEST", eventSource)

	if *symbol.source.Name() != test_name || reflect.TypeOf(symbol.source) != reflect.TypeOf(eventSource) {
		t.Fatalf(`Wrong values in subscription symbol`)
	}
}

func TestOrderSourceType(t *testing.T) {
	eventSource := order.NtvL2()
	symbol := NewIndexedEventSubscriptionSymbolWithOrderSource("TEST", eventSource)

	if *symbol.source.Name() != *eventSource.Name() || reflect.TypeOf(symbol.source) != reflect.TypeOf(eventSource) {
		t.Fatalf(`Wrong values in subscription symbol`)
	}
}
