package ipf

import (
	"fmt"
	"math"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// InstrumentProfileField is a standard field of an instrument profile. The fields read and write profiles by field
// name, in the text form of the IPF format or in the numeric form; see also GetField and SetField, which accept the
// names of custom fields too.
type InstrumentProfileField int

// The standard fields, in the order of the IPF format.
const (
	FieldType InstrumentProfileField = iota
	FieldSymbol
	FieldDescription
	FieldLocalSymbol
	FieldLocalDescription
	FieldCountry
	FieldOPOL
	FieldExchangeData
	FieldExchanges
	FieldCurrency
	FieldBaseCurrency
	FieldCFI
	FieldISIN
	FieldSEDOL
	FieldCUSIP
	FieldICB
	FieldSIC
	FieldMultiplier
	FieldProduct
	FieldUnderlying
	FieldSPC
	FieldAdditionalUnderlyings
	FieldMMY
	FieldExpiration
	FieldLastTrade
	FieldStrike
	FieldOptionType
	FieldExpirationStyle
	FieldSettlementStyle
	FieldPriceIncrements
	FieldTradingHours
)

var fieldNames = [...]string{"TYPE", "SYMBOL", "DESCRIPTION", "LOCAL_SYMBOL", "LOCAL_DESCRIPTION", "COUNTRY", "OPOL",
	"EXCHANGE_DATA", "EXCHANGES", "CURRENCY", "BASE_CURRENCY", "CFI", "ISIN", "SEDOL", "CUSIP", "ICB", "SIC",
	"MULTIPLIER", "PRODUCT", "UNDERLYING", "SPC", "ADDITIONAL_UNDERLYINGS", "MMY", "EXPIRATION", "LAST_TRADE", "STRIKE",
	"OPTION_TYPE", "EXPIRATION_STYLE", "SETTLEMENT_STYLE", "PRICE_INCREMENTS", "TRADING_HOURS"}

// Fields returns all standard fields in the order of the IPF format.
func Fields() []InstrumentProfileField {
	fields := make([]InstrumentProfileField, len(fieldNames))
	for i := range fields {
		fields[i] = InstrumentProfileField(i)
	}
	return fields
}

// FindField returns the standard field with the name, such as "SYMBOL"; the names are case sensitive. It returns
// false for any other name, which is a custom field.
func FindField(name string) (InstrumentProfileField, bool) {
	for i, n := range fieldNames {
		if n == name {
			return InstrumentProfileField(i), true
		}
	}
	return 0, false
}

func (f InstrumentProfileField) valid() bool {
	return f >= 0 && int(f) < len(fieldNames)
}

// Name returns the name of the field in the IPF format, such as "SYMBOL".
func (f InstrumentProfileField) Name() string {
	if !f.valid() {
		return fmt.Sprintf("InstrumentProfileField(%d)", int(f))
	}
	return fieldNames[f]
}

// String returns the name of the field.
func (f InstrumentProfileField) String() string {
	return f.Name()
}

// IsNumeric reports whether the field has a numeric value: ICB, SIC, MULTIPLIER, SPC, EXPIRATION, LAST_TRADE and
// STRIKE. EXPIRATION and LAST_TRADE are dates, whose numeric value is a day id.
func (f InstrumentProfileField) IsNumeric() bool {
	switch f {
	case FieldICB, FieldSIC, FieldMultiplier, FieldSPC, FieldExpiration, FieldLastTrade, FieldStrike:
		return true
	}
	return false
}

// IsDate reports whether the field is a date: EXPIRATION and LAST_TRADE. The numeric value of a date is a day id, the
// number of days since January 1, 1970; the text is yyyy-MM-dd, or empty for 0.
func (f InstrumentProfileField) IsDate() bool {
	return f == FieldExpiration || f == FieldLastTrade
}

func (f InstrumentProfileField) stringField(p *events.InstrumentProfile) (get func() *string, set func(*string)) {
	switch f {
	case FieldType:
		return p.InstrumentType, p.SetInstrumentType
	case FieldSymbol:
		return p.Symbol, p.SetSymbol
	case FieldDescription:
		return p.Description, p.SetDescription
	case FieldLocalSymbol:
		return p.LocalSymbol, p.SetLocalSymbol
	case FieldLocalDescription:
		return p.LocalDescription, p.SetLocalDescription
	case FieldCountry:
		return p.Country, p.SetCountry
	case FieldOPOL:
		return p.Opol, p.SetOpol
	case FieldExchangeData:
		return p.ExchangeData, p.SetExchangeData
	case FieldExchanges:
		return p.Exchanges, p.SetExchanges
	case FieldCurrency:
		return p.Currency, p.SetCurrency
	case FieldBaseCurrency:
		return p.BaseCurrency, p.SetBaseCurrency
	case FieldCFI:
		return p.Cfi, p.SetCfi
	case FieldISIN:
		return p.Isin, p.SetIsin
	case FieldSEDOL:
		return p.Sedol, p.SetSedol
	case FieldCUSIP:
		return p.Cusip, p.SetCusip
	case FieldProduct:
		return p.Product, p.SetProduct
	case FieldUnderlying:
		return p.Underlying, p.SetUnderlying
	case FieldAdditionalUnderlyings:
		return p.AdditionalUnderlyings, p.SetAdditionalUnderlyings
	case FieldMMY:
		return p.Mmy, p.SetMmy
	case FieldOptionType:
		return p.OptionType, p.SetOptionType
	case FieldExpirationStyle:
		return p.ExpirationStyle, p.SetExpirationStyle
	case FieldSettlementStyle:
		return p.SettlementStyle, p.SetSettlementStyle
	case FieldPriceIncrements:
		return p.PriceIncrements, p.SetPriceIncrements
	case FieldTradingHours:
		return p.TradingHours, p.SetTradingHours
	}
	return nil, nil
}

func (f InstrumentProfileField) check(p *events.InstrumentProfile) error {
	if p == nil {
		return fmt.Errorf("ipf: nil instrument profile")
	}
	if !f.valid() {
		return fmt.Errorf("ipf: unknown field %v", f)
	}
	return nil
}

// Get returns the value of the field in the text form of the IPF format: a number like FormatNumber does, a date like
// FormatDate does.
func (f InstrumentProfileField) Get(p *events.InstrumentProfile) (string, error) {
	if err := f.check(p); err != nil {
		return "", err
	}
	if f.IsDate() {
		value, _ := f.GetNumeric(p)
		return FormatDate(int32(value))
	}
	if f.IsNumeric() {
		value, _ := f.GetNumeric(p)
		return FormatNumber(value)
	}
	get, _ := f.stringField(p)
	if s := get(); s != nil {
		return *s, nil
	}
	return "", nil
}

// Set sets the value of the field from the text form of the IPF format; the text of a number or a date is parsed like
// ParseNumber and ParseDate do, and a wrong format is an error.
func (f InstrumentProfileField) Set(p *events.InstrumentProfile, value string) error {
	if err := f.check(p); err != nil {
		return err
	}
	if f.IsDate() {
		day, err := ParseDate(value)
		if err != nil {
			return err
		}
		return f.SetNumeric(p, float64(day))
	}
	if f.IsNumeric() {
		number, err := ParseNumber(value)
		if err != nil {
			return err
		}
		return f.SetNumeric(p, number)
	}
	_, set := f.stringField(p)
	set(&value)
	return nil
}

// GetNumeric returns the value of a numeric field (see IsNumeric); a date gives its day id. It is an error for the
// other fields.
func (f InstrumentProfileField) GetNumeric(p *events.InstrumentProfile) (float64, error) {
	if err := f.check(p); err != nil {
		return 0, err
	}
	switch f {
	case FieldICB:
		return float64(p.Icb()), nil
	case FieldSIC:
		return float64(p.Sic()), nil
	case FieldMultiplier:
		return p.Multiplier(), nil
	case FieldSPC:
		return p.Spc(), nil
	case FieldExpiration:
		return float64(p.Expiration()), nil
	case FieldLastTrade:
		return float64(p.LastTrade()), nil
	case FieldStrike:
		return p.Strike(), nil
	}
	return 0, fmt.Errorf("ipf: field %v is not numeric", f)
}

// SetNumeric sets the value of a numeric field (see IsNumeric); ICB, SIC and the dates are 32-bit integers and keep the
// integer part, saturated to the int32 range (NaN gives 0), as a Java (int) cast does. It is an error for the other
// fields.
func (f InstrumentProfileField) SetNumeric(p *events.InstrumentProfile, value float64) error {
	if err := f.check(p); err != nil {
		return err
	}
	switch f {
	case FieldICB:
		p.SetIcb(int64(toInt32(value)))
	case FieldSIC:
		p.SetSic(int64(toInt32(value)))
	case FieldMultiplier:
		p.SetMultiplier(value)
	case FieldSPC:
		p.SetSpc(value)
	case FieldExpiration:
		p.SetExpiration(int64(toInt32(value)))
	case FieldLastTrade:
		p.SetLastTrade(int64(toInt32(value)))
	case FieldStrike:
		p.SetStrike(value)
	default:
		return fmt.Errorf("ipf: field %v is not numeric", f)
	}
	return nil
}

// GetField returns the value of a standard field (see InstrumentProfileField.Get) or, for any other name, of a custom
// field; a custom field that is not set gives an empty string.
func GetField(p *events.InstrumentProfile, name string) (string, error) {
	if f, ok := FindField(name); ok {
		return f.Get(p)
	}
	if p == nil {
		return "", fmt.Errorf("ipf: nil instrument profile")
	}
	return p.CustomFields()[name], nil
}

// SetField sets the value of a standard field (see InstrumentProfileField.Set) or, for any other name, of a custom
// field; an empty value removes a custom field.
func SetField(p *events.InstrumentProfile, name string, value string) error {
	if f, ok := FindField(name); ok {
		return f.Set(p, value)
	}
	if p == nil {
		return fmt.Errorf("ipf: nil instrument profile")
	}
	setCustomField(p, name, value)
	return nil
}

func setCustomField(p *events.InstrumentProfile, name string, value string) {
	fields := make(map[string]string, len(p.CustomFields())+1)
	for k, v := range p.CustomFields() {
		fields[k] = v
	}
	if value == "" {
		delete(fields, name)
	} else {
		fields[name] = value
	}
	p.SetCustomFields(fields)
}

// GetNumericField returns the value of a numeric standard field (see InstrumentProfileField.GetNumeric) or, for any
// other name, of a custom field: a custom field that is not set gives 0, a date (yyyy-MM-dd) gives its day id, any
// other text is parsed as a number. It is an error for the standard fields that are not numeric.
func GetNumericField(p *events.InstrumentProfile, name string) (float64, error) {
	if f, ok := FindField(name); ok {
		return f.GetNumeric(p)
	}
	if p == nil {
		return 0, fmt.Errorf("ipf: nil instrument profile")
	}
	value := p.CustomFields()[name]
	if value == "" {
		return 0, nil
	}
	if looksLikeDate(value) {
		day, err := ParseDate(value)
		return float64(day), err
	}
	return ParseNumber(value)
}

// SetNumericField sets the value of a numeric standard field (see InstrumentProfileField.SetNumeric) or, for any
// other name, of a custom field, which stores the number as FormatNumber formats it (0 removes the field). It is an
// error for the standard fields that are not numeric.
func SetNumericField(p *events.InstrumentProfile, name string, value float64) error {
	if f, ok := FindField(name); ok {
		return f.SetNumeric(p, value)
	}
	if p == nil {
		return fmt.Errorf("ipf: nil instrument profile")
	}
	text, err := FormatNumber(value)
	if err != nil {
		return err
	}
	setCustomField(p, name, text)
	return nil
}

// GetDateField returns the day id (days since January 1, 1970) of a numeric standard field (see
// InstrumentProfileField.GetNumeric; a numeric field that is not a date gives its integer part) or, for any other name,
// of a custom field, parsed like ParseDate does; a custom field that is not set gives 0. It is an error for the standard
// fields that are not numeric.
func GetDateField(p *events.InstrumentProfile, name string) (int32, error) {
	if f, ok := FindField(name); ok {
		value, err := f.GetNumeric(p)
		return int32(value), err
	}
	if p == nil {
		return 0, fmt.Errorf("ipf: nil instrument profile")
	}
	return ParseDate(p.CustomFields()[name])
}

// SetDateField sets the day id (days since January 1, 1970) of a numeric standard field (see
// InstrumentProfileField.SetNumeric) or, for any other name, of a custom field, which stores the date as FormatDate
// formats it (0 removes the field). It is an error for the standard fields that are not numeric.
func SetDateField(p *events.InstrumentProfile, name string, dayID int32) error {
	if f, ok := FindField(name); ok {
		return f.SetNumeric(p, float64(dayID))
	}
	if p == nil {
		return fmt.Errorf("ipf: nil instrument profile")
	}
	text, err := FormatDate(dayID)
	if err != nil {
		return err
	}
	setCustomField(p, name, text)
	return nil
}

func looksLikeDate(value string) bool {
	return len(value) == 10 && value[4] == '-' && value[7] == '-'
}

// FormatNumber formats a number like the numeric fields of the IPF format, without digit grouping; 0 gives an empty
// string, so a numeric field of 0 is empty in the text form.
func FormatNumber(value float64) (string, error) {
	return native.FormatProfileNumber(value)
}

// ParseNumber parses a number of the IPF format; an empty text gives 0.
func ParseNumber(text string) (float64, error) {
	return native.ParseProfileNumber(text)
}

// FormatDate formats a day id (days since January 1, 1970) as yyyy-MM-dd; 0 gives an empty string.
func FormatDate(dayID int32) (string, error) {
	return native.FormatProfileDate(dayID)
}

// ParseDate parses a yyyy-MM-dd date of the IPF format to a day id (days since January 1, 1970); an empty text gives
// 0.
func ParseDate(text string) (int32, error) {
	return native.ParseProfileDate(text)
}

// toInt32 converts like a Java (int) cast: the integer part, saturated to the int32 range, NaN gives 0.
func toInt32(value float64) int32 {
	switch {
	case value != value:
		return 0
	case value >= math.MaxInt32:
		return math.MaxInt32
	case value <= math.MinInt32:
		return math.MinInt32
	}
	return int32(value)
}
