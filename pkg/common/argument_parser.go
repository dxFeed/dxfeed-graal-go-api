package common

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/Osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"strings"
)

func ParseSymbols(value string) []any {
	var symbols []any
	symbolStr := strings.Split(value, ",")
	for _, element := range symbolStr {
		element = strings.TrimSpace(element)
		if element == "all" {
			symbols = append(symbols, Osub.NewWildcardSymbol())
		} else {
			symbols = append(symbols, element)
		}
	}
	return symbols
}

func ParseEventTypes(value string) []eventcodes.EventCode {
	if value == "all" {
		return []eventcodes.EventCode{eventcodes.Quote, eventcodes.TimeAndSale, eventcodes.Profile}
	}
	typeStr := strings.Split(value, ",")
	var types []eventcodes.EventCode

	for _, element := range typeStr {
		element = strings.TrimSpace(element)
		if strings.ToLower(element) == "quote" {
			types = append(types, eventcodes.Quote)
		} else if strings.ToLower(element) == "timeandsale" {
			types = append(types, eventcodes.TimeAndSale)
		} else if strings.ToLower(element) == "profile" {
			types = append(types, eventcodes.Profile)
		}
	}
	return types
}
