// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package side defines the order side enumeration.
package side

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
)

// Side is the side of an order: undefined, buy or sell.
type Side int64

// String returns the name of the side.
func (s Side) String() string {
	switch s {
	case Undefined:
		return "Undefined"
	case Buy:
		return "Buy"
	case Sell:
		return "Sell"
	default:
		return fmt.Sprintf("Side: Wrong value %d", s)
	}
}

const (
	// Undefined means the side is not defined.
	Undefined = 0
	// Buy marks a buy (bid) order.
	Buy = 1
	// Sell marks a sell (ask) order.
	Sell = 2
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(Undefined, []int64{Undefined, Buy, Sell})
)

// SideValueOf returns the Side for a numeric value.
func SideValueOf(value int64) Side {
	return Side(allValues[value])
}
