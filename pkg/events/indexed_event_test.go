package events_test

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/marketmaker"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/optionsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/series"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/theoprice"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
)

// The indexed and time series events implement IndexedEvent.
var (
	_ events.IndexedEvent = (*candle.Candle)(nil)
	_ events.IndexedEvent = (*candle.DailyCandle)(nil)
	_ events.IndexedEvent = (*greeks.Greeks)(nil)
	_ events.IndexedEvent = (*order.Order)(nil)
	_ events.IndexedEvent = (*order.AnalyticOrder)(nil)
	_ events.IndexedEvent = (*order.SpreadOrder)(nil)
	_ events.IndexedEvent = (*order.OtcMarketsOrder)(nil)
	_ events.IndexedEvent = (*optionsale.OptionSale)(nil)
	_ events.IndexedEvent = (*marketmaker.MarketMaker)(nil)
	_ events.IndexedEvent = (*series.Series)(nil)
	_ events.IndexedEvent = (*theoprice.TheoPrice)(nil)
	_ events.IndexedEvent = (*timeandsale.TimeAndSale)(nil)
	_ events.IndexedEvent = (*underlying.Underlying)(nil)
)
