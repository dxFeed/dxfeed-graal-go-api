// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"strings"
)

// option is an option of a tool: its names, the first of which is the main one, and whether it takes a value.
type option struct {
	names    []string
	hasValue bool
}

// The options of the tools; each tool accepts some of them.
var (
	quietOption           = option{names: []string{"--quiet", "-q", "--quite"}}
	tapeOption            = option{names: []string{"--tape", "-t"}, hasValue: true}
	fromTimeOption        = option{names: []string{"--from-time", "-f"}, hasValue: true}
	propertiesOption      = option{names: []string{"--properties", "-p"}, hasValue: true}
	forceStreamOption     = option{names: []string{"--force-stream"}}
	ignoreExchangesOption = option{names: []string{"--ignore-exchanges"}, hasValue: true}
)

// toolArguments are the parsed arguments of a tool: the positional arguments and the options that are given.
type toolArguments struct {
	positional []string
	values     map[string]string // by the main name of the option; a flag has the value ""
}

// parseArguments parses the arguments that follow the name of a tool. An argument that starts with "-" is an option,
// which must be one of the accepted ones; an option with a value takes the next argument, and the other arguments are
// positional. A missing value or an unknown option is an error.
func parseArguments(args []string, accepted ...option) (toolArguments, error) {
	byName := map[string]option{}
	for _, o := range accepted {
		for _, name := range o.names {
			byName[name] = o
		}
	}
	parsed := toolArguments{values: map[string]string{}}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			parsed.positional = append(parsed.positional, arg)
			continue
		}
		o, ok := byName[arg]
		if !ok {
			return toolArguments{}, fmt.Errorf("unknown option %s", arg)
		}
		value := ""
		if o.hasValue {
			if i+1 >= len(args) {
				return toolArguments{}, fmt.Errorf("option %s needs a value", arg)
			}
			i++
			value = args[i]
		}
		parsed.values[o.names[0]] = value
	}
	return parsed, nil
}

// has reports whether the option is given.
func (a toolArguments) has(o option) bool {
	_, ok := a.values[o.names[0]]
	return ok
}

// value returns the value of the option, or nil if it is not given.
func (a toolArguments) value(o option) *string {
	if value, ok := a.values[o.names[0]]; ok {
		return &value
	}
	return nil
}

// properties returns the properties of the --properties option, a comma-separated list of name=value pairs.
func (a toolArguments) properties() (map[string]string, error) {
	properties := map[string]string{}
	list := a.value(propertiesOption)
	if list == nil {
		return properties, nil
	}
	for _, pair := range strings.Split(*list, ",") {
		name, value, ok := strings.Cut(pair, "=")
		if !ok || name == "" {
			return nil, fmt.Errorf("option %s: %q is not name=value", propertiesOption.names[0], pair)
		}
		properties[name] = value
	}
	return properties, nil
}
