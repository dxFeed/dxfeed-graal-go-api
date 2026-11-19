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
