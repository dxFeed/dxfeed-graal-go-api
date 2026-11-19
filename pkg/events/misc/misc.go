// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package misc provides the miscellaneous events: Message and TextMessage carry application messages,
// Configuration and TextConfiguration carry the configuration of a symbol.
package misc

import (
	"strconv"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/formatutil"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/timeutil"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// symbolTime holds the event symbol and the event time that every event of this package has.
type symbolTime struct {
	eventSymbol *string
	eventTime   int64
}

// EventSymbol returns the symbol of this event.
func (e *symbolTime) EventSymbol() *string {
	return e.eventSymbol
}

// SetEventSymbol sets the symbol of this event.
func (e *symbolTime) SetEventSymbol(eventSymbol string) {
	e.eventSymbol = &eventSymbol
}

// EventTime returns the time, in milliseconds, when the event was received or last modified by the source.
func (e *symbolTime) EventTime() int64 {
	return e.eventTime
}

// SetEventTime sets the event time in milliseconds.
func (e *symbolTime) SetEventTime(eventTime int64) {
	e.eventTime = eventTime
}

// timeSequence holds the time and the sequence of an event packed in one value.
type timeSequence struct {
	timeSequence int64
}

// TimeSequence returns the time and the sequence of this event packed in one value.
func (e *timeSequence) TimeSequence() int64 {
	return e.timeSequence
}

// SetTimeSequence sets the time and the sequence of this event packed in one value.
func (e *timeSequence) SetTimeSequence(value int64) {
	e.timeSequence = value
}

// Time returns the time of this event in milliseconds.
func (e *timeSequence) Time() int64 {
	return timeutil.TimeOfTimeSequence(e.timeSequence)
}

// SetTime sets the time of this event in milliseconds; the sequence is kept.
func (e *timeSequence) SetTime(timeMillis int64) {
	e.timeSequence = timeutil.WithTime(e.timeSequence, timeMillis)
}

// Sequence returns the sequence number of this event, distinguishing events with the same time.
func (e *timeSequence) Sequence() int32 {
	return timeutil.SequenceOfTimeSequence(e.timeSequence)
}

// SetSequence sets the sequence number of this event and returns an error if it is out of range.
func (e *timeSequence) SetSequence(sequence int32) error {
	value, err := timeutil.WithSequence(e.timeSequence, sequence)
	e.timeSequence = value
	return err
}

func (e *timeSequence) timeSequenceString() string {
	return ", time=" + formatutil.FormatTime(e.Time()) +
		", sequence=" + strconv.FormatInt(int64(e.Sequence()), 10)
}

// Message is an application message with a string attachment. Messages are not conflated: every message published
// is delivered.
type Message struct {
	symbolTime
	attachment *string
}

// NewMessage creates a Message for the given event symbol.
func NewMessage(eventSymbol string) *Message {
	return &Message{symbolTime: symbolTime{eventSymbol: &eventSymbol}}
}

// Type returns the event code of this event type.
func (m *Message) Type() eventcodes.EventCode {
	return eventcodes.Message
}

// Attachment returns the attachment of the message, or nil if it has none.
func (m *Message) Attachment() *string {
	return m.attachment
}

// SetAttachment sets the attachment of the message; nil removes it.
func (m *Message) SetAttachment(attachment *string) {
	m.attachment = attachment
}

// String returns a text representation of the event.
func (m *Message) String() string {
	return "Message{" + formatutil.FormatString(m.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(m.EventTime()) +
		", attachment=" + formatutil.FormatString(m.attachment) +
		"}"
}

// TextMessage is an application text message with a time and a sequence. Messages are not conflated: every message
// published is delivered.
type TextMessage struct {
	symbolTime
	timeSequence
	text *string
}

// NewTextMessage creates a TextMessage for the given event symbol.
func NewTextMessage(eventSymbol string) *TextMessage {
	return &TextMessage{symbolTime: symbolTime{eventSymbol: &eventSymbol}}
}

// Type returns the event code of this event type.
func (m *TextMessage) Type() eventcodes.EventCode {
	return eventcodes.TextMessage
}

// Text returns the text of the message, or nil if it has none.
func (m *TextMessage) Text() *string {
	return m.text
}

// SetText sets the text of the message.
func (m *TextMessage) SetText(text *string) {
	m.text = text
}

// String returns a text representation of the event.
func (m *TextMessage) String() string {
	return "TextMessage{" + formatutil.FormatString(m.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(m.EventTime()) +
		m.timeSequenceString() +
		", text=" + formatutil.FormatString(m.text) +
		"}"
}

// Configuration is the configuration of a symbol: a versioned string attachment. Only the last configuration is
// kept for a symbol.
type Configuration struct {
	symbolTime
	version    int32
	attachment *string
}

// NewConfiguration creates a Configuration for the given event symbol.
func NewConfiguration(eventSymbol string) *Configuration {
	return &Configuration{symbolTime: symbolTime{eventSymbol: &eventSymbol}}
}

// Type returns the event code of this event type.
func (c *Configuration) Type() eventcodes.EventCode {
	return eventcodes.Configuration
}

// Version returns the version of the configuration.
func (c *Configuration) Version() int32 {
	return c.version
}

// SetVersion sets the version of the configuration.
func (c *Configuration) SetVersion(version int32) {
	c.version = version
}

// Attachment returns the attachment of the configuration, or nil if it has none.
func (c *Configuration) Attachment() *string {
	return c.attachment
}

// SetAttachment sets the attachment of the configuration; nil removes it.
func (c *Configuration) SetAttachment(attachment *string) {
	c.attachment = attachment
}

// String returns a text representation of the event.
func (c *Configuration) String() string {
	return "Configuration{" + formatutil.FormatString(c.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(c.EventTime()) +
		", version=" + strconv.FormatInt(int64(c.version), 10) +
		", attachment=" + formatutil.FormatString(c.attachment) +
		"}"
}

// TextConfiguration is the configuration of a symbol as versioned text with a time and a sequence. Only the last
// configuration is kept for a symbol.
type TextConfiguration struct {
	symbolTime
	timeSequence
	version int32
	text    *string
}

// NewTextConfiguration creates a TextConfiguration for the given event symbol.
func NewTextConfiguration(eventSymbol string) *TextConfiguration {
	return &TextConfiguration{symbolTime: symbolTime{eventSymbol: &eventSymbol}}
}

// Type returns the event code of this event type.
func (c *TextConfiguration) Type() eventcodes.EventCode {
	return eventcodes.TextConfiguration
}

// Version returns the version of the configuration.
func (c *TextConfiguration) Version() int32 {
	return c.version
}

// SetVersion sets the version of the configuration.
func (c *TextConfiguration) SetVersion(version int32) {
	c.version = version
}

// Text returns the text of the configuration, or nil if it has none.
func (c *TextConfiguration) Text() *string {
	return c.text
}

// SetText sets the text of the configuration.
func (c *TextConfiguration) SetText(text *string) {
	c.text = text
}

// String returns a text representation of the event.
func (c *TextConfiguration) String() string {
	return "TextConfiguration{" + formatutil.FormatString(c.EventSymbol()) +
		", eventTime=" + formatutil.FormatTime(c.EventTime()) +
		c.timeSequenceString() +
		", version=" + strconv.FormatInt(int64(c.version), 10) +
		", text=" + formatutil.FormatString(c.text) +
		"}"
}
