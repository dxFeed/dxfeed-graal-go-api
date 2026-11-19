// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package common

import "fmt"

// JavaError is a Java exception thrown inside the dxFeed Graal Native SDK. Calls that fail with one return it; use
// errors.As to read it:
//
//	var javaErr common.JavaError
//	if errors.As(err, &javaErr) {
//		log.Print(javaErr.ClassName)
//	}
type JavaError struct {
	// ClassName is the fully qualified class name of the exception, for example java.lang.IllegalArgumentException.
	ClassName string
	// Message is the message of the exception; it may be empty.
	Message string
	// StackTrace is the printed stack trace of the exception.
	StackTrace string
}

// Error returns the class name and the message; the stack trace is left out.
func (e JavaError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("java: %s", e.ClassName)
	}
	return fmt.Sprintf("java: %s: %s", e.ClassName, e.Message)
}
