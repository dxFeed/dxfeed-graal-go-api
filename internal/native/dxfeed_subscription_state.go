package native

/*
#include "dxfg_api.h"
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
)

// Symbols returns a copy of the symbols of the subscription, without the decoration of a time series subscription.
func (s DXFeedSubscription) Symbols() ([]osub.Symbol, error) {
	return s.symbols(func(thread *isolateThread, list **C.dxfg_v4_symbol_list) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_getSymbols(thread.ptr, s.ptr, list)
	})
}

// DecoratedSymbols returns a copy of the symbols of the subscription as the feed sees them.
func (s DXFeedSubscription) DecoratedSymbols() ([]osub.Symbol, error) {
	return s.symbols(func(thread *isolateThread, list **C.dxfg_v4_symbol_list) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_getDecoratedSymbols(thread.ptr, s.ptr, list)
	})
}

func (s DXFeedSubscription) symbols(call func(thread *isolateThread, list **C.dxfg_v4_symbol_list) C.int32_t) (
	[]osub.Symbol, error) {
	var result []osub.Symbol
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_symbol_list
		if err := checkStatus(thread, call(thread, &list)); err != nil {
			return err
		}
		symbols, err := takeSymbolList(thread, list)
		result = symbols
		return err
	})
	return result, err
}

// SetSymbols replaces the symbols of the subscription.
func (s DXFeedSubscription) SetSymbols(symbols ...osub.Symbol) error {
	return s.withSymbols(symbols, func(thread *isolateThread, list *C.dxfg_v4_symbol_list) C.int32_t {
		return C.dxfg_v4_DXFeedSubscription_setSymbols(thread.ptr, s.ptr, list)
	})
}

// EventTypes returns the event classes of the subscription.
func (s DXFeedSubscription) EventTypes() ([]int32, error) {
	var result []int32
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var list *C.dxfg_v4_event_clazz_list
		if err := checkStatus(thread, C.dxfg_v4_DXFeedSubscription_getEventTypes(thread.ptr, s.ptr, &list)); err != nil {
			return err
		}
		result = takeEventClazzList(thread, list)
		return nil
	})
	return result, err
}

// ContainsEventType reports whether the subscription has the event class.
func (s DXFeedSubscription) ContainsEventType(eventType int32) (bool, error) {
	var contains C.int32_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_containsEventType(thread.ptr, s.ptr,
			C.dxfg_v4_event_clazz_t(eventType), &contains))
	})
	return contains != 0, err
}

// AggregationPeriodMillis returns the aggregation period of the subscription in milliseconds, and false if none was
// set.
func (s DXFeedSubscription) AggregationPeriodMillis() (int64, bool, error) {
	var millis C.int64_t
	set := false
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var period *C.dxfg_v4_time_period_t
		if err := checkStatus(thread, C.dxfg_v4_DXFeedSubscription_getAggregationPeriod(thread.ptr, s.ptr, &period)); err != nil {
			return err
		}
		if period == nil {
			return nil
		}
		// Only the presence of the period is needed: the millis come from the other function.
		if err := checkStatus(thread, C.dxfg_v4_JavaObjectHandler_release(thread.ptr,
			(*C.dxfg_v4_java_object_handler)(unsafe.Pointer(period)))); err != nil {
			return err
		}
		set = true
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_getAggregationPeriodMillis(thread.ptr, s.ptr, &millis))
	})
	return int64(millis), set, err
}

// SetAggregationPeriodMillis sets the aggregation period of the subscription in milliseconds.
func (s DXFeedSubscription) SetAggregationPeriodMillis(millis int64) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_setAggregationPeriodMillis(thread.ptr, s.ptr,
			C.int64_t(millis)))
	})
}

// ClearAggregationPeriod withdraws the aggregation period of the subscription.
func (s DXFeedSubscription) ClearAggregationPeriod() error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_setAggregationPeriod(thread.ptr, s.ptr, nil))
	})
}

// AggregationPeriodInfo returns the bounds of the aggregation period applied to the delivered events in milliseconds,
// negative while unknown.
func (s DXFeedSubscription) AggregationPeriodInfo() (minMillis, maxMillis int64, err error) {
	var lo, hi C.int64_t
	err = executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_getAggregationPeriodInfo(thread.ptr, s.ptr, &lo, &hi))
	})
	return int64(lo), int64(hi), err
}

// EventsBatchLimit returns the events batch limit of the subscription.
func (s DXFeedSubscription) EventsBatchLimit() (int32, error) {
	var limit C.int32_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_getEventsBatchLimit(thread.ptr, s.ptr, &limit))
	})
	return int32(limit), err
}

// SetEventsBatchLimit sets the events batch limit of the subscription.
func (s DXFeedSubscription) SetEventsBatchLimit(limit int32) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_setEventsBatchLimit(thread.ptr, s.ptr, C.int32_t(limit)))
	})
}

// SetExecutor sets the executor of the event notifications of the subscription; nil restores the feed's executor.
func (s DXFeedSubscription) SetExecutor(executor *ExecutorHandle) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_setExecutor(thread.ptr, s.ptr, executor.ptr()))
	})
}

// TotalDroppedEvents returns the number of the events dropped because of a buffer overflow.
func (s DXFeedSubscription) TotalDroppedEvents() (int64, error) {
	var dropped C.int64_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_getTotalDroppedEvents(thread.ptr, s.ptr, &dropped))
	})
	return int64(dropped), err
}
