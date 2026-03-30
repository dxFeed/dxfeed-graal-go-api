package summary_test

import (
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/market"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/summary"
)

func TestSummaryLocalHub(t *testing.T) {
	ep, err := api.CreateEndpoint(api.LocalHub)
	mustNoErr(t, err)
	defer func() { _ = ep.CloseAndAwaitTermination() }()

	feed, err := ep.GetFeed()
	mustNoErr(t, err)
	pub, err := ep.GetPublisher()
	mustNoErr(t, err)

	sub, err := feed.CreateSubscription(eventcodes.Summary)
	mustNoErr(t, err)
	defer sub.Close()

	ch := make(chan []interface{}, 4)
	mustNoErr(t, sub.AddListener(&chanListener{ch: ch}))

	const sym = "HUB_SUM_TEST"
	mustNoErr(t, sub.AddSymbol(sym))

	// publish
	want := summary.NewSummary(sym)
	want.SetEventTime(1_700_000_000_000)
	want.SetDayId(19_000)
	want.SetDayOpenPrice(100.5)
	want.SetDayHighPrice(101.25)
	want.SetDayLowPrice(99.75)
	want.SetDayClosePrice(100.875)
	want.SetDayClosePriceType(market.Final)
	want.SetPrevDayId(18_999)
	want.SetPrevDayClosePrice(99.5)
	want.SetPrevDayClosePriceType(market.Regular)
	want.SetPrevDayVolume(1_234_567.89)
	want.SetOpenInterest(42_000)

	mustNoErr(t, pub.Publish([]interface{}{want}))
	_ = ep.AwaitProcessed()

	// receive & verify
	got := receiveOne(t, ch, 10*time.Second)

	if *got.EventSymbol() != *want.EventSymbol() {
		t.Errorf("EventSymbol: got %q, want %q", *got.EventSymbol(), *want.EventSymbol())
	}

	for _, tc := range []struct {
		name string
		got  int64
		want int64
	}{
		{"DayId", int64(got.DayId()), int64(want.DayId())},
		{"PrevDayId", int64(got.PrevDayId()), int64(want.PrevDayId())},
		{"OpenInterest", int64(got.OpenInterest()), int64(want.OpenInterest())},
		{"Flags", int64(got.Flags()), int64(want.Flags())},
	} {
		if tc.got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, tc.got, tc.want)
		}
	}

	for _, tc := range []struct {
		name string
		got  float64
		want float64
	}{
		{"DayOpenPrice", got.DayOpenPrice(), want.DayOpenPrice()},
		{"DayHighPrice", got.DayHighPrice(), want.DayHighPrice()},
		{"DayLowPrice", got.DayLowPrice(), want.DayLowPrice()},
		{"DayClosePrice", got.DayClosePrice(), want.DayClosePrice()},
		{"PrevDayClosePrice", got.PrevDayClosePrice(), want.PrevDayClosePrice()},
		{"PrevDayVolume", got.PrevDayVolume(), want.PrevDayVolume()},
	} {
		assertF64(t, tc.name, tc.got, tc.want)
	}

	if got.DayClosePriceType() != want.DayClosePriceType() {
		t.Errorf("DayClosePriceType: got %v, want %v", got.DayClosePriceType(), want.DayClosePriceType())
	}
	if got.PrevDayClosePriceType() != want.PrevDayClosePriceType() {
		t.Errorf("PrevDayClosePriceType: got %v, want %v", got.PrevDayClosePriceType(), want.PrevDayClosePriceType())
	}
}

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

type chanListener struct {
	ch chan []interface{}
}

func (l *chanListener) Update(events []interface{}) { l.ch <- events }

var _ common.EventListener = (*chanListener)(nil)

func receiveOne(t *testing.T, ch <-chan []interface{}, timeout time.Duration) *summary.Summary {
	t.Helper()
	select {
	case batch := <-ch:
		if len(batch) != 1 {
			t.Fatalf("expected 1 event, got %d", len(batch))
		}
		s, ok := batch[0].(*summary.Summary)
		if !ok {
			t.Fatalf("expected *summary.Summary, got %T", batch[0])
		}
		return s
	case <-time.After(timeout):
		t.Fatal("timeout waiting for Summary event")
		return nil
	}
}

func assertF64(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.IsNaN(want) && math.IsNaN(got) {
		return
	}
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s: got %g, want %g", name, got, want)
	}
}
