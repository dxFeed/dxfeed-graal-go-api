package side

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
)

type Side int64

func (s Side) String() string {
	switch s {
	case Undefined:
		return "Undefined"
	case Buy:
		return "Buy"
	case Sell:
		return "Sell"
	default:
		return fmt.Sprintf("Side: Wrong value %d", s)
	}
}

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
