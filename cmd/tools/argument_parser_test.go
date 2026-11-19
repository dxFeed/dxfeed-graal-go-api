// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseArguments(t *testing.T) {
	accepted := []option{quietOption, tapeOption, propertiesOption}
	a, err := parseArguments([]string{"a", "-t", "out.txt", "Quote", "--quite", "-p", "x=1", "AAPL"}, accepted...)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a", "Quote", "AAPL"}; !reflect.DeepEqual(a.positional, want) {
		t.Errorf("positional %q, want %q", a.positional, want)
	}
	if tape := a.value(tapeOption); tape == nil || *tape != "out.txt" {
		t.Errorf("tape %v, want out.txt", tape)
	}
	if !a.has(quietOption) {
		t.Error("--quite is not taken as --quiet")
	}
	if properties, err := a.properties(); err != nil || !reflect.DeepEqual(properties, map[string]string{"x": "1"}) {
		t.Errorf("properties %v, %v; want x=1", properties, err)
	}
	if a.has(fromTimeOption) || a.value(fromTimeOption) != nil {
		t.Error("an option that is not given is reported")
	}
}

func TestParseArgumentsErrors(t *testing.T) {
	for _, test := range []struct {
		args []string
		want string
	}{
		{[]string{"a", "Quote", "AAPL", "-p"}, "option -p needs a value"},
		{[]string{"a", "Quote", "AAPL", "--tape"}, "option --tape needs a value"},
		{[]string{"a", "Quote", "AAPL", "--bogus"}, "unknown option --bogus"},
		{[]string{"a", "Quote", "AAPL", "-f", "20230901Z"}, "unknown option -f"},
	} {
		_, err := parseArguments(test.args, quietOption, tapeOption, propertiesOption)
		if err == nil || err.Error() != test.want {
			t.Errorf("%q: error %v, want %q", test.args, err, test.want)
		}
	}
}

func TestProperties(t *testing.T) {
	for _, list := range []string{"x", "=1", "x=1,y"} {
		a, err := parseArguments([]string{"-p", list}, propertiesOption)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := a.properties(); err == nil || !strings.Contains(err.Error(), "is not name=value") {
			t.Errorf("%q: error %v, want not name=value", list, err)
		}
	}
	a, err := parseArguments([]string{"--properties", "x=1,y=a=b,z="}, propertiesOption)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"x": "1", "y": "a=b", "z": ""}
	if properties, err := a.properties(); err != nil || !reflect.DeepEqual(properties, want) {
		t.Errorf("properties %v, %v; want %v", properties, err, want)
	}
}
