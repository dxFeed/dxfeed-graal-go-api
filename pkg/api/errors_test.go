package api

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// TestSubscriptionSymbolKinds adds and removes every kind of symbol, one by one and as a list, and checks that an
// invalid symbol is an error, not a panic or a NULL passed to the SDK. A value of a type that is not an osub.Symbol
// does not compile.
func TestSubscriptionSymbolKinds(t *testing.T) {
	feed, _ := newLocalHub(t)
	sub, err := feed.CreateSubscription(eventcodes.Order)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer func() { _ = sub.Close() }()

	symbols := []osub.Symbol{
		osub.String("AAPL"),
		osub.NewWildcardSymbol(),
		osub.NewTimeSeriesSubscriptionSymbol("IBM", 1000),
		osub.NewIndexedEventSubscriptionSymbol(osub.String("MSFT"), order.Default()),
		osub.NewIndexedEventSubscriptionSymbol(osub.NewTimeSeriesSubscriptionSymbol("TSLA", 0), order.Default()),
	}
	for _, symbol := range symbols {
		if err := sub.AddSymbol(symbol); err != nil {
			t.Errorf("AddSymbol(%T): %v", symbol, err)
		}
		if err := sub.RemoveSymbol(symbol); err != nil {
			t.Errorf("RemoveSymbol(%T): %v", symbol, err)
		}
	}
	if err := sub.AddSymbols(symbols...); err != nil {
		t.Errorf("AddSymbols: %v", err)
	}
	if err := sub.RemoveSymbols(symbols...); err != nil {
		t.Errorf("RemoveSymbols: %v", err)
	}
	if err := sub.Clear(); err != nil {
		t.Errorf("Clear: %v", err)
	}

	if err := sub.AddSymbol(nil); err == nil {
		t.Error("AddSymbol(nil) returned no error")
	}
	if err := sub.AddSymbols(osub.String("AAPL"), nil); err == nil {
		t.Error("AddSymbols with a nil symbol returned no error")
	}
	if err := sub.RemoveSymbols(osub.NewTimeSeriesSubscriptionSymbol(42, 0)); err == nil {
		t.Error("RemoveSymbols of a time series symbol of an int returned no error")
	}
	for _, symbol := range []osub.Symbol{(*osub.WildcardSymbol)(nil), (*osub.TimeSeriesSubscriptionSymbol)(nil),
		(*osub.IndexedEventSubscriptionSymbol)(nil), osub.NewIndexedEventSubscriptionSymbol("AAPL", nil),
		(*osub.CandleSymbol)(nil), osub.Candle(nil)} {
		if err := sub.AddSymbol(symbol); err == nil {
			t.Errorf("AddSymbol(%#v) returned no error", symbol)
		}
	}
}

// panickingListener panics on its first update and passes the later ones to a channel.
type panickingListener struct {
	once   sync.Once
	events chan interface{}
}

func (l *panickingListener) Update(events []interface{}) {
	panicked := false
	l.once.Do(func() { panicked = true })
	if panicked {
		panic("listener failure")
	}
	for _, event := range events {
		l.events <- event
	}
}

// TestListenerPanicIsRecovered checks that a panic of a listener is logged and does not end the process, and that the
// listener still receives later events.
func TestListenerPanicIsRecovered(t *testing.T) {
	logged := &syncBuffer{}
	defer log.SetOutput(log.Writer())
	log.SetOutput(logged)

	feed, publisher := newLocalHub(t)
	sub, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	defer func() { _ = sub.Close() }()
	listener := &panickingListener{events: make(chan interface{}, 16)}
	if err := sub.AddListener(listener); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	if err := sub.AddSymbol(osub.String("PANIC/LISTENER")); err != nil {
		t.Fatalf("add symbol: %v", err)
	}

	deadline := time.After(10 * time.Second)
	for price := 1.0; ; price++ {
		// A quote equal to the last one may be dropped, so every quote has a new price.
		q := quote.NewQuote("PANIC/LISTENER")
		q.SetBidPrice(price)
		if err := publisher.Publish([]interface{}{q}); err != nil {
			t.Fatalf("publish: %v", err)
		}
		select {
		case <-listener.events:
			if !strings.Contains(logged.String(), "listener failure") {
				t.Fatalf("the panic was not logged; log: %q", logged.String())
			}
			return
		case <-time.After(100 * time.Millisecond):
		case <-deadline:
			t.Fatal("no event received after the listener panicked")
		}
	}
}

// syncBuffer is a buffer that the logger of a callback thread and the test can use at once.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
