// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package order provides the order book events: [Order] (an order or a price level of a source), [AnalyticOrder] (an
// order with the iceberg information) and [SpreadOrder] (an order of a spread), and their attributes: [Source],
// [Action], [Scope], [IcebergType].
//
// Orders are indexed events. A subscription gets the orders of the sources given with
// osub.IndexedEventSubscriptionSymbol, for example order.NtvL2() for the price levels of Nasdaq TotalView, or of all
// sources. [Source] has a function per known source and [ValueOfName] for the others.
package order
