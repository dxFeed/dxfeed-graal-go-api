package main

import (
	"strings"
)

type DXArguments struct {
	args []string
}

func (a DXArguments) isQuite() bool {
	for _, arg := range a.args {
		if arg == "-q" || arg == "--quite" {
			return true
		}
	}
	return false
}

func (a DXArguments) tape() *string {
	for index, arg := range a.args {
		if arg == "-t" || arg == "--tape" {
			if index+1 >= len(a.args) {
				panic("Check value after -t parameter")
			}
			return &a.args[index+1]
		}
	}
	return nil
}

func (a DXArguments) forceStream() bool {
	for _, arg := range a.args {
		if arg == "--force-stream" {
			return true
		}
	}
	return false
}

func (a DXArguments) properties() map[string]string {
	properties := map[string]string{}
	for index, arg := range a.args {
		if arg == "-p" || arg == "--properties" {
			valueArg := a.args[index+1]
			for _, property := range strings.Split(valueArg, ",") {
				propPair := strings.Split(property, "=")
				properties[propPair[0]] = propPair[1]
			}
		}
	}
	return properties
}

func (a DXArguments) arguments() []string {
	var arguments []string

	for index, arg := range a.args {
		if index < 2 {
			continue
		}
		if !strings.HasPrefix(arg, "-") {
			arguments = append(arguments, arg)
		}
	}
	return arguments
}
