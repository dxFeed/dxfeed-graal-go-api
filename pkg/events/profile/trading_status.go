package profile

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type TradingStatus int64

const (
	UndefinedTradingStatus = 0
	HaltedTradingStatus    = 1
	ActiveTradingStatus    = 2
)

var (
	allTradingStatusValues = mathutil.CreateEnumBitMaskArrayByValue(UndefinedTradingStatus,
		[]int64{UndefinedTradingStatus, HaltedTradingStatus, ActiveTradingStatus})
)

func TradingStatusValueOf(value int64) TradingStatus {
	return TradingStatus(allTradingStatusValues[value])
}
