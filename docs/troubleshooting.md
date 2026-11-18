# Troubleshooting

## Build

- **`cannot find -lDxFeedGraalNativeSdk`**, or undefined references to `dxfg_` functions: cgo does not find the
  library of the SDK. Install it with `go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk` and set
  `CGO_LDFLAGS` to the value it prints ([quick-start.md](quick-start.md#install)). A path with a space goes in single
  quotes as a whole, as the installer prints it.
- **`64-bit mode not compiled in`** (Windows): a 32-bit GCC is found first. Use the 64-bit GCC of MSYS2 and put
  `C:\msys64\mingw64\bin` first in `PATH` ([HOW_TO_DEVELOP.md](../HOW_TO_DEVELOP.md#windows)).
- **`C compiler "gcc" not found`**, or a build that skips cgo: install the compiler of your system (see
  [Requirements](../README.md#requirements)) and keep `CGO_ENABLED=1`.

## Start-up

The program builds, but stops at start or its first call fails:

- **Exit status `0xc0000135`** (Windows), **`cannot open shared object file`** (Linux) or **`Library not loaded`**
  (macOS): the program does not find the library at run time. On Windows add the directory that `install-sdk` prints
  to `PATH`, or copy `DxFeedGraalNativeSdk.dll` next to the executable; `go test` and `go run` start binaries from a
  temporary directory, so `PATH` is the way for them. On Linux and macOS build with the `-Wl,-rpath,<dir>` flag that
  the installer prints, or copy the library next to the executable ([quick-start.md](quick-start.md#install)).
- **`dxFeed Graal Native SDK "x.y.z" is loaded, but this API needs ...`**: the library of another SDK version is
  found, for example after an upgrade of the module. Run `install-sdk` again: each SDK version goes into its own
  directory, so `CGO_LDFLAGS` and `PATH` must point at the new one.

### Isolate arguments

**`The isolate arguments could not be parsed.`**: an argument of `api.SetIsolateArguments` is not accepted, for
example `-Xmx64q` or an unknown `-XX:` option. The SDK prints the reason to the standard error output, and every call
that needs the SDK returns the error. **`the isolate is already created`** (`api.ErrIsolateCreated`): the arguments
were set after the first call of the API; set them at the start of `main`
([configuration.md](configuration.md#isolate-arguments)).

## No events

The endpoint is created and the subscription has a listener, but no events come:

- **Address**: `DXEndpoint.State` and a connection state listener (`DXEndpoint.AddListener`) show whether the endpoint
  connects; `DXEndpoint.LastConnectionError` gives the last error of its connections. The log of the SDK names the
  address it connects to.
- **dxLink**: the experimental dxLink support must be enabled before the first endpoint is created
  ([configuration.md](configuration.md#addresses)).
- **Symbols and event types**: the symbol must exist in the feed. Time series and indexed events need the symbols of
  package `osub` (`osub.TimeSeriesSubscriptionSymbol`, `osub.IndexedEventSubscriptionSymbol`), and a candle symbol
  names its period, for example `AAPL{=1d}`.
- **Entitlements**: a feed without the rights for a symbol or an event type sends nothing for it. Check the token in
  the address ([configuration.md](configuration.md#addresses)).
- **The program ends**: the events come to the listener while the program runs; a `main` that returns right after
  `AddSymbols` receives nothing. Wait for a signal or a timer.
