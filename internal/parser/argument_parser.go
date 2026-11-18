package parser

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"strings"
)

// ParseSymbols parses a comma-separated list of symbols; "all" is the wildcard symbol.
func ParseSymbols(value string) ([]osub.Symbol, error) {
	value = strings.TrimSpace(value)

	if value == "all" {
		return []osub.Symbol{osub.NewWildcardSymbol()}, nil
	}
	symbols, err := native.ParseSymbols(value)
	if err != nil {
		return nil, err
	}
	return osub.Strings(symbols...), nil
}

func ParseEventTypes(value string) []eventcodes.EventCode {
	codes := map[string]eventcodes.EventCode{
		"quote":         eventcodes.Quote,
		"summary":       eventcodes.Summary,
		"greeks":        eventcodes.Greeks,
		"underlying":    eventcodes.Underlying,
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
