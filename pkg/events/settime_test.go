package events_test

import (
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/candle"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/trade"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
)

// SetTime packs the timestamp as (seconds << 32) | (millis << 22) | sequence.
// If the `millis << 22` shift is evaluated in int32 width before being widened
// to int64, it overflows for any millis >= 512 (512 == 1<<31 >> 22): the value
// goes negative and sign-extension smears 1-bits over the seconds field.
//
// These cases straddle that boundary, so a regression reintroducing the
// narrow shift fails on the ms >= 512 rows while still passing the low ones.
var setTimeCases = []struct {
	name string
	time int64
}{
	{"ms=0", 1_700_000_123_000},
	{"ms=456", 1_700_000_123_456},
	{"ms=511_boundary", 1_700_000_123_511},
	{"ms=512_overflow", 1_700_000_123_512},
	{"ms=999_max", 1_700_000_123_999},
	{"epoch_ms=999", 999},
}

const maxSequence = (1 << 22) - 1

// timeSequenceEvent is implemented by every event that packs time and sequence
// into a single index field.
type timeSequenceEvent interface {
	SetTime(int64)
	Time() int64
}

func assertTimeRoundTrip(t *testing.T, e timeSequenceEvent, wantSeq int64, seqOf func() int64) {
	t.Helper()
	for _, tc := range setTimeCases {
		e.SetTime(tc.time)
		if got := e.Time(); got != tc.time {
			t.Errorf("%s: Time() = %d, want %d", tc.name, got, tc.time)
		}
		if got := seqOf(); got != wantSeq {
			t.Errorf("%s: SetTime clobbered sequence: got %d, want %d", tc.name, got, wantSeq)
		}
	}
}

func TestGreeksSetTime(t *testing.T) {
	e := greeks.NewGreeks("S")
	if err := e.SetSequence(maxSequence); err != nil {
		t.Fatal(err)
	}
	assertTimeRoundTrip(t, e, maxSequence, func() int64 { return int64(e.Sequence()) })
}

func TestUnderlyingSetTime(t *testing.T) {
	e := underlying.NewUnderlying("S")
	if err := e.SetSequence(maxSequence); err != nil {
		t.Fatal(err)
	}
	assertTimeRoundTrip(t, e, maxSequence, func() int64 { return int64(e.Sequence()) })
}

func TestCandleSetTime(t *testing.T) {
	e := candle.NewCandle("S")
	if err := e.SetSequence(maxSequence); err != nil {
		t.Fatal(err)
	}
	assertTimeRoundTrip(t, e, maxSequence, e.Sequence)
}

func TestOrderSetTime(t *testing.T) {
	e := order.NewOrder("S")
	if err := e.SetSequence(maxSequence); err != nil {
		t.Fatal(err)
	}
	assertTimeRoundTrip(t, e, maxSequence, e.Sequence)
}

func TestTimeAndSaleSetTime(t *testing.T) {
	e := timeandsale.NewTimeAndSale("S")
	if err := e.SetSequence(maxSequence); err != nil {
		t.Fatal(err)
	}
	assertTimeRoundTrip(t, e, maxSequence, e.Sequence)
}

func TestTradeSetTime(t *testing.T) {
	e := trade.NewTrade("S")
	if err := e.SetSequence(maxSequence); err != nil {
		t.Fatal(err)
	}
	assertTimeRoundTrip(t, e, maxSequence, e.Sequence)
}
