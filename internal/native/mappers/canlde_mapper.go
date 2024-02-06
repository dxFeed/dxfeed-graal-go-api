package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"unsafe"
)

type CandleMapper struct {
}

func (c CandleMapper) GoEvent(nativeEvent unsafe.Pointer) interface{} {
	candleNative := (*C.dxfg_candle_t)(nativeEvent)

	newCandle := candle.NewCandle(C.GoString(candleNative.event_symbol))
	newCandle.SetEventTime(int64(candleNative.event_time))
	newCandle.SetEventFlags(int32(candleNative.event_flags))
	newCandle.SetIndex(int64(candleNative.index))
	newCandle.SetCount(int64(candleNative.count))

	newCandle.SetOpen(float64(candleNative.open))
	newCandle.SetHigh(float64(candleNative.high))
	newCandle.SetLow(float64(candleNative.low))
	newCandle.SetClose(float64(candleNative.close))
	newCandle.SetVolume(float64(candleNative.volume))
	newCandle.SetVwap(float64(candleNative.vwap))
	newCandle.SetBidVolume(float64(candleNative.bid_volume))

	newCandle.SetAskVolume(float64(candleNative.ask_volume))
	newCandle.SetImpVolatility(float64(candleNative.imp_volatility))
	newCandle.SetOpenInterest(float64(candleNative.open_interest))
	return newCandle
}

func (c CandleMapper) CEvent(event interface{}) unsafe.Pointer {
	candleEvent := event.(*candle.Candle)

	native := (*C.dxfg_candle_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_candle_t{}))))
	native.event_type.clazz = C.DXFG_EVENT_CANDLE
	native.event_symbol = C.CString(*candleEvent.EventSymbol().Symbol())
	native.event_time = C.int64_t(candleEvent.EventTime())
	native.event_flags = C.int32_t(candleEvent.EventFlags())
	native.index = C.int64_t(candleEvent.Index())
	native.count = C.int64_t(candleEvent.Count())
	native.open = C.double(candleEvent.Open())
	native.high = C.double(candleEvent.High())
	native.low = C.double(candleEvent.Low())
	native.close = C.double(candleEvent.Close())
	native.volume = C.double(candleEvent.Volume())
	native.vwap = C.double(candleEvent.Vwap())
	native.bid_volume = C.double(candleEvent.BidVolume())
	native.ask_volume = C.double(candleEvent.AskVolume())
	native.imp_volatility = C.double(candleEvent.ImpVolatility())
	native.open_interest = C.double(candleEvent.OpenInterest())

	return unsafe.Pointer(native)
}
