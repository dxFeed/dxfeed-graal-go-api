package native

import "testing"

func TestJavaHandleFreeTwice(t *testing.T) {
	handle := NewJavaHandle(nil)
	if handle.Ptr() != nil {
		t.Fatal("a nil handle is not nil")
	}
	if err := handle.Free(); err != nil {
		t.Fatal(err)
	}
	endpoint, err := NewDXEndpointHandle(5)
	if err != nil {
		t.Fatal(err)
	}
	if err := endpoint.Close(); err != nil {
		t.Fatal(err)
	}
	if endpoint.self.Ptr() != nil {
		t.Fatal("the handle is not released after Close")
	}
	if err := endpoint.Free(); err != nil {
		t.Fatalf("a second release: %v", err)
	}
}
