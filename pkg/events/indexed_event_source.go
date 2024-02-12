package events

type EventSourceType int32

const (
	IndexedEventSourceType = iota
	OrderSourceType
)

type IndexedEventSourceInterface interface {
	Name() *string
	Id() int64
	Type() EventSourceType
}

type IndexedEventSource struct {
	id   int64
	name string
}

var defaultIndexedEventSource = NewIndexedEventSource(0, "DEFAULT")

func NewIndexedEventSource(id int64, name string) *IndexedEventSource {
	return &IndexedEventSource{id, name}
}

func (source IndexedEventSource) Name() *string {
	return &source.name
}

func (source IndexedEventSource) Id() int64 {
	return source.id
}

func (source IndexedEventSource) Type() EventSourceType {
	return IndexedEventSourceType
}

func DefaultIndexedEventSource() *IndexedEventSource {
	return defaultIndexedEventSource
}
