package order

import (
	"bytes"
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"regexp"
	"sync"
)

type Source struct {
	events.IndexedEventSource
	pubFlags  int64
	isBuiltin bool
}

const (
	pubOrder           = 0x0001
	pubAnalyticOrder   = 0x0002
	pubOtcMarketsOrder = 0x0004
	pubSpreadOrder     = 0x0008
	fullOrderBook      = 0x0010
	flagsSize          = 5
)

var (
	sourcesById   sync.Map
	sourcesByName sync.Map
)

func newOrderSourceWithIDName(identifier int64, name *string) *Source {
	return &Source{IndexedEventSource: *events.NewIndexedEventSource(identifier, *name), pubFlags: 0, isBuiltin: false}
}

func newOrderSourceWithIDNameFlags(identifier int64, name string, pubflags int64) *Source {
	err := checkIdAndName(identifier, name)
	if err != nil {
		panic(err)
	}
	// Flag FullOrderBook requires that source must be publishable.
	if (pubflags&fullOrderBook) != 0 &&
		(pubflags&(pubOrder|pubAnalyticOrder|pubSpreadOrder)) == 0 {
		panic("unpublishable full order book order")
	}
	value := &Source{IndexedEventSource: *events.NewIndexedEventSource(identifier, name), pubFlags: pubflags, isBuiltin: true}
	_, loadedById := sourcesById.LoadOrStore(identifier, value)
	if loadedById {
		panic(fmt.Sprintf("duplicate id %d", identifier))
	}
	_, loadedByName := sourcesByName.LoadOrStore(name, value)
	if loadedByName {
		panic(fmt.Sprintf("duplicate name %s", name))
	}
	return value
}

func checkIdAndName(identifier int64, name string) error {
	switch {
	case identifier < 0:
		return fmt.Errorf("id is negative")
	case identifier > 0 && identifier < 0x20 && !IsSpecialSourceId(identifier):
		return fmt.Errorf("id is not marked as special")
	case identifier > 0x20:
		decodedName, err := orderSourceDecodeName(identifier)
		if err != nil {
			return fmt.Errorf("id does not match name")
		}
		composeId, err := orderSourceComposeId(name)
		if err != nil {
			return fmt.Errorf("id does not match name")
		}
		if (identifier != composeId) && (name != *decodedName) {
			return fmt.Errorf("id does not match name")
		}
	default:
		break
	}
	return nil
}

func orderSourceDecodeName(identifier int64) (*string, error) {
	if identifier == 0 {
		return nil, fmt.Errorf("source name must contain from 1 to 4 characters. Current %d", identifier)
	}
	var buffer bytes.Buffer

	for index := 24; index >= 0; index -= 8 {
		if identifier>>index == 0 { // Skip highest contiguous zeros.
			continue
		}
		char := rune((identifier >> index) & 0xff)
		str := fmt.Sprintf("%c", char)
		if !orderSourceCheck(str) {
			return nil, fmt.Errorf("source name must contain only alphanumeric characters")
		}
		buffer.WriteRune(char)
	}
	resultString := buffer.String()
	return &resultString, nil
}

func newOrderSourceName(name string, pubflags int64) *Source {
	id, err := orderSourceComposeId(name)
	if err != nil {
		return nil
	}
	return newOrderSourceWithIDNameFlags(id, name, pubflags)
}

func orderSourceComposeId(name string) (int64, error) {
	var sourceId int64
	count := len(name)
	if count == 0 || count > 4 {
		return 0, nil
	}
	for _, ch := range name {
		str := fmt.Sprintf("%c", ch)
		notAlpha := orderSourceCheck(str)
		if !notAlpha {
			return 0, fmt.Errorf("source name must contain only alphanumeric characters. Current %s", str)
		}
		sourceId = (sourceId << 8) | int64(ch)
	}

	return sourceId, nil
}

func orderSourceCheck(char string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(char)
}

func IsSpecialSourceId(value int64) bool {
	return value >= 1 && value <= 6
}

func ValueOfIdentifier(identifier int64) (*Source, error) {
	value, ok := sourcesById.Load(identifier)
	if ok {
		return value.(*Source), nil
	} else {
		name, err := orderSourceDecodeName(identifier)
		if err != nil {
			return nil, err
		}
		source := newOrderSourceWithIDName(identifier, name)
		return source, nil
	}

}

func ValueOfName(name string) (*Source, error) {
	value, ok := sourcesByName.Load(name)
	if ok {
		return value.(*Source), nil
	} else {
		identifier, err := orderSourceComposeId(name)
		if err != nil {
			return nil, err
		}
		source := newOrderSourceWithIDName(identifier, &name)
		return source, nil
	}
}

func (source Source) Type() events.EventSourceType {
	return events.OrderSourceType
}
