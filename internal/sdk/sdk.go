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
const Version = "3.18.0"

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
	"windows/amd64": {"graal-native-sdk-3.18.0-amd64-windows.zip",
		"3e372e15eb2083f3b044e761b4056b313ccb63e80ef642745cdecee8c6fe605c"},
	"linux/amd64": {"graal-native-sdk-3.18.0-amd64-linux.zip",
		"4f393b33fea08cbbb775a5b6819e7582e5767ac8968c2f2151810f873f1b1d21"},
	"linux/arm64": {"graal-native-sdk-3.18.0-aarch64-linux.zip",
		"30b3bee21307150894011e59a3d520d973b974b0bce6f9c8b9accb821f3b6853"},
	"darwin/arm64": {"graal-native-sdk-3.18.0-aarch64-osx.zip",
		"77529976ed6759c7b0a757c39e4da19b5e7530bb47553774108406958481aac4"},
	"darwin/amd64": {"graal-native-sdk-3.18.0-x86_64-osx.zip",
		"0fe56d40c23d900f705108d180e3e1a75898437a40f6da62fc6185236538077e"},
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

// parseVersion parses "major.minor.patch"; a suffix after the patch number ("3.18.0-SNAPSHOT") is ignored.
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
