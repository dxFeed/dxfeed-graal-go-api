package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
)

// newCSymbol builds a subscription symbol in C memory, so the SDK never sees Go memory; free it with freeCSymbol.
// It accepts a string, osub.String, *osub.CandleSymbol, *candle.CandleSymbol, *osub.WildcardSymbol,
// *osub.TimeSeriesSubscriptionSymbol and *osub.IndexedEventSubscriptionSymbol; anything else is an error.
func newCSymbol(symbol any) (*C.dxfg_v4_symbol_t, error) {
	switch value := symbol.(type) {
	case osub.String:
		return newCSymbol(string(value))
	case string:
		s := (*C.dxfg_v4_string_symbol_t)(C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_string_symbol_t{}))))
		s.supper._type = C.STRING
		s.symbol = C.CString(value)
		return &s.supper, nil
	case *osub.CandleSymbol:
		if value == nil {
			return nil, fmt.Errorf("nil %T symbol", symbol)
		}
		return newCSymbol(value.CandleSymbol())
	case *candle.CandleSymbol:
		if value == nil {
			return nil, fmt.Errorf("nil %T symbol", symbol)
		}
		s := (*C.dxfg_candle_symbol_t)(C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_candle_symbol_t{}))))
		s.supper._type = C.CANDLE
		s.symbol = C.CString(value.String())
		return &s.supper, nil
	case *osub.WildcardSymbol:
		if value == nil {
			return nil, fmt.Errorf("nil %T symbol", symbol)
		}
		s := (*C.dxfg_v4_wildcard_symbol_t)(C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_wildcard_symbol_t{}))))
		s.supper._type = C.WILDCARD
		return &s.supper, nil
	case *osub.TimeSeriesSubscriptionSymbol:
		if value == nil {
			return nil, fmt.Errorf("nil %T symbol", symbol)
		}
		inner, err := newCSymbol(value.Symbol())
		if err != nil {
			return nil, err
		}
		s := (*C.dxfg_v4_time_series_subscription_symbol_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_time_series_subscription_symbol_t{}))))
		s.supper._type = C.TIME_SERIES_SUBSCRIPTION
		s.symbol = inner
		s.from_time = C.int64_t(value.FromTime())
		return &s.supper, nil
	case *osub.IndexedEventSubscriptionSymbol:
		if value == nil {
			return nil, fmt.Errorf("nil %T symbol", symbol)
		}
		source, err := newCIndexedEventSource(value.Source())
		if err != nil {
			return nil, err
		}
		inner, err := newCSymbol(value.Symbol())
		if err != nil {
			freeCIndexedEventSource(source)
			return nil, err
		}
		s := (*C.dxfg_v4_indexed_event_subscription_symbol_t)(
			C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_indexed_event_subscription_symbol_t{}))))
		s.supper._type = C.INDEXED_EVENT_SUBSCRIPTION
		s.symbol = inner
		s.source = source
		return &s.supper, nil
	default:
		return nil, fmt.Errorf("unsupported symbol %T", symbol)
	}
}

func newCIndexedEventSource(source events.IndexedEventSourceInterface) (*C.dxfg_v4_indexed_event_source_t, error) {
	if source == nil {
		return nil, fmt.Errorf("indexed event subscription symbol without a source")
	}
	var sourceType C.dxfg_v4_indexed_event_source_type_t
	switch source.Type() {
	case events.IndexedEventSourceType:
		sourceType = C.INDEXED_EVENT_SOURCE
	case events.OrderSourceType:
		sourceType = C.ORDER_SOURCE
	default:
		return nil, fmt.Errorf("unsupported indexed event source type %d", source.Type())
	}
	name := source.Name()
	if name == nil {
		return nil, fmt.Errorf("indexed event source %d without a name", source.Id())
	}
	s := (*C.dxfg_v4_indexed_event_source_t)(C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_indexed_event_source_t{}))))
	s._type = sourceType
	s.id = C.int32_t(source.Id())
	s.name = C.CString(*name)
	return s, nil
}

func freeCIndexedEventSource(s *C.dxfg_v4_indexed_event_source_t) {
	if s == nil {
		return
	}
	C.free(unsafe.Pointer(s.name))
	C.free(unsafe.Pointer(s))
}

// freeCSymbol frees a symbol built by newCSymbol, its nested symbol and source included.
func freeCSymbol(s *C.dxfg_v4_symbol_t) {
	if s == nil {
		return
	}
	switch s._type {
	case C.STRING:
		C.free(unsafe.Pointer((*C.dxfg_v4_string_symbol_t)(unsafe.Pointer(s)).symbol))
	case C.CANDLE:
		C.free(unsafe.Pointer((*C.dxfg_candle_symbol_t)(unsafe.Pointer(s)).symbol))
	case C.TIME_SERIES_SUBSCRIPTION:
		freeCSymbol((*C.dxfg_v4_time_series_subscription_symbol_t)(unsafe.Pointer(s)).symbol)
	case C.INDEXED_EVENT_SUBSCRIPTION:
		indexed := (*C.dxfg_v4_indexed_event_subscription_symbol_t)(unsafe.Pointer(s))
		freeCSymbol(indexed.symbol)
		freeCIndexedEventSource(indexed.source)
	}
	C.free(unsafe.Pointer(s))
}

// newCSymbolList builds a list of symbols in C memory; free it with freeCSymbolList.
func newCSymbolList(symbols []osub.Symbol) (*C.dxfg_v4_symbol_list, error) {
	list := (*C.dxfg_v4_symbol_list)(C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_symbol_list{}))))
	if len(symbols) == 0 {
		return list, nil
	}
	var element *C.dxfg_v4_symbol_t
	list.elements = (**C.dxfg_v4_symbol_t)(C.calloc(C.size_t(len(symbols)), C.size_t(unsafe.Sizeof(element))))
	list.size = C.int32_t(len(symbols))
	elements := unsafe.Slice(list.elements, len(symbols))
	for i, symbol := range symbols {
		s, err := newCSymbol(symbol)
		if err != nil {
			freeCSymbolList(list)
			return nil, err
		}
		elements[i] = s
	}
	return list, nil
}

// freeCSymbolList frees a list built by newCSymbolList; it also frees a list that newCSymbolList left incomplete.
func freeCSymbolList(list *C.dxfg_v4_symbol_list) {
	if list == nil {
		return
	}
	if list.elements != nil {
		for _, s := range unsafe.Slice(list.elements, int(list.size)) {
			freeCSymbol(s)
		}
		C.free(unsafe.Pointer(list.elements))
	}
	C.free(unsafe.Pointer(list))
}

// stringSymbolText returns the text of a string symbol built by newCSymbol, and false for a symbol of another kind.
func stringSymbolText(symbol *C.dxfg_v4_symbol_t) (string, bool) {
	if symbol._type != C.STRING {
		return "", false
	}
	return C.GoString((*C.dxfg_v4_string_symbol_t)(unsafe.Pointer(symbol)).symbol), true
}
