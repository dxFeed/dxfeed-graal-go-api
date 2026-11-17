# Release notes

Each version lists one item per Jira ticket. Changes merged since the last release wait in
[release-notes/](release-notes/README.md) as one file per ticket and move here when a version is released.

Versions before v0.4.0 were tagged without the `v` prefix, so Go tools see them only as pseudo-versions.

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
