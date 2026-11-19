// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package eventcodes

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestNamesMatchConstants checks that every EventCode constant of this package has its own name as the name of the
// event type, so a new constant cannot be left without a name.
func TestNamesMatchConstants(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "event_code.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	constants := 0
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			for _, name := range spec.(*ast.ValueSpec).Names {
				constants++
				code, ok := ByName(name.Name)
				if !ok {
					t.Errorf("ByName(%q) found nothing", name.Name)
					continue
				}
				if code.String() != name.Name {
					t.Errorf("%s.String() = %q", name.Name, code.String())
				}
			}
		}
	}
	if constants != len(All()) {
		t.Errorf("%d constants, All returns %d codes", constants, len(All()))
	}
	for i, code := range All() {
		if int(code) != i {
			t.Errorf("All()[%d] = %d, want the codes in the order of the values", i, code)
		}
	}
}

func TestString(t *testing.T) {
	for code, want := range map[EventCode]string{Quote: "Quote", TimeAndSale: "TimeAndSale",
		OrderImbalance: "OrderImbalance", -1: "EventCode(-1)", 27: "EventCode(27)"} {
		if got := code.String(); got != want {
			t.Errorf("EventCode(%d).String() = %q, want %q", int32(code), got, want)
		}
	}
}

func TestByNameUnknown(t *testing.T) {
	for _, name := range []string{"", "quote", "Qoute", "EventCode(0)"} {
		if code, ok := ByName(name); ok {
			t.Errorf("ByName(%q) = %v, want none", name, code)
		}
	}
}
