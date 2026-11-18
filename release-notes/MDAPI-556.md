* **\[MDAPI-556]\[GO]** Package `logging` configures the logging of the SDK: `Init`, `InitFile` and `InitFiles` set the
  levels and the log files, `SetLevel` and `SetErrLevel` change the levels, and `SetListener` gives the log records of
  the SDK and of the Go API (a panic of a listener, events that the Go API skips) to a Go function.
  * `api.SetIsolateArguments` sets the arguments of the Graal isolate, such as `-Xmx512m` or `-Dname=value`, before
    the first call of the API; later it returns `api.ErrIsolateCreated`.
