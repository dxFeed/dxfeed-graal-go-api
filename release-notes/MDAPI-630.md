* **\[MDAPI-630]\[GO]** dxFeed Graal Native SDK updated to 3.19.0, the minimum version of the native library now.
  * A library of SDK 3.18 is refused at start-up: install the new one with
    `go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk` and point `CGO_LDFLAGS` and the run-time library
    path to its directory (the directory name holds the SDK version).
  * `api.GetBuildInfo()` tells how the library was built and what runs it: the QDS version, the vendor and version of
    the runtime, the target machine of the build and the garbage collectors, as text for display and records. It
    returns an error that wraps `api.ErrBuildInfoUnsupported` when the library does not give them.
  * `tools nettest` writes them to `tool.qd_version` and `tool.build` of its summary file, which were `null`, and
    passes the CLI vectors of the NetTest kit of SDK 3.19.0.
  * The SDK migrated to QDS 3.358. Closing an endpoint built with an executor of `api.NewInPlaceExecutor()`
    with `CloseAndAwaitTermination` no longer waits forever for the tasks nobody runs: the call runs them
    while it waits, so the listeners of that endpoint may be called on the closing goroutine.
