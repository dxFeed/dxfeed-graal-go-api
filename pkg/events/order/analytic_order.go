package order

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/mathutil"
)

const (
	icebergTypeMask  = 3
	icebergTypeShift = 0
)

type AnalyticOrder struct {
	Order
	icebergPeakSize     float64
	icebergHiddenSize   float64
	icebergExecutedSize float64
	icebergFlags        int32
}

func NewAnalyticOrder(eventSymbol string) *AnalyticOrder {
	return &AnalyticOrder{Order: *NewOrder(eventSymbol)}
}

func (o *AnalyticOrder) IcebergPeakSize() float64 {
	return o.icebergPeakSize
}

func (o *AnalyticOrder) SetIcebergPeakSize(icebergPeakSize float64) {
	o.icebergPeakSize = icebergPeakSize
}

func (o *AnalyticOrder) IcebergHiddenSize() float64 {
	return o.icebergHiddenSize
}

func (o *AnalyticOrder) SetIcebergHiddenSize(icebergHiddenSize float64) {
	o.icebergHiddenSize = icebergHiddenSize
}

func (o *AnalyticOrder) IcebergExecutedSize() float64 {
	return o.icebergExecutedSize
}

func (o *AnalyticOrder) SetIcebergExecutedSize(icebergExecutedSize float64) {
	o.icebergExecutedSize = icebergExecutedSize
}

func (o *AnalyticOrder) IcebergFlags() int32 {
	return o.icebergFlags
}

func (o *AnalyticOrder) SetIcebergFlags(icebergFlags int32) {
	o.icebergFlags = icebergFlags
}

func (o *AnalyticOrder) IcebergType() IcebergType {
	return ValueOf(int32(mathutil.GetBits(int64(o.IcebergFlags()), icebergTypeMask, icebergTypeShift)))
}

func (o *AnalyticOrder) SetIcebergType(value IcebergType) {
	o.SetIcebergFlags(int32(mathutil.SetBits(int64(o.IcebergFlags()), icebergTypeMask, icebergTypeShift, int64(value))))
}

func (o *AnalyticOrder) Type() eventcodes.EventCode {
	return eventcodes.AnalyticOrder
}

func (o *AnalyticOrder) String() string {
	return "AnalyticOrder{" + o.baseFieldsToString() +
		", icebergPeakSize=" + formatutil.FormatFloat64(o.IcebergPeakSize()) +
		", icebergHiddenSize=" + formatutil.FormatFloat64(o.IcebergHiddenSize()) +
		", icebergExecutedSize=" + formatutil.FormatFloat64(o.IcebergExecutedSize()) +
		", icebergType=" + formatutil.FormatInt64(int64((o.IcebergType()))) +
		"}"
}
