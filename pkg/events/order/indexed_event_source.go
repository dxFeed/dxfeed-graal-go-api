package order

type IndexedEventSource struct {
	identifier int64
	name       *string
}

var (
	defaultIndexedEventSource = newDefaultIndexedEventSource()
)

func newDefaultIndexedEventSource() *IndexedEventSource {
	defaultStr := "DEFAULT"
	return &IndexedEventSource{identifier: 0, name: &defaultStr}
}

func DefaultIndexedEventSource() *IndexedEventSource {
	return defaultIndexedEventSource
}
