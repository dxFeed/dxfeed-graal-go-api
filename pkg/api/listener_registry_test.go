package api

import (
	"errors"
	"testing"
)

// fakeNative stands for a native listener in the registry tests.
type fakeNative struct {
	detached int
}

func detachFake(n *fakeNative) error {
	n.detached++
	return nil
}

// TestListenerRegistryRemoveDuringAttach checks that a listener removed while its native listener is being attached
// is detached once, after the attach.
func TestListenerRegistryRemoveDuringAttach(t *testing.T) {
	r := &listenerRegistry[string, *fakeNative]{kind: "test listener"}
	n := &fakeNative{}
	if err := r.add("a", func() (*fakeNative, error) {
		r.remove("a", detachFake)
		if n.detached != 0 {
			t.Error("detached before the attach finished")
		}
		return n, nil
	}, detachFake); err != nil {
		t.Fatalf("add: %v", err)
	}
	if n.detached != 1 {
		t.Errorf("detached %d times, want 1", n.detached)
	}
	if len(r.entries) != 0 {
		t.Errorf("%d entries left, want none", len(r.entries))
	}
}

// TestListenerRegistryCloseDuringAttach checks that a listener whose attach races close is detached by add, and that
// close does not return it.
func TestListenerRegistryCloseDuringAttach(t *testing.T) {
	r := &listenerRegistry[string, *fakeNative]{kind: "test listener"}
	n := &fakeNative{}
	var closed []*fakeNative
	if err := r.add("a", func() (*fakeNative, error) {
		closed = r.close()
		return n, nil
	}, detachFake); err != nil {
		t.Fatalf("add: %v", err)
	}
	if len(closed) != 0 || n.detached != 1 {
		t.Errorf("close returned %d natives and the native was detached %d times; want 0 and 1", len(closed),
			n.detached)
	}
	if err := r.add("b", func() (*fakeNative, error) {
		t.Error("attach called after close")
		return nil, nil
	}, detachFake); err != nil {
		t.Errorf("add after close: %v", err)
	}
}

// TestListenerRegistryAddRemoveClose checks an attach error, a listener added twice and removed once, a remove of a
// listener that was not added, and close.
func TestListenerRegistryAddRemoveClose(t *testing.T) {
	r := &listenerRegistry[string, *fakeNative]{kind: "test listener"}
	failure := errors.New("attach failure")
	if err := r.add("a", func() (*fakeNative, error) { return nil, failure }, detachFake); !errors.Is(err, failure) {
		t.Errorf("add with a failing attach = %v, want %v", err, failure)
	}
	if len(r.entries) != 0 {
		t.Fatalf("%d entries after a failed attach, want none", len(r.entries))
	}
	first, second := &fakeNative{}, &fakeNative{}
	for _, n := range []*fakeNative{first, second} {
		n := n
		if err := r.add("a", func() (*fakeNative, error) { return n, nil }, detachFake); err != nil {
			t.Fatalf("add: %v", err)
		}
	}
	r.remove("not added", detachFake)
	r.remove("a", detachFake)
	if first.detached != 1 || second.detached != 0 {
		t.Errorf("detached %d and %d times, want 1 and 0", first.detached, second.detached)
	}
	if natives := r.close(); len(natives) != 1 || natives[0] != second {
		t.Errorf("close returned %v, want the second native", natives)
	}
}
