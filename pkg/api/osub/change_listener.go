// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package osub

// ObservableSubscriptionChangeListener is notified of the changes of the symbols of a subscription: add it with
// DXFeedSubscription.AddChangeListener or ObservableSubscription.AddChangeListener in package api. The calls for one
// subscription do not overlap. They come on the goroutine that changes the subscription or on an SDK thread, while
// the SDK holds the lock of the subscription: a listener must return quickly. A panic in a listener is logged and
// does not end the process.
//
// The symbols are copies that the listener may keep. For a time series subscription they are
// *TimeSeriesSubscriptionSymbol; a string or candle symbol is String.
type ObservableSubscriptionChangeListener interface {
	// SymbolsAdded is called with the symbols added to the subscription. When the listener is added to a subscription
	// that has symbols, it is called with them before the add returns. The subscription of a publisher may call it
	// with symbols it already reported: the data provider is to publish their latest events again.
	SymbolsAdded(symbols []Symbol)
	// SymbolsRemoved is called with the symbols removed from the subscription.
	SymbolsRemoved(symbols []Symbol)
	// SubscriptionClosed is called when the subscription is closed or the listener is removed from it.
	SubscriptionClosed()
}

// ConfigurationChangeListener is implemented by an ObservableSubscriptionChangeListener that is also notified when
// the configuration of the subscription, such as its aggregation period, changes. ConfigurationChanged is called on
// the goroutine that changes the configuration, before that call returns.
type ConfigurationChangeListener interface {
	ConfigurationChanged()
}
