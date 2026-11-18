* **\[MDAPI-550]\[GO]\[BREAKING]** The public API is cleaned up before v1.
  * **\[BREAKING]** The package `pkg/api/Osub` is renamed to `pkg/api/osub`.
  * **\[BREAKING]** The helper packages `pkg/mathutil`, `pkg/timeutil`, `pkg/formatutil` and `pkg/parser` are moved to
    `internal/` and are no longer part of the public API.
  * **\[BREAKING]** `api.CreateEndpoint` is removed; use `api.NewEndpoint`.
  * **\[BREAKING]** `DXEndpoint.AwaitProcessed`, `AwaitNotConnected` and `CloseAndAwaitTermination` take a
    `context.Context`. If it ends first, the call returns `ctx.Err()` and the wait goes on in the background;
    `CloseAndAwaitTermination` closes the endpoint even with a context that is already done.
  * `DXFeedSubscription.IsClosed` reports whether `Close` succeeded; before, it always returned `true`.
