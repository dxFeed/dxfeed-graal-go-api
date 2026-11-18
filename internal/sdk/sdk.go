// Package sdk pins the dxFeed Graal Native SDK release that the API is built against.
//
// The package does not use cgo, so the installer command can use it before the native library is installed.
// Updating the SDK means changing Version and Archives here and replacing the headers in internal/native/include.
package sdk

import (
	"fmt"
	"strconv"
	"strings"
)

// Version is the release of the dxFeed Graal Native SDK whose headers are in internal/native/include.
const Version = "3.16.0"

// LibraryName is the name of the native library without the platform prefix and extension,
// as passed to the linker with -l.
const LibraryName = "DxFeedGraalNativeSdk"

// Archive is a release archive of the SDK for one platform.
type Archive struct {
	// Name is the file name of the archive in the SDK release.
	Name string
	// SHA256 is the checksum of the archive, in lower-case hex.
	SHA256 string
}

// Archives maps a Go platform (GOOS/GOARCH) to the SDK archive for it.
var Archives = map[string]Archive{
	"windows/amd64": {"graal-native-sdk-3.16.0-amd64-windows.zip",
		"a1866cb8176ddc2428f882be50e1b877ba94f8a9f6646f5f8439e936ea4913a5"},
	"linux/amd64": {"graal-native-sdk-3.16.0-amd64-linux.zip",
		"b7eef6512bcc059f4c3629f288c7a8abfded0c7efa880aa63060ba5431484350"},
	"linux/arm64": {"graal-native-sdk-3.16.0-aarch64-linux.zip",
		"c9d3e2841309ba8de8cb8412dae03d263ccc0a299933bb4758a32abe81c4a077"},
	"darwin/arm64": {"graal-native-sdk-3.16.0-aarch64-osx.zip",
		"516046ec33d3fdddf0cfd86968fe3b76b8d3030f3677bd04de43fee9e96bbc8d"},
	"darwin/amd64": {"graal-native-sdk-3.16.0-x86_64-osx.zip",
		"2a37f734742bd02139e9fec634b2356369c5b002233a33580c4f5c01c7529bfd"},
}

// DownloadURL returns the address of an archive of the pinned SDK release.
func DownloadURL(archive Archive) string {
	return "https://github.com/dxFeed/dxfeed-graal-native-sdk/releases/download/v" + Version + "/" + archive.Name
}

// CheckVersion reports whether a library of the given SDK release can be used with the pinned headers:
// the major and minor versions must be equal, and the patch version must not be older than the pinned one.
func CheckVersion(libraryVersion string) error {
	want, err := parseVersion(Version)
	if err != nil {
		return fmt.Errorf("pinned SDK version %q: %w", Version, err)
	}
	got, err := parseVersion(libraryVersion)
	if err != nil || got[0] != want[0] || got[1] != want[1] || got[2] < want[2] {
		return fmt.Errorf("dxFeed Graal Native SDK %q is loaded, but this API needs %d.%d.x not older than %s; "+
			"install it with: go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk",
			libraryVersion, want[0], want[1], Version)
	}
	return nil
}

// parseVersion parses "major.minor.patch"; a suffix after the patch number ("3.16.0-SNAPSHOT") is ignored.
func parseVersion(s string) ([3]int, error) {
	var v [3]int
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return v, fmt.Errorf("not a major.minor.patch version")
	}
	if i := strings.IndexAny(parts[2], "-+"); i >= 0 {
		parts[2] = parts[2][:i]
	}
	for i, p := range parts {
		if p == "" || strings.Trim(p, "0123456789") != "" {
			return v, fmt.Errorf("not a major.minor.patch version")
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return v, fmt.Errorf("not a major.minor.patch version")
		}
		v[i] = n
	}
	return v, nil
}
