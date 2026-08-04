package underlying_test

import (
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/underlying"
)

func TestUnderlyingLocalHub(t *testing.T) {
	ep, err := api.CreateEndpoint(api.LocalHub)
	mustNoErr(t, err)
	defer func() { _ = ep.CloseAndAwaitTermination() }()

	feed, err := ep.GetFeed()
	mustNoErr(t, err)
	pub, err := ep.GetPublisher()
	mustNoErr(t, err)

	sub, err := feed.CreateSubscription(eventcodes.Underlying)
	mustNoErr(t, err)
	defer sub.Close()

	ch := make(chan []interface{}, 4)
	mustNoErr(t, sub.AddListener(&chanListener{ch: ch}))

	const sym = "HUB_UNDERLYING_TEST"
	mustNoErr(t, sub.AddSymbol(sym))

	// publish
	want := underlying.NewUnderlying(sym)
	want.SetEventTime(1_700_000_000_000)
	want.SetEventFlags(0x15)
	mustNoErr(t, want.SetSequence(17))
	// Millis part must be >= 512 so that a 32-bit-wide `millis << 22` in
	// SetTime would overflow and corrupt the index.
	want.SetTime(1_700_000_123_999)
	want.SetVolatility(0.28)
	want.SetFrontVolatility(0.31)
	want.SetBackVolatility(0.25)
	want.SetCallVolume(1234)
	want.SetPutVolume(4321)
	want.SetPutCallRatio(3.5)

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
		{"Index", got.Index(), want.Index()},
		{"Time", got.Time(), want.Time()},
		{"Sequence", int64(got.Sequence()), int64(want.Sequence())},
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
		{"Volatility", got.Volatility(), want.Volatility()},
		{"FrontVolatility", got.FrontVolatility(), want.FrontVolatility()},
		{"BackVolatility", got.BackVolatility(), want.BackVolatility()},
		{"CallVolume", got.CallVolume(), want.CallVolume()},
		{"PutVolume", got.PutVolume(), want.PutVolume()},
		{"OptionVolume", got.OptionVolume(), want.OptionVolume()},
		{"PutCallRatio", got.PutCallRatio(), want.PutCallRatio()},
	} {
		assertF64(t, tc.name, tc.got, tc.want)
	}
}

func TestOptionVolume(t *testing.T) {
	for _, tc := range []struct {
		name            string
		call, put, want float64
	}{
		{"both", 100, 200, 300},
		{"noPut", 100, math.NaN(), 100},
		{"noCall", math.NaN(), 200, 200},
		{"neither", math.NaN(), math.NaN(), math.NaN()},
	} {
		u := underlying.NewUnderlying("SYM")
		u.SetCallVolume(tc.call)
		u.SetPutVolume(tc.put)
		assertF64(t, tc.name, u.OptionVolume(), tc.want)
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

func receiveOne(t *testing.T, ch <-chan []interface{}, timeout time.Duration) *underlying.Underlying {
	t.Helper()
	select {
	case batch := <-ch:
		if len(batch) != 1 {
			t.Fatalf("expected 1 event, got %d", len(batch))
		}
		u, ok := batch[0].(*underlying.Underlying)
		if !ok {
			t.Fatalf("expected *underlying.Underlying, got %T", batch[0])
		}
		return u
	case <-time.After(timeout):
		t.Fatal("timeout waiting for Underlying event")
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
