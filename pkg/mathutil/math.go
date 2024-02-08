package mathutil

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func FloorDivInt[T Integer](x, y T) T {
	r := x / y
	// if the signs are different and modulo not zero, round down
	if x^y < 0 && r*y != x {
		r--
	}
	return r
}

func FloorModInt[T Integer](x, y T) T {
	return x - FloorDivInt(x, y)*y
}

func RoundUp[T Integer](numToRound T, multiple T) T {
	if multiple == 0 {
		return numToRound
	}
	remainder := numToRound % multiple
	if remainder == 0 {
		return numToRound
	}
	return numToRound + multiple - remainder
}

func MinInt[T Integer](a, b T) T {
	if a <= b {
		return a
	}
	return b
}

func MaxInt[T Integer](a, b T) T {
	if a >= b {
		return a
	}
	return b
}

func Div[T Integer](a, b T) T {
	if a >= 0 {
		return a / b
	} else {
		if b >= 0 {
			return (a+1)/b - 1
		} else {
			return (a+1)/b + 1
		}
	}
}

func Abs[T Integer](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
