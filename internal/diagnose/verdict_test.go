// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package diagnose

import (
	"reflect"
	"strings"
	"testing"
)

func TestComputeVerdict(t *testing.T) {
	types, symbols := []string{"Quote", "Trade"}, []string{"AAPL", "IBM"}
	all := map[Pair]bool{{"Quote", "AAPL"}: true, {"Quote", "IBM"}: true, {"Trade", "AAPL"}: true, {"Trade", "IBM"}: true}

	if v := ComputeVerdict(Input{Types: types, Symbols: symbols, Received: all}); v.Kind != Data || v.Hint != "" ||
		v.Kind.ExitCode() != 0 {
		t.Errorf("all pairs: %+v", v)
	}
	v := ComputeVerdict(Input{Types: types, Symbols: symbols, Received: map[Pair]bool{{"Quote", "AAPL"}: true}})
	if want := []string{"Quote IBM", "Trade AAPL", "Trade IBM"}; v.Kind != PartialData || v.Hint != HintPartial ||
		!reflect.DeepEqual(v.Missing, want) || v.Kind.ExitCode() != 2 {
		t.Errorf("some pairs: %+v", v)
	}
	v = ComputeVerdict(Input{Types: []string{"Quote"}, Symbols: []string{Wildcard},
		Received: map[Pair]bool{{"Quote", "MSFT"}: true}})
	if v.Kind != Data {
		t.Errorf("wildcard: %+v", v)
	}

	for _, test := range []struct {
		name string
		in   Input
		kind Kind
		hint string
		exit int
	}{
		{"connected, no data", Input{WasConnected: true}, ConnectedNoData, HintNoData, 3},
		{"connected, error", Input{WasConnected: true, HasLastError: true, LastError: "x"}, ConnectedNoData,
			HintNoData + " Last connection error: x", 3},
		{"DNS error", Input{HasLastError: true, LastErrorIsDNS: true, LastError: "y"}, CannotResolve,
			HintResolve + " Last error: y", 5},
		{"no host resolved", Input{NoHostResolved: true}, CannotResolve, HintResolve, 5},
		{"connect error, no host resolved", Input{HasLastError: true, LastError: "z", NoHostResolved: true},
			CannotConnect, "Last connection error: z", 4},
		{"nothing", Input{}, CannotConnect, HintNoError, 4},
	} {
		test.in.Types, test.in.Symbols = types, symbols
		v := ComputeVerdict(test.in)
		if v.Kind != test.kind || v.Hint != test.hint || v.Missing != nil || v.Kind.ExitCode() != test.exit {
			t.Errorf("%s: %+v, exit %d", test.name, v, v.Kind.ExitCode())
		}
	}
	if !strings.Contains(HintNoData, "IPF") || strings.Contains(HintNoData, "Last connection error") {
		t.Error("the hint of no data")
	}
	for kind, name := range map[Kind]string{Data: "DATA", PartialData: "PARTIAL DATA",
		ConnectedNoData: "CONNECTED, NO DATA", CannotConnect: "CANNOT CONNECT", CannotResolve: "CANNOT RESOLVE"} {
		if kind.String() != name {
			t.Errorf("%d: %q, want %q", kind, kind.String(), name)
		}
	}
}
