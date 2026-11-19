// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

import "C"
import (
	"unsafe"
)

type MapperInterface interface {
	GoEvent(native unsafe.Pointer) interface{}
	// CEvent maps the event to C memory; free it with FreeCEvent.
	CEvent(event interface{}) unsafe.Pointer
	// FreeCEvent frees an event made by CEvent together with its strings.
	FreeCEvent(event unsafe.Pointer)
}
