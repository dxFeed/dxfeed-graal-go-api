package timeandsale

import "github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"

// Type is the time and sale type: new, correction or cancel.
type Type int64

const (
	// TypeNew marks a newly reported trade.
	TypeNew = 0
	// TypeCorrection marks a correction of a previously reported trade.
	TypeCorrection = 1
	// TypeCancel marks the cancellation of a previously reported trade.
	TypeCancel = 2
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(TypeNew,
		[]int64{TypeNew, TypeCorrection, TypeCancel})
)

// TypeValueOf returns the Type for the given integer value.
func TypeValueOf(value int64) Type {
	return Type(allValues[value])
}
