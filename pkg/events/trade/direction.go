package trade

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type Direction int32

const (
	Undefined = iota
	Down
	ZeroDown
	Zero
	ZeroUp
	Up
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(Undefined, []int64{Undefined, Down, ZeroDown, Zero, ZeroUp, Up})
)

func ValueOf(value int64) Direction {
	return Direction(allValues[value])
}
