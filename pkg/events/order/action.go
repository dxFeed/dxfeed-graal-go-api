package order

import "github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"

// Action is the kind of change an order event describes, as in the dxFeed OrderAction enum.
type Action int32

const (
	// ActionUndefined means the action is not defined.
	ActionUndefined = iota
	// ActionNew means a new order is added to the book.
	ActionNew
	// ActionReplace means the order is replaced by a new one.
	ActionReplace
	// ActionModify means the order is modified in place (for example, its size changes).
	ActionModify
	// ActionDelete means the order is deleted from the book.
	ActionDelete
	// ActionPartial means the order is partially executed.
	ActionPartial
	// ActionExecute means the order is fully executed.
	ActionExecute
	// ActionTrade means a trade not related to a particular order took place.
	ActionTrade
	// ActionBurst is the BUST action of the dxFeed OrderAction enum: an earlier trade is cancelled (busted).
	ActionBurst
)

var (
	allValues = mathutil.CreateEnumBitMaskArrayByValue(ActionUndefined,
		[]int64{ActionUndefined, ActionNew, ActionReplace, ActionModify, ActionDelete, ActionPartial, ActionExecute, ActionTrade, ActionBurst})
)

// ActionValueOf returns the Action for the given numeric value; unknown values map to ActionUndefined.
func ActionValueOf(value int32) Action {
	return Action(allValues[value])
}
