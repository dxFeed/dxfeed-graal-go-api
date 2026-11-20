// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/toolargs"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
)

type Qds struct{}

func (c Qds) ShortDescription() string {
	return "Runs the QDS tools of the SDK with their command line."
}

// Run runs the QDS tools. As the C++ tool, it takes -h only as the first argument: the other arguments belong to the
// QDS tools. With a command line the QDS tools end the process with their exit status, so Run returns only without
// one or with an error.
func (c Qds) Run(args []string) int {
	if len(args) > 0 && toolargs.IsHelp(args[0]) {
		fmt.Print(qdsHelp)
		return 0
	}
	qdsArgs, a, err := toolargs.ParseQds(args)
	var properties map[string]string
	if err == nil {
		properties, err = a.Properties()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "qds: %v\n\n%s", err, qdsHelp)
		return 2
	}
	if err := setToolProperties(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during qds: %v\n", err)
		return 1
	}
	for key, value := range properties {
		if err := api.SetSystemProperty(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "Error during qds: SetSystemProperty %s: %v\n", key, err)
			return 1
		}
	}
	if len(qdsArgs) == 0 {
		fmt.Print(qdsHelp)
		fmt.Println()
	}
	if err := native.ToolsMain(qdsArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error during qds: %v\n", err)
		return 1
	}
	return 0
}
