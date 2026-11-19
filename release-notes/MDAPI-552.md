* **\[MDAPI-552]\[GO]** `DXFeedSubscription` covers more of the SDK 3.16 subscription: `Symbols`,
  `DecoratedSymbols`, `SetSymbols`, `EventTypes` and `ContainsEventType` are added.
  * Delivery settings: `SetAggregationPeriod`, `ClearAggregationPeriod`, `AggregationPeriod`, `AggregationPeriodInfo`,
    `SetEventsBatchLimit`, `EventsBatchLimit`, `SetExecutor` and `TotalDroppedEvents`.
  * Change listeners: `DXFeedSubscription.AddChangeListener` and `RemoveChangeListener` take an
    `osub.ObservableSubscriptionChangeListener`; a listener that also implements `osub.ConfigurationChangeListener` is
    told about configuration changes.
  * `DXPublisher.Subscription(eventType)` returns the `ObservableSubscription` of the symbols the feeds subscribe to,
    with `IsClosed`, `EventTypes`, `ContainsEventType` and change listeners. `DXEndpoint.GetPublisher` now returns the
    same publisher on every call.
