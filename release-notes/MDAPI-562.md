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
