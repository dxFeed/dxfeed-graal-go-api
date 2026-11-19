// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package order

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

const (
	icebergTypeMask  = 3
	icebergTypeShift = 0
)

// AnalyticOrder is an Order extended with analytic information, such as iceberg order fields.
type AnalyticOrder struct {
	Order
	icebergPeakSize     float64
	icebergHiddenSize   float64
	icebergExecutedSize float64
	icebergFlags        int32
}

// NewAnalyticOrder creates an AnalyticOrder for the given event symbol.
func NewAnalyticOrder(eventSymbol string) *AnalyticOrder {
	return &AnalyticOrder{Order: *NewOrder(eventSymbol)}
}

// IcebergPeakSize returns the peak (visible) size of the iceberg order.
func (o *AnalyticOrder) IcebergPeakSize() float64 {
	return o.icebergPeakSize
}

// SetIcebergPeakSize sets the peak (visible) size of the iceberg order.
func (o *AnalyticOrder) SetIcebergPeakSize(icebergPeakSize float64) {
	o.icebergPeakSize = icebergPeakSize
}

// IcebergHiddenSize returns the hidden size of the iceberg order.
func (o *AnalyticOrder) IcebergHiddenSize() float64 {
	return o.icebergHiddenSize
}

// SetIcebergHiddenSize sets the hidden size of the iceberg order.
func (o *AnalyticOrder) SetIcebergHiddenSize(icebergHiddenSize float64) {
	o.icebergHiddenSize = icebergHiddenSize
}

// IcebergExecutedSize returns the executed size of the iceberg order.
func (o *AnalyticOrder) IcebergExecutedSize() float64 {
	return o.icebergExecutedSize
}

// SetIcebergExecutedSize sets the executed size of the iceberg order.
func (o *AnalyticOrder) SetIcebergExecutedSize(icebergExecutedSize float64) {
	o.icebergExecutedSize = icebergExecutedSize
}

// IcebergFlags returns the raw iceberg flags bit field.
func (o *AnalyticOrder) IcebergFlags() int32 {
	return o.icebergFlags
}

// SetIcebergFlags sets the raw iceberg flags bit field.
func (o *AnalyticOrder) SetIcebergFlags(icebergFlags int32) {
	o.icebergFlags = icebergFlags
}

// IcebergType returns the iceberg type encoded in the iceberg flags.
func (o *AnalyticOrder) IcebergType() IcebergType {
	return ValueOf(int32(mathutil.GetBits(int64(o.IcebergFlags()), icebergTypeMask, icebergTypeShift)))
}

// SetIcebergType sets the iceberg type in the iceberg flags.
func (o *AnalyticOrder) SetIcebergType(value IcebergType) {
	o.SetIcebergFlags(int32(mathutil.SetBits(int64(o.IcebergFlags()), icebergTypeMask, icebergTypeShift, int64(value))))
}

// Type returns the event code of this event type.
func (o *AnalyticOrder) Type() eventcodes.EventCode {
	return eventcodes.AnalyticOrder
}

// String returns a text representation of the event.
func (o *AnalyticOrder) String() string {
	return "AnalyticOrder{" + o.baseFieldsToString() +
		", icebergPeakSize=" + formatutil.FormatFloat64(o.IcebergPeakSize()) +
		", icebergHiddenSize=" + formatutil.FormatFloat64(o.IcebergHiddenSize()) +
		", icebergExecutedSize=" + formatutil.FormatFloat64(o.IcebergExecutedSize()) +
		", icebergType=" + formatutil.FormatInt64(int64((o.IcebergType()))) +
		"}"
}
