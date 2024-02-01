package order

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type Scope int64

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

func ScopeValueOf(value int64) Scope {
	return Scope(scopeAllValues[value])
}
