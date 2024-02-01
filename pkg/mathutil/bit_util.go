package mathutil

import "fmt"

func GetBits(flags int64, mask int64, shift int64) int64 {
	return (flags >> shift) & mask
}

func SetBits(flags int64, mask int64, shift int64, bits int64) int64 {
	return (flags & ^(mask << shift)) | ((bits & mask) << shift)
}

func CheckChar(value int64, mask int64, name string) error {
	if (value & ^mask) != 0 {
		return fmt.Errorf("invalid %s: %d", name, value)
	}
	return nil
}
