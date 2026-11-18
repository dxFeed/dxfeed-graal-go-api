package order_test

import (
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
)

// Time and sequence live in timeSequence, as in Java and C++; the index (order source and ID) must not change.
func TestTimeAndSequenceUseTimeSequence(t *testing.T) {
	ntv, err := order.ValueOfName("NTV")
	if err != nil {
		t.Fatal(err)
	}
	o := order.NewOrder("AAPL")
	o.SetOrderSource(ntv)
	if err := o.SetIndex(o.Index() | 12345); err != nil {
		t.Fatal(err)
	}
	index := o.Index()

	o.SetTime(1_700_000_000_123)
	if err := o.SetSequence(7); err != nil {
		t.Fatal(err)
	}

	if o.Index() != index {
		t.Errorf("Index = %#x after SetTime/SetSequence, want %#x", o.Index(), index)
	}
	if src, err := o.OrderSource(); err != nil || *src.Name() != "NTV" {
		t.Errorf("OrderSource = %v, %v; want NTV", src, err)
	}
	if want := int64(1_700_000_000<<32 | 123<<22 | 7); o.TimeSequence() != want {
		t.Errorf("TimeSequence = %#x, want %#x", o.TimeSequence(), want)
	}
}

func TestTimeAndSequenceDecodeTimeSequence(t *testing.T) {
	for _, tc := range []struct {
		name string
		base interface {
			SetTimeSequence(int64)
			Time() int64
			Sequence() int64
		}
	}{
		{"Order", order.NewOrder("AAPL")},
		{"AnalyticOrder", order.NewAnalyticOrder("AAPL")},
		{"SpreadOrder", order.NewSpreadOrder("AAPL")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.base.SetTimeSequence(1_700_000_000<<32 | 123<<22 | 7)
			if tc.base.Time() != 1_700_000_000_123 || tc.base.Sequence() != 7 {
				t.Errorf("Time, Sequence = %d, %d; want 1700000000123, 7", tc.base.Time(), tc.base.Sequence())
			}
		})
	}
}
