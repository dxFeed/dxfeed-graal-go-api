package native

import (
	"dxfeed-graal-go-api/pkg/events"
	"testing"
)

func TestCreateEmptyEventClazzList(t *testing.T) {
	l := createEventClazzList()

	if l.size != 0 {
		t.Fatalf(`createEventClazzList() should return a list with size 0, but size is %d`, (int)(l.size))
	}

	if l.elements != nil {
		t.Fatalf(`createEventClazzList() should return list with elements == nil, but elements = %p`, l.elements)
	}
}

func TestCreateEventClazzListWithRepeatingElements(t *testing.T) {
	l := createEventClazzList((int32)(events.Candle), (int32)(events.Candle), (int32)(events.Order))

	if l.size != 3 {
		t.Fatalf(`createEventClazzList() must not remove repetitive elements`)
	}
}
