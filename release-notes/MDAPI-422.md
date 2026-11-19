* **\[MDAPI-422]\[GO]** Added the feed queries of Graal Native SDK 3.16: last events, events if subscribed,
  promises, attaching subscriptions and the time series subscription.
  * `DXFeed.GetLastEvent`, `GetLastEvents`, `GetLastEventIfSubscribed`, `GetIndexedEventsIfSubscribed` and
    `GetTimeSeriesIfSubscribed` read the local cache of the feed.
  * `DXFeed.GetLastEventPromise`, `GetLastEventsPromises`, `GetIndexedEventsPromise` and `GetTimeSeriesPromise`
    return a `Promise` that `Await(ctx)` waits for; a done context or closing the endpoint cancels the request. The
    `print-last-events` sample shows them.
  * `NewDXFeedSubscription` creates a detached subscription; `DXFeed.AttachSubscription`, `DetachSubscription`,
    `DetachSubscriptionAndClear` and `DXFeedSubscription.Attach` and `Detach` attach it to feeds and detach it.
  * `DXFeed.CreateTimeSeriesSubscription` returns a `DXFeedTimeSeriesSubscription` whose `SetFromTime` sets the time
    from which it receives the events.
