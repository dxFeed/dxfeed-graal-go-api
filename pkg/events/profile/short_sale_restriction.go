package profile

import "github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"

// ShortSaleRestriction is the short sale restriction status of a security instrument.
type ShortSaleRestriction int64

const (
	// UndefinedShortSaleRestriction means the restriction status is unknown or undefined.
	UndefinedShortSaleRestriction = 0
	// ActiveShortSaleRestriction means short sale is restricted.
	ActiveShortSaleRestriction = 1
	// InactiveShortSaleRestriction means short sale is not restricted.
	InactiveShortSaleRestriction = 2
)

var (
	allShortSaleRestrictionValues = mathutil.CreateEnumBitMaskArrayByValue(UndefinedShortSaleRestriction,
		[]int64{UndefinedShortSaleRestriction, ActiveShortSaleRestriction, InactiveShortSaleRestriction})
)

// ShortSaleRestrictionValueOf returns the ShortSaleRestriction for the given integer value.
func ShortSaleRestrictionValueOf(value int64) ShortSaleRestriction {
	return ShortSaleRestriction(allShortSaleRestrictionValues[value])
}
