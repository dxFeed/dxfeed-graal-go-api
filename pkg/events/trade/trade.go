package trade

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

type Trade struct {
	TradeBase
}

func NewTrade(eventSymbol string) *Trade {
	return &Trade{TradeBase: *NewTradeBase(eventSymbol)}
}

func (t *Trade) Type() eventcodes.EventCode {
	return eventcodes.Trade
}

func (t *Trade) String() string {
	return "Trade{" +
		t.TradeBase.String() +
		"}"
}
