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
		return []any{Osub.WildcardSymbol{}}
	}
	values, _ := native.ParseSymbols(value)
	return values
}

func ParseEventTypes(value string) []eventcodes.EventCode {
	if value == "all" {
		return []eventcodes.EventCode{
			eventcodes.Quote,
			eventcodes.TimeAndSale,
			eventcodes.Profile,
			eventcodes.Order,
			eventcodes.SpreadOrder,
		}
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
		} else if strings.ToLower(element) == "order" {
			types = append(types, eventcodes.Order)
		} else if strings.ToLower(element) == "spreadorder" {
			types = append(types, eventcodes.SpreadOrder)
		}
	}
	return types
}
