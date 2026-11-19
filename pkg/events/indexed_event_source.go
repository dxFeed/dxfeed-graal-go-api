// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package events

// EventSourceType is the kind of an indexed event source.
type EventSourceType int32

const (
	// IndexedEventSourceType is the kind of a plain indexed event source.
	IndexedEventSourceType = iota
	// OrderSourceType is the kind of an order event source.
	OrderSourceType
)

// IndexedEventSourceInterface is implemented by sources of indexed events.
type IndexedEventSourceInterface interface {
	Name() *string
	Id() int64
	Type() EventSourceType
}

// IndexedEventSource identifies a source of indexed events by numeric id and name.
type IndexedEventSource struct {
	id   int64
	name string
}

var defaultIndexedEventSource = NewIndexedEventSource(0, "DEFAULT")

// NewIndexedEventSource creates an indexed event source with the given id and name.
func NewIndexedEventSource(id int64, name string) *IndexedEventSource {
	return &IndexedEventSource{id, name}
}

// Name returns the name of the source.
func (source IndexedEventSource) Name() *string {
	return &source.name
}

// Id returns the numeric identifier of the source.
func (source IndexedEventSource) Id() int64 {
	return source.id
}

// Type returns the kind of the source.
func (source IndexedEventSource) Type() EventSourceType {
	return IndexedEventSourceType
}

// DefaultIndexedEventSource returns the default indexed event source, which has id 0 and name DEFAULT.
func DefaultIndexedEventSource() *IndexedEventSource {
	return defaultIndexedEventSource
}
