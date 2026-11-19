# Configuration

The API runs the Java code of the dxFeed Graal Native SDK, so it is configured as the dxFeed Java API is: with the
address of an endpoint, the properties of the endpoints and the Java system properties of the SDK.

## Addresses

`DXEndpoint.Connect` takes the address of the data source:

- `demo.dxfeed.com:7300`: a QD endpoint, host and port.
- `dxlink:wss://demo.dxfeed.com/market-data/dxlink-ws`: a dxLink endpoint. dxLink support is experimental: before the
  first endpoint is created, set the system property `dxfeed.experimental.dxlink.enable` to `true` and `scheme` to
  `ext:opt:sysprops,resource:dxlink.xml`, as the [print-quote-events](../cmd/samples/print-quote-events/main.go)
  sample does.

An authorization token goes into the address in brackets: `host:port[login=entitle:<token>]` for a QD endpoint and
`dxlink:wss://...[login=dxlink:<token>]` for dxLink.

## Endpoint properties

`api.NewEndpoint` creates an endpoint with the default properties. The properties of the endpoints, such as
`dxfeed.address` or `dxfeed.aggregationPeriod`, are set with `api.NewEndpointWithProperties` or with the builder:

```text
endpoint, err := api.NewEndpointBuilder().
	WithRole(api.Feed).
	WithPropertiesFile("dxfeed.properties").
	WithProperty("dxfeed.aggregationPeriod", "1s").
	Build()
```

`WithPropertiesFile` reads a file in the format of Java properties files. The settings apply in the order of the calls,
so a later one wins over an earlier one with the same key: here the period of the example wins over the file.
`Build` connects the endpoint at once if `dxfeed.address` (for a feed) or `dxpublisher.address` (for a publisher) is
set. `SupportsProperty` tells whether the builder knows a property.

## System properties

The Java system properties configure the SDK as a whole: the logging, the dxLink support, the defaults of the
endpoints. `api.SetSystemProperty` sets one and `api.GetSystemProperty` reads it. Some properties are read only once,
when the SDK first needs them: set those at the start of the program, before the endpoints are created.

## Isolate arguments

The SDK runs in a Graal isolate that the first call of the API creates. `api.SetIsolateArguments` sets its arguments
before that, as a native executable takes its command line:

- `-Xmx512m` limits the heap of the SDK;
- `-Dname=value` sets a system property before any Java code runs, for the properties that are read too early for
  `api.SetSystemProperty`.

Call it at the start of `main`. After the isolate is created it returns `api.ErrIsolateCreated`. An argument the SDK
does not accept makes the creation of the isolate fail: see [troubleshooting.md](troubleshooting.md#isolate-arguments).

## Logging

By default the SDK writes its log to the standard output. Package
[logging](https://pkg.go.dev/github.com/dxfeed/dxfeed-graal-go-api/pkg/logging) sets the levels and the log files
(`logging.Init`, `logging.InitFile`, `logging.InitFiles`), changes the levels at run time (`logging.SetLevel`,
`logging.SetErrLevel`) and gives the records to a Go function (`logging.SetListener`), for example to pass them to
`log/slog`.
