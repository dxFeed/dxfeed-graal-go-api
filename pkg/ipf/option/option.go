// Package option groups option instrument profiles into chains: an [OptionChain] of a product or an underlying holds
// its [OptionSeries], and a series holds the calls and puts of one expiration and contract terms by strike.
//
//	reader, _ := ipf.NewInstrumentProfileReader()
//	profiles, _ := reader.ReadFromFile(address)
//	chains := option.BuildChains(profiles)
//	for _, series := range chains["AAPL"].Series() {
//		fmt.Println(series, series.Strikes())
//	}
package option

import (
	"cmp"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
)

// seriesKey are the attributes that tell the series apart; it is comparable, so it keys the series of a chain.
type seriesKey struct {
	expiration, lastTrade int64
	multiplier, spc       uint64 // math.Float64bits of normalizeFloat, so that NaN equals NaN
	additionalUnderlyings string
	mmy, optionType       string
	expirationStyle       string
	settlementStyle       string
	cfi                   string
}

// normalizeFloat makes the numbers that the series treat as equal equal for ==: -0 and 0, and all NaNs.
func normalizeFloat(v float64) float64 {
	switch {
	case v == 0:
		return 0
	case math.IsNaN(v):
		return math.NaN()
	}
	return v
}

func compareKeys(a, b seriesKey) int {
	return firstNonZero(
		cmp.Compare(a.expiration, b.expiration),
		cmp.Compare(a.lastTrade, b.lastTrade),
		cmp.Compare(math.Float64frombits(a.multiplier), math.Float64frombits(b.multiplier)),
		cmp.Compare(math.Float64frombits(a.spc), math.Float64frombits(b.spc)),
		strings.Compare(a.additionalUnderlyings, b.additionalUnderlyings),
		strings.Compare(a.mmy, b.mmy),
		strings.Compare(a.optionType, b.optionType),
		strings.Compare(a.expirationStyle, b.expirationStyle),
		strings.Compare(a.settlementStyle, b.settlementStyle),
		strings.Compare(a.cfi, b.cfi),
	)
}

func firstNonZero(results ...int) int {
	for _, r := range results {
		if r != 0 {
			return r
		}
	}
	return 0
}

// OptionSeries are the options of a chain with the same expiration and contract terms: the calls and the puts by
// strike. The series of a chain are ordered by expiration, last trading day, multiplier, SPC, additional underlyings,
// MMY, option type, expiration style, settlement style and CFI.
type OptionSeries struct {
	key     seriesKey
	calls   map[float64]*events.InstrumentProfile
	puts    map[float64]*events.InstrumentProfile
	strikes []float64 // sorted, nil until computed
}

// Expiration returns the day id of the expiration (days since January 1, 1970).
func (s *OptionSeries) Expiration() int64 { return s.key.expiration }

// LastTrade returns the day id of the last trading day, or 0.
func (s *OptionSeries) LastTrade() int64 { return s.key.lastTrade }

// Multiplier returns the market value multiplier.
func (s *OptionSeries) Multiplier() float64 { return math.Float64frombits(s.key.multiplier) }

// SPC returns the shares per contract of the underlying.
func (s *OptionSeries) SPC() float64 { return math.Float64frombits(s.key.spc) }

// AdditionalUnderlyings returns the text of the additional underlyings, or "".
func (s *OptionSeries) AdditionalUnderlyings() string { return s.key.additionalUnderlyings }

// MMY returns the maturity month-year as in the FIX protocol, or "".
func (s *OptionSeries) MMY() string { return s.key.mmy }

// OptionType returns the type of the options, such as STD or LEAP, or "".
func (s *OptionSeries) OptionType() string { return s.key.optionType }

// ExpirationStyle returns the expiration cycle style, such as Weeklys or Quarterlys, or "".
func (s *OptionSeries) ExpirationStyle() string { return s.key.expirationStyle }

// SettlementStyle returns the settlement price determination style, such as Open or Close, or "".
func (s *OptionSeries) SettlementStyle() string { return s.key.settlementStyle }

// CFI returns the CFI code of the series: the CFI of its options with X in place of the call or put letter, such as
// OXASPS.
func (s *OptionSeries) CFI() string { return s.key.cfi }

// Calls returns the calls by strike; the map is a copy.
func (s *OptionSeries) Calls() map[float64]*events.InstrumentProfile {
	return maps.Clone(s.calls)
}

// Puts returns the puts by strike; the map is a copy.
func (s *OptionSeries) Puts() map[float64]*events.InstrumentProfile {
	return maps.Clone(s.puts)
}

// Strikes returns the strikes of the calls and puts in ascending order; the slice is a copy.
func (s *OptionSeries) Strikes() []float64 {
	return slices.Clone(s.sortedStrikes())
}

func (s *OptionSeries) sortedStrikes() []float64 {
	if s.strikes == nil {
		strikes := make([]float64, 0, len(s.calls)+len(s.puts))
		for strike := range s.calls {
			strikes = append(strikes, strike)
		}
		for strike := range s.puts {
			if _, ok := s.calls[strike]; !ok {
				strikes = append(strikes, strike)
			}
		}
		slices.Sort(strikes)
		s.strikes = strikes
	}
	return s.strikes
}

// NStrikesAround returns up to n strikes in ascending order around a strike: about n/2 below it and the rest from it
// up, fewer at the ends of the strikes. It returns nil for n <= 0.
func (s *OptionSeries) NStrikesAround(n int, strike float64) []float64 {
	if n <= 0 {
		return nil
	}
	strikes := s.sortedStrikes()
	index, _ := slices.BinarySearch(strikes, strike)
	from := max(0, index-n/2)
	to := min(len(strikes), from+n)
	return slices.Clone(strikes[from:to])
}

// String returns the attributes of the series, such as "expiration=20261218, multiplier=100, cfi=OXASPS"; the
// attributes that are not set are left out.
func (s *OptionSeries) String() string {
	var b strings.Builder
	b.WriteString("expiration=" + yearMonthDay(s.key.expiration))
	if s.key.lastTrade != 0 {
		b.WriteString(", lastTrade=" + yearMonthDay(s.key.lastTrade))
	}
	if m := s.Multiplier(); m != 0 {
		b.WriteString(", multiplier=" + strconv.FormatFloat(m, 'g', -1, 64))
	}
	if spc := s.SPC(); spc != 0 {
		b.WriteString(", spc=" + strconv.FormatFloat(spc, 'g', -1, 64))
	}
	for _, a := range []struct{ name, value string }{
		{"additionalUnderlyings", s.key.additionalUnderlyings}, {"mmy", s.key.mmy}, {"optionType", s.key.optionType},
		{"expirationStyle", s.key.expirationStyle}, {"settlementStyle", s.key.settlementStyle},
	} {
		if a.value != "" {
			b.WriteString(", " + a.name + "=" + a.value)
		}
	}
	b.WriteString(", cfi=" + s.key.cfi)
	return b.String()
}

// yearMonthDay formats a day id as yyyymmdd.
func yearMonthDay(dayID int64) string {
	return time.Unix(dayID*86400, 0).UTC().Format("20060102")
}

func (s *OptionSeries) add(isCall bool, strike float64, profile *events.InstrumentProfile) {
	m := s.puts
	if isCall {
		m = s.calls
	}
	if _, ok := m[strike]; !ok {
		s.strikes = nil
	}
	m[strike] = profile
}

// OptionChain holds the option series of a product or an underlying.
type OptionChain struct {
	symbol string
	series map[seriesKey]*OptionSeries
}

// Symbol returns the product or underlying symbol of the chain.
func (c *OptionChain) Symbol() string {
	return c.symbol
}

// Series returns the series of the chain in their order.
func (c *OptionChain) Series() []*OptionSeries {
	result := make([]*OptionSeries, 0, len(c.series))
	for _, s := range c.series {
		result = append(result, s)
	}
	slices.SortFunc(result, func(a, b *OptionSeries) int { return compareKeys(a.key, b.key) })
	return result
}

func (c *OptionChain) add(key seriesKey, isCall bool, strike float64, profile *events.InstrumentProfile) {
	s := c.series[key]
	if s == nil {
		s = &OptionSeries{key: key, calls: map[float64]*events.InstrumentProfile{},
			puts: map[float64]*events.InstrumentProfile{}}
		c.series[key] = s
	}
	s.add(isCall, strike, profile)
}

// OptionChainsBuilder groups option profiles into chains, one per product and one per underlying.
type OptionChainsBuilder struct {
	chains map[string]*OptionChain
}

// NewOptionChainsBuilder creates a builder without chains.
func NewOptionChainsBuilder() *OptionChainsBuilder {
	return &OptionChainsBuilder{chains: map[string]*OptionChain{}}
}

// BuildChains returns the chains of the option profiles by product and underlying symbol; the other profiles are
// skipped.
func BuildChains(profiles []*events.InstrumentProfile) map[string]*OptionChain {
	b := NewOptionChainsBuilder()
	for _, p := range profiles {
		b.Add(p)
	}
	return b.Chains()
}

func value(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Add adds an option profile to the chain of its product and to the chain of its underlying, if they are set. It
// skips a profile that is not of the type OPTION, whose CFI does not start with OC (a call) or OP (a put), that has
// no expiration or whose strike is not finite. A later option with the same series and strike replaces the earlier
// one.
func (b *OptionChainsBuilder) Add(profile *events.InstrumentProfile) {
	if profile == nil || value(profile.InstrumentType()) != ipf.TypeOption {
		return
	}
	cfi := value(profile.Cfi())
	isCall := strings.HasPrefix(cfi, "OC")
	if !isCall && !strings.HasPrefix(cfi, "OP") {
		return
	}
	strike := profile.Strike()
	if profile.Expiration() == 0 || math.IsNaN(strike) || math.IsInf(strike, 0) {
		return
	}
	key := seriesKey{
		expiration:            profile.Expiration(),
		lastTrade:             profile.LastTrade(),
		multiplier:            math.Float64bits(normalizeFloat(profile.Multiplier())),
		spc:                   math.Float64bits(normalizeFloat(profile.Spc())),
		additionalUnderlyings: value(profile.AdditionalUnderlyings()),
		mmy:                   value(profile.Mmy()),
		optionType:            value(profile.OptionType()),
		expirationStyle:       value(profile.ExpirationStyle()),
		settlementStyle:       value(profile.SettlementStyle()),
		cfi:                   cfi[:1] + "X" + cfi[2:],
	}
	strike = normalizeFloat(strike)
	for _, symbol := range []string{value(profile.Product()), value(profile.Underlying())} {
		if symbol != "" {
			b.chain(symbol).add(key, isCall, strike, profile)
		}
	}
}

func (b *OptionChainsBuilder) chain(symbol string) *OptionChain {
	c := b.chains[symbol]
	if c == nil {
		c = &OptionChain{symbol: symbol, series: map[seriesKey]*OptionSeries{}}
		b.chains[symbol] = c
	}
	return c
}

// Chains returns the chains by product and underlying symbol. The map is a copy, but the chains are those of the
// builder: adding more profiles changes them.
func (b *OptionChainsBuilder) Chains() map[string]*OptionChain {
	return maps.Clone(b.chains)
}
