package market

type PriceType int32

const (
	Regular PriceType = iota
	Indicative
	Preliminary
	Final
)

func (p PriceType) Code() int32 {
	return int32(p)
}

func PriceTypeValueOf(code int64) PriceType {
	if code >= 0 && code <= 3 {
		return PriceType(code)
	}
	return Regular
}

func (p PriceType) String() string {
	switch p {
	case Regular:
		return "REGULAR"
	case Indicative:
		return "INDICATIVE"
	case Preliminary:
		return "PRELIMINARY"
	case Final:
		return "FINAL"
	default:
		return "REGULAR"
	}
}
