// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

// ToolsMain runs the QDS tools of the SDK with the arguments: the name of a tool and its arguments. With at least one
// argument the Java tools end the whole process with the exit code 0 on success and 1 on failure, so the call does
// not return then; with none it prints their usage to the standard error stream and returns nil.
func ToolsMain(args []string) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		list := newStringList(args)
		defer freeStringList(list)
		return checkStatus(thread, C.dxfg_v4_Tools_main(thread.ptr, list))
	})
}
