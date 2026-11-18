package native

import (
	"reflect"
	"sort"
	"testing"
)

func TestStringListRoundTrip(t *testing.T) {
	for _, values := range [][]string{nil, {""}, {"AAPL", "", "IBM", "Привет"}} {
		list := newStringList(values)
		got := goStrings(list)
		freeStringList(list)
		if len(values) == 0 {
			if got != nil {
				t.Fatalf("goStrings(newStringList(%q)) = %q, want nil", values, got)
			}
			continue
		}
		if !reflect.DeepEqual(got, values) {
			t.Fatalf("goStrings(newStringList(%q)) = %q", values, got)
		}
	}
}

func TestParseSymbols(t *testing.T) {
	got, err := ParseSymbols("AAPL,IBM")
	if err != nil {
		t.Fatalf("ParseSymbols: %v", err)
	}
	// The SDK returns a set: the order of the symbols is not defined.
	symbols := append([]string(nil), got...)
	sort.Strings(symbols)
	if want := []string{"AAPL", "IBM"}; !reflect.DeepEqual(symbols, want) {
		t.Fatalf("ParseSymbols = %v, want %v in any order", got, want)
	}
}

func TestParseTime(t *testing.T) {
	got, err := ParseTime("19700101-000001Z")
	if err != nil {
		t.Fatalf("ParseTime: %v", err)
	}
	if got != 1000 {
		t.Fatalf("ParseTime = %d, want 1000", got)
	}
	if _, err := ParseTime("not a time"); err == nil {
		t.Fatal("ParseTime of an invalid value returned no error")
	}
}
