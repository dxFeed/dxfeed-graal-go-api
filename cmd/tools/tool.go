package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type Tool interface {
	Run(args []string)
	ShortDescription() string
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
