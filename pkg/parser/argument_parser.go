package parser

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/Osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"strings"
)

func ParseSymbols(value string) []any {
	value = strings.TrimSpace(value)

	if value == "all" {
		return []any{Osub.NewWildcardSymbol()}
	}
	values, _ := native.ParseSymbols(value)
	return values
}

func ParseEventTypes(value string) []eventcodes.EventCode {
	codes := map[string]eventcodes.EventCode{
		"quote":         eventcodes.Quote,
		"timeandsale":   eventcodes.TimeAndSale,
		"profile":       eventcodes.Profile,
		"order":         eventcodes.Order,
		"spreadorder":   eventcodes.SpreadOrder,
		"candle":        eventcodes.Candle,
		"trade":         eventcodes.Trade,
		"tradeeth":      eventcodes.TradeETH,
		"analyticorder": eventcodes.AnalyticOrder,
	}
	var values []eventcodes.EventCode
	if value == "all" {
		for k := range codes {
			values = append(values, codes[k])
		}
		return values
	}

	typeStr := strings.Split(value, ",")
	for _, element := range typeStr {
		element = strings.ToLower(strings.TrimSpace(element))
		eventCode, exists := codes[element]
		if exists {
			values = append(values, eventCode)
		}
	}
	return values
}

func ParseTime(time string) (int64, error) {
	value, err := native.ParseTime(time)
	return value, err
}
