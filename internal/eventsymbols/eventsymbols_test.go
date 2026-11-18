package eventsymbols

import (
	"math"
	"strings"
	"testing"
)

func checkGet(t *testing.T, symbol, key string, want string, wantOk bool) {
	t.Helper()
	got, ok := GetAttributeStringByKey(symbol, key)
	if got != want || ok != wantOk {
		t.Errorf("GetAttributeStringByKey(%q, %q) = (%q, %v), want (%q, %v)", symbol, key, got, ok, want, wantOk)
	}
}

func checkEq(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}

func checkNonAttributedSymbol(t *testing.T, s string) {
	t.Helper()
	checkEq(t, "GetBaseSymbol("+s+")", GetBaseSymbol(s), s)
	checkEq(t, "ChangeBaseSymbol("+s+")", ChangeBaseSymbol(s, "MSFT"), "MSFT")
	if HasExchangeCode(s) {
		t.Errorf("HasExchangeCode(%q) = true", s)
	}
	if GetExchangeCode(s) != 0 {
		t.Errorf("GetExchangeCode(%q) != 0", s)
	}
	checkEq(t, "ChangeExchangeCode(0)", ChangeExchangeCode(s, 0), s)
	checkEq(t, "ChangeExchangeCode(C)", ChangeExchangeCode(s, 'C'), s+"&C")

	checkGet(t, s, "key", "", false)
	checkEq(t, "Remove key", RemoveAttributeStringByKey(s, "key"), s)
	checkEq(t, "Change key", ChangeAttributeStringByKey(s, "key", "val"), s+"{key=val}")

	checkGet(t, s, "", "", false)
	checkEq(t, "Remove empty key", RemoveAttributeStringByKey(s, ""), s)
	checkEq(t, "Change empty key", ChangeAttributeStringByKey(s, "", "val"), s+"{=val}")
	checkEq(t, "Change empty key and value", ChangeAttributeStringByKey(s, "", ""), s+"{=}")
}

func TestEmptySymbol(t *testing.T) {
	checkNonAttributedSymbol(t, "")
}

func TestRegularSymbol(t *testing.T) {
	checkNonAttributedSymbol(t, "IBM")
}

func TestBrokenSymbol(t *testing.T) {
	for _, s := range []string{"{", "}", "{}", "A{", "A}", "A{}", "{B", "}B", "{}B"} {
		checkNonAttributedSymbol(t, s)
	}
}

// An embedded NUL is an ordinary character (the C++ test also checks a null pointer, which is an empty string here).
func TestEmbeddedNul(t *testing.T) {
	symbol := "A\x00B&N{key=v\x00l}"

	checkEq(t, "GetBaseSymbol", GetBaseSymbol(symbol), "A\x00B")
	if !HasExchangeCode(symbol) {
		t.Error("HasExchangeCode = false")
	}
	if GetExchangeCode(symbol) != 'N' {
		t.Error("GetExchangeCode != N")
	}
	checkEq(t, "ChangeBaseSymbol", ChangeBaseSymbol(symbol, "MSFT"), "MSFT&N{key=v\x00l}")
	checkGet(t, symbol, "key", "v\x00l", true)
	checkEq(t, "Remove", RemoveAttributeStringByKey(symbol, "key"), "A\x00B&N")
	checkEq(t, "Change", ChangeAttributeStringByKey(symbol, "key", "x"), "A\x00B&N{key=x}")
}

func TestRegional(t *testing.T) {
	s := "GE&N"

	checkEq(t, "GetBaseSymbol", GetBaseSymbol("AAPL&Q"), "AAPL")
	checkEq(t, "GetBaseSymbol", GetBaseSymbol(s), "GE")
	checkEq(t, "ChangeBaseSymbol", ChangeBaseSymbol(s, "F"), "F&N")
	if !HasExchangeCode(s) || GetExchangeCode(s) != 'N' {
		t.Error("exchange code of GE&N is not N")
	}
	checkEq(t, "ChangeExchangeCode(0)", ChangeExchangeCode(s, 0), "GE")
	checkEq(t, "ChangeExchangeCode(Q)", ChangeExchangeCode(s, 'Q'), "GE&Q")

	checkGet(t, s, "tho", "", false)
	checkEq(t, "Remove", RemoveAttributeStringByKey(s, "tho"), s)
	checkEq(t, "Change", ChangeAttributeStringByKey(s, "tho", "true"), "GE&N{tho=true}")

	checkGet(t, s, "", "", false)
	checkEq(t, "Remove empty", RemoveAttributeStringByKey(s, ""), s)
	checkEq(t, "Change empty", ChangeAttributeStringByKey(s, "", "yes"), "GE&N{=yes}")
	checkEq(t, "Change empty, empty", ChangeAttributeStringByKey(s, "", ""), "GE&N{=}")
}

func TestOneAttribute(t *testing.T) {
	s := "/ES{tho=true}"

	checkEq(t, "GetBaseSymbol", GetBaseSymbol(s), "/ES")
	checkEq(t, "ChangeBaseSymbol", ChangeBaseSymbol(s, "/NQ"), "/NQ{tho=true}")
	if HasExchangeCode(s) || GetExchangeCode(s) != 0 {
		t.Error("/ES{tho=true} has an exchange code")
	}
	checkEq(t, "ChangeExchangeCode(0)", ChangeExchangeCode(s, 0), s)
	checkEq(t, "ChangeExchangeCode(G)", ChangeExchangeCode(s, 'G'), "/ES&G{tho=true}")

	checkGet(t, s, "tho", "true", true)
	checkEq(t, "Remove tho", RemoveAttributeStringByKey(s, "tho"), "/ES")
	checkEq(t, "Change tho", ChangeAttributeStringByKey(s, "tho", "false"), "/ES{tho=false}")

	checkGet(t, s, "t", "", false)
	checkEq(t, "Remove t", RemoveAttributeStringByKey(s, "t"), s)
	checkEq(t, "Change t", ChangeAttributeStringByKey(s, "t", "MIN"), "/ES{t=MIN,tho=true}")

	checkGet(t, s, "zap", "", false)
	checkEq(t, "Remove zap", RemoveAttributeStringByKey(s, "zap"), s)
	checkEq(t, "Change zap", ChangeAttributeStringByKey(s, "zap", "15"), "/ES{tho=true,zap=15}")
	checkEq(t, "Change zap empty", ChangeAttributeStringByKey(s, "zap", ""), "/ES{tho=true,zap=}")

	checkGet(t, s, "", "", false)
	checkEq(t, "Remove empty", RemoveAttributeStringByKey(s, ""), s)
	checkEq(t, "Change empty", ChangeAttributeStringByKey(s, "", "code"), "/ES{=code,tho=true}")
	checkEq(t, "Change empty, empty", ChangeAttributeStringByKey(s, "", ""), "/ES{=,tho=true}")
}

func TestTwoAttributes(t *testing.T) {
	s := "A{c=1,e=3}"

	checkEq(t, "GetBaseSymbol", GetBaseSymbol(s), "A")
	checkEq(t, "ChangeBaseSymbol", ChangeBaseSymbol(s, "B"), "B{c=1,e=3}")
	if HasExchangeCode(s) || GetExchangeCode(s) != 0 {
		t.Error("A{c=1,e=3} has an exchange code")
	}
	checkEq(t, "ChangeExchangeCode(0)", ChangeExchangeCode(s, 0), s)
	checkEq(t, "ChangeExchangeCode(D)", ChangeExchangeCode(s, 'D'), "A&D{c=1,e=3}")

	checkGet(t, s, "b", "", false)
	checkEq(t, "Remove b", RemoveAttributeStringByKey(s, "b"), s)
	checkEq(t, "Change b", ChangeAttributeStringByKey(s, "b", "2"), "A{b=2,c=1,e=3}")
	checkEq(t, "Change b empty", ChangeAttributeStringByKey(s, "b", ""), "A{b=,c=1,e=3}")

	checkGet(t, s, "c", "1", true)
	checkEq(t, "Remove c", RemoveAttributeStringByKey(s, "c"), "A{e=3}")
	checkEq(t, "Change c", ChangeAttributeStringByKey(s, "c", "2"), "A{c=2,e=3}")
	checkEq(t, "Change c empty", ChangeAttributeStringByKey(s, "c", ""), "A{c=,e=3}")

	checkGet(t, s, "d", "", false)
	checkEq(t, "Remove d", RemoveAttributeStringByKey(s, "d"), s)
	checkEq(t, "Change d", ChangeAttributeStringByKey(s, "d", "4"), "A{c=1,d=4,e=3}")
	checkEq(t, "Change d empty", ChangeAttributeStringByKey(s, "d", ""), "A{c=1,d=,e=3}")

	checkGet(t, s, "e", "3", true)
	checkEq(t, "Remove e", RemoveAttributeStringByKey(s, "e"), "A{c=1}")
	checkEq(t, "Change e", ChangeAttributeStringByKey(s, "e", "0"), "A{c=1,e=0}")
	checkEq(t, "Change e empty", ChangeAttributeStringByKey(s, "e", ""), "A{c=1,e=}")

	checkGet(t, s, "t", "", false)
	checkEq(t, "Remove t", RemoveAttributeStringByKey(s, "t"), s)
	checkEq(t, "Change t", ChangeAttributeStringByKey(s, "t", "5"), "A{c=1,e=3,t=5}")
	checkEq(t, "Change t empty", ChangeAttributeStringByKey(s, "t", ""), "A{c=1,e=3,t=}")

	checkGet(t, s, "", "", false)
	checkEq(t, "Remove empty", RemoveAttributeStringByKey(s, ""), s)
	checkEq(t, "Change empty", ChangeAttributeStringByKey(s, "", "-1"), "A{=-1,c=1,e=3}")
	checkEq(t, "Change empty, empty", ChangeAttributeStringByKey(s, "", ""), "A{=,c=1,e=3}")
}

// The expected values are printed by the Java API 3.356 for the same arguments.
func TestHasAttributesAsInJava(t *testing.T) {
	for _, c := range []struct {
		symbol string
		want   bool
	}{
		{"", false},
		{"AAPL", false},
		{"AAPL{=5m}", true},
		{"AAPL&Q{=5m,tho=true}", true},
		{"{}", false},
		{"A{}", false},
		{"{a}", true},
		{"AAPL{", false},
		{"AAPL}", false},
	} {
		if got := HasAttributes(c.symbol); got != c.want {
			t.Errorf("HasAttributes(%q) = %v, want %v", c.symbol, got, c.want)
		}
	}
}

func TestBuildSpreadSymbolAsInJava(t *testing.T) {
	inf := math.Inf(1)
	for _, c := range []struct {
		spread map[string]float64
		want   string
	}{
		{map[string]float64{}, ""},
		{map[string]float64{"AAPL": 1}, "AAPL"},
		{map[string]float64{"AAPL": 2}, "=2*AAPL"},
		{map[string]float64{"AAPL": 1, "MSFT": -1}, "=AAPL-MSFT"},
		{map[string]float64{"MSFT": 1, "AAPL": 1, "IBM": -1, "GOOG": -1}, "=AAPL+MSFT-IBM-GOOG"},
		{map[string]float64{"A": 0.5, "B": -2.5, "C": 0}, "=0.5*A-2.5*B"},
		{map[string]float64{"A": 1e-5, "B": 3e9, "C": 1e7}, "=3.0E9*B+10000000*C+1.0E-5*A"},
		{map[string]float64{"A": inf, "B": -inf}, "=Infinity*A-Infinity*B"},
		{map[string]float64{"A": 2147483647.0, "B": 2147483648.0}, "=2.147483648E9*B+2147483647*A"},
		{map[string]float64{"A": -1, "B": -1}, "=-B-A"},
		{map[string]float64{"A": 0.0, "B": math.Copysign(0, -1)}, ""},
	} {
		if got := BuildSpreadSymbol(c.spread); got != c.want {
			t.Errorf("BuildSpreadSymbol(%v) = %q, want %q", c.spread, got, c.want)
		}
	}

	// Java cannot order NaN ratios; here they go last, by symbol.
	got := BuildSpreadSymbol(map[string]float64{"B": math.NaN(), "A": math.NaN(), "C": -1})
	if got != "=-C+NaN*A+NaN*B" {
		t.Errorf("BuildSpreadSymbol with NaN = %q", got)
	}
}

func TestGetExchangesByPatternAsInJava(t *testing.T) {
	if ExchangesPattern != "[0-9A-Za-z]" {
		t.Errorf("ExchangesPattern = %q", ExchangesPattern)
	}
	if SupportedExchanges != "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" {
		t.Errorf("SupportedExchanges = %q", SupportedExchanges)
	}
	for _, c := range []struct{ class, want string }{
		{"", ""},
		{"A-C", "ABC"},
		{"[0-9]", "0123456789"},
		{`\d`, "0123456789"},
		{"XYZ", "XYZ"},
		{"A-Za-z0-9", SupportedExchanges},
		{ExchangesPattern, SupportedExchanges},
	} {
		got, err := GetExchangesByPattern(c.class)
		if err != nil || got != c.want {
			t.Errorf("GetExchangesByPattern(%q) = (%q, %v), want %q", c.class, got, err, c.want)
		}
	}
	for _, c := range []struct{ class, message string }{
		{"[!-#]", `unsupported exchange codes: !, ", #`},
		// A dash at the end is not a range.
		{"[a-]", "unsupported exchange codes: -"},
		// A reversed range is rejected, as in Java.
		{"[z-a]", "the end of a range is less than its start"},
		{"z-a", "the end of a range is less than its start"},
		{"[^9-0]", "the end of a range is less than its start"},
	} {
		_, err := GetExchangesByPattern(c.class)
		if err == nil || !strings.Contains(err.Error(), c.message) {
			t.Errorf("GetExchangesByPattern(%q) error = %v, want a message with %q", c.class, err, c.message)
		}
	}
}
