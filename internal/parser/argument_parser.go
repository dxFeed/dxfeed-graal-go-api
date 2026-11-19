// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
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

// ParseEventTypes parses a comma-separated list of event types, such as Quote,TimeAndSale; the names are
// case-insensitive, and "all" is every event type except OrderBase, which has no events of its own. An unknown name is
// an error that names it.
func ParseEventTypes(value string) ([]eventcodes.EventCode, error) {
	if strings.EqualFold(strings.TrimSpace(value), "all") {
		var codes []eventcodes.EventCode
		for _, code := range eventcodes.All() {
			if code != eventcodes.OrderBase {
				codes = append(codes, code)
			}
		}
		return codes, nil
	}
	var codes []eventcodes.EventCode
	var unknown []string
	for _, name := range strings.Split(value, ",") {
		name = strings.TrimSpace(name)
		code, ok := eventTypeByName(name)
		if !ok {
			unknown = append(unknown, strconv.Quote(name))
			continue
		}
		codes = append(codes, code)
	}
	if len(unknown) > 0 {
		return nil, fmt.Errorf("unknown event type %s", strings.Join(unknown, ", "))
	}
	return codes, nil
}

// eventTypeByName finds the event code by the name of its type, ignoring the case.
func eventTypeByName(name string) (eventcodes.EventCode, bool) {
	for _, code := range eventcodes.All() {
		if strings.EqualFold(code.String(), name) {
			return code, true
		}
	}
	return 0, false
}
func ParseTime(time string) (int64, error) {
	value, err := native.ParseTime(time)
	return value, err
}
