package order

import "github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"

type Action int32

const (
	ActionUndefined = iota
	ActionNew
	ActionReplace
	ActionModify
	ActionDelete
	ActionPartial
	ActionExecute
	ActionTrade
	ActionBurst
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(ActionUndefined,
		[]int64{ActionUndefined, ActionNew, ActionReplace, ActionModify, ActionDelete, ActionPartial, ActionExecute, ActionTrade, ActionBurst})
)

func ActionValueOf(value int32) Action {
	return Action(allValues[value])
}
