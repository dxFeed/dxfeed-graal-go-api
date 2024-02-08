package trade

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"

type TradeETH struct {
	TradeBase
}

func NewTradeETH(eventSymbol string) *TradeETH {
	return &TradeETH{TradeBase: *NewTradeBase(eventSymbol)}
}

func (t *TradeETH) Type() eventcodes.EventCode {
	return eventcodes.TradeETH
}

func (t *TradeETH) String() string {
	return "TradeETH{" +
		t.TradeBase.String() +
		"}"
}
