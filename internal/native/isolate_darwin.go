// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build darwin

package native

// The directory of the library comes from CGO_LDFLAGS (-L<dir> -Wl,-rpath,<dir>), see cmd/install-sdk;
// build with the origin tag to let a program find the library next to it.

/*
#cgo LDFLAGS: -lDxFeedGraalNativeSdk
*/
import "C"
