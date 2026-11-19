// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"os"
)

func main() {
	tool, status := selectTool(os.Args, os.Stdout, os.Stderr)
	if tool == nil {
		os.Exit(status)
	}
	// Enable experimental feature, and set scheme for dxLink.
	for _, property := range [][2]string{
		{"dxfeed.experimental.dxlink.enable", "true"},
		{"scheme", "ext:opt:sysprops,resource:dxlink.xml"},
	} {
		if err := api.SetSystemProperty(property[0], property[1]); err != nil {
			fmt.Printf("Error setting %s: %v\n", property[0], err)
			os.Exit(1)
		}
	}
	tool.Run(os.Args)
}
