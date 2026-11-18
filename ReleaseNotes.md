# Release notes

Each version lists one item per Jira ticket. Changes merged since the last release wait in
[release-notes/](release-notes/README.md) as one file per ticket and move here when a version is released.

Versions before v0.4.0 were tagged without the `v` prefix, so Go tools see them only as pseudo-versions.

## v0.5.0

* **\[MDAPI-537]\[GO]** Added a contributing guide, a development guide (Windows setup with MinGW-w64 included) and
  issue and pull request templates; README links the public issue tracker for questions and bug reports.
* **\[MDAPI-547]\[GO]\[BREAKING]** A project builds with the API after `go get` and a documented per-user setup,
  without `sudo` and without writing into the Go module cache.
  * The headers of the dxFeed Graal Native SDK are part of the module, and the SDK version is pinned in one place.
    On the first call the API checks the version of the loaded native library and returns an error naming the needed
    version if it does not match.
  * **\[BREAKING]** The native library is installed with
    `go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk` into a per-user directory and linked through
    `CGO_LDFLAGS`; `go run .` in the module directory no longer installs it. See the README for the settings of
    each operating system.
* **\[MDAPI-548]\[GO]** The API calls the current (`dxfg_v4`) functions of the dxFeed Graal Native SDK instead of
  the ones that SDK 3.16 deprecates and SDK 4.0.0 removes.
  * Failures of adding or removing symbols, adding a listener and publishing events are returned as errors; before,
    they were lost.
  * Publishing a `Profile` without a status reason, or another event with an unset string field, no longer crashes
    the process.
* **\[MDAPI-549]\[GO]\[BREAKING]** The API returns errors instead of panicking.
  * A failure to create the isolate, or a native library of another version, is returned by the first call and every
    later one; before, every call panicked.
  * Java exceptions are returned as `common.JavaError` with the class name, the message and the stack trace; read it
    with `errors.As`. The error text is the class name and the message; the stack trace is in `JavaError.StackTrace`.
  * A panic of an event or connection state listener is logged and no longer ends the process.
  * Events of a type that the API does not support yet are skipped, and each such type is logged once.
  * Subscription symbols are passed to the SDK in C memory: time series and indexed event symbols follow the cgo
    pointer rules, and an unsupported symbol in `AddSymbols` or `RemoveSymbols` is an error.
  * **\[BREAKING]** `DXFeedSubscription.Clear` and `DXFeedSubscription.Close` return an `error`.
  * **\[BREAKING]** `api.SetSystemProperty` returns an `error`, and `api.GetSystemProperty` returns `(string, error)`.
  * **\[BREAKING]** `parser.ParseSymbols` returns `([]any, error)`; before, a parse failure gave an empty list.
* **\[MDAPI-550]\[GO]\[BREAKING]** The public API is cleaned up before v1.
  * **\[BREAKING]** The package `pkg/api/Osub` is renamed to `pkg/api/osub`.
  * **\[BREAKING]** The helper packages `pkg/mathutil`, `pkg/timeutil`, `pkg/formatutil` and `pkg/parser` are moved to
    `internal/` and are no longer part of the public API.
  * **\[BREAKING]** `api.CreateEndpoint` is removed; use `api.NewEndpoint`.
  * **\[BREAKING]** `DXEndpoint.AwaitProcessed`, `AwaitNotConnected` and `CloseAndAwaitTermination` take a
    `context.Context`. If it ends first, the call returns `ctx.Err()` and the wait goes on in the background;
    `CloseAndAwaitTermination` closes the endpoint even with a context that is already done.
  * `DXFeedSubscription.IsClosed` reports whether `Close` succeeded; before, it always returned `true`.
  * The connection states `common.NotConnected`, `Connecting`, `Connected` and `Closed` are exported, so a state
    listener can compare them.
  * Every exported identifier has a doc comment, and runnable examples show the endpoint, the subscription, the
    publisher and `common.JavaError`.
  * **\[BREAKING]** `DXFeedSubscription.AddSymbol`, `AddSymbols`, `RemoveSymbol` and `RemoveSymbols` take `osub.Symbol`
    instead of `any`, so a symbol of a wrong type does not compile. `osub.Symbol` is implemented by `osub.String`,
    `*osub.WildcardSymbol`, `*osub.TimeSeriesSubscriptionSymbol` and `*osub.IndexedEventSubscriptionSymbol`. Pass a
    string as `osub.String("AAPL")`, or several at once as `sub.AddSymbols(osub.Strings("AAPL", "IBM")...)`.
* **\[MDAPI-567]\[GO]** Fixed `Profile.IsShortSaleRestricted`: it checks the short sale restriction, as in Java and C++.
  Before, it compared the trading status and returned true for a halted instrument.
* **\[MDAPI-568]\[GO]** Fixed `Time`, `Sequence`, `SetTime` and `SetSequence` of `Order`, `AnalyticOrder` and
  `SpreadOrder`: they now use the time-sequence field, as in Java and C++. Before, received orders reported time 0,
  and `SetTime` or `SetSequence` overwrote the index with the order source.
* **\[MDAPI-569]\[GO]** Fixed `DXFeedSubscription.RemoveListener`: the removed listener no longer receives events, and
  its native listener is released. A listener may remove itself from `Update`. `Close` releases the remaining
  listeners.
* **\[MDAPI-570]\[GO]** Fixed a data race in the connection state listeners of `DXEndpoint`: `AddListener` and
  `RemoveListener` may now be called from any goroutine, also while the connection state changes or from a listener.
* **\[MDAPI-571]\[GO]** Fixed a memory leak in `DXPublisher.Publish`: the C memory of the published events is now
  freed after the call.

## v0.4.0

* **\[MDAPI-524]\[GO]** dxFeed Graal Native SDK updated to 3.16.0.
* **\[MDAPI-534]\[GO]** Fixed event codes after `AnalyticOrder`, which were shifted by one against the native SDK:
  `SpreadOrder` subscribed to OTC Markets orders and `Series` to spread orders, and publishing a `SpreadOrder` panicked.
  * **BREAKING:** the values of `eventcodes.SpreadOrder`, `Series` and `OptionSale` change; code that stores the
    numbers must be updated.
  * Added the codes `OtcMarketsOrder`, `TextMessage`, `MarketMaker`, `TextConfiguration`, `NuamOrder`,
    `NuamTimeAndSale`, `NuamTrade` and `OrderImbalance`; their Go event types are not implemented yet.
  * `CreateSubscription` and `Publish` return an error for an event type without a Go implementation, and received
    events of such a type are skipped instead of panicking.
* **\[MDAPI-535]\[GO]** Fixed `NewEndpointWithProperties`: connection state listeners added with `AddListener` are now
  called, as for `NewEndpoint`.
  * Fixed a memory leak: the property keys and values passed to `NewEndpointWithProperties` were never freed
    ([#22](https://github.com/dxFeed/dxfeed-graal-go-api/issues/22), reported by @OvOhao).
* **\[MDAPI-536]\[GO]** Releases are tagged `vX.Y.Z`, so `go get` and the Go module proxy resolve them; `v0.4.0` is
  the first. Pushing a version tag publishes a GitHub Release with the notes of that version from `ReleaseNotes.md`.
  * Added `ReleaseNotes.md` with the history of earlier versions. Changes since the last release are kept as one file
    per Jira ticket in `release-notes/` and checked in CI.
  * `HOW_TO_RELEASE.md` describes the release steps and the versioning policy.
* **\[MDAPI-538]\[GO]** The macOS tools archive is now named `tools-osx-aarch64-<version>.zip`: it is built on Apple
  silicon, while earlier releases named the same arm64 build `osx-x64`.
  * Releases add a Linux arm64 tools archive, `tools-linux-aarch64-<version>.zip`.
  * The minimum Go version in `go.mod` is now 1.21, the version the tools already required.
  * Tools: `latencytest` and `perftest` check the `Connect` error; `connect` and `dump` return a failed subscription
    as an error instead of panicking; `connect` reports `AddSymbol` errors for time series subscriptions; error
    messages no longer end with a stray `e`. Samples check the `AddListener` error.

## 0.3.0

* **\[MDAPI-419]\[GO]** Fixed the build of applications
  ([#21](https://github.com/dxFeed/dxfeed-graal-go-api/pull/21)).
  * dxFeed Graal Native SDK updated to 3.2.13.
  * Fixed the `InstrumentProfile` structure and its mapping.
  * Updated the links to the dxLink demo endpoints.

## 0.2.1

* The native SDK is downloaded from GitHub releases.

## 0.2.0

* Added the `Underlying` event ([#20](https://github.com/dxFeed/dxfeed-graal-go-api/pull/20)).

## 0.1.0

* Added the ignore exchanges option to the tools ([#15](https://github.com/dxFeed/dxfeed-graal-go-api/pull/15)).
* Added the `Summary` and `Greeks` events ([#19](https://github.com/dxFeed/dxfeed-graal-go-api/pull/19)).
* Renamed `dxednpoint.go` to `dxendpoint.go` ([#17](https://github.com/dxFeed/dxfeed-graal-go-api/pull/17)).

## 0.0.5

* Updated the GraalVM version of the native SDK.
* Added the ignore exchanges option and dxLink support to the tools.

## 0.0.4

* Added `IndexedEventSubscriptionSymbol` ([#13](https://github.com/dxFeed/dxfeed-graal-go-api/pull/13)).
* Added the latency test tool ([#14](https://github.com/dxFeed/dxfeed-graal-go-api/pull/14)).

## 0.0.3

* Added event mappers and the `Order`, `SpreadOrder`, `AnalyticOrder`, `Candle` and `Trade` events
  ([#5](https://github.com/dxFeed/dxfeed-graal-go-api/pull/5), [#8](https://github.com/dxFeed/dxfeed-graal-go-api/pull/8),
  [#9](https://github.com/dxFeed/dxfeed-graal-go-api/pull/9), [#10](https://github.com/dxFeed/dxfeed-graal-go-api/pull/10),
  [#11](https://github.com/dxFeed/dxfeed-graal-go-api/pull/11)).
* Added the endpoint connection state listener ([#6](https://github.com/dxFeed/dxfeed-graal-go-api/pull/6)).
* Lowered the required Go version ([#7](https://github.com/dxFeed/dxfeed-graal-go-api/pull/7)).
* Fixed a generic function in `mathutil` ([#12](https://github.com/dxFeed/dxfeed-graal-go-api/pull/12)).

## 0.0.2

* Added `InstrumentProfileReader` ([#2](https://github.com/dxFeed/dxfeed-graal-go-api/pull/2)).
* Symbols are parsed by the native SDK ([#4](https://github.com/dxFeed/dxfeed-graal-go-api/pull/4)).

## 0.0.1

* First release.
