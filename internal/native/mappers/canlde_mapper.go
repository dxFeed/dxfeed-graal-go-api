package mappers

/*
#include "dxfg_api.h"
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
	n := (*C.dxfg_v4_candle_t)(nativeEvent)
	e := candle.NewCandle(C.GoString(n.event_symbol))
	candleFromC(e, n)
	return e
}

func (c CandleMapper) CEvent(event interface{}) unsafe.Pointer {
	n := (*C.dxfg_v4_candle_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_candle_t{})))
	n.event_type.clazz = C.DXFG_EVENT_CANDLE
	candleToC(event.(*candle.Candle), n)
	return unsafe.Pointer(n)
}

// FreeCEvent frees an event made by CEvent.
func (CandleMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_candle_t)(event)
	cFree(unsafe.Pointer(n.event_symbol))
	cFree(event)
}

// DailyCandleMapper maps DailyCandle: the native structure wraps a candle.
type DailyCandleMapper struct{}

func (DailyCandleMapper) GoEvent(nativeEvent unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_daily_candle_t)(nativeEvent)
	e := candle.NewDailyCandle(C.GoString(n.candle.event_symbol))
	candleFromC(&e.Candle, &n.candle)
	return e
}

func (DailyCandleMapper) CEvent(event interface{}) unsafe.Pointer {
	n := (*C.dxfg_v4_daily_candle_t)(cCalloc(unsafe.Sizeof(C.dxfg_v4_daily_candle_t{})))
	n.candle.event_type.clazz = C.DXFG_EVENT_DAILY_CANDLE
	candleToC(&event.(*candle.DailyCandle).Candle, &n.candle)
	return unsafe.Pointer(n)
}

// FreeCEvent frees an event made by CEvent.
func (DailyCandleMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_daily_candle_t)(event)
	cFree(unsafe.Pointer(n.candle.event_symbol))
	cFree(event)
}

// candleFromC copies the fields of a native candle, except the event symbol, to e.
func candleFromC(e *candle.Candle, n *C.dxfg_v4_candle_t) {
	e.SetEventTime(int64(n.event_time))
	e.SetEventFlags(int32(n.event_flags))
	e.SetIndex(int64(n.index))
	e.SetCount(int64(n.count))
	e.SetOpen(float64(n.open))
	e.SetHigh(float64(n.high))
	e.SetLow(float64(n.low))
	e.SetClose(float64(n.close))
	e.SetVolume(float64(n.volume))
	e.SetVwap(float64(n.vwap))
	e.SetBidVolume(float64(n.bid_volume))
	e.SetAskVolume(float64(n.ask_volume))
	e.SetImpVolatility(float64(n.imp_volatility))
	e.SetOpenInterest(float64(n.open_interest))
}

// candleToC copies the fields of e to a native candle; the event symbol is allocated with cString.
func candleToC(e *candle.Candle, n *C.dxfg_v4_candle_t) {
	n.event_symbol = cString(*e.EventSymbol().Symbol())
	n.event_time = C.int64_t(e.EventTime())
	n.event_flags = C.int32_t(e.EventFlags())
	n.index = C.int64_t(e.Index())
	n.count = C.int64_t(e.Count())
	n.open = C.double(e.Open())
	n.high = C.double(e.High())
	n.low = C.double(e.Low())
	n.close = C.double(e.Close())
	n.volume = C.double(e.Volume())
	n.vwap = C.double(e.Vwap())
	n.bid_volume = C.double(e.BidVolume())
	n.ask_volume = C.double(e.AskVolume())
	n.imp_volatility = C.double(e.ImpVolatility())
	n.open_interest = C.double(e.OpenInterest())
}
