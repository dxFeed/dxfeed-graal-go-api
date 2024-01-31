package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"unsafe"
)

type ProfileMapper struct {
}

func (m ProfileMapper) CEvent(event interface{}) unsafe.Pointer {
	profile := event.(*profile.Profile)
	p := (*C.dxfg_profile_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_profile_t{}))))
	p.market_event.event_type.clazz = C.DXFG_EVENT_PROFILE
	p.market_event.event_symbol = C.CString(*profile.EventSymbol())
	p.market_event.event_time = C.int64_t(profile.EventTime())
	if profile.Description() != nil {
		p.description = C.CString(*profile.Description())
	}
	if profile.StatusReason() != nil {
		p.status_reason = C.CString(*profile.StatusReason())
	}
	p.halt_end_time = C.int64_t(profile.HaltEndTime())
	p.halt_start_time = C.int64_t(profile.HaltStartTime())
	p.halt_end_time = C.int64_t(profile.HaltEndTime())
	p.high_limit_price = C.double(profile.HighLimitPrice())
	p.low_limit_price = C.double(profile.LowLimitPrice())
	p.high_52_week_price = C.double(profile.High52WeekPrice())
	p.low_52_week_price = C.double(profile.Low52WeekPrice())
	p.beta = C.double(profile.Beta())
	p.earnings_per_share = C.double(profile.EarningsPerShare())
	p.dividend_frequency = C.double(profile.DividendFrequency())
	p.ex_dividend_amount = C.double(profile.ExDividendAmount())
	p.ex_dividend_day_id = C.int32_t(profile.ExDividendDayId())
	p.shares = C.double(profile.Shares())
	p.free_float = C.double(profile.FreeFloat())
	p.flags = C.int32_t(profile.Flags())
	return unsafe.Pointer(p)
}

func (pr ProfileMapper) GoEvent(nativeEvent unsafe.Pointer) interface{} {
	native := (*C.dxfg_profile_t)(nativeEvent)

	p := profile.NewProfile(C.GoString(native.market_event.event_symbol))
	p.SetEventTime(int64(native.market_event.event_time))
	p.SetDescription(convertString(native.description))
	p.SetStatusReason(convertString(native.status_reason))
	p.SetHaltStartTime(int64(native.halt_start_time))
	p.SetHaltEndTime(int64(native.halt_end_time))
	p.SetHighLimitPrice(float64(native.high_limit_price))
	p.SetLowLimitPrice(float64(native.low_limit_price))
	p.SetHigh52WeekPrice(float64(native.high_52_week_price))
	p.SetLow52WeekPrice(float64(native.low_52_week_price))
	p.SetBeta(float64(native.beta))
	p.SetEarningsPerShare(float64(native.earnings_per_share))
	p.SetDividendFrequency(float64(native.dividend_frequency))
	p.SetExDividendAmount(float64(native.ex_dividend_amount))
	p.SetExDividendDayId(int32(native.ex_dividend_day_id))
	p.SetShares(float64(native.shares))
	p.SetFreeFloat(float64(native.free_float))
	p.SetFlags(int32(native.flags))
	return p
}
