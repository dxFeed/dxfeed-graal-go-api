package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/summary"
)

type SummaryMapper struct{}

func (m SummaryMapper) CEvent(event interface{}) unsafe.Pointer {
	s := event.(*summary.Summary)

	n := (*C.dxfg_summary_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_summary_t{}))))
	n.market_event.event_type.clazz = C.DXFG_EVENT_SUMMARY
	n.market_event.event_symbol = C.CString(*s.EventSymbol())
	n.market_event.event_time = C.int64_t(s.EventTime())
	n.day_id = C.int32_t(s.DayId())
	n.day_open_price = C.double(s.DayOpenPrice())
	n.day_high_price = C.double(s.DayHighPrice())
	n.day_low_price = C.double(s.DayLowPrice())
	n.day_close_price = C.double(s.DayClosePrice())
	n.prev_day_id = C.int32_t(s.PrevDayId())
	n.prev_day_close_price = C.double(s.PrevDayClosePrice())
	n.prev_day_volume = C.double(s.PrevDayVolume())
	n.open_interest = C.int64_t(s.OpenInterest())
	n.flags = C.int32_t(s.Flags())
	return unsafe.Pointer(n)
}

func (m SummaryMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_summary_t)(native)

	s := summary.NewSummary(C.GoString(n.market_event.event_symbol))
	s.SetEventTime(int64(n.market_event.event_time))
	s.SetDayId(int32(n.day_id))
	s.SetDayOpenPrice(float64(n.day_open_price))
	s.SetDayHighPrice(float64(n.day_high_price))
	s.SetDayLowPrice(float64(n.day_low_price))
	s.SetDayClosePrice(float64(n.day_close_price))
	s.SetPrevDayId(int32(n.prev_day_id))
	s.SetPrevDayClosePrice(float64(n.prev_day_close_price))
	s.SetPrevDayVolume(float64(n.prev_day_volume))
	s.SetOpenInterest(int64(n.open_interest))
	s.SetFlags(int32(n.flags))
	return s
}
