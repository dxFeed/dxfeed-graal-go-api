package utils

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
