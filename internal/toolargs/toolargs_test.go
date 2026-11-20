// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package toolargs

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	accepted := []Option{Quiet, Tape, Properties, Source}
	a, err := Parse([]string{"a", "-t", "out.txt", "Quote", "--quite", "-p", "x=1", "AAPL", "-s", "NTV"}, accepted...)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a", "Quote", "AAPL"}; !reflect.DeepEqual(a.Positional, want) {
		t.Errorf("positional %q, want %q", a.Positional, want)
	}
	if tape := a.Value(Tape); tape == nil || *tape != "out.txt" {
		t.Errorf("tape %v, want out.txt", tape)
	}
	if source := a.Value(Source); source == nil || *source != "NTV" {
		t.Errorf("source %v, want NTV", source)
	}
	if !a.Has(Quiet) {
		t.Error("--quite is not taken as --quiet")
	}
	if properties, err := a.Properties(); err != nil || !reflect.DeepEqual(properties, map[string]string{"x": "1"}) {
		t.Errorf("properties %v, %v; want x=1", properties, err)
	}
	if a.Has(FromTime) || a.Value(FromTime) != nil {
		t.Error("an option that is not given is reported")
	}
}

func TestParseLastValueCounts(t *testing.T) {
	a, err := Parse([]string{"-t", "a.txt", "--tape", "b.txt", "-"}, Tape)
	if err != nil {
		t.Fatal(err)
	}
	if tape := a.Value(Tape); tape == nil || *tape != "b.txt" {
		t.Errorf("tape %v, want b.txt", tape)
	}
	if want := []string{"-"}; !reflect.DeepEqual(a.Positional, want) {
		t.Errorf("positional %q, want %q", a.Positional, want)
	}
}

func TestParseErrors(t *testing.T) {
	for _, test := range []struct {
		args []string
		want string
	}{
		{[]string{"a", "Quote", "AAPL", "-p"}, "option -p needs a value"},
		{[]string{"a", "Quote", "AAPL", "--tape"}, "option --tape needs a value"},
		{[]string{"a", "Quote", "AAPL", "--bogus"}, "unknown option --bogus"},
		{[]string{"a", "Quote", "AAPL", "-f", "20230901Z"}, "unknown option -f"},
		{[]string{"a", "Quote", "AAPL", "--tape=out.txt"}, "unknown option --tape=out.txt"},
	} {
		_, err := Parse(test.args, Quiet, Tape, Properties)
		if err == nil || err.Error() != test.want {
			t.Errorf("%q: error %v, want %q", test.args, err, test.want)
		}
	}
}

func TestProperties(t *testing.T) {
	for _, list := range []string{"x", "=1", "x=1,y"} {
		a, err := Parse([]string{"-p", list}, Properties)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := a.Properties(); err == nil || !strings.Contains(err.Error(), "is not name=value") {
			t.Errorf("%q: error %v, want not name=value", list, err)
		}
	}
	a, err := Parse([]string{"--properties", "x=1,y=a=b,z="}, Properties)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"x": "1", "y": "a=b", "z": ""}
	if properties, err := a.Properties(); err != nil || !reflect.DeepEqual(properties, want) {
		t.Errorf("properties %v, %v; want %v", properties, err, want)
	}
	if properties, err := (Args{}).Properties(); err != nil || len(properties) != 0 {
		t.Errorf("properties without the option %v, %v; want none", properties, err)
	}
}

func TestCommon(t *testing.T) {
	for _, test := range []struct {
		positional             []string
		defaultTypes, defaultS string
		want                   [3]string
		err                    string
	}{
		{[]string{"a", "Quote", "AAPL"}, "", "", [3]string{"a", "Quote", "AAPL"}, ""},
		{[]string{"a"}, "all", "all", [3]string{"a", "all", "all"}, ""},
		{[]string{"a", "Trade"}, "all", "all", [3]string{"a", "Trade", "all"}, ""},
		{[]string{"a", "Quote"}, "", "", [3]string{}, "the symbols argument is missing"},
		{[]string{}, "all", "all", [3]string{}, "the address argument is missing"},
		{[]string{"a", "Quote", "AAPL", "extra"}, "", "", [3]string{}, "unexpected argument extra"},
	} {
		address, types, symbols, err := Args{Positional: test.positional}.Common(test.defaultTypes, test.defaultS)
		if test.err != "" {
			if err == nil || err.Error() != test.err {
				t.Errorf("%q: error %v, want %q", test.positional, err, test.err)
			}
			continue
		}
		if err != nil || [3]string{address, types, symbols} != test.want {
			t.Errorf("%q: %q %q %q, %v; want %q", test.positional, address, types, symbols, err, test.want)
		}
	}
}

func TestPositiveSeconds(t *testing.T) {
	for value, want := range map[string]int{"1": 1, "2": 2, "3600": 3600} {
		if got, err := PositiveSeconds(Interval, value); err != nil || got != want {
			t.Errorf("%q: %d, %v; want %d", value, got, err, want)
		}
	}
	for _, value := range []string{"0", "-1", "1.5", "x", ""} {
		if _, err := PositiveSeconds(Interval, value); err == nil ||
			!strings.Contains(err.Error(), "option --interval") {
			t.Errorf("%q: error %v, want an error naming --interval", value, err)
		}
	}
}

func TestIsHelp(t *testing.T) {
	for arg, want := range map[string]bool{"-h": true, "--help": true, "-help": true, "help": false, "-H": false} {
		if IsHelp(arg) != want {
			t.Errorf("IsHelp(%q) = %v", arg, !want)
		}
	}
}
