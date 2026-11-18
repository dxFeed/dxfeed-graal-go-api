package candle

import (
	"math"
	"testing"
)

// sameDouble compares as Java Double.compare equality: NaN equals NaN, 0.0 and -0.0 differ.
func sameDouble(a, b float64) bool {
	return doubleEquals(a, b)
}

func typeByName(t *testing.T, name string) CandleType {
	t.Helper()
	for _, typ := range CandleTypes() {
		if typ.Name() == name {
			return typ
		}
	}
	t.Fatalf("unknown candle type name: %s", name)
	return CandleTypeTick
}

func TestParseCandlePeriodMatchesJava(t *testing.T) {
	for _, c := range periodParseCases {
		period, err := ParseCandlePeriod(c.input)
		if !c.ok {
			if err == nil {
				t.Errorf("ParseCandlePeriod(%q) did not fail", c.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseCandlePeriod(%q): %v", c.input, err)
			continue
		}
		if !sameDouble(period.Value(), c.value) {
			t.Errorf("ParseCandlePeriod(%q).Value() = %v, want %v", c.input, period.Value(), c.value)
		}
		if period.Type().Name() != c.typeName {
			t.Errorf("ParseCandlePeriod(%q).Type() = %s, want %s", c.input, period.Type().Name(), c.typeName)
		}
		if period.String() != c.str {
			t.Errorf("ParseCandlePeriod(%q).String() = %q, want %q", c.input, period.String(), c.str)
		}
		if period.PeriodIntervalMillis() != c.millis {
			t.Errorf("ParseCandlePeriod(%q).PeriodIntervalMillis() = %d, want %d", c.input, period.PeriodIntervalMillis(), c.millis)
		}
	}
}

func TestNewCandlePeriodMatchesJavaStringAndInterval(t *testing.T) {
	for _, c := range periodValueOfCases {
		period := NewCandlePeriod(c.value, typeByName(t, c.typeName))
		if period.String() != c.str {
			t.Errorf("NewCandlePeriod(%v, %s).String() = %q, want %q", c.value, c.typeName, period.String(), c.str)
		}
		if period.PeriodIntervalMillis() != c.millis {
			t.Errorf("NewCandlePeriod(%v, %s).PeriodIntervalMillis() = %d, want %d", c.value, c.typeName,
				period.PeriodIntervalMillis(), c.millis)
		}
	}
}

func TestCandlePeriodEqualityFollowsDoubleCompare(t *testing.T) {
	nan := math.NaN()

	if NewCandlePeriod(0.0, CandleTypeMinute).Equal(NewCandlePeriod(math.Copysign(0, -1), CandleTypeMinute)) {
		t.Error("0.0 equals -0.0")
	}
	if !NewCandlePeriod(nan, CandleTypeMinute).Equal(NewCandlePeriod(nan, CandleTypeMinute)) {
		t.Error("NaN does not equal NaN")
	}
	if !mustPeriod(t, "5.0m").Equal(mustPeriod(t, "5m")) {
		t.Error("5.0m does not equal 5m")
	}
	if !NewCandlePeriod(1.0, CandleTypeDay).Equal(CandlePeriodDay) {
		t.Error("1.0 day does not equal CandlePeriodDay")
	}
	if !NewCandlePeriod(1.0, CandleTypeTick).Equal(CandlePeriodTick) {
		t.Error("1.0 tick does not equal CandlePeriodTick")
	}
}

func TestParseCandleTypeMatchesJava(t *testing.T) {
	for _, c := range typeParseCases {
		typ, err := ParseCandleType(c.input)
		if c.typeName == "" {
			if err == nil {
				t.Errorf("ParseCandleType(%q) did not fail", c.input)
			}
			continue
		}
		if err != nil || typ.Name() != c.typeName {
			t.Errorf("ParseCandleType(%q) = %v, %v, want %s", c.input, typ.Name(), err, c.typeName)
		}
	}
}

func TestParseCandlePriceLevelMatchesJava(t *testing.T) {
	for _, c := range priceLevelParseCases {
		level, err := ParseCandlePriceLevel(c.input)
		if !c.ok {
			if err == nil {
				t.Errorf("ParseCandlePriceLevel(%q) did not fail", c.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseCandlePriceLevel(%q): %v", c.input, err)
			continue
		}
		if level.String() != c.str {
			t.Errorf("ParseCandlePriceLevel(%q).String() = %q, want %q", c.input, level.String(), c.str)
		}
		if level.Equal(CandlePriceLevelDefault) != c.isDefault {
			t.Errorf("ParseCandlePriceLevel(%q) is default = %v, want %v", c.input, !c.isDefault, c.isDefault)
		}
	}
}

func TestNewCandlePriceLevelMatchesJava(t *testing.T) {
	for _, c := range priceLevelValueOfCases {
		level, err := NewCandlePriceLevel(c.value)
		if c.str == "" {
			if err == nil {
				t.Errorf("NewCandlePriceLevel(%v) did not fail", c.value)
			}
			continue
		}
		if err != nil || level.String() != c.str {
			t.Errorf("NewCandlePriceLevel(%v) = %q, %v, want %q", c.value, level.String(), err, c.str)
		}
	}
}

// Checked against Java QD 3.356: Double.parseDouble trims all characters <= ' ' and takes a type suffix after an
// exponent; CandlePeriod.parse reads the value only as far as its own grammar goes.
func TestCandleAttributesParseControlCharacterTrimAndTypeSuffixesAsJava(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"\t5\n", "5"},
		{"\n1.5\r", "1.5"},
		{"\x0b2\x1f", "2"},
		{"1e2d", "100"},
		{"1E2F", "100"},
		{"0x1p3f", "8"},
		{"0x1.8p1D", "3"},
	} {
		level, err := ParseCandlePriceLevel(c.input)
		if err != nil || level.String() != c.want {
			t.Errorf("ParseCandlePriceLevel(%q) = %q, %v, want %q", c.input, level.String(), err, c.want)
		}
	}

	for _, input := range []string{"1e2dd", "1e2 d", "0x1p3ff", "1 5", "\t\n"} {
		if _, err := ParseCandlePriceLevel(input); err == nil {
			t.Errorf("ParseCandlePriceLevel(%q) did not fail", input)
		}
	}
	for _, input := range []string{"1e2dm", "0x1p3fm", "\t5m"} {
		if _, err := ParseCandlePeriod(input); err == nil {
			t.Errorf("ParseCandlePeriod(%q) did not fail", input)
		}
	}
}

func TestZeroCandlePeriodPrintsTickValue(t *testing.T) {
	var period CandlePeriod

	// The C++ default period has no type and prints "0"; the zero Go type is the tick.
	if period.String() != "0t" || period.PeriodIntervalMillis() != 0 {
		t.Errorf("zero CandlePeriod = %q, %d", period.String(), period.PeriodIntervalMillis())
	}
}

func TestCandleAttributeStrings(t *testing.T) {
	if CandlePeriodDefault.String() != "t" || CandlePriceDefault.String() != "last" ||
		CandleSessionDefault.String() != "false" || CandleAlignmentDefault.String() != "m" ||
		CandlePriceLevelDefault.String() != "NaN" || CandleDataTypeDefault.String() != "" {
		t.Error("the default attributes have unexpected strings")
	}
	for _, c := range []struct {
		parse func(string) (string, error)
		input string
		want  string
	}{
		{func(s string) (string, error) { v, e := ParseCandlePrice(s); return v.String(), e }, "b", "bid"},
		{func(s string) (string, error) { v, e := ParseCandlePrice(s); return v.String(), e }, "s", "s"},
		{func(s string) (string, error) { v, e := ParseCandlePrice(s); return v.String(), e }, "", ""},
		{func(s string) (string, error) { v, e := ParseCandlePrice(s); return v.String(), e }, "x", ""},
		{func(s string) (string, error) { v, e := ParseCandleSession(s); return v.String(), e }, "T", "true"},
		{func(s string) (string, error) { v, e := ParseCandleSession(s); return v.String(), e }, "", ""},
		{func(s string) (string, error) { v, e := ParseCandleAlignment(s); return v.String(), e }, "S", "s"},
		{func(s string) (string, error) { v, e := ParseCandleAlignment(s); return v.String(), e }, "", ""},
	} {
		got, err := c.parse(c.input)
		if (c.want == "") != (err != nil) || got != c.want && err == nil {
			t.Errorf("parse(%q) = %q, %v, want %q", c.input, got, err, c.want)
		}
	}
	if got := CandleDataTypeFromSymbol("A{dt=tho}"); got != CandleDataTypeTHO {
		t.Errorf("CandleDataTypeFromSymbol = %q", got)
	}
	if got := NormalizeCandleDataTypeAttribute("A{dt=}"); got != "A" {
		t.Errorf("NormalizeCandleDataTypeAttribute = %q", got)
	}
}
