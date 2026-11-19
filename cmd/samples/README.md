# Samples

Each directory is a program; run it from the root of the repository with `go run ./cmd/samples/<name>`, and `-h`
lists its flags. The Graal Native SDK must be installed, as [README.md](../../README.md#installation) describes.
Most samples connect to the dxFeed demo feed `demo.dxfeed.com:7300` and need the network; the ones marked
*offline* do not, and the tests in [samples_test.go](samples_test.go) run them on every CI platform.

## Feed and subscription

| Sample                                         | Shows                                                                     |
|------------------------------------------------|---------------------------------------------------------------------------|
| [quick-start](quick-start/main.go)             | the program of README.md: the quotes of AAPL and IBM                      |
| [print-quote-events](print-quote-events/main.go) | the quotes with the `DXFeed` singleton and `dxfeed.properties`          |
| [connect](connect/main.go)                     | events of any types and symbols, optionally from a time                   |
| [quote-and-trade](quote-and-trade/main.go)     | listeners, a promise and snapshots of the quotes and the trades           |
| [print-indexed-events](print-indexed-events/main.go) | the orders of one source and the connection state                   |
| [print-candles](print-candles/main.go)         | the candles of a candle symbol from a time                                |
| [fetch-daily-candles](fetch-daily-candles/main.go) | the daily candles of the last 20 days as a time series promise        |
| [print-last-events](print-last-events/main.go) | the last events and a time series with promises                           |
| [last-event-console](last-event-console/main.go) | the last events of the symbols typed in the console                     |
| [auth](auth/main.go)                           | an authorization token renewed periodically                               |
| [endpoint-reconnect](endpoint-reconnect/main.go) | the connectors and the connection errors of an endpoint; reconnecting   |

## Order books

| Sample                                                  | Shows                                                            |
|---------------------------------------------------------|------------------------------------------------------------------|
| [market-depth-model](market-depth-model/main.go)        | the order book of a symbol with `MarketDepthModel`               |
| [multiple-market-depth](multiple-market-depth/main.go)  | the order books of several symbols                               |
| [price-level-book](price-level-book/main.go)            | the price levels aggregated from the orders                      |
| [inc-order-snapshot](inc-order-snapshot/main.go)        | the snapshot and the updates of the orders with `IndexedTxModel` |
| [simple-order-book](simple-order-book/main.go)          | an order book built from the transactions of `IndexedTxModel`    |
| [price-level-service-client](price-level-service-client/main.go) | the price levels of the ORCS price level service (token) |

## Tape files and publishing

| Sample                                         | Shows                                                                     |
|------------------------------------------------|---------------------------------------------------------------------------|
| [write-tape-file](write-tape-file/main.go)     | *offline*: publishing quotes into a tape file                             |
| [file-parser](file-parser/main.go)             | *offline*: reading the events of a tape file                              |
| [convert-tape-file](convert-tape-file/main.go) | *offline*: converting a tape file into another                            |
| [publish-profiles](publish-profiles/main.go)   | *offline*: a publisher that answers the subscriptions with profiles       |
| [request-profile](request-profile/main.go)     | *offline*: requesting a profile from publish-profiles                     |

## Instrument profiles, schedule and history

| Sample                                           | Shows                                                                   |
|--------------------------------------------------|-------------------------------------------------------------------------|
| [feed-ipf-connect](feed-ipf-connect/main.go)     | reading instrument profiles and subscribing to their symbols            |
| [live-ipf](live-ipf/main.go)                     | the live updates of instrument profiles                                 |
| [option-chain](option-chain/main.go)             | the option chains of a symbol, their series and strikes                 |
| [schedule](schedule/main.go)                     | the trading days and sessions of an instrument                          |
| [ondemand-replay](ondemand-replay/main.go)       | replaying historical data with `OnDemandService`                        |
| [candle-data-download](candle-data-download/main.go) | downloading candles from the candle web service (credentials)       |
