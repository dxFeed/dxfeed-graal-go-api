package utils

const (
	NanosInMillis = 1_000_000
)

func GetNanosFromMillisAndNanoPart(timeMillis int64, timeNanoPart int32) int64 {
	return timeMillis*NanosInMillis + int64(timeNanoPart)
}
