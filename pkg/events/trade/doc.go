// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package trade provides the [Trade] and [TradeETH] events: the last trade of a symbol, in the regular and in the
// extended trading hours, with its price, size, [Direction] (the tick of the price) and the volume of the day. Both are
// lasting events; [TradeBase] holds their common fields.
package trade
