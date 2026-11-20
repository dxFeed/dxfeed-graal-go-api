// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"os"
	"sync"
)

func main() {
	tool, status := selectTool(os.Args, os.Stdout, os.Stderr)
	if tool == nil {
		os.Exit(status)
	}
	os.Exit(tool.Run(os.Args[2:]))
}

// toolProperties are the system properties of the tools, as in the tools of the C++ API: the wildcard symbol "all",
// the experimental dxLink and its scheme. newEndpoint sets them before the -p properties, so that the SDK does not
// start before the arguments are checked.
var toolProperties = [][2]string{
	{"dxfeed.wildcard.enable", "true"},
	{"dxfeed.experimental.dxlink.enable", "true"},
	{"scheme", "ext:opt:sysprops,resource:dxlink.xml"},
}

var setToolProperties = sync.OnceValue(func() error {
	for _, property := range toolProperties {
		if err := api.SetSystemProperty(property[0], property[1]); err != nil {
			return fmt.Errorf("SetSystemProperty %s: %w", property[0], err)
		}
	}
	return nil
})
