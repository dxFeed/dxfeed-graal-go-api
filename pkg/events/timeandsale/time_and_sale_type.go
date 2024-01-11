package timeandsale

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type Type int64

const (
	TypeNew        = 0
	TypeCorrection = 1
	TypeCancel     = 2
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(TypeNew,
		[]int64{TypeNew, TypeCorrection, TypeCancel})
)

func TypeValueOf(value int64) Type {
	return Type(allValues[value])
}
