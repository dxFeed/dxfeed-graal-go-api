package profile

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type ShortSaleRestriction int64

const (
	UndefinedShortSaleRestriction = 0
	ActiveShortSaleRestriction    = 1
	InactiveShortSaleRestriction  = 2
)

var (
	allShortSaleRestrictionValues = mathutil.CreateEnumBitMaskArrayByValue(UndefinedShortSaleRestriction,
		[]int64{UndefinedShortSaleRestriction, ActiveShortSaleRestriction, InactiveShortSaleRestriction})
)

func ShortSaleRestrictionValueOf(value int64) ShortSaleRestriction {
	return ShortSaleRestriction(allShortSaleRestrictionValues[value])
}
