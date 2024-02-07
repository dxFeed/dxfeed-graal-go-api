package mathutil

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func FloorDivInt64(x, y int64) int64 {
	r := x / y
	// if the signs are different and modulo not zero, round down
	if x^y < 0 && r*y != x {
		r--
	}
	return r
}

func FloorModInt64(x, y int64) int64 {
	return x - FloorDivInt64(x, y)*y
}

func RoundUp(numToRound int64, multiple int64) int64 {
	if multiple == 0 {
		return numToRound
	}
	remainder := numToRound % multiple
	if remainder == 0 {
		return numToRound
	}
	return numToRound + multiple - remainder
}

func MinInt64(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
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
