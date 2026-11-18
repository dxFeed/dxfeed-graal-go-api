package trade

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"

// TradeETH is a trade event for the extended trading hours: the latest trade including extended hours and daily statistics.
type TradeETH struct {
	// TradeBase is the base for the Trade and TradeETH events: the price, size and daily statistics of the latest trade.
	TradeBase
}

// NewTradeETH creates a new TradeETH for the given event symbol.
func NewTradeETH(eventSymbol string) *TradeETH {
	return &TradeETH{TradeBase: *NewTradeBase(eventSymbol)}
}

// Type returns the event code of this event.
func (t *TradeETH) Type() eventcodes.EventCode {
	return eventcodes.TradeETH
}

// String returns a string representation of this event.
func (t *TradeETH) String() string {
	return "TradeETH{" +
		t.TradeBase.String() +
		"}"
}
