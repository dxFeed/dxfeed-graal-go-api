* **\[MDAPI-422]\[GO]** Added the feed queries of Graal Native SDK 3.16: last events, events if subscribed and
  promises.
  * `DXFeed.GetLastEvent`, `GetLastEvents`, `GetLastEventIfSubscribed`, `GetIndexedEventsIfSubscribed` and
    `GetTimeSeriesIfSubscribed` read the local cache of the feed.
  * `DXFeed.GetLastEventPromise`, `GetLastEventsPromises`, `GetIndexedEventsPromise` and `GetTimeSeriesPromise`
    return a `Promise` that `Await(ctx)` waits for; a done context or closing the endpoint cancels the request. The
    `print-last-events` sample shows them.
