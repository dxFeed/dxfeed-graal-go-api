# How to develop

The API is a cgo wrapper over the [Graal Native SDK](https://github.com/dxFeed/dxfeed-graal-native-sdk), a C library.
Building it needs Go, a C compiler and the SDK library, which an installer command downloads into the user cache.

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
go run ./cmd/install-sdk
```

The SDK headers are in the repository ([internal/native/include](internal/native/include)), and the SDK version with
the checksums of its archives is pinned in [internal/sdk/sdk.go](internal/sdk/sdk.go).
[cmd/install-sdk](cmd/install-sdk/main.go) downloads the library of that version for the current OS and architecture
into the user cache directory and prints the settings: `CGO_LDFLAGS` with the directory of the library and, on
Windows, the `PATH` entry. Set them in the shell or in the run and test configurations of your IDE, as
[docs/quick-start.md](docs/quick-start.md#install) describes. Run the command again after the SDK version changes;
each version has its own directory.
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

CI also checks the API against the latest release tag with gorelease; an incompatible change needs a `[BREAKING]`
release note fragment in the same pull request ([release-notes/README.md](release-notes/README.md)). gorelease
compares the last commit and refuses a working tree with any change, untracked files included, so run it after a
commit, with the tags and `origin/main` fetched (`git fetch --tags origin`):

```bash
go run ./internal/cmd/release api
```

The API calls only the current (`dxfg_v4`) SDK functions: a call of a deprecated one fails the build.

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
   `PATH`, and `go test` and `go run` start binaries from a temporary directory. Add the directory that
   `go run ./cmd/install-sdk` prints to `PATH` in the shell, or in the run and test configurations of your IDE.
   Without it, tests fail at start with exit status `0xc0000135` (DLL not found).

Linux and macOS binaries need nothing like this: the `-Wl,-rpath` flag in `CGO_LDFLAGS` points them at the library.

### Line endings

The repository stores LF line endings, and [.gitattributes](.gitattributes) keeps Go files LF in the working tree as
well, whatever `core.autocrlf` says. Otherwise `gofmt -l` would list every file on Windows. A clone made before
`.gitattributes` was added keeps its CRLF files until they are checked out again, for example with
`git rm --cached -r -q . && git reset --hard` on a clean working tree.

## Update the Graal Native SDK

1. Change `Version` and the archive names and SHA-256 checksums in [internal/sdk/sdk.go](internal/sdk/sdk.go) to those
   of the new [SDK release](https://github.com/dxFeed/dxfeed-graal-native-sdk/releases) (the release page shows the
   checksum of each archive), and replace the headers in [internal/native/include](internal/native/include) with
   those of the release.
2. Run `go run ./cmd/install-sdk`, update `CGO_LDFLAGS` and `PATH` to the new directory, then build and run the tests
   on every platform you can; CI covers the rest.
3. Check the SDK release notes for changes in the C API (`internal/native` calls it) and in the event types.
4. Add a release note fragment, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Release

See [HOW_TO_RELEASE.md](HOW_TO_RELEASE.md).
