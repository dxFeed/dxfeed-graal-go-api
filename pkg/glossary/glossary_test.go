package glossary

import (
	"maps"
	"math"
	"slices"
	"strings"
	"testing"
)

func near(a, b float64) bool {
	return math.Abs(a-b) <= 1e-12*math.Max(1, math.Max(math.Abs(a), math.Abs(b))) || a == b
}

// ok returns the value; an error fails the test with a panic.
func ok[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func TestCFI(t *testing.T) {
	cfi := ok(ParseCFI("ESNTPB"))
	if cfi.Code() != "ESNTPB" || cfi.String() != "ESNTPB" {
		t.Errorf("code = %q", cfi.Code())
	}
	intCode := ok(cfi.IntCode())
	if byIntCode := ok(CFIOfIntCode(intCode)); byIntCode != cfi {
		t.Errorf("CFI of the int code %d = %v", intCode, byIntCode)
	}
	if cfi.Category() != 'E' || cfi.Group() != 'S' {
		t.Errorf("category, group = %c, %c", cfi.Category(), cfi.Group())
	}
	if !cfi.IsEquity() || cfi.IsDebtInstrument() || cfi.IsEntitlement() || cfi.IsOption() || cfi.IsFuture() ||
		cfi.IsOther() {
		t.Error("predicates of an equity are wrong")
	}
	if option := ok(ParseCFI("OPEFPS")); !option.IsOption() || option.IsEquity() || option.IsFuture() {
		t.Error("predicates of an option are wrong")
	}
	if short := ok(ParseCFI("ES")); short.Code() != "ESXXXX" {
		t.Errorf("short code = %q", short.Code())
	}
	if _, err := ParseCFI("!!!!!!"); err == nil {
		t.Error("an invalid code did not fail")
	}

	values := ok(cfi.Decipher())
	if len(values) != 6 || values[0].Code != 'E' || values[1].Code != 'S' {
		t.Fatalf("deciphered values = %+v", values)
	}
	for _, v := range values {
		if v.Name == "" || v.Attribute.Name == "" {
			t.Errorf("value without a name: %+v", v)
		}
	}
	if description := ok(cfi.Describe()); description == "" {
		t.Error("empty description")
	}

	var empty CFI
	if empty.Code() != "XXXXXX" || ok(ParseCFI("XXXXXX")) != empty || ok(ParseCFI("")) != empty {
		t.Error("the empty CFI is wrong")
	}
	if !strings.HasPrefix(empty.String(), "X") || empty.IsEquity() {
		t.Error("the empty CFI has a category")
	}
}

func TestPriceIncrements(t *testing.T) {
	crypto := ok(ParsePriceIncrements("1E-18"))
	if crypto.Text() != "1E-18" || ok(PriceIncrementsOf(1e-18)) != crypto {
		t.Errorf("text = %q", crypto.Text())
	}
	if p := ok(crypto.Precision()); p != 18 {
		t.Errorf("precision = %d", p)
	}
	if p := ok(crypto.PrecisionForPrice(0.12345)); p != 18 {
		t.Errorf("precision for a price = %d", p)
	}
	for _, price := range []float64{math.NaN(), 1e-10, 1e10} {
		if i := ok(crypto.IncrementForPrice(price, 0)); !near(i, 1e-18) {
			t.Errorf("increment for %v = %v", price, i)
		}
	}
	for _, c := range []struct {
		price     float64
		direction int
		want      float64
	}{
		{0.0001, 1, 1.00000000000001e-4}, {1e-9, 1, 1.000000001e-9}, {1e-18, 1, 2e-18}, {1e-18, -1, 0},
		{0, 1, 1e-18}, {0, -1, -1e-18}, {1e-15, -1, 9.99e-16}, {1, 1, 1}, {1e4, -1, 1e4},
	} {
		if got := ok(crypto.IncrementPrice(c.price, c.direction)); !near(got, c.want) {
			t.Errorf("IncrementPrice(%v, %d) = %v, want %v", c.price, c.direction, got, c.want)
		}
	}
	fine := ok(ParsePriceIncrements("0.00000001"))
	if got := ok(fine.IncrementPrice(0.60618514, 1)); !near(got, 0.60618515) {
		t.Errorf("IncrementPrice = %v", got)
	}

	for _, c := range []struct {
		value float64
		want  int
	}{{100000000.123456789, 0}, {100.12345678, 4}, {0.12345678912345678, 7}, {0.000100001, 9}} {
		if got := ok(ok(PriceIncrementsOf(c.value)).Precision()); got != c.want {
			t.Errorf("precision of %v = %d, want %d", c.value, got, c.want)
		}
	}
	ranged := ok(ParsePriceIncrements("1E-10 1E-5; 0.01"))
	if p := ok(ranged.PrecisionForPrice(1.00000000001e-5)); p != 10 {
		t.Errorf("precision below the limit = %d", p)
	}
	if p := ok(ranged.PrecisionForPrice(1.0000000000100001e-5)); p != 2 {
		t.Errorf("precision above the limit = %d", p)
	}

	tiered := ok(ParsePriceIncrements("0.0001 1; 0.01"))
	if values := ok(tiered.Values()); !slices.Equal(values, []float64{0.0001, 1, 0.01}) {
		t.Errorf("values = %v", values)
	}
	if i := ok(tiered.Increment()); i != 0.0001 {
		t.Errorf("first increment = %v", i)
	}
	for _, c := range []struct {
		price     float64
		direction int
		want      float64
	}{{0.5, 0, 0.0001}, {5, 0, 0.01}, {1, 1, 0.01}, {1, -1, 0.0001}} {
		if got := ok(tiered.IncrementForPrice(c.price, c.direction)); got != c.want {
			t.Errorf("IncrementForPrice(%v, %d) = %v, want %v", c.price, c.direction, got, c.want)
		}
	}
	if got := ok(tiered.IncrementPriceWithStep(5, 1, 0.05)); !near(got, 5.05) {
		t.Errorf("IncrementPriceWithStep = %v", got)
	}

	whole := ok(PriceIncrementsOf(1))
	if got := ok(whole.RoundPrice(2.4)); got != 2 {
		t.Errorf("RoundPrice = %v", got)
	}
	if got := ok(whole.RoundPriceWithDirection(2.4, 7)); got != 3 {
		t.Errorf("RoundPriceWithDirection up = %v", got)
	}
	if got := ok(whole.RoundPriceWithDirection(2.6, -1)); got != 2 {
		t.Errorf("RoundPriceWithDirection down = %v", got)
	}
	modes := map[RoundingMode][2]float64{
		RoundingUp: {3, -3}, RoundingDown: {2, -2}, RoundingCeiling: {3, -2}, RoundingFloor: {2, -3},
		RoundingHalfUp: {3, -3}, RoundingHalfDown: {2, -2}, RoundingHalfEven: {2, -2},
	}
	for mode, want := range modes {
		for i, price := range []float64{2.5, -2.5} {
			if got := ok(whole.RoundPriceWithMode(price, mode)); got != want[i] {
				t.Errorf("RoundPriceWithMode(%v, %d) = %v, want %v", price, mode, got, want[i])
			}
		}
	}
	if _, err := whole.RoundPriceWithMode(2.5, RoundingUnnecessary); err == nil {
		t.Error("RoundingUnnecessary of a price that needs rounding did not fail")
	}
	if got := ok(whole.RoundPriceWithMode(2, RoundingUnnecessary)); got != 2 {
		t.Errorf("RoundingUnnecessary of a valid price = %v", got)
	}
	if _, err := ParsePriceIncrements("bad"); err == nil {
		t.Error("an invalid text did not fail")
	}

	var empty PriceIncrements
	if empty.Text() != "" || ok(ParsePriceIncrements("")) != empty {
		t.Error("the empty price increments are wrong")
	}
}

func TestAdditionalUnderlyings(t *testing.T) {
	fromText := ok(ParseAdditionalUnderlyings("FIS 53; US$ 45.46"))
	if fromText.Text() != "FIS 53; US$ 45.46" || fromText.String() != fromText.Text() {
		t.Errorf("text = %q", fromText.Text())
	}
	want := map[string]float64{"FIS": 53, "US$": 45.46}
	if m := ok(fromText.Map()); !maps.Equal(m, want) {
		t.Errorf("map = %v", m)
	}
	if spc := ok(fromText.SPC("US$")); spc != 45.46 {
		t.Errorf("SPC = %v", spc)
	}
	if spc := ok(fromText.SPC("IBM")); spc != 0 {
		t.Errorf("SPC of a missing symbol = %v", spc)
	}
	if spc := ok(SPCFromText("FIS 53; US$ 45.46", "US$")); spc != 45.46 {
		t.Errorf("SPCFromText = %v", spc)
	}
	if spc := ok(SPCFromText("FIS 53; US$ 45.46", "IBM")); spc != 0 {
		t.Errorf("SPCFromText of a missing symbol = %v", spc)
	}
	if fromMap := ok(AdditionalUnderlyingsOf(map[string]float64{"US$": 45.46, "FIS": 53})); fromMap != fromText {
		t.Errorf("from the map = %q", fromMap.Text())
	}
	var empty AdditionalUnderlyings
	if got := ok(AdditionalUnderlyingsOf(nil)); got != empty || len(ok(empty.Map())) != 0 {
		t.Error("the empty additional underlyings are wrong")
	}
	if got := ok(ParseAdditionalUnderlyings("")); got != empty {
		t.Error("an empty text is not empty")
	}
	if _, err := ParseAdditionalUnderlyings("FIS x"); err == nil {
		t.Error("an invalid text did not fail")
	}
	if _, err := ParseAdditionalUnderlyings("A\x00 1"); err == nil {
		t.Error("a NUL character did not fail")
	}
}
