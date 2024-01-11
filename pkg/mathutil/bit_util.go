package mathutil

func GetBits(flags int64, mask int64, shift int64) int64 {
	return (flags >> shift) & mask
}

func SetBits(flags int64, mask int64, shift int64, bits int64) int64 {
	return (flags & ^(mask << shift)) | ((bits & mask) << shift)
}
