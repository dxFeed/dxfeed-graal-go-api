* **\[MDAPI-422]\[GO]** Added the synchronous feed queries of Graal Native SDK 3.16: last events and events if
  subscribed.
  * `DXFeed.GetLastEvent`, `GetLastEvents`, `GetLastEventIfSubscribed`, `GetIndexedEventsIfSubscribed` and
    `GetTimeSeriesIfSubscribed` read the local cache of the feed.
