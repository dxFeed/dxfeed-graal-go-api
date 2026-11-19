// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"math"
	"sort"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/order"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/side"
)

// OrderEvent is implemented by the order events that a MarketDepthModel keeps: *order.Order, *order.AnalyticOrder,
// *order.OtcMarketsOrder, *order.SpreadOrder and the other order types.
type OrderEvent interface {
	events.IndexedEvent
	// Price returns the price of the order.
	Price() float64
	// Size returns the size of the order.
	Size() float64
	// HasSize reports whether the size is neither zero nor NaN.
	HasSize() bool
	// Side returns the side of the order: buy or sell.
	Side() side.Side
	// Scope returns the scope of the order.
	Scope() order.Scope
	// TimeSequence returns the time and the sequence of the order.
	TimeSequence() int64
	// ExchangeCode returns the exchange code of the order.
	ExchangeCode() rune
}

// bookEntry is an order of a book with the source of the transaction that brought it.
type bookEntry[E OrderEvent] struct {
	order  E
	source int64
}

// orderBook is one side of a book: the orders sorted by compare, and the first depthLimit of them as the listener
// saw them last. The model guards it.
type orderBook[E OrderEvent] struct {
	compare func(o1, o2 E) int
	orders  []*bookEntry[E]
	// snapshot is the first depthLimit orders, as of the last call of toSlice.
	snapshot   []E
	depthLimit int
	changed    bool
}

func newOrderBook[E OrderEvent](compare func(o1, o2 E) int, depthLimit int) orderBook[E] {
	return orderBook[E]{compare: compare, depthLimit: depthLimit}
}

// unbounded reports whether the book has no depth limit: 0 or math.MaxInt.
func (b *orderBook[E]) unbounded() bool {
	return b.depthLimit <= 0 || b.depthLimit == math.MaxInt
}

// withinDepthLimit reports whether the order is not below the last order that the listener saw.
func (b *orderBook[E]) withinDepthLimit(o E) bool {
	return len(b.snapshot) == 0 || b.compare(b.snapshot[len(b.snapshot)-1], o) >= 0
}

// markChanged marks the book as changed if an insertion or a removal of the order changes what the listener sees.
func (b *orderBook[E]) markChanged(o E) {
	if !b.changed && (b.unbounded() || len(b.orders) <= b.depthLimit || b.withinDepthLimit(o)) {
		b.changed = true
	}
}

func (b *orderBook[E]) setDepthLimit(depthLimit int) {
	if b.depthLimit == depthLimit {
		return
	}
	b.depthLimit = depthLimit
	b.changed = true
}

// search returns the position of the order, or of the place for it.
func (b *orderBook[E]) search(o E) int {
	return sort.Search(len(b.orders), func(i int) bool { return b.compare(b.orders[i].order, o) >= 0 })
}

// insert adds the entry; it returns false if the book has an equal order.
func (b *orderBook[E]) insert(e *bookEntry[E]) bool {
	i := b.search(e.order)
	if i < len(b.orders) && b.compare(b.orders[i].order, e.order) == 0 {
		return false
	}
	b.orders = append(b.orders, nil)
	copy(b.orders[i+1:], b.orders[i:])
	b.orders[i] = e
	b.markChanged(e.order)
	return true
}

// erase removes the order equal to the order of the entry; it returns false if there is none.
func (b *orderBook[E]) erase(e *bookEntry[E]) bool {
	i := b.search(e.order)
	if i == len(b.orders) || b.compare(b.orders[i].order, e.order) != 0 {
		return false
	}
	b.orders = append(b.orders[:i], b.orders[i+1:]...)
	b.markChanged(e.order)
	return true
}

// clearBySource removes the orders of the source.
func (b *orderBook[E]) clearBySource(source int64) {
	kept := b.orders[:0]
	for _, e := range b.orders {
		if e.source != source {
			kept = append(kept, e)
		}
	}
	if len(kept) == len(b.orders) {
		return
	}
	clear(b.orders[len(kept):])
	b.orders = kept
	b.changed = true
}

// toSlice returns the first depthLimit orders; the call makes them the snapshot that the listener saw.
func (b *orderBook[E]) toSlice() []E {
	if b.changed {
		b.changed = false
		n := len(b.orders)
		if !b.unbounded() && n > b.depthLimit {
			n = b.depthLimit
		}
		b.snapshot = b.snapshot[:0]
		for _, e := range b.orders[:n] {
			b.snapshot = append(b.snapshot, e.order)
		}
	}
	return append(make([]E, 0, len(b.snapshot)), b.snapshot...)
}

// compareOrders orders the orders of one price: the individual orders after the others, by time and index; the
// others by size (the largest first), time, scope, exchange, market maker and index.
func compareOrders[E OrderEvent](o1, o2 E) int {
	individual1, individual2 := o1.Scope() == order.ScopeOrder, o2.Scope() == order.ScopeOrder
	switch {
	case individual1 && individual2:
		if c := compareInt(o1.TimeSequence(), o2.TimeSequence()); c != 0 {
			return c
		}
		return compareInt(o1.Index(), o2.Index())
	case individual1:
		return 1
	case individual2:
		return -1
	}
	if c := compareFloat(o2.Size(), o1.Size()); c != 0 {
		return c
	}
	if c := compareInt(o1.TimeSequence(), o2.TimeSequence()); c != 0 {
		return c
	}
	if c := compareInt(o1.Scope(), o2.Scope()); c != 0 {
		return c
	}
	if c := compareInt(o1.ExchangeCode(), o2.ExchangeCode()); c != 0 {
		return c
	}
	if m1, ok := any(o1).(interface{ MarketMaker() *string }); ok {
		if m2, ok := any(o2).(interface{ MarketMaker() *string }); ok {
			if c := compareInt(stringOrEmpty(m1.MarketMaker()), stringOrEmpty(m2.MarketMaker())); c != 0 {
				return c
			}
		}
	}
	return compareInt(o1.Index(), o2.Index())
}

// compareBuyOrders orders the buy orders: the highest price first.
func compareBuyOrders[E OrderEvent](o1, o2 E) int {
	if o1.Price() < o2.Price() {
		return 1
	}
	if o1.Price() > o2.Price() {
		return -1
	}
	return compareOrders(o1, o2)
}

// compareSellOrders orders the sell orders: the lowest price first.
func compareSellOrders[E OrderEvent](o1, o2 E) int {
	if o1.Price() < o2.Price() {
		return -1
	}
	if o1.Price() > o2.Price() {
		return 1
	}
	return compareOrders(o1, o2)
}

func compareInt[T ~int32 | ~int64 | ~string](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// compareFloat compares as Double.compare in Java: NaN is above all the numbers and equal to itself, -0 is below 0.
func compareFloat(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	nanA, nanB := math.IsNaN(a), math.IsNaN(b)
	switch {
	case nanA && nanB:
		return 0
	case nanA:
		return 1
	case nanB:
		return -1
	}
	return compareInt(signBit(b), signBit(a))
}

func signBit(value float64) int32 {
	if math.Signbit(value) {
		return 1
	}
	return 0
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
