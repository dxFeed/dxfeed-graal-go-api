package glossary

import (
	"fmt"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// RoundingMode is the way of rounding a price that is between two valid prices.
type RoundingMode int32

// The rounding modes, as in java.math.RoundingMode.
const (
	// RoundingUp rounds away from zero.
	RoundingUp RoundingMode = iota
	// RoundingDown rounds towards zero.
	RoundingDown
	// RoundingCeiling rounds towards positive infinity.
	RoundingCeiling
	// RoundingFloor rounds towards negative infinity.
	RoundingFloor
	// RoundingHalfUp rounds to the nearest price, away from zero if both are equally near.
	RoundingHalfUp
	// RoundingHalfDown rounds to the nearest price, towards zero if both are equally near.
	RoundingHalfDown
	// RoundingHalfEven rounds to the nearest price, to the even one if both are equally near.
	RoundingHalfEven
	// RoundingUnnecessary fails unless the price is valid already.
	RoundingUnnecessary
)

// PriceIncrements are the minimum price changes of an instrument, which may depend on the price: the text "0.01" is
// one increment for all prices, "0.0001 1; 0.01" is 0.0001 for prices up to 1 and 0.01 above. The zero value has no
// increments.
type PriceIncrements struct {
	text string
}

// ParsePriceIncrements parses the text of price increments, as in the PRICE_INCREMENTS field of a profile. The text is
// normalized, so equal increments are equal values whatever the spacing of the text.
func ParsePriceIncrements(text string) (PriceIncrements, error) {
	values, err := native.PriceIncrementsValues(text)
	if err != nil {
		return PriceIncrements{}, err
	}
	return PriceIncrementsOf(values...)
}

// PriceIncrementsOf returns the price increments of their internal representation: an increment, or increments
// alternating with the price limits up to which they apply, such as 0.0001, 1, 0.01.
func PriceIncrementsOf(values ...float64) (PriceIncrements, error) {
	text, err := native.PriceIncrementsFromValues(values)
	if err != nil {
		return PriceIncrements{}, err
	}
	return PriceIncrements{text}, nil
}

// Text returns the text of the price increments.
func (p PriceIncrements) Text() string {
	return p.text
}

// String returns the text of the price increments.
func (p PriceIncrements) String() string {
	return p.text
}

// Values returns the internal representation of the price increments, as PriceIncrementsOf takes it.
func (p PriceIncrements) Values() ([]float64, error) {
	return native.PriceIncrementsValues(p.text)
}

func (p PriceIncrements) call(op native.PriceIncrementsOp, args native.PriceIncrementsArgs) (float64, error) {
	return native.PriceIncrementsCall(p.text, op, args)
}

// Increment returns the first price increment, which applies to the lowest prices.
func (p PriceIncrements) Increment() (float64, error) {
	return p.call(native.PriceIncrement, native.PriceIncrementsArgs{})
}

// IncrementForPrice returns the price increment for a price. At a limit between two increments the direction
// chooses: a positive one the increment above, a negative one the increment below, 0 the smaller one.
func (p PriceIncrements) IncrementForPrice(price float64, direction int) (float64, error) {
	return p.call(native.PriceIncrementForPrice, native.PriceIncrementsArgs{Price: price, Direction: sign(direction)})
}

// Precision returns the number of decimal digits that represent any valid price.
func (p PriceIncrements) Precision() (int, error) {
	precision, err := p.call(native.PricePrecision, native.PriceIncrementsArgs{})
	return int(precision), err
}

// PrecisionForPrice returns the number of decimal digits that represent the valid prices near a price.
func (p PriceIncrements) PrecisionForPrice(price float64) (int, error) {
	precision, err := p.call(native.PricePrecisionForPrice, native.PriceIncrementsArgs{Price: price})
	return int(precision), err
}

// RoundPrice rounds a price to the nearest valid price.
func (p PriceIncrements) RoundPrice(price float64) (float64, error) {
	return p.call(native.RoundPrice, native.PriceIncrementsArgs{Price: price})
}

// RoundPriceWithDirection rounds a price to a valid price: up for a positive direction, down for a negative one, to
// the nearest for 0.
func (p PriceIncrements) RoundPriceWithDirection(price float64, direction int) (float64, error) {
	return p.call(native.RoundPriceWithDirection, native.PriceIncrementsArgs{Price: price, Direction: sign(direction)})
}

// RoundPriceWithMode rounds a price to a valid price in the rounding mode; RoundingUnnecessary fails for a price that
// is not valid.
func (p PriceIncrements) RoundPriceWithMode(price float64, mode RoundingMode) (float64, error) {
	if mode < RoundingUp || mode > RoundingUnnecessary {
		return 0, fmt.Errorf("glossary: unknown rounding mode %d", mode)
	}
	return p.call(native.RoundPriceWithMode, native.PriceIncrementsArgs{Price: price, RoundingMode: int32(mode)})
}

// IncrementPrice returns the next valid price up for a positive direction or down for a negative one; the price is
// rounded in that direction first. A direction of 0 is an error.
func (p PriceIncrements) IncrementPrice(price float64, direction int) (float64, error) {
	return p.call(native.IncrementPrice, native.PriceIncrementsArgs{Price: price, Direction: sign(direction)})
}

// IncrementPriceWithStep moves a price by a step up for a positive direction or down for a negative one, and rounds
// the result to a valid price in that direction; a step of 0 moves by the price increment. A direction of 0 and a
// negative or NaN step are errors.
func (p PriceIncrements) IncrementPriceWithStep(price float64, direction int, step float64) (float64, error) {
	return p.call(native.IncrementPriceWithStep,
		native.PriceIncrementsArgs{Price: price, Direction: sign(direction), Step: step})
}

func sign(direction int) int32 {
	switch {
	case direction > 0:
		return 1
	case direction < 0:
		return -1
	}
	return 0
}
