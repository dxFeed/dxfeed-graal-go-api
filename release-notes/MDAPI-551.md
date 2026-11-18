* **\[MDAPI-551]\[GO]** `DXEndpoint` covers more of the SDK 3.16 endpoint: `api.NewEndpointBuilder()` builds endpoints
  with a role, a name, properties, a properties file and an executor; `User`, `Password`, `Reconnect`, `Disconnect`,
  `DisconnectAndClear`, `State`, `IsClosed`, `Role` and `EventTypes` are added.
  * The requested aggregation period: `SetRequestedAggregationPeriod`, `ClearRequestedAggregationPeriod`,
    `RequestedAggregationPeriod` and `AggregationPeriodInfo`.
  * Executors of the notifications: `api.NewFixedThreadPoolExecutor` and `api.NewInPlaceExecutor`, whose
    `ProcessAllPendingTasks` delivers the events on the calling goroutine.
