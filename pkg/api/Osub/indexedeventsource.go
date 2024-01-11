package Osub

type IndexedEventSource struct {
	id   int
	name string
}

var DefaultIndexedEventSource = NewIndexedEventSource(0, "DEFAULT")

func NewIndexedEventSource(id int, name string) IndexedEventSource {
	return IndexedEventSource{id, name}
}

func (source IndexedEventSource) GetName() string {
	return source.name
}
