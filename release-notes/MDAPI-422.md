* **\[MDAPI-422]\[GO]** Added the feed queries of Graal Native SDK 3.16: last events, events if subscribed, promises
  and the time-series subscription.
  * `DXFeed.GetLastEvent`, `GetLastEvents`, `GetLastEventIfSubscribed`, `GetIndexedEventsIfSubscribed` and
    `GetTimeSeriesIfSubscribed` read the local cache of the feed.
