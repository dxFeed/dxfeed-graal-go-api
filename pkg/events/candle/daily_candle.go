package candle

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"

// DailyCandle is a Candle of a daily period: it has every field of Candle and the event code DailyCandle.
//
// A feed delivers daily candles as Candle events: subscribe to eventcodes.Candle with a daily candle symbol, such as
// "AAPL{=d}". A DailyCandle can be published; subscribers of Candle receive it as a Candle.
type DailyCandle struct {
	Candle
}

// NewDailyCandle creates a DailyCandle for the given candle symbol, such as "AAPL{=d}".
func NewDailyCandle(eventSymbol string) *DailyCandle {
	return &DailyCandle{Candle: *NewCandle(eventSymbol)}
}

// Type returns the event code of the daily candle event.
func (c *DailyCandle) Type() eventcodes.EventCode {
	return eventcodes.DailyCandle
}

// String returns a textual representation of the daily candle.
func (c *DailyCandle) String() string {
	return "DailyCandle{" + c.fieldsString() + "}"
}
