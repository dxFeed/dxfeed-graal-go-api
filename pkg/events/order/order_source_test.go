package order

import "testing"

func TestOrderSourceComposeId(t *testing.T) {
	l, _ := OrderSourceComposeId("1234")
	expected := int64(825373492)
	if l != expected {
		t.Fatalf(`OrderSourceComposeId should be "%d". But it equals "%d"`, expected, l)
	}

	if !OrderSourceCheck("a") {
		t.Fatalf(`Not char`)
	}

}
