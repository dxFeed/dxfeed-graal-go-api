* **\[MDAPI-547]\[GO]\[BREAKING]** A project builds with the API after `go get` and a documented per-user setup,
  without `sudo` and without writing into the Go module cache.
  * The headers of the dxFeed Graal Native SDK are part of the module, and the SDK version is pinned in one place.
    On the first call the API checks the version of the loaded native library and returns an error naming the needed
    version if it does not match.
  * **\[BREAKING]** The native library is installed with
    `go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk` into a per-user directory and linked through
    `CGO_LDFLAGS`; `go run .` in the module directory no longer installs it. See the README for the settings of
    each operating system.
