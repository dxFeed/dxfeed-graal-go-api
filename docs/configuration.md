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

## Signals

On Linux and macOS the SDK writes to its connections on threads of its own. When the peer has closed a connection,
such a write raises `SIGPIPE`, and Go ends the process on a signal of a thread that it did not create, unless the
program is notified of it. So, when it creates the isolate, the API registers `SIGPIPE` and `SIGXFSZ` (file size
limit) with [os/signal](https://pkg.go.dev/os/signal) `Notify`, and the write fails with an error that the SDK
handles, as in a JVM.

This changes the program as a whole:

- a write of the program to a closed standard output or error (for example, `program | head`) returns `EPIPE`
  instead of ending the program, so check the errors of these writes or exit on them;
- `signal.Reset` of these signals, or `signal.Reset()` without arguments, brings the ending of the process back; the
  program's own `signal.Notify` of them is fine, and so is `signal.Ignore` after the first call of the API (before
  it, the registration of the API cancels the `Ignore`).

## Logging

By default the SDK writes its log to the standard output. Package
[logging](https://pkg.go.dev/github.com/dxfeed/dxfeed-graal-go-api/pkg/logging) sets the levels and the log files
(`logging.Init`, `logging.InitFile`, `logging.InitFiles`), changes the levels at run time (`logging.SetLevel`,
`logging.SetErrLevel`) and gives the records to a Go function (`logging.SetListener`), for example to pass them to
`log/slog`.
