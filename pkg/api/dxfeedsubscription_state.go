// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

// The special values of DXFeedSubscription.SetEventsBatchLimit.
const (
	// EventsBatchLimitOptimal lets the SDK choose the number of events in one notification; it is the default.
	EventsBatchLimitOptimal int32 = 0
	// EventsBatchLimitMax passes all available events in one notification.
	EventsBatchLimitMax int32 = math.MaxInt32
)

// Symbols returns a copy of the symbols of the subscription. A string or candle symbol is an osub.String, and the
// symbol nested in a time series or indexed event subscription symbol is a plain string.
func (s *DXFeedSubscription) Symbols() ([]osub.Symbol, error) {
	return s.sub.Symbols()
}

// DecoratedSymbols returns a copy of the symbols of the subscription as the feed sees them. It differs from Symbols
// only for a subscription that decorates its symbols, such as a time series subscription in the Java API, which wraps
// each symbol into a time series subscription symbol; see Symbols for the types of the symbols.
func (s *DXFeedSubscription) DecoratedSymbols() ([]osub.Symbol, error) {
	return s.sub.DecoratedSymbols()
}

// SetSymbols replaces the symbols of the subscription; with no symbols it is the same as Clear. See DXFeedSubscription
// for the symbol types.
func (s *DXFeedSubscription) SetSymbols(symbols ...osub.Symbol) error {
	return s.sub.SetSymbols(symbols...)
}

// EventTypes returns the event types of the subscription.
func (s *DXFeedSubscription) EventTypes() ([]eventcodes.EventCode, error) {
	classes, err := s.sub.EventTypes()
	if err != nil {
		return nil, err
	}
	return eventCodes(classes), nil
}

// ContainsEventType reports whether the subscription has the event type.
func (s *DXFeedSubscription) ContainsEventType(eventType eventcodes.EventCode) (bool, error) {
	return s.sub.ContainsEventType(eventType.NativeCode())
}

// AggregationPeriod returns the aggregation period of the subscription; ok is false when none was set.
func (s *DXFeedSubscription) AggregationPeriod() (period time.Duration, ok bool, err error) {
	millis, ok, err := s.sub.AggregationPeriodMillis()
	if err != nil || !ok {
		return 0, false, err
	}
	return time.Duration(millis) * time.Millisecond, true, nil
}

// SetAggregationPeriod limits the rate of the notifications of the event listeners: the events of each symbol are
// conflated over the period. The new period takes effect after the current one elapses, and the effective period is
// the maximum of this one and the aggregation period of the endpoint. The SDK keeps whole milliseconds: a shorter
// period is rounded up to 1 ms. The period must be between 0 and 24 hours.
func (s *DXFeedSubscription) SetAggregationPeriod(period time.Duration) error {
	millis, err := aggregationPeriodMillis(period)
	if err != nil {
		return err
	}
	return s.sub.SetAggregationPeriodMillis(millis)
}

// ClearAggregationPeriod withdraws the aggregation period set with SetAggregationPeriod.
func (s *DXFeedSubscription) ClearAggregationPeriod() error {
	return s.sub.ClearAggregationPeriod()
}

// AggregationPeriodInfo returns the aggregation period that was applied to the events delivered to the subscription.
func (s *DXFeedSubscription) AggregationPeriodInfo() (AggregationPeriodInfo, error) {
	lo, hi, err := s.sub.AggregationPeriodInfo()
	if err != nil {
		return AggregationPeriodInfo{}, err
	}
	return newAggregationPeriodInfo(lo, hi), nil
}

// EventsBatchLimit returns the maximum number of events in one notification of the event listeners; see
// SetEventsBatchLimit.
func (s *DXFeedSubscription) EventsBatchLimit() (int32, error) {
	return s.sub.EventsBatchLimit()
}

// SetEventsBatchLimit sets the maximum number of events in one notification of the event listeners. It is not a strict
// guarantee: a notification may exceed it. EventsBatchLimitOptimal (the default) lets the SDK choose, and
// EventsBatchLimitMax passes all available events at once. The limit must not be negative.
func (s *DXFeedSubscription) SetEventsBatchLimit(limit int32) error {
	if limit < 0 {
		return fmt.Errorf("dxfeed: events batch limit %d is negative", limit)
	}
	return s.sub.SetEventsBatchLimit(limit)
}

// SetExecutor sets the executor that runs the event listeners of the subscription; nil restores the executor of the
// endpoint. The subscription keeps using the executor after the executor is closed. It fails if the executor is
// closed or is a nil pointer.
func (s *DXFeedSubscription) SetExecutor(executor Executor) error {
	if executor == nil {
		return s.sub.SetExecutor(nil)
	}
	h := executor.handle()
	if h == nil || h.IsFreed() {
		return errors.New("dxfeed: the executor of the subscription is closed")
	}
	return s.sub.SetExecutor(h)
}

// TotalDroppedEvents returns the number of the events that were dropped because of a buffer overflow since the
// subscription was created. It is experimental in the SDK and may change.
func (s *DXFeedSubscription) TotalDroppedEvents() (int64, error) {
	return s.sub.TotalDroppedEvents()
}
