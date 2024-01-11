package main

import (
	"fmt"
	"strings"
)

type Tool interface {
	Run(args []string)
	ShortDescription() string
}

func Run(args []string) {
	if len(args) == 1 {
		fmt.Println(`
Usage: tools <tool> [...]
Where <tool> is one of:`)
		for key, tool := range tools() {
			fmt.Printf("\t%-15s   -%s\n", key, tool.ShortDescription())
		}
	} else {
		tool := createTool(args)
		Tool.Run(tool, args)
	}
}

func createTool(args []string) Tool {
	tools := tools()
	value := strings.ToLower(args[1])
	return tools[value]
}

func tools() map[string]Tool {
	return map[string]Tool{
		"connect":  Connect{},
		"dump":     Dump{},
		"perftest": PerfTest{},
	}
}
