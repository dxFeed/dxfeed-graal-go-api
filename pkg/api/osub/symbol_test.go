// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package osub

import (
	"reflect"
	"testing"
)

type ticker string

func TestStrings(t *testing.T) {
	if got := Strings[string](); len(got) != 0 {
		t.Fatalf("Strings() = %v, want empty", got)
	}
	got := Strings(ticker("AAPL"), ticker("IBM"))
	if want := []Symbol{String("AAPL"), String("IBM")}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Strings = %#v, want %#v", got, want)
	}
}
