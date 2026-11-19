// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package events holds what the event types have in common: the [EventType] interface, the indexed events
// ([IndexedEvent], [IndexedEventSource] and the [EventFlag] values), and [InstrumentProfile], the description of an
// instrument that package ipf reads (not an event). The event types themselves are in the packages under events, one
// package per event type or family: quote, trade, summary, profile, order, timeandsale, candle, greeks, theoprice,
// underlying, series, misc; their codes for a subscription are in package eventcodes.
//
// # Kinds of events
//
//   - Lasting events (Quote, Trade, TradeETH, Summary, Profile and others) are the last state of a symbol: a new
//     subscription gets the last known value at once, then the updates.
//   - Indexed events (Order, AnalyticOrder, SpreadOrder, Series, and the time series events below) have several values
//     per symbol, told apart by an index and a source. The [EventFlag] values of an event mark a transaction
//     ([TxPending]), a snapshot ([SnapshotBegin], [SnapshotEnd], [SnapshotSnip]) and a removed event ([RemoveEvent]);
//     [IsSnapshotBegin], [IsSnapshotEnd], [IsSnapshotSnip] and [IsRemove] test them.
//   - Time series events (Candle, TimeAndSale, Greeks, TheoPrice, Underlying) are indexed by time; a subscription with
//     a time series symbol gets their history from a given time.
//
// The events are pointers to the types of these packages, such as *quote.Quote; [EventType.Type] gives the code of
// the type.
package events
