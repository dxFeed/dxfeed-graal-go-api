package order

import "testing"

func TestOrderSourceComposeId(t *testing.T) {
	l, _ := orderSourceComposeId("1234")
	expected := int64(825373492)
	if l != expected {
		t.Fatalf(`orderSourceComposeId should be "%d". But it equals "%d"`, expected, l)
	}
	if !orderSourceCheck("a") {
		t.Fatalf(`Not char`)
	}
}

func TestGetOrder(t *testing.T) {
	value, err := ValueOfIdentifier(1)
	if err != nil {
		t.Fatalf(`Get Source with error "%v".`, err)
	}
	constV := CompsoiteBid()

	if constV != value {
		t.Fatalf(`Value %v doesn't equal "%v".`, constV, value)
	}

	valueByName, err := ValueOfName("COMPOSITE_BID")
	if err != nil {
		t.Fatalf(`Get Source with error "%v".`, err)
	}
	if constV.name != valueByName.name {
		t.Fatalf(`Value %v doesn't equal "%v".`, constV, value)
	}
}
