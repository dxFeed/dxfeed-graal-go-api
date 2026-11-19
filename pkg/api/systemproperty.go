// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// SetSystemProperty sets a Java system property of the SDK, for example "dxfeed.address".
func SetSystemProperty(key string, value string) error {
	return native.SetSystemProperty(key, value)
}

// GetSystemProperty returns a Java system property of the SDK; an unset property is an empty string.
func GetSystemProperty(key string) (string, error) {
	return native.GetSystemProperty(key)
}

// ErrIsolateCreated is returned by [SetIsolateArguments] after the first call that needed the SDK: the Graal isolate
// that runs the Java code of the SDK is created then, and its arguments cannot change anymore.
var ErrIsolateCreated = native.ErrIsolateCreated

// SetIsolateArguments sets the arguments that the Graal isolate of the SDK is created with, as a native executable
// takes its command line: for example "-Xmx512m" (the maximum heap size) or "-Dname=value" (a Java system property
// that is set before the first Java code runs; [SetSystemProperty] works only later). Call it before any other
// function of the API that calls the SDK, for example at the start of main; it does not create the isolate itself.
// A repeated call replaces the arguments, a call without arguments clears them.
//
// It returns [ErrIsolateCreated] if the isolate is already created (or its creation has started), and an error if an
// argument contains a NUL character; the stored arguments do not change then. An argument that the SDK does not accept
// (for example "-Xmx64q" or an unknown "-XX:" option) makes the creation of the isolate fail: the SDK prints the
// reason to the standard error output, and every call that needs the SDK returns the error.
func SetIsolateArguments(args ...string) error {
	return native.SetIsolateArguments(args)
}
