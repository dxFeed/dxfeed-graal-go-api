* **\[MDAPI-538]\[GO]** The macOS tools archive is now named `tools-osx-aarch64-<version>.zip`: it is built on Apple
  silicon, while earlier releases named the same arm64 build `osx-x64`.
  * Releases add a Linux arm64 tools archive, `tools-linux-aarch64-<version>.zip`.
  * The minimum Go version in `go.mod` is now 1.21, the version the tools already required.
  * Tools: `latencytest` and `perftest` check the `Connect` error; `connect` and `dump` return a failed subscription
    as an error instead of panicking; `connect` reports `AddSymbol` errors for time series subscriptions; error
    messages no longer end with a stray `e`. Samples check the `AddListener` error.
