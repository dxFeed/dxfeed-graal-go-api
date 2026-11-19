// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

//go:build darwin && origin

package native

/*
#cgo LDFLAGS: -Wl,-rpath,@executable_path/
*/
import "C"
