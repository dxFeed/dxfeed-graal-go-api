// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package order

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
)

// Scope is the scope of an order: composite, regional, aggregate or individual order.
type Scope int32

// String returns the name of the scope.
func (s Scope) String() string {
	switch s {
	case ScopeComposite:
		return "Composite"
	case ScopeRegional:
		return "Regional"
	case ScopeAggregate:
		return "Aggregate"
	case ScopeOrder:
		return "Order"
	default:
		return fmt.Sprintf("Scope: Wrong value %d", s)

	}
}

const (
	// ScopeComposite is the best bid or ask price on the symbol across all exchanges.
	ScopeComposite = iota
	// ScopeRegional is the best bid or ask price on the symbol on a particular exchange.
	ScopeRegional
	// ScopeAggregate is an aggregate information for a price level.
	ScopeAggregate
	// ScopeOrder is an individual order on the market.
	ScopeOrder
)

var (
	scopeAllValues = mathutil.CreateEnumBitMaskArrayByValue(ScopeComposite,
		[]int64{ScopeComposite, ScopeRegional, ScopeAggregate, ScopeOrder})
)

// ScopeValueOf returns the Scope for the given numeric value.
func ScopeValueOf(value int32) Scope {
	return Scope(scopeAllValues[value])
}
