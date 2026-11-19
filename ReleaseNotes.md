# Release notes

Each version lists one item per Jira ticket. Changes merged since the last release wait in
[release-notes/](release-notes/README.md) as one file per ticket and move here when a version is released.

Versions before v0.4.0 were tagged without the `v` prefix, so Go tools see them only as pseudo-versions.

## v0.7.0

* **\[MDAPI-562]\[GO]\[BREAKING]** dxFeed Graal Native SDK updated to 3.17.0, the minimum version of the native
  library now.
  * **\[BREAKING]** A library of SDK 3.16 is refused at start-up: install the new one with
    `go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk` and point `CGO_LDFLAGS` and the run-time library
    path to its directory (the directory name holds the SDK version).
  * `api.GetRuntimeStatistics()` returns the runtime statistics of the Graal isolate and of its process as
    `api.RuntimeStatistics`: the heap (used, committed, maximum, and the eden, survivor and old pools), the counts and
    the times of the young and full collections, the threads, the CPU time of the process, the uptime, the allocated
    bytes and the available processors. A value that the platform or the garbage collector does not give is `-1`.
  * `DXEndpoint.LastConnectionError()` and `MessageConnectorInfo.LastError` no longer lose a TLS failure of a
    connection (for example, a `tls+` address of a plain port); before, it was reported only by the next connect
    attempt (fixed in the SDK).
* **\[MDAPI-563]\[GO]** Go samples for the remaining samples of the C++ API.
  * Tape file samples: `write-tape-file` publishes quotes into a tape file, `file-parser` prints the events of a
    tape file, `convert-tape-file` converts one tape file into another. `publish-profiles` serves Profile events of
    `:TEST` symbols, and `request-profile` requests one of them. The tests run these samples, without network, on
    every CI platform.
  * API samples: `connect` prints events of any types and symbols, optionally from a time; `quote-and-trade` shows
    listeners, a promise and snapshots of quotes and trades; `fetch-daily-candles` fetches the daily candles of the
    last 20 days; `last-event-console` keeps the last events of many symbols in memory and prints those of the symbols
    typed in the console.
  * Order book samples: `multiple-market-depth` keeps the order books of several symbols; `price-level-book`
    aggregates the orders of a source into price levels. `cmd/samples/README.md` lists all the samples by topic.
* **\[MDAPI-573]\[GO]** CI checks the exported API against the latest release with gorelease: an incompatible change
  needs a release note marked `[BREAKING]`, and a release needs a version high enough for its changes.
* **\[MDAPI-580]\[GO]** Fixed the end of a program by `SIGPIPE` (exit status 141) on Linux when the peer of a
  connection closed it while the SDK wrote to it, for example a publisher whose client disconnected at once.
  * The API registers `SIGPIPE` and `SIGXFSZ` with `os/signal` on Linux and macOS, so a write of the program to a
    closed standard output returns `EPIPE` instead of ending the program (docs/configuration.md, "Signals").

## v0.6.0

* **\[MDAPI-422]\[GO]** Added the feed queries of Graal Native SDK 3.16: last events, events if subscribed,
  promises, attaching subscriptions and the time series subscription.
  * `DXFeed.GetLastEvent`, `GetLastEvents`, `GetLastEventIfSubscribed`, `GetIndexedEventsIfSubscribed` and
    `GetTimeSeriesIfSubscribed` read the local cache of the feed.
  * `DXFeed.GetLastEventPromise`, `GetLastEventsPromises`, `GetIndexedEventsPromise` and `GetTimeSeriesPromise`
    return a `Promise` that `Await(ctx)` waits for; a done context or closing the endpoint cancels the request. The
    `print-last-events` sample shows them.
  * `NewDXFeedSubscription` creates a detached subscription; `DXFeed.AttachSubscription`, `DetachSubscription`,
    `DetachSubscriptionAndClear` and `DXFeedSubscription.Attach` and `Detach` attach it to feeds and detach it.
  * `DXFeed.CreateTimeSeriesSubscription` returns a `DXFeedTimeSeriesSubscription` whose `SetFromTime` sets the time
    from which it receives the events.
* **\[MDAPI-551]\[GO]** `DXEndpoint` covers more of the SDK 3.16 endpoint: `api.NewEndpointBuilder()` builds endpoints
  with a role, a name, properties, a properties file and an executor; `User`, `Password`, `Reconnect`, `Disconnect`,
  `DisconnectAndClear`, `State`, `IsClosed`, `Role` and `EventTypes` are added.
  * The requested aggregation period: `SetRequestedAggregationPeriod`, `ClearRequestedAggregationPeriod`,
    `RequestedAggregationPeriod` and `AggregationPeriodInfo`.
  * Executors of the notifications: `api.NewFixedThreadPoolExecutor` and `api.NewInPlaceExecutor`, whose
    `ProcessAllPendingTasks` delivers the events on the calling goroutine.
  * The message connectors and the connection errors: `Connectors` returns a snapshot of each connector,
    `LastConnectionError` the newest failed connect or dropped connection as an `*api.ConnectionError`, and
    `AddMessageConnectorListener` reports the state changes of the connectors. The
    `endpoint-reconnect` sample shows them.
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
* **\[MDAPI-553]\[GO]** Added the remaining event types of Graal Native SDK 3.16 and the event flags API.
  * `events.EventFlag` (`TxPending`, `RemoveEvent`, `SnapshotBegin`, `SnapshotEnd`, `SnapshotSnip` and others), the
    `events.IndexedEvent` interface and the `IsSnapshotBegin`, `IsSnapshotEnd`, `IsSnapshotSnip`,
    `IsSnapshotEndOrSnip`, `IsPending` and `IsRemove` helpers.
  * New events: `series.Series`, `theoprice.TheoPrice`, `candle.DailyCandle`, and `misc.Message`,
    `misc.TextMessage`, `misc.Configuration`, `misc.TextConfiguration`.
  * New market events: `optionsale.OptionSale`, `marketmaker.MarketMaker`, `order.OtcMarketsOrder` (with
    `order.OtcMarketsPriceType`) and `orderimbalance.OrderImbalance` (with `orderimbalance.ImbalanceSide` and
    `orderimbalance.AuctionType`).
  * New Nuam Exchange events: `order.NuamOrder` (with `order.NuamOrderType` and `order.NuamTimeInForceType`),
    `timeandsale.NuamTimeAndSale` and `trade.NuamTrade`.
  * `TimeAndSale.String` prints the buyer and the seller instead of the addresses of their strings.
* **\[MDAPI-554]\[GO]** Added candle symbols with attributes.
  * `candle.CandleSymbol` parses and formats the symbols of the Java and C++ APIs: `candle.ParseCandleSymbol` (an error
    for an invalid attribute), `candle.NewCandleSymbol` and `candle.NewCandleSymbolWithAttributes`; the string form is
    normalized. `Symbol` and `SetSymbol` are deprecated.
  * The attributes: `candle.CandlePeriod` and `candle.CandleType`, `candle.CandlePrice`, `candle.CandleSession`,
    `candle.CandleAlignment`, `candle.CandleExchange`, `candle.CandlePriceLevel` and `candle.CandleDataType`.
  * `osub.Candle` makes a `candle.CandleSymbol` a subscription symbol, also inside a `TimeSeriesSubscriptionSymbol`;
    the new `print-candles` sample prints the candles of a candle symbol.
* **\[MDAPI-555]\[GO]** Instrument profiles: token, all fields by name, writer, collector, live connection, glossary and
  option chains.
  * `ipf.InstrumentProfileReader.ReadFromFileWithToken` reads profiles with an `auth.Token`.
  * `ipf.InstrumentProfileField`, `ipf.GetField`, `SetField`, `GetNumericField`, `SetNumericField`, `GetDateField` and
    `SetDateField` read and write the standard and the custom fields of a profile by name; `ipf.FormatNumber`,
    `ParseNumber`, `FormatDate`, `ParseDate` and the `ipf.Type*` instrument types.
  * `ipf.WriteToFile` writes profiles to a file in the IPF format, plain or compressed.
  * `live.InstrumentProfileCollector` collects profiles and notifies `InstrumentProfileUpdateListener`s;
    `live.InstrumentProfileConnection` keeps it up to date from a file, a URL or a server (`CreateConnection`, `Start`,
    `WaitUntilCompleted`, `State`, state change listeners, update period); the new `live-ipf` sample prints the updates.
  * The new `glossary` package: `PriceIncrements` (increments by price, precision, rounding with a direction or a
    `RoundingMode`, price increments), `CFI` (category, group, predicates, integer code, decipher, describe) and
    `AdditionalUnderlyings` (map, SPC).
  * The new `ipf/option` package groups option profiles into `OptionChain`s of `OptionSeries` by product and underlying
    (`BuildChains`, `OptionChainsBuilder`; calls, puts, strikes, `NStrikesAround`); the new `option-chain` sample
    prints the series of a chain.
* **\[MDAPI-556]\[GO]** Package `logging` configures the logging of the SDK: `Init`, `InitFile` and `InitFiles` set the
  levels and the log files, `SetLevel` and `SetErrLevel` change the levels, and `SetListener` gives the log records of
  the SDK and of the Go API (a panic of a listener, events that the Go API skips) to a Go function.
  * `api.SetIsolateArguments` sets the arguments of the Graal isolate, such as `-Xmx512m` or `-Dname=value`, before
    the first call of the API; later it returns `api.ErrIsolateCreated`.
  * Package `auth` creates, parses and prints the authorization tokens (`Basic`, `Bearer`, custom) for the addresses
    of the endpoints; the `auth` sample connects with a token and renews it.
  * Package `dxtime` formats and parses the times as the SDK does (`DefaultFormat`, `GMTFormat`, `NewFormat` of a
    time zone, with milliseconds, the zone or full ISO 8601) and the time periods (`ParsePeriod`, `FormatPeriod`).
* **\[MDAPI-557]\[GO]** Added the `schedule` package: trading schedules, days and sessions.
  * `schedule.FromProfile`, `FromProfileAndVenue` and `FromDefinition` return a `Schedule`; `TradingVenues`,
    `DownloadDefaults` and `SetDefaults` manage the venues and the default schedule data.
  * `Schedule`, `Day` and `Session` find the day or the session at a time, the sessions of a day, and the previous
    or next day or session that a `DayFilter` or `SessionFilter` accepts, such as `TradingDays` or `RegularSessions`;
    a filter may test the days or sessions with a Go function.
  * The new `schedule` sample prints the schedule of an instrument at a time.
* **\[MDAPI-558]\[GO]** Added the event models of the new package `model`.
  * `model.NewIndexedTxModelBuilder` builds an `IndexedTxModel` of the indexed events of a symbol, such as the
    orders of an order book: it passes the complete transactions of its sources to a `TxModelListener`, with batch
    and snapshot processing; `Attach`, `Detach` and `SetSources` change its subscription.
  * `model.NewTimeSeriesTxModelBuilder` builds a `TimeSeriesTxModel` of the time series events of a symbol from a
    time, such as candles: it passes their complete transactions to a `TimeSeriesTxModelListener`; `SetFromTime`
    changes the time of its subscription. The new interface `events.TimeSeriesEvent` describes these events.
  * `model.NewMarketDepthModelBuilder` builds a `MarketDepthModel`, the order book of a symbol on top of
    `IndexedTxModel`: it passes the best buy and sell orders of its sources to a `MarketDepthModelListener`, up to
    a depth limit and at most once per aggregation period; `SetDepthLimit` and `SetAggregationPeriod` change them.
  * The new `inc-order-snapshot` sample prints the snapshot and the incremental updates of the orders of a source.
  * The new `market-depth-model` and `simple-order-book` samples print the order book of a symbol.
* **\[MDAPI-560]\[GO]** Added `OnDemandService` and `HistoryEndpoint`.
  * `DXEndpoint.OnDemandService` returns the on-demand service of a `Feed` or `OnDemandFeed` endpoint: `Replay`,
    `ReplayWithSpeed`, `Pause`, `SetSpeed`, `StopAndResume`, `StopAndClear`, the state getters and property change
    listeners; the new `ondemand-replay` sample replays historical quotes.
  * `candlewebservice.HistoryEndpoint` downloads candles and other time series events from the candle web service:
    `NewHistoryEndpointBuilder` with the address, the user name and the password or a token, the compression and the
    format, then `GetTimeSeries`; the new `candle-data-download` sample prints the downloaded candles.
* **\[MDAPI-561]\[GO]** Added the client of the ORCS price level service.
  * `orcs.NewPriceLevelService` connects to the service: `GetOrders`, `GetQuotes` and their `WithCaller` forms
    request the price levels (orders) and the quotes of a candle symbol with the granularity attribute
    (`orcs.GranularityAttributeKey`) within a time range, `GetAuthOrderSource` returns the order sources and the
    symbols available to the client.
  * `orcs.NewAuthOrderSource` creates an `AuthOrderSource` from the symbols by the order source identifier.
  * `orcs.ValidatePriceLevels` checks the price levels as `PriceLevelChecker` does; the new
    `price-level-service-client` sample requests and checks them.
* **\[MDAPI-566]\[GO]** The packages have overviews on pkg.go.dev: package `api` introduces the module (endpoint roles,
  subscriptions and symbols, event types, time series, publishing, instrument profiles, threads and errors), the other
  packages describe their concepts and link the related types.
  * Examples on pkg.go.dev: time series and indexed subscriptions, the connection state listener, candle symbols, auth
    tokens, the time format and periods of `dxtime`, and a `log/slog` listener of `logging`.
  * The guides `docs/quick-start.md`, `docs/configuration.md` and `docs/troubleshooting.md` describe the installation,
    the configuration and the usual problems; a test checks the links of the Markdown files.

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
