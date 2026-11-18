# Release notes

Each version lists one item per Jira ticket. Changes merged since the last release wait in
[release-notes/](release-notes/README.md) as one file per ticket and move here when a version is released.

Versions before v0.4.0 were tagged without the `v` prefix, so Go tools see them only as pseudo-versions.

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
