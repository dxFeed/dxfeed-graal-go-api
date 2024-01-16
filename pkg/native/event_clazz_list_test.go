package native

import "testing"

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
	l := createEventClazzList(0, 0, 12)

	if l.size != 3 {
		t.Fatalf(`createEventClazzList() must not remove repetitive elements`)
	}
}
