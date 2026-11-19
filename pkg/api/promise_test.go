package api

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
)

func awaitContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestGetLastEventPromise(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "PROMISE/A"
	published := quote.NewQuote(symbol)
	published.SetBidPrice(7.5)
	subscribeAndPublish(t, feed, publisher, eventcodes.Quote, osub.Strings(symbol), []interface{}{published}, 1)

	promise, err := feed.GetLastEventPromise(eventcodes.Quote, osub.String(symbol))
	if err != nil {
		t.Fatalf("GetLastEventPromise: %v", err)
	}
	got, err := promise.Await(awaitContext(t))
	if err != nil {
		t.Fatalf("Await: %v", err)
	}
	if q, ok := got.(*quote.Quote); !ok || q.BidPrice() != 7.5 || *q.EventSymbol() != symbol {
		t.Fatalf("Await = %v, want the published quote", got)
	}
	if !promise.IsDone() || promise.IsCancelled() {
		t.Fatalf("IsDone = %v, IsCancelled = %v; want true, false", promise.IsDone(), promise.IsCancelled())
	}
	if again, err := promise.Result(); err != nil || again != got {
		t.Fatalf("Result = %v, %v; want the awaited quote", again, err)
	}
	if err := promise.Cancel(); err != nil {
		t.Fatalf("Cancel of a done promise: %v", err)
	}
}

func TestPromiseAwaitCancelsOnContext(t *testing.T) {
	feed, _ := newLocalHub(t)
	promise, err := feed.GetLastEventPromise(eventcodes.Quote, osub.String("PROMISE/NONE"))
	if err != nil {
		t.Fatalf("GetLastEventPromise: %v", err)
	}
	if _, err := promise.Result(); !errors.Is(err, ErrPromiseNotDone) {
		t.Fatalf("Result of a pending promise = %v, want ErrPromiseNotDone", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if _, err := promise.Await(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Await = %v, want context.DeadlineExceeded", err)
	}
	select {
	case <-promise.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("the promise is not done 10 s after the cancellation")
	}
	if !promise.IsCancelled() {
		t.Fatal("IsCancelled = false after Await was cancelled")
	}
	if _, err := promise.Result(); !errors.Is(err, ErrPromiseCancelled) {
		t.Fatalf("Result = %v, want ErrPromiseCancelled", err)
	}
}

func TestGetLastEventsPromises(t *testing.T) {
	feed, publisher := newLocalHub(t)
	symbols := []string{"PROMISE/B", "PROMISE/C"}
	var published []interface{}
	for i, s := range symbols {
		q := quote.NewQuote(s)
		q.SetBidPrice(float64(i + 1))
		published = append(published, q)
	}
	subscribeAndPublish(t, feed, publisher, eventcodes.Quote, osub.Strings(symbols...), published, len(symbols))

	promises, err := feed.GetLastEventsPromises(eventcodes.Quote, osub.Strings(symbols...)...)
	if err != nil {
		t.Fatalf("GetLastEventsPromises: %v", err)
	}
	if len(promises) != len(symbols) {
		t.Fatalf("%d promises, want %d", len(promises), len(symbols))
	}
	for i, p := range promises {
		got, err := p.Await(awaitContext(t))
		if err != nil {
			t.Fatalf("Await %d: %v", i, err)
		}
		if q, ok := got.(*quote.Quote); !ok || *q.EventSymbol() != symbols[i] || q.BidPrice() != float64(i+1) {
			t.Fatalf("promise %d = %v, want the quote of %s", i, got, symbols[i])
		}
	}
	if empty, err := feed.GetLastEventsPromises(eventcodes.Quote); err != nil || len(empty) != 0 {
		t.Fatalf("GetLastEventsPromises without symbols = %v, %v; want no promises", empty, err)
	}
}

func TestGetTimeSeriesPromise(t *testing.T) {
	feed, publisher := newLocalHub(t)
	const symbol = "PROMISE/TS"
	base := time.Now().Add(-time.Hour).UnixMilli() / 1000 * 1000
	var published []interface{}
	for i := 2; i >= 0; i-- {
		e := timeandsale.NewTimeAndSale(symbol)
		e.SetTime(base + int64(i)*1000)
		published = append(published, e)
	}
	published[0].(*timeandsale.TimeAndSale).SetEventFlags(events.SnapshotBegin.Flag())
	published[2].(*timeandsale.TimeAndSale).SetEventFlags(events.SnapshotEnd.Flag())

	promise, err := feed.GetTimeSeriesPromise(eventcodes.TimeAndSale, osub.String(symbol), base, math.MaxInt64)
	if err != nil {
		t.Fatalf("GetTimeSeriesPromise: %v", err)
	}
	subscribeAndPublish(t, feed, publisher, eventcodes.TimeAndSale,
		[]osub.Symbol{osub.NewTimeSeriesSubscriptionSymbol(osub.String(symbol), base)}, published, len(published))
	got, err := promise.Await(awaitContext(t))
	if err != nil {
		t.Fatalf("Await: %v", err)
	}
	if len(got) != len(published) {
		t.Fatalf("%d events, want %d", len(got), len(published))
	}
	for i, e := range got {
		if ts, ok := e.(*timeandsale.TimeAndSale); !ok || ts.Time() != base+int64(i)*1000 {
			t.Fatalf("event %d = %v, want time %d", i, e, base+int64(i)*1000)
		}
	}
}
