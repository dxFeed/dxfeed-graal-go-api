package order

type IndexedEventSource struct {
	identifier int64
	name       *string
}

func DefaultIndexedEventSource() *IndexedEventSource {
	defaultStr := "DEFAULT"
	return &IndexedEventSource{identifier: 0, name: &defaultStr}
}
