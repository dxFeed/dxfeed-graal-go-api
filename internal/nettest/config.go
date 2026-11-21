// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package nettest holds the parts of the NetTest tool that do not call the SDK: the command line, the event pools,
// the generated symbols, the schedule of the rate limit and the output lines of the NetTest specification nettest/1
// (docs/wrapper-guide/nettest.md of the dxFeed Graal Native SDK sources from 3.18.0).
package nettest

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SpecVersion is the version of the NetTest specification that the tool follows.
const SpecVersion = "nettest/1"

// The limits of the options.
const (
	MaxSymbols   = 100000000
	MaxInstances = 1000
	MaxBatch     = 4194304
	MaxRate      = 1e12
)

// The modes.
const (
	Producer = "p"
	Consumer = "c"
)

// The roles of the endpoints.
const (
	RolePublisher       = "PUBLISHER"
	RoleStreamPublisher = "STREAM_PUBLISHER"
	RoleFeed            = "FEED"
	RoleStreamFeed      = "STREAM_FEED"
)

// The events.
const (
	EventTrade       = "Trade"
	EventTimeAndSale = "TimeAndSale"
)

// Property is a system property of the command line, -D<key>=<value>.
type Property struct{ Key, Value string }

// Config is a parsed command line.
type Config struct {
	// Help tells that a help option came where an option may come: the tool prints the usage and exits with 0. An
	// option that takes a value takes the next argument whatever it is, as in the SDK C tool, so "-n --help" is a name.
	Help bool
	// Mode is Producer or Consumer.
	Mode    string
	Address string
	// Symbols is the number of the symbols; SymbolFile is the IPF file they come from, empty when they are generated.
	Symbols    int
	SymbolFile string
	Instances  int
	// Rate is the total events per second of all instances, 0 without a limit.
	Rate  float64
	Batch int
	Event string
	// Collector is "stream", "ticker" or empty.
	Collector string
	Role      string
	Reuse     bool
	Period    time.Duration
	// Warmup and Duration are 0 when absent.
	Warmup   time.Duration
	Duration time.Duration
	Name     string
	// Properties are the -D properties in the command-line order.
	Properties []Property
	// Summary is the path of the summary file, empty without one.
	Summary string
	// MetricsAddress is the host and the port of the metrics, empty without them.
	MetricsAddress string
}

// Conflated tells whether the consumer receives conflated events: a FEED role with Trade.
func (c *Config) Conflated() bool {
	return c.Mode == Consumer && c.Role == RoleFeed && c.Event == EventTrade
}

// ServerAddress tells whether the address listens (":port") instead of connecting.
func (c *Config) ServerAddress() bool {
	return strings.HasPrefix(c.Address, ":")
}

// Property returns the value of the last -D property with the key.
func (c *Config) Property(key string) (string, bool) {
	for i := len(c.Properties) - 1; i >= 0; i-- {
		if c.Properties[i].Key == key {
			return c.Properties[i].Value, true
		}
	}
	return "", false
}

// Settings returns the settings of the command line, named as in the config of the summary file; an absent value is
// nil.
func (c *Config) Settings() map[string]any {
	null := func(ok bool, value any) any {
		if ok {
			return value
		}
		return nil
	}
	ms := func(d time.Duration) any { return null(d > 0, d.Milliseconds()) }
	producer := c.Mode == Producer
	source := "generated"
	if c.SymbolFile != "" {
		source = "file"
	}
	side := "producer"
	if !producer {
		side = "consumer"
	}
	return map[string]any{
		"side":           side,
		"address":        c.Address,
		"symbols":        c.Symbols,
		"symbols_source": source,
		"instances":      c.Instances,
		"rate_limit":     null(producer && c.Rate > 0, c.Rate),
		"batch":          null(producer, c.Batch),
		"event":          c.Event,
		"collector":      null(c.Collector != "", c.Collector),
		"role":           c.Role,
		"conflated":      c.Conflated(),
		"reuse":          null(producer, c.Reuse),
		"period_ms":      c.Period.Milliseconds(),
		"warmup_ms":      ms(c.Warmup),
		"duration_ms":    ms(c.Duration),
		"name":           c.Name,
	}
}

// SummaryConfig returns the config of the summary file: the settings without the side and the name, the rate limit as
// an integer, the reuse that takes effect (the API copies the events, so false on p) and the nanosecond part.
func (c *Config) SummaryConfig(nanoTime bool) map[string]any {
	config := c.Settings()
	delete(config, "side")
	delete(config, "name")
	if c.Mode == Producer && c.Rate > 0 {
		config["rate_limit"] = int64(math.Round(c.Rate))
	}
	config["reuse_effective"] = nil
	if c.Mode == Producer {
		config["reuse_effective"] = false
	}
	config["nano_time"] = nanoTime
	config["isolate_arguments"] = []string{}
	return config
}

// UsageError is an invalid command line: the tool prints "<tool>: <text>" and exits with 2.
type UsageError struct{ Text string }

func (e *UsageError) Error() string { return e.Text }

func usage(format string, args ...any) error { return &UsageError{fmt.Sprintf(format, args...)} }

// ErrNoArguments is the error of an empty command line: the tool prints its usage and exits with 2.
var ErrNoArguments = &UsageError{"the mode is missing"}

// IsHelp tells whether the argument is a help option.
func IsHelp(arg string) bool {
	return arg == "--help" || arg == "-help" || arg == "-h"
}

// valueOptions are the options with a value, the long ones also as --option=value; the value is the canonical name.
var valueOptions = map[string]string{
	"-S": "-S", "-C": "-C",
	"-r": "-r", "--rate-limit": "-r",
	"-s": "-s", "--stat": "-s",
	"-n": "-n", "--name": "-n",
	"-c": "-c", "--collector": "-c",
	"--event": "--event", "--batch": "--batch", "--role": "--role",
	"--duration": "--duration", "--warmup": "--warmup",
	"--summary": "--summary", "--metrics-port": "--metrics-port",
}

// producerOptions are the options of the mode p only, by their canonical names.
var producerOptions = map[string]bool{"-r": true, "--batch": true, "--reuse": true}

var (
	timePattern = regexp.MustCompile(`^([0-9]+)(ms|s|m)$`)
	ratePattern = regexp.MustCompile(`^([0-9]+[.]?[0-9]*|[.][0-9]+)([kMG]?)$`)
)

// Parse parses the command line that follows the name of the tool. The first error found wins; a help option ends the
// parsing with Help set, so an error before it wins and the arguments after it are not checked, as in the SDK C tool.
func Parse(args []string) (Config, error) {
	c := Config{Instances: 1, Batch: 100, Event: EventTrade, Period: 10 * time.Second, Name: "NetTest"}
	if len(args) == 0 {
		return c, ErrNoArguments
	}
	if IsHelp(args[0]) {
		c.Help = true
		return c, nil
	}
	c.Mode = args[0]
	if c.Mode != Producer && c.Mode != Consumer {
		return c, usage("unknown mode '%s'; the mode is p or c", c.Mode)
	}
	symbols, role := "", ""
	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch {
		case IsHelp(arg):
			c.Help = true
			return c, nil
		case strings.HasPrefix(arg, "-D"):
			key, value, ok := strings.Cut(arg[2:], "=")
			if !ok || key == "" {
				return c, usage("-D needs <key>=<value>: '%s'", arg)
			}
			c.Properties = append(c.Properties, Property{key, value})
			continue
		case strings.HasPrefix(arg, "-X"):
			return c, usage("unsupported option: %s", arg)
		case arg == "-w" || arg == "--wildcard":
			return c, usage("unsupported QD flag: %s (the wildcard subscription)", arg)
		case arg == "--legacy":
			return c, usage("unsupported option: %s", arg)
		case arg == "--reuse":
			if c.Mode != Producer {
				return c, usage("%s is an option of the mode p", arg)
			}
			c.Reuse = true
			continue
		case !strings.HasPrefix(arg, "-") || arg == "-":
			if c.Address != "" {
				return c, usage("unexpected argument '%s' after the address '%s'", arg, c.Address)
			}
			c.Address = arg
			continue
		}
		name, value, inline := arg, "", false
		if strings.HasPrefix(arg, "--") {
			name, value, inline = strings.Cut(arg, "=")
		}
		canonical, ok := valueOptions[name]
		if !ok {
			if inline && (name == "--reuse" || name == "--legacy" || name == "--wildcard" || name == "--help") {
				return c, usage("%s takes no value: '%s'", name, arg)
			}
			return c, usage("unknown flag %s", name)
		}
		if producerOptions[canonical] && c.Mode != Producer {
			return c, usage("%s is an option of the mode p", name)
		}
		if !inline {
			if i+1 >= len(args) {
				return c, usage("missing value of %s", name)
			}
			i++
			value = args[i]
		}
		var err error
		switch canonical {
		case "-S":
			symbols = value
			c.Symbols, c.SymbolFile, err = parseSymbols(value)
		case "-C":
			c.Instances, err = parseNumber(name, value, MaxInstances)
		case "-r":
			c.Rate, err = parseRate(name, value)
		case "-s":
			c.Period, err = parseTime(name, value)
		case "-n":
			if value == "" {
				err = usage("%s needs a non-empty name", name)
			}
			c.Name = value
		case "-c":
			switch value {
			case "stream", "ticker":
				c.Collector = value
			case "history":
				err = usage("-c history is reserved")
			default:
				err = usage("%s is stream or ticker: '%s'", name, value)
			}
		case "--event":
			switch {
			case strings.EqualFold(value, EventTrade):
				c.Event = EventTrade
			case strings.EqualFold(value, EventTimeAndSale):
				c.Event = EventTimeAndSale
			default:
				err = usage("%s is Trade or TimeAndSale: '%s'", name, value)
			}
		case "--batch":
			c.Batch, err = parseNumber(name, value, MaxBatch)
		case "--role":
			role = value
			if c.Mode == Producer && value != RolePublisher && value != RoleStreamPublisher {
				err = usage("%s is PUBLISHER or STREAM_PUBLISHER: '%s'", name, value)
			} else if c.Mode == Consumer && value != RoleFeed && value != RoleStreamFeed {
				err = usage("%s is FEED or STREAM_FEED: '%s'", name, value)
			}
		case "--duration":
			c.Duration, err = parseTime(name, value)
		case "--warmup":
			c.Warmup, err = parseTime(name, value)
		case "--summary":
			c.Summary = value
			err = CheckSummaryPath(value)
		case "--metrics-port":
			c.MetricsAddress, err = ParseMetricsAddress(value)
		}
		if err != nil {
			return c, err
		}
	}
	switch {
	case c.Address == "":
		return c, usage("the address is missing")
	case symbols == "":
		return c, usage("-S is missing")
	case c.Instances > 1 && c.ServerAddress():
		return c, usage("-C %d cannot be used with the server address '%s': every instance would bind the same port",
			c.Instances, c.Address)
	}
	return c, c.selectRole(role)
}

// selectRole sets the role from -c and --role and checks that the consumer still receives the contract of -c.
func (c *Config) selectRole(role string) error {
	producer := c.Mode == Producer
	if c.Collector == "ticker" && c.Event != EventTrade {
		return usage("-c ticker needs --event Trade: TimeAndSale has no ticker contract")
	}
	switch {
	case role != "":
		c.Role = role
	case producer && c.Collector == "stream":
		c.Role = RoleStreamPublisher
	case producer:
		c.Role = RolePublisher
	case c.Collector == "ticker":
		c.Role = RoleFeed
	default:
		c.Role = RoleStreamFeed
	}
	conflict := false
	switch c.Collector {
	case "stream":
		conflict = c.Role == RoleFeed && c.Event == EventTrade
	case "ticker":
		conflict = c.Role == RoleStreamPublisher || c.Role == RoleStreamFeed
	}
	if conflict {
		return usage("--role %s conflicts with -c %s for Trade", c.Role, c.Collector)
	}
	return nil
}

func parseSymbols(value string) (int, string, error) {
	switch {
	case strings.ContainsAny(value, "/"):
		return 0, "", usage("unsupported QD flag: the symbol split '/r' or '/s' in '-S %s'", value)
	case strings.ContainsAny(value, "{"):
		return 0, "", usage("unsupported QD flag: the symbol length '{len,max}' in '-S %s'", value)
	}
	total, file, hasFile := strings.Cut(value, "@")
	if hasFile && file == "" {
		return 0, "", usage("-S has no file name after '@': '%s'", value)
	}
	n, err := strconv.Atoi(total)
	if err != nil || n < 1 || n > MaxSymbols {
		return 0, "", usage("-S needs <total> or <total>@<file>, where <total> is a number from 1 to %d: '%s'",
			MaxSymbols, value)
	}
	return n, file, nil
}

func parseNumber(name, value string, max int) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 || n > max {
		return 0, usage("%s needs a number from 1 to %d: '%s'", name, max, value)
	}
	return n, nil
}

// parseTime parses a positive integer with the unit ms, s or m.
func parseTime(name, value string) (time.Duration, error) {
	bad := usage("%s needs a positive time with ms, s or m (for example 10s): '%s'", name, value)
	m := timePattern.FindStringSubmatch(value)
	if m == nil {
		return 0, bad
	}
	n, err := strconv.ParseInt(m[1], 10, 64)
	unit := map[string]time.Duration{"ms": time.Millisecond, "s": time.Second, "m": time.Minute}[m[2]]
	if err != nil || n < 1 || n > math.MaxInt64/int64(unit) {
		return 0, bad
	}
	return time.Duration(n) * unit, nil
}

// parseRate parses */<rate>: a positive decimal number with at most one dot and an optional k, M or G.
func parseRate(name, value string) (float64, error) {
	if strings.Contains(value, ";") {
		return 0, usage("unsupported QD flag: the rate limiter bucket ';<time>' in '%s %s'", name, value)
	}
	text, ok := strings.CutPrefix(value, "*/")
	if !ok {
		return 0, usage("%s needs */<rate> (for example */500k): '%s'", name, value)
	}
	bad := usage("%s has a bad rate (a positive number with an optional k, M or G): '%s'", name, value)
	m := ratePattern.FindStringSubmatch(text)
	if m == nil {
		return 0, bad
	}
	rate, err := strconv.ParseFloat(m[1], 64)
	if err != nil || rate <= 0 {
		return 0, bad
	}
	// A rate below 1 or above 1000G is a bad rate too, as in the SDK tools.
	rate *= map[string]float64{"": 1, "k": 1e3, "M": 1e6, "G": 1e9}[m[2]]
	if rate < 1 || rate > MaxRate {
		return 0, bad
	}
	return rate, nil
}
