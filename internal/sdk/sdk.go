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
const Version = "3.19.0"

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
	"windows/amd64": {"graal-native-sdk-3.19.0-amd64-windows.zip",
		"6b42c8663956c2e39f334f9efaa0a3a1cd4d979c0dd9b03b0874ccfa63a1995f"},
	"linux/amd64": {"graal-native-sdk-3.19.0-amd64-linux.zip",
		"03a62e732defc22e0724ba1e776cdc80e2d25650b7f3602965ac10f4730e4b25"},
	"linux/arm64": {"graal-native-sdk-3.19.0-aarch64-linux.zip",
		"5fb2c19c15ad70f731b6875fadb9505e81104fc1cffdf33761a9f1bc5f870ea5"},
	"darwin/arm64": {"graal-native-sdk-3.19.0-aarch64-osx.zip",
		"ae64ac996faec2612e202be33cdb5711ad2b96082c6ba8759675ef33ea796bdd"},
	"darwin/amd64": {"graal-native-sdk-3.19.0-x86_64-osx.zip",
		"f4f68ef8cbe3333c5987cb8fdcdb48a1b1c41dc2f557838e539b3df8ebb54f4d"},
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

// parseVersion parses "major.minor.patch"; a suffix after the patch number ("3.19.0-SNAPSHOT") is ignored.
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
