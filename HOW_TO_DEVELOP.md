# How to develop

The API is a cgo wrapper over the [Graal Native SDK](https://github.com/dxFeed/dxfeed-graal-native-sdk), a C library.
Building it needs Go, a C compiler and the SDK, which a script downloads into the repository.

Supported platforms: Windows x64, Linux x64 and arm64, macOS arm64 and x64.

## Prerequisites

- **Go:** the version in [go.mod](go.mod) or later.
- **Git.**
- **A C compiler** for cgo:
    - **Linux:** GCC (`sudo apt-get install build-essential` on Debian and Ubuntu).
    - **macOS:** Clang from the Xcode Command Line Tools (`xcode-select --install`).
    - **Windows:** a 64-bit GCC-compatible compiler, see [Windows](#windows). MSVC does not work with cgo.

## Get the sources and the SDK

```bash
git clone https://github.com/dxFeed/dxfeed-graal-go-api.git
cd dxfeed-graal-go-api
go run .
```

`go run .` runs [install_dependencies.go](install_dependencies.go): it downloads the SDK release for the current OS and
architecture into `internal/native/graal` (ignored by git). Run it from the repository root, and again after the SDK
version changes.

## Build and test

```bash
go build ./...
go vet ./...
go test ./...
```

Before a pull request, run the checks that CI runs ([build.yml](.github/workflows/build.yml)):

```bash
gofmt -l .
go run honnef.co/go/tools/cmd/staticcheck@2026.2.2 ./...
go run ./internal/cmd/release check
go test -race ./...
```

`gofmt -l` must print nothing. The race detector also needs cgo, so it runs with the same compiler.

The C compiler warns about deprecated SDK functions that the API still calls; `CGO_CFLAGS=-w` hides these warnings.

`make build` builds the tools into `bin/` and copies the SDK library next to them, as the release archives have it
(on Windows it needs `make` and `cp`, for example from MSYS2 or Git Bash).

## Windows

cgo needs a GCC-compatible compiler: MinGW-w64 GCC (from [MSYS2](https://www.msys2.org/)) or Clang in MinGW mode
([llvm-mingw](https://github.com/mstorsjo/llvm-mingw)). The SDK is a plain C DLL, so either works with it.

1. Install MSYS2 and, in its shell, the 64-bit GCC:

   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```

2. Add `C:\msys64\mingw64\bin` to `PATH`: GCC runs its other tools from there, so setting `CC` alone is not enough.
   Then point Go at the compiler once; `go env -w` keeps the settings for all your Go builds:

   ```powershell
   go env -w CGO_ENABLED=1 CC=C:\msys64\mingw64\bin\gcc.exe
   ```

   Set `CC` explicitly and put the MSYS2 directory first: a 32-bit GCC found earlier in `PATH` (for example, from an
   old MinGW) fails the build with `64-bit mode not compiled in`.

3. Make the SDK DLL findable at run time. Windows has no rpath: a program finds the DLL next to its executable or in
   `PATH`, and `go test` and `go run` start binaries from a temporary directory. Add the SDK directory to `PATH` in the
   shell, or in the run and test configurations of your IDE:

   ```powershell
   $env:PATH = "$PWD\internal\native\graal;$env:PATH"
   ```

   Prefer this per checkout to a global `PATH` entry: checkouts on different SDK versions would load the wrong DLL.
   Without it, tests fail at start with exit status `0xc0000135` (DLL not found).

Linux and macOS binaries need nothing like this: they carry an rpath to `internal/native/graal`.

### Line endings

The repository stores LF line endings, and [.gitattributes](.gitattributes) keeps Go files LF in the working tree as
well, whatever `core.autocrlf` says. Otherwise `gofmt -l` would list every file on Windows. A clone made before
`.gitattributes` was added keeps its CRLF files until they are checked out again, for example with
`git rm --cached -r -q . && git reset --hard` on a clean working tree.

## Update the Graal Native SDK

1. Change `version` in [install_dependencies.go](install_dependencies.go) to the new
   [SDK release](https://github.com/dxFeed/dxfeed-graal-native-sdk/releases).
2. Run `go run .`, then build and run the tests on every platform you can; CI covers the rest.
3. Check the SDK release notes for changes in the C API (`internal/native` calls it) and in the event types.
4. Add a release note fragment, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Release

See [HOW_TO_RELEASE.md](HOW_TO_RELEASE.md).
