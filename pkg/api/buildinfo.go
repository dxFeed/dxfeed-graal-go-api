// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// ErrBuildInfoUnsupported is wrapped by the error of [GetBuildInfo] when the loaded library does not give its build
// information.
var ErrBuildInfoUnsupported = native.ErrBuildInfoUnsupported

// BuildInfo tells how the dxFeed Graal Native SDK library was built and what runs it (see [GetBuildInfo]).
//
// The values are text for display and records, such as the NetTest summary: their text may change with the runtime
// or the build, so do not parse them. A value the library does not know is empty.
type BuildInfo struct {
	// QDVersion is the version of QDS that the library is built with.
	QDVersion string
	// Vendor is the vendor and the version of the runtime of the library.
	Vendor string
	// March is the target machine of the build (-march of native-image).
	March string
	// GC is the names of the garbage collectors of the runtime, separated by ", ".
	GC string
}

// GetBuildInfo reads the build information of the library. The call creates the isolate if it is not created yet.
func GetBuildInfo() (BuildInfo, error) {
	d, err := native.GetBuildInfo()
	return BuildInfo(d), err
}
