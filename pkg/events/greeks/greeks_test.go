package greeks_test

import (
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
)

func TestGreeksLocalHub(t *testing.T) {
	ep, err := api.CreateEndpoint(api.LocalHub)
	mustNoErr(t, err)
	defer func() { _ = ep.CloseAndAwaitTermination() }()

	feed, err := ep.GetFeed()
	mustNoErr(t, err)
	pub, err := ep.GetPublisher()
	mustNoErr(t, err)

	sub, err := feed.CreateSubscription(eventcodes.Greeks)
	mustNoErr(t, err)
	defer sub.Close()

	ch := make(chan []interface{}, 4)
	mustNoErr(t, sub.AddListener(&chanListener{ch: ch}))

	const sym = "HUB_GREEKS_TEST"
	mustNoErr(t, sub.AddSymbol(sym))

	// publish
	want := greeks.NewGreeks(sym)
	want.SetEventTime(1_700_000_000_000)
	want.SetEventFlags(0x15)
	mustNoErr(t, want.SetSequence(17))
	want.SetTime(1_700_000_123_456)
	want.SetPrice(12.34)
	want.SetVolatility(0.28)
	want.SetDelta(0.55)
	want.SetGamma(0.03)
	want.SetTheta(-0.02)
	want.SetRho(0.01)
	want.SetVega(0.15)

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
		{"Price", got.Price(), want.Price()},
		{"Volatility", got.Volatility(), want.Volatility()},
		{"Delta", got.Delta(), want.Delta()},
		{"Gamma", got.Gamma(), want.Gamma()},
		{"Theta", got.Theta(), want.Theta()},
		{"Rho", got.Rho(), want.Rho()},
		{"Vega", got.Vega(), want.Vega()},
	} {
		assertF64(t, tc.name, tc.got, tc.want)
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

func receiveOne(t *testing.T, ch <-chan []interface{}, timeout time.Duration) *greeks.Greeks {
	t.Helper()
	select {
	case batch := <-ch:
		if len(batch) != 1 {
			t.Fatalf("expected 1 event, got %d", len(batch))
		}
		g, ok := batch[0].(*greeks.Greeks)
		if !ok {
			t.Fatalf("expected *greeks.Greeks, got %T", batch[0])
		}
		return g
	case <-time.After(timeout):
		t.Fatal("timeout waiting for Greeks event")
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
