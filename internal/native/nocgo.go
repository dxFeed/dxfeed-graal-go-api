// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build !cgo

package native

// The API calls the Graal Native SDK, a C library, through cgo. Without cgo (CGO_ENABLED=0, no C compiler found, or
// the Go Playground), the build fails here with one error whose identifier tells what to do. The other files of the
// package need cgo and are left out of such a build.
var _ = dxfeed_graal_go_api_requires_cgo_set_CGO_ENABLED_1_and_install_a_C_compiler
