// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/parser"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/sdk"
	"github.com/dxfeed/dxfeed-graal-go-api/internal/toolargs"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

type Tool interface {
	// Run runs the tool with the arguments that follow its name and returns the exit status; a tool that runs until it
	// is stopped returns only after an error.
	Run(args []string) int
	ShortDescription() string
}

// usageError is an error of the arguments of a tool: it gives the help and the exit status 2.
type usageError struct{ err error }

func (e usageError) Error() string { return e.err.Error() }

// runTool parses the arguments of a tool and runs it. With -h, --help or -help anywhere, or with no positional
// arguments, it prints the help to stdout and returns 0. A wrong argument prints the error and the help to stderr and
// returns 2, as the flag package does; an error of the run is printed to stderr and returns 1.
func runTool(name string, args []string, help string, accepted []toolargs.Option,
	run func(a toolargs.Args) error) int {
	for _, arg := range args {
		if toolargs.IsHelp(arg) {
			fmt.Print(help)
			return 0
		}
	}
	parsed, err := toolargs.Parse(args, accepted...)
	switch {
	case err != nil:
		err = usageError{err}
	case len(parsed.Positional) == 0:
		fmt.Print(help)
		return 0
	default:
		err = run(parsed)
	}
	var usage usageError
	switch {
	case err == nil:
		return 0
	case errors.As(err, &usage):
		fmt.Fprintf(os.Stderr, "%s: %v\n\n%s", name, err, help)
		return 2
	default:
		fmt.Fprintf(os.Stderr, "Error during %s: %v\n", name, err)
		return 1
	}
}

// commonArgs are the arguments of the tools that connect to an address and subscribe.
type commonArgs struct {
	address    string
	types      []eventcodes.EventCode
	symbols    []osub.Symbol
	properties map[string]string
}

// parseCommon parses the address, the event types, the symbols and the -p properties. A missing types or symbols
// argument gets its default; a default of "" makes the argument required.
func parseCommon(a toolargs.Args, defaultTypes, defaultSymbols string) (commonArgs, error) {
	address, typesArg, symbolsArg, err := a.Common(defaultTypes, defaultSymbols)
	if err != nil {
		return commonArgs{}, usageError{err}
	}
	// The properties go first: parsing the symbols starts the SDK, which logs.
	properties, err := a.Properties()
	if err != nil {
		return commonArgs{}, usageError{err}
	}
	types, err := parser.ParseEventTypes(typesArg)
	if err != nil {
		return commonArgs{}, usageError{err}
	}
	symbols, err := parser.ParseSymbols(symbolsArg)
	if err != nil {
		return commonArgs{}, usageError{fmt.Errorf("symbols: %w", err)}
	}
	return commonArgs{address: address, types: types, symbols: symbols, properties: properties}, nil
}

// newEndpoint applies the -p properties as the tools of every API do: as system properties and as the properties of
// the endpoint it creates.
func newEndpoint(role common.Role, properties map[string]string) (*api.DXEndpoint, error) {
	for key, value := range properties {
		if err := api.SetSystemProperty(key, value); err != nil {
			return nil, fmt.Errorf("SetSystemProperty %s: %w", key, err)
		}
	}
	endpoint, err := api.NewEndpointWithProperties(role, properties)
	if err != nil {
		return nil, fmt.Errorf("NewEndpoint: %w", err)
	}
	return endpoint, nil
}

// selectTool returns the tool named by args[1]. Without a tool name, or with --help, -h or help, it prints the usage
// to stdout and returns no tool and the exit status 0; for an unknown name it prints an error and the usage to stderr
// and returns the exit status 2, as the flag package does for a wrong flag. It does not call the SDK, so the usage does
// not start its isolate and does not fail with it.
func selectTool(args []string, stdout, stderr io.Writer) (Tool, int) {
	if len(args) < 2 {
		printUsage(stdout)
		return nil, 0
	}
	name := strings.ToLower(args[1])
	switch name {
	case "--help", "-help", "-h", "help":
		printUsage(stdout)
		return nil, 0
	case "--version", "-version":
		fmt.Fprintf(stdout, "dxFeed Graal Go API tools %s, dxFeed Graal Native SDK %s\n", moduleVersion(), sdk.Version)
		return nil, 0
	}
	tool, ok := tools()[name]
	if !ok {
		fmt.Fprintf(stderr, "Unknown tool: %s\n", args[1])
		printUsage(stderr)
		return nil, 2
	}
	return tool, 0
}

// moduleVersion returns the version of the module the tools are built from, or "(devel)" for a build of a checkout.
func moduleVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "(devel)"
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, `
Usage: tools <tool> [...]
       tools --version
Where <tool> is one of:`)
	all := tools()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(w, "\t%-15s   -%s\n", name, all[name].ShortDescription())
	}
}

func tools() map[string]Tool {
	return map[string]Tool{
		"connect":     Connect{},
		"dump":        Dump{},
		"perftest":    PerfTest{},
		"latencytest": LatencyTest{},
	}
}
