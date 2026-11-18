package trade

import "github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"

// Direction is the tick direction of a trade relative to the previous trade price.
type Direction int32

const (
	// Undefined means the direction is unknown or undefined.
	Undefined = iota
	// Down means the price went down.
	Down
	// ZeroDown means the price is unchanged after a previous downtick.
	ZeroDown
	// Zero means the price is unchanged.
	Zero
	// ZeroUp means the price is unchanged after a previous uptick.
	ZeroUp
	// Up means the price went up.
	Up
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(Undefined, []int64{Undefined, Down, ZeroDown, Zero, ZeroUp, Up})
)

// ValueOf returns the Direction for the given integer value.
func ValueOf(value int64) Direction {
	return Direction(allValues[value])
}
