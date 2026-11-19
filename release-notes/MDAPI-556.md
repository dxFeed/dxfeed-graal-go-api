* **\[MDAPI-556]\[GO]** Package `logging` configures the logging of the SDK: `Init`, `InitFile` and `InitFiles` set the
  levels and the log files, `SetLevel` and `SetErrLevel` change the levels, and `SetListener` gives the log records of
  the SDK and of the Go API (a panic of a listener, events that the Go API skips) to a Go function.
  * `api.SetIsolateArguments` sets the arguments of the Graal isolate, such as `-Xmx512m` or `-Dname=value`, before
    the first call of the API; later it returns `api.ErrIsolateCreated`.
  * Package `auth` creates, parses and prints the authorization tokens (`Basic`, `Bearer`, custom) for the addresses
    of the endpoints; the `auth` sample connects with a token and renews it.
  * Package `dxtime` formats and parses the times as the SDK does (`DefaultFormat`, `GMTFormat`, `NewFormat` of a
    time zone, with milliseconds, the zone or full ISO 8601) and the time periods (`ParsePeriod`, `FormatPeriod`).
