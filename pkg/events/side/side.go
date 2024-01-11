package side

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type Side int64

const (
	Undefined = 0
	Buy       = 1
	Sell      = 2
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(Undefined, []int64{Undefined, Buy, Sell})
)

func SideValueOf(value int64) Side {
	return Side(allValues[value])
}
