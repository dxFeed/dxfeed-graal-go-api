// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package order

import "github.com/dxfeed/dxfeed-graal-go-api/internal/mathutil"

// IcebergType is the type of an iceberg order.
type IcebergType int32

const (
	// Undefined means the iceberg type is not defined.
	Undefined = iota
	// Native means the iceberg order is native, that is, provided by the exchange.
	Native
	// Synthetic means the iceberg order is synthetic, that is, derived from the order book data.
	Synthetic
)

var (
	icebergTypeAllValues = mathutil.CreateEnumBitMaskArrayByValue(Undefined,
		[]int64{Undefined, Native, Synthetic})
)

// ValueOf returns the IcebergType for the given numeric value; unknown values map to Undefined.
func ValueOf(value int32) IcebergType {
	return IcebergType(icebergTypeAllValues[value])
}
