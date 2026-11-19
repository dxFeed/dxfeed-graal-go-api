// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package market provides market-related enumerations such as PriceType.
package market

// PriceType is the type of a price value (regular, indicative, preliminary or final) used by candle price types.
type PriceType int32

const (
	// Regular is the ordinary price of a trade or quote.
	Regular PriceType = iota
	// Indicative is a non-binding indicative price, for example an opening indication.
	Indicative
	// Preliminary is a price that may still be revised.
	Preliminary
	// Final is a price that is final for the period.
	Final
)

// Code returns the numeric code of the price type.
func (p PriceType) Code() int32 {
	return int32(p)
}

// PriceTypeValueOf returns the PriceType for a numeric code, or Regular if the code is unknown.
func PriceTypeValueOf(code int64) PriceType {
	if code >= 0 && code <= 3 {
		return PriceType(code)
	}
	return Regular
}

// String returns the name of the price type.
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
