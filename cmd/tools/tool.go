// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/parser"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
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

// runTool parses the arguments of a tool and runs it. With fewer than three positional arguments it prints the help
// to stdout and returns 0. A wrong argument prints the error and the help to stderr and returns 2, as the flag package
// does; an error of the run is printed to stderr and returns 1.
func runTool(name string, args []string, help string, accepted []option, run func(a toolArguments) error) int {
	parsed, err := parseArguments(args, accepted...)
	switch {
	case err != nil:
		err = usageError{err}
	case len(parsed.positional) < 3:
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

// parseCommon parses the address, the event types and the symbols, the positional arguments of every tool.
func parseCommon(a toolArguments) (string, []eventcodes.EventCode, []osub.Symbol, error) {
	if len(a.positional) > 3 {
		return "", nil, nil, usageError{fmt.Errorf("unexpected argument %s", a.positional[3])}
	}
	symbols, err := parser.ParseSymbols(a.positional[2])
	if err != nil {
		return "", nil, nil, usageError{fmt.Errorf("symbols: %w", err)}
	}
	return a.positional[0], parser.ParseEventTypes(a.positional[1]), symbols, nil
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
	}
	tool, ok := tools()[name]
	if !ok {
		fmt.Fprintf(stderr, "Unknown tool: %s\n", args[1])
		printUsage(stderr)
		return nil, 2
	}
	return tool, 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, `
Usage: tools <tool> [...]
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
