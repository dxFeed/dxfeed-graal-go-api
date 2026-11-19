package ipf

import (
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/auth"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

const testProfiles = "testdata/profiles.ipf"

func readTestProfiles(t *testing.T) []*events.InstrumentProfile {
	t.Helper()
	reader, err := NewInstrumentProfileReader()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	profiles, err := reader.ReadFromFile(testProfiles)
	if err != nil {
		t.Fatalf("read %s: %v", testProfiles, err)
	}
	return profiles
}

func findProfile(t *testing.T, profiles []*events.InstrumentProfile, symbol string) *events.InstrumentProfile {
	t.Helper()
	for _, p := range profiles {
		if *p.Symbol() == symbol {
			return p
		}
	}
	t.Fatalf("no profile %s", symbol)
	return nil
}

func dayID(t *testing.T, date string) int32 {
	t.Helper()
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		t.Fatal(err)
	}
	return int32(d.Unix() / 86400)
}

func TestAllFieldsRead(t *testing.T) {
	profiles := readTestProfiles(t)
	if len(profiles) != 9 {
		t.Fatalf("read %d profiles, want 9", len(profiles))
	}
	p := findProfile(t, profiles, "TEST")
	want := map[InstrumentProfileField]string{
		FieldType:                  "STOCK",
		FieldSymbol:                "TEST",
		FieldDescription:           "Test Corp. common stock",
		FieldLocalSymbol:           "TST",
		FieldLocalDescription:      "Test local",
		FieldCountry:               "US",
		FieldOPOL:                  "XNAS",
		FieldExchangeData:          "data",
		FieldExchanges:             "ARCX;XNAS",
		FieldCurrency:              "USD",
		FieldBaseCurrency:          "EUR",
		FieldCFI:                   "ESNTPR",
		FieldISIN:                  "US0000000001",
		FieldSEDOL:                 "0000001",
		FieldCUSIP:                 "000000001",
		FieldICB:                   "9537",
		FieldSIC:                   "7372",
		FieldMultiplier:            "1",
		FieldProduct:               "PRD",
		FieldUnderlying:            "UND",
		FieldSPC:                   "100",
		FieldAdditionalUnderlyings: "UND2 50",
		FieldMMY:                   "202612",
		FieldExpiration:            "2026-12-18",
		FieldLastTrade:             "2026-12-17",
		FieldStrike:                "12.5",
		FieldOptionType:            "STD",
		FieldExpirationStyle:       "Monthly",
		FieldSettlementStyle:       "Close",
		FieldPriceIncrements:       "0.0001 1; 0.01",
		FieldTradingHours:          "NYSE()",
	}
	if len(want) != len(Fields()) {
		t.Fatalf("test covers %d fields, there are %d", len(want), len(Fields()))
	}
	for _, f := range Fields() {
		got, err := GetField(p, f.Name())
		if err != nil {
			t.Fatalf("%v: %v", f, err)
		}
		if got != want[f] {
			t.Errorf("%v = %q, want %q", f, got, want[f])
		}
	}
	if got, _ := GetField(p, "SECTOR"); got != "Technology" {
		t.Errorf("custom SECTOR = %q, want Technology", got)
	}
	if got, _ := GetField(p, "MISSING"); got != "" {
		t.Errorf("missing custom field = %q, want empty", got)
	}
	if p.Expiration() != int64(dayID(t, "2026-12-18")) || p.Strike() != 12.5 || p.Icb() != 9537 {
		t.Errorf("typed values: expiration %d, strike %v, icb %d", p.Expiration(), p.Strike(), p.Icb())
	}
	if got, err := GetDateField(p, "LAST_TRADE"); err != nil || got != dayID(t, "2026-12-17") {
		t.Errorf("GetDateField(LAST_TRADE) = %d, %v", got, err)
	}
}

func TestSetFields(t *testing.T) {
	p := events.NewInstrumentProfile()
	for name, value := range map[string]string{"SYMBOL": "ABC", "STRIKE": "1.25", "EXPIRATION": "2027-01-15",
		"ICB": "42", "SECTOR": "Energy", "LISTED": "2020-02-03", "LOT": "0.5"} {
		if err := SetField(p, name, value); err != nil {
			t.Fatalf("SetField(%s, %s): %v", name, value, err)
		}
	}
	if *p.Symbol() != "ABC" || p.Strike() != 1.25 || p.Expiration() != int64(dayID(t, "2027-01-15")) || p.Icb() != 42 {
		t.Fatalf("set values: %s %v %d %d", *p.Symbol(), p.Strike(), p.Expiration(), p.Icb())
	}
	if got, err := GetNumericField(p, "LOT"); err != nil || got != 0.5 {
		t.Errorf("GetNumericField(LOT) = %v, %v", got, err)
	}
	if got, err := GetNumericField(p, "LISTED"); err != nil || got != float64(dayID(t, "2020-02-03")) {
		t.Errorf("GetNumericField(LISTED) = %v, %v", got, err)
	}
	if got, err := GetDateField(p, "LISTED"); err != nil || got != dayID(t, "2020-02-03") {
		t.Errorf("GetDateField(LISTED) = %v, %v", got, err)
	}
	if err := SetNumericField(p, "LOT", 2); err != nil {
		t.Fatal(err)
	}
	if got, _ := GetField(p, "LOT"); got != "2" {
		t.Errorf("LOT after SetNumericField = %q, want 2", got)
	}
	if err := SetDateField(p, "LISTED", dayID(t, "2021-03-04")); err != nil {
		t.Fatal(err)
	}
	if got, _ := GetField(p, "LISTED"); got != "2021-03-04" {
		t.Errorf("LISTED after SetDateField = %q", got)
	}
	if err := SetField(p, "SECTOR", ""); err != nil {
		t.Fatal(err)
	}
	if _, ok := p.CustomFields()["SECTOR"]; ok {
		t.Error("an empty value did not remove the custom field")
	}
	if got, err := GetNumericField(p, "MISSING"); err != nil || got != 0 {
		t.Errorf("missing numeric custom field = %v, %v", got, err)
	}
}

func TestFieldErrors(t *testing.T) {
	p := events.NewInstrumentProfile()
	if err := SetField(p, "STRIKE", "abc"); err == nil {
		t.Error("SetField(STRIKE, abc) did not fail")
	}
	if err := SetField(p, "EXPIRATION", "18.12.2026"); err == nil {
		t.Error("SetField(EXPIRATION, 18.12.2026) did not fail")
	}
	if _, err := GetNumericField(p, "SYMBOL"); err == nil {
		t.Error("GetNumericField(SYMBOL) did not fail")
	}
	if err := SetNumericField(p, "SYMBOL", 1); err == nil {
		t.Error("SetNumericField(SYMBOL) did not fail")
	}
	if _, err := GetField(nil, "SECTOR"); err == nil {
		t.Error("GetField(nil) did not fail")
	}
	if _, err := InstrumentProfileField(99).Get(p); err == nil {
		t.Error("unknown field did not fail")
	}
	if _, ok := FindField("symbol"); ok {
		t.Error("field names must be case sensitive")
	}
	if f, ok := FindField("TRADING_HOURS"); !ok || f != FieldTradingHours || f.String() != "TRADING_HOURS" {
		t.Errorf("FindField(TRADING_HOURS) = %v, %v", f, ok)
	}
	if !FieldExpiration.IsDate() || !FieldExpiration.IsNumeric() || FieldStrike.IsDate() || FieldSymbol.IsNumeric() {
		t.Error("IsDate or IsNumeric is wrong")
	}
}

func TestWriteAndReadBack(t *testing.T) {
	profiles := readTestProfiles(t)
	for _, name := range []string{"profiles.ipf", "profiles.ipf.gz", "profiles.zip"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			if err := WriteToFile(path, profiles); err != nil {
				t.Fatalf("write: %v", err)
			}
			reader, err := NewInstrumentProfileReader()
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			back, err := reader.ReadFromFile(path)
			if err != nil {
				t.Fatalf("read back: %v", err)
			}
			if len(back) != len(profiles) {
				t.Fatalf("read back %d profiles, want %d", len(back), len(profiles))
			}
			for i := range profiles {
				for _, f := range Fields() {
					want, _ := f.Get(profiles[i])
					got, _ := f.Get(back[i])
					if got != want {
						t.Errorf("profile %d %v = %q, want %q", i, f, got, want)
					}
				}
				if got, want := back[i].CustomFields()["SECTOR"], profiles[i].CustomFields()["SECTOR"]; got != want {
					t.Errorf("profile %d SECTOR = %q, want %q", i, got, want)
				}
			}
		})
	}
	if err := WriteToFile(filepath.Join(t.TempDir(), "x.ipf"), []*events.InstrumentProfile{nil}); err == nil {
		t.Error("a nil profile did not fail")
	}
}

func TestReadFromFileWithToken(t *testing.T) {
	content, err := os.ReadFile(testProfiles)
	if err != nil {
		t.Fatal(err)
	}
	var mu sync.Mutex
	var authorization string
	lastAuthorization := func() string {
		mu.Lock()
		defer mu.Unlock()
		return authorization
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		authorization = r.Header.Get("Authorization")
		mu.Unlock()
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write(content)
	}))
	defer server.Close()

	bearer, err := auth.NewBearer("secret")
	if err != nil {
		t.Fatal(err)
	}
	basic, err := auth.NewBasicUser("user", "password")
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range []*auth.Token{bearer, basic} {
		reader, err := NewInstrumentProfileReader()
		if err != nil {
			t.Fatal(err)
		}
		profiles, err := reader.ReadFromFileWithToken(server.URL+"/profiles.ipf", token)
		_ = reader.Close()
		if err != nil {
			t.Fatalf("%s: %v", token.Scheme(), err)
		}
		if len(profiles) != 9 {
			t.Errorf("%s: read %d profiles, want 9", token.Scheme(), len(profiles))
		}
		if got := lastAuthorization(); got != token.HTTPAuthorization() {
			t.Errorf("Authorization = %q, want %q", got, token.HTTPAuthorization())
		}
	}

	reader, err := NewInstrumentProfileReader()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if _, err := reader.ReadFromFileWithToken(server.URL+"/profiles.ipf", nil); err == nil {
		t.Error("a read without a token did not get the 401")
	}
	if got := lastAuthorization(); got != "" {
		t.Errorf("a nil token sent Authorization %q", got)
	}
}

func TestWriteGoProfile(t *testing.T) {
	p := events.NewInstrumentProfile()
	for name, value := range map[string]string{"TYPE": TypeFuture, "SYMBOL": "/ABCH27", "DESCRIPTION": "Été, \"quoted\"",
		"MULTIPLIER": "1000", "EXPIRATION": "2027-03-19", "ZETA": "z", "ALPHA": "a", "MIDDLE": "Ünïcode"} {
		if err := SetField(p, name, value); err != nil {
			t.Fatal(err)
		}
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "one.ipf")
	if err := WriteToFile(path, []*events.InstrumentProfile{p}); err != nil {
		t.Fatalf("write: %v", err)
	}
	reader, err := NewInstrumentProfileReader()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	back, err := reader.ReadFromFile(path)
	if err != nil || len(back) != 1 {
		t.Fatalf("read back: %d profiles, %v", len(back), err)
	}
	for _, f := range Fields() {
		want, _ := f.Get(p)
		if got, _ := f.Get(back[0]); got != want {
			t.Errorf("%v = %q, want %q", f, got, want)
		}
	}
	if !reflect.DeepEqual(back[0].CustomFields(), p.CustomFields()) {
		t.Errorf("custom fields = %v, want %v", back[0].CustomFields(), p.CustomFields())
	}

	if err := WriteToFile(filepath.Join(dir, "empty.ipf"), nil); err != nil {
		t.Errorf("write of no profiles: %v", err)
	}
	if err := WriteToFile(filepath.Join(dir, "missing", "x.ipf"), []*events.InstrumentProfile{p}); err == nil {
		t.Error("a write to a missing directory did not fail")
	}
	p.SetIcb(1 << 32)
	if err := WriteToFile(filepath.Join(dir, "range.ipf"), []*events.InstrumentProfile{p}); err == nil {
		t.Error("an ICB out of the int32 range did not fail")
	}
}

func TestNumericSaturation(t *testing.T) {
	p := events.NewInstrumentProfile()
	for value, want := range map[float64]int64{1e10: math.MaxInt32, -1e10: math.MinInt32, math.NaN(): 0, 7.9: 7} {
		if err := FieldICB.SetNumeric(p, value); err != nil {
			t.Fatal(err)
		}
		if p.Icb() != want {
			t.Errorf("SetNumeric(ICB, %v) = %d, want %d", value, p.Icb(), want)
		}
	}
	if err := FieldStrike.SetNumeric(p, 0); err != nil {
		t.Fatal(err)
	}
	if got, _ := FieldStrike.Get(p); got != "" {
		t.Errorf("STRIKE 0 = %q, want empty", got)
	}
}
