// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package toolargs parses the command lines of the tools in cmd/tools. It does not use cgo, so its tests run without
// the native library.
//
// The rules follow the tools of the dxFeed Graal C++ API: an argument that starts with "-" is an option, the other
// arguments are positional and keep their order, and the options may stand anywhere. An option with a value takes the
// next argument; a flag takes none. Unlike the C++ tools, an unknown option is an error, so that a mistyped option
// does not go unnoticed.
package toolargs

import (
	"fmt"
	"strconv"
	"strings"
)

// Option is an option of a tool: its names, the first of which is the main one, and whether it takes a value.
type Option struct {
	Names    []string
	HasValue bool
}

// The options of the tools; each tool accepts some of them.
var (
	Quiet           = Option{Names: []string{"--quiet", "-q", "--quite"}}
	Tape            = Option{Names: []string{"--tape", "-t"}, HasValue: true}
	FromTime        = Option{Names: []string{"--from-time", "-f"}, HasValue: true}
	Source          = Option{Names: []string{"--source", "-s"}, HasValue: true}
	Properties      = Option{Names: []string{"--properties", "-p"}, HasValue: true}
	ForceStream     = Option{Names: []string{"--force-stream"}}
	IgnoreExchanges = Option{Names: []string{"--ignore-exchanges"}, HasValue: true}
	CPUUsageByCore  = Option{Names: []string{"--cpu-usage-by-core"}}
	DetachListener  = Option{Names: []string{"--detach-listener"}}
	Interval        = Option{Names: []string{"--interval"}, HasValue: true}
)

// IsHelp reports whether an argument asks for the help: -h, --help or -help.
func IsHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "-help"
}

// Args are the parsed arguments of a tool: the positional arguments and the options that are given.
type Args struct {
	Positional []string
	values     map[string]string // by the main name of the option; a flag has the value ""
}

// Parse parses the arguments that follow the name of a tool. An argument that starts with "-" (except "-" itself) is
// an option, which must be one of the accepted ones; an option with a value takes the next argument, and the other
// arguments are positional. A missing value or an unknown option is an error. When an option is repeated, the last
// value counts.
func Parse(args []string, accepted ...Option) (Args, error) {
	byName := map[string]Option{}
	for _, o := range accepted {
		for _, name := range o.Names {
			byName[name] = o
		}
	}
	parsed := Args{values: map[string]string{}}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			parsed.Positional = append(parsed.Positional, arg)
			continue
		}
		o, ok := byName[arg]
		if !ok {
			return Args{}, fmt.Errorf("unknown option %s", arg)
		}
		value := ""
		if o.HasValue {
			if i+1 >= len(args) {
				return Args{}, fmt.Errorf("option %s needs a value", arg)
			}
			i++
			value = args[i]
		}
		parsed.values[o.Names[0]] = value
	}
	return parsed, nil
}

// Has reports whether the option is given.
func (a Args) Has(o Option) bool {
	_, ok := a.values[o.Names[0]]
	return ok
}

// Value returns the value of the option, or nil if it is not given.
func (a Args) Value(o Option) *string {
	if value, ok := a.values[o.Names[0]]; ok {
		return &value
	}
	return nil
}

// Properties returns the properties of the --properties option, a comma-separated list of name=value pairs; the
// value may contain "=". It returns an empty map when the option is not given.
func (a Args) Properties() (map[string]string, error) {
	properties := map[string]string{}
	list := a.Value(Properties)
	if list == nil {
		return properties, nil
	}
	for _, pair := range strings.Split(*list, ",") {
		name, value, ok := strings.Cut(pair, "=")
		if !ok || name == "" {
			return nil, fmt.Errorf("option %s: %q is not name=value", Properties.Names[0], pair)
		}
		properties[name] = value
	}
	return properties, nil
}

// Common returns the address, the event types and the symbols: the positional arguments of the tools that connect
// to an address. The address is required; a missing types or symbols argument gets its default, or is an error when
// the default is "". More than three positional arguments are an error.
func (a Args) Common(defaultTypes, defaultSymbols string) (address, types, symbols string, err error) {
	p := a.Positional
	if len(p) > 3 {
		return "", "", "", fmt.Errorf("unexpected argument %s", p[3])
	}
	values := []string{"", defaultTypes, defaultSymbols}
	copy(values, p)
	for i, name := range []string{"address", "types", "symbols"} {
		if values[i] == "" {
			return "", "", "", fmt.Errorf("the %s argument is missing", name)
		}
	}
	return values[0], values[1], values[2], nil
}

// MaxSeconds is the greatest value of PositiveSeconds: a day.
const MaxSeconds = 86400

// PositiveSeconds parses the value of an option that is a whole number of seconds from 1 to MaxSeconds.
func PositiveSeconds(o Option, value string) (int, error) {
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds < 1 || seconds > MaxSeconds {
		return 0, fmt.Errorf("option %s: %q is not a whole number of seconds from 1 to %d", o.Names[0], value,
			MaxSeconds)
	}
	return seconds, nil
}
