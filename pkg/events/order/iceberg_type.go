package order

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type IcebergType int32

const (
	Undefined = iota
	Native
	Synthetic
)

var (
	icebergTypeAllValues = mathutil.CreateEnumBitMaskArrayByValue(Undefined,
		[]int64{Undefined, Native, Synthetic})
)

func ValueOf(value int32) IcebergType {
	return IcebergType(icebergTypeAllValues[value])
}
