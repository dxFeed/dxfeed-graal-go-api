package option

import (
	"math"
	"slices"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
)

func readProfiles(t *testing.T) []*events.InstrumentProfile {
	t.Helper()
	reader, err := ipf.NewInstrumentProfileReader()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	profiles, err := reader.ReadFromFile("../testdata/profiles.ipf")
	if err != nil {
		t.Fatal(err)
	}
	return profiles
}

func symbols(m map[float64]*events.InstrumentProfile) map[float64]string {
	result := make(map[float64]string, len(m))
	for strike, p := range m {
		result[strike] = *p.Symbol()
	}
	return result
}

func TestBuildChainsFromFile(t *testing.T) {
	chains := BuildChains(readProfiles(t))
	// The options have the underlying TEST and no product; the stock and the future are skipped.
	if len(chains) != 1 || chains["TEST"] == nil {
		t.Fatalf("chains = %v", chains)
	}
	chain := chains["TEST"]
	if chain.Symbol() != "TEST" {
		t.Errorf("symbol = %q", chain.Symbol())
	}
	series := chain.Series()
	if len(series) != 2 {
		t.Fatalf("%d series", len(series))
	}
	dec, jan := series[0], series[1]
	if dec.String() != "expiration=20261218, lastTrade=20261218, multiplier=100, spc=100, mmy=202612, "+
		"optionType=STD, expirationStyle=Monthly, settlementStyle=Close, cfi=OXASPS" {
		t.Errorf("series = %s", dec)
	}
	if dec.Expiration() != 20805 || dec.LastTrade() != 20805 || dec.Multiplier() != 100 || dec.SPC() != 100 ||
		dec.MMY() != "202612" || dec.OptionType() != "STD" || dec.ExpirationStyle() != "Monthly" ||
		dec.SettlementStyle() != "Close" || dec.CFI() != "OXASPS" || dec.AdditionalUnderlyings() != "" {
		t.Errorf("attributes of the series = %s, expiration %d", dec, dec.Expiration())
	}
	if !slices.Equal(dec.Strikes(), []float64{10, 15, 20}) {
		t.Errorf("strikes = %v", dec.Strikes())
	}
	calls, puts := symbols(dec.Calls()), symbols(dec.Puts())
	if len(calls) != 3 || calls[20] != ".TEST261218C20" || len(puts) != 2 || puts[15] != ".TEST261218P15" {
		t.Errorf("calls = %v, puts = %v", calls, puts)
	}
	if jan.Expiration() <= dec.Expiration() || jan.MMY() != "202701" || !slices.Equal(jan.Strikes(), []float64{10}) {
		t.Errorf("second series = %s, strikes %v", jan, jan.Strikes())
	}

	for _, c := range []struct {
		n      int
		strike float64
		want   []float64
	}{
		{1, 15, []float64{15}}, {2, 15, []float64{10, 15}}, {3, 15, []float64{10, 15, 20}},
		{2, 12, []float64{10, 15}}, {2, 0, []float64{10, 15}}, {2, 100, []float64{20}},
		{10, 15, []float64{10, 15, 20}}, {0, 15, nil}, {-1, 15, nil},
	} {
		if got := dec.NStrikesAround(c.n, c.strike); !slices.Equal(got, c.want) {
			t.Errorf("NStrikesAround(%d, %v) = %v, want %v", c.n, c.strike, got, c.want)
		}
	}

	// The copies do not change the series.
	dec.Strikes()[0] = 99
	delete(dec.Calls(), 10)
	if dec.Strikes()[0] != 10 || len(dec.Calls()) != 3 {
		t.Error("a copy changed the series")
	}
}

func option(t *testing.T, fields map[string]string) *events.InstrumentProfile {
	t.Helper()
	p := events.NewInstrumentProfile()
	base := map[string]string{"TYPE": ipf.TypeOption, "CFI": "OCASPS", "EXPIRATION": "2026-12-18", "STRIKE": "10"}
	for name, value := range base {
		if _, ok := fields[name]; !ok {
			fields[name] = value
		}
	}
	for name, value := range fields {
		if err := ipf.SetField(p, name, value); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

func TestBuilderRules(t *testing.T) {
	b := NewOptionChainsBuilder()
	b.Add(nil)
	b.Add(option(t, map[string]string{"SYMBOL": "stock", "TYPE": ipf.TypeStock, "UNDERLYING": "U"}))
	b.Add(option(t, map[string]string{"SYMBOL": "future CFI", "CFI": "FFICSX", "UNDERLYING": "U"}))
	b.Add(option(t, map[string]string{"SYMBOL": "no expiration", "EXPIRATION": "", "UNDERLYING": "U"}))
	b.Add(option(t, map[string]string{"SYMBOL": "no chain"}))
	nan := option(t, map[string]string{"SYMBOL": "NaN strike", "UNDERLYING": "U"})
	nan.SetStrike(math.NaN())
	b.Add(nan)
	inf := option(t, map[string]string{"SYMBOL": "Inf strike", "UNDERLYING": "U"})
	inf.SetStrike(math.Inf(1))
	b.Add(inf)
	if chains := b.Chains(); len(chains) != 0 {
		t.Fatalf("skipped profiles made chains %v", chains)
	}

	// An option goes to the chain of its product and to the chain of its underlying.
	b.Add(option(t, map[string]string{"SYMBOL": "call", "PRODUCT": "P", "UNDERLYING": "U"}))
	b.Add(option(t, map[string]string{"SYMBOL": "put", "CFI": "OPASPS", "UNDERLYING": "U"}))
	// A later option with the same strike replaces the earlier one.
	b.Add(option(t, map[string]string{"SYMBOL": "call 2", "UNDERLYING": "U"}))
	// Another multiplier is another series; a NaN multiplier is one series.
	for _, symbol := range []string{"nan 1", "nan 2"} {
		p := option(t, map[string]string{"SYMBOL": symbol, "UNDERLYING": "U"})
		p.SetMultiplier(math.NaN())
		b.Add(p)
	}
	chains := b.Chains()
	if len(chains) != 2 || len(chains["P"].Series()) != 1 {
		t.Fatalf("chains = %v", chains)
	}
	if calls := symbols(chains["P"].Series()[0].Calls()); calls[10] != "call" {
		t.Errorf("calls of the product = %v", calls)
	}
	series := chains["U"].Series()
	if len(series) != 2 {
		t.Fatalf("%d series of the underlying: %v", len(series), series)
	}
	var plain, withNaN *OptionSeries
	for _, s := range series {
		if math.IsNaN(s.Multiplier()) {
			withNaN = s
		} else {
			plain = s
		}
	}
	if plain == nil || withNaN == nil {
		t.Fatalf("series = %v", series)
	}
	if calls, puts := symbols(plain.Calls()), symbols(plain.Puts()); calls[10] != "call 2" || puts[10] != "put" {
		t.Errorf("calls = %v, puts = %v", calls, puts)
	}
	if calls := symbols(withNaN.Calls()); calls[10] != "nan 2" {
		t.Errorf("calls of the NaN series = %v", calls)
	}
	if plain.String() != "expiration=20261218, cfi=OXASPS" {
		t.Errorf("series without attributes = %s", plain)
	}
}

func TestSeriesOrder(t *testing.T) {
	b := NewOptionChainsBuilder()
	for _, f := range []map[string]string{
		{"SYMBOL": "c", "EXPIRATION": "2027-01-15"},
		{"SYMBOL": "b", "MMY": "B"},
		{"SYMBOL": "a", "MMY": "A"},
		{"SYMBOL": "d", "LAST_TRADE": "2026-12-17"},
	} {
		f["UNDERLYING"] = "U"
		b.Add(option(t, f))
	}
	var got []string
	for _, s := range b.Chains()["U"].Series() {
		for _, p := range s.Calls() {
			got = append(got, *p.Symbol())
		}
	}
	// By expiration, then by last trading day (none first), then by MMY.
	if want := []string{"a", "b", "d", "c"}; !slices.Equal(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}
