// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package candle_test

import (
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
)

// The symbol of a candle names the instrument and the aggregation: here 5-minute candles of the regular trading
// hours. The same string goes into the subscription, for example as
// osub.NewTimeSeriesSubscriptionSymbol("AAPL{=5m,tho=true}", fromTime).
func ExampleCandle_EventSymbol() {
	c := candle.NewCandle("AAPL{=5m,tho=true}")
	fmt.Println(c.EventSymbol())
	// Output: AAPL{=5m,tho=true}
}
