package candle

import (
	"fmt"
	"math"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/eventsymbols"
)

func mustPeriod(t *testing.T, s string) CandlePeriod {
	t.Helper()
	p, err := ParseCandlePeriod(s)
	if err != nil {
		t.Fatalf("ParseCandlePeriod(%q): %v", s, err)
	}
	return p
}

func mustSymbolWith(t *testing.T, symbol string, attributes ...CandleSymbolAttribute) string {
	t.Helper()
	c, err := NewCandleSymbolWithAttributes(symbol, attributes...)
	if err != nil {
		t.Fatalf("NewCandleSymbolWithAttributes(%q): %v", symbol, err)
	}
	return c.String()
}

func TestCandlePeriods(t *testing.T) {
	for _, c := range []struct {
		input string
		want  CandlePeriod
	}{
		{"DAY", CandlePeriodDay},
		{"1DAY", CandlePeriodDay},
		{"1day", CandlePeriodDay},
		{"1Day", CandlePeriodDay},
		{"Day", CandlePeriodDay},
		{"Days", CandlePeriodDay},
		{"d", CandlePeriodDay},
		{"1d", CandlePeriodDay},

		{"TICK", CandlePeriodTick},
		{"tick", CandlePeriodTick},
		{"Ticks", CandlePeriodTick},
		{"1tick", CandlePeriodTick},
		{"1T", CandlePeriodTick},

		{"5MINUTE", NewCandlePeriod(5, CandleTypeMinute)},
		{"5min", NewCandlePeriod(5, CandleTypeMinute)},
		{"5Min", NewCandlePeriod(5, CandleTypeMinute)},
		{"5m", NewCandlePeriod(5, CandleTypeMinute)},
	} {
		if got := mustPeriod(t, c.input); !got.Equal(c.want) {
			t.Errorf("ParseCandlePeriod(%q) = %v, want %v", c.input, got, c.want)
		}
	}

	type key struct {
		bits uint64
		typ  CandleType
	}
	set := map[key]bool{}
	val := 0.0
	for i := 0; i < 20; i++ {
		for _, typ := range CandleTypes() {
			p := NewCandlePeriod(val, typ)

			k := key{math.Float64bits(val), typ}
			if set[k] {
				t.Errorf("the period %v is not different", p)
			}
			set[k] = true

			s1 := p.String()
			if got := mustPeriod(t, s1); !got.Equal(p) {
				t.Errorf("ParseCandlePeriod(%q) = %v, want %v", s1, got, p)
			}

			s2 := fmt.Sprintf("%f%s", val, typ.Name())
			if got := mustPeriod(t, s2); !got.Equal(p) {
				t.Errorf("ParseCandlePeriod(%q) = %v, want %v (val = %v, type = %s)", s2, got, p, val, typ.Name())
			}
		}
		val += 0.25
	}
}

func TestNormalization(t *testing.T) {
	for _, c := range []struct {
		symbol     string
		attributes []CandleSymbolAttribute
		want       string
	}{
		{"IBM", []CandleSymbolAttribute{NewCandleExchange('E')}, "IBM&E"},
		{"IBM", []CandleSymbolAttribute{CandleExchangeComposite}, "IBM"},

		{"IBM", []CandleSymbolAttribute{CandlePeriodDay}, "IBM{=d}"},
		{"IBM", []CandleSymbolAttribute{CandlePeriodTick}, "IBM"},

		{"IBM", []CandleSymbolAttribute{CandlePriceAsk}, "IBM{price=ask}"},
		{"IBM", []CandleSymbolAttribute{CandlePriceLast}, "IBM"},

		{"IBM", []CandleSymbolAttribute{CandleSessionRegular}, "IBM{tho=true}"},
		{"IBM", []CandleSymbolAttribute{CandleSessionAny}, "IBM"},

		{"EUR/USD{source=bank}", []CandleSymbolAttribute{CandlePriceBid, NewCandlePeriod(2, CandleTypeHour)},
			"EUR/USD{=2h,price=bid,source=bank}"},
		{"IBM{aa=zz,price=b}", []CandleSymbolAttribute{NewCandlePeriod(15, CandleTypeMinute)},
			"IBM{=15m,aa=zz,price=bid}"},
	} {
		if got := mustSymbolWith(t, c.symbol, c.attributes...); got != c.want {
			t.Errorf("NewCandleSymbolWithAttributes(%q, %v) = %q, want %q", c.symbol, c.attributes, got, c.want)
		}
	}
}

func TestDoubleChecks(t *testing.T) {
	minutes15 := NewCandlePeriod(15, CandleTypeMinute)

	if got := mustSymbolWith(t, "AAPL", minutes15); got != "AAPL{=15m}" {
		t.Errorf("got %q", got)
	}
	if got := mustSymbolWith(t, "AAPL", CandlePriceBid, minutes15); got != "AAPL{=15m,price=bid}" {
		t.Errorf("got %q", got)
	}
	attributes := []CandleSymbolAttribute{minutes15}
	if got := mustSymbolWith(t, "AAPL", attributes...); got != "AAPL{=15m}" {
		t.Errorf("got %q", got)
	}
}

func TestCandlePeriodStringRoundTripsSmallValues(t *testing.T) {
	period := NewCandlePeriod(1.0e-8, CandleTypePrice)

	if got := mustPeriod(t, period.String()); !got.Equal(period) {
		t.Errorf("ParseCandlePeriod(%q) = %v, want %v", period.String(), got, period)
	}
	if period.Value() != 1.0e-8 {
		t.Errorf("Value() = %v", period.Value())
	}
}

func TestCandlePriceLevelFollowsJavaValidationAndEquality(t *testing.T) {
	nan, err := NewCandlePriceLevel(math.NaN())
	if err != nil || !nan.Equal(CandlePriceLevelDefault) || !CandlePriceLevelDefault.Equal(CandlePriceLevelDefault) {
		t.Errorf("NaN is not the default level: %v, %v", nan, err)
	}
	if got := CandlePriceLevelDefault.ChangeAttributeForSymbol("AAPL{pl=1}"); got != "AAPL" {
		t.Errorf("default level: got %q", got)
	}
	one, err := NewCandlePriceLevel(1.0)
	if err != nil || one.String() != "1" {
		t.Errorf("level 1.0 = %q, %v", one, err)
	}
	zero, _ := NewCandlePriceLevel(0.0)
	tiny, _ := NewCandlePriceLevel(1.0e-20)
	if zero.Equal(one) || zero.Equal(tiny) {
		t.Error("0.0 equals 1.0 or 1.0e-20")
	}

	for _, v := range []float64{-1.0, math.Copysign(0, -1), math.Inf(1)} {
		if _, err := NewCandlePriceLevel(v); err == nil {
			t.Errorf("NewCandlePriceLevel(%v) did not fail", v)
		}
	}
	if _, err := ParseCandlePriceLevel("1garbage"); err == nil {
		t.Error(`ParseCandlePriceLevel("1garbage") did not fail`)
	}
}

func TestCandleSessionAttributeLookupUsesParseBooleanSemantics(t *testing.T) {
	for _, c := range []struct {
		symbol string
		want   CandleSession
	}{
		{"AAPL{tho=true}", CandleSessionRegular},
		{"AAPL{tho=TRUE}", CandleSessionRegular},
		{"AAPL{tho=t}", CandleSessionDefault},
		{"AAPL{tho=garbage}", CandleSessionDefault},
	} {
		if got := CandleSessionFromSymbol(c.symbol); got != c.want {
			t.Errorf("CandleSessionFromSymbol(%q) = %v, want %v", c.symbol, got, c.want)
		}
	}
}

func TestCandleSymbolReportsInvalidAttributesAsErrors(t *testing.T) {
	if _, err := ParseCandleSymbol("AAPL{=garbage}"); err == nil {
		t.Error(`ParseCandleSymbol("AAPL{=garbage}") did not fail`)
	}
}

// NewCandleSymbol never fails: the invalid attribute stays in the string and has its default value.
func TestNewCandleSymbolIsLenient(t *testing.T) {
	c := NewCandleSymbol("AAPL{=garbage}")

	if c.String() != "AAPL{=garbage}" || c.BaseSymbol() != "AAPL" || !c.Period().Equal(CandlePeriodDefault) {
		t.Errorf("NewCandleSymbol(AAPL{=garbage}) = %q, base %q, period %v", c, c.BaseSymbol(), c.Period())
	}
}

func TestZeroValuesOfTheAttributesAreTheDefaultOnes(t *testing.T) {
	var session CandleSession
	if session != CandleSessionAny || session.String() != "false" {
		t.Errorf("zero CandleSession = %v", session)
	}
	var level CandlePriceLevel
	if !level.Equal(CandlePriceLevelDefault) || level.String() != "NaN" {
		t.Errorf("zero CandlePriceLevel = %v", level)
	}
	var exchange CandleExchange
	if exchange != CandleExchangeComposite || exchange.String() != "COMPOSITE" {
		t.Errorf("zero CandleExchange = %v", exchange)
	}
	var dataType CandleDataType
	if dataType != CandleDataTypeDefault {
		t.Errorf("zero CandleDataType = %q", dataType)
	}
}

// Where the C++ API deliberately differs from Java: Java keeps the default attributes (its normalization drops the
// result of the removal), C++ removes them. The attributes themselves are the same as in Java (checked below).
var cxxStrOverrides = map[string]string{
	"A{=t}":                                 "A",
	"A{=1t}":                                "A",
	"A{tho=false}":                          "A",
	"A{tho=t}":                              "A",
	"A{tho=yes}":                            "A",
	"A{tho=}":                               "A",
	"A{price=last}":                         "A",
	"A{price=l}":                            "A",
	"A{price=la}":                           "A",
	"A{a=m}":                                "A",
	"A{a=M}":                                "A",
	"A{pl=NaN}":                             "A",
	"A{=t,pl=NaN}":                          "A",
	"A{pl=NaN,tho=false,a=m,price=last,=t}": "A",
	"A{tho=false,tho=true}":                 "A",
	"A{=t,=5m}":                             "A",
	"A&Q{=t}":                               "A&Q",
	"{=t}":                                  "",
	"A{{=t}":                                "A{",
	"A{=t,tho=true}":                        "A{tho=true}",
	"A{=t,tho=false}":                       "A",
	"A{=1t,tho=false}":                      "A",
	"A{tho=false,=1t}":                      "A",
	"A{=1t,pl=NaN,a=m}":                     "A",
	"A{=1t,pl=NaN,a=s}":                     "A{a=s}",
	"A{pl=NaN,price=ASK}":                   "A{price=ask}",
	"A{pl=NaN,price=last,a=s}":              "A{a=s}",
	"A{pl=NaN,price=last,a=s,tho=TRUE}":     "A{a=s,tho=true}",
	"B{=Tick}":                              "B",
	"B{=Ticks}":                             "B",
	"A{tho=false,tho=false}":                "A",
	"A{tho=FALSE}":                          "A",
	"A{price=last,price=bid}":               "A",
	"A{a=m,a=s}":                            "A",
	"A{=t,=m}":                              "A",
	"A{a=m,=t,tho=false}":                   "A",
	"A{a=m,=1t}":                            "A",
	"A{=1t,a=s}":                            "A{a=s}",
	"A{=1.0t}":                              "A",
}

func TestParseCandleSymbolMatchesJavaAttributesAndNormalizesIdempotently(t *testing.T) {
	for _, c := range symbolValueOfCases {
		symbol, err := ParseCandleSymbol(c.input)
		if !c.ok {
			if err == nil {
				t.Errorf("ParseCandleSymbol(%q) did not fail", c.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseCandleSymbol(%q): %v", c.input, err)
			continue
		}

		expected := c.str
		if overridden, ok := cxxStrOverrides[c.input]; ok {
			expected = overridden
		}
		check := func(name, got, want string) {
			if got != want {
				t.Errorf("ParseCandleSymbol(%q): %s = %q, want %q", c.input, name, got, want)
			}
		}
		check("String", symbol.String(), expected)
		check("BaseSymbol", symbol.BaseSymbol(), c.base)
		check("Price", symbol.Price().String(), c.price)
		check("Session", symbol.Session().String(), c.session)
		check("Period", symbol.Period().String(), c.period)
		check("Alignment", symbol.Alignment().String(), c.alignment)
		check("PriceLevel", symbol.PriceLevel().String(), c.priceLevel)

		again, err := ParseCandleSymbol(symbol.String())
		if err != nil {
			t.Errorf("ParseCandleSymbol(%q): %v", symbol.String(), err)
			continue
		}
		check("the second normalization", again.String(), symbol.String())
	}
}

func TestCandleSymbolKeepsAllAttributes(t *testing.T) {
	level, _ := NewCandlePriceLevel(0.5)
	original, err := NewCandleSymbolWithAttributes("EUR/USD&E{source=bank}", CandlePriceBid,
		NewCandlePeriod(2, CandleTypeHour), CandleSessionRegular, CandleAlignmentSession, level)
	if err != nil {
		t.Fatal(err)
	}
	if original.String() != "EUR/USD&E{=2h,a=s,pl=0.5,price=bid,source=bank,tho=true}" {
		t.Errorf("String() = %q", original.String())
	}

	check := func(other *CandleSymbol) {
		t.Helper()
		if !other.Equal(original) || other.String() != original.String() || other.BaseSymbol() != original.BaseSymbol() ||
			other.Exchange() != original.Exchange() || other.Price() != original.Price() ||
			other.Session() != original.Session() || !other.Period().Equal(original.Period()) ||
			other.Alignment() != original.Alignment() || !other.PriceLevel().Equal(original.PriceLevel()) {
			t.Errorf("%+v differs from %+v", other, original)
		}
	}

	copied := *original
	check(&copied)

	// The copy holds the same fields as a symbol parsed from its string.
	reparsed, err := ParseCandleSymbol(original.String())
	if err != nil {
		t.Fatal(err)
	}
	check(reparsed)
	if original.Exchange().String() != "E" || original.BaseSymbol() != "EUR/USD" {
		t.Errorf("exchange %q, base %q", original.Exchange(), original.BaseSymbol())
	}

	assigned := NewCandleSymbol("IBM")
	*assigned = *original
	check(assigned)
}

func TestCandleSymbolNormalizationKeepsADefaultAttributeThatTheBaseSymbolWouldTakeForItsOwn(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"A{=t}{=t}", "A{=t}{=t}"},
		{"A{=t}{=1t}", "A{=t}{=t}"},
		{"X{=t,=bpr}{=t,=bpr}", "X{=t,=bpr}{=t}"},
		{"IB{=1M,=x}}{pl=NaN}", "IB{=1M,=x}}{pl=NaN}"},
		{"A{price=bid}{=t,price=bid}", "A{price=bid}{price=bid}"},
		{"A{b=1}{price=last}", "A{b=1}{price=last}"},
		{"A{=t,}", "A{=t,}"},
		{"A{&}{=t}", "A{&}{=t}"},
	} {
		symbol, err := ParseCandleSymbol(c.input)
		if err != nil {
			t.Errorf("ParseCandleSymbol(%q): %v", c.input, err)
			continue
		}
		if symbol.String() != c.want {
			t.Errorf("ParseCandleSymbol(%q) = %q, want %q", c.input, symbol, c.want)
		}
		again, err := ParseCandleSymbol(symbol.String())
		if err != nil || again.String() != c.want {
			t.Errorf("ParseCandleSymbol(%q) = %v, %v, want %q", symbol, again, err, c.want)
		}
		if again.BaseSymbol() != symbol.BaseSymbol() {
			t.Errorf("the base symbol of %q is %q, of its reparsed form %q", c.input, symbol.BaseSymbol(), again.BaseSymbol())
		}
	}
}

// The test of MarketEventSymbols in the C++ API that uses CandleAlignment; it is here since the eventsymbols package
// cannot import this one.
func TestNormalizeAttributeForSymbol(t *testing.T) {
	if got := NormalizeCandleAlignmentAttribute("AAPL&Q{a=m}"); got != "AAPL&Q" {
		t.Errorf("got %q", got)
	}
	if got := eventsymbols.GetBaseSymbol("AAPL&Q{a=m}"); got != "AAPL" {
		t.Errorf("got %q", got)
	}
}

func TestCandleSymbolDeprecatedAccessors(t *testing.T) {
	c := NewCandleSymbol("AAPL{=5Min}")
	if *c.Symbol() != "AAPL{=5m}" {
		t.Errorf("Symbol() = %q", *c.Symbol())
	}
	symbol := "IBM{price=BID}"
	c.SetSymbol(&symbol)
	if c.String() != "IBM{price=bid}" || c.Price() != CandlePriceBid {
		t.Errorf("after SetSymbol: %q, %v", c, c.Price())
	}
	c.SetSymbol(nil)
	if c.String() != "" {
		t.Errorf("after SetSymbol(nil): %q", c)
	}
}
