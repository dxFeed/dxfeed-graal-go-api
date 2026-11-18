* **\[MDAPI-560]\[GO]** Added `OnDemandService` and `HistoryEndpoint`.
  * `DXEndpoint.OnDemandService` returns the on-demand service of a `Feed` or `OnDemandFeed` endpoint: `Replay`,
    `ReplayWithSpeed`, `Pause`, `SetSpeed`, `StopAndResume`, `StopAndClear`, the state getters and property change
    listeners; the new `ondemand-replay` sample replays historical quotes.
