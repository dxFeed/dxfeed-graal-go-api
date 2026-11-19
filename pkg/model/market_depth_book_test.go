package model

import (
	"math"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
)

func bookOrder(t *testing.T, index int64, price float64, source int64) *bookEntry[*order.Order] {
	t.Helper()
	o := order.NewOrder("TEST")
	if err := o.SetIndex(index); err != nil {
		t.Fatalf("SetIndex: %v", err)
	}
	o.SetSide(side.Buy)
	o.SetPrice(price)
	o.SetSize(1)
	return &bookEntry[*order.Order]{order: o, source: source}
}

func indices(orders []*order.Order) []int64 {
	result := make([]int64, len(orders))
	for i, o := range orders {
		result[i] = o.Index()
	}
	return result
}

// TestOrderBookClearBySource follows the C++ test of the removal of the orders of a source from a book.
func TestOrderBookClearBySource(t *testing.T) {
	const ntv, ise = 1, 2
	t.Run("two orders of the source", func(t *testing.T) {
		b := newOrderBook(compareBuyOrders[*order.Order], 0)
		b.insert(bookOrder(t, 1, 2, ntv))
		b.insert(bookOrder(t, 0, 1, ntv))
		b.toSlice()
		b.clearBySource(ntv)
		if !b.changed {
			t.Error("not changed")
		}
		if got := b.toSlice(); len(got) != 0 {
			t.Errorf("orders: %v, want none", indices(got))
		}
	})
	t.Run("orders of the source before an order of another source", func(t *testing.T) {
		b := newOrderBook(compareBuyOrders[*order.Order], 0)
		b.insert(bookOrder(t, 2, 3, ntv))
		b.insert(bookOrder(t, 1, 2, ntv))
		b.insert(bookOrder(t, 0, 1, ise))
		b.clearBySource(ntv)
		if got := indices(b.toSlice()); len(got) != 1 || got[0] != 0 {
			t.Errorf("orders: %v, want [0]", got)
		}
	})
	t.Run("no orders of the source", func(t *testing.T) {
		b := newOrderBook(compareBuyOrders[*order.Order], 0)
		b.insert(bookOrder(t, 0, 1, ise))
		b.toSlice()
		b.clearBySource(ntv)
		if b.changed {
			t.Error("changed")
		}
		if got := b.toSlice(); len(got) != 1 {
			t.Errorf("orders: %v, want one", indices(got))
		}
	})
}

func TestOrderBookOrdering(t *testing.T) {
	buy := newOrderBook(compareBuyOrders[*order.Order], 2)
	sell := newOrderBook(compareSellOrders[*order.Order], 0)
	for i, price := range []float64{2, 3, 1, math.NaN()} {
		buy.insert(bookOrder(t, int64(i), price, 0))
		sell.insert(bookOrder(t, int64(i), price, 0))
	}
	if got := indices(buy.toSlice()); len(got) != 2 || got[0] != 1 || got[1] != 0 {
		t.Errorf("buy orders: %v, want [1 0] (prices 3, 2)", got)
	}
	if got := indices(sell.toSlice()); len(got) != 4 || got[0] != 2 || got[1] != 0 || got[2] != 1 {
		t.Errorf("sell orders: %v, want [2 0 1 ...] (prices 1, 2, 3)", got)
	}
	if buy.insert(bookOrder(t, 1, 3, 0)) {
		t.Error("an equal order is inserted")
	}
	if !buy.erase(bookOrder(t, 1, 3, 0)) || buy.erase(bookOrder(t, 1, 3, 0)) {
		t.Error("erase of an order: wrong result")
	}
}

func TestCompareFloat(t *testing.T) {
	nan := math.NaN()
	for _, c := range []struct {
		a, b float64
		want int
	}{
		{1, 2, -1}, {2, 1, 1}, {1, 1, 0}, {nan, nan, 0}, {nan, 1, 1}, {1, nan, -1}, {math.Copysign(0, -1), 0, -1},
		{0, math.Copysign(0, -1), 1},
	} {
		if got := compareFloat(c.a, c.b); got != c.want {
			t.Errorf("compareFloat(%v, %v): %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestGoroutineID(t *testing.T) {
	id := goroutineID()
	if id <= 0 {
		t.Fatalf("goroutineID: %d", id)
	}
	other := make(chan int64)
	go func() { other <- goroutineID() }()
	if got := <-other; got == id || got <= 0 {
		t.Fatalf("goroutineID of another goroutine: %d, this one %d", got, id)
	}
}

// aggregateOrder returns an order of the given scope with the fields that compareOrders looks at.
func aggregateOrder(t *testing.T, index int64, scope int, size float64) *order.Order {
	t.Helper()
	o := order.NewOrder("TEST")
	if err := o.SetIndex(index); err != nil {
		t.Fatalf("SetIndex: %v", err)
	}
	o.SetScope(order.Scope(scope))
	o.SetSize(size)
	return o
}

// TestCompareOrdersTiebreaks checks the order of the orders of one price: each rule decides when the previous ones
// find the orders equal.
func TestCompareOrdersTiebreaks(t *testing.T) {
	maker1, maker2 := "MM1", "MM2"
	withMaker := func(o *order.Order, maker *string) *order.Order {
		o.SetMarketMaker(maker)
		return o
	}
	withExchange := func(t *testing.T, o *order.Order, code rune) *order.Order {
		t.Helper()
		if err := o.SetExchangeCode(code); err != nil {
			t.Fatalf("SetExchangeCode: %v", err)
		}
		return o
	}
	withTime := func(o *order.Order, timeSequence int64) *order.Order {
		o.SetTimeSequence(timeSequence)
		return o
	}
	for _, c := range []struct {
		name   string
		o1, o2 *order.Order
		want   int
	}{
		{"individual after aggregate", aggregateOrder(t, 1, order.ScopeOrder, 1),
			aggregateOrder(t, 2, order.ScopeAggregate, 1), 1},
		{"aggregate before individual", aggregateOrder(t, 1, order.ScopeAggregate, 1),
			aggregateOrder(t, 2, order.ScopeOrder, 1), -1},
		{"individual by time", withTime(aggregateOrder(t, 2, order.ScopeOrder, 1), 1),
			withTime(aggregateOrder(t, 1, order.ScopeOrder, 1), 2), -1},
		{"individual by index", aggregateOrder(t, 1, order.ScopeOrder, 1), aggregateOrder(t, 2, order.ScopeOrder, 1), -1},
		{"larger size first", aggregateOrder(t, 2, order.ScopeAggregate, 2), aggregateOrder(t, 1, order.ScopeAggregate, 1), -1},
		{"smaller size last", aggregateOrder(t, 1, order.ScopeAggregate, 1), aggregateOrder(t, 2, order.ScopeAggregate, 2), 1},
		{"by time", withTime(aggregateOrder(t, 2, order.ScopeAggregate, 1), 1),
			withTime(aggregateOrder(t, 1, order.ScopeAggregate, 1), 2), -1},
		{"by scope", aggregateOrder(t, 2, order.ScopeRegional, 1), aggregateOrder(t, 1, order.ScopeAggregate, 1), -1},
		{"by exchange", withExchange(t, aggregateOrder(t, 2, order.ScopeAggregate, 1), 'A'),
			withExchange(t, aggregateOrder(t, 1, order.ScopeAggregate, 1), 'B'), -1},
		{"by market maker", withMaker(aggregateOrder(t, 2, order.ScopeAggregate, 1), &maker1),
			withMaker(aggregateOrder(t, 1, order.ScopeAggregate, 1), &maker2), -1},
		{"market maker after none", withMaker(aggregateOrder(t, 1, order.ScopeAggregate, 1), &maker1),
			aggregateOrder(t, 2, order.ScopeAggregate, 1), 1},
		{"by index", aggregateOrder(t, 1, order.ScopeAggregate, 1), aggregateOrder(t, 2, order.ScopeAggregate, 1), -1},
		{"equal", aggregateOrder(t, 1, order.ScopeAggregate, 1), aggregateOrder(t, 1, order.ScopeAggregate, 1), 0},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := compareOrders(c.o1, c.o2); got != c.want {
				t.Errorf("compareOrders: %d, want %d", got, c.want)
			}
			if got := compareOrders(c.o2, c.o1); got != -c.want {
				t.Errorf("compareOrders reversed: %d, want %d", got, -c.want)
			}
		})
	}
}
