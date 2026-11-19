// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package common

// EventListener receives the events of a subscription in batches; see DXFeedSubscription.AddListener in package api.
// It is called from an SDK thread, and the events are of the subscribed types, such as *quote.Quote.
type EventListener interface {
	Update([]interface{})
}
