// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

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
const Version = "3.17.0"

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
	"windows/amd64": {"graal-native-sdk-3.17.0-amd64-windows.zip",
		"a91bb975b0b093185addb0b694eedb6b6eb31dd638a45f65c7b84012ca9a0ef5"},
	"linux/amd64": {"graal-native-sdk-3.17.0-amd64-linux.zip",
		"d5ebd45f11175d5ff7f1d08b11e66e77c620358e0b518b4e41578fdca01940af"},
	"linux/arm64": {"graal-native-sdk-3.17.0-aarch64-linux.zip",
		"d46e01d5b3d6209f33fc6032c9f44b29253fb5659e48ba3ac387edae1aad4d8c"},
	"darwin/arm64": {"graal-native-sdk-3.17.0-aarch64-osx.zip",
		"0e01dfe6b90c514dcc3248786c53f8ab0a758bd8896e9e4c407e1a209cc4d6be"},
	"darwin/amd64": {"graal-native-sdk-3.17.0-x86_64-osx.zip",
		"27d56659453d41bfd5cd695165e1aef7aafa747f407a8a1cdd5d45259a15995e"},
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

// parseVersion parses "major.minor.patch"; a suffix after the patch number ("3.17.0-SNAPSHOT") is ignored.
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
