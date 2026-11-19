* **\[MDAPI-560]\[GO]** Added `OnDemandService` and `HistoryEndpoint`.
  * `DXEndpoint.OnDemandService` returns the on-demand service of a `Feed` or `OnDemandFeed` endpoint: `Replay`,
    `ReplayWithSpeed`, `Pause`, `SetSpeed`, `StopAndResume`, `StopAndClear`, the state getters and property change
    listeners; the new `ondemand-replay` sample replays historical quotes.
  * `candlewebservice.HistoryEndpoint` downloads candles and other time series events from the candle web service:
    `NewHistoryEndpointBuilder` with the address, the user name and the password or a token, the compression and the
    format, then `GetTimeSeries`; the new `candle-data-download` sample prints the downloaded candles.
