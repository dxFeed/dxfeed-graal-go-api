package glossary

import "github.com/dxfeed/dxfeed-graal-go-api/internal/native"

// AdditionalUnderlyings are the underlyings of an option besides the main one, each with its shares per contract
// (SPC), such as "FIS 53; US$ 45.46". The zero value has none.
type AdditionalUnderlyings struct {
	text string
}

// ParseAdditionalUnderlyings parses the text of additional underlyings, as in the ADDITIONAL_UNDERLYINGS field of a
// profile.
func ParseAdditionalUnderlyings(text string) (AdditionalUnderlyings, error) {
	normalized, err := native.NormalizeAdditionalUnderlyings(text)
	if err != nil {
		return AdditionalUnderlyings{}, err
	}
	return AdditionalUnderlyings{normalized}, nil
}

// AdditionalUnderlyingsOf returns the additional underlyings of a map from a symbol to its SPC.
func AdditionalUnderlyingsOf(spcs map[string]float64) (AdditionalUnderlyings, error) {
	text, err := native.AdditionalUnderlyingsFromMap(spcs)
	if err != nil {
		return AdditionalUnderlyings{}, err
	}
	return AdditionalUnderlyings{text}, nil
}

// Text returns the text of the additional underlyings, ordered by symbol.
func (a AdditionalUnderlyings) Text() string {
	return a.text
}

// String returns the text of the additional underlyings.
func (a AdditionalUnderlyings) String() string {
	return a.text
}

// Map returns the map from a symbol to its SPC.
func (a AdditionalUnderlyings) Map() (map[string]float64, error) {
	return native.AdditionalUnderlyingsMap(a.text)
}

// SPC returns the SPC of a symbol, or 0 if it is not an additional underlying.
func (a AdditionalUnderlyings) SPC(symbol string) (float64, error) {
	return native.AdditionalUnderlyingsSPC(a.text, symbol)
}

// SPCFromText returns the SPC of a symbol in the text of additional underlyings, or 0 if it is not there; it does
// not parse the rest of the text.
func SPCFromText(text, symbol string) (float64, error) {
	return native.AdditionalUnderlyingsSPCFromText(text, symbol)
}
