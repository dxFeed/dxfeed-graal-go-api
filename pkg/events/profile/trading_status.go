// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package profile

import "github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"

// TradingStatus is the trading status of a security instrument.
type TradingStatus int64

const (
	// UndefinedTradingStatus means the trading status is unknown or undefined.
	UndefinedTradingStatus = 0
	// HaltedTradingStatus means trading is halted.
	HaltedTradingStatus = 1
	// ActiveTradingStatus means trading is active.
	ActiveTradingStatus = 2
)

var (
	allTradingStatusValues = mathutil.CreateEnumBitMaskArrayByValue(UndefinedTradingStatus,
		[]int64{UndefinedTradingStatus, HaltedTradingStatus, ActiveTradingStatus})
)

// TradingStatusValueOf returns the TradingStatus for the given integer value.
func TradingStatusValueOf(value int64) TradingStatus {
	return TradingStatus(allTradingStatusValues[value])
}
