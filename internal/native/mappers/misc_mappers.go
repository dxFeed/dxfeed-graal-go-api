// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package mappers

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/misc"
)

// MessageMapper maps the Message event; the native attachment is a C string, NULL for no attachment.
type MessageMapper struct{}

func (MessageMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*misc.Message)
	n := (*C.dxfg_v4_message_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_message_t{})))
	n.event_type.clazz = C.DXFG_EVENT_MESSAGE
	n.event_symbol = a.cString(*e.EventSymbol())
	n.event_time = C.int64_t(e.EventTime())
	n.attachment = unsafe.Pointer(a.cStringOrNil(e.Attachment()))
	return unsafe.Pointer(n)
}

func (MessageMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_message_t)(native)
	e := misc.NewMessage(C.GoString(n.event_symbol))
	e.SetEventTime(int64(n.event_time))
	e.SetAttachment(convertString((*C.char)(n.attachment)))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (MessageMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_message_t)(event)
	cFree(unsafe.Pointer(n.event_symbol))
	cFree(n.attachment)
	cFree(event)
}

// TextMessageMapper maps the TextMessage event.
type TextMessageMapper struct{}

func (TextMessageMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*misc.TextMessage)
	n := (*C.dxfg_v4_text_message_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_text_message_t{})))
	n.event_type.clazz = C.DXFG_EVENT_TEXT_MESSAGE
	n.event_symbol = a.cString(*e.EventSymbol())
	n.event_time = C.int64_t(e.EventTime())
	n.time_sequence = C.int64_t(e.TimeSequence())
	n.text = a.cStringOrNil(e.Text())
	return unsafe.Pointer(n)
}

func (TextMessageMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_text_message_t)(native)
	e := misc.NewTextMessage(C.GoString(n.event_symbol))
	e.SetEventTime(int64(n.event_time))
	e.SetTimeSequence(int64(n.time_sequence))
	e.SetText(convertString(n.text))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (TextMessageMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_text_message_t)(event)
	cFree(unsafe.Pointer(n.event_symbol))
	cFree(unsafe.Pointer(n.text))
	cFree(event)
}

// ConfigurationMapper maps the Configuration event; the native attachment is a C string, NULL for no attachment.
type ConfigurationMapper struct{}

func (ConfigurationMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*misc.Configuration)
	n := (*C.dxfg_v4_configuration_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_configuration_t{})))
	n.event_type.clazz = C.DXFG_EVENT_CONFIGURATION
	n.event_symbol = a.cString(*e.EventSymbol())
	n.event_time = C.int64_t(e.EventTime())
	n.version = C.int32_t(e.Version())
	n.attachment = unsafe.Pointer(a.cStringOrNil(e.Attachment()))
	return unsafe.Pointer(n)
}

func (ConfigurationMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_configuration_t)(native)
	e := misc.NewConfiguration(C.GoString(n.event_symbol))
	e.SetEventTime(int64(n.event_time))
	e.SetVersion(int32(n.version))
	e.SetAttachment(convertString((*C.char)(n.attachment)))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (ConfigurationMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_configuration_t)(event)
	cFree(unsafe.Pointer(n.event_symbol))
	cFree(n.attachment)
	cFree(event)
}

// TextConfigurationMapper maps the TextConfiguration event.
type TextConfigurationMapper struct{}

func (TextConfigurationMapper) CEvent(a *Allocator, event interface{}) unsafe.Pointer {
	e := event.(*misc.TextConfiguration)
	n := (*C.dxfg_v4_text_configuration_t)(a.calloc(unsafe.Sizeof(C.dxfg_v4_text_configuration_t{})))
	n.event_type.clazz = C.DXFG_EVENT_TEXT_CONFIGURATION
	n.event_symbol = a.cString(*e.EventSymbol())
	n.event_time = C.int64_t(e.EventTime())
	n.time_sequence = C.int64_t(e.TimeSequence())
	n.version = C.int32_t(e.Version())
	n.text = a.cStringOrNil(e.Text())
	return unsafe.Pointer(n)
}

func (TextConfigurationMapper) GoEvent(native unsafe.Pointer) interface{} {
	n := (*C.dxfg_v4_text_configuration_t)(native)
	e := misc.NewTextConfiguration(C.GoString(n.event_symbol))
	e.SetEventTime(int64(n.event_time))
	e.SetTimeSequence(int64(n.time_sequence))
	e.SetVersion(int32(n.version))
	e.SetText(convertString(n.text))
	return e
}

// FreeCEvent frees an event made by CEvent.
func (TextConfigurationMapper) FreeCEvent(event unsafe.Pointer) {
	n := (*C.dxfg_v4_text_configuration_t)(event)
	cFree(unsafe.Pointer(n.event_symbol))
	cFree(unsafe.Pointer(n.text))
	cFree(event)
}
