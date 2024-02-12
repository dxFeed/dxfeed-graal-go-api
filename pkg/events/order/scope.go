package order

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
)

type Scope int32

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
	ScopeComposite = iota
	ScopeRegional
	ScopeAggregate
	ScopeOrder
)

var (
	scopeAllValues = mathutil.CreateEnumBitMaskArrayByValue(ScopeComposite,
		[]int64{ScopeComposite, ScopeRegional, ScopeAggregate, ScopeOrder})
)

func ScopeValueOf(value int32) Scope {
	return Scope(scopeAllValues[value])
}
