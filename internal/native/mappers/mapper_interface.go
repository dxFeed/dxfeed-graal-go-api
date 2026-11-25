// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

import "C"
import (
	"unsafe"
)

type MapperInterface interface {
	GoEvent(native unsafe.Pointer) interface{}
	// CEvent maps the event to C memory taken from the allocator; with a nil allocator free it with FreeCEvent,
	// otherwise with Allocator.Free.
	CEvent(a *Allocator, event interface{}) unsafe.Pointer
	// FreeCEvent frees an event made by CEvent with a nil allocator together with its strings.
	FreeCEvent(event unsafe.Pointer)
}
