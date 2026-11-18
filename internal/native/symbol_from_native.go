package native

/*
#include "dxfg_api.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

// goSymbol copies a symbol that the SDK returned into Go memory. A string or candle symbol becomes osub.String; the
// symbol nested in a time series or indexed event subscription symbol becomes a plain string, as the constructors of
// those symbols usually get it.
func goSymbol(s *C.dxfg_v4_symbol_t) (osub.Symbol, error) {
	if s == nil {
		return nil, fmt.Errorf("native: nil symbol")
	}
	switch s._type {
	case C.STRING:
		return osub.String(C.GoString((*C.dxfg_v4_string_symbol_t)(unsafe.Pointer(s)).symbol)), nil
	case C.CANDLE:
		return osub.String(C.GoString((*C.dxfg_v4_candle_symbol_t)(unsafe.Pointer(s)).symbol)), nil
	case C.WILDCARD:
		return osub.NewWildcardSymbol(), nil
	case C.TIME_SERIES_SUBSCRIPTION:
		ts := (*C.dxfg_v4_time_series_subscription_symbol_t)(unsafe.Pointer(s))
		inner, err := goNestedSymbol(ts.symbol)
		if err != nil {
			return nil, err
		}
		return osub.NewTimeSeriesSubscriptionSymbol(inner, int64(ts.from_time)), nil
	case C.INDEXED_EVENT_SUBSCRIPTION:
		indexed := (*C.dxfg_v4_indexed_event_subscription_symbol_t)(unsafe.Pointer(s))
		inner, err := goNestedSymbol(indexed.symbol)
		if err != nil {
			return nil, err
		}
		source, err := goIndexedEventSource(indexed.source)
		if err != nil {
			return nil, err
		}
		return osub.NewIndexedEventSubscriptionSymbol(inner, source), nil
	default:
		return nil, fmt.Errorf("native: unsupported symbol type %d", int(s._type))
	}
}

// goNestedSymbol returns a nested string or candle symbol as a plain string, and any other symbol as goSymbol does.
func goNestedSymbol(s *C.dxfg_v4_symbol_t) (any, error) {
	symbol, err := goSymbol(s)
	if err != nil {
		return nil, err
	}
	if text, ok := symbol.(osub.String); ok {
		return string(text), nil
	}
	return symbol, nil
}

func goIndexedEventSource(s *C.dxfg_v4_indexed_event_source_t) (events.IndexedEventSourceInterface, error) {
	if s == nil {
		return nil, fmt.Errorf("native: indexed event subscription symbol without a source")
	}
	switch s._type {
	case C.ORDER_SOURCE:
		return order.ValueOfIdentifier(int64(s.id))
	case C.INDEXED_EVENT_SOURCE:
		return events.NewIndexedEventSource(int64(s.id), C.GoString(s.name)), nil
	default:
		return nil, fmt.Errorf("native: unsupported indexed event source type %d", int(s._type))
	}
}

// takeSymbolList copies a list of symbols that the SDK returned and releases the list.
func takeSymbolList(thread *isolateThread, list *C.dxfg_v4_symbol_list) ([]osub.Symbol, error) {
	if list == nil {
		return nil, nil
	}
	var symbols []osub.Symbol
	var err error
	if list.size > 0 && list.elements != nil {
		symbols = make([]osub.Symbol, 0, int(list.size))
		for _, s := range unsafe.Slice(list.elements, int(list.size)) {
			symbol, convErr := goSymbol(s)
			if convErr != nil {
				err = convErr
				symbols = nil
				break
			}
			symbols = append(symbols, symbol)
		}
	}
	if releaseErr := checkStatus(thread, C.dxfg_v4_symbol_list_release(thread.ptr, list)); err == nil {
		err = releaseErr
	}
	if err != nil {
		return nil, err
	}
	return symbols, nil
}
