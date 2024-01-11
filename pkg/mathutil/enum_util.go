package mathutil

func CreateEnumBitMaskArrayByValue(defaultValue int64, allCases []int64) []int64 {
	length := int(RoundUp(int64(len(allCases)), 2))

	result := make([]int64, length)

	for i := 0; i < length; i++ {
		if i >= len(allCases) {
			result[i] = defaultValue
		} else {
			result[i] = allCases[i]
		}
	}
	return result
}
