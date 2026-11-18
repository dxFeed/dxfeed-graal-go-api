// Package trade provides the Trade and TradeETH events with the latest trade price, size and daily statistics.
package trade

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// Trade is a trade event: the price and size of the latest regular trading hours trade and daily statistics.
type Trade struct {
	// TradeBase is the base for the Trade and TradeETH events: the price, size and daily statistics of the latest trade.
	TradeBase
}

// NewTrade creates a new Trade for the given event symbol.
func NewTrade(eventSymbol string) *Trade {
	return &Trade{TradeBase: *NewTradeBase(eventSymbol)}
}

// Type returns the event code of this event.
func (t *Trade) Type() eventcodes.EventCode {
	return eventcodes.Trade
}

// String returns a string representation of this event.
func (t *Trade) String() string {
	return "Trade{" +
		t.TradeBase.String() +
		"}"
}
