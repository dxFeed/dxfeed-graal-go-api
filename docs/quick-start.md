# Quick start

This guide installs the API and runs a first program. The [README](../README.md#requirements) lists the supported
systems; [troubleshooting.md](troubleshooting.md) helps when a step fails.

## Install

The API is a cgo package: it needs a C compiler (see [Requirements](../README.md#requirements)) and the dxFeed Graal
Native SDK library, which an installer command downloads into your user cache directory. No `sudo` is needed, and
nothing is written into the Go module cache.

1. Add the module to your project:

   ```bash
   go get github.com/dxfeed/dxfeed-graal-go-api@latest
   ```

2. Install the SDK version that this version of the API needs. The command downloads the release archive from
   [dxfeed-graal-native-sdk releases](https://github.com/dxFeed/dxfeed-graal-native-sdk/releases), checks its SHA-256
   checksum and prints the settings for your system:

   ```bash
   go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk
   ```

   Run it again after upgrading the module: each SDK version goes into its own directory. `-dir` chooses another
   directory, and `-print-dir` prints only the directory, for scripts.

3. Point cgo at the library with `CGO_LDFLAGS`, using the directory that the installer printed (`<dir>` below):

   * **Linux and macOS:**

     ```bash
     export CGO_LDFLAGS="-L<dir> -Wl,-rpath,<dir>"
     ```

     The run path lets the program find the library without further settings.

   * **Windows (PowerShell):** Windows has no run path; the program finds `DxFeedGraalNativeSdk.dll` next to its
     executable or in `PATH`:

     ```powershell
     $env:CGO_LDFLAGS = "-L<dir>"
     $env:PATH = "<dir>;" + $env:PATH
     ```

   `go env -w CGO_LDFLAGS="..."` keeps the setting for all your Go builds instead. A flag with a space in its path
   goes in single quotes as a whole, for example `'-LC:/Users/John Doe/AppData/Local/...'`; the installer prints it
   so.

4. Build and run your program as usual with `go build` or `go run`. To ship it, copy the library from `<dir>` next to
   the executable: on Linux the binaries also look there, on macOS build with `-tags origin` and
   `CGO_LDFLAGS_ALLOW=-Wl,-rpath,@executable_path/` for that.

At start-up the API checks the version of the loaded library. If it does not match, the first call and every later one
return an error naming the needed version.

> :warning: Before `0.5.0`, the SDK was installed with `go run .` in the module directory, which needed `sudo` on
> Linux. Versions `0.2.0` and earlier fetched it from a retired repository and cannot be installed; upgrade instead.

### Windows: the C compiler

cgo on Windows needs the 64-bit GCC of [MSYS2](https://www.msys2.org/) (`mingw-w64-x86_64-gcc`) with
`C:\msys64\mingw64\bin` in `PATH`; a 32-bit GCC found earlier in `PATH` fails the build. The steps are in
[HOW_TO_DEVELOP.md](../HOW_TO_DEVELOP.md#windows).

## A first program

The program of the [README](../README.md#implementation-details) connects to the dxFeed demo feed and prints the quotes
of AAPL and IBM. It is the [quick-start](../cmd/samples/quick-start/main.go) sample: in a clone of this repository, run

```bash
go run ./cmd/samples/quick-start
```

The quotes come after a few lines of the log of the SDK. In your own module, copy the program into `main.go` and run it
with `go run .`.

## Next steps

- [configuration.md](configuration.md): addresses, tokens, endpoint and system properties, logging.
- The package overviews and examples on [pkg.go.dev](https://pkg.go.dev/github.com/dxfeed/dxfeed-graal-go-api/pkg/api).
- More programs in [cmd/samples](../cmd/samples).
