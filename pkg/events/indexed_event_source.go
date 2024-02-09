package events

type IndexedEventSourceInterface interface {
	Name() *string
	Id() int64
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

func DefaultIndexedEventSource() *IndexedEventSource {
	return defaultIndexedEventSource
}
