package events_test

import (
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

func TestEventFlagValues(t *testing.T) {
	cases := []struct {
		flag  events.EventFlag
		value int32
		name  string
	}{
		{events.TxPending, 0x01, "TX_PENDING"},
		{events.RemoveEvent, 0x02, "REMOVE_EVENT"},
		{events.SnapshotBegin, 0x04, "SNAPSHOT_BEGIN"},
		{events.SnapshotEnd, 0x08, "SNAPSHOT_END"},
		{events.SnapshotSnip, 0x10, "SNAPSHOT_SNIP"},
		{events.SnapshotMode, 0x40, "SNAPSHOT_MODE"},
		{events.RemoveSymbol, 0x80, "REMOVE_SYMBOL"},
		{events.EventFlag(0x20), 0x20, "0x20"},
	}
	for _, c := range cases {
		if c.flag.Flag() != c.value || c.flag.String() != c.name {
			t.Errorf("flag %#x %q, want %#x %q", c.flag.Flag(), c.flag.String(), c.value, c.name)
		}
	}
}

func TestEventFlagSetClearIn(t *testing.T) {
	flags := events.SnapshotBegin.Set(events.TxPending.Set(0))
	if flags != 0x05 {
		t.Fatalf("flags %#x, want 0x5", flags)
	}
	if !events.TxPending.In(flags) || !events.SnapshotBegin.In(flags) || events.SnapshotEnd.In(flags) {
		t.Errorf("In is wrong for %#x", flags)
	}
	if flags = events.TxPending.Clear(flags); flags != 0x04 {
		t.Errorf("Clear left %#x, want 0x4", flags)
	}
}

func TestEventFlagHelpers(t *testing.T) {
	var e events.IndexedEvent = order.NewOrder("A")
	check := func(flags int32, begin, end, snip, pending, remove bool) {
		t.Helper()
		e.SetEventFlags(flags)
		got := []bool{events.IsSnapshotBegin(e), events.IsSnapshotEnd(e), events.IsSnapshotSnip(e),
			events.IsSnapshotEndOrSnip(e), events.IsPending(e), events.IsRemove(e)}
		want := []bool{begin, end, snip, end || snip, pending, remove}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("flags %#x: helpers %v, want %v", flags, got, want)
				return
			}
		}
	}
	check(0, false, false, false, false, false)
	check(events.SnapshotBegin.Flag(), true, false, false, false, false)
	check(events.SnapshotEnd.Flag()|events.RemoveEvent.Flag(), false, true, false, false, true)
	check(events.SnapshotSnip.Flag()|events.TxPending.Flag(), false, false, true, true, false)
}
